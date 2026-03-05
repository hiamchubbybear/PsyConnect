

package jitdec

import (
    `unsafe`

    `github.com/bytedance/sonic/loader`
)

//go:nosplit
func pbool(v bool) uintptr {
    return freezeValue(unsafe.Pointer(&v))
}

//go:nosplit
func ptodec(p loader.Function) _Decoder {
    return *(*_Decoder)(unsafe.Pointer(&p))
}

func assert_eq(v int64, exp int64, msg string) {
    if v != exp {
        panic(msg)
    }
}
