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

type MessagerClient struct {
	client *LwM2MClient
	state  connState
	mute   bool

	// service layer delegator
	deviceCtrlDelegator DeviceControlClient
	bootstrapDelegator  BootstrapClient
	reporterDelegator   ReportingClient
}

func NewMessager(c *LwM2MClient) *MessagerClient {
	m := &MessagerClient{
		mute:   false,
		state:  disconnected,
		client: c,
	}

	m.deviceCtrlDelegator = c.controller
	m.bootstrapDelegator = c.bootstrapper
	m.reporterDelegator = c.reporter

	return m
}

func (m *MessagerClient) Start() {
	s := m.conn()

	// add a callback to trigger auto registration
	// procedure when transport layer started.
	s.OnStart(func(server coap.CoapServer) {
		m.state = connected
		log.Infoln("lwm2m client connected")
	})

	s.OnObserve(func(observationId string, msg *coap.Message) {
		log.Infoln("observe request received for", observationId)
		// TODO: extract attributes
		m.reporterDelegator.OnObserve(observationId, nil)
	})

	s.OnObserveCancel(func(observationId string, msg *coap.Message) {
		log.Infoln("observe request received for", observationId)
		m.reporterDelegator.OnCancelObservation(observationId)
	})

	s.OnError(func(err error) {
		log.Errorln("err received:", err)
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

func (m *MessagerClient) PauseAcceptRequests() {
	m.mute = true
}

func (m *MessagerClient) ResumeAcceptRequests() {
	m.mute = false
}

func (m *MessagerClient) muted() bool {
	return m.mute
}

func (m *MessagerClient) conn() coap.CoapServer {
	return m.client.coapConn
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

	err := m.bootstrapper().OnFinish()
	msg := m.NewAckPiggyback(req, GetErrorCode(err), coap.NewEmptyPayload())

	return coap.NewResponseWithMessage(msg)
}

////// device management and service enablement handlers

func (m *MessagerClient) onServerCreate(req coap.CoapRequest) coap.CoapResponse {
	log.Debugln("receive create request:", req.GetMessage().GetURIPath())

	objectId := m.getOID(req)
	err := m.devController().OnCreate(objectId, String(""))
	msg := m.NewAckPiggyback(req, GetErrorCode(err), coap.NewEmptyPayload())

	log.Debugln("create request done:", msg)

	return coap.NewResponseWithMessage(msg)
}

func (m *MessagerClient) onServerRead(req coap.CoapRequest) coap.CoapResponse {
	log.Debugln("receive read request:", req.GetMessage().GetURIPath())

	oid := m.getOID(req)
	oiId := m.getOIID(req)
	rid := m.getRID(req)
	riId := m.getRIID(req)

	var payload coap.MessagePayload
	value, err := m.devController().OnRead(oid, oiId, rid, riId)
	if err == ErrorNone {
		buf := endec.EncodeValue(rid, value.Class().Multiple(), value)
		payload = coap.NewBytesPayload(buf)
	}

	msg := m.NewAckPiggyback(req, GetErrorCode(err), payload)
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
	msg := m.NewAckPiggyback(req, GetErrorCode(err), coap.NewEmptyPayload())

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
	msg := m.NewAckPiggyback(req, GetErrorCode(err), coap.NewEmptyPayload())

	return coap.NewResponseWithMessage(msg)
}

func (m *MessagerClient) onServerExecute(req coap.CoapRequest) coap.CoapResponse {
	log.Debugln("receive execute request:", req.GetMessage().GetURIPath())

	oid := m.getOID(req)
	oiId := m.getOIID(req)
	rid := m.getRID(req)

	err := m.devController().OnExecute(oid, oiId, rid, "")
	msg := m.NewAckPiggyback(req, GetErrorCode(err), coap.NewEmptyPayload())

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

func (m *MessagerClient) NewConRequestOpaque(method coap.Code, uri string, payload []byte) coap.CoapRequest {
	req := m.NewRequest(coap.MessageConfirmable, method, coap.MediaTypeOpaqueVndOmaLwm2m, uri)
	req.SetPayload(payload)
	return req
}

func (m *MessagerClient) NewRequest(t uint8, c coap.Code, mt coap.MediaType, uri string) coap.CoapRequest {
	req := coap.NewRequest(t, c, coap.GenerateMessageID())
	req.SetRequestURI(uri)
	req.SetMediaType(mt)
	return req
}

func (m *MessagerClient) SendRequest(req coap.CoapRequest) (coap.CoapResponse, error) {
	rsp, err := m.conn().Send(req)
	if err != nil {
		//log.Println(err)
		return nil, err
	}

	return rsp, nil
}

func (m *MessagerClient) SendNotify(observationId string, data []byte) error {
	m.conn().NotifyChange(observationId, string(data), false)
	return nil
}

func (m *MessagerClient) Connected() bool {
	return m.state == connected
}
