



package neon

//go:nosplit
func i64toa(out *byte, val int64) (ret int) {
    return __i64toa(out, val)
}

//go:nosplit
//go:noescape

func __i64toa(out *byte, val int64) (ret int)
