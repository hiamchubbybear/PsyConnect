



package neon

//go:nosplit
func u64toa(out *byte, val uint64) (ret int) {
    return __u64toa(out, val)
}

//go:nosplit
//go:noescape

func __u64toa(out *byte, val uint64) (ret int)
