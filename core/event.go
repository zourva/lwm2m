package core

import log "github.com/sirupsen/logrus"

// EventType defines exposed lifecycle
// event of a lwM2M client or a server.
type EventType = int

const (
	EventClientBootstrapped   EventType = iota // issued when client is bootstrapped
	EventClientRegistered                      // issued when client is registered
	EventClientRegUpdated                      // issued when client registration info is updated
	EventClientUnregistered                    // issued when client is unregistered
	EventClientDevInfoChanged                  // issued when any resource of client is operated
	EventClientObserved                        // issued when any resource of client is observed
	EventClientReported                        // issued when any resource of client is changed and reported
	EventClientAbnormal                        // issued when any error happened

	EventServerStarted
	EventServerStopped
)

type Event interface {
	// Name returns name of this event.
	Name() string

	// Type returns type of this event.
	Type() EventType

	// Message provides more details about the event.
	Message() string
}

type EventHandler = func(event Event)

// BaseEvent implements base event.
type BaseEvent struct {
	name string
	evt  EventType
	msg  string
}

func (e *BaseEvent) Name() string {
	return e.name
}

func (e *BaseEvent) Type() EventType {
	return e.evt
}

func (e *BaseEvent) Message() string {
	return e.msg
}

func optString(opt, def string) string {
	if len(opt) != 0 {
		return opt
	}

	return def
}

func NewBaseEvent(evt EventType, defName, defMsg string, args ...string) *BaseEvent {
	name := defName
	if len(args) >= 1 {
		name = optString(args[0], defName)
	}

	msg := defMsg
	if len(args) >= 2 {
		msg = optString(args[1], defMsg)
	}

	return &BaseEvent{
		evt:  evt,
		name: name,
		msg:  msg,
	}
}

// EventGenerator
//
//	name: args[0]
//	msg: args[1]
type EventGenerator func(args ...string) Event

type EventManager struct {
	listeners map[EventType]EventHandler
	creators  map[EventType]EventGenerator
}

func NewEventManager() *EventManager {
	em := &EventManager{
		listeners: make(map[EventType]EventHandler),
		creators:  make(map[EventType]EventGenerator),
	}

	return em
}

func (em *EventManager) AddListener(et EventType, h EventHandler) {
	em.listeners[et] = h
}

func (em *EventManager) EmitEvent(evt EventType) {
	if handler, ok := em.listeners[evt]; ok {
		handler(em.createEvent(evt))
	}
}

func (em *EventManager) RegisterCreator(evt EventType, gen EventGenerator) {
	em.creators[evt] = gen
}

func (em *EventManager) createEvent(evt EventType) Event {
	if creator, ok := em.creators[evt]; ok {
		return creator()
	}

	log.Errorln("event is not supported in client side:", evt)
	return nil
}
