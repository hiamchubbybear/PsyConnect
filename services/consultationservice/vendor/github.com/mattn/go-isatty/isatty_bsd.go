//go:build (darwin || freebsd || openbsd || netbsd || dragonfly || hurd) && !appengine && !tinygo
// +build darwin freebsd openbsd netbsd dragonfly hurd
// +build !appengine
// +build !tinygo

package isatty

import "golang.org/x/sys/unix"


func IsTerminal(fd uintptr) bool {
	_, err := unix.IoctlGetTermios(int(fd), unix.TIOCGETA)
	return err == nil
}



func IsCygwinTerminal(fd uintptr) bool {
	return false
}
