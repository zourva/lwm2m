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

// Options defines client options.
type Options struct {
	registry ObjectRegistry
	store    ObjectInstanceStore
	//provider      OperatorProvider
	//storage       InstanceStorageManager
	serverAddress []string
	localAddress  string
}

type Option func(*Options)

// WithLocalAddress provides local address as a hint.
// If not provided or the hinted address cannot be set
// the default address ":0" is used.
func WithLocalAddress(local string) Option {
	return func(s *Options) {
		s.localAddress = local
	}
}

// WithServerAddresses provides server list in an ";"
// separated string, e.g.:
//
//	1.0.0.1:5683;1.0.0.2:5683
//
// When not provided, "127.0.0.1:5683" is used.
func WithServerAddresses(addrString string) Option {
	return func(s *Options) {
		servers := strings.Split(addrString, ";")
		for _, server := range servers {
			s.serverAddress = append(s.serverAddress, server)
		}
	}
}

// WithObjectStore provides an object instance persistent
// layer accessor. If it is not provided, the default
// in-memory object instance store is used.
//func WithObjectStore(store ObjectInstanceStore) Option {
//	return func(s *Options) {
//		s.store = store
//	}
//}

func WithObjectClassRegistry(registry ObjectRegistry) Option {
	return func(s *Options) {
		s.registry = registry
	}
}

// New returns a LwM2M client with the mandatory name
// and other options, or nil when any failure occurred.
func New(name string, store ObjectInstanceStore, opts ...Option) *LwM2MClient {
	c := &LwM2MClient{
		name:    name,
		store:   store,
		options: &Options{},
		machine: meta.NewStateMachine[state](name, time.Second),
	}

	c.initiateBootstrap(bootstrapReasonStartup)

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

	coapConn coap.Server
	messager *MessagerClient // messager to communicate with server

	// lifecycle event manager
	evtMgr *EventManager

	// delegators
	bootstrapper *Bootstrapper // instance created for the latest try
	registrar    *Registrar    // instance created for the latest try
	controller   *DeviceController
	reporter     *Reporter

	bootstrapPending atomic.Bool
	registerPending  atomic.Bool
	updatePending    atomic.Bool
}

func (c *LwM2MClient) initialize() error {
	c.makeDefaults()
	c.coapConn = coap.NewServer(c.name, c.options.localAddress, c.options.serverAddress[0])
	c.messager = NewMessager(c)
	//c.bootstrapper = NewBootstrapper(c)
	//c.registrar = NewRegistrar(c)
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

	return nil
}

// get bootstrap server account from store, if any
func (c *LwM2MClient) getBootstrapServerAccount() *BootstrapServerAccount {
	instances := c.store.GetInstances(OmaObjectSecurity)
	for _, instance := range instances {
		for _, f := range instance.Fields(LwM2MSecurityBootstrapServer) {
			if f.Get().(bool) { // BootstrapServer == true
				return &BootstrapServerAccount{
					SecurityObjectInstance: instance,
				}
			}
		}
	}

	return nil
}

func (c *LwM2MClient) doBootstrap() {
	c.clearBootstrapPending()

	log.Infoln("client is ready to bootstrap")

	c.messager.PauseUserPlane()

	// always create a new bootstrapper
	c.bootstrapper = NewBootstrapper(c)
	c.bootstrapper.SetBootstrapServerBootstrapInfo(
		&BootstrapServerBootstrapInfo{
			BootstrapServerAccount: c.getBootstrapServerAccount(),
		},
	)
	c.bootstrapper.Start()

	c.machine.MoveToState(bootstrapping)
}

func (c *LwM2MClient) doRegister() {
	log.Infoln("client is ready to register")

	//c.resetReportFailCounter()
	c.messager.PauseUserPlane()

	// always create a new bootstrapper
	c.registrar = NewRegistrar(c)
	c.registrar.Start()

	c.machine.MoveToState(registering)
}

func (c *LwM2MClient) enableService() {
	c.messager.ResumeUserPlane()
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
	log.Traceln("client is servicing")

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
	if c.reporter.FailureCounter() > 3 {
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

	if len(c.options.serverAddress) == 0 {
		c.options.serverAddress[0] = defaultServerAddr
	}

	if len(c.options.localAddress) == 0 {
		c.options.localAddress = ":0"
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

// Start runs the client's state-driven loop.
func (c *LwM2MClient) Start() bool {
	c.messager.Start()
	return c.machine.Startup()
}

func (c *LwM2MClient) Stop() {
	//c.messager.Stop()
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
