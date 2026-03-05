
package zstd

import (
	"io"
	"sync"

	"github.com/klauspost/compress/zstd"
)



type Codec struct {
	
	
	
	Level int

	encoderPool sync.Pool 
}


func (c *Codec) Code() int8 { return 4 }


func (c *Codec) Name() string { return "zstd" }


func (c *Codec) NewReader(r io.Reader) io.ReadCloser {
	p := new(reader)
	if p.dec, _ = decoderPool.Get().(*zstd.Decoder); p.dec != nil {
		p.dec.Reset(r)
	} else {
		z, err := zstd.NewReader(r,
			zstd.WithDecoderConcurrency(1),
		)
		if err != nil {
			p.err = err
		} else {
			p.dec = z
		}
	}
	return p
}

func (c *Codec) level() int {
	if c.Level != 0 {
		return c.Level
	}
	return 3
}

func (c *Codec) zstdLevel() zstd.EncoderLevel {
	return zstd.EncoderLevelFromZstd(c.level())
}

var decoderPool sync.Pool 

type reader struct {
	dec *zstd.Decoder
	err error
}


func (r *reader) Close() error {
	if r.dec != nil {
		r.dec.Reset(devNull{}) 
		decoderPool.Put(r.dec)
		r.dec = nil
		r.err = io.ErrClosedPipe
	}
	return nil
}


func (r *reader) Read(p []byte) (int, error) {
	if r.err != nil {
		return 0, r.err
	}
	if r.dec == nil {
		return 0, io.EOF
	}
	return r.dec.Read(p)
}


func (r *reader) WriteTo(w io.Writer) (int64, error) {
	if r.err != nil {
		return 0, r.err
	}
	if r.dec == nil {
		return 0, io.ErrClosedPipe
	}
	return r.dec.WriteTo(w)
}


func (c *Codec) NewWriter(w io.Writer) io.WriteCloser {
	p := new(writer)
	if enc, _ := c.encoderPool.Get().(*zstd.Encoder); enc == nil {
		z, err := zstd.NewWriter(w,
			zstd.WithEncoderLevel(c.zstdLevel()),
			zstd.WithEncoderConcurrency(1),
			zstd.WithZeroFrames(true),
		)
		if err != nil {
			p.err = err
		} else {
			p.enc = z
		}
	} else {
		p.enc = enc
		p.enc.Reset(w)
	}
	p.c = c
	return p
}

type writer struct {
	c   *Codec
	enc *zstd.Encoder
	err error
}


func (w *writer) Close() error {
	if w.enc != nil {
		
		
		
		err := w.enc.Close()
		if err != nil {
			w.err = err
		}
		w.enc.Reset(devNull{}) 
		w.c.encoderPool.Put(w.enc)
		w.enc = nil
		return err
	}
	return w.err
}


func (w *writer) Write(p []byte) (int, error) {
	if w.err != nil {
		return 0, w.err
	}
	if w.enc == nil {
		return 0, io.ErrClosedPipe
	}
	return w.enc.Write(p)
}


func (w *writer) ReadFrom(r io.Reader) (int64, error) {
	if w.err != nil {
		return 0, w.err
	}
	if w.enc == nil {
		return 0, io.ErrClosedPipe
	}
	return w.enc.ReadFrom(r)
}

type devNull struct{}

func (devNull) Read([]byte) (int, error)  { return 0, io.EOF }
func (devNull) Write([]byte) (int, error) { return 0, nil }
