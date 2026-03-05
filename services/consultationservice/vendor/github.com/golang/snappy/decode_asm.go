



// +build !appengine
// +build gc
// +build !noasm
// +build amd64 arm64

package snappy



//go:noescape
func decode(dst, src []byte) int
