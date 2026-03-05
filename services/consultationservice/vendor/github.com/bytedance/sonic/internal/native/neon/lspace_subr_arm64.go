// +build !noasm !appengine


package neon

//go:nosplit
//go:noescape

func __lspace_entry__() uintptr

var (
    _subr__lspace uintptr = __lspace_entry__() + 0
)

const (
    _stack__lspace = 32
)

var (
    _ = _subr__lspace
)

const (
    _ = _stack__lspace
)
