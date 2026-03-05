//go:build (!amd64 && !arm64) || go1.25 || !go1.17 || (arm64 && !go1.20)
// +build !amd64,!arm64 go1.25 !go1.17 arm64,!go1.20



package decoder

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"unsafe"

	"github.com/bytedance/sonic/internal/decoder/consts"
	"github.com/bytedance/sonic/internal/native/types"
	"github.com/bytedance/sonic/option"
	"github.com/bytedance/sonic/internal/compat"
)

func init() {
     compat.Warn("sonic/decoder")
}

const (
     _F_use_int64       = consts.F_use_int64
     _F_disable_urc     = consts.F_disable_unknown
     _F_disable_unknown = consts.F_disable_unknown
     _F_copy_string     = consts.F_copy_string
 
     _F_use_number      = consts.F_use_number
     _F_validate_string = consts.F_validate_string
     _F_allow_control   = consts.F_allow_control
     _F_no_validate_json = consts.F_no_validate_json
     _F_case_sensitive  = consts.F_case_sensitive
)

type Options uint64

const (
     OptionUseInt64         Options = 1 << _F_use_int64
     OptionUseNumber        Options = 1 << _F_use_number
     OptionUseUnicodeErrors Options = 1 << _F_disable_urc
     OptionDisableUnknown   Options = 1 << _F_disable_unknown
     OptionCopyString       Options = 1 << _F_copy_string
     OptionValidateString   Options = 1 << _F_validate_string
     OptionNoValidateJSON   Options = 1 << _F_no_validate_json
     OptionCaseSensitive    Options = 1 << _F_case_sensitive
)

func (self *Decoder) SetOptions(opts Options) {
     if (opts & OptionUseNumber != 0) && (opts & OptionUseInt64 != 0) {
         panic("can't set OptionUseInt64 and OptionUseNumber both!")
     }
     self.f = uint64(opts)
}



type Decoder struct {
     i int
     f uint64
     s string
}


func NewDecoder(s string) *Decoder {
     return &Decoder{s: s}
}


func (self *Decoder) Pos() int {
     return self.i
}

func (self *Decoder) Reset(s string) {
     self.s = s
     self.i = 0
     
}


func (self *Decoder) CheckTrailings() error {
     pos := self.i
     buf := self.s
     
     if pos != len(buf) {
         for pos < len(buf) && (types.SPACE_MASK & (1 << buf[pos])) != 0 {
             pos++
         }
     }

     
     if pos == len(buf) {
         return nil
     }

     
     return nil
}




func (self *Decoder) Decode(val interface{}) error {
    r := bytes.NewBufferString(self.s)
   dec := json.NewDecoder(r)
   if (self.f & uint64(OptionUseNumber)) != 0  {
       dec.UseNumber()
   }
   if (self.f & uint64(OptionDisableUnknown)) != 0  {
       dec.DisallowUnknownFields()
   }
   return dec.Decode(val)
}



func (self *Decoder) UseInt64() {
     self.f  |= 1 << _F_use_int64
     self.f &^= 1 << _F_use_number
}



func (self *Decoder) UseNumber() {
     self.f &^= 1 << _F_use_int64
     self.f  |= 1 << _F_use_number
}



func (self *Decoder) UseUnicodeErrors() {
     self.f |= 1 << _F_disable_urc
}




func (self *Decoder) DisallowUnknownFields() {
     self.f |= 1 << _F_disable_unknown
}


func (self *Decoder) CopyString() {
     self.f |= 1 << _F_copy_string
}




func (self *Decoder) ValidateString() {
     self.f |= 1 << _F_validate_string
}






func Pretouch(vt reflect.Type, opts ...option.CompileOption) error {
     return nil
}

type StreamDecoder = json.Decoder




func NewStreamDecoder(r io.Reader) *StreamDecoder {
   return json.NewDecoder(r)
}


type SyntaxError json.SyntaxError


func (s SyntaxError) Description() string {
     return (*json.SyntaxError)(unsafe.Pointer(&s)).Error()
}

func (s SyntaxError) Error() string {
     return (*json.SyntaxError)(unsafe.Pointer(&s)).Error()
}


type MismatchTypeError json.UnmarshalTypeError
