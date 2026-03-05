

package abi

import (
    _ `unsafe`

    `github.com/bytedance/sonic/loader/internal/rt`
)

const (
    _G_stackguard0 = 0x10
)

var (
    F_morestack_noctxt = uintptr(rt.FuncAddr(morestack_noctxt))
)

//go:linkname morestack_noctxt runtime.morestack_noctxt
func morestack_noctxt()

