// +build !noasm !appengine


package neon

//go:nosplit
//go:noescape

func __get_by_path_entry__() uintptr

var (
    _subr__get_by_path uintptr = __get_by_path_entry__() + 48
)

const (
    _stack__get_by_path = 208
)

var (
    _ = _subr__get_by_path
)

const (
    _ = _stack__get_by_path
)
