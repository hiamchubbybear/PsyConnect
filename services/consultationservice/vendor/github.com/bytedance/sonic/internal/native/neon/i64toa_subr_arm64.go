// +build !noasm !appengine


package neon

//go:nosplit
//go:noescape

func __i64toa_entry__() uintptr

var (
    _subr__i64toa uintptr = __i64toa_entry__() + 48
)

const (
    _stack__i64toa = 32
)

var (
    _ = _subr__i64toa
)

const (
    _ = _stack__i64toa
)
