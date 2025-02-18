package client

import (
	log "github.com/sirupsen/logrus"
	"github.com/zourva/lwm2m/coap"
	. "github.com/zourva/lwm2m/core"
	"github.com/zourva/pareto/box/meta"
	"strings"
	"sync/atomic"
	"time"
)

// New returns a LwM2M client with the mandatory name
// and other options, or nil when any failure occurred.
func New(name string, store ObjectInstanceStore, opts ...Option) *LwM2MClient {
	c := &LwM2MClient{
		name:    name,
		store:   store,
		options: &Options{},
		machine: meta.NewStateMachine[state](name, time.Second),
	}

	for _, f := range opts {
		f(c.options)
	}

	if err := c.initialize(); err != nil {
		log.Errorln("initialize client failed:", err)
		return nil
	}

	return c
}

// LwM2MClient implements client side
// functionalities and exposes API to
// applications using callbacks.
type LwM2MClient struct {
	//endpoint assigned when provision
	name    string
	options *Options

	machine *meta.StateMachine[state]

	// store to save object instances
	// loaded from local persistent storage.
	store ObjectInstanceStore

	//coapConn coap.Client
	// messagerc - messager client
	messagerc *MessagerClient // messager to communicate with server

	// lifecycle event manager
	evtMgr *EventManager

	// delegators
	bootstrapper *Bootstrapper // instance created for the latest try
	registrar    *Registrar    // instance created for the latest try
	reporter     *Reporter
	controller   DeviceControlClient

	bootstrapPending atomic.Bool
	registerPending  atomic.Bool
	updatePending    atomic.Bool
}

func (c *LwM2MClient) initialize() error {
	c.makeDefaults()
	//c.messager = NewMessager(c)
	//c.bootstrapper = NewBootstrapper(c)
	//c.registrar = NewRegistrar(c)
	c.controller = NewDeviceController(c)
	c.reporter = NewReporter(c)
	c.machine.RegisterStates([]*meta.State[state]{
		{Name: initiating, Action: c.onInitiating},
		{Name: bootstrapping, Action: c.onBootstrapping},
		{Name: registering, Action: c.onRegistering},
		{Name: servicing, Action: c.onServicing},
		{Name: exiting, Action: c.onExiting},
	})

	c.evtMgr = NewEventManager()
	c.evtMgr.RegisterCreator(EventClientBeforeBootstrap, NewBootstrappedEvent)
	c.evtMgr.RegisterCreator(EventClientBootstrapped, NewBootstrappedEvent)
	c.evtMgr.RegisterCreator(EventClientBeforeRegister, NewBootstrappedEvent)
	c.evtMgr.RegisterCreator(EventClientRegistered, NewRegisteredEvent)
	c.evtMgr.RegisterCreator(EventClientBeforeUpdate, NewRegUpdatedEvent)
	c.evtMgr.RegisterCreator(EventClientRegUpdated, NewRegUpdatedEvent)
	c.evtMgr.RegisterCreator(EventClientBeforeUnregister, NewUnregisteredEvent)
	c.evtMgr.RegisterCreator(EventClientUnregistered, NewUnregisteredEvent)
	c.evtMgr.RegisterCreator(EventClientBeforeDevInfoChange, NewDeviceChangedEvent)
	c.evtMgr.RegisterCreator(EventClientDevInfoChanged, NewDeviceChangedEvent)
	c.evtMgr.RegisterCreator(EventClientBeforeObserve, NewInfoObservedEvent)
	c.evtMgr.RegisterCreator(EventClientObserved, NewInfoObservedEvent)
	c.evtMgr.RegisterCreator(EventClientBeforeObserveCancel, NewInfoObservedEvent)
	c.evtMgr.RegisterCreator(EventClientObserveCancelled, NewInfoObservedEvent)
	c.evtMgr.RegisterCreator(EventClientBeforeReport, NewInfoReportedEvent)
	c.evtMgr.RegisterCreator(EventClientReported, NewInfoReportedEvent)
	c.evtMgr.RegisterCreator(EventClientAbnormal, NewAbnormalEvent)

	if err := c.store.Load(); err != nil {
		log.Errorln("load object instances failed:", err)
		return err
	}

	if c.hasRegistrationServer() {
		// has been bootstrap
		c.initiateRegister()
	} else {
		c.initiateBootstrap(bootstrapReasonStartup)
	}

	return nil
}

func (c *LwM2MClient) messager() *MessagerClient {
	return c.messagerc
}

func (c *LwM2MClient) hasRegistrationServer() bool {
	instances := c.store.GetInstances(OmaObjectSecurity)
	for _, instance := range instances {
		//op := instance.Class().Operator()
		//fields, err := op.GetAll(instance, LwM2MSecurityBootstrapServer)
		//if err != nil {
		//	continue
		//}
		//field := fields.SingleField()
		//if !field.Get().(bool) { // BootstrapServer != true

		isBootstrapServer := FieldValue[bool](instance, LwM2MSecurityBootstrapServer)
		if !isBootstrapServer { // BootstrapServer == true
			return true
		}
	}
	return false
}

func (c *LwM2MClient) getBearerFromURISchema(uri string) (bearer, address string, secured bool) {
	if strings.HasPrefix(uri, coap.UdpCoapSchema) {
		return coap.UDPBearer, strings.TrimPrefix(uri, coap.UdpCoapSchema), false
	} else if strings.HasPrefix(uri, coap.DtlsCoapSchema) {
		return coap.UDPBearer, strings.TrimPrefix(uri, coap.DtlsCoapSchema), true
	} else if strings.HasPrefix(uri, coap.TcpCoapSchema) {
		return coap.TCPBearer, strings.TrimPrefix(uri, coap.TcpCoapSchema), false
	} else if strings.HasPrefix(uri, coap.TlsCoapSchema) {
		return coap.TCPBearer, strings.TrimPrefix(uri, coap.TlsCoapSchema), true
	}

	// udp without security by default
	return coap.UDPBearer, strings.TrimPrefix(uri, coap.UdpCoapSchema), false
}

// get bootstrap server account from store, if any
func (c *LwM2MClient) getBootstrapInfos() (*BootstrapServerBootstrapInfo, *ServerInfo) {
	instances := c.store.GetInstances(OmaObjectSecurity)
	for _, instance := range instances {
		isBootstrapServer := FieldValue[bool](instance, LwM2MSecurityBootstrapServer)
		if isBootstrapServer { // BootstrapServer == true
			uri := FieldValue[string](instance, LwM2MSecurityLwM2MServerURI)

			securityMode := FieldValue[int](instance, LwM2MSecuritySecurityMode)
			publicKeyOrIdentity := FieldValue[[]byte](instance, LwM2MSecurityPublicKeyOrIdentity)
			secretKey := FieldValue[[]byte](instance, LwM2MSecuritySecretKey)
			serverPublicKey := FieldValue[[]byte](instance, LwM2MSecurityServerPublicKeyOrIdentity)

			bootstrapInfo := &BootstrapServerBootstrapInfo{
				BootstrapServerAccount: &BootstrapServerAccount{
					SecurityObjectInstance: instance,
				}}

			network, address, secured := c.getBearerFromURISchema(uri)
			if secured && securityMode == SecurityModeNoSec {
				log.Errorln("security mode conflicts with bootstap server uri schema")
				return nil, nil
			}

			serverInfo := &ServerInfo{
				network:             network,
				address:             address,
				securityMode:        securityMode,
				publicKeyOrIdentity: publicKeyOrIdentity,
				serverPublicKey:     serverPublicKey,
				secretKey:           secretKey,
			}

			return bootstrapInfo, serverInfo
		}
	}

	return nil, nil
}

func (c *LwM2MClient) getRegistrationServers() []*regServerInfo {
	var list []*regServerInfo
	var ms = make(map[int]*regServerInfo)

	// extract server shortId<->info index from storage
	instances := c.store.GetInstances(OmaObjectSecurity)
	for _, instance := range instances {
		// bootstrapServer is true, otherwise false
		//instance.Class().Operator().Get(instance, LwM2MSecurityBootstrapServer, LwM2MSecurityLwM2MServerURI)
		isBootstrapServer := FieldValue[bool](instance, LwM2MSecurityBootstrapServer)
		if !isBootstrapServer {
			uri := FieldValue[string](instance, LwM2MSecurityLwM2MServerURI)
			shortId := FieldValue[int](instance, LwM2MSecurityShortServerID)

			securityMode := FieldValue[int](instance, LwM2MSecuritySecurityMode)
			publicKeyOrIdentity := FieldValue[[]byte](instance, LwM2MSecurityPublicKeyOrIdentity)
			secretKey := FieldValue[[]byte](instance, LwM2MSecuritySecretKey)
			serverPublicKey := FieldValue[[]byte](instance, LwM2MSecurityServerPublicKeyOrIdentity)

			network, address, secured := c.getBearerFromURISchema(uri)
			if secured && securityMode == SecurityModeNoSec {
				log.Errorln("security mode conflicts with server uri schema")
				// TODO expose error
			}

			ms[shortId] = &regServerInfo{
				ServerInfo: ServerInfo{
					network:             network,
					address:             address,
					securityMode:        securityMode,
					publicKeyOrIdentity: publicKeyOrIdentity,
					serverPublicKey:     serverPublicKey,
					secretKey:           secretKey,
				},
				lifetime:          defaultLifetime,
				blocking:          true,
				bootstrap:         true,
				priorityOrder:     1,
				initRegDelay:      defInitRegistrationDelay,
				commRetryLimit:    defCommRetryCount,
				commRetryDelay:    defCommRetryTimer,
				commSeqRetryDelay: defCommSeqDelayTimer,
				commSeqRetryLimit: defCommSeqRetryCount,
			}
		}
	}

	// get a target reg server from reg server list
	servers := c.store.GetInstances(OmaObjectServer)
	for _, server := range servers {
		// TODO: retrieve from server
		shortId := FieldValue[int](server, LwM2MServerShortServerID)

		if s, ok := ms[shortId]; ok {
			lifetime := FieldValue[int](server, LwM2MServerLifetime)
			s.lifetime = uint64(lifetime)
			list = append(list, s)
		}
	}

	return list
}

func (c *LwM2MClient) doBootstrap() {
	c.clearBootstrapPending()

	log.Infoln("client is ready to bootstrap")

	//c.messager().PauseUserPlane()

	bootstrapInfo, serverInfo := c.getBootstrapInfos()

	// always create a new bootstrapper
	opts := []BootstrapOption{
		WithBootstrapInfo(bootstrapInfo),
		WithServerInfo(serverInfo),
	}
	c.bootstrapper = NewBootstrapper(c, opts...)
	c.bootstrapper.Start()

	c.machine.MoveToState(bootstrapping)
}

func (c *LwM2MClient) doRegister() {
	c.clearRegisterPending()

	log.Infoln("client is ready to register")

	//c.resetReportFailCounter()
	//c.messager().PauseUserPlane()

	// always create a new registrar
	if c.registrar != nil {
		c.registrar.Stop()
	}
	c.registrar = NewRegistrar(c)
	c.registrar.Start()

	c.machine.MoveToState(registering)
}

func (c *LwM2MClient) enableService() {
	c.messagerc = c.registrar.messager

	//c.messager().ResumeUserPlane()
	c.machine.MoveToState(servicing)
}

func (c *LwM2MClient) onInitiating(_ any) {
	// determine bootstrap or registration procedure
	if c.bootstrapRequired() {
		c.doBootstrap()
	} else {
		c.doRegister()
	}
}

func (c *LwM2MClient) onBootstrapping(_ any) {
	if c.bootstrapper.Bootstrapped() {
		log.Infoln("client bootstrapped")
		c.evtMgr.EmitEvent(EventClientBootstrapped)
		c.bootstrapper.Stop()
		c.doRegister()
	} else {
		//restart bootstrap if timeout
		if c.bootstrapper.Timeout() {
			c.bootstrapper.Stop()
			c.initiateBootstrap(bootstrapReasonBootFail)
			log.Infof("client bootstrap timeout, retry")
		}
	}
}

func (c *LwM2MClient) onRegistering(_ any) {
	if c.registrar.Registered() {
		log.Infoln("client registered")
		c.evtMgr.EmitEvent(EventClientRegistered)
		// registrar is long-running, so not stopped
		//c.registrar.Stop()
		c.enableService()
	} else {
		//restart registration if timeout
		if c.registrar.Timeout() {
			c.registrar.Stop()
			c.initiateBootstrap(bootstrapReasonRegFail)
			log.Infof("client register timeout, retry bootstrapping")
		}
	}
}

func (c *LwM2MClient) onServicing(_ any) {
	//log.Traceln("client is servicing")

	// check bootstrap first
	if c.bootstrapRequired() {
		log.Infoln("client arranged a new bootstrap")
		c.machine.MoveToState(initiating)
		return
	}

	// check registration
	if c.registerRequired() {
		log.Infoln("client arranged a new registration")
		c.machine.MoveToState(initiating)
		return
	}

	// checking health of components
	failed := c.reporter.FailureCounter()
	if failed > 3 {
		log.Errorf("client reported failure(%d) times exceed %d, "+
			"enter the re-registration process.", failed, 3)
		c.reporter.resetFailCounter()
		c.initiateRegister()
		return
	}
}

func (c *LwM2MClient) onExiting(_ any) {
	if err := c.registrar.Deregister(); err != nil {
		log.Errorln("client unregister failed:", err)
		return
	}

	log.Infoln("client is unregistered")
	c.evtMgr.EmitEvent(EventClientUnregistered)
	c.registrar.Stop()
}

func (c *LwM2MClient) makeDefaults() {
	if c.options.registry == nil {
		c.options.registry = NewObjectRegistry()
	}

	//if len(c.options.serverAddress) == 0 {
	//	c.options.serverAddress[0] = defaultServerAddr
	//}

	if len(c.options.localAddress) == 0 {
		c.options.localAddress = ":0"
	}

	if c.options.sendTimeout == 0 {
		c.options.sendTimeout = coap.DefaultTimeout
	}
	if c.options.recvTimeout == 0 {
		c.options.recvTimeout = coap.DefaultTimeout
	}
}

func (c *LwM2MClient) bootstrapRequired() bool {
	return c.bootstrapPending.Load()
}

func (c *LwM2MClient) registerRequired() bool {
	return c.registerPending.Load()
}

func (c *LwM2MClient) updateRequired() bool {
	return c.updatePending.Load()
}

func (c *LwM2MClient) initiateBootstrap(reason bootstrapReason) {
	c.bootstrapPending.Store(true)
	c.machine.MoveToState(initiating)
	log.Infof("initiating a bootstrap with reason: %d", reason)
}

func (c *LwM2MClient) initiateRegister() {
	// redundant clear
	c.clearBootstrapPending()
	c.registerPending.Store(true)
	c.machine.MoveToState(initiating)
}

func (c *LwM2MClient) clearBootstrapPending() {
	c.bootstrapPending.Store(false)
}

func (c *LwM2MClient) clearRegisterPending() {
	c.registerPending.Store(false)
}

// Start runs the client's state-driven loop.
func (c *LwM2MClient) Start() bool {
	return c.machine.Startup()
}

func (c *LwM2MClient) Stop() {
	//c.messager().Stop()
	c.machine.Shutdown()
	_ = c.store.StorageManager().Close()
}

func (c *LwM2MClient) Servicing() bool {
	return c.machine.GetState() == servicing
}

func (c *LwM2MClient) Notify(somebody string, something []byte) error {
	return c.reporter.Notify(somebody, something)
}

func (c *LwM2MClient) Send(data []byte) ([]byte, error) {
	return c.reporter.Send(data)
}

func (c *LwM2MClient) OnEvent(et EventType, h EventHandler) {
	c.evtMgr.AddListener(et, h)
}

func (c *LwM2MClient) SetOperator(oid ObjectID, operator Operator) {
	c.store.SetOperator(oid, operator)
}

func (c *LwM2MClient) SetOperators(operators OperatorMap) {
	c.store.SetOperators(operators)
}

func (c *LwM2MClient) EnableInstance(oid ObjectID, ids ...InstanceID) {
	c.store.EnableInstance(oid, ids...)
}

func (c *LwM2MClient) EnableInstances(m InstanceIdsMap) {
	c.store.EnableInstances(m)
}
