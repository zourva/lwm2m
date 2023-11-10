package client

import (
	log "github.com/sirupsen/logrus"
	"github.com/zourva/lwm2m/coap"
	"github.com/zourva/lwm2m/coap2"
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
		machine: meta.NewStateMachine(name, time.Second),
	}

	c.bootstrapPending.Store(true)

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
	// name of endpoint, globally unique, assigned when provision
	name    string
	machine *meta.StateMachine
	options *Options

	// store to save object instances
	// loaded from local persistent storage.
	store ObjectInstanceStore

	coapConn coap.CoapServer

	// messager to communicate with server
	messager *MessagerClient

	// lifecycle event manager
	evtMgr *EventManager

	// delegators
	bootstrapper *Bootstrapper
	registrar    *Registrar
	controller   *DeviceController
	reporter     *InfoReporter

	bootstrapPending  atomic.Bool
	regRestartCounter int
}

func (c *LwM2MClient) initialize() error {
	c.makeDefaults()
	c.coapConn = coap2.NewServer(c.name, c.options.localAddress, c.options.serverAddress[0])
	c.messager = NewMessager(c)
	c.bootstrapper = NewBootstrapper(c)
	c.registrar = NewRegistrar(c)
	c.reporter = NewReporter(c)
	c.machine.RegisterStates([]*meta.State{
		{Name: initial, Action: c.onInitial},
		{Name: bootstrapping, Action: c.onBootstrapping},
		{Name: registering, Action: c.onRegistering},
		{Name: servicing, Action: c.onServicing},
		{Name: exiting, Action: c.onExiting},
	})

	c.machine.SetStartingState(initial)
	c.machine.SetStoppingState(exiting)

	c.evtMgr = NewEventManager()
	c.evtMgr.RegisterCreator(EventClientBootstrapped, NewBootstrappedEvent)
	c.evtMgr.RegisterCreator(EventClientRegistered, NewRegisteredEvent)
	c.evtMgr.RegisterCreator(EventClientRegUpdated, NewRegUpdatedEvent)
	c.evtMgr.RegisterCreator(EventClientUnregistered, NewUnregisteredEvent)
	c.evtMgr.RegisterCreator(EventClientDevInfoChanged, NewDeviceChangedEvent)
	c.evtMgr.RegisterCreator(EventClientObserved, NewInfoObservedEvent)
	c.evtMgr.RegisterCreator(EventClientReported, NewInfoReportedEvent)
	c.evtMgr.RegisterCreator(EventClientAbnormal, NewAbnormalEvent)

	if err := c.store.Load(); err != nil {
		log.Errorln("load object instances failed:", err)
		return err
	}

	c.bootstrapper.SetBootstrapServerBootstrapInfo(
		&BootstrapServerBootstrapInfo{
			BootstrapServerAccount: c.getBootstrapServerAccount(),
		},
	)

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
	if c.bootstrapper.Start() {
		log.Infoln("client is ready to bootstrap")
		// clear pending state
		c.clearBootstrapPending()

		// stop accept requests from servers
		c.messager.PauseAcceptRequests()

		// trap into bootstrapper
		c.machine.MoveToState(bootstrapping)
	} else {
		//try start again next iteration
		log.Errorln("start bootstrapper failed, will try again")
	}
}

func (c *LwM2MClient) doRegister() {
	if c.registrar.Start() {
		log.Infoln("client is ready to register")
		c.resetRegRestartCounter()

		// stop accept requests from servers
		c.messager.PauseAcceptRequests()

		c.machine.MoveToState(registering)
	} else {
		//try start again next iteration
	}
}

func (c *LwM2MClient) enableService() {
	c.messager.ResumeAcceptRequests()
	c.machine.MoveToState(servicing)
}

func (c *LwM2MClient) onInitial(args any) {
	// determine bootstrap or registration procedure
	if c.bootstrapRequired() {
		c.doBootstrap()
	} else {
		c.doRegister()
	}
}

func (c *LwM2MClient) onBootstrapping(args any) {
	if c.bootstrapper.bootstrapped() {
		log.Infoln("client is bootstrapped")
		c.evtMgr.EmitEvent(EventClientBootstrapped)
		c.bootstrapper.Stop()
		c.doRegister()
	} else {
		//restart bootstrap if timeout
		if c.bootstrapper.timeout() {
			c.bootstrapper.Stop()
			c.machine.MoveToState(initial)
			log.Infof("client bootstrap timeout, retry")
		}
	}
}

func (c *LwM2MClient) onRegistering(args any) {
	if c.registrar.registered() {
		log.Infoln("client is registered")
		c.evtMgr.EmitEvent(EventClientRegistered)
		// registrar is long-running, so not stopped
		c.enableService()
	} else {
		//restart registration if timeout
		if c.registrar.timeout() {
			c.registrar.Stop()
			c.machine.MoveToState(initial)
			log.Infof("client registration timeout, retry")
		}
	}
}

func (c *LwM2MClient) onServicing(args any) {
	log.Traceln("client is servicing")
	if c.regRestartCounter > 3 || c.bootstrapRequired() {
		log.Infoln("client arranged for a new register or bootstrap")
		c.machine.MoveToState(initial)
	}
}

func (c *LwM2MClient) onExiting(args any) {
	if err := c.registrar.Deregister(); err != nil {
		log.Errorln("client unregister failed:", err)
		return
	}

	if c.registrar.unregistered() {
		log.Infoln("client is unregistered")
		c.evtMgr.EmitEvent(EventClientUnregistered)
		c.registrar.Stop()
	}
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

func (c *LwM2MClient) requestBootstrap(reason bootstrapReason) {
	c.bootstrapPending.Store(true)
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
	c.store.StorageManager().Close()
}

func (c *LwM2MClient) Notify(data []byte) error {
	return nil
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

func (c *LwM2MClient) incRegRestartCounter() {
	c.regRestartCounter++
}

func (c *LwM2MClient) resetRegRestartCounter() {
	c.regRestartCounter = 0
}
