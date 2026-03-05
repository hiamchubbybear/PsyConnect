// +build !noasm !appengine


package neon

//go:nosplit
//go:noescape

func __u64toa_entry__() uintptr

var (
    _subr__u64toa uintptr = __u64toa_entry__() + 48
)

const (
    _stack__u64toa = 32
)

var (
    _ = _subr__u64toa
)

const (
    _ = _stack__u64toa
)
