package coap

import (
	"net"
	"strconv"
	"strings"
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
	SetProxyURI(uri string)
	SetMediaType(mt MediaType)
	GetConnection() *net.UDPConn
	GetAddress() *net.UDPAddr
	GetAttributes() map[string]string
	GetAttribute(o string) string
	GetAttributeAsInt(o string) int
	GetMessage() *Message
	SetPayload([]byte)
	SetStringPayload(s string)
	SetRequestURI(uri string)
	SetConfirmable(con bool)
	SetToken(t string)
	GetURIQuery(q string) string
	SetURIQuery(k string, v string)
}

// DefaultCoapRequest wraps a CoAP Message as a Request
// Provides various methods which proxies the Message object methods
type DefaultCoapRequest struct {
	msg    *Message
	attrs  map[string]string
	conn   *net.UDPConn
	addr   *net.UDPAddr
	server *Server
}

func (c *DefaultCoapRequest) SetProxyURI(uri string) {
	c.msg.AddOption(OptionProxyURI, uri)
}

func (c *DefaultCoapRequest) SetMediaType(mt MediaType) {
	c.msg.AddOption(OptionContentFormat, mt)
}

func (c *DefaultCoapRequest) GetConnection() *net.UDPConn {
	return c.conn
}

func (c *DefaultCoapRequest) GetAddress() *net.UDPAddr {
	return c.addr
}

func (c *DefaultCoapRequest) GetAttributes() map[string]string {
	return c.attrs
}

func (c *DefaultCoapRequest) GetAttribute(o string) string {
	return c.attrs[o]
}

func (c *DefaultCoapRequest) GetAttributeAsInt(o string) int {
	attr := c.GetAttribute(o)
	i, _ := strconv.Atoi(attr)

	return i
}

func (c *DefaultCoapRequest) GetMessage() *Message {
	return c.msg
}

func (c *DefaultCoapRequest) SetStringPayload(s string) {
	c.msg.Payload = NewPlainTextPayload(s)
}

func (c *DefaultCoapRequest) SetPayload(b []byte) {
	c.msg.Payload = NewBytesPayload(b)
}

func (c *DefaultCoapRequest) SetRequestURI(uri string) {
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

func (c *DefaultCoapRequest) GetURIQuery(q string) string {
	qs := c.GetMessage().GetOptionsAsString(OptionURIQuery)

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

func (c *DefaultCoapRequest) SetURIQuery(k string, v string) {
	c.GetMessage().AddOption(OptionURIQuery, k+"="+v)
}
