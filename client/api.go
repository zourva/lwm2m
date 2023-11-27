package client

import (
	"github.com/zourva/lwm2m/core"
	"strings"
	"time"
)

// Client defines api for application layer to use.
//
// Procedures, initiated by the client side applications
// will be exposed, including:
//
//	Send of Information Reporting
//
// Procedures, initiated by the server but not terminated
// within the LwM2M protocol layer, will be exposed to
// client side applications by methods defined here too:
//
//	Create/Read/Write/Execute of DeviceManagement & Service Enablement
//	Observe of Information Reporting
//
// Event listeners are also supported to acquire client states including
// bootstrapping results, registration results etc.
type Client interface {
	// OnEvent adds an event listener for the specified
	// event, and overwrites the old if already exists.
	OnEvent(et core.EventType, h core.EventHandler)

	Send(data []byte) ([]byte, error)
	Notify(somebody string, something []byte) error
}

// Options defines client options.
type Options struct {
	registry core.ObjectRegistry
	store    core.ObjectInstanceStore
	//provider      OperatorProvider
	//storage       InstanceStorageManager
	serverAddress []string
	localAddress  string
	//dtlsConf      *piondtls.Config

	// send timeout
	timeout time.Duration
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

func WithServerSendTimeout(timeout time.Duration) Option {
	return func(s *Options) {
		s.timeout = timeout
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

func WithObjectClassRegistry(registry core.ObjectRegistry) Option {
	return func(s *Options) {
		s.registry = registry
	}
}
