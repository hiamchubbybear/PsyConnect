



package json

import (
	"bytes"
	"strconv"
)





func parseNumber(input []byte) (int, bool) {
	var n int

	s := input
	if len(s) == 0 {
		return 0, false
	}

	
	if s[0] == '-' {
		s = s[1:]
		n++
		if len(s) == 0 {
			return 0, false
		}
	}

	
	switch {
	case s[0] == '0':
		s = s[1:]
		n++

	case '1' <= s[0] && s[0] <= '9':
		s = s[1:]
		n++
		for len(s) > 0 && '0' <= s[0] && s[0] <= '9' {
			s = s[1:]
			n++
		}

	default:
		return 0, false
	}

	
	if len(s) >= 2 && s[0] == '.' && '0' <= s[1] && s[1] <= '9' {
		s = s[2:]
		n += 2
		for len(s) > 0 && '0' <= s[0] && s[0] <= '9' {
			s = s[1:]
			n++
		}
	}

	
	
	if len(s) >= 2 && (s[0] == 'e' || s[0] == 'E') {
		s = s[1:]
		n++
		if s[0] == '+' || s[0] == '-' {
			s = s[1:]
			n++
			if len(s) == 0 {
				return 0, false
			}
		}
		for len(s) > 0 && '0' <= s[0] && s[0] <= '9' {
			s = s[1:]
			n++
		}
	}

	
	if n < len(input) && isNotDelim(input[n]) {
		return 0, false
	}

	return n, true
}



type numberParts struct {
	neg  bool
	intp []byte
	frac []byte
	exp  []byte
}




func parseNumberParts(input []byte) (numberParts, bool) {
	var neg bool
	var intp []byte
	var frac []byte
	var exp []byte

	s := input
	if len(s) == 0 {
		return numberParts{}, false
	}

	
	if s[0] == '-' {
		neg = true
		s = s[1:]
		if len(s) == 0 {
			return numberParts{}, false
		}
	}

	
	switch {
	case s[0] == '0':
		
		s = s[1:]

	case '1' <= s[0] && s[0] <= '9':
		intp = s
		n := 1
		s = s[1:]
		for len(s) > 0 && '0' <= s[0] && s[0] <= '9' {
			s = s[1:]
			n++
		}
		intp = intp[:n]

	default:
		return numberParts{}, false
	}

	
	if len(s) >= 2 && s[0] == '.' && '0' <= s[1] && s[1] <= '9' {
		frac = s[1:]
		n := 1
		s = s[2:]
		for len(s) > 0 && '0' <= s[0] && s[0] <= '9' {
			s = s[1:]
			n++
		}
		frac = frac[:n]
	}

	
	
	if len(s) >= 2 && (s[0] == 'e' || s[0] == 'E') {
		s = s[1:]
		exp = s
		n := 0
		if s[0] == '+' || s[0] == '-' {
			s = s[1:]
			n++
			if len(s) == 0 {
				return numberParts{}, false
			}
		}
		for len(s) > 0 && '0' <= s[0] && s[0] <= '9' {
			s = s[1:]
			n++
		}
		exp = exp[:n]
	}

	return numberParts{
		neg:  neg,
		intp: intp,
		frac: bytes.TrimRight(frac, "0"), 
		exp:  exp,
	}, true
}




func normalizeToIntString(n numberParts) (string, bool) {
	intpSize := len(n.intp)
	fracSize := len(n.frac)

	if intpSize == 0 && fracSize == 0 {
		return "0", true
	}

	var exp int
	if len(n.exp) > 0 {
		i, err := strconv.ParseInt(string(n.exp), 10, 32)
		if err != nil {
			return "", false
		}
		exp = int(i)
	}

	var num []byte
	if exp >= 0 {
		
		

		
		
		if fracSize > exp {
			return "", false
		}

		
		
		
		const maxDigits = 20 
		if intpSize+exp > maxDigits {
			return "", false
		}

		
		num = n.intp[:len(n.intp):len(n.intp)]
		num = append(num, n.frac...)
		for i := 0; i < exp-fracSize; i++ {
			num = append(num, '0')
		}
	} else {
		

		
		if fracSize > 0 {
			return "", false
		}

		
		
		index := intpSize + exp
		if index < 0 {
			return "", false
		}

		num = n.intp
		
		
		for i := index; i < intpSize; i++ {
			if num[i] != '0' {
				return "", false
			}
		}
		num = num[:index]
	}

	if n.neg {
		return "-" + string(num), true
	}
	return string(num), true
}
