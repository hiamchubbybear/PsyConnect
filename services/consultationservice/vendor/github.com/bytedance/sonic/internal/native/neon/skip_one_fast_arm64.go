



package neon

//go:nosplit
func skip_one_fast(s *string, p *int) (ret int) {
    return __skip_one_fast(s, p)
}

//go:nosplit
//go:noescape

func __skip_one_fast(s *string, p *int) (ret int)
