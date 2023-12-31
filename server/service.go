package server

import (
	log "github.com/sirupsen/logrus"
	"github.com/zourva/lwm2m/core"
)

// RegisteredClientObserver defines lifecycle event
// observers/callbacks for a LwM2M client.
type RegisteredClientObserver interface {
	// Bootstrapped invoked after client is bootstrapped
	Bootstrapped(epName string)

	// Registered invoked after client is registered
	Registered(c core.RegisteredClient)

	// Updated invoked after client registration info is updated
	Updated(c core.RegisteredClient)

	// Unregistered invoked after client unregistered
	Unregistered(c core.RegisteredClient)

	// DeviceOperated invoked after any resource of and object is operated
	//DeviceOperated(c core.RegisteredClient, objs []core.ObjectInstance)
}

// DefaultEventObserver implements RegisteredClientObserver
// and provides a dummy operation for each event.
type DefaultEventObserver struct {
}

func (d *DefaultEventObserver) Bootstrapped(epName string) {
	log.Infof("client %s is bootstrappd", epName)
	return
}

func (d *DefaultEventObserver) Registered(c core.RegisteredClient) {
	log.Infof("client %s is registered", c.Name())
	return
}

func (d *DefaultEventObserver) Updated(c core.RegisteredClient) {
	log.Infof("registration inof of client %s is updated", c.Name())
	return
}

func (d *DefaultEventObserver) Unregistered(c core.RegisteredClient) {
	log.Infof("client %s is deregistered", c.Name())
	return
}

func (d *DefaultEventObserver) DeviceOperated(c core.RegisteredClient, objs []core.ObjectInstance) {
	log.Infof("client %s objects and resources %v is manipulated", c.Name(), objs)
	return
}

var _ RegisteredClientObserver = &DefaultEventObserver{}

func NewDefaultEventObserver() *DefaultEventObserver {
	return &DefaultEventObserver{}
}
