



package protolazy

import (
	"sync/atomic"
	"unsafe"
)

func atomicLoadIndex(p **[]IndexEntry) *[]IndexEntry {
	return (*[]IndexEntry)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(p))))
}
func atomicStoreIndex(p **[]IndexEntry, v *[]IndexEntry) {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(p)), unsafe.Pointer(v))
}
