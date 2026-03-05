

package caching

import (
    `sync`
    `sync/atomic`
    `unsafe`

    `github.com/bytedance/sonic/internal/rt`
)



const (
    _LoadFactor   = 0.5
    _InitCapacity = 4096    
)

type _ProgramMap struct {
    n uint64
    m uint32
    b []_ProgramEntry
}

type _ProgramEntry struct {
    vt *rt.GoType
    fn interface{}
}

func newProgramMap() *_ProgramMap {
    return &_ProgramMap {
        n: 0,
        m: _InitCapacity - 1,
        b: make([]_ProgramEntry, _InitCapacity),
    }
}

func (self *_ProgramMap) copy() *_ProgramMap {
    fork := &_ProgramMap{
        n: self.n,
        m: self.m,
        b: make([]_ProgramEntry, len(self.b)),
    }
    for i, f := range self.b {
        fork.b[i] = f
    }
    return fork
}

func (self *_ProgramMap) get(vt *rt.GoType) interface{} {
    i := self.m + 1
    p := vt.Hash & self.m

    
    for ; i > 0; i-- {
        if b := self.b[p]; b.vt == vt {
            return b.fn
        } else if b.vt == nil {
            break
        } else {
            p = (p + 1) & self.m
        }
    }

    
    return nil
}

func (self *_ProgramMap) add(vt *rt.GoType, fn interface{}) *_ProgramMap {
    p := self.copy()
    f := float64(atomic.LoadUint64(&p.n) + 1) / float64(p.m + 1)

    
    if f > _LoadFactor {
        p = p.rehash()
    }

    
    p.insert(vt, fn)
    return p
}

func (self *_ProgramMap) rehash() *_ProgramMap {
    c := (self.m + 1) << 1
    r := &_ProgramMap{m: c - 1, b: make([]_ProgramEntry, int(c))}

    
    for i := uint32(0); i <= self.m; i++ {
        if b := self.b[i]; b.vt != nil {
            r.insert(b.vt, b.fn)
        }
    }

    
    return r
}

func (self *_ProgramMap) insert(vt *rt.GoType, fn interface{}) {
    h := vt.Hash
    p := h & self.m

    
    for i := uint32(0); i <= self.m; i++ {
        if b := &self.b[p]; b.vt != nil {
            p += 1
            p &= self.m
        } else {
            b.vt = vt
            b.fn = fn
            atomic.AddUint64(&self.n, 1)
            return
        }
    }

    
    panic("no available slots")
}



type ProgramCache struct {
    m sync.Mutex
    p unsafe.Pointer
}

func CreateProgramCache() *ProgramCache {
    return &ProgramCache {
        m: sync.Mutex{},
        p: unsafe.Pointer(newProgramMap()),
    }
}

func (self *ProgramCache) Get(vt *rt.GoType) interface{} {
    return (*_ProgramMap)(atomic.LoadPointer(&self.p)).get(vt)
}

func (self *ProgramCache) Compute(vt *rt.GoType, compute func(*rt.GoType, ... interface{}) (interface{}, error), ex ...interface{}) (interface{}, error) {
    var err error
    var val interface{}

    
    self.m.Lock()
    defer self.m.Unlock()

    
    if val = self.Get(vt); val != nil {
        return val, nil
    }

    
    if val, err = compute(vt, ex...); err != nil {
        return nil, err
    }

    
    atomic.StorePointer(&self.p, unsafe.Pointer((*_ProgramMap)(atomic.LoadPointer(&self.p)).add(vt, val)))
    return val, nil
}
