// +build !amd64,!arm64 go1.25 !go1.17 arm64,!go1.20



package ast

import (
    `encoding/json`
    `unicode/utf8`

    `github.com/bytedance/sonic/internal/native/types`
    `github.com/bytedance/sonic/internal/rt`
    `github.com/bytedance/sonic/internal/compat`
)

func init() {
    compat.Warn("sonic/ast")
}

func quote(buf *[]byte, val string) {
    quoteString(buf, val)
}


func unquote(src string) (string, types.ParsingError) {
    sp := rt.IndexChar(src, -1)
    out, ok := unquoteBytes(rt.BytesFrom(sp, len(src)+2, len(src)+2))
    if !ok {
        return "", types.ERR_INVALID_ESCAPE
    }
    return rt.Mem2Str(out), 0
}


func (self *Parser) decodeValue() (val types.JsonState) {
    e, v := decodeValue(self.s, self.p, self.dbuf == nil)
    if e < 0 {
        return v
    }
    self.p = e
    return v
}

func (self *Parser) skip() (int, types.ParsingError) {
    e, s := skipValue(self.s, self.p)
    if e < 0 {
        return self.p, types.ParsingError(-e)
    }
    self.p = e
    return s, 0
}

func (self *Parser) skipFast() (int, types.ParsingError) {
    e, s := skipValueFast(self.s, self.p)
    if e < 0 {
        return self.p, types.ParsingError(-e)
    }
    self.p = e
    return s, 0
}

func (self *Node) encodeInterface(buf *[]byte) error {
    out, err := json.Marshal(self.packAny())
    if err != nil {
        return err
    }
    *buf = append(*buf, out...)
    return nil
}

func (self *Parser) getByPath(validate bool, path ...interface{}) (int, types.ParsingError) {
    for _, p := range path {
        if idx, ok := p.(int); ok && idx >= 0 {
            if err := self.searchIndex(idx); err != 0 {
                return self.p, err
            }
        } else if key, ok := p.(string); ok {
            if err := self.searchKey(key); err != 0 {
                return self.p, err
            }
        } else {
            panic("path must be either int(>=0) or string")
        }
    }

    var start int
    var e types.ParsingError
    if validate {
        start, e = self.skip()
    } else {
        start, e = self.skipFast()
    }
    if e != 0 {
        return self.p, e
    }
    return start, 0
}

func validate_utf8(str string) bool {
    return utf8.ValidString(str)
}
