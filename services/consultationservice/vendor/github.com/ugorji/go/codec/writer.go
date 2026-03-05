


package codec

import "io"


type encWriter interface {
	writeb([]byte)
	writestr(string)
	writeqstr(string) 
	writen1(byte)

	
	writen2(byte, byte)
	writen4([4]byte)
	writen8([8]byte)

	end()
}



type bufioEncWriter struct {
	w io.Writer

	buf []byte

	n int

	b [16]byte 
}

func (z *bufioEncWriter) reset(w io.Writer, bufsize int, blist *bytesFreelist) {
	z.w = w
	z.n = 0
	if bufsize <= 0 {
		bufsize = defEncByteBufSize
	}
	
	if bufsize <= 8 {
		bufsize = 8
	}
	if cap(z.buf) < bufsize {
		if len(z.buf) > 0 && &z.buf[0] != &z.b[0] {
			blist.put(z.buf)
		}
		if len(z.b) > bufsize {
			z.buf = z.b[:]
		} else {
			z.buf = blist.get(bufsize)
		}
	}
	z.buf = z.buf[:cap(z.buf)]
}

func (z *bufioEncWriter) flushErr() (err error) {
	n, err := z.w.Write(z.buf[:z.n])
	z.n -= n
	if z.n > 0 {
		if err == nil {
			err = io.ErrShortWrite
		}
		if n > 0 {
			copy(z.buf, z.buf[n:z.n+n])
		}
	}
	return err
}

func (z *bufioEncWriter) flush() {
	halt.onerror(z.flushErr())
}

func (z *bufioEncWriter) writeb(s []byte) {
LOOP:
	a := len(z.buf) - z.n
	if len(s) > a {
		z.n += copy(z.buf[z.n:], s[:a])
		s = s[a:]
		z.flush()
		goto LOOP
	}
	z.n += copy(z.buf[z.n:], s)
}

func (z *bufioEncWriter) writestr(s string) {
	
LOOP:
	a := len(z.buf) - z.n
	if len(s) > a {
		z.n += copy(z.buf[z.n:], s[:a])
		s = s[a:]
		z.flush()
		goto LOOP
	}
	z.n += copy(z.buf[z.n:], s)
}

func (z *bufioEncWriter) writeqstr(s string) {
	
	
	

	if z.n+len(s)+2 > len(z.buf) {
		z.flush()
	}
	setByteAt(z.buf, uint(z.n), '"')
	
	z.n++
LOOP:
	a := len(z.buf) - z.n
	if len(s)+1 > a {
		z.n += copy(z.buf[z.n:], s[:a])
		s = s[a:]
		z.flush()
		goto LOOP
	}
	z.n += copy(z.buf[z.n:], s)
	setByteAt(z.buf, uint(z.n), '"')
	
	z.n++
}

func (z *bufioEncWriter) writen1(b1 byte) {
	if 1 > len(z.buf)-z.n {
		z.flush()
	}
	setByteAt(z.buf, uint(z.n), b1)
	
	z.n++
}
func (z *bufioEncWriter) writen2(b1, b2 byte) {
	if 2 > len(z.buf)-z.n {
		z.flush()
	}
	setByteAt(z.buf, uint(z.n+1), b2)
	setByteAt(z.buf, uint(z.n), b1)
	
	
	z.n += 2
}

func (z *bufioEncWriter) writen4(b [4]byte) {
	if 4 > len(z.buf)-z.n {
		z.flush()
	}
	
	
	
	
	copy(z.buf[z.n:], b[:])
	z.n += 4
}

func (z *bufioEncWriter) writen8(b [8]byte) {
	if 8 > len(z.buf)-z.n {
		z.flush()
	}
	copy(z.buf[z.n:], b[:])
	z.n += 8
}

func (z *bufioEncWriter) endErr() (err error) {
	if z.n > 0 {
		err = z.flushErr()
	}
	return
}




type bytesEncAppender struct {
	b   []byte
	out *[]byte
}

func (z *bytesEncAppender) writeb(s []byte) {
	z.b = append(z.b, s...)
}
func (z *bytesEncAppender) writestr(s string) {
	z.b = append(z.b, s...)
}
func (z *bytesEncAppender) writeqstr(s string) {
	z.b = append(append(append(z.b, '"'), s...), '"')
	
	
	
}
func (z *bytesEncAppender) writen1(b1 byte) {
	z.b = append(z.b, b1)
}
func (z *bytesEncAppender) writen2(b1, b2 byte) {
	z.b = append(z.b, b1, b2)
}

func (z *bytesEncAppender) writen4(b [4]byte) {
	z.b = append(z.b, b[:]...)
	
}

func (z *bytesEncAppender) writen8(b [8]byte) {
	z.b = append(z.b, b[:]...)
	
}

func (z *bytesEncAppender) endErr() error {
	*(z.out) = z.b
	return nil
}
func (z *bytesEncAppender) reset(in []byte, out *[]byte) {
	z.b = in[:0]
	z.out = out
}



type encWr struct {
	wb bytesEncAppender
	wf *bufioEncWriter

	bytes bool 

	
	

	js bool 
	be bool 

	c containerState

	calls uint16
	seq   uint16 
}






func (z *encWr) writeb(s []byte) {
	if z.bytes {
		z.wb.writeb(s)
	} else {
		z.wf.writeb(s)
	}
}
func (z *encWr) writestr(s string) {
	if z.bytes {
		z.wb.writestr(s)
	} else {
		z.wf.writestr(s)
	}
}




func (z *encWr) WriteStr(s string) {
	if z.bytes {
		z.wb.writestr(s)
	} else {
		z.wf.writestr(s)
	}
}

func (z *encWr) writen1(b1 byte) {
	if z.bytes {
		z.wb.writen1(b1)
	} else {
		z.wf.writen1(b1)
	}
}

func (z *encWr) writen2(b1, b2 byte) {
	if z.bytes {
		
		z.wb.b = append(z.wb.b, b1, b2)
	} else {
		z.wf.writen2(b1, b2)
	}
}

func (z *encWr) writen4(b [4]byte) {
	if z.bytes {
		
		z.wb.b = append(z.wb.b, b[:]...)
		
	} else {
		z.wf.writen4(b)
	}
}
func (z *encWr) writen8(b [8]byte) {
	if z.bytes {
		
		z.wb.writen8(b)
	} else {
		z.wf.writen8(b)
	}
}

func (z *encWr) writeqstr(s string) {
	if z.bytes {
		
		z.wb.b = append(append(append(z.wb.b, '"'), s...), '"')
	} else {
		z.wf.writeqstr(s)
	}
}

func (z *encWr) endErr() error {
	if z.bytes {
		return z.wb.endErr()
	}
	return z.wf.endErr()
}

func (z *encWr) end() {
	halt.onerror(z.endErr())
}

var _ encWriter = (*encWr)(nil)
