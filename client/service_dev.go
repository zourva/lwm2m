package client

import (
	"fmt"
	"github.com/zourva/lwm2m/core"
)

type DeviceController struct {
	client *LwM2MClient //lwm2m context
}

var (
	errInvalidObjectId = fmt.Errorf("invalid oid")
	errNotFound        = fmt.Errorf("not found")
	errNotExists       = fmt.Errorf("not exists")
)

var _ core.DeviceControlClient = &DeviceController{}

func NewDeviceController(c *LwM2MClient) *DeviceController {
	return &DeviceController{
		client: c,
	}
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

	//registry := d.client.store.ObjectRegistry()
	//cls := registry.GetObject(oid)
	//obj := core.NewObjectInstance(cls)
	//obj.SetId(oid)
	//
	//res := obj.Class().Resource(fid)
	//
	//core.NewResourceField2(fiid, res, val)
	//obj.AddField()
	//
	////obj := core.NewObjectInstance()
	//d.client.store.GetInstanceManager(oid).Add()

	return nil
}

func (d *DeviceController) OnRead(oid core.ObjectID, instId core.InstanceID, resId core.ResourceID, resInstId core.InstanceID) ([]byte, error) {
	//TODO implement me
	//panic("implement me")

	if oid == core.NoneID {
		return nil, errInvalidObjectId
	}

	objs := d.client.store.GetInstanceManager(oid)
	if objs != nil {
		if instId == core.NoneID {
			return objs.MarshalJSON()
		}

		inst := objs.Get(instId)
		if inst != nil {
			if resId == core.NoneID {
				return inst.MarshalJSON()
			}
			res := inst.Fields(resId)
			if res != nil {
				if resInstId == core.NoneID {
					return res.MarshalJSON()
				}

				field := res.Field(resInstId)
				if field != nil {
					return field.MarshalJSON()
				}
			}
		}
	}

	return nil, errNotExists
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
