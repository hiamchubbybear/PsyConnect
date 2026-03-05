







package arch

import (
	"github.com/twitchyliquid64/golang-asm/obj"
	"github.com/twitchyliquid64/golang-asm/obj/riscv"
)




func IsRISCV64AMO(op obj.As) bool {
	switch op {
	case riscv.ASCW, riscv.ASCD, riscv.AAMOSWAPW, riscv.AAMOSWAPD, riscv.AAMOADDW, riscv.AAMOADDD,
		riscv.AAMOANDW, riscv.AAMOANDD, riscv.AAMOORW, riscv.AAMOORD, riscv.AAMOXORW, riscv.AAMOXORD,
		riscv.AAMOMINW, riscv.AAMOMIND, riscv.AAMOMINUW, riscv.AAMOMINUD,
		riscv.AAMOMAXW, riscv.AAMOMAXD, riscv.AAMOMAXUW, riscv.AAMOMAXUD:
		return true
	}
	return false
}
