package core

import "bytes"

// Field defines the runtime
// representation of resource.
type Field interface {
	Class() Resource
	InstanceID() InstanceID
	Value
}

// ResourceField implements Field.
type ResourceField struct {
	id     ResourceID //deprecated
	class  Resource
	instId InstanceID
	value  Value
}

func NewResourceField(id ResourceID, value Value) *ResourceField {
	return &ResourceField{
		id:    id,
		value: value,
	}
}

func (v *ResourceField) InstanceID() InstanceID {
	return v.instId
}

func (v *ResourceField) Class() Resource {
	return v.class
}

func (v *ResourceField) ToBytes() []byte {
	return v.value.ToBytes()
}

func (v *ResourceField) ContainedType() ValueType {
	return ValueTypeResource
}

func (v *ResourceField) Type() ValueType {
	return ValueTypeResource
}

func (v *ResourceField) ToString() string {
	return v.value.ToString()
}

func (v *ResourceField) Get() interface{} {
	return v.value.Get()
}

func NewMultipleResourceValue(id ResourceID, value []*ResourceField) Value {
	return &MultipleResourceValue{
		id:        id,
		instances: value,
	}
}

type MultipleResourceValue struct {
	id        ResourceID
	instances []*ResourceField
}

func (v *MultipleResourceValue) ToBytes() []byte {
	return []byte{}
}

func (v *MultipleResourceValue) ContainedType() ValueType {
	return ValueTypeResource
}

func (v *MultipleResourceValue) Type() ValueType {
	return ValueTypeMultiResource
}

func (v *MultipleResourceValue) ToString() string {
	var buf bytes.Buffer

	for _, res := range v.instances {
		buf.WriteString(res.ToString())
		buf.WriteString(",")
	}
	return buf.String()
}

func (v *MultipleResourceValue) Get() any {
	return v.instances
}
