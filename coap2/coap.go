package coap2

import (
	piondtls "github.com/pion/dtls/v2"
	"github.com/plgd-dev/go-coap/v3/dtls"
	"github.com/plgd-dev/go-coap/v3/dtls/server"
	"github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/message/codes"
	"github.com/plgd-dev/go-coap/v3/mux"
	"github.com/plgd-dev/go-coap/v3/net"
	"github.com/plgd-dev/go-coap/v3/options"
	udpClient "github.com/plgd-dev/go-coap/v3/udp/client"
	log "github.com/sirupsen/logrus"
)

type MediaType = message.MediaType
type Message = mux.Message

type Handler = mux.HandlerFunc
type Interceptor = mux.Handler
type Router = mux.Router

type ResponseWriter = mux.ResponseWriter

type Request interface {
	Connection() *net.UDPConn
	Attributes() map[string]string
	Attribute(o string) string
	AttributeAsInt(o string) int
	Message() *Message
	UriQuery(q string) string

	SetProxyUri(uri string)
	SetPayload([]byte)
	SetStringPayload(s string)
	SetRequestUri(uri string)
	SetConfirmable(con bool)
	SetToken(t string)
	SetUriQuery(k string, v string)
}

type Response interface {
	Message() *Message
	Error() error
	Payload() []byte
	UriQuery(q string) string
}

type PatternHandler = func(Request) Response
type RouteHandler = PatternHandler

//func ListenAndServe(s *Server, network, addr string) error {
//	l, err := net.NewListenUDP(network, addr)
//	if err != nil {
//		return err
//	}
//	defer func() {
//		if errC := l.Close(); errC != nil && err == nil {
//			err = errC
//		}
//	}()
//	return s.Serve(l)
//}

func ListenAndServeDTLS(s *Server, network, addr string) error {
	l, err := net.NewDTLSListener(network, addr, &piondtls.Config{})
	if err != nil {
		return err
	}
	defer func() {
		if errC := l.Close(); errC != nil && err == nil {
			err = errC
		}
	}()
	return s.Serve(l)
}

type request struct {
	Request
	msg *mux.Message
}

type response struct {
	Response
}

type Client struct {
	*udpClient.Conn
}

func NewClient(server string, dtlsConf *piondtls.Config) *Client {
	c := &Client{}

	dial, err := dtls.Dial(server, dtlsConf)
	if err != nil {
		log.Fatalf("error dialing dtls: %v", err)
	}

	c.Conn = dial

	return c
}

func (s *Client) Route(method codes.Code, pattern string, h PatternHandler) error {
	return nil
}

type Server struct {
	*server.Server
	router *Router
}

func NewServer(r *Router) *Server {
	s := &Server{
		router: r,
		//Server: udp.NewServer(options.WithMux(r)),
		Server: dtls.NewServer(options.WithMux(r)),
	}

	return s
}

func (s *Server) rrWrapper(h PatternHandler, w mux.ResponseWriter, r *mux.Message) {
	rsp := h(&request{msg: r})

	// get msg from pool and return afterward
	rspMsg := w.Conn().AcquireMessage(r.Context())
	defer w.Conn().ReleaseMessage(rspMsg)

	rspMsg.SetCode(rsp.Message().Code())
	rspMsg.SetToken(rsp.Message().Token())
	mt, _ := rsp.Message().ContentFormat()

	rspMsg.SetContentFormat(mt)
	rspMsg.SetBody(rsp.Message().Body())
	err := w.Conn().WriteMessage(rspMsg)
	if err != nil {
		log.Errorf("coap cannot set response: %v", err)
	}
}

func (s *Server) regHandler(path string, h PatternHandler) error {
	return s.router.Handle(path, mux.HandlerFunc(func(w mux.ResponseWriter, r *mux.Message) {
		s.rrWrapper(h, w, r)
	}))
}

func (s *Server) Get(path string, h PatternHandler) error {
	return s.regHandler(path, h)
}

func (s *Server) Delete(path string, h PatternHandler) error {
	return s.regHandler(path, h)
}

func (s *Server) Put(path string, h PatternHandler) error {
	return s.regHandler(path, h)
}

func (s *Server) Post(path string, h PatternHandler) error {
	return s.regHandler(path, h)
}

func (s *Server) Options(path string, h PatternHandler) error {
	return s.regHandler(path, h)
}

func (s *Server) Patch(path string, h PatternHandler) error {
	return s.regHandler(path, h)
}

func NewRouter() *Router {
	return mux.NewRouter()
}
