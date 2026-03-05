// +build !noasm !appengine


package neon

//go:nosplit
//go:noescape

func __vunsigned_entry__() uintptr

var (
    _subr__vunsigned uintptr = __vunsigned_entry__() + 0
)

const (
    _stack__vunsigned = 32
)

var (
    _ = _subr__vunsigned
)

const (
    _ = _stack__vunsigned
)
