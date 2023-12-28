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
	gonet "net"
)

type Client interface {
	Peer

	// Send sends request to the server
	// currently connected and expects
	// a response from remote.
	Send(req Request) (Response, error)
	Notify(key string, value []byte) error
	Close() error
}

type coapClient struct {
	*peer
	delegate *udpclt.Conn
}

var _ Client = &coapClient{}

func Dial(server string, opts ...PeerOption) (Client, error) {
	c := &coapClient{
		peer: newPeer(NewRouter()),
	}

	for _, fn := range opts {
		fn(c.peer)
	}

	if c.dtlsOn {
		dial, err := dtls.Dial(server, c.dtlsConf)
		if err != nil {
			log.Errorf("error dialing dtls: %v", err)
			return nil, err
		}

		//dial.NetConn().SetReadDeadline()
		//dial.NetConn().SetReadDeadline()
		//err := c.delegate.Session().NetConn().(*piondtls.Conn).SetReadBuffer(c.readBufferSize)
		//if err != nil {
		//	log.Errorf("error set reader buffer size: %v", err)
		//	return nil, err
		//}

		c.delegate = dial
	} else {
		dial, err := udp.Dial(server)
		if err != nil {
			log.Errorf("error dialing dtls: %v", err)
			return nil, err
		}

		err = c.delegate.Session().NetConn().(*gonet.UDPConn).SetWriteBuffer(c.writeBufferSize)
		if err != nil {
			log.Errorf("error set write buffer size: %v", err)
			return nil, err
		}

		err = c.delegate.Session().NetConn().(*gonet.UDPConn).SetReadBuffer(c.readBufferSize)
		if err != nil {
			log.Errorf("error set read buffer size: %v", err)
			return nil, err
		}

		c.delegate = dial
	}

	return c, nil
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

func (s *coapClient) Close() error {
	return s.delegate.Close()
}
