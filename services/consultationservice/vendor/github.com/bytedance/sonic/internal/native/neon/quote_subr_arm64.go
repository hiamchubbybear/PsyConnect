// +build !noasm !appengine


package neon

//go:nosplit
//go:noescape

func __quote_entry__() uintptr

var (
    _subr__quote uintptr = __quote_entry__() + 32
)

const (
    _stack__quote = 32
)

var (
    _ = _subr__quote
)

const (
    _ = _stack__quote
)
