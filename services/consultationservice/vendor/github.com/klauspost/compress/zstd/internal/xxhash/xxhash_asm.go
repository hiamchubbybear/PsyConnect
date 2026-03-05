//go:build (amd64 || arm64) && !appengine && gc && !purego && !noasm
// +build amd64 arm64
// +build !appengine
// +build gc
// +build !purego
// +build !noasm

package xxhash



//go:noescape
func Sum64(b []byte) uint64

//go:noescape
func writeBlocks(s *Digest, b []byte) int
