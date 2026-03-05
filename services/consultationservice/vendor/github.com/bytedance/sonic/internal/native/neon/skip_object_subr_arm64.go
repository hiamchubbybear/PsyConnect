// +build !noasm !appengine


package neon

//go:nosplit
//go:noescape

func __skip_object_entry__() uintptr

var (
    _subr__skip_object uintptr = __skip_object_entry__() + 48
)

const (
    _stack__skip_object = 240
)

var (
    _ = _subr__skip_object
)

const (
    _ = _stack__skip_object
)
