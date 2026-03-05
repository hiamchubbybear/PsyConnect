




package flate

import (
	"encoding/binary"
	"fmt"
	"math/bits"
)

type fastEnc interface {
	Encode(dst *tokens, src []byte)
	Reset()
}

func newFastEnc(level int) fastEnc {
	switch level {
	case 1:
		return &fastEncL1{fastGen: fastGen{cur: maxStoreBlockSize}}
	case 2:
		return &fastEncL2{fastGen: fastGen{cur: maxStoreBlockSize}}
	case 3:
		return &fastEncL3{fastGen: fastGen{cur: maxStoreBlockSize}}
	case 4:
		return &fastEncL4{fastGen: fastGen{cur: maxStoreBlockSize}}
	case 5:
		return &fastEncL5{fastGen: fastGen{cur: maxStoreBlockSize}}
	case 6:
		return &fastEncL6{fastGen: fastGen{cur: maxStoreBlockSize}}
	default:
		panic("invalid level specified")
	}
}

const (
	tableBits       = 15             
	tableSize       = 1 << tableBits 
	tableShift      = 32 - tableBits 
	baseMatchOffset = 1              
	baseMatchLength = 3              
	maxMatchOffset  = 1 << 15        

	bTableBits   = 17                                               
	bTableSize   = 1 << bTableBits                                  
	allocHistory = maxStoreBlockSize * 5                            
	bufferReset  = (1 << 31) - allocHistory - maxStoreBlockSize - 1 
)

const (
	prime3bytes = 506832829
	prime4bytes = 2654435761
	prime5bytes = 889523592379
	prime6bytes = 227718039650203
	prime7bytes = 58295818150454627
	prime8bytes = 0xcf1bbcdcb7a56463
)

func load3232(b []byte, i int32) uint32 {
	return binary.LittleEndian.Uint32(b[i:])
}

func load6432(b []byte, i int32) uint64 {
	return binary.LittleEndian.Uint64(b[i:])
}

type tableEntry struct {
	offset int32
}




type fastGen struct {
	hist []byte
	cur  int32
}

func (e *fastGen) addBlock(src []byte) int32 {
	
	if len(e.hist)+len(src) > cap(e.hist) {
		if cap(e.hist) == 0 {
			e.hist = make([]byte, 0, allocHistory)
		} else {
			if cap(e.hist) < maxMatchOffset*2 {
				panic("unexpected buffer size")
			}
			
			offset := int32(len(e.hist)) - maxMatchOffset
			
			*(*[maxMatchOffset]byte)(e.hist) = *(*[maxMatchOffset]byte)(e.hist[offset:])
			e.cur += offset
			e.hist = e.hist[:maxMatchOffset]
		}
	}
	s := int32(len(e.hist))
	e.hist = append(e.hist, src...)
	return s
}

type tableEntryPrev struct {
	Cur  tableEntry
	Prev tableEntry
}



func hash7(u uint64, h uint8) uint32 {
	return uint32(((u << (64 - 56)) * prime7bytes) >> ((64 - h) & reg8SizeMask64))
}





func hashLen(u uint64, length, mls uint8) uint32 {
	switch mls {
	case 3:
		return (uint32(u<<8) * prime3bytes) >> (32 - length)
	case 5:
		return uint32(((u << (64 - 40)) * prime5bytes) >> (64 - length))
	case 6:
		return uint32(((u << (64 - 48)) * prime6bytes) >> (64 - length))
	case 7:
		return uint32(((u << (64 - 56)) * prime7bytes) >> (64 - length))
	case 8:
		return uint32((u * prime8bytes) >> (64 - length))
	default:
		return (uint32(u) * prime4bytes) >> (32 - length)
	}
}




func (e *fastGen) matchlen(s, t int32, src []byte) int32 {
	if debugDecode {
		if t >= s {
			panic(fmt.Sprint("t >=s:", t, s))
		}
		if int(s) >= len(src) {
			panic(fmt.Sprint("s >= len(src):", s, len(src)))
		}
		if t < 0 {
			panic(fmt.Sprint("t < 0:", t))
		}
		if s-t > maxMatchOffset {
			panic(fmt.Sprint(s, "-", t, "(", s-t, ") > maxMatchLength (", maxMatchOffset, ")"))
		}
	}
	s1 := int(s) + maxMatchLength - 4
	if s1 > len(src) {
		s1 = len(src)
	}

	
	return int32(matchLen(src[s:s1], src[t:]))
}



func (e *fastGen) matchlenLong(s, t int32, src []byte) int32 {
	if debugDeflate {
		if t >= s {
			panic(fmt.Sprint("t >=s:", t, s))
		}
		if int(s) >= len(src) {
			panic(fmt.Sprint("s >= len(src):", s, len(src)))
		}
		if t < 0 {
			panic(fmt.Sprint("t < 0:", t))
		}
		if s-t > maxMatchOffset {
			panic(fmt.Sprint(s, "-", t, "(", s-t, ") > maxMatchLength (", maxMatchOffset, ")"))
		}
	}
	
	return int32(matchLen(src[s:], src[t:]))
}


func (e *fastGen) Reset() {
	if cap(e.hist) < allocHistory {
		e.hist = make([]byte, 0, allocHistory)
	}
	
	
	if e.cur <= bufferReset {
		e.cur += maxMatchOffset + int32(len(e.hist))
	}
	e.hist = e.hist[:0]
}



func matchLen(a, b []byte) int {
	var checked int

	for len(a) >= 8 {
		if diff := binary.LittleEndian.Uint64(a) ^ binary.LittleEndian.Uint64(b); diff != 0 {
			return checked + (bits.TrailingZeros64(diff) >> 3)
		}
		checked += 8
		a = a[8:]
		b = b[8:]
	}
	b = b[:len(a)]
	for i := range a {
		if a[i] != b[i] {
			return i + checked
		}
	}
	return len(a) + checked
}
