package coap

import (
	"github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/message/pool"
)

type Response interface {
	// Code returns response code
	Code() Code
	Body() []byte

	// Length returns body length.
	Length() int64

	// LocationPath returns option result of LocationPath.
	LocationPath() string
	SetLocationPath(s string)

	message() *Message
}

func NewResponse(msg *pool.Message) Response {
	rsp := &response{
		msg: &Message{
			Message:     msg,
			RouteParams: nil,
		},
		body: nil,
	}

	if msg != nil {
		rsp.body, _ = msg.ReadBody()
	}

	return rsp
}

type response struct {
	msg  *Message
	body []byte
}

func (r *response) Code() Code {
	return Code(r.msg.Code())
}

func (r *response) message() *Message {
	return r.msg
}

func (r *response) Body() []byte {
	return r.body
}

func (r *response) Length() int64 {
	size, _ := r.msg.BodySize()
	return size
}

func (r *response) LocationPath() string {
	p, _ := r.msg.Options().LocationPath()
	return p
}

func (r *response) SetLocationPath(s string) {
	r.msg.SetOptionString(message.LocationPath, s)

	//buf := make([]byte, 1024)
	//_, _, err := r.msg.Options().SetLocationPath(buf, s)
	//if err != nil {
	//	return
	//}
}
