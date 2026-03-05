





package sse

import (
    `unsafe`

    `github.com/bytedance/sonic/internal/rt`
)

var F_lookup_small_key func(key unsafe.Pointer, table unsafe.Pointer, lowerOff int) (ret int)

var S_lookup_small_key uintptr

//go:nosplit
func lookup_small_key(key *string, table *[]byte, lowerOff int) (ret int) {
    return F_lookup_small_key(rt.NoEscape(unsafe.Pointer(key)), rt.NoEscape(unsafe.Pointer(table)), lowerOff)
}

