package core

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

	Fields(rid ResourceID) []Field
	SetFields(rid ResourceID, fields []Field)

	AllFields() map[ResourceID][]Field
	SetAllFields(all map[ResourceID][]Field)

	// SingleField equals Field(id, 0)
	SingleField(id ResourceID) Field

	// SetSingleField overwrites Field(id, 0)
	SetSingleField(f Field)
}

type InstanceMap = map[InstanceID]ObjectInstance

type InstanceIdsMap = map[ObjectID][]InstanceID

type BaseInstance struct {
	class    Object     //object class
	instId   InstanceID //object instance id
	operator Operator   //copy from class definition

	resources map[ResourceID][]Field
}

func NewObjectInstance(class Object) ObjectInstance {
	i := &BaseInstance{
		class:     class,
		instId:    0,
		operator:  class.Operator(),
		resources: make(map[ResourceID][]Field),
	}

	return i
}

func (o *BaseInstance) findInstance(instances []Field, iid InstanceID) Field {
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

func (o *BaseInstance) Fields(id ResourceID) []Field {
	return o.resources[id]
}

func (o *BaseInstance) AllFields() map[ResourceID][]Field {
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
	o.resources[f.Class().Id()] = append(o.resources[f.Class().Id()], f)
}

func (o *BaseInstance) SetSingleField(f Field) {
	if len(o.resources[f.Class().Id()]) == 0 {
		o.AddField(f)
	} else {
		o.resources[f.Class().Id()][0] = f
	}
}

func (o *BaseInstance) SetFields(rid ResourceID, fields []Field) {
	o.resources[rid] = fields
}

func (o *BaseInstance) SetAllFields(all map[ResourceID][]Field) {
	o.resources = all
}
