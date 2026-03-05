




package s2

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
)

var (
	
	ErrCorrupt = errors.New("s2: corrupt input")
	
	ErrCRC = errors.New("s2: corrupt input, crc mismatch")
	
	ErrTooLarge = errors.New("s2: decoded block is too large")
	
	ErrUnsupported = errors.New("s2: unsupported input")
)


func DecodedLen(src []byte) (int, error) {
	v, _, err := decodedLen(src)
	return v, err
}



func decodedLen(src []byte) (blockLen, headerLen int, err error) {
	v, n := binary.Uvarint(src)
	if n <= 0 || v > 0xffffffff {
		return 0, 0, ErrCorrupt
	}

	const wordSize = 32 << (^uint(0) >> 32 & 1)
	if wordSize == 32 && v > 0x7fffffff {
		return 0, 0, ErrTooLarge
	}
	return int(v), n, nil
}

const (
	decodeErrCodeCorrupt = 1
)






func Decode(dst, src []byte) ([]byte, error) {
	dLen, s, err := decodedLen(src)
	if err != nil {
		return nil, err
	}
	if dLen <= cap(dst) {
		dst = dst[:dLen]
	} else {
		dst = make([]byte, dLen)
	}
	if s2Decode(dst, src[s:]) != 0 {
		return nil, ErrCorrupt
	}
	return dst, nil
}






func s2DecodeDict(dst, src []byte, dict *Dict) int {
	if dict == nil {
		return s2Decode(dst, src)
	}
	const debug = false
	const debugErrs = debug

	if debug {
		fmt.Println("Starting decode, dst len:", len(dst))
	}
	var d, s, length int
	offset := len(dict.dict) - dict.repeat

	
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
			if debug {
				fmt.Println("literals, length:", length, "d-after:", d+length)
			}
			if length > len(dst)-d || length > len(src)-s || (strconv.IntSize == 32 && length <= 0) {
				if debugErrs {
					fmt.Println("corrupt literal: length:", length, "d-left:", len(dst)-d, "src-left:", len(src)-s)
				}
				return decodeErrCodeCorrupt
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

		if offset <= 0 || length > len(dst)-d {
			if debugErrs {
				fmt.Println("match error; offset:", offset, "length:", length, "dst-left:", len(dst)-d)
			}
			return decodeErrCodeCorrupt
		}

		
		if d < offset {
			if d > MaxDictSrcOffset {
				if debugErrs {
					fmt.Println("dict after", MaxDictSrcOffset, "d:", d, "offset:", offset, "length:", length)
				}
				return decodeErrCodeCorrupt
			}
			startOff := len(dict.dict) - offset + d
			if startOff < 0 || startOff+length > len(dict.dict) {
				if debugErrs {
					fmt.Printf("offset (%d) + length (%d) bigger than dict (%d)\n", offset, length, len(dict.dict))
				}
				return decodeErrCodeCorrupt
			}
			if debug {
				fmt.Println("dict copy, length:", length, "offset:", offset, "d-after:", d+length, "dict start offset:", startOff)
			}
			copy(dst[d:d+length], dict.dict[startOff:])
			d += length
			continue
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
					if debugErrs {
						fmt.Println("src went oob")
					}
					return decodeErrCodeCorrupt
				}
				x = uint32(src[s-1])
			case x == 61:
				s += 3
				if uint(s) > uint(len(src)) { 
					if debugErrs {
						fmt.Println("src went oob")
					}
					return decodeErrCodeCorrupt
				}
				x = uint32(src[s-2]) | uint32(src[s-1])<<8
			case x == 62:
				s += 4
				if uint(s) > uint(len(src)) { 
					if debugErrs {
						fmt.Println("src went oob")
					}
					return decodeErrCodeCorrupt
				}
				x = uint32(src[s-3]) | uint32(src[s-2])<<8 | uint32(src[s-1])<<16
			case x == 63:
				s += 5
				if uint(s) > uint(len(src)) { 
					if debugErrs {
						fmt.Println("src went oob")
					}
					return decodeErrCodeCorrupt
				}
				x = uint32(src[s-4]) | uint32(src[s-3])<<8 | uint32(src[s-2])<<16 | uint32(src[s-1])<<24
			}
			length = int(x) + 1
			if length > len(dst)-d || length > len(src)-s || (strconv.IntSize == 32 && length <= 0) {
				if debugErrs {
					fmt.Println("corrupt literal: length:", length, "d-left:", len(dst)-d, "src-left:", len(src)-s)
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
				if debugErrs {
					fmt.Println("src went oob")
				}
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
						if debugErrs {
							fmt.Println("src went oob")
						}
						return decodeErrCodeCorrupt
					}
					length = int(uint32(src[s-1])) + 4
				case 6:
					s += 2
					if uint(s) > uint(len(src)) { 
						if debugErrs {
							fmt.Println("src went oob")
						}
						return decodeErrCodeCorrupt
					}
					length = int(uint32(src[s-2])|(uint32(src[s-1])<<8)) + (1 << 8)
				case 7:
					s += 3
					if uint(s) > uint(len(src)) { 
						if debugErrs {
							fmt.Println("src went oob")
						}
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
				if debugErrs {
					fmt.Println("src went oob")
				}
				return decodeErrCodeCorrupt
			}
			length = 1 + int(src[s-3])>>2
			offset = int(uint32(src[s-2]) | uint32(src[s-1])<<8)

		case tagCopy4:
			s += 5
			if uint(s) > uint(len(src)) { 
				if debugErrs {
					fmt.Println("src went oob")
				}
				return decodeErrCodeCorrupt
			}
			length = 1 + int(src[s-5])>>2
			offset = int(uint32(src[s-4]) | uint32(src[s-3])<<8 | uint32(src[s-2])<<16 | uint32(src[s-1])<<24)
		}

		if offset <= 0 || length > len(dst)-d {
			if debugErrs {
				fmt.Println("match error; offset:", offset, "length:", length, "dst-left:", len(dst)-d)
			}
			return decodeErrCodeCorrupt
		}

		
		if d < offset {
			if d > MaxDictSrcOffset {
				if debugErrs {
					fmt.Println("dict after", MaxDictSrcOffset, "d:", d, "offset:", offset, "length:", length)
				}
				return decodeErrCodeCorrupt
			}
			rOff := len(dict.dict) - (offset - d)
			if debug {
				fmt.Println("starting dict entry from dict offset", len(dict.dict)-rOff)
			}
			if rOff+length > len(dict.dict) {
				if debugErrs {
					fmt.Println("err: END offset", rOff+length, "bigger than dict", len(dict.dict), "dict offset:", rOff, "length:", length)
				}
				return decodeErrCodeCorrupt
			}
			if rOff < 0 {
				if debugErrs {
					fmt.Println("err: START offset", rOff, "less than 0", len(dict.dict), "dict offset:", rOff, "length:", length)
				}
				return decodeErrCodeCorrupt
			}
			copy(dst[d:d+length], dict.dict[rOff:])
			d += length
			continue
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
		if debugErrs {
			fmt.Println("wanted length", len(dst), "got", d)
		}
		return decodeErrCodeCorrupt
	}
	return 0
}
