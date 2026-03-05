

package reflect2

import (
	"unsafe"
)



//go:noescape
//go:linkname mapiterinit reflect.mapiterinit
func mapiterinit(rtype unsafe.Pointer, m unsafe.Pointer) (val *hiter)

func (type2 *UnsafeMapType) UnsafeIterate(obj unsafe.Pointer) MapIterator {
	return &UnsafeMapIterator{
		hiter:      mapiterinit(type2.rtype, *(*unsafe.Pointer)(obj)),
		pKeyRType:  type2.pKeyRType,
		pElemRType: type2.pElemRType,
	}
}