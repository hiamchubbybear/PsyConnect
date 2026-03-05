

package caching

import (
    `unsafe`

    `github.com/bytedance/sonic/internal/rt`
)

var (
    V_strhash = rt.UnpackEface(rt.Strhash)
    S_strhash = *(*uintptr)(V_strhash.Value)
)

func StrHash(s string) uint64 {
    if v := rt.Strhash(unsafe.Pointer(&s), 0); v == 0 {
        return 1
    } else {
        return uint64(v)
    }
}
