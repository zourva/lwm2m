package core

import (
	"testing"
)

func TestOpaque(t *testing.T) {
	//k := 23232

	//kl := unsafe.Pointer(&k)

	o := Opaque([]byte("dsadsa"))

	t.Logf("%v", o.ToBytes())
	t.Logf("%v", o.ToString())

	j, _ := o.MarshalJSON()
	t.Logf("%v", j)
}
