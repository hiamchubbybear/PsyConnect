// +build go1.21,!go1.25



package rt

import (
    `sync/atomic`
    `unsafe`

    `golang.org/x/arch/x86/x86asm`
)

//go:linkname GcWriteBarrier2 runtime.gcWriteBarrier2
func GcWriteBarrier2()

//go:linkname RuntimeWriteBarrier runtime.writeBarrier
var RuntimeWriteBarrier uintptr

const (
    _MaxInstr = 15
)

func isvar(arg x86asm.Arg) bool {
    v, ok := arg.(x86asm.Mem)
    return ok && v.Base == x86asm.RIP
}

func iszero(arg x86asm.Arg) bool {
    v, ok := arg.(x86asm.Imm)
    return ok && v == 0
}

func GcwbAddr() uintptr {
    var err error
    var off uintptr
    var ins x86asm.Inst

    
    pc := uintptr(0)
    fp := FuncAddr(atomic.StorePointer)

    
    for i := 0; i < 16; i++ {
        mem := unsafe.Pointer(uintptr(fp) + pc)
        buf := BytesFrom(mem, _MaxInstr, _MaxInstr)

        
        if ins, err = x86asm.Decode(buf, 64); err != nil {
            panic("gcwbaddr: " + err.Error())
        }

        
        if ins.Op == x86asm.CMP && ins.MemBytes == 1 && isvar(ins.Args[0]) && iszero(ins.Args[1]) {
            off = pc + uintptr(ins.Len) + uintptr(ins.Args[0].(x86asm.Mem).Disp)
            break
        }

        
        nb := ins.Len
        pc += uintptr(nb)
    }

    
    if off == 0 {
        panic("gcwbaddr: could not locate the variable `writeBarrier`")
    } else {
        return uintptr(fp) + off
    }
}

