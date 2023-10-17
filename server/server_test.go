package server

import (
	"github.com/stretchr/testify/assert"
	"github.com/zourva/lwm2m/core"
	"github.com/zourva/lwm2m/preset"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	repo := core.NewClassStore(&preset.OMAObjectInfoProvider{})
	srv := New("Test Server",
		WithRegistrationInfoStore(NewInMemorySessionStore()),
		WithObjectFactory(core.NewObjectFactory(repo)))

	go srv.Serve()

	time.AfterFunc(5*time.Second, func() {
		srv.Shutdown()
	})

	assert.NotNil(t, srv)
}
