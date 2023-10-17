package core

import (
	log "github.com/sirupsen/logrus"
)

type ObjectInstanceMap map[InstanceID]Object

type InstanceStore struct {
	instances ObjectInstanceMap
}

func (i *InstanceStore) Add(id InstanceID, object Object) {
	object.SetInstanceID(id)
	i.instances[id] = object
}

func (i *InstanceStore) Get(id InstanceID) Object {
	return i.instances[id]
}

func (i *InstanceStore) Empty() bool {
	return len(i.instances) == 0
}

func (i *InstanceStore) GetAll() map[InstanceID]Object {
	return i.instances
}

// Size returns number instances we have.
func (i *InstanceStore) Size() int {
	return len(i.instances)
}

// NextId returns next instance id using Size.
//
// This holds since instance ids are incrementally
// allocated from 0.
func (i *InstanceStore) NextId() InstanceID {
	return InstanceID(i.Size())
}

// ObjectStore implements a data repository
// for storing instances of all enabled objects.
//
// Object instances are indexed from 0 to keep
// accordance with Instance ID.
//
// If an object is defined with no multiple objects
// available, then an Instance ID of 0 is assigned.
//
// ObjectStore uses a reader to load existing object
// instances from some persist external storages,
// uses a factory to create new instances dynamically,
// and uses a writer to persist back.
type ObjectStore struct {
	objects map[ObjectID]*InstanceStore
	//lock sync.Mutex

	factory  ObjectFactory
	accessor ObjectPersistence
}

func NewObjectStore(accessor ObjectPersistence, factory ObjectFactory) *ObjectStore {
	os := &ObjectStore{
		objects:  make(map[ObjectID]*InstanceStore),
		accessor: accessor,
		factory:  factory,
	}

	return os
}

// SaveInstance saves the object to the instance table.
// It replaces the old one if the given instance id already exists.
func (s *ObjectStore) SaveInstance(obj Object) {
	class := obj.GetClass()
	if c, ok := s.objects[class.Id()]; ok {
		c.Add(c.NextId(), obj)
	} else {
		newStore := &InstanceStore{
			instances: make(ObjectInstanceMap),
		}

		newStore.Add(newStore.NextId(), obj)

		s.objects[class.Id()] = newStore
	}

	log.Tracef("save instance %d for %s",
		obj.InstanceID(), obj.GetClass().Name())
}

func (s *ObjectStore) GetAllInstances() map[ObjectID]*InstanceStore {
	return s.objects
}

// GetInstances returns instances of an object
// class or nil if not exist.
func (s *ObjectStore) GetInstances(id ObjectID) ObjectInstanceMap {
	if v, ok := s.objects[id]; ok {
		return v.instances
	}

	return nil
}

// GetInstance returns instance of an object
// class or nil if not exist.
func (s *ObjectStore) GetInstance(oid ObjectID, inst InstanceID) Object {
	return s.objects[oid].Get(inst)
}

// CreateInstance creates a new instance but not saved,
// and thus its instance id is not valid.
func (s *ObjectStore) CreateInstance(oid ObjectID) Object {
	obj := s.factory.Create(oid)
	log.Tracef("create instance for %s", obj.GetClass().Name())
	return obj
}

// ForkInstance creates an instance and saved it to the store.
func (s *ObjectStore) ForkInstance(oid ObjectID) Object {
	inst := s.CreateInstance(oid)
	s.SaveInstance(inst)

	log.Infof("add instance %d for %d-%s",
		inst.InstanceID(), inst.GetClass().Id(), inst.GetClass().Name())

	return inst
}

// GetSingleInstance returns instance 0 of
// an object class or nil if not exist.
func (s *ObjectStore) GetSingleInstance(oid ObjectID) Object {
	return s.GetInstance(oid, 0)
}

// Load reads all instances of all supported
// objects from the external storage.
func (s *ObjectStore) Load() error {
	if s.accessor == nil {
		log.Infoln("accessor is not provided, use preset")
		s.loadPreset()
	}

	s.clear()

	load, err := s.accessor.Load()
	if err != nil {
		return err
	}

	if load == nil {
		log.Infoln("objects and instances loaded are not valid, use preset")
		s.loadPreset()
	} else {
		s.objects = load
	}

	return nil
}

// Flush writes all instances of all
// objects to the external storage.
func (s *ObjectStore) Flush() error {
	if s.accessor == nil {
		return nil
	}

	return s.accessor.Flush(s.objects)
}

func (s *ObjectStore) clear() {
	for object := range s.objects {
		delete(s.objects, object)
	}

	s.objects = make(map[ObjectID]*InstanceStore, 0)
}

// loadPreset creates object instances
// based on initial preset classes.
func (s *ObjectStore) loadPreset() {
	s.ForkInstance(OmaObjectSecurity)
	s.ForkInstance(OmaObjectSecurity)
	s.ForkInstance(OmaObjectSecurity)
	s.ForkInstance(OmaObjectServer)
	s.ForkInstance(OmaObjectAccessControl)
	s.ForkInstance(OmaObjectAccessControl)
	s.ForkInstance(OmaObjectAccessControl)
	s.ForkInstance(OmaObjectDevice)
	s.ForkInstance(OmaObjectConnMonitor)
	s.ForkInstance(OmaObjectFirmwareUpdate)
	s.ForkInstance(OmaObjectLocation)
	s.ForkInstance(OmaObjectConnStats)
}
