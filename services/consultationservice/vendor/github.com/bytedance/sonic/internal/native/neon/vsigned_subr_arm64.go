// +build !noasm !appengine


package neon

//go:nosplit
//go:noescape

func __vsigned_entry__() uintptr

var (
    _subr__vsigned uintptr = __vsigned_entry__() + 0
)

const (
    _stack__vsigned = 32
)

var (
    _ = _subr__vsigned
)

const (
    _ = _stack__vsigned
)
