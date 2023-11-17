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
	*BaseMessager
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
	s.OnStart(func(server coap.Server) {
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

func (m *MessagerClient) conn() coap.Server {
	return m.client.coapConn
}

func (m *MessagerClient) bootstrapper() BootstrapClient {
	return m.bootstrapDelegator
}

func (m *MessagerClient) devController() DeviceControlClient {
	return m.deviceCtrlDelegator
}

func (m *MessagerClient) getOID(req coap.Request) ObjectID {
	objectId := req.GetAttributeAsInt("oid")
	return ObjectID(objectId)
}

// if not provided, return NoneID
func (m *MessagerClient) getOIID(req coap.Request) InstanceID {
	instanceId := NoneID

	instance := req.GetAttribute("oiid")
	if instance != "" {
		instanceId = InstanceID(req.GetAttributeAsInt("oiid"))
	}

	return instanceId
}

// if not provided, return NoneID
func (m *MessagerClient) getRID(req coap.Request) ResourceID {
	resourceId := NoneID

	resource := req.GetAttribute("rid")
	if resource != "" {
		resourceId = ResourceID(req.GetAttributeAsInt("rid"))
	}

	return resourceId
}

// if not provided, return NoneID
func (m *MessagerClient) getRIId(req coap.Request) InstanceID {
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

func (m *MessagerClient) onBootstrapRead(req coap.Request) coap.Response {
	panic("implement me")
}

func (m *MessagerClient) onBootstrapWrite(req coap.Request) coap.Response {
	panic("implement me")
}

func (m *MessagerClient) onBootstrapDelete(req coap.Request) coap.Response {
	panic("implement me")
}

func (m *MessagerClient) onBootstrapDiscover(req coap.Request) coap.Response {
	panic("implement me")
}

func (m *MessagerClient) onBootstrapFinish(req coap.Request) coap.Response {
	log.Debugln("receive bootstrap finish")

	err := m.bootstrapper().OnFinish()

	return m.NewPiggybackedResponse(req, GetErrorCode(err), coap.NewEmptyPayload())
}

////// device management and service enablement handlers

func (m *MessagerClient) onServerCreate(req coap.Request) coap.Response {
	log.Debugln("receive create request:", req.GetMessage().GetURIPath())

	objectId := m.getOID(req)
	err := m.devController().OnCreate(objectId, String(""))

	return m.NewPiggybackedResponse(req, GetErrorCode(err), coap.NewEmptyPayload())
}

func (m *MessagerClient) onServerRead(req coap.Request) coap.Response {
	log.Debugln("receive read request:", req.GetMessage().GetURIPath())

	oid := m.getOID(req)
	oiId := m.getOIID(req)
	rid := m.getRID(req)
	riId := m.getRIId(req)

	var payload coap.Payload
	value, err := m.devController().OnRead(oid, oiId, rid, riId)
	if err == ErrorNone {
		buf := endec.EncodeValue(rid, value.Class().Multiple(), value)
		payload = coap.NewBytesPayload(buf)
	}

	rsp := m.NewPiggybackedResponse(req, GetErrorCode(err), payload)
	rsp.GetMessage().AddOption(coap.OptionContentFormat, m.getMediaTypeFromValue(value))

	return rsp
}

func (m *MessagerClient) onServerDelete(req coap.Request) coap.Response {
	log.Debugln("receive delete request:", req.GetMessage().GetURIPath())

	oid := m.getOID(req)
	oiId := m.getOIID(req)
	rid := m.getRID(req)
	riId := m.getRIId(req)

	err := m.devController().OnDelete(oid, oiId, rid, riId)

	return m.NewPiggybackedResponse(req, GetErrorCode(err), coap.NewEmptyPayload())
}

func (m *MessagerClient) onServerDiscover(req coap.Request) {
	log.Debugln("receive discover request:", req.GetMessage().GetURIPath())
}

func (m *MessagerClient) onServerWrite(req coap.Request) coap.Response {
	log.Debugln("receive write request:", req.GetMessage().GetURIPath())

	oid := m.getOID(req)
	oiId := m.getOIID(req)
	rid := m.getRID(req)
	riId := m.getRIId(req)

	err := m.devController().OnWrite(oid, oiId, rid, riId, String(""))

	return m.NewPiggybackedResponse(req, GetErrorCode(err), coap.NewEmptyPayload())
}

func (m *MessagerClient) onServerExecute(req coap.Request) coap.Response {
	log.Debugln("receive execute request:", req.GetMessage().GetURIPath())

	oid := m.getOID(req)
	oiId := m.getOIID(req)
	rid := m.getRID(req)

	err := m.devController().OnExecute(oid, oiId, rid, "")

	return m.NewPiggybackedResponse(req, GetErrorCode(err), coap.NewEmptyPayload())
}

func (m *MessagerClient) onServerObserve() {
	log.Println("Observe Request")
}

// NewPiggybackedResponse creates an ACK-piggybacked response.
//
//	Client              Server
//	   |                  |
//	   |   CON [0x7d34]   |
//	   +----------------->|
//	   |                  |
//	   |   ACK [0x7d34]   |
//	   |<-----------------+
//	   |                  |
func (m *MessagerClient) NewPiggybackedResponse(req coap.Request, code coap.Code, payload coap.Payload) coap.Response {
	msg := coap.NewMessageOfType(coap.MessageAcknowledgment, req.GetMessage().Id)
	msg.Token = req.GetMessage().Token
	msg.Code = code

	if payload != nil {
		msg.Payload = payload
	}

	log.Debugln("new piggybacked response:", msg)

	return coap.NewResponseWithMessage(msg)
}

func (m *MessagerClient) Send(req coap.Request) (coap.Response, error) {
	rsp, err := m.conn().Send(req)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

func (m *MessagerClient) Notify(observationId string, data []byte) error {
	m.conn().NotifyChange(observationId, string(data), false)
	return nil
}

func (m *MessagerClient) Connected() bool {
	return m.state == connected
}
