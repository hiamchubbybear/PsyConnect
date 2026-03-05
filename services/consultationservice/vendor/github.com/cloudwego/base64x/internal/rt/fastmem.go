

package rt

import (
	"unsafe"
)







//go:nosplit

func NoEscape(p unsafe.Pointer) unsafe.Pointer {
    x := uintptr(p)
    return unsafe.Pointer(x ^ 0)
}

//go:nosplit
func MoreStack(size uintptr)

//go:nosplit
func Add(ptr unsafe.Pointer, off uintptr) unsafe.Pointer {
    return unsafe.Pointer(uintptr(ptr) + off)
}
