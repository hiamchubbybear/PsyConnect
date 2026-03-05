package characters

import (
	"unicode/utf8"
)

type utf8Err struct {
	Index int
	Size  int
}

func (u utf8Err) Zero() bool {
	return u.Size == 0
}




















func Utf8TomlValidAlreadyEscaped(p []byte) (err utf8Err) {
	
	offset := 0
	for len(p) >= 8 {
		
		
		
		
		first32 := uint32(p[0]) | uint32(p[1])<<8 | uint32(p[2])<<16 | uint32(p[3])<<24
		second32 := uint32(p[4]) | uint32(p[5])<<8 | uint32(p[6])<<16 | uint32(p[7])<<24
		if (first32|second32)&0x80808080 != 0 {
			
			break
		}

		for i, b := range p[:8] {
			if InvalidAscii(b) {
				err.Index = offset + i
				err.Size = 1
				return
			}
		}

		p = p[8:]
		offset += 8
	}
	n := len(p)
	for i := 0; i < n; {
		pi := p[i]
		if pi < utf8.RuneSelf {
			if InvalidAscii(pi) {
				err.Index = offset + i
				err.Size = 1
				return
			}
			i++
			continue
		}
		x := first[pi]
		if x == xx {
			
			err.Index = offset + i
			err.Size = 1
			return
		}
		size := int(x & 7)
		if i+size > n {
			
			err.Index = offset + i
			err.Size = n - i
			return
		}
		accept := acceptRanges[x>>4]
		if c := p[i+1]; c < accept.lo || accept.hi < c {
			err.Index = offset + i
			err.Size = 2
			return
		} else if size == 2 {
		} else if c := p[i+2]; c < locb || hicb < c {
			err.Index = offset + i
			err.Size = 3
			return
		} else if size == 3 {
		} else if c := p[i+3]; c < locb || hicb < c {
			err.Index = offset + i
			err.Size = 4
			return
		}
		i += size
	}
	return
}


func Utf8ValidNext(p []byte) int {
	c := p[0]

	if c < utf8.RuneSelf {
		if InvalidAscii(c) {
			return 0
		}
		return 1
	}

	x := first[c]
	if x == xx {
		
		return 0
	}
	size := int(x & 7)
	if size > len(p) {
		
		return 0
	}
	accept := acceptRanges[x>>4]
	if c := p[1]; c < accept.lo || accept.hi < c {
		return 0
	} else if size == 2 {
	} else if c := p[2]; c < locb || hicb < c {
		return 0
	} else if size == 3 {
	} else if c := p[3]; c < locb || hicb < c {
		return 0
	}

	return size
}



type acceptRange struct {
	lo uint8 
	hi uint8 
}


var acceptRanges = [16]acceptRange{
	0: {locb, hicb},
	1: {0xA0, hicb},
	2: {locb, 0x9F},
	3: {0x90, hicb},
	4: {locb, 0x8F},
}


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
	
	locb = 0b10000000
	hicb = 0b10111111

	
	
	
	
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
