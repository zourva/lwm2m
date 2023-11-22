package server

import (
	log "github.com/sirupsen/logrus"
	"github.com/zourva/lwm2m/core"
)

type ReportingServerDelegator struct {
	server  *LwM2MServer
	service ReportingService
}

func (r *ReportingServerDelegator) Observe(oid core.ObjectID, oiId core.InstanceID, rid core.ResourceID, riId core.InstanceID, attrs map[string]any) error {
	//TODO implement me
	panic("implement me")
}

func (r *ReportingServerDelegator) CancelObservation(oid core.ObjectID, oiId core.InstanceID, rid core.ResourceID, riId core.InstanceID) error {
	//TODO implement me
	panic("implement me")
}

func (r *ReportingServerDelegator) ObserveComposite() error {
	//TODO implement me
	panic("implement me")
}

func (r *ReportingServerDelegator) CancelObservationComposite() error {
	//TODO implement me
	panic("implement me")
}

func (r *ReportingServerDelegator) OnNotify(c core.RegisteredClient, value []byte) error {
	log.Tracef("receive Notify operation data %d bytes", len(value))

	if r.service.Notify != nil {
		_, err := r.service.Notify(c, value)
		return err
	}

	return nil
}

func (r *ReportingServerDelegator) OnSend(c core.RegisteredClient, value []byte) ([]byte, error) {
	log.Tracef("receive Send operation data %d bytes", len(value))

	if r.service.Send != nil {
		return r.service.Send(c, value)
	}

	return nil, nil
}

func NewReportingServerDelegator(server *LwM2MServer, service ReportingService) core.ReportingServer {
	s := &ReportingServerDelegator{
		server:  server,
		service: service,
	}
	return s
}
