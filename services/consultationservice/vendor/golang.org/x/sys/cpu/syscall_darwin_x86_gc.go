






//go:build darwin && amd64 && gc

package cpu

import (
	"syscall"
	"unsafe"
)

type _C_int int32


func darwinOSRelease(release *[256]byte) error {
	
	const (
		CTL_KERN       = 0x1
		KERN_OSRELEASE = 0x2
	)

	mib := []_C_int{CTL_KERN, KERN_OSRELEASE}
	n := unsafe.Sizeof(*release)

	return sysctl(mib, &release[0], &n, nil, 0)
}

type Errno = syscall.Errno

var _zero uintptr 


func sysctl(mib []_C_int, old *byte, oldlen *uintptr, new *byte, newlen uintptr) error {
	var _p0 unsafe.Pointer
	if len(mib) > 0 {
		_p0 = unsafe.Pointer(&mib[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	if _, _, err := syscall_syscall6(
		libc_sysctl_trampoline_addr,
		uintptr(_p0),
		uintptr(len(mib)),
		uintptr(unsafe.Pointer(old)),
		uintptr(unsafe.Pointer(oldlen)),
		uintptr(unsafe.Pointer(new)),
		uintptr(newlen),
	); err != 0 {
		return err
	}

	return nil
}

var libc_sysctl_trampoline_addr uintptr


func darwinSysctlEnabled(name []byte) bool {
	out := int32(0)
	nout := unsafe.Sizeof(out)
	if ret := sysctlbyname(&name[0], (*byte)(unsafe.Pointer(&out)), &nout, nil, 0); ret != nil {
		return false
	}
	return out > 0
}

//go:cgo_import_dynamic libc_sysctl sysctl "/usr/lib/libSystem.B.dylib"

var libc_sysctlbyname_trampoline_addr uintptr


func sysctlbyname(name *byte, old *byte, oldlen *uintptr, new *byte, newlen uintptr) error {
	if _, _, err := syscall_syscall6(
		libc_sysctlbyname_trampoline_addr,
		uintptr(unsafe.Pointer(name)),
		uintptr(unsafe.Pointer(old)),
		uintptr(unsafe.Pointer(oldlen)),
		uintptr(unsafe.Pointer(new)),
		uintptr(newlen),
		0,
	); err != 0 {
		return err
	}

	return nil
}

//go:cgo_import_dynamic libc_sysctlbyname sysctlbyname "/usr/lib/libSystem.B.dylib"


func syscall_syscall6(fn, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err Errno)

//go:linkname syscall_syscall6 syscall.syscall6
