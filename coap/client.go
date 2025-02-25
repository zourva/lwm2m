package coap

import (
	"bytes"
	"context"
	"fmt"
	"github.com/plgd-dev/go-coap/v3/dtls"
	"github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/message/codes"
	"github.com/plgd-dev/go-coap/v3/mux"
	"github.com/plgd-dev/go-coap/v3/options"
	"github.com/plgd-dev/go-coap/v3/tcp"
	"github.com/plgd-dev/go-coap/v3/udp"
	log "github.com/sirupsen/logrus"
	gonet "net"
	"time"
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
	bearer mux.Conn
}

var _ Client = &coapClient{}

// Dial connects to server using the given bearer, address and options.
// Supported bearer includes: udp(coap)/tcp(coap)/mqtt/http.
func Dial(bearer, address string, opts ...PeerOption) (Client, error) {
	c := &coapClient{
		peer: newPeer(NewRouter()),
	}

	for _, fn := range opts {
		fn(c.peer)
	}

	var err error
	switch bearer {
	case UDPBearer: //coap over udp/dtls
		err = c.dialUdp(address)
	case TCPBearer: //coap over tcp/tls
		err = c.dialTcp(address)
	case MQTTBearer: //mqtt/mqtts over tcp/tls
		err = c.dialMqtt(address)
	case HTTPBearer: //http/https over tcp/tls
		err = c.dialHttp(address)
	default:
		log.Errorf("error dialing: unsupported bearer: %s", bearer)
		return nil, fmt.Errorf("unsupported bearer: %s", bearer)
	}

	return c, err
}

func (s *coapClient) dialHttp(address string) error {
	if s.tlsOn {
		//TODO
	} else {
		//TODO
	}

	return fmt.Errorf("http bearer is not supported yet")
}

func (s *coapClient) dialMqtt(address string) error {
	if s.tlsOn {
		//TODO
	} else {
		//TODO
	}

	return fmt.Errorf("mqtt bearer is not supported yet")
}

func (s *coapClient) dialTcp(address string) error {
	// In TCP dialing, disable the block-wise option.
	opts := []tcp.Option{options.WithBlockwise(false, 0, 0)}
	if s.tlsOn {
		opts = append(opts, options.WithTLS(s.tlsConf))
	}

	dial, err := tcp.Dial(address, opts...)
	if err != nil {
		log.Errorf("error dialing tcp: %v", err)
		return err
	}

	s.bearer = dial

	return nil
}

func (s *coapClient) dialUdp(address string) error {
	if s.tlsOn {
		dial, err := dtls.Dial(address, s.dtlsConf, options.WithMux(s.Router()),
			options.WithTransmission(1, 500*time.Millisecond, 4),
			options.WithPeriodicRunner(func(f func(now time.Time) bool) {
				go func() {
					for f(time.Now()) {
						time.Sleep(1 * time.Second)
					}
				}()
			}))
		if err != nil {
			log.Errorf("error dialing dtls: %v", err)
			return err
		}

		//err := c.bearer.Session().NetConn().(*piondtls.Conn).SetReadBuffer(c.readBufferSize)
		//if err != nil {
		//	log.Errorf("error set reader buffer size: %v", err)
		//	return nil, err
		//}

		s.bearer = dial
	} else {
		dial, err := udp.Dial(address, options.WithMux(s.Router()),
			options.WithTransmission(1, 400*time.Millisecond, 4),
			options.WithPeriodicRunner(func(f func(now time.Time) bool) {
				go func() {
					for f(time.Now()) {
						time.Sleep(1 * time.Second)
					}
				}()
			}))
		if err != nil {
			log.Errorf("error dialing dtls: %v", err)
			return err
		}

		err = dial.Session().NetConn().(*gonet.UDPConn).SetWriteBuffer(s.writeBufferSize)
		if err != nil {
			log.Errorf("error set write buffer size: %v", err)
			return err
		}

		err = dial.Session().NetConn().(*gonet.UDPConn).SetReadBuffer(s.readBufferSize)
		if err != nil {
			log.Errorf("error set read buffer size: %v", err)
			return err
		}

		s.bearer = dial
	}

	return nil
}

func (s *coapClient) Send(req Request) (Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), req.Timeout())
	defer cancel()

	req.message().SetContext(ctx)
	msg := req.message().Message
	rsp, err := s.bearer.Do(msg)

	log.Tracef("make request to %v, req: %v, rsp: %v",
		s.bearer.RemoteAddr(), msg, rsp)

	return NewResponse(rsp), err
}

func (s *coapClient) Notify(observationId string, data []byte) error {
	m := s.bearer.AcquireMessage(s.bearer.Context())
	defer s.bearer.ReleaseMessage(m)
	m.SetCode(codes.Content)
	//m.SetToken(token)
	m.SetBody(bytes.NewReader(data))
	m.SetContentFormat(message.TextPlain)
	//m.SetObserve(uint32(obs))

	err := s.bearer.WriteMessage(m)
	log.Tracef("notify %v of observation %s, msg: %v, err: %v",
		s.bearer.RemoteAddr(), observationId, m, err)

	return err
}

func (s *coapClient) Close() error {
	return s.bearer.Close()
}
