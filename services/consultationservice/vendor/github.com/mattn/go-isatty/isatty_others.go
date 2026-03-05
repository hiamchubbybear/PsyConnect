//go:build (appengine || js || nacl || tinygo || wasm) && !windows
// +build appengine js nacl tinygo wasm
// +build !windows

package isatty



func IsTerminal(fd uintptr) bool {
	return false
}



func IsCygwinTerminal(fd uintptr) bool {
	return false
}
