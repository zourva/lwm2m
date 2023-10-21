package core

// ObjectID represents a LwM2M Object Type
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

// Object describes a LwM2M Object depicted in
// OMA-TS-LightweightM2M_Core-V1_2_1-20221209-A
// Appendix D.1 Object Template.
type Object interface {
	Name() string
	Id() ObjectID
	Version() string
	LwM2MVersion() string
	URN() string
	Multiple() bool
	Mandatory() bool
	Description() string

	// Resources returns all resource fields
	// defined for this Object.
	Resources() []Resource

	// Resource returns the resource field
	// identified by the given id.
	Resource(n ResourceID) Resource

	// Operator returns operators can be
	// applied against those resource fields.
	Operator() Operator
}

// ObjectClass implements Object.
type ObjectClass struct {
	id           ObjectID
	name         string
	version      string //the 1st version must be "1.0"
	multiple     bool
	mandatory    bool
	lwM2MVersion string //since v1.1, e.g. version "1.1"
	urn          string //since v1.1, e.g. urn:oma:lwm2m:oma:1
	description  string

	resources []Resource

	// delayed initialization
	operator Operator
}

func (o *ObjectClass) SetId(id ObjectID) {
	o.id = id
}

func (o *ObjectClass) SetName(name string) {
	o.name = name
}

func (o *ObjectClass) SetVersion(version string) {
	o.version = version
}

func (o *ObjectClass) SetMultiple(multiple bool) {
	o.multiple = multiple
}

func (o *ObjectClass) SetMandatory(mandatory bool) {
	o.mandatory = mandatory
}

func (o *ObjectClass) SetLwM2MVersion(lwM2MVersion string) {
	o.lwM2MVersion = lwM2MVersion
}

func (o *ObjectClass) SetUrn(urn string) {
	o.urn = urn
}

func (o *ObjectClass) SetDescription(description string) {
	o.description = description
}

func (o *ObjectClass) SetResources(r []Resource) {
	o.resources = r
}

func (o *ObjectClass) Name() string {
	return o.name
}

func (o *ObjectClass) Id() ObjectID {
	return o.id
}

func (o *ObjectClass) Version() string {
	return o.version
}

func (o *ObjectClass) LwM2MVersion() string {
	return o.lwM2MVersion
}

func (o *ObjectClass) URN() string {
	return o.urn
}

func (o *ObjectClass) Multiple() bool {
	return o.multiple
}

func (o *ObjectClass) Mandatory() bool {
	return o.mandatory
}

func (o *ObjectClass) Description() string {
	return o.description
}

func (o *ObjectClass) Resources() []Resource {
	return o.resources
}

func (o *ObjectClass) Resource(n ResourceID) Resource {
	for _, res := range o.resources {
		if res.Id() == n {
			return res
		}
	}
	return nil
}

func (o *ObjectClass) Operator() Operator {
	return o.operator
}
