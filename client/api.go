package client

import "github.com/zourva/lwm2m/core"

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
