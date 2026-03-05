



//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris || zos

package unix


func (fds *FdSet) Set(fd int) {
	fds.Bits[fd/NFDBITS] |= (1 << (uintptr(fd) % NFDBITS))
}


func (fds *FdSet) Clear(fd int) {
	fds.Bits[fd/NFDBITS] &^= (1 << (uintptr(fd) % NFDBITS))
}


func (fds *FdSet) IsSet(fd int) bool {
	return fds.Bits[fd/NFDBITS]&(1<<(uintptr(fd)%NFDBITS)) != 0
}


func (fds *FdSet) Zero() {
	for i := range fds.Bits {
		fds.Bits[i] = 0
	}
}
