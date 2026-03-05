




package s2

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"runtime"
	"sync"
)


type ErrCantSeek struct {
	Reason string
}


func (e ErrCantSeek) Error() string {
	return fmt.Sprintf("s2: Can't seek because %s", e.Reason)
}




func NewReader(r io.Reader, opts ...ReaderOption) *Reader {
	nr := Reader{
		r:        r,
		maxBlock: maxBlockSize,
	}
	for _, opt := range opts {
		if err := opt(&nr); err != nil {
			nr.err = err
			return &nr
		}
	}
	nr.maxBufSize = MaxEncodedLen(nr.maxBlock) + checksumSize
	if nr.lazyBuf > 0 {
		nr.buf = make([]byte, MaxEncodedLen(nr.lazyBuf)+checksumSize)
	} else {
		nr.buf = make([]byte, MaxEncodedLen(defaultBlockSize)+checksumSize)
	}
	nr.readHeader = nr.ignoreStreamID
	nr.paramsOK = true
	return &nr
}


type ReaderOption func(*Reader) error









func ReaderMaxBlockSize(blockSize int) ReaderOption {
	return func(r *Reader) error {
		if blockSize > maxBlockSize || blockSize <= 0 {
			return errors.New("s2: block size too large. Must be <= 4MB and > 0")
		}
		if r.lazyBuf == 0 && blockSize < defaultBlockSize {
			r.lazyBuf = blockSize
		}
		r.maxBlock = blockSize
		return nil
	}
}






func ReaderAllocBlock(blockSize int) ReaderOption {
	return func(r *Reader) error {
		if blockSize > maxBlockSize || blockSize < 1024 {
			return errors.New("s2: invalid ReaderAllocBlock. Must be <= 4MB and >= 1024")
		}
		r.lazyBuf = blockSize
		return nil
	}
}




func ReaderIgnoreStreamIdentifier() ReaderOption {
	return func(r *Reader) error {
		r.ignoreStreamID = true
		return nil
	}
}






func ReaderSkippableCB(id uint8, fn func(r io.Reader) error) ReaderOption {
	return func(r *Reader) error {
		if id < 0x80 || id > 0xfd {
			return fmt.Errorf("ReaderSkippableCB: Invalid id provided, must be 0x80-0xfd (inclusive)")
		}
		r.skippableCB[id] = fn
		return nil
	}
}


func ReaderIgnoreCRC() ReaderOption {
	return func(r *Reader) error {
		r.ignoreCRC = true
		return nil
	}
}


type Reader struct {
	r           io.Reader
	err         error
	decoded     []byte
	buf         []byte
	skippableCB [0x80]func(r io.Reader) error
	blockStart  int64 
	index       *Index

	
	i, j int
	
	maxBlock int
	
	maxBufSize int
	
	lazyBuf        int
	readHeader     bool
	paramsOK       bool
	snappyFrame    bool
	ignoreStreamID bool
	ignoreCRC      bool
}




func (r *Reader) GetBufferCapacity() int {
	return cap(r.buf)
}



func (r *Reader) ensureBufferSize(n int) bool {
	if n > r.maxBufSize {
		r.err = ErrCorrupt
		return false
	}
	if cap(r.buf) >= n {
		return true
	}
	
	r.buf = make([]byte, n)
	return true
}




func (r *Reader) Reset(reader io.Reader) {
	if !r.paramsOK {
		return
	}
	r.index = nil
	r.r = reader
	r.err = nil
	r.i = 0
	r.j = 0
	r.blockStart = 0
	r.readHeader = r.ignoreStreamID
}

func (r *Reader) readFull(p []byte, allowEOF bool) (ok bool) {
	if _, r.err = io.ReadFull(r.r, p); r.err != nil {
		if r.err == io.ErrUnexpectedEOF || (r.err == io.EOF && !allowEOF) {
			r.err = ErrCorrupt
		}
		return false
	}
	return true
}





func (r *Reader) skippable(tmp []byte, n int, allowEOF bool, id uint8) (ok bool) {
	if id < 0x80 {
		r.err = fmt.Errorf("interbal error: skippable id < 0x80")
		return false
	}
	if fn := r.skippableCB[id-0x80]; fn != nil {
		rd := io.LimitReader(r.r, int64(n))
		r.err = fn(rd)
		if r.err != nil {
			return false
		}
		_, r.err = io.CopyBuffer(ioutil.Discard, rd, tmp)
		return r.err == nil
	}
	if rs, ok := r.r.(io.ReadSeeker); ok {
		_, err := rs.Seek(int64(n), io.SeekCurrent)
		if err == nil {
			return true
		}
		if err == io.ErrUnexpectedEOF || (r.err == io.EOF && !allowEOF) {
			r.err = ErrCorrupt
			return false
		}
	}
	for n > 0 {
		if n < len(tmp) {
			tmp = tmp[:n]
		}
		if _, r.err = io.ReadFull(r.r, tmp); r.err != nil {
			if r.err == io.ErrUnexpectedEOF || (r.err == io.EOF && !allowEOF) {
				r.err = ErrCorrupt
			}
			return false
		}
		n -= len(tmp)
	}
	return true
}


func (r *Reader) Read(p []byte) (int, error) {
	if r.err != nil {
		return 0, r.err
	}
	for {
		if r.i < r.j {
			n := copy(p, r.decoded[r.i:r.j])
			r.i += n
			return n, nil
		}
		if !r.readFull(r.buf[:4], true) {
			return 0, r.err
		}
		chunkType := r.buf[0]
		if !r.readHeader {
			if chunkType != chunkTypeStreamIdentifier {
				r.err = ErrCorrupt
				return 0, r.err
			}
			r.readHeader = true
		}
		chunkLen := int(r.buf[1]) | int(r.buf[2])<<8 | int(r.buf[3])<<16

		
		
		switch chunkType {
		case chunkTypeCompressedData:
			r.blockStart += int64(r.j)
			
			if chunkLen < checksumSize {
				r.err = ErrCorrupt
				return 0, r.err
			}
			if !r.ensureBufferSize(chunkLen) {
				if r.err == nil {
					r.err = ErrUnsupported
				}
				return 0, r.err
			}
			buf := r.buf[:chunkLen]
			if !r.readFull(buf, false) {
				return 0, r.err
			}
			checksum := uint32(buf[0]) | uint32(buf[1])<<8 | uint32(buf[2])<<16 | uint32(buf[3])<<24
			buf = buf[checksumSize:]

			n, err := DecodedLen(buf)
			if err != nil {
				r.err = err
				return 0, r.err
			}
			if r.snappyFrame && n > maxSnappyBlockSize {
				r.err = ErrCorrupt
				return 0, r.err
			}

			if n > len(r.decoded) {
				if n > r.maxBlock {
					r.err = ErrCorrupt
					return 0, r.err
				}
				r.decoded = make([]byte, n)
			}
			if _, err := Decode(r.decoded, buf); err != nil {
				r.err = err
				return 0, r.err
			}
			if !r.ignoreCRC && crc(r.decoded[:n]) != checksum {
				r.err = ErrCRC
				return 0, r.err
			}
			r.i, r.j = 0, n
			continue

		case chunkTypeUncompressedData:
			r.blockStart += int64(r.j)
			
			if chunkLen < checksumSize {
				r.err = ErrCorrupt
				return 0, r.err
			}
			if !r.ensureBufferSize(chunkLen) {
				if r.err == nil {
					r.err = ErrUnsupported
				}
				return 0, r.err
			}
			buf := r.buf[:checksumSize]
			if !r.readFull(buf, false) {
				return 0, r.err
			}
			checksum := uint32(buf[0]) | uint32(buf[1])<<8 | uint32(buf[2])<<16 | uint32(buf[3])<<24
			
			n := chunkLen - checksumSize
			if r.snappyFrame && n > maxSnappyBlockSize {
				r.err = ErrCorrupt
				return 0, r.err
			}
			if n > len(r.decoded) {
				if n > r.maxBlock {
					r.err = ErrCorrupt
					return 0, r.err
				}
				r.decoded = make([]byte, n)
			}
			if !r.readFull(r.decoded[:n], false) {
				return 0, r.err
			}
			if !r.ignoreCRC && crc(r.decoded[:n]) != checksum {
				r.err = ErrCRC
				return 0, r.err
			}
			r.i, r.j = 0, n
			continue

		case chunkTypeStreamIdentifier:
			
			if chunkLen != len(magicBody) {
				r.err = ErrCorrupt
				return 0, r.err
			}
			if !r.readFull(r.buf[:len(magicBody)], false) {
				return 0, r.err
			}
			if string(r.buf[:len(magicBody)]) != magicBody {
				if string(r.buf[:len(magicBody)]) != magicBodySnappy {
					r.err = ErrCorrupt
					return 0, r.err
				} else {
					r.snappyFrame = true
				}
			} else {
				r.snappyFrame = false
			}
			continue
		}

		if chunkType <= 0x7f {
			
			
			r.err = ErrUnsupported
			return 0, r.err
		}
		
		
		if chunkLen > maxChunkSize {
			
			r.err = ErrUnsupported
			return 0, r.err
		}

		
		if !r.skippable(r.buf, chunkLen, false, chunkType) {
			return 0, r.err
		}
	}
}







func (r *Reader) DecodeConcurrent(w io.Writer, concurrent int) (written int64, err error) {
	if r.i > 0 || r.j > 0 || r.blockStart > 0 {
		return 0, errors.New("DecodeConcurrent called after ")
	}
	if concurrent <= 0 {
		concurrent = runtime.NumCPU()
	}

	
	var errMu sync.Mutex
	var aErr error
	setErr := func(e error) (ok bool) {
		errMu.Lock()
		defer errMu.Unlock()
		if e == nil {
			return aErr == nil
		}
		if aErr == nil {
			aErr = e
		}
		return false
	}
	hasErr := func() (ok bool) {
		errMu.Lock()
		v := aErr != nil
		errMu.Unlock()
		return v
	}

	var aWritten int64
	toRead := make(chan []byte, concurrent)
	writtenBlocks := make(chan []byte, concurrent)
	queue := make(chan chan []byte, concurrent)
	reUse := make(chan chan []byte, concurrent)
	for i := 0; i < concurrent; i++ {
		toRead <- make([]byte, 0, r.maxBufSize)
		writtenBlocks <- make([]byte, 0, r.maxBufSize)
		reUse <- make(chan []byte, 1)
	}
	
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for toWrite := range queue {
			entry := <-toWrite
			reUse <- toWrite
			if hasErr() {
				writtenBlocks <- entry
				continue
			}
			n, err := w.Write(entry)
			want := len(entry)
			writtenBlocks <- entry
			if err != nil {
				setErr(err)
				continue
			}
			if n != want {
				setErr(io.ErrShortWrite)
				continue
			}
			aWritten += int64(n)
		}
	}()

	
	defer func() {
		close(queue)
		if r.err != nil {
			err = r.err
			setErr(r.err)
		}
		wg.Wait()
		if err == nil {
			err = aErr
		}
		written = aWritten
	}()

	for !hasErr() {
		if !r.readFull(r.buf[:4], true) {
			if r.err == io.EOF {
				r.err = nil
			}
			return 0, r.err
		}
		chunkType := r.buf[0]
		if !r.readHeader {
			if chunkType != chunkTypeStreamIdentifier {
				r.err = ErrCorrupt
				return 0, r.err
			}
			r.readHeader = true
		}
		chunkLen := int(r.buf[1]) | int(r.buf[2])<<8 | int(r.buf[3])<<16

		
		
		switch chunkType {
		case chunkTypeCompressedData:
			r.blockStart += int64(r.j)
			
			if chunkLen < checksumSize {
				r.err = ErrCorrupt
				return 0, r.err
			}
			if chunkLen > r.maxBufSize {
				r.err = ErrCorrupt
				return 0, r.err
			}
			orgBuf := <-toRead
			buf := orgBuf[:chunkLen]

			if !r.readFull(buf, false) {
				return 0, r.err
			}

			checksum := uint32(buf[0]) | uint32(buf[1])<<8 | uint32(buf[2])<<16 | uint32(buf[3])<<24
			buf = buf[checksumSize:]

			n, err := DecodedLen(buf)
			if err != nil {
				r.err = err
				return 0, r.err
			}
			if r.snappyFrame && n > maxSnappyBlockSize {
				r.err = ErrCorrupt
				return 0, r.err
			}

			if n > r.maxBlock {
				r.err = ErrCorrupt
				return 0, r.err
			}
			wg.Add(1)

			decoded := <-writtenBlocks
			entry := <-reUse
			queue <- entry
			go func() {
				defer wg.Done()
				decoded = decoded[:n]
				_, err := Decode(decoded, buf)
				toRead <- orgBuf
				if err != nil {
					writtenBlocks <- decoded
					setErr(err)
					return
				}
				if !r.ignoreCRC && crc(decoded) != checksum {
					writtenBlocks <- decoded
					setErr(ErrCRC)
					return
				}
				entry <- decoded
			}()
			continue

		case chunkTypeUncompressedData:

			
			if chunkLen < checksumSize {
				r.err = ErrCorrupt
				return 0, r.err
			}
			if chunkLen > r.maxBufSize {
				r.err = ErrCorrupt
				return 0, r.err
			}
			
			orgBuf := <-writtenBlocks
			buf := orgBuf[:checksumSize]
			if !r.readFull(buf, false) {
				return 0, r.err
			}
			checksum := uint32(buf[0]) | uint32(buf[1])<<8 | uint32(buf[2])<<16 | uint32(buf[3])<<24
			
			n := chunkLen - checksumSize

			if r.snappyFrame && n > maxSnappyBlockSize {
				r.err = ErrCorrupt
				return 0, r.err
			}
			if n > r.maxBlock {
				r.err = ErrCorrupt
				return 0, r.err
			}
			
			buf = orgBuf[:n]
			if !r.readFull(buf, false) {
				return 0, r.err
			}

			if !r.ignoreCRC && crc(buf) != checksum {
				r.err = ErrCRC
				return 0, r.err
			}
			entry := <-reUse
			queue <- entry
			entry <- buf
			continue

		case chunkTypeStreamIdentifier:
			
			if chunkLen != len(magicBody) {
				r.err = ErrCorrupt
				return 0, r.err
			}
			if !r.readFull(r.buf[:len(magicBody)], false) {
				return 0, r.err
			}
			if string(r.buf[:len(magicBody)]) != magicBody {
				if string(r.buf[:len(magicBody)]) != magicBodySnappy {
					r.err = ErrCorrupt
					return 0, r.err
				} else {
					r.snappyFrame = true
				}
			} else {
				r.snappyFrame = false
			}
			continue
		}

		if chunkType <= 0x7f {
			
			
			r.err = ErrUnsupported
			return 0, r.err
		}
		
		
		if chunkLen > maxChunkSize {
			
			r.err = ErrUnsupported
			return 0, r.err
		}

		
		if !r.skippable(r.buf, chunkLen, false, chunkType) {
			return 0, r.err
		}
	}
	return 0, r.err
}






func (r *Reader) Skip(n int64) error {
	if n < 0 {
		return errors.New("attempted negative skip")
	}
	if r.err != nil {
		return r.err
	}

	for n > 0 {
		if r.i < r.j {
			
			
			left := int64(r.j - r.i)
			if left >= n {
				tmp := int64(r.i) + n
				if tmp > math.MaxInt32 {
					return errors.New("s2: internal overflow in skip")
				}
				r.i = int(tmp)
				return nil
			}
			n -= int64(r.j - r.i)
			r.i = r.j
		}

		
		if !r.readFull(r.buf[:4], true) {
			if r.err == io.EOF {
				r.err = io.ErrUnexpectedEOF
			}
			return r.err
		}
		chunkType := r.buf[0]
		if !r.readHeader {
			if chunkType != chunkTypeStreamIdentifier {
				r.err = ErrCorrupt
				return r.err
			}
			r.readHeader = true
		}
		chunkLen := int(r.buf[1]) | int(r.buf[2])<<8 | int(r.buf[3])<<16

		
		
		switch chunkType {
		case chunkTypeCompressedData:
			r.blockStart += int64(r.j)
			
			if chunkLen < checksumSize {
				r.err = ErrCorrupt
				return r.err
			}
			if !r.ensureBufferSize(chunkLen) {
				if r.err == nil {
					r.err = ErrUnsupported
				}
				return r.err
			}
			buf := r.buf[:chunkLen]
			if !r.readFull(buf, false) {
				return r.err
			}
			checksum := uint32(buf[0]) | uint32(buf[1])<<8 | uint32(buf[2])<<16 | uint32(buf[3])<<24
			buf = buf[checksumSize:]

			dLen, err := DecodedLen(buf)
			if err != nil {
				r.err = err
				return r.err
			}
			if dLen > r.maxBlock {
				r.err = ErrCorrupt
				return r.err
			}
			
			if int64(dLen) > n {
				if len(r.decoded) < dLen {
					r.decoded = make([]byte, dLen)
				}
				if _, err := Decode(r.decoded, buf); err != nil {
					r.err = err
					return r.err
				}
				if crc(r.decoded[:dLen]) != checksum {
					r.err = ErrCorrupt
					return r.err
				}
			} else {
				
				n -= int64(dLen)
				r.blockStart += int64(dLen)
				dLen = 0
			}
			r.i, r.j = 0, dLen
			continue
		case chunkTypeUncompressedData:
			r.blockStart += int64(r.j)
			
			if chunkLen < checksumSize {
				r.err = ErrCorrupt
				return r.err
			}
			if !r.ensureBufferSize(chunkLen) {
				if r.err != nil {
					r.err = ErrUnsupported
				}
				return r.err
			}
			buf := r.buf[:checksumSize]
			if !r.readFull(buf, false) {
				return r.err
			}
			checksum := uint32(buf[0]) | uint32(buf[1])<<8 | uint32(buf[2])<<16 | uint32(buf[3])<<24
			
			n2 := chunkLen - checksumSize
			if n2 > len(r.decoded) {
				if n2 > r.maxBlock {
					r.err = ErrCorrupt
					return r.err
				}
				r.decoded = make([]byte, n2)
			}
			if !r.readFull(r.decoded[:n2], false) {
				return r.err
			}
			if int64(n2) < n {
				if crc(r.decoded[:n2]) != checksum {
					r.err = ErrCorrupt
					return r.err
				}
			}
			r.i, r.j = 0, n2
			continue
		case chunkTypeStreamIdentifier:
			
			if chunkLen != len(magicBody) {
				r.err = ErrCorrupt
				return r.err
			}
			if !r.readFull(r.buf[:len(magicBody)], false) {
				return r.err
			}
			if string(r.buf[:len(magicBody)]) != magicBody {
				if string(r.buf[:len(magicBody)]) != magicBodySnappy {
					r.err = ErrCorrupt
					return r.err
				}
			}

			continue
		}

		if chunkType <= 0x7f {
			
			r.err = ErrUnsupported
			return r.err
		}
		if chunkLen > maxChunkSize {
			r.err = ErrUnsupported
			return r.err
		}
		
		
		if !r.skippable(r.buf, chunkLen, false, chunkType) {
			return r.err
		}
	}
	return nil
}



type ReadSeeker struct {
	*Reader
	readAtMu sync.Mutex
}














func (r *Reader) ReadSeeker(random bool, index []byte) (*ReadSeeker, error) {
	
	if len(index) != 0 {
		if r.index == nil {
			r.index = &Index{}
		}
		if _, err := r.index.Load(index); err != nil {
			return nil, ErrCantSeek{Reason: "loading index returned: " + err.Error()}
		}
	}

	
	rs, ok := r.r.(io.ReadSeeker)
	if !ok {
		if !random {
			return &ReadSeeker{Reader: r}, nil
		}
		return nil, ErrCantSeek{Reason: "input stream isn't seekable"}
	}

	if r.index != nil {
		
		return &ReadSeeker{Reader: r}, nil
	}

	
	r.index = &Index{}

	
	pos, err := rs.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, ErrCantSeek{Reason: "seeking input returned: " + err.Error()}
	}
	err = r.index.LoadStream(rs)
	if err != nil {
		if err == ErrUnsupported {
			
			if !random {
				_, err = rs.Seek(pos, io.SeekStart)
				if err != nil {
					return nil, ErrCantSeek{Reason: "resetting stream returned: " + err.Error()}
				}
				r.index = nil
				return &ReadSeeker{Reader: r}, nil
			}
			return nil, ErrCantSeek{Reason: "input stream does not contain an index"}
		}
		return nil, ErrCantSeek{Reason: "reading index returned: " + err.Error()}
	}

	
	_, err = rs.Seek(pos, io.SeekStart)
	if err != nil {
		return nil, ErrCantSeek{Reason: "seeking input returned: " + err.Error()}
	}
	return &ReadSeeker{Reader: r}, nil
}


func (r *ReadSeeker) Seek(offset int64, whence int) (int64, error) {
	if r.err != nil {
		if !errors.Is(r.err, io.EOF) {
			return 0, r.err
		}
		
		r.err = nil
	}

	
	absOffset := offset

	switch whence {
	case io.SeekStart:
	case io.SeekCurrent:
		absOffset = r.blockStart + int64(r.i) + offset
	case io.SeekEnd:
		if r.index == nil {
			return 0, ErrUnsupported
		}
		absOffset = r.index.TotalUncompressed + offset
	default:
		r.err = ErrUnsupported
		return 0, r.err
	}

	if absOffset < 0 {
		return 0, errors.New("seek before start of file")
	}

	if !r.readHeader {
		
		_, r.err = r.Read([]byte{})
		if r.err != nil {
			return 0, r.err
		}
	}

	
	
	if absOffset >= r.blockStart && absOffset < r.blockStart+int64(r.j) {
		r.i = int(absOffset - r.blockStart)
		return r.blockStart + int64(r.i), nil
	}

	rs, ok := r.r.(io.ReadSeeker)
	if r.index == nil || !ok {
		currOffset := r.blockStart + int64(r.i)
		if absOffset >= currOffset {
			err := r.Skip(absOffset - currOffset)
			return r.blockStart + int64(r.i), err
		}
		return 0, ErrUnsupported
	}

	
	c, u, err := r.index.Find(absOffset)
	if err != nil {
		return r.blockStart + int64(r.i), err
	}

	
	_, err = rs.Seek(c, io.SeekStart)
	if err != nil {
		return 0, err
	}

	r.i = r.j                     
	r.blockStart = u - int64(r.j) 
	if u < absOffset {
		
		return absOffset, r.Skip(absOffset - u)
	}
	if u > absOffset {
		return 0, fmt.Errorf("s2 seek: (internal error) u (%d) > absOffset (%d)", u, absOffset)
	}
	return absOffset, nil
}























func (r *ReadSeeker) ReadAt(p []byte, offset int64) (int, error) {
	r.readAtMu.Lock()
	defer r.readAtMu.Unlock()
	_, err := r.Seek(offset, io.SeekStart)
	if err != nil {
		return 0, err
	}
	n := 0
	for n < len(p) {
		n2, err := r.Read(p[n:])
		if err != nil {
			
			return n + n2, err
		}
		n += n2
	}
	return n, nil
}


func (r *Reader) ReadByte() (byte, error) {
	if r.err != nil {
		return 0, r.err
	}
	if r.i < r.j {
		c := r.decoded[r.i]
		r.i++
		return c, nil
	}
	var tmp [1]byte
	for i := 0; i < 10; i++ {
		n, err := r.Read(tmp[:])
		if err != nil {
			return 0, err
		}
		if n == 1 {
			return tmp[0], nil
		}
	}
	return 0, io.ErrNoProgress
}







func (r *Reader) SkippableCB(id uint8, fn func(r io.Reader) error) error {
	if id < 0x80 || id > chunkTypePadding {
		return fmt.Errorf("ReaderSkippableCB: Invalid id provided, must be 0x80-0xfe (inclusive)")
	}
	r.skippableCB[id] = fn
	return nil
}
