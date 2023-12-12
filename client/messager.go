package client

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zourva/lwm2m/coap"
	. "github.com/zourva/lwm2m/core"
	"github.com/zourva/lwm2m/utils"
	"strconv"
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
	coap.Client
	lwM2MClient *LwM2MClient

	state connState
	mute  bool

	// service layer delegator
	deviceCtrlDelegator DeviceControlClient
	bootstrapDelegator  BootstrapClient
	reporterDelegator   ReportingClient
}

func NewMessager(c *LwM2MClient) *MessagerClient {
	m := &MessagerClient{
		//Client:      coap.NewClient(c.options.serverAddress[0], coap.WithDTLSConfig(c.options.dtlsConf)),
		mute:        false,
		state:       disconnected,
		lwM2MClient: c,
	}

	m.deviceCtrlDelegator = c.controller
	m.bootstrapDelegator = c.bootstrapper
	m.reporterDelegator = c.reporter

	return m
}

func (m *MessagerClient) Close() error {
	if m.Client != nil {
		log.Debugf("close established connection...")

		err := m.Client.Close()
		if err != nil {
			log.Errorf("close established connection failed, err:%v", err)
			return err
		}
		log.Debugf("close established connection successfully")
	}

	return nil
}

func (m *MessagerClient) Dial(addr string, opts ...coap.PeerOption) error {
	log.Debugf("dial new connection(to:%s)...", addr)
	cli, err := coap.Dial(addr, opts...)
	if err != nil {
		return err
	}
	m.Client = cli
	return nil
}

func (m *MessagerClient) Redial(addr string, opts ...coap.PeerOption) error {
	if m.Client != nil {
		if err := m.Close(); err != nil {
			return err
		}
	}

	return m.Dial(addr, opts...)
}

func (m *MessagerClient) Start() {
	//s := m.conn()
	//// add a callback to trigger auto registration
	//// procedure when transport layer started.
	//s.OnStart(func(server coap.Server) {
	//	m.state = connected
	//	log.Infoln("lwm2m client connected")
	//})
	//
	//s.OnObserve(func(observationId string, msg *coap.Message) {
	//	log.Infoln("observe request received for", observationId)
	//	// TODO: extract attributes
	//	m.reporterDelegator.OnObserve(observationId, nil)
	//})
	//
	//s.OnObserveCancel(func(observationId string, msg *coap.Message) {
	//	log.Infoln("observe request received for", observationId)
	//	m.reporterDelegator.OnCancelObservation(observationId)
	//})
	//
	//s.OnError(func(err error) {
	//	log.Errorln("err received:", err)
	//})

	// for device control interface methods
	m.Get("/{oid:[0-9]+}/{oiid:[0-9]+}/{rid:[0-9]+}/{riid:[0-9]+}", m.onServerRead)
	m.Get("/{oid:[0-9]+}/{oiid:[0-9]+}/{rid:[0-9]+}", m.onServerRead)
	m.Get("/{oid:[0-9]+}/{oiid:[0-9]+}", m.onServerRead)
	m.Get("/{oid:[0-9]+}", m.onServerRead)

	m.Put("/{oid:[0-9]+}/{oiid:[0-9]+}/{rid:[0-9]+}/{riid:[0-9]+}", m.onServerWrite)
	m.Put("/{oid:[0-9]+}/{oiid:[0-9]+}/{rid:[0-9]+}", m.onServerWrite)
	m.Put("/{oid:[0-9]+}/{oiid:[0-9]+}", m.onServerWrite)

	m.Delete("/{oid:[0-9]+}/{oiid:[0-9]+}/{rid:[0-9]+}/{riid:[0-9]+}", m.onServerDelete)
	m.Delete("/{oid:[0-9]+}/{oiid:[0-9]+}", m.onServerDelete)

	m.Post("/{oid:[0-9]+}/{oiid:[0-9]+}/{rid:[0-9]+}", m.onServerExecute)
	m.Post("/{oid:[0-9]+}", m.onServerCreate)

	m.Post("/bs", m.onBootstrapFinish)
}

// PauseUserPlane stops accepting requests from servers.
func (m *MessagerClient) PauseUserPlane() {
	m.mute = true
}

// ResumeUserPlane resumes accepting requests from servers.
func (m *MessagerClient) ResumeUserPlane() {
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

func (m *MessagerClient) getOID(req coap.Request) ObjectID {
	objectId := req.Attribute("oid")
	oid, _ := strconv.Atoi(objectId)
	return ObjectID(oid)
}

// if not provided, return NoneID
func (m *MessagerClient) getOIID(req coap.Request) InstanceID {
	instanceId := NoneID

	instance := req.Attribute("oiid")
	if instance != "" {
		oiId, _ := strconv.Atoi(instance)
		instanceId = InstanceID(oiId)
	}

	return instanceId
}

// if not provided, return NoneID
func (m *MessagerClient) getRID(req coap.Request) ResourceID {
	resourceId := NoneID

	resource := req.Attribute("rid")
	if resource != "" {
		rid, _ := strconv.Atoi(resource)
		resourceId = ResourceID(rid)
	}

	return resourceId
}

// if not provided, return NoneID
func (m *MessagerClient) getRIId(req coap.Request) InstanceID {
	instanceId := NoneID

	instance := req.Attribute("riid")
	if instance != "" {
		riId, _ := strconv.Atoi(instance)
		instanceId = InstanceID(riId)
	}

	return instanceId
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

	return m.NewAckResponse(req, GetErrorCode(err))
}

////// device management and service enablement handlers

func (m *MessagerClient) onServerCreate(req coap.Request) coap.Response {
	log.Debugln("receive create request:", req.Path())

	objectId := m.getOID(req)
	err := m.devController().OnCreate(objectId, String(""))

	return m.NewAckResponse(req, GetErrorCode(err))
}

func (m *MessagerClient) onServerRead(req coap.Request) coap.Response {
	log.Debugln("receive read request:", req.Path())

	oid := m.getOID(req)
	oiId := m.getOIID(req)
	rid := m.getRID(req)
	riId := m.getRIId(req)

	value, err := m.devController().OnRead(oid, oiId, rid, riId)
	rsp := m.NewAckPiggybackedResponse(req, GetErrorCode(err), value)

	log.Debugf("on read response:%s", value)
	return rsp
}

func (m *MessagerClient) onServerDelete(req coap.Request) coap.Response {
	log.Debugln("receive delete request:", req.Path())

	oid := m.getOID(req)
	oiId := m.getOIID(req)
	rid := m.getRID(req)
	riId := m.getRIId(req)

	err := m.devController().OnDelete(oid, oiId, rid, riId)

	return m.NewAckResponse(req, GetErrorCode(err))
}

func (m *MessagerClient) onServerDiscover(req coap.Request) {
	log.Debugln("receive discover request:", req.Path())
}

func (m *MessagerClient) onServerWrite(req coap.Request) coap.Response {
	log.Debugln("receive write request:", req.Path())

	oid := m.getOID(req)
	oiId := m.getOIID(req)
	rid := m.getRID(req)
	riId := m.getRIId(req)

	value := req.Body()
	err := m.devController().OnWrite(oid, oiId, rid, riId, value)

	return m.NewAckResponse(req, GetErrorCode(err))
}

func (m *MessagerClient) onServerExecute(req coap.Request) coap.Response {
	log.Debugln("receive execute request:", req.Path())

	oid := m.getOID(req)
	oiId := m.getOIID(req)
	rid := m.getRID(req)

	err := m.devController().OnExecute(oid, oiId, rid, "")

	return m.NewAckResponse(req, GetErrorCode(err))
}

func (m *MessagerClient) onServerObserve() {
	log.Println("Observe Request")
}

func (m *MessagerClient) Connected() bool {
	return m.state == connected
}

func (m *MessagerClient) Register(info *regInfo) error {
	// send request
	req := m.NewPostRequestCoReLink(RegisterUri, []byte(info.objects))
	req.AddQuery("ep", info.name)
	req.AddQuery("lt", utils.IntToStr(info.lifetime))
	req.AddQuery("lwm2m", lwM2MVersion)
	req.AddQuery("b", info.mode)

	log.Infof("send register(%s) request...", info.name)
	rsp, err := m.Send(req)
	if err != nil {
		log.Errorf("send register(%s) request failed:%v", info.name, err)
		return err
	}

	// check response code
	if rsp.Code().Created() {
		// save location for update or de-register operation
		info.location = rsp.LocationPath()
		log.Infof("register(%s) done with assigned location:%s", info.name, info.location)
		return nil
	}

	log.Errorf("register(%s) request failed:%s", info.name, rsp.Code().String())

	return errors.New(rsp.Code().String())
}

func (m *MessagerClient) Update(info *regInfo, params ...string) error {
	uri := RegisterUri + fmt.Sprintf("/%s", info.location)
	req := m.NewPostRequestCoReLink(uri, nil)

	for _, param := range params {
		if param == "lt" {
			req.AddQuery("lt", utils.IntToStr(info.lifetime))
		} else if param == "objlink" {
			req.SetBody([]byte(info.objects))
		}
	}

	rsp, err := m.Send(req)
	if err != nil {
		log.Errorln("send update request failed:", err)
		return err
	}

	// check response code
	if rsp.Code().Changed() {
		log.Infoln("update done on", uri)
		return nil
	}

	log.Errorln("update request failed:", rsp.Code().String())

	return errors.New(rsp.Code().String())
}

func (m *MessagerClient) Deregister(info *regInfo) error {
	uri := RegisterUri + fmt.Sprintf("/%s", info.location)
	req := m.NewDeleteRequestPlain(uri)
	rsp, err := m.Send(req)
	if err != nil {
		log.Errorln("send de-register request failed:", err)
		return err
	}

	// check response code
	if rsp.Code().Deleted() {
		log.Infoln("deregister done on", uri)
		return nil
	}

	log.Errorln("de-register request failed:", rsp.Code().String())

	return errors.New(rsp.Code().String())
}
