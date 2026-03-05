

package jit

import (
    `reflect`
    `unsafe`

    `github.com/bytedance/sonic/internal/rt`
    `github.com/twitchyliquid64/golang-asm/obj`
)

func Func(f interface{}) obj.Addr {
    if p := rt.UnpackEface(f); p.Type.Kind() != reflect.Func {
        panic("f is not a function")
    } else {
        return Imm(*(*int64)(p.Value))
    }
}

func Type(t reflect.Type) obj.Addr {
    return Gtype(rt.UnpackType(t))
}

func Itab(i *rt.GoType, t reflect.Type) obj.Addr {
    return Imm(int64(uintptr(unsafe.Pointer(rt.GetItab(rt.IfaceType(i), rt.UnpackType(t), false)))))
}

func Gitab(i *rt.GoItab) obj.Addr {
    return Imm(int64(uintptr(unsafe.Pointer(i))))
}

func Gtype(t *rt.GoType) obj.Addr {
    return Imm(int64(uintptr(unsafe.Pointer(t))))
}
