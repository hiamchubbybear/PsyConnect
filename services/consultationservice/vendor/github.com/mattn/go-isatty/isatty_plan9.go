//go:build plan9
// +build plan9

package isatty

import (
	"syscall"
)


func IsTerminal(fd uintptr) bool {
	path, err := syscall.Fd2path(int(fd))
	if err != nil {
		return false
	}
	return path == "/dev/cons" || path == "/mnt/term/dev/cons"
}



func IsCygwinTerminal(fd uintptr) bool {
	return false
}
