package coap

import (
	"net"
	"strconv"
	"strings"
	"time"
)

// NewRequest creates a new coap Request
func NewRequest(t uint8, methodCode Code, id uint16) Request {
	msg := NewMessage(t, methodCode, id)
	msg.Token = []byte(GenerateToken(8))

	return &DefaultCoapRequest{
		msg: msg,
	}
}

func NewConfirmableGetRequest() Request {
	msg := NewMessage(MessageConfirmable, Get, GenerateMessageID())
	msg.Token = []byte(GenerateToken(8))

	return &DefaultCoapRequest{
		msg: msg,
	}
}

func NewConfirmablePostRequest() Request {
	msg := NewMessage(MessageConfirmable, Post, GenerateMessageID())
	msg.Token = []byte(GenerateToken(8))

	return &DefaultCoapRequest{
		msg: msg,
	}
}

func NewConfirmablePutRequest() Request {
	msg := NewMessage(MessageConfirmable, Put, GenerateMessageID())
	msg.Token = []byte(GenerateToken(8))

	return &DefaultCoapRequest{
		msg: msg,
	}
}

func NewConfirmableDeleteRequest() Request {
	msg := NewMessage(MessageConfirmable, Delete, GenerateMessageID())
	msg.Token = []byte(GenerateToken(8))

	return &DefaultCoapRequest{
		msg: msg,
	}
}

// NewRequestFromMessage creates a new request messages from a CoAP Message.
func NewRequestFromMessage(msg *Message) Request {
	return &DefaultCoapRequest{
		msg: msg,
	}
}

func NewClientRequestFromMessage(msg *Message, attrs map[string]string, conn *net.UDPConn, addr *net.UDPAddr) Request {
	return &DefaultCoapRequest{
		msg:   msg,
		attrs: attrs,
		conn:  conn,
		addr:  addr,
	}
}

type Request interface {
	Conn() *net.UDPConn
	Address() *net.UDPAddr
	Attributes() map[string]string
	Attribute(o string) string
	AttributeAsInt(o string) int
	Message() *Message
	UriQuery(q string) string

	SetProxyUri(uri string)
	SetMediaType(mt MediaType)
	SetPayload([]byte)
	SetStringPayload(s string)
	SetRequestUri(uri string)
	SetConfirmable(con bool)
	SetToken(t string)
	SetUriQuery(k string, v string)
	SetTimeout(to time.Duration)
	GetTimeout() time.Duration
}

// DefaultCoapRequest wraps a CoAP Message as a Request
// Provides various methods which proxies the Message object methods
type DefaultCoapRequest struct {
	msg    *Message
	attrs  map[string]string
	conn   *net.UDPConn
	addr   *net.UDPAddr
	server *Server
	to     time.Duration
}

func (c *DefaultCoapRequest) SetProxyUri(uri string) {
	c.msg.AddOption(OptionProxyURI, uri)
}

func (c *DefaultCoapRequest) SetMediaType(mt MediaType) {
	c.msg.AddOption(OptionContentFormat, mt)
}

func (c *DefaultCoapRequest) Conn() *net.UDPConn {
	return c.conn
}

func (c *DefaultCoapRequest) Address() *net.UDPAddr {
	return c.addr
}

func (c *DefaultCoapRequest) Attributes() map[string]string {
	return c.attrs
}

func (c *DefaultCoapRequest) Attribute(o string) string {
	return c.attrs[o]
}

func (c *DefaultCoapRequest) AttributeAsInt(o string) int {
	attr := c.Attribute(o)
	i, _ := strconv.Atoi(attr)

	return i
}

func (c *DefaultCoapRequest) Message() *Message {
	return c.msg
}

func (c *DefaultCoapRequest) SetStringPayload(s string) {
	c.msg.Payload = NewPlainTextPayload(s)
}

func (c *DefaultCoapRequest) SetPayload(b []byte) {
	c.msg.Payload = NewBytesPayload(b)
}

func (c *DefaultCoapRequest) SetRequestUri(uri string) {
	c.msg.AddOptions(NewPathOptions(uri))
}

func (c *DefaultCoapRequest) SetConfirmable(con bool) {
	if con {
		c.msg.Type = MessageConfirmable
	} else {
		c.msg.Type = MessageNonConfirmable
	}
}

func (c *DefaultCoapRequest) SetToken(t string) {
	c.msg.Token = []byte(t)
}

func (c *DefaultCoapRequest) UriQuery(q string) string {
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

func (c *DefaultCoapRequest) SetUriQuery(k string, v string) {
	c.Message().AddOption(OptionURIQuery, k+"="+v)
}

func (c *DefaultCoapRequest) SetTimeout(t time.Duration) {
	c.to = t
}

func (c *DefaultCoapRequest) GetTimeout() time.Duration {
	return c.to
}
