package server

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zourva/lwm2m/coap"
	"github.com/zourva/lwm2m/coap2"
	. "github.com/zourva/lwm2m/core"
	"net"
	"strconv"
	"time"
)

// MessagerServer encapsulates and hide
// transport layer details.
//
// Put them together here to make it
// easier when replacing COAP layer later.
type MessagerServer struct {
	*BaseMessager
	server *LwM2MServer //server context

	// transport layer
	network string
	address string
	//router  *coap2.Router
	//server  *coap2.Server

	coapConn coap.Server
}

func NewMessager(server *LwM2MServer) *MessagerServer {
	m := &MessagerServer{
		server:  server,
		network: server.network,
		address: server.address,
	}

	//m.router = coap2.NewRouter()
	//m.server = coap2.NewServer(m.router)
	m.coapConn = server.coapConn

	return m
}

func (m *MessagerServer) Start() {
	// setup hooks
	//m.router.Use(m.logInterceptor)

	// register route handlers
	//_ = m.server.Post("/bs", m.onClientBootstrap)        //POST
	//_ = m.server.Get("/bspack", m.onClientBootstrapPack) //GET
	//_ = m.server.Post("/rd", m.onClientRegister)         //POST
	//_ = m.server.Put("/rd/:id", m.onClientUpdate)        //PUT
	//_ = m.server.Delete("/rd/:id", m.onClientDeregister) //DELETE
	//_ = m.server.Post("/dp", m.onSendInfo)               //POST

	//go func() {
	//	err := coap2.ListenAndServe(m.server, m.network, m.address)
	//	if err != nil {
	//		log.Fatalln("lwm2m messager start failed:", err)
	//	} else {
	//		log.Infoln("lwm2m messager started at", m.address)
	//	}
	//}()

	// register route handlers
	m.coapConn.Post("/bs", m.onClientBootstrap)
	m.coapConn.Get("/bspack", m.onClientBootstrapPack)
	m.coapConn.Post("/rd", m.onClientRegister)
	m.coapConn.Put("/rd/:id", m.onClientUpdate)
	m.coapConn.Delete("/rd/:id", m.onClientDeregister)
	m.coapConn.Post("/dp", m.onSendInfo)

	go m.coapConn.Start()
}

func (m *MessagerServer) Stop() {
	//m.server.Stop()
	log.Infoln("lwm2m messager stopped")
}

// handle request parameters like:
//
//	 uri:
//		/bs?ep={Endpoint Client Name}&pct={Preferred Content Format}
func (m *MessagerServer) onClientBootstrap(req coap.Request) coap.Response {
	log.Debugf("receive Bootstrap-Request operation, size=%d bytes",
		req.Message().Payload.Length())

	ep := req.UriQuery("ep")
	addr := req.Address().String()
	err := m.server.bootstrapDelegator.OnRequest(ep, addr)
	code := coap.CodeChanged
	if err != nil {
		log.Errorf("error bootstrap client %s: %v", ep, err)
		code = GetErrorCode(err)
	}

	log.Debugf("Bootstrap-Request operation processed")

	return m.NewPiggybackedResponse(req, code, nil)
}

// handle request parameters like:
//
//	 uri:
//		/bspack?ep={Endpoint Client Name}
func (m *MessagerServer) onClientBootstrapPack(req coap.Request) coap.Response {
	log.Debugf("receive Bootstrap-Pack-Request operation, size=%d bytes",
		req.Message().Payload.Length())

	ep := req.UriQuery("ep")
	rsp, err := m.server.bootstrapDelegator.OnPackRequest(ep)
	code := coap.CodeContent
	if err != nil {
		log.Errorf("error bootstrap pack client %s: %v", ep, err)
		code = GetErrorCode(err)
	}

	log.Debugf("Bootstrap-Pack-Request operation processed")

	return m.NewPiggybackedResponse(req, code, coap.NewBytesPayload(rsp))
}

// handle request parameters like:
//
//	uri: /rd?ep={Endpoint Client Name}&lt={Lifetime}
//	        &lwm2m={version}&b={binding}&Q&sms={MSISDN}&pid={ProfileID}
//	   b/Q/sms/pid are optional.
//	body: </1/0>,... which is optional.
func (m *MessagerServer) onClientRegister(req coap.Request) coap.Response {
	log.Debugf("receive Register operation, size=%d bytes",
		req.Message().Payload.Length())

	//The Media-Type of the registration message, if used,
	//MUST be the CoRE Link Format (application/link-format)
	if mts := req.Message().GetOptions(coap.OptionContentFormat); len(mts) > 0 {
		if mts[0].IntValue() != int(coap.MediaTypeApplicationLinkFormat) {
			return m.NewPiggybackedResponse(req, coap.CodeUnsupportedContentFormat, coap.NewEmptyPayload())
		}
	}

	ep := req.UriQuery("ep")
	lt, _ := strconv.Atoi(req.UriQuery("lt"))
	lwm2m := req.UriQuery("lwm2m")
	binding := req.UriQuery("b")

	now := time.Now()
	list := coap.CoreResourcesFromString(req.Message().Payload.String())
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

	clientId, err := m.server.registerDelegator.OnRegister(info)
	code := coap.CodeCreated
	if err != nil {
		log.Errorf("error registering client %s: %v", ep, err)
		code = GetErrorCode(err)
	}

	rsp := m.NewPiggybackedResponse(req, code, coap.NewEmptyPayload())
	rsp.Message().AddOption(coap.OptionLocationPath, "rd/"+clientId)

	log.Debugf("Register operation processed")

	//enable device management when registration succeeded
	m.server.manager.Enable(clientId)

	return rsp
}

// handle request with parameters like:
//
//	uri: /rd/{location}?lt={Lifetime}&b={binding}&Q&sms={MSISDN}
//		where lt/b/Q/sms are optional.
//	body: </1/0>,... which is optional.
func (m *MessagerServer) onClientUpdate(req coap.Request) coap.Response {
	log.Debugf("receive Update operation, size=%d bytes",
		req.Message().Payload.Length())

	if mts := req.Message().GetOptions(coap.OptionContentFormat); len(mts) > 0 {
		if mts[0].IntValue() != int(coap.MediaTypeApplicationLinkFormat) {
			return m.NewPiggybackedResponse(req, coap.CodeUnsupportedContentFormat, coap.NewEmptyPayload())
		}
	}

	// get location from uri
	loc := req.Attribute("id")
	info := &RegistrationInfo{
		Location:   loc,
		UpdateTime: time.Now(),
	}

	//binding := req.UriQuery("b")
	//info.BindingMode = binding
	if len(req.UriQuery("lt")) > 0 {
		lt, _ := strconv.Atoi(req.UriQuery("lt"))
		info.Lifetime = lt
		info.RegRenewTime = info.UpdateTime
	}

	list := coap.CoreResourcesFromString(req.Message().Payload.String())
	if len(list) > 0 {
		info.ObjectInstances = list
	}

	err := m.server.registerDelegator.OnUpdate(info)
	code := coap.CodeChanged
	if err != nil {
		log.Errorf("error updating client %s: %v", info.Name, err)
		code = GetErrorCode(err)
	}

	log.Debugf("Update operation processed")

	return m.NewPiggybackedResponse(req, code, coap.NewEmptyPayload())
}

// handle request with parameters like:
//
//	uri: /rd/{location}
func (m *MessagerServer) onClientDeregister(req coap.Request) coap.Response {
	log.Debugf("receive Deregister operation, size=%d bytes",
		req.Message().Payload.Length())

	id := req.Attribute("id")
	m.server.registerDelegator.OnDeregister(id)

	log.Debugf("Deregister operation processed")

	m.server.manager.Disable(id)

	return m.NewPiggybackedResponse(req, coap.CodeDeleted, coap.NewEmptyPayload())
}

// handle request with parameters like:
//
//	uri: /dp
//	body: implementation-specific.
func (m *MessagerServer) onSendInfo(req coap.Request) coap.Response {
	data := req.Message().Payload.GetBytes()
	// check resource contained in reported list
	// check server granted read access

	log.Debugf("receive info via Send operation, size=%d bytes", len(data))

	// get registered client bound to this info
	c := m.server.manager.GetByAddr(req.Address().String())
	if c == nil {
		log.Errorf("not registered or address changed, " +
			"a new registration is needed and the info sent is ignored")
		return m.NewPiggybackedResponse(req, coap.CodeUnauthorized, coap.NewEmptyPayload())
	}

	// commit to application layer
	rsp, err := m.server.reportDelegator.OnSend(c, data)
	if err != nil {
		log.Errorf("error recv client info: %v", err)
		return m.NewPiggybackedResponse(req, coap.CodeInternalServerError, coap.NewEmptyPayload())
	}

	log.Debugln("process info via Send operation done")

	return m.NewPiggybackedResponse(req, coap.CodeChanged, coap.NewBytesPayload(rsp))
}

func (m *MessagerServer) BootstrapDiscover(peer string, oid ObjectID) ([]*coap.CoreResource, error) {
	return m.Discover(Percent, oid, NoneID, NoneID, 1)
}

func (m *MessagerServer) BootstrapWrite(peer string, oid ObjectID, oiId InstanceID, rid ResourceID, value Value) error {
	return m.Write(peer, oid, oiId, rid, NoneID, value)
}

func (m *MessagerServer) BootstrapDelete(peer string, oid ObjectID, oiId InstanceID) error {
	return m.Delete(peer, oid, oiId, NoneID, NoneID)
}

func (m *MessagerServer) BootstrapFinish(peer string) error {
	req := m.NewConRequestPlainText(coap.Get, BootstrapFinishUri)
	rsp, err := m.SendRequest(peer, req)
	if err != nil {
		log.Errorln("bootstrap finish operation failed:", err)
		return err
	}

	// check response code
	if rsp.Message().Code == coap.CodeChanged {
		log.Debugf("bootstrap finish operation done")
		return nil
	}

	return GetCodeError(rsp.Message().Code)
}

func (m *MessagerServer) Observe(peer string, oid ObjectID, oiId InstanceID, rid ResourceID,
	riId InstanceID, attrs NotificationAttrs, h ObserveHandler) error {
	uri := m.makeAccessPath(oid, oiId, rid, riId)
	req := m.NewConRequestPlainText(coap.Get, uri)
	req.SetObserve(true)
	for k, v := range attrs {
		req.SetUriQuery(k, v)
	}

	rsp, err := m.SendRequest(peer, req)
	if err != nil {
		log.Errorln("observe operation failed:", err)
		return err
	}

	// check response code
	if rsp.Message().Code == coap.CodeContent {
		log.Debugf("observe client %s at %s done", peer, uri)
		return nil
	}

	return GetCodeError(rsp.Message().Code)
}

func (m *MessagerServer) CancelObservation(peer string, oid ObjectID, oiId InstanceID, rid ResourceID, riId InstanceID) error {
	uri := m.makeAccessPath(oid, oiId, rid, riId)
	req := m.NewConRequestPlainText(coap.Get, uri)
	req.SetObserve(false)
	rsp, err := m.SendRequest(peer, req)
	if err != nil {
		log.Errorln("cancel observation operation failed:", err)
		return err
	}

	// check response code
	if rsp.Message().Code == coap.CodeContent {
		log.Debugf("cancel observation of client %s at %s done", peer, uri)
		return nil
	}

	return GetCodeError(rsp.Message().Code)
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
	req := m.NewConRequestPlainText(coap.Get, uri)
	rsp, err := m.SendRequest(peer, req)
	if err != nil {
		log.Errorln("read operation failed:", err)
		return nil, err
	}

	// check response code
	if rsp.Message().Code == coap.CodeContent {
		log.Debugf("read operation against %s done", uri)
		return rsp.Message().Payload.GetBytes(), nil
	}

	return nil, GetCodeError(rsp.Message().Code)
}

func (m *MessagerServer) Discover(peer string, oid ObjectID, oiId InstanceID, rid ResourceID, depth int) ([]*coap.CoreResource, error) {
	panic("implement me")
}

func (m *MessagerServer) Write(peer string, oid ObjectID, oiId InstanceID, rid ResourceID, riId InstanceID, value Value) error {
	uri := m.makeAccessPath(oid, oiId, rid, riId)
	req := m.NewConRequestPlainText(coap.Put, uri)
	req.SetPayload(value.ToBytes())
	rsp, err := m.SendRequest(peer, req)
	if err != nil {
		log.Errorln("write operation failed:", err)
		return err
	}

	// check response code
	if rsp.Message().Code == coap.CodeChanged {
		log.Debugf("write operation against %s done", uri)
		return nil
	}

	return GetCodeError(rsp.Message().Code)
}

func (m *MessagerServer) Execute(peer string, oid ObjectID, id InstanceID, rid ResourceID, args string) error {
	return nil
}

func (m *MessagerServer) Delete(peer string, oid ObjectID, oiId InstanceID, rid ResourceID, riId InstanceID) error {
	uri := m.makeAccessPath(oid, oiId, rid, riId)
	req := m.NewConRequestPlainText(coap.Put, uri)
	rsp, err := m.SendRequest(peer, req)
	if err != nil {
		log.Errorln("delete operation failed:", err)
		return err
	}

	// check response code
	if rsp.Message().Code == coap.CodeDeleted {
		log.Debugf("delete operation against %s done", uri)
		return nil
	}

	return GetCodeError(rsp.Message().Code)
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

func (m *MessagerServer) SendRequest(peer string, req coap.Request) (coap.Response, error) {
	clientAddr, _ := net.ResolveUDPAddr("udp", peer)

	rsp, err := m.coapConn.SendTo(req, clientAddr)
	if err != nil {
		//log.Println(err)
		return nil, err
	}

	return rsp, nil
}

func (m *MessagerServer) statsInterceptor(next coap2.Interceptor) coap2.Interceptor {
	return coap2.Handler(func(w coap2.ResponseWriter, r *coap2.Message) {
		//m.count++
		next.ServeCOAP(w, r)
	})
}

func (m *MessagerServer) logInterceptor(next coap2.Interceptor) coap2.Interceptor {
	return coap2.Handler(func(w coap2.ResponseWriter, r *coap2.Message) {
		log.Debugf("recv msg from %v, content: %v", w.Conn().RemoteAddr(), r.String())
		next.ServeCOAP(w, r)
	})
}

func (m *MessagerServer) patternedRouteHandler(w coap2.ResponseWriter, r *coap2.Message) {

}
