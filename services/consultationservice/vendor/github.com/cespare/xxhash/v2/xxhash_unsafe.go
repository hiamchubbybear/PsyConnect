//go:build !appengine
// +build !appengine




package xxhash

import (
	"unsafe"
)


























func Sum64String(s string) uint64 {
	b := *(*[]byte)(unsafe.Pointer(&sliceHeader{s, len(s)}))
	return Sum64(b)
}



func (d *Digest) WriteString(s string) (n int, err error) {
	d.Write(*(*[]byte)(unsafe.Pointer(&sliceHeader{s, len(s)})))
	
	
	
	return len(s), nil
}



type sliceHeader struct {
	s   string
	cap int
}
