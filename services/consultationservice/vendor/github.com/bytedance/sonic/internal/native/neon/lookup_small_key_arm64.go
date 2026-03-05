

package neon

import (
	"unsafe"

	"github.com/bytedance/sonic/internal/rt"
)

//go:nosplit
func lookup_small_key(key *string, table *[]byte, lowerOff int) (ret int) {
    return __lookup_small_key(rt.NoEscape(unsafe.Pointer(key)), rt.NoEscape(unsafe.Pointer(table)), lowerOff)
}

//go:nosplit
func __lookup_small_key(key unsafe.Pointer, table unsafe.Pointer, lowerOff int) (ret int)
