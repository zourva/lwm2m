package core

import (
	"bytes"
	"github.com/zourva/pareto/endec/senml"
	"strconv"
)

// Field defines the runtime
// representation of resource.
type Field interface {
	Parent() ObjectInstance
	Class() Resource
	InstanceID() InstanceID
	Value

	AppendSENML(dst []senml.Record) []senml.Record
	MarshalJSON() ([]byte, error)
}

var _ Field = &ResourceField{}

type Fields map[InstanceID]Field

func NewFields() Fields {
	//return orderedmap.New[string, string]()
	return make(Fields)
}

func (f *Fields) Add(s Field) {
	(*f)[s.InstanceID()] = s
}

func (f *Fields) Update(s Field) {
	(*f)[s.InstanceID()] = s
}

func (f *Fields) Delete(id InstanceID) {
	delete(*f, id)
}

func (f *Fields) Field(id InstanceID) Field {
	//for _, v := range f {
	//	if v.InstanceID() == id {
	//		return v
	//	}
	//}
	v, ok := (*f)[id]
	if ok {
		return v
	}
	return nil
}

func (f *Fields) AppendSENML(dst []senml.Record) []senml.Record {
	for _, v := range *f {
		dst = v.AppendSENML(dst)
	}

	return dst
}

func (f *Fields) MarshalJSON() ([]byte, error) {
	var pack senml.Pack
	records := f.AppendSENML(nil)

	if len(records) > 0 {
		parent := (*f)[0].Parent()
		bname := GenBaseName(parent)

		records[0].BaseName = bname
	}

	pack.Records = records
	return senml.Encode(pack, senml.JSON)
}

// ResourceField implements Field.
type ResourceField struct {
	id         ResourceID //deprecated
	class      Resource
	instanceId InstanceID
	parent     ObjectInstance
	value      Value
}

func NewResourceField(id ResourceID, value Value) *ResourceField {
	return &ResourceField{
		id:    id,
		value: value,
	}
}

func NewResourceField2(parent ObjectInstance, instId InstanceID, class Resource, value Value) *ResourceField {
	return &ResourceField{
		id:         class.Id(),
		class:      class,
		instanceId: instId,
		value:      value,
		parent:     parent,
	}
}

func (v *ResourceField) AppendSENML(dst []senml.Record) []senml.Record {
	name := strconv.Itoa(int(v.Class().Id()))
	instId := v.InstanceID()
	if instId != 0 {
		// if instId == 0, don't modify name
		name = name + `/` + strconv.Itoa(int(instId))
	}

	r := fieldValueToSenmlRecord(v)
	r.Name = name

	return append(dst, *r)
}

func (v *ResourceField) MarshalJSON() ([]byte, error) {
	var pack senml.Pack
	records := v.AppendSENML(nil)
	if v.parent != nil {
		if len(records) > 0 {
			parent := v.parent
			bname := GenBaseName(parent)

			records[0].BaseName = bname
		}
	}
	pack.Records = records
	return senml.Encode(pack, senml.JSON)
}

func (v *ResourceField) Parent() ObjectInstance {
	return v.parent
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
