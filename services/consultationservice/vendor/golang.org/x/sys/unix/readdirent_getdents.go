



//go:build aix || dragonfly || freebsd || linux || netbsd || openbsd

package unix


func ReadDirent(fd int, buf []byte) (n int, err error) {
	return Getdents(fd, buf)
}
