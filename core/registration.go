package core

import (
	"github.com/zourva/lwm2m/coap"
	"time"
)

type RegistrationClient interface {
	// Register registers a client to lwM2M Server(s).
	Register() error

	// Deregister unregisters a client from the lwM2M
	// servers to which it previously registered.
	Deregister() error

	// Update updates registration info of a client.
	Update(params ...any) error
}

type RegistrationServer interface {
	OnRegister(*RegistrationInfo) (string, error)
	OnUpdate(*RegistrationInfo) error
	OnDeregister(location string)
}

type RegisteredClient interface {
	Name() string
	Address() string
	Location() string
	Timeout() bool
	Update(info *RegistrationInfo)
	GetObjectClass(t ObjectID) Object
	RegistrationInfo() *RegistrationInfo
	DeviceControlProxy
}

// RegistrationInfo defines registered client
// info passed from protocol layer to service layer.
type RegistrationInfo struct {
	// optional endpoint name, must be globally unique if provided
	Name string `msgpack:"name"`

	// mandatory ip:port tuple or MSISDN
	Address string `msgpack:"address"`

	// mandatory lifetime in seconds, 2592000(30 days) by default
	Lifetime int `msgpack:"lifetime"`

	// mandatory protocol version for compatability
	LwM2MVersion string `msgpack:"lwM2MVersion"`

	// optional binding mode, U by default
	BindingMode BindingMode `msgpack:"bindingMode"`

	// mandatory objects and instances, excluding
	// object 0, 21, and 23
	ObjectInstances []*coap.CoreResource `msgpack:"objectInstances"`

	Location       string    `msgpack:"location"`
	RegisterTime   time.Time `msgpack:"registerTime"`
	DeregisterTime time.Time `msgpack:"deregisterTime"`
	UpdateTime     time.Time `msgpack:"updateTime"`
}

func (r *RegistrationInfo) Update(info *RegistrationInfo) {
	// TODO: update other fields
	r.Address = info.Address
	r.LwM2MVersion = info.LwM2MVersion
	r.BindingMode = info.BindingMode
	r.ObjectInstances = info.ObjectInstances

	r.UpdateTime = time.Now()
}
