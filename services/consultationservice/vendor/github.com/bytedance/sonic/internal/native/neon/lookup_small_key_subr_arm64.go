// +build !noasm !appengine


package neon

//go:nosplit
//go:noescape

func __lookup_small_key_entry__() uintptr

var (
    _subr__lookup_small_key uintptr = __lookup_small_key_entry__() + 32
)

const (
    _stack__lookup_small_key = 32
)

var (
    _ = _subr__lookup_small_key
)

const (
    _ = _stack__lookup_small_key
)
