package core

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/zourva/pareto/endec/senml"
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
	ObjectRegistry() ObjectRegistry
	//OperatorProvider() OperatorProvider
	StorageManager() InstanceStorageManager
	SetObjectRegistry(r ObjectRegistry)
	SetStorageManager(m InstanceStorageManager)

	EnableInstances(m InstanceIdsMap)
	EnableInstance(oid ObjectID, ids ...InstanceID)
	SetOperators(operators OperatorMap)
	SetOperator(id ObjectID, operator Operator)

	//SpawnInstance(oid ObjectID) ObjectInstance

	//GetInstanceManager returns the instance manager
	//for the given object and create a new one if not found.
	GetInstanceManager(id ObjectID) (*InstanceManager, error)
	GetInstanceManagers() InstanceManagers
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
	//Open the underlying storage medium.
	Open() error

	//Close the underlying storage after flushing.
	Close() error

	//Bind binds with an object instance store, making
	//itself an underlying physical storage backing the
	//given store.
	Bind(store ObjectInstanceStore)

	//Load all instances into a map.
	Load() error

	//Flush all instances to storage.
	Flush() error
}

type StoreOption = func(o ObjectInstanceStore)

func WithStorageManager(s InstanceStorageManager) StoreOption {
	return func(o ObjectInstanceStore) {
		o.SetStorageManager(s)
	}
}

func WithOperators(s OperatorMap) StoreOption {
	return func(o ObjectInstanceStore) {
		o.SetOperators(s)
	}
}

func WithOperator(id ObjectID, operator Operator) StoreOption {
	return func(o ObjectInstanceStore) {
		o.SetOperator(id, operator)
	}
}

func WithEnableInstances(s InstanceIdsMap) StoreOption {
	return func(o ObjectInstanceStore) {
		o.EnableInstances(s)
	}
}

func WithEnableInstance(oid ObjectID, ids ...InstanceID) StoreOption {
	return func(o ObjectInstanceStore) {
		o.EnableInstance(oid, ids...)
	}
}

// NewObjectInstanceStore returns nil if neither a
// storage manager nor an operators provider is valid.
// If both are provided, provider is used first.
func NewObjectInstanceStore(r ObjectRegistry, opts ...StoreOption) ObjectInstanceStore {
	if r == nil {
		log.Errorln("registry must not be nil")
		return nil
	}

	os := &objectInstanceStore{
		registry: r,
		enabled:  make(InstanceIdsMap),
		//operators: make(OperatorMap),
		managers: make(map[ObjectID]*InstanceManager),
	}

	for _, fn := range opts {
		fn(os)
	}

	return os
}

type objectInstanceStore struct {
	registry ObjectRegistry
	managers InstanceManagers
	//operators OperatorMap    //operators bound
	enabled InstanceIdsMap //instances enabled

	storage InstanceStorageManager
}

func (s *objectInstanceStore) SetObjectRegistry(r ObjectRegistry) { s.registry = r }
func (s *objectInstanceStore) SetStorageManager(m InstanceStorageManager) {
	m.Bind(s)
	s.storage = m
}

func (s *objectInstanceStore) SetOperator(id ObjectID, operator Operator) {
	//s.operators[id] = operator
	if class := s.registry.GetObject(id); class != nil {
		class.SetOperator(operator)
	}
}

func (s *objectInstanceStore) SetOperators(operators OperatorMap) {
	// add operators by merge
	for id, operator := range operators {
		if operator == nil { // no overwrite if invalid
			continue
		}
		s.SetOperator(id, operator)
	}
}

func (s *objectInstanceStore) EnableInstance(oid ObjectID, ids ...InstanceID) {
	list, ok := s.enabled[oid]
	if !ok {
		s.enabled[oid] = append(s.enabled[oid], ids...)
	}

	s.mergeList(list, ids)
}

func (s *objectInstanceStore) EnableInstances(mapIds InstanceIdsMap) {
	for id, instIds := range mapIds {
		s.EnableInstance(id, instIds...)
	}
}

func (s *objectInstanceStore) ObjectRegistry() ObjectRegistry {
	return s.registry
}

func (s *objectInstanceStore) StorageManager() InstanceStorageManager {
	return s.storage
}

func (s *objectInstanceStore) GetInstanceManager(id ObjectID) (*InstanceManager, error) {
	im, ok := s.managers[id]
	if !ok {
		// 1. 确定是否允许创建
		obj := s.registry.GetObject(id)
		if obj != nil {
			im = NewInstanceManager()
			s.managers[id] = im
		} else {
			// 不支持的类型，不要创建
			log.Errorf("unsupported create object(%d) from registry", id)
			return nil, Forbidden
		}
	}

	return im, nil
}

func (s *objectInstanceStore) GetInstanceManagers() InstanceManagers {
	return s.managers
}

// GetInstances returns instances of an object
// class or nil if not exist.
func (s *objectInstanceStore) GetInstances(id ObjectID) InstanceMap {
	if v, ok := s.managers[id]; ok {
		return v.instances
	}

	return nil
}

// GetInstance returns instance of an object
// class or nil if not exist.
func (s *objectInstanceStore) GetInstance(oid ObjectID, inst InstanceID) ObjectInstance {
	if v, ok := s.managers[oid]; ok {
		return v.Get(inst)
	}
	return nil
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
		if err := s.storage.Load(); err != nil {
			return err
		}

		log.Infoln("object instances loaded, size:", len(s.managers))
		return nil
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

	return s.storage.Flush()
}

// spawnInstance creates an instance and saved it to the store.
// When creating, it spawns a new instance first, then
// invoke the constructor to initialize the instance.
func (s *objectInstanceStore) spawnInstance(oid ObjectID) ObjectInstance {
	class := s.registry.GetObject(oid)
	inst := NewObjectInstance(class)
	if err := inst.Construct(); err != nil {
		log.Warnf("instance construction for object %d failed", oid)
		return nil
	}

	//log.Infof("add instance %d for %d-%s",
	//	inst.Id(), class.Id(), class.Name())
	return inst
}

// LoadExplicit creates object instances
// based on descriptors passed in.
func (s *objectInstanceStore) loadPreset() error {
	objectCount := 0
	instanceCount := 0

	for oid, instances := range s.enabled {
		objectCount++

		if len(instances) == 0 {
			instances = append(instances, 0)
		}

		for range instances {
			instanceCount++
			inst := s.spawnInstance(oid)
			if inst == nil {
				return errors.New("create instance failed")
			}
		}
	}

	log.Infof("%d preset objects and %d instances loaded",
		objectCount, instanceCount)

	return nil
}

func (s *objectInstanceStore) clear() {
	for object := range s.managers {
		delete(s.managers, object)
	}

	s.managers = make(map[ObjectID]*InstanceManager, 0)
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

//
//func (s *objectInstanceStore) getOperator(oid ObjectID) Operator {
//	if operator := s.operators[oid]; operator != nil {
//		log.Warnf("instance operator for object %d is not enabled", oid)
//		return operator
//	}
//
//	return nil
//}

type InstanceManager struct {
	instances InstanceMap
}

type InstanceManagers = map[ObjectID]*InstanceManager

func (i *InstanceManager) Get(id InstanceID) ObjectInstance {
	if v, ok := i.instances[id]; ok {
		return v
	}
	return nil
}

func (i *InstanceManager) Upsert(object ObjectInstance) error {
	//if err := object.Construct(); err != nil {
	//	log.Errorf("call ObjectInstance::Construct failed, %v", err)
	//	return err
	//}

	i.instances[object.Id()] = object
	return nil
}

func (i *InstanceManager) Delete(id InstanceID) error {
	if v, ok := i.instances[id]; ok {
		err := v.Destruct()
		if err != nil {
			log.Errorf("call ObjectInstance::Destruct failed, %v", err)
			return err
		}
		delete(i.instances, id)
	}

	return nil
}

func (i *InstanceManager) Empty() bool {
	return len(i.instances) == 0
}

func (i *InstanceManager) GetAll() InstanceMap {
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

func (i *InstanceManager) MarshalJSON() ([]byte, error) {
	var pack senml.Pack
	var records []senml.Record

	for _, v := range i.instances {
		records = v.AppendSENML(records)
	}

	pack.Records = records
	return senml.Encode(pack, senml.JSON)
}

func NewInstanceManager() *InstanceManager {
	return &InstanceManager{
		instances: make(InstanceMap),
	}
}
