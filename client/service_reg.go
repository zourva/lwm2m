package client

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	. "github.com/zourva/lwm2m/core"
	"github.com/zourva/pareto/box/meta"
	"math"
	"sort"
	"sync/atomic"
	"time"
)

const (
	defaultLifetime          = 2592000 //30 days = 3600 * 24 * 30 seconds
	defInitRegistrationDelay = 0       // Initial Registration Delay Timer, seconds
	defCommRetryCount        = 4       // Communication Retry Count, attempts within a retry sequence
	defCommRetryTimer        = 1       // Communication Retry Timer, seconds
	defCommSeqDelayTimer     = 30      // Communication Sequence Delay Timer
	defCommSeqRetryCount     = 1       // Communication Sequence Retry Count
)

// regServerInfo defines tracking info
// to each registration server.
type regServerInfo struct {
	ServerInfo

	lifetime          uint64
	blocking          bool //pin current server
	bootstrap         bool //initiate a new bootstrap when failure
	priorityOrder     int
	initRegDelay      uint64
	commRetryLimit    uint64
	commRetryDelay    uint64
	commSeqRetryDelay uint64
	commSeqRetryLimit uint64

	retryCount     uint64 //register count within a sequence
	retrySequences uint64 //sequence count
	//exponential uint64
}

// Communication Retry Timer * 2^(Current Attempt-1)
func (r *regServerInfo) backoff() uint64 {
	//r.exponential <<= r.retryCount - 1
	//	return r.commRetryDelay * r.exponential
	delay := r.commRetryDelay * uint64(math.Pow(2, float64(r.retryCount))-1)
	log.Warnf("backoff delay %d seconds for retry %d", delay, r.retryCount)
	return delay
}

func (r *regServerInfo) reset() {
	r.retryCount = 0
	r.retrySequences = 0
	//r.exponential = 1
}

// regInfo maintains client side
// registration info.
type regInfo struct {
	name     string
	lifetime uint64 //MUST equal to Server.Lifetime resource
	mode     BindingMode
	objects  string
	//objects  []*coap.CoreResource
	//smsNumber
	//profileID

	//temporary id assigned by server
	//when registration completed
	location     string
	lifetimeLeft time.Duration
}

func (r *regInfo) setLifetime(lifetime uint64) {
	r.lifetime = lifetime
	r.lifetimeLeft = time.Duration(lifetime) * time.Second

}

// return false to invoke renew lifetime.
func (r *regInfo) decreaseLifetime(duration time.Duration) bool {
	if r.lifetimeLeft < duration {
		return false //renew lifetime
	}

	r.lifetimeLeft -= duration
	return true
}

//func (r *regInfo) buildObjectInstances(oo map[ObjectID]*InstanceManager)  {
//
//}

// Registrar implements application layer logic
// for client registration procedure at server side.
type Registrar struct {
	*meta.StateMachine[state]
	client *LwM2MClient //lwm2m context

	messager *MessagerClient
	//messagers []*MessagerClient

	regInfo   *regInfo
	servers   []*regServerInfo
	current   int
	nextDelay uint64

	fail atomic.Bool

	// update
	timer    *time.Timer
	duration time.Duration //update duration
}

func NewRegistrar(client *LwM2MClient) *Registrar {
	s := &Registrar{
		StateMachine: meta.NewStateMachine[state]("registrar", time.Second),
		client:       client,
		//messager:     client.messager,
		nextDelay: 0,
		current:   0,
		duration:  time.Second * 15,
	}

	s.timer = time.NewTimer(s.duration)
	s.timer.Stop() //stop to wait for rescheduling
	s.servers = client.getRegistrationServers()
	s.regInfo = &regInfo{
		name:     client.name,
		lifetime: defaultLifetime, //delay init to lifetime of selected server
		mode:     BindingModeUDP,
		objects:  s.buildObjectInstancesList(),
	}

	s.RegisterStates([]*meta.State[state]{
		{Name: initiating, Action: s.onInitiating},
		{Name: registering, Action: s.onRegistering},
		{Name: registered, Action: s.onRegistered},
		{Name: exiting, Action: s.onExiting},
	})

	return s
}

func (r *Registrar) singleObjectInst(oid ObjectID) ObjectInstance {
	return r.client.store.GetSingleInstance(oid)
}

// sort according to Registration Priority Order
// when multiple lwM2M servers exists.
func (r *Registrar) sortServers() {
	sort.Slice(r.servers, func(i, j int) bool {
		return r.servers[i].priorityOrder < r.servers[j].priorityOrder
	})
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

func (r *Registrar) buildObjectInstancesList() string {
	var buf bytes.Buffer

	all := r.client.store.GetInstanceManagers()
	for oid, store := range all {
		if store.Empty() {
			buf.WriteString(fmt.Sprintf("</%d>,", oid))
		} else {
			for _, inst := range store.GetAll() {
				buf.WriteString(fmt.Sprintf("</%d/%d>,", oid, inst.Id()))
			}
		}
	}

	return buf.String()
}

func (r *Registrar) enablePeriodicUpdate() {
	r.timer = time.AfterFunc(r.duration, func() {
		var params []string
		if !r.regInfo.decreaseLifetime(r.duration) {
			params = append(params, "lt")
		}

		err := r.Update(params...)
		if err != nil {
			log.Errorf("registrar update failed %v, re-register", err)
			r.client.initiateRegister()
			return
		}

		log.Tracef("registrar update successfully")

		if len(params) > 0 {
			r.regInfo.setLifetime(r.regInfo.lifetime)
		}

		r.timer.Reset(r.duration)
	})
}

func (r *Registrar) Timeout() bool {
	return r.fail.Load()
}

func (r *Registrar) onInitiating(_ any) {
	r.sortServers()

	// delay for "Initial Registration Delay Timer"
	r.addDelay(r.currentServer().initRegDelay)

	r.MoveToState(registering)
}

func (r *Registrar) onRegistering(_ any) {
	server := r.currentServer()
	r.addDelay(r.nextDelay)

	messager, err := dial(r.client, &server.ServerInfo)
	if err == nil {
		r.messager = messager

		err = r.Register()

		if err == nil {

			log.Infof("register to %s done", server.address)

			if !r.hasMoreServers() {
				r.MoveToState(registered)
				r.enablePeriodicUpdate()
				return
			}

			r.selectNextServer()
			r.nextDelay = r.currentServer().initRegDelay
			log.Infof("proceed with next server: %s", server.address)
			return
		}
	}

	log.Errorf("register to %s failed: %v", server.address, err)

	// register to current server failed
	server.retryCount++
	if server.retryCount <= server.commRetryLimit {
		// update delay and try again within current retry sequence
		r.nextDelay = server.backoff()
		return
	}

	// retry sequence exhausted
	if server.blocking {
		server.retrySequences++
		if server.retrySequences <= server.commSeqRetryLimit {
			//starts a new retry sequence to current blocked server
			r.nextDelay = server.commSeqRetryDelay
			log.Warnf("retry sequence %d exhausted, retrying to the same server %s after %d seconds",
				server.retrySequences, server.address, r.nextDelay)
			return
		} else {
			if server.bootstrap {
				// initiate a new bootstrap
			} else {
				// impl-dependent, also initiate a new bootstrap
			}

			//server.reset()
			// always initiate a new bootstrap when run out of retry sequences
			r.fail.Store(true)
			log.Infoln("retry failed, a new bootstrap needed")
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

func (r *Registrar) onRegistered(_ any) {
	//wait for client to retrieve state and give further command
}

func (r *Registrar) onExiting(_ any) {
	log.Infof("registrar exiting")
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
	// update reg info
	r.regInfo.setLifetime(r.currentServer().lifetime)

	return r.messager.Register(r.regInfo)
}

// Update requests with parameters like:
//
//	method: POST
//	uri: /{location}?lt={Lifetime}&b={binding}&Q&sms={MSISDN}
//		where location has a format of /rd/{id} and b/Q/sms are optional.
//	body: </1/0>,... which is optional.
func (r *Registrar) Update(params ...string) error {
	return r.messager.Update(r.regInfo, params...)
}

// Deregister request with parameters like:
//
//	 method: DELETE
//	 uri: /{location}
//		 where location has a format of /rd/{id}
func (r *Registrar) Deregister() error {
	return r.messager.Deregister(r.regInfo)
}

func (r *Registrar) Registered() bool {
	state := r.GetState()
	return state == registered
}

func (r *Registrar) Start() bool {
	return r.Startup()
}

func (r *Registrar) Stop() {
	r.timer.Stop()
	r.Shutdown()
}
