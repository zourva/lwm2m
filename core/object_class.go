package core

// ObjectID represents an LwM2M Object Type
type ObjectID = uint16

const (
	OmaObjectSecurity       ObjectID = 0
	OmaObjectServer         ObjectID = 1
	OmaObjectAccessControl  ObjectID = 2
	OmaObjectDevice         ObjectID = 3
	OmaObjectConnMonitor    ObjectID = 4
	OmaObjectFirmwareUpdate ObjectID = 5
	OmaObjectLocation       ObjectID = 6
	OmaObjectConnStats      ObjectID = 7
	OmaObjectOSCORE         ObjectID = 21
	OmaObjectCOSE           ObjectID = 23
	OmaObjectMQTTServer     ObjectID = 24
	OmaObjectGateway        ObjectID = 25
	OmaObjectGatewayRouting ObjectID = 26
	OmaObject5GNRConn       ObjectID = 27
)

// ObjectClass describes an LwM2M Object depicted in
// OMA-TS-LightweightM2M_Core-V1_2_1-20221209-A
// Appendix D.1 Object Template.
type ObjectClass interface {
	Name() string
	Id() ObjectID
	Version() string
	LwM2MVersion() string
	URN() string
	Multiple() bool
	Mandatory() bool
	Description() string

	// Resources returns all resource classes
	// defined for this ObjectClass.
	Resources() []Resource

	// Resource returns the resource class
	// identified by the given id.
	Resource(n ResourceID) Resource
}

type ObjectClassImpl struct {
	id           ObjectID
	name         string
	version      string //the 1st version must be "1.0"
	multiple     bool
	mandatory    bool
	lwM2MVersion string //since v1.1, e.g. version "1.1"
	urn          string //since v1.1, e.g. urn:oma:lwm2m:oma:1
	description  string

	resources []Resource
}

func (o *ObjectClassImpl) SetId(id ObjectID) {
	o.id = id
}

func (o *ObjectClassImpl) SetName(name string) {
	o.name = name
}

func (o *ObjectClassImpl) SetVersion(version string) {
	o.version = version
}

func (o *ObjectClassImpl) SetMultiple(multiple bool) {
	o.multiple = multiple
}

func (o *ObjectClassImpl) SetMandatory(mandatory bool) {
	o.mandatory = mandatory
}

func (o *ObjectClassImpl) SetLwM2MVersion(lwM2MVersion string) {
	o.lwM2MVersion = lwM2MVersion
}

func (o *ObjectClassImpl) SetUrn(urn string) {
	o.urn = urn
}

func (o *ObjectClassImpl) SetDescription(description string) {
	o.description = description
}

func (o *ObjectClassImpl) SetResources(r []Resource) {
	o.resources = r
}

func (o *ObjectClassImpl) Name() string {
	return o.name
}

func (o *ObjectClassImpl) Id() ObjectID {
	return o.id
}

func (o *ObjectClassImpl) Version() string {
	return o.version
}

func (o *ObjectClassImpl) LwM2MVersion() string {
	return o.lwM2MVersion
}

func (o *ObjectClassImpl) URN() string {
	return o.urn
}

func (o *ObjectClassImpl) Multiple() bool {
	return o.multiple
}

func (o *ObjectClassImpl) Mandatory() bool {
	return o.mandatory
}

func (o *ObjectClassImpl) Description() string {
	return o.description
}

func (o *ObjectClassImpl) Resources() []Resource {
	return o.resources
}

func (o *ObjectClassImpl) Resource(n ResourceID) Resource {
	for _, res := range o.resources {
		if res.Id() == n {
			return res
		}
	}
	return nil
}
