

package api

import (
    `reflect`

    `github.com/bytedance/sonic/internal/native`
    `github.com/bytedance/sonic/internal/native/types`
	`github.com/bytedance/sonic/internal/decoder/consts`
	`github.com/bytedance/sonic/internal/decoder/errors`
    `github.com/bytedance/sonic/internal/rt`
    `github.com/bytedance/sonic/option`
)

const (
	_F_allow_control = consts.F_allow_control
	_F_copy_string = consts.F_copy_string
	_F_disable_unknown = consts.F_disable_unknown
	_F_disable_urc = consts.F_disable_urc
	_F_use_int64 = consts.F_use_int64
	_F_use_number = consts.F_use_number
	_F_validate_string = consts.F_validate_string
    _F_case_sensitive = consts.F_case_sensitive

	_MaxStack = consts.MaxStack

	OptionUseInt64 	       = consts.OptionUseInt64
	OptionUseNumber        = consts.OptionUseNumber
    OptionUseUnicodeErrors = consts.OptionUseUnicodeErrors
    OptionDisableUnknown   = consts.OptionDisableUnknown
    OptionCopyString       = consts.OptionCopyString
    OptionValidateString   = consts.OptionValidateString
    OptionNoValidateJSON   = consts.OptionNoValidateJSON
    OptionCaseSensitive    = consts.OptionCaseSensitive
)

type (
	Options = consts.Options
	MismatchTypeError = errors.MismatchTypeError
	SyntaxError = errors.SyntaxError
)

func (self *Decoder) SetOptions(opts Options) {
    if (opts & consts.OptionUseNumber != 0) && (opts & consts.OptionUseInt64 != 0) {
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

    
    return SyntaxError {
        Src  : buf,
        Pos  : pos,
        Code : types.ERR_INVALID_CHAR,
    }
}




func (self *Decoder) Decode(val interface{}) error {
	return decodeImpl(&self.s, &self.i, self.f, val)
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
	return pretouchImpl(vt, opts...)
}



func Skip(data []byte) (start int, end int) {
    s := rt.Mem2Str(data)
    p := 0
    m := types.NewStateMachine()
    ret := native.SkipOne(&s, &p, m, uint64(0))
    types.FreeStateMachine(m) 
    return ret, p
}
