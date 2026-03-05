



package neon

import (
    `unsafe`

    
)

//go:nosplit
func html_escape(sp unsafe.Pointer, nb int, dp unsafe.Pointer, dn *int) (ret int) {
    return __html_escape(sp, nb, dp, dn)
}

//go:nosplit
//go:noescape

func __html_escape(sp unsafe.Pointer, nb int, dp unsafe.Pointer, dn *int) (ret int)
