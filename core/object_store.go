package core

import (
	log "github.com/sirupsen/logrus"
)

type InstanceEnabled struct {
	Object    ObjectID
	Instances []InstanceID
}

// ObjectInstanceStore implements a data repository
// for storing instances of all enabled objects.
//
// Object instances are indexed from 0 to keep
// accordance with Instance ID.
//
// If an object is defined with no multiple objects
// available, then the default Instance ID 0 is assigned.
//
// ObjectInstanceStore can use an instance operators provider,
// registered for each object class, to delegate operations
// such as create/read/write/delete/execute etc., or can use
// a storage manager to load object instances from or to flush
// instances to some external storage.
type ObjectInstanceStore interface {
	ObjectRegistry() ObjectRegistry
	OperatorProvider() InstanceOperatorProvider
	StorageManager() InstanceStorageManager

	SetObjectRegistry(r ObjectRegistry)
	SetOperatorProvider(p InstanceOperatorProvider)
	SetStorageManager(m InstanceStorageManager)
	SetPresetObjects(objects map[ObjectID][]InstanceID)

	CreateInstance(oid ObjectID) ObjectInstance
	GetAllInstances() map[ObjectID]*InstanceManager
	GetInstances(id ObjectID) InstanceMap

	// GetInstance returns instance of an object
	// class or nil if not exist.
	GetInstance(oid ObjectID, inst InstanceID) ObjectInstance

	// GetSingleInstance returns instance 0 of
	// an object class or nil if not exist.
	GetSingleInstance(oid ObjectID) ObjectInstance

	// Load reads all instances of all
	// objects from external storage.
	Load() error

	// Flush writes all instances of all
	// objects to external storage.
	Flush() error
}

// NewObjectInstanceStore returns nil if neither a
// storage manager nor an operators provider is valid.
// If both are provided, provider is used first.
func NewObjectInstanceStore(r ObjectRegistry) ObjectInstanceStore {
	if r == nil {
		log.Errorln("registry must not be nil")
		return nil
	}

	os := &ObjectInstanceStoreImpl{
		registry: r,
		objects:  make(map[ObjectID]*InstanceManager),
	}

	return os
}

type ObjectInstanceStoreImpl struct {
	registry ObjectRegistry
	storage  InstanceStorageManager

	operators InstanceOperatorProvider
	objects   map[ObjectID]*InstanceManager

	preset map[ObjectID][]InstanceID
}

func (s *ObjectInstanceStoreImpl) ObjectRegistry() ObjectRegistry {
	return s.registry
}

func (s *ObjectInstanceStoreImpl) OperatorProvider() InstanceOperatorProvider {
	return s.operators
}

func (s *ObjectInstanceStoreImpl) StorageManager() InstanceStorageManager {
	return s.storage
}

func (s *ObjectInstanceStoreImpl) SetObjectRegistry(r ObjectRegistry) {
	s.registry = r
}

func (s *ObjectInstanceStoreImpl) SetOperatorProvider(p InstanceOperatorProvider) {
	s.operators = p
}

func (s *ObjectInstanceStoreImpl) SetStorageManager(m InstanceStorageManager) {
	s.storage = m
}

func (s *ObjectInstanceStoreImpl) SetPresetObjects(objects map[ObjectID][]InstanceID) {
	s.preset = objects
}

func (s *ObjectInstanceStoreImpl) GetAllInstances() map[ObjectID]*InstanceManager {
	return s.objects
}

// GetInstances returns instances of an object
// class or nil if not exist.
func (s *ObjectInstanceStoreImpl) GetInstances(id ObjectID) InstanceMap {
	if v, ok := s.objects[id]; ok {
		return v.instances
	}

	return nil
}

// GetInstance returns instance of an object
// class or nil if not exist.
func (s *ObjectInstanceStoreImpl) GetInstance(oid ObjectID, inst InstanceID) ObjectInstance {
	return s.objects[oid].Get(inst)
}

// CreateInstance creates an instance and saved it to the store.
func (s *ObjectInstanceStoreImpl) CreateInstance(oid ObjectID) ObjectInstance {
	class := s.registry.GetClass(oid)
	operator := s.operators.Get(oid)
	if operator == nil {
		log.Warnf("instance operator for object %d is not enabled", oid)
		return nil
	}

	inst := operator.Create(class)
	if inst == nil {
		log.Warnf("instance creation for object %d failed", oid)
		return nil
	}

	return inst
	//if c, ok := s.objects[class.Id()]; ok {
	//	c.Add(c.NextId(), inst)
	//} else {
	//	newStore := &InstanceManager{
	//		instances: make(InstanceMap),
	//	}
	//
	//	newStore.Add(newStore.NextId(), inst)
	//
	//	s.objects[class.Id()] = newStore
	//}
	//
	//log.Infof("add instance %d for %d-%s",
	//	inst.InstanceID(), class.Id(), class.Name())
	//return inst
}

// GetSingleInstance returns instance 0 of
// an object class or nil if not exist.
func (s *ObjectInstanceStoreImpl) GetSingleInstance(oid ObjectID) ObjectInstance {
	return s.GetInstance(oid, 0)
}

// Load reads all instances of all supported
// objects from the built-in object classes.
func (s *ObjectInstanceStoreImpl) Load() error {
	if s.storage != nil {
		if load, err := s.storage.Load(); err != nil {
			return err
		} else {
			s.objects = load
			log.Infoln("objects and instances loaded")
			return nil
		}
	}

	s.clear()

	return s.loadPreset()
}

// Flush writes all instances of all
// objects to the external storage.
func (s *ObjectInstanceStoreImpl) Flush() error {
	if s.storage == nil {
		return nil
	}

	return s.storage.Flush(s.objects)
}

// LoadExplicit creates object instances
// based on descriptors passed in.
func (s *ObjectInstanceStoreImpl) loadPreset() error {
	objectCount := 0
	instanceCount := 0

	for oid, instances := range s.preset {
		objectCount++

		if len(instances) == 0 {
			instances = append(instances, 0)
		}

		for range instances {
			instanceCount++
			s.CreateInstance(oid)
		}
	}

	log.Infof("%d preset objects and %d instances spawned",
		objectCount, instanceCount)

	return nil
}

func (s *ObjectInstanceStoreImpl) clear() {
	for object := range s.objects {
		delete(s.objects, object)
	}

	s.objects = make(map[ObjectID]*InstanceManager, 0)
}

type InstanceMap map[InstanceID]ObjectInstance

type InstanceManager struct {
	instances InstanceMap
}

func (i *InstanceManager) Add(id InstanceID, object ObjectInstance) {
	object.SetInstanceID(id)
	i.instances[id] = object
}

func (i *InstanceManager) Get(id InstanceID) ObjectInstance {
	return i.instances[id]
}

func (i *InstanceManager) Empty() bool {
	return len(i.instances) == 0
}

func (i *InstanceManager) GetAll() map[InstanceID]ObjectInstance {
	return i.instances
}

// Size returns number instances we have.
func (i *InstanceManager) Size() int {
	return len(i.instances)
}

// NextId returns next instance id using Size.
//
// This holds since instance ids are incrementally
// allocated from 0.
func (i *InstanceManager) NextId() InstanceID {
	return InstanceID(i.Size())
}
