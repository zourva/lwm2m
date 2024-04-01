package server

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zourva/lwm2m/coap"
	. "github.com/zourva/lwm2m/core"
	"github.com/zourva/pareto/endec/senml"
	"strconv"
	"time"
)

// MessagerServer encapsulates and hide
// transport layer details.
//
// Put them together here to make it
// easier when replacing COAP layer later.
type MessagerServer struct {
	coap.Server
	lwM2MServer *LwM2MServer //server context

	// transport layer
	network string
	address string
}

func NewMessager(s *LwM2MServer) *MessagerServer {
	m := &MessagerServer{
		Server:      coap.NewServer(s.network, s.address, coap.WithDTLSConfig(s.dtlsConf)),
		lwM2MServer: s,
		network:     s.network,
		address:     s.address,
	}

	m.Router().Use(m.logInterceptor)

	return m
}

func (m *MessagerServer) Start() {
	// register route handlers
	_ = m.Server.Post("/bs", m.onClientBootstrap)         //POST
	_ = m.Server.Get("/bspack", m.onClientBootstrapPack)  //GET
	_ = m.Server.Post("/rd", m.onClientRegister)          //POST
	_ = m.Server.Post("/rd/{id}", m.onClientUpdate)       //POST
	_ = m.Server.Delete("/rd/{id}", m.onClientDeregister) //DELETE
	_ = m.Server.Post("/dp", m.onSendInfo)                //POST

	go func() {
		err := m.Serve()
		if err != nil {
			log.Fatalln("lwm2m messager start failed:", err)
		}
	}()

	log.Infoln("lwm2m messager started at", m.address)
}

func (m *MessagerServer) Stop() {
	m.Shutdown()
	log.Infoln("lwm2m messager stopped")
}

func (m *MessagerServer) checkClientRequest(req coap.Request) error {
	ep := req.Query("ep")
	id := req.SecurityIdentity()

	if m.lwM2MServer.dtlsConf != nil && len(ep) != 0 {
		//If the OSCORE Sender ID is not set to Endpoint Client Name, then the LwM2M Server MUST compare the received
		//Endpoint Client Name identifier with the OSCORE Sender ID of the LwM2M Client. This comparison may either be an
		//equality match or may involve a dedicated lookup table to ensure that LwM2M Clients cannot intentionally or due to
		//misconfiguration impersonate other LwM2M Clients. The LwM2M Server MUST respond with a "4.00 Bad Request" to
		//the LwM2M Client if these fields do not match.
		if ep != id {
			err := BadRequest
			//code := coap.CodeBadRequest
			log.Errorf("error bootstrap client: %v, ep(%s) != CN(%s)", err, ep, id)
			//return m.NewAckResponse(req, code)

			return err
		}
	}

	return nil
}

// handle request parameters like:
//
//	 uri:
//		/bs?ep={Endpoint Client Name}&pct={Preferred Content Format}
func (m *MessagerServer) onClientBootstrap(req coap.Request) coap.Response {
	log.Debugf("receive Bootstrap-Request operation, size=%d bytes", req.Length())

	if err := m.checkClientRequest(req); err != nil {
		code := coap.CodeBadRequest
		return m.NewAckResponse(req, code)
	}

	ep := req.Query("ep")
	addr := req.Address().String()
	err := m.lwM2MServer.bootstrapDelegator.OnRequest(ep, addr)
	code := coap.CodeChanged
	if err != nil {
		log.Errorf("error bootstrap client %s: %v", ep, err)
		code = GetErrorCode(err)
	}

	log.Debugf("Bootstrap-Request operation processed")

	return m.NewAckResponse(req, code)
}

// handle request parameters like:
//
//	 uri:
//		/bspack?ep={Endpoint Client Name}
func (m *MessagerServer) onClientBootstrapPack(req coap.Request) coap.Response {
	log.Debugf("receive Bootstrap-Pack-Request operation, size=%d bytes", req.Length())

	if err := m.checkClientRequest(req); err != nil {
		code := coap.CodeBadRequest
		return m.NewAckResponse(req, code)
	}

	ep := req.Query("ep")
	rspPayload, err := m.lwM2MServer.bootstrapDelegator.OnPackRequest(ep)
	code := coap.CodeContent
	if err != nil {
		log.Errorf("error bootstrap pack client %s: %v", ep, err)
		code = GetErrorCode(err)
	}

	log.Debugf("Bootstrap-Pack-Request operation processed")

	return m.NewAckPiggybackedResponse(req, code, rspPayload)
}

// handle request parameters like:
//
//	uri: /rd?ep={Endpoint Client Name}&lt={Lifetime}
//	        &lwm2m={version}&b={binding}&Q&sms={MSISDN}&pid={ProfileID}
//	   b/Q/sms/pid are optional.
//	body: </1/0>,... which is optional.
func (m *MessagerServer) onClientRegister(req coap.Request) coap.Response {
	log.Debugf("receive Register operation, size=%d bytes", req.Length())

	if err := m.checkClientRequest(req); err != nil {
		code := coap.CodeBadRequest
		return m.NewAckResponse(req, code)
	}

	//The Media-Type of the registration message, if used,
	//MUST be the CoRE Link Format (application/link-format)
	if !req.IsCoRELinkContent() {
		return m.NewAckResponse(req, coap.CodeUnsupportedMediaType)
	}

	ep := req.Query("ep")
	lt, _ := strconv.Atoi(req.Query("lt"))
	lwm2m := req.Query("lwm2m")
	binding := req.Query("b")

	now := time.Now()
	list := coap.ParseCoRELinkString(string(req.Body()))
	info := &RegistrationInfo{
		Name:            ep,
		Address:         req.Address().String(),
		Lifetime:        lt,
		LwM2MVersion:    lwm2m,
		BindingMode:     binding,
		ObjectInstances: list,
		Location:        "",
		RegisterTime:    now,
		RegRenewTime:    now,
		UpdateTime:      now,
	}

	clientId, err := m.lwM2MServer.registerDelegator.OnRegister(info)
	code := coap.CodeCreated
	if err != nil {
		log.Errorf("error registering client %s: %v", ep, err)
		code = GetErrorCode(err)
	}

	rsp := m.NewAckResponse(req, code)
	rsp.SetLocationPath("rd/" + clientId)

	log.Debugf("Register operation processed")

	//enable device management when registration succeeded
	m.lwM2MServer.manager.Enable(clientId)

	return rsp
}

// handle request with parameters like:
//
//	uri: /rd/{location}?lt={Lifetime}&b={binding}&Q&sms={MSISDN}
//		where lt/b/Q/sms are optional.
//	body: </1/0>,... which is optional.
func (m *MessagerServer) onClientUpdate(req coap.Request) coap.Response {
	log.Debugf("receive Update operation, size=%d bytes", req.Length())

	if !req.IsCoRELinkContent() {
		return m.NewAckResponse(req, coap.CodeUnsupportedMediaType)
	}

	// get location from uri
	loc := req.Attribute("id")

	// see m.onClientRegister()
	// ep = SecurityIdentity
	ep := req.SecurityIdentity()
	info := &RegistrationInfo{
		Name:       ep,
		Address:    req.Address().String(),
		Location:   loc,
		UpdateTime: time.Now(),
	}

	//binding := req.Query("b")
	//info.BindingMode = binding
	if len(req.Query("lt")) > 0 {
		lt, _ := strconv.Atoi(req.Query("lt"))
		info.Lifetime = lt
		info.RegRenewTime = info.UpdateTime
	}

	list := coap.ParseCoRELinkString(string(req.Body()))
	if len(list) > 0 {
		info.ObjectInstances = list
	}

	err := m.lwM2MServer.registerDelegator.OnUpdate(info)
	code := coap.CodeChanged
	if err != nil {
		log.Errorf("error updating client %s: %v", info.Name, err)
		code = GetErrorCode(err)
	}

	log.Debugf("Update operation processed")

	return m.NewAckResponse(req, code)
}

// handle request with parameters like:
//
//	uri: /rd/{location}
func (m *MessagerServer) onClientDeregister(req coap.Request) coap.Response {
	log.Debugf("receive Deregister operation, size=%d bytes", req.Length())

	id := req.Attribute("id")
	m.lwM2MServer.registerDelegator.OnDeregister(id)

	log.Debugf("Deregister operation processed")

	m.lwM2MServer.manager.Disable(id)

	return m.NewAckResponse(req, coap.CodeDeleted)
}

// handle request with parameters like:
//
//	uri: /dp
//	body: implementation-specific.
func (m *MessagerServer) onSendInfo(req coap.Request) coap.Response {
	data := req.Body()
	// check resource contained in reported list
	// check server granted read access

	log.Tracef("receive info via Send operation, size=%d bytes", len(data))

	// get registered client bound to this info
	c := m.lwM2MServer.manager.GetByAddr(req.Address().String())
	if c == nil {
		log.Errorf("not registered or address changed, " +
			"a new registration is needed and the info sent is ignored")
		return m.NewAckResponse(req, coap.CodeUnauthorized)
	}

	// commit to application layer
	rsp, err := m.lwM2MServer.reportDelegator.OnSend(c, data)
	if err != nil {
		log.Errorf("error recv client info: %v", err)
		return m.NewAckResponse(req, coap.CodeInternalServerError)
	}

	log.Tracef("process info via Send operation done")

	return m.NewAckPiggybackedResponse(req, coap.CodeChanged, rsp)
}

func (m *MessagerServer) BootstrapDiscover(peer string, oid ObjectID) ([]*coap.CoREResource, error) {
	return m.Discover(Percent, oid, NoneID, NoneID, 1)
}

func (m *MessagerServer) BootstrapWrite(peer string, oid ObjectID, oiId InstanceID, rid ResourceID, value Value) error {
	_, err := m.Write(peer, oid, oiId, rid, NoneID, value)
	return err
}

func (m *MessagerServer) BootstrapDelete(peer string, oid ObjectID, oiId InstanceID) error {
	return m.Delete(peer, oid, oiId, NoneID, NoneID)
}

func (m *MessagerServer) BootstrapFinish(peer string) error {
	req := m.NewGetRequestPlain(BootstrapFinishUri)
	rsp, err := m.SendTo(peer, req)
	if err != nil {
		log.Errorln("bootstrap finish operation failed:", err)
		return err
	}

	// check response code
	if rsp.Code().Changed() {
		log.Debugf("bootstrap finish operation done")
		return nil
	}

	return GetCodeError(rsp.Code())
}

func (m *MessagerServer) Observe(peer string, oid ObjectID, oiId InstanceID, rid ResourceID,
	riId InstanceID, attrs NotificationAttrs, h ObserveHandler) error {
	uri := m.makeAccessPath(oid, oiId, rid, riId)
	req := m.NewGetRequestPlain(uri)
	req.SetObserve(true)
	for k, v := range attrs {
		req.AddQuery(k, v)
	}

	rsp, err := m.SendTo(peer, req)
	if err != nil {
		log.Errorln("observe operation failed:", err)
		return err
	}

	// check response code
	if rsp.Code().Content() {
		log.Debugf("observe client %s at %s done", peer, uri)
		return nil
	}

	return GetCodeError(rsp.Code())
}

func (m *MessagerServer) CancelObservation(peer string, oid ObjectID, oiId InstanceID, rid ResourceID, riId InstanceID) error {
	uri := m.makeAccessPath(oid, oiId, rid, riId)
	req := m.NewGetRequestPlain(uri)
	req.SetObserve(false)
	rsp, err := m.SendTo(peer, req)
	if err != nil {
		log.Errorln("cancel observation operation failed:", err)
		return err
	}

	// check response code
	if rsp.Code().Content() {
		log.Debugf("cancel observation of client %s at %s done", peer, uri)
		return nil
	}

	return GetCodeError(rsp.Code())
}

func (m *MessagerServer) ObserveComposite(peer string, t coap.MediaType, body []byte, h ObserveHandler) ([]byte, error) {
	//if contentType == coap.MediaTypeApplicationSenMLJson {
	//
	//}
	return nil, nil
}

func (m *MessagerServer) CancelObservationComposite(peer string, t coap.MediaType, body []byte) error {
	return nil
}

func (m *MessagerServer) Create(peer string, oid ObjectID, value Value) error {
	return nil
}

func (m *MessagerServer) Read(peer string, oid ObjectID, oiId InstanceID, rid ResourceID, riId InstanceID) ([]byte, error) {
	uri := m.makeAccessPath(oid, oiId, rid, riId)
	req := m.NewGetRequestPlain(uri)
	rsp, err := m.SendTo(peer, req)
	if err != nil {
		log.Errorln("read operation failed:", err)
		return nil, err
	}

	// check response code
	if rsp.Code().Content() {
		log.Debugf("read operation against %s done", uri)
		return rsp.Body(), nil
	}

	return nil, GetCodeError(rsp.Code())
}

func (m *MessagerServer) Discover(peer string, oid ObjectID, oiId InstanceID, rid ResourceID, depth int) ([]*coap.CoREResource, error) {
	panic("implement me")
}

func (m *MessagerServer) makeSenmlBody(oid ObjectID, oiId InstanceID, rid ResourceID, riId InstanceID, value Value) ([]byte, error) {
	pack := senml.Pack{
		Records: make([]senml.Record, 1),
	}

	record := &pack.Records[0]
	SenmlRecordSetFieldValue(record, value)
	record.Name = fmt.Sprintf("/%d/%d/%d/%d", oid, oiId, rid, riId)

	data, err := senml.Encode(pack, senml.JSON)
	if err == nil {
		return data, nil
	}

	return nil, err
}

func (m *MessagerServer) unpackSenmlBody(oid ObjectID, oiId InstanceID, rid ResourceID, riId InstanceID, body []byte) ([]byte, error) {
	pack, err := senml.Decode(body, senml.JSON)
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("/%d/%d/%d/%d", oid, oiId, rid, riId)

	for _, r := range pack.Records {
		if r.Name == name {
			if r.OpaqueValue != nil {
				return Opaque([]byte(*r.OpaqueValue)).ToBytes(), nil
			} else if r.Value != nil {
				return Float64(*r.Value).ToBytes(), nil
			} else if r.StringValue != nil {
				return String(*r.StringValue).ToBytes(), nil
			} else if r.BoolValue != nil {
				return Boolean(*r.BoolValue).ToBytes(), nil
			}
		}
	}
	return nil, nil // no ack
}

func (m *MessagerServer) Write(peer string, oid ObjectID, oiId InstanceID, rid ResourceID, riId InstanceID, value Value) ([]byte, error) {
	uri := m.makeAccessPath(oid, oiId, rid, riId)
	body, err := m.makeSenmlBody(oid, oiId, rid, riId, value)
	if err != nil {
		log.Errorln("make m2m msg failed:", err)
		return nil, err
	}
	req := m.NewPutRequestPlain(uri, body)
	rsp, err := m.SendTo(peer, req)
	if err != nil {
		log.Errorln("write operation failed:", err)
		return nil, err
	}

	// check response code
	if rsp.Code().Changed() {
		log.Debugf("write operation against %s done", uri)
		return m.unpackSenmlBody(oid, oiId, rid, riId, rsp.Body())
	}

	return nil, GetCodeError(rsp.Code())
}

func (m *MessagerServer) Execute(peer string, oid ObjectID, id InstanceID, rid ResourceID, args string) error {
	return nil
}

func (m *MessagerServer) Delete(peer string, oid ObjectID, oiId InstanceID, rid ResourceID, riId InstanceID) error {
	uri := m.makeAccessPath(oid, oiId, rid, riId)
	req := m.NewPutRequestPlain(uri, nil)
	rsp, err := m.SendTo(peer, req)
	if err != nil {
		log.Errorln("delete operation failed:", err)
		return err
	}

	// check response code
	if rsp.Code().Deleted() {
		log.Debugf("delete operation against %s done", uri)
		return nil
	}

	return GetCodeError(rsp.Code())
}

func (m *MessagerServer) makeAccessPath(oid ObjectID, oiId InstanceID, rid ResourceID, riId InstanceID) string {
	optionIds := []uint16{oiId, rid, riId}

	uri := fmt.Sprintf("/%d", oid)
	for _, id := range optionIds {
		if id == NoneID {
			break
		}

		uri += fmt.Sprintf("/%d", id)
	}

	return uri
}

func (m *MessagerServer) statsInterceptor(next coap.Interceptor) coap.Interceptor {
	return coap.Handler(func(w coap.ResponseWriter, r *coap.Message) {
		//m.count++
		next.ServeCOAP(w, r)
	})
}

func (m *MessagerServer) logInterceptor(next coap.Interceptor) coap.Interceptor {
	return coap.Handler(func(w coap.ResponseWriter, r *coap.Message) {
		//if r.Code() != codes.NotFound {
		log.Tracef("recv msg from %v, content: %v", w.Conn().RemoteAddr(), r.String())
		//} else {
		//	if r.Code() == codes.NotFound {
		//
		//	}
		//}
		next.ServeCOAP(w, r)
	})
}
