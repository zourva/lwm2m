package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack/v5"
	. "github.com/zourva/lwm2m/core"
)

// RegInfoStore defines storage
// operations for client registration info.
type RegInfoStore interface {
	Init()
	Close()

	//Get returns the registration info of the client
	//identified by global unique name.
	Get(name string) *RegistrationInfo

	//Save saves the registration info of a client to the store.
	Save(info *RegistrationInfo) error

	//Delete deletes the registration info of a client.
	Delete(name string)

	//Update updates the registration info of a client.
	Update(info *RegistrationInfo) error
}

type InMemoryRegInfoStore struct {
	savedClients map[string]*RegistrationInfo
}

func (db *InMemoryRegInfoStore) Init() {
}

func (db *InMemoryRegInfoStore) Close() {
}

func (db *InMemoryRegInfoStore) Get(name string) *RegistrationInfo {
	return db.savedClients[name]
}

func (db *InMemoryRegInfoStore) Save(c *RegistrationInfo) error {
	if c == nil {
		log.Errorln("invalid registration info")
		return errors.New("invalid registration info")
	}

	db.savedClients[c.Name] = c
	return nil
}

func (db *InMemoryRegInfoStore) Delete(name string) {
	delete(db.savedClients, name)
}

func (db *InMemoryRegInfoStore) Update(info *RegistrationInfo) error {
	old := db.Get(info.Name)
	if old == nil {
		return errors.New("registration info not found")
	}

	old.Update(info)
	return db.Save(old)
}

func NewInMemorySessionStore() *InMemoryRegInfoStore {
	return &InMemoryRegInfoStore{
		savedClients: make(map[string]*RegistrationInfo),
	}
}

type RedisRegInfoStore struct {
	client *redis.Client
}

func NewRedisStore(addr, pwd string) *RedisRegInfoStore {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
		DB:       0, // use default DB
	})

	return &RedisRegInfoStore{
		client: rdb,
	}
}

func (db *RedisRegInfoStore) Init() {

}

func (db *RedisRegInfoStore) Close() {

}

func (db *RedisRegInfoStore) makePrimaryKey(name string) string {
	// create a flattened key: dev_reg_{client_name}
	return fmt.Sprintf("dev_reg_%s", name)
}

func (db *RedisRegInfoStore) makeAddrIndexKey(address string) string {
	// create a flattened key: dev_reg_idx_addr_{address}
	return fmt.Sprintf("dev_reg_idx_addr_%s", address)
}

func (db *RedisRegInfoStore) getSessionByKey(key string) *RegistrationInfo {
	ctx := context.Background()
	val, err := db.client.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Infof("session %s is not found", key)
		return nil
	}

	if err != nil {
		log.Errorln("redis get failed:", err)
		return nil
	}

	s := &RegistrationInfo{}
	err = msgpack.Unmarshal([]byte(val), s)
	if err != nil {
		log.Errorln("getSessionByKey unmarshal failed:", err)
		return nil
	}

	return s
}

func (db *RedisRegInfoStore) delSessionByKey(key string) {
	ctx := context.Background()
	_, err := db.client.Del(ctx, key).Result()
	if err != nil {
		log.Errorln("redis del failed:", err)
	}
}

//func (db *RedisRegInfoStore) GetByAddr(addr string) *RegistrationInfo {
//	ctx := context.Background()
//	key, err := db.client.Get(ctx, addr).Result()
//	if err == redis.Nil {
//		log.Infof("session addr index %s is not found", addr)
//		return nil
//	}
//
//	return db.getSessionByKey(key)
//}
//
//func (db *RedisRegInfoStore) DeleteByAddr(addr string) {
//	ctx := context.Background()
//	key, err := db.client.Get(ctx, addr).Result()
//	if err == redis.Nil {
//		log.Infof("session addr index %s is not found", addr)
//		return
//	}
//
//	db.delSessionByKey(key)
//}

func (db *RedisRegInfoStore) Get(name string) *RegistrationInfo {
	return db.getSessionByKey(db.makePrimaryKey(name))
}

func (db *RedisRegInfoStore) Save(c *RegistrationInfo) error {
	//index := db.makeAddrIndexKey(c.Address)
	//_, err := db.client.Set(context.Background(), index, key, 0).Result()
	//if err != nil {
	//	log.Errorln("redis make index failed:", err)
	//	return err
	//}

	val, err := msgpack.Marshal(c)
	if err != nil {
		log.Errorln("msgpack marshal failed:", err)
		return err
	}

	key := db.makePrimaryKey(c.Name)
	_, err = db.client.Set(context.Background(), key, val, 0).Result()
	if err != nil {
		log.Errorln("redis set failed:", err)
		return err
	}

	return nil
}

func (db *RedisRegInfoStore) Delete(name string) {
	db.delSessionByKey(db.makePrimaryKey(name))
}

func (db *RedisRegInfoStore) Update(info *RegistrationInfo) error {
	old := db.Get(info.Name)
	if old == nil {
		return nil
	}

	old.Update(info)

	//write back the updated one
	err := db.Save(old)
	if err != nil {
		log.Errorln("update registration info failed:", err)
		return err
	}

	return nil
}
