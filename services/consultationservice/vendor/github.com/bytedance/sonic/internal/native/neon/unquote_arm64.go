



package neon

import (
    `unsafe`

    
)

//go:nosplit
func unquote(sp unsafe.Pointer, nb int, dp unsafe.Pointer, ep *int, flags uint64) (ret int) {
    return __unquote(sp, nb, dp, ep, flags)
}

//go:nosplit
//go:noescape

func __unquote(sp unsafe.Pointer, nb int, dp unsafe.Pointer, ep *int, flags uint64) (ret int)
