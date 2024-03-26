package core

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zourva/pareto/endec/senml"
	"strconv"
	"strings"
)

var typeMap = map[string]ValueType{
	"multiple":      ValueTypeMultiple,
	"string":        ValueTypeString,
	"byte":          ValueTypeByte,
	"int":           ValueTypeInteger,
	"int32":         ValueTypeInteger32,
	"int64":         ValueTypeInteger64,
	"float":         ValueTypeFloat,
	"float64":       ValueTypeFloat64,
	"bool":          ValueTypeBoolean,
	"opaque":        ValueTypeOpaque,
	"time":          ValueTypeTime,
	"objectlink":    ValueTypeObjectLink,
	"object":        ValueTypeObject,
	"resource":      ValueTypeResource,
	"multiresource": ValueTypeMultiResource,
}

var opsMap = map[string]OpCode{
	"N":  OpNone,
	"R":  OpRead,
	"W":  OpWrite,
	"RW": OpReadWrite,
	"E":  OpExecute,
}

func optSetString(val any, setter func(string)) {
	if val == nil {
		return
	}
	setter(val.(string))
}

func optSetBool(val any, setter func(bool)) {
	if val == nil {
		return
	}
	realVal, _ := val.(bool)
	setter(realVal)
}

func ParseResources(resources []any) []Resource {
	var arrays []Resource
	for _, res := range resources {
		r := res.(map[string]any)

		rc := &ResourceImpl{}
		rc.SetId(ResourceID(r["Id"].(float64)))
		optSetString(r["Name"], rc.SetName)
		optSetBool(r["Multiple"], rc.SetMultiple)
		optSetBool(r["Mandatory"], rc.SetMandatory)
		optSetString(r["RangeOrEnums"], rc.SetRangeOrEnums)
		rc.SetType(typeMap[r["ResourceType"].(string)])
		rc.SetOperations(opsMap[r["Operations"].(string)])

		arrays = append(arrays, rc)
	}

	return arrays
}

func ParseObject(objJSON string) Object {
	var objectMap map[string]any
	if err := json.Unmarshal([]byte(objJSON), &objectMap); err != nil {
		log.Errorln("unmarshal object class json failed", err)
		return nil
	}

	class := &ObjectClass{}
	class.SetId(ObjectID(objectMap["Id"].(float64)))
	optSetString(objectMap["Name"], class.SetName)
	optSetString(objectMap["Version"], class.SetVersion)
	optSetString(objectMap["LwM2MVersion"], class.SetLwM2MVersion)
	optSetString(objectMap["URN"], class.SetUrn)
	optSetString(objectMap["Description"], class.SetDescription)
	optSetBool(objectMap["Multiple"], class.SetMultiple)
	optSetBool(objectMap["Mandatory"], class.SetMandatory)

	if objectMap["Resources"] != nil {
		resources := objectMap["Resources"].([]any)
		resClasses := ParseResources(resources)
		class.SetResources(resClasses)
	}

	class.SetOperator(NewBaseOperator())

	return class
}

func ParsePathToNumbers(src string, sep string) ([]uint16, error) {
	var ids []uint16

	if len(src) == 0 {
		return nil, fmt.Errorf("empty source string")
	}

	keys := strings.Split(src, sep)
	if keys[0] == "" && len(keys) > 1 {
		// 跳过 第一个空的分割
		keys = keys[1:]
	}

	for _, k := range keys {
		if id, err := strconv.Atoi(k); err != nil {
			log.Errorf("path to ids failed:%v, err:%v", src, err)
			return nil, err
		} else {
			ids = append(ids, uint16(id))
		}
	}

	return ids, nil
}

func SenmlRecordSetFieldValue(r *senml.Record, src Field) {
	kind := src.Class().Type()
	switch kind {
	case ValueTypeEmpty:
	case ValueTypeMultiple: //return &MultipleValue{}
	case ValueTypeString:
		tmp := src.ToString()
		r.StringValue = &tmp
	case ValueTypeByte:
		tmp := src.ToString()
		r.StringValue = &tmp
	case ValueTypeInteger:
		tmp := float64(src.Get().(int))
		r.Value = &tmp
	case ValueTypeInteger32:
		tmp := float64(src.Get().(int32))
		r.Value = &tmp
	case ValueTypeInteger64:
		tmp := float64(src.Get().(int64))
		r.Value = &tmp
	case ValueTypeFloat:
		tmp := float64(src.Get().(float32))
		r.Value = &tmp
	case ValueTypeFloat64:
		tmp := float64(src.Get().(float64))
		r.Value = &tmp
	case ValueTypeBoolean:
		tmp := src.Get().(bool)
		r.BoolValue = &tmp
	case ValueTypeOpaque:
		tmp := string(src.Get().([]byte))
		r.OpaqueValue = &tmp
	case ValueTypeTime:
	case ValueTypeObjectLink:
	case ValueTypeObject:
	case ValueTypeResource:
	case ValueTypeMultiResource:
	}
	//return r
}

func FieldValueToSenmlRecord(src Field) *senml.Record {
	r := &senml.Record{}
	SenmlRecordSetFieldValue(r, src)
	return r
}

func SenmlRecordToFieldValue(kind ValueType, val *senml.Record) Value {
	switch kind {
	case ValueTypeEmpty:
		return Empty()
	case ValueTypeMultiple: //return &MultipleValue{}
	case ValueTypeString:
		return String(*val.StringValue)
	case ValueTypeByte:
		return ByteVal([]byte(*val.StringValue)...)
	case ValueTypeInteger:
		return Integer(int(*val.Value))
	case ValueTypeInteger32:
		return Integer(int(*val.Value))
	case ValueTypeInteger64:
		return Integer(int(*val.Value))
	case ValueTypeFloat:
		return Float(float32(*val.Value))
	case ValueTypeFloat64:
		return Float64(*val.Value)
	case ValueTypeBoolean:
		return Boolean(*val.BoolValue)
	case ValueTypeOpaque:
		return Opaque([]byte(*val.OpaqueValue))
	case ValueTypeTime:
		return nil
	case ValueTypeObjectLink:
		return nil
	case ValueTypeObject:
		return nil
	case ValueTypeResource:
		return nil
	case ValueTypeMultiResource:
		return nil
	}
	return nil
}
