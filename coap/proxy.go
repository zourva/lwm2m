package coap

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
)

// Proxy Filter
type ProxyFilter func(*Message, *net.UDPAddr) bool

func NullProxyFilter(*Message, *net.UDPAddr) bool {
	return true
}

type ProxyHandler func(c Server, msg *Message, conn *net.UDPConn, addr *net.UDPAddr)

// The default handler when proxying is disabled
func NullProxyHandler(c Server, msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	SendMessageTo(&MessageContext{server: c, msg: ProxyingNotSupportedMessage(msg.Id, MessageAcknowledgment), conn: NewUDPConnection(conn), addr: addr})
}

func COAPProxyHandler(c Server, msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	proxyURI := msg.GetOption(OptionProxyURI).StringValue()

	parsedURL, err := url.Parse(proxyURI)
	if err != nil {
		log.Println("Error parsing proxy URI")
		SendMessageTo(&MessageContext{server: c, msg: BadGatewayMessage(msg.Id, MessageAcknowledgment), conn: NewUDPConnection(conn), addr: addr})
		return
	}

	client := NewCoapClient("Proxy Client")
	client.OnStart(func(server Server) {
		client.Dial(parsedURL.Host)

		msg.RemoveOptions(OptionProxyURI)
		req := NewRequestFromMessage(msg)
		req.SetRequestUri(parsedURL.RequestURI())

		response, err := client.Send(req)
		if err != nil {
			SendMessageTo(&MessageContext{server: c, msg: BadGatewayMessage(msg.Id, MessageAcknowledgment), conn: NewUDPConnection(conn), addr: addr})
			client.Stop()
			return
		}

		_, err = SendMessageTo(&MessageContext{server: c, msg: response.Message(), conn: NewUDPConnection(conn), addr: addr})
		if err != nil {
			log.Println("Error occured responding to proxy request")
			client.Stop()
			return
		}
		client.Stop()

	})
	client.Start()
}

// Handles requests for proxying from CoAP to HTTP
func HTTPProxyHandler(c Server, msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	proxyURI := msg.GetOption(OptionProxyURI).StringValue()
	requestMethod := msg.Code

	client := &http.Client{}
	req, err := http.NewRequest(MethodString(Code(msg.GetMethod())), proxyURI, nil)
	if err != nil {
		SendMessageTo(&MessageContext{server: c, msg: BadGatewayMessage(msg.Id, MessageAcknowledgment), conn: NewUDPConnection(conn), addr: addr})
		return
	}

	etag := msg.GetOption(OptionEtag)
	if etag != nil {
		req.Header.Add("ETag", etag.StringValue())
	}

	// TODO: Set timeout handler, and on timeout return 5.04
	resp, err := client.Do(req)
	if err != nil {
		SendMessageTo(&MessageContext{server: c, msg: BadGatewayMessage(msg.Id, MessageAcknowledgment), conn: NewUDPConnection(conn), addr: addr})
		return
	}

	defer resp.Body.Close()

	contents, _ := ioutil.ReadAll(resp.Body)
	msg.Payload = NewBytesPayload(contents)
	respMsg := NewRequestFromMessage(msg)

	if requestMethod == Get {
		etag := resp.Header.Get("ETag")
		if etag != "" {
			msg.AddOption(OptionEtag, etag)
		}
	}

	// TODO: Check payload length against Size1 options
	if len(respMsg.Message().Payload.String()) > MaxPacketSize {
		SendMessageTo(&MessageContext{server: c, msg: BadGatewayMessage(msg.Id, MessageAcknowledgment), conn: NewUDPConnection(conn), addr: addr})
		return
	}

	_, err = SendMessageTo(&MessageContext{server: c, msg: respMsg.Message(), conn: NewUDPConnection(conn), addr: addr})
	if err != nil {
		println(err.Error())
	}
}

// Handles requests for proxying from HTTP to CoAP
func HTTPCOAPProxyHandler(msg *Message, conn *net.UDPConn, addr *net.UDPAddr) {
	log.Println("HttpCoapProxyHandler Proxy Handler")
}
