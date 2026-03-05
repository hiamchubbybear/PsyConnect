


//go:build !safe && !codec.safe && !appengine && go1.9 && gc
// +build !safe,!codec.safe,!appengine,go1.9,gc

package codec

import (
	"reflect"
	_ "runtime" 
	"unsafe"
)








const (
	
	mapMaxElemSize = 128
)

func unsafeGrowslice(typ unsafe.Pointer, old unsafeSlice, cap, incr int) (v unsafeSlice) {
	return growslice(typ, old, cap+incr)
}












func mapStoresElemIndirect(elemsize uintptr) bool {
	return elemsize > mapMaxElemSize
}

func mapSet(m, k, v reflect.Value, keyFastKind mapKeyFastKind, valIsIndirect, valIsRef bool) {
	var urv = (*unsafeReflectValue)(unsafe.Pointer(&k))
	var kptr = unsafeMapKVPtr(urv)
	urv = (*unsafeReflectValue)(unsafe.Pointer(&v))
	var vtyp = urv.typ
	var vptr = unsafeMapKVPtr(urv)

	urv = (*unsafeReflectValue)(unsafe.Pointer(&m))
	mptr := rvRefPtr(urv)

	var vvptr unsafe.Pointer

	
	
	
	

	if valIsIndirect {
		vvptr = mapassign(urv.typ, mptr, kptr)
		typedmemmove(vtyp, vvptr, vptr)
		
		return
	}

	switch keyFastKind {
	case mapKeyFastKind32:
		vvptr = mapassign_fast32(urv.typ, mptr, *(*uint32)(kptr))
	case mapKeyFastKind32ptr:
		vvptr = mapassign_fast32ptr(urv.typ, mptr, *(*unsafe.Pointer)(kptr))
	case mapKeyFastKind64:
		vvptr = mapassign_fast64(urv.typ, mptr, *(*uint64)(kptr))
	case mapKeyFastKind64ptr:
		vvptr = mapassign_fast64ptr(urv.typ, mptr, *(*unsafe.Pointer)(kptr))
	case mapKeyFastKindStr:
		vvptr = mapassign_faststr(urv.typ, mptr, *(*string)(kptr))
	default:
		vvptr = mapassign(urv.typ, mptr, kptr)
	}

	
	
	

	typedmemmove(vtyp, vvptr, vptr)
}

func mapGet(m, k, v reflect.Value, keyFastKind mapKeyFastKind, valIsIndirect, valIsRef bool) (_ reflect.Value) {
	var urv = (*unsafeReflectValue)(unsafe.Pointer(&k))
	var kptr = unsafeMapKVPtr(urv)
	urv = (*unsafeReflectValue)(unsafe.Pointer(&m))
	mptr := rvRefPtr(urv)

	var vvptr unsafe.Pointer
	var ok bool

	
	

	switch keyFastKind {
	case mapKeyFastKind32, mapKeyFastKind32ptr:
		vvptr, ok = mapaccess2_fast32(urv.typ, mptr, *(*uint32)(kptr))
	case mapKeyFastKind64, mapKeyFastKind64ptr:
		vvptr, ok = mapaccess2_fast64(urv.typ, mptr, *(*uint64)(kptr))
	case mapKeyFastKindStr:
		vvptr, ok = mapaccess2_faststr(urv.typ, mptr, *(*string)(kptr))
	default:
		vvptr, ok = mapaccess2(urv.typ, mptr, kptr)
	}

	if !ok {
		return
	}

	urv = (*unsafeReflectValue)(unsafe.Pointer(&v))

	if keyFastKind != 0 && valIsIndirect {
		urv.ptr = *(*unsafe.Pointer)(vvptr)
	} else if helperUnsafeDirectAssignMapEntry || valIsRef {
		urv.ptr = vvptr
	} else {
		typedmemmove(urv.typ, urv.ptr, vvptr)
	}

	return v
}

//go:linkname unsafeZeroArr runtime.zeroVal
var unsafeZeroArr [1024]byte





//go:linkname mapassign_fast32 runtime.mapassign_fast32
//go:noescape
func mapassign_fast32(typ unsafe.Pointer, m unsafe.Pointer, key uint32) unsafe.Pointer

//go:linkname mapassign_fast32ptr runtime.mapassign_fast32ptr
//go:noescape
func mapassign_fast32ptr(typ unsafe.Pointer, m unsafe.Pointer, key unsafe.Pointer) unsafe.Pointer

//go:linkname mapassign_fast64 runtime.mapassign_fast64
//go:noescape
func mapassign_fast64(typ unsafe.Pointer, m unsafe.Pointer, key uint64) unsafe.Pointer

//go:linkname mapassign_fast64ptr runtime.mapassign_fast64ptr
//go:noescape
func mapassign_fast64ptr(typ unsafe.Pointer, m unsafe.Pointer, key unsafe.Pointer) unsafe.Pointer

//go:linkname mapassign_faststr runtime.mapassign_faststr
//go:noescape
func mapassign_faststr(typ unsafe.Pointer, m unsafe.Pointer, s string) unsafe.Pointer

//go:linkname mapaccess2_fast32 runtime.mapaccess2_fast32
//go:noescape
func mapaccess2_fast32(typ unsafe.Pointer, m unsafe.Pointer, key uint32) (val unsafe.Pointer, ok bool)

//go:linkname mapaccess2_fast64 runtime.mapaccess2_fast64
//go:noescape
func mapaccess2_fast64(typ unsafe.Pointer, m unsafe.Pointer, key uint64) (val unsafe.Pointer, ok bool)

//go:linkname mapaccess2_faststr runtime.mapaccess2_faststr
//go:noescape
func mapaccess2_faststr(typ unsafe.Pointer, m unsafe.Pointer, key string) (val unsafe.Pointer, ok bool)
