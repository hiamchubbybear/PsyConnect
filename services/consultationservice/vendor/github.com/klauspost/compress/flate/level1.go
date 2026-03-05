package flate

import (
	"encoding/binary"
	"fmt"
	"math/bits"
)




type fastEncL1 struct {
	fastGen
	table [tableSize]tableEntry
}


func (e *fastEncL1) Encode(dst *tokens, src []byte) {
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
			nextHash := hashLen(cv, tableBits, hashBytes)
			candidate = e.table[nextHash]
			nextS = s + doEvery + (s-nextEmit)>>skipLog
			if nextS > sLimit {
				goto emitRemainder
			}

			now := load6432(src, nextS)
			e.table[nextHash] = tableEntry{offset: s + e.cur}
			nextHash = hashLen(now, tableBits, hashBytes)

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
				e.table[nextHash] = tableEntry{offset: nextS + e.cur}
				break
			}
			cv = now
			s = nextS
		}

		
		
		
		for {
			
			

			
			t := candidate.offset - e.cur
			var l = int32(4)
			if false {
				l = e.matchlenLong(s+4, t+4, src) + 4
			} else {
				
				a := src[s+4:]
				b := src[t+4:]
				for len(a) >= 8 {
					if diff := binary.LittleEndian.Uint64(a) ^ binary.LittleEndian.Uint64(b); diff != 0 {
						l += int32(bits.TrailingZeros64(diff) >> 3)
						break
					}
					l += 8
					a = a[8:]
					b = b[8:]
				}
				if len(a) < 8 {
					b = b[:len(a)]
					for i := range a {
						if a[i] != b[i] {
							break
						}
						l++
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
				dst.AddMatchLong(l, uint32(s-t-baseMatchOffset))
			} else {
				
				xoffset := uint32(s - t - baseMatchOffset)
				xlength := l
				oc := offsetCode(xoffset)
				xoffset |= oc << 16
				for xlength > 0 {
					xl := xlength
					if xl > 258 {
						if xl > 258+baseMatchLength {
							xl = 258
						} else {
							xl = 258 - baseMatchLength
						}
					}
					xlength -= xl
					xl -= baseMatchLength
					dst.extraHist[lengthCodes1[uint8(xl)]]++
					dst.offHist[oc]++
					dst.tokens[dst.n] = token(matchType | uint32(xl)<<lengthShift | xoffset)
					dst.n++
				}
			}
			s += l
			nextEmit = s
			if nextS >= s {
				s = nextS + 1
			}
			if s >= sLimit {
				
				if int(s+l+8) < len(src) {
					cv := load6432(src, s)
					e.table[hashLen(cv, tableBits, hashBytes)] = tableEntry{offset: s + e.cur}
				}
				goto emitRemainder
			}

			
			
			
			
			
			
			x := load6432(src, s-2)
			o := e.cur + s - 2
			prevHash := hashLen(x, tableBits, hashBytes)
			e.table[prevHash] = tableEntry{offset: o}
			x >>= 16
			currHash := hashLen(x, tableBits, hashBytes)
			candidate = e.table[currHash]
			e.table[currHash] = tableEntry{offset: o + 2}

			offset := s - (candidate.offset - e.cur)
			if offset > maxMatchOffset || uint32(x) != load3232(src, candidate.offset-e.cur) {
				cv = x >> 8
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
