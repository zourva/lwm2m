package server

import "github.com/zourva/lwm2m/core"

// API defines methods for application layer to use.
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
type API interface {
	GetClient(name string) core.RegisteredClient

	// OnEvent adds an event listener.
	Listen(et core.EventType, h core.EventHandler)

	// SetOnInfoSent sets the callback to be invoked
	// when info is received from Send operation of reporting interface.
	SetOnInfoSent(handler InfoReportHandler)

	// SetOnInfoNotified sets the callback to be invoked
	// when info is received from Notify operation of reporting interface.
	SetOnInfoNotified(handler InfoReportHandler)

	SetOnClientRegistered(handler ClientRegHandler)
	SetOnClientRegUpdated(handler ClientRegHandler)
	SetOnClientUnregistered(handler ClientRegHandler)
}
