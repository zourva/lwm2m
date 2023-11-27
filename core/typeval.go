package core

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type ValueType = byte

const (
	ValueTypeEmpty         ValueType = 0
	ValueTypeMultiple      ValueType = 1
	ValueTypeString        ValueType = 2
	ValueTypeByte          ValueType = 3
	ValueTypeInteger       ValueType = 4
	ValueTypeInteger32     ValueType = 5
	ValueTypeInteger64     ValueType = 6
	ValueTypeFloat         ValueType = 7
	ValueTypeFloat64       ValueType = 8
	ValueTypeBoolean       ValueType = 9
	ValueTypeOpaque        ValueType = 10
	ValueTypeTime          ValueType = 11
	ValueTypeObjectLink    ValueType = 12
	ValueTypeObject        ValueType = 13
	ValueTypeResource      ValueType = 14
	ValueTypeMultiResource ValueType = 15
)

// Value defines a generic way
// (like reflection) to encode/decode
// resource data types and values.
type Value interface {
	// Type returns type of the underlying value.
	Type() ValueType

	// ContainedType returns the value type
	// of each element contained in this value.
	// It's the same as Type() when the underlying
	// value is not ValueTypeMultiple nor ValueTypeMultiResource.
	ContainedType() ValueType

	// Get returns the underlying value.
	Get() any

	// ToBytes returns the underlying bytes
	// array of this value.
	ToBytes() []byte

	// ToString returns a human-readable
	// string representation of this value.
	ToString() string

	MarshalJSON() ([]byte, error)
}

var _ = []Value{
	&MultipleValue{},
	&StringValue{},
	&IntegerValue{},
	&TimeValue{},
	&FloatValue{},
	&Float64Value{},
	&BooleanValue{},
	&EmptyValue{},
	&OpaqueValue{},
	&ByteValue{},
}

type MultipleValue struct {
	values        []Value
	containedType ValueType
}

func (v *MultipleValue) MarshalJSON() ([]byte, error) {
	buf := v.ToString()
	return []byte(`"` + buf + `"`), nil
}

func (v *MultipleValue) ContainedType() ValueType {
	return v.containedType
}

func (v *MultipleValue) ToBytes() []byte {
	return []byte("")
}

func (v *MultipleValue) Type() ValueType {
	return ValueTypeMultiple
}

func (v *MultipleValue) Get() any {
	return v.values
}

func (v *MultipleValue) ToString() string {
	return ""
}

type StringValue struct {
	value string
}

func (v *StringValue) MarshalJSON() ([]byte, error) {
	buf := v.ToString()
	return []byte(`"` + buf + `"`), nil
}

func (v *StringValue) ToBytes() []byte {
	return []byte(v.value)
}

func (v *StringValue) Type() ValueType {
	return ValueTypeString
}

func (v *StringValue) ContainedType() ValueType {
	return ValueTypeString
}

func (v *StringValue) Get() any {
	return v.value
}

func (v *StringValue) ToString() string {
	return v.value
}

type IntegerValue struct {
	value int
}

func (v *IntegerValue) MarshalJSON() ([]byte, error) {
	buf := v.ToString()
	return []byte(`"` + buf + `"`), nil
}

func (v *IntegerValue) ToBytes() []byte {
	sz, _ := GetValueByteLength(v.value)
	buf := new(bytes.Buffer)
	if sz == 1 {
		_ = binary.Write(buf, binary.LittleEndian, int8(v.value))
	} else if sz == 2 {
		_ = binary.Write(buf, binary.LittleEndian, int16(v.value))
	} else if sz == 4 {
		_ = binary.Write(buf, binary.LittleEndian, int32(v.value))
	} else if sz == 8 {
		_ = binary.Write(buf, binary.LittleEndian, int64(v.value))
	}
	return buf.Bytes()
}

func (v *IntegerValue) Type() ValueType {
	return ValueTypeInteger
}

func (v *IntegerValue) ContainedType() ValueType {
	return ValueTypeInteger
}

func (v *IntegerValue) Get() any {
	return v.value
}

func (v *IntegerValue) ToString() string {
	return strconv.Itoa(v.value)
}

type TimeValue struct {
	value time.Time
}

func (v *TimeValue) MarshalJSON() ([]byte, error) {
	buf := v.ToString()
	return []byte(`"` + buf + `"`), nil
}

func (v *TimeValue) ToBytes() []byte {
	return []byte(strconv.FormatInt(v.value.Unix(), 10))
}

func (v *TimeValue) Type() ValueType {
	return ValueTypeTime
}

func (v *TimeValue) ContainedType() ValueType {
	return ValueTypeTime
}

func (v *TimeValue) Get() any {
	return v.value
}

func (v *TimeValue) ToString() string {
	return strconv.FormatInt(v.value.Unix(), 10)
}

type FloatValue struct {
	value float32
}

func (v *FloatValue) MarshalJSON() ([]byte, error) {
	buf := v.ToString()
	return []byte(`"` + buf + `"`), nil
}

func (v *FloatValue) ToBytes() []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, v.value)

	return buf.Bytes()
}

func (v *FloatValue) Type() ValueType {
	return ValueTypeFloat
}

func (v *FloatValue) ContainedType() ValueType {
	return ValueTypeFloat
}

func (v *FloatValue) Get() any {
	return v.value
}

func (v *FloatValue) ToString() string {
	return strconv.FormatFloat(float64(v.value), 'g', 1, 32)
}

type Float64Value struct {
	value float64
}

func (v *Float64Value) MarshalJSON() ([]byte, error) {
	buf := v.ToString()
	return []byte(`"` + buf + `"`), nil
}

func (v *Float64Value) ToBytes() []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, v.value)

	return buf.Bytes()
}

func (v *Float64Value) Type() ValueType {
	return ValueTypeFloat64
}

func (v *Float64Value) ContainedType() ValueType {
	return ValueTypeFloat64
}

func (v *Float64Value) Get() any {
	return v.value
}

func (v *Float64Value) ToString() string {
	return strconv.FormatFloat(v.value, 'g', 1, 64)
}

type BooleanValue struct {
	value bool
}

func (v *BooleanValue) MarshalJSON() ([]byte, error) {
	buf := v.ToString()
	return []byte(`"` + buf + `"`), nil
}

func (v *BooleanValue) ToBytes() []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, v.value)

	return buf.Bytes()
}

func (v *BooleanValue) Type() ValueType {
	return ValueTypeBoolean
}

func (v *BooleanValue) ContainedType() ValueType {
	return ValueTypeBoolean
}

func (v *BooleanValue) Get() any {
	return v.value
}

func (v *BooleanValue) ToString() string {
	if v.value {
		return "1"
	} else {
		return "0"
	}
}

func Empty() Value {
	return &EmptyValue{}
}

type EmptyValue struct {
}

func (v *EmptyValue) MarshalJSON() ([]byte, error) {
	buf := v.ToString()
	return []byte(`"` + buf + `"`), nil
}

func (v *EmptyValue) ToBytes() []byte {
	return []byte("")
}

func (v *EmptyValue) Type() ValueType {
	return ValueTypeEmpty
}

func (v *EmptyValue) ContainedType() ValueType {
	return ValueTypeEmpty
}

func (v *EmptyValue) Get() any {
	return ""
}

func (v *EmptyValue) ToString() string {
	return ""
}

type OpaqueValue struct {
	value []byte
}

func (v *OpaqueValue) MarshalJSON() ([]byte, error) {
	buf := v.ToString()
	return []byte(`"` + buf + `"`), nil
}

func (v *OpaqueValue) Type() ValueType {
	return ValueTypeOpaque
}

func (v *OpaqueValue) ContainedType() ValueType {
	return ValueTypeOpaque
}

func (v *OpaqueValue) Get() any {
	return v.value
}

func (v *OpaqueValue) ToBytes() []byte {
	return v.value
}

func (v *OpaqueValue) ToString() string {
	return string(v.value)
}

type ByteValue struct {
	value byte
}

func (v *ByteValue) MarshalJSON() ([]byte, error) {
	buf := v.ToString()
	return []byte(`"` + buf + `"`), nil
}

func (v *ByteValue) Type() ValueType {
	return ValueTypeByte
}

func (v *ByteValue) ContainedType() ValueType {
	return ValueTypeByte
}

func (v *ByteValue) Get() any {
	return v.value
}

func (v *ByteValue) ToBytes() []byte {
	var ret []byte
	return append(ret, v.value)
}

func (v *ByteValue) ToString() string {
	return fmt.Sprintf("0x%x", v.value)
}

func String(v ...string) Value {
	if len(v) > 1 {
		var vs []Value

		for _, o := range v {
			vs = append(vs, String(o))
		}
		return Multiple(ValueTypeString, vs...)
	} else {
		return &StringValue{
			value: v[0],
		}
	}
}

func Integer(v ...int) Value {
	if len(v) > 1 {
		var vs []Value

		for _, o := range v {
			vs = append(vs, Integer(o))
		}
		return Multiple(ValueTypeInteger, vs...)
	} else {
		return &IntegerValue{
			value: v[0],
		}
	}
}

func Time(v ...time.Time) Value {
	if len(v) > 1 {
		var vs []Value

		for _, o := range v {
			vs = append(vs, Time(o))
		}
		return Multiple(ValueTypeTime, vs...)
	} else {
		return &TimeValue{
			value: v[0],
		}
	}
}

func Float(v ...float32) Value {
	if len(v) > 1 {
		var vs []Value

		for _, o := range v {
			vs = append(vs, Float(o))
		}
		return Multiple(ValueTypeFloat, vs...)
	} else {
		return &FloatValue{
			value: v[0],
		}
	}
}

func Float64(v ...float64) Value {
	if len(v) > 1 {
		var vs []Value

		for _, o := range v {
			vs = append(vs, Float64(o))
		}
		return Multiple(ValueTypeFloat64, vs...)
	} else {
		return &Float64Value{
			value: v[0],
		}
	}
}

func Boolean(v ...bool) Value {
	if len(v) > 1 {
		var vs []Value

		for _, o := range v {
			vs = append(vs, Boolean(o))
		}
		return Multiple(ValueTypeBoolean, vs...)
	} else {
		return &BooleanValue{
			value: v[0],
		}
	}
}

func Opaque(v ...[]byte) Value {
	if len(v) > 1 {
		var vs []Value

		for _, o := range v {
			vs = append(vs, Opaque(o))
		}
		return Multiple(ValueTypeOpaque, vs...)
	} else {
		return &OpaqueValue{
			value: v[0],
		}
	}
}

func ByteVal(v ...byte) Value {
	if len(v) > 1 {
		var vs []Value

		for _, o := range v {
			vs = append(vs, ByteVal(o))
		}
		return Multiple(ValueTypeByte, vs...)
	} else {
		return &ByteValue{
			value: v[0],
		}
	}
}

func Multiple(ct ValueType, v ...Value) Value {
	return &MultipleValue{
		values:        v,
		containedType: ct,
	}
}

func MultipleIntegers(v ...Value) Value {
	return &MultipleValue{
		values:        v,
		containedType: ValueTypeInteger,
	}
}

func ValueByType(t ValueType, val []byte) Value {
	var value Value

	switch t {
	case ValueTypeString:
		value = String(string(val))
		break
	}

	return value
}

func GetValueByteLength(val any) (uint32, error) {
	if _, ok := val.(int); ok {
		v := val.(int)
		if v > 127 || v < -128 {
			if v > 32767 || v < -32768 {
				if v > 2147483647 || v < -2147483648 {
					return 8, nil
				} else {
					return 4, nil
				}
			} else {
				return 2, nil
			}
		} else {
			return 1, nil
		}
	} else if _, ok := val.(bool); ok {
		return 1, nil
	} else if _, ok := val.(string); ok {
		v := val.(string)

		return uint32(len(v)), nil
	} else if _, ok := val.(float64); ok {
		v := val.(float64)

		if v > +3.4e+38 || v < -3.4e+38 {
			return 8, nil
		} else {
			return 4, nil
		}
	} else if _, ok := val.(time.Time); ok {
		return 8, nil
	} else if _, ok := val.([]byte); ok {
		v := val.([]byte)
		return uint32(len(v)), nil
	} else {
		return 0, errors.New("unknown type")
	}
}

func BytesToIntegerValue(b []byte) (conv Value) {
	intLen := len(b)

	if intLen == 1 {
		conv = Integer(int(b[0]))
	} else if intLen == 2 {
		conv = Integer(int(b[1]) | (int(b[0]) << 8))
	} else if intLen == 4 {
		conv = Integer(int(b[3]) | (int(b[2]) << 8) | (int(b[1]) << 16) | (int(b[0]) << 24))
	} else if intLen == 8 {
		conv = Integer(int(b[7]) | (int(b[6]) << 8) | (int(b[5]) << 16) | (int(b[4]) << 24) | (int(b[3]) << 32) | (int(b[2]) << 40) | (int(b[1]) << 48) | (int(b[0]) << 56))
	} else {
		// Error
	}
	return
}
