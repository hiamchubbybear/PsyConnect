








package mem

import (
	"fmt"
	"sync"
	"sync/atomic"
)












type Buffer interface {
	
	
	ReadOnlyData() []byte
	
	Ref()
	
	
	Free()
	
	Len() int

	split(n int) (left, right Buffer)
	read(buf []byte) (int, Buffer)
}

var (
	bufferPoolingThreshold = 1 << 10

	bufferObjectPool = sync.Pool{New: func() any { return new(buffer) }}
	refObjectPool    = sync.Pool{New: func() any { return new(atomic.Int32) }}
)




func IsBelowBufferPoolingThreshold(size int) bool {
	return size <= bufferPoolingThreshold
}

type buffer struct {
	origData *[]byte
	data     []byte
	refs     *atomic.Int32
	pool     BufferPool
}

func newBuffer() *buffer {
	return bufferObjectPool.Get().(*buffer)
}









func NewBuffer(data *[]byte, pool BufferPool) Buffer {
	
	
	
	
	if pool == nil || IsBelowBufferPoolingThreshold(cap(*data)) {
		return (SliceBuffer)(*data)
	}
	b := newBuffer()
	b.origData = data
	b.data = *data
	b.pool = pool
	b.refs = refObjectPool.Get().(*atomic.Int32)
	b.refs.Add(1)
	return b
}







func Copy(data []byte, pool BufferPool) Buffer {
	if IsBelowBufferPoolingThreshold(len(data)) {
		buf := make(SliceBuffer, len(data))
		copy(buf, data)
		return buf
	}

	buf := pool.Get(len(data))
	copy(*buf, data)
	return NewBuffer(buf, pool)
}

func (b *buffer) ReadOnlyData() []byte {
	if b.refs == nil {
		panic("Cannot read freed buffer")
	}
	return b.data
}

func (b *buffer) Ref() {
	if b.refs == nil {
		panic("Cannot ref freed buffer")
	}
	b.refs.Add(1)
}

func (b *buffer) Free() {
	if b.refs == nil {
		panic("Cannot free freed buffer")
	}

	refs := b.refs.Add(-1)
	switch {
	case refs > 0:
		return
	case refs == 0:
		if b.pool != nil {
			b.pool.Put(b.origData)
		}

		refObjectPool.Put(b.refs)
		b.origData = nil
		b.data = nil
		b.refs = nil
		b.pool = nil
		bufferObjectPool.Put(b)
	default:
		panic("Cannot free freed buffer")
	}
}

func (b *buffer) Len() int {
	return len(b.ReadOnlyData())
}

func (b *buffer) split(n int) (Buffer, Buffer) {
	if b.refs == nil {
		panic("Cannot split freed buffer")
	}

	b.refs.Add(1)
	split := newBuffer()
	split.origData = b.origData
	split.data = b.data[n:]
	split.refs = b.refs
	split.pool = b.pool

	b.data = b.data[:n]

	return b, split
}

func (b *buffer) read(buf []byte) (int, Buffer) {
	if b.refs == nil {
		panic("Cannot read freed buffer")
	}

	n := copy(buf, b.data)
	if n == len(b.data) {
		b.Free()
		return n, nil
	}

	b.data = b.data[n:]
	return n, b
}

func (b *buffer) String() string {
	return fmt.Sprintf("mem.Buffer(%p, data: %p, length: %d)", b, b.ReadOnlyData(), len(b.ReadOnlyData()))
}



func ReadUnsafe(dst []byte, buf Buffer) (int, Buffer) {
	return buf.read(dst)
}




func SplitUnsafe(buf Buffer, n int) (left, right Buffer) {
	return buf.split(n)
}

type emptyBuffer struct{}

func (e emptyBuffer) ReadOnlyData() []byte {
	return nil
}

func (e emptyBuffer) Ref()  {}
func (e emptyBuffer) Free() {}

func (e emptyBuffer) Len() int {
	return 0
}

func (e emptyBuffer) split(int) (left, right Buffer) {
	return e, e
}

func (e emptyBuffer) read([]byte) (int, Buffer) {
	return 0, e
}



type SliceBuffer []byte


func (s SliceBuffer) ReadOnlyData() []byte { return s }


func (s SliceBuffer) Ref() {}


func (s SliceBuffer) Free() {}


func (s SliceBuffer) Len() int { return len(s) }

func (s SliceBuffer) split(n int) (left, right Buffer) {
	return s[:n], s[n:]
}

func (s SliceBuffer) read(buf []byte) (int, Buffer) {
	n := copy(buf, s)
	if n == len(s) {
		return n, nil
	}
	return n, s[n:]
}
