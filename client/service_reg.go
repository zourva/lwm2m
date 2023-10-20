package client

import (
	"bytes"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zourva/lwm2m/coap"
	. "github.com/zourva/lwm2m/core"
	"github.com/zourva/pareto/box/meta"
	"math"
	"sort"
	"time"
)

type regState = int

const (
	rsUnknown regState = iota
	rsRegistering
	rsRegisterDone
	rsUpdating
	rsUpdateDone
	rsUnregistering
	rsUnregisterDone
)

const (
	defInitRegistrationDelay = 0         // Initial Registration Delay Timer, seconds
	defCommRetryCount        = 5         // Communication Retry Count, attempts within a retry sequence
	defCommRetryTimer        = 60        // Communication Retry Timer, seconds
	defCommSeqDelayTimer     = 24 * 3600 // Communication Sequence Delay Timer
	defCommSeqRetryCount     = 1         // Communication Sequence Retry Count
)

type regServerInfo struct {
	blocking          bool //pin current server
	bootstrap         bool //initiate a new bootstrap when failure
	address           string
	priorityOrder     int
	initRegDelay      uint64
	commRetryLimit    uint64
	commRetryDelay    uint64
	commSeqRetryDelay uint64
	commSeqRetryLimit uint64

	retryCount     uint64
	retrySequences uint64
	registered     bool
}

// Communication Retry Timer * 2^(Current Attempt-1)
func (r *regServerInfo) backoff() uint64 {
	return r.commRetryDelay * uint64(math.Pow(2, float64(r.retryCount))-1)
}

func (r *regServerInfo) reset() {
	r.retryCount = 0
	r.retrySequences = 0
}

// Registrar implements application layer logic
// for client registration procedure at server side.
type Registrar struct {
	machine *meta.StateMachine
	client  *LwM2MClient

	//router interface providing
	//uplink accessibility
	messager Messager

	//uri location assigned
	//when registration complete
	location string

	//state tracking
	state regState

	servers   []*regServerInfo
	current   int
	nextDelay uint64
}

func NewRegistrar(client *LwM2MClient) *Registrar {
	s := &Registrar{
		machine:   meta.NewStateMachine("registrar", time.Second),
		client:    client,
		messager:  client.messager,
		location:  "",
		state:     rsUnknown,
		nextDelay: 0,
		current:   0,
	}

	s.servers = []*regServerInfo{
		{
			blocking:          true,
			bootstrap:         true,
			address:           "127.0.0.1:5683",
			priorityOrder:     1,
			initRegDelay:      defInitRegistrationDelay,
			commRetryLimit:    defCommRetryCount,
			commRetryDelay:    defCommRetryTimer,
			commSeqRetryDelay: defCommSeqDelayTimer,
			commSeqRetryLimit: defCommSeqRetryCount,
		},
	}

	s.machine.RegisterStates([]*meta.State{
		{Name: initial, Action: s.onInitial},
		{Name: registering, Action: s.onRegistering},
		{Name: monitoring, Action: s.onMonitoring},
		{Name: updating, Action: s.onUpdating},
		{Name: unregistering, Action: s.onUnregistering},
		{Name: exiting, Action: s.onExiting},
	})

	s.machine.SetStartingState(initial)
	s.machine.SetStoppingState(exiting)

	return s
}

func (r *Registrar) getState() regState {
	return r.state
}

func (r *Registrar) setState(s regState) {
	r.state = s
}

func (r *Registrar) singleObjectInst(oid ObjectID) ObjectInstance {
	return r.client.store.GetSingleInstance(oid)
}

func (r *Registrar) resMgr(oid ObjectID) *ResInstManager {
	return r.singleObjectInst(oid).ResInstManager()
}

// sort according to Registration Priority Order
// when multiple lwM2M servers exists.
func (r *Registrar) sortServers() {
	sort.Slice(r.servers, func(i, j int) bool {
		return r.servers[i].priorityOrder < r.servers[j].priorityOrder
	})
}

func (r *Registrar) initiateBootstrap() {
	r.machine.Pause()
	r.client.RequestBootstrap(bsRegRetryFailure)
}

func (r *Registrar) currentServer() *regServerInfo {
	return r.servers[r.current]
}

func (r *Registrar) hasMoreServers() bool {
	return r.current < len(r.servers)-1
}

func (r *Registrar) selectNextServer() {
	r.current += 1
	r.current %= len(r.servers)
}

func (r *Registrar) addDelay(delaySec uint64) {
	time.Sleep(time.Duration(delaySec) * time.Second)
}

func (r *Registrar) registrationInfoChanged() bool {
	return false
}

func (r *Registrar) onInitial(args any) {
	r.sortServers()

	// delay for "Initial Registration Delay Timer"
	r.addDelay(r.currentServer().initRegDelay)

	r.machine.MoveToState(registering)
}

func (r *Registrar) onRegistering(args any) {
	server := r.currentServer()

	r.addDelay(r.nextDelay)

	if err := r.Register(); err != nil {
		log.Errorf("register to %s failed: %v", server.address, err)
		return
	}

	if r.registered() {
		log.Infof("register to %s done", server.address)

		if !r.hasMoreServers() {
			r.machine.MoveToState(monitoring)
			return
		}

		r.selectNextServer()
		r.nextDelay = r.currentServer().initRegDelay
		log.Infof("proceed with next server: %s", server.address)
		return
	}

	// register to current server failed
	server.retryCount++
	if server.retryCount <= server.commRetryLimit {
		// update delay and try again within current retry sequence
		r.nextDelay = server.backoff()
		return
	}

	// retry sequence exhausted
	if server.blocking {
		if server.retrySequences <= server.commSeqRetryLimit {
			//starts a new retry sequence to current blocked server
			r.nextDelay = server.commSeqRetryDelay
			return
		} else {
			if server.bootstrap {
				// initiate a new bootstrap
			} else {
				// impl-dependent, also initiate a new bootstrap
			}

			//server.reset()
			// always initiate a new bootstrap when run out of retry sequences
			r.initiateBootstrap()
			log.Infoln("retry failed, a new bootstrap is requested")
			return
		}
	} else {
		// impl-dependent non-blocking retry, proceeds with the next server
		r.selectNextServer()
		next := r.currentServer()
		next.reset()
		r.nextDelay = next.initRegDelay
		log.Infoln("retry nonblocking registration to next server", next)
	}
}

func (r *Registrar) onMonitoring(args any) {
	//TODO: collect changed registration info
	if r.registrationInfoChanged() {
		r.machine.MoveToState(updating)
	}
}
func (r *Registrar) onUpdating(args any) {
	if r.updated() {
		r.machine.MoveToState(monitoring)
		log.Infoln("registration info updated")
	}
}

func (r *Registrar) onUnregistering(args any) {
	if r.unregistered() {
		r.machine.MoveToState(exiting)
		log.Infoln("client unregistered")
	}
}

func (r *Registrar) onExiting(args any) {
	//clear thing
}

func (r *Registrar) registered() bool {
	return r.getState() == rsRegisterDone
}

func (r *Registrar) unregistered() bool {
	return r.getState() == rsUnregisterDone
}

func (r *Registrar) updated() bool {
	return r.getState() == rsUpdateDone
}

// Register encapsulates request payload containing objects
// and instances and requests according to the following:
//
//	method: POST
//	uri: /rd?ep={Endpoint Client Name}&lt={Lifetime}
//	        &lwm2m={version}&b={binding}&Q&sms={MSISDN}&pid={ProfileID}
//	   b/Q/sms/pid are optional.
//	body: </1/0>,... which is optional.
func (r *Registrar) Register() error {
	r.setState(rsRegistering)

	// send request
	req := r.messager.NewConRequestPlainText(coap.Post, registerUri)
	req.SetURIQuery("ep", r.client.Name())
	req.SetURIQuery("lt", defaultLifetime)
	req.SetURIQuery("lwm2m", lwM2MVersion)
	req.SetURIQuery("b", BindingModeUDP)
	req.SetStringPayload(r.buildObjectInstancesList())
	rsp, err := r.messager.SendRequest(req)
	if err != nil {
		log.Errorln("send register request failed:", err)
		return err
	}

	// check response code
	if rsp.GetMessage().Code == coap.CodeCreated {
		r.setState(rsRegisterDone)

		// save location for update or de-register operation
		r.location = rsp.GetMessage().GetLocationPath()
		log.Infoln("register done with assigned location:", r.location)
		return nil
	}

	return errors.New(rsp.GetMessage().GetCodeString())
}

// Deregister request with parameters like:
//
//	 method: DELETE
//	 uri: /{location}
//		 where location has a format of /rd/{id}
func (r *Registrar) Deregister() error {
	r.setState(rsUnregistering)

	uri := registerUri + fmt.Sprintf("/%s", r.location)
	req := r.messager.NewConRequestPlainText(coap.Delete, uri)
	rsp, err := r.messager.SendRequest(req)
	if err != nil {
		log.Errorln("send de-register request failed:", err)
		return err
	}

	// check response code
	if rsp.GetMessage().Code == coap.CodeDeleted {
		log.Infoln("deregister done on", uri)
		r.setState(rsUnregisterDone)
		return nil
	}

	return errors.New(rsp.GetMessage().GetCodeString())
}

// Update requests with parameters like:
//
//	method: POST
//	uri: /{location}?lt={Lifetime}&b={binding}&Q&sms={MSISDN}
//		where location has a format of /rd/{id} and b/Q/sms are optional.
//	body: </1/0>,... which is optional.
func (r *Registrar) Update(params ...any) error {
	r.setState(rsUpdating)

	uri := registerUri + fmt.Sprintf("/%s", r.location)
	req := r.messager.NewConRequestPlainText(coap.Post, uri)
	req.SetStringPayload(r.buildObjectInstancesList())
	rsp, err := r.messager.SendRequest(req)
	if err != nil {
		log.Errorln("send update request failed:", err)
		return err
	}

	// check response code
	if rsp.GetMessage().Code == coap.CodeChanged {
		log.Infoln("update done on", uri)
		r.setState(rsUpdateDone)
		return nil
	}

	return nil
}

func (r *Registrar) Start() bool {
	return r.machine.Startup()
}

func (r *Registrar) Stop() {
	r.machine.Shutdown()
}

func (r *Registrar) buildObjectInstancesList() string {
	var buf bytes.Buffer

	all := r.client.store.GetAllInstances()
	for oid, store := range all {
		if store.Empty() {
			buf.WriteString(fmt.Sprintf("</%d>,", oid))
		} else {
			for _, inst := range store.GetAll() {
				buf.WriteString(fmt.Sprintf("</%d/%d>,", oid, inst.InstanceID()))
			}
		}
	}

	return buf.String()
}
