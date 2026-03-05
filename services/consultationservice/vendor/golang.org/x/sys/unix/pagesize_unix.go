



//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris || zos



package unix

import "syscall"

func Getpagesize() int {
	return syscall.Getpagesize()
}
