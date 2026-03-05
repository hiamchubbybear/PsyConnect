// +build !noasm !appengine


package neon

//go:nosplit
//go:noescape

func __parse_with_padding_entry__() uintptr

var (
    _subr__parse_with_padding uintptr = __parse_with_padding_entry__() + 248
)

const (
    _stack__parse_with_padding = 160
)

var (
    _ = _subr__parse_with_padding
)

const (
    _ = _stack__parse_with_padding
)
