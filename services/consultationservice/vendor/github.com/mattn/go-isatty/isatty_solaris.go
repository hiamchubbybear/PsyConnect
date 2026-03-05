//go:build solaris && !appengine
// +build solaris,!appengine

package isatty

import (
	"golang.org/x/sys/unix"
)



func IsTerminal(fd uintptr) bool {
	_, err := unix.IoctlGetTermio(int(fd), unix.TCGETA)
	return err == nil
}



func IsCygwinTerminal(fd uintptr) bool {
	return false
}
