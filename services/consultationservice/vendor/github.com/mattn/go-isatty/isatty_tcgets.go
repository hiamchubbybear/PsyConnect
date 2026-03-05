//go:build (linux || aix || zos) && !appengine && !tinygo
// +build linux aix zos
// +build !appengine
// +build !tinygo

package isatty

import "golang.org/x/sys/unix"


func IsTerminal(fd uintptr) bool {
	_, err := unix.IoctlGetTermios(int(fd), unix.TCGETS)
	return err == nil
}



func IsCygwinTerminal(fd uintptr) bool {
	return false
}
