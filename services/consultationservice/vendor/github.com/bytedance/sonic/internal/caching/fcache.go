

package caching

import (
    `strings`
    `unsafe`

    `github.com/bytedance/sonic/internal/rt`
)

type FieldMap struct {
    N uint64
    b unsafe.Pointer
    m map[string]int
}

type FieldEntry struct {
    ID   int
    Name string
    Hash uint64
}

const (
    FieldMap_N     = int64(unsafe.Offsetof(FieldMap{}.N))
    FieldMap_b     = int64(unsafe.Offsetof(FieldMap{}.b))
	FieldEntrySize = int64(unsafe.Sizeof(FieldEntry{}))
)

func newBucket(n int) unsafe.Pointer {
    v := make([]FieldEntry, n)
    return (*rt.GoSlice)(unsafe.Pointer(&v)).Ptr
}

func CreateFieldMap(n int) *FieldMap {
    return &FieldMap {
        N: uint64(n * 2),
        b: newBucket(n * 2),    
        m: make(map[string]int, n * 2),
    }
}

func (self *FieldMap) At(p uint64) *FieldEntry {
    off := uintptr(p) * uintptr(FieldEntrySize)
    return (*FieldEntry)(unsafe.Pointer(uintptr(self.b) + off))
}




func (self *FieldMap) Get(name string) int {
    h := StrHash(name)
    p := h % self.N
    s := self.At(p)

    
    for s.Hash != 0 {
        if s.Hash == h && s.Name == name {
            return s.ID
        } else {
            p = (p + 1) % self.N
            s = self.At(p)
        }
    }

    
    return -1
}

func (self *FieldMap) Set(name string, i int) {
    h := StrHash(name)
    p := h % self.N
    s := self.At(p)

    
    for s.Hash != 0 {
        p = (p + 1) % self.N
        s = self.At(p)
    }

    
    s.ID   = i
    s.Hash = h
    s.Name = name

    
    key := strings.ToLower(name)
    if v, ok := self.m[key]; !ok || i < v {
        self.m[key] = i
    }
}

func (self *FieldMap) GetCaseInsensitive(name string) int {
    if i, ok := self.m[strings.ToLower(name)]; ok {
        return i
    } else {
        return -1
    }
}
