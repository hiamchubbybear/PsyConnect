



package zstd

import (
	"crypto/rand"
	"fmt"
	"io"
	"math"
	rdebug "runtime/debug"
	"sync"

	"github.com/klauspost/compress/zstd/internal/xxhash"
)







type Encoder struct {
	o        encoderOptions
	encoders chan encoder
	state    encoderState
	init     sync.Once
}

type encoder interface {
	Encode(blk *blockEnc, src []byte)
	EncodeNoHist(blk *blockEnc, src []byte)
	Block() *blockEnc
	CRC() *xxhash.Digest
	AppendCRC([]byte) []byte
	WindowSize(size int64) int32
	UseBlock(*blockEnc)
	Reset(d *dict, singleBlock bool)
}

type encoderState struct {
	w                io.Writer
	filling          []byte
	current          []byte
	previous         []byte
	encoder          encoder
	writing          *blockEnc
	err              error
	writeErr         error
	nWritten         int64
	nInput           int64
	frameContentSize int64
	headerWritten    bool
	eofWritten       bool
	fullFrameWritten bool

	
	wg sync.WaitGroup
	
	wWg sync.WaitGroup
}



func NewWriter(w io.Writer, opts ...EOption) (*Encoder, error) {
	initPredefined()
	var e Encoder
	e.o.setDefault()
	for _, o := range opts {
		err := o(&e.o)
		if err != nil {
			return nil, err
		}
	}
	if w != nil {
		e.Reset(w)
	}
	return &e, nil
}

func (e *Encoder) initialize() {
	if e.o.concurrent == 0 {
		e.o.setDefault()
	}
	e.encoders = make(chan encoder, e.o.concurrent)
	for i := 0; i < e.o.concurrent; i++ {
		enc := e.o.encoder()
		e.encoders <- enc
	}
}



func (e *Encoder) Reset(w io.Writer) {
	s := &e.state
	s.wg.Wait()
	s.wWg.Wait()
	if cap(s.filling) == 0 {
		s.filling = make([]byte, 0, e.o.blockSize)
	}
	if e.o.concurrent > 1 {
		if cap(s.current) == 0 {
			s.current = make([]byte, 0, e.o.blockSize)
		}
		if cap(s.previous) == 0 {
			s.previous = make([]byte, 0, e.o.blockSize)
		}
		s.current = s.current[:0]
		s.previous = s.previous[:0]
		if s.writing == nil {
			s.writing = &blockEnc{lowMem: e.o.lowMem}
			s.writing.init()
		}
		s.writing.initNewEncode()
	}
	if s.encoder == nil {
		s.encoder = e.o.encoder()
	}
	s.filling = s.filling[:0]
	s.encoder.Reset(e.o.dict, false)
	s.headerWritten = false
	s.eofWritten = false
	s.fullFrameWritten = false
	s.w = w
	s.err = nil
	s.nWritten = 0
	s.nInput = 0
	s.writeErr = nil
	s.frameContentSize = 0
}






func (e *Encoder) ResetContentSize(w io.Writer, size int64) {
	e.Reset(w)
	if size >= 0 {
		e.state.frameContentSize = size
	}
}






func (e *Encoder) Write(p []byte) (n int, err error) {
	s := &e.state
	for len(p) > 0 {
		if len(p)+len(s.filling) < e.o.blockSize {
			if e.o.crc {
				_, _ = s.encoder.CRC().Write(p)
			}
			s.filling = append(s.filling, p...)
			return n + len(p), nil
		}
		add := p
		if len(p)+len(s.filling) > e.o.blockSize {
			add = add[:e.o.blockSize-len(s.filling)]
		}
		if e.o.crc {
			_, _ = s.encoder.CRC().Write(add)
		}
		s.filling = append(s.filling, add...)
		p = p[len(add):]
		n += len(add)
		if len(s.filling) < e.o.blockSize {
			return n, nil
		}
		err := e.nextBlock(false)
		if err != nil {
			return n, err
		}
		if debugAsserts && len(s.filling) > 0 {
			panic(len(s.filling))
		}
	}
	return n, nil
}



func (e *Encoder) nextBlock(final bool) error {
	s := &e.state
	
	s.wg.Wait()
	if s.err != nil {
		return s.err
	}
	if len(s.filling) > e.o.blockSize {
		return fmt.Errorf("block > maxStoreBlockSize")
	}
	if !s.headerWritten {
		
		if final && len(s.filling) == 0 && !e.o.fullZero {
			s.headerWritten = true
			s.fullFrameWritten = true
			s.eofWritten = true
			return nil
		}
		if final && len(s.filling) > 0 {
			s.current = e.EncodeAll(s.filling, s.current[:0])
			var n2 int
			n2, s.err = s.w.Write(s.current)
			if s.err != nil {
				return s.err
			}
			s.nWritten += int64(n2)
			s.nInput += int64(len(s.filling))
			s.current = s.current[:0]
			s.filling = s.filling[:0]
			s.headerWritten = true
			s.fullFrameWritten = true
			s.eofWritten = true
			return nil
		}

		var tmp [maxHeaderSize]byte
		fh := frameHeader{
			ContentSize:   uint64(s.frameContentSize),
			WindowSize:    uint32(s.encoder.WindowSize(s.frameContentSize)),
			SingleSegment: false,
			Checksum:      e.o.crc,
			DictID:        e.o.dict.ID(),
		}

		dst, err := fh.appendTo(tmp[:0])
		if err != nil {
			return err
		}
		s.headerWritten = true
		s.wWg.Wait()
		var n2 int
		n2, s.err = s.w.Write(dst)
		if s.err != nil {
			return s.err
		}
		s.nWritten += int64(n2)
	}
	if s.eofWritten {
		
		final = false
	}

	if len(s.filling) == 0 {
		
		if final {
			enc := s.encoder
			blk := enc.Block()
			blk.reset(nil)
			blk.last = true
			blk.encodeRaw(nil)
			s.wWg.Wait()
			_, s.err = s.w.Write(blk.output)
			s.nWritten += int64(len(blk.output))
			s.eofWritten = true
		}
		return s.err
	}

	
	if e.o.concurrent == 1 {
		src := s.filling
		s.nInput += int64(len(s.filling))
		if debugEncoder {
			println("Adding sync block,", len(src), "bytes, final:", final)
		}
		enc := s.encoder
		blk := enc.Block()
		blk.reset(nil)
		enc.Encode(blk, src)
		blk.last = final
		if final {
			s.eofWritten = true
		}

		s.err = blk.encode(src, e.o.noEntropy, !e.o.allLitEntropy)
		if s.err != nil {
			return s.err
		}
		_, s.err = s.w.Write(blk.output)
		s.nWritten += int64(len(blk.output))
		s.filling = s.filling[:0]
		return s.err
	}

	
	s.filling, s.current, s.previous = s.previous[:0], s.filling, s.current
	s.nInput += int64(len(s.current))
	s.wg.Add(1)
	go func(src []byte) {
		if debugEncoder {
			println("Adding block,", len(src), "bytes, final:", final)
		}
		defer func() {
			if r := recover(); r != nil {
				s.err = fmt.Errorf("panic while encoding: %v", r)
				rdebug.PrintStack()
			}
			s.wg.Done()
		}()
		enc := s.encoder
		blk := enc.Block()
		enc.Encode(blk, src)
		blk.last = final
		if final {
			s.eofWritten = true
		}
		
		s.wWg.Wait()
		if s.writeErr != nil {
			s.err = s.writeErr
			return
		}
		
		blk.swapEncoders(s.writing)
		
		enc.UseBlock(s.writing)
		s.writing = blk
		s.wWg.Add(1)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					s.writeErr = fmt.Errorf("panic while encoding/writing: %v", r)
					rdebug.PrintStack()
				}
				s.wWg.Done()
			}()
			s.writeErr = blk.encode(src, e.o.noEntropy, !e.o.allLitEntropy)
			if s.writeErr != nil {
				return
			}
			_, s.writeErr = s.w.Write(blk.output)
			s.nWritten += int64(len(blk.output))
		}()
	}(s.current)
	return nil
}






func (e *Encoder) ReadFrom(r io.Reader) (n int64, err error) {
	if debugEncoder {
		println("Using ReadFrom")
	}

	
	if len(e.state.filling) > 0 {
		if err := e.nextBlock(false); err != nil {
			return 0, err
		}
	}
	e.state.filling = e.state.filling[:e.o.blockSize]
	src := e.state.filling
	for {
		n2, err := r.Read(src)
		if e.o.crc {
			_, _ = e.state.encoder.CRC().Write(src[:n2])
		}
		
		src = src[n2:]
		n += int64(n2)
		switch err {
		case io.EOF:
			e.state.filling = e.state.filling[:len(e.state.filling)-len(src)]
			if debugEncoder {
				println("ReadFrom: got EOF final block:", len(e.state.filling))
			}
			return n, nil
		case nil:
		default:
			if debugEncoder {
				println("ReadFrom: got error:", err)
			}
			e.state.err = err
			return n, err
		}
		if len(src) > 0 {
			if debugEncoder {
				println("ReadFrom: got space left in source:", len(src))
			}
			continue
		}
		err = e.nextBlock(false)
		if err != nil {
			return n, err
		}
		e.state.filling = e.state.filling[:e.o.blockSize]
		src = e.state.filling
	}
}




func (e *Encoder) Flush() error {
	s := &e.state
	if len(s.filling) > 0 {
		err := e.nextBlock(false)
		if err != nil {
			return err
		}
	}
	s.wg.Wait()
	s.wWg.Wait()
	if s.err != nil {
		return s.err
	}
	return s.writeErr
}




func (e *Encoder) Close() error {
	s := &e.state
	if s.encoder == nil {
		return nil
	}
	err := e.nextBlock(true)
	if err != nil {
		return err
	}
	if s.frameContentSize > 0 {
		if s.nInput != s.frameContentSize {
			return fmt.Errorf("frame content size %d given, but %d bytes was written", s.frameContentSize, s.nInput)
		}
	}
	if e.state.fullFrameWritten {
		return s.err
	}
	s.wg.Wait()
	s.wWg.Wait()

	if s.err != nil {
		return s.err
	}
	if s.writeErr != nil {
		return s.writeErr
	}

	
	if e.o.crc && s.err == nil {
		
		var tmp [4]byte
		_, s.err = s.w.Write(s.encoder.AppendCRC(tmp[:0]))
		s.nWritten += 4
	}

	
	if s.err == nil && e.o.pad > 0 {
		add := calcSkippableFrame(s.nWritten, int64(e.o.pad))
		frame, err := skippableFrame(s.filling[:0], add, rand.Reader)
		if err != nil {
			return err
		}
		_, s.err = s.w.Write(frame)
	}
	return s.err
}







func (e *Encoder) EncodeAll(src, dst []byte) []byte {
	if len(src) == 0 {
		if e.o.fullZero {
			
			fh := frameHeader{
				ContentSize:   0,
				WindowSize:    MinWindowSize,
				SingleSegment: true,
				
				Checksum: false,
				DictID:   0,
			}
			dst, _ = fh.appendTo(dst)

			
			var blk blockHeader
			blk.setSize(0)
			blk.setType(blockTypeRaw)
			blk.setLast(true)
			dst = blk.appendTo(dst)
		}
		return dst
	}
	e.init.Do(e.initialize)
	enc := <-e.encoders
	defer func() {
		
		
		e.encoders <- enc
	}()
	
	single := len(src) <= e.o.windowSize && len(src) > MinWindowSize
	if e.o.single != nil {
		single = *e.o.single
	}
	fh := frameHeader{
		ContentSize:   uint64(len(src)),
		WindowSize:    uint32(enc.WindowSize(int64(len(src)))),
		SingleSegment: single,
		Checksum:      e.o.crc,
		DictID:        e.o.dict.ID(),
	}

	
	if len(dst) == 0 && cap(dst) == 0 && len(src) < 1<<20 && !e.o.lowMem {
		dst = make([]byte, 0, len(src))
	}
	dst, err := fh.appendTo(dst)
	if err != nil {
		panic(err)
	}

	
	if len(src) <= e.o.blockSize {
		enc.Reset(e.o.dict, true)
		
		if e.o.crc {
			_, _ = enc.CRC().Write(src)
		}
		blk := enc.Block()
		blk.last = true
		if e.o.dict == nil {
			enc.EncodeNoHist(blk, src)
		} else {
			enc.Encode(blk, src)
		}

		
		
		oldout := blk.output
		
		blk.output = dst

		err := blk.encode(src, e.o.noEntropy, !e.o.allLitEntropy)
		if err != nil {
			panic(err)
		}
		dst = blk.output
		blk.output = oldout
	} else {
		enc.Reset(e.o.dict, false)
		blk := enc.Block()
		for len(src) > 0 {
			todo := src
			if len(todo) > e.o.blockSize {
				todo = todo[:e.o.blockSize]
			}
			src = src[len(todo):]
			if e.o.crc {
				_, _ = enc.CRC().Write(todo)
			}
			blk.pushOffsets()
			enc.Encode(blk, todo)
			if len(src) == 0 {
				blk.last = true
			}
			err := blk.encode(todo, e.o.noEntropy, !e.o.allLitEntropy)
			if err != nil {
				panic(err)
			}
			dst = append(dst, blk.output...)
			blk.reset(nil)
		}
	}
	if e.o.crc {
		dst = enc.AppendCRC(dst)
	}
	
	if e.o.pad > 0 {
		add := calcSkippableFrame(int64(len(dst)), int64(e.o.pad))
		dst, err = skippableFrame(dst, add, rand.Reader)
		if err != nil {
			panic(err)
		}
	}
	return dst
}



func (e *Encoder) MaxEncodedSize(size int) int {
	frameHeader := 4 + 2 
	if e.o.dict != nil {
		frameHeader += 4
	}
	
	if size < 256 {
		frameHeader++
	} else if size < 65536+256 {
		frameHeader += 2
	} else if size < math.MaxInt32 {
		frameHeader += 4
	} else {
		frameHeader += 8
	}
	
	if e.o.crc {
		frameHeader += 4
	}

	
	
	blocks := (size + e.o.blockSize) / e.o.blockSize

	
	maxSz := frameHeader + 3*blocks + size
	if e.o.pad > 1 {
		maxSz += calcSkippableFrame(int64(maxSz), int64(e.o.pad))
	}
	return maxSz
}
