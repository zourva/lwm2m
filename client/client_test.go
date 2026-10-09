package client

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zourva/lwm2m/core"
	"github.com/zourva/lwm2m/storage"
	"github.com/zourva/pareto/box/meta"
)

// TestOnServicingEmitsAbnormalOnReRegister 验证 onServicing 在重注册阈值触发时
// 会 emit EventClientAbnormal。
func TestOnServicingEmitsAbnormalOnReRegister(t *testing.T) {
	c := &LwM2MClient{
		name:     "ut",
		evtMgr:   core.NewEventManager(),
		reporter: NewReporter(nil),
		machine:  meta.NewStateMachine[state]("ut", time.Second),
	}
	c.reporter.failure.Store(4) // 超过阈值 3
	c.evtMgr.RegisterCreator(core.EventClientAbnormal, NewAbnormalEvent)

	got := make(chan core.Event, 1)
	c.OnEvent(core.EventClientAbnormal, func(e core.Event) { got <- e })

	c.onServicing(nil)

	select {
	case e := <-got:
		assert.Equal(t, core.EventClientAbnormal, e.Type())
	case <-time.After(2 * time.Second):
		t.Fatal("EventClientAbnormal was not emitted on re-register")
	}
}

// TestLwM2MClientStopNilSafe 验证 Stop 在 registrar/bootstrapper/messagerc 为 nil、
// 且 store 已带 StorageManager 时不会 panic/挂死。
// 不调用 Start，因此不会触发网络 bootstrap/register，是确定性的单元级回归测试。
func TestLwM2MClientStopNilSafe(t *testing.T) {
	store := newStopTestStore(t)

	c := &LwM2MClient{
		name:     "ut",
		evtMgr:   core.NewEventManager(),
		reporter: NewReporter(nil),
		machine:  meta.NewStateMachine[state]("ut", time.Second),
		store:    store,
	}

	done := make(chan struct{})
	go func() {
		c.Stop()
		close(done)
	}()

	select {
	case <-done:
		// OK: Stop 立即返回，不会阻塞 agent 退出
	case <-time.After(5 * time.Second):
		t.Fatal("Stop hung: agent would not exit on Ctrl+C")
	}
}

// newStopTestStore 复用本包 service_dev_test.go 中的测试设施。
func newStopTestStore(t *testing.T) core.ObjectInstanceStore {
	conf := NewConfCenter()
	reg := core.NewObjectRegistry()
	db := storage.NewConfStorage(conf)
	require.NoError(t, db.Open())

	store := core.NewObjectInstanceStore(reg)
	store.SetStorageManager(db)
	store.SetOperators(newEnabledOperators(db))
	return store
}
