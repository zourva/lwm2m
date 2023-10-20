package server

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zourva/lwm2m/coap"
	. "github.com/zourva/lwm2m/core"
	"strconv"
	"strings"
	"time"
)

// RegisteredClient manages the lifecycle of a client
// on server side from register to deregister.
//
// NOTE: not goroutine-safe.
type RegisteredClient struct {
	// registration info of a client.
	regInfo *RegistrationInfo

	// each client has its own enabled objects that are told
	// to server when the client is registering or updating.
	//objectStore ObjectInstanceStore

	registry ObjectRegistry

	// object instance ids when reported or updated
	instances map[ObjectID]map[InstanceID]RegisteredObject

	messager *Messager
}

// NewClient creates a new session for a registered client
// using the given registration information.
func NewClient(info *RegistrationInfo, registry ObjectRegistry) *RegisteredClient {
	session := &RegisteredClient{
		regInfo:   info,
		registry:  registry,
		instances: make(map[ObjectID]map[InstanceID]RegisteredObject),
	}

	// predefined object classes
	//reg := NewObjectRegistry(preset.NewOMAObjectClassInfoProvider())
	//store := NewObjectInstanceStore(reg, nil)
	//session.objectStore = store
	session.createObjects(info.ObjectInstances)

	return session
}

func (c *RegisteredClient) Name() string {
	return c.regInfo.Name
}

func (c *RegisteredClient) Address() string {
	return c.regInfo.Address
}

func (c *RegisteredClient) Location() string {
	return c.regInfo.Location
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
func (c *RegisteredClient) Update(info *RegistrationInfo) {
	c.regInfo.Update(info)
}

//
//func (c *RegisteredClient) SetObjects(objects map[ObjectID]ObjectInstance) {
//	c.enabledObjects = objects
//}

// GetObjectClass returns Object class definition for the given id.
func (c *RegisteredClient) GetObjectClass(t ObjectID) Object {
	return c.registry.GetClass(t)
}

func (c *RegisteredClient) Create(oid ObjectID, newValue Value) error {
	return nil
}

func (c *RegisteredClient) Read(oid ObjectID, instId InstanceID, resId ResourceID, resInstId InstanceID) error {
	uri := c.makeAccessPath(oid, instId, resId, resInstId)
	mt := coap.MediaTypeTextPlainVndOmaLwm2m
	if c.GetObjectClass(oid).Resource(resId).Multiple() {
		mt = coap.MediaTypeTlvVndOmaLwm2m
	}

	req := c.messager.NewRequest(coap.MessageConfirmable, coap.Get, coap.GenerateMessageID(), mt, uri)
	rsp, err := c.messager.SendRequestToClient(c.Address(), req)
	if err != nil {
		log.Errorln("read operation failed:", err)
		return err
	}

	// TODO: parse response
	//	responseValue, _ := utils.DecodeResourceValue()
	log.Infoln("read operation done, rsp:", rsp)

	return nil
}

func (c *RegisteredClient) Write(oid ObjectID, instId InstanceID, resId ResourceID, resInstId InstanceID, newValue Value) error {
	return nil
}

func (c *RegisteredClient) Delete(oid ObjectID, instId InstanceID, resId ResourceID, resInstId InstanceID) error {
	return nil
}

func (c *RegisteredClient) Execute(oid ObjectID, instId InstanceID, resId ResourceID, args string) error {
	return nil
}

func (c *RegisteredClient) Discover(oid ObjectID, instId InstanceID, resId ResourceID, depth int) error {
	return nil
}

func (c *RegisteredClient) makeAccessPath(oid ObjectID, oiId InstanceID, rid ResourceID, riId InstanceID) string {
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
func (c *RegisteredClient) createObjects(objInstances []*coap.CoreResource) {
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

		class := c.registry.GetClass(oid)
		m[InstanceID(instanceId)] = NewRegisteredObject(class, InstanceID(instanceId))
	}
}

//func (c *RegisteredClient) ReadResource(obj ObjectID, objInst InstanceID, res ResourceID) (Value, error) {
//	clientAddr, _ := net.ResolveUDPAddr("udp", c.Address())
//
//	uri := fmt.Sprintf("/%d/%d/%d", obj, objInst, res)
//	req := coap.NewRequest(coap.MessageConfirmable, coap.Get, coap.GenerateMessageID())
//	req.SetRequestURI(uri)
//
//	resourceDefinition := c.GetObject(obj).GetClass().Resource(res)
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
//		res, response.GetMessage().Payload.ToBytes(), resourceDefinition)
//
//	return responseValue, nil
//}

func (c *RegisteredClient) timeout() bool {
	// TODO: configurable session timeout
	return time.Since(c.regInfo.UpdateTime) > 30*time.Minute
}
