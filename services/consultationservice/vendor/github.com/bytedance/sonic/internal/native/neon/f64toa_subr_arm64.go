// +build !noasm !appengine


package neon

//go:nosplit
//go:noescape

func __f64toa_entry__() uintptr

var (
    _subr__f64toa uintptr = __f64toa_entry__() + 0
)

const (
    _stack__f64toa = 32
)

var (
    _ = _subr__f64toa
)

const (
    _ = _stack__f64toa
)
