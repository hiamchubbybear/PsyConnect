package flate

import "fmt"

type fastEncL4 struct {
	fastGen
	table  [tableSize]tableEntry
	bTable [tableSize]tableEntry
}

func (e *fastEncL4) Encode(dst *tokens, src []byte) {
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
				e.bTable[i] = tableEntry{}
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
			v := e.bTable[i].offset
			if v <= minOff {
				v = 0
			} else {
				v = v - e.cur + maxMatchOffset
			}
			e.bTable[i].offset = v
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
	for {
		const skipLog = 6
		const doEvery = 1

		nextS := s
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
			e.bTable[nextHashL] = entry

			t = lCandidate.offset - e.cur
			if s-t < maxMatchOffset && uint32(cv) == load3232(src, lCandidate.offset-e.cur) {
				
				break
			}

			t = sCandidate.offset - e.cur
			if s-t < maxMatchOffset && uint32(cv) == load3232(src, sCandidate.offset-e.cur) {
				
				lCandidate = e.bTable[hash7(next, tableBits)]

				
				lOff := nextS - (lCandidate.offset - e.cur)
				if lOff < maxMatchOffset && load3232(src, lCandidate.offset-e.cur) == uint32(next) {
					l1, l2 := matchLen(src[s+4:], src[t+4:]), matchLen(src[nextS+4:], src[nextS-lOff+4:])
					if l2 > l1 {
						s = nextS
						t = lCandidate.offset - e.cur
					}
				}
				break
			}
			cv = next
		}

		
		
		

		
		l := e.matchlenLong(s+4, t+4, src) + 4

		
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
		if debugDeflate {
			if t >= s {
				panic("s-t")
			}
			if (s - t) > maxMatchOffset {
				panic(fmt.Sprintln("mmo", t))
			}
			if l < baseMatchLength {
				panic("bml")
			}
		}

		dst.AddMatchLong(l, uint32(s-t-baseMatchOffset))
		s += l
		nextEmit = s
		if nextS >= s {
			s = nextS + 1
		}

		if s >= sLimit {
			
			if int(s+8) < len(src) {
				cv := load6432(src, s)
				e.table[hashLen(cv, tableBits, hashShortBytes)] = tableEntry{offset: s + e.cur}
				e.bTable[hash7(cv, tableBits)] = tableEntry{offset: s + e.cur}
			}
			goto emitRemainder
		}

		
		if true {
			i := nextS
			if i < s-1 {
				cv := load6432(src, i)
				t := tableEntry{offset: i + e.cur}
				t2 := tableEntry{offset: t.offset + 1}
				e.bTable[hash7(cv, tableBits)] = t
				e.bTable[hash7(cv>>8, tableBits)] = t2
				e.table[hashLen(cv>>8, tableBits, hashShortBytes)] = t2

				i += 3
				for ; i < s-1; i += 3 {
					cv := load6432(src, i)
					t := tableEntry{offset: i + e.cur}
					t2 := tableEntry{offset: t.offset + 1}
					e.bTable[hash7(cv, tableBits)] = t
					e.bTable[hash7(cv>>8, tableBits)] = t2
					e.table[hashLen(cv>>8, tableBits, hashShortBytes)] = t2
				}
			}
		}

		
		
		x := load6432(src, s-1)
		o := e.cur + s - 1
		prevHashS := hashLen(x, tableBits, hashShortBytes)
		prevHashL := hash7(x, tableBits)
		e.table[prevHashS] = tableEntry{offset: o}
		e.bTable[prevHashL] = tableEntry{offset: o}
		cv = x >> 8
	}

emitRemainder:
	if int(nextEmit) < len(src) {
		
		if dst.n == 0 {
			return
		}

		emitLiteral(dst, src[nextEmit:])
	}
}
