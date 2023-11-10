package client

const (
	defaultLocalAddr  = ":0"
	defaultServerAddr = "127.0.0.1:5683"
)

const (
	initial = "initial"

	bootstrapping = "bootstrapping"
	bootstrapped  = "bootstrapped"
	networking    = "networking"
	servicing     = "servicing"

	registering   = "registering" //long duration state
	registered    = "registered"  //should enable update sub-procedure
	monitoring    = "monitoring"
	updating      = "updating" //long duration state
	updated       = "updated"  //transient state
	unregistering = "unregistering"
	unregistered  = "unregistered"

	exiting = "exiting"
)

const (
	defaultLifetime = "2592000" //30 days = 3600 * 24 * 30 seconds
	lwM2MVersion    = "1.1"
)
