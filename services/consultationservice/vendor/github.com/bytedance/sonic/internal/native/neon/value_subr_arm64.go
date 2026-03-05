// +build !noasm !appengine


package neon

//go:nosplit
//go:noescape

func __value_entry__() uintptr

var (
    _subr__value uintptr = __value_entry__() + 32
)

const (
    _stack__value = 112
)

var (
    _ = _subr__value
)

const (
    _ = _stack__value
)
