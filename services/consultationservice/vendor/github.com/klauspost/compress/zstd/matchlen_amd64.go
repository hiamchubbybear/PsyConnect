//go:build amd64 && !appengine && !noasm && gc
// +build amd64,!appengine,!noasm,gc




package zstd







//go:noescape
func matchLen(a []byte, b []byte) int
