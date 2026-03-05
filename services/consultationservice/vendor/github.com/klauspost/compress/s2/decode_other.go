




//go:build (!amd64 && !arm64) || appengine || !gc || noasm
// +build !amd64,!arm64 appengine !gc noasm

package s2

import (
	"fmt"
	"strconv"
)






func s2Decode(dst, src []byte) int {
	const debug = false
	if debug {
		fmt.Println("Starting decode, dst len:", len(dst))
	}
	var d, s, length int
	offset := 0

	
	for s < len(src)-5 {
		
		
		
		switch src[s] & 0x03 {
		case tagLiteral:
			x := uint32(src[s] >> 2)
			switch {
			case x < 60:
				s++
			case x == 60:
				s += 2
				x = uint32(src[s-1])
			case x == 61:
				in := src[s : s+3]
				x = uint32(in[1]) | uint32(in[2])<<8
				s += 3
			case x == 62:
				in := src[s : s+4]
				
				x = uint32(in[0]) | uint32(in[1])<<8 | uint32(in[2])<<16 | uint32(in[3])<<24
				x >>= 8
				s += 4
			case x == 63:
				in := src[s : s+5]
				x = uint32(in[1]) | uint32(in[2])<<8 | uint32(in[3])<<16 | uint32(in[4])<<24
				s += 5
			}
			length = int(x) + 1
			if length > len(dst)-d || length > len(src)-s || (strconv.IntSize == 32 && length <= 0) {
				if debug {
					fmt.Println("corrupt: lit size", length)
				}
				return decodeErrCodeCorrupt
			}
			if debug {
				fmt.Println("literals, length:", length, "d-after:", d+length)
			}

			copy(dst[d:], src[s:s+length])
			d += length
			s += length
			continue

		case tagCopy1:
			s += 2
			toffset := int(uint32(src[s-2])&0xe0<<3 | uint32(src[s-1]))
			length = int(src[s-2]) >> 2 & 0x7
			if toffset == 0 {
				if debug {
					fmt.Print("(repeat) ")
				}
				
				switch length {
				case 5:
					length = int(src[s]) + 4
					s += 1
				case 6:
					in := src[s : s+2]
					length = int(uint32(in[0])|(uint32(in[1])<<8)) + (1 << 8)
					s += 2
				case 7:
					in := src[s : s+3]
					length = int((uint32(in[2])<<16)|(uint32(in[1])<<8)|uint32(in[0])) + (1 << 16)
					s += 3
				default: 
				}
			} else {
				offset = toffset
			}
			length += 4
		case tagCopy2:
			in := src[s : s+3]
			offset = int(uint32(in[1]) | uint32(in[2])<<8)
			length = 1 + int(in[0])>>2
			s += 3

		case tagCopy4:
			in := src[s : s+5]
			offset = int(uint32(in[1]) | uint32(in[2])<<8 | uint32(in[3])<<16 | uint32(in[4])<<24)
			length = 1 + int(in[0])>>2
			s += 5
		}

		if offset <= 0 || d < offset || length > len(dst)-d {
			if debug {
				fmt.Println("corrupt: match, length", length, "offset:", offset, "dst avail:", len(dst)-d, "dst pos:", d)
			}

			return decodeErrCodeCorrupt
		}

		if debug {
			fmt.Println("copy, length:", length, "offset:", offset, "d-after:", d+length)
		}

		
		
		if offset > length {
			copy(dst[d:d+length], dst[d-offset:])
			d += length
			continue
		}

		
		
		
		
		
		
		
		a := dst[d : d+length]
		b := dst[d-offset:]
		b = b[:len(a)]
		for i := range a {
			a[i] = b[i]
		}
		d += length
	}

	
	for s < len(src) {
		switch src[s] & 0x03 {
		case tagLiteral:
			x := uint32(src[s] >> 2)
			switch {
			case x < 60:
				s++
			case x == 60:
				s += 2
				if uint(s) > uint(len(src)) { 
					return decodeErrCodeCorrupt
				}
				x = uint32(src[s-1])
			case x == 61:
				s += 3
				if uint(s) > uint(len(src)) { 
					return decodeErrCodeCorrupt
				}
				x = uint32(src[s-2]) | uint32(src[s-1])<<8
			case x == 62:
				s += 4
				if uint(s) > uint(len(src)) { 
					return decodeErrCodeCorrupt
				}
				x = uint32(src[s-3]) | uint32(src[s-2])<<8 | uint32(src[s-1])<<16
			case x == 63:
				s += 5
				if uint(s) > uint(len(src)) { 
					return decodeErrCodeCorrupt
				}
				x = uint32(src[s-4]) | uint32(src[s-3])<<8 | uint32(src[s-2])<<16 | uint32(src[s-1])<<24
			}
			length = int(x) + 1
			if length > len(dst)-d || length > len(src)-s || (strconv.IntSize == 32 && length <= 0) {
				if debug {
					fmt.Println("corrupt: lit size", length)
				}
				return decodeErrCodeCorrupt
			}
			if debug {
				fmt.Println("literals, length:", length, "d-after:", d+length)
			}

			copy(dst[d:], src[s:s+length])
			d += length
			s += length
			continue

		case tagCopy1:
			s += 2
			if uint(s) > uint(len(src)) { 
				return decodeErrCodeCorrupt
			}
			length = int(src[s-2]) >> 2 & 0x7
			toffset := int(uint32(src[s-2])&0xe0<<3 | uint32(src[s-1]))
			if toffset == 0 {
				if debug {
					fmt.Print("(repeat) ")
				}
				
				switch length {
				case 5:
					s += 1
					if uint(s) > uint(len(src)) { 
						return decodeErrCodeCorrupt
					}
					length = int(uint32(src[s-1])) + 4
				case 6:
					s += 2
					if uint(s) > uint(len(src)) { 
						return decodeErrCodeCorrupt
					}
					length = int(uint32(src[s-2])|(uint32(src[s-1])<<8)) + (1 << 8)
				case 7:
					s += 3
					if uint(s) > uint(len(src)) { 
						return decodeErrCodeCorrupt
					}
					length = int(uint32(src[s-3])|(uint32(src[s-2])<<8)|(uint32(src[s-1])<<16)) + (1 << 16)
				default: 
				}
			} else {
				offset = toffset
			}
			length += 4
		case tagCopy2:
			s += 3
			if uint(s) > uint(len(src)) { 
				return decodeErrCodeCorrupt
			}
			length = 1 + int(src[s-3])>>2
			offset = int(uint32(src[s-2]) | uint32(src[s-1])<<8)

		case tagCopy4:
			s += 5
			if uint(s) > uint(len(src)) { 
				return decodeErrCodeCorrupt
			}
			length = 1 + int(src[s-5])>>2
			offset = int(uint32(src[s-4]) | uint32(src[s-3])<<8 | uint32(src[s-2])<<16 | uint32(src[s-1])<<24)
		}

		if offset <= 0 || d < offset || length > len(dst)-d {
			if debug {
				fmt.Println("corrupt: match, length", length, "offset:", offset, "dst avail:", len(dst)-d, "dst pos:", d)
			}
			return decodeErrCodeCorrupt
		}

		if debug {
			fmt.Println("copy, length:", length, "offset:", offset, "d-after:", d+length)
		}

		
		
		if offset > length {
			copy(dst[d:d+length], dst[d-offset:])
			d += length
			continue
		}

		
		
		
		
		
		
		
		a := dst[d : d+length]
		b := dst[d-offset:]
		b = b[:len(a)]
		for i := range a {
			a[i] = b[i]
		}
		d += length
	}

	if d != len(dst) {
		return decodeErrCodeCorrupt
	}
	return 0
}
