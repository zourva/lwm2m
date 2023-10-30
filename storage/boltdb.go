package storage

import (
	"github.com/asdine/storm/v3"
	log "github.com/sirupsen/logrus"
	"github.com/zourva/lwm2m/core"
	bolt "go.etcd.io/bbolt"
	"time"
)

type BoltDBStorage struct {
	name string
	db   *storm.DB
}

func NewBoltDBStorage(name string) *BoltDBStorage {
	db, err := storm.Open(name,
		storm.BoltOptions(0755, &bolt.Options{Timeout: 10 * time.Second}))
	defer db.Close()

	if err != nil {
		log.Fatalf("open boltdb %s failed: %v", name, err)
	}

	storage := &BoltDBStorage{
		db:   db,
		name: name,
	}

	return storage
}

func (s *BoltDBStorage) getObject(oid core.ObjectID) core.Object {
	//return core.NewObject()
	panic("implement me")
}

func (s *BoltDBStorage) getResource(oid core.ObjectID, rid core.ResourceID) core.Resource {
	//return core.NewResource()
	panic("implement me")
}

func (s *BoltDBStorage) Load() (map[core.ObjectID]*core.InstanceManager, error) {
	var all []Instance
	err := s.db.All(&all)
	if err != nil {
		return nil, err
	}

	var m = make(map[core.ObjectID]*core.InstanceManager)
	for _, tuple := range all {
		im, ok := m[tuple.OId]
		if !ok {
			im = core.NewInstanceManager()
			m[tuple.OId] = im
		}

		instance := im.Get(tuple.OIId)
		if instance == nil {
			class := s.getObject(tuple.OId)
			instance = core.NewObjectInstance(class, tuple.OIId, nil)
			im.Add(instance)
		}

		resource := s.getResource(tuple.OId, tuple.RId)
		// todo: check type
		f := core.NewResourceField(tuple.RId, core.String(tuple.Value.(string)))
		if resource.Multiple() {
			instance.AddField(f)
		} else {
			instance.SetSingleField(f)
		}
	}

	return m, nil
}

func (s *BoltDBStorage) Flush(objects map[core.ObjectID]*core.InstanceManager) error {
	return nil
}
