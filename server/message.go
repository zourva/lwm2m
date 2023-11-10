package server

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zourva/lwm2m/coap"
	. "github.com/zourva/lwm2m/core"
	"net"
	"strconv"
	"time"
)

// ServerMessager encapsulates and hide
// transport layer details.
//
// Put them together here to make it
// easier when replacing COAP layer later.
type ServerMessager struct {
	server *LwM2MServer

	// application layer
	bootstrapService BootstrapServer
	registerService  RegistrationServer
	reportingService ReportingServer
	//m.deviceControlService = NewDeviceControlService(s)

	// session layer
	coapConn coap.CoapServer
}

func NewMessager(server *LwM2MServer) *ServerMessager {
	m := &ServerMessager{
		server:   server,
		coapConn: server.coapConn,
	}

	m.bootstrapService = NewBootstrapService(server)
	m.registerService = NewRegistrationService(server)
	m.reportingService = NewInfoReportingService(server)

	//m.deviceControlService = NewDeviceControlService(server)

	return m
}

// handle request parameters like:
//
//	 uri:
//		/bs?ep={Endpoint Client Name}&pct={Preferred Content Format}
func (m *ServerMessager) onClientBootstrap(req coap.CoapRequest) coap.CoapResponse {
	ep := req.GetURIQuery("ep")
	addr := req.GetAddress().String()
	err := m.bootstrapService.OnRequest(ep, addr)
	if err != nil {
		log.Errorf("error bootstrap client %s: %v", ep, err)
		msg := m.createRspMsg(req, coap.MessageAcknowledgment, GetErrorCode(err))
		return coap.NewResponseWithMessage(msg)
	}

	return nil
}

// handle request parameters like:
//
//	 uri:
//		/bspack?ep={Endpoint Client Name}
func (m *ServerMessager) onClientBootstrapPack(req coap.CoapRequest) coap.CoapResponse {
	panic("implement me")
}

// handle request parameters like:
//
//	uri: /rd?ep={Endpoint Client Name}&lt={Lifetime}
//	        &lwm2m={version}&b={binding}&Q&sms={MSISDN}&pid={ProfileID}
//	   b/Q/sms/pid are optional.
//	body: </1/0>,... which is optional.
func (m *ServerMessager) onClientRegister(req coap.CoapRequest) coap.CoapResponse {
	ep := req.GetURIQuery("ep")
	lt, _ := strconv.Atoi(req.GetURIQuery("lt"))
	lwm2m := req.GetURIQuery("lwm2m")
	binding := req.GetURIQuery("b")

	list := coap.CoreResourcesFromString(req.GetMessage().Payload.String())
	info := &RegistrationInfo{
		Name:            ep,
		Address:         req.GetAddress().String(),
		Lifetime:        lt,
		LwM2MVersion:    lwm2m,
		BindingMode:     binding,
		ObjectInstances: list,
		Location:        "",
		RegisterTime:    time.Now(),
		UpdateTime:      time.Now(),
	}

	clientId, err := m.registerService.OnRegister(info)
	if err != nil {
		log.Errorf("error registering client %s: %v", ep, err)
		msg := m.createRspMsg(req, coap.MessageAcknowledgment, coap.CodeInternalServerError)
		return coap.NewResponseWithMessage(msg)
	}

	//s.options.lcHandler.OnClientRegistered()

	msg := m.createRspMsg(req, coap.MessageAcknowledgment, coap.CodeCreated)
	msg.AddOption(coap.OptionLocationPath, "rd/"+clientId)
	return coap.NewResponseWithMessage(msg)
}

// handle request with parameters like:
//
//	uri: /{location}?lt={Lifetime}&b={binding}&Q&sms={MSISDN}
//		where location has a format of /rd/{id} and b/Q/sms are optional.
//	body: </1/0>,... which is optional.
func (m *ServerMessager) onClientUpdate(req coap.CoapRequest) coap.CoapResponse {
	id := req.GetAttribute("id")
	lt, _ := strconv.Atoi(req.GetURIQuery("lt"))
	binding := req.GetURIQuery("b")

	list := coap.CoreResourcesFromString(req.GetMessage().Payload.String())
	info := &RegistrationInfo{
		Location:        id,
		Lifetime:        lt,
		BindingMode:     binding,
		ObjectInstances: list,
		UpdateTime:      time.Now(),
	}

	err := m.registerService.OnUpdate(info)
	if err != nil {
		log.Errorf("error updating client %s: %v", info.Name, err)
		msg := m.createRspMsg(req, coap.MessageAcknowledgment, coap.CodeInternalServerError)
		return coap.NewResponseWithMessage(msg)
	}

	msg := m.createRspMsg(req, coap.MessageAcknowledgment, coap.CodeChanged)
	return coap.NewResponseWithMessage(msg)
}

// handle request with parameters like:
//
//	uri: /{location}
//	 where location has a format of /rd/{id}
func (m *ServerMessager) onClientDeregister(req coap.CoapRequest) coap.CoapResponse {
	id := req.GetAttribute("id")

	m.registerService.OnDeregister(id)

	msg := m.createRspMsg(req, coap.MessageAcknowledgment, coap.CodeDeleted)
	return coap.NewResponseWithMessage(msg)
}

// handle request with parameters like:
//
//	uri: /dp
//	body: implementation-specific.
func (m *ServerMessager) onSendInfo(req coap.CoapRequest) coap.CoapResponse {
	data := req.GetMessage().Payload.GetBytes()
	// check resource contained in reported list
	// check server granted read access

	// get registered client bound to this info
	c := m.server.GetClientManager().GetByAddr(req.GetAddress().String())
	if c == nil {
		log.Errorf("not registered or address changed, " +
			"a new registration is needed and the info sent is ignored")
		msg := m.createRspMsg(req, coap.MessageAcknowledgment, coap.CodeUnauthorized)
		return coap.NewResponseWithMessage(msg)
	}

	// commit to application layer
	rsp, err := m.reportingService.OnSend(c, data)
	if err != nil {
		log.Errorf("error recv client info: %v", err)
		msg := m.createRspMsg(req, coap.MessageAcknowledgment, coap.CodeInternalServerError)
		return coap.NewResponseWithMessage(msg)
	}

	msg := m.createRspMsg(req, coap.MessageAcknowledgment, coap.CodeChanged)
	msg.Payload = coap.NewBytesPayload(rsp)
	return coap.NewResponseWithMessage(msg)
}

func (m *ServerMessager) Read(peer string, oid ObjectID, oiId InstanceID, rid ResourceID, riId InstanceID) ([]byte, error) {
	uri := m.makeAccessPath(oid, oiId, rid, riId)
	req := m.NewConRequestPlainText(coap.Get, uri)
	rsp, err := m.SendRequest(peer, req)
	if err != nil {
		log.Errorln("read operation failed:", err)
		return nil, err
	}

	// check response code
	if rsp.GetMessage().Code == coap.CodeContent {
		log.Debugf("read operation against %s done", uri)

		return rsp.GetMessage().Payload.GetBytes(), nil
	}

	return nil, GetCodeError(rsp.GetMessage().Code)
}

func (m *ServerMessager) Discover(peer string, oid ObjectID) ([]byte, error) {
	panic("implement me")
}

func (m *ServerMessager) Write(peer string, oid ObjectID, oiId InstanceID, rid ResourceID, value Value) error {
	uri := m.makeAccessPath(oid, oiId, rid, NoneID)
	req := m.NewConRequestPlainText(coap.Put, uri)
	req.SetPayload(value.ToBytes())
	rsp, err := m.SendRequest(peer, req)
	if err != nil {
		log.Errorln("write operation failed:", err)
		return err
	}

	// check response code
	if rsp.GetMessage().Code == coap.CodeChanged {
		log.Debugf("write operation against %s done", uri)
		return nil
	}

	return GetCodeError(rsp.GetMessage().Code)
}

func (m *ServerMessager) Delete(peer string, oid ObjectID, oiId InstanceID) error {
	uri := m.makeAccessPath(oid, oiId, NoneID, NoneID)
	req := m.NewConRequestPlainText(coap.Put, uri)
	rsp, err := m.SendRequest(peer, req)
	if err != nil {
		log.Errorln("delete operation failed:", err)
		return err
	}

	// check response code
	if rsp.GetMessage().Code == coap.CodeDeleted {
		log.Debugf("delete operation against %s done", uri)
		return nil
	}

	return GetCodeError(rsp.GetMessage().Code)
}

func (m *ServerMessager) Finish(peer string) error {
	req := m.NewConRequestPlainText(coap.Get, BootstrapFinishUri)
	rsp, err := m.SendRequest(peer, req)
	if err != nil {
		log.Errorln("bootstrap finish operation failed:", err)
		return err
	}

	// check response code
	if rsp.GetMessage().Code == coap.CodeChanged {
		log.Debugf("bootstrap finish operation done")
		return nil
	}

	return GetCodeError(rsp.GetMessage().Code)
}

func (m *ServerMessager) makeAccessPath(oid ObjectID, oiId InstanceID, rid ResourceID, riId InstanceID) string {
	optionIds := []uint16{oiId, rid, riId}

	uri := fmt.Sprintf("/%d", oid)
	for _, id := range optionIds {
		if id == NoneID {
			break
		}

		uri += fmt.Sprintf("/%d", id)
	}

	return uri
}

func (m *ServerMessager) createRspMsg(req coap.CoapRequest, mt uint8, code coap.Code) *coap.Message {
	msg := coap.NewMessageOfType(mt, req.GetMessage().MessageID)
	msg.Token = req.GetMessage().Token
	msg.Code = code

	return msg
}

func (m *ServerMessager) NewAckPiggyback(req coap.CoapRequest, code coap.Code, payload coap.MessagePayload) *coap.Message {
	msg := coap.NewMessageOfType(coap.MessageAcknowledgment, req.GetMessage().MessageID)
	msg.Token = req.GetMessage().Token
	msg.Code = code

	if payload != nil {
		msg.Payload = payload
	}

	return msg
}

func (m *ServerMessager) NewConRequestPlainText(method coap.Code, uri string) coap.CoapRequest {
	return m.NewRequest(coap.MessageConfirmable, method, coap.MediaTypeTextPlain, uri)
}

func (m *ServerMessager) NewConRequestOpaque(method coap.Code, uri string, payload []byte) coap.CoapRequest {
	req := m.NewRequest(coap.MessageConfirmable, method, coap.MediaTypeOpaqueVndOmaLwm2m, uri)
	req.SetPayload(payload)
	return req
}

func (m *ServerMessager) NewRequest(t uint8, c coap.Code, mt coap.MediaType, uri string) coap.CoapRequest {
	//req := coap.NewRequest(coap.MessageConfirmable, coap.Get, coap.GenerateMessageID())
	req := coap.NewRequest(t, c, coap.GenerateMessageID())
	req.SetRequestURI(uri)
	req.SetMediaType(mt)
	return req
}

func (m *ServerMessager) SendRequest(peer string, req coap.CoapRequest) (coap.CoapResponse, error) {
	clientAddr, _ := net.ResolveUDPAddr("udp", peer)

	rsp, err := m.coapConn.SendTo(req, clientAddr)
	if err != nil {
		//log.Println(err)
		return nil, err
	}

	return rsp, nil
}

func (m *ServerMessager) SendRequestToClient(peer string, req coap.CoapRequest) ([]byte, error) {
	clientAddr, _ := net.ResolveUDPAddr("udp", peer)

	response, err := m.coapConn.SendTo(req, clientAddr)
	if err != nil {
		log.Errorf("send to peer %s failed: %v", peer, err)
		return nil, err
	}

	return response.GetMessage().Payload.GetBytes(), nil
}
