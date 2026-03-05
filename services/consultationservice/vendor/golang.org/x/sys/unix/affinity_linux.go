





package unix

import (
	"math/bits"
	"unsafe"
)

const cpuSetSize = _CPU_SETSIZE / _NCPUBITS


type CPUSet [cpuSetSize]cpuMask

func schedAffinity(trap uintptr, pid int, set *CPUSet) error {
	_, _, e := RawSyscall(trap, uintptr(pid), uintptr(unsafe.Sizeof(*set)), uintptr(unsafe.Pointer(set)))
	if e != 0 {
		return errnoErr(e)
	}
	return nil
}



func SchedGetaffinity(pid int, set *CPUSet) error {
	return schedAffinity(SYS_SCHED_GETAFFINITY, pid, set)
}



func SchedSetaffinity(pid int, set *CPUSet) error {
	return schedAffinity(SYS_SCHED_SETAFFINITY, pid, set)
}


func (s *CPUSet) Zero() {
	for i := range s {
		s[i] = 0
	}
}

func cpuBitsIndex(cpu int) int {
	return cpu / _NCPUBITS
}

func cpuBitsMask(cpu int) cpuMask {
	return cpuMask(1 << (uint(cpu) % _NCPUBITS))
}


func (s *CPUSet) Set(cpu int) {
	i := cpuBitsIndex(cpu)
	if i < len(s) {
		s[i] |= cpuBitsMask(cpu)
	}
}


func (s *CPUSet) Clear(cpu int) {
	i := cpuBitsIndex(cpu)
	if i < len(s) {
		s[i] &^= cpuBitsMask(cpu)
	}
}


func (s *CPUSet) IsSet(cpu int) bool {
	i := cpuBitsIndex(cpu)
	if i < len(s) {
		return s[i]&cpuBitsMask(cpu) != 0
	}
	return false
}


func (s *CPUSet) Count() int {
	c := 0
	for _, b := range s {
		c += bits.OnesCount64(uint64(b))
	}
	return c
}
