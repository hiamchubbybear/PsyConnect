


//go:build go1.10 && (safe || codec.safe || appengine)
// +build go1.10
// +build safe codec.safe appengine

package codec

import "reflect"

func makeMapReflect(t reflect.Type, size int) reflect.Value {
	return reflect.MakeMapWithSize(t, size)
}
