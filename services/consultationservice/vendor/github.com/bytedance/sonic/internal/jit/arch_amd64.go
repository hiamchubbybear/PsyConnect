

package jit

import (
	"unsafe"

	"github.com/twitchyliquid64/golang-asm/asm/arch"
	"github.com/twitchyliquid64/golang-asm/obj"
)

var (
    _AC = arch.Set("amd64")
)

func As(op string) obj.As {
    if ret, ok := _AC.Instructions[op]; ok {
        return ret
    } else {
        panic("invalid instruction: " + op)
    }
}

func ImmPtr(imm unsafe.Pointer) obj.Addr {
    return obj.Addr {
        Type   : obj.TYPE_CONST,
        Offset : int64(uintptr(imm)),
    }
}

func Imm(imm int64) obj.Addr {
    return obj.Addr {
        Type   : obj.TYPE_CONST,
        Offset : imm,
    }
}

func Reg(reg string) obj.Addr {
    if ret, ok := _AC.Register[reg]; ok {
        return obj.Addr{Reg: ret, Type: obj.TYPE_REG}
    } else {
        panic("invalid register name: " + reg)
    }
}

func Ptr(reg obj.Addr, offs int64) obj.Addr {
    return obj.Addr {
        Reg    : reg.Reg,
        Type   : obj.TYPE_MEM,
        Offset : offs,
    }
}

func Sib(reg obj.Addr, idx obj.Addr, scale int16, offs int64) obj.Addr {
    return obj.Addr {
        Reg    : reg.Reg,
        Index  : idx.Reg,
        Scale  : scale,
        Type   : obj.TYPE_MEM,
        Offset : offs,
    }
}
