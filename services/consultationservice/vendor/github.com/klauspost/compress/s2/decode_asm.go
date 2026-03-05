




//go:build (amd64 || arm64) && !appengine && gc && !noasm
// +build amd64 arm64
// +build !appengine
// +build gc
// +build !noasm

package s2



//go:noescape
func s2Decode(dst, src []byte) int
