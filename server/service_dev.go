package server

import (
	"github.com/zourva/lwm2m/coap"
	. "github.com/zourva/lwm2m/core"
)

// DeviceControlDelegator implements the server-side operations
// for Device Management and Service Enablement Interface.
type DeviceControlDelegator struct {
	server  *LwM2MServer //server context
	service DeviceControlService
}

func (d *DeviceControlDelegator) Create(oid ObjectID, newValue Value) error {
	return nil
}

func (d *DeviceControlDelegator) Read(oid ObjectID, oiId InstanceID, rid ResourceID, riId InstanceID) ([]byte, error) {
	return nil, nil
}

func (d *DeviceControlDelegator) Write(oid ObjectID, instId InstanceID, resId ResourceID, resInstId InstanceID, newValue Value) error {
	return nil
}

func (d *DeviceControlDelegator) Delete(oid ObjectID, instId InstanceID, resId ResourceID, resInstId InstanceID) error {
	return nil
}

func (d *DeviceControlDelegator) Execute(oid ObjectID, instId InstanceID, resId ResourceID, args string) error {
	return nil
}

func (d *DeviceControlDelegator) Discover(oid ObjectID, instId InstanceID, resId ResourceID, depth int) ([]*coap.CoreResource, error) {
	return nil, nil
}

func NewDeviceControlServerDelegator(server *LwM2MServer) DeviceControlServer {
	return &DeviceControlDelegator{
		server:  server,
		service: server.deviceDelegator,
	}
}
