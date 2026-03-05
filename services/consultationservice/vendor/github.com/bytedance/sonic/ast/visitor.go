

package ast

import (
    `encoding/json`
    `errors`

    `github.com/bytedance/sonic/internal/native/types`
)












type Visitor interface {

    
    OnNull() error

    
    OnBool(v bool) error

    
    OnString(v string) error

    
    OnInt64(v int64, n json.Number) error

    
    OnFloat64(v float64, n json.Number) error

    
    
    
    
    
    
    
    
    
    
    
    
    
    OnObjectBegin(capacity int) error

    
    OnObjectKey(key string) error

    
    OnObjectEnd() error

    
    
    
    
    
    
    
    
    
    
    
    
    
    OnArrayBegin(capacity int) error

    
    OnArrayEnd() error
}



type VisitorOptions struct {
    
    
    
    OnlyNumber bool
}

var defaultVisitorOptions = &VisitorOptions{}





func Preorder(str string, visitor Visitor, opts *VisitorOptions) error {
    if opts == nil {
        opts = defaultVisitorOptions
    }
    
    
    var (
        optDecodeNumber = !opts.OnlyNumber
    )

    tv := &traverser{
        parser: Parser{
            s:         str,
            noLazy:    true,
            skipValue: false,
        },
        visitor: visitor,
    }

    if optDecodeNumber {
        tv.parser.decodeNumber(true)
    }

    err := tv.decodeValue()

    if optDecodeNumber {
        tv.parser.decodeNumber(false)
    }
    return err
}

type traverser struct {
    parser  Parser
    visitor Visitor
}


func (self *traverser) decodeValue() error {
    switch val := self.parser.decodeValue(); val.Vt {
    case types.V_EOF:
        return types.ERR_EOF
    case types.V_NULL:
        return self.visitor.OnNull()
    case types.V_TRUE:
        return self.visitor.OnBool(true)
    case types.V_FALSE:
        return self.visitor.OnBool(false)
    case types.V_STRING:
        return self.decodeString(val.Iv, val.Ep)
    case types.V_DOUBLE:
        return self.visitor.OnFloat64(val.Dv,
            json.Number(self.parser.s[val.Ep:self.parser.p]))
    case types.V_INTEGER:
        return self.visitor.OnInt64(val.Iv,
            json.Number(self.parser.s[val.Ep:self.parser.p]))
    case types.V_ARRAY:
        return self.decodeArray()
    case types.V_OBJECT:
        return self.decodeObject()
    default:
        return types.ParsingError(-val.Vt)
    }
}


func (self *traverser) decodeArray() error {
    sp := self.parser.p
    ns := len(self.parser.s)

    
    if err := self.visitor.OnArrayBegin(_DEFAULT_NODE_CAP); err != nil {
        if err == VisitOPSkip {
            
            self.parser.p -= 1
            if _, e := self.parser.skipFast(); e != 0 {
                return e
            }
            return self.visitor.OnArrayEnd()
        }
        return err
    }

    
    self.parser.p = self.parser.lspace(sp)
    if self.parser.p >= ns {
        return types.ERR_EOF
    }

    
    if self.parser.s[self.parser.p] == ']' {
        self.parser.p++
        return self.visitor.OnArrayEnd()
    }

    for {
        
        if err := self.decodeValue(); err != nil {
            return err
        }
        self.parser.p = self.parser.lspace(self.parser.p)

        
        if self.parser.p >= ns {
            return types.ERR_EOF
        }

        
        switch self.parser.s[self.parser.p] {
        case ',':
            self.parser.p++
        case ']':
            self.parser.p++
            return self.visitor.OnArrayEnd()
        default:
            return types.ERR_INVALID_CHAR
        }
    }
}


func (self *traverser) decodeObject() error {
    sp := self.parser.p
    ns := len(self.parser.s)

    
    if err := self.visitor.OnObjectBegin(_DEFAULT_NODE_CAP); err != nil {
        if err == VisitOPSkip {
            
            self.parser.p -= 1
            if _, e := self.parser.skipFast(); e != 0 {
                return e
            }
            return self.visitor.OnObjectEnd()
        }
        return err
    }

    
    self.parser.p = self.parser.lspace(sp)
    if self.parser.p >= ns {
        return types.ERR_EOF
    }

    
    if self.parser.s[self.parser.p] == '}' {
        self.parser.p++
        return self.visitor.OnObjectEnd()
    }

    for {
        var njs types.JsonState
        var err types.ParsingError

        
        if njs = self.parser.decodeValue(); njs.Vt != types.V_STRING {
            return types.ERR_INVALID_CHAR
        }

        
        idx := self.parser.p - 1
        key := self.parser.s[njs.Iv:idx]

        
        if njs.Ep != -1 {
            if key, err = unquote(key); err != 0 {
                return err
            }
        }

        if err := self.visitor.OnObjectKey(key); err != nil {
            return err
        }

        
        if err = self.parser.delim(); err != 0 {
            return err
        }

        
        if err := self.decodeValue(); err != nil {
            return err
        }

        self.parser.p = self.parser.lspace(self.parser.p)

        
        if self.parser.p >= ns {
            return types.ERR_EOF
        }

        
        switch self.parser.s[self.parser.p] {
        case ',':
            self.parser.p++
        case '}':
            self.parser.p++
            return self.visitor.OnObjectEnd()
        default:
            return types.ERR_INVALID_CHAR
        }
    }
}


func (self *traverser) decodeString(iv int64, ep int) error {
    p := self.parser.p - 1
    s := self.parser.s[iv:p]

    
    if ep == -1 {
        return self.visitor.OnString(s)
    }

    
    out, err := unquote(s)
    if err != 0 {
        return err
    }
    return self.visitor.OnString(out)
}



var VisitOPSkip = errors.New("")
