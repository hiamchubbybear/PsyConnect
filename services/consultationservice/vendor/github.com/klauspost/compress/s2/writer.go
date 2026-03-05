




package s2

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"runtime"
	"sync"
)

const (
	levelUncompressed = iota + 1
	levelFast
	levelBetter
	levelBest
)








func NewWriter(w io.Writer, opts ...WriterOption) *Writer {
	w2 := Writer{
		blockSize:   defaultBlockSize,
		concurrency: runtime.GOMAXPROCS(0),
		randSrc:     rand.Reader,
		level:       levelFast,
	}
	for _, opt := range opts {
		if err := opt(&w2); err != nil {
			w2.errState = err
			return &w2
		}
	}
	w2.obufLen = obufHeaderLen + MaxEncodedLen(w2.blockSize)
	w2.paramsOK = true
	w2.ibuf = make([]byte, 0, w2.blockSize)
	w2.buffers.New = func() interface{} {
		return make([]byte, w2.obufLen)
	}
	w2.Reset(w)
	return &w2
}


type Writer struct {
	errMu    sync.Mutex
	errState error

	
	ibuf []byte

	blockSize     int
	obufLen       int
	concurrency   int
	written       int64
	uncompWritten int64 
	output        chan chan result
	buffers       sync.Pool
	pad           int

	writer    io.Writer
	randSrc   io.Reader
	writerWg  sync.WaitGroup
	index     Index
	customEnc func(dst, src []byte) int

	
	wroteStreamHeader bool
	paramsOK          bool
	snappy            bool
	flushOnWrite      bool
	appendIndex       bool
	level             uint8
}

type result struct {
	b []byte
	
	startOffset int64
}



func (w *Writer) err(err error) error {
	w.errMu.Lock()
	errSet := w.errState
	if errSet == nil && err != nil {
		w.errState = err
		errSet = err
	}
	w.errMu.Unlock()
	return errSet
}



func (w *Writer) Reset(writer io.Writer) {
	if !w.paramsOK {
		return
	}
	
	if w.output != nil {
		close(w.output)
		w.writerWg.Wait()
		w.output = nil
	}
	w.errState = nil
	w.ibuf = w.ibuf[:0]
	w.wroteStreamHeader = false
	w.written = 0
	w.writer = writer
	w.uncompWritten = 0
	w.index.reset(w.blockSize)

	
	if writer == nil {
		return
	}
	
	if w.concurrency == 1 {
		return
	}

	toWrite := make(chan chan result, w.concurrency)
	w.output = toWrite
	w.writerWg.Add(1)

	
	go func() {
		defer w.writerWg.Done()

		
		for write := range toWrite {
			
			input := <-write
			in := input.b
			if len(in) > 0 {
				if w.err(nil) == nil {
					
					toWrite := in[:len(in):len(in)]
					
					n, err := writer.Write(toWrite)
					if err == nil && n != len(toWrite) {
						err = io.ErrShortBuffer
					}
					_ = w.err(err)
					w.err(w.index.add(w.written, input.startOffset))
					w.written += int64(n)
				}
			}
			if cap(in) >= w.obufLen {
				w.buffers.Put(in)
			}
			
			
			close(write)
		}
	}()
}


func (w *Writer) Write(p []byte) (nRet int, errRet error) {
	if err := w.err(nil); err != nil {
		return 0, err
	}
	if w.flushOnWrite {
		return w.write(p)
	}
	
	for len(p) > (cap(w.ibuf)-len(w.ibuf)) && w.err(nil) == nil {
		var n int
		if len(w.ibuf) == 0 {
			
			
			n, _ = w.write(p)
		} else {
			n = copy(w.ibuf[len(w.ibuf):cap(w.ibuf)], p)
			w.ibuf = w.ibuf[:len(w.ibuf)+n]
			w.write(w.ibuf)
			w.ibuf = w.ibuf[:0]
		}
		nRet += n
		p = p[n:]
	}
	if err := w.err(nil); err != nil {
		return nRet, err
	}
	
	n := copy(w.ibuf[len(w.ibuf):cap(w.ibuf)], p)
	w.ibuf = w.ibuf[:len(w.ibuf)+n]
	nRet += n
	return nRet, nil
}






func (w *Writer) ReadFrom(r io.Reader) (n int64, err error) {
	if err := w.err(nil); err != nil {
		return 0, err
	}
	if len(w.ibuf) > 0 {
		err := w.Flush()
		if err != nil {
			return 0, err
		}
	}
	if br, ok := r.(byter); ok {
		buf := br.Bytes()
		if err := w.EncodeBuffer(buf); err != nil {
			return 0, err
		}
		return int64(len(buf)), w.Flush()
	}
	for {
		inbuf := w.buffers.Get().([]byte)[:w.blockSize+obufHeaderLen]
		n2, err := io.ReadFull(r, inbuf[obufHeaderLen:])
		if err != nil {
			if err == io.ErrUnexpectedEOF {
				err = io.EOF
			}
			if err != io.EOF {
				return n, w.err(err)
			}
		}
		if n2 == 0 {
			break
		}
		n += int64(n2)
		err2 := w.writeFull(inbuf[:n2+obufHeaderLen])
		if w.err(err2) != nil {
			break
		}

		if err != nil {
			
			break
		}
	}

	return n, w.err(nil)
}




func (w *Writer) AddSkippableBlock(id uint8, data []byte) (err error) {
	if err := w.err(nil); err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}
	if id < 0x80 || id > chunkTypePadding {
		return fmt.Errorf("invalid skippable block id %x", id)
	}
	if len(data) > maxChunkSize {
		return fmt.Errorf("skippable block excessed maximum size")
	}
	var header [4]byte
	chunkLen := 4 + len(data)
	header[0] = id
	header[1] = uint8(chunkLen >> 0)
	header[2] = uint8(chunkLen >> 8)
	header[3] = uint8(chunkLen >> 16)
	if w.concurrency == 1 {
		write := func(b []byte) error {
			n, err := w.writer.Write(b)
			if err = w.err(err); err != nil {
				return err
			}
			if n != len(data) {
				return w.err(io.ErrShortWrite)
			}
			w.written += int64(n)
			return w.err(nil)
		}
		if !w.wroteStreamHeader {
			w.wroteStreamHeader = true
			if w.snappy {
				if err := write([]byte(magicChunkSnappy)); err != nil {
					return err
				}
			} else {
				if err := write([]byte(magicChunk)); err != nil {
					return err
				}
			}
		}
		if err := write(header[:]); err != nil {
			return err
		}
		if err := write(data); err != nil {
			return err
		}
	}

	
	if !w.wroteStreamHeader {
		w.wroteStreamHeader = true
		hWriter := make(chan result)
		w.output <- hWriter
		if w.snappy {
			hWriter <- result{startOffset: w.uncompWritten, b: []byte(magicChunkSnappy)}
		} else {
			hWriter <- result{startOffset: w.uncompWritten, b: []byte(magicChunk)}
		}
	}

	
	inbuf := w.buffers.Get().([]byte)[:4]
	copy(inbuf, header[:])
	inbuf = append(inbuf, data...)

	output := make(chan result, 1)
	
	w.output <- output
	output <- result{startOffset: w.uncompWritten, b: inbuf}

	return nil
}











func (w *Writer) EncodeBuffer(buf []byte) (err error) {
	if err := w.err(nil); err != nil {
		return err
	}

	if w.flushOnWrite {
		_, err := w.write(buf)
		return err
	}
	
	if len(w.ibuf) > 0 {
		err := w.Flush()
		if err != nil {
			return err
		}
	}
	if w.concurrency == 1 {
		_, err := w.writeSync(buf)
		return err
	}

	
	if !w.wroteStreamHeader {
		w.wroteStreamHeader = true
		hWriter := make(chan result)
		w.output <- hWriter
		if w.snappy {
			hWriter <- result{startOffset: w.uncompWritten, b: []byte(magicChunkSnappy)}
		} else {
			hWriter <- result{startOffset: w.uncompWritten, b: []byte(magicChunk)}
		}
	}

	for len(buf) > 0 {
		
		uncompressed := buf
		if len(uncompressed) > w.blockSize {
			uncompressed = uncompressed[:w.blockSize]
		}
		buf = buf[len(uncompressed):]
		
		obuf := w.buffers.Get().([]byte)[:len(uncompressed)+obufHeaderLen]
		output := make(chan result)
		
		w.output <- output
		res := result{
			startOffset: w.uncompWritten,
		}
		w.uncompWritten += int64(len(uncompressed))
		go func() {
			checksum := crc(uncompressed)

			
			chunkType := uint8(chunkTypeUncompressedData)
			chunkLen := 4 + len(uncompressed)

			
			n := binary.PutUvarint(obuf[obufHeaderLen:], uint64(len(uncompressed)))
			n2 := w.encodeBlock(obuf[obufHeaderLen+n:], uncompressed)

			
			if n2 > 0 {
				chunkType = uint8(chunkTypeCompressedData)
				chunkLen = 4 + n + n2
				obuf = obuf[:obufHeaderLen+n+n2]
			} else {
				
				copy(obuf[obufHeaderLen:], uncompressed)
			}

			
			obuf[0] = chunkType
			obuf[1] = uint8(chunkLen >> 0)
			obuf[2] = uint8(chunkLen >> 8)
			obuf[3] = uint8(chunkLen >> 16)
			obuf[4] = uint8(checksum >> 0)
			obuf[5] = uint8(checksum >> 8)
			obuf[6] = uint8(checksum >> 16)
			obuf[7] = uint8(checksum >> 24)

			
			res.b = obuf
			output <- res
		}()
	}
	return nil
}

func (w *Writer) encodeBlock(obuf, uncompressed []byte) int {
	if w.customEnc != nil {
		if ret := w.customEnc(obuf, uncompressed); ret >= 0 {
			return ret
		}
	}
	if w.snappy {
		switch w.level {
		case levelFast:
			return encodeBlockSnappy(obuf, uncompressed)
		case levelBetter:
			return encodeBlockBetterSnappy(obuf, uncompressed)
		case levelBest:
			return encodeBlockBestSnappy(obuf, uncompressed)
		}
		return 0
	}
	switch w.level {
	case levelFast:
		return encodeBlock(obuf, uncompressed)
	case levelBetter:
		return encodeBlockBetter(obuf, uncompressed)
	case levelBest:
		return encodeBlockBest(obuf, uncompressed, nil)
	}
	return 0
}

func (w *Writer) write(p []byte) (nRet int, errRet error) {
	if err := w.err(nil); err != nil {
		return 0, err
	}
	if w.concurrency == 1 {
		return w.writeSync(p)
	}

	
	for len(p) > 0 {
		if !w.wroteStreamHeader {
			w.wroteStreamHeader = true
			hWriter := make(chan result)
			w.output <- hWriter
			if w.snappy {
				hWriter <- result{startOffset: w.uncompWritten, b: []byte(magicChunkSnappy)}
			} else {
				hWriter <- result{startOffset: w.uncompWritten, b: []byte(magicChunk)}
			}
		}

		var uncompressed []byte
		if len(p) > w.blockSize {
			uncompressed, p = p[:w.blockSize], p[w.blockSize:]
		} else {
			uncompressed, p = p, nil
		}

		
		
		inbuf := w.buffers.Get().([]byte)[:len(uncompressed)+obufHeaderLen]
		obuf := w.buffers.Get().([]byte)[:w.obufLen]
		copy(inbuf[obufHeaderLen:], uncompressed)
		uncompressed = inbuf[obufHeaderLen:]

		output := make(chan result)
		
		w.output <- output
		res := result{
			startOffset: w.uncompWritten,
		}
		w.uncompWritten += int64(len(uncompressed))

		go func() {
			checksum := crc(uncompressed)

			
			chunkType := uint8(chunkTypeUncompressedData)
			chunkLen := 4 + len(uncompressed)

			
			n := binary.PutUvarint(obuf[obufHeaderLen:], uint64(len(uncompressed)))
			n2 := w.encodeBlock(obuf[obufHeaderLen+n:], uncompressed)

			
			if n2 > 0 {
				chunkType = uint8(chunkTypeCompressedData)
				chunkLen = 4 + n + n2
				obuf = obuf[:obufHeaderLen+n+n2]
			} else {
				
				obuf, inbuf = inbuf, obuf
			}

			
			obuf[0] = chunkType
			obuf[1] = uint8(chunkLen >> 0)
			obuf[2] = uint8(chunkLen >> 8)
			obuf[3] = uint8(chunkLen >> 16)
			obuf[4] = uint8(checksum >> 0)
			obuf[5] = uint8(checksum >> 8)
			obuf[6] = uint8(checksum >> 16)
			obuf[7] = uint8(checksum >> 24)

			
			res.b = obuf
			output <- res

			
			w.buffers.Put(inbuf)
		}()
		nRet += len(uncompressed)
	}
	return nRet, nil
}





func (w *Writer) writeFull(inbuf []byte) (errRet error) {
	if err := w.err(nil); err != nil {
		return err
	}

	if w.concurrency == 1 {
		_, err := w.writeSync(inbuf[obufHeaderLen:])
		return err
	}

	
	if !w.wroteStreamHeader {
		w.wroteStreamHeader = true
		hWriter := make(chan result)
		w.output <- hWriter
		if w.snappy {
			hWriter <- result{startOffset: w.uncompWritten, b: []byte(magicChunkSnappy)}
		} else {
			hWriter <- result{startOffset: w.uncompWritten, b: []byte(magicChunk)}
		}
	}

	
	obuf := w.buffers.Get().([]byte)[:w.obufLen]
	uncompressed := inbuf[obufHeaderLen:]

	output := make(chan result)
	
	w.output <- output
	res := result{
		startOffset: w.uncompWritten,
	}
	w.uncompWritten += int64(len(uncompressed))

	go func() {
		checksum := crc(uncompressed)

		
		chunkType := uint8(chunkTypeUncompressedData)
		chunkLen := 4 + len(uncompressed)

		
		n := binary.PutUvarint(obuf[obufHeaderLen:], uint64(len(uncompressed)))
		n2 := w.encodeBlock(obuf[obufHeaderLen+n:], uncompressed)

		
		if n2 > 0 {
			chunkType = uint8(chunkTypeCompressedData)
			chunkLen = 4 + n + n2
			obuf = obuf[:obufHeaderLen+n+n2]
		} else {
			
			obuf, inbuf = inbuf, obuf
		}

		
		obuf[0] = chunkType
		obuf[1] = uint8(chunkLen >> 0)
		obuf[2] = uint8(chunkLen >> 8)
		obuf[3] = uint8(chunkLen >> 16)
		obuf[4] = uint8(checksum >> 0)
		obuf[5] = uint8(checksum >> 8)
		obuf[6] = uint8(checksum >> 16)
		obuf[7] = uint8(checksum >> 24)

		
		res.b = obuf
		output <- res

		
		w.buffers.Put(inbuf)
	}()
	return nil
}

func (w *Writer) writeSync(p []byte) (nRet int, errRet error) {
	if err := w.err(nil); err != nil {
		return 0, err
	}
	if !w.wroteStreamHeader {
		w.wroteStreamHeader = true
		var n int
		var err error
		if w.snappy {
			n, err = w.writer.Write([]byte(magicChunkSnappy))
		} else {
			n, err = w.writer.Write([]byte(magicChunk))
		}
		if err != nil {
			return 0, w.err(err)
		}
		if n != len(magicChunk) {
			return 0, w.err(io.ErrShortWrite)
		}
		w.written += int64(n)
	}

	for len(p) > 0 {
		var uncompressed []byte
		if len(p) > w.blockSize {
			uncompressed, p = p[:w.blockSize], p[w.blockSize:]
		} else {
			uncompressed, p = p, nil
		}

		obuf := w.buffers.Get().([]byte)[:w.obufLen]
		checksum := crc(uncompressed)

		
		chunkType := uint8(chunkTypeUncompressedData)
		chunkLen := 4 + len(uncompressed)

		
		n := binary.PutUvarint(obuf[obufHeaderLen:], uint64(len(uncompressed)))
		n2 := w.encodeBlock(obuf[obufHeaderLen+n:], uncompressed)

		if n2 > 0 {
			chunkType = uint8(chunkTypeCompressedData)
			chunkLen = 4 + n + n2
			obuf = obuf[:obufHeaderLen+n+n2]
		} else {
			obuf = obuf[:8]
		}

		
		obuf[0] = chunkType
		obuf[1] = uint8(chunkLen >> 0)
		obuf[2] = uint8(chunkLen >> 8)
		obuf[3] = uint8(chunkLen >> 16)
		obuf[4] = uint8(checksum >> 0)
		obuf[5] = uint8(checksum >> 8)
		obuf[6] = uint8(checksum >> 16)
		obuf[7] = uint8(checksum >> 24)

		n, err := w.writer.Write(obuf)
		if err != nil {
			return 0, w.err(err)
		}
		if n != len(obuf) {
			return 0, w.err(io.ErrShortWrite)
		}
		w.err(w.index.add(w.written, w.uncompWritten))
		w.written += int64(n)
		w.uncompWritten += int64(len(uncompressed))

		if chunkType == chunkTypeUncompressedData {
			
			n, err := w.writer.Write(uncompressed)
			if err != nil {
				return 0, w.err(err)
			}
			if n != len(uncompressed) {
				return 0, w.err(io.ErrShortWrite)
			}
			w.written += int64(n)
		}
		w.buffers.Put(obuf)
		
		nRet += len(uncompressed)
	}
	return nRet, nil
}



func (w *Writer) Flush() error {
	if err := w.err(nil); err != nil {
		return err
	}

	
	if len(w.ibuf) != 0 {
		if !w.wroteStreamHeader {
			_, err := w.writeSync(w.ibuf)
			w.ibuf = w.ibuf[:0]
			return w.err(err)
		} else {
			_, err := w.write(w.ibuf)
			w.ibuf = w.ibuf[:0]
			err = w.err(err)
			if err != nil {
				return err
			}
		}
	}
	if w.output == nil {
		return w.err(nil)
	}

	
	res := make(chan result)
	w.output <- res
	
	res <- result{b: nil, startOffset: w.uncompWritten}
	
	<-res
	return w.err(nil)
}




func (w *Writer) Close() error {
	_, err := w.closeIndex(w.appendIndex)
	return err
}



func (w *Writer) CloseIndex() ([]byte, error) {
	return w.closeIndex(true)
}

func (w *Writer) closeIndex(idx bool) ([]byte, error) {
	err := w.Flush()
	if w.output != nil {
		close(w.output)
		w.writerWg.Wait()
		w.output = nil
	}

	var index []byte
	if w.err(err) == nil && w.writer != nil {
		
		if idx {
			compSize := int64(-1)
			if w.pad <= 1 {
				compSize = w.written
			}
			index = w.index.appendTo(w.ibuf[:0], w.uncompWritten, compSize)
			
			if w.appendIndex {
				w.written += int64(len(index))
			}
		}

		if w.pad > 1 {
			tmp := w.ibuf[:0]
			if len(index) > 0 {
				
				tmp = w.buffers.Get().([]byte)[:0]
				defer w.buffers.Put(tmp)
			}
			add := calcSkippableFrame(w.written, int64(w.pad))
			frame, err := skippableFrame(tmp, add, w.randSrc)
			if err = w.err(err); err != nil {
				return nil, err
			}
			n, err2 := w.writer.Write(frame)
			if err2 == nil && n != len(frame) {
				err2 = io.ErrShortWrite
			}
			_ = w.err(err2)
		}
		if len(index) > 0 && w.appendIndex {
			n, err2 := w.writer.Write(index)
			if err2 == nil && n != len(index) {
				err2 = io.ErrShortWrite
			}
			_ = w.err(err2)
		}
	}
	err = w.err(errClosed)
	if err == errClosed {
		return index, nil
	}
	return nil, err
}





func calcSkippableFrame(written, wantMultiple int64) int {
	if wantMultiple <= 0 {
		panic("wantMultiple <= 0")
	}
	if written < 0 {
		panic("written < 0")
	}
	leftOver := written % wantMultiple
	if leftOver == 0 {
		return 0
	}
	toAdd := wantMultiple - leftOver
	for toAdd < skippableFrameHeader {
		toAdd += wantMultiple
	}
	return int(toAdd)
}



func skippableFrame(dst []byte, total int, r io.Reader) ([]byte, error) {
	if total == 0 {
		return dst, nil
	}
	if total < skippableFrameHeader {
		return dst, fmt.Errorf("s2: requested skippable frame (%d) < 4", total)
	}
	if int64(total) >= maxBlockSize+skippableFrameHeader {
		return dst, fmt.Errorf("s2: requested skippable frame (%d) >= max 1<<24", total)
	}
	
	dst = append(dst, chunkTypePadding)
	f := uint32(total - skippableFrameHeader)
	
	dst = append(dst, uint8(f), uint8(f>>8), uint8(f>>16))
	
	start := len(dst)
	dst = append(dst, make([]byte, f)...)
	_, err := io.ReadFull(r, dst[start:])
	return dst, err
}

var errClosed = errors.New("s2: Writer is closed")


type WriterOption func(*Writer) error





func WriterConcurrency(n int) WriterOption {
	return func(w *Writer) error {
		if n <= 0 {
			return errors.New("concurrency must be at least 1")
		}
		w.concurrency = n
		return nil
	}
}



func WriterAddIndex() WriterOption {
	return func(w *Writer) error {
		w.appendIndex = true
		return nil
	}
}




func WriterBetterCompression() WriterOption {
	return func(w *Writer) error {
		w.level = levelBetter
		return nil
	}
}




func WriterBestCompression() WriterOption {
	return func(w *Writer) error {
		w.level = levelBest
		return nil
	}
}




func WriterUncompressed() WriterOption {
	return func(w *Writer) error {
		w.level = levelUncompressed
		return nil
	}
}












func WriterBlockSize(n int) WriterOption {
	return func(w *Writer) error {
		if w.snappy && n > maxSnappyBlockSize || n < minBlockSize {
			return errors.New("s2: block size too large. Must be <= 64K and >=4KB on for snappy compatible output")
		}
		if n > maxBlockSize || n < minBlockSize {
			return errors.New("s2: block size too large. Must be <= 4MB and >=4KB")
		}
		w.blockSize = n
		return nil
	}
}







func WriterPadding(n int) WriterOption {
	return func(w *Writer) error {
		if n <= 0 {
			return fmt.Errorf("s2: padding must be at least 1")
		}
		
		if n == 1 {
			w.pad = 0
		}
		if n > maxBlockSize {
			return fmt.Errorf("s2: padding must less than 4MB")
		}
		w.pad = n
		return nil
	}
}



func WriterPaddingSrc(reader io.Reader) WriterOption {
	return func(w *Writer) error {
		w.randSrc = reader
		return nil
	}
}




func WriterSnappyCompat() WriterOption {
	return func(w *Writer) error {
		w.snappy = true
		if w.blockSize > 64<<10 {
			
			
			w.blockSize = (64 << 10) - 8
		}
		return nil
	}
}







func WriterFlushOnWrite() WriterOption {
	return func(w *Writer) error {
		w.flushOnWrite = true
		return nil
	}
}







func WriterCustomEncoder(fn func(dst, src []byte) int) WriterOption {
	return func(w *Writer) error {
		w.customEnc = fn
		return nil
	}
}
