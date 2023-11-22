package server

import (
	log "github.com/sirupsen/logrus"
	. "github.com/zourva/lwm2m/core"
	"time"
)

// BootstrapContext defines context
// that is used for a client to do bootstrap.
//
// The server will start a separate goroutine
// to run bootstrap sub procedure.
//
// As the bootstrap process is comprised multi-stage
// message exchanging passes, the context provides
// method to respond to client when BootstrapRequest
// is done, and provides methods to access client
// resources in case needed.
type BootstrapContext interface {
	Name() string
	Address() string
	Stale() bool

	//RespondBootstrapRequest(error)

	// Read implements BootstrapRead operation
	//  method: GET
	//  format: TLV, LwM2M CBOR, SenML CBOR or SenML JSON
	//  path: /{Object ID} in LwM2M 1.1 and thereafter, Object ID MUST be '1'
	//     (Server Object) or '2' (Access Control Object)
	//  code may be responded:
	//    2.05 Content
	//    4.00 Bad Request
	//    4.01 Unauthorized
	//    4.04 Not Found
	//    4.05 Method Not Allowed
	//    4.06 Not Acceptable
	Read(oid ObjectID) ([]byte, error)

	// Discover implements BootstrapDiscover operation
	//  method: GET
	//  path: /{Object ID}
	//  code may be responded:
	//    2.05 Content
	//    4.00 Bad Request
	//    4.04 Not Found
	Discover(oid ObjectID) ([]byte, error)

	// Write implements BootstrapWrite operation
	//  method: PUT
	//  path: /{Object ID}
	//        /{Object ID}/{optional Object Instance ID}
	//        /{Object ID}/{optional Object Instance ID}/{optional Resource ID}
	//  code may be responded:
	//    2.04 Changed
	//    4.00 Bad Request
	//    4.15 Unsupported content format
	Write(oid ObjectID, oiId InstanceID, rid ResourceID, value Value) error

	// Delete implements BootstrapDelete operation
	//  method: DELETE
	//  path: /{Object ID}/{Object Instance ID}
	//  code may be responded:
	//    2.02 Deleted
	//    4.00 Bad Request
	Delete(oid ObjectID, oiId InstanceID) error

	// Finish implements BootstrapFinish operation
	//  method: POST
	//  path: /bs
	//  code may be responded:
	//    2.04 Changed
	//    4.00 Bad Request
	//    4.06 Not Acceptable
	Finish() error
}

// BootstrapServerDelegator delegates application layer logic
// for client bootstrap procedure at server side.
type BootstrapServerDelegator struct {
	server *LwM2MServer
	// TODO: lock protection?
	clients map[string]BootstrapContext //name -> addr
	service BootstrapService
}

func (b *BootstrapServerDelegator) OnRequest(name, addr string) error {
	if b.service == nil {
		return NotImplemented
	}

	//TODO: existence check
	if b.get(name) != nil {
		return NotAcceptable
	}

	ctx := &bootstrapContext{
		owner:  b,
		name:   name,
		addr:   addr,
		create: time.Now(),
	}

	b.save(ctx)

	err := b.service.Bootstrap(ctx)
	if err != nil {
		return err
	}

	time.AfterFunc(500*time.Millisecond, func() {
		err = b.service.Bootstrapping(ctx)
		if err != nil {
			log.Errorln("bootstrap failed in procedure:", err)
		}
	})

	log.Infof("bootstrap request from %s accepted", name)

	return nil
}

func (b *BootstrapServerDelegator) OnPackRequest(name string) ([]byte, error) {
	if b.service == nil {
		return nil, NotImplemented
	}

	ctx := &bootstrapContext{
		owner:  b,
		name:   name,
		create: time.Now(),
	}

	pack, err := b.service.BootstrapPack(ctx)
	if err != nil {
		return nil, err
	}

	log.Infof("bootstrap-pack-request from %s accepted", name)

	return pack, nil
}

func (b *BootstrapServerDelegator) save(ctx BootstrapContext) {
	b.clients[ctx.Name()] = ctx
}

func (b *BootstrapServerDelegator) get(name string) BootstrapContext {
	return b.clients[name]
}

func NewBootstrapServerDelegator(server *LwM2MServer, service BootstrapService) BootstrapServer {
	s := &BootstrapServerDelegator{
		server:  server,
		service: service,
		clients: make(map[string]BootstrapContext),
	}

	return s
}

type bootstrapContext struct {
	owner  *BootstrapServerDelegator
	name   string
	addr   string
	create time.Time
}

func (b *bootstrapContext) Name() string {
	return b.name
}

func (b *bootstrapContext) Address() string {
	return b.addr
}

func (b *bootstrapContext) Stale() bool {
	return time.Now().Sub(b.create) > 10*time.Minute
}

func (b *bootstrapContext) RespondBootstrapRequest(err error) {
	//TODO implement me
	panic("implement me")
}

func (b *bootstrapContext) Read(oid ObjectID) ([]byte, error) {
	return b.owner.server.messager.Read(b.addr, oid, NoneID, NoneID, NoneID)
}

func (b *bootstrapContext) Discover(oid ObjectID) ([]byte, error) {
	return b.owner.server.messager.Discover(b.addr, oid)
}

func (b *bootstrapContext) Write(oid ObjectID, oiId InstanceID, rid ResourceID, value Value) error {
	return b.owner.server.messager.Write(b.addr, oid, oiId, rid, value)
}

func (b *bootstrapContext) Delete(oid ObjectID, oiId InstanceID) error {
	return b.owner.server.messager.Delete(b.addr, oid, oiId)
}

func (b *bootstrapContext) Finish() error {
	return b.owner.server.messager.Finish(b.addr)
}
