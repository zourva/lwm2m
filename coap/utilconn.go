package coap

import (
	"net"
)

type CoapResponseChannel struct {
	Response Response
	Error    error
}

func doSendMessage(c Server, msg *Message, conn Connection, addr *net.UDPAddr, ch chan *CoapResponseChannel) {
	resp := &CoapResponseChannel{}

	b, err := MessageToBytes(msg)
	if err != nil {
		resp.Error = err
		ch <- resp
	}

	_, err = conn.WriteTo(b, addr)
	if err != nil {
		resp.Error = err
		ch <- resp
	}

	if msg.Type == MessageNonConfirmable {
		resp.Response = NewResponse(NewEmptyMessage(msg.Id), nil)
		ch <- resp
	}

	AddResponseChannel(c, msg.Id, ch)
}

// SendMessageTo sends a CoAP Message to UDP address
func SendMessageTo(c Server, msg *Message, conn Connection, addr *net.UDPAddr) (Response, error) {
	if conn == nil {
		return nil, ErrNilConn
	}

	if msg == nil {
		return nil, ErrNilMessage
	}

	if addr == nil {
		return nil, ErrNilAddr
	}

	ch := NewResponseChannel()
	go doSendMessage(c, msg, conn, addr, ch)
	respCh := <-ch

	return respCh.Response, respCh.Error
}

func MessageSizeAllowed(req Request) bool {
	msg := req.GetMessage()
	b, _ := MessageToBytes(msg)

	if len(b) > 65536 {
		return false
	}

	return true
}
