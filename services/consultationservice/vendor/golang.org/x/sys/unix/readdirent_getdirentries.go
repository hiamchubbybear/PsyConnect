



//go:build darwin || zos

package unix

import "unsafe"


func ReadDirent(fd int, buf []byte) (n int, err error) {
	
	
	
	
	var base = (*uintptr)(unsafe.Pointer(new(uint64)))
	return Getdirentries(fd, buf, base)
}
