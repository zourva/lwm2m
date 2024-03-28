package coap

import (
	"github.com/plgd-dev/go-coap/v3/message/codes"
	"github.com/plgd-dev/go-coap/v3/mux"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRouterAddHandle(t *testing.T) {
	router := NewRouter()
	tests := []struct {
		method  codes.Code
		pattern string
		h       Handler
	}{
		{method: codes.GET, pattern: "/1/2/3"},
		{method: codes.PUT, pattern: "/1/2/3"},
		{method: codes.POST, pattern: "/1/2/3"},
		{method: codes.DELETE, pattern: "/1/2/3"},
	}

	for _, c := range tests {
		err := router.Handle(c.method, c.pattern, func(w mux.ResponseWriter, r *mux.Message) {
			t.Logf("%v, %v\n", w.Conn().RemoteAddr(), r.String())
		})

		assert.Nil(t, err)
	}

	for _, c := range tests {
		err := router.Handle(c.method, c.pattern, func(w mux.ResponseWriter, r *mux.Message) {
			t.Logf("%v, %v\n", w.Conn().RemoteAddr(), r.String())
		})

		assert.NotNil(t, err)
	}

	t.Logf("%v", router)
}
