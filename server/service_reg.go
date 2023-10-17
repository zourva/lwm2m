package server

import (
	"errors"
	. "github.com/zourva/lwm2m/core"
	"strings"
)

// RegistrationService implements application layer logic
// for client registration procedure at server side.
type RegistrationService struct {
	clientMgr RegisteredClientManager
}

func NewRegistrationService(server *LwM2MServer) RegistrationServer {
	s := &RegistrationService{
		clientMgr: server.manager,
	}

	return s
}

// OnRegister registers a client and returns the assigned
// location mapping to the unique endpoint client name.
func (s *RegistrationService) OnRegister(info *RegistrationInfo) (string, error) {
	if err := s.validateRegInfo(info); err != nil {
		return "", err
	}

	// existence check: removes the old one
	if ssn := s.clientMgr.GetByAddr(info.Address); ssn != nil {
		s.clientMgr.DeleteByLocation(ssn.Name())
	}

	// create and save the session
	session := s.clientMgr.Create(info)

	return session.Location(), nil
}

func (s *RegistrationService) OnUpdate(info *RegistrationInfo) error {
	return s.clientMgr.Update(info)
}

func (s *RegistrationService) OnDeregister(location string) {
	s.clientMgr.DeleteByLocation(location)
}

func (s *RegistrationService) validateRegInfo(info *RegistrationInfo) error {
	if len(info.Name) != 0 {
		// TODO: unique check when necessary
		// urn:uuid:########-####-####-####-############
	}

	// version check
	if len(info.LwM2MVersion) == 0 ||
		(info.LwM2MVersion != "1.0" &&
			info.LwM2MVersion != "1.1") {
		return errors.New("unsupported LwM2M version")
	}

	// object list check
	// LwM2M Security Object ID:0, LwM2M OSCORE Object ID:21,
	// and LwM2M COSE Object ID:23, MUST NOT be part of this list
	for _, o := range info.ObjectInstances {
		t := o.Target[1:len(o.Target)]

		// remove root path
		if strings.Contains(t, " ") {
			t = strings.Split(t, " ")[1]
		}

		// t has format: /1/0> or </1/0>
		if strings.Contains(t, "/0/") ||
			strings.Contains(t, "/21/") ||
			strings.Contains(t, "/23/") {
			return errors.New("unexpected object instances")
		}
	}

	return nil
}
