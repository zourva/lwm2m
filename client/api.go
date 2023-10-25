package client

import "github.com/zourva/lwm2m/core"

// API defines methods for application layer to use.
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
type API interface {
	// OnEvent adds an event listener for the specified
	// event, and overwrites the old if already exists.
	OnEvent(et core.EventType, h core.EventHandler)

	Send(data []byte) ([]byte, error)

	// SetBootstrapServerAccount set the pre-provisioned bootstrap
	// server account as depicted:
	// In order for the LwM2M Client and the LwM2M Bootstrap-Server
	// to establish a connection on the Bootstrap Interface, either in
	// Client Initiated Bootstrap mode or in Server Initiated Bootstrap
	// mode, the LwM2M Client MUST have an LwM2M Bootstrap-Server Account pre-provisioned.
	SetBootstrapServerAccount(account *core.BootstrapServerAccount)
}
