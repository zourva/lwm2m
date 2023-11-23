package client

import (
	"bytes"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zourva/lwm2m/coap"
	. "github.com/zourva/lwm2m/core"
	"github.com/zourva/lwm2m/utils"
	"github.com/zourva/pareto/box/meta"
	"math"
	"sort"
	"sync/atomic"
	"time"
)

const (
	defaultLifetime          = 2592000   //30 days = 3600 * 24 * 30 seconds
	defInitRegistrationDelay = 0         // Initial Registration Delay Timer, seconds
	defCommRetryCount        = 5         // Communication Retry Count, attempts within a retry sequence
	defCommRetryTimer        = 60        // Communication Retry Timer, seconds
	defCommSeqDelayTimer     = 24 * 3600 // Communication Sequence Delay Timer
	defCommSeqRetryCount     = 1         // Communication Sequence Retry Count
)

// regServerInfo defines tracking info
// to each registration server.
type regServerInfo struct {
	lifetime          uint64
	blocking          bool //pin current server
	bootstrap         bool //initiate a new bootstrap when failure
	address           string
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
	return r.commRetryDelay * uint64(math.Pow(2, float64(r.retryCount))-1)
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
	location string

	lifetimeUpdate atomic.Bool
	objectsChanged atomic.Bool
}

//func (r *regInfo) buildObjectInstances(oo map[ObjectID]*InstanceManager)  {
//
//}

// Registrar implements application layer logic
// for client registration procedure at server side.
type Registrar struct {
	*meta.StateMachine[state]
	client *LwM2MClient //lwm2m context

	messager Messager

	regInfo   *regInfo
	servers   []*regServerInfo
	current   int
	nextDelay uint64

	fail atomic.Bool

	// update
	timer        *time.Timer
	duration     time.Duration //update duration
	lifetimeLeft time.Duration //lifetime
}

func NewRegistrar(client *LwM2MClient) *Registrar {
	s := &Registrar{
		StateMachine: meta.NewStateMachine[state]("registrar", time.Second),
		client:       client,
		messager:     client.messager,
		nextDelay:    0,
		current:      0,
		duration:     time.Second * 15,
	}

	s.timer = time.NewTimer(s.duration)
	s.timer.Stop() //stop to wait for rescheduling
	s.servers = client.getRegistrationServers()
	s.regInfo = &regInfo{
		name:     client.name,
		lifetime: defaultLifetime,
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
	r.lifetimeLeft = time.Duration(r.regInfo.lifetime) * time.Second
	time.AfterFunc(r.duration, func() {
		var params []string
		if r.lifetimeLeft < r.duration {
			//re-rent lifetime
			r.lifetimeLeft = time.Duration(r.regInfo.lifetime) * time.Second
			params = append(params, "lt")
		} else {
			r.lifetimeLeft -= r.duration
		}

		err := r.Update(params...)
		if err != nil {
			log.Errorf("registrar update failed %v, re-register", err)
			r.client.initiateRegister()
			return
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

	err := r.Register()

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
	r.regInfo.lifetime = r.currentServer().lifetime

	// send request
	req := r.messager.NewConRequestPlainText(coap.Post, RegisterUri)
	req.SetUriQuery("ep", r.regInfo.name)
	req.SetUriQuery("lt", utils.IntToStr(r.regInfo.lifetime))
	req.SetUriQuery("lwm2m", lwM2MVersion)
	req.SetUriQuery("b", r.regInfo.mode)
	req.SetStringPayload(r.regInfo.objects)
	rsp, err := r.messager.Send(req)
	if err != nil {
		log.Errorln("send register request failed:", err)
		return err
	}

	// check response code
	if rsp.Message().Code == coap.CodeCreated {
		// save location for update or de-register operation
		r.regInfo.location = rsp.Message().GetLocationPath()
		log.Infoln("register done with assigned location:", r.regInfo.location)
		return nil
	}

	log.Errorln("register request failed:", coap.CodeString(rsp.Message().Code))

	return errors.New(rsp.Message().GetCodeString())
}

// Deregister request with parameters like:
//
//	 method: DELETE
//	 uri: /{location}
//		 where location has a format of /rd/{id}
func (r *Registrar) Deregister() error {
	uri := RegisterUri + fmt.Sprintf("/%s", r.regInfo.location)
	req := r.messager.NewConRequestPlainText(coap.Delete, uri)
	rsp, err := r.messager.Send(req)
	if err != nil {
		log.Errorln("send de-register request failed:", err)
		return err
	}

	// check response code
	if rsp.Message().Code == coap.CodeDeleted {
		log.Infoln("deregister done on", uri)
		return nil
	}

	log.Errorln("de-register request failed:", coap.CodeString(rsp.Message().Code))

	return errors.New(coap.CodeString(rsp.Message().Code))
}

// Update requests with parameters like:
//
//	method: POST
//	uri: /{location}?lt={Lifetime}&b={binding}&Q&sms={MSISDN}
//		where location has a format of /rd/{id} and b/Q/sms are optional.
//	body: </1/0>,... which is optional.
func (r *Registrar) Update(params ...string) error {
	uri := RegisterUri + fmt.Sprintf("/%s", r.regInfo.location)
	req := r.messager.NewConRequestPlainText(coap.Post, uri)

	for _, param := range params {
		if param == "lt" {
			req.SetUriQuery("lt", utils.IntToStr(r.regInfo.lifetime))
		} else if param == "objlink" {
			req.SetStringPayload(r.regInfo.objects)
		}
	}

	rsp, err := r.messager.Send(req)
	if err != nil {
		log.Errorln("send update request failed:", err)
		return err
	}

	// check response code
	if rsp.Message().Code == coap.CodeChanged {
		log.Infoln("update done on", uri)
		return nil
	}

	log.Errorln("update request failed:", coap.CodeString(rsp.Message().Code))

	return errors.New(coap.CodeString(rsp.Message().Code))
}

func (r *Registrar) Registered() bool {
	return r.GetState() == registered
}

func (r *Registrar) Start() bool {
	return r.Startup()
}

func (r *Registrar) Stop() {
	r.timer.Stop()
	r.Shutdown()
}
