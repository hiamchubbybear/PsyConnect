// +build go1.17,!go1.21















package jitdec

import (
    `strconv`
    `unsafe`

    `github.com/bytedance/sonic/internal/rt`
    `github.com/bytedance/sonic/internal/jit`
    `github.com/twitchyliquid64/golang-asm/obj`
    `github.com/twitchyliquid64/golang-asm/obj/x86`
)

var (
    _V_writeBarrier = jit.Imm(int64(uintptr(unsafe.Pointer(&rt.RuntimeWriteBarrier))))

    _F_gcWriteBarrierAX = jit.Func(rt.GcWriteBarrierAX)
)

func (self *_Assembler) WritePtrAX(i int, rec obj.Addr, saveDI bool) {
    self.Emit("MOVQ", _V_writeBarrier, _R9)
    self.Emit("CMPL", jit.Ptr(_R9, 0), jit.Imm(0))
    self.Sjmp("JE", "_no_writeBarrier" + strconv.Itoa(i) + "_{n}")
    if saveDI {
        self.save(_DI)
    }
    self.Emit("LEAQ", rec, _DI)
    self.call(_F_gcWriteBarrierAX)  
    if saveDI {
        self.load(_DI)
    }    
    self.Sjmp("JMP", "_end_writeBarrier" + strconv.Itoa(i) + "_{n}")
    self.Link("_no_writeBarrier" + strconv.Itoa(i) + "_{n}")
    self.Emit("MOVQ", _AX, rec)
    self.Link("_end_writeBarrier" + strconv.Itoa(i) + "_{n}")
}

func (self *_Assembler) WriteRecNotAX(i int, ptr obj.Addr, rec obj.Addr, saveDI bool, saveAX bool) {
    if rec.Reg == x86.REG_AX || rec.Index == x86.REG_AX {
        panic("rec contains AX!")
    }
    self.Emit("MOVQ", _V_writeBarrier, _R9)
    self.Emit("CMPL", jit.Ptr(_R9, 0), jit.Imm(0))
    self.Sjmp("JE", "_no_writeBarrier" + strconv.Itoa(i) + "_{n}")
    if saveAX {
        self.Emit("XCHGQ", ptr, _AX)
    } else {
        self.Emit("MOVQ", ptr, _AX)
    }
    if saveDI {
        self.save(_DI)
    }
    self.Emit("LEAQ", rec, _DI)
    self.call(_F_gcWriteBarrierAX) 
    if saveDI {
        self.load(_DI)
    } 
    if saveAX {
        self.Emit("XCHGQ", ptr, _AX)
    }    
    self.Sjmp("JMP", "_end_writeBarrier" + strconv.Itoa(i) + "_{n}")
    self.Link("_no_writeBarrier" + strconv.Itoa(i) + "_{n}")
    self.Emit("MOVQ", ptr, rec)
    self.Link("_end_writeBarrier" + strconv.Itoa(i) + "_{n}")
}


func (self *_ValueDecoder) WritePtrAX(i int, rec obj.Addr, saveDI bool) {
    self.Emit("MOVQ", _V_writeBarrier, _R9)
    self.Emit("CMPL", jit.Ptr(_R9, 0), jit.Imm(0))
    self.Sjmp("JE", "_no_writeBarrier" + strconv.Itoa(i) + "_{n}")
    if saveDI {
        self.save(_DI)
    }
    self.Emit("LEAQ", rec, _DI)
    self.call(_F_gcWriteBarrierAX)
    if saveDI {
        self.load(_DI)
    }    
    self.Sjmp("JMP", "_end_writeBarrier" + strconv.Itoa(i) + "_{n}")
    self.Link("_no_writeBarrier" + strconv.Itoa(i) + "_{n}")
    self.Emit("MOVQ", _AX, rec)
    self.Link("_end_writeBarrier" + strconv.Itoa(i) + "_{n}")
}

func (self *_ValueDecoder) WriteRecNotAX(i int, ptr obj.Addr, rec obj.Addr, saveDI bool) {
    if rec.Reg == x86.REG_AX || rec.Index == x86.REG_AX {
        panic("rec contains AX!")
    }
    self.Emit("MOVQ", _V_writeBarrier, _AX)
    self.Emit("CMPL", jit.Ptr(_AX, 0), jit.Imm(0))
    self.Sjmp("JE", "_no_writeBarrier" + strconv.Itoa(i) + "_{n}")
    self.Emit("MOVQ", ptr, _AX)
    if saveDI {
        self.save(_DI)
    }
    self.Emit("LEAQ", rec, _DI)
    self.call(_F_gcWriteBarrierAX)
    if saveDI {
        self.load(_DI)
    }    
    self.Sjmp("JMP", "_end_writeBarrier" + strconv.Itoa(i) + "_{n}")
    self.Link("_no_writeBarrier" + strconv.Itoa(i) + "_{n}")
    self.Emit("MOVQ", ptr, rec)
    self.Link("_end_writeBarrier" + strconv.Itoa(i) + "_{n}")
}
