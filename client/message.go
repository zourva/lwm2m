package client

import (
	log "github.com/sirupsen/logrus"
	"github.com/zourva/lwm2m/coap"
	. "github.com/zourva/lwm2m/core"
	"github.com/zourva/lwm2m/endec"
)

type connState = int

const (
	disconnected connState = iota
	connecting
	connected
	reconnectPending
	reconnected
)

// Messager hides details using coap binding.
// All the LwM2M operations using CoAP layer
// MUST be Confirmable CoAP messages, except as follows:
type Messager interface {
	NewRequest(t uint8, m coap.Code, mt coap.MediaType, uri string) coap.CoapRequest
	NewConRequestPlainText(method coap.Code, uri string) coap.CoapRequest
	NewAckPiggyback(coap.CoapRequest, coap.Code, coap.MessagePayload) *coap.Message
	SendRequest(req coap.CoapRequest) (coap.CoapResponse, error)
	//SetCallback()
}

var errorCodesMapping = map[ErrorType]coap.Code{
	ErrorNone:                coap.CodeEmpty,
	BadRequest:               coap.CodeBadRequest,
	Unauthorized:             coap.CodeUnauthorized,
	BadOption:                coap.CodeBadOption,
	Forbidden:                coap.CodeForbidden,
	Conflict:                 coap.CodeConflict,
	NotFound:                 coap.CodeNotFound,
	MethodNotAllowed:         coap.CodeMethodNotAllowed,
	NotAcceptable:            coap.CodeNotAcceptable,
	RequestEntityIncomplete:  coap.CodeRequestEntityIncomplete,
	PreconditionFailed:       coap.CodePreconditionFailed,
	RequestEntityTooLarge:    coap.CodeRequestEntityTooLarge,
	UnsupportedContentFormat: coap.CodeUnsupportedContentFormat,
}

type MessagerClient struct {
	coapConn coap.CoapServer
	state    connState
	mute     bool

	// service layer delegator
	deviceCtrlDelegator DeviceControlClient
	bootstrapDelegator  BootstrapClient
}

func NewMessager(c *LwM2MClient) *MessagerClient {
	m := &MessagerClient{
		mute:     false,
		state:    disconnected,
		coapConn: coap.NewServer(c.name, c.options.localAddress, c.options.serverAddress[0]),
	}

	m.deviceCtrlDelegator = c.controller
	m.bootstrapDelegator = c.bootstrapper

	return m
}

func (m *MessagerClient) Start() {
	s := m.coapConn

	// add a callback to trigger auto registration
	// procedure when transport layer started.
	s.OnStart(func(server coap.CoapServer) {
		m.state = connected
		log.Infoln("lwm2m client connected")
	})

	s.OnObserve(func(resource string, msg *coap.Message) {
		log.Infoln("observe request received for", resource)
	})

	// for device control interface methods
	s.Get("/:oid/:oiid/:rid/:riid", m.onServerRead)
	s.Get("/:oid/:oiid/:rid", m.onServerRead)
	s.Get("/:oid/:oiid", m.onServerRead)
	s.Get("/:oid", m.onServerRead)

	s.Put("/:oid/:oiid/:rid/:riid", m.onServerWrite)
	s.Put("/:oid/:oiid/:rid", m.onServerWrite)
	s.Put("/:oid/:oiid", m.onServerWrite)

	s.Delete("/:oid/:oiid/:rid/:riid", m.onServerDelete)
	s.Delete("/:oid/:oiid", m.onServerDelete)

	s.Post("/:oid/:oiid/:rid", m.onServerExecute)
	s.Post("/:oid", m.onServerCreate)

	s.Post("/bs", m.onBootstrapFinish)

	// this method does not hold
	s.Start()
}

func (m *MessagerClient) Pause() {
	m.mute = true
}

func (m *MessagerClient) Resume() {
	m.mute = false
}

func (m *MessagerClient) muted() bool {
	return m.mute
}

func (m *MessagerClient) bootstrapper() BootstrapClient {
	return m.bootstrapDelegator
}

func (m *MessagerClient) devController() DeviceControlClient {
	return m.deviceCtrlDelegator
}

func (m *MessagerClient) getOID(req coap.CoapRequest) ObjectID {
	objectId := req.GetAttributeAsInt("oid")
	return ObjectID(objectId)
}

// if not provided, return NoneID
func (m *MessagerClient) getOIID(req coap.CoapRequest) InstanceID {
	instanceId := NoneID

	instance := req.GetAttribute("oiid")
	if instance != "" {
		instanceId = InstanceID(req.GetAttributeAsInt("oiid"))
	}

	return instanceId
}

// if not provided, return NoneID
func (m *MessagerClient) getRID(req coap.CoapRequest) ResourceID {
	resourceId := NoneID

	resource := req.GetAttribute("rid")
	if resource != "" {
		resourceId = ResourceID(req.GetAttributeAsInt("rid"))
	}

	return resourceId
}

// if not provided, return NoneID
func (m *MessagerClient) getRIID(req coap.CoapRequest) InstanceID {
	instanceId := NoneID

	instance := req.GetAttribute("riid")
	if instance != "" {
		instanceId = InstanceID(req.GetAttributeAsInt("riid"))
	}

	return instanceId
}

func (m *MessagerClient) getMediaTypeFromValue(v Value) coap.MediaType {
	if v.Type() == ValueTypeMultiple {
		return coap.MediaTypeTlvVndOmaLwm2m
	} else {
		return coap.MediaTypeTextPlain
	}
}

func (m *MessagerClient) getErrorCode(err ErrorType) coap.Code {
	return errorCodesMapping[err]
}

////// bootstrap procedure handlers

func (m *MessagerClient) onBootstrapRead(req coap.CoapRequest) coap.CoapResponse {
	panic("implement me")
}

func (m *MessagerClient) onBootstrapWrite(req coap.CoapRequest) coap.CoapResponse {
	panic("implement me")
}

func (m *MessagerClient) onBootstrapDelete(req coap.CoapRequest) coap.CoapResponse {
	panic("implement me")
}

func (m *MessagerClient) onBootstrapDiscover(req coap.CoapRequest) coap.CoapResponse {
	panic("implement me")
}

func (m *MessagerClient) onBootstrapFinish(req coap.CoapRequest) coap.CoapResponse {
	log.Debugln("receive bootstrap finish")

	err := m.bootstrapper().OnBootstrapFinish()
	code := m.getErrorCode(err)
	msg := m.NewAckPiggyback(req, code, coap.NewEmptyPayload())

	return coap.NewResponseWithMessage(msg)
}

////// device management and service enablement handlers

func (m *MessagerClient) onServerCreate(req coap.CoapRequest) coap.CoapResponse {
	log.Debugln("receive create request:", req.GetMessage().GetURIPath())

	objectId := m.getOID(req)
	err := m.devController().OnCreate(objectId, String(""))
	code := m.getErrorCode(err)

	msg := m.NewAckPiggyback(req, code, coap.NewEmptyPayload())

	log.Debugln("create request done:", msg)

	return coap.NewResponseWithMessage(msg)
}

func (m *MessagerClient) onServerRead(req coap.CoapRequest) coap.CoapResponse {
	log.Debugln("receive read request:", req.GetMessage().GetURIPath())

	oid := m.getOID(req)
	oiId := m.getOIID(req)
	rid := m.getRID(req)
	riId := m.getRIID(req)

	var code coap.Code
	var payload coap.MessagePayload

	value, err := m.devController().OnRead(oid, oiId, rid, riId)
	if err == ErrorNone {
		buf := endec.EncodeValue(rid, value.Class().Multiple(), value)
		payload = coap.NewBytesPayload(buf)
	}

	msg := m.NewAckPiggyback(req, code, payload)
	msg.AddOption(coap.OptionContentFormat, m.getMediaTypeFromValue(value))

	return coap.NewResponseWithMessage(msg)
}

func (m *MessagerClient) onServerDelete(req coap.CoapRequest) coap.CoapResponse {
	log.Debugln("receive delete request:", req.GetMessage().GetURIPath())

	oid := m.getOID(req)
	oiId := m.getOIID(req)
	rid := m.getRID(req)
	riId := m.getRIID(req)

	err := m.devController().OnDelete(oid, oiId, rid, riId)
	code := m.getErrorCode(err)

	msg := m.NewAckPiggyback(req, code, coap.NewEmptyPayload())

	return coap.NewResponseWithMessage(msg)
}

func (m *MessagerClient) onServerDiscover(req coap.CoapRequest) {
	log.Debugln("receive discover request:", req.GetMessage().GetURIPath())
}

func (m *MessagerClient) onServerWrite(req coap.CoapRequest) coap.CoapResponse {
	log.Debugln("receive write request:", req.GetMessage().GetURIPath())

	oid := m.getOID(req)
	oiId := m.getOIID(req)
	rid := m.getRID(req)
	riId := m.getRIID(req)

	err := m.devController().OnWrite(oid, oiId, rid, riId, String(""))
	code := m.getErrorCode(err)
	msg := m.NewAckPiggyback(req, code, coap.NewEmptyPayload())

	return coap.NewResponseWithMessage(msg)
}

func (m *MessagerClient) onServerExecute(req coap.CoapRequest) coap.CoapResponse {
	log.Debugln("receive execute request:", req.GetMessage().GetURIPath())

	oid := m.getOID(req)
	oiId := m.getOIID(req)
	rid := m.getRID(req)

	err := m.devController().OnExecute(oid, oiId, rid, "")
	code := m.getErrorCode(err)
	msg := m.NewAckPiggyback(req, code, coap.NewEmptyPayload())

	return coap.NewResponseWithMessage(msg)
}

func (m *MessagerClient) onServerObserve() {
	log.Println("Observe Request")
}

func (m *MessagerClient) NewAckPiggyback(req coap.CoapRequest, code coap.Code, payload coap.MessagePayload) *coap.Message {
	msg := coap.NewMessageOfType(coap.MessageAcknowledgment, req.GetMessage().MessageID)
	msg.Token = req.GetMessage().Token
	msg.Code = code

	if payload != nil {
		msg.Payload = payload
	}

	return msg
}

func (m *MessagerClient) NewConRequestPlainText(method coap.Code, uri string) coap.CoapRequest {
	return m.NewRequest(coap.MessageConfirmable, method, coap.MediaTypeTextPlain, uri)
}

func (m *MessagerClient) NewRequest(t uint8, c coap.Code, mt coap.MediaType, uri string) coap.CoapRequest {
	req := coap.NewRequest(t, c, coap.GenerateMessageID())
	req.SetRequestURI(uri)
	req.SetMediaType(mt)
	return req
}

func (m *MessagerClient) SendRequest(req coap.CoapRequest) (coap.CoapResponse, error) {
	rsp, err := m.coapConn.Send(req)
	if err != nil {
		//log.Println(err)
		return nil, err
	}

	return rsp, nil
}

func (m *MessagerClient) Connected() bool {
	return m.state == connected
}
