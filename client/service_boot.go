package client

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zourva/lwm2m/coap"
	"github.com/zourva/lwm2m/core"
	"github.com/zourva/pareto/box/meta"
	"time"
)

type bootstrapReason = int32

const (
	bootstrapReasonStartup  bootstrapReason = iota
	bootstrapReasonBootFail                 //bootstrap interface failure
	bootstrapReasonRegFail                  //register interface failure
	bootstrapReasonRptFail                  //report interface failure
	bootstrapReasonDmtFail                  //device management interface failure
)

// Bootstrapper implements the
// "Client Initiated Bootstrap" mode
// defined in Bootstrap interface.
type Bootstrapper struct {
	*meta.StateMachine[state]

	client   *LwM2MClient
	messager coap.Client

	// TODO: get from bootstrap info
	lastAttempt time.Time

	bootSeverBootInfo *core.BootstrapServerBootstrapInfo
	serverBootInfo    *core.ServerBootstrapInfo
}

// Request implements BootstrapRequest operation
//
//	method: POST
//	path: /bs?ep={Endpoint Client Name}&pct={Preferred Content Format}
//	code may be responded:
//	 2.04 Changed Bootstrap-Request is completed successfully
//	 4.00 Bad Request Unknown Endpoint Client Name
//		   Endpoint Client Name does not match with CN field of X.509 Certificates
//	 4.15 Unsupported content format The specified format is not supported
func (r *Bootstrapper) Request() error {
	//r.setState(bsBootstrapping)

	req := r.messager.NewPostRequestPlain(core.BoostrapUri, nil)
	req.AddQuery("ep", r.client.name)
	//req.AddQuery("pct", fmt.Sprintf("%d", coap.MediaTypeVndOmaLwm2mCbor))
	rsp, err := r.messager.Send(req)
	if err != nil {
		log.Errorln("send bootstrap request failed:", err)
		return err
	}

	//check response code
	if rsp.Code().Changed() {
		log.Infoln("bootstrap request accepted, progressing")
		return nil
	}

	return errors.New(rsp.Code().String())
}

// PackRequest implements BootstrapPackRequest operation
//
//	method: GET
//	format: SenML CBOR, SenML JSON, or LwM2M CBOR
//	path: /bspack?ep={Endpoint Client Name}
//	code may be responded:
//	 2.05 Content The response includes the Bootstrap-Pack.
//	 4.00 Bad Request Undetermined error occurred
//	 4.01 Unauthorized Access Right Permission Denied
//	 4.04 Not Found URI of "Bootstrap-Pack-Request" operation is not found
//	 4.05 Method Not Allowed The LwM2M Client is not allowed for "Bootstrap-Pack-Request" operation
//	 4.06 Not Acceptable The specified Content-Format is not supported
//	 5.01 Not Implemented The operation is not implemented.
func (r *Bootstrapper) PackRequest() error {
	req := r.messager.NewGetRequestPlain(core.BootstrapPackUri)
	req.AddQuery("ep", r.client.name)
	rsp, err := r.messager.Send(req)
	if err != nil {
		log.Errorln("bootstrap pack request failed:", err)
		return err
	}

	//check response code
	if rsp.Code().Content() {
		log.Infof("bootstrap pack request done with %d bytes response", rsp.Length())

		objs, err := core.ParseObjectInstancesWithJSON(r.client.store.ObjectRegistry(), string(rsp.Body()))
		if err != nil {
			log.Errorf("bootstrap pack parse object failed")
			return fmt.Errorf("parse failed")
		}

		for _, o := range objs {
			// save register server
			im := r.client.store.GetInstanceManager(o.Class().Id())
			if err = im.Upsert(o); err != nil {
				log.Errorf("bootstrap save pack response info failed, err:%v", err)
				return err
			}
		}

		// TODO: save cookies from server
		return nil
	}

	return errors.New(rsp.Code().String())
}

func (r *Bootstrapper) OnRead() (*core.ResourceField, error) {
	//TODO implement me
	//codes may respond:
	//2.05 Content "Read" operation is completed successfully
	//4.00 Bad Request Undetermined error occurred
	//4.01 Unauthorized Access Right Permission Denied
	//4.04 Not Found URI of "Read" operation is not found
	//4.05 Method Not Allowed Target is not allowed for "Read" operation
	//4.06 Not Acceptable None of the preferred Content-Formats can be returned
	panic("implement me")
}

func (r *Bootstrapper) OnWrite() error {
	//TODO implement me
	//codes may respond:
	//2.04 Changed "Write" operation is completed successfully
	//4.00 Bad Request The format of data to be written is different
	//4.15 Unsupported content format The specified format is not supported
	panic("implement me")
}

func (r *Bootstrapper) OnDelete() error {
	//TODO implement me
	//codes may respond:
	//2.02 Deleted "Delete" operation is completed successfully
	//4.00 Bad Request Bad or unknown URI provided
	panic("implement me")
}

func (r *Bootstrapper) OnDiscover() (*core.ResourceField, error) {
	//TODO implement me
	//codes may respond:
	//2.05 Content "Discover" operation is completed successfully
	//4.00 Bad Request Undetermined error occurred
	//4.04 Not Found URI of "Discover" operation is not found
	panic("implement me")
}

func (r *Bootstrapper) OnFinish() error {
	//2.04 Changed Bootstrap-Finished is completed successfully
	//4.00 Bad Request Bad URI provided
	//4.06 Not Acceptable Inconsistent loaded configuration
	//r.setState(bsBootstrapDone)

	return core.ErrorNone
}

func (r *Bootstrapper) BootstrapServerBootstrapInfo() *core.BootstrapServerBootstrapInfo {
	return r.bootSeverBootInfo
}

// SetBootstrapServerBootstrapInfo set the pre-provisioned bootstrap
// server account as depicted:
// In order for the LwM2M Client and the LwM2M Bootstrap-Server
// to establish a connection on the Bootstrap Interface, either in
// Client Initiated Bootstrap mode or in Server Initiated Bootstrap
// mode, the LwM2M Client MUST have an LwM2M Bootstrap-Server Account pre-provisioned.
func (r *Bootstrapper) SetBootstrapServerBootstrapInfo(info *core.BootstrapServerBootstrapInfo) {
	r.bootSeverBootInfo = info
}

var _ core.BootstrapClient = &Bootstrapper{}

func NewBootstrapper(client *LwM2MClient) *Bootstrapper {
	s := &Bootstrapper{
		StateMachine: meta.NewStateMachine[state]("bootstrapper", time.Second),
		client:       client,
		messager:     nil,
	}

	s.RegisterStates([]*meta.State[state]{
		{Name: initiating, Action: s.onInitiating},
		{Name: bootstrapping, Action: s.onBootstrapping},
		{Name: bootstrapped, Action: s.onBootstrapped},
		{Name: exiting, Action: s.onExiting},
	})

	return s
}

func (r *Bootstrapper) Start() bool {
	r.lastAttempt = time.Now() //.Add(60 * time.Second)
	return r.Startup()
}

func (r *Bootstrapper) Stop() {
	r.Shutdown()
}

func (r *Bootstrapper) Bootstrapped() bool {
	return r.GetState() == bootstrapped
}

func (r *Bootstrapper) Timeout() bool {
	return time.Since(r.lastAttempt) > 60*time.Second
	//return time.Now().Sub(r.lastAttempt)
}

func (r *Bootstrapper) onInitiating(_ any) {
	messager, err := coap.Dial(r.client.options.serverAddress[0], coap.WithDTLSConfig(r.client.options.dtlsConf))
	if err != nil {
		log.Errorf("bootstrap dial failed: %v", err)
		return
	}
	r.messager = messager

	log.Infof("bootstrap initiated")

	r.MoveToState(bootstrapping)
}

// NOTE: not used for packed request.
func (r *Bootstrapper) onBootstrapping(_ any) {
	//wait for Bootstrap-Finish
	//if r.Bootstrapped() {
	//	log.Infof("bootstrap done")
	//	return
	//}

	if err := r.PackRequest(); err != nil {
		log.Errorf("bootstrap failed: %v", err)
		return
	}

	log.Infof("bootstrap requested")

	r.MoveToState(bootstrapped)
}

func (r *Bootstrapper) onBootstrapped(_ any) {
	//wait for client to retrieve state and give further command
}

func (r *Bootstrapper) onExiting(_ any) {
	log.Infof("bootstraper exiting")
}
