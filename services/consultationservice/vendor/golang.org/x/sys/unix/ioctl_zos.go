



//go:build zos && s390x

package unix

import (
	"runtime"
	"unsafe"
)






func IoctlSetInt(fd int, req int, value int) error {
	return ioctl(fd, req, uintptr(value))
}




func IoctlSetWinsize(fd int, req int, value *Winsize) error {
	
	
	return ioctlPtr(fd, req, unsafe.Pointer(value))
}




func IoctlSetTermios(fd int, req int, value *Termios) error {
	if (req != TCSETS) && (req != TCSETSW) && (req != TCSETSF) {
		return ENOSYS
	}
	err := Tcsetattr(fd, int(req), value)
	runtime.KeepAlive(value)
	return err
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
	if req != TCGETS {
		return &value, ENOSYS
	}
	err := Tcgetattr(fd, &value)
	return &value, err
}
