



package neon

import (
    `unsafe`

    
)

//go:nosplit
func lspace(sp unsafe.Pointer, nb int, off int) (ret int) {
    return __lspace(sp, nb, off)
}

//go:nosplit
//go:noescape

func __lspace(sp unsafe.Pointer, nb int, off int) (ret int)
