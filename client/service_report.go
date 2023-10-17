package client

import (
	"github.com/zourva/lwm2m/core"
	"github.com/zourva/pareto/box/meta"
)

type Reporter struct {
	*meta.StateMachine
	client   *LwM2MClient
	messager Messager
}

func NewReporter(c *LwM2MClient) *Reporter {
	return &Reporter{
		client:   c,
		messager: c.messager,
	}
}

func (r *Reporter) OnObserve(oid core.ObjectID, oiId core.InstanceID, rid core.ResourceID, riId core.InstanceID, attrs map[string]any) core.ErrorType {
	//TODO implement me
	panic("implement me")
}

func (r *Reporter) OnCancelObservation(oid core.ObjectID, oiId core.InstanceID, rid core.ResourceID, riId core.InstanceID) core.ErrorType {
	//TODO implement me
	panic("implement me")
}

func (r *Reporter) OnObserveComposite() core.ErrorType {
	//TODO implement me
	panic("implement me")
}

func (r *Reporter) OnCancelObservationComposite() core.ErrorType {
	//TODO implement me
	panic("implement me")
}

func (r *Reporter) Notify(updated *core.Value) error {
	//TODO implement me
	panic("implement me")
}

func (r *Reporter) Send(updated *core.Value) error {
	//TODO implement me
	panic("implement me")
}

var _ core.ReportingClient = &Reporter{}
