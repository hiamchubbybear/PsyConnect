// +build !noasm !appengine


package neon

//go:nosplit
//go:noescape

func __vstring_entry__() uintptr

var (
    _subr__vstring uintptr = __vstring_entry__() + 32
)

const (
    _stack__vstring = 48
)

var (
    _ = _subr__vstring
)

const (
    _ = _stack__vstring
)
