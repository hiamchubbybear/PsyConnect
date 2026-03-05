



package cpu

import (
	"syscall"
	"unsafe"
)










































const (
	
	riscv_HWPROBE_KEY_IMA_EXT_0   = 0x4
	riscv_HWPROBE_IMA_C           = 0x2
	riscv_HWPROBE_IMA_V           = 0x4
	riscv_HWPROBE_EXT_ZBA         = 0x8
	riscv_HWPROBE_EXT_ZBB         = 0x10
	riscv_HWPROBE_EXT_ZBS         = 0x20
	riscv_HWPROBE_KEY_CPUPERF_0   = 0x5
	riscv_HWPROBE_MISALIGNED_FAST = 0x3
	riscv_HWPROBE_MISALIGNED_MASK = 0x7
)

const (
	
	sys_RISCV_HWPROBE = 258
)


type riscvHWProbePairs struct {
	key   int64
	value uint64
}

const (
	
	hwcap_RISCV_ISA_C = 1 << ('C' - 'A')
)

func doinit() {
	
	
	
	

	pairs := []riscvHWProbePairs{
		{riscv_HWPROBE_KEY_IMA_EXT_0, 0},
		{riscv_HWPROBE_KEY_CPUPERF_0, 0},
	}

	
	if riscvHWProbe(pairs, 0) {
		if pairs[0].key != -1 {
			v := uint(pairs[0].value)
			RISCV64.HasC = isSet(v, riscv_HWPROBE_IMA_C)
			RISCV64.HasV = isSet(v, riscv_HWPROBE_IMA_V)
			RISCV64.HasZba = isSet(v, riscv_HWPROBE_EXT_ZBA)
			RISCV64.HasZbb = isSet(v, riscv_HWPROBE_EXT_ZBB)
			RISCV64.HasZbs = isSet(v, riscv_HWPROBE_EXT_ZBS)
		}
		if pairs[1].key != -1 {
			v := pairs[1].value & riscv_HWPROBE_MISALIGNED_MASK
			RISCV64.HasFastMisaligned = v == riscv_HWPROBE_MISALIGNED_FAST
		}
	}

	
	

	if !RISCV64.HasC {
		RISCV64.HasC = isSet(hwCap, hwcap_RISCV_ISA_C)
	}
}

func isSet(hwc uint, value uint) bool {
	return hwc&value != 0
}






func riscvHWProbe(pairs []riscvHWProbePairs, flags uint) bool {
	var _zero uintptr
	var p0 unsafe.Pointer
	if len(pairs) > 0 {
		p0 = unsafe.Pointer(&pairs[0])
	} else {
		p0 = unsafe.Pointer(&_zero)
	}

	_, _, e1 := syscall.Syscall6(sys_RISCV_HWPROBE, uintptr(p0), uintptr(len(pairs)), uintptr(0), uintptr(0), uintptr(flags), 0)
	return e1 == 0
}
