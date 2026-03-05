package reflect2

import "unsafe"

//go:linkname unsafe_New reflect.unsafe_New
func unsafe_New(rtype unsafe.Pointer) unsafe.Pointer

//go:linkname typedmemmove reflect.typedmemmove
func typedmemmove(rtype unsafe.Pointer, dst, src unsafe.Pointer)

//go:linkname unsafe_NewArray reflect.unsafe_NewArray
func unsafe_NewArray(rtype unsafe.Pointer, length int) unsafe.Pointer



//go:linkname typedslicecopy reflect.typedslicecopy
//go:noescape
func typedslicecopy(elemType unsafe.Pointer, dst, src sliceHeader) int

//go:linkname mapassign reflect.mapassign
//go:noescape
func mapassign(rtype unsafe.Pointer, m unsafe.Pointer, key unsafe.Pointer, val unsafe.Pointer)

//go:linkname mapaccess reflect.mapaccess
//go:noescape
func mapaccess(rtype unsafe.Pointer, m unsafe.Pointer, key unsafe.Pointer) (val unsafe.Pointer)

//go:noescape
//go:linkname mapiternext reflect.mapiternext
func mapiternext(it *hiter)

//go:linkname ifaceE2I reflect.ifaceE2I
func ifaceE2I(rtype unsafe.Pointer, src interface{}, dst unsafe.Pointer)




type hiter struct {
	key         unsafe.Pointer
	value       unsafe.Pointer
	t           unsafe.Pointer
	h           unsafe.Pointer
	buckets     unsafe.Pointer
	bptr        unsafe.Pointer
	overflow    *[]unsafe.Pointer
	oldoverflow *[]unsafe.Pointer
	startBucket uintptr
	offset      uint8
	wrapped     bool
	B           uint8
	i           uint8
	bucket      uintptr
	checkBucket uintptr
}








func add(p unsafe.Pointer, x uintptr, whySafe string) unsafe.Pointer {
	return unsafe.Pointer(uintptr(p) + x)
}








func arrayAt(p unsafe.Pointer, i int, eltSize uintptr, whySafe string) unsafe.Pointer {
	return add(p, uintptr(i)*eltSize, "i < len")
}
