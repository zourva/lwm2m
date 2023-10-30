package storage

import "github.com/zourva/lwm2m/core"

type Object struct {
	Id           core.ObjectID `storm:"id"`
	Name         string
	Multiple     bool
	Mandatory    bool
	Version      string
	LwM2MVersion string
	URN          string
}

type Resource struct {
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

type Instance struct {
	Pk    int             `storm:"id,increment"` //not used
	OId   core.ObjectID   `storm:"index"`
	OIId  core.InstanceID `storm:"index"`
	RId   core.ResourceID `storm:"index"`
	RIId  core.InstanceID `storm:"index"`
	Value any             `storm:"value"`
}

type Observation struct {
	Pk    int `storm:"id,increment"` //not used
	OId   core.ObjectID
	OIId  core.InstanceID
	RId   core.ResourceID
	RIId  core.InstanceID
	Attrs map[string]any
}
