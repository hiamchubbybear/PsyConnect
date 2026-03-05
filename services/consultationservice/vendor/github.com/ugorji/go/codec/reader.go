


package codec

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)



type decReader interface {
	
	
	
	readx(n uint) []byte

	readb([]byte)

	readn1() byte
	readn2() [2]byte
	readn3() [3]byte
	readn4() [4]byte
	readn8() [8]byte
	

	
	

	numread() uint 

	
	skipWhitespace() (token byte)

	
	
	
	jsonReadNum() []byte

	
	
	jsonReadAsisChars() []byte

	
	

	
	

	
	readUntil(stop byte) (out []byte)
}



type unreadByteStatus uint8



const (
	unreadByteUndefined unreadByteStatus = iota
	unreadByteCanRead
	unreadByteCanUnread
)






type ioReaderByteScanner interface {
	io.Reader
	io.ByteScanner
	
	
	
}



type ioReaderByteScannerT struct {
	r io.Reader

	l  byte             
	ls unreadByteStatus 

	_ [2]byte 
	b [4]byte 
}

func (z *ioReaderByteScannerT) ReadByte() (c byte, err error) {
	if z.ls == unreadByteCanRead {
		z.ls = unreadByteCanUnread
		c = z.l
	} else {
		_, err = z.Read(z.b[:1])
		c = z.b[0]
	}
	return
}

func (z *ioReaderByteScannerT) UnreadByte() (err error) {
	switch z.ls {
	case unreadByteCanUnread:
		z.ls = unreadByteCanRead
	case unreadByteCanRead:
		err = errDecUnreadByteLastByteNotRead
	case unreadByteUndefined:
		err = errDecUnreadByteNothingToRead
	default:
		err = errDecUnreadByteUnknown
	}
	return
}

func (z *ioReaderByteScannerT) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return
	}
	var firstByte bool
	if z.ls == unreadByteCanRead {
		z.ls = unreadByteCanUnread
		p[0] = z.l
		if len(p) == 1 {
			n = 1
			return
		}
		firstByte = true
		p = p[1:]
	}
	n, err = z.r.Read(p)
	if n > 0 {
		if err == io.EOF && n == len(p) {
			err = nil 
		}
		z.l = p[n-1]
		z.ls = unreadByteCanUnread
	}
	if firstByte {
		n++
	}
	return
}

func (z *ioReaderByteScannerT) reset(r io.Reader) {
	z.r = r
	z.ls = unreadByteUndefined
	z.l = 0
}


type ioDecReader struct {
	rr ioReaderByteScannerT 

	n uint 

	blist *bytesFreelist

	bufr []byte              
	br   ioReaderByteScanner 
	bb   *bufio.Reader       

	x [64 + 40]byte 
}

func (z *ioDecReader) reset(r io.Reader, bufsize int, blist *bytesFreelist) {
	z.blist = blist
	z.n = 0
	z.bufr = z.blist.check(z.bufr, 256)
	z.br = nil

	var ok bool

	if bufsize <= 0 {
		z.br, ok = r.(ioReaderByteScanner)
		if !ok {
			z.rr.reset(r)
			z.br = &z.rr
		}
		return
	}

	

	
	
	switch bb := r.(type) {
	case *strings.Reader:
		z.br = bb
	case *bytes.Buffer:
		z.br = bb
	case *bytes.Reader:
		z.br = bb
	case *bufio.Reader:
		if bb.Size() == bufsize {
			z.br = bb
		}
	}

	if z.br == nil {
		if z.bb != nil && z.bb.Size() == bufsize {
			z.bb.Reset(r)
		} else {
			z.bb = bufio.NewReaderSize(r, bufsize)
		}
		z.br = z.bb
	}
}

func (z *ioDecReader) numread() uint {
	return z.n
}

func (z *ioDecReader) readn1() (b uint8) {
	b, err := z.br.ReadByte()
	halt.onerror(err)
	z.n++
	return
}

func (z *ioDecReader) readn2() (bs [2]byte) {
	z.readb(bs[:])
	return
}

func (z *ioDecReader) readn3() (bs [3]byte) {
	z.readb(bs[:])
	return
}

func (z *ioDecReader) readn4() (bs [4]byte) {
	z.readb(bs[:])
	return
}

func (z *ioDecReader) readn8() (bs [8]byte) {
	z.readb(bs[:])
	return
}

func (z *ioDecReader) readx(n uint) (bs []byte) {
	if n == 0 {
		return zeroByteSlice
	}
	if n < uint(len(z.x)) {
		bs = z.x[:n]
	} else {
		bs = make([]byte, n)
	}
	nn, err := readFull(z.br, bs)
	z.n += nn
	halt.onerror(err)
	return
}

func (z *ioDecReader) readb(bs []byte) {
	if len(bs) == 0 {
		return
	}
	nn, err := readFull(z.br, bs)
	z.n += nn
	halt.onerror(err)
}













func (z *ioDecReader) jsonReadNum() (bs []byte) {
	z.unreadn1()
	z.bufr = z.bufr[:0]
LOOP:
	
	i, err := z.br.ReadByte()
	if err == io.EOF {
		return z.bufr
	}
	if err != nil {
		halt.onerror(err)
	}
	z.n++
	if isNumberChar(i) {
		z.bufr = append(z.bufr, i)
		goto LOOP
	}
	z.unreadn1()
	return z.bufr
}

func (z *ioDecReader) jsonReadAsisChars() (bs []byte) {
	z.bufr = z.bufr[:0]
LOOP:
	i := z.readn1()
	z.bufr = append(z.bufr, i)
	if i == '"' || i == '\\' {
		return z.bufr
	}
	goto LOOP
}

func (z *ioDecReader) skipWhitespace() (token byte) {
LOOP:
	token = z.readn1()
	if isWhitespaceChar(token) {
		goto LOOP
	}
	return
}












func (z *ioDecReader) readUntil(stop byte) []byte {
	z.bufr = z.bufr[:0]
LOOP:
	token := z.readn1()
	if token == stop {
		return z.bufr
	}
	z.bufr = append(z.bufr, token)
	goto LOOP
}

func (z *ioDecReader) unreadn1() {
	err := z.br.UnreadByte()
	halt.onerror(err)
	z.n--
}













type bytesDecReader struct {
	b []byte 
	c uint   
}

func (z *bytesDecReader) reset(in []byte) {
	z.b = in[:len(in):len(in)] 
	z.c = 0
}

func (z *bytesDecReader) numread() uint {
	return z.c
}





func (z *bytesDecReader) readx(n uint) (bs []byte) {
	
	
	
	bs = z.b[z.c : z.c+n]
	z.c += n
	return
}

func (z *bytesDecReader) readb(bs []byte) {
	copy(bs, z.readx(uint(len(bs))))
}
























func (z *bytesDecReader) readn1() (v uint8) {
	v = z.b[z.c]
	z.c++
	return
}

func (z *bytesDecReader) readn2() (bs [2]byte) {
	
	
	
	bs = okBytes2(z.b[z.c : z.c+2])
	z.c += 2
	return
}

func (z *bytesDecReader) readn3() (bs [3]byte) {
	
	bs = okBytes3(z.b[z.c : z.c+3])
	z.c += 3
	return
}

func (z *bytesDecReader) readn4() (bs [4]byte) {
	
	bs = okBytes4(z.b[z.c : z.c+4])
	z.c += 4
	return
}

func (z *bytesDecReader) readn8() (bs [8]byte) {
	
	bs = okBytes8(z.b[z.c : z.c+8])
	z.c += 8
	return
}

func (z *bytesDecReader) jsonReadNum() []byte {
	z.c-- 
	i := z.c
LOOP:
	
	if i < uint(len(z.b)) && isNumberChar(z.b[i]) {
		i++
		goto LOOP
	}
	z.c, i = i, z.c
	
	
	return z.b[i:z.c]
}

func (z *bytesDecReader) jsonReadAsisChars() []byte {
	i := z.c
LOOP:
	token := z.b[i]
	i++
	if token == '"' || token == '\\' {
		z.c, i = i, z.c
		return byteSliceOf(z.b, i, z.c)
		
	}
	goto LOOP
}

func (z *bytesDecReader) skipWhitespace() (token byte) {
	i := z.c
LOOP:
	token = z.b[i]
	if isWhitespaceChar(token) {
		i++
		goto LOOP
	}
	z.c = i + 1
	return
}

func (z *bytesDecReader) readUntil(stop byte) (out []byte) {
	i := z.c
LOOP:
	if z.b[i] == stop {
		out = byteSliceOf(z.b, z.c, i)
		
		z.c = i + 1
		return
	}
	i++
	goto LOOP
}



type decRd struct {
	rb bytesDecReader
	ri *ioDecReader

	decReader

	bytes bool 

	
	

	mtr bool 
	str bool 

	be   bool 
	js   bool 
	jsms bool 
	cbor bool 

	cbreak bool 

}
















func (z *decRd) numread() uint {
	if z.bytes {
		return z.rb.numread()
	}
	return z.ri.numread()
}

func (z *decRd) readn1() (v uint8) {
	if z.bytes {
		
		
		
		
		v = z.rb.b[z.rb.c]
		z.rb.c++
		return
	}
	return z.ri.readn1()
}






















type devNullReader struct{}

func (devNullReader) Read(p []byte) (int, error) { return 0, io.EOF }
func (devNullReader) Close() error               { return nil }

func readFull(r io.Reader, bs []byte) (n uint, err error) {
	var nn int
	for n < uint(len(bs)) && err == nil {
		nn, err = r.Read(bs[n:])
		if nn > 0 {
			if err == io.EOF {
				
				err = nil
			}
			n += uint(nn)
		}
	}
	
	
	return
}

var _ decReader = (*decRd)(nil)
