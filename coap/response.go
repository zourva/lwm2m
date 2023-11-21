package coap

import "strings"

func NoResponse() Response {
	return NilResponse{}
}

type Response interface {
	Message() *Message
	Error() error
	Payload() []byte
	UriQuery(q string) string
}

type NilResponse struct {
}

func (c NilResponse) Message() *Message {
	return nil
}

func (c NilResponse) Error() error {
	return nil
}

func (c NilResponse) Payload() []byte {
	return nil
}

func (c NilResponse) UriQuery(q string) string {
	return ""
}

// NewResponse creates a new Response object with a Message object and any error messages
func NewResponse(msg *Message, err error) Response {
	resp := &DefaultResponse{
		msg: msg,
		err: err,
	}

	return resp
}

// NewResponseWithMessage creates a new response object with a Message object
func NewResponseWithMessage(msg *Message) Response {
	resp := &DefaultResponse{
		msg: msg,
	}

	return resp
}

type DefaultResponse struct {
	msg *Message
	err error
}

func (c *DefaultResponse) Message() *Message {
	return c.msg
}

func (c *DefaultResponse) Error() error {
	return c.err
}

func (c *DefaultResponse) Payload() []byte {
	return c.Message().Payload.GetBytes()
}

func (c *DefaultResponse) UriQuery(q string) string {
	qs := c.Message().GetOptionsAsString(OptionURIQuery)

	for _, o := range qs {
		ps := strings.Split(o, "=")
		if len(ps) == 2 {
			if ps[0] == q {
				return ps[1]
			}
		}
	}
	return ""
}
