



package avx2

import (
    `unsafe`

    `github.com/cloudwego/base64x/internal/rt`
)

var F_b64decode func(out unsafe.Pointer, src unsafe.Pointer, len int, mod int) (ret int)

var S_b64decode uintptr

//go:nosplit
func B64decode(out *[]byte, src unsafe.Pointer, len int, mode int) (ret int) {
    return F_b64decode(rt.NoEscape(unsafe.Pointer(out)), rt.NoEscape(unsafe.Pointer(src)), len, mode)
}

