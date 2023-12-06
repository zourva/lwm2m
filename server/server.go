package server

import (
	piondtls "github.com/pion/dtls/v2"
	log "github.com/sirupsen/logrus"
	. "github.com/zourva/lwm2m/core"
)

const (
	defaultNetwork = "udp"
	defaultAddress = ":5683"
)

func New(name string, opts ...Option) *LwM2MServer {
	s := &LwM2MServer{
		name: name,
	}

	for _, f := range opts {
		f(s)
	}

	s.makeDefaults()
	//s.coapConn = coap.NewServer(name, s.address)
	s.stats = NewStatistics()
	s.manager = NewRegisteredClientManager(s)
	s.bootstrapDelegator = NewBootstrapServerDelegator(s)
	s.registerDelegator = NewRegistrationServerDelegator(s)
	s.reportDelegator = NewReportingServerDelegator(s)
	//s.deviceDelegator = NewDeviceControlServerDelegator(s)

	//s.evtMgr = NewEventManager()
	//s.evtMgr.RegisterCreator(EventServerStarted, NewServerStartedEvent)
	//s.evtMgr.RegisterCreator(EventServerStopped, NewServerStoppedEvent)

	log.Infoln("lwm2m server created")

	return s
}

type LwM2MServer struct {
	name     string
	network  string
	address  string
	registry ObjectRegistry
	store    RegInfoStore
	provider GuidProvider
	dtlsConf *piondtls.Config

	observer RegisteredClientObserver
	manager  RegisteredClientManager

	// delegator layer
	bootstrapDelegator *BootstrapServerDelegator
	registerDelegator  *RegisterServiceDelegator
	reportDelegator    *ReportingServerDelegator
	//deviceDelegator    *DeviceControlDelegator

	// service layer
	bootstrapService BootstrapService
	registerService  RegistrationService
	reportService    ReportingService
	//controlService   DeviceControlService

	// session layer
	//coapConn coap.Server
	messager *MessagerServer
	stats    *Statistics
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

//func (s *LwM2MServer) EnableDeviceControlService(controlService DeviceControlService) {
//	s.controlService = controlService
//}

func (s *LwM2MServer) Serve() {
	s.messager = NewMessager(s)
	s.messager.Start()
	s.manager.Start()
	//s.evtMgr.EmitEvent(EventServerStarted)
	log.Infoln("lwm2m server started")
}

// Shutdown shuts down the server gracefully.
func (s *LwM2MServer) Shutdown() {
	s.manager.Stop()
	s.messager.Stop()
	//s.evtMgr.EmitEvent(EventServerStopped)
	log.Infoln("lwm2m server stopped")
}

func (s *LwM2MServer) GetClient(name string) RegisteredClient {
	return s.manager.Get(name)
}

func (s *LwM2MServer) Listen(et EventType, h EventHandler) {
	//s.evtMgr.AddListener(et, h)
}

func (s *LwM2MServer) makeDefaults() {
	if len(s.network) == 0 {
		s.network = defaultNetwork
	}

	if len(s.address) == 0 {
		s.address = defaultAddress
	}

	if s.registry == nil {
		s.registry = NewObjectRegistry()
	}

	if s.store == nil {
		s.store = NewInMemorySessionStore()
	}

	//if s.stats == nil {
	//	s.stats = &Statistics{}
	//}

	if s.observer == nil {
		s.observer = NewDefaultEventObserver()
	}

	if s.provider == nil {
		s.provider = NewUrnUuidProvider()
	}

}
