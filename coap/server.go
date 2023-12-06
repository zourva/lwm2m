package coap

import (
	"context"
	"github.com/plgd-dev/go-coap/v3/dtls"
	"github.com/plgd-dev/go-coap/v3/dtls/server"
	coapnet "github.com/plgd-dev/go-coap/v3/net"
	"github.com/plgd-dev/go-coap/v3/options"
	"github.com/plgd-dev/go-coap/v3/udp"
	udpclt "github.com/plgd-dev/go-coap/v3/udp/client"
	udpsrv "github.com/plgd-dev/go-coap/v3/udp/server"
	log "github.com/sirupsen/logrus"
	"sync"
)

type Server interface {
	Peer
	// Serve serves coap service and holds if succeeded.
	Serve() error
	Shutdown()

	// SendTo send request to the remote peer identified by addr.
	SendTo(addr string, req Request) (Response, error)
}

type coapServer struct {
	*peer
	network string
	address string

	udpListener *coapnet.UDPConn
	udpDelegate *udpsrv.Server

	dtlsListener server.Listener
	dtlsDelegate *server.Server

	conns sync.Map
}

func NewServer(network, addr string, opts ...PeerOption) Server {
	r := NewRouter()
	s := &coapServer{
		peer:    newPeer(r),
		network: network,
		address: addr,
	}

	for _, fn := range opts {
		fn(s.peer)
	}

	if s.dtlsOn {
		s.dtlsDelegate = dtls.NewServer(options.WithMux(r),
			options.WithOnNewConn(s.newConnCallback))

		l, err := coapnet.NewDTLSListener(network, addr, s.dtlsConf)
		if err != nil {
			log.Errorln("new listener failed:", err)
			return nil
		}
		s.dtlsListener = l
	} else {
		s.udpDelegate = udp.NewServer(options.WithMux(r),
			options.WithOnNewConn(s.newConnCallback))

		l, err := coapnet.NewListenUDP(network, addr)
		if err != nil {
			log.Errorln("new listener failed:", err)
			return nil
		}

		s.udpListener = l
	}

	return s
}

func (s *coapServer) newConnCallback(cc *udpclt.Conn) {
	s.conns.Store(cc.RemoteAddr().String(), cc)
	log.Infof("connection accepted: %s-%p", cc.RemoteAddr().String(), cc)

	cc.AddOnClose(func() {
		log.Infof("connection released: %s-%p", cc.RemoteAddr().String(), cc)
		s.conns.Delete(cc.RemoteAddr().String())
	})
}

func (s *coapServer) Serve() error {
	if s.dtlsOn {
		return s.dtlsDelegate.Serve(s.dtlsListener)
	} else {
		return s.udpDelegate.Serve(s.udpListener)
	}
}

func (s *coapServer) Shutdown() {
	if s.dtlsOn {
		_ = s.dtlsListener.Close()
	} else {
		_ = s.udpListener.Close()
	}
}

func (s *coapServer) SendTo(addr string, req Request) (Response, error) {
	c, ok := s.conns.Load(addr)
	if !ok {
		log.Errorf("remote peer address %s is not found", addr)
		return nil, nil
	}

	cc := c.(*udpclt.Conn)

	ctx, cancel := context.WithTimeout(context.Background(), req.Timeout())
	defer cancel()

	req.message().SetContext(ctx)
	rsp, err := cc.Do(req.message().Message)
	return NewResponse(rsp), err
}
