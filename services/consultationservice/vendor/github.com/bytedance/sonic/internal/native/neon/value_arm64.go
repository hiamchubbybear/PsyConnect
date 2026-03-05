


package neon

import (
    `unsafe`
    `github.com/bytedance/sonic/internal/native/types`
)

//go:nosplit
func value(s unsafe.Pointer, n int, p int, v *types.JsonState, flags uint64) (ret int) {
    return __value(s, n, p, v, flags)
}

//go:nosplit
//go:noescape

func __value(s unsafe.Pointer, n int, p int, v *types.JsonState, flags uint64) (ret int)

