



package neon

import (
    `unsafe`

    
)

//go:nosplit
func quote(sp unsafe.Pointer, nb int, dp unsafe.Pointer, dn *int, flags uint64) (ret int) {
    return __quote(sp, nb, dp, dn, flags)
}

//go:nosplit
//go:noescape

func __quote(sp unsafe.Pointer, nb int, dp unsafe.Pointer, dn *int, flags uint64) (ret int)
