




package gin














func cleanPath(p string) string {
	const stackBufSize = 128
	
	if p == "" {
		return "/"
	}

	
	
	buf := make([]byte, 0, stackBufSize)

	n := len(p)

	
	
	

	
	r := 1
	w := 1

	if p[0] != '/' {
		r = 0

		if n+1 > stackBufSize {
			buf = make([]byte, n+1)
		} else {
			buf = buf[:n+1]
		}
		buf[0] = '/'
	}

	trailing := n > 1 && p[n-1] == '/'

	
	
	
	

	for r < n {
		switch {
		case p[r] == '/':
			
			r++

		case p[r] == '.' && r+1 == n:
			trailing = true
			r++

		case p[r] == '.' && p[r+1] == '/':
			
			r += 2

		case p[r] == '.' && p[r+1] == '.' && (r+2 == n || p[r+2] == '/'):
			
			r += 3

			if w > 1 {
				
				w--

				if len(buf) == 0 {
					for w > 1 && p[w] != '/' {
						w--
					}
				} else {
					for w > 1 && buf[w] != '/' {
						w--
					}
				}
			}

		default:
			
			
			if w > 1 {
				bufApp(&buf, p, w, '/')
				w++
			}

			
			for r < n && p[r] != '/' {
				bufApp(&buf, p, w, p[r])
				w++
				r++
			}
		}
	}

	
	if trailing && w > 1 {
		bufApp(&buf, p, w, '/')
		w++
	}

	
	
	
	if len(buf) == 0 {
		return p[:w]
	}
	return string(buf[:w])
}



func bufApp(buf *[]byte, s string, w int, c byte) {
	b := *buf
	if len(b) == 0 {
		
		
		
		if s[w] == c {
			return
		}

		
		
		length := len(s)
		if length > cap(b) {
			*buf = make([]byte, length)
		} else {
			*buf = (*buf)[:length]
		}
		b = *buf

		copy(b, s[:w])
	}
	b[w] = c
}
