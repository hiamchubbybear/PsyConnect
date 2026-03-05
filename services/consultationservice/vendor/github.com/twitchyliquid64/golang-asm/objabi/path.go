



package objabi

import "strings"






func PathToPrefix(s string) string {
	slash := strings.LastIndex(s, "/")
	
	n := 0
	for r := 0; r < len(s); r++ {
		if c := s[r]; c <= ' ' || (c == '.' && r > slash) || c == '%' || c == '"' || c >= 0x7F {
			n++
		}
	}

	
	if n == 0 {
		return s
	}

	
	const hex = "0123456789abcdef"
	p := make([]byte, 0, len(s)+2*n)
	for r := 0; r < len(s); r++ {
		if c := s[r]; c <= ' ' || (c == '.' && r > slash) || c == '%' || c == '"' || c >= 0x7F {
			p = append(p, '%', hex[c>>4], hex[c&0xF])
		} else {
			p = append(p, c)
		}
	}

	return string(p)
}
