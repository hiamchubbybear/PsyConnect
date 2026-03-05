


package codec

import (
	"bufio"
	"errors"
	"io"
	"net/rpc"
)

var (
	errRpcIsClosed = errors.New("rpc - connection has been closed")
	errRpcNoConn   = errors.New("rpc - no connection")

	rpcSpaceArr = [1]byte{' '}
)


type Rpc interface {
	ServerCodec(conn io.ReadWriteCloser, h Handle) rpc.ServerCodec
	ClientCodec(conn io.ReadWriteCloser, h Handle) rpc.ClientCodec
}


type RPCOptions struct {
	
	
	
	
	
	RPCNoBuffer bool
}


type rpcCodec struct {
	c io.Closer
	r io.Reader
	w io.Writer
	f ioFlusher

	dec *Decoder
	enc *Encoder
	h   Handle

	cls atomicClsErr
}

func newRPCCodec(conn io.ReadWriteCloser, h Handle) rpcCodec {
	return newRPCCodec2(conn, conn, conn, h)
}

func newRPCCodec2(r io.Reader, w io.Writer, c io.Closer, h Handle) rpcCodec {
	bh := h.getBasicHandle()
	
	
	
	f, ok := w.(ioFlusher)
	if !bh.RPCNoBuffer {
		if bh.WriterBufferSize <= 0 {
			if !ok { 
				bw := bufio.NewWriter(w)
				f, w = bw, bw
			}
		}
		if bh.ReaderBufferSize <= 0 {
			if _, ok = w.(ioBuffered); !ok {
				r = bufio.NewReader(r)
			}
		}
	}
	return rpcCodec{
		c:   c,
		w:   w,
		r:   r,
		f:   f,
		h:   h,
		enc: NewEncoder(w, h),
		dec: NewDecoder(r, h),
	}
}

func (c *rpcCodec) write(obj ...interface{}) (err error) {
	err = c.ready()
	if err != nil {
		return
	}
	if c.f != nil {
		defer func() {
			flushErr := c.f.Flush()
			if flushErr != nil && err == nil {
				err = flushErr
			}
		}()
	}

	for _, o := range obj {
		err = c.enc.Encode(o)
		if err != nil {
			return
		}
		
		
		
		if c.h.isJson() {
			_, err = c.w.Write(rpcSpaceArr[:])
			if err != nil {
				return
			}
		}
	}
	return
}

func (c *rpcCodec) read(obj interface{}) (err error) {
	err = c.ready()
	if err == nil {
		
		if obj == nil {
			
			err = c.dec.swallowErr()
		} else {
			err = c.dec.Decode(obj)
		}
	}
	return
}

func (c *rpcCodec) Close() (err error) {
	if c.c != nil {
		cls := c.cls.load()
		if !cls.closed {
			cls.err = c.c.Close()
			cls.closed = true
			c.cls.store(cls)
		}
		err = cls.err
	}
	return
}

func (c *rpcCodec) ready() (err error) {
	if c.c == nil {
		err = errRpcNoConn
	} else {
		cls := c.cls.load()
		if cls.closed {
			if err = cls.err; err == nil {
				err = errRpcIsClosed
			}
		}
	}
	return
}

func (c *rpcCodec) ReadResponseBody(body interface{}) error {
	return c.read(body)
}



type goRpcCodec struct {
	rpcCodec
}

func (c *goRpcCodec) WriteRequest(r *rpc.Request, body interface{}) error {
	return c.write(r, body)
}

func (c *goRpcCodec) WriteResponse(r *rpc.Response, body interface{}) error {
	return c.write(r, body)
}

func (c *goRpcCodec) ReadResponseHeader(r *rpc.Response) error {
	return c.read(r)
}

func (c *goRpcCodec) ReadRequestHeader(r *rpc.Request) error {
	return c.read(r)
}

func (c *goRpcCodec) ReadRequestBody(body interface{}) error {
	return c.read(body)
}





type goRpc struct{}


































var GoRpc goRpc

func (x goRpc) ServerCodec(conn io.ReadWriteCloser, h Handle) rpc.ServerCodec {
	return &goRpcCodec{newRPCCodec(conn, h)}
}

func (x goRpc) ClientCodec(conn io.ReadWriteCloser, h Handle) rpc.ClientCodec {
	return &goRpcCodec{newRPCCodec(conn, h)}
}
