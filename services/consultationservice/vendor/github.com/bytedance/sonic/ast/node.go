

package ast

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/bytedance/sonic/internal/native/types"
	"github.com/bytedance/sonic/internal/rt"
)

const (
    _V_NONE         types.ValueType = 0
    _V_NODE_BASE    types.ValueType = 1 << 5
    _V_LAZY         types.ValueType = 1 << 7
    _V_RAW          types.ValueType = 1 << 8
    _V_NUMBER                       = _V_NODE_BASE + 1
    _V_ANY                          = _V_NODE_BASE + 2
    _V_ARRAY_LAZY                   = _V_LAZY | types.V_ARRAY
    _V_OBJECT_LAZY                  = _V_LAZY | types.V_OBJECT
    _MASK_LAZY                      = _V_LAZY - 1
    _MASK_RAW                      = _V_RAW - 1
)

const (
    V_NONE   = 0
    V_ERROR  = 1
    V_NULL   = int(types.V_NULL)
    V_TRUE   = int(types.V_TRUE)
    V_FALSE  = int(types.V_FALSE)
    V_ARRAY  = int(types.V_ARRAY)
    V_OBJECT = int(types.V_OBJECT)
    V_STRING = int(types.V_STRING)
    V_NUMBER = int(_V_NUMBER)
    V_ANY    = int(_V_ANY)
)

type Node struct {
    t types.ValueType
    l uint
    p unsafe.Pointer
    m *sync.RWMutex
}



func (self *Node) UnmarshalJSON(data []byte) (err error) {
    *self = NewRaw(string(data))
    return self.Check()
}

















func (self Node) Type() int {
    return int(self.t & _MASK_LAZY & _MASK_RAW)
}













func (self *Node) TypeSafe() int {
    return int(self.loadt() & _MASK_LAZY & _MASK_RAW)
}

func (self *Node) itype() types.ValueType {
    return self.t & _MASK_LAZY & _MASK_RAW
}


func (self *Node) Exists() bool {
    if self == nil {
        return false
    }
    t := self.loadt()
    return t != V_ERROR && t != _V_NONE
}


func (self *Node) Valid() bool {
    if self == nil {
        return false
    }
    return self.loadt() != V_ERROR
}




func (self *Node)  Check() error {
    if self == nil {
        return ErrNotExist
    } else if self.loadt() != V_ERROR {
        return nil
    } else {
        return self
    }
}




func (self Node) IsRaw() bool {
    return self.t & _V_RAW != 0
}


func (self *Node) isRaw() bool {
    return self.loadt() & _V_RAW != 0
}

func (self *Node) isLazy() bool {
    return self != nil && self.t & _V_LAZY != 0
}

func (self *Node) isAny() bool {
    return self != nil && self.loadt() == _V_ANY
}




func (self *Node) Raw() (string, error) {
    if self == nil {
        return "", ErrNotExist
    }
    lock := self.rlock()
    if !self.isRaw() {
        if lock {
            self.runlock()
        }
        buf, err := self.MarshalJSON()
        return rt.Mem2Str(buf), err
    }
    ret := self.toString()
    if lock {
        self.runlock()
    }
    return ret, nil
}

func (self *Node) checkRaw() error {
    if err := self.Check(); err != nil {
        return err
    }
    if self.isRaw() {
        self.parseRaw(false)
    }
    return self.Check()
}




func (self *Node) Bool() (bool, error) {
    if err := self.checkRaw(); err != nil {
        return false, err
    }
    switch self.t {
        case types.V_TRUE  : return true , nil
        case types.V_FALSE : return false, nil
        case types.V_NULL  : return false, nil
        case _V_NUMBER     : 
            if i, err := self.toInt64(); err == nil {
                return i != 0, nil
            } else if f, err := self.toFloat64(); err == nil {
                return f != 0, nil
            } else {
                return false, err
            }
        case types.V_STRING: return strconv.ParseBool(self.toString())
        case _V_ANY        :   
            any := self.packAny()     
            switch v := any.(type) {
                case bool   : return v, nil
                case int    : return v != 0, nil
                case int8   : return v != 0, nil
                case int16  : return v != 0, nil
                case int32  : return v != 0, nil
                case int64  : return v != 0, nil
                case uint   : return v != 0, nil
                case uint8  : return v != 0, nil
                case uint16 : return v != 0, nil
                case uint32 : return v != 0, nil
                case uint64 : return v != 0, nil
                case float32: return v != 0, nil
                case float64: return v != 0, nil
                case string : return strconv.ParseBool(v)
                case json.Number: 
                    if i, err := v.Int64(); err == nil {
                        return i != 0, nil
                    } else if f, err := v.Float64(); err == nil {
                        return f != 0, nil
                    } else {
                        return false, err
                    }
                default: return false, ErrUnsupportType
            }
        default            : return false, ErrUnsupportType
    }
}




func (self *Node) Int64() (int64, error) {
    if err := self.checkRaw(); err != nil {
        return 0, err
    }
    switch self.t {
        case _V_NUMBER, types.V_STRING :
            if i, err := self.toInt64(); err == nil {
                return i, nil
            } else if f, err := self.toFloat64(); err == nil {
                return int64(f), nil
            } else {
                return 0, err
            }
        case types.V_TRUE     : return 1, nil
        case types.V_FALSE    : return 0, nil
        case types.V_NULL     : return 0, nil
        case _V_ANY           :  
            any := self.packAny()
            switch v := any.(type) {
                case bool   : if v { return 1, nil } else { return 0, nil }
                case int    : return int64(v), nil
                case int8   : return int64(v), nil
                case int16  : return int64(v), nil
                case int32  : return int64(v), nil
                case int64  : return int64(v), nil
                case uint   : return int64(v), nil
                case uint8  : return int64(v), nil
                case uint16 : return int64(v), nil
                case uint32 : return int64(v), nil
                case uint64 : return int64(v), nil
                case float32: return int64(v), nil
                case float64: return int64(v), nil
                case string : 
                    if i, err := strconv.ParseInt(v, 10, 64); err == nil {
                        return i, nil
                    } else if f, err := strconv.ParseFloat(v, 64); err == nil {
                        return int64(f), nil
                    } else {
                        return 0, err
                    }
                case json.Number: 
                    if i, err := v.Int64(); err == nil {
                        return i, nil
                    } else if f, err := v.Float64(); err == nil {
                        return int64(f), nil
                    } else {
                        return 0, err
                    }
                default: return 0, ErrUnsupportType
            }
        default               : return 0, ErrUnsupportType
    }
}


func (self *Node) StrictInt64() (int64, error) {
    if err := self.checkRaw(); err != nil {
        return 0, err
    }
    switch self.t {
        case _V_NUMBER        : return self.toInt64()
        case _V_ANY           :  
            any := self.packAny()
            switch v := any.(type) {
                case int   : return int64(v), nil
                case int8  : return int64(v), nil
                case int16 : return int64(v), nil
                case int32 : return int64(v), nil
                case int64 : return int64(v), nil
                case uint  : return int64(v), nil
                case uint8 : return int64(v), nil
                case uint16: return int64(v), nil
                case uint32: return int64(v), nil
                case uint64: return int64(v), nil
                case json.Number: 
                    if i, err := v.Int64(); err == nil {
                        return i, nil
                    } else {
                        return 0, err
                    }
                default: return 0, ErrUnsupportType
            }
        default               : return 0, ErrUnsupportType
    }
}

func castNumber(v bool) json.Number {
    if v {
        return json.Number("1")
    } else {
        return json.Number("0")
    }
}




func (self *Node) Number() (json.Number, error) {
    if err := self.checkRaw(); err != nil {
        return json.Number(""), err
    }
    switch self.t {
        case _V_NUMBER        : return self.toNumber(), nil
        case types.V_STRING : 
            if _, err := self.toInt64(); err == nil {
                return self.toNumber(), nil
            } else if _, err := self.toFloat64(); err == nil {
                return self.toNumber(), nil
            } else {
                return json.Number(""), err
            }
        case types.V_TRUE     : return json.Number("1"), nil
        case types.V_FALSE    : return json.Number("0"), nil
        case types.V_NULL     : return json.Number("0"), nil
        case _V_ANY           :        
            any := self.packAny()
            switch v := any.(type) {
                case bool   : return castNumber(v), nil
                case int    : return castNumber(v != 0), nil
                case int8   : return castNumber(v != 0), nil
                case int16  : return castNumber(v != 0), nil
                case int32  : return castNumber(v != 0), nil
                case int64  : return castNumber(v != 0), nil
                case uint   : return castNumber(v != 0), nil
                case uint8  : return castNumber(v != 0), nil
                case uint16 : return castNumber(v != 0), nil
                case uint32 : return castNumber(v != 0), nil
                case uint64 : return castNumber(v != 0), nil
                case float32: return castNumber(v != 0), nil
                case float64: return castNumber(v != 0), nil
                case string : 
                    if _, err := strconv.ParseFloat(v, 64); err == nil {
                        return json.Number(v), nil
                    } else {
                        return json.Number(""), err
                    }
                case json.Number: return v, nil
                default: return json.Number(""), ErrUnsupportType
            }
        default               : return json.Number(""), ErrUnsupportType
    }
}


func (self *Node) StrictNumber() (json.Number, error) {
    if err := self.checkRaw(); err != nil {
        return json.Number(""), err
    }
    switch self.t {
        case _V_NUMBER        : return self.toNumber()  , nil
        case _V_ANY        :        
            if v, ok := self.packAny().(json.Number); ok {
                return v, nil
            } else {
                return json.Number(""), ErrUnsupportType
            }
        default               : return json.Number(""), ErrUnsupportType
    }
}




func (self *Node) String() (string, error) {
    if err := self.checkRaw(); err != nil {
        return "", err
    }
    switch self.t {
        case types.V_NULL    : return "" , nil
        case types.V_TRUE    : return "true" , nil
        case types.V_FALSE   : return "false", nil
        case types.V_STRING, _V_NUMBER  : return self.toString(), nil
        case _V_ANY          :        
        any := self.packAny()
        switch v := any.(type) {
            case bool   : return strconv.FormatBool(v), nil
            case int    : return strconv.Itoa(v), nil
            case int8   : return strconv.Itoa(int(v)), nil
            case int16  : return strconv.Itoa(int(v)), nil
            case int32  : return strconv.Itoa(int(v)), nil
            case int64  : return strconv.Itoa(int(v)), nil
            case uint   : return strconv.Itoa(int(v)), nil
            case uint8  : return strconv.Itoa(int(v)), nil
            case uint16 : return strconv.Itoa(int(v)), nil
            case uint32 : return strconv.Itoa(int(v)), nil
            case uint64 : return strconv.Itoa(int(v)), nil
            case float32: return strconv.FormatFloat(float64(v), 'g', -1, 64), nil
            case float64: return strconv.FormatFloat(float64(v), 'g', -1, 64), nil
            case string : return v, nil 
            case json.Number: return v.String(), nil
            default: return "", ErrUnsupportType
        }
        default              : return ""     , ErrUnsupportType
    }
}



func (self *Node) StrictString() (string, error) {
    if err := self.checkRaw(); err != nil {
        return "", err
    }
    switch self.t {
        case types.V_STRING  : return self.toString(), nil
        case _V_ANY          :        
            if v, ok := self.packAny().(string); ok {
                return v, nil
            } else {
                return "", ErrUnsupportType
            }
        default              : return "", ErrUnsupportType
    }
}




func (self *Node) Float64() (float64, error) {
    if err := self.checkRaw(); err != nil {
        return 0.0, err
    }
    switch self.t {
        case _V_NUMBER, types.V_STRING : return self.toFloat64()
        case types.V_TRUE    : return 1.0, nil
        case types.V_FALSE   : return 0.0, nil
        case types.V_NULL    : return 0.0, nil
        case _V_ANY          :        
            any := self.packAny()
            switch v := any.(type) {
                case bool    : 
                    if v {
                        return 1.0, nil
                    } else {
                        return 0.0, nil
                    }
                case int    : return float64(v), nil
                case int8   : return float64(v), nil
                case int16  : return float64(v), nil
                case int32  : return float64(v), nil
                case int64  : return float64(v), nil
                case uint   : return float64(v), nil
                case uint8  : return float64(v), nil
                case uint16 : return float64(v), nil
                case uint32 : return float64(v), nil
                case uint64 : return float64(v), nil
                case float32: return float64(v), nil
                case float64: return float64(v), nil
                case string : 
                    if f, err := strconv.ParseFloat(v, 64); err == nil {
                        return float64(f), nil
                    } else {
                        return 0, err
                    }
                case json.Number: 
                    if f, err := v.Float64(); err == nil {
                        return float64(f), nil
                    } else {
                        return 0, err
                    }
                default     : return 0, ErrUnsupportType
            }
        default             : return 0.0, ErrUnsupportType
    }
}


func (self *Node) StrictFloat64() (float64, error) {
    if err := self.checkRaw(); err != nil {
        return 0.0, err
    }
    switch self.t {
        case _V_NUMBER       : return self.toFloat64()
        case _V_ANY        :        
            any := self.packAny()
            switch v := any.(type) {
                case float32 : return float64(v), nil
                case float64 : return float64(v), nil
                default      : return 0, ErrUnsupportType
            }
        default              : return 0.0, ErrUnsupportType
    }
}





func (self *Node) Len() (int, error) {
    if err := self.checkRaw(); err != nil {
        return 0, err
    }
    if self.t == types.V_ARRAY || self.t == types.V_OBJECT || self.t == _V_ARRAY_LAZY || self.t == _V_OBJECT_LAZY || self.t == types.V_STRING {
        return int(self.l), nil
    } else if self.t == _V_NONE || self.t == types.V_NULL {
        return 0, nil
    } else {
        return 0, ErrUnsupportType
    }
}

func (self *Node) len() int {
    return int(self.l)
}


func (self *Node) Cap() (int, error) {
    if err := self.checkRaw(); err != nil {
        return 0, err
    }
    switch self.t {
    case types.V_ARRAY: return (*linkedNodes)(self.p).Cap(), nil
    case types.V_OBJECT: return (*linkedPairs)(self.p).Cap(), nil
    case _V_ARRAY_LAZY: return (*parseArrayStack)(self.p).v.Cap(), nil
    case _V_OBJECT_LAZY: return (*parseObjectStack)(self.p).v.Cap(), nil
    case _V_NONE, types.V_NULL: return 0, nil
    default: return 0, ErrUnsupportType
    }
}




func (self *Node) Set(key string, node Node) (bool, error) {
    if err := self.checkRaw(); err != nil {
        return false, err
    }
    if err := node.Check(); err != nil {
        return false, err 
    }
    
    if self.t == _V_NONE || self.t == types.V_NULL {
        *self = NewObject([]Pair{NewPair(key, node)})
        return false, nil
    } else if self.itype() != types.V_OBJECT {
        return false, ErrUnsupportType
    }

    p := self.Get(key)

    if !p.Exists() {
        
        if self.len() == 0 {
            *self = newObject(new(linkedPairs))
        }
        s := (*linkedPairs)(self.p)
        s.Push(NewPair(key, node))
        self.l++
        return false, nil

    } else if err := p.Check(); err != nil {
        return false, err
    } 

    *p = node
    return true, nil
}


func (self *Node) SetAny(key string, val interface{}) (bool, error) {
    return self.Set(key, NewAny(val))
}


func (self *Node) Unset(key string) (bool, error) {
    if err := self.should(types.V_OBJECT); err != nil {
        return false, err
    }
    
    if err := self.skipAllKey(); err != nil {
        return false, err
    }
    p, i := self.skipKey(key)
    if !p.Exists() {
        return false, nil
    } else if err := p.Check(); err != nil {
        return false, err
    }
    self.removePairAt(i)
    return true, nil
}




func (self *Node) SetByIndex(index int, node Node) (bool, error) {
    if err := self.checkRaw(); err != nil {
        return false, err 
    }
    if err := node.Check(); err != nil {
        return false, err 
    }

    if index == 0 && (self.t == _V_NONE || self.t == types.V_NULL) {
        *self = NewArray([]Node{node})
        return false, nil
    }

    p := self.Index(index)
    if !p.Exists() {
        return false, ErrNotExist
    } else if err := p.Check(); err != nil {
        return false, err
    }

    *p = node
    return true, nil
}


func (self *Node) SetAnyByIndex(index int, val interface{}) (bool, error) {
    return self.SetByIndex(index, NewAny(val))
}





func (self *Node) UnsetByIndex(index int) (bool, error) {
    if err := self.checkRaw(); err != nil {
        return false, err
    }

    var p *Node
    it := self.itype()

    if it == types.V_ARRAY {
        if err := self.skipAllIndex(); err != nil {
            return false, err
        }
        p = self.nodeAt(index)
    } else if it == types.V_OBJECT {
        if err := self.skipAllKey(); err != nil {
            return false, err
        }
        pr := self.pairAt(index)
        if pr == nil {
           return false, ErrNotExist
        }
        p = &pr.Value
    } else {
        return false, ErrUnsupportType
    }

    if !p.Exists() {
        return false, ErrNotExist
    }

    
    if index == self.len() - 1 {
        return true, self.Pop()
    }

    
    if it == types.V_ARRAY {
        self.removeNode(index)
    }else if it == types.V_OBJECT {
        self.removePair(index)
    }
    return true, nil
}




func (self *Node) Add(node Node) error {
    if err := self.checkRaw(); err != nil {
        return err
    }

    if self != nil && (self.t == _V_NONE || self.t == types.V_NULL) {
        *self = NewArray([]Node{node})
        return nil
    }
    if err := self.should(types.V_ARRAY); err != nil {
        return err
    }

    s, err := self.unsafeArray()
    if err != nil {
        return err
    }

    
    s.Push(node)
    self.l++
    return nil
}


func (self *Node) Pop() error {
    if err := self.checkRaw(); err != nil {
        return err
    }

    if it := self.itype(); it == types.V_ARRAY {
        s, err := self.unsafeArray()
        if err != nil {
            return err
        }
        
        for i := s.Len()-1; i >= 0; i-- {
            if s.At(i).Exists() {
                s.Pop()
                self.l--
                break
            }
            s.Pop()
        }

    } else if it == types.V_OBJECT {
        s, err := self.unsafeMap()
        if err != nil {
            return err
        }
        
        for i := s.Len()-1; i >= 0; i-- {
            if p := s.At(i); p != nil && p.Value.Exists() {
                s.Pop()
                self.l--
                break
            }
            s.Pop()
        }

    } else {
        return ErrUnsupportType
    }

    return nil
}





func (self *Node) Move(dst, src int) error {
    if err := self.should(types.V_ARRAY); err != nil {
        return err
    }

    s, err := self.unsafeArray()
    if err != nil {
        return err
    }

    
    if l :=  s.Len(); self.len() != l {
        di, si := dst, src
        
        for i := 0; i < l; i++ {
            if s.At(i).Exists() {
                di--
                si--
            }
            if di == -1 {
                dst = i
                di--
            } 
            if si == -1 {
                src = i
                si--
            }
            if di == -2 && si == -2 {
                break
            }
        }
    }

    s.MoveOne(src, dst)
    return nil
}


func (self *Node) AddAny(val interface{}) error {
    return self.Add(NewAny(val))
}






func (self *Node) GetByPath(path ...interface{}) *Node {
    if !self.Valid() {
        return self
    }
    var s = self
    for _, p := range path {
        switch p := p.(type) {
        case int:
            s = s.Index(p)
            if !s.Valid() {
                return s
            }
        case string:
            s = s.Get(p)
            if !s.Valid() {
                return s
            }
        default:
            panic("path must be either int or string")
        }
    }
    return s
}


func (self *Node) Get(key string) *Node {
    if err := self.should(types.V_OBJECT); err != nil {
        return unwrapError(err)
    }
    n, _ := self.skipKey(key)
    return n
}



func (self *Node) Index(idx int) *Node {
    if err := self.checkRaw(); err != nil {
        return unwrapError(err)
    }

    it := self.itype()
    if it == types.V_ARRAY {
        return self.skipIndex(idx)

    }else if it == types.V_OBJECT {
        pr := self.skipIndexPair(idx)
        if pr == nil {
           return newError(_ERR_NOT_FOUND, "value not exists")
        }
        return &pr.Value

    } else {
        return newError(_ERR_UNSUPPORT_TYPE, fmt.Sprintf("unsupported type: %v", self.itype()))
    }
}



func (self *Node) IndexPair(idx int) *Pair {
    if err := self.should(types.V_OBJECT); err != nil {
        return nil
    }
    return self.skipIndexPair(idx)
}

func (self *Node) indexOrGet(idx int, key string) (*Node, int) {
    if err := self.should(types.V_OBJECT); err != nil {
        return unwrapError(err), idx
    }

    pr := self.skipIndexPair(idx)
    if pr != nil && pr.Key == key {
        return &pr.Value, idx
    }

    return self.skipKey(key)
}



func (self *Node) IndexOrGet(idx int, key string) *Node {
    node, _ := self.indexOrGet(idx, key)
    return node
}



func (self *Node) IndexOrGetWithIdx(idx int, key string) (*Node, int) {
    return self.indexOrGet(idx, key)
}




func (self *Node) Map() (map[string]interface{}, error) {
    if self.isAny() {
        any := self.packAny()
        if v, ok := any.(map[string]interface{}); ok {
            return v, nil
        } else {
            return nil, ErrUnsupportType
        }
    }
    if err := self.should(types.V_OBJECT); err != nil {
        return nil, err
    }
    if err := self.loadAllKey(false); err != nil {
        return nil, err
    }
    return self.toGenericObject()
}


func (self *Node) MapUseNumber() (map[string]interface{}, error) {
    if self.isAny() {
        any := self.packAny()
        if v, ok := any.(map[string]interface{}); ok {
            return v, nil
        } else {
            return nil, ErrUnsupportType
        }
    }
    if err := self.should(types.V_OBJECT); err != nil {
        return nil, err
    }
    if err := self.loadAllKey(false); err != nil {
        return nil, err
    }
    return self.toGenericObjectUseNumber()
}



func (self *Node) MapUseNode() (map[string]Node, error) {
    if self.isAny() {
        any := self.packAny()
        if v, ok := any.(map[string]Node); ok {
            return v, nil
        } else {
            return nil, ErrUnsupportType
        }
    }
    if err := self.should(types.V_OBJECT); err != nil {
        return nil, err
    }
    if err := self.skipAllKey(); err != nil {
        return nil, err
    }
    return self.toGenericObjectUseNode()
}















//go:nocheckptr
func (self *Node) unsafeMap() (*linkedPairs, error) {
    if err := self.skipAllKey(); err != nil {
        return nil, err
    }
    if self.p == nil {
        *self = newObject(new(linkedPairs))
    }
    return (*linkedPairs)(self.p), nil
}



func (self *Node) SortKeys(recurse bool) error {
    
    if err := self.checkRaw(); err != nil {
        return err
    }
    if self.itype() == types.V_OBJECT {
        return self.sortKeys(recurse)
    } else if self.itype() == types.V_ARRAY {
        var err error
        err2 := self.ForEach(func(path Sequence, node *Node) bool {
            it := node.itype()
            if it == types.V_ARRAY || it == types.V_OBJECT {
                err = node.SortKeys(recurse)
                if err != nil {
                    return false
                }
            }
            return true
        })
        if err != nil {
            return err
        }
        return err2
    } else {
        return nil
    }
}

func (self *Node) sortKeys(recurse bool) (err error) {
    
    if err := self.checkRaw(); err != nil {
        return err
    }
    ps, err := self.unsafeMap()
    if err != nil {
        return err
    }
    ps.Sort()
    if recurse {
        var sc Scanner
        sc = func(path Sequence, node *Node) bool {
            if node.itype() == types.V_OBJECT {
                if err := node.sortKeys(recurse); err != nil {
                    return false
                }
            }
            if node.itype() == types.V_ARRAY {
                if err := node.ForEach(sc); err != nil {
                    return false
                }
            }
            return true
        }
        if err := self.ForEach(sc); err != nil {
            return err
        }
    }
    return nil
}


func (self *Node) Array() ([]interface{}, error) {
    if self.isAny() {
        any := self.packAny()
        if v, ok := any.([]interface{}); ok {
            return v, nil
        } else {
            return nil, ErrUnsupportType
        }
    }
    if err := self.should(types.V_ARRAY); err != nil {
        return nil, err
    }
    if err := self.loadAllIndex(false); err != nil {
        return nil, err
    }
    return self.toGenericArray()
}


func (self *Node) ArrayUseNumber() ([]interface{}, error) {
    if self.isAny() {
        any := self.packAny()
        if v, ok := any.([]interface{}); ok {
            return v, nil
        } else {
            return nil, ErrUnsupportType
        }
    }
    if err := self.should(types.V_ARRAY); err != nil {
        return nil, err
    }
    if err := self.loadAllIndex(false); err != nil {
        return nil, err
    }
    return self.toGenericArrayUseNumber()
}



func (self *Node) ArrayUseNode() ([]Node, error) {
    if self.isAny() {
        any := self.packAny()
        if v, ok := any.([]Node); ok {
            return v, nil
        } else {
            return nil, ErrUnsupportType
        }
    }
    if err := self.should(types.V_ARRAY); err != nil {
        return nil, err
    }
    if err := self.skipAllIndex(); err != nil {
        return nil, err
    }
    return self.toGenericArrayUseNode()
}
















func (self *Node) unsafeArray() (*linkedNodes, error) {
    if err := self.skipAllIndex(); err != nil {
        return nil, err
    }
    if self.p == nil {
        *self = newArray(new(linkedNodes))
    }
    return (*linkedNodes)(self.p), nil
}




func (self *Node) Interface() (interface{}, error) {
    if err := self.checkRaw(); err != nil {
        return nil, err
    }
    switch self.t {
        case V_ERROR         : return nil, self.Check()
        case types.V_NULL    : return nil, nil
        case types.V_TRUE    : return true, nil
        case types.V_FALSE   : return false, nil
        case types.V_ARRAY   : return self.toGenericArray()
        case types.V_OBJECT  : return self.toGenericObject()
        case types.V_STRING  : return self.toString(), nil
        case _V_NUMBER       : 
            v, err := self.toFloat64()
            if err != nil {
                return nil, err
            }
            return v, nil
        case _V_ARRAY_LAZY   :
            if err := self.loadAllIndex(false); err != nil {
                return nil, err
            }
            return self.toGenericArray()
        case _V_OBJECT_LAZY  :
            if err := self.loadAllKey(false); err != nil {
                return nil, err
            }
            return self.toGenericObject()
        case _V_ANY:
            switch v := self.packAny().(type) {
                case Node : return v.Interface()
                case *Node: return v.Interface()
                default   : return v, nil
            }
        default              : return nil,  ErrUnsupportType
    }
}

func (self *Node) packAny() interface{} {
    return *(*interface{})(self.p)
}



func (self *Node) InterfaceUseNumber() (interface{}, error) {
    if err := self.checkRaw(); err != nil {
        return nil, err
    }
    switch self.t {
        case V_ERROR         : return nil, self.Check()
        case types.V_NULL    : return nil, nil
        case types.V_TRUE    : return true, nil
        case types.V_FALSE   : return false, nil
        case types.V_ARRAY   : return self.toGenericArrayUseNumber()
        case types.V_OBJECT  : return self.toGenericObjectUseNumber()
        case types.V_STRING  : return self.toString(), nil
        case _V_NUMBER       : return self.toNumber(), nil
        case _V_ARRAY_LAZY   :
            if err := self.loadAllIndex(false); err != nil {
                return nil, err
            }
            return self.toGenericArrayUseNumber()
        case _V_OBJECT_LAZY  :
            if err := self.loadAllKey(false); err != nil {
                return nil, err
            }
            return self.toGenericObjectUseNumber()
        case _V_ANY          : return self.packAny(), nil
        default              : return nil, ErrUnsupportType
    }
}



func (self *Node) InterfaceUseNode() (interface{}, error) {
    if err := self.checkRaw(); err != nil {
        return nil, err
    }
    switch self.t {
        case types.V_ARRAY   : return self.toGenericArrayUseNode()
        case types.V_OBJECT  : return self.toGenericObjectUseNode()
        case _V_ARRAY_LAZY   :
            if err := self.skipAllIndex(); err != nil {
                return nil, err
            }
            return self.toGenericArrayUseNode()
        case _V_OBJECT_LAZY  :
            if err := self.skipAllKey(); err != nil {
                return nil, err
            }
            return self.toGenericObjectUseNode()
        default              : return *self, self.Check()
    }
}



func (self *Node) LoadAll() error {
    return self.Load()
}



func (self *Node) Load() error {
    switch self.t {
        case _V_ARRAY_LAZY: self.loadAllIndex(true)
        case _V_OBJECT_LAZY: self.loadAllKey(true)
        case V_ERROR: return self
        case V_NONE: return nil
    }
    if self.m == nil {
        self.m = new(sync.RWMutex)
    }
    return self.checkRaw()
}



func (self *Node) should(t types.ValueType) error {
    if err := self.checkRaw(); err != nil {
        return err
    }
    if  self.itype() != t {
        return ErrUnsupportType
    }
    return nil
}

func (self *Node) nodeAt(i int) *Node {
    var p *linkedNodes
    if self.isLazy() {
        _, stack := self.getParserAndArrayStack()
        p = &stack.v
    } else {
        p = (*linkedNodes)(self.p)
        if l := p.Len(); l != self.len() {
            
            for j:=0; j<l; j++ {
                v := p.At(j)
                if v.Exists() {
                    i--
                }
                if i < 0 {
                    return v
                }
            }
            return nil
        } 
    }
    return p.At(i)
}

func (self *Node) pairAt(i int) *Pair {
    var p *linkedPairs
    if self.isLazy() {
        _, stack := self.getParserAndObjectStack()
        p = &stack.v
    } else {
        p = (*linkedPairs)(self.p)
        if l := p.Len(); l != self.len() {
            
            for j:=0; j<l; j++ {
                v := p.At(j)
                if v != nil && v.Value.Exists() {
                    i--
                }
                if i < 0 {
                    return v
                }
            }
           return nil
       } 
    }
    return p.At(i)
}

func (self *Node) skipAllIndex() error {
    if !self.isLazy() {
        return nil
    }
    var err types.ParsingError
    parser, stack := self.getParserAndArrayStack()
    parser.skipValue = true
    parser.noLazy = true
    *self, err = parser.decodeArray(&stack.v)
    if err != 0 {
        return parser.ExportError(err)
    }
    return nil
}

func (self *Node) skipAllKey() error {
    if !self.isLazy() {
        return nil
    }
    var err types.ParsingError
    parser, stack := self.getParserAndObjectStack()
    parser.skipValue = true
    parser.noLazy = true
    *self, err = parser.decodeObject(&stack.v)
    if err != 0 {
        return parser.ExportError(err)
    }
    return nil
}

func (self *Node) skipKey(key string) (*Node, int) {
    nb := self.len()
    lazy := self.isLazy()

    if nb > 0 {
        
        var p *Pair
        var i int
        if lazy {
            s := (*parseObjectStack)(self.p)
            p, i = s.v.Get(key)
        } else {
            p, i = (*linkedPairs)(self.p).Get(key)
        }

        if p != nil {
            return &p.Value, i
        }
    }

    
    if !lazy {
        return nil, -1
    }

    
    for last, i := self.skipNextPair(), nb; last != nil; last, i = self.skipNextPair(), i+1 {
        if last.Value.Check() != nil {
            return &last.Value, -1
        }
        if last.Key == key {
            return &last.Value, i
        }
    }

    return nil, -1
}

func (self *Node) skipIndex(index int) *Node {
    nb := self.len()
    if nb > index {
        v := self.nodeAt(index)
        return v
    }
    if !self.isLazy() {
        return nil
    }

    
    for last := self.skipNextNode(); last != nil; last = self.skipNextNode(){
        if last.Check() != nil {
            return last
        }
        if self.len() > index {
            return last
        }
    }

    return nil
}

func (self *Node) skipIndexPair(index int) *Pair {
    nb := self.len()
    if nb > index {
        return self.pairAt(index)
    }
    if !self.isLazy() {
        return nil
    }

    
    for last := self.skipNextPair(); last != nil; last = self.skipNextPair(){
        if last.Value.Check() != nil {
            return last
        }
        if self.len() > index {
            return last
        }
    }

    return nil
}

func (self *Node) loadAllIndex(loadOnce bool) error {
    if !self.isLazy() {
        return nil
    }
    var err types.ParsingError
    parser, stack := self.getParserAndArrayStack()
    if !loadOnce {
        parser.noLazy = true
    } else {
        parser.loadOnce = true
    }
    *self, err = parser.decodeArray(&stack.v)
    if err != 0 {
        return parser.ExportError(err)
    }
    return nil
}

func (self *Node) loadAllKey(loadOnce bool) error {
    if !self.isLazy() {
        return nil
    }
    var err types.ParsingError
    parser, stack := self.getParserAndObjectStack()
    if !loadOnce {
        parser.noLazy = true
        *self, err = parser.decodeObject(&stack.v)
    } else {
        parser.loadOnce = true
        *self, err = parser.decodeObject(&stack.v)
    }
    if err != 0 {
        return parser.ExportError(err)
    }
    return nil
}

func (self *Node) removeNode(i int) {
    node := self.nodeAt(i)
    if node == nil {
        return
    }
    *node = Node{}
    
    self.l--
}

func (self *Node) removePair(i int) {
    last := self.pairAt(i)
    if last == nil {
        return
    }
    *last = Pair{}
    
    self.l--
}

func (self *Node) removePairAt(i int) {
    p := (*linkedPairs)(self.p).At(i)
    if p == nil {
        return
    }
    *p = Pair{}
    
    self.l--
}

func (self *Node) toGenericArray() ([]interface{}, error) {
    nb := self.len()
    if nb == 0 {
        return []interface{}{}, nil
    }
    ret := make([]interface{}, 0, nb)
    
    
    it := self.values()
    for v := it.next(); v != nil; v = it.next() {
        vv, err := v.Interface()
        if err != nil {
            return nil, err
        }
        ret = append(ret, vv)
    }

    
    return ret, nil
}

func (self *Node) toGenericArrayUseNumber() ([]interface{}, error) {
    nb := self.len()
    if nb == 0 {
        return []interface{}{}, nil
    }
    ret := make([]interface{}, 0, nb)

    
    it := self.values()
    for v := it.next(); v != nil; v = it.next() {
        vv, err := v.InterfaceUseNumber()
        if err != nil {
            return nil, err
        }
        ret = append(ret, vv)
    }

    
    return ret, nil
}

func (self *Node) toGenericArrayUseNode() ([]Node, error) {
    var nb = self.len()
    if nb == 0 {
        return []Node{}, nil
    }

    var s = (*linkedNodes)(self.p)
    var out = make([]Node, nb)
    s.ToSlice(out)

    return out, nil
}

func (self *Node) toGenericObject() (map[string]interface{}, error) {
    nb := self.len()
    if nb == 0 {
        return map[string]interface{}{}, nil
    }
    ret := make(map[string]interface{}, nb)

    
    it := self.properties()
    for v := it.next(); v != nil; v = it.next() {
        vv, err := v.Value.Interface()
        if err != nil {
            return nil, err
        }
        ret[v.Key] = vv
    }

    
    return ret, nil
}


func (self *Node) toGenericObjectUseNumber() (map[string]interface{}, error) {
    nb := self.len()
    if nb == 0 {
        return map[string]interface{}{}, nil
    }
    ret := make(map[string]interface{}, nb)

    
    it := self.properties()
    for v := it.next(); v != nil; v = it.next() {
        vv, err := v.Value.InterfaceUseNumber()
        if err != nil {
            return nil, err
        }
        ret[v.Key] = vv
    }

    
    return ret, nil
}

func (self *Node) toGenericObjectUseNode() (map[string]Node, error) {
    var nb = self.len()
    if nb == 0 {
        return map[string]Node{}, nil
    }

    var s = (*linkedPairs)(self.p)
    var out = make(map[string]Node, nb)
    s.ToMap(out)

    
    return out, nil
}



var (
    nullNode  = Node{t: types.V_NULL}
    trueNode  = Node{t: types.V_TRUE}
    falseNode = Node{t: types.V_FALSE}
)



func NewRaw(json string) Node {
    parser := NewParserObj(json)
    start, err := parser.skip()
    if err != 0 {
        return *newError(err, err.Message()) 
    }
    it := switchRawType(parser.s[start])
    if it == _V_NONE {
        return Node{}
    }
    return newRawNode(parser.s[start:parser.p], it, false)
}




func NewRawConcurrentRead(json string) Node {
    parser := NewParserObj(json)
    start, err := parser.skip()
    if err != 0 {
        return *newError(err, err.Message()) 
    }
    it := switchRawType(parser.s[start])
    if it == _V_NONE {
        return Node{}
    }
    return newRawNode(parser.s[start:parser.p], it, true)
}



func NewAny(any interface{}) Node {
    switch n := any.(type) {
    case Node:
        return n
    case *Node:
        return *n
    default:
        return Node{
            t: _V_ANY,
            p: unsafe.Pointer(&any),
        }
    }
}


func NewBytes(src []byte) Node {
    if len(src) == 0 {
        panic("empty src bytes")
    }
    out := rt.EncodeBase64ToString(src)
    return NewString(out)
}


func NewNull() Node {
    return Node{
        p: nil,
        t: types.V_NULL,
    }
}




func NewBool(v bool) Node {
    var t = types.V_FALSE
    if v {
        t = types.V_TRUE
    }
    return Node{
        p: nil,
        t: t,
    }
}



func NewNumber(v string) Node {
    return Node{
        l: uint(len(v)),
        p: rt.StrPtr(v),
        t: _V_NUMBER,
    }
}

func (node *Node) toNumber() json.Number {
    return json.Number(rt.StrFrom(node.p, int64(node.l)))
}

func (self *Node) toString() string {
    return rt.StrFrom(self.p, int64(self.l))
}

func (node *Node) toFloat64() (float64, error) {
    ret, err := node.toNumber().Float64()
    if err != nil {
        return 0, err
    }
    return ret, nil
}

func (node *Node) toInt64() (int64, error) {
    ret,err := node.toNumber().Int64()
    if err != nil {
        return 0, err
    }
    return ret, nil
}

func newBytes(v []byte) Node {
    return Node{
        t: types.V_STRING,
        p: mem2ptr(v),
        l: uint(len(v)),
    }
}





func NewString(v string) Node {
    return Node{
        t: types.V_STRING,
        p: rt.StrPtr(v),
        l: uint(len(v)),
    }
}



func NewArray(v []Node) Node {
    s := new(linkedNodes)
    s.FromSlice(v)
    return newArray(s)
}

const _Threshold_Index = 16

func newArray(v *linkedNodes) Node {
    return Node{
        t: types.V_ARRAY,
        l: uint(v.Len()),
        p: unsafe.Pointer(v),
    }
}

func (self *Node) setArray(v *linkedNodes) {
    self.t = types.V_ARRAY
    self.l = uint(v.Len())
    self.p = unsafe.Pointer(v)
}



func NewObject(v []Pair) Node {
    s := new(linkedPairs)
    s.FromSlice(v)
    return newObject(s)
}

func newObject(v *linkedPairs) Node {
    if v.size > _Threshold_Index {
        v.BuildIndex()
    }
    return Node{
        t: types.V_OBJECT,
        l: uint(v.Len()),
        p: unsafe.Pointer(v),
    }
}

func (self *Node) setObject(v *linkedPairs) {
    if v.size > _Threshold_Index {
        v.BuildIndex()
    }
    self.t = types.V_OBJECT
    self.l = uint(v.Len())
    self.p = unsafe.Pointer(v)
}

func (self *Node) parseRaw(full bool) {
    lock := self.lock()
    defer self.unlock()
    if !self.isRaw() {
        return
    }
    raw := self.toString()
    parser := NewParserObj(raw)
    var e types.ParsingError
    if full {
        parser.noLazy = true
        *self, e = parser.Parse()
    } else if lock {
        var n Node
        parser.noLazy = true
        parser.loadOnce = true
        n, e = parser.Parse()
        self.assign(n)
    } else {
        *self, e = parser.Parse()
    }
    if e != 0 {
        *self = *newSyntaxError(parser.syntaxError(e))
    }
}

func (self *Node) assign(n Node) {
    self.l = n.l
    self.p = n.p
    atomic.StoreInt64(&self.t, n.t)
}
