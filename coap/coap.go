package coap

import (
	"fmt"
	"github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/message/codes"
	"github.com/plgd-dev/go-coap/v3/mux"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	DefaultTimeout = 5 * time.Second
)

type Message = mux.Message
type Options = message.Options

type Handler = mux.HandlerFunc
type Interceptor = mux.Handler

//type Router = mux.Router

type ResponseWriter = mux.ResponseWriter

type PatternHandler = func(Request) Response
type RouteHandler = PatternHandler

func NewRouter() *Router {
	r := &Router{
		Router: mux.NewRouter(),
		z:      make(map[string]*Route),
	}

	r.Router.DefaultHandleFunc(
		func(w ResponseWriter, m *Message) {
			if m.Type() == message.Acknowledgement {
				// ignore if it's ack
				return
			}
			if err := w.SetResponse(codes.NotFound, message.TextPlain, nil); err != nil {
				log.Errorf("router handler: cannot set response: %v", err)
			}
		})

	return r
}

type Route struct {
	h     Handler
	proxy map[codes.Code]Handler
}

type Router struct {
	*mux.Router

	z map[string]*Route
}

func (r *Router) Handle(method codes.Code, pattern string, handler Handler) error {
	route, ok := r.z[pattern]
	if ok {
		h, ok := route.proxy[method]
		if ok {
			return fmt.Errorf("duplicate path registration, method:%d, pattern:%s, already-handler:%p, new-handler%p",
				method, pattern, h, handler)
		}
		route.proxy[method] = handler
		return nil
	} else {
		route = &Route{
			//h:   r.hook(),
			proxy: make(map[codes.Code]Handler),
		}
		route.h = r.hook(route)
		route.proxy[method] = handler
	}

	err := r.Router.Handle(pattern, route.h)
	if err != nil {
		return err
	}

	r.z[pattern] = route
	return nil
}

func (r *Router) hook(route *Route) Handler {
	return func(w ResponseWriter, req *Message) {
		method := req.Code()
		h, ok := route.proxy[method]
		if ok {
			h(w, req)
		}
	}
}
