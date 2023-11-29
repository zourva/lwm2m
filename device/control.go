package device

import (
	log "github.com/sirupsen/logrus"
	"github.com/zourva/lwm2m/core"
	"github.com/zourva/lwm2m/server"
	"time"
)

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
			core.MinimumPeriod: "30",
			core.MaximumPeriod: "3600",
		},
		func(notifiedData []byte) {
			//pack := senml.Decode(notifiedData, senml.JSON)
			log.Infof("connectivity state changed for client %s: %v", c.Name(), string(notifiedData))
		})
	if err != nil {
		log.Errorln("observe client connectivity monitoring failed:", err)
		return
	}

	// observe all resources of connectivity statistics object
	err = c.Observe(core.OmaObjectConnStats,
		core.NotificationAttrs{
			core.MinimumPeriod: "30",
			core.MaximumPeriod: "3600",
		},
		func(notifiedData []byte) {
			log.Infof("connectivity stats changed for client %s: %v", c.Name(), string(notifiedData))
		})
	if err != nil {
		log.Errorln("observe client connectivity statistics failed:", err)
		return
	}

	ticker := time.NewTicker(d.duration)
	for {
		select {
		case <-ticker.C:
		}
	}
}
