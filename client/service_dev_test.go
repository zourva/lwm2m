package client

import (
	"github.com/stretchr/testify/assert"
	"github.com/zourva/lwm2m/core"
	"github.com/zourva/lwm2m/storage"
	"testing"
)

func TestOnRead(t *testing.T) {
	reg := core.NewObjectRegistry()
	db := storage.NewDBStorage("conf.db")
	if err := db.Open(); err != nil {
		t.Fatalf("create lwm2m client failed")
	}

	store := core.NewObjectInstanceStore(reg)
	store.SetStorageManager(db)

	c := &LwM2MClient{store: store}
	d := &DeviceController{client: c}

	if err := c.store.Load(); err != nil {
		t.Fatalf("load object instances failed:%v", err)
	}
	read := func(oid, instId, resId, rfid uint16) {
		rsp, err := d.OnRead(oid, instId, resId, rfid)
		assert.Nil(t, err)
		t.Logf("%s", rsp)
	}

	read(0, core.NoneID, core.NoneID, core.NoneID)
	read(0, 1, core.NoneID, core.NoneID)
	read(0, 1, 0, core.NoneID)
	read(1, 0, 0, core.NoneID)
	read(2, 0, 2, 101)
}
