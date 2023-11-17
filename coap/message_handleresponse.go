package coap

import (
	"log"
	"net"
)

func handleResponse(s Server, msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	if msg.GetOption(OptionObserve) != nil {
		handleAcknowledgeObserveRequest(s, msg)
		return
	}

	ch := GetResponseChannel(s, msg.Id)
	if ch != nil {
		resp := &CoapResponseChannel{
			Response: NewResponse(msg, nil),
		}
		ch <- resp
		DeleteResponseChannel(s, msg.Id)
	} else {
		log.Println("Channel is nil", msg.Id)
	}
}

func handleAcknowledgeObserveRequest(s Server, msg *Message) {
	s.GetEvents().Notify(msg.GetURIPath(), msg.Payload, msg)
}
