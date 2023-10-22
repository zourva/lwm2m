package server

import "github.com/zourva/lwm2m/core"

type LifecycleHandler interface {
	OnClientRegistered(c core.RegisteredClient)
}

type DefaultLifecycleHandler struct {
	LifecycleHandler
}

func (a *DefaultLifecycleHandler) OnClientRegistered(session core.RegisteredClient) {

}

func (a *DefaultLifecycleHandler) OnClientUpdated(session core.RegisteredClient) {

}

func (a *DefaultLifecycleHandler) OnClientDeregistered(session core.RegisteredClient) {

}

func NewDefaultLifecycleHandler() LifecycleHandler {
	return &DefaultLifecycleHandler{}
}
