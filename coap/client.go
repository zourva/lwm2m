package coap

import (
	"bytes"
	"context"
	"github.com/plgd-dev/go-coap/v3/dtls"
	"github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/message/codes"
	"github.com/plgd-dev/go-coap/v3/udp"
	udpclt "github.com/plgd-dev/go-coap/v3/udp/client"
	log "github.com/sirupsen/logrus"
)

type Client interface {
	Peer

	// Send sends request to the server
	// currently connected and expects
	// a response from remote.
	Send(req Request) (Response, error)
	Notify(key string, value []byte) error
}

type coapClient struct {
	*peer
	delegate *udpclt.Conn
}

func NewClient(server string, opts ...PeerOption) Client {
	c := &coapClient{
		peer: newPeer(NewRouter()),
	}

	for _, fn := range opts {
		fn(c.peer)
	}

	if c.dtlsOn {
		dial, err := dtls.Dial(server, c.dtlsConf)
		if err != nil {
			log.Fatalf("error dialing dtls: %v", err)
		}

		//dial.NetConn().SetReadDeadline()
		//dial.NetConn().SetReadDeadline()

		c.delegate = dial
	} else {
		dial, err := udp.Dial(server)
		if err != nil {
			log.Fatalf("error dialing udp: %v", err)
		}

		c.delegate = dial
	}

	return c
}

func (s *coapClient) Send(req Request) (Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), req.Timeout())
	defer cancel()

	req.message().SetContext(ctx)
	rsp, err := s.delegate.Do(req.message().Message)
	return NewResponse(rsp), err
}

func (s *coapClient) Notify(observationId string, data []byte) error {
	m := s.delegate.AcquireMessage(s.delegate.Context())
	defer s.delegate.ReleaseMessage(m)
	m.SetCode(codes.Content)
	//m.SetToken(token)
	m.SetBody(bytes.NewReader(data))
	m.SetContentFormat(message.TextPlain)
	//m.SetObserve(uint32(obs))

	return s.delegate.WriteMessage(m)
}
