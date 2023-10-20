package client

import . "github.com/zourva/lwm2m/core"

type BootstrappedEvent struct {
	*BaseEvent
}

func NewBootstrappedEvent(args ...string) Event {
	return &BootstrappedEvent{
		BaseEvent: newEvt(EventClientBootstrapped, "bootstrapped", "", args...),
	}
}

type RegisteredEvent struct {
	*BaseEvent
}

func NewRegisteredEvent(args ...string) Event {
	return &RegisteredEvent{
		BaseEvent: newEvt(EventClientRegistered, "registered", "", args...),
	}
}

type RegUpdatedEvent struct {
	*BaseEvent
}

func NewRegUpdatedEvent(args ...string) Event {
	return &RegUpdatedEvent{
		BaseEvent: newEvt(EventClientRegUpdated, "registration updated", "", args...),
	}
}

type UnregisteredEvent struct {
	*BaseEvent
}

func NewUnregisteredEvent(args ...string) Event {
	return &UnregisteredEvent{
		BaseEvent: newEvt(EventClientUnregistered, "unregistered", "", args...),
	}
}

type DeviceChangedEvent struct {
	*BaseEvent
}

func NewDeviceChangedEvent(args ...string) Event {
	return &DeviceChangedEvent{
		BaseEvent: newEvt(EventClientDevInfoChanged, "device control", "", args...),
	}
}

type InfoObservedEvent struct {
	*BaseEvent
}

func NewInfoObservedEvent(args ...string) Event {
	return &InfoObservedEvent{
		BaseEvent: newEvt(EventClientObserved, "observe", "", args...),
	}
}

type InfoReportedEvent struct {
	*BaseEvent
}

func NewInfoReportedEvent(args ...string) Event {
	return &InfoReportedEvent{
		BaseEvent: newEvt(EventClientReported, "report", "", args...),
	}
}

type AbnormalEvent struct {
	*BaseEvent
}

func NewAbnormalEvent(args ...string) Event {
	return &AbnormalEvent{
		BaseEvent: newEvt(EventClientAbnormal, "abnormal", "", args...),
	}
}

func optString(opt, def string) string {
	if len(opt) != 0 {
		return opt
	}

	return def
}

func newEvt(evt EventType, defName, defMsg string, args ...string) *BaseEvent {
	name := defName
	if len(args) >= 1 {
		name = optString(args[0], defName)
	}

	msg := defMsg
	if len(args) >= 2 {
		msg = optString(args[1], defMsg)
	}

	return NewBaseEvent(evt, name, msg)
}
