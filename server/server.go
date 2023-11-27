package server

import (
	log "github.com/sirupsen/logrus"
	"github.com/zourva/lwm2m/coap"
	. "github.com/zourva/lwm2m/core"
)

const (
	defaultAddress = ":5683"
)

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
	s.clientManager = NewRegisteredClientManager(s)
	s.bootstrapDelegator = NewBootstrapServerDelegator(s)
	s.registerDelegator = NewRegistrationServerDelegator(s)
	s.reportDelegator = NewReportingServerDelegator(s)
	s.deviceDelegator = NewDeviceControlServerDelegator(s)

	s.evtMgr = NewEventManager()
	s.evtMgr.RegisterCreator(EventServerStarted, NewServerStartedEvent)
	s.evtMgr.RegisterCreator(EventServerStopped, NewServerStoppedEvent)

	log.Infoln("lwm2m server created")

	return s
}

type LwM2MServer struct {
	name    string
	options *Options

	// session layer
	coapConn coap.Server
	messager *MessagerServer

	evtMgr *EventManager

	clientManager RegisteredClientManager

	// delegator layer
	bootstrapDelegator BootstrapServer
	registerDelegator  RegistrationServer
	reportDelegator    *ReportingServerDelegator
	deviceDelegator    *DeviceControlDelegator

	// service layer
	bootstrapService BootstrapService
	registerService  RegistrationService
	reportService    ReportingService
	controlService   DeviceControlService
}

func (s *LwM2MServer) EnableBootstrapService(bootstrapService BootstrapService) {
	s.bootstrapService = bootstrapService
}

func (s *LwM2MServer) EnableRegistrationService(registerService RegistrationService) {
	s.registerService = registerService
}

func (s *LwM2MServer) EnableReportingService(reportService ReportingService) {
	s.reportService = reportService
}

func (s *LwM2MServer) EnableDeviceControlService(controlService DeviceControlService) {
	s.controlService = controlService
}

func (s *LwM2MServer) Serve() {
	s.messager = NewMessager(s)
	s.messager.Start()
	s.clientManager.Start()
	s.evtMgr.EmitEvent(EventServerStarted)
	log.Infoln("lwm2m server started")
}

// Shutdown shuts down the server gracefully.
func (s *LwM2MServer) Shutdown() {
	s.clientManager.Stop()
	s.messager.Stop()
	s.evtMgr.EmitEvent(EventServerStopped)
	log.Infoln("lwm2m server stopped")
}

func (s *LwM2MServer) GetClient(name string) RegisteredClient {
	return s.clientManager.Get(name)
}

func (s *LwM2MServer) Listen(et EventType, h EventHandler) {
	s.evtMgr.AddListener(et, h)
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

	//if s.options.stats == nil {
	//	s.options.stats = &DefaultStatistics{}
	//}

	if s.options.observer == nil {
		s.options.observer = NewDefaultEventObserver()
	}

	if s.options.provider == nil {
		s.options.provider = NewUrnUuidProvider()
	}
}
