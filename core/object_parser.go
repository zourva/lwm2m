package core

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
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
	optSetBool(objectMap["Multiple"], class.SetMandatory)
	optSetBool(objectMap["Mandatory"], class.SetMandatory)

	if objectMap["Resources"] != nil {
		resources := objectMap["Resources"].([]any)
		resClasses := ParseResources(resources)
		class.SetResources(resClasses)
	}

	return class
}
