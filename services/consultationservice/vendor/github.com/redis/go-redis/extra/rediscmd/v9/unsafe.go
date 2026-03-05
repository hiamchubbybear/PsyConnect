//go:build !appengine
// +build !appengine

package rediscmd

import "unsafe"


func String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}


func Bytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}
