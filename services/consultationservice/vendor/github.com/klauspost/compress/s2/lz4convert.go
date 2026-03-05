



package s2

import (
	"encoding/binary"
	"errors"
	"fmt"
)



type LZ4Converter struct {
}


var ErrDstTooSmall = errors.New("s2: destination too small")





func (l *LZ4Converter) ConvertBlock(dst, src []byte) ([]byte, int, error) {
	if len(src) == 0 {
		return dst, 0, nil
	}
	const debug = false
	const inline = true
	const lz4MinMatch = 4

	s, d := 0, len(dst)
	dst = dst[:cap(dst)]
	if !debug && hasAmd64Asm {
		res, sz := cvtLZ4BlockAsm(dst[d:], src)
		if res < 0 {
			const (
				errCorrupt     = -1
				errDstTooSmall = -2
			)
			switch res {
			case errCorrupt:
				return nil, 0, ErrCorrupt
			case errDstTooSmall:
				return nil, 0, ErrDstTooSmall
			default:
				return nil, 0, fmt.Errorf("unexpected result: %d", res)
			}
		}
		if d+sz > len(dst) {
			return nil, 0, ErrDstTooSmall
		}
		return dst[:d+sz], res, nil
	}

	dLimit := len(dst) - 10
	var lastOffset uint16
	var uncompressed int
	if debug {
		fmt.Printf("convert block start: len(src): %d, len(dst):%d \n", len(src), len(dst))
	}

	for {
		if s >= len(src) {
			return dst[:d], 0, ErrCorrupt
		}
		
		token := src[s]
		ll := int(token >> 4)
		ml := int(lz4MinMatch + (token & 0xf))

		
		if token >= 0xf0 {
			for {
				s++
				if s >= len(src) {
					if debug {
						fmt.Printf("error reading ll: s (%d) >= len(src) (%d)\n", s, len(src))
					}
					return dst[:d], 0, ErrCorrupt
				}
				val := src[s]
				ll += int(val)
				if val != 255 {
					break
				}
			}
		}
		
		if s+ll >= len(src) {
			if debug {
				fmt.Printf("error literals: s+ll (%d+%d) >= len(src) (%d)\n", s, ll, len(src))
			}
			return nil, 0, ErrCorrupt
		}
		s++
		if ll > 0 {
			if d+ll > dLimit {
				return nil, 0, ErrDstTooSmall
			}
			if debug {
				fmt.Printf("emit %d literals\n", ll)
			}
			d += emitLiteralGo(dst[d:], src[s:s+ll])
			s += ll
			uncompressed += ll
		}

		
		if s == len(src) && ml == lz4MinMatch {
			break
		}
		
		if s >= len(src)-2 {
			if debug {
				fmt.Printf("s (%d) >= len(src)-2 (%d)", s, len(src)-2)
			}
			return nil, 0, ErrCorrupt
		}
		offset := binary.LittleEndian.Uint16(src[s:])
		s += 2
		if offset == 0 {
			if debug {
				fmt.Printf("error: offset 0, ml: %d, len(src)-s: %d\n", ml, len(src)-s)
			}
			return nil, 0, ErrCorrupt
		}
		if int(offset) > uncompressed {
			if debug {
				fmt.Printf("error: offset (%d)> uncompressed (%d)\n", offset, uncompressed)
			}
			return nil, 0, ErrCorrupt
		}

		if ml == lz4MinMatch+15 {
			for {
				if s >= len(src) {
					if debug {
						fmt.Printf("error reading ml: s (%d) >= len(src) (%d)\n", s, len(src))
					}
					return nil, 0, ErrCorrupt
				}
				val := src[s]
				s++
				ml += int(val)
				if val != 255 {
					if s >= len(src) {
						if debug {
							fmt.Printf("error reading ml: s (%d) >= len(src) (%d)\n", s, len(src))
						}
						return nil, 0, ErrCorrupt
					}
					break
				}
			}
		}
		if offset == lastOffset {
			if debug {
				fmt.Printf("emit repeat, length: %d, offset: %d\n", ml, offset)
			}
			if !inline {
				d += emitRepeat16(dst[d:], offset, ml)
			} else {
				length := ml
				dst := dst[d:]
				for len(dst) > 5 {
					
					length -= 4
					if length <= 4 {
						dst[0] = uint8(length)<<2 | tagCopy1
						dst[1] = 0
						d += 2
						break
					}
					if length < 8 && offset < 2048 {
						
						dst[1] = uint8(offset)
						dst[0] = uint8(offset>>8)<<5 | uint8(length)<<2 | tagCopy1
						d += 2
						break
					}
					if length < (1<<8)+4 {
						length -= 4
						dst[2] = uint8(length)
						dst[1] = 0
						dst[0] = 5<<2 | tagCopy1
						d += 3
						break
					}
					if length < (1<<16)+(1<<8) {
						length -= 1 << 8
						dst[3] = uint8(length >> 8)
						dst[2] = uint8(length >> 0)
						dst[1] = 0
						dst[0] = 6<<2 | tagCopy1
						d += 4
						break
					}
					const maxRepeat = (1 << 24) - 1
					length -= 1 << 16
					left := 0
					if length > maxRepeat {
						left = length - maxRepeat + 4
						length = maxRepeat - 4
					}
					dst[4] = uint8(length >> 16)
					dst[3] = uint8(length >> 8)
					dst[2] = uint8(length >> 0)
					dst[1] = 0
					dst[0] = 7<<2 | tagCopy1
					if left > 0 {
						d += 5 + emitRepeat16(dst[5:], offset, left)
						break
					}
					d += 5
					break
				}
			}
		} else {
			if debug {
				fmt.Printf("emit copy, length: %d, offset: %d\n", ml, offset)
			}
			if !inline {
				d += emitCopy16(dst[d:], offset, ml)
			} else {
				length := ml
				dst := dst[d:]
				for len(dst) > 5 {
					
					if length > 64 {
						off := 3
						if offset < 2048 {
							
							dst[1] = uint8(offset)
							dst[0] = uint8(offset>>8)<<5 | uint8(8-4)<<2 | tagCopy1
							length -= 8
							off = 2
						} else {
							
							
							dst[2] = uint8(offset >> 8)
							dst[1] = uint8(offset)
							dst[0] = 59<<2 | tagCopy2
							length -= 60
						}
						
						d += off + emitRepeat16(dst[off:], offset, length)
						break
					}
					if length >= 12 || offset >= 2048 {
						
						dst[2] = uint8(offset >> 8)
						dst[1] = uint8(offset)
						dst[0] = uint8(length-1)<<2 | tagCopy2
						d += 3
						break
					}
					
					dst[1] = uint8(offset)
					dst[0] = uint8(offset>>8)<<5 | uint8(length-4)<<2 | tagCopy1
					d += 2
					break
				}
			}
			lastOffset = offset
		}
		uncompressed += ml
		if d > dLimit {
			return nil, 0, ErrDstTooSmall
		}
	}

	return dst[:d], uncompressed, nil
}





func (l *LZ4Converter) ConvertBlockSnappy(dst, src []byte) ([]byte, int, error) {
	if len(src) == 0 {
		return dst, 0, nil
	}
	const debug = false
	const lz4MinMatch = 4

	s, d := 0, len(dst)
	dst = dst[:cap(dst)]
	
	if !debug && hasAmd64Asm {
		res, sz := cvtLZ4BlockSnappyAsm(dst[d:], src)
		if res < 0 {
			const (
				errCorrupt     = -1
				errDstTooSmall = -2
			)
			switch res {
			case errCorrupt:
				return nil, 0, ErrCorrupt
			case errDstTooSmall:
				return nil, 0, ErrDstTooSmall
			default:
				return nil, 0, fmt.Errorf("unexpected result: %d", res)
			}
		}
		if d+sz > len(dst) {
			return nil, 0, ErrDstTooSmall
		}
		return dst[:d+sz], res, nil
	}

	dLimit := len(dst) - 10
	var uncompressed int
	if debug {
		fmt.Printf("convert block start: len(src): %d, len(dst):%d \n", len(src), len(dst))
	}

	for {
		if s >= len(src) {
			return nil, 0, ErrCorrupt
		}
		
		token := src[s]
		ll := int(token >> 4)
		ml := int(lz4MinMatch + (token & 0xf))

		
		if token >= 0xf0 {
			for {
				s++
				if s >= len(src) {
					if debug {
						fmt.Printf("error reading ll: s (%d) >= len(src) (%d)\n", s, len(src))
					}
					return nil, 0, ErrCorrupt
				}
				val := src[s]
				ll += int(val)
				if val != 255 {
					break
				}
			}
		}
		
		if s+ll >= len(src) {
			if debug {
				fmt.Printf("error literals: s+ll (%d+%d) >= len(src) (%d)\n", s, ll, len(src))
			}
			return nil, 0, ErrCorrupt
		}
		s++
		if ll > 0 {
			if d+ll > dLimit {
				return nil, 0, ErrDstTooSmall
			}
			if debug {
				fmt.Printf("emit %d literals\n", ll)
			}
			d += emitLiteralGo(dst[d:], src[s:s+ll])
			s += ll
			uncompressed += ll
		}

		
		if s == len(src) && ml == lz4MinMatch {
			break
		}
		
		if s >= len(src)-2 {
			if debug {
				fmt.Printf("s (%d) >= len(src)-2 (%d)", s, len(src)-2)
			}
			return nil, 0, ErrCorrupt
		}
		offset := binary.LittleEndian.Uint16(src[s:])
		s += 2
		if offset == 0 {
			if debug {
				fmt.Printf("error: offset 0, ml: %d, len(src)-s: %d\n", ml, len(src)-s)
			}
			return nil, 0, ErrCorrupt
		}
		if int(offset) > uncompressed {
			if debug {
				fmt.Printf("error: offset (%d)> uncompressed (%d)\n", offset, uncompressed)
			}
			return nil, 0, ErrCorrupt
		}

		if ml == lz4MinMatch+15 {
			for {
				if s >= len(src) {
					if debug {
						fmt.Printf("error reading ml: s (%d) >= len(src) (%d)\n", s, len(src))
					}
					return nil, 0, ErrCorrupt
				}
				val := src[s]
				s++
				ml += int(val)
				if val != 255 {
					if s >= len(src) {
						if debug {
							fmt.Printf("error reading ml: s (%d) >= len(src) (%d)\n", s, len(src))
						}
						return nil, 0, ErrCorrupt
					}
					break
				}
			}
		}
		if debug {
			fmt.Printf("emit copy, length: %d, offset: %d\n", ml, offset)
		}
		length := ml
		
		for length > 0 {
			if d >= dLimit {
				return nil, 0, ErrDstTooSmall
			}

			
			if length > 64 {
				
				dst[d+2] = uint8(offset >> 8)
				dst[d+1] = uint8(offset)
				dst[d+0] = 63<<2 | tagCopy2
				length -= 64
				d += 3
				continue
			}
			if length >= 12 || offset >= 2048 || length < 4 {
				
				dst[d+2] = uint8(offset >> 8)
				dst[d+1] = uint8(offset)
				dst[d+0] = uint8(length-1)<<2 | tagCopy2
				d += 3
				break
			}
			
			dst[d+1] = uint8(offset)
			dst[d+0] = uint8(offset>>8)<<5 | uint8(length-4)<<2 | tagCopy1
			d += 2
			break
		}
		uncompressed += ml
		if d > dLimit {
			return nil, 0, ErrDstTooSmall
		}
	}

	return dst[:d], uncompressed, nil
}



func emitRepeat16(dst []byte, offset uint16, length int) int {
	
	length -= 4
	if length <= 4 {
		dst[0] = uint8(length)<<2 | tagCopy1
		dst[1] = 0
		return 2
	}
	if length < 8 && offset < 2048 {
		
		dst[1] = uint8(offset)
		dst[0] = uint8(offset>>8)<<5 | uint8(length)<<2 | tagCopy1
		return 2
	}
	if length < (1<<8)+4 {
		length -= 4
		dst[2] = uint8(length)
		dst[1] = 0
		dst[0] = 5<<2 | tagCopy1
		return 3
	}
	if length < (1<<16)+(1<<8) {
		length -= 1 << 8
		dst[3] = uint8(length >> 8)
		dst[2] = uint8(length >> 0)
		dst[1] = 0
		dst[0] = 6<<2 | tagCopy1
		return 4
	}
	const maxRepeat = (1 << 24) - 1
	length -= 1 << 16
	left := 0
	if length > maxRepeat {
		left = length - maxRepeat + 4
		length = maxRepeat - 4
	}
	dst[4] = uint8(length >> 16)
	dst[3] = uint8(length >> 8)
	dst[2] = uint8(length >> 0)
	dst[1] = 0
	dst[0] = 7<<2 | tagCopy1
	if left > 0 {
		return 5 + emitRepeat16(dst[5:], offset, left)
	}
	return 5
}








func emitCopy16(dst []byte, offset uint16, length int) int {
	
	if length > 64 {
		off := 3
		if offset < 2048 {
			
			dst[1] = uint8(offset)
			dst[0] = uint8(offset>>8)<<5 | uint8(8-4)<<2 | tagCopy1
			length -= 8
			off = 2
		} else {
			
			
			dst[2] = uint8(offset >> 8)
			dst[1] = uint8(offset)
			dst[0] = 59<<2 | tagCopy2
			length -= 60
		}
		
		return off + emitRepeat16(dst[off:], offset, length)
	}
	if length >= 12 || offset >= 2048 {
		
		dst[2] = uint8(offset >> 8)
		dst[1] = uint8(offset)
		dst[0] = uint8(length-1)<<2 | tagCopy2
		return 3
	}
	
	dst[1] = uint8(offset)
	dst[0] = uint8(offset>>8)<<5 | uint8(length-4)<<2 | tagCopy1
	return 2
}







func emitLiteralGo(dst, lit []byte) int {
	if len(lit) == 0 {
		return 0
	}
	i, n := 0, uint(len(lit)-1)
	switch {
	case n < 60:
		dst[0] = uint8(n)<<2 | tagLiteral
		i = 1
	case n < 1<<8:
		dst[1] = uint8(n)
		dst[0] = 60<<2 | tagLiteral
		i = 2
	case n < 1<<16:
		dst[2] = uint8(n >> 8)
		dst[1] = uint8(n)
		dst[0] = 61<<2 | tagLiteral
		i = 3
	case n < 1<<24:
		dst[3] = uint8(n >> 16)
		dst[2] = uint8(n >> 8)
		dst[1] = uint8(n)
		dst[0] = 62<<2 | tagLiteral
		i = 4
	default:
		dst[4] = uint8(n >> 24)
		dst[3] = uint8(n >> 16)
		dst[2] = uint8(n >> 8)
		dst[1] = uint8(n)
		dst[0] = 63<<2 | tagLiteral
		i = 5
	}
	return i + copy(dst[i:], lit)
}
