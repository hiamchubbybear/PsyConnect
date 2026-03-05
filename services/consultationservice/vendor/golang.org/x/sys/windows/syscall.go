



//go:build windows



















package windows 

import (
	"bytes"
	"strings"
	"syscall"
	"unsafe"
)




func ByteSliceFromString(s string) ([]byte, error) {
	if strings.IndexByte(s, 0) != -1 {
		return nil, syscall.EINVAL
	}
	a := make([]byte, len(s)+1)
	copy(a, s)
	return a, nil
}




func BytePtrFromString(s string) (*byte, error) {
	a, err := ByteSliceFromString(s)
	if err != nil {
		return nil, err
	}
	return &a[0], nil
}



func ByteSliceToString(s []byte) string {
	if i := bytes.IndexByte(s, 0); i != -1 {
		s = s[:i]
	}
	return string(s)
}




func BytePtrToString(p *byte) string {
	if p == nil {
		return ""
	}
	if *p == 0 {
		return ""
	}

	
	n := 0
	for ptr := unsafe.Pointer(p); *(*byte)(ptr) != 0; n++ {
		ptr = unsafe.Pointer(uintptr(ptr) + 1)
	}

	return string(unsafe.Slice(p, n))
}



var _zero uintptr

func (ts *Timespec) Unix() (sec int64, nsec int64) {
	return int64(ts.Sec), int64(ts.Nsec)
}

func (tv *Timeval) Unix() (sec int64, nsec int64) {
	return int64(tv.Sec), int64(tv.Usec) * 1000
}

func (ts *Timespec) Nano() int64 {
	return int64(ts.Sec)*1e9 + int64(ts.Nsec)
}

func (tv *Timeval) Nano() int64 {
	return int64(tv.Sec)*1e9 + int64(tv.Usec)*1000
}
