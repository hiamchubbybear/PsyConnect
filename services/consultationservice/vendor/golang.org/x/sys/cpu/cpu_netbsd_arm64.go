



package cpu

import (
	"syscall"
	"unsafe"
)




const (
	_CTL_QUERY = -2

	_SYSCTL_VERS_1 = 0x1000000
)

var _zero uintptr

func sysctl(mib []int32, old *byte, oldlen *uintptr, new *byte, newlen uintptr) (err error) {
	var _p0 unsafe.Pointer
	if len(mib) > 0 {
		_p0 = unsafe.Pointer(&mib[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	_, _, errno := syscall.Syscall6(
		syscall.SYS___SYSCTL,
		uintptr(_p0),
		uintptr(len(mib)),
		uintptr(unsafe.Pointer(old)),
		uintptr(unsafe.Pointer(oldlen)),
		uintptr(unsafe.Pointer(new)),
		uintptr(newlen))
	if errno != 0 {
		return errno
	}
	return nil
}

type sysctlNode struct {
	Flags          uint32
	Num            int32
	Name           [32]int8
	Ver            uint32
	__rsvd         uint32
	Un             [16]byte
	_sysctl_size   [8]byte
	_sysctl_func   [8]byte
	_sysctl_parent [8]byte
	_sysctl_desc   [8]byte
}

func sysctlNodes(mib []int32) ([]sysctlNode, error) {
	var olen uintptr

	
	
	mib = append(mib, _CTL_QUERY)
	qnode := sysctlNode{Flags: _SYSCTL_VERS_1}
	qp := (*byte)(unsafe.Pointer(&qnode))
	sz := unsafe.Sizeof(qnode)
	if err := sysctl(mib, nil, &olen, qp, sz); err != nil {
		return nil, err
	}

	
	nodes := make([]sysctlNode, olen/sz)
	np := (*byte)(unsafe.Pointer(&nodes[0]))
	if err := sysctl(mib, np, &olen, qp, sz); err != nil {
		return nil, err
	}

	return nodes, nil
}

func nametomib(name string) ([]int32, error) {
	
	var parts []string
	last := 0
	for i := 0; i < len(name); i++ {
		if name[i] == '.' {
			parts = append(parts, name[last:i])
			last = i + 1
		}
	}
	parts = append(parts, name[last:])

	mib := []int32{}
	
	for partno, part := range parts {
		nodes, err := sysctlNodes(mib)
		if err != nil {
			return nil, err
		}
		for _, node := range nodes {
			n := make([]byte, 0)
			for i := range node.Name {
				if node.Name[i] != 0 {
					n = append(n, byte(node.Name[i]))
				}
			}
			if string(n) == part {
				mib = append(mib, int32(node.Num))
				break
			}
		}
		if len(mib) != partno+1 {
			return nil, err
		}
	}

	return mib, nil
}


type aarch64SysctlCPUID struct {
	midr      uint64 
	revidr    uint64 
	mpidr     uint64 
	aa64dfr0  uint64 
	aa64dfr1  uint64 
	aa64isar0 uint64 
	aa64isar1 uint64 
	aa64mmfr0 uint64 
	aa64mmfr1 uint64 
	aa64mmfr2 uint64 
	aa64pfr0  uint64 
	aa64pfr1  uint64 
	aa64zfr0  uint64 
	mvfr0     uint32 
	mvfr1     uint32 
	mvfr2     uint32 
	pad       uint32
	clidr     uint64 
	ctr       uint64 
}

func sysctlCPUID(name string) (*aarch64SysctlCPUID, error) {
	mib, err := nametomib(name)
	if err != nil {
		return nil, err
	}

	out := aarch64SysctlCPUID{}
	n := unsafe.Sizeof(out)
	_, _, errno := syscall.Syscall6(
		syscall.SYS___SYSCTL,
		uintptr(unsafe.Pointer(&mib[0])),
		uintptr(len(mib)),
		uintptr(unsafe.Pointer(&out)),
		uintptr(unsafe.Pointer(&n)),
		uintptr(0),
		uintptr(0))
	if errno != 0 {
		return nil, errno
	}
	return &out, nil
}

func doinit() {
	cpuid, err := sysctlCPUID("machdep.cpu0.cpu_id")
	if err != nil {
		setMinimalFeatures()
		return
	}
	parseARM64SystemRegisters(cpuid.aa64isar0, cpuid.aa64isar1, cpuid.aa64pfr0)

	Initialized = true
}
