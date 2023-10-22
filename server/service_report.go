package server

import (
	log "github.com/sirupsen/logrus"
	"github.com/zourva/lwm2m/core"
)

type ReportingService struct {
	server *LwM2MServer
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

func (r *ReportingService) OnNotify(c core.RegisteredClient, value []byte) error {
	log.Debugf("receive Notify operation data %d bytes", len(value))

	if r.server.onSent != nil {
		return r.server.onNotified(c, value)
	}

	return nil
}

func (r *ReportingService) OnSend(c core.RegisteredClient, value []byte) ([]byte, error) {
	log.Debugf("receive Send operation data %d bytes", len(value))

	if r.server.onSent != nil {
		return r.server.onSent(c, value)
	}

	return nil, nil
}

func NewInfoReportingService(server *LwM2MServer) core.ReportingServer {
	s := &ReportingService{
		server: server,
	}
	return s
}
