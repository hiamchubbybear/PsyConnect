


//go:build go1.10 && !safe && !codec.safe && !appengine
// +build go1.10,!safe,!codec.safe,!appengine

package codec

import (
	"reflect"
	"unsafe"
)

func makeMapReflect(typ reflect.Type, size int) (rv reflect.Value) {
	t := (*unsafeIntf)(unsafe.Pointer(&typ)).ptr
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	urv.typ = t
	urv.flag = uintptr(reflect.Map)
	urv.ptr = makemap(t, size, nil)
	return
}

//go:linkname makemap runtime.makemap
//go:noescape
func makemap(typ unsafe.Pointer, size int, h unsafe.Pointer) unsafe.Pointer
