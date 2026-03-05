




package s2

import (
	"bytes"
	"fmt"
	"math/bits"
)



func hash4(u uint64, h uint8) uint32 {
	const prime4bytes = 2654435761
	return (uint32(u) * prime4bytes) >> ((32 - h) & 31)
}



func hash5(u uint64, h uint8) uint32 {
	const prime5bytes = 889523592379
	return uint32(((u << (64 - 40)) * prime5bytes) >> ((64 - h) & 63))
}



func hash7(u uint64, h uint8) uint32 {
	const prime7bytes = 58295818150454627
	return uint32(((u << (64 - 56)) * prime7bytes) >> ((64 - h) & 63))
}



func hash8(u uint64, h uint8) uint32 {
	const prime8bytes = 0xcf1bbcdcb7a56463
	return uint32((u * prime8bytes) >> ((64 - h) & 63))
}









func encodeBlockBetterGo(dst, src []byte) (d int) {
	
	
	
	sLimit := len(src) - inputMargin
	if len(src) < minNonLiteralBlockSize {
		return 0
	}

	
	const (
		
		lTableBits    = 17
		maxLTableSize = 1 << lTableBits

		
		sTableBits    = 14
		maxSTableSize = 1 << sTableBits
	)

	var lTable [maxLTableSize]uint32
	var sTable [maxSTableSize]uint32

	
	dstLimit := len(src) - len(src)>>5 - 6

	
	nextEmit := 0

	
	
	s := 1
	cv := load64(src, s)

	
	repeat := 0

	for {
		candidateL := 0
		nextS := 0
		for {
			
			nextS = s + (s-nextEmit)>>7 + 1
			if nextS > sLimit {
				goto emitRemainder
			}
			hashL := hash7(cv, lTableBits)
			hashS := hash4(cv, sTableBits)
			candidateL = int(lTable[hashL])
			candidateS := int(sTable[hashS])
			lTable[hashL] = uint32(s)
			sTable[hashS] = uint32(s)

			valLong := load64(src, candidateL)
			valShort := load64(src, candidateS)

			
			if cv == valLong {
				break
			}
			if cv == valShort {
				candidateL = candidateS
				break
			}

			
			const checkRep = 1
			
			
			
			const wantRepeatBytes = 6
			const repeatMask = ((1 << (wantRepeatBytes * 8)) - 1) << (8 * checkRep)
			if false && repeat > 0 && cv&repeatMask == load64(src, s-repeat)&repeatMask {
				base := s + checkRep
				
				for i := base - repeat; base > nextEmit && i > 0 && src[i-1] == src[base-1]; {
					i--
					base--
				}
				d += emitLiteral(dst[d:], src[nextEmit:base])

				
				candidate := s - repeat + wantRepeatBytes + checkRep
				s += wantRepeatBytes + checkRep
				for s < len(src) {
					if len(src)-s < 8 {
						if src[s] == src[candidate] {
							s++
							candidate++
							continue
						}
						break
					}
					if diff := load64(src, s) ^ load64(src, candidate); diff != 0 {
						s += bits.TrailingZeros64(diff) >> 3
						break
					}
					s += 8
					candidate += 8
				}
				
				d += emitRepeat(dst[d:], repeat, s-base)
				nextEmit = s
				if s >= sLimit {
					goto emitRemainder
				}
				
				index0 := base + 1
				index1 := s - 2

				for index0 < index1 {
					cv0 := load64(src, index0)
					cv1 := load64(src, index1)
					lTable[hash7(cv0, lTableBits)] = uint32(index0)
					sTable[hash4(cv0>>8, sTableBits)] = uint32(index0 + 1)

					lTable[hash7(cv1, lTableBits)] = uint32(index1)
					sTable[hash4(cv1>>8, sTableBits)] = uint32(index1 + 1)
					index0 += 2
					index1 -= 2
				}

				cv = load64(src, s)
				continue
			}

			
			if uint32(cv) == uint32(valLong) {
				break
			}

			
			if uint32(cv) == uint32(valShort) {
				
				hashL = hash7(cv>>8, lTableBits)
				candidateL = int(lTable[hashL])
				lTable[hashL] = uint32(s + 1)
				if uint32(cv>>8) == load32(src, candidateL) {
					s++
					break
				}
				
				candidateL = candidateS
				break
			}

			cv = load64(src, nextS)
			s = nextS
		}

		
		for candidateL > 0 && s > nextEmit && src[candidateL-1] == src[s-1] {
			candidateL--
			s--
		}

		
		if d+(s-nextEmit) > dstLimit {
			return 0
		}

		base := s
		offset := base - candidateL

		
		s += 4
		candidateL += 4
		for s < len(src) {
			if len(src)-s < 8 {
				if src[s] == src[candidateL] {
					s++
					candidateL++
					continue
				}
				break
			}
			if diff := load64(src, s) ^ load64(src, candidateL); diff != 0 {
				s += bits.TrailingZeros64(diff) >> 3
				break
			}
			s += 8
			candidateL += 8
		}

		if offset > 65535 && s-base <= 5 && repeat != offset {
			
			s = nextS + 1
			if s >= sLimit {
				goto emitRemainder
			}
			cv = load64(src, s)
			continue
		}

		d += emitLiteral(dst[d:], src[nextEmit:base])
		if repeat == offset {
			d += emitRepeat(dst[d:], offset, s-base)
		} else {
			d += emitCopy(dst[d:], offset, s-base)
			repeat = offset
		}

		nextEmit = s
		if s >= sLimit {
			goto emitRemainder
		}

		if d > dstLimit {
			
			return 0
		}

		
		index0 := base + 1
		index1 := s - 2

		cv0 := load64(src, index0)
		cv1 := load64(src, index1)
		lTable[hash7(cv0, lTableBits)] = uint32(index0)
		sTable[hash4(cv0>>8, sTableBits)] = uint32(index0 + 1)

		
		lTable[hash7(cv1, lTableBits)] = uint32(index1)
		sTable[hash4(cv1>>8, sTableBits)] = uint32(index1 + 1)
		index0 += 1
		index1 -= 1
		cv = load64(src, s)

		
		
		index2 := (index0 + index1 + 1) >> 1
		for index2 < index1 {
			lTable[hash7(load64(src, index0), lTableBits)] = uint32(index0)
			lTable[hash7(load64(src, index2), lTableBits)] = uint32(index2)
			index0 += 2
			index2 += 2
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









func encodeBlockBetterSnappyGo(dst, src []byte) (d int) {
	
	
	
	sLimit := len(src) - inputMargin
	if len(src) < minNonLiteralBlockSize {
		return 0
	}

	
	const (
		
		lTableBits    = 16
		maxLTableSize = 1 << lTableBits

		
		sTableBits    = 14
		maxSTableSize = 1 << sTableBits
	)

	var lTable [maxLTableSize]uint32
	var sTable [maxSTableSize]uint32

	
	dstLimit := len(src) - len(src)>>5 - 6

	
	nextEmit := 0

	
	
	s := 1
	cv := load64(src, s)

	
	repeat := 0
	const maxSkip = 100

	for {
		candidateL := 0
		nextS := 0
		for {
			
			nextS = (s-nextEmit)>>7 + 1
			if nextS > maxSkip {
				nextS = s + maxSkip
			} else {
				nextS += s
			}

			if nextS > sLimit {
				goto emitRemainder
			}
			hashL := hash7(cv, lTableBits)
			hashS := hash4(cv, sTableBits)
			candidateL = int(lTable[hashL])
			candidateS := int(sTable[hashS])
			lTable[hashL] = uint32(s)
			sTable[hashS] = uint32(s)

			if uint32(cv) == load32(src, candidateL) {
				break
			}

			
			if uint32(cv) == load32(src, candidateS) {
				
				hashL = hash7(cv>>8, lTableBits)
				candidateL = int(lTable[hashL])
				lTable[hashL] = uint32(s + 1)
				if uint32(cv>>8) == load32(src, candidateL) {
					s++
					break
				}
				
				candidateL = candidateS
				break
			}

			cv = load64(src, nextS)
			s = nextS
		}

		
		for candidateL > 0 && s > nextEmit && src[candidateL-1] == src[s-1] {
			candidateL--
			s--
		}

		
		if d+(s-nextEmit) > dstLimit {
			return 0
		}

		base := s
		offset := base - candidateL

		
		s += 4
		candidateL += 4
		for s < len(src) {
			if len(src)-s < 8 {
				if src[s] == src[candidateL] {
					s++
					candidateL++
					continue
				}
				break
			}
			if diff := load64(src, s) ^ load64(src, candidateL); diff != 0 {
				s += bits.TrailingZeros64(diff) >> 3
				break
			}
			s += 8
			candidateL += 8
		}

		if offset > 65535 && s-base <= 5 && repeat != offset {
			
			s = nextS + 1
			if s >= sLimit {
				goto emitRemainder
			}
			cv = load64(src, s)
			continue
		}

		d += emitLiteral(dst[d:], src[nextEmit:base])
		d += emitCopyNoRepeat(dst[d:], offset, s-base)
		repeat = offset

		nextEmit = s
		if s >= sLimit {
			goto emitRemainder
		}

		if d > dstLimit {
			
			return 0
		}

		
		index0 := base + 1
		index1 := s - 2

		cv0 := load64(src, index0)
		cv1 := load64(src, index1)
		lTable[hash7(cv0, lTableBits)] = uint32(index0)
		sTable[hash4(cv0>>8, sTableBits)] = uint32(index0 + 1)

		lTable[hash7(cv1, lTableBits)] = uint32(index1)
		sTable[hash4(cv1>>8, sTableBits)] = uint32(index1 + 1)
		index0 += 1
		index1 -= 1
		cv = load64(src, s)

		
		
		index2 := (index0 + index1 + 1) >> 1
		for index2 < index1 {
			lTable[hash7(load64(src, index0), lTableBits)] = uint32(index0)
			lTable[hash7(load64(src, index2), lTableBits)] = uint32(index2)
			index0 += 2
			index2 += 2
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









func encodeBlockBetterDict(dst, src []byte, dict *Dict) (d int) {
	
	
	
	
	const (
		
		lTableBits    = 17
		maxLTableSize = 1 << lTableBits

		
		sTableBits    = 14
		maxSTableSize = 1 << sTableBits

		maxAhead = 8 

		debug = false
	)

	sLimit := len(src) - inputMargin
	if sLimit > MaxDictSrcOffset-maxAhead {
		sLimit = MaxDictSrcOffset - maxAhead
	}
	if len(src) < minNonLiteralBlockSize {
		return 0
	}

	dict.initBetter()

	var lTable [maxLTableSize]uint32
	var sTable [maxSTableSize]uint32

	
	dstLimit := len(src) - len(src)>>5 - 6

	
	nextEmit := 0

	
	
	s := 0
	cv := load64(src, s)

	
	repeat := len(dict.dict) - dict.repeat

	
searchDict:
	for {
		candidateL := 0
		nextS := 0
		for {
			
			nextS = s + (s-nextEmit)>>7 + 1
			if nextS > sLimit {
				break searchDict
			}
			hashL := hash7(cv, lTableBits)
			hashS := hash4(cv, sTableBits)
			candidateL = int(lTable[hashL])
			candidateS := int(sTable[hashS])
			dictL := int(dict.betterTableLong[hashL])
			dictS := int(dict.betterTableShort[hashS])
			lTable[hashL] = uint32(s)
			sTable[hashS] = uint32(s)

			valLong := load64(src, candidateL)
			valShort := load64(src, candidateS)

			
			if s != 0 {
				if cv == valLong {
					goto emitMatch
				}
				if cv == valShort {
					candidateL = candidateS
					goto emitMatch
				}
			}

			
			if repeat >= s+4 {
				candidate := len(dict.dict) - repeat + s
				if candidate > 0 && uint32(cv) == load32(dict.dict, candidate) {
					
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
					
					index0 := base + 1
					index1 := s - 2

					cv = load64(src, s)
					for index0 < index1 {
						cv0 := load64(src, index0)
						cv1 := load64(src, index1)
						lTable[hash7(cv0, lTableBits)] = uint32(index0)
						sTable[hash4(cv0>>8, sTableBits)] = uint32(index0 + 1)

						lTable[hash7(cv1, lTableBits)] = uint32(index1)
						sTable[hash4(cv1>>8, sTableBits)] = uint32(index1 + 1)
						index0 += 2
						index1 -= 2
					}
					continue
				}
			}
			
			if s == 0 {
				cv = load64(src, nextS)
				s = nextS
				continue
			}

			
			if uint32(cv) == uint32(valLong) {
				goto emitMatch
			}

			
			if uint32(cv) == load32(dict.dict, dictL) {
				candidateL = dictL
				goto emitDict
			}

			
			if uint32(cv) == uint32(valShort) {
				
				hashL = hash7(cv>>8, lTableBits)
				candidateL = int(lTable[hashL])
				lTable[hashL] = uint32(s + 1)
				if uint32(cv>>8) == load32(src, candidateL) {
					s++
					goto emitMatch
				}
				
				candidateL = candidateS
				goto emitMatch
			}
			if uint32(cv) == load32(dict.dict, dictS) {
				
				hashL = hash7(cv>>8, lTableBits)
				candidateL = int(lTable[hashL])
				lTable[hashL] = uint32(s + 1)
				if uint32(cv>>8) == load32(src, candidateL) {
					s++
					goto emitMatch
				}
				candidateL = dictS
				goto emitDict
			}
			cv = load64(src, nextS)
			s = nextS
		}
	emitDict:
		{
			if debug {
				if load32(dict.dict, candidateL) != load32(src, s) {
					panic("dict emit mismatch")
				}
			}
			
			
			for candidateL > 0 && s > nextEmit && dict.dict[candidateL-1] == src[s-1] {
				candidateL--
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
				offset := s + (len(dict.dict)) - candidateL

				
				s += 4
				candidateL += 4
				for s <= len(src)-8 && len(dict.dict)-candidateL >= 8 {
					if diff := load64(src, s) ^ load64(dict.dict, candidateL); diff != 0 {
						s += bits.TrailingZeros64(diff) >> 3
						break
					}
					s += 8
					candidateL += 8
				}

				if repeat == offset {
					if debug {
						fmt.Println("emitted dict repeat, length", s-base, "offset:", offset, "s:", s, "dict offset:", candidateL)
					}
					d += emitRepeat(dst[d:], offset, s-base)
				} else {
					if debug {
						fmt.Println("emitted dict copy, length", s-base, "offset:", offset, "s:", s, "dict offset:", candidateL)
					}
					
					if s <= sLimit || s-base < 8 {
						d += emitCopy(dst[d:], offset, s-base)
					} else {
						
						d += emitCopy(dst[d:], offset, 4)
						d += emitRepeat(dst[d:], offset, s-base-4)
					}
					repeat = offset
				}
				if false {
					
					if s <= candidateL {
						panic("s <= candidate")
					}
					a := src[base:s]
					b := dict.dict[base-repeat : base-repeat+(s-base)]
					if !bytes.Equal(a, b) {
						panic("mismatch")
					}
				}

				nextEmit = s
				if s >= sLimit {
					break searchDict
				}

				if d > dstLimit {
					
					return 0
				}

				
				index0 := base + 1
				index1 := s - 2

				cv0 := load64(src, index0)
				cv1 := load64(src, index1)
				lTable[hash7(cv0, lTableBits)] = uint32(index0)
				sTable[hash4(cv0>>8, sTableBits)] = uint32(index0 + 1)

				lTable[hash7(cv1, lTableBits)] = uint32(index1)
				sTable[hash4(cv1>>8, sTableBits)] = uint32(index1 + 1)
				index0 += 1
				index1 -= 1
				cv = load64(src, s)

				
				for index0 < index1 {
					lTable[hash7(load64(src, index0), lTableBits)] = uint32(index0)
					lTable[hash7(load64(src, index1), lTableBits)] = uint32(index1)
					index0 += 2
					index1 -= 2
				}
			}
			continue
		}
	emitMatch:

		
		for candidateL > 0 && s > nextEmit && src[candidateL-1] == src[s-1] {
			candidateL--
			s--
		}

		
		if d+(s-nextEmit) > dstLimit {
			return 0
		}

		base := s
		offset := base - candidateL

		
		s += 4
		candidateL += 4
		for s < len(src) {
			if len(src)-s < 8 {
				if src[s] == src[candidateL] {
					s++
					candidateL++
					continue
				}
				break
			}
			if diff := load64(src, s) ^ load64(src, candidateL); diff != 0 {
				s += bits.TrailingZeros64(diff) >> 3
				break
			}
			s += 8
			candidateL += 8
		}

		if offset > 65535 && s-base <= 5 && repeat != offset {
			
			s = nextS + 1
			if s >= sLimit {
				goto emitRemainder
			}
			cv = load64(src, s)
			continue
		}

		d += emitLiteral(dst[d:], src[nextEmit:base])
		if debug && nextEmit != s {
			fmt.Println("emitted ", s-nextEmit, "literals")
		}
		if repeat == offset {
			if debug {
				fmt.Println("emitted match repeat, length", s-base, "offset:", offset, "s:", s)
			}
			d += emitRepeat(dst[d:], offset, s-base)
		} else {
			if debug {
				fmt.Println("emitted match copy, length", s-base, "offset:", offset, "s:", s)
			}
			d += emitCopy(dst[d:], offset, s-base)
			repeat = offset
		}

		nextEmit = s
		if s >= sLimit {
			goto emitRemainder
		}

		if d > dstLimit {
			
			return 0
		}

		
		index0 := base + 1
		index1 := s - 2

		cv0 := load64(src, index0)
		cv1 := load64(src, index1)
		lTable[hash7(cv0, lTableBits)] = uint32(index0)
		sTable[hash4(cv0>>8, sTableBits)] = uint32(index0 + 1)

		lTable[hash7(cv1, lTableBits)] = uint32(index1)
		sTable[hash4(cv1>>8, sTableBits)] = uint32(index1 + 1)
		index0 += 1
		index1 -= 1
		cv = load64(src, s)

		
		
		index2 := (index0 + index1 + 1) >> 1
		for index2 < index1 {
			lTable[hash7(load64(src, index0), lTableBits)] = uint32(index0)
			lTable[hash7(load64(src, index2), lTableBits)] = uint32(index2)
			index0 += 2
			index2 += 2
		}
	}

	
	if repeat > s {
		repeat = 0
	}

	
	sLimit = len(src) - inputMargin
	if s >= sLimit {
		goto emitRemainder
	}
	cv = load64(src, s)
	if debug {
		fmt.Println("now", s, "->", sLimit, "out:", d, "left:", len(src)-s, "nextemit:", nextEmit, "dstLimit:", dstLimit, "s:", s)
	}
	for {
		candidateL := 0
		nextS := 0
		for {
			
			nextS = s + (s-nextEmit)>>7 + 1
			if nextS > sLimit {
				goto emitRemainder
			}
			hashL := hash7(cv, lTableBits)
			hashS := hash4(cv, sTableBits)
			candidateL = int(lTable[hashL])
			candidateS := int(sTable[hashS])
			lTable[hashL] = uint32(s)
			sTable[hashS] = uint32(s)

			valLong := load64(src, candidateL)
			valShort := load64(src, candidateS)

			
			if cv == valLong {
				break
			}
			if cv == valShort {
				candidateL = candidateS
				break
			}

			
			const checkRep = 1
			
			
			
			const wantRepeatBytes = 6
			const repeatMask = ((1 << (wantRepeatBytes * 8)) - 1) << (8 * checkRep)
			if false && repeat > 0 && cv&repeatMask == load64(src, s-repeat)&repeatMask {
				base := s + checkRep
				
				for i := base - repeat; base > nextEmit && i > 0 && src[i-1] == src[base-1]; {
					i--
					base--
				}
				d += emitLiteral(dst[d:], src[nextEmit:base])

				
				candidate := s - repeat + wantRepeatBytes + checkRep
				s += wantRepeatBytes + checkRep
				for s < len(src) {
					if len(src)-s < 8 {
						if src[s] == src[candidate] {
							s++
							candidate++
							continue
						}
						break
					}
					if diff := load64(src, s) ^ load64(src, candidate); diff != 0 {
						s += bits.TrailingZeros64(diff) >> 3
						break
					}
					s += 8
					candidate += 8
				}
				
				d += emitRepeat(dst[d:], repeat, s-base)
				nextEmit = s
				if s >= sLimit {
					goto emitRemainder
				}
				
				index0 := base + 1
				index1 := s - 2

				for index0 < index1 {
					cv0 := load64(src, index0)
					cv1 := load64(src, index1)
					lTable[hash7(cv0, lTableBits)] = uint32(index0)
					sTable[hash4(cv0>>8, sTableBits)] = uint32(index0 + 1)

					lTable[hash7(cv1, lTableBits)] = uint32(index1)
					sTable[hash4(cv1>>8, sTableBits)] = uint32(index1 + 1)
					index0 += 2
					index1 -= 2
				}

				cv = load64(src, s)
				continue
			}

			
			if uint32(cv) == uint32(valLong) {
				break
			}

			
			if uint32(cv) == uint32(valShort) {
				
				hashL = hash7(cv>>8, lTableBits)
				candidateL = int(lTable[hashL])
				lTable[hashL] = uint32(s + 1)
				if uint32(cv>>8) == load32(src, candidateL) {
					s++
					break
				}
				
				candidateL = candidateS
				break
			}

			cv = load64(src, nextS)
			s = nextS
		}

		
		for candidateL > 0 && s > nextEmit && src[candidateL-1] == src[s-1] {
			candidateL--
			s--
		}

		
		if d+(s-nextEmit) > dstLimit {
			return 0
		}

		base := s
		offset := base - candidateL

		
		s += 4
		candidateL += 4
		for s < len(src) {
			if len(src)-s < 8 {
				if src[s] == src[candidateL] {
					s++
					candidateL++
					continue
				}
				break
			}
			if diff := load64(src, s) ^ load64(src, candidateL); diff != 0 {
				s += bits.TrailingZeros64(diff) >> 3
				break
			}
			s += 8
			candidateL += 8
		}

		if offset > 65535 && s-base <= 5 && repeat != offset {
			
			s = nextS + 1
			if s >= sLimit {
				goto emitRemainder
			}
			cv = load64(src, s)
			continue
		}

		d += emitLiteral(dst[d:], src[nextEmit:base])
		if repeat == offset {
			d += emitRepeat(dst[d:], offset, s-base)
		} else {
			d += emitCopy(dst[d:], offset, s-base)
			repeat = offset
		}

		nextEmit = s
		if s >= sLimit {
			goto emitRemainder
		}

		if d > dstLimit {
			
			return 0
		}

		
		index0 := base + 1
		index1 := s - 2

		cv0 := load64(src, index0)
		cv1 := load64(src, index1)
		lTable[hash7(cv0, lTableBits)] = uint32(index0)
		sTable[hash4(cv0>>8, sTableBits)] = uint32(index0 + 1)

		lTable[hash7(cv1, lTableBits)] = uint32(index1)
		sTable[hash4(cv1>>8, sTableBits)] = uint32(index1 + 1)
		index0 += 1
		index1 -= 1
		cv = load64(src, s)

		
		
		index2 := (index0 + index1 + 1) >> 1
		for index2 < index1 {
			lTable[hash7(load64(src, index0), lTableBits)] = uint32(index0)
			lTable[hash7(load64(src, index2), lTableBits)] = uint32(index2)
			index0 += 2
			index2 += 2
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
