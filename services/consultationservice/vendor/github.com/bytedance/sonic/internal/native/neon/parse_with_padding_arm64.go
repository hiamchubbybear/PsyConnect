

package neon

import (
	"unsafe"

	"github.com/bytedance/sonic/internal/rt"
)
 
//go:nosplit
func parse_with_padding(parser unsafe.Pointer) (ret int) {
    return __parse_with_padding(rt.NoEscape(parser))
}

func __parse_with_padding(parser unsafe.Pointer) (ret int) 
