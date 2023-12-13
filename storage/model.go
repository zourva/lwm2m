package storage

import (
	"github.com/zourva/lwm2m/core"
)

type ObjectDescriptor struct {
	PK           int           `storm:"id,increment"` //not used
	Id           core.ObjectID `storm:"unique"`
	Name         string
	Multiple     bool
	Mandatory    bool
	Version      string
	LwM2MVersion string
	URN          string
	Resources    []struct {
		Id             int
		Name           string
		Operations     string
		Multiple       bool
		Mandatory      bool
		ResourceType   string
		RangeOrEnums   string
		ValueValidator string
	}
}

type ObjectRecord struct {
	Pk      int    `storm:"id,increment"` //not used
	Unique  uint32 `storm:"unique"`
	Content any
}

func newObjectRecord(unique uint32, content any) *ObjectRecord {
	if content == nil {
		// make default []
		content = make([]any, 0)
	}
	return &ObjectRecord{
		Unique:  unique,
		Content: content,
	}
}

type DBObject struct {
	Id           core.ObjectID `storm:"id"`
	Name         string
	Multiple     bool
	Mandatory    bool
	Version      string
	LwM2MVersion string
	URN          string
}

type DBResource struct {
	Pk           int             `storm:"id,increment"` //not used
	OId          core.ObjectID   `storm:"index"`
	Id           core.ResourceID `storm:"index"`
	Name         string
	Operations   string
	Multiple     bool
	Mandatory    bool
	ResourceType string
	RangeOrEnums string
}

type DBInstance struct {
	Pk    int             `storm:"id,increment"` //not used
	OId   core.ObjectID   `storm:"index"`
	OIId  core.InstanceID `storm:"index"`
	RId   core.ResourceID `storm:"index"`
	RIId  core.InstanceID `storm:"index"`
	Value core.Field
}

type DBObservation struct {
	Pk    int `storm:"id,increment"` //not used
	OId   core.ObjectID
	OIId  core.InstanceID
	RId   core.ResourceID
	RIId  core.InstanceID
	Attrs map[string]any
}
