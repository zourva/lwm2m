package server

import . "github.com/zourva/lwm2m/core"

type ServerStartedEvent struct {
	*BaseEvent
}

func NewServerStartedEvent(args ...string) Event {
	return &ServerStartedEvent{
		BaseEvent: NewBaseEvent(EventServerStarted, "server started", "", args...),
	}
}

type ServerStoppedEvent struct {
	*BaseEvent
}

func NewServerStoppedEvent(args ...string) Event {
	return &ServerStoppedEvent{
		BaseEvent: NewBaseEvent(EventServerStopped, "server stopped", "", args...),
	}
}

type ClientBootstrappedEvent struct {
	*BaseEvent
}

func NewClientBootstrappedEvent(args ...string) Event {
	return &ClientBootstrappedEvent{
		BaseEvent: NewBaseEvent(EventClientBootstrapped, "client bootstrapped", "", args...),
	}
}

type ClientRegisteredEvent struct {
	*BaseEvent
}

func NewClientRegisteredEvent(args ...string) Event {
	return &ClientRegisteredEvent{
		BaseEvent: NewBaseEvent(EventClientRegistered, "client registered", "", args...),
	}
}

type ClientRegUpdatedEvent struct {
	*BaseEvent
}

func NewClientRegUpdatedEvent(args ...string) Event {
	return &ClientRegUpdatedEvent{
		BaseEvent: NewBaseEvent(EventClientRegUpdated, "client registration updated", "", args...),
	}
}

type ClientUnregisteredEvent struct {
	*BaseEvent
}

func NewClientUnregisteredEvent(args ...string) Event {
	return &ClientUnregisteredEvent{
		BaseEvent: NewBaseEvent(EventClientUnregistered, "client unregistered", "", args...),
	}
}
