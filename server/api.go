package server

import (
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
	Register(info *core.RegistrationInfo) ([]byte, error)
	Update(info *core.RegistrationInfo) ([]byte, error)
	Unregister(info *core.RegistrationInfo) ([]byte, error)
}

type ReportingService interface {
	// Send invoked when info is received from send operation of reporting interface.
	Send(c core.RegisteredClient, data []byte) ([]byte, error)

	// Notify invoked when info is received from notify operation of reporting interface.
	Notify(c core.RegisteredClient, data []byte) ([]byte, error)
}

type DeviceControlService interface {
}

// Options defines options application can
// provide when creating a LwM2M server.
type Options struct {
	registry core.ObjectRegistry
	provider GuidProvider //
	store    RegInfoStore //registered client info store
	address  string       //binding address
	stats    Statistics
	observer ClientEventObserver
}

type Option func(*Options)

func WithBindingAddress(addr string) Option {
	return func(s *Options) {
		s.address = addr
	}
}

func WithClientEventObserver(observer ClientEventObserver) Option {
	return func(s *Options) {
		s.observer = observer
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

func WithObjectClassRegistry(registry core.ObjectRegistry) Option {
	return func(s *Options) {
		s.registry = registry
	}
}
