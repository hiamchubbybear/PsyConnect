



package neon

//go:nosplit
func f64toa(out *byte, val float64) (ret int) {
    return __f64toa(out, val)
}

//go:nosplit
//go:noescape

func __f64toa(out *byte, val float64) (ret int) 
