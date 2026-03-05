




package s2

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/bits"
)

func load32(b []byte, i int) uint32 {
	return binary.LittleEndian.Uint32(b[i:])
}

func load64(b []byte, i int) uint64 {
	return binary.LittleEndian.Uint64(b[i:])
}



func hash6(u uint64, h uint8) uint32 {
	const prime6bytes = 227718039650203
	return uint32(((u << (64 - 48)) * prime6bytes) >> ((64 - h) & 63))
}

func encodeGo(dst, src []byte) []byte {
	if n := MaxEncodedLen(len(src)); n < 0 {
		panic(ErrTooLarge)
	} else if len(dst) < n {
		dst = make([]byte, n)
	}

	
	d := binary.PutUvarint(dst, uint64(len(src)))

	if len(src) == 0 {
		return dst[:d]
	}
	if len(src) < minNonLiteralBlockSize {
		d += emitLiteral(dst[d:], src)
		return dst[:d]
	}
	n := encodeBlockGo(dst[d:], src)
	if n > 0 {
		d += n
		return dst[:d]
	}
	
	d += emitLiteral(dst[d:], src)
	return dst[:d]
}









func encodeBlockGo(dst, src []byte) (d int) {
	
	const (
		tableBits    = 14
		maxTableSize = 1 << tableBits

		debug = false
	)

	var table [maxTableSize]uint32

	
	
	
	sLimit := len(src) - inputMargin

	
	dstLimit := len(src) - len(src)>>5 - 5

	
	nextEmit := 0

	
	
	s := 1
	cv := load64(src, s)

	
	repeat := 1

	for {
		candidate := 0
		for {
			
			nextS := s + (s-nextEmit)>>6 + 4
			if nextS > sLimit {
				goto emitRemainder
			}
			hash0 := hash6(cv, tableBits)
			hash1 := hash6(cv>>8, tableBits)
			candidate = int(table[hash0])
			candidate2 := int(table[hash1])
			table[hash0] = uint32(s)
			table[hash1] = uint32(s + 1)
			hash2 := hash6(cv>>16, tableBits)

			
			const checkRep = 1
			if uint32(cv>>(checkRep*8)) == load32(src, s-repeat+checkRep) {
				base := s + checkRep
				
				for i := base - repeat; base > nextEmit && i > 0 && src[i-1] == src[base-1]; {
					i--
					base--
				}
				d += emitLiteral(dst[d:], src[nextEmit:base])

				
				candidate := s - repeat + 4 + checkRep
				s += 4 + checkRep
				for s <= sLimit {
					if diff := load64(src, s) ^ load64(src, candidate); diff != 0 {
						s += bits.TrailingZeros64(diff) >> 3
						break
					}
					s += 8
					candidate += 8
				}
				if debug {
					
					if s <= candidate {
						panic("s <= candidate")
					}
					a := src[base:s]
					b := src[base-repeat : base-repeat+(s-base)]
					if !bytes.Equal(a, b) {
						panic("mismatch")
					}
				}
				if nextEmit > 0 {
					
					d += emitRepeat(dst[d:], repeat, s-base)
				} else {
					
					d += emitCopy(dst[d:], repeat, s-base)
				}
				nextEmit = s
				if s >= sLimit {
					goto emitRemainder
				}

				cv = load64(src, s)
				continue
			}

			if uint32(cv) == load32(src, candidate) {
				break
			}
			candidate = int(table[hash2])
			if uint32(cv>>8) == load32(src, candidate2) {
				table[hash2] = uint32(s + 2)
				candidate = candidate2
				s++
				break
			}
			table[hash2] = uint32(s + 2)
			if uint32(cv>>16) == load32(src, candidate) {
				s += 2
				break
			}

			cv = load64(src, nextS)
			s = nextS
		}

		
		
		for candidate > 0 && s > nextEmit && src[candidate-1] == src[s-1] {
			candidate--
			s--
		}

		
		if d+(s-nextEmit) > dstLimit {
			return 0
		}

		
		
		

		d += emitLiteral(dst[d:], src[nextEmit:s])

		
		
		
		
		
		
		
		
		for {
			
			
			base := s
			repeat = base - candidate

			
			s += 4
			candidate += 4
			for s <= len(src)-8 {
				if diff := load64(src, s) ^ load64(src, candidate); diff != 0 {
					s += bits.TrailingZeros64(diff) >> 3
					break
				}
				s += 8
				candidate += 8
			}

			d += emitCopy(dst[d:], repeat, s-base)
			if debug {
				
				if s <= candidate {
					panic("s <= candidate")
				}
				a := src[base:s]
				b := src[base-repeat : base-repeat+(s-base)]
				if !bytes.Equal(a, b) {
					panic("mismatch")
				}
			}

			nextEmit = s
			if s >= sLimit {
				goto emitRemainder
			}

			if d > dstLimit {
				
				return 0
			}
			
			x := load64(src, s-2)
			m2Hash := hash6(x, tableBits)
			currHash := hash6(x>>16, tableBits)
			candidate = int(table[currHash])
			table[m2Hash] = uint32(s - 2)
			table[currHash] = uint32(s)
			if debug && s == candidate {
				panic("s == candidate")
			}
			if uint32(x>>16) != load32(src, candidate) {
				cv = load64(src, s+1)
				s++
				break
			}
		}
	}

emitRemainder:
	if nextEmit < len(src) {
		
		if d+len(src)-nextEmit > dstLimit {
			return 0
		}
		d += emitLiteral(dst[d:], src[nextEmit:])
	}
	return d
}

func encodeBlockSnappyGo(dst, src []byte) (d int) {
	
	const (
		tableBits    = 14
		maxTableSize = 1 << tableBits
	)

	var table [maxTableSize]uint32

	
	
	
	sLimit := len(src) - inputMargin

	
	dstLimit := len(src) - len(src)>>5 - 5

	
	nextEmit := 0

	
	
	s := 1
	cv := load64(src, s)

	
	repeat := 1

	for {
		candidate := 0
		for {
			
			nextS := s + (s-nextEmit)>>6 + 4
			if nextS > sLimit {
				goto emitRemainder
			}
			hash0 := hash6(cv, tableBits)
			hash1 := hash6(cv>>8, tableBits)
			candidate = int(table[hash0])
			candidate2 := int(table[hash1])
			table[hash0] = uint32(s)
			table[hash1] = uint32(s + 1)
			hash2 := hash6(cv>>16, tableBits)

			
			const checkRep = 1
			if uint32(cv>>(checkRep*8)) == load32(src, s-repeat+checkRep) {
				base := s + checkRep
				
				for i := base - repeat; base > nextEmit && i > 0 && src[i-1] == src[base-1]; {
					i--
					base--
				}
				d += emitLiteral(dst[d:], src[nextEmit:base])

				
				candidate := s - repeat + 4 + checkRep
				s += 4 + checkRep
				for s <= sLimit {
					if diff := load64(src, s) ^ load64(src, candidate); diff != 0 {
						s += bits.TrailingZeros64(diff) >> 3
						break
					}
					s += 8
					candidate += 8
				}

				d += emitCopyNoRepeat(dst[d:], repeat, s-base)
				nextEmit = s
				if s >= sLimit {
					goto emitRemainder
				}

				cv = load64(src, s)
				continue
			}

			if uint32(cv) == load32(src, candidate) {
				break
			}
			candidate = int(table[hash2])
			if uint32(cv>>8) == load32(src, candidate2) {
				table[hash2] = uint32(s + 2)
				candidate = candidate2
				s++
				break
			}
			table[hash2] = uint32(s + 2)
			if uint32(cv>>16) == load32(src, candidate) {
				s += 2
				break
			}

			cv = load64(src, nextS)
			s = nextS
		}

		
		for candidate > 0 && s > nextEmit && src[candidate-1] == src[s-1] {
			candidate--
			s--
		}

		
		if d+(s-nextEmit) > dstLimit {
			return 0
		}

		
		
		

		d += emitLiteral(dst[d:], src[nextEmit:s])

		
		
		
		
		
		
		
		
		for {
			
			
			base := s
			repeat = base - candidate

			
			s += 4
			candidate += 4
			for s <= len(src)-8 {
				if diff := load64(src, s) ^ load64(src, candidate); diff != 0 {
					s += bits.TrailingZeros64(diff) >> 3
					break
				}
				s += 8
				candidate += 8
			}

			d += emitCopyNoRepeat(dst[d:], repeat, s-base)
			if false {
				
				a := src[base:s]
				b := src[base-repeat : base-repeat+(s-base)]
				if !bytes.Equal(a, b) {
					panic("mismatch")
				}
			}

			nextEmit = s
			if s >= sLimit {
				goto emitRemainder
			}

			if d > dstLimit {
				
				return 0
			}
			
			x := load64(src, s-2)
			m2Hash := hash6(x, tableBits)
			currHash := hash6(x>>16, tableBits)
			candidate = int(table[currHash])
			table[m2Hash] = uint32(s - 2)
			table[currHash] = uint32(s)
			if uint32(x>>16) != load32(src, candidate) {
				cv = load64(src, s+1)
				s++
				break
			}
		}
	}

emitRemainder:
	if nextEmit < len(src) {
		
		if d+len(src)-nextEmit > dstLimit {
			return 0
		}
		d += emitLiteral(dst[d:], src[nextEmit:])
	}
	return d
}









func encodeBlockDictGo(dst, src []byte, dict *Dict) (d int) {
	
	const (
		tableBits    = 14
		maxTableSize = 1 << tableBits
		maxAhead     = 8 

		debug = false
	)
	dict.initFast()

	var table [maxTableSize]uint32

	
	
	
	sLimit := len(src) - inputMargin
	if sLimit > MaxDictSrcOffset-maxAhead {
		sLimit = MaxDictSrcOffset - maxAhead
	}

	
	dstLimit := len(src) - len(src)>>5 - 5

	
	nextEmit := 0

	
	s := 0

	
	repeat := len(dict.dict) - dict.repeat
	cv := load64(src, 0)

	
searchDict:
	for {
		
		nextS := s + (s-nextEmit)>>6 + 4
		hash0 := hash6(cv, tableBits)
		hash1 := hash6(cv>>8, tableBits)
		if nextS > sLimit {
			if debug {
				fmt.Println("slimit reached", s, nextS)
			}
			break searchDict
		}
		candidateDict := int(dict.fastTable[hash0])
		candidateDict2 := int(dict.fastTable[hash1])
		candidate2 := int(table[hash1])
		candidate := int(table[hash0])
		table[hash0] = uint32(s)
		table[hash1] = uint32(s + 1)
		hash2 := hash6(cv>>16, tableBits)

		
		const checkRep = 1

		if repeat > s {
			candidate := len(dict.dict) - repeat + s
			if repeat-s >= 4 && uint32(cv) == load32(dict.dict, candidate) {
				
				base := s
				for i := candidate; base > nextEmit && i > 0 && dict.dict[i-1] == src[base-1]; {
					i--
					base--
				}
				d += emitLiteral(dst[d:], src[nextEmit:base])
				if debug && nextEmit != base {
					fmt.Println("emitted ", base-nextEmit, "literals")
				}
				s += 4
				candidate += 4
				for candidate < len(dict.dict)-8 && s <= len(src)-8 {
					if diff := load64(src, s) ^ load64(dict.dict, candidate); diff != 0 {
						s += bits.TrailingZeros64(diff) >> 3
						break
					}
					s += 8
					candidate += 8
				}
				d += emitRepeat(dst[d:], repeat, s-base)
				if debug {
					fmt.Println("emitted dict repeat length", s-base, "offset:", repeat, "s:", s)
				}
				nextEmit = s
				if s >= sLimit {
					break searchDict
				}
				cv = load64(src, s)
				continue
			}
		} else if uint32(cv>>(checkRep*8)) == load32(src, s-repeat+checkRep) {
			base := s + checkRep
			
			for i := base - repeat; base > nextEmit && i > 0 && src[i-1] == src[base-1]; {
				i--
				base--
			}
			d += emitLiteral(dst[d:], src[nextEmit:base])
			if debug && nextEmit != base {
				fmt.Println("emitted ", base-nextEmit, "literals")
			}

			
			candidate := s - repeat + 4 + checkRep
			s += 4 + checkRep
			for s <= sLimit {
				if diff := load64(src, s) ^ load64(src, candidate); diff != 0 {
					s += bits.TrailingZeros64(diff) >> 3
					break
				}
				s += 8
				candidate += 8
			}
			if debug {
				
				if s <= candidate {
					panic("s <= candidate")
				}
				a := src[base:s]
				b := src[base-repeat : base-repeat+(s-base)]
				if !bytes.Equal(a, b) {
					panic("mismatch")
				}
			}

			if nextEmit > 0 {
				
				d += emitRepeat(dst[d:], repeat, s-base)
			} else {
				
				d += emitCopy(dst[d:], repeat, s-base)
			}

			nextEmit = s
			if s >= sLimit {
				break searchDict
			}
			if debug {
				fmt.Println("emitted reg repeat", s-base, "s:", s)
			}
			cv = load64(src, s)
			continue searchDict
		}
		if s == 0 {
			cv = load64(src, nextS)
			s = nextS
			continue searchDict
		}
		
		if uint32(cv) == load32(src, candidate) {
			goto emitMatch
		}
		candidate = int(table[hash2])
		if uint32(cv>>8) == load32(src, candidate2) {
			table[hash2] = uint32(s + 2)
			candidate = candidate2
			s++
			goto emitMatch
		}

		
		if cv == load64(dict.dict, candidateDict) {
			table[hash2] = uint32(s + 2)
			goto emitDict
		}

		candidateDict = int(dict.fastTable[hash2])
		
		if candidateDict2 >= 1 {
			if cv^load64(dict.dict, candidateDict2-1) < (1 << 8) {
				table[hash2] = uint32(s + 2)
				candidateDict = candidateDict2
				s++
				goto emitDict
			}
		}

		table[hash2] = uint32(s + 2)
		if uint32(cv>>16) == load32(src, candidate) {
			s += 2
			goto emitMatch
		}
		if candidateDict >= 2 {
			
			if cv^load64(dict.dict, candidateDict-2) < (1 << 16) {
				s += 2
				goto emitDict
			}
		}

		cv = load64(src, nextS)
		s = nextS
		continue searchDict

	emitDict:
		{
			if debug {
				if load32(dict.dict, candidateDict) != load32(src, s) {
					panic("dict emit mismatch")
				}
			}
			
			
			for candidateDict > 0 && s > nextEmit && dict.dict[candidateDict-1] == src[s-1] {
				candidateDict--
				s--
			}

			
			if d+(s-nextEmit) > dstLimit {
				return 0
			}

			
			
			

			d += emitLiteral(dst[d:], src[nextEmit:s])
			if debug && nextEmit != s {
				fmt.Println("emitted ", s-nextEmit, "literals")
			}
			{
				
				
				base := s
				repeat = s + (len(dict.dict)) - candidateDict

				
				s += 4
				candidateDict += 4
				for s <= len(src)-8 && len(dict.dict)-candidateDict >= 8 {
					if diff := load64(src, s) ^ load64(dict.dict, candidateDict); diff != 0 {
						s += bits.TrailingZeros64(diff) >> 3
						break
					}
					s += 8
					candidateDict += 8
				}

				
				if s <= sLimit || s-base < 8 {
					d += emitCopy(dst[d:], repeat, s-base)
				} else {
					
					d += emitCopy(dst[d:], repeat, 4)
					d += emitRepeat(dst[d:], repeat, s-base-4)
				}
				if false {
					
					if s <= candidate {
						panic("s <= candidate")
					}
					a := src[base:s]
					b := dict.dict[base-repeat : base-repeat+(s-base)]
					if !bytes.Equal(a, b) {
						panic("mismatch")
					}
				}
				if debug {
					fmt.Println("emitted dict copy, length", s-base, "offset:", repeat, "s:", s)
				}
				nextEmit = s
				if s >= sLimit {
					break searchDict
				}

				if d > dstLimit {
					
					return 0
				}

				
				x := load64(src, s-2)
				m2Hash := hash6(x, tableBits)
				currHash := hash6(x>>8, tableBits)
				table[m2Hash] = uint32(s - 2)
				table[currHash] = uint32(s - 1)
				cv = load64(src, s)
			}
			continue
		}
	emitMatch:

		
		
		for candidate > 0 && s > nextEmit && src[candidate-1] == src[s-1] {
			candidate--
			s--
		}

		
		if d+(s-nextEmit) > dstLimit {
			return 0
		}

		
		
		

		d += emitLiteral(dst[d:], src[nextEmit:s])
		if debug && nextEmit != s {
			fmt.Println("emitted ", s-nextEmit, "literals")
		}
		
		
		
		
		
		
		
		
		for {
			
			
			base := s
			repeat = base - candidate

			
			s += 4
			candidate += 4
			for s <= len(src)-8 {
				if diff := load64(src, s) ^ load64(src, candidate); diff != 0 {
					s += bits.TrailingZeros64(diff) >> 3
					break
				}
				s += 8
				candidate += 8
			}

			d += emitCopy(dst[d:], repeat, s-base)
			if debug {
				
				if s <= candidate {
					panic("s <= candidate")
				}
				a := src[base:s]
				b := src[base-repeat : base-repeat+(s-base)]
				if !bytes.Equal(a, b) {
					panic("mismatch")
				}
			}
			if debug {
				fmt.Println("emitted src copy, length", s-base, "offset:", repeat, "s:", s)
			}
			nextEmit = s
			if s >= sLimit {
				break searchDict
			}

			if d > dstLimit {
				
				return 0
			}
			
			x := load64(src, s-2)
			m2Hash := hash6(x, tableBits)
			currHash := hash6(x>>16, tableBits)
			candidate = int(table[currHash])
			table[m2Hash] = uint32(s - 2)
			table[currHash] = uint32(s)
			if debug && s == candidate {
				panic("s == candidate")
			}
			if uint32(x>>16) != load32(src, candidate) {
				cv = load64(src, s+1)
				s++
				break
			}
		}
	}

	
	if repeat > s {
		repeat = 0
	}

	
	sLimit = len(src) - inputMargin
	if s >= sLimit {
		goto emitRemainder
	}
	if debug {
		fmt.Println("non-dict matching at", s, "repeat:", repeat)
	}
	cv = load64(src, s)
	if debug {
		fmt.Println("now", s, "->", sLimit, "out:", d, "left:", len(src)-s, "nextemit:", nextEmit, "dstLimit:", dstLimit, "s:", s)
	}
	for {
		candidate := 0
		for {
			
			nextS := s + (s-nextEmit)>>6 + 4
			if nextS > sLimit {
				goto emitRemainder
			}
			hash0 := hash6(cv, tableBits)
			hash1 := hash6(cv>>8, tableBits)
			candidate = int(table[hash0])
			candidate2 := int(table[hash1])
			table[hash0] = uint32(s)
			table[hash1] = uint32(s + 1)
			hash2 := hash6(cv>>16, tableBits)

			
			const checkRep = 1
			if repeat > 0 && uint32(cv>>(checkRep*8)) == load32(src, s-repeat+checkRep) {
				base := s + checkRep
				
				for i := base - repeat; base > nextEmit && i > 0 && src[i-1] == src[base-1]; {
					i--
					base--
				}
				d += emitLiteral(dst[d:], src[nextEmit:base])
				if debug && nextEmit != base {
					fmt.Println("emitted ", base-nextEmit, "literals")
				}
				
				candidate := s - repeat + 4 + checkRep
				s += 4 + checkRep
				for s <= sLimit {
					if diff := load64(src, s) ^ load64(src, candidate); diff != 0 {
						s += bits.TrailingZeros64(diff) >> 3
						break
					}
					s += 8
					candidate += 8
				}
				if debug {
					
					if s <= candidate {
						panic("s <= candidate")
					}
					a := src[base:s]
					b := src[base-repeat : base-repeat+(s-base)]
					if !bytes.Equal(a, b) {
						panic("mismatch")
					}
				}
				if nextEmit > 0 {
					
					d += emitRepeat(dst[d:], repeat, s-base)
				} else {
					
					d += emitCopy(dst[d:], repeat, s-base)
				}
				if debug {
					fmt.Println("emitted src repeat length", s-base, "offset:", repeat, "s:", s)
				}
				nextEmit = s
				if s >= sLimit {
					goto emitRemainder
				}

				cv = load64(src, s)
				continue
			}

			if uint32(cv) == load32(src, candidate) {
				break
			}
			candidate = int(table[hash2])
			if uint32(cv>>8) == load32(src, candidate2) {
				table[hash2] = uint32(s + 2)
				candidate = candidate2
				s++
				break
			}
			table[hash2] = uint32(s + 2)
			if uint32(cv>>16) == load32(src, candidate) {
				s += 2
				break
			}

			cv = load64(src, nextS)
			s = nextS
		}

		
		
		for candidate > 0 && s > nextEmit && src[candidate-1] == src[s-1] {
			candidate--
			s--
		}

		
		if d+(s-nextEmit) > dstLimit {
			return 0
		}

		
		
		

		d += emitLiteral(dst[d:], src[nextEmit:s])
		if debug && nextEmit != s {
			fmt.Println("emitted ", s-nextEmit, "literals")
		}
		
		
		
		
		
		
		
		
		for {
			
			
			base := s
			repeat = base - candidate

			
			s += 4
			candidate += 4
			for s <= len(src)-8 {
				if diff := load64(src, s) ^ load64(src, candidate); diff != 0 {
					s += bits.TrailingZeros64(diff) >> 3
					break
				}
				s += 8
				candidate += 8
			}

			d += emitCopy(dst[d:], repeat, s-base)
			if debug {
				
				if s <= candidate {
					panic("s <= candidate")
				}
				a := src[base:s]
				b := src[base-repeat : base-repeat+(s-base)]
				if !bytes.Equal(a, b) {
					panic("mismatch")
				}
			}
			if debug {
				fmt.Println("emitted src copy, length", s-base, "offset:", repeat, "s:", s)
			}
			nextEmit = s
			if s >= sLimit {
				goto emitRemainder
			}

			if d > dstLimit {
				
				return 0
			}
			
			x := load64(src, s-2)
			m2Hash := hash6(x, tableBits)
			currHash := hash6(x>>16, tableBits)
			candidate = int(table[currHash])
			table[m2Hash] = uint32(s - 2)
			table[currHash] = uint32(s)
			if debug && s == candidate {
				panic("s == candidate")
			}
			if uint32(x>>16) != load32(src, candidate) {
				cv = load64(src, s+1)
				s++
				break
			}
		}
	}

emitRemainder:
	if nextEmit < len(src) {
		
		if d+len(src)-nextEmit > dstLimit {
			return 0
		}
		d += emitLiteral(dst[d:], src[nextEmit:])
		if debug && nextEmit != s {
			fmt.Println("emitted ", len(src)-nextEmit, "literals")
		}
	}
	return d
}
