//go:build go1.21 && !go1.25
// +build go1.21,!go1.25















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

    _F_gcWriteBarrier2 = jit.Func(rt.GcWriteBarrier2)
)

func (self *Assembler) WritePtr(i int, ptr obj.Addr, old obj.Addr) {
    if old.Reg == x86.REG_AX || old.Index == x86.REG_AX {
        panic("rec contains AX!")
    }
    self.Emit("MOVQ", _V_writeBarrier, _BX)
    self.Emit("CMPL", jit.Ptr(_BX, 0), jit.Imm(0))
    self.Sjmp("JE", "_no_writeBarrier" + strconv.Itoa(i) + "_{n}")
    self.xsave(_SP_q)
    self.Emit("MOVQ", _F_gcWriteBarrier2, _BX)  
    self.Rjmp("CALL", _BX)  
    self.Emit("MOVQ", ptr, jit.Ptr(_SP_q, 0))
    self.Emit("MOVQ", old, _AX)
    self.Emit("MOVQ", _AX, jit.Ptr(_SP_q, 8))
    self.xload(_SP_q)  
    self.Link("_no_writeBarrier" + strconv.Itoa(i) + "_{n}")
    self.Emit("MOVQ", ptr, old)
}
