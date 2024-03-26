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

type ObjectInstanceHelper interface {
	Field(rid ResourceID, riId InstanceID) Field
	AddField(f Field)
	DelField(rid ResourceID, riId InstanceID)

	Fields(rid ResourceID) *Fields
	SetFields(rid ResourceID, fields *Fields)

	AllFields() *Resources
	SetAllFields(all *Resources)

	// SingleField equals Field(id, 0)
	SingleField(id ResourceID) Field

	// SetSingleField overwrites Field(id, 0)
	SetSingleField(f Field)
}

type Senmler interface {
	AppendSENML(dst []senml.Record) []senml.Record
}

type Marshaler interface {
	MarshalJSON() ([]byte, error)
	String() string
}

// ObjectInstance defines an instance of
// an Object at runtime.
//
//	1 ObjectID -> 1 Object Instance Store
//	1 Object Instance Store -> 0/1/* Object Instances mapped by id
type ObjectInstance interface {
	Class() Object
	Construct() error //shortcut method
	Destruct() error  //shortcut method

	Helper() ObjectInstanceHelper

	Id() InstanceID
	SetId(id InstanceID)

	Senmler
	Marshaler
}

func FieldValue[T bool | int | string | []byte](inst ObjectInstance, id ResourceID) T {
	f := inst.Helper().SingleField(id)
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

	resources *Resources
}

var _ ObjectInstance = &BaseInstance{}

func NewObjectInstance(class Object) ObjectInstance {
	i := &BaseInstance{
		class:    class,
		instId:   0,
		operator: class.Operator(),
	}

	i.resources = NewResources(i)
	return i
}

func (o *BaseInstance) Construct() error             { return o.Class().Operator().Construct(o) }
func (o *BaseInstance) Destruct() error              { return o.Class().Operator().Destruct(o) }
func (o *BaseInstance) Helper() ObjectInstanceHelper { return o }

// SingleField returns instance 0 of a resource or nil if not exist.
func (o *BaseInstance) SingleField(id ResourceID) Field {
	fields := o.resources.Fields(id)
	if fields != nil {
		field := fields.SingleField()
		if field.Class().Multiple() {
			//
		}
		return field
	}

	return nil
}

// Field returns instance of a resource or nil if not exist.
func (o *BaseInstance) Field(id ResourceID, iid InstanceID) Field {
	fields := o.resources.Fields(id)
	if fields != nil {
		return fields.Field(iid)
	}

	return nil
}

func (o *BaseInstance) DelField(rid ResourceID, riId InstanceID) {
	fields := o.resources.Fields(rid)
	if fields != nil {
		fields.Delete(riId)
	}
}

func (o *BaseInstance) Fields(id ResourceID) *Fields {
	fields := o.resources.Fields(id)
	return fields
}

func (o *BaseInstance) AllFields() *Resources {
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
	fields := o.resources.Fields(rid)
	if fields == nil {
		// add new field
		//err := o.operator.Add(o, rid, f.InstanceID(), f)
		//if err != nil {
		//	return
		//}

		fields = NewFields(o)
		fields.Add(f)
		o.resources.Add(rid, fields)
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

func (o *BaseInstance) SetFields(rid ResourceID, fields *Fields) { o.resources.Update(rid, fields) }
func (o *BaseInstance) SetAllFields(all *Resources)              { o.resources = all }

func (o *BaseInstance) String() string {
	tmp, _ := o.MarshalJSON()
	return string(tmp)
}

func (o *BaseInstance) AppendSENML(dst []senml.Record) []senml.Record {
	bname := GenBaseName(o)

	baseIdx := len(dst)
	//resources, err := o.operator.GetAll(o, NoneID)
	//if err != nil {
	//	return dst
	//}
	resources := o.resources
	dst = resources.AppendSENML(dst)

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

func ForeachSenmlJSON(strJson string, iter func(oid, iid, rid, fid uint16, value *senml.Record) error) error {
	var err error
	var normalize senml.Pack

	normalize, err = senml.DecodeAndNormalize([]byte(strJson), senml.JSON)
	if err != nil {
		log.Errorf("senml decode failed, err:%v", err)
		return err
	}

	var ids []uint16
	for i := 0; i < len(normalize.Records); i++ {
		r := &normalize.Records[i]
		if ids, err = ParsePathToNumbers(r.Name, "/"); err != nil || len(ids) < 3 {
			return fmt.Errorf("invalid path:%s, err:%v", r.Name, err)
		}

		oid, iid, rid, fid := ids[0], ids[1], ids[2], uint16(0)
		if len(ids) > 3 {
			fid = ids[3]
		}

		err = iter(oid, iid, rid, fid, r)
		if err != nil {
			return err
		}
	}

	return nil
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

		oid, iid, rid, riid := ids[0], ids[1], ids[2], uint16(0)
		if len(ids) > 3 {
			riid = ids[3]
		}

		if curObj == nil || curObj.Class().Id() != oid || curObj.Id() != iid {
			if curObj != nil {
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
		res := curObj.Class().Resource(rid)
		val := SenmlRecordToFieldValue(res.Type(), r)

		field := NewResourceField2(curObj, riid, res, val)
		//curObj.AddField(field)
		curObj.Class().Operator().Add(curObj, rid, riid, field)
	}

	if curObj != nil {
		// append
		objects = append(objects, curObj)
	}
	return objects, nil
}
