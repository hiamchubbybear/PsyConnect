



//go:build go1.21 && (aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris || zos)

package unix

import (
	"syscall"
	"unsafe"
)

//go:linkname runtime_getAuxv runtime.getAuxv
func runtime_getAuxv() []uintptr





func Auxv() ([][2]uintptr, error) {
	vec := runtime_getAuxv()
	vecLen := len(vec)

	if vecLen == 0 {
		return nil, syscall.ENOENT
	}

	if vecLen%2 != 0 {
		return nil, syscall.EINVAL
	}

	result := make([]uintptr, vecLen)
	copy(result, vec)
	return unsafe.Slice((*[2]uintptr)(unsafe.Pointer(&result[0])), vecLen/2), nil
}
