



package zstd

import "fmt"

const (
	betterLongTableBits = 19                       
	betterLongTableSize = 1 << betterLongTableBits 
	betterLongLen       = 8                        

	
	
	
	
	betterShortTableBits = 13                        
	betterShortTableSize = 1 << betterShortTableBits 
	betterShortLen       = 5                         

	betterLongTableShardCnt  = 1 << (betterLongTableBits - dictShardBits)    
	betterLongTableShardSize = betterLongTableSize / betterLongTableShardCnt 

	betterShortTableShardCnt  = 1 << (betterShortTableBits - dictShardBits)     
	betterShortTableShardSize = betterShortTableSize / betterShortTableShardCnt 
)

type prevEntry struct {
	offset int32
	prev   int32
}







type betterFastEncoder struct {
	fastBase
	table     [betterShortTableSize]tableEntry
	longTable [betterLongTableSize]prevEntry
}

type betterFastEncoderDict struct {
	betterFastEncoder
	dictTable            []tableEntry
	dictLongTable        []prevEntry
	shortTableShardDirty [betterShortTableShardCnt]bool
	longTableShardDirty  [betterLongTableShardCnt]bool
	allDirty             bool
}


func (e *betterFastEncoder) Encode(blk *blockEnc, src []byte) {
	const (
		
		
		inputMargin            = 8 + 2
		minNonLiteralBlockSize = 16
	)

	
	for e.cur >= e.bufferReset-int32(len(e.hist)) {
		if len(e.hist) == 0 {
			e.table = [betterShortTableSize]tableEntry{}
			e.longTable = [betterLongTableSize]prevEntry{}
			e.cur = e.maxMatchOff
			break
		}
		
		minOff := e.cur + int32(len(e.hist)) - e.maxMatchOff
		for i := range e.table[:] {
			v := e.table[i].offset
			if v < minOff {
				v = 0
			} else {
				v = v - e.cur + e.maxMatchOff
			}
			e.table[i].offset = v
		}
		for i := range e.longTable[:] {
			v := e.longTable[i].offset
			v2 := e.longTable[i].prev
			if v < minOff {
				v = 0
				v2 = 0
			} else {
				v = v - e.cur + e.maxMatchOff
				if v2 < minOff {
					v2 = 0
				} else {
					v2 = v2 - e.cur + e.maxMatchOff
				}
			}
			e.longTable[i] = prevEntry{
				offset: v,
				prev:   v2,
			}
		}
		e.cur = e.maxMatchOff
		break
	}

	s := e.addBlock(src)
	blk.size = len(src)
	if len(src) < minNonLiteralBlockSize {
		blk.extraLits = len(src)
		blk.literals = blk.literals[:len(src)]
		copy(blk.literals, src)
		return
	}

	
	src = e.hist
	sLimit := int32(len(src)) - inputMargin
	
	
	const stepSize = 1

	const kSearchStrength = 9

	
	nextEmit := s
	cv := load6432(src, s)

	
	offset1 := int32(blk.recentOffsets[0])
	offset2 := int32(blk.recentOffsets[1])

	addLiterals := func(s *seq, until int32) {
		if until == nextEmit {
			return
		}
		blk.literals = append(blk.literals, src[nextEmit:until]...)
		s.litLen = uint32(until - nextEmit)
	}
	if debugEncoder {
		println("recent offsets:", blk.recentOffsets)
	}

encodeLoop:
	for {
		var t int32
		
		canRepeat := len(blk.sequences) > 2
		var matched int32

		for {
			if debugAsserts && canRepeat && offset1 == 0 {
				panic("offset0 was 0")
			}

			nextHashL := hashLen(cv, betterLongTableBits, betterLongLen)
			nextHashS := hashLen(cv, betterShortTableBits, betterShortLen)
			candidateL := e.longTable[nextHashL]
			candidateS := e.table[nextHashS]

			const repOff = 1
			repIndex := s - offset1 + repOff
			off := s + e.cur
			e.longTable[nextHashL] = prevEntry{offset: off, prev: candidateL.offset}
			e.table[nextHashS] = tableEntry{offset: off, val: uint32(cv)}

			if canRepeat {
				if repIndex >= 0 && load3232(src, repIndex) == uint32(cv>>(repOff*8)) {
					
					var seq seq
					lenght := 4 + e.matchlen(s+4+repOff, repIndex+4, src)

					seq.matchLen = uint32(lenght - zstdMinMatch)

					
					
					start := s + repOff
					
					
					startLimit := nextEmit + 1

					tMin := s - e.maxMatchOff
					if tMin < 0 {
						tMin = 0
					}
					for repIndex > tMin && start > startLimit && src[repIndex-1] == src[start-1] && seq.matchLen < maxMatchLength-zstdMinMatch-1 {
						repIndex--
						start--
						seq.matchLen++
					}
					addLiterals(&seq, start)

					
					seq.offset = 1
					if debugSequences {
						println("repeat sequence", seq, "next s:", s)
					}
					blk.sequences = append(blk.sequences, seq)

					
					index0 := s + repOff
					s += lenght + repOff

					nextEmit = s
					if s >= sLimit {
						if debugEncoder {
							println("repeat ended", s, lenght)

						}
						break encodeLoop
					}
					
					for index0 < s-1 {
						cv0 := load6432(src, index0)
						cv1 := cv0 >> 8
						h0 := hashLen(cv0, betterLongTableBits, betterLongLen)
						off := index0 + e.cur
						e.longTable[h0] = prevEntry{offset: off, prev: e.longTable[h0].offset}
						e.table[hashLen(cv1, betterShortTableBits, betterShortLen)] = tableEntry{offset: off + 1, val: uint32(cv1)}
						index0 += 2
					}
					cv = load6432(src, s)
					continue
				}
				const repOff2 = 1

				
				
				
				if false && repIndex >= 0 && load6432(src, repIndex) == load6432(src, s+repOff) {
					
					var seq seq
					lenght := 8 + e.matchlen(s+8+repOff2, repIndex+8, src)

					seq.matchLen = uint32(lenght - zstdMinMatch)

					
					
					start := s + repOff2
					
					
					startLimit := nextEmit + 1

					tMin := s - e.maxMatchOff
					if tMin < 0 {
						tMin = 0
					}
					for repIndex > tMin && start > startLimit && src[repIndex-1] == src[start-1] && seq.matchLen < maxMatchLength-zstdMinMatch-1 {
						repIndex--
						start--
						seq.matchLen++
					}
					addLiterals(&seq, start)

					
					seq.offset = 2
					if debugSequences {
						println("repeat sequence 2", seq, "next s:", s)
					}
					blk.sequences = append(blk.sequences, seq)

					index0 := s + repOff2
					s += lenght + repOff2
					nextEmit = s
					if s >= sLimit {
						if debugEncoder {
							println("repeat ended", s, lenght)

						}
						break encodeLoop
					}

					
					for index0 < s-1 {
						cv0 := load6432(src, index0)
						cv1 := cv0 >> 8
						h0 := hashLen(cv0, betterLongTableBits, betterLongLen)
						off := index0 + e.cur
						e.longTable[h0] = prevEntry{offset: off, prev: e.longTable[h0].offset}
						e.table[hashLen(cv1, betterShortTableBits, betterShortLen)] = tableEntry{offset: off + 1, val: uint32(cv1)}
						index0 += 2
					}
					cv = load6432(src, s)
					
					offset1, offset2 = offset2, offset1
					continue
				}
			}
			
			coffsetL := candidateL.offset - e.cur
			coffsetLP := candidateL.prev - e.cur

			
			if s-coffsetL < e.maxMatchOff && cv == load6432(src, coffsetL) {
				
				matched = e.matchlen(s+8, coffsetL+8, src) + 8
				t = coffsetL
				if debugAsserts && s <= t {
					panic(fmt.Sprintf("s (%d) <= t (%d)", s, t))
				}
				if debugAsserts && s-t > e.maxMatchOff {
					panic("s - t >e.maxMatchOff")
				}
				if debugMatches {
					println("long match")
				}

				if s-coffsetLP < e.maxMatchOff && cv == load6432(src, coffsetLP) {
					
					prevMatch := e.matchlen(s+8, coffsetLP+8, src) + 8
					if prevMatch > matched {
						matched = prevMatch
						t = coffsetLP
					}
					if debugAsserts && s <= t {
						panic(fmt.Sprintf("s (%d) <= t (%d)", s, t))
					}
					if debugAsserts && s-t > e.maxMatchOff {
						panic("s - t >e.maxMatchOff")
					}
					if debugMatches {
						println("long match")
					}
				}
				break
			}

			
			if s-coffsetLP < e.maxMatchOff && cv == load6432(src, coffsetLP) {
				
				matched = e.matchlen(s+8, coffsetLP+8, src) + 8
				t = coffsetLP
				if debugAsserts && s <= t {
					panic(fmt.Sprintf("s (%d) <= t (%d)", s, t))
				}
				if debugAsserts && s-t > e.maxMatchOff {
					panic("s - t >e.maxMatchOff")
				}
				if debugMatches {
					println("long match")
				}
				break
			}

			coffsetS := candidateS.offset - e.cur

			
			if s-coffsetS < e.maxMatchOff && uint32(cv) == candidateS.val {
				
				matched = e.matchlen(s+4, coffsetS+4, src) + 4

				
				const checkAt = 1
				cv := load6432(src, s+checkAt)
				nextHashL = hashLen(cv, betterLongTableBits, betterLongLen)
				candidateL = e.longTable[nextHashL]
				coffsetL = candidateL.offset - e.cur

				
				e.longTable[nextHashL] = prevEntry{offset: s + checkAt + e.cur, prev: candidateL.offset}
				if s-coffsetL < e.maxMatchOff && cv == load6432(src, coffsetL) {
					
					matchedNext := e.matchlen(s+8+checkAt, coffsetL+8, src) + 8
					if matchedNext > matched {
						t = coffsetL
						s += checkAt
						matched = matchedNext
						if debugMatches {
							println("long match (after short)")
						}
						break
					}
				}

				
				coffsetL = candidateL.prev - e.cur
				if s-coffsetL < e.maxMatchOff && cv == load6432(src, coffsetL) {
					
					matchedNext := e.matchlen(s+8+checkAt, coffsetL+8, src) + 8
					if matchedNext > matched {
						t = coffsetL
						s += checkAt
						matched = matchedNext
						if debugMatches {
							println("prev long match (after short)")
						}
						break
					}
				}
				t = coffsetS
				if debugAsserts && s <= t {
					panic(fmt.Sprintf("s (%d) <= t (%d)", s, t))
				}
				if debugAsserts && s-t > e.maxMatchOff {
					panic("s - t >e.maxMatchOff")
				}
				if debugAsserts && t < 0 {
					panic("t<0")
				}
				if debugMatches {
					println("short match")
				}
				break
			}

			
			s += stepSize + ((s - nextEmit) >> (kSearchStrength - 1))
			if s >= sLimit {
				break encodeLoop
			}
			cv = load6432(src, s)
		}

		
		if s+matched < sLimit {
			
			
			
			
			const skipBeginning = 3

			nextHashL := hashLen(load6432(src, s+matched), betterLongTableBits, betterLongLen)
			s2 := s + skipBeginning
			cv := load3232(src, s2)
			candidateL := e.longTable[nextHashL]
			coffsetL := candidateL.offset - e.cur - matched + skipBeginning
			if coffsetL >= 0 && coffsetL < s2 && s2-coffsetL < e.maxMatchOff && cv == load3232(src, coffsetL) {
				
				matchedNext := e.matchlen(s2+4, coffsetL+4, src) + 4
				if matchedNext > matched {
					t = coffsetL
					s = s2
					matched = matchedNext
					if debugMatches {
						println("long match at end-of-match")
					}
				}
			}

			
			if true {
				coffsetL = candidateL.prev - e.cur - matched + skipBeginning
				if coffsetL >= 0 && coffsetL < s2 && s2-coffsetL < e.maxMatchOff && cv == load3232(src, coffsetL) {
					
					matchedNext := e.matchlen(s2+4, coffsetL+4, src) + 4
					if matchedNext > matched {
						t = coffsetL
						s = s2
						matched = matchedNext
						if debugMatches {
							println("prev long match at end-of-match")
						}
					}
				}
			}
		}
		
		offset2 = offset1
		offset1 = s - t

		if debugAsserts && s <= t {
			panic(fmt.Sprintf("s (%d) <= t (%d)", s, t))
		}

		if debugAsserts && canRepeat && int(offset1) > len(src) {
			panic("invalid offset")
		}

		
		l := matched

		
		tMin := s - e.maxMatchOff
		if tMin < 0 {
			tMin = 0
		}
		for t > tMin && s > nextEmit && src[t-1] == src[s-1] && l < maxMatchLength {
			s--
			t--
			l++
		}

		
		var seq seq
		seq.litLen = uint32(s - nextEmit)
		seq.matchLen = uint32(l - zstdMinMatch)
		if seq.litLen > 0 {
			blk.literals = append(blk.literals, src[nextEmit:s]...)
		}
		seq.offset = uint32(s-t) + 3
		s += l
		if debugSequences {
			println("sequence", seq, "next s:", s)
		}
		blk.sequences = append(blk.sequences, seq)
		nextEmit = s
		if s >= sLimit {
			break encodeLoop
		}

		
		index0 := s - l + 1
		for index0 < s-1 {
			cv0 := load6432(src, index0)
			cv1 := cv0 >> 8
			h0 := hashLen(cv0, betterLongTableBits, betterLongLen)
			off := index0 + e.cur
			e.longTable[h0] = prevEntry{offset: off, prev: e.longTable[h0].offset}
			e.table[hashLen(cv1, betterShortTableBits, betterShortLen)] = tableEntry{offset: off + 1, val: uint32(cv1)}
			index0 += 2
		}

		cv = load6432(src, s)
		if !canRepeat {
			continue
		}

		
		for {
			o2 := s - offset2
			if load3232(src, o2) != uint32(cv) {
				
				break
			}

			
			nextHashL := hashLen(cv, betterLongTableBits, betterLongLen)
			nextHashS := hashLen(cv, betterShortTableBits, betterShortLen)

			
			
			l := 4 + e.matchlen(s+4, o2+4, src)

			e.longTable[nextHashL] = prevEntry{offset: s + e.cur, prev: e.longTable[nextHashL].offset}
			e.table[nextHashS] = tableEntry{offset: s + e.cur, val: uint32(cv)}
			seq.matchLen = uint32(l) - zstdMinMatch
			seq.litLen = 0

			
			seq.offset = 1
			s += l
			nextEmit = s
			if debugSequences {
				println("sequence", seq, "next s:", s)
			}
			blk.sequences = append(blk.sequences, seq)

			
			offset1, offset2 = offset2, offset1
			if s >= sLimit {
				
				break encodeLoop
			}
			cv = load6432(src, s)
		}
	}

	if int(nextEmit) < len(src) {
		blk.literals = append(blk.literals, src[nextEmit:]...)
		blk.extraLits = len(src) - int(nextEmit)
	}
	blk.recentOffsets[0] = uint32(offset1)
	blk.recentOffsets[1] = uint32(offset2)
	if debugEncoder {
		println("returning, recent offsets:", blk.recentOffsets, "extra literals:", blk.extraLits)
	}
}




func (e *betterFastEncoder) EncodeNoHist(blk *blockEnc, src []byte) {
	e.ensureHist(len(src))
	e.Encode(blk, src)
}


func (e *betterFastEncoderDict) Encode(blk *blockEnc, src []byte) {
	const (
		
		
		inputMargin            = 8 + 2
		minNonLiteralBlockSize = 16
	)

	
	for e.cur >= e.bufferReset-int32(len(e.hist)) {
		if len(e.hist) == 0 {
			for i := range e.table[:] {
				e.table[i] = tableEntry{}
			}
			for i := range e.longTable[:] {
				e.longTable[i] = prevEntry{}
			}
			e.cur = e.maxMatchOff
			e.allDirty = true
			break
		}
		
		minOff := e.cur + int32(len(e.hist)) - e.maxMatchOff
		for i := range e.table[:] {
			v := e.table[i].offset
			if v < minOff {
				v = 0
			} else {
				v = v - e.cur + e.maxMatchOff
			}
			e.table[i].offset = v
		}
		for i := range e.longTable[:] {
			v := e.longTable[i].offset
			v2 := e.longTable[i].prev
			if v < minOff {
				v = 0
				v2 = 0
			} else {
				v = v - e.cur + e.maxMatchOff
				if v2 < minOff {
					v2 = 0
				} else {
					v2 = v2 - e.cur + e.maxMatchOff
				}
			}
			e.longTable[i] = prevEntry{
				offset: v,
				prev:   v2,
			}
		}
		e.allDirty = true
		e.cur = e.maxMatchOff
		break
	}

	s := e.addBlock(src)
	blk.size = len(src)
	if len(src) < minNonLiteralBlockSize {
		blk.extraLits = len(src)
		blk.literals = blk.literals[:len(src)]
		copy(blk.literals, src)
		return
	}

	
	src = e.hist
	sLimit := int32(len(src)) - inputMargin
	
	
	const stepSize = 1

	const kSearchStrength = 9

	
	nextEmit := s
	cv := load6432(src, s)

	
	offset1 := int32(blk.recentOffsets[0])
	offset2 := int32(blk.recentOffsets[1])

	addLiterals := func(s *seq, until int32) {
		if until == nextEmit {
			return
		}
		blk.literals = append(blk.literals, src[nextEmit:until]...)
		s.litLen = uint32(until - nextEmit)
	}
	if debugEncoder {
		println("recent offsets:", blk.recentOffsets)
	}

encodeLoop:
	for {
		var t int32
		
		canRepeat := len(blk.sequences) > 2
		var matched int32

		for {
			if debugAsserts && canRepeat && offset1 == 0 {
				panic("offset0 was 0")
			}

			nextHashL := hashLen(cv, betterLongTableBits, betterLongLen)
			nextHashS := hashLen(cv, betterShortTableBits, betterShortLen)
			candidateL := e.longTable[nextHashL]
			candidateS := e.table[nextHashS]

			const repOff = 1
			repIndex := s - offset1 + repOff
			off := s + e.cur
			e.longTable[nextHashL] = prevEntry{offset: off, prev: candidateL.offset}
			e.markLongShardDirty(nextHashL)
			e.table[nextHashS] = tableEntry{offset: off, val: uint32(cv)}
			e.markShortShardDirty(nextHashS)

			if canRepeat {
				if repIndex >= 0 && load3232(src, repIndex) == uint32(cv>>(repOff*8)) {
					
					var seq seq
					lenght := 4 + e.matchlen(s+4+repOff, repIndex+4, src)

					seq.matchLen = uint32(lenght - zstdMinMatch)

					
					
					start := s + repOff
					
					
					startLimit := nextEmit + 1

					tMin := s - e.maxMatchOff
					if tMin < 0 {
						tMin = 0
					}
					for repIndex > tMin && start > startLimit && src[repIndex-1] == src[start-1] && seq.matchLen < maxMatchLength-zstdMinMatch-1 {
						repIndex--
						start--
						seq.matchLen++
					}
					addLiterals(&seq, start)

					
					seq.offset = 1
					if debugSequences {
						println("repeat sequence", seq, "next s:", s)
					}
					blk.sequences = append(blk.sequences, seq)

					
					index0 := s + repOff
					s += lenght + repOff

					nextEmit = s
					if s >= sLimit {
						if debugEncoder {
							println("repeat ended", s, lenght)

						}
						break encodeLoop
					}
					
					for index0 < s-1 {
						cv0 := load6432(src, index0)
						cv1 := cv0 >> 8
						h0 := hashLen(cv0, betterLongTableBits, betterLongLen)
						off := index0 + e.cur
						e.longTable[h0] = prevEntry{offset: off, prev: e.longTable[h0].offset}
						e.markLongShardDirty(h0)
						h1 := hashLen(cv1, betterShortTableBits, betterShortLen)
						e.table[h1] = tableEntry{offset: off + 1, val: uint32(cv1)}
						e.markShortShardDirty(h1)
						index0 += 2
					}
					cv = load6432(src, s)
					continue
				}
				const repOff2 = 1

				
				
				
				if false && repIndex >= 0 && load6432(src, repIndex) == load6432(src, s+repOff) {
					
					var seq seq
					lenght := 8 + e.matchlen(s+8+repOff2, repIndex+8, src)

					seq.matchLen = uint32(lenght - zstdMinMatch)

					
					
					start := s + repOff2
					
					
					startLimit := nextEmit + 1

					tMin := s - e.maxMatchOff
					if tMin < 0 {
						tMin = 0
					}
					for repIndex > tMin && start > startLimit && src[repIndex-1] == src[start-1] && seq.matchLen < maxMatchLength-zstdMinMatch-1 {
						repIndex--
						start--
						seq.matchLen++
					}
					addLiterals(&seq, start)

					
					seq.offset = 2
					if debugSequences {
						println("repeat sequence 2", seq, "next s:", s)
					}
					blk.sequences = append(blk.sequences, seq)

					index0 := s + repOff2
					s += lenght + repOff2
					nextEmit = s
					if s >= sLimit {
						if debugEncoder {
							println("repeat ended", s, lenght)

						}
						break encodeLoop
					}

					
					for index0 < s-1 {
						cv0 := load6432(src, index0)
						cv1 := cv0 >> 8
						h0 := hashLen(cv0, betterLongTableBits, betterLongLen)
						off := index0 + e.cur
						e.longTable[h0] = prevEntry{offset: off, prev: e.longTable[h0].offset}
						e.markLongShardDirty(h0)
						h1 := hashLen(cv1, betterShortTableBits, betterShortLen)
						e.table[h1] = tableEntry{offset: off + 1, val: uint32(cv1)}
						e.markShortShardDirty(h1)
						index0 += 2
					}
					cv = load6432(src, s)
					
					offset1, offset2 = offset2, offset1
					continue
				}
			}
			
			coffsetL := candidateL.offset - e.cur
			coffsetLP := candidateL.prev - e.cur

			
			if s-coffsetL < e.maxMatchOff && cv == load6432(src, coffsetL) {
				
				matched = e.matchlen(s+8, coffsetL+8, src) + 8
				t = coffsetL
				if debugAsserts && s <= t {
					panic(fmt.Sprintf("s (%d) <= t (%d)", s, t))
				}
				if debugAsserts && s-t > e.maxMatchOff {
					panic("s - t >e.maxMatchOff")
				}
				if debugMatches {
					println("long match")
				}

				if s-coffsetLP < e.maxMatchOff && cv == load6432(src, coffsetLP) {
					
					prevMatch := e.matchlen(s+8, coffsetLP+8, src) + 8
					if prevMatch > matched {
						matched = prevMatch
						t = coffsetLP
					}
					if debugAsserts && s <= t {
						panic(fmt.Sprintf("s (%d) <= t (%d)", s, t))
					}
					if debugAsserts && s-t > e.maxMatchOff {
						panic("s - t >e.maxMatchOff")
					}
					if debugMatches {
						println("long match")
					}
				}
				break
			}

			
			if s-coffsetLP < e.maxMatchOff && cv == load6432(src, coffsetLP) {
				
				matched = e.matchlen(s+8, coffsetLP+8, src) + 8
				t = coffsetLP
				if debugAsserts && s <= t {
					panic(fmt.Sprintf("s (%d) <= t (%d)", s, t))
				}
				if debugAsserts && s-t > e.maxMatchOff {
					panic("s - t >e.maxMatchOff")
				}
				if debugMatches {
					println("long match")
				}
				break
			}

			coffsetS := candidateS.offset - e.cur

			
			if s-coffsetS < e.maxMatchOff && uint32(cv) == candidateS.val {
				
				matched = e.matchlen(s+4, coffsetS+4, src) + 4

				
				const checkAt = 1
				cv := load6432(src, s+checkAt)
				nextHashL = hashLen(cv, betterLongTableBits, betterLongLen)
				candidateL = e.longTable[nextHashL]
				coffsetL = candidateL.offset - e.cur

				
				e.longTable[nextHashL] = prevEntry{offset: s + checkAt + e.cur, prev: candidateL.offset}
				e.markLongShardDirty(nextHashL)
				if s-coffsetL < e.maxMatchOff && cv == load6432(src, coffsetL) {
					
					matchedNext := e.matchlen(s+8+checkAt, coffsetL+8, src) + 8
					if matchedNext > matched {
						t = coffsetL
						s += checkAt
						matched = matchedNext
						if debugMatches {
							println("long match (after short)")
						}
						break
					}
				}

				
				coffsetL = candidateL.prev - e.cur
				if s-coffsetL < e.maxMatchOff && cv == load6432(src, coffsetL) {
					
					matchedNext := e.matchlen(s+8+checkAt, coffsetL+8, src) + 8
					if matchedNext > matched {
						t = coffsetL
						s += checkAt
						matched = matchedNext
						if debugMatches {
							println("prev long match (after short)")
						}
						break
					}
				}
				t = coffsetS
				if debugAsserts && s <= t {
					panic(fmt.Sprintf("s (%d) <= t (%d)", s, t))
				}
				if debugAsserts && s-t > e.maxMatchOff {
					panic("s - t >e.maxMatchOff")
				}
				if debugAsserts && t < 0 {
					panic("t<0")
				}
				if debugMatches {
					println("short match")
				}
				break
			}

			
			s += stepSize + ((s - nextEmit) >> (kSearchStrength - 1))
			if s >= sLimit {
				break encodeLoop
			}
			cv = load6432(src, s)
		}
		
		if s+matched < sLimit {
			nextHashL := hashLen(load6432(src, s+matched), betterLongTableBits, betterLongLen)
			cv := load3232(src, s)
			candidateL := e.longTable[nextHashL]
			coffsetL := candidateL.offset - e.cur - matched
			if coffsetL >= 0 && coffsetL < s && s-coffsetL < e.maxMatchOff && cv == load3232(src, coffsetL) {
				
				matchedNext := e.matchlen(s+4, coffsetL+4, src) + 4
				if matchedNext > matched {
					t = coffsetL
					matched = matchedNext
					if debugMatches {
						println("long match at end-of-match")
					}
				}
			}

			
			if true {
				coffsetL = candidateL.prev - e.cur - matched
				if coffsetL >= 0 && coffsetL < s && s-coffsetL < e.maxMatchOff && cv == load3232(src, coffsetL) {
					
					matchedNext := e.matchlen(s+4, coffsetL+4, src) + 4
					if matchedNext > matched {
						t = coffsetL
						matched = matchedNext
						if debugMatches {
							println("prev long match at end-of-match")
						}
					}
				}
			}
		}
		
		offset2 = offset1
		offset1 = s - t

		if debugAsserts && s <= t {
			panic(fmt.Sprintf("s (%d) <= t (%d)", s, t))
		}

		if debugAsserts && canRepeat && int(offset1) > len(src) {
			panic("invalid offset")
		}

		
		l := matched

		
		tMin := s - e.maxMatchOff
		if tMin < 0 {
			tMin = 0
		}
		for t > tMin && s > nextEmit && src[t-1] == src[s-1] && l < maxMatchLength {
			s--
			t--
			l++
		}

		
		var seq seq
		seq.litLen = uint32(s - nextEmit)
		seq.matchLen = uint32(l - zstdMinMatch)
		if seq.litLen > 0 {
			blk.literals = append(blk.literals, src[nextEmit:s]...)
		}
		seq.offset = uint32(s-t) + 3
		s += l
		if debugSequences {
			println("sequence", seq, "next s:", s)
		}
		blk.sequences = append(blk.sequences, seq)
		nextEmit = s
		if s >= sLimit {
			break encodeLoop
		}

		
		index0 := s - l + 1
		for index0 < s-1 {
			cv0 := load6432(src, index0)
			cv1 := cv0 >> 8
			h0 := hashLen(cv0, betterLongTableBits, betterLongLen)
			off := index0 + e.cur
			e.longTable[h0] = prevEntry{offset: off, prev: e.longTable[h0].offset}
			e.markLongShardDirty(h0)
			h1 := hashLen(cv1, betterShortTableBits, betterShortLen)
			e.table[h1] = tableEntry{offset: off + 1, val: uint32(cv1)}
			e.markShortShardDirty(h1)
			index0 += 2
		}

		cv = load6432(src, s)
		if !canRepeat {
			continue
		}

		
		for {
			o2 := s - offset2
			if load3232(src, o2) != uint32(cv) {
				
				break
			}

			
			nextHashL := hashLen(cv, betterLongTableBits, betterLongLen)
			nextHashS := hashLen(cv, betterShortTableBits, betterShortLen)

			
			
			l := 4 + e.matchlen(s+4, o2+4, src)

			e.longTable[nextHashL] = prevEntry{offset: s + e.cur, prev: e.longTable[nextHashL].offset}
			e.markLongShardDirty(nextHashL)
			e.table[nextHashS] = tableEntry{offset: s + e.cur, val: uint32(cv)}
			e.markShortShardDirty(nextHashS)
			seq.matchLen = uint32(l) - zstdMinMatch
			seq.litLen = 0

			
			seq.offset = 1
			s += l
			nextEmit = s
			if debugSequences {
				println("sequence", seq, "next s:", s)
			}
			blk.sequences = append(blk.sequences, seq)

			
			offset1, offset2 = offset2, offset1
			if s >= sLimit {
				
				break encodeLoop
			}
			cv = load6432(src, s)
		}
	}

	if int(nextEmit) < len(src) {
		blk.literals = append(blk.literals, src[nextEmit:]...)
		blk.extraLits = len(src) - int(nextEmit)
	}
	blk.recentOffsets[0] = uint32(offset1)
	blk.recentOffsets[1] = uint32(offset2)
	if debugEncoder {
		println("returning, recent offsets:", blk.recentOffsets, "extra literals:", blk.extraLits)
	}
}


func (e *betterFastEncoder) Reset(d *dict, singleBlock bool) {
	e.resetBase(d, singleBlock)
	if d != nil {
		panic("betterFastEncoder: Reset with dict")
	}
}


func (e *betterFastEncoderDict) Reset(d *dict, singleBlock bool) {
	e.resetBase(d, singleBlock)
	if d == nil {
		return
	}
	
	if len(e.dictTable) != len(e.table) || d.id != e.lastDictID {
		if len(e.dictTable) != len(e.table) {
			e.dictTable = make([]tableEntry, len(e.table))
		}
		end := int32(len(d.content)) - 8 + e.maxMatchOff
		for i := e.maxMatchOff; i < end; i += 4 {
			const hashLog = betterShortTableBits

			cv := load6432(d.content, i-e.maxMatchOff)
			nextHash := hashLen(cv, hashLog, betterShortLen)      
			nextHash1 := hashLen(cv>>8, hashLog, betterShortLen)  
			nextHash2 := hashLen(cv>>16, hashLog, betterShortLen) 
			nextHash3 := hashLen(cv>>24, hashLog, betterShortLen) 
			e.dictTable[nextHash] = tableEntry{
				val:    uint32(cv),
				offset: i,
			}
			e.dictTable[nextHash1] = tableEntry{
				val:    uint32(cv >> 8),
				offset: i + 1,
			}
			e.dictTable[nextHash2] = tableEntry{
				val:    uint32(cv >> 16),
				offset: i + 2,
			}
			e.dictTable[nextHash3] = tableEntry{
				val:    uint32(cv >> 24),
				offset: i + 3,
			}
		}
		e.lastDictID = d.id
		e.allDirty = true
	}

	
	if len(e.dictLongTable) != len(e.longTable) || d.id != e.lastDictID {
		if len(e.dictLongTable) != len(e.longTable) {
			e.dictLongTable = make([]prevEntry, len(e.longTable))
		}
		if len(d.content) >= 8 {
			cv := load6432(d.content, 0)
			h := hashLen(cv, betterLongTableBits, betterLongLen)
			e.dictLongTable[h] = prevEntry{
				offset: e.maxMatchOff,
				prev:   e.dictLongTable[h].offset,
			}

			end := int32(len(d.content)) - 8 + e.maxMatchOff
			off := 8 
			for i := e.maxMatchOff + 1; i < end; i++ {
				cv = cv>>8 | (uint64(d.content[off]) << 56)
				h := hashLen(cv, betterLongTableBits, betterLongLen)
				e.dictLongTable[h] = prevEntry{
					offset: i,
					prev:   e.dictLongTable[h].offset,
				}
				off++
			}
		}
		e.lastDictID = d.id
		e.allDirty = true
	}

	
	{
		dirtyShardCnt := 0
		if !e.allDirty {
			for i := range e.shortTableShardDirty {
				if e.shortTableShardDirty[i] {
					dirtyShardCnt++
				}
			}
		}
		const shardCnt = betterShortTableShardCnt
		const shardSize = betterShortTableShardSize
		if e.allDirty || dirtyShardCnt > shardCnt*4/6 {
			copy(e.table[:], e.dictTable)
			for i := range e.shortTableShardDirty {
				e.shortTableShardDirty[i] = false
			}
		} else {
			for i := range e.shortTableShardDirty {
				if !e.shortTableShardDirty[i] {
					continue
				}

				copy(e.table[i*shardSize:(i+1)*shardSize], e.dictTable[i*shardSize:(i+1)*shardSize])
				e.shortTableShardDirty[i] = false
			}
		}
	}
	{
		dirtyShardCnt := 0
		if !e.allDirty {
			for i := range e.shortTableShardDirty {
				if e.shortTableShardDirty[i] {
					dirtyShardCnt++
				}
			}
		}
		const shardCnt = betterLongTableShardCnt
		const shardSize = betterLongTableShardSize
		if e.allDirty || dirtyShardCnt > shardCnt*4/6 {
			copy(e.longTable[:], e.dictLongTable)
			for i := range e.longTableShardDirty {
				e.longTableShardDirty[i] = false
			}
		} else {
			for i := range e.longTableShardDirty {
				if !e.longTableShardDirty[i] {
					continue
				}

				copy(e.longTable[i*shardSize:(i+1)*shardSize], e.dictLongTable[i*shardSize:(i+1)*shardSize])
				e.longTableShardDirty[i] = false
			}
		}
	}
	e.cur = e.maxMatchOff
	e.allDirty = false
}

func (e *betterFastEncoderDict) markLongShardDirty(entryNum uint32) {
	e.longTableShardDirty[entryNum/betterLongTableShardSize] = true
}

func (e *betterFastEncoderDict) markShortShardDirty(entryNum uint32) {
	e.shortTableShardDirty[entryNum/betterShortTableShardSize] = true
}
