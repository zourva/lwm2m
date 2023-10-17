package server

import . "github.com/zourva/lwm2m/core"

type DeviceControlServer interface {
	Create(client *RegisteredClient, oid ObjectID, newValue Value) error
	Read(client *RegisteredClient, oid ObjectID, instId InstanceID, resId ResourceID, resInstId InstanceID) error
	Write(client *RegisteredClient, oid ObjectID, instId InstanceID, resId ResourceID, resInstId InstanceID, newValue Value) error
	Delete(client *RegisteredClient, oid ObjectID, instId InstanceID, resId ResourceID, resInstId InstanceID) error
	Execute(client *RegisteredClient, oid ObjectID, instId InstanceID, resId ResourceID, args string) error
	Discover(client *RegisteredClient, oid ObjectID, instId InstanceID, resId ResourceID, depth int) error
	//ReadComposite()
	//WriteComposite()
	//WriteAttributes()
}

// DeviceControlService implements the server-side operations
// for Device Management and Service Enablement Interface.
type DeviceControlService struct {
	clientMgr RegisteredClientManager
}

func (d *DeviceControlService) Create(client *RegisteredClient, oid ObjectID, newValue Value) error {
	return client.Create(oid, newValue)
}

func (d *DeviceControlService) Read(client *RegisteredClient, oid ObjectID,
	instId InstanceID, resId ResourceID, resInstId InstanceID) error {
	return client.Read(oid, instId, resId, resInstId)
}

func (d *DeviceControlService) Write(client *RegisteredClient, oid ObjectID,
	instId InstanceID, resId ResourceID, resInstId InstanceID, newValue Value) error {
	return client.Write(oid, instId, resId, resInstId, newValue)
}

func (d *DeviceControlService) Delete(client *RegisteredClient, oid ObjectID,
	instId InstanceID, resId ResourceID, resInstId InstanceID) error {
	return client.Delete(oid, instId, resId, resInstId)
}

func (d *DeviceControlService) Execute(client *RegisteredClient, oid ObjectID,
	instId InstanceID, resId ResourceID, args string) error {
	return client.Execute(oid, instId, resId, args)
}

func (d *DeviceControlService) Discover(client *RegisteredClient, oid ObjectID,
	instId InstanceID, resId ResourceID, depth int) error {
	return client.Discover(oid, instId, resId, depth)
}

func NewDeviceControlService(server *LwM2MServer) DeviceControlServer {
	return &DeviceControlService{
		clientMgr: server.manager,
	}
}
