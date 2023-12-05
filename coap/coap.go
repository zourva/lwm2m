package coap

import (
	"github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/mux"
	"time"
)

const (
	DefaultTimeout = 5 * time.Second
)

type Message = mux.Message
type Options = message.Options

type Handler = mux.HandlerFunc
type Interceptor = mux.Handler
type Router = mux.Router

type ResponseWriter = mux.ResponseWriter

type PatternHandler = func(Request) Response
type RouteHandler = PatternHandler

func NewRouter() *Router {
	return mux.NewRouter()
}
