

package alg

import (
	"encoding"
	"reflect"
	"strconv"
	"sync"
	"unsafe"

	"github.com/bytedance/sonic/internal/encoder/vars"
	"github.com/bytedance/sonic/internal/rt"
)

type _MapPair struct {
    k string  
    v unsafe.Pointer
    m [32]byte
}

type MapIterator struct {
    It rt.GoMapIterator     
    kv rt.GoSlice           
    ki int
}

var (
    iteratorPool = sync.Pool{}
    iteratorPair = rt.UnpackType(reflect.TypeOf(_MapPair{}))
)

func init() {
    if unsafe.Offsetof(MapIterator{}.It) != 0 {
        panic("_MapIterator.it is not the first field")
    }
}


func newIterator() *MapIterator {
    if v := iteratorPool.Get(); v == nil {
        return new(MapIterator)
    } else {
        return resetIterator(v.(*MapIterator))
    }
}

func resetIterator(p *MapIterator) *MapIterator {
    p.ki = 0
    p.It = rt.GoMapIterator{}
    p.kv.Len = 0
    return p
}

func (self *MapIterator) at(i int) *_MapPair {
    return (*_MapPair)(unsafe.Pointer(uintptr(self.kv.Ptr) + uintptr(i) * unsafe.Sizeof(_MapPair{})))
}

func (self *MapIterator) add() (p *_MapPair) {
    p = self.at(self.kv.Len)
    self.kv.Len++
    return
}

func (self *MapIterator) data() (p []_MapPair) {
    *(*rt.GoSlice)(unsafe.Pointer(&p)) = self.kv
    return
}

func (self *MapIterator) append(t *rt.GoType, k unsafe.Pointer, v unsafe.Pointer) (err error) {
    p := self.add()
    p.v = v

    
    if tk := t.Kind(); tk != reflect.String {
        return self.appendGeneric(p, t, tk, k)
    }

    
    p.k = *(*string)(k)
    return nil
}

func (self *MapIterator) appendGeneric(p *_MapPair, t *rt.GoType, v reflect.Kind, k unsafe.Pointer) error {
    switch v {
        case reflect.Int       : p.k = rt.Mem2Str(strconv.AppendInt(p.m[:0], int64(*(*int)(k)), 10))      ; return nil
        case reflect.Int8      : p.k = rt.Mem2Str(strconv.AppendInt(p.m[:0], int64(*(*int8)(k)), 10))     ; return nil
        case reflect.Int16     : p.k = rt.Mem2Str(strconv.AppendInt(p.m[:0], int64(*(*int16)(k)), 10))    ; return nil
        case reflect.Int32     : p.k = rt.Mem2Str(strconv.AppendInt(p.m[:0], int64(*(*int32)(k)), 10))    ; return nil
        case reflect.Int64     : p.k = rt.Mem2Str(strconv.AppendInt(p.m[:0], int64(*(*int64)(k)), 10))           ; return nil
        case reflect.Uint      : p.k = rt.Mem2Str(strconv.AppendUint(p.m[:0], uint64(*(*uint)(k)), 10))    ; return nil
        case reflect.Uint8     : p.k = rt.Mem2Str(strconv.AppendUint(p.m[:0], uint64(*(*uint8)(k)), 10))   ; return nil
        case reflect.Uint16    : p.k = rt.Mem2Str(strconv.AppendUint(p.m[:0], uint64(*(*uint16)(k)), 10))  ; return nil
        case reflect.Uint32    : p.k = rt.Mem2Str(strconv.AppendUint(p.m[:0], uint64(*(*uint32)(k)), 10))  ; return nil
        case reflect.Uint64    : p.k = rt.Mem2Str(strconv.AppendUint(p.m[:0], uint64(*(*uint64)(k)), 10))          ; return nil
        case reflect.Uintptr   : p.k = rt.Mem2Str(strconv.AppendUint(p.m[:0], uint64(*(*uintptr)(k)), 10)) ; return nil
        case reflect.Interface : return self.appendInterface(p, t, k)
        case reflect.Struct, reflect.Ptr : return self.appendConcrete(p, t, k)
        default                : panic("unexpected map key type")
    }
}

func (self *MapIterator) appendConcrete(p *_MapPair, t *rt.GoType, k unsafe.Pointer) (err error) {
    
    if !t.Indirect() {
        k = *(*unsafe.Pointer)(k)
    }
    eface := rt.GoEface{Value: k, Type: t}.Pack()
    out, err := eface.(encoding.TextMarshaler).MarshalText()
    if err != nil {
        return err
    }
    p.k = rt.Mem2Str(out)
    return
}

func (self *MapIterator) appendInterface(p *_MapPair, t *rt.GoType, k unsafe.Pointer) (err error) {
    if len(rt.IfaceType(t).Methods) == 0 {
        panic("unexpected map key type")
    } else if p.k, err = asText(k); err == nil {
        return nil
    } else {
        return
    }
}

func IteratorStop(p *MapIterator) {
    iteratorPool.Put(p)
}

func IteratorNext(p *MapIterator) {
    i := p.ki
    t := &p.It

    
    if i < 0 {
        rt.Mapiternext(t)
        return
    }

    
    if p.ki >= p.kv.Len {
        t.K = nil
        t.V = nil
        return
    }

    
    t.K = unsafe.Pointer(&p.at(p.ki).k)
    t.V = p.at(p.ki).v
    p.ki++
}

func IteratorStart(t *rt.GoMapType, m unsafe.Pointer, fv uint64) (*MapIterator, error) {
    it := newIterator()
    rt.Mapiterinit(t, m, &it.It)
    count := rt.Maplen(m)

    
    if count == 0 || (fv & (1<<BitSortMapKeys)) == 0 {
        it.ki = -1
        return it, nil
    }

    
    if count > it.kv.Cap {
        it.kv = rt.GrowSlice(iteratorPair, it.kv, count)
    }

    
    for ; it.It.K != nil; rt.Mapiternext(&it.It) {
        if err := it.append(t.Key, it.It.K, it.It.V); err != nil {
            IteratorStop(it)
            return nil, err
        }
    }

    
    if it.ki = 1; count > 1 {
        radixQsort(it.data(), 0, maxDepth(it.kv.Len))
    }

    
    it.It.V = it.at(0).v
    it.It.K = unsafe.Pointer(&it.at(0).k)
    return it, nil
}

func asText(v unsafe.Pointer) (string, error) {
	text := rt.AssertI2I(rt.UnpackType(vars.EncodingTextMarshalerType), *(*rt.GoIface)(v))
	r, e := (*(*encoding.TextMarshaler)(unsafe.Pointer(&text))).MarshalText()
	return rt.Mem2Str(r), e
}
