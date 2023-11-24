package server

import (
	log "github.com/sirupsen/logrus"
	"github.com/zourva/lwm2m/core"
	"sync"
	"time"
)

// RegisteredClientManager manages sessions of clients,
// based on RegisteredClient, that are registered to this server.
type RegisteredClientManager interface {
	Add(info *core.RegistrationInfo) core.RegisteredClient
	Get(name string) core.RegisteredClient
	GetByAddr(addr string) core.RegisteredClient
	GetByLocation(location string) core.RegisteredClient
	Update(info *core.RegistrationInfo) error
	Delete(name string)
	DeleteByLocation(location string)

	Start()
	Stop()

	// Enable enables management of the registered
	// client identified by location.
	Enable(location string)

	// Disable disables management of the registered
	// client identified by location.
	Disable(location string)
}

func NewRegisteredClientManager(server *LwM2MServer) RegisteredClientManager {
	r := &sessionManager{
		server:   server,
		store:    server.options.store,
		provider: server.options.provider,
		registry: server.options.registry,

		sessions:  make(map[string]core.RegisteredClient),
		indexAddr: make(map[string]core.RegisteredClient),
		indexLoc:  make(map[string]core.RegisteredClient),
		quit:      make(chan bool),
	}

	return r
}

// sessionManager implements RegisteredClientManager.
type sessionManager struct {
	server *LwM2MServer // server context

	sessions  map[string]core.RegisteredClient // index ep name -> session
	indexAddr map[string]core.RegisteredClient // index addr -> session
	indexLoc  map[string]core.RegisteredClient // index location -> session
	store     RegInfoStore                     //registration info store
	lock      sync.Mutex                       //TODO: optimize with lock-free

	provider GuidProvider // session id generator
	registry core.ObjectRegistry

	quit chan bool
}

func (r *sessionManager) Start() {
	go r.loop()
}

func (r *sessionManager) Stop() {
	r.quit <- true
}

func (r *sessionManager) Enable(location string) {
	client := r.GetByLocation(location)
	if client == nil {
		log.Warnf("enable client by location %s ignored due to not found", location)
		return
	}

	client.Enable()
	r.server.evtMgr.EmitEvent(core.EventClientRegistered, client)
}

func (r *sessionManager) Disable(location string) {
	client := r.GetByLocation(location)
	if client == nil {
		log.Warnf("disable client by location %s ignored due to not found", location)
		return
	}

	client.Disable()
	r.server.evtMgr.EmitEvent(core.EventClientUnregistered, client)
}

// loop used to maintain session states.
func (r *sessionManager) loop() {
	timer := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-timer.C:
			r.removeStale()
		case <-r.quit:
			log.Infoln("route table check loop quits")
			return
		}
	}
}

// check and remove stale sessions
func (r *sessionManager) removeStale() {
	r.lock.Lock()
	defer r.lock.Unlock()

	for _, session := range r.sessions {
		if session.Timeout() {
			r.delete(session)
		}
	}
}

// this method is not protected, should be guaranteed by callers.
func (r *sessionManager) delete(session core.RegisteredClient) {
	delete(r.sessions, session.Name())
	delete(r.indexLoc, session.Location())
	delete(r.indexAddr, session.Address())
}

func (r *sessionManager) genLocation(epName string) string {
	// TODO: shorten the name and establish the mappings.
	return r.provider.GetGuidWithHint(epName)
}

// GetByAddr returns session by peer ip:port address.
func (r *sessionManager) GetByAddr(addr string) core.RegisteredClient {
	r.lock.Lock()
	defer r.lock.Unlock()

	return r.indexAddr[addr]
}

// GetByLocation returns session by assigned location.
// Used when updating or deletion.
func (r *sessionManager) GetByLocation(location string) core.RegisteredClient {
	r.lock.Lock()
	defer r.lock.Unlock()

	return r.indexLoc[location]
}

// Get returns session by its endpoint client name.
func (r *sessionManager) Get(name string) core.RegisteredClient {
	r.lock.Lock()
	defer r.lock.Unlock()

	return r.sessions[name]
}

// Add creates a new session using the given registration info
// saves it to the internal store, and establishes related indices.
func (r *sessionManager) Add(info *core.RegistrationInfo) core.RegisteredClient {
	r.lock.Lock()
	defer r.lock.Unlock()

	info.Location = r.genLocation(info.Name)
	session := NewRegisteredClient(r.server, info, r.registry)

	err := r.store.Save(session.RegistrationInfo())
	if err != nil {
		return nil
	}

	r.sessions[session.Name()] = session
	r.indexLoc[session.Location()] = session
	r.indexAddr[session.Address()] = session

	log.Infof("a new client %s registered, location = %s", info.Name, info.Location)

	return session
}

// Update updates a session using the given new registration info.
func (r *sessionManager) Update(info *core.RegistrationInfo) error {
	session := r.GetByLocation(info.Location)
	if session == nil {
		return core.NotFound
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	if info.Address != session.Address() {
		delete(r.indexAddr, session.Address())

		//re-create index using the new info
		r.indexAddr[info.Address] = session
	}

	session.Update(info)

	if err := r.store.Update(session.RegistrationInfo()); err != nil {
		//rollback the index updating ?
		log.Errorln("update registration info failed:", err)
		return err
	}

	return nil
}

func (r *sessionManager) Delete(name string) {
	if session := r.Get(name); session != nil {
		r.lock.Lock()
		defer r.lock.Unlock()

		r.store.Delete(session.Name())
		r.delete(session)
	}
}

// DeleteByLocation removes a session including its related
// access indices and registration info.
func (r *sessionManager) DeleteByLocation(location string) {
	if session := r.GetByLocation(location); session != nil {
		r.lock.Lock()
		defer r.lock.Unlock()

		r.store.Delete(session.Name())
		r.delete(session)
	}
}
