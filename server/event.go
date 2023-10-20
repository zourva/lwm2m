package server

import . "github.com/zourva/lwm2m/core"

type ServerStartedEvent struct {
	*BaseEvent
}

func NewServerStartedEvent(args ...string) Event {
	return &ServerStartedEvent{
		BaseEvent: newEvt(EventServerStarted, "server started", "", args...),
	}
}

type ServerStoppedEvent struct {
	*BaseEvent
}

func NewServerStoppedEvent(args ...string) Event {
	return &ServerStoppedEvent{
		BaseEvent: newEvt(EventServerStopped, "server stopped", "", args...),
	}
}

func optString(opt, def string) string {
	if len(opt) != 0 {
		return opt
	}

	return def
}

func newEvt(evt EventType, defName, defMsg string, args ...string) *BaseEvent {
	name := optString(args[0], defName)
	msg := optString(args[1], defMsg)
	return NewBaseEvent(evt, name, msg)
}
