// +build go1.20



package rt

import "unsafe"



func GrowSlice(et *GoType, old GoSlice, newCap int) GoSlice {
	if newCap < old.Len {
		panic("growslice's newCap is smaller than old length")
	}
	s := growslice(old.Ptr, newCap, old.Cap, newCap - old.Len, et)
	s.Len = old.Len
	return s
}

//go:linkname growslice runtime.growslice

func growslice(oldPtr unsafe.Pointer, newLen, oldCap, num int, et *GoType) GoSlice
