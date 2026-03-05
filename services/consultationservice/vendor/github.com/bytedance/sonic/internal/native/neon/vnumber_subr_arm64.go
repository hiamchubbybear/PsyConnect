// +build !noasm !appengine


package neon

//go:nosplit
//go:noescape

func __vnumber_entry__() uintptr

var (
    _subr__vnumber uintptr = __vnumber_entry__() + 0
)

const (
    _stack__vnumber = 112
)

var (
    _ = _subr__vnumber
)

const (
    _ = _stack__vnumber
)
