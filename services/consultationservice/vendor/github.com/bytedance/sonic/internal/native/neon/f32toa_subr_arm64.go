// +build !noasm !appengine


package neon

//go:nosplit
//go:noescape

func __f32toa_entry__() uintptr

var (
    _subr__f32toa uintptr = __f32toa_entry__() + 0
)

const (
    _stack__f32toa = 32
)

var (
    _ = _subr__f32toa
)

const (
    _ = _stack__f32toa
)
