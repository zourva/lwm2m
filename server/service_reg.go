package server

import (
	"errors"
	. "github.com/zourva/lwm2m/core"
	"strings"
)

// RegisterServiceDelegator delegates application layer logic
// for client registration procedure at server side.
type RegisterServiceDelegator struct {
	server *LwM2MServer
}

func NewRegistrationServerDelegator(server *LwM2MServer) *RegisterServiceDelegator {
	s := &RegisterServiceDelegator{
		server: server,
	}

	return s
}

// OnRegister registers a client and returns the assigned
// location mapping to the unique endpoint client name.
func (s *RegisterServiceDelegator) OnRegister(info *RegistrationInfo) (string, error) {
	if err := s.validateRegInfo(info); err != nil {
		return "", err
	}

	if s.server.registerService != nil {
		if err := s.server.registerService.Register(info); err != nil {
			return "", err
		}
	}

	// existence check: removes the old one
	client := s.server.manager.GetByAddr(info.Address)
	if client != nil {
		s.server.manager.DeleteByLocation(client.Name())
	}

	// create and save the session
	client = s.server.manager.Add(info)

	return client.Location(), nil
}

func (s *RegisterServiceDelegator) OnUpdate(info *RegistrationInfo) error {
	if s.server.registerService != nil {
		if err := s.server.registerService.Update(info); err != nil {
			return err
		}
	}

	return s.server.manager.Update(info)

	// TODO: forward to application layer
}

func (s *RegisterServiceDelegator) OnDeregister(location string) {
	found := s.server.manager.GetByLocation(location)
	if found == nil {
		return
	}

	if s.server.registerService != nil {
		s.server.registerService.Unregister(found.RegistrationInfo())
	}

	s.server.manager.DeleteByLocation(location)
}

func (s *RegisterServiceDelegator) validateRegInfo(info *RegistrationInfo) error {
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
