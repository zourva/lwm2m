package client

import "github.com/zourva/lwm2m/core"

type DeviceController struct {
}

var _ core.DeviceControlClient = &DeviceController{}

func NewDeviceController() *DeviceController {
	return &DeviceController{}
}

func (d *DeviceController) validateOID(oid core.ObjectID) bool {
	return true
}

func (d *DeviceController) validateOIID(oid core.ObjectID, oiId core.InstanceID) bool {
	return true
}

func (d *DeviceController) validateRID(rid core.ResourceID) bool {
	return true
}

func (d *DeviceController) validateRIID(rid core.ResourceID, riId core.InstanceID) bool {
	return true
}

func (d *DeviceController) preCheck(oid core.ObjectID, oiId core.InstanceID, rid core.ResourceID, riId core.InstanceID) error {
	//check existence

	//check access control

	return core.ErrorNone
}

func (d *DeviceController) OnCreate(oid core.ObjectID, newValue core.Value) error {
	//TODO implement me
	panic("implement me")
}

func (d *DeviceController) OnRead(oid core.ObjectID, instId core.InstanceID, resId core.ResourceID, resInstId core.InstanceID) (*core.ResourceField, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DeviceController) OnWrite(oid core.ObjectID, instId core.InstanceID, resId core.ResourceID, resInstId core.InstanceID, newValue core.Value) error {
	//TODO implement me
	panic("implement me")
}

func (d *DeviceController) OnDelete(oid core.ObjectID, instId core.InstanceID, resId core.ResourceID, resInstId core.InstanceID) error {
	//TODO implement me
	panic("implement me")
}

func (d *DeviceController) OnExecute(oid core.ObjectID, instId core.InstanceID, resId core.ResourceID, args string) error {
	// check if executable

	//TODO implement me
	panic("implement me")
}

func (d *DeviceController) OnDiscover(oid core.ObjectID, instId core.InstanceID, resId core.ResourceID, depth int) error {
	//TODO implement me
	panic("implement me")
}
