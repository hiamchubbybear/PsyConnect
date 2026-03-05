package snappy

import (
	"io"
	"sync"

	"github.com/klauspost/compress/s2"
	"github.com/klauspost/compress/snappy"
)



type Framing int

const (
	Framed Framing = iota
	Unframed
)


type Compression int

const (
	DefaultCompression Compression = iota
	FasterCompression
	BetterCompression
	BestCompression
)

var (
	readerPool sync.Pool
	writerPool sync.Pool
)



type Codec struct {
	
	
	
	Framing Framing

	
	Compression Compression
}


func (c *Codec) Code() int8 { return 2 }


func (c *Codec) Name() string { return "snappy" }


func (c *Codec) NewReader(r io.Reader) io.ReadCloser {
	x, _ := readerPool.Get().(*xerialReader)
	if x != nil {
		x.Reset(r)
	} else {
		x = &xerialReader{
			reader: r,
			decode: snappy.Decode,
		}
	}
	return &reader{xerialReader: x}
}


func (c *Codec) NewWriter(w io.Writer) io.WriteCloser {
	x, _ := writerPool.Get().(*xerialWriter)
	if x != nil {
		x.Reset(w)
	} else {
		x = &xerialWriter{writer: w}
	}
	x.framed = c.Framing == Framed
	switch c.Compression {
	case FasterCompression:
		x.encode = s2.EncodeSnappy
	case BetterCompression:
		x.encode = s2.EncodeSnappyBetter
	case BestCompression:
		x.encode = s2.EncodeSnappyBest
	default:
		x.encode = snappy.Encode 
	}
	return &writer{xerialWriter: x}
}

type reader struct{ *xerialReader }

func (r *reader) Close() (err error) {
	if x := r.xerialReader; x != nil {
		r.xerialReader = nil
		x.Reset(nil)
		readerPool.Put(x)
	}
	return
}

type writer struct{ *xerialWriter }

func (w *writer) Close() (err error) {
	if x := w.xerialWriter; x != nil {
		w.xerialWriter = nil
		err = x.Flush()
		x.Reset(nil)
		writerPool.Put(x)
	}
	return
}
