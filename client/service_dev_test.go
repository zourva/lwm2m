package client

import (
	"github.com/stretchr/testify/assert"
	"github.com/zourva/lwm2m/core"
	"github.com/zourva/lwm2m/storage"
	"github.com/zourva/pareto/endec/senml"
	"testing"
)

func newEnabledOperators(db *storage.DBStorage) core.OperatorMap {
	dbOp := storage.NewDBOperator(db)
	enabledOperators := core.OperatorMap{
		//core.OmaObjectSecurity:       &objects.SecurityDefaultOperator{},
		//core.OmaObjectServer:         &objects.ServerDefaultOperator{},
		//core.OmaObjectAccessControl:  &objects.AccessControlDefaultOperator{},
		//core.OmaObjectDevice:         &objects.DeviceDefaultOperator{},
		//core.OmaObjectConnMonitor:    &objects.ConnMonitorDefaultOperator{},
		//core.OmaObjectFirmwareUpdate: &objects.FirmwareDefaultOperator{},
		//core.OmaObjectLocation:       &objects.LocationDefaultOperator{},
		//core.OmaObjectConnStats:      &objects.ConnStatsDefaultOperator{},
		core.OmaObjectSecurity:       dbOp,
		core.OmaObjectServer:         dbOp,
		core.OmaObjectAccessControl:  dbOp,
		core.OmaObjectDevice:         dbOp,
		core.OmaObjectConnMonitor:    dbOp,
		core.OmaObjectFirmwareUpdate: dbOp,
		core.OmaObjectLocation:       dbOp,
		core.OmaObjectConnStats:      dbOp,
		//n1.VehicleStatus:             nil,
		//n1.UplinkTransferInfo:        core.NewBaseOperator(),
	}
	return enabledOperators
}

func TestOnCreate(t *testing.T) {
	reg := core.NewObjectRegistry()
	db := storage.NewDBStorage("conf.db")
	if err := db.Open(); err != nil {
		t.Fatalf("create lwm2m client failed")
	}

	enabledOperators := newEnabledOperators(db)
	store := core.NewObjectInstanceStore(reg)
	store.SetStorageManager(db)
	store.SetOperators(enabledOperators)

	c := &LwM2MClient{store: store}
	d := &DeviceController{client: c}

	if err := c.store.Load(); err != nil {
		t.Fatalf("load object instances failed:%v", err)
	}
	create := func(oid uint16, value string) {
		err := d.OnCreate(oid, []byte(value))
		assert.Nil(t, err)
		t.Logf("%s", value)
	}

	tests := []string{
		`[{"bn":"/0/1/","n":"2","v":2},{"n":"3","vd":"-----BEGIN CERTIFICATE-----\nMIIBsjCCAWACFHuYO/9xs/9p7CGT+cr4TWDy8OoKMAoGCCqGSM49BAMCMFUxCzAJ\nBgNVBAYTAkJSMQ8wDQYDVQQIDAZQYXJhbmExETAPBgNVBAcMCEN1cml0aWJhMQww\nCgYDVQQKDANEaXMxFDASBgNVBAMMC2V4YW1wbGUuY29tMB4XDTIzMTIwNjExMjcy\nMVoXDTMzMTIwMzExMjcyMVowdTELMAkGA1UEBhMCQlIxDzANBgNVBAgMBlBhcmFu\nYTERMA8GA1UEBwwIQ3VyaXRpYmExDDAKBgNVBAoMA0RpczEQMA4GA1UEAwwHREYx\nMDAwMTEiMCAGCSqGSIb3DQEJARYTY2xpZW50MUBleGFtcGxlLmNvbTBOMBAGByqG\nSM49AgEGBSuBBAAhAzoABAaCIEBLaIbYYVR07Inoo0TTloZTf6spVvRfUoP5SlTz\nA5G0ivJYljgQR5e6D/xXiRBQhgkVJaR5MAoGCCqGSM49BAMCA0AAMD0CHF9PXREP\niXtlKe6Aap7U7/TkdAqs6Y2zKe2/6ewCHQD9bGbWltuzrO8Z8HKbX/dE6NkGd/Ou\nStKehKd8\n-----END CERTIFICATE-----"},{"n":"5","vd":"-----BEGIN EC PRIVATE KEY-----\nMGgCAQEEHEWtYTxhhOo52Kkpd8Uo0Rl9xDRwMsnLNGVsjDygBwYFK4EEACGhPAM6\nAAQGgiBAS2iG2GFUdOyJ6KNE05aGU3+rKVb0X1KD+UpU8wORtIryWJY4EEeXug/8\nV4kQUIYJFSWkeQ==\n-----END EC PRIVATE KEY-----"},{"n":"10","v":101},{"n":"0","vs":"obts.ibrifuture.com:5684"},{"n":"1","vb":false}]`,
		`[{"bn":"/1/0/","n":"3","v":6000},{"n":"5","v":86400},{"n":"6","vb":true},{"n":"7","vs":"U"},{"n":"0","v":101},{"n":"1","v":86400},{"n":"2","v":300}]`,
		`[{"bn":"/2/0/","n":"1","v":0},{"n":"2/101","v":15},{"n":"3","v":101},{"n":"0","v":1}]`,
		`[{"bn":"/0/1/","n":"3","vd":"-----BEGIN CERTIFICATE-----\nMIIBsjCCAWACFHuYO/9xs/9p7CGT+cr4TWDy8OoKMAoGCCqGSM49BAMCMFUxCzAJ\nBgNVBAYTAkJSMQ8wDQYDVQQIDAZQYXJhbmExETAPBgNVBAcMCEN1cml0aWJhMQww\nCgYDVQQKDANEaXMxFDASBgNVBAMMC2V4YW1wbGUuY29tMB4XDTIzMTIwNjExMjcy\nMVoXDTMzMTIwMzExMjcyMVowdTELMAkGA1UEBhMCQlIxDzANBgNVBAgMBlBhcmFu\nYTERMA8GA1UEBwwIQ3VyaXRpYmExDDAKBgNVBAoMA0RpczEQMA4GA1UEAwwHREYx\nMDAwMTEiMCAGCSqGSIb3DQEJARYTY2xpZW50MUBleGFtcGxlLmNvbTBOMBAGByqG\nSM49AgEGBSuBBAAhAzoABAaCIEBLaIbYYVR07Inoo0TTloZTf6spVvRfUoP5SlTz\nA5G0ivJYljgQR5e6D/xXiRBQhgkVJaR5MAoGCCqGSM49BAMCA0AAMD0CHF9PXREP\niXtlKe6Aap7U7/TkdAqs6Y2zKe2/6ewCHQD9bGbWltuzrO8Z8HKbX/dE6NkGd/Ou\nStKehKd8\n-----END CERTIFICATE-----"},{"n":"5","vd":"-----BEGIN EC PRIVATE KEY-----\nMGgCAQEEHEWtYTxhhOo52Kkpd8Uo0Rl9xDRwMsnLNGVsjDygBwYFK4EEACGhPAM6\nAAQGgiBAS2iG2GFUdOyJ6KNE05aGU3+rKVb0X1KD+UpU8wORtIryWJY4EEeXug/8\nV4kQUIYJFSWkeQ==\n-----END EC PRIVATE KEY-----"},{"n":"10","v":101},{"n":"0","vs":"obts.ibrifuture.com:5684"},{"n":"1","vb":false},{"n":"2","v":2}]`,
		`[{"bn":"/1/0/","n":"0","v":102}]`,
		`[{"bn":"/2/0/","n":"2/102","v":114}]`,
	}

	create(0, tests[0])
	create(1, tests[1])
	create(2, tests[2])
	create(0, tests[3])
	create(1, tests[4])
	create(2, tests[5])
}

func TestOnRead(t *testing.T) {
	reg := core.NewObjectRegistry()
	db := storage.NewDBStorage("conf.db")
	if err := db.Open(); err != nil {
		t.Fatalf("create lwm2m client failed")
	}

	enabledOperators := newEnabledOperators(db)
	store := core.NewObjectInstanceStore(reg)
	store.SetStorageManager(db)
	store.SetOperators(enabledOperators)

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
	read(1, core.NoneID, core.NoneID, core.NoneID)
	read(2, core.NoneID, core.NoneID, core.NoneID)
	read(0, 1, core.NoneID, core.NoneID)
	read(0, 1, 0, core.NoneID)
	read(1, 0, 0, core.NoneID)
	read(2, 0, 2, 102)
}

func TestOnWrite(t *testing.T) {

	reg := core.NewObjectRegistry()
	db := storage.NewDBStorage("conf.db")
	if err := db.Open(); err != nil {
		t.Fatalf("create lwm2m client failed")
	}

	enabledOperators := newEnabledOperators(db)
	store := core.NewObjectInstanceStore(reg)
	store.SetStorageManager(db)
	store.SetOperators(enabledOperators)

	c := &LwM2MClient{store: store}
	d := &DeviceController{client: c}

	if err := c.store.Load(); err != nil {
		t.Fatalf("load object instances failed:%v", err)
	}
	write := func(oid, instId, resId, rfid uint16, value string) {

		err := d.OnWrite(oid, instId, resId, rfid, []byte(value))
		assert.Nil(t, err)

		rsp, err := d.OnRead(oid, instId, resId, rfid)
		assert.Nil(t, err)
		left, _ := senml.Decode(rsp, senml.JSON)
		right, _ := senml.Decode([]byte(value), senml.JSON)

		assert.Equal(t, left, right)
		t.Logf("%s", rsp)
	}

	tests := []string{
		`[{"bn":"/0/1/","n":"2","v":2},{"n":"3","vd":"-----BEGIN CERTIFICATE-----\nMIIBsjCCAWACFHuYO/9xs/9p7CGT+cr4TWDy8OoKMAoGCCqGSM49BAMCMFUxCzAJ\nBgNVBAYTAkJSMQ8wDQYDVQQIDAZQYXJhbmExETAPBgNVBAcMCEN1cml0aWJhMQww\nCgYDVQQKDANEaXMxFDASBgNVBAMMC2V4YW1wbGUuY29tMB4XDTIzMTIwNjExMjcy\nMVoXDTMzMTIwMzExMjcyMVowdTELMAkGA1UEBhMCQlIxDzANBgNVBAgMBlBhcmFu\nYTERMA8GA1UEBwwIQ3VyaXRpYmExDDAKBgNVBAoMA0RpczEQMA4GA1UEAwwHREYx\nMDAwMTEiMCAGCSqGSIb3DQEJARYTY2xpZW50MUBleGFtcGxlLmNvbTBOMBAGByqG\nSM49AgEGBSuBBAAhAzoABAaCIEBLaIbYYVR07Inoo0TTloZTf6spVvRfUoP5SlTz\nA5G0ivJYljgQR5e6D/xXiRBQhgkVJaR5MAoGCCqGSM49BAMCA0AAMD0CHF9PXREP\niXtlKe6Aap7U7/TkdAqs6Y2zKe2/6ewCHQD9bGbWltuzrO8Z8HKbX/dE6NkGd/Ou\nStKehKd8\n-----END CERTIFICATE-----"},{"n":"5","vd":"-----BEGIN EC PRIVATE KEY-----\nMGgCAQEEHEWtYTxhhOo52Kkpd8Uo0Rl9xDRwMsnLNGVsjDygBwYFK4EEACGhPAM6\nAAQGgiBAS2iG2GFUdOyJ6KNE05aGU3+rKVb0X1KD+UpU8wORtIryWJY4EEeXug/8\nV4kQUIYJFSWkeQ==\n-----END EC PRIVATE KEY-----"},{"n":"10","v":101},{"n":"0","vs":"obts.ibrifuture.com:5684"},{"n":"1","vb":false}]`,
		`[{"bn":"/1/0/","n":"3","v":6000},{"n":"5","v":86400},{"n":"6","vb":true},{"n":"7","vs":"U"},{"n":"0","v":101},{"n":"1","v":86400},{"n":"2","v":300}]`,
		`[{"bn":"/2/0/","n":"1","v":0},{"n":"2/101","v":15},{"n":"3","v":101},{"n":"0","v":1}]`,
		`[{"bn":"/0/1/","n":"3","vd":"-----BEGIN CERTIFICATE-----\nMIIBsjCCAWACFHuYO/9xs/9p7CGT+cr4TWDy8OoKMAoGCCqGSM49BAMCMFUxCzAJ\nBgNVBAYTAkJSMQ8wDQYDVQQIDAZQYXJhbmExETAPBgNVBAcMCEN1cml0aWJhMQww\nCgYDVQQKDANEaXMxFDASBgNVBAMMC2V4YW1wbGUuY29tMB4XDTIzMTIwNjExMjcy\nMVoXDTMzMTIwMzExMjcyMVowdTELMAkGA1UEBhMCQlIxDzANBgNVBAgMBlBhcmFu\nYTERMA8GA1UEBwwIQ3VyaXRpYmExDDAKBgNVBAoMA0RpczEQMA4GA1UEAwwHREYx\nMDAwMTEiMCAGCSqGSIb3DQEJARYTY2xpZW50MUBleGFtcGxlLmNvbTBOMBAGByqG\nSM49AgEGBSuBBAAhAzoABAaCIEBLaIbYYVR07Inoo0TTloZTf6spVvRfUoP5SlTz\nA5G0ivJYljgQR5e6D/xXiRBQhgkVJaR5MAoGCCqGSM49BAMCA0AAMD0CHF9PXREP\niXtlKe6Aap7U7/TkdAqs6Y2zKe2/6ewCHQD9bGbWltuzrO8Z8HKbX/dE6NkGd/Ou\nStKehKd8\n-----END CERTIFICATE-----"},{"n":"5","vd":"-----BEGIN EC PRIVATE KEY-----\nMGgCAQEEHEWtYTxhhOo52Kkpd8Uo0Rl9xDRwMsnLNGVsjDygBwYFK4EEACGhPAM6\nAAQGgiBAS2iG2GFUdOyJ6KNE05aGU3+rKVb0X1KD+UpU8wORtIryWJY4EEeXug/8\nV4kQUIYJFSWkeQ==\n-----END EC PRIVATE KEY-----"},{"n":"10","v":101},{"n":"0","vs":"obts.ibrifuture.com:5684"},{"n":"1","vb":false},{"n":"2","v":2}]`,
		`[{"bn":"/0/1/","n":"0","vs":"obts.ibrifuture.com:5684"}]`,
		`[{"bn":"/1/0/","n":"0","v":102}]`,
		`[{"bn":"/2/0/","n":"2/102","v":114}]`,
	}

	//write(0, core.NoneID, core.NoneID, core.NoneID, tests[0])
	//write(1, core.NoneID, core.NoneID, core.NoneID, tests[1])
	//write(2, core.NoneID, core.NoneID, core.NoneID, tests[2])
	//write(0, 1, core.NoneID, core.NoneID, tests[3])
	write(0, 1, 0, core.NoneID, tests[4])
	write(1, 0, 0, core.NoneID, tests[5])
	write(2, 0, 2, 101, tests[6])
}

func TestOnDelete(t *testing.T) {
	reg := core.NewObjectRegistry()
	db := storage.NewDBStorage("conf.db")
	if err := db.Open(); err != nil {
		t.Fatalf("create lwm2m client failed")
	}

	enabledOperators := newEnabledOperators(db)
	store := core.NewObjectInstanceStore(reg)
	store.SetStorageManager(db)
	store.SetOperators(enabledOperators)

	c := &LwM2MClient{store: store}
	d := &DeviceController{client: c}

	if err := c.store.Load(); err != nil {
		t.Fatalf("load object instances failed:%v", err)
	}
	del := func(oid, iid, rid, riid uint16) {
		err := d.OnDelete(oid, iid, rid, riid)
		assert.Nil(t, err)
	}

	//del(0, core.NoneID, core.NoneID, core.NoneID)
	//del(1, core.NoneID, core.NoneID, core.NoneID)
	//del(2, core.NoneID, core.NoneID, core.NoneID)
	//del(0, 1, core.NoneID, core.NoneID)
	//del(0, 1, 0, core.NoneID)
	//del(1, 0, 0, core.NoneID)
	del(2, 0, 2, 102)
}
