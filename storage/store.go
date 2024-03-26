package storage

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	. "github.com/zourva/lwm2m/core"
	"github.com/zourva/lwm2m/objects"
	"github.com/zourva/pareto/endec/senml"
	"strconv"
	"strings"
)

type Configer interface {
	All(string) (map[string]string, error)
	Get(string) (string, error)
	Set(string, string) error
	Delete(string)
	// Open() error
	// Flush()
	// Close()
}

type Store struct {
	name  string
	store ObjectInstanceStore //the store bounded with

	conf Configer
}

// NewConfStorage pass ":memory:" to use in memory.
func NewConfStorage(conf Configer) *Store {
	storage := &Store{
		conf: conf,
	}

	_ = storage.ImportPreset()

	return storage
}

func (s *Store) Open() error {
	return nil
}

func (s *Store) Close() error {
	if err := s.Flush(); err != nil {
		return err
	}

	return nil
}

func (s *Store) Bind(store ObjectInstanceStore) {
	s.store = store
}

func (s *Store) Load() error {
	var all map[string]string
	var err error
	if all, err = s.conf.All(RecordPrefix); err != nil {
		return err
	}

	for _, value := range all {
		instance := s.deserialize(value)
		if instance != nil {
			mgr, err := s.store.GetInstanceManager(instance.Class().Id())
			if err != nil {
				log.Errorf("get instance manager failed, %d", instance.Class().Id())
				return err
			}
			mgr.Upsert(instance)
			//s.store.GetInstanceManager(instance.Class().Id()).Upsert(instance)
		}
	}

	return nil
}

func (s *Store) Flush() error {
	total := 0
	for _, im := range s.store.GetInstanceManagers() {
		for _, instance := range im.GetAll() {
			if err := s.InsertInstanceResources(instance); err != nil {
				return InternalServerError
			}
			total++
		}
	}

	log.Debugf("boltdb upsert ObjectRecord total records: %d", total)
	return nil
}

func (s *Store) deserialize(value string) ObjectInstance {
	registry := s.store.ObjectRegistry()

	instance, err := ParseObjectInstancesWithJSON(registry, value)
	if err != nil {
		log.Errorf("parse object with json failed, err:%v, string:%s", err, value)
		return nil
	}

	if len(instance) > 0 {
		return instance[0]

	}
	return nil
}

func (s *Store) InsertInstanceResources(instance ObjectInstance) error {
	key, _, _ := genKey(instance, 0, 0)

	value, err := instance.MarshalJSON()
	if err != nil {
		return InternalServerError
	}

	err = s.conf.Set(key, string(value))
	if err != nil {
		return InternalServerError
	}

	return ErrorNone
}

func (s *Store) DeleteInstanceResources(instance ObjectInstance) error {
	key, _, _ := genKey(instance, 0, 0)

	s.conf.Delete(key)
	return ErrorNone
}

func (s *Store) UpdateResourceInstance() error {
	return ErrorNone
}

func (s *Store) GetResourceInstance(inst ObjectInstance, rid ResourceID, riId InstanceID) (Field, error) {
	var err error
	var pack senml.Pack

	key, _, subName := genKey(inst, rid, riId)

	value, err := s.conf.Get(key)
	if err != nil {
		log.Errorf("get key(%s) failed,err:%v", key, err)
		return nil, err
	}

	pack, err = senml.Decode([]byte(value), senml.JSON)
	if err != nil {
		log.Errorf("decode senml value of key(%s) failed,err:%v, value:%s", key, err, value)
		return nil, err
	}

	for i := 0; i < len(pack.Records); i++ {
		r := &pack.Records[i]
		if r.Name == subName {
			res := inst.Class().Resource(rid)
			kind := res.Type()
			v := SenmlRecordToFieldValue(kind, r)

			field := NewResourceField2(inst, riId, res, v)
			return field, ErrorNone
		}
	}

	return nil, NotFound
}

func (s *Store) ExecuteResourceInstance(id ObjectID, id2 InstanceID, rid ResourceID, id3 InstanceID) error {
	return ErrorNone
}

const (
	RecordPrefix = "record"
	OMAPrefix    = "oma"
	Separator    = "."
)

// genKey
// return
// - path key:     prefix./oid/iid/ 							(e.g. /0/0)
// - full path:    prefix./oid/iid/rid or /oid/iid/rid/riid	(e.g. /0/0/0 or /0/0/0/101
// - sub name: 	   rid or rid/riid 					(e.g. 0      or 0/101)
// value: [{"bn":"/0/0/","n":"0","vs":"obts.ibrifuture.com:5684"},{"n":"1","vb":true}]
func genKey(inst ObjectInstance, rid ResourceID, riId InstanceID) (string, string, string) {
	oid, iid := inst.Class().Id(), inst.Id()
	key := strconv.AppendUint([]byte(`/`), uint64(oid), 10)
	key = append(key, '/')
	key = strconv.AppendUint(key, uint64(iid), 10)
	key = append(key, '/')

	sname := strconv.AppendUint([]byte(nil), uint64(rid), 10)
	if riId != 0 && riId != NoneID {
		sname = append(sname, '/')
		sname = strconv.AppendUint(sname, uint64(riId), 10)
	}

	full := append(key, sname...)

	return RecordPrefix + Separator + string(key),
		RecordPrefix + Separator + string(full),
		string(sname)
}

// InsertResourceInstance
// key:   		/oid/iid/ 							(e.g. /0/0)
// subName: 	rid or rid/riid 					(e.g. 0      or 0/101)
// fullName:    /oid/iid/rid or /oid/iid/rid/riid	(e.g. /0/0/0 or /0/0/0/101
// value: [{"bn":"/0/0/","n":"0","vs":"obts.ibrifuture.com:5684"},{"n":"1","vb":true}]
func (s *Store) InsertResourceInstance(inst ObjectInstance, rid ResourceID, riId InstanceID, field Field) error {
	var data []byte
	var err error
	var pack senml.Pack

	key, _, subName := genKey(inst, rid, riId)
	value, err := s.conf.Get(key)
	if err != nil {
		if err == NotFound {
			// add new field
			data, err = inst.MarshalJSON()
			goto set
		}
		log.Errorf("get key(%s) failed,err:%v", key, err)
		return err
	}

	pack, err = senml.Decode([]byte(value), senml.JSON)
	if err != nil {
		log.Errorf("decode senml value of key(%s) failed,err:%v, value:%s", key, err, value)
		return err
	}

	for i := 0; i < len(pack.Records); i++ {
		r := &pack.Records[i]
		if r.Name == subName {
			SenmlRecordSetFieldValue(r, field)
			goto encode
		}
	}

	pack.Records = field.AppendSENML(pack.Records)

encode:
	data, err = senml.Encode(pack, senml.JSON)
	if err != nil {
		log.Errorf("encode senml value of key(%s) failed, err:%v", key, err)
		return err
	}

set:
	err = s.conf.Set(key, string(data))
	if err != nil {
		log.Errorf("set value of key(%s) failed, err:%v", key, err)
		return err
	}
	return ErrorNone
}

func (s *Store) DeleteResourceInstance(id ObjectID, id2 InstanceID, rid ResourceID, id3 InstanceID) error {
	return ErrorNone
}

func (s *Store) ImportPreset() error {
	genKey := func(d *ObjectDescriptor) string {
		return OMAPrefix + Separator + strconv.Itoa(int(d.Id))
	}
	descriptors := objects.GetOMAObjectDescriptors()
	for _, desc := range descriptors {
		var od = &ObjectDescriptor{}
		err := json.Unmarshal([]byte(desc), od)
		if err != nil {
			log.Errorln("boltdb import descriptors failed:", err)
			return err
		}

		key := genKey(od)
		value := strings.Join(strings.Fields(desc), "")
		//value = strings.ReplaceAll(value, "\n", "")
		s.conf.Set(key, value)

		log.Infoln("boltdb imported object", od.Name)
	}

	return nil
}

type StoreOperator struct {
	*BaseOperator
	storage *Store
}

func (o *StoreOperator) Construct(inst ObjectInstance) error {
	return o.storage.InsertInstanceResources(inst)
}

func (o *StoreOperator) Destruct(inst ObjectInstance) error {
	return o.storage.DeleteInstanceResources(inst)
}

func (o *StoreOperator) Add(inst ObjectInstance, rid ResourceID, riId InstanceID, field Field) error {
	inst.Helper().AddField(field)

	return o.storage.InsertResourceInstance(inst, rid, riId, field)
}

func (o *StoreOperator) Update(inst ObjectInstance, rid ResourceID, riId InstanceID, field Field) error {
	return o.storage.UpdateResourceInstance()
}

func (o *StoreOperator) GetAll(inst ObjectInstance, rid ResourceID) (*Fields, error) {
	// 直接从内存中取
	fields := inst.Helper().Fields(rid)
	if fields == nil {
		return nil, NotFound
	}

	return fields, nil
}

func (o *StoreOperator) Get(inst ObjectInstance, rid ResourceID, riId InstanceID) (Field, error) {
	// 直接从内存中取
	field := inst.Helper().Field(rid, riId)
	if field == nil {
		return nil, NotFound
	}
	return field, nil
	//return o.storage.GetResourceInstance(inst, rid, riId)
}

func (o *StoreOperator) Delete(inst ObjectInstance, rid ResourceID, riId InstanceID) error {
	return o.storage.DeleteResourceInstance(o.Class().Id(), inst.Id(), rid, riId)
}

func (o *StoreOperator) Execute(inst ObjectInstance, rid ResourceID, riId InstanceID) error {
	return o.storage.ExecuteResourceInstance(o.Class().Id(), inst.Id(), rid, riId)
}

func NewConfOperator(conf *Store) Operator {
	return &StoreOperator{
		BaseOperator: NewBaseOperator(),
		storage:      conf,
	}
}
