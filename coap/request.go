package coap

import (
	"bytes"
	"github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/mux"
	"net"
	"strings"
	"time"
)

type Request interface {
	Address() net.Addr
	SetAddress(addr net.Addr)

	Query(q string) string
	AddQuery(k string, v string)

	// Attribute returns attributes extracted
	// from location path with keyed pattern.
	Attribute(key string) string

	// Path returns uri of the request.
	Path() string

	Body() []byte
	SetBody([]byte)
	IsCoRELinkContent() bool

	// Timeout returns duration to elapse
	// before make the request timeout.
	Timeout() time.Duration
	SetTimeout(to time.Duration)

	// Length returns body length.
	Length() int64
	Options() Options
	ContentFormat() MediaType
	SetObserve(on bool)

	SecurityIdentity() string

	message() *Message
}

func NewRequest(msg *mux.Message) Request {
	req := &request{
		msg:     msg,
		body:    nil,
		timeout: DefaultTimeout,
	}

	if msg != nil {
		req.body, _ = msg.ReadBody()
	}

	return req
}

type request struct {
	addr    net.Addr
	msg     *Message
	body    []byte
	timeout time.Duration
}

func (r *request) message() *Message {
	return r.msg
}

func (r *request) Address() net.Addr {
	return r.addr
}

func (r *request) SetAddress(addr net.Addr) {
	r.addr = addr
}

func (r *request) Query(q string) string {
	qs, _ := r.msg.Queries()
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

func (r *request) AddQuery(k string, v string) {
	r.msg.AddQuery(k + "=" + v)
}

func (r *request) Attribute(o string) string {
	return r.msg.RouteParams.Vars[o]
}

func (r *request) Body() []byte {
	return r.body
}

func (r *request) SetBody(b []byte) {
	r.msg.SetBody(bytes.NewReader(b))
}

func (r *request) SetObserve(on bool) {
	if on {
		r.msg.SetObserve(0)
	} else {
		r.msg.SetObserve(1)
	}
}

func (r *request) Timeout() time.Duration {
	return r.timeout
}

func (r *request) SetTimeout(to time.Duration) {
	r.timeout = to
}

func (r *request) Path() string {
	path, _ := r.msg.Path()
	return path
}

func (r *request) Length() int64 {
	size, _ := r.msg.BodySize()
	return size
}

func (r *request) Options() Options {
	return r.msg.Options()
}

func (r *request) ContentFormat() MediaType {
	m, _ := r.Options().ContentFormat()
	return m
}

func (r *request) IsCoRELinkContent() bool {
	return r.ContentFormat() == message.AppLinkFormat
}

func (r *request) SetContentFormat(mt MediaType) {
	r.msg.SetContentFormat(mt)
}

func (r *request) SecurityIdentity() string {
	id, ok := r.message().Context().Value(keyClientCertCommonName).(string)
	if !ok {
		return ""
	}
	return id
}
