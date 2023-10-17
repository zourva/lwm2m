package core

import "bytes"

// ResourceInstance contains a class
// describing the meta type, an instance
// id and a value.
type ResourceInstance interface {
	Class() Resource
	InstanceID() InstanceID
	Value
}

type ResInstManager struct {
	resources map[ResourceID][]ResourceInstance
}

func NewResInstManager() *ResInstManager {
	return &ResInstManager{
		resources: make(map[ResourceID][]ResourceInstance),
	}
}

func (m *ResInstManager) findInstance(instances []ResourceInstance, iid InstanceID) ResourceInstance {
	for _, i := range instances {
		if i.InstanceID() == iid {
			return i
		}
	}

	return nil
}

// GetSingleResource returns instance 0 of a resource or nil if not exist.
func (m *ResInstManager) GetSingleResource(id ResourceID) ResourceInstance {
	if instances, ok := m.resources[id]; ok {
		inst := m.findInstance(instances, 0)
		if inst.Class().Multiple() {
			//
		}
		return inst
	}

	return nil
}

// GetResource returns instance of a resource or nil if not exist.
func (m *ResInstManager) GetResource(id ResourceID, iid InstanceID) ResourceInstance {
	if instances, ok := m.resources[id]; ok {
		return m.findInstance(instances, iid)
	}

	return nil
}

func (m *ResInstManager) GetResources(id ResourceID) []ResourceInstance {
	return m.resources[id]
}

func (m *ResInstManager) GetAllResources() map[ResourceID][]ResourceInstance {
	return m.resources
}

type ResourceValue struct {
	id     ResourceID //deprecated
	class  Resource
	instId InstanceID
	value  Value
}

func NewResourceValue(id ResourceID, value Value) Value {
	return &ResourceValue{
		id:    id,
		value: value,
	}
}

func (v ResourceValue) InstanceID() InstanceID {
	return v.instId
}

func (v ResourceValue) Class() Resource {
	return v.class
}

// Deprecated
func (v ResourceValue) GetId() ResourceID {
	return v.id
}

func (v ResourceValue) ToBytes() []byte {
	return v.value.ToBytes()
}

func (v ResourceValue) ContainedType() ValueType {
	return ValueTypeResource
}

func (v ResourceValue) Type() ValueType {
	return ValueTypeResource
}

func (v ResourceValue) ToString() string {
	return v.value.ToString()
}

func (v ResourceValue) Get() interface{} {
	return v.value.Get()
}

func NewMultipleResourceValue(id ResourceID, value []*ResourceValue) Value {
	return &MultipleResourceValue{
		id:        id,
		instances: value,
	}
}

type MultipleResourceValue struct {
	id        ResourceID
	instances []*ResourceValue
}

func (v MultipleResourceValue) ToBytes() []byte {
	return []byte{}
}

func (v MultipleResourceValue) ContainedType() ValueType {
	return ValueTypeResource
}

func (v MultipleResourceValue) Type() ValueType {
	return ValueTypeMultiResource
}

func (v MultipleResourceValue) ToString() string {
	var buf bytes.Buffer

	for _, res := range v.instances {
		buf.WriteString(res.ToString())
		buf.WriteString(",")
	}
	return buf.String()
}

func (v MultipleResourceValue) Get() any {
	return v.instances
}
