// +build !amd64,!arm64 go1.25 !go1.17 arm64,!go1.20



package encoder

import (
    `io`
    `bytes`
    `encoding/json`
    `reflect`

    `github.com/bytedance/sonic/option`
    `github.com/bytedance/sonic/internal/compat`
)

func init() {
    compat.Warn("sonic/encoder")
}


const EnableFallback = true


type Options uint64

const (
    bitSortMapKeys          = iota
    bitEscapeHTML          
    bitCompactMarshaler
    bitNoQuoteTextMarshaler
    bitNoNullSliceOrMap
    bitValidateString
    bitNoValidateJSONMarshaler
    bitNoEncoderNewline

    
    bitPointerValue = 63
)

const (
    
    
    
    SortMapKeys          Options = 1 << bitSortMapKeys

    
    
    
    EscapeHTML           Options = 1 << bitEscapeHTML

    
    
    CompactMarshaler     Options = 1 << bitCompactMarshaler

    
    
    NoQuoteTextMarshaler Options = 1 << bitNoQuoteTextMarshaler

    
    
    NoNullSliceOrMap     Options = 1 << bitNoNullSliceOrMap

    
    
    ValidateString       Options = 1 << bitValidateString

    
    
    NoValidateJSONMarshaler Options = 1 << bitNoValidateJSONMarshaler

    
    NoEncoderNewline Options = 1 << bitNoEncoderNewline
  
    
    CompatibleWithStd Options = SortMapKeys | EscapeHTML | CompactMarshaler
)


type Encoder struct {
    Opts Options
    prefix string
    indent string
}


func (self *Encoder) Encode(v interface{}) ([]byte, error) {
    if self.indent != "" || self.prefix != "" { 
        return EncodeIndented(v, self.prefix, self.indent, self.Opts)
    }
    return Encode(v, self.Opts)
}


func (self *Encoder) SortKeys() *Encoder {
    self.Opts |= SortMapKeys
    return self
}


func (self *Encoder) SetEscapeHTML(f bool) {
    if f {
        self.Opts |= EscapeHTML
    } else {
        self.Opts &= ^EscapeHTML
    }
}


func (self *Encoder) SetValidateString(f bool) {
    if f {
        self.Opts |= ValidateString
    } else {
        self.Opts &= ^ValidateString
    }
}


func (self *Encoder) SetNoValidateJSONMarshaler(f bool) {
    if f {
        self.Opts |= NoValidateJSONMarshaler
    } else {
        self.Opts &= ^NoValidateJSONMarshaler
    }
}


func (self *Encoder) SetNoEncoderNewline(f bool) {
    if f {
        self.Opts |= NoEncoderNewline
    } else {
        self.Opts &= ^NoEncoderNewline
    }
}


func (self *Encoder) SetCompactMarshaler(f bool) {
    if f {
        self.Opts |= CompactMarshaler
    } else {
        self.Opts &= ^CompactMarshaler
    }
}


func (self *Encoder) SetNoQuoteTextMarshaler(f bool) {
    if f {
        self.Opts |= NoQuoteTextMarshaler
    } else {
        self.Opts &= ^NoQuoteTextMarshaler
    }
}




func (enc *Encoder) SetIndent(prefix, indent string) {
    enc.prefix = prefix
    enc.indent = indent
}


func Quote(s string) string {
    
    if s == "" {
        return `""`
    }

    out, _ := json.Marshal(s)
    return string(out)
}


func Encode(val interface{}, opts Options) ([]byte, error) {
   return json.Marshal(val)
}



func EncodeInto(buf *[]byte, val interface{}, opts Options) error {
    if buf == nil {
        panic("user-supplied buffer buf is nil")
    }
    w := bytes.NewBuffer(*buf)
    enc := json.NewEncoder(w)
    enc.SetEscapeHTML((opts & EscapeHTML) != 0)
    err := enc.Encode(val)
    *buf = w.Bytes()
    l := len(*buf)
    if l > 0 && (opts & NoEncoderNewline != 0) && (*buf)[l-1] == '\n' {
        *buf = (*buf)[:l-1]
    }
    return err
}







func HTMLEscape(dst []byte, src []byte) []byte {
   d := bytes.NewBuffer(dst)
   json.HTMLEscape(d, src)
   return d.Bytes()
}




func EncodeIndented(val interface{}, prefix string, indent string, opts Options) ([]byte, error) {
   w := bytes.NewBuffer([]byte{})
   enc := json.NewEncoder(w)
   enc.SetEscapeHTML((opts & EscapeHTML) != 0)
   enc.SetIndent(prefix, indent)
   err := enc.Encode(val)
   out := w.Bytes()
   return out, err
}






func Pretouch(vt reflect.Type, opts ...option.CompileOption) error {
   return nil
}






func Valid(data []byte) (ok bool, start int) {
   return json.Valid(data), 0
}


type StreamEncoder = json.Encoder




func NewStreamEncoder(w io.Writer) *StreamEncoder {
   return json.NewEncoder(w)
}

