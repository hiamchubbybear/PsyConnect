









//go:build aix && gccgo

package cpu

import (
	"syscall"
)


func gccgoGetsystemcfg(label uint32) (r uint64)

func callgetsystemcfg(label int) (r1 uintptr, e1 syscall.Errno) {
	r1 = uintptr(gccgoGetsystemcfg(uint32(label)))
	e1 = syscall.GetErrno()
	return
}
