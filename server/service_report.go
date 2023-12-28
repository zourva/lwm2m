package server

import (
	"github.com/zourva/lwm2m/core"
)

type ReportingServerDelegator struct {
	server *LwM2MServer
}

func (r *ReportingServerDelegator) OnNotify(c core.RegisteredClient, value []byte) error {
	//log.Tracef("receive Notify operation data %d bytes", len(value))

	if r.server.reportService.Notify != nil {
		_, err := r.server.reportService.Notify(c, value)
		return err
	}

	return nil
}

func (r *ReportingServerDelegator) OnSend(c core.RegisteredClient, value []byte) ([]byte, error) {
	//log.Tracef("receive Send operation data %d bytes", len(value))

	if r.server.reportService.Send != nil {
		return r.server.reportService.Send(c, value)
	}

	return nil, nil
}

func NewReportingServerDelegator(server *LwM2MServer) *ReportingServerDelegator {
	s := &ReportingServerDelegator{
		server: server,
	}
	return s
}
