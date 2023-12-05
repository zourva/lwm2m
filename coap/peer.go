package coap

import (
	"bytes"
	"context"
	piondtls "github.com/pion/dtls/v2"
	"github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/message/codes"
	"github.com/plgd-dev/go-coap/v3/message/pool"
	"github.com/plgd-dev/go-coap/v3/mux"
	log "github.com/sirupsen/logrus"
)

type PeerOption func(peer Peer)

func WithDTLSConfig(dtlsConf *piondtls.Config) PeerOption {
	return func(peer Peer) {
		peer.EnableDTLS(dtlsConf)
	}
}

// Peer defines a LwM2M peer which
// may be run as a server or a client or both.
type Peer interface {
	Router() *Router

	Post(pattern string, h PatternHandler) error
	Get(pattern string, h PatternHandler) error
	Put(pattern string, h PatternHandler) error
	Delete(pattern string, h PatternHandler) error
	Options(pattern string, h PatternHandler) error
	Patch(pattern string, h PatternHandler) error

	NewGetRequestPlain(uri string) Request
	NewDeleteRequestPlain(uri string) Request
	NewPutRequestPlain(uri string, body []byte) Request
	NewPostRequestPlain(uri string, body []byte) Request
	NewPostRequestOpaque(uri string, body []byte) Request
	NewPostRequestCoReLink(uri string, body []byte) Request

	NewAckResponse(req Request, code Code) Response
	NewAckPiggybackedResponse(req Request, code Code, body []byte) Response

	EnableDTLS(conf *piondtls.Config)
}

func newPeer(router *Router) *peer {
	return &peer{
		//pool:   pool.New(msgPoolSize, math.MaxUint16),
		router: router,
	}
}

type peer struct {
	//pool *pool.Pool
	router *Router

	dtlsOn   bool
	dtlsConf *piondtls.Config
}

func (p *peer) Router() *Router {
	return p.router
}

func (p *peer) NewGetRequestPlain(uri string) Request {
	return p.NewConfirmableRequest(Get, message.TextPlain, uri, nil)
}

func (p *peer) NewDeleteRequestPlain(uri string) Request {
	return p.NewConfirmableRequest(Delete, message.TextPlain, uri, nil)
}

func (p *peer) NewPutRequestPlain(uri string, body []byte) Request {
	return p.NewConfirmableRequest(Put, message.TextPlain, uri, body)
}

func (p *peer) NewPostRequestPlain(uri string, body []byte) Request {
	return p.NewConfirmableRequest(Post, message.TextPlain, uri, body)
}

func (p *peer) NewPostRequestOpaque(uri string, body []byte) Request {
	return p.NewConfirmableRequest(Post, message.AppLwm2mCbor, uri, body)
}

func (p *peer) NewPostRequestCoReLink(uri string, body []byte) Request {
	return p.NewConfirmableRequest(Post, message.AppLinkFormat, uri, body)
}

func (p *peer) NewConfirmableRequest(m Code, mt MediaType, uri string, body []byte) Request {
	msg := pool.NewMessage(context.Background())
	req := NewRequest(
		&mux.Message{Message: msg, RouteParams: nil},
	)

	token, _ := message.GetToken()

	switch m {
	case Get:
		_ = req.message().SetupGet(uri, token)
	case Put:
		_ = req.message().SetupPut(uri, token, mt, bytes.NewReader(body))
	case Post:
		_ = req.message().SetupPost(uri, token, mt, bytes.NewReader(body))
	case Delete:
		_ = req.message().SetupDelete(uri, token)
	default:
		_ = req.message().SetupGet(uri, token)
	}

	return req
}

// NewAckResponse creates an ACK response.
//
//	Client              Server
//	   |                  |
//	   |   CON [0x7d34]   |
//	   +----------------->|
//	   |                  |
//	   |   ACK [0x7d34]   |
//	   |<-----------------+
//	   |                  |
func (p *peer) NewAckResponse(req Request, code Code) Response {
	return p.NewAckPiggybackedResponse(req, code, nil)
}

// NewAckPiggybackedResponse creates an ACK-piggybacked response.
func (p *peer) NewAckPiggybackedResponse(req Request, code Code, body []byte) Response {
	msg := pool.NewMessage(req.message().Context())
	msg.SetType(message.Acknowledgement)
	msg.SetCode(codes.Code(code))
	msg.SetMessageID(req.message().MessageID())
	msg.SetToken(req.message().Token())

	if body != nil {
		msg.SetBody(bytes.NewReader(body))
	}

	log.Traceln("new response:", msg)

	return NewResponse(msg)
}

func (p *peer) EnableDTLS(conf *piondtls.Config) {
	if conf == nil {
		return
	}

	p.dtlsConf = conf
	p.dtlsOn = true
}

func (p *peer) rrWrapper(fn PatternHandler, w mux.ResponseWriter, r *mux.Message) {
	//wrap request received
	req := NewRequest(r)
	req.SetAddress(w.Conn().RemoteAddr())

	//invoke handler
	rsp := fn(req)

	//write response to send
	err := w.Conn().WriteMessage(rsp.message().Message)
	if err != nil {
		log.Errorf("coap cannot write response: %v", err)
	}
}

func (p *peer) regHandler(path string, h PatternHandler) error {
	return p.router.Handle(path, mux.HandlerFunc(func(w mux.ResponseWriter, r *mux.Message) {
		p.rrWrapper(h, w, r)
	}))
}

func (p *peer) Get(path string, h PatternHandler) error {
	return p.regHandler(path, h)
}

func (p *peer) Delete(path string, h PatternHandler) error {
	return p.regHandler(path, h)
}

func (p *peer) Put(path string, h PatternHandler) error {
	return p.regHandler(path, h)
}

func (p *peer) Post(path string, h PatternHandler) error {
	return p.regHandler(path, h)
}

func (p *peer) Options(path string, h PatternHandler) error {
	return p.regHandler(path, h)
}

func (p *peer) Patch(path string, h PatternHandler) error {
	return p.regHandler(path, h)
}
