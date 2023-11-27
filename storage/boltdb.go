package storage

import (
	"encoding/json"
	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
	log "github.com/sirupsen/logrus"
	. "github.com/zourva/lwm2m/core"
	"github.com/zourva/lwm2m/objects"
	bolt "go.etcd.io/bbolt"
	"time"
)

type DBStorage struct {
	name  string
	store ObjectInstanceStore //the store bounded with

	db *storm.DB
}

// NewDBStorage pass ":memory:" to use in memory.
func NewDBStorage(name string) *DBStorage {
	storage := &DBStorage{
		name: name,
	}

	return storage
}

func (s *DBStorage) Open() error {
	db, err := storm.Open(s.name,
		storm.BoltOptions(0755, &bolt.Options{Timeout: 10 * time.Second}))

	if err != nil {
		log.Errorf("open boltdb %s failed: %v", s.name, err)
		return err
	}

	s.db = db

	if err = db.Init(&ObjectDescriptor{}); err != nil {
		return err
	}

	if err = db.Init(&ObjectRecord{}); err != nil {
		return err
	}

	err = s.ImportPreset()
	if err != nil {
		return err
	}

	return nil
}

func (s *DBStorage) Close() error {
	if err := s.Flush(); err != nil {
		return err
	}

	return s.db.Close()
}

func (s *DBStorage) Bind(store ObjectInstanceStore) {
	s.store = store
}

func (s *DBStorage) Load() error {
	var all []ObjectRecord
	if err := s.db.All(&all); err != nil {
		return err
	}

	for _, record := range all {
		instance := s.deserialize(&record)
		if instance != nil {
			s.store.GetInstanceManager(instance.Class().Id()).Add(instance)
		}
	}

	return nil
}
func (s *DBStorage) uniqueId(instance ObjectInstance) uint32 {
	unique := uint32(instance.Class().Id()+1)<<16 | uint32(instance.Id()+1)
	return unique
}
func (s *DBStorage) Flush() error {
	var err error
	total := 0
	for _, im := range s.store.GetInstanceManagers() {
		for _, instance := range im.GetAll() {
			record := s.serialize(instance)
			upsert := s.db.Update

			unique := s.uniqueId(instance)
			err = s.db.One("unique", unique, &ObjectRecord{})
			if err != nil {
				upsert = s.db.Save
			}

			err = upsert(record)
			if err != nil {
				log.Errorln("boltdb import ObjectRecord failed:", err)
				return err
			}

			total++
		}
	}

	log.Debugf("boltdb upsert ObjectRecord total records: %d", total)
	return nil
}

func (s *DBStorage) serialize(instance ObjectInstance) *ObjectRecord {
	record := &ObjectRecord{
		Unique:  s.uniqueId(instance),
		Content: instance.String(),
	}

	return record
}

func (s *DBStorage) deserialize(record *ObjectRecord) ObjectInstance {
	registry := s.store.ObjectRegistry()

	instance, err := ParseObjectInstancesWithJSON(registry, record.Content)
	if err != nil {
		log.Errorf("parse object with json failed, err:%v, string:%s", err, record.Content)
		return nil
	}

	return instance[0]
}

func (s *DBStorage) getDBObject(oid ObjectID) *DBObject {
	val := &DBObject{}
	err := s.db.One("Id", oid, val)
	if err != nil {
		log.Errorf("query object %d failed: %v", oid, err)
		return nil
	}

	return val
}

func (s *DBStorage) getDBResource(oid ObjectID, rid ResourceID) *DBResource {
	val := &DBResource{}
	query := s.db.Select(q.Eq("OId", oid), q.Eq("Id", rid))
	err := query.First(val)
	if err != nil {
		log.Errorf("query resource %d for object %d failed: %v", oid, rid, err)
		return nil
	}

	return val
}

func (s *DBStorage) getDBResourceInstance(oid ObjectID, oiId InstanceID,
	rid ResourceID, riId InstanceID) *DBInstance {
	val := &DBInstance{}
	query := s.db.Select(q.Eq("OId", oid), q.Eq("OIId", oiId),
		q.Eq("RId", rid), q.Eq("RIId", riId))
	err := query.First(val)
	if err != nil {
		log.Errorf("query instance for /%d/%d/%d/%d failed: %v", oid, oiId, rid, riId, err)
		return nil
	}

	return val
}

func (s *DBStorage) getDBObservation(oid ObjectID, oiId InstanceID,
	rid ResourceID, riId InstanceID) *DBObservation {
	val := &DBObservation{}
	query := s.db.Select(q.Eq("OId", oid), q.Eq("OIId", oiId),
		q.Eq("RId", rid), q.Eq("RIId", riId))
	err := query.First(val)
	if err != nil {
		log.Errorf("query observation for /%d/%d/%d/%d failed: %v", oid, oiId, rid, riId, err)
		return nil
	}

	return val
}

func (s *DBStorage) InsertInstanceResources(inst ObjectInstance) error {
	//tx, err := s.db.Begin(true)
	//if err != nil {
	//	log.Errorln("InsertInstanceResources begin transaction failed:", err)
	//	return InternalServerError
	//}
	//defer tx.Rollback()
	//
	//for rid, fields := range inst.AllFields() {
	//	for riId, field := range fields {
	//		err = s.db.Save(&DBInstance{
	//			OId:   inst.Class().Id(),
	//			OIId:  inst.Id(),
	//			RId:   rid,
	//			RIId:  InstanceID(riId),
	//			Value: field,
	//		})
	//		if err != nil {
	//			log.Errorln("InsertInstanceResources save failed:", err)
	//			return InternalServerError
	//		}
	//	}
	//}
	//
	//if err = tx.Commit(); err != nil {
	//	log.Errorln("InsertInstanceResources commit transaction failed:", err)
	//	return InternalServerError
	//}

	return ErrorNone
}

func (s *DBStorage) DeleteInstanceResources(inst ObjectInstance) error {
	//tx, err := s.db.Begin(true)
	//if err != nil {
	//	log.Errorln("DeleteInstanceResources begin transaction failed:", err)
	//	return InternalServerError
	//}
	//defer tx.Rollback()
	//
	//s.db.DeleteStruct()
	//
	//if err = tx.Commit(); err != nil {
	//	log.Errorln("DeleteInstanceResources commit transaction failed:", err)
	//	return InternalServerError
	//}

	return ErrorNone
}

func (s *DBStorage) UpdateResourceInstance() error {
	return ErrorNone
}

func (s *DBStorage) ExecuteResourceInstance(id ObjectID, id2 InstanceID, rid ResourceID, id3 InstanceID) error {
	return ErrorNone
}

func (s *DBStorage) InsertResourceInstance(inst ObjectInstance) error {
	return ErrorNone
}

func (s *DBStorage) DeleteResourceInstance(id ObjectID, id2 InstanceID, rid ResourceID, id3 InstanceID) error {
	return ErrorNone
}

func (s *DBStorage) GetObject(id ObjectID) *ObjectDescriptor {
	var val *ObjectDescriptor
	err := s.db.One("Id", id, val)
	if err != nil {
		return nil
	}

	return val
}

func (s *DBStorage) ImportPreset() error {
	descriptors := objects.GetOMAObjectDescriptors()
	for _, desc := range descriptors {
		var od = &ObjectDescriptor{}
		err := json.Unmarshal([]byte(desc), od)
		if err != nil {
			log.Errorln("boltdb import descriptors failed:", err)
			return err
		}

		physicalId := od.Id + 1
		err = s.db.One("Id", physicalId, &ObjectDescriptor{})
		if err == nil { //exist, skip
			continue
		}

		od.Id += 1
		err = s.db.Save(od)
		if err != nil {
			log.Errorln("boltdb import descriptors failed:", err)
			return err
		}

		log.Infoln("boltdb imported object", od.Name)
	}

	return nil
}

type DBOperator struct {
	*BaseOperator
	storage *DBStorage
}

func (o *DBOperator) Construct(inst ObjectInstance) error {
	return o.storage.InsertInstanceResources(inst)
}

func (o *DBOperator) Destruct(inst ObjectInstance) error {
	return o.storage.DeleteInstanceResources(inst)
}

func (o *DBOperator) Add(inst ObjectInstance, rid ResourceID, riId InstanceID, field Field) error {
	return o.storage.InsertResourceInstance(inst)
}

func (o *DBOperator) Update(inst ObjectInstance, rid ResourceID, riId InstanceID, field Field) error {
	return o.storage.UpdateResourceInstance()
}

func (o *DBOperator) Get(inst ObjectInstance, rid ResourceID, riId InstanceID) (Field, error) {
	_ = o.storage.getDBResourceInstance(o.Class().Id(), inst.Id(), rid, riId)
	return nil, ErrorNone
}

func (o *DBOperator) Delete(inst ObjectInstance, rid ResourceID, riId InstanceID) error {
	return o.storage.DeleteResourceInstance(o.Class().Id(), inst.Id(), rid, riId)
}

func (o *DBOperator) Execute(inst ObjectInstance, rid ResourceID, riId InstanceID) error {
	return o.storage.ExecuteResourceInstance(o.Class().Id(), inst.Id(), rid, riId)
}

func NewDBOperator(db *DBStorage) Operator {
	return &DBOperator{
		BaseOperator: NewBaseOperator(),
		storage:      db,
	}
}
