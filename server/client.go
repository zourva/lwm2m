package server

import (
	"fmt"
	"github.com/zourva/lwm2m/coap"
	. "github.com/zourva/lwm2m/core"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

// NewRegisteredClient creates a new session, using the given
// registration information, representing the registered client.
func NewRegisteredClient(server *LwM2MServer, info *RegistrationInfo, registry ObjectRegistry) RegisteredClient {
	client := &registeredClient{
		regInfo:   info,
		registry:  registry,
		server:    server,
		instances: make(map[ObjectID]map[InstanceID]RegisteredObject),
	}

	client.createObjects(info.ObjectInstances)

	return client
}

// registeredClient manages the lifecycle of a client
// on server side from register to deregister.
//
// NOTE: not goroutine-safe.
type registeredClient struct {
	server   *LwM2MServer //server context
	regInfo  *RegistrationInfo
	registry ObjectRegistry

	// TODO: no need to instantiate, just keep the CoRE-Link
	// object instance ids when reported or updated
	instances map[ObjectID]map[InstanceID]RegisteredObject

	enabled atomic.Bool
}

func (c *registeredClient) Enabled() bool {
	return c.enabled.Load()
}

func (c *registeredClient) Enable() {
	c.enabled.Store(true)
}

func (c *registeredClient) Disable() {
	c.enabled.Store(false)
}

func (c *registeredClient) RegistrationInfo() *RegistrationInfo {
	return c.regInfo
}

func (c *registeredClient) Name() string {
	return c.regInfo.Name
}

func (c *registeredClient) Address() string {
	return c.regInfo.Address
}

func (c *registeredClient) Location() string {
	return c.regInfo.Location
}

// Timeout returns true if a duration of lifetime
// elapsed since last renewal update of lifetime.
func (c *registeredClient) Timeout() bool {
	duration := c.regInfo.UpdateTime.Sub(c.regInfo.RegRenewTime)
	return duration > time.Duration(c.regInfo.Lifetime)*time.Second
}

// Update updates parameters defined in
// OMA-TS-LightweightM2M_Core-V1_2_1-20221209-A
// 6.2.2. including:
//
//	Lifetime
//	Binding Mode
//	SMS Number
//	Objects and Object Instances
//	Profile ID
func (c *registeredClient) Update(info *RegistrationInfo) {
	c.regInfo.Update(info)
}

//
//func (c *registeredClient) SetObjects(objects map[ObjectID]ObjectInstance) {
//	c.enabledObjects = objects
//}

// GetObjectClass returns Object class definition for the given id.
func (c *registeredClient) GetObjectClass(t ObjectID) Object {
	return c.registry.GetObject(t)
}

func (c *registeredClient) Create(oid ObjectID, newValue Value) error {
	return c.server.messager.Create(c.Address(), oid, newValue)
}

func (c *registeredClient) Read(oid ObjectID, oiId InstanceID, rid ResourceID, riId InstanceID) ([]byte, error) {
	return c.server.messager.Read(c.Address(), oid, oiId, rid, riId)
}

func (c *registeredClient) Write(oid ObjectID, oiId InstanceID, rid ResourceID, riId InstanceID, newValue Value) ([]byte, error) {
	return c.server.messager.Write(c.Address(), oid, oiId, rid, riId, newValue)
}

func (c *registeredClient) Delete(oid ObjectID, oiId InstanceID, rid ResourceID, riId InstanceID) error {
	return c.server.messager.Delete(c.Address(), oid, oiId, rid, riId)
}

func (c *registeredClient) Execute(oid ObjectID, oiId InstanceID, rid ResourceID, args string) error {
	return c.server.messager.Execute(c.Address(), oid, oiId, rid, args)
}

func (c *registeredClient) Discover(oid ObjectID, oiId InstanceID, rid ResourceID, depth int) ([]*coap.CoREResource, error) {
	return c.server.messager.Discover(c.Address(), oid, oiId, rid, depth)
}

func (c *registeredClient) Observe(oid ObjectID, attrs NotificationAttrs, h ObserveHandler, moreIds ...uint16) error {
	oiId, rid, riId := NoneID, NoneID, NoneID
	if len(moreIds) > 0 {
		oiId = moreIds[0]
	}

	if len(moreIds) > 1 {
		rid = moreIds[1]
	}

	if len(moreIds) > 2 {
		riId = moreIds[2]
	}

	return c.server.messager.Observe(c.Address(), oid, oiId, rid, riId, attrs, h)
}

func (c *registeredClient) CancelObservation(oid ObjectID, oiId InstanceID, rid ResourceID, riId InstanceID) error {
	return c.server.messager.CancelObservation(c.Address(), oid, oiId, rid, riId)
}

func (c *registeredClient) ObserveComposite(contentType coap.MediaType, reqBody []byte, h ObserveHandler) ([]byte, error) {
	return c.server.messager.ObserveComposite(c.Address(), contentType, reqBody, h)
}

func (c *registeredClient) CancelObservationComposite(contentType coap.MediaType, reqBody []byte) error {
	return c.server.messager.CancelObservationComposite(c.Address(), contentType, reqBody)
}

func (c *registeredClient) makeAccessPath(oid ObjectID, oiId InstanceID, rid ResourceID, riId InstanceID) string {
	optionIds := []uint16{oiId, rid, riId}

	uri := fmt.Sprintf("/%d", oid)
	for _, id := range optionIds {
		if oiId == NoneID {
			break
		}

		uri += fmt.Sprintf("/%d", id)
	}

	return uri
}

// creates registered object instances based on paths.
// from CoaP POST body
//
//	</lwm2m>;rt="oma.lwm2m", </lwm2m /1/0>,</lwm2m /1/1>,
//	</lwm2m /2/0>,</lwm2m /2/1>,</lwm2m /2/2>,</lwm2m/2/3>,
//	</lwm2m /2/4>,</lwm2m /3/0>,</lwm2m /4/0>,</lwm2m /5>
//	or
//	</>;ct=110, </1/0>,</1/1>,</2/0>,</2/1>,</2/2>,</2/3>,</2/4>,</3/0>,</4/0>,</5>
func (c *registeredClient) createObjects(objInstances []*coap.CoREResource) {
	for _, o := range objInstances {
		t := o.Target[1:len(o.Target)]

		// remove root path
		if strings.Contains(t, " ") {
			t = strings.Split(t, " ")[1]
		}

		// t has format: /1/0> or </1/0>
		sp := strings.Split(t, "/")
		objectId, _ := strconv.Atoi(sp[0])

		oid := ObjectID(objectId)

		// create instance id map if new
		m, ok := c.instances[oid]
		if !ok {
			m = make(map[InstanceID]RegisteredObject)
		}

		instanceId := 0
		if len(sp) > 1 {
			instanceId, _ = strconv.Atoi(sp[1])
		}

		class := c.registry.GetObject(oid)
		m[InstanceID(instanceId)] = NewRegisteredObject(class, InstanceID(instanceId))
	}
}

//func (c *registeredClient) ReadResource(obj ObjectID, objInst Id, res ResourceID) (Value, error) {
//	clientAddr, _ := net.ResolveUDPAddr("udp", c.Address())
//
//	uri := fmt.Sprintf("/%d/%d/%d", obj, objInst, res)
//	req := coap.NewRequest(coap.MessageConfirmable, coap.Get, coap.GenerateMessageID())
//	req.SetRequestUri(uri)
//
//	resourceDefinition := c.GetObject(obj).Class().Resource(res)
//	if resourceDefinition.Multiple() {
//		req.SetMediaType(coap.MediaTypeTlvVndOmaLwm2m)
//	} else {
//		req.SetMediaType(coap.MediaTypeTextPlainVndOmaLwm2m)
//	}
//
//	response, err := c.coapChannel.SendTo(req, clientAddr)
//	if err != nil {
//		log.Println(err)
//		return nil, err
//	}
//	responseValue, _ := utils.DecodeResourceValue(
//		res, response.Message().Payload.ToBytes(), resourceDefinition)
//
//	return responseValue, nil
//}
