package flate

import "fmt"

type fastEncL6 struct {
	fastGen
	table  [tableSize]tableEntry
	bTable [tableSize]tableEntryPrev
}

func (e *fastEncL6) Encode(dst *tokens, src []byte) {
	const (
		inputMargin            = 12 - 1
		minNonLiteralBlockSize = 1 + 1 + inputMargin
		hashShortBytes         = 4
	)
	if debugDeflate && e.cur < 0 {
		panic(fmt.Sprint("e.cur < 0: ", e.cur))
	}

	
	for e.cur >= bufferReset {
		if len(e.hist) == 0 {
			for i := range e.table[:] {
				e.table[i] = tableEntry{}
			}
			for i := range e.bTable[:] {
				e.bTable[i] = tableEntryPrev{}
			}
			e.cur = maxMatchOffset
			break
		}
		
		minOff := e.cur + int32(len(e.hist)) - maxMatchOffset
		for i := range e.table[:] {
			v := e.table[i].offset
			if v <= minOff {
				v = 0
			} else {
				v = v - e.cur + maxMatchOffset
			}
			e.table[i].offset = v
		}
		for i := range e.bTable[:] {
			v := e.bTable[i]
			if v.Cur.offset <= minOff {
				v.Cur.offset = 0
				v.Prev.offset = 0
			} else {
				v.Cur.offset = v.Cur.offset - e.cur + maxMatchOffset
				if v.Prev.offset <= minOff {
					v.Prev.offset = 0
				} else {
					v.Prev.offset = v.Prev.offset - e.cur + maxMatchOffset
				}
			}
			e.bTable[i] = v
		}
		e.cur = maxMatchOffset
	}

	s := e.addBlock(src)

	
	
	if len(src) < minNonLiteralBlockSize {
		
		
		dst.n = uint16(len(src))
		return
	}

	
	src = e.hist
	nextEmit := s

	
	
	
	sLimit := int32(len(src) - inputMargin)

	
	cv := load6432(src, s)
	
	repeat := int32(1)
	for {
		const skipLog = 7
		const doEvery = 1

		nextS := s
		var l int32
		var t int32
		for {
			nextHashS := hashLen(cv, tableBits, hashShortBytes)
			nextHashL := hash7(cv, tableBits)
			s = nextS
			nextS = s + doEvery + (s-nextEmit)>>skipLog
			if nextS > sLimit {
				goto emitRemainder
			}
			
			sCandidate := e.table[nextHashS]
			lCandidate := e.bTable[nextHashL]
			next := load6432(src, nextS)
			entry := tableEntry{offset: s + e.cur}
			e.table[nextHashS] = entry
			eLong := &e.bTable[nextHashL]
			eLong.Cur, eLong.Prev = entry, eLong.Cur

			
			nextHashS = hashLen(next, tableBits, hashShortBytes)
			nextHashL = hash7(next, tableBits)

			t = lCandidate.Cur.offset - e.cur
			if s-t < maxMatchOffset {
				if uint32(cv) == load3232(src, lCandidate.Cur.offset-e.cur) {
					

					
					e.table[nextHashS] = tableEntry{offset: nextS + e.cur}
					eLong := &e.bTable[nextHashL]
					eLong.Cur, eLong.Prev = tableEntry{offset: nextS + e.cur}, eLong.Cur

					
					t2 := lCandidate.Prev.offset - e.cur
					if s-t2 < maxMatchOffset && uint32(cv) == load3232(src, lCandidate.Prev.offset-e.cur) {
						l = e.matchlen(s+4, t+4, src) + 4
						ml1 := e.matchlen(s+4, t2+4, src) + 4
						if ml1 > l {
							t = t2
							l = ml1
							break
						}
					}
					break
				}
				
				t = lCandidate.Prev.offset - e.cur
				if s-t < maxMatchOffset && uint32(cv) == load3232(src, lCandidate.Prev.offset-e.cur) {
					
					e.table[nextHashS] = tableEntry{offset: nextS + e.cur}
					eLong := &e.bTable[nextHashL]
					eLong.Cur, eLong.Prev = tableEntry{offset: nextS + e.cur}, eLong.Cur
					break
				}
			}

			t = sCandidate.offset - e.cur
			if s-t < maxMatchOffset && uint32(cv) == load3232(src, sCandidate.offset-e.cur) {
				
				l = e.matchlen(s+4, t+4, src) + 4

				
				lCandidate = e.bTable[nextHashL]

				
				e.table[nextHashS] = tableEntry{offset: nextS + e.cur}
				eLong := &e.bTable[nextHashL]
				eLong.Cur, eLong.Prev = tableEntry{offset: nextS + e.cur}, eLong.Cur

				
				const repOff = 1
				t2 := s - repeat + repOff
				if load3232(src, t2) == uint32(cv>>(8*repOff)) {
					ml := e.matchlen(s+4+repOff, t2+4, src) + 4
					if ml > l {
						t = t2
						l = ml
						s += repOff
						
						break
					}
				}

				
				t2 = lCandidate.Cur.offset - e.cur
				if nextS-t2 < maxMatchOffset {
					if load3232(src, lCandidate.Cur.offset-e.cur) == uint32(next) {
						ml := e.matchlen(nextS+4, t2+4, src) + 4
						if ml > l {
							t = t2
							s = nextS
							l = ml
							
						}
					}
					
					t2 = lCandidate.Prev.offset - e.cur
					if nextS-t2 < maxMatchOffset && load3232(src, lCandidate.Prev.offset-e.cur) == uint32(next) {
						ml := e.matchlen(nextS+4, t2+4, src) + 4
						if ml > l {
							t = t2
							s = nextS
							l = ml
							break
						}
					}
				}
				break
			}
			cv = next
		}

		
		
		

		
		if l == 0 {
			l = e.matchlenLong(s+4, t+4, src) + 4
		} else if l == maxMatchLength {
			l += e.matchlenLong(s+l, t+l, src)
		}

		
		if sAt := s + l; sAt < sLimit {
			
			
			
			
			
			const skipBeginning = 2
			eLong := &e.bTable[hash7(load6432(src, sAt), tableBits)]
			
			t2 := eLong.Cur.offset - e.cur - l + skipBeginning
			s2 := s + skipBeginning
			off := s2 - t2
			if off < maxMatchOffset {
				if off > 0 && t2 >= 0 {
					if l2 := e.matchlenLong(s2, t2, src); l2 > l {
						t = t2
						l = l2
						s = s2
					}
				}
				
				t2 = eLong.Prev.offset - e.cur - l + skipBeginning
				off := s2 - t2
				if off > 0 && off < maxMatchOffset && t2 >= 0 {
					if l2 := e.matchlenLong(s2, t2, src); l2 > l {
						t = t2
						l = l2
						s = s2
					}
				}
			}
		}

		
		for t > 0 && s > nextEmit && src[t-1] == src[s-1] {
			s--
			t--
			l++
		}
		if nextEmit < s {
			if false {
				emitLiteral(dst, src[nextEmit:s])
			} else {
				for _, v := range src[nextEmit:s] {
					dst.tokens[dst.n] = token(v)
					dst.litHist[v]++
					dst.n++
				}
			}
		}
		if false {
			if t >= s {
				panic(fmt.Sprintln("s-t", s, t))
			}
			if (s - t) > maxMatchOffset {
				panic(fmt.Sprintln("mmo", s-t))
			}
			if l < baseMatchLength {
				panic("bml")
			}
		}

		dst.AddMatchLong(l, uint32(s-t-baseMatchOffset))
		repeat = s - t
		s += l
		nextEmit = s
		if nextS >= s {
			s = nextS + 1
		}

		if s >= sLimit {
			
			for i := nextS + 1; i < int32(len(src))-8; i += 2 {
				cv := load6432(src, i)
				e.table[hashLen(cv, tableBits, hashShortBytes)] = tableEntry{offset: i + e.cur}
				eLong := &e.bTable[hash7(cv, tableBits)]
				eLong.Cur, eLong.Prev = tableEntry{offset: i + e.cur}, eLong.Cur
			}
			goto emitRemainder
		}

		
		if true {
			for i := nextS + 1; i < s-1; i += 2 {
				cv := load6432(src, i)
				t := tableEntry{offset: i + e.cur}
				t2 := tableEntry{offset: t.offset + 1}
				eLong := &e.bTable[hash7(cv, tableBits)]
				eLong2 := &e.bTable[hash7(cv>>8, tableBits)]
				e.table[hashLen(cv, tableBits, hashShortBytes)] = t
				eLong.Cur, eLong.Prev = t, eLong.Cur
				eLong2.Cur, eLong2.Prev = t2, eLong2.Cur
			}
		}

		
		
		cv = load6432(src, s)
	}

emitRemainder:
	if int(nextEmit) < len(src) {
		
		if dst.n == 0 {
			return
		}

		emitLiteral(dst, src[nextEmit:])
	}
}
