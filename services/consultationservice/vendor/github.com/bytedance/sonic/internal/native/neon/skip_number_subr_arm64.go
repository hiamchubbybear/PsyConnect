// +build !noasm !appengine


package neon

//go:nosplit
//go:noescape

func __skip_number_entry__() uintptr

var (
    _subr__skip_number uintptr = __skip_number_entry__() + 32
)

const (
    _stack__skip_number = 48
)

var (
    _ = _subr__skip_number
)

const (
    _ = _stack__skip_number
)
