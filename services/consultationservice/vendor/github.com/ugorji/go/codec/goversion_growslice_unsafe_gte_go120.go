


//go:build go1.20 && !safe && !codec.safe && !appengine
// +build go1.20,!safe,!codec.safe,!appengine

package codec

import (
	_ "reflect" 
	"unsafe"
)

func growslice(typ unsafe.Pointer, old unsafeSlice, num int) (s unsafeSlice) {
	
	num -= old.Cap - old.Len
	s = rtgrowslice(old.Data, old.Cap+num, old.Cap, num, typ)
	s.Len = old.Len
	return
}

//go:linkname rtgrowslice runtime.growslice
//go:noescape
func rtgrowslice(oldPtr unsafe.Pointer, newLen, oldCap, num int, typ unsafe.Pointer) unsafeSlice




