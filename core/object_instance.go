package core

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zourva/pareto/endec/senml"
	"strconv"
)

type InstanceID = uint16

const (
	NoneID uint16 = 0xFFFF

	DefaultId uint16 = 0
)

// ObjectInstance defines an instance of
// an Object at runtime.
//
//	1 ObjectID -> 1 Object Instance Store
//	1 Object Instance Store -> 0/1/* Object Instances mapped by id
type ObjectInstance interface {
	Class() Object
	Construct() error //shortcut method
	Destruct() error  //shortcut method

	Id() InstanceID
	SetId(id InstanceID)

	Field(rid ResourceID, riId InstanceID) Field
	AddField(f Field)

	Fields(rid ResourceID) Fields
	SetFields(rid ResourceID, fields Fields)

	AllFields() map[ResourceID]Fields
	SetAllFields(all map[ResourceID]Fields)

	// SingleField equals Field(id, 0)
	SingleField(id ResourceID) Field

	// SetSingleField overwrites Field(id, 0)
	SetSingleField(f Field)

	AppendSENML(dst []senml.Record) []senml.Record
	MarshalJSON() ([]byte, error)
	String() string
}

func FieldValue[T bool | int | string | []byte](inst ObjectInstance, id ResourceID) T {
	f := inst.SingleField(id)
	if f != nil {
		return f.Get().(T)
	}

	var v T
	return v
}

type InstanceMap = map[InstanceID]ObjectInstance

type InstanceIdsMap = map[ObjectID][]InstanceID

type BaseInstance struct {
	class    Object     //object class
	instId   InstanceID //object instance id
	operator Operator   //copy from class definition

	resources map[ResourceID]Fields
}

var _ ObjectInstance = &BaseInstance{}

func NewObjectInstance(class Object) ObjectInstance {
	i := &BaseInstance{
		class:     class,
		instId:    0,
		operator:  class.Operator(),
		resources: make(map[ResourceID]Fields),
	}

	return i
}

func (o *BaseInstance) findInstance(instances Fields, iid InstanceID) Field {
	for _, i := range instances {
		if i.InstanceID() == iid {
			return i
		}
	}

	return nil
}

func (o *BaseInstance) Construct() error {
	return o.Class().Operator().Construct(o)
}

func (o *BaseInstance) Destruct() error {
	return o.Class().Operator().Destruct(o)
}

// SingleField returns instance 0 of a resource or nil if not exist.
func (o *BaseInstance) SingleField(id ResourceID) Field {
	if instances, ok := o.resources[id]; ok {
		inst := o.findInstance(instances, 0)
		if inst.Class().Multiple() {
			//
		}
		return inst
	}

	return nil
}

// Field returns instance of a resource or nil if not exist.
func (o *BaseInstance) Field(id ResourceID, iid InstanceID) Field {
	if instances, ok := o.resources[id]; ok {
		return o.findInstance(instances, iid)
	}

	return nil
}

func (o *BaseInstance) Fields(id ResourceID) Fields {
	return o.resources[id]
}

func (o *BaseInstance) AllFields() map[ResourceID]Fields {
	return o.resources
}

func (o *BaseInstance) Class() Object {
	return o.class
}

func (o *BaseInstance) Id() InstanceID {
	return o.instId
}

func (o *BaseInstance) SetId(id InstanceID) {
	o.instId = id
}

func (o *BaseInstance) AddField(f Field) {
	rid := f.Class().Id()
	fields, ok := o.resources[rid]
	if !ok {
		// add new field
		//err := o.operator.Add(o, rid, f.InstanceID(), f)
		//if err != nil {
		//	return
		//}

		fields = NewFields()
		o.resources[rid] = fields
		fields.Add(f)
	} else {
		// update field
		//err := o.operator.Update(o, rid, f.InstanceID(), f)
		//if err != nil {
		//	return
		//}

		fields.Update(f)
	}
}

func (o *BaseInstance) SetSingleField(f Field) {
	//if len(o.resources[f.Class().Id()]) == 0 {
	//	o.AddField(f)
	//} else {
	//	o.resources[f.Class().Id()].Add(f)
	//}
	o.AddField(f)
}

func (o *BaseInstance) SetFields(rid ResourceID, fields Fields) {
	o.resources[rid] = fields
}

func (o *BaseInstance) SetAllFields(all map[ResourceID]Fields) {
	o.resources = all
}

func (o *BaseInstance) String() string {
	tmp, _ := o.MarshalJSON()
	return string(tmp)
}

func (o *BaseInstance) AppendSENML(dst []senml.Record) []senml.Record {
	bname := GenBaseName(o)

	baseIdx := len(dst)
	for _, fields := range o.AllFields() {
		dst = fields.AppendSENML(dst)
	}

	if baseIdx < len(dst) {
		dst[baseIdx].BaseName = bname
	}
	return dst
}

func (o *BaseInstance) MarshalJSON() ([]byte, error) {
	var pack senml.Pack
	records := o.AppendSENML(nil)
	pack.Records = records

	return senml.Encode(pack, senml.JSON)
}

func GenBaseName(o ObjectInstance) string {
	oid, iid := o.Class().Id(), o.Id()
	bname := `/` + strconv.Itoa(int(oid)) + `/` + strconv.Itoa(int(iid)) + `/`
	return bname
}

func NewObjectInstance2(oid ObjectID, iid InstanceID, registry ObjectRegistry) (ObjectInstance, error) {
	class := registry.GetObject(oid)
	if class == nil {
		log.Errorf("unsupported object(%d)", oid)
		return nil, fmt.Errorf("unsupported object(%d)", oid)
	}

	// new objects
	instance := NewObjectInstance(class)
	instance.SetId(iid)
	return instance, nil
}

func ParseObjectInstancesWithJSON(registry ObjectRegistry, str string) ([]ObjectInstance, error) {
	var err error
	var normalize senml.Pack

	normalize, err = senml.DecodeAndNormalize([]byte(str), senml.JSON)
	if err != nil {
		log.Errorf("senml decode failed, err:%v", err)
		return nil, err
	}

	var ids []uint16
	var curObj ObjectInstance
	var objects []ObjectInstance
	for i := 0; i < len(normalize.Records); i++ {
		r := &normalize.Records[i]
		if ids, err = ParsePathToNumbers(r.Name, "/"); err != nil || len(ids) < 3 {
			return nil, fmt.Errorf("invalid path:%s, err:%v", r.Name, err)
		}

		oid, iid, fid, fiid := ids[0], ids[1], ids[2], uint16(0)
		if len(ids) > 3 {
			fiid = ids[3]
		}

		if curObj == nil || curObj.Class().Id() != oid || curObj.Id() != iid {
			if curObj != nil {
				err = curObj.Construct()
				if err != nil {
					return nil, err
				}
				// append
				objects = append(objects, curObj)
				curObj = nil
			}
			curObj, err = NewObjectInstance2(oid, iid, registry)
			if err != nil {
				return nil, err
			}
		}

		// add field
		res := curObj.Class().Resource(fid)
		val := SenmlRecordToFieldValue(res.Type(), r)

		field := NewResourceField2(curObj, fiid, res, val)
		curObj.AddField(field)
	}

	if curObj != nil {
		err = curObj.Construct()
		if err != nil {
			return nil, err
		}
		// append
		objects = append(objects, curObj)
	}
	return objects, nil
}
