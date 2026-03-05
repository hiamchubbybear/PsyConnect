package flate

import (
	"io"
	"math"
	"sync"
)

const (
	maxStatelessBlock = math.MaxInt16
	
	maxStatelessDict = 8 << 10

	slTableBits  = 13
	slTableSize  = 1 << slTableBits
	slTableShift = 32 - slTableBits
)

type statelessWriter struct {
	dst    io.Writer
	closed bool
}

func (s *statelessWriter) Close() error {
	if s.closed {
		return nil
	}
	s.closed = true
	
	return StatelessDeflate(s.dst, nil, true, nil)
}

func (s *statelessWriter) Write(p []byte) (n int, err error) {
	err = StatelessDeflate(s.dst, p, false, nil)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (s *statelessWriter) Reset(w io.Writer) {
	s.dst = w
	s.closed = false
}






func NewStatelessWriter(dst io.Writer) io.WriteCloser {
	return &statelessWriter{dst: dst}
}


var bitWriterPool = sync.Pool{
	New: func() interface{} {
		return newHuffmanBitWriter(nil)
	},
}






func StatelessDeflate(out io.Writer, in []byte, eof bool, dict []byte) error {
	var dst tokens
	bw := bitWriterPool.Get().(*huffmanBitWriter)
	bw.reset(out)
	defer func() {
		
		bw.reset(nil)
		bitWriterPool.Put(bw)
	}()
	if eof && len(in) == 0 {
		
		
		bw.writeStoredHeader(0, true)
		bw.flush()
		return bw.err
	}

	
	if len(dict) > maxStatelessDict {
		dict = dict[len(dict)-maxStatelessDict:]
	}

	
	var inDict []byte

	for len(in) > 0 {
		todo := in
		if len(inDict) > 0 {
			if len(todo) > maxStatelessBlock-maxStatelessDict {
				todo = todo[:maxStatelessBlock-maxStatelessDict]
			}
		} else if len(todo) > maxStatelessBlock-len(dict) {
			todo = todo[:maxStatelessBlock-len(dict)]
		}
		inOrg := in
		in = in[len(todo):]
		uncompressed := todo
		if len(dict) > 0 {
			
			bufLen := len(todo) + len(dict)
			combined := make([]byte, bufLen)
			copy(combined, dict)
			copy(combined[len(dict):], todo)
			todo = combined
		}
		
		if len(inDict) == 0 {
			statelessEnc(&dst, todo, int16(len(dict)))
		} else {
			statelessEnc(&dst, inDict[:maxStatelessDict+len(todo)], maxStatelessDict)
		}
		isEof := eof && len(in) == 0

		if dst.n == 0 {
			bw.writeStoredHeader(len(uncompressed), isEof)
			if bw.err != nil {
				return bw.err
			}
			bw.writeBytes(uncompressed)
		} else if int(dst.n) > len(uncompressed)-len(uncompressed)>>4 {
			
			bw.writeBlockHuff(isEof, uncompressed, len(in) == 0)
		} else {
			bw.writeBlockDynamic(&dst, isEof, uncompressed, len(in) == 0)
		}
		if len(in) > 0 {
			
			inDict = inOrg[len(uncompressed)-maxStatelessDict:]
			dict = nil
			dst.Reset()
		}
		if bw.err != nil {
			return bw.err
		}
	}
	if !eof {
		
		bw.writeStoredHeader(0, false)
	}
	bw.flush()
	return bw.err
}

func hashSL(u uint32) uint32 {
	return (u * 0x1e35a7bd) >> slTableShift
}

func load3216(b []byte, i int16) uint32 {
	
	b = b[i:]
	b = b[:4]
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

func load6416(b []byte, i int16) uint64 {
	
	b = b[i:]
	b = b[:8]
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
}

func statelessEnc(dst *tokens, src []byte, startAt int16) {
	const (
		inputMargin            = 12 - 1
		minNonLiteralBlockSize = 1 + 1 + inputMargin
	)

	type tableEntry struct {
		offset int16
	}

	var table [slTableSize]tableEntry

	
	
	if len(src)-int(startAt) < minNonLiteralBlockSize {
		
		
		dst.n = 0
		return
	}
	
	if startAt > 0 {
		cv := load3232(src, 0)
		for i := int16(0); i < startAt; i++ {
			table[hashSL(cv)] = tableEntry{offset: i}
			cv = (cv >> 8) | (uint32(src[i+4]) << 24)
		}
	}

	s := startAt + 1
	nextEmit := startAt
	
	
	
	sLimit := int16(len(src) - inputMargin)

	
	cv := load3216(src, s)

	for {
		const skipLog = 5
		const doEvery = 2

		nextS := s
		var candidate tableEntry
		for {
			nextHash := hashSL(cv)
			candidate = table[nextHash]
			nextS = s + doEvery + (s-nextEmit)>>skipLog
			if nextS > sLimit || nextS <= 0 {
				goto emitRemainder
			}

			now := load6416(src, nextS)
			table[nextHash] = tableEntry{offset: s}
			nextHash = hashSL(uint32(now))

			if cv == load3216(src, candidate.offset) {
				table[nextHash] = tableEntry{offset: nextS}
				break
			}

			
			cv = uint32(now)
			s = nextS
			nextS++
			candidate = table[nextHash]
			now >>= 8
			table[nextHash] = tableEntry{offset: s}

			if cv == load3216(src, candidate.offset) {
				table[nextHash] = tableEntry{offset: nextS}
				break
			}
			cv = uint32(now)
			s = nextS
		}

		
		
		
		for {
			
			

			
			t := candidate.offset
			l := int16(matchLen(src[s+4:], src[t+4:]) + 4)

			
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

			
			dst.AddMatchLong(int32(l), uint32(s-t-baseMatchOffset))
			s += l
			nextEmit = s
			if nextS >= s {
				s = nextS + 1
			}
			if s >= sLimit {
				goto emitRemainder
			}

			
			
			
			
			
			
			x := load6416(src, s-2)
			o := s - 2
			prevHash := hashSL(uint32(x))
			table[prevHash] = tableEntry{offset: o}
			x >>= 16
			currHash := hashSL(uint32(x))
			candidate = table[currHash]
			table[currHash] = tableEntry{offset: o + 2}

			if uint32(x) != load3216(src, candidate.offset) {
				cv = uint32(x >> 8)
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
