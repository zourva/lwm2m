package core

import "github.com/zourva/lwm2m/coap"

// Messager hides details using coap binding.
// All the LwM2M operations using CoAP layer
// MUST be Confirmable CoAP messages, except as follows:
type Messager interface {
	NewRequest(t uint8, m coap.Code, mt coap.MediaType, uri string) coap.CoapRequest
	NewConRequestPlainText(method coap.Code, uri string) coap.CoapRequest
	NewConRequestOpaque(method coap.Code, uri string, payload []byte) coap.CoapRequest
	NewAckPiggyback(coap.CoapRequest, coap.Code, coap.MessagePayload) *coap.Message
	SendRequest(req coap.CoapRequest) (coap.CoapResponse, error)
	SendNotify(observationId string, data []byte) error
	//SetCallback()
}
