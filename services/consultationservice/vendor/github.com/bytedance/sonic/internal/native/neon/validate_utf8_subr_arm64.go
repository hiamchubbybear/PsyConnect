// +build !noasm !appengine


package neon

//go:nosplit
//go:noescape

func __validate_utf8_entry__() uintptr

var (
    _subr__validate_utf8 uintptr = __validate_utf8_entry__() + 0
)

const (
    _stack__validate_utf8 = 64
)

var (
    _ = _subr__validate_utf8
)

const (
    _ = _stack__validate_utf8
)
