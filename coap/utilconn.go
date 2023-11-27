package coap

import (
	"net"
	"time"
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

type MessageContext struct {
	server  Server
	msg     *Message
	conn    Connection
	addr    *net.UDPAddr
	timeout time.Duration
}

// SendMessageTo sends a CoAP Message to UDP address
func SendMessageTo(msgCtx *MessageContext) (Response, error) {
	c, msg, conn, addr, timeout := msgCtx.server, msgCtx.msg, msgCtx.conn, msgCtx.addr, msgCtx.timeout

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

	if timeout == 0 {
		// 针对每个请求如果没有设置请求超时时间，则读取server默认配置的超时时间
		timeout = c.GetRecvTimeout()
	}

	select {
	case respCh := <-ch:
		return respCh.Response, respCh.Error
	case <-time.After(timeout):
		DeleteResponseChannel(c, msg.Id)
		return nil, ErrRecvTimeout
	}
}

func MessageSizeAllowed(req Request) bool {
	msg := req.Message()
	b, _ := MessageToBytes(msg)

	if len(b) > 65536 {
		return false
	}

	return true
}
