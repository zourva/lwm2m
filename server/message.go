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

// ServerMessager encapsulates and hide
// transport layer details.
//
// Put them together here to make it
// easier when replacing COAP layer later.
type ServerMessager struct {
	*BaseMessager
	registerManager RegisterManager

	// application layer
	bootstrapService BootstrapServer
	registerService  RegistrationServer
	reportingService ReportingServer

	// transport layer
	network string
	address string
	router  *coap2.Router
	server  *coap2.Server

	coapConn coap.Server
}

func NewMessager(server *LwM2MServer) *ServerMessager {
	m := &ServerMessager{
		network: "udp",
		address: server.options.address,
	}

	//m.router = coap2.NewRouter()
	//m.server = coap2.NewServer(m.router)
	m.coapConn = server.coapConn
	m.registerManager = server.GetClientManager()
	m.bootstrapService = NewBootstrapService(server)
	m.registerService = NewRegistrationService(server)
	m.reportingService = NewInfoReportingService(server)

	//m.deviceControlService = NewDeviceControlService(server)

	return m
}

func (m *ServerMessager) Start() {
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

func (m *ServerMessager) Stop() {
	m.server.Stop()
	log.Infoln("lwm2m messager stopped")
}

// handle request parameters like:
//
//	 uri:
//		/bs?ep={Endpoint Client Name}&pct={Preferred Content Format}
func (m *ServerMessager) onClientBootstrap(req coap.Request) coap.Response {
	ep := req.UriQuery("ep")
	addr := req.Address().String()
	err := m.bootstrapService.OnRequest(ep, addr)
	code := coap.CodeChanged
	if err != nil {
		log.Errorf("error bootstrap client %s: %v", ep, err)
		code = GetErrorCode(err)
	}

	return m.NewPiggybackedResponse(req, code, nil)
}

// handle request parameters like:
//
//	 uri:
//		/bspack?ep={Endpoint Client Name}
func (m *ServerMessager) onClientBootstrapPack(req coap.Request) coap.Response {
	ep := req.UriQuery("ep")
	rsp, err := m.bootstrapService.OnPackRequest(ep)
	code := coap.CodeContent
	if err != nil {
		log.Errorf("error bootstrap pack client %s: %v", ep, err)
		code = GetErrorCode(err)
	}

	return m.NewPiggybackedResponse(req, code, coap.NewBytesPayload(rsp))
}

// handle request parameters like:
//
//	uri: /rd?ep={Endpoint Client Name}&lt={Lifetime}
//	        &lwm2m={version}&b={binding}&Q&sms={MSISDN}&pid={ProfileID}
//	   b/Q/sms/pid are optional.
//	body: </1/0>,... which is optional.
func (m *ServerMessager) onClientRegister(req coap.Request) coap.Response {
	ep := req.UriQuery("ep")
	lt, _ := strconv.Atoi(req.UriQuery("lt"))
	lwm2m := req.UriQuery("lwm2m")
	binding := req.UriQuery("b")

	list := coap.CoreResourcesFromString(req.Message().Payload.String())
	info := &RegistrationInfo{
		Name:            ep,
		Address:         req.Address().String(),
		Lifetime:        lt,
		LwM2MVersion:    lwm2m,
		BindingMode:     binding,
		ObjectInstances: list,
		Location:        "",
		RegisterTime:    time.Now(),
		UpdateTime:      time.Now(),
	}

	clientId, err := m.registerService.OnRegister(info)
	code := coap.CodeCreated
	if err != nil {
		log.Errorf("error registering client %s: %v", ep, err)
		code = GetErrorCode(err)
	}

	rsp := m.NewPiggybackedResponse(req, code, coap.NewEmptyPayload())
	rsp.Message().AddOption(coap.OptionLocationPath, "rd/"+clientId)

	return rsp
}

// handle request with parameters like:
//
//	uri: /{location}?lt={Lifetime}&b={binding}&Q&sms={MSISDN}
//		where location has a format of /rd/{id} and b/Q/sms are optional.
//	body: </1/0>,... which is optional.
func (m *ServerMessager) onClientUpdate(req coap.Request) coap.Response {
	id := req.Attribute("id")
	lt, _ := strconv.Atoi(req.UriQuery("lt"))
	binding := req.UriQuery("b")

	list := coap.CoreResourcesFromString(req.Message().Payload.String())
	info := &RegistrationInfo{
		Location:        id,
		Lifetime:        lt,
		BindingMode:     binding,
		ObjectInstances: list,
		UpdateTime:      time.Now(),
	}

	err := m.registerService.OnUpdate(info)
	code := coap.CodeChanged
	if err != nil {
		log.Errorf("error updating client %s: %v", info.Name, err)
		code = GetErrorCode(err)
	}

	return m.NewPiggybackedResponse(req, code, coap.NewEmptyPayload())
}

// handle request with parameters like:
//
//	uri: /{location}
//	 where location has a format of /rd/{id}
func (m *ServerMessager) onClientDeregister(req coap.Request) coap.Response {
	id := req.Attribute("id")

	m.registerService.OnDeregister(id)

	return m.NewPiggybackedResponse(req, coap.CodeDeleted, coap.NewEmptyPayload())
}

// handle request with parameters like:
//
//	uri: /dp
//	body: implementation-specific.
func (m *ServerMessager) onSendInfo(req coap.Request) coap.Response {
	data := req.Message().Payload.GetBytes()
	// check resource contained in reported list
	// check server granted read access

	// get registered client bound to this info
	c := m.registerManager.GetByAddr(req.Address().String())
	if c == nil {
		log.Errorf("not registered or address changed, " +
			"a new registration is needed and the info sent is ignored")
		return m.NewPiggybackedResponse(req, coap.CodeUnauthorized, coap.NewEmptyPayload())
	}

	// commit to application layer
	rsp, err := m.reportingService.OnSend(c, data)
	if err != nil {
		log.Errorf("error recv client info: %v", err)
		return m.NewPiggybackedResponse(req, coap.CodeInternalServerError, coap.NewEmptyPayload())
	}

	return m.NewPiggybackedResponse(req, coap.CodeChanged, coap.NewBytesPayload(rsp))
}

func (m *ServerMessager) Read(peer string, oid ObjectID, oiId InstanceID, rid ResourceID, riId InstanceID) ([]byte, error) {
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

func (m *ServerMessager) Discover(peer string, oid ObjectID) ([]byte, error) {
	panic("implement me")
}

func (m *ServerMessager) Write(peer string, oid ObjectID, oiId InstanceID, rid ResourceID, value Value) error {
	uri := m.makeAccessPath(oid, oiId, rid, NoneID)
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

func (m *ServerMessager) Delete(peer string, oid ObjectID, oiId InstanceID) error {
	uri := m.makeAccessPath(oid, oiId, NoneID, NoneID)
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

func (m *ServerMessager) Finish(peer string) error {
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

func (m *ServerMessager) makeAccessPath(oid ObjectID, oiId InstanceID, rid ResourceID, riId InstanceID) string {
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

func (m *ServerMessager) SendRequest(peer string, req coap.Request) (coap.Response, error) {
	clientAddr, _ := net.ResolveUDPAddr("udp", peer)

	rsp, err := m.coapConn.SendTo(req, clientAddr)
	if err != nil {
		//log.Println(err)
		return nil, err
	}

	return rsp, nil
}

func (m *ServerMessager) statsInterceptor(next coap2.Interceptor) coap2.Interceptor {
	return coap2.Handler(func(w coap2.ResponseWriter, r *coap2.Message) {
		//m.count++
		next.ServeCOAP(w, r)
	})
}

func (m *ServerMessager) logInterceptor(next coap2.Interceptor) coap2.Interceptor {
	return coap2.Handler(func(w coap2.ResponseWriter, r *coap2.Message) {
		log.Debugf("recv msg from %v, content: %v", w.Conn().RemoteAddr(), r.String())
		next.ServeCOAP(w, r)
	})
}

func (m *ServerMessager) patternedRouteHandler(w coap2.ResponseWriter, r *coap2.Message) {

}
