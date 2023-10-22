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
