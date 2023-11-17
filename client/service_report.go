package client

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zourva/lwm2m/coap"
	"github.com/zourva/lwm2m/core"
	"github.com/zourva/pareto/box/meta"
	"time"
)

type InfoReporter struct {
	*meta.StateMachine
	client   *LwM2MClient
	messager core.Messager

	observer *Observer
}

func NewReporter(c *LwM2MClient) *InfoReporter {
	return &InfoReporter{
		client:   c,
		messager: c.messager,
		observer: newObserver(),
	}
}

func (r *InfoReporter) OnObserve(observationId string, attrs map[string]any) error {
	r.observer.add(observationId, attrs, nil)
	return core.ErrorNone
}

func (r *InfoReporter) OnCancelObservation(observationId string) error {
	r.observer.delete(observationId)
	return core.ErrorNone
}

func (r *InfoReporter) OnObserveComposite() error {
	//TODO implement me
	panic("implement me")
}

func (r *InfoReporter) OnCancelObservationComposite() error {
	//TODO implement me
	panic("implement me")
}

func (r *InfoReporter) Notify(observationId string, value []byte) error {
	observation := r.observer.get(observationId)
	if observation == nil {
		log.Traceln("observation is not found for", observationId)
		// TODO: stop notification
		//return errors.New("invalid observation id")
	}

	return r.messager.Notify(observationId, value)
}

func (r *InfoReporter) Send(value []byte) ([]byte, error) {
	req := r.messager.NewConRequestOpaque(coap.Post, core.SendReportUri, value)
	rsp, err := r.messager.Send(req)
	if err != nil {
		r.client.incRegRestartCounter()
		log.Errorln("send opaque request failed: %v, ", err)
		return nil, err
	}

	r.client.resetRegRestartCounter()

	// check response code
	if rsp.GetMessage().Code == coap.CodeChanged {
		log.Traceln("send opaque request done")
		return rsp.GetPayload(), nil
	}

	return nil, errors.New(coap.CoapCodeToString(rsp.GetMessage().Code))
}

var _ core.ReportingClient = &InfoReporter{}

type Observation struct {
	oid   core.ObjectID
	oiId  core.InstanceID
	rid   core.ResourceID
	riId  core.InstanceID
	attrs map[string]any

	key   string //joined path
	token []byte

	addTime int64 //time when this observation is added, in seconds
}

// Observer manages observations across registrations.
type Observer struct {
	observations map[string]*Observation
}

func (o *Observer) makeKey(oid core.ObjectID, oiId core.InstanceID, rid core.ResourceID, riId core.InstanceID) string {
	return fmt.Sprintf("/%x/%x/%x/%x", oid, oiId, rid, riId)
}

func (o *Observer) get(key string) *Observation {
	if observation, ok := o.observations[key]; ok {
		return observation
	}

	return nil
}

func (o *Observer) add(key string, attrs map[string]any, token []byte) {
	//key := o.makeKey(oid, oiId, rid, riId)
	o.observations[key] = &Observation{
		//oid:     oid,
		//oiId:    oiId,
		//rid:     rid,
		//riId:    riId,
		attrs:   attrs,
		key:     key,
		token:   token,
		addTime: time.Now().Unix(),
	}
}

func (o *Observer) delete(key string /*token []byte*/) {
	delete(o.observations, key)
	//if observation, ok := o.observations[key]; ok {
	//	if bytes.Equal(observation.token, token) {
	//		delete(o.observations, key)
	//	} else {
	//		//forgotten by peer, reset by replace
	//		log.Warnln("observation deletion ignored since token does not match")
	//	}
	//}
}

func newObserver() *Observer {
	observer := &Observer{
		observations: make(map[string]*Observation),
	}
	return observer
}
