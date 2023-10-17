package preset

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/zourva/lwm2m/core"
)

type ObjectClassMap = map[core.ObjectID]core.ObjectClass

var oma ObjectClassMap
var groups map[core.ProviderType]ObjectClassMap

var typeMap = map[string]core.ValueType{
	"multiple":      core.ValueTypeMultiple,
	"string":        core.ValueTypeString,
	"byte":          core.ValueTypeByte,
	"int":           core.ValueTypeInteger,
	"int32":         core.ValueTypeInteger32,
	"int64":         core.ValueTypeInteger64,
	"float":         core.ValueTypeFloat,
	"float64":       core.ValueTypeFloat64,
	"bool":          core.ValueTypeBoolean,
	"opaque":        core.ValueTypeOpaque,
	"time":          core.ValueTypeTime,
	"objectlink":    core.ValueTypeObjectLink,
	"object":        core.ValueTypeObject,
	"resource":      core.ValueTypeResource,
	"multiresource": core.ValueTypeMultiResource,
}

var opsMap = map[string]core.OpCode{
	"N":  core.OpNone,
	"R":  core.OpRead,
	"W":  core.OpWrite,
	"RW": core.OpReadWrite,
	"E":  core.OpExecute,
}

func init() {
	oma = make(ObjectClassMap)
	groups = map[core.ProviderType]ObjectClassMap{
		core.OMAObjects: oma,
	}

	buildPreset()
}

func buildPreset() {
	parseObjectClass(securityDescriptor, oma)
	parseObjectClass(serverDescriptor, oma)
	parseObjectClass(accessControlDescriptor, oma)
	parseObjectClass(deviceDescriptor, oma)
	parseObjectClass(connMonitorDescriptor, oma)
	parseObjectClass(firmwareUpdateDescriptor, oma)
	parseObjectClass(locationDescriptor, oma)
	parseObjectClass(connStatsDescriptor, oma)
}

func GetAllPresetClasses(group core.ProviderType) ObjectClassMap {
	return groups[group]
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

func parseResourceClass(resources []any) []core.Resource {
	var arrays []core.Resource
	for _, res := range resources {
		r := res.(map[string]any)

		rc := &core.ResourceImpl{}
		rc.SetId(core.ResourceID(r["Id"].(float64)))
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

func parseObjectClass(objJSON string, groupMap map[core.ObjectID]core.ObjectClass) {
	var objectMap map[string]any
	if err := json.Unmarshal([]byte(objJSON), &objectMap); err != nil {
		log.Errorln("unmarshal object class json failed", err)
		return
	}

	class := &core.ObjectClassImpl{}
	class.SetId(core.ObjectID(objectMap["Id"].(float64)))
	optSetString(objectMap["Name"], class.SetName)
	optSetString(objectMap["Version"], class.SetVersion)
	optSetString(objectMap["LwM2MVersion"], class.SetLwM2MVersion)
	optSetString(objectMap["URN"], class.SetUrn)
	optSetString(objectMap["Description"], class.SetDescription)
	optSetBool(objectMap["Multiple"], class.SetMandatory)
	optSetBool(objectMap["Mandatory"], class.SetMandatory)

	if objectMap["Resources"] != nil {
		resources := objectMap["Resources"].([]any)
		resClasses := parseResourceClass(resources)
		class.SetResources(resClasses)
	}

	groupMap[class.Id()] = class
}
