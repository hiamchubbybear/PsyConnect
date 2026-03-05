












//go:build aix || darwin || openbsd || freebsd || netbsd || dragonfly || zos
// +build aix darwin openbsd freebsd netbsd dragonfly zos

package afero

import (
	"syscall"
)

const BADFD = syscall.EBADF
