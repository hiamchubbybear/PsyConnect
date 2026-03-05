

package abi

import (
    `fmt`
    `reflect`
    `sort`
    `strings`

    `github.com/bytedance/sonic/loader/internal/rt`
)

type FunctionLayout struct {
    FP   uint32
    Args []Parameter
    Rets []Parameter
}

func (self FunctionLayout) String() string {
    return self.formatFn()
}

func (self FunctionLayout) ArgSize() uint32 {
    size := uintptr(0) 
    for _, arg := range self.Args {
        size += arg.Type.Size()
    }
    return uint32(size)
}

type slot struct {
    p bool
    m uint32
}

func (self FunctionLayout) StackMap() *rt.StackMap {
    var st []slot
    var mb rt.StackMapBuilder

    
    for _, v := range self.Args {
        st = append(st, slot {
            m: v.Mem,
            p: v.IsPointer,
        })
    }

    
    for _, v := range self.Rets {
        if !v.InRegister {
            st = append(st, slot {
                m: v.Mem,
                p: v.IsPointer,
            })
        }
    }

    
    sort.Slice(st, func(i int, j int) bool {
        return st[i].m < st[j].m
    })

    
    for _, v := range st {
        mb.AddField(v.p)
    }

    
    return mb.Build()
}

func (self FunctionLayout) formatFn() string {
    fp := self.FP
    return fmt.Sprintf("\n%#04x\nRets:\n%s\nArgs:\n%s", fp, self.formatSeq(self.Rets, &fp), self.formatSeq(self.Args, &fp))
}

func (self FunctionLayout) formatSeq(v []Parameter, fp *uint32) string {
    nb := len(v)
    mm := make([]string, 0, len(v))

    
    for i := nb-1; i >=0; i-- {
        *fp -= PtrSize
        mm = append(mm, fmt.Sprintf("%#04x %s", *fp, v[i].String()))
    }

    
    return strings.Join(mm, "\n")
}

type Frame struct {
    desc      *FunctionLayout
    locals    []bool
    ccall     bool
}

func NewFrame(desc *FunctionLayout, locals []bool, ccall bool) Frame {
    fr := Frame{}
    fr.desc = desc
    fr.locals = locals
    fr.ccall = ccall
    return fr
}

func (self *Frame) String() string {
    out := self.desc.String()

    off := -8
    out += fmt.Sprintf("\n%#4x [Return PC]", off)
    off -= 8
    out += fmt.Sprintf("\n%#4x [RBP]", off)
    off -= 8

    for _, v := range ReservedRegs(self.ccall) {
        out += fmt.Sprintf("\n%#4x [%v]", off, v)
        off -= PtrSize
    }

    for _, b := range self.locals {
        out += fmt.Sprintf("\n%#4x [%v]", off, b)
        off -= PtrSize
    }

    return out
}

func (self *Frame) Prev() uint32 {
    return self.Size() + PtrSize
}

func (self *Frame) Size() uint32 {
    return uint32(self.Offs() + PtrSize)
}

func (self *Frame) Offs() uint32 {
    return uint32(len(ReservedRegs(self.ccall)) * PtrSize + len(self.locals)*PtrSize)
}

func (self *Frame) ArgPtrs() *rt.StackMap {
    return self.desc.StackMap()
}

func (self *Frame) LocalPtrs() *rt.StackMap {
    var m rt.StackMapBuilder
    for _, b := range self.locals {
        m.AddFields(len(ReservedRegs(self.ccall)), b)
    }
    return m.Build()
}

func alignUp(n uint32, a int) uint32 {
    return (uint32(n) + uint32(a) - 1) &^ (uint32(a) - 1)
}

func isPointer(vt reflect.Type) bool {
    switch vt.Kind() {
        case reflect.Bool          : fallthrough
        case reflect.Int           : fallthrough
        case reflect.Int8          : fallthrough
        case reflect.Int16         : fallthrough
        case reflect.Int32         : fallthrough
        case reflect.Int64         : fallthrough
        case reflect.Uint          : fallthrough
        case reflect.Uint8         : fallthrough
        case reflect.Uint16        : fallthrough
        case reflect.Uint32        : fallthrough
        case reflect.Uint64        : fallthrough
        case reflect.Float32       : fallthrough
        case reflect.Float64       : fallthrough
        case reflect.Uintptr       : return false
        case reflect.Chan          : fallthrough
        case reflect.Func          : fallthrough
        case reflect.Map           : fallthrough
        case reflect.Ptr           : fallthrough
        case reflect.UnsafePointer : return true
        case reflect.Complex64     : fallthrough
        case reflect.Complex128    : fallthrough
        case reflect.Array         : fallthrough
        case reflect.Struct        : panic("abi: unsupported types")
        default                    : panic("abi: invalid value type")
    }
}