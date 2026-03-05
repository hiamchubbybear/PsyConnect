package gzip

import (
	"io"
	"sync"

	"github.com/klauspost/compress/gzip"
)

var (
	readerPool sync.Pool
)



type Codec struct {
	
	
	
	
	Level int

	writerPool sync.Pool
}


func (c *Codec) Code() int8 { return 1 }


func (c *Codec) Name() string { return "gzip" }


func (c *Codec) NewReader(r io.Reader) io.ReadCloser {
	var err error
	z, _ := readerPool.Get().(*gzip.Reader)
	if z != nil {
		err = z.Reset(r)
	} else {
		z, err = gzip.NewReader(r)
	}
	if err != nil {
		if z != nil {
			readerPool.Put(z)
		}
		return &errorReader{err: err}
	}
	return &reader{Reader: z}
}


func (c *Codec) NewWriter(w io.Writer) io.WriteCloser {
	x := c.writerPool.Get()
	z, _ := x.(*gzip.Writer)
	if z == nil {
		x, err := gzip.NewWriterLevel(w, c.level())
		if err != nil {
			return &errorWriter{err: err}
		}
		z = x
	} else {
		z.Reset(w)
	}
	return &writer{codec: c, Writer: z}
}

func (c *Codec) level() int {
	if c.Level != 0 {
		return c.Level
	}
	return gzip.DefaultCompression
}

type reader struct{ *gzip.Reader }

func (r *reader) Close() (err error) {
	if z := r.Reader; z != nil {
		r.Reader = nil
		err = z.Close()
		
		
		
		
		
		
		
		z.Reset(emptyReader{})
		readerPool.Put(z)
	}
	return
}

type writer struct {
	codec *Codec
	*gzip.Writer
}

func (w *writer) Close() (err error) {
	if z := w.Writer; z != nil {
		w.Writer = nil
		err = z.Close()
		z.Reset(nil)
		w.codec.writerPool.Put(z)
	}
	return
}

type emptyReader struct{}

func (emptyReader) ReadByte() (byte, error) { return 0, io.EOF }

func (emptyReader) Read([]byte) (int, error) { return 0, io.EOF }

type errorReader struct{ err error }

func (r *errorReader) Close() error { return r.err }

func (r *errorReader) Read([]byte) (int, error) { return 0, r.err }

type errorWriter struct{ err error }

func (w *errorWriter) Close() error { return w.err }

func (w *errorWriter) Write([]byte) (int, error) { return 0, w.err }
