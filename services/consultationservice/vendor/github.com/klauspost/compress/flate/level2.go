package flate

import "fmt"




type fastEncL2 struct {
	fastGen
	table [bTableSize]tableEntry
}



func (e *fastEncL2) Encode(dst *tokens, src []byte) {
	const (
		inputMargin            = 12 - 1
		minNonLiteralBlockSize = 1 + 1 + inputMargin
		hashBytes              = 5
	)

	if debugDeflate && e.cur < 0 {
		panic(fmt.Sprint("e.cur < 0: ", e.cur))
	}

	
	for e.cur >= bufferReset {
		if len(e.hist) == 0 {
			for i := range e.table[:] {
				e.table[i] = tableEntry{}
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
		
		const skipLog = 5
		const doEvery = 2

		nextS := s
		var candidate tableEntry
		for {
			nextHash := hashLen(cv, bTableBits, hashBytes)
			s = nextS
			nextS = s + doEvery + (s-nextEmit)>>skipLog
			if nextS > sLimit {
				goto emitRemainder
			}
			candidate = e.table[nextHash]
			now := load6432(src, nextS)
			e.table[nextHash] = tableEntry{offset: s + e.cur}
			nextHash = hashLen(now, bTableBits, hashBytes)

			offset := s - (candidate.offset - e.cur)
			if offset < maxMatchOffset && uint32(cv) == load3232(src, candidate.offset-e.cur) {
				e.table[nextHash] = tableEntry{offset: nextS + e.cur}
				break
			}

			
			cv = now
			s = nextS
			nextS++
			candidate = e.table[nextHash]
			now >>= 8
			e.table[nextHash] = tableEntry{offset: s + e.cur}

			offset = s - (candidate.offset - e.cur)
			if offset < maxMatchOffset && uint32(cv) == load3232(src, candidate.offset-e.cur) {
				break
			}
			cv = now
		}

		
		
		

		
		
		
		
		
		
		
		
		for {
			
			

			
			t := candidate.offset - e.cur
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

			dst.AddMatchLong(l, uint32(s-t-baseMatchOffset))
			s += l
			nextEmit = s
			if nextS >= s {
				s = nextS + 1
			}

			if s >= sLimit {
				
				if int(s+l+8) < len(src) {
					cv := load6432(src, s)
					e.table[hashLen(cv, bTableBits, hashBytes)] = tableEntry{offset: s + e.cur}
				}
				goto emitRemainder
			}

			
			for i := s - l + 2; i < s-5; i += 7 {
				x := load6432(src, i)
				nextHash := hashLen(x, bTableBits, hashBytes)
				e.table[nextHash] = tableEntry{offset: e.cur + i}
				
				x >>= 16
				nextHash = hashLen(x, bTableBits, hashBytes)
				e.table[nextHash] = tableEntry{offset: e.cur + i + 2}
				
				x >>= 16
				nextHash = hashLen(x, bTableBits, hashBytes)
				e.table[nextHash] = tableEntry{offset: e.cur + i + 4}
			}

			
			
			
			
			
			
			x := load6432(src, s-2)
			o := e.cur + s - 2
			prevHash := hashLen(x, bTableBits, hashBytes)
			prevHash2 := hashLen(x>>8, bTableBits, hashBytes)
			e.table[prevHash] = tableEntry{offset: o}
			e.table[prevHash2] = tableEntry{offset: o + 1}
			currHash := hashLen(x>>16, bTableBits, hashBytes)
			candidate = e.table[currHash]
			e.table[currHash] = tableEntry{offset: o + 2}

			offset := s - (candidate.offset - e.cur)
			if offset > maxMatchOffset || uint32(x>>16) != load3232(src, candidate.offset-e.cur) {
				cv = x >> 24
				s++
				break
			}
		}
	}

emitRemainder:
	if int(nextEmit) < len(src) {
		
		if dst.n == 0 {
			return
		}

		emitLiteral(dst, src[nextEmit:])
	}
}
