package server

import "github.com/zourva/lwm2m/core"

type ReportingService struct {
}

func (r *ReportingService) Observe(oid core.ObjectID, oiId core.InstanceID, rid core.ResourceID, riId core.InstanceID, attrs map[string]any) error {
	//TODO implement me
	panic("implement me")
}

func (r *ReportingService) CancelObservation(oid core.ObjectID, oiId core.InstanceID, rid core.ResourceID, riId core.InstanceID) error {
	//TODO implement me
	panic("implement me")
}

func (r *ReportingService) ObserveComposite() error {
	//TODO implement me
	panic("implement me")
}

func (r *ReportingService) CancelObservationComposite() error {
	//TODO implement me
	panic("implement me")
}

func (r *ReportingService) OnNotify(value []byte) error {
	//TODO implement me
	panic("implement me")
}

func (r *ReportingService) OnSend(value []byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func NewInfoReportingService(server *LwM2MServer) core.ReportingServer {
	s := &ReportingService{}
	return s
}
