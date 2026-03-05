



















package atomic

import (
	"sync/atomic"
	"unsafe"
)


type UnsafePointer struct {
	_ nocmp 

	v unsafe.Pointer
}


func NewUnsafePointer(val unsafe.Pointer) *UnsafePointer {
	return &UnsafePointer{v: val}
}


func (p *UnsafePointer) Load() unsafe.Pointer {
	return atomic.LoadPointer(&p.v)
}


func (p *UnsafePointer) Store(val unsafe.Pointer) {
	atomic.StorePointer(&p.v, val)
}


func (p *UnsafePointer) Swap(val unsafe.Pointer) (old unsafe.Pointer) {
	return atomic.SwapPointer(&p.v, val)
}


func (p *UnsafePointer) CAS(old, new unsafe.Pointer) (swapped bool) {
	return atomic.CompareAndSwapPointer(&p.v, old, new)
}
