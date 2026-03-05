


//go:build go1.9 && !go1.20 && !safe && !codec.safe && !appengine
// +build go1.9,!go1.20,!safe,!codec.safe,!appengine

package codec

import (
	_ "runtime" 
	"unsafe"
)

//go:linkname growslice runtime.growslice
//go:noescape
func growslice(typ unsafe.Pointer, old unsafeSlice, num int) unsafeSlice
