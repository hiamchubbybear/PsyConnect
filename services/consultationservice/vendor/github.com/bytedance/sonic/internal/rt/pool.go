package rt

import (
	"unsafe"
)

type SlicePool struct {
	pool  unsafe.Pointer
	len   int
	index int
	typ   uintptr
}

func NewPool(typ *GoType, size int) SlicePool {
	return SlicePool{pool: newarray(typ, size), len: size, typ: uintptr(unsafe.Pointer(typ))}
}

func (self *SlicePool) GetSlice(size int) unsafe.Pointer {
	
	if size > self.Remain() {
		return newarray(AsGoType(self.typ), size)
	}

	ptr := PtrAdd(self.pool, uintptr(self.index)* AsGoType(self.typ).Size)
	self.index += size
	return ptr
}

func (self *SlicePool) Remain() int {
	return self.len - self.index
}
