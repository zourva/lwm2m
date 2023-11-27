package core

import (
	"bytes"
	"encoding/json"
	"strconv"
)

// Field defines the runtime
// representation of resource.
type Field interface {
	Class() Resource
	InstanceID() InstanceID
	Value

	MarshalJSON() ([]byte, error)
}

var _ Field = &ResourceField{}

// ResourceField implements Field.
type ResourceField struct {
	id         ResourceID //deprecated
	class      Resource
	instanceId InstanceID
	value      Value
}

func NewResourceField(id ResourceID, value Value) *ResourceField {
	return &ResourceField{
		id:    id,
		value: value,
	}
}

func NewResourceField2(instId InstanceID, class Resource, value Value) *ResourceField {
	return &ResourceField{
		id:         class.Id(),
		class:      class,
		instanceId: instId,
		value:      value,
	}
}

func (v *ResourceField) MarshalJSON() ([]byte, error) {
	buf := []byte(`{`)
	buf = append(buf, `"id":`+strconv.Itoa(int(v.id))...)
	buf = append(buf, `,"instanceId":`+strconv.Itoa(int(v.instanceId))...)

	data, _ := json.Marshal(v.class)
	buf = append(buf, `,"class":`+string(data)...)
	data, _ = v.value.MarshalJSON()
	buf = append(buf, `,"value":`+string(data)...)
	buf = append(buf, `}`...)

	return buf, nil
}

func (v *ResourceField) InstanceID() InstanceID {
	return v.instanceId
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

func (v *MultipleResourceValue) MarshalJSON() ([]byte, error) {
	buf := v.ToString()
	return []byte(`"` + buf + `"`), nil
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
