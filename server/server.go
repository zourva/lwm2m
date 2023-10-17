package server

import (
	log "github.com/sirupsen/logrus"
	"github.com/zourva/lwm2m/coap"
	. "github.com/zourva/lwm2m/core"
	"github.com/zourva/lwm2m/preset"
)

const (
	defaultAddress = ":5683"
)

type Options struct {
	provider GuidProvider
	store    RegInfoStore
	factory  ObjectFactory
	stats    Statistics

	address string //binding address

	lcHandler LifecycleHandler
}

type Option func(*Options)

func WithBindingAddress(addr string) Option {
	return func(s *Options) {
		s.address = addr
	}
}

func WithLifecycleHandler(lh LifecycleHandler) Option {
	return func(s *Options) {
		s.lcHandler = lh
	}
}

func WithGuidProvider(provider GuidProvider) Option {
	return func(s *Options) {
		s.provider = provider
	}
}

func WithRegistrationInfoStore(store RegInfoStore) Option {
	return func(s *Options) {
		s.store = store
	}
}

func WithObjectFactory(factory ObjectFactory) Option {
	return func(s *Options) {
		s.factory = factory
	}
}

func New(name string, opts ...Option) *LwM2MServer {
	s := &LwM2MServer{
		name: name,
		opts: &Options{},
	}

	for _, f := range opts {
		f(s.opts)
	}

	s.makeDefaults()
	s.coapConn = coap.NewCoapServer(name, s.opts.address)
	s.manager = NewSessionManager(s)
	s.messager = NewMessageHandler(s)

	log.Infoln("lwm2m server created")

	return s
}

type LwM2MServer struct {
	name string
	opts *Options

	manager RegisteredClientManager

	// application layer
	registrationService RegistrationServer
	deviceMgmtService   DeviceControlProxy

	// session layer
	coapConn coap.CoapServer
	messager *Messager
}

func (s *LwM2MServer) Serve() {
	// setup hooks
	s.coapConn.OnMessage(func(msg *coap.Message, inbound bool) {
		s.opts.stats.IncrementRequestCount()
	})

	// register route handlers
	s.coapConn.Post("/rd", s.messager.onClientRegister)
	s.coapConn.Put("/rd/:id", s.messager.onClientUpdate)
	s.coapConn.Delete("/rd/:id", s.messager.onClientDeregister)

	go s.coapConn.Start()

	log.Infoln("lwm2m server started at", s.coapConn.GetLocalAddress().String())
}

// Shutdown shuts down the server gracefully.
func (s *LwM2MServer) Shutdown() {
	s.coapConn.Stop()
	//s.ClearSessions()

	log.Infoln("lwm2m server stopped")
}

func (s *LwM2MServer) GetServerStats() Statistics {
	return s.opts.stats
}

func (s *LwM2MServer) makeDefaults() {
	if len(s.opts.address) == 0 {
		s.opts.address = defaultAddress
	}

	if s.opts.store == nil {
		s.opts.store = NewInMemorySessionStore()
	}

	if s.opts.factory == nil {
		repo := NewClassStore(preset.NewOMAObjectInfoProvider())
		s.opts.factory = NewObjectFactory(repo)
	}

	if s.opts.stats == nil {
		s.opts.stats = &DefaultStatistics{}
	}

	if s.opts.lcHandler == nil {
		s.opts.lcHandler = NewDefaultLifecycleHandler()
	}

	if s.opts.provider == nil {
		s.opts.provider = NewUrnUuidProvider()
	}
}
