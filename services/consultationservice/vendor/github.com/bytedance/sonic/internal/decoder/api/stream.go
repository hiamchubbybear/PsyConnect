

package api

import (
    `bytes`
    `io`
    `sync`

    `github.com/bytedance/sonic/internal/native`
    `github.com/bytedance/sonic/internal/native/types`
    `github.com/bytedance/sonic/internal/rt`
    `github.com/bytedance/sonic/option`
)

var (
    minLeftBufferShift uint = 1
)


type StreamDecoder struct {
    r       io.Reader
    buf     []byte
    scanp   int
    scanned int64
    err     error
    Decoder
}

var bufPool = sync.Pool{
    New: func () interface{} {
        return make([]byte, 0, option.DefaultDecoderBufferSize)
    },
}

func freeBytes(buf []byte) {
    if rt.CanSizeResue(cap(buf)) {
        bufPool.Put(buf[:0])
    }
}




func NewStreamDecoder(r io.Reader) *StreamDecoder {
    return &StreamDecoder{r : r}
}





func (self *StreamDecoder) Decode(val interface{}) (err error) {
    
    if self.More() {
        var s = self.scanp
    try_skip:
        var e = len(self.buf)
        var src = rt.Mem2Str(self.buf[s:e])
        
        var x = 0;
        if y := native.SkipOneFast(&src, &x); y < 0 {
            if self.readMore()  {
                goto try_skip
            } else {
                err = SyntaxError{e, self.s, types.ParsingError(-s), ""}
                self.setErr(err)
                return
            }
        } else {
            s = y + s
            e = x + s
        }
        
        
        self.Decoder.Reset(string(self.buf[s:e]))
        err = self.Decoder.Decode(val)
        if err != nil {
            self.setErr(err)
            return 
        }

        self.scanp = e
        _, empty := self.scan()
        if empty {
            
            mem := self.buf
            self.buf = nil
            freeBytes(mem)
        } else {
            
            n := copy(self.buf, self.buf[self.scanp:])
            self.buf = self.buf[:n]
        }   

        self.scanned += int64(self.scanp)
        self.scanp = 0
    }    

    return self.err
}



func (self *StreamDecoder) InputOffset() int64 {
    return self.scanned + int64(self.scanp)
}



func (self *StreamDecoder) Buffered() io.Reader {
    return bytes.NewReader(self.buf[self.scanp:])
}



func (self *StreamDecoder) More() bool {
    if self.err != nil {
        return false
    }
    c, err := self.peek()
    return err == nil && c != ']' && c != '}'
}



func (self *StreamDecoder) readMore() bool {
    if self.err != nil {
        return false
    }

    var err error
    var n int
    for {
        
        l := len(self.buf)
        realloc(&self.buf)

        n, err = self.r.Read(self.buf[l:cap(self.buf)])
        self.buf = self.buf[: l+n]

        self.scanp = l
        _, empty := self.scan()
        if !empty {
            return true
        }

        
        if err != nil  {
            self.setErr(err)
            return false
        }
    }
}

func (self *StreamDecoder) setErr(err error) {
    self.err = err
    mem := self.buf[:0]
    self.buf = nil
    freeBytes(mem)
}

func (self *StreamDecoder) peek() (byte, error) {
    var err error
    for {
        c, empty := self.scan()
        if !empty {
            return byte(c), nil
        }
        
        if err != nil {
            self.setErr(err)
            return 0, err
        }
        err = self.refill()
    }
}

func (self *StreamDecoder) scan() (byte, bool) {
    for i := self.scanp; i < len(self.buf); i++ {
        c := self.buf[i]
        if isSpace(c) {
            continue
        }
        self.scanp = i
        return c, false
    }
    return 0, true
}

func isSpace(c byte) bool {
    return types.SPACE_MASK & (1 << c) != 0
}

func (self *StreamDecoder) refill() error {
    
    
    if self.scanp > 0 {
        self.scanned += int64(self.scanp)
        n := copy(self.buf, self.buf[self.scanp:])
        self.buf = self.buf[:n]
        self.scanp = 0
    }

    
    realloc(&self.buf)

    
    n, err := self.r.Read(self.buf[len(self.buf):cap(self.buf)])
    self.buf = self.buf[0 : len(self.buf)+n]

    return err
}

func realloc(buf *[]byte) bool {
    l := uint(len(*buf))
    c := uint(cap(*buf))
    if c == 0 {
       *buf = bufPool.Get().([]byte)
       return true
    }
    if c - l <= c >> minLeftBufferShift {
        e := l+(l>>minLeftBufferShift)
        if e <= c {
            e = c*2
        }
        tmp := make([]byte, l, e)
        copy(tmp, *buf)
        *buf = tmp
        return true
    }
    return false
}

