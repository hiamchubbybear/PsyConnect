

package resolver

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type FieldOpts int
type OffsetType int

const (
    F_omitempty FieldOpts = 1 << iota
    F_stringize
    F_omitzero
)

const (
    F_offset OffsetType = iota
    F_deref
)

type Offset struct {
    Size uintptr
    Kind OffsetType
    Type reflect.Type
}

type FieldMeta struct {
    Name string
    Path []Offset
    Opts FieldOpts
    Type reflect.Type
    IsZero func(reflect.Value) bool
}

func (self *FieldMeta) String() string {
    var path []string
    var opts []string

    
    for _, off := range self.Path {
        if off.Kind == F_offset {
            path = append(path, fmt.Sprintf("%d", off.Size))
        } else {
            path = append(path, fmt.Sprintf("%d.(*%s)", off.Size, off.Type))
        }
    }

    
    if (self.Opts & F_stringize) != 0 {
        opts = append(opts, "string")
    }

    
    if (self.Opts & F_omitempty) != 0 {
        opts = append(opts, "omitempty")
    }

    
    return fmt.Sprintf(
        "{Field \"%s\" @ %s, opts=%s, type=%s}",
        self.Name,
        strings.Join(path, "."),
        strings.Join(opts, ","),
        self.Type,
    )
}

func (self *FieldMeta) optimize() {
    var n int
    var v uintptr

    
    for _, o := range self.Path {
        if v += o.Size; o.Kind == F_deref {
            self.Path[n].Size    = v
            self.Path[n].Type, v = o.Type, 0
            self.Path[n].Kind, n = F_deref, n + 1
        }
    }

    
    if v != 0 {
        self.Path[n].Size = v
        self.Path[n].Type = nil
        self.Path[n].Kind = F_offset
        n++
    }

    
    if n != 0 {
        self.Path = self.Path[:n]
    } else {
        self.Path = []Offset{{Kind: F_offset}}
    }
}

func resolveFields(vt reflect.Type) []FieldMeta {
    tfv := typeFields(vt)
    ret := []FieldMeta(nil)

    
    for _, fv := range tfv.list {
        
        ret = append(ret, FieldMeta{})
        fm := &ret[len(ret)-1]

        item := vt
        path := []Offset(nil)

        
        if fv.quoted {
            fm.Opts |= F_stringize
        }

        
        if fv.omitEmpty {
            fm.Opts |= F_omitempty
        }

        
        handleOmitZero(fv, fm)

        
        for _, i := range fv.index {
            kind := F_offset
            fval := item.Field(i)
            item  = fval.Type

            
            if item.Kind() == reflect.Ptr {
                kind = F_deref
                item = item.Elem()
            }

            
            path = append(path, Offset {
                Kind: kind,
                Type: item,
                Size: fval.Offset,
            })
        }

        
        idx := len(path) - 1
        fvt := path[idx].Type

        
        if path[idx].Kind == F_deref {
            fvt = reflect.PtrTo(fvt)
            path[idx].Kind = F_offset
        }

        fm.Type = fvt
        fm.Path = path
        fm.Name = fv.name
    }

    
    for i := range ret {
        ret[i].optimize()
    }

    
    return ret
}

var (
    fieldLock  = sync.RWMutex{}
    fieldCache = map[reflect.Type][]FieldMeta{}
)

func ResolveStruct(vt reflect.Type) []FieldMeta {
    var ok bool
    var fm []FieldMeta

    
    fieldLock.RLock()
    fm, ok = fieldCache[vt]
    fieldLock.RUnlock()

    
    if ok {
        return fm
    }

    
    fieldLock.Lock()
    defer fieldLock.Unlock()

    
    if fm, ok = fieldCache[vt]; ok {
        return fm
    }

    
    fm = resolveFields(vt)
    fieldCache[vt] = fm
    return fm
}
