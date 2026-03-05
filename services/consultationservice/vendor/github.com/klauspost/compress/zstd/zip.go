


package zstd

import (
	"errors"
	"io"
	"sync"
)



const ZipMethodWinZip = 93




const ZipMethodPKWare = 20


var zipReaderPool = sync.Pool{New: func() interface{} {
	z, err := NewReader(nil, WithDecoderLowmem(true), WithDecoderMaxWindow(128<<20), WithDecoderConcurrency(1))
	if err != nil {
		panic(err)
	}
	return z
}}


func newZipReader(opts ...DOption) func(r io.Reader) io.ReadCloser {
	pool := &zipReaderPool
	if len(opts) > 0 {
		opts = append([]DOption{WithDecoderLowmem(true), WithDecoderMaxWindow(128 << 20)}, opts...)
		
		opts = append(opts, WithDecoderConcurrency(1))
		
		pool = &sync.Pool{}
	}
	return func(r io.Reader) io.ReadCloser {
		dec, ok := pool.Get().(*Decoder)
		if ok {
			dec.Reset(r)
		} else {
			d, err := NewReader(r, opts...)
			if err != nil {
				panic(err)
			}
			dec = d
		}
		return &pooledZipReader{dec: dec, pool: pool}
	}
}

type pooledZipReader struct {
	mu   sync.Mutex 
	pool *sync.Pool
	dec  *Decoder
}

func (r *pooledZipReader) Read(p []byte) (n int, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.dec == nil {
		return 0, errors.New("read after close or EOF")
	}
	dec, err := r.dec.Read(p)
	if err == io.EOF {
		r.dec.Reset(nil)
		r.pool.Put(r.dec)
		r.dec = nil
	}
	return dec, err
}

func (r *pooledZipReader) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	var err error
	if r.dec != nil {
		err = r.dec.Reset(nil)
		r.pool.Put(r.dec)
		r.dec = nil
	}
	return err
}

type pooledZipWriter struct {
	mu   sync.Mutex 
	enc  *Encoder
	pool *sync.Pool
}

func (w *pooledZipWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.enc == nil {
		return 0, errors.New("Write after Close")
	}
	return w.enc.Write(p)
}

func (w *pooledZipWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	var err error
	if w.enc != nil {
		err = w.enc.Close()
		w.pool.Put(w.enc)
		w.enc = nil
	}
	return err
}



func ZipCompressor(opts ...EOption) func(w io.Writer) (io.WriteCloser, error) {
	var pool sync.Pool
	return func(w io.Writer) (io.WriteCloser, error) {
		enc, ok := pool.Get().(*Encoder)
		if ok {
			enc.Reset(w)
		} else {
			var err error
			enc, err = NewWriter(w, opts...)
			if err != nil {
				return nil, err
			}
		}
		return &pooledZipWriter{enc: enc, pool: &pool}, nil
	}
}






func ZipDecompressor(opts ...DOption) func(r io.Reader) io.ReadCloser {
	return newZipReader(opts...)
}
