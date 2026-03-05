



//go:build aix || solaris

package unix

import (
	"unsafe"
)






func IoctlSetInt(fd int, req int, value int) error {
	return ioctl(fd, req, uintptr(value))
}





func IoctlSetPointerInt(fd int, req int, value int) error {
	v := int32(value)
	return ioctlPtr(fd, req, unsafe.Pointer(&v))
}




func IoctlSetWinsize(fd int, req int, value *Winsize) error {
	
	
	return ioctlPtr(fd, req, unsafe.Pointer(value))
}




func IoctlSetTermios(fd int, req int, value *Termios) error {
	
	return ioctlPtr(fd, req, unsafe.Pointer(value))
}






func IoctlGetInt(fd int, req int) (int, error) {
	var value int
	err := ioctlPtr(fd, req, unsafe.Pointer(&value))
	return value, err
}

func IoctlGetWinsize(fd int, req int) (*Winsize, error) {
	var value Winsize
	err := ioctlPtr(fd, req, unsafe.Pointer(&value))
	return &value, err
}

func IoctlGetTermios(fd int, req int) (*Termios, error) {
	var value Termios
	err := ioctlPtr(fd, req, unsafe.Pointer(&value))
	return &value, err
}
