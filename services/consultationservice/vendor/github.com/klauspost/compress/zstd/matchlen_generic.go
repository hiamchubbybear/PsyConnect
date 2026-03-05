//go:build !amd64 || appengine || !gc || noasm
// +build !amd64 appengine !gc noasm




package zstd

import (
	"encoding/binary"
	"math/bits"
)



func matchLen(a, b []byte) (n int) {
	for ; len(a) >= 8 && len(b) >= 8; a, b = a[8:], b[8:] {
		diff := binary.LittleEndian.Uint64(a) ^ binary.LittleEndian.Uint64(b)
		if diff != 0 {
			return n + bits.TrailingZeros64(diff)>>3
		}
		n += 8
	}

	for i := range a {
		if a[i] != b[i] {
			break
		}
		n++
	}
	return n

}
