package core

type InstanceID = uint16

const (
	NoneID uint16 = 0xFFFF
)

// ObjectInstance defines an instance of
// an Object at runtime.
//
//	1 ObjectID -> 1 Object Instance Store
//	1 Object Instance Store -> 0/1/* Object Instances mapped by id
type ObjectInstance interface {
	Class() Object

	Id() InstanceID

	SetId(id InstanceID)

	Field(id ResourceID, iid InstanceID) Field

	Fields(id ResourceID) []Field

	AllFields() map[ResourceID][]Field

	// SingleField equals Field(id, 0)
	SingleField(id ResourceID) Field
}

type InstanceMap = map[InstanceID]ObjectInstance

type InstanceIdsMap = map[ObjectID][]InstanceID

type objectInstance struct {
	class  Object     //object class
	instId InstanceID //object instance id

	resources map[ResourceID][]Field
}

func newObjectInstance(class Object, id InstanceID) *objectInstance {
	return &objectInstance{
		class:     class,
		instId:    id,
		resources: make(map[ResourceID][]Field),
	}
}

func (o *objectInstance) findInstance(instances []Field, iid InstanceID) Field {
	for _, i := range instances {
		if i.InstanceID() == iid {
			return i
		}
	}

	return nil
}

// GetSingleResource returns instance 0 of a resource or nil if not exist.
func (o *objectInstance) SingleField(id ResourceID) Field {
	if instances, ok := o.resources[id]; ok {
		inst := o.findInstance(instances, 0)
		if inst.Class().Multiple() {
			//
		}
		return inst
	}

	return nil
}

// GetResource returns instance of a resource or nil if not exist.
func (o *objectInstance) Field(id ResourceID, iid InstanceID) Field {
	if instances, ok := o.resources[id]; ok {
		return o.findInstance(instances, iid)
	}

	return nil
}

func (o *objectInstance) Fields(id ResourceID) []Field {
	return o.resources[id]
}

func (o *objectInstance) AllFields() map[ResourceID][]Field {
	return o.resources
}

func (o *objectInstance) Class() Object {
	return o.class
}

func (o *objectInstance) Id() InstanceID {
	return o.instId
}

func (o *objectInstance) SetId(id InstanceID) {
	o.instId = id
}
