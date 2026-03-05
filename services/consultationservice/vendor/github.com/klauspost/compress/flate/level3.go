package flate

import "fmt"


type fastEncL3 struct {
	fastGen
	table [1 << 16]tableEntryPrev
}


func (e *fastEncL3) Encode(dst *tokens, src []byte) {
	const (
		inputMargin            = 12 - 1
		minNonLiteralBlockSize = 1 + 1 + inputMargin
		tableBits              = 16
		tableSize              = 1 << tableBits
		hashBytes              = 5
	)

	if debugDeflate && e.cur < 0 {
		panic(fmt.Sprint("e.cur < 0: ", e.cur))
	}

	
	for e.cur >= bufferReset {
		if len(e.hist) == 0 {
			for i := range e.table[:] {
				e.table[i] = tableEntryPrev{}
			}
			e.cur = maxMatchOffset
			break
		}
		
		minOff := e.cur + int32(len(e.hist)) - maxMatchOffset
		for i := range e.table[:] {
			v := e.table[i]
			if v.Cur.offset <= minOff {
				v.Cur.offset = 0
			} else {
				v.Cur.offset = v.Cur.offset - e.cur + maxMatchOffset
			}
			if v.Prev.offset <= minOff {
				v.Prev.offset = 0
			} else {
				v.Prev.offset = v.Prev.offset - e.cur + maxMatchOffset
			}
			e.table[i] = v
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
		const skipLog = 7
		nextS := s
		var candidate tableEntry
		for {
			nextHash := hashLen(cv, tableBits, hashBytes)
			s = nextS
			nextS = s + 1 + (s-nextEmit)>>skipLog
			if nextS > sLimit {
				goto emitRemainder
			}
			candidates := e.table[nextHash]
			now := load6432(src, nextS)

			
			minOffset := e.cur + s - (maxMatchOffset - 4)
			e.table[nextHash] = tableEntryPrev{Prev: candidates.Cur, Cur: tableEntry{offset: s + e.cur}}

			
			candidate = candidates.Cur
			if candidate.offset < minOffset {
				cv = now
				
				continue
			}

			if uint32(cv) == load3232(src, candidate.offset-e.cur) {
				if candidates.Prev.offset < minOffset || uint32(cv) != load3232(src, candidates.Prev.offset-e.cur) {
					break
				}
				
				offset := s - (candidate.offset - e.cur)
				o2 := s - (candidates.Prev.offset - e.cur)
				l1, l2 := matchLen(src[s+4:], src[s-offset+4:]), matchLen(src[s+4:], src[s-o2+4:])
				if l2 > l1 {
					candidate = candidates.Prev
				}
				break
			} else {
				
				
				candidate = candidates.Prev
				if candidate.offset > minOffset && uint32(cv) == load3232(src, candidate.offset-e.cur) {
					break
				}
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
				t += l
				
				if int(t+8) < len(src) && t > 0 {
					cv = load6432(src, t)
					nextHash := hashLen(cv, tableBits, hashBytes)
					e.table[nextHash] = tableEntryPrev{
						Prev: e.table[nextHash].Cur,
						Cur:  tableEntry{offset: e.cur + t},
					}
				}
				goto emitRemainder
			}

			
			for i := s - l + 2; i < s-5; i += 6 {
				nextHash := hashLen(load6432(src, i), tableBits, hashBytes)
				e.table[nextHash] = tableEntryPrev{
					Prev: e.table[nextHash].Cur,
					Cur:  tableEntry{offset: e.cur + i}}
			}
			
			
			x := load6432(src, s-2)
			prevHash := hashLen(x, tableBits, hashBytes)

			e.table[prevHash] = tableEntryPrev{
				Prev: e.table[prevHash].Cur,
				Cur:  tableEntry{offset: e.cur + s - 2},
			}
			x >>= 8
			prevHash = hashLen(x, tableBits, hashBytes)

			e.table[prevHash] = tableEntryPrev{
				Prev: e.table[prevHash].Cur,
				Cur:  tableEntry{offset: e.cur + s - 1},
			}
			x >>= 8
			currHash := hashLen(x, tableBits, hashBytes)
			candidates := e.table[currHash]
			cv = x
			e.table[currHash] = tableEntryPrev{
				Prev: candidates.Cur,
				Cur:  tableEntry{offset: s + e.cur},
			}

			
			candidate = candidates.Cur
			minOffset := e.cur + s - (maxMatchOffset - 4)

			if candidate.offset > minOffset {
				if uint32(cv) == load3232(src, candidate.offset-e.cur) {
					
					continue
				}
				candidate = candidates.Prev
				if candidate.offset > minOffset && uint32(cv) == load3232(src, candidate.offset-e.cur) {
					
					continue
				}
			}
			cv = x >> 8
			s++
			break
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
