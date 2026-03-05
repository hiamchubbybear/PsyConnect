

package ast

import (
	"unsafe"

	"github.com/bytedance/sonic/internal/rt"
)

//go:nosplit
func mem2ptr(s []byte) unsafe.Pointer {
    return (*rt.GoSlice)(unsafe.Pointer(&s)).Ptr
}

//go:linkname unquoteBytes encoding/json.unquoteBytes
func unquoteBytes(s []byte) (t []byte, ok bool)
