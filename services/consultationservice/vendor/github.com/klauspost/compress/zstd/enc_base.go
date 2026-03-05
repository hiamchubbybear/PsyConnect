package zstd

import (
	"fmt"
	"math/bits"

	"github.com/klauspost/compress/zstd/internal/xxhash"
)

const (
	dictShardBits = 6
)

type fastBase struct {
	
	cur int32
	
	maxMatchOff int32
	bufferReset int32
	hist        []byte
	crc         *xxhash.Digest
	tmp         [8]byte
	blk         *blockEnc
	lastDictID  uint32
	lowMem      bool
}


func (e *fastBase) CRC() *xxhash.Digest {
	return e.crc
}


func (e *fastBase) AppendCRC(dst []byte) []byte {
	crc := e.crc.Sum(e.tmp[:0])
	dst = append(dst, crc[7], crc[6], crc[5], crc[4])
	return dst
}



func (e *fastBase) WindowSize(size int64) int32 {
	if size > 0 && size < int64(e.maxMatchOff) {
		b := int32(1) << uint(bits.Len(uint(size)))
		
		if b < 1024 {
			b = 1024
		}
		return b
	}
	return e.maxMatchOff
}


func (e *fastBase) Block() *blockEnc {
	return e.blk
}

func (e *fastBase) addBlock(src []byte) int32 {
	if debugAsserts && e.cur > e.bufferReset {
		panic(fmt.Sprintf("ecur (%d) > buffer reset (%d)", e.cur, e.bufferReset))
	}
	
	if len(e.hist)+len(src) > cap(e.hist) {
		if cap(e.hist) == 0 {
			e.ensureHist(len(src))
		} else {
			if cap(e.hist) < int(e.maxMatchOff+maxCompressedBlockSize) {
				panic(fmt.Errorf("unexpected buffer cap %d, want at least %d with window %d", cap(e.hist), e.maxMatchOff+maxCompressedBlockSize, e.maxMatchOff))
			}
			
			offset := int32(len(e.hist)) - e.maxMatchOff
			copy(e.hist[0:e.maxMatchOff], e.hist[offset:])
			e.cur += offset
			e.hist = e.hist[:e.maxMatchOff]
		}
	}
	s := int32(len(e.hist))
	e.hist = append(e.hist, src...)
	return s
}


func (e *fastBase) ensureHist(n int) {
	if cap(e.hist) >= n {
		return
	}
	l := e.maxMatchOff
	if (e.lowMem && e.maxMatchOff > maxCompressedBlockSize) || e.maxMatchOff <= maxCompressedBlockSize {
		l += maxCompressedBlockSize
	} else {
		l += e.maxMatchOff
	}
	
	if l < 1<<20 && !e.lowMem {
		l = 1 << 20
	}
	
	if l < int32(n) {
		l = int32(n)
	}
	e.hist = make([]byte, 0, l)
}



func (e *fastBase) UseBlock(enc *blockEnc) {
	enc.reset(e.blk)
	e.blk = enc
}

func (e *fastBase) matchlen(s, t int32, src []byte) int32 {
	if debugAsserts {
		if s < 0 {
			err := fmt.Sprintf("s (%d) < 0", s)
			panic(err)
		}
		if t < 0 {
			err := fmt.Sprintf("s (%d) < 0", s)
			panic(err)
		}
		if s-t > e.maxMatchOff {
			err := fmt.Sprintf("s (%d) - t (%d) > maxMatchOff (%d)", s, t, e.maxMatchOff)
			panic(err)
		}
		if len(src)-int(s) > maxCompressedBlockSize {
			panic(fmt.Sprintf("len(src)-s (%d) > maxCompressedBlockSize (%d)", len(src)-int(s), maxCompressedBlockSize))
		}
	}
	return int32(matchLen(src[s:], src[t:]))
}


func (e *fastBase) resetBase(d *dict, singleBlock bool) {
	if e.blk == nil {
		e.blk = &blockEnc{lowMem: e.lowMem}
		e.blk.init()
	} else {
		e.blk.reset(nil)
	}
	e.blk.initNewEncode()
	if e.crc == nil {
		e.crc = xxhash.New()
	} else {
		e.crc.Reset()
	}
	e.blk.dictLitEnc = nil
	if d != nil {
		low := e.lowMem
		if singleBlock {
			e.lowMem = true
		}
		e.ensureHist(d.ContentSize() + maxCompressedBlockSize)
		e.lowMem = low
	}

	
	
	if e.cur < e.bufferReset {
		e.cur += e.maxMatchOff + int32(len(e.hist))
	}
	e.hist = e.hist[:0]
	if d != nil {
		
		for i, off := range d.offsets {
			e.blk.recentOffsets[i] = uint32(off)
			e.blk.prevRecentOffsets[i] = e.blk.recentOffsets[i]
		}
		
		e.blk.dictLitEnc = d.litEnc
		e.hist = append(e.hist, d.content...)
	}
}
