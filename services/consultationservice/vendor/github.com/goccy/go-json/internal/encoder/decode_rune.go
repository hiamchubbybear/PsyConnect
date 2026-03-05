package encoder

import "unicode/utf8"

const (
	
	locb = 128 
	hicb = 191 

	
	
	
	
	xx = 0xF1 
	as = 0xF0 
	s1 = 0x02 
	s2 = 0x13 
	s3 = 0x03 
	s4 = 0x23 
	s5 = 0x34 
	s6 = 0x04 
	s7 = 0x44 
)


var first = [256]uint8{
	
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, 
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, 
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, 
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, 
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, 
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, 
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, 
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, 
	
	xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, 
	xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, 
	xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, 
	xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, 
	xx, xx, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, 
	s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, 
	s2, s3, s3, s3, s3, s3, s3, s3, s3, s3, s3, s3, s3, s4, s3, s3, 
	s5, s6, s6, s6, s7, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, 
}

const (
	lineSep      = byte(168) 
	paragraphSep = byte(169) 
)

type decodeRuneState int

const (
	validUTF8State decodeRuneState = iota
	runeErrorState
	lineSepState
	paragraphSepState
)

func decodeRuneInString(s string) (decodeRuneState, int) {
	n := len(s)
	s0 := s[0]
	x := first[s0]
	if x >= as {
		
		
		
		mask := rune(x) << 31 >> 31 
		if rune(s[0])&^mask|utf8.RuneError&mask == utf8.RuneError {
			return runeErrorState, 1
		}
		return validUTF8State, 1
	}
	sz := int(x & 7)
	if n < sz {
		return runeErrorState, 1
	}
	s1 := s[1]
	switch x >> 4 {
	case 0:
		if s1 < locb || hicb < s1 {
			return runeErrorState, 1
		}
	case 1:
		if s1 < 0xA0 || hicb < s1 {
			return runeErrorState, 1
		}
	case 2:
		if s1 < locb || 0x9F < s1 {
			return runeErrorState, 1
		}
	case 3:
		if s1 < 0x90 || hicb < s1 {
			return runeErrorState, 1
		}
	case 4:
		if s1 < locb || 0x8F < s1 {
			return runeErrorState, 1
		}
	}
	if sz <= 2 {
		return validUTF8State, 2
	}
	s2 := s[2]
	if s2 < locb || hicb < s2 {
		return runeErrorState, 1
	}
	if sz <= 3 {
		
		if s0 == 226 && s1 == 128 {
			switch s2 {
			case lineSep:
				return lineSepState, 3
			case paragraphSep:
				return paragraphSepState, 3
			}
		}
		return validUTF8State, 3
	}
	s3 := s[3]
	if s3 < locb || hicb < s3 {
		return runeErrorState, 1
	}
	return validUTF8State, 4
}
