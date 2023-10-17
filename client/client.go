package client

import (
	log "github.com/sirupsen/logrus"
	. "github.com/zourva/lwm2m/core"
	"github.com/zourva/lwm2m/preset"
	"github.com/zourva/pareto/box/meta"
	"strings"
	"sync/atomic"
	"time"
)

// Options defines client options.
type Options struct {
	persistence   ObjectPersistence
	factory       ObjectFactory
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

// WithObjectPersistence provides an object instance persistence
// layer accessor. If it is not provided, the factory-built-in
// version, i.e. initial version, of object instances are used.
func WithObjectPersistence(p ObjectPersistence) Option {
	return func(s *Options) {
		s.persistence = p
	}
}

// New returns an LwM2M client with the mandatory name
// and other options, or nil when any failure occurred.
func New(name string, opts ...Option) *LwM2MClient {
	c := &LwM2MClient{
		StateMachine: meta.NewStateMachine(name, time.Second),
		name:         name,
		opts:         &Options{},
	}

	for _, f := range opts {
		f(c.opts)
	}

	if err := c.initialize(); err != nil {
		log.Errorln("initialize client failed:", err)
		return nil
	}

	return c
}

type LwM2MClient struct {
	*meta.StateMachine
	opts *Options

	// name of endpoint
	name string

	// store to save object instances
	// loaded from local persistent storage.
	store *ObjectStore

	// messager to communicate with server
	messager *MessagerClient

	// delegators
	registrar  *Registrar
	controller *DeviceController
	reporter   *Reporter

	bootstrapPending atomic.Bool
}

func (c *LwM2MClient) initialize() error {
	c.makeDefaults()
	c.messager = NewMessager(c)
	c.registrar = NewRegistrar(c)
	c.reporter = NewReporter(c)
	c.RegisterStates([]*meta.State{
		{Name: initial, Action: c.onInitial},
		{Name: bootstrapping, Action: c.onBootstrapping},
		{Name: networking, Action: c.onNetworking},
		{Name: servicing, Action: c.onServicing},
		{Name: exiting, Action: c.onExiting},
	})

	c.SetStartingState(initial)
	c.SetStoppingState(exiting)

	c.store = NewObjectStore(c.opts.persistence, c.opts.factory)

	return c.store.Load()
}

func (c *LwM2MClient) RequestBootstrap(reason bootstrapReason) {
	c.bootstrapPending.Store(true)
}

func (c *LwM2MClient) onInitial(args any) {
	// determine bootstrap or registration procedure
	if c.bootstrapRequired() {
		//c.bootstrapper.Start()
	} else {
		if c.registrar.Startup() {
			log.Infoln("client is ready to register")
			c.MoveToState(networking)
		}
	}
}

func (c *LwM2MClient) onBootstrapping(args any) {
	//
}

func (c *LwM2MClient) onNetworking(args any) {
	if c.registrar.registered() {
		log.Infoln("client is registered")
		c.MoveToState(servicing)
	}
}

func (c *LwM2MClient) onServicing(args any) {
	//do nothing
	log.Traceln("client is servicing")
}

func (c *LwM2MClient) onExiting(args any) {
	err := c.registrar.Deregister()
	if err != nil {
		log.Errorln("client unregister failed, will try again:", err)
		return
	}

	if c.registrar.unregistered() {
		log.Infoln("client is unregistered")
		c.registrar.Shutdown()
	}
}

func (c *LwM2MClient) makeDefaults() {
	if c.opts.factory == nil {
		repo := NewClassStore(preset.NewOMAObjectInfoProvider())
		c.opts.factory = NewObjectFactory(repo)
	}

	if c.opts.persistence == nil {
		//support no persistence
		//use preset version each time
	}
}

func (c *LwM2MClient) bootstrapRequired() bool {
	return c.bootstrapPending.Load()
}

func (c *LwM2MClient) Registrar() RegistrationClient {
	return c.registrar
}

func (c *LwM2MClient) Reporter() ReportingClient {
	return c.reporter
}

// Start runs the client's state-driven loop.
func (c *LwM2MClient) Start() {
	c.messager.Start()
	c.Startup()
}

func (c *LwM2MClient) Stop() {
	//c.messager.Stop()
	c.Shutdown()
}

func (c *LwM2MClient) Name() string {
	return c.name
}
