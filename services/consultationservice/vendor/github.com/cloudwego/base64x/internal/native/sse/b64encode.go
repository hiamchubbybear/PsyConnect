



package sse

import (
    `unsafe`

    `github.com/cloudwego/base64x/internal/rt`
)

var F_b64encode func(out unsafe.Pointer, src unsafe.Pointer, mod int)

var S_b64encode uintptr

//go:nosplit
func B64encode(out *[]byte, src *[]byte, mode int) {
    F_b64encode(rt.NoEscape(unsafe.Pointer(out)), rt.NoEscape(unsafe.Pointer(src)), mode)
}
