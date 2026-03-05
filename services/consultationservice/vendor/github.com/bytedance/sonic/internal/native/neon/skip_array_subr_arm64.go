// +build !noasm !appengine


package neon

//go:nosplit
//go:noescape

func __skip_array_entry__() uintptr

var (
    _subr__skip_array uintptr = __skip_array_entry__() + 48
)

const (
    _stack__skip_array = 240
)

var (
    _ = _subr__skip_array
)

const (
    _ = _stack__skip_array
)
