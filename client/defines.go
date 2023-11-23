package client

const (
	lwM2MVersion      = "1.1"
	defaultServerAddr = "127.0.0.1:5683"
	defaultLocalAddr  = ":0"
)

// client state
type state = int32

const (
	initiating state = iota
	bootstrapping
	bootstrapped
	//networking
	registering
	registered
	servicing
	updating
	updated
	unregistering
	unregistered
	reporting
	exiting
)

var stateNameMapping = map[state]string{
	initiating:    "initiating",
	bootstrapping: "bootstrapping",
	bootstrapped:  "bootstrapped",
	//networking:    "networking",
	servicing:     "servicing",   //reporting
	registering:   "registering", //long duration state
	registered:    "registered",  //should enable update sub-procedure
	updating:      "updating",    //long duration state
	updated:       "updated",     //transient state
	unregistering: "unregistering",
	unregistered:  "unregistered",
	reporting:     "reporting",
	exiting:       "exiting",
}

func stateName(s state) string {
	return stateNameMapping[s]
}
