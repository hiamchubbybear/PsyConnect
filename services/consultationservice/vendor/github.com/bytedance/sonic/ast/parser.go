

package ast

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/bytedance/sonic/internal/native/types"
	"github.com/bytedance/sonic/internal/rt"
)

const (
    _DEFAULT_NODE_CAP int = 16
    _APPEND_GROW_SHIFT = 1
)

const (
    _ERR_NOT_FOUND      types.ParsingError = 33
    _ERR_UNSUPPORT_TYPE types.ParsingError = 34
)

var (
    
    ErrNotExist error = newError(_ERR_NOT_FOUND, "value not exists")

    
    ErrUnsupportType error = newError(_ERR_UNSUPPORT_TYPE, "unsupported type")
)

type Parser struct {
    p           int
    s           string
    noLazy      bool
    loadOnce  bool
    skipValue   bool
    dbuf        *byte
}



func (self *Parser) delim() types.ParsingError {
    n := len(self.s)
    p := self.lspace(self.p)

    
    if p >= n {
        return types.ERR_EOF
    }

    
    if self.s[p] != ':' {
        return types.ERR_INVALID_CHAR
    }

    
    self.p = p + 1
    return 0
}

func (self *Parser) object() types.ParsingError {
    n := len(self.s)
    p := self.lspace(self.p)

    
    if p >= n {
        return types.ERR_EOF
    }

    
    if self.s[p] != '{' {
        return types.ERR_INVALID_CHAR
    }

    
    self.p = p + 1
    return 0
}

func (self *Parser) array() types.ParsingError {
    n := len(self.s)
    p := self.lspace(self.p)

    
    if p >= n {
        return types.ERR_EOF
    }

    
    if self.s[p] != '[' {
        return types.ERR_INVALID_CHAR
    }

    
    self.p = p + 1
    return 0
}

func (self *Parser) lspace(sp int) int {
    ns := len(self.s)
    for ; sp<ns && isSpace(self.s[sp]); sp+=1 {}

    return sp
}

func (self *Parser) backward() {
    for ; self.p >= 0 && isSpace(self.s[self.p]); self.p-=1 {}
}

func (self *Parser) decodeArray(ret *linkedNodes) (Node, types.ParsingError) {
    sp := self.p
    ns := len(self.s)

    
    if self.p = self.lspace(sp); self.p >= ns {
        return Node{}, types.ERR_EOF
    }

    
    if self.s[self.p] == ']' {
        self.p++
        return Node{t: types.V_ARRAY}, 0
    }

    
    for {
        var val Node
        var err types.ParsingError

        if self.skipValue {
            
            var start int
            if start, err = self.skipFast(); err != 0 {
                return Node{}, err
            }
            if self.p > ns {
                return Node{}, types.ERR_EOF
            }
            t := switchRawType(self.s[start])
            if t == _V_NONE {
                return Node{}, types.ERR_INVALID_CHAR
            }
            val = newRawNode(self.s[start:self.p], t, false)
        }else{
            
            if val, err = self.Parse(); err != 0 {
                return Node{}, err
            }
        }

        
        ret.Push(val)
        self.p = self.lspace(self.p)

        
        if self.p >= ns {
            return Node{}, types.ERR_EOF
        }

        
        switch self.s[self.p] {
            case ',' : self.p++
            case ']' : self.p++; return newArray(ret), 0
            default:
                
                
                
                return Node{}, types.ERR_INVALID_CHAR
        }
    }
}

func (self *Parser) decodeObject(ret *linkedPairs) (Node, types.ParsingError) {
    sp := self.p
    ns := len(self.s)

    
    if self.p = self.lspace(sp); self.p >= ns {
        return Node{}, types.ERR_EOF
    }

    
    if self.s[self.p] == '}' {
        self.p++
        return Node{t: types.V_OBJECT}, 0
    }

    
    for {
        var val Node
        var njs types.JsonState
        var err types.ParsingError

        
        if njs = self.decodeValue(); njs.Vt != types.V_STRING {
            return Node{}, types.ERR_INVALID_CHAR
        }

        
        idx := self.p - 1
        key := self.s[njs.Iv:idx]

        
        if njs.Ep != -1 {
            if key, err = unquote(key); err != 0 {
                return Node{}, err
            }
        }

        
        if err = self.delim(); err != 0 {
            return Node{}, err
        }

        
        if self.skipValue {
            
            var start int
            if start, err = self.skipFast(); err != 0 {
                return Node{}, err
            }
            if self.p > ns {
                return Node{}, types.ERR_EOF
            }
            t := switchRawType(self.s[start])
            if t == _V_NONE {
                return Node{}, types.ERR_INVALID_CHAR
            }
            val = newRawNode(self.s[start:self.p], t, false)
        } else {
            
            if val, err = self.Parse(); err != 0 {
                return Node{}, err
            }
        }

        
        
        ret.Push(NewPair(key, val))
        self.p = self.lspace(self.p)

        
        if self.p >= ns {
            return Node{}, types.ERR_EOF
        }

        
        switch self.s[self.p] {
            case ',' : self.p++
            case '}' : self.p++; return newObject(ret), 0
        default:
            
            
            
            return Node{}, types.ERR_INVALID_CHAR
        }
    }
}

func (self *Parser) decodeString(iv int64, ep int) (Node, types.ParsingError) {
    p := self.p - 1
    s := self.s[iv:p]

    
    if ep == -1 {
        return NewString(s), 0
    }

    
    out, err := unquote(s)

    
    if err != 0 {
        return Node{}, err
    } else {
        return newBytes(rt.Str2Mem(out)), 0
    }
}



func (self *Parser) Pos() int {
    return self.p
}





func (self *Parser) Parse() (Node, types.ParsingError) {
    switch val := self.decodeValue(); val.Vt {
        case types.V_EOF     : return Node{}, types.ERR_EOF
        case types.V_NULL    : return nullNode, 0
        case types.V_TRUE    : return trueNode, 0
        case types.V_FALSE   : return falseNode, 0
        case types.V_STRING  : return self.decodeString(val.Iv, val.Ep)
        case types.V_ARRAY:
            s := self.p - 1;
            if p := skipBlank(self.s, self.p); p >= self.p && self.s[p] == ']' {
                self.p = p + 1
                return Node{t: types.V_ARRAY}, 0
            }
            if self.noLazy {
                if self.loadOnce {
                    self.noLazy = false
                }
                return self.decodeArray(new(linkedNodes))
            }
            
            if self.loadOnce {
                self.p = s
                s, e := self.skipFast()
                if e != 0 {
                    return Node{}, e
                }
                return newRawNode(self.s[s:self.p], types.V_ARRAY, true), 0
            }
            return newLazyArray(self), 0
        case types.V_OBJECT:
            s := self.p - 1;
            if p := skipBlank(self.s, self.p); p >= self.p && self.s[p] == '}' {
                self.p = p + 1
                return Node{t: types.V_OBJECT}, 0
            }
            
            if self.noLazy {
                if self.loadOnce {
                    self.noLazy = false
                }
                return self.decodeObject(new(linkedPairs))
            }
            if self.loadOnce {
                self.p = s
                s, e := self.skipFast()
                if e != 0 {
                    return Node{}, e
                }
                return newRawNode(self.s[s:self.p], types.V_OBJECT, true), 0
            }
            return newLazyObject(self), 0
        case types.V_DOUBLE  : return NewNumber(self.s[val.Ep:self.p]), 0
        case types.V_INTEGER : return NewNumber(self.s[val.Ep:self.p]), 0
        default              : return Node{}, types.ParsingError(-val.Vt)
    }
}

func (self *Parser) searchKey(match string) types.ParsingError {
    ns := len(self.s)
    if err := self.object(); err != 0 {
        return err
    }

    
    if self.p = self.lspace(self.p); self.p >= ns {
        return types.ERR_EOF
    }

    
    if self.s[self.p] == '}' {
        self.p++
        return _ERR_NOT_FOUND
    }

    var njs types.JsonState
    var err types.ParsingError
    
    for {

        
        if njs = self.decodeValue(); njs.Vt != types.V_STRING {
            return types.ERR_INVALID_CHAR
        }

        
        idx := self.p - 1
        key := self.s[njs.Iv:idx]

        
        if njs.Ep != -1 {
            if key, err = unquote(key); err != 0 {
                return err
            }
        }

        
        if err = self.delim(); err != 0 {
            return err
        }

        
        if key != match {
            if _, err = self.skipFast(); err != 0 {
                return err
            }
        } else {
            return 0
        }

        
        self.p = self.lspace(self.p)
        if self.p >= ns {
            return types.ERR_EOF
        }

        
        switch self.s[self.p] {
        case ',':
            self.p++
        case '}':
            self.p++
            return _ERR_NOT_FOUND
        default:
            return types.ERR_INVALID_CHAR
        }
    }
}

func (self *Parser) searchIndex(idx int) types.ParsingError {
    ns := len(self.s)
    if err := self.array(); err != 0 {
        return err
    }

    
    if self.p = self.lspace(self.p); self.p >= ns {
        return types.ERR_EOF
    }

    
    if self.s[self.p] == ']' {
        self.p++
        return _ERR_NOT_FOUND
    }

    var err types.ParsingError
    
    for i := 0; i < idx; i++ {

        
        if _, err = self.skipFast(); err != 0 {
            return err
        }

        
        self.p = self.lspace(self.p)
        if self.p >= ns {
            return types.ERR_EOF
        }

        
        switch self.s[self.p] {
        case ',':
            self.p++
        case ']':
            self.p++
            return _ERR_NOT_FOUND
        default:
            return types.ERR_INVALID_CHAR
        }
    }

    return 0
}

func (self *Node) skipNextNode() *Node {
    if !self.isLazy() {
        return nil
    }

    parser, stack := self.getParserAndArrayStack()
    ret := &stack.v
    sp := parser.p
    ns := len(parser.s)

    
    if parser.p = parser.lspace(sp); parser.p >= ns {
        return newSyntaxError(parser.syntaxError(types.ERR_EOF))
    }

    
    if parser.s[parser.p] == ']' {
        parser.p++
        self.setArray(ret)
        return nil
    }

    var val Node
    
    if start, err := parser.skipFast(); err != 0 {
        return newSyntaxError(parser.syntaxError(err))
    } else {
        t := switchRawType(parser.s[start])
        if t == _V_NONE {
            return newSyntaxError(parser.syntaxError(types.ERR_INVALID_CHAR))
        }
        val = newRawNode(parser.s[start:parser.p], t, false)
    }

    
    ret.Push(val)
    self.l++
    parser.p = parser.lspace(parser.p)

    
    if parser.p >= ns {
        return newSyntaxError(parser.syntaxError(types.ERR_EOF))
    }

    
    switch parser.s[parser.p] {
    case ',':
        parser.p++
        return ret.At(ret.Len()-1)
    case ']':
        parser.p++
        self.setArray(ret)
        return ret.At(ret.Len()-1)
    default:
        return newSyntaxError(parser.syntaxError(types.ERR_INVALID_CHAR))
    }
}

func (self *Node) skipNextPair() (*Pair) {
    if !self.isLazy() {
        return nil
    }

    parser, stack := self.getParserAndObjectStack()
    ret := &stack.v
    sp := parser.p
    ns := len(parser.s)

    
    if parser.p = parser.lspace(sp); parser.p >= ns {
        return newErrorPair(parser.syntaxError(types.ERR_EOF))
    }

    
    if parser.s[parser.p] == '}' {
        parser.p++
        self.setObject(ret)
        return nil
    }

    
    var val Node
    var njs types.JsonState
    var err types.ParsingError

    
    if njs = parser.decodeValue(); njs.Vt != types.V_STRING {
        return newErrorPair(parser.syntaxError(types.ERR_INVALID_CHAR))
    }

    
    idx := parser.p - 1
    key := parser.s[njs.Iv:idx]

    
    if njs.Ep != -1 {
        if key, err = unquote(key); err != 0 {
            return newErrorPair(parser.syntaxError(err))
        }
    }

    
    if err = parser.delim(); err != 0 {
        return newErrorPair(parser.syntaxError(err))
    }

    
    if start, err := parser.skipFast(); err != 0 {
        return newErrorPair(parser.syntaxError(err))
    } else {
        t := switchRawType(parser.s[start])
        if t == _V_NONE {
            return newErrorPair(parser.syntaxError(types.ERR_INVALID_CHAR))
        }
        val = newRawNode(parser.s[start:parser.p], t, false)
    }

    
    ret.Push(NewPair(key, val))
    self.l++
    parser.p = parser.lspace(parser.p)

    
    if parser.p >= ns {
        return newErrorPair(parser.syntaxError(types.ERR_EOF))
    }

    
    switch parser.s[parser.p] {
    case ',':
        parser.p++
        return ret.At(ret.Len()-1)
    case '}':
        parser.p++
        self.setObject(ret)
        return ret.At(ret.Len()-1)
    default:
        return newErrorPair(parser.syntaxError(types.ERR_INVALID_CHAR))
    }
}





func Loads(src string) (int, interface{}, error) {
    ps := &Parser{s: src}
    np, err := ps.Parse()

    
    if err != 0 {
        return 0, nil, ps.ExportError(err)
    } else {
        x, err := np.Interface()
        if err != nil {
            return 0, nil, err
        }
        return ps.Pos(), x, nil
    }
}


func LoadsUseNumber(src string) (int, interface{}, error) {
    ps := &Parser{s: src}
    np, err := ps.Parse()

    
    if err != 0 {
        return 0, nil, err
    } else {
        x, err := np.InterfaceUseNumber()
        if err != nil {
            return 0, nil, err
        }
        return ps.Pos(), x, nil
    }
}


func NewParser(src string) *Parser {
    return &Parser{s: src}
}


func NewParserObj(src string) Parser {
    return Parser{s: src}
}




func (self *Parser) decodeNumber(decode bool) {
    if !decode && self.dbuf != nil {
        types.FreeDbuf(self.dbuf)
        self.dbuf = nil
        return
    }
    if decode && self.dbuf == nil {
        self.dbuf = types.NewDbuf()
    }
}


func (self *Parser) ExportError(err types.ParsingError) error {
    if err == _ERR_NOT_FOUND {
        return ErrNotExist
    }
    return fmt.Errorf("%q", SyntaxError{
        Pos : self.p,
        Src : self.s,
        Code: err,
    }.Description())
}

func backward(src string, i int) int {
    for ; i>=0 && isSpace(src[i]); i-- {}
    return i
}


func newRawNode(str string, typ types.ValueType, lock bool) Node {
    ret := Node{
        t: typ | _V_RAW,
        p: rt.StrPtr(str),
        l: uint(len(str)),
    }
    if lock {
        ret.m = new(sync.RWMutex)
    }
    return ret
}

var typeJumpTable = [256]types.ValueType{
    '"' : types.V_STRING,
    '-' : _V_NUMBER,
    '0' : _V_NUMBER,
    '1' : _V_NUMBER,
    '2' : _V_NUMBER,
    '3' : _V_NUMBER,
    '4' : _V_NUMBER,
    '5' : _V_NUMBER,
    '6' : _V_NUMBER,
    '7' : _V_NUMBER,
    '8' : _V_NUMBER,
    '9' : _V_NUMBER,
    '[' : types.V_ARRAY,
    'f' : types.V_FALSE,
    'n' : types.V_NULL,
    't' : types.V_TRUE,
    '{' : types.V_OBJECT,
}

func switchRawType(c byte) types.ValueType {
    return typeJumpTable[c]
}

func (self *Node) loadt() types.ValueType {
    return (types.ValueType)(atomic.LoadInt64(&self.t))
}

func (self *Node) lock() bool {
    if m := self.m; m != nil {
        m.Lock()
        return true
    }
    return false
}

func (self *Node) unlock() {
    if m := self.m; m != nil {
        m.Unlock()
    }
}

func (self *Node) rlock() bool {
    if m := self.m; m != nil {
        m.RLock()
        return true
    }
    return false
}

func (self *Node) runlock() {
    if m := self.m; m != nil {
        m.RUnlock()
    }
}
