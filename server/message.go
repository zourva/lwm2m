package server

import (
	log "github.com/sirupsen/logrus"
	"github.com/zourva/lwm2m/coap"
	. "github.com/zourva/lwm2m/core"
	"net"
	"strconv"
	"time"
)

// Messager encapsulates and hide
// transport layer details.
type Messager struct {
	// application layer
	//bootstrapService
	registrationService RegistrationServer
	reportingService    ReportingServer
	//m.deviceControlService = NewDeviceControlService(s)

	// session layer
	coapConn coap.CoapServer
}

func NewMessageHandler(s *LwM2MServer) *Messager {
	m := &Messager{
		coapConn: s.coapConn,
	}

	//s.bootstrapService = NewBootstrapService()
	m.registrationService = NewRegistrationService(s)
	m.reportingService = NewInfoReportingService(s)

	//m.deviceControlService = NewDeviceControlService(s)

	return m
}

// handle request parameters like:
//
//	uri: /rd?ep={Endpoint Client Name}&lt={Lifetime}
//	        &lwm2m={version}&b={binding}&Q&sms={MSISDN}&pid={ProfileID}
//	   b/Q/sms/pid are optional.
//	body: </1/0>,... which is optional.
func (m *Messager) onClientRegister(req coap.CoapRequest) coap.CoapResponse {
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

	clientId, err := m.registrationService.OnRegister(info)
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
func (m *Messager) onClientUpdate(req coap.CoapRequest) coap.CoapResponse {
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

	err := m.registrationService.OnUpdate(info)
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
func (m *Messager) onClientDeregister(req coap.CoapRequest) coap.CoapResponse {
	id := req.GetAttribute("id")

	m.registrationService.OnDeregister(id)

	msg := m.createRspMsg(req, coap.MessageAcknowledgment, coap.CodeDeleted)
	return coap.NewResponseWithMessage(msg)
}

// handle request with parameters like:
//
//	uri: /dp
//	body: implementation-specific.
func (m *Messager) onSendInfo(req coap.CoapRequest) coap.CoapResponse {
	data := req.GetMessage().Payload.GetBytes()
	// check resource contained in reported list
	// check server granted read access

	// commit to application layer
	rsp, err := m.reportingService.OnSend(data)
	if err != nil {
		log.Errorf("error recv client info: %v", err)
		msg := m.createRspMsg(req, coap.MessageAcknowledgment, coap.CodeInternalServerError)
		return coap.NewResponseWithMessage(msg)
	}

	msg := m.createRspMsg(req, coap.MessageAcknowledgment, coap.CodeContent)
	msg.Payload = coap.NewBytesPayload(rsp)
	return coap.NewResponseWithMessage(msg)
}

func (m *Messager) createRspMsg(req coap.CoapRequest, mt uint8, code coap.Code) *coap.Message {
	msg := coap.NewMessageOfType(mt, req.GetMessage().MessageID)
	msg.Token = req.GetMessage().Token
	msg.Code = code

	return msg
}

func (m *Messager) NewRequest(t uint8, c coap.Code, id uint16, mt coap.MediaType, uri string) coap.CoapRequest {
	//req := coap.NewRequest(coap.MessageConfirmable, coap.Get, coap.GenerateMessageID())
	req := coap.NewRequest(t, c, id)
	req.SetRequestURI(uri)
	req.SetMediaType(mt)
	return req
}

func (m *Messager) SendRequestToClient(peer string, req coap.CoapRequest) ([]byte, error) {
	clientAddr, _ := net.ResolveUDPAddr("udp", peer)

	response, err := m.coapConn.SendTo(req, clientAddr)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return response.GetMessage().Payload.GetBytes(), nil
}
