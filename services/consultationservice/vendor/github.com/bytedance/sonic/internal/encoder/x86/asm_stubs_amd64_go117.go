//go:build go1.17 && !go1.21
// +build go1.17,!go1.21















package x86

import (
	"strconv"
	"unsafe"

	"github.com/bytedance/sonic/internal/jit"
	"github.com/bytedance/sonic/internal/rt"
	"github.com/twitchyliquid64/golang-asm/obj"
	"github.com/twitchyliquid64/golang-asm/obj/x86"
)

var (
    _V_writeBarrier = jit.Imm(int64(uintptr(unsafe.Pointer(&rt.RuntimeWriteBarrier))))

    _F_gcWriteBarrierAX = jit.Func(rt.GcWriteBarrierAX)
)

func (self *Assembler) WritePtr(i int, ptr obj.Addr, rec obj.Addr) {
    if rec.Reg == x86.REG_AX || rec.Index == x86.REG_AX {
        panic("rec contains AX!")
    }
    self.Emit("MOVQ", _V_writeBarrier, _BX)
    self.Emit("CMPL", jit.Ptr(_BX, 0), jit.Imm(0))
    self.Sjmp("JE", "_no_writeBarrier" + strconv.Itoa(i) + "_{n}")
    self.xsave(_DI)
    self.Emit("MOVQ", ptr, _AX)
    self.Emit("LEAQ", rec, _DI)
    self.Emit("MOVQ", _F_gcWriteBarrierAX, _BX)  
    self.Rjmp("CALL", _BX)  
    self.xload(_DI)  
    self.Sjmp("JMP", "_end_writeBarrier" + strconv.Itoa(i) + "_{n}")
    self.Link("_no_writeBarrier" + strconv.Itoa(i) + "_{n}")
    self.Emit("MOVQ", ptr, rec)
    self.Link("_end_writeBarrier" + strconv.Itoa(i) + "_{n}")
}
