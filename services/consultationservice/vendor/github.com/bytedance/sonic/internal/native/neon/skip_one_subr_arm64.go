// +build !noasm !appengine


package neon

//go:nosplit
//go:noescape

func __skip_one_entry__() uintptr

var (
    _subr__skip_one uintptr = __skip_one_entry__() + 48
)

const (
    _stack__skip_one = 192
)

var (
    _ = _subr__skip_one
)

const (
    _ = _stack__skip_one
)
