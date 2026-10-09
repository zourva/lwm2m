package core

import log "github.com/sirupsen/logrus"

// EventType defines exposed lifecycle
// event of a lwM2M client or a server.
type EventType = int

const (
	EventClientDummy               EventType = iota
	EventClientBeforeBootstrap               // issued before client starts bootstrapping
	EventClientBootstrapped                  // issued when client is bootstrapped
	EventClientBeforeRegister                //
	EventClientRegistered                    // issued when client is registered on both side
	EventClientBeforeUpdate                  //
	EventClientRegUpdated                    // issued when client registration info is updated
	EventClientBeforeUnregister              //
	EventClientUnregistered                  // issued when client is unregistered on both side
	EventClientBeforeDevInfoChange           //
	EventClientDevInfoChanged                // issued when any resource of client is operated
	EventClientBeforeObserve                 //
	EventClientObserved                      // issued when any resource of client is observed
	EventClientBeforeObserveCancel           //
	EventClientObserveCancelled              //
	EventClientBeforeReport                  //
	EventClientReported                      // issued when any resource of client is changed and reported
	EventClientAbnormal                      // issued when any error happened

	EventServerStarted
	EventServerStopped
)

const DummyMsg = ""

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

func NewBaseEvent(evt EventType, name, defMsg, userMsg string) *BaseEvent {
	msg := defMsg
	if len(userMsg) >= 0 {
		msg = userMsg
	}

	return &BaseEvent{
		evt:  evt,
		name: name,
		msg:  msg,
	}
}

// EventGenerator
// userMsg optional
type EventGenerator func(userMsg string) Event

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

// EmitEvent triggers the callback registered on evt.
func (em *EventManager) EmitEvent(evt EventType, msg string) {
	if handler, ok := em.listeners[evt]; ok {
		handler(em.createEvent(evt, msg))
	}
}

func (em *EventManager) RegisterCreator(evt EventType, gen EventGenerator) {
	em.creators[evt] = gen
}

func (em *EventManager) createEvent(evt EventType, msg string) Event {
	if creator, ok := em.creators[evt]; ok {
		return creator(msg)
	}

	log.Errorln("event is not supported in client side:", evt)
	return nil
}
