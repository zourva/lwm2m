package client

import (
	"fmt"
	"github.com/zourva/lwm2m/core"
	"github.com/zourva/pareto/box/meta"
	"time"
)

type InfoReporter struct {
	*meta.StateMachine
	client   *LwM2MClient
	messager Messager

	observer *Observer
}

func NewReporter(c *LwM2MClient) *InfoReporter {
	return &InfoReporter{
		client:   c,
		messager: c.messager,
	}
}

func (r *InfoReporter) OnObserve(oid core.ObjectID, oiId core.InstanceID, rid core.ResourceID, riId core.InstanceID, attrs map[string]any) core.ErrorType {
	//TODO implement me
	panic("implement me")
}

func (r *InfoReporter) OnCancelObservation(oid core.ObjectID, oiId core.InstanceID, rid core.ResourceID, riId core.InstanceID) core.ErrorType {
	//TODO implement me
	panic("implement me")
}

func (r *InfoReporter) OnObserveComposite() core.ErrorType {
	//TODO implement me
	panic("implement me")
}

func (r *InfoReporter) OnCancelObservationComposite() core.ErrorType {
	//TODO implement me
	panic("implement me")
}

func (r *InfoReporter) Notify(updated *core.Value) error {
	//TODO implement me
	panic("implement me")
}

func (r *InfoReporter) Send(value []byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

var _ core.ReportingClient = &InfoReporter{}

type Observation struct {
	oid   core.ObjectID
	oiId  core.InstanceID
	rid   core.ResourceID
	riId  core.InstanceID
	attrs map[string]any

	key string

	addTime int64 //time when this observation is added, in seconds
}

// Observer manages persistent observations across registrations.
type Observer struct {
	observations map[string]*Observation
}

func (o *Observer) makeKey(oid core.ObjectID, oiId core.InstanceID, rid core.ResourceID, riId core.InstanceID) string {
	return fmt.Sprintf("/%x/%x/%x/%x", oid, oiId, rid, riId)
}

func (o *Observer) Add(oid core.ObjectID, oiId core.InstanceID, rid core.ResourceID, riId core.InstanceID, attrs map[string]any) {
	key := o.makeKey(oid, oiId, rid, riId)
	o.observations[key] = &Observation{
		oid:     oid,
		oiId:    oiId,
		rid:     rid,
		riId:    riId,
		attrs:   attrs,
		key:     key,
		addTime: time.Now().Unix(),
	}
}
