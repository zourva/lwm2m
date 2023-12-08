package client

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zourva/lwm2m/core"
	"sync/atomic"
	"time"
)

// Reporter implements information reporting
// functionalities of the lwm2m protocol.
// Unlike other interface implementations,
// reporter has no standalone loop embedded.
type Reporter struct {
	client *LwM2MClient
	//messager coap.Client

	observer *Observer

	// statistics
	failure atomic.Int32
}

func NewReporter(c *LwM2MClient) *Reporter {
	r := &Reporter{
		client: c,
		//messager: c.messager(),
		observer: newObserver(),
	}

	r.failure.Store(0)

	return r
}

func (r *Reporter) messager() *MessagerClient {
	return r.client.messager()
}

func (r *Reporter) OnObserve(observationId string, attrs core.NotificationAttrs) error {
	r.observer.add(observationId, attrs, nil)
	return core.ErrorNone
}

func (r *Reporter) OnCancelObservation(observationId string) error {
	r.observer.delete(observationId)
	return core.ErrorNone
}

func (r *Reporter) OnObserveComposite() error {
	//TODO implement me
	panic("implement me")
}

func (r *Reporter) OnCancelObservationComposite() error {
	//TODO implement me
	panic("implement me")
}

func (r *Reporter) Notify(observationId string, value []byte) error {
	observation := r.observer.get(observationId)
	if observation == nil {
		log.Traceln("observation is not found for", observationId)
		// TODO: stop notification
		//return errors.New("invalid observation id")
	}

	return r.messager().Notify(observationId, value)
}

func (r *Reporter) Send(value []byte) ([]byte, error) {
	req := r.messager().NewPostRequestOpaque(core.SendReportUri, value)
	rsp, err := r.messager().Send(req)
	if err != nil {
		r.incrementFailCounter()
		log.Errorf("send opaque request failed: %v ", err)
		return nil, err
	}

	// check response code
	if rsp.Code().Changed() {
		r.resetFailCounter()
		log.Traceln("send opaque request done")
		return rsp.Body(), nil
	}

	return nil, errors.New(rsp.Code().String())
}

func (r *Reporter) FailureCounter() int32 {
	return r.failure.Load()
}

func (r *Reporter) incrementFailCounter() {
	r.failure.Add(1)
}

func (r *Reporter) resetFailCounter() {
	r.failure.Store(0)
}

var _ core.ReportingClient = &Reporter{}

type Observation struct {
	oid   core.ObjectID
	oiId  core.InstanceID
	rid   core.ResourceID
	riId  core.InstanceID
	attrs core.NotificationAttrs

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

func (o *Observer) add(key string, attrs core.NotificationAttrs, token []byte) {
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
