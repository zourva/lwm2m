package server

import (
	log "github.com/sirupsen/logrus"
	"github.com/zourva/lwm2m/core"
	"sync"
	"time"
)

type RegisteredClientManager interface {
	Create(info *core.RegistrationInfo) *RegisteredClient

	Get(name string) *RegisteredClient
	GetByAddr(addr string) *RegisteredClient
	GetByLocation(location string) *RegisteredClient

	Update(info *core.RegistrationInfo) error

	Delete(name string)
	DeleteByLocation(location string)
}

// SessionManager manages client
// sessions locally.
type SessionManager struct {
	sessions  map[string]*RegisteredClient // index ep name -> session
	indexAddr map[string]*RegisteredClient // index addr -> session
	indexLoc  map[string]*RegisteredClient // index location -> session
	store     RegInfoStore                 //registration info store
	lock      sync.Mutex                   //TODO: optimize with lock-free

	provider GuidProvider

	quit chan bool
}

func NewSessionManager(server *LwM2MServer) RegisteredClientManager {
	r := &SessionManager{
		store:    server.opts.store,
		provider: server.opts.provider,

		sessions:  make(map[string]*RegisteredClient),
		indexAddr: make(map[string]*RegisteredClient),
		indexLoc:  make(map[string]*RegisteredClient),
		quit:      make(chan bool),
	}

	return r
}

func (r *SessionManager) Start() {
	go r.loop()
}

func (r *SessionManager) Stop() {
	r.quit <- true
}

// loop used to maintain session states.
func (r *SessionManager) loop() {
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
func (r *SessionManager) removeStale() {
	r.lock.Lock()
	defer r.lock.Unlock()

	for _, session := range r.sessions {
		if session.timeout() {
			r.delete(session)
		}
	}
}

// this method is not protected, should be guaranteed by callers.
func (r *SessionManager) delete(session *RegisteredClient) {
	delete(r.sessions, session.Name())
	delete(r.indexLoc, session.Location())
	delete(r.indexAddr, session.Address())
}

func (r *SessionManager) genLocation(epName string) string {
	// TODO: shorten the name and establish the mappings.
	return r.provider.GetGuidWithHint(epName)
}

// GetByAddr returns session by peer ip:port address.
func (r *SessionManager) GetByAddr(addr string) *RegisteredClient {
	r.lock.Lock()
	defer r.lock.Unlock()

	return r.indexAddr[addr]
}

// GetByLocation returns session by assigned location.
// Used when updating or deletion.
func (r *SessionManager) GetByLocation(location string) *RegisteredClient {
	r.lock.Lock()
	defer r.lock.Unlock()

	return r.indexLoc[location]
}

// Get returns session by its endpoint client name.
func (r *SessionManager) Get(name string) *RegisteredClient {
	r.lock.Lock()
	defer r.lock.Unlock()

	return r.sessions[name]
}

// Create creates a new session using the given registration info
// saves it to the internal store, and establishes related indices.
func (r *SessionManager) Create(info *core.RegistrationInfo) *RegisteredClient {
	r.lock.Lock()
	defer r.lock.Unlock()

	info.Location = r.genLocation(info.Name)
	session := NewClient(info)

	err := r.store.Save(session.regInfo)
	if err != nil {
		return nil
	}

	r.sessions[session.Name()] = session
	r.indexLoc[session.Location()] = session
	r.indexAddr[session.Address()] = session

	return session
}

// Update updates a session using the given new registration info.
func (r *SessionManager) Update(info *core.RegistrationInfo) error {
	session := r.GetByLocation(info.Location)
	if session == nil {
		return core.Errors(core.ClientNotFound)
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	if info.Address != session.Address() {
		delete(r.indexAddr, session.Address())

		//re-create index using the new info
		r.indexAddr[info.Address] = session
	}

	session.Update(info)

	if err := r.store.Update(session.regInfo); err != nil {
		//rollback the index updating ?
		log.Errorln("update registration info failed:", err)
		return err
	}

	return nil
}

func (r *SessionManager) Delete(name string) {
	if session := r.Get(name); session != nil {
		r.lock.Lock()
		defer r.lock.Unlock()

		r.store.Delete(session.Name())
		r.delete(session)
	}
}

// DeleteByLocation removes a session including its related
// access indices and registration info.
func (r *SessionManager) DeleteByLocation(location string) {
	if session := r.GetByLocation(location); session != nil {
		r.lock.Lock()
		defer r.lock.Unlock()

		r.store.Delete(session.Name())
		r.delete(session)
	}
}
