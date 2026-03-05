

package unix

import "unsafe"


func PtraceGetRegSetArm64(pid, addr int, regsout *PtraceRegsArm64) error {
	iovec := Iovec{(*byte)(unsafe.Pointer(regsout)), uint64(unsafe.Sizeof(*regsout))}
	return ptracePtr(PTRACE_GETREGSET, pid, uintptr(addr), unsafe.Pointer(&iovec))
}


func PtraceSetRegSetArm64(pid, addr int, regs *PtraceRegsArm64) error {
	iovec := Iovec{(*byte)(unsafe.Pointer(regs)), uint64(unsafe.Sizeof(*regs))}
	return ptracePtr(PTRACE_SETREGSET, pid, uintptr(addr), unsafe.Pointer(&iovec))
}
