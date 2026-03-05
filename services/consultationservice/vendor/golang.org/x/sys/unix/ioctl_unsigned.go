



//go:build darwin || dragonfly || freebsd || hurd || linux || netbsd || openbsd

package unix

import (
	"unsafe"
)






func IoctlSetInt(fd int, req uint, value int) error {
	return ioctl(fd, req, uintptr(value))
}





func IoctlSetPointerInt(fd int, req uint, value int) error {
	v := int32(value)
	return ioctlPtr(fd, req, unsafe.Pointer(&v))
}




func IoctlSetWinsize(fd int, req uint, value *Winsize) error {
	
	
	return ioctlPtr(fd, req, unsafe.Pointer(value))
}




func IoctlSetTermios(fd int, req uint, value *Termios) error {
	
	return ioctlPtr(fd, req, unsafe.Pointer(value))
}






func IoctlGetInt(fd int, req uint) (int, error) {
	var value int
	err := ioctlPtr(fd, req, unsafe.Pointer(&value))
	return value, err
}

func IoctlGetWinsize(fd int, req uint) (*Winsize, error) {
	var value Winsize
	err := ioctlPtr(fd, req, unsafe.Pointer(&value))
	return &value, err
}

func IoctlGetTermios(fd int, req uint) (*Termios, error) {
	var value Termios
	err := ioctlPtr(fd, req, unsafe.Pointer(&value))
	return &value, err
}
