

package mem

import (
	"io"
)

const (
	
	readAllBufSize = 32 * 1024
)















type BufferSlice []Buffer







func (s BufferSlice) Len() int {
	var length int
	for _, b := range s {
		length += b.Len()
	}
	return length
}


func (s BufferSlice) Ref() {
	for _, b := range s {
		b.Ref()
	}
}


func (s BufferSlice) Free() {
	for _, b := range s {
		b.Free()
	}
}





func (s BufferSlice) CopyTo(dst []byte) int {
	off := 0
	for _, b := range s {
		off += copy(dst[off:], b.ReadOnlyData())
	}
	return off
}



func (s BufferSlice) Materialize() []byte {
	l := s.Len()
	if l == 0 {
		return nil
	}
	out := make([]byte, l)
	s.CopyTo(out)
	return out
}







func (s BufferSlice) MaterializeToBuffer(pool BufferPool) Buffer {
	if len(s) == 1 {
		s[0].Ref()
		return s[0]
	}
	sLen := s.Len()
	if sLen == 0 {
		return emptyBuffer{}
	}
	buf := pool.Get(sLen)
	s.CopyTo(*buf)
	return NewBuffer(buf, pool)
}



func (s BufferSlice) Reader() Reader {
	s.Ref()
	return &sliceReader{
		data: s,
		len:  s.Len(),
	}
}





type Reader interface {
	io.Reader
	io.ByteReader
	
	
	Close() error
	
	Remaining() int
}

type sliceReader struct {
	data BufferSlice
	len  int
	
	bufferIdx int
}

func (r *sliceReader) Remaining() int {
	return r.len
}

func (r *sliceReader) Close() error {
	r.data.Free()
	r.data = nil
	r.len = 0
	return nil
}

func (r *sliceReader) freeFirstBufferIfEmpty() bool {
	if len(r.data) == 0 || r.bufferIdx != len(r.data[0].ReadOnlyData()) {
		return false
	}

	r.data[0].Free()
	r.data = r.data[1:]
	r.bufferIdx = 0
	return true
}

func (r *sliceReader) Read(buf []byte) (n int, _ error) {
	if r.len == 0 {
		return 0, io.EOF
	}

	for len(buf) != 0 && r.len != 0 {
		
		
		data := r.data[0].ReadOnlyData()
		copied := copy(buf, data[r.bufferIdx:])
		r.len -= copied       
		r.bufferIdx += copied 
		n += copied           
		buf = buf[copied:]    

		
		
		r.freeFirstBufferIfEmpty()
	}

	return n, nil
}

func (r *sliceReader) ReadByte() (byte, error) {
	if r.len == 0 {
		return 0, io.EOF
	}

	
	
	for r.freeFirstBufferIfEmpty() {
	}

	b := r.data[0].ReadOnlyData()[r.bufferIdx]
	r.len--
	r.bufferIdx++
	
	r.freeFirstBufferIfEmpty()
	return b, nil
}

var _ io.Writer = (*writer)(nil)

type writer struct {
	buffers *BufferSlice
	pool    BufferPool
}

func (w *writer) Write(p []byte) (n int, err error) {
	b := Copy(p, w.pool)
	*w.buffers = append(*w.buffers, b)
	return b.Len(), nil
}





func NewWriter(buffers *BufferSlice, pool BufferPool) io.Writer {
	return &writer{buffers: buffers, pool: pool}
}









func ReadAll(r io.Reader, pool BufferPool) (BufferSlice, error) {
	var result BufferSlice
	if wt, ok := r.(io.WriterTo); ok {
		
		
		
		
		w := NewWriter(&result, pool)
		_, err := wt.WriteTo(w)
		return result, err
	}
nextBuffer:
	for {
		buf := pool.Get(readAllBufSize)
		
		
		*buf = (*buf)[:cap(*buf)]
		usedCap := 0
		for {
			n, err := r.Read((*buf)[usedCap:])
			usedCap += n
			if err != nil {
				if usedCap == 0 {
					
					pool.Put(buf)
				} else {
					*buf = (*buf)[:usedCap]
					result = append(result, NewBuffer(buf, pool))
				}
				if err == io.EOF {
					err = nil
				}
				return result, err
			}
			if len(*buf) == usedCap {
				result = append(result, NewBuffer(buf, pool))
				continue nextBuffer
			}
		}
	}
}
