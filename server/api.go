package server

import (
	"github.com/zourva/lwm2m/coap"
	"github.com/zourva/lwm2m/core"
)

// Server defines api for application layer to use.
//
// Procedures, initiated by the client but not terminated
// within the LwM2M protocol layer, will be exposed to
// server side applications by methods defined here, including:
//
//	Send or Notify of Information Reporting
//
// Procedures, initiated by the server but not terminated
// within the LwM2M protocol layer, will be exposed to
// client side applications by methods defined here too:
//
//	Create/Read/Write/Execute of DeviceManagement & Service Enablement
//	Observe of Information Reporting
//
// EventType listeners are also supported to acquire client states including
// bootstrapping results, registration results etc.
type Server interface {
	Serve()
	Shutdown()
	GetClient(name string) core.RegisteredClient
	Listen(et core.EventType, h core.EventHandler)
}

type BootstrapService interface {
	Bootstrap(ctx BootstrapContext) error
	Bootstrapping(ctx BootstrapContext) error
	BootstrapPack(ctx BootstrapContext) ([]byte, error)
}

type RegistrationService interface {
	Register(info *core.RegistrationInfo) error
	Update(info *core.RegistrationInfo) error
	Unregister(info *core.RegistrationInfo)
}

type ReportingService interface {
	// Send invoked when info is received from send operation of reporting interface.
	Send(c core.RegisteredClient, data []byte) ([]byte, error)

	// Notify invoked when info is received from notify operation of reporting interface.
	Notify(c core.RegisteredClient, data []byte) ([]byte, error)
}

type DeviceControlService interface {
}

type Option func(s *LwM2MServer)

func WithBindingAddress(network, addr string) Option {
	return func(s *LwM2MServer) {
		s.network = network
		s.address = addr
	}
}

func WithClientEventObserver(observer RegisteredClientObserver) Option {
	return func(s *LwM2MServer) {
		s.observer = observer
	}
}

func WithGuidProvider(provider GuidProvider) Option {
	return func(s *LwM2MServer) {
		s.provider = provider
	}
}

func WithRegistrationInfoStore(store RegInfoStore) Option {
	return func(s *LwM2MServer) {
		s.store = store
	}
}

func WithObjectClassRegistry(registry core.ObjectRegistry) Option {
	return func(s *LwM2MServer) {
		s.registry = registry
	}
}

func WithSecurityConfig(kind coap.SecurityLayer, conf any) Option {
	return func(s *LwM2MServer) {
		s.secureLayer = kind
		s.secureConf = conf
	}
}
