package server

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	srv := New(WithRegistrationInfoStore(NewInMemorySessionStore()))

	go srv.Serve()

	time.AfterFunc(5*time.Second, func() {
		srv.Shutdown()
	})

	assert.NotNil(t, srv)
}
