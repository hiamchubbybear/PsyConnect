



package neon

//go:nosplit
func validate_utf8_fast(s *string)  (ret int) {
    return __validate_utf8_fast(s)
}

//go:nosplit
//go:noescape

func __validate_utf8_fast(s *string)  (ret int)
