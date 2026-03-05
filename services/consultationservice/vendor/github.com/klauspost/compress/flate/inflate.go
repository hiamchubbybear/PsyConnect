






package flate

import (
	"bufio"
	"compress/flate"
	"fmt"
	"io"
	"math/bits"
	"sync"
)

const (
	maxCodeLen     = 16 
	maxCodeLenMask = 15 
	
	
	
	maxNumLit  = 286
	maxNumDist = 30
	numCodes   = 19 

	debugDecode = false
)


type lengthExtra struct {
	length, extra uint8
}

var decCodeToLen = [32]lengthExtra{{length: 0x0, extra: 0x0}, {length: 0x1, extra: 0x0}, {length: 0x2, extra: 0x0}, {length: 0x3, extra: 0x0}, {length: 0x4, extra: 0x0}, {length: 0x5, extra: 0x0}, {length: 0x6, extra: 0x0}, {length: 0x7, extra: 0x0}, {length: 0x8, extra: 0x1}, {length: 0xa, extra: 0x1}, {length: 0xc, extra: 0x1}, {length: 0xe, extra: 0x1}, {length: 0x10, extra: 0x2}, {length: 0x14, extra: 0x2}, {length: 0x18, extra: 0x2}, {length: 0x1c, extra: 0x2}, {length: 0x20, extra: 0x3}, {length: 0x28, extra: 0x3}, {length: 0x30, extra: 0x3}, {length: 0x38, extra: 0x3}, {length: 0x40, extra: 0x4}, {length: 0x50, extra: 0x4}, {length: 0x60, extra: 0x4}, {length: 0x70, extra: 0x4}, {length: 0x80, extra: 0x5}, {length: 0xa0, extra: 0x5}, {length: 0xc0, extra: 0x5}, {length: 0xe0, extra: 0x5}, {length: 0xff, extra: 0x0}, {length: 0x0, extra: 0x0}, {length: 0x0, extra: 0x0}, {length: 0x0, extra: 0x0}}

var bitMask32 = [32]uint32{
	0, 1, 3, 7, 0xF, 0x1F, 0x3F, 0x7F, 0xFF,
	0x1FF, 0x3FF, 0x7FF, 0xFFF, 0x1FFF, 0x3FFF, 0x7FFF, 0xFFFF,
	0x1ffff, 0x3ffff, 0x7FFFF, 0xfFFFF, 0x1fFFFF, 0x3fFFFF, 0x7fFFFF, 0xffFFFF,
	0x1ffFFFF, 0x3ffFFFF, 0x7ffFFFF, 0xfffFFFF, 0x1fffFFFF, 0x3fffFFFF, 0x7fffFFFF,
} 


var fixedOnce sync.Once
var fixedHuffmanDecoder huffmanDecoder


type CorruptInputError = flate.CorruptInputError


type InternalError string

func (e InternalError) Error() string { return "flate: internal error: " + string(e) }




type ReadError = flate.ReadError




type WriteError = flate.WriteError




type Resetter interface {
	
	
	Reset(r io.Reader, dict []byte) error
}





















const (
	huffmanChunkBits  = 9
	huffmanNumChunks  = 1 << huffmanChunkBits
	huffmanCountMask  = 15
	huffmanValueShift = 4
)

type huffmanDecoder struct {
	maxRead  int                       
	chunks   *[huffmanNumChunks]uint16 
	links    [][]uint16                
	linkMask uint32                    
}






func (h *huffmanDecoder) init(lengths []int) bool {
	
	
	
	const sanity = false

	if h.chunks == nil {
		h.chunks = &[huffmanNumChunks]uint16{}
	}
	if h.maxRead != 0 {
		*h = huffmanDecoder{chunks: h.chunks, links: h.links}
	}

	
	
	var count [maxCodeLen]int
	var min, max int
	for _, n := range lengths {
		if n == 0 {
			continue
		}
		if min == 0 || n < min {
			min = n
		}
		if n > max {
			max = n
		}
		count[n&maxCodeLenMask]++
	}

	
	
	
	
	
	
	
	if max == 0 {
		return true
	}

	code := 0
	var nextcode [maxCodeLen]int
	for i := min; i <= max; i++ {
		code <<= 1
		nextcode[i&maxCodeLenMask] = code
		code += count[i&maxCodeLenMask]
	}

	
	
	
	
	
	if code != 1<<uint(max) && !(code == 1 && max == 1) {
		if debugDecode {
			fmt.Println("coding failed, code, max:", code, max, code == 1<<uint(max), code == 1 && max == 1, "(one should be true)")
		}
		return false
	}

	h.maxRead = min
	chunks := h.chunks[:]
	for i := range chunks {
		chunks[i] = 0
	}

	if max > huffmanChunkBits {
		numLinks := 1 << (uint(max) - huffmanChunkBits)
		h.linkMask = uint32(numLinks - 1)

		
		link := nextcode[huffmanChunkBits+1] >> 1
		if cap(h.links) < huffmanNumChunks-link {
			h.links = make([][]uint16, huffmanNumChunks-link)
		} else {
			h.links = h.links[:huffmanNumChunks-link]
		}
		for j := uint(link); j < huffmanNumChunks; j++ {
			reverse := int(bits.Reverse16(uint16(j)))
			reverse >>= uint(16 - huffmanChunkBits)
			off := j - uint(link)
			if sanity && h.chunks[reverse] != 0 {
				panic("impossible: overwriting existing chunk")
			}
			h.chunks[reverse] = uint16(off<<huffmanValueShift | (huffmanChunkBits + 1))
			if cap(h.links[off]) < numLinks {
				h.links[off] = make([]uint16, numLinks)
			} else {
				links := h.links[off][:0]
				h.links[off] = links[:numLinks]
			}
		}
	} else {
		h.links = h.links[:0]
	}

	for i, n := range lengths {
		if n == 0 {
			continue
		}
		code := nextcode[n]
		nextcode[n]++
		chunk := uint16(i<<huffmanValueShift | n)
		reverse := int(bits.Reverse16(uint16(code)))
		reverse >>= uint(16 - n)
		if n <= huffmanChunkBits {
			for off := reverse; off < len(h.chunks); off += 1 << uint(n) {
				
				
				
				
				
				if sanity && h.chunks[off] != 0 {
					panic("impossible: overwriting existing chunk")
				}
				h.chunks[off] = chunk
			}
		} else {
			j := reverse & (huffmanNumChunks - 1)
			if sanity && h.chunks[j]&huffmanCountMask != huffmanChunkBits+1 {
				
				
				panic("impossible: not an indirect chunk")
			}
			value := h.chunks[j] >> huffmanValueShift
			linktab := h.links[value]
			reverse >>= huffmanChunkBits
			for off := reverse; off < len(linktab); off += 1 << uint(n-huffmanChunkBits) {
				if sanity && linktab[off] != 0 {
					panic("impossible: overwriting existing chunk")
				}
				linktab[off] = chunk
			}
		}
	}

	if sanity {
		
		
		
		for i, chunk := range h.chunks {
			if chunk == 0 {
				
				
				
				if code == 1 && i%2 == 1 {
					continue
				}
				panic("impossible: missing chunk")
			}
		}
		for _, linktab := range h.links {
			for _, chunk := range linktab {
				if chunk == 0 {
					panic("impossible: missing chunk")
				}
			}
		}
	}

	return true
}




type Reader interface {
	io.Reader
	io.ByteReader
}


type decompressor struct {
	
	r       Reader
	roffset int64

	
	h1, h2 huffmanDecoder

	
	bits     *[maxNumLit + maxNumDist]int
	codebits *[numCodes]int

	
	dict dictDecoder

	
	
	step      func(*decompressor)
	stepState int
	err       error
	toRead    []byte
	hl, hd    *huffmanDecoder
	copyLen   int
	copyDist  int

	
	buf [4]byte

	
	b uint32

	nb    uint
	final bool
}

func (f *decompressor) nextBlock() {
	for f.nb < 1+2 {
		if f.err = f.moreBits(); f.err != nil {
			return
		}
	}
	f.final = f.b&1 == 1
	f.b >>= 1
	typ := f.b & 3
	f.b >>= 2
	f.nb -= 1 + 2
	switch typ {
	case 0:
		f.dataBlock()
		if debugDecode {
			fmt.Println("stored block")
		}
	case 1:
		
		f.hl = &fixedHuffmanDecoder
		f.hd = nil
		f.huffmanBlockDecoder()()
		if debugDecode {
			fmt.Println("predefinied huffman block")
		}
	case 2:
		
		if f.err = f.readHuffman(); f.err != nil {
			break
		}
		f.hl = &f.h1
		f.hd = &f.h2
		f.huffmanBlockDecoder()()
		if debugDecode {
			fmt.Println("dynamic huffman block")
		}
	default:
		
		if debugDecode {
			fmt.Println("reserved data block encountered")
		}
		f.err = CorruptInputError(f.roffset)
	}
}

func (f *decompressor) Read(b []byte) (int, error) {
	for {
		if len(f.toRead) > 0 {
			n := copy(b, f.toRead)
			f.toRead = f.toRead[n:]
			if len(f.toRead) == 0 {
				return n, f.err
			}
			return n, nil
		}
		if f.err != nil {
			return 0, f.err
		}
		f.step(f)
		if f.err != nil && len(f.toRead) == 0 {
			f.toRead = f.dict.readFlush() 
		}
	}
}


func (f *decompressor) WriteTo(w io.Writer) (int64, error) {
	total := int64(0)
	flushed := false
	for {
		if len(f.toRead) > 0 {
			n, err := w.Write(f.toRead)
			total += int64(n)
			if err != nil {
				f.err = err
				return total, err
			}
			if n != len(f.toRead) {
				return total, io.ErrShortWrite
			}
			f.toRead = f.toRead[:0]
		}
		if f.err != nil && flushed {
			if f.err == io.EOF {
				return total, nil
			}
			return total, f.err
		}
		if f.err == nil {
			f.step(f)
		}
		if len(f.toRead) == 0 && f.err != nil && !flushed {
			f.toRead = f.dict.readFlush() 
			flushed = true
		}
	}
}

func (f *decompressor) Close() error {
	if f.err == io.EOF {
		return nil
	}
	return f.err
}




var codeOrder = [...]int{16, 17, 18, 0, 8, 7, 9, 6, 10, 5, 11, 4, 12, 3, 13, 2, 14, 1, 15}

func (f *decompressor) readHuffman() error {
	
	for f.nb < 5+5+4 {
		if err := f.moreBits(); err != nil {
			return err
		}
	}
	nlit := int(f.b&0x1F) + 257
	if nlit > maxNumLit {
		if debugDecode {
			fmt.Println("nlit > maxNumLit", nlit)
		}
		return CorruptInputError(f.roffset)
	}
	f.b >>= 5
	ndist := int(f.b&0x1F) + 1
	if ndist > maxNumDist {
		if debugDecode {
			fmt.Println("ndist > maxNumDist", ndist)
		}
		return CorruptInputError(f.roffset)
	}
	f.b >>= 5
	nclen := int(f.b&0xF) + 4
	
	f.b >>= 4
	f.nb -= 5 + 5 + 4

	
	for i := 0; i < nclen; i++ {
		for f.nb < 3 {
			if err := f.moreBits(); err != nil {
				return err
			}
		}
		f.codebits[codeOrder[i]] = int(f.b & 0x7)
		f.b >>= 3
		f.nb -= 3
	}
	for i := nclen; i < len(codeOrder); i++ {
		f.codebits[codeOrder[i]] = 0
	}
	if !f.h1.init(f.codebits[0:]) {
		if debugDecode {
			fmt.Println("init codebits failed")
		}
		return CorruptInputError(f.roffset)
	}

	
	
	for i, n := 0, nlit+ndist; i < n; {
		x, err := f.huffSym(&f.h1)
		if err != nil {
			return err
		}
		if x < 16 {
			
			f.bits[i] = x
			i++
			continue
		}
		
		var rep int
		var nb uint
		var b int
		switch x {
		default:
			return InternalError("unexpected length code")
		case 16:
			rep = 3
			nb = 2
			if i == 0 {
				if debugDecode {
					fmt.Println("i==0")
				}
				return CorruptInputError(f.roffset)
			}
			b = f.bits[i-1]
		case 17:
			rep = 3
			nb = 3
			b = 0
		case 18:
			rep = 11
			nb = 7
			b = 0
		}
		for f.nb < nb {
			if err := f.moreBits(); err != nil {
				if debugDecode {
					fmt.Println("morebits:", err)
				}
				return err
			}
		}
		rep += int(f.b & uint32(1<<(nb&regSizeMaskUint32)-1))
		f.b >>= nb & regSizeMaskUint32
		f.nb -= nb
		if i+rep > n {
			if debugDecode {
				fmt.Println("i+rep > n", i, rep, n)
			}
			return CorruptInputError(f.roffset)
		}
		for j := 0; j < rep; j++ {
			f.bits[i] = b
			i++
		}
	}

	if !f.h1.init(f.bits[0:nlit]) || !f.h2.init(f.bits[nlit:nlit+ndist]) {
		if debugDecode {
			fmt.Println("init2 failed")
		}
		return CorruptInputError(f.roffset)
	}

	
	
	
	
	if f.h1.maxRead < f.bits[endBlockMarker] {
		f.h1.maxRead = f.bits[endBlockMarker]
	}
	if !f.final {
		
		
		
		f.h1.maxRead += 10
	}

	return nil
}


func (f *decompressor) dataBlock() {
	
	
	left := (f.nb) & 7
	f.nb -= left
	f.b >>= left

	offBytes := f.nb >> 3
	
	f.buf[0] = uint8(f.b)
	f.buf[1] = uint8(f.b >> 8)
	f.buf[2] = uint8(f.b >> 16)
	f.buf[3] = uint8(f.b >> 24)

	f.roffset += int64(offBytes)
	f.nb, f.b = 0, 0

	
	nr, err := io.ReadFull(f.r, f.buf[offBytes:4])
	f.roffset += int64(nr)
	if err != nil {
		f.err = noEOF(err)
		return
	}
	n := uint16(f.buf[0]) | uint16(f.buf[1])<<8
	nn := uint16(f.buf[2]) | uint16(f.buf[3])<<8
	if nn != ^n {
		if debugDecode {
			ncomp := ^n
			fmt.Println("uint16(nn) != uint16(^n)", nn, ncomp)
		}
		f.err = CorruptInputError(f.roffset)
		return
	}

	if n == 0 {
		f.toRead = f.dict.readFlush()
		f.finishBlock()
		return
	}

	f.copyLen = int(n)
	f.copyData()
}



func (f *decompressor) copyData() {
	buf := f.dict.writeSlice()
	if len(buf) > f.copyLen {
		buf = buf[:f.copyLen]
	}

	cnt, err := io.ReadFull(f.r, buf)
	f.roffset += int64(cnt)
	f.copyLen -= cnt
	f.dict.writeMark(cnt)
	if err != nil {
		f.err = noEOF(err)
		return
	}

	if f.dict.availWrite() == 0 || f.copyLen > 0 {
		f.toRead = f.dict.readFlush()
		f.step = (*decompressor).copyData
		return
	}
	f.finishBlock()
}

func (f *decompressor) finishBlock() {
	if f.final {
		if f.dict.availRead() > 0 {
			f.toRead = f.dict.readFlush()
		}
		f.err = io.EOF
	}
	f.step = (*decompressor).nextBlock
}


func noEOF(e error) error {
	if e == io.EOF {
		return io.ErrUnexpectedEOF
	}
	return e
}

func (f *decompressor) moreBits() error {
	c, err := f.r.ReadByte()
	if err != nil {
		return noEOF(err)
	}
	f.roffset++
	f.b |= uint32(c) << (f.nb & regSizeMaskUint32)
	f.nb += 8
	return nil
}


func (f *decompressor) huffSym(h *huffmanDecoder) (int, error) {
	
	
	
	
	n := uint(h.maxRead)
	
	
	
	nb, b := f.nb, f.b
	for {
		for nb < n {
			c, err := f.r.ReadByte()
			if err != nil {
				f.b = b
				f.nb = nb
				return 0, noEOF(err)
			}
			f.roffset++
			b |= uint32(c) << (nb & regSizeMaskUint32)
			nb += 8
		}
		chunk := h.chunks[b&(huffmanNumChunks-1)]
		n = uint(chunk & huffmanCountMask)
		if n > huffmanChunkBits {
			chunk = h.links[chunk>>huffmanValueShift][(b>>huffmanChunkBits)&h.linkMask]
			n = uint(chunk & huffmanCountMask)
		}
		if n <= nb {
			if n == 0 {
				f.b = b
				f.nb = nb
				if debugDecode {
					fmt.Println("huffsym: n==0")
				}
				f.err = CorruptInputError(f.roffset)
				return 0, f.err
			}
			f.b = b >> (n & regSizeMaskUint32)
			f.nb = nb - n
			return int(chunk >> huffmanValueShift), nil
		}
	}
}

func makeReader(r io.Reader) Reader {
	if rr, ok := r.(Reader); ok {
		return rr
	}
	return bufio.NewReader(r)
}

func fixedHuffmanDecoderInit() {
	fixedOnce.Do(func() {
		
		var bits [288]int
		for i := 0; i < 144; i++ {
			bits[i] = 8
		}
		for i := 144; i < 256; i++ {
			bits[i] = 9
		}
		for i := 256; i < 280; i++ {
			bits[i] = 7
		}
		for i := 280; i < 288; i++ {
			bits[i] = 8
		}
		fixedHuffmanDecoder.init(bits[:])
	})
}

func (f *decompressor) Reset(r io.Reader, dict []byte) error {
	*f = decompressor{
		r:        makeReader(r),
		bits:     f.bits,
		codebits: f.codebits,
		h1:       f.h1,
		h2:       f.h2,
		dict:     f.dict,
		step:     (*decompressor).nextBlock,
	}
	f.dict.init(maxMatchOffset, dict)
	return nil
}









func NewReader(r io.Reader) io.ReadCloser {
	fixedHuffmanDecoderInit()

	var f decompressor
	f.r = makeReader(r)
	f.bits = new([maxNumLit + maxNumDist]int)
	f.codebits = new([numCodes]int)
	f.step = (*decompressor).nextBlock
	f.dict.init(maxMatchOffset, nil)
	return &f
}








func NewReaderDict(r io.Reader, dict []byte) io.ReadCloser {
	fixedHuffmanDecoderInit()

	var f decompressor
	f.r = makeReader(r)
	f.bits = new([maxNumLit + maxNumDist]int)
	f.codebits = new([numCodes]int)
	f.step = (*decompressor).nextBlock
	f.dict.init(maxMatchOffset, dict)
	return &f
}
