package client

import . "github.com/zourva/lwm2m/core"

type BootstrappedEvent struct {
	*BaseEvent
}

func NewBootstrappedEvent(args ...string) Event {
	return &BootstrappedEvent{
		BaseEvent: NewBaseEvent(EventClientBootstrapped, "bootstrapped", "", args...),
	}
}

type RegisteredEvent struct {
	*BaseEvent
}

func NewRegisteredEvent(args ...string) Event {
	return &RegisteredEvent{
		BaseEvent: NewBaseEvent(EventClientRegistered, "registered", "", args...),
	}
}

type RegUpdatedEvent struct {
	*BaseEvent
}

func NewRegUpdatedEvent(args ...string) Event {
	return &RegUpdatedEvent{
		BaseEvent: NewBaseEvent(EventClientRegUpdated, "registration updated", "", args...),
	}
}

type UnregisteredEvent struct {
	*BaseEvent
}

func NewUnregisteredEvent(args ...string) Event {
	return &UnregisteredEvent{
		BaseEvent: NewBaseEvent(EventClientUnregistered, "unregistered", "", args...),
	}
}

type DeviceChangedEvent struct {
	*BaseEvent
}

func NewDeviceChangedEvent(args ...string) Event {
	return &DeviceChangedEvent{
		BaseEvent: NewBaseEvent(EventClientDevInfoChanged, "device control", "", args...),
	}
}

type InfoObservedEvent struct {
	*BaseEvent
}

func NewInfoObservedEvent(args ...string) Event {
	return &InfoObservedEvent{
		BaseEvent: NewBaseEvent(EventClientObserved, "observe", "", args...),
	}
}

type InfoReportedEvent struct {
	*BaseEvent
}

func NewInfoReportedEvent(args ...string) Event {
	return &InfoReportedEvent{
		BaseEvent: NewBaseEvent(EventClientReported, "report", "", args...),
	}
}

type AbnormalEvent struct {
	*BaseEvent
}

func NewAbnormalEvent(args ...string) Event {
	return &AbnormalEvent{
		BaseEvent: NewBaseEvent(EventClientAbnormal, "abnormal", "", args...),
	}
}
