package core

import (
	log "github.com/sirupsen/logrus"
	"github.com/zourva/lwm2m/coap"
)

// Messager hides details using coap binding.
// All the LwM2M operations using CoAP layer
// MUST be Confirmable CoAP messages.
type Messager interface {
	NewRequest(t coap.MessageType, m coap.Code, mt coap.MediaType, uri string) coap.Request
	NewConRequestPlainText(method coap.Code, uri string) coap.Request
	NewConRequestOpaque(method coap.Code, uri string, payload []byte) coap.Request
	NewGetRequest(uri string) coap.Request
	NewPiggybackedResponse(coap.Request, coap.Code, coap.Payload) coap.Response
	Send(req coap.Request) (coap.Response, error)
	Notify(key string, value []byte) error
}

type BaseMessager struct {
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
func (m *BaseMessager) NewPiggybackedResponse(req coap.Request, code coap.Code, payload coap.Payload) coap.Response {
	msg := coap.NewMessageOfType(coap.MessageAcknowledgment, req.Message().Id)
	msg.Token = req.Message().Token
	msg.Code = code

	if payload != nil {
		msg.Payload = payload
	}

	log.Debugln("new piggybacked response:", msg)

	return coap.NewResponseWithMessage(msg)
}

func (m *BaseMessager) NewGetRequest(uri string) coap.Request {
	return m.NewRequest(coap.MessageConfirmable, coap.Get, coap.MediaTypeTextPlain, uri)
}

func (m *BaseMessager) NewConRequestPlainText(method coap.Code, uri string) coap.Request {
	return m.NewRequest(coap.MessageConfirmable, method, coap.MediaTypeTextPlain, uri)
}

func (m *BaseMessager) NewConRequestOpaque(method coap.Code, uri string, payload []byte) coap.Request {
	req := m.NewRequest(coap.MessageConfirmable, method, coap.MediaTypeOpaqueVndOmaLwm2m, uri)
	req.SetPayload(payload)
	return req
}

func (m *BaseMessager) NewRequest(t uint8, c coap.Code, mt coap.MediaType, uri string) coap.Request {
	req := coap.NewRequest(t, c, coap.GenerateMessageID())
	req.SetRequestUri(uri)
	req.SetMediaType(mt)
	return req
}

func (m *BaseMessager) Send(req coap.Request) (coap.Response, error) {
	return nil, nil
}

func (m *BaseMessager) Notify(observationId string, data []byte) error {
	return nil
}
