// +build !noasm !appengine


package neon

//go:nosplit
//go:noescape

func __unquote_entry__() uintptr

var (
    _subr__unquote uintptr = __unquote_entry__() + 32
)

const (
    _stack__unquote = 112
)

var (
    _ = _subr__unquote
)

const (
    _ = _stack__unquote
)
