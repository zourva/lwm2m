package server

type LifecycleHandler interface {
	OnClientRegistered(session *RegisteredClient)
}

type DefaultLifecycleHandler struct {
	LifecycleHandler
}

func (a *DefaultLifecycleHandler) OnClientRegistered(session *RegisteredClient) {

}

func (a *DefaultLifecycleHandler) OnClientUpdated(session *RegisteredClient) {

}

func (a *DefaultLifecycleHandler) OnClientDeregistered(session *RegisteredClient) {

}

func NewDefaultLifecycleHandler() LifecycleHandler {
	return &DefaultLifecycleHandler{}
}
