



package neon

//go:nosplit
func f32toa(out *byte, val float32) (ret int) {
    return __f32toa(out, val)
}

//go:nosplit
//go:noescape

func __f32toa(out *byte, val float32) (ret int) 
