package server

import (
	log "github.com/sirupsen/logrus"
	"github.com/zourva/lwm2m/coap"
	. "github.com/zourva/lwm2m/core"
)

const (
	defaultAddress = ":5683"
)

type Options struct {
	registry ObjectRegistry
	provider GuidProvider //
	store    RegInfoStore //registered client info store
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

func WithObjectClassRegistry(registry ObjectRegistry) Option {
	return func(s *Options) {
		s.registry = registry
	}
}

func New(name string, opts ...Option) *LwM2MServer {
	s := &LwM2MServer{
		name:    name,
		options: &Options{},
	}

	for _, f := range opts {
		f(s.options)
	}

	s.makeDefaults()
	s.coapConn = coap.NewCoapServer(name, s.options.address)
	s.registerManager = NewSessionManager(s)
	s.messager = NewMessager(s)

	s.evtMgr = NewEventManager()
	s.evtMgr.RegisterCreator(EventServerStarted, NewServerStartedEvent)
	s.evtMgr.RegisterCreator(EventServerStopped, NewServerStoppedEvent)

	log.Infoln("lwm2m server created")

	return s
}

type ClientBootstrapHandler = func(ctx BootstrapContext) error
type ClientRegHandler = func(info *RegistrationInfo) ([]byte, error)
type InfoReportHandler = func(c RegisteredClient, data []byte) ([]byte, error)

type LwM2MServer struct {
	name    string
	options *Options

	// session layer
	coapConn coap.CoapServer
	messager *ServerMessager

	evtMgr *EventManager

	onSent         InfoReportHandler
	onNotified     InfoReportHandler
	onRegistered   ClientRegHandler
	onRegUpdated   ClientRegHandler
	onUnregistered ClientRegHandler

	onBootstrapInit ClientBootstrapHandler
	onBootstrapping ClientBootstrapHandler

	registerManager RegisterManager
}

func (s *LwM2MServer) Serve() {
	// setup hooks
	s.coapConn.OnMessage(func(msg *coap.Message, inbound bool) {
		s.options.stats.IncrementRequestCount()
	})

	// register route handlers
	s.coapConn.Post("/bs", s.messager.onClientBootstrap)
	s.coapConn.Get("/bspack", s.messager.onClientBootstrapPack)

	s.coapConn.Post("/rd", s.messager.onClientRegister)
	s.coapConn.Put("/rd/:id", s.messager.onClientUpdate)
	s.coapConn.Delete("/rd/:id", s.messager.onClientDeregister)

	s.coapConn.Post("/dp", s.messager.onSendInfo)

	go s.coapConn.Start()

	s.evtMgr.EmitEvent(EventServerStarted)

	log.Infoln("lwm2m server started at", s.coapConn.GetLocalAddress().String())
}

// Shutdown shuts down the server gracefully.
func (s *LwM2MServer) Shutdown() {
	s.coapConn.Stop()
	//s.ClearSessions()

	s.evtMgr.EmitEvent(EventServerStopped)

	log.Infoln("lwm2m server stopped")
}

func (s *LwM2MServer) GetClientManager() RegisterManager {
	return s.registerManager
}

func (s *LwM2MServer) GetClient(name string) RegisteredClient {
	return s.registerManager.Get(name)
}

func (s *LwM2MServer) Listen(et EventType, h EventHandler) {
	s.evtMgr.AddListener(et, h)
}

func (s *LwM2MServer) SetOnInfoSent(handler InfoReportHandler) {
	s.onSent = handler
}

func (s *LwM2MServer) SetOnInfoNotified(handler InfoReportHandler) {
	s.onNotified = handler
}

func (s *LwM2MServer) SetOnClientBootstrapInit(handler ClientBootstrapHandler) {
	s.onBootstrapInit = handler
}

func (s *LwM2MServer) SetOnClientBootstrapping(handler ClientBootstrapHandler) {
	s.onBootstrapping = handler
}

func (s *LwM2MServer) SetOnClientRegistered(handler ClientRegHandler) {
	s.onRegistered = handler
}

func (s *LwM2MServer) SetOnClientRegUpdated(handler ClientRegHandler) {
	s.onRegUpdated = handler
}

func (s *LwM2MServer) SetOnClientUnregistered(handler ClientRegHandler) {
	s.onUnregistered = handler
}

func (s *LwM2MServer) makeDefaults() {
	if len(s.options.address) == 0 {
		s.options.address = defaultAddress
	}

	if s.options.registry == nil {
		s.options.registry = NewObjectRegistry()
	}

	if s.options.store == nil {
		s.options.store = NewInMemorySessionStore()
	}

	if s.options.stats == nil {
		s.options.stats = &DefaultStatistics{}
	}

	if s.options.lcHandler == nil {
		s.options.lcHandler = NewDefaultLifecycleHandler()
	}

	if s.options.provider == nil {
		s.options.provider = NewUrnUuidProvider()
	}
}
