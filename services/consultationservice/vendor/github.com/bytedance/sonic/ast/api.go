//go:build (amd64 && go1.17 && !go1.25) || (arm64 && go1.20 && !go1.25)
// +build amd64,go1.17,!go1.25 arm64,go1.20,!go1.25



package ast

import (
    `runtime`
    `unsafe`

    `github.com/bytedance/sonic/encoder`
    `github.com/bytedance/sonic/internal/native`
    `github.com/bytedance/sonic/internal/native/types`
    `github.com/bytedance/sonic/internal/rt`
    uq `github.com/bytedance/sonic/unquote`
    `github.com/bytedance/sonic/utf8`
)

var typeByte = rt.UnpackEface(byte(0)).Type

//go:nocheckptr
func quote(buf *[]byte, val string) {
    *buf = append(*buf, '"')
    if len(val) == 0 {
        *buf = append(*buf, '"')
        return
    }

    sp := rt.IndexChar(val, 0)
    nb := len(val)
    b := (*rt.GoSlice)(unsafe.Pointer(buf))

    
    for nb > 0 {
        
        dp := unsafe.Pointer(uintptr(b.Ptr) + uintptr(b.Len))
        dn := b.Cap - b.Len
        
        ret := native.Quote(sp, nb, dp, &dn, 0)
        
        b.Len += dn

        
        if ret >= 0 {
            break
        }

        
        *b = rt.GrowSlice(typeByte, *b, b.Cap*2)
        
        ret = ^ret
        
        nb -= ret
        sp = unsafe.Pointer(uintptr(sp) + uintptr(ret))
    }

    runtime.KeepAlive(buf)
    runtime.KeepAlive(sp)
    *buf = append(*buf, '"')
}

func unquote(src string) (string, types.ParsingError) {
    return uq.String(src)
}

func (self *Parser) decodeValue() (val types.JsonState) {
    sv := (*rt.GoString)(unsafe.Pointer(&self.s))
    flag := types.F_USE_NUMBER
    if self.dbuf != nil {
        flag = 0
        val.Dbuf = self.dbuf
        val.Dcap = types.MaxDigitNums
    }
    self.p = native.Value(sv.Ptr, sv.Len, self.p, &val, uint64(flag))
    return
}

func (self *Parser) skip() (int, types.ParsingError) {
    fsm := types.NewStateMachine()
    start := native.SkipOne(&self.s, &self.p, fsm, 0)
    types.FreeStateMachine(fsm)

    if start < 0 {
        return self.p, types.ParsingError(-start)
    }
    return start, 0
}

func (self *Node) encodeInterface(buf *[]byte) error {
    
    return encoder.EncodeInto(buf, self.packAny(), encoder.NoEncoderNewline)
}

func (self *Parser) skipFast() (int, types.ParsingError) {
    start := native.SkipOneFast(&self.s, &self.p)
    if start < 0 {
        return self.p, types.ParsingError(-start)
    }
    return start, 0
}

func (self *Parser) getByPath(validate bool, path ...interface{}) (int, types.ParsingError) {
    var fsm *types.StateMachine
    if validate {
        fsm = types.NewStateMachine()
    }
    start := native.GetByPath(&self.s, &self.p, &path, fsm)
    if validate {
        types.FreeStateMachine(fsm)
    }
    runtime.KeepAlive(path)
    if start < 0 {
        return self.p, types.ParsingError(-start)
    }
    return start, 0
}

func validate_utf8(str string) bool {
    return utf8.ValidateString(str)
}
