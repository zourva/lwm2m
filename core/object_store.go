package core

import (
	log "github.com/sirupsen/logrus"
)

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
	//ObjectRegistry() ObjectRegistry
	//OperatorProvider() OperatorProvider
	//StorageManager() InstanceStorageManager
	//SetObjectRegistry(r ObjectRegistry)
	//SetStorageManager(m InstanceStorageManager)

	EnableInstances(m InstanceIdsMap)
	EnableInstance(oid ObjectID, ids ...InstanceID)
	SetOperators(operators OperatorMap)
	SetOperator(id ObjectID, operator Operator)

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

type InstanceStorageManager interface {
	Load() (map[ObjectID]*InstanceManager, error)
	Flush(objects map[ObjectID]*InstanceManager) error
}

// NewObjectInstanceStore returns nil if neither a
// storage manager nor an operators provider is valid.
// If both are provided, provider is used first.
func NewObjectInstanceStore(r ObjectRegistry) ObjectInstanceStore {
	if r == nil {
		log.Errorln("registry must not be nil")
		return nil
	}

	os := &objectInstanceStore{
		registry: r,
		objects:  make(map[ObjectID]*InstanceManager),
	}

	return os
}

type objectInstanceStore struct {
	registry ObjectRegistry
	storage  InstanceStorageManager

	operators OperatorMap
	objects   map[ObjectID]*InstanceManager

	preset InstanceIdsMap
}

func (s *objectInstanceStore) SetObjectRegistry(r ObjectRegistry) {
	s.registry = r
}

func (s *objectInstanceStore) SetStorageManager(m InstanceStorageManager) {
	s.storage = m
}

func (s *objectInstanceStore) SetOperator(id ObjectID, operator Operator) {
	s.operators[id] = operator
}

func (s *objectInstanceStore) SetOperators(operators OperatorMap) {
	// add operators by merge
	for id, operator := range operators {
		s.operators[id] = operator
	}
}

func (s *objectInstanceStore) EnableInstance(oid ObjectID, ids ...InstanceID) {
	list, ok := s.preset[oid]
	if !ok {
		s.preset[oid] = append(s.preset[oid], ids...)
	}

	s.mergeList(list, ids)
}

func (s *objectInstanceStore) EnableInstances(mapIds InstanceIdsMap) {
	for id, instIds := range mapIds {
		s.EnableInstance(id, instIds...)
	}
}

func (s *objectInstanceStore) GetAllInstances() map[ObjectID]*InstanceManager {
	return s.objects
}

// GetInstances returns instances of an object
// class or nil if not exist.
func (s *objectInstanceStore) GetInstances(id ObjectID) InstanceMap {
	if v, ok := s.objects[id]; ok {
		return v.instances
	}

	return nil
}

// GetInstance returns instance of an object
// class or nil if not exist.
func (s *objectInstanceStore) GetInstance(oid ObjectID, inst InstanceID) ObjectInstance {
	return s.objects[oid].Get(inst)
}

// CreateInstance creates an instance and saved it to the store.
func (s *objectInstanceStore) CreateInstance(oid ObjectID) ObjectInstance {
	class := s.registry.GetObject(oid)
	operator := s.getOperator(oid)
	if operator == nil {
		return nil
	}

	inst := operator.Construct(class)
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
	//	inst.Id(), class.Id(), class.Name())
	//return inst
}

// GetSingleInstance returns instance 0 of
// an object class or nil if not exist.
func (s *objectInstanceStore) GetSingleInstance(oid ObjectID) ObjectInstance {
	return s.GetInstance(oid, 0)
}

// Load reads all instances of all supported
// objects from the built-in object classes.
func (s *objectInstanceStore) Load() error {
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
func (s *objectInstanceStore) Flush() error {
	if s.storage == nil {
		return nil
	}

	return s.storage.Flush(s.objects)
}

// LoadExplicit creates object instances
// based on descriptors passed in.
func (s *objectInstanceStore) loadPreset() error {
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

func (s *objectInstanceStore) clear() {
	for object := range s.objects {
		delete(s.objects, object)
	}

	s.objects = make(map[ObjectID]*InstanceManager, 0)
}

func (s *objectInstanceStore) mergeList(dstList, newList []InstanceID) {
	for _, id := range dstList {
		for _, instanceID := range newList {
			if id != instanceID {
				dstList = append(dstList, instanceID)
			}
		}
	}
}

func (s *objectInstanceStore) getOperator(oid ObjectID) Operator {
	if operator := s.operators[oid]; operator != nil {
		log.Warnf("instance operator for object %d is not enabled", oid)
		return operator
	}

	return nil
}

type InstanceManager struct {
	instances InstanceMap
}

func (i *InstanceManager) Add(id InstanceID, object ObjectInstance) {
	object.SetId(id)
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
