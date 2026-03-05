


//go:build !safe && !codec.safe && !appengine && go1.9 && !gc
// +build !safe,!codec.safe,!appengine,go1.9,!gc

package codec

import (
	"reflect"
	_ "runtime" 
	"unsafe"
)

var unsafeZeroArr [1024]byte




func unsafeGrowslice(typ unsafe.Pointer, old unsafeSlice, cap, incr int) (v unsafeSlice) {
	size := rtsize2(typ)
	if size == 0 {
		return unsafeSlice{unsafe.Pointer(&unsafeZeroArr[0]), old.Len, cap + incr}
	}
	newcap := int(growCap(uint(cap), uint(size), uint(incr)))
	v = unsafeSlice{Data: newarray(typ, newcap), Len: old.Len, Cap: newcap}
	if old.Len > 0 {
		typedslicecopy(typ, v, old)
	}
	
	return
}










func mapStoresElemIndirect(elemsize uintptr) bool { return false }

func mapSet(m, k, v reflect.Value, _ mapKeyFastKind, _, valIsRef bool) {
	var urv = (*unsafeReflectValue)(unsafe.Pointer(&k))
	var kptr = unsafeMapKVPtr(urv)
	urv = (*unsafeReflectValue)(unsafe.Pointer(&v))
	var vtyp = urv.typ
	var vptr = unsafeMapKVPtr(urv)

	urv = (*unsafeReflectValue)(unsafe.Pointer(&m))
	mptr := rvRefPtr(urv)

	vvptr := mapassign(urv.typ, mptr, kptr)
	typedmemmove(vtyp, vvptr, vptr)
}

func mapGet(m, k, v reflect.Value, _ mapKeyFastKind, _, valIsRef bool) (_ reflect.Value) {
	var urv = (*unsafeReflectValue)(unsafe.Pointer(&k))
	var kptr = unsafeMapKVPtr(urv)
	urv = (*unsafeReflectValue)(unsafe.Pointer(&m))
	mptr := rvRefPtr(urv)

	vvptr, ok := mapaccess2(urv.typ, mptr, kptr)

	if !ok {
		return
	}

	urv = (*unsafeReflectValue)(unsafe.Pointer(&v))

	if helperUnsafeDirectAssignMapEntry || valIsRef {
		urv.ptr = vvptr
	} else {
		typedmemmove(urv.typ, urv.ptr, vvptr)
	}

	return v
}
