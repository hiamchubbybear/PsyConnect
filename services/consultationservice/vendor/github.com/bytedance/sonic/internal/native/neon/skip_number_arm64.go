



package neon

//go:nosplit
func skip_number(s *string, p *int) (ret int) {
    return __skip_number(s, p)
}

//go:nosplit
//go:noescape

func __skip_number(s *string, p *int) (ret int)
