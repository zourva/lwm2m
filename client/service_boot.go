package client

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/zourva/lwm2m/coap"
	"github.com/zourva/lwm2m/core"
	"github.com/zourva/pareto/box/meta"
	"time"
)

type bootstrapReason = int32

const (
	bsRegRetryFailure bootstrapReason = iota
)

type bootstrapState = int32

const (
	bsNone bootstrapState = iota
	bsBootstrapping
	bsBootstrapDone
)

type Bootstrapper struct {
	client  *LwM2MClient
	machine *meta.StateMachine

	messager Messager

	//state tracking
	state bootstrapState

	bootstrapInfo *core.BootstrapInfo
}

// BootstrapRequest implements BootstrapRequest operation
//
//	method: POST
//	path: /bs?ep={Endpoint Client Name}&pct={Preferred Content Format}
//	code may be responded:
//	 2.04 Changed Bootstrap-Request is completed successfully
//	 4.00 Bad Request Unknown Endpoint Client Name
//		   Endpoint Client Name does not match with CN field of X.509 Certificates
//	 4.15 Unsupported content format The specified format is not supported
func (r *Bootstrapper) BootstrapRequest() error {
	r.setState(bsBootstrapping)

	req := r.messager.NewConRequestPlainText(coap.Post, boostrapUri)
	req.SetURIQuery("ep", r.client.Name())
	rsp, err := r.messager.SendRequest(req)
	if err != nil {
		log.Errorln("send bootstrap request failed:", err)
		return err
	}

	//check response code
	if rsp.GetMessage().Code == coap.CodeChanged {
		log.Infoln("bootstrap request accepted, progressing")
		return nil
	}

	return errors.New(rsp.GetMessage().GetCodeString())
}

// BootstrapPackRequest implements BootstrapPackRequest operation
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
func (r *Bootstrapper) BootstrapPackRequest() error {
	req := r.messager.NewConRequestPlainText(coap.Get, bootstrapPackUri)
	req.SetURIQuery("ep", r.client.Name())
	rsp, err := r.messager.SendRequest(req)
	if err != nil {
		log.Errorln("bootstrap pack request failed:", err)
		return err
	}

	//check response code
	if rsp.GetMessage().Code == coap.CodeContent {
		log.Infoln("bootstrap pack request done")
		return nil
	}

	return errors.New(rsp.GetMessage().GetCodeString())
}

func (r *Bootstrapper) OnBootstrapRead() (*core.ResourceValue, core.ErrorType) {
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

func (r *Bootstrapper) OnBootstrapWrite() core.ErrorType {
	//TODO implement me
	//codes may respond:
	//2.04 Changed "Write" operation is completed successfully
	//4.00 Bad Request The format of data to be written is different
	//4.15 Unsupported content format The specified format is not supported
	panic("implement me")
}

func (r *Bootstrapper) OnBootstrapDelete() core.ErrorType {
	//TODO implement me
	//codes may respond:
	//2.02 Deleted "Delete" operation is completed successfully
	//4.00 Bad Request Bad or unknown URI provided
	panic("implement me")
}

func (r *Bootstrapper) OnBootstrapDiscover() (*core.ResourceValue, core.ErrorType) {
	//TODO implement me
	//codes may respond:
	//2.05 Content "Discover" operation is completed successfully
	//4.00 Bad Request Undetermined error occurred
	//4.04 Not Found URI of "Discover" operation is not found
	panic("implement me")
}

func (r *Bootstrapper) OnBootstrapFinish() core.ErrorType {
	//2.04 Changed Bootstrap-Finished is completed successfully
	//4.00 Bad Request Bad URI provided
	//4.06 Not Acceptable Inconsistent loaded configuration
	r.setState(bsBootstrapDone)

	return core.ErrorNone
}

func (r *Bootstrapper) BootstrapInfo() *core.BootstrapInfo {
	//TODO implement me
	panic("implement me")
}

func (r *Bootstrapper) SecurityCredentials() *core.SecurityCredentials {
	//TODO implement me
	panic("implement me")
}

var _ core.BootstrapClient = &Bootstrapper{}

func NewBootstrapper(client *LwM2MClient) *Bootstrapper {
	s := &Bootstrapper{
		client:   client,
		machine:  meta.NewStateMachine("bootstrapper", time.Second),
		messager: client.messager,
	}

	s.bootstrapInfo = &core.BootstrapInfo{}

	s.machine.RegisterStates([]*meta.State{
		{Name: initial, Action: s.onInitial},
		{Name: bootstrapping, Action: s.onBootstrapping},
		{Name: exiting, Action: s.onExiting},
	})

	s.machine.SetStartingState(initial)
	s.machine.SetStoppingState(exiting)

	return s
}

func (r *Bootstrapper) Start() bool {
	return r.machine.Startup()
}

func (r *Bootstrapper) getState() bootstrapState {
	return r.state
}

func (r *Bootstrapper) setState(s bootstrapState) {
	r.state = s
}

func (r *Bootstrapper) bootstrapped() bool {
	return r.getState() == bsBootstrapDone
}

func (r *Bootstrapper) onInitial(args any) {
	if err := r.BootstrapRequest(); err != nil {
		log.Errorf("bootstrap failed: %v", err)
		return
	}

	r.machine.MoveToState(bootstrapping)
}

func (r *Bootstrapper) onBootstrapping(args any) {
	//wait for Bootstrap-Finish
	if r.bootstrapped() {
		log.Infof("bootstrap done")
		return
	}
}

func (r *Bootstrapper) onExiting(args any) {

}
