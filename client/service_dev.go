package client

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zourva/lwm2m/core"
	"github.com/zourva/pareto/endec/senml"
)

type DeviceController struct {
	client *LwM2MClient //lwm2m context
}

var (
	errInvalidObjectId = fmt.Errorf("invalid oid")
	errNotFound        = fmt.Errorf("not found")
	errNotExists       = fmt.Errorf("not exists")
	errNoPermission    = fmt.Errorf("no permission")
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

func (d *DeviceController) errorConvert(value []byte, err error) ([]byte, error) {
	if err == nil {
		return value, err
	}
	return value, core.InternalServerError
}

func (d *DeviceController) OnRead(oid core.ObjectID, instId core.InstanceID, resId core.ResourceID, resInstId core.InstanceID) ([]byte, error) {
	if oid == core.NoneID {
		log.Errorf("read failed, invalid oid:%d", oid)
		return nil, core.BadRequest
	}

	objs := d.client.store.GetInstanceManager(oid)
	if objs != nil {
		if instId == core.NoneID {
			return d.errorConvert(objs.MarshalJSON())
		}

		inst := objs.Get(instId)
		if inst != nil {
			if resId == core.NoneID {
				return d.errorConvert(inst.MarshalJSON())
			}
			res := inst.Fields(resId)
			if res != nil {
				if resInstId == core.NoneID {
					return d.errorConvert(res.MarshalJSON())
				}

				field := res.Field(resInstId)
				if field != nil {
					return d.errorConvert(field.MarshalJSON())
				}
			}
		}
	}

	return nil, core.NotFound
}

func (d *DeviceController) OnWrite(oid core.ObjectID, instId core.InstanceID, resId core.ResourceID, resInstId core.InstanceID, newValue []byte) error {
	if oid == core.NoneID || instId == core.NoneID {
		log.Errorf("write failed, invalid object id(%d) or instance id(%d)", oid, instId)
		return core.BadRequest
	}

	var err error
	var normalize senml.Pack
	objmgr := d.client.store.GetInstanceManager(oid)
	instance := objmgr.Get(instId)
	if instance == nil {
		instance, err = core.NewObjectInstance2(oid, instId, d.client.store.ObjectRegistry())
		if err != nil {
			log.Errorf("write failed: %v", err)
			return core.NotImplemented
		}
	}

	normalize, err = senml.DecodeAndNormalize(newValue, senml.JSON)
	if err != nil {
		log.Errorf("write failed: %v", err)
		return core.BadRequest
	}

	var ids []uint16
	for i := 0; i < len(normalize.Records); i++ {
		r := &normalize.Records[i]
		if ids, err = core.ParsePathToNumbers(r.Name, "/"); err != nil || len(ids) < 3 {
			log.Errorf("write failed: invalid path:%s, err:%v", r.Name, err)
			return core.BadRequest
		}

		xoid, xiid, xrid, xriid := ids[0], ids[1], ids[2], uint16(0)
		if len(ids) > 3 {
			xriid = ids[3]
		}

		if instance.Class().Id() != xoid || instance.Id() != xiid {
			log.Errorf("write failed: multiple oids or iids specified")
			return core.NotAcceptable
		}

		// add field
		res := instance.Class().Resource(xrid)
		val := core.SenmlRecordToFieldValue(res.Type(), r)

		field := core.NewResourceField2(instance, xriid, res, val)
		instance.AddField(field)
	}

	objmgr.Add(instance)
	err = d.client.store.Flush()

	return err
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
