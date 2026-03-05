

package encoder

import (
	"bytes"
	"encoding/json"
	"reflect"
	"runtime"
	"unsafe"

	"github.com/bytedance/sonic/utf8"
	"github.com/bytedance/sonic/internal/encoder/alg"
	"github.com/bytedance/sonic/internal/encoder/vars"
	"github.com/bytedance/sonic/internal/rt"
	"github.com/bytedance/sonic/option"
)


type Options uint64

const (
    
    
    
    SortMapKeys          Options = 1 << alg.BitSortMapKeys

    
    
    
    EscapeHTML           Options = 1 << alg.BitEscapeHTML

    
    
    CompactMarshaler     Options = 1 << alg.BitCompactMarshaler

    
    
    NoQuoteTextMarshaler Options = 1 << alg.BitNoQuoteTextMarshaler

    
    
    
    NoNullSliceOrMap     Options = 1 << alg.BitNoNullSliceOrMap

    
    
    ValidateString       Options = 1 << alg.BitValidateString

    
    
    NoValidateJSONMarshaler Options = 1 << alg.BitNoValidateJSONMarshaler

    
    NoEncoderNewline Options = 1 << alg.BitNoEncoderNewline
  
    
    CompatibleWithStd Options = SortMapKeys | EscapeHTML | CompactMarshaler

    
    EncodeNullForInfOrNan Options = 1 << alg.BitEncodeNullForInfOrNan
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
    buf := make([]byte, 0, len(s)+2)
    buf = alg.Quote(buf, s, false)
    return rt.Mem2Str(buf)
}


func Encode(val interface{}, opts Options) ([]byte, error) {
    var ret []byte

    buf := vars.NewBytes()
    err := encodeIntoCheckRace(buf, val, opts)

    
    if err != nil {
        vars.FreeBytes(buf)
        return nil, err
    }

    
    old := buf
    *buf = encodeFinish(*old, opts)
    pbuf := ((*rt.GoSlice)(unsafe.Pointer(buf))).Ptr
    pold := ((*rt.GoSlice)(unsafe.Pointer(old))).Ptr

    
    if pbuf != pold {
        vars.FreeBytes(old)
        return *buf, nil
    }

    
    if rt.CanSizeResue(cap(*buf)) {
        ret = make([]byte, len(*buf))
        copy(ret, *buf)
        vars.FreeBytes(buf)
    } else {
        ret = *buf
    }
    
    
    return ret, nil
}



func EncodeInto(buf *[]byte, val interface{}, opts Options) error {
    err := encodeIntoCheckRace(buf, val, opts)
    if err != nil {
        return err
    }
    *buf = encodeFinish(*buf, opts)
    return err
}

func encodeInto(buf *[]byte, val interface{}, opts Options) error {
    stk := vars.NewStack()
    efv := rt.UnpackEface(val)
    err := encodeTypedPointer(buf, efv.Type, &efv.Value, stk, uint64(opts))

    
    if err != nil {
        vars.ResetStack(stk)
    }
    vars.FreeStack(stk)

    
    runtime.KeepAlive(buf)
    runtime.KeepAlive(efv)
    return err
}

func encodeFinish(buf []byte, opts Options) []byte {
    if opts & EscapeHTML != 0 {
        buf = HTMLEscape(nil, buf)
    }
    if (opts & ValidateString != 0) && !utf8.Validate(buf) {
        buf = utf8.CorrectWith(nil, buf, `\ufffd`)
    }
    return buf
}








func HTMLEscape(dst []byte, src []byte) []byte {
    return alg.HtmlEscape(dst, src)
}




func EncodeIndented(val interface{}, prefix string, indent string, opts Options) ([]byte, error) {
    var err error
    var buf *bytes.Buffer

    
    out := vars.NewBytes()
    err = EncodeInto(out, val, opts)

    
    if err != nil {
        vars.FreeBytes(out)
        return nil, err
    }

    
    buf = vars.NewBuffer()
    err = json.Indent(buf, *out, prefix, indent)
    vars.FreeBytes(out)

    
    if err != nil {
        vars.FreeBuffer(buf)
        return nil, err
    }

    
    var ret []byte
    if rt.CanSizeResue(cap(buf.Bytes())) {
        ret = make([]byte, buf.Len())
        copy(ret, buf.Bytes())
        
        vars.FreeBuffer(buf)
    } else {
        ret = buf.Bytes()
    }
    
    return ret, nil
}






func Pretouch(vt reflect.Type, opts ...option.CompileOption) error {
    cfg := option.DefaultCompileOptions()
    for _, opt := range opts {
        opt(&cfg)
    }
    return pretouchRec(map[reflect.Type]uint8{vt: 0}, cfg)
}






func Valid(data []byte) (ok bool, start int) {
    return alg.Valid(data)
}
