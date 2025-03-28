package coap

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	piondtls "github.com/pion/dtls/v2"
	"github.com/plgd-dev/go-coap/v3/dtls"
	"github.com/plgd-dev/go-coap/v3/dtls/server"
	"github.com/plgd-dev/go-coap/v3/message/pool"
	coapnet "github.com/plgd-dev/go-coap/v3/net"
	"github.com/plgd-dev/go-coap/v3/options"
	"github.com/plgd-dev/go-coap/v3/tcp"
	tcpclt "github.com/plgd-dev/go-coap/v3/tcp/client"
	tcpsrv "github.com/plgd-dev/go-coap/v3/tcp/server"
	"github.com/plgd-dev/go-coap/v3/udp"
	udpclt "github.com/plgd-dev/go-coap/v3/udp/client"
	udpsrv "github.com/plgd-dev/go-coap/v3/udp/server"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

const (
	keyClientSecurityIdentity = "securityId"
)

type bearerDescriptor struct {
	create func() error
	close  func() error
	serve  func() error
}

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

	bearers map[string]*bearerDescriptor

	// udp and dtls are not unified,
	// so we need to maintain them separately
	udpListener *coapnet.UDPConn
	udpDelegate *udpsrv.Server

	dtlsListener server.Listener
	dtlsDelegate *server.Server

	tcpListener *coapnet.TCPListener
	tlsListener *coapnet.TLSListener
	tcpDelegate *tcpsrv.Server

	conns sync.Map
}

func NewServer(network, addr string, opts ...PeerOption) Server {
	r := NewRouter()
	s := &coapServer{
		peer:    newPeer(r),
		network: network,
		address: addr,
		bearers: make(map[string]*bearerDescriptor),
	}

	s.bearers = map[string]*bearerDescriptor{
		UDPBearer: {
			create: s.newUdp,
			serve:  s.serveUdp,
			close:  s.closeUdp,
		},
		TCPBearer: {
			create: s.newTcp,
			serve:  s.serveTcp,
			close:  s.closeTcp,
		},
	}

	bearer, ok := s.bearers[network]
	if !ok {
		log.Errorf("unsupported bearer: %s", network)
		return nil
	}

	for _, fn := range opts {
		fn(s.peer)
	}

	err := bearer.create()
	if err != nil {
		return nil
	}

	return s
}

func (s *coapServer) newUdp() error {
	if s.tlsOn {
		s.dtlsDelegate = dtls.NewServer(options.WithMux(s.peer.router),
			options.WithOnNewConn(s.newUdpConnCallback),
			options.WithTransmission(1, 500*time.Millisecond, 4),
			options.WithPeriodicRunner(func(f func(now time.Time) bool) {
				go func() {
					for f(time.Now()) {
						time.Sleep(1 * time.Second)
					}
				}()
			}))

		l, err := coapnet.NewDTLSListener(s.network, s.address, s.dtlsConf)
		if err != nil {
			log.Errorln("new listener failed:", err)
			return err
		}

		s.dtlsListener = l
	} else {
		s.udpDelegate = udp.NewServer(options.WithMux(s.peer.router),
			options.WithOnNewConn(s.newUdpConnCallback),
			options.WithTransmission(1, 400*time.Millisecond, 4),
			options.WithPeriodicRunner(func(f func(now time.Time) bool) {
				go func() {
					for f(time.Now()) {
						time.Sleep(1 * time.Second)
					}
				}()
			}))

		l, err := coapnet.NewListenUDP(s.network, s.address)
		if err != nil {
			log.Errorln("new listener failed:", err)
			return err
		}

		s.udpListener = l
	}

	return nil
}

func (s *coapServer) newTcp() error {
	if s.tlsOn {
		s.tcpDelegate = tcp.NewServer(options.WithMux(s.peer.router),
			options.WithOnNewConn(s.newTcpConnCallback),
			options.WithPeriodicRunner(func(f func(now time.Time) bool) {
				go func() {
					for f(time.Now()) {
						time.Sleep(1 * time.Second)
					}
				}()
			}))

		l, err := coapnet.NewTLSListener(s.network, s.address, s.tlsConf)
		if err != nil {
			log.Errorln("new tls listener failed:", err)
			return err
		}

		s.tlsListener = l
	} else {
		s.tcpDelegate = tcp.NewServer(options.WithMux(s.peer.router),
			options.WithOnNewConn(s.newTcpConnCallback),
			options.WithPeriodicRunner(func(f func(now time.Time) bool) {
				go func() {
					for f(time.Now()) {
						time.Sleep(1 * time.Second)
					}
				}()
			}))

		l, err := coapnet.NewTCPListener(s.network, s.address)
		if err != nil {
			log.Errorln("new tcp listener failed:", err)
			return err
		}

		s.tcpListener = l
	}

	return nil
}

func (s *coapServer) newTcpConnCallback(cc *tcpclt.Conn) {
	s.conns.Store(cc.RemoteAddr().String(), cc)
	log.Infof("connection accepted: %s-%p", cc.RemoteAddr().String(), cc)

	if s.tlsConf != nil {
		state := cc.NetConn().(*tls.Conn).ConnectionState()
		if state.PeerCertificates != nil { // certificate mode
			clientCert := state.PeerCertificates[0]
			if len(clientCert.Subject.CommonName) != 0 {
				cc.SetContextValue(keyClientSecurityIdentity, clientCert.Subject.CommonName)
			}
		} else { // psk mode or raw public key mode
			log.Fatalf("TLS must have common name provided")
			//log.Warnf("TLS must have common name provided")
		}
	}

	cc.AddOnClose(func() {
		log.Infof("connection released: %s-%p", cc.RemoteAddr().String(), cc)
		s.conns.Delete(cc.RemoteAddr().String())
	})
}

func (s *coapServer) newUdpConnCallback(cc *udpclt.Conn) {
	s.conns.Store(cc.RemoteAddr().String(), cc)
	log.Infof("connection accepted: %s-%p", cc.RemoteAddr().String(), cc)

	// save  if dtls enabled
	if s.dtlsConf != nil {
		state := cc.NetConn().(*piondtls.Conn).ConnectionState()
		if state.PeerCertificates != nil { // certificate mode
			if clientCert, err := x509.ParseCertificate(state.PeerCertificates[0]); err == nil {
				if len(clientCert.Subject.CommonName) != 0 {
					cc.SetContextValue(keyClientSecurityIdentity, clientCert.Subject.CommonName)
				}
			}
		} else { // psk mode or raw public key mode
			if state.IdentityHint != nil {
				cc.SetContextValue(keyClientSecurityIdentity, state.IdentityHint)
			}
		}
	}

	cc.AddOnClose(func() {
		log.Infof("connection released: %s-%p", cc.RemoteAddr().String(), cc)
		s.conns.Delete(cc.RemoteAddr().String())
	})
}

func (s *coapServer) serveUdp() error {
	if s.tlsOn {
		return s.dtlsDelegate.Serve(s.dtlsListener)
	} else {
		return s.udpDelegate.Serve(s.udpListener)
	}
}

func (s *coapServer) serveTcp() error {
	if s.tlsOn {
		return s.tcpDelegate.Serve(s.tlsListener)
	} else {
		return s.tcpDelegate.Serve(s.tcpListener)
	}
}

func (s *coapServer) Serve() error {
	return s.bearers[s.network].serve()
}

func (s *coapServer) closeUdp() error {
	if s.tlsOn {
		return s.dtlsListener.Close()
	} else {
		return s.udpListener.Close()
	}
}

func (s *coapServer) closeTcp() error {
	if s.tlsOn {
		return s.tlsListener.Close()
	} else {
		return s.tcpListener.Close()
	}
}

func (s *coapServer) Shutdown() {
	_ = s.bearers[s.network].close()
}

func (s *coapServer) SendTo(addr string, req Request) (Response, error) {
	c, ok := s.conns.Load(addr)
	if !ok {
		log.Errorf("remote peer address %s is not found", addr)
		return nil, fmt.Errorf("remote peer address %s is not found", addr)
	}

	ctx, cancel := context.WithTimeout(context.Background(), req.Timeout())
	defer cancel()

	req.message().SetContext(ctx)

	var rsp *pool.Message
	var err error
	switch s.network {
	case UDPBearer:
		cc := c.(*udpclt.Conn)
		rsp, err = cc.Do(req.message().Message)
	case TCPBearer:
		cc := c.(*tcpclt.Conn)
		rsp, err = cc.Do(req.message().Message)
	}

	return NewResponse(rsp), err
}
