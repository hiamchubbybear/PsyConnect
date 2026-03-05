



//go:build arm && gc && linux

package unix

import "syscall"



func seek(fd int, offset int64, whence int) (newoffset int64, err syscall.Errno)
