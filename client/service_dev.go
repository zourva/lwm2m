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

// OnPack
// The LwM2M client MUST delete the existing Objects and their Instances in the client with the Objects and their
// Instances given in the "Bootstrap-Pack" if the Object IDs are the same. Object IDs available in the LwM2M Client which
// are not provided in the "Bootstrap-Pack" MUST NOT be deleted. The two exceptions are the LwM2M Bootstrap-Server
// Account, potentially including an associated Instance of an OSCORE Object ID:21, and the single Instance of the mandatory Device Object (ID:3), which are not affected by any Delete operation. Thus, when checking the configuration
// consistency, the LwM2M Client MUST ensure that the LwM2M Bootstrap-Server Account is still present.
func (d *DeviceController) OnPack(newValue []byte) error {
	return core.ForeachSenmlJSON(string(newValue), func(oid, iid, rid, fid uint16, r *senml.Record) error {
		mgr, err := d.client.store.GetInstanceManager(oid)
		if err != nil {
			return core.NotFound
		}

		instance := mgr.Get(iid)
		if instance == nil {
			instance, err = core.NewObjectInstance2(oid, iid, d.client.store.ObjectRegistry())
			if err != nil {
				log.Errorf("create object instance failed, %v", err)
				return core.InternalServerError
			}
			_ = mgr.Upsert(instance)
		}

		res := instance.Class().Resource(rid)
		val := core.SenmlRecordToFieldValue(res.Type(), r)
		field := core.NewResourceField2(instance, fid, res, val)
		instance.Class().Operator().Add(instance, rid, fid, field)
		return nil
	})
}

func (d *DeviceController) OnCreate(specifyOId core.ObjectID, newValue []byte) error {
	if specifyOId == core.NoneID {
		log.Errorf("create failed, the object id(%d) not specified", specifyOId)
		return core.BadRequest
	}

	return core.ForeachSenmlJSON(string(newValue), func(oid, iid, rid, fid uint16, r *senml.Record) error {
		if specifyOId != oid {
			log.Errorf("create failed, the oid(%d) is not specialed(%d)", oid, specifyOId)
			return core.BadRequest
		}

		mgr, err := d.client.store.GetInstanceManager(oid)
		if err != nil {
			return core.NotFound
		}

		instance := mgr.Get(iid)
		if instance == nil {
			instance, err = core.NewObjectInstance2(oid, iid, d.client.store.ObjectRegistry())
			if err != nil {
				log.Errorf("create object instance failed, %v", err)
				return core.InternalServerError
			}

			_ = mgr.Upsert(instance)
		}

		res := instance.Class().Resource(rid)
		val := core.SenmlRecordToFieldValue(res.Type(), r)
		field := core.NewResourceField2(instance, fid, res, val)
		instance.Class().Operator().Add(instance, rid, fid, field)
		return nil
	})
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

	objs, err := d.client.store.GetInstanceManager(oid)
	if err != nil {
		return nil, core.NotFound
	}

	if instId == core.NoneID {
		return d.errorConvert(objs.MarshalJSON())
	}

	inst := objs.Get(instId)
	if inst != nil {
		if resId == core.NoneID {
			return d.errorConvert(inst.MarshalJSON())
		}
		if resInstId == core.NoneID {
			res, err := inst.Class().Operator().GetAll(inst, resId)
			if err != nil {
				return nil, err
			}
			return d.errorConvert(res.MarshalJSON())
		}

		field, err := inst.Class().Operator().Get(inst, resId, resInstId)
		if err != nil {
			return d.errorConvert(field.MarshalJSON())
		}
	}

	return nil, core.NotFound
}

func (d *DeviceController) OnWrite(
	oid core.ObjectID,
	instId core.InstanceID,
	resId core.ResourceID,
	resInstId core.InstanceID,
	newValue []byte) error {
	if oid == core.NoneID || instId == core.NoneID {
		log.Errorf("write failed, invalid object id(%d) or instance id(%d)", oid, instId)
		return core.BadRequest
	}

	objs, err := d.client.store.GetInstanceManager(oid)
	if err != nil {
		return core.NotFound
	}

	instance := objs.Get(instId)
	if instance == nil {
		instance, err = core.NewObjectInstance2(oid, instId, d.client.store.ObjectRegistry())
		if err != nil {
			log.Errorf("write failed: %v", err)
			return core.NotImplemented
		}
	}

	return core.ForeachSenmlJSON(string(newValue), func(oid, iid, rid, fid uint16, r *senml.Record) error {
		if instance.Class().Id() != oid || instance.Id() != iid {
			log.Errorf("write failed: multiple oids or iids specified")
			return core.NotAcceptable
		}

		// add field
		res := instance.Class().Resource(rid)
		if res.Operations()&core.OpWrite != core.OpWrite {
			log.Errorf("write failed: %s", core.Forbidden)
			return core.Forbidden
		}

		val := core.SenmlRecordToFieldValue(res.Type(), r)
		field := core.NewResourceField2(instance, fid, res, val)
		return instance.Class().Operator().Add(instance, rid, fid, field)
	})
}

func (d *DeviceController) OnDelete(oid core.ObjectID, instId core.InstanceID, resId core.ResourceID, resInstId core.InstanceID) error {
	if oid == core.NoneID || instId == core.NoneID {
		log.Errorf("delete failed, invalid object id(%d) or instance id(%d)", oid, instId)
		return core.BadRequest
	}

	objs, err := d.client.store.GetInstanceManager(oid)
	if err != nil {
		return core.NotFound
	}

	instance := objs.Get(instId)
	if instance == nil {
		log.Warnf("delete failed: not found")
		return core.NotFound
	}

	err = instance.Class().Operator().Delete(instance, resId, resInstId)
	if err != nil {
		log.Warnf("delete failed:%v", err)
		return core.InternalServerError
	}

	if resId == core.NoneID || resInstId == core.NoneID {
		_ = objs.Delete(instId)
	}

	log.Debugf("delete(%d,%d) successfully", oid, instId)
	return err
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
