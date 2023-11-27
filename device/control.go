package device

import (
	log "github.com/sirupsen/logrus"
	"github.com/zourva/lwm2m/core"
	"github.com/zourva/lwm2m/server"
	"time"
)

type ConnState struct {
}

type ConnStatsConfig struct {
}

type VersionInfo struct {
	//
}

type VersionInfoProvider interface {
	VersionInfo(c core.RegisteredClient) VersionInfo
}

type PeriodicController struct {
	*server.DefaultEventObserver
	duration time.Duration
}

func NewPeriodicController() server.RegisteredClientObserver {
	pc := &PeriodicController{
		duration: time.Second * 5,
	}

	return pc
}

var _ server.DeviceControlService = &PeriodicController{}

func (d *PeriodicController) Registered(c core.RegisteredClient) {
	d.DefaultEventObserver.Registered(c)

	// TODO: refactor using TimingWheels
	go d.processClient(c)
}

func (d *PeriodicController) Unregistered(c core.RegisteredClient) {
	//
}

func (d *PeriodicController) processClient(c core.RegisteredClient) {
	//read device info to check version update
	//read, err := c.Read(core.OmaObjectDevice, 0, core.DeviceManufacturer, 0)
	//if err != nil {
	//	return
	//}
	//
	//// unmarshal and compare
	//c.Write(core.OmaObjectFirmwareUpdate, 0)
	//c.Write(core.OmaObjectFirmwareUpdate, 0)
	//c.Write(core.OmaObjectFirmwareUpdate, 0)
	//c.Write(core.OmaObjectFirmwareUpdate, 0)

	// observe all resources of connectivity monitoring object
	err := c.Observe(core.OmaObjectConnMonitor,
		core.NotificationAttrs{
			core.MinimumPeriod: 5,
			core.MaximumPeriod: 3600,
		},
		func(notifiedData []byte) {

		}, core.ConnectivityMonitoringIPAddresses)
	if err != nil {
		log.Errorln("observe client connectivity monitoring failed:", err)
		return
	}

	// observe all resources of connectivity statistics object
	err = c.Observe(core.OmaObjectConnStats,
		core.NotificationAttrs{
			core.MinimumPeriod: 5,
			core.MaximumPeriod: 3600,
		},
		func(notifiedData []byte) {

		})
	if err != nil {
		log.Errorln("observe client connectivity monitoring failed:", err)
		return
	}

	ticker := time.NewTicker(d.duration)
	for {
		select {
		case <-ticker.C:
		}
	}
}
