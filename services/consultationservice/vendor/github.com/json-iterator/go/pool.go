package jsoniter

import (
	"io"
)


type IteratorPool interface {
	BorrowIterator(data []byte) *Iterator
	ReturnIterator(iter *Iterator)
}


type StreamPool interface {
	BorrowStream(writer io.Writer) *Stream
	ReturnStream(stream *Stream)
}

func (cfg *frozenConfig) BorrowStream(writer io.Writer) *Stream {
	stream := cfg.streamPool.Get().(*Stream)
	stream.Reset(writer)
	return stream
}

func (cfg *frozenConfig) ReturnStream(stream *Stream) {
	stream.out = nil
	stream.Error = nil
	stream.Attachment = nil
	cfg.streamPool.Put(stream)
}

func (cfg *frozenConfig) BorrowIterator(data []byte) *Iterator {
	iter := cfg.iteratorPool.Get().(*Iterator)
	iter.ResetBytes(data)
	return iter
}

func (cfg *frozenConfig) ReturnIterator(iter *Iterator) {
	iter.Error = nil
	iter.Attachment = nil
	cfg.iteratorPool.Put(iter)
}
