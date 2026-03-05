//go:build (amd64 && go1.16 && !go1.25) || (arm64 && go1.20 && !go1.25)
// +build amd64,go1.16,!go1.25 arm64,go1.20,!go1.25



package alg

import (
	"runtime"
	"unsafe"

	"github.com/bytedance/sonic/internal/native"
	"github.com/bytedance/sonic/internal/native/types"
	"github.com/bytedance/sonic/internal/rt"
)






func Valid(data []byte) (ok bool, start int) {
    n := len(data)
    if n == 0 {
        return false, -1
    }
    s := rt.Mem2Str(data)
    p := 0
    m := types.NewStateMachine()
    ret := native.ValidateOne(&s, &p, m, 0)
    types.FreeStateMachine(m)

    if ret < 0 {
        return false, p-1
    }

    
    for ;p < n; p++ {
        if (types.SPACE_MASK & (1 << data[p])) == 0 {
            return false, p
        }
    }

    return true, ret
}

var typeByte = rt.UnpackEface(byte(0)).Type

//go:nocheckptr
func Quote(buf []byte, val string, double bool) []byte {
	if len(val) == 0 {
		if double {
			return append(buf, `"\"\""`...)
		}
		return append(buf, `""`...)
	}

	if double {
		buf = append(buf, `"\"`...)
	} else {
		buf = append(buf, `"`...)
	}
	sp := rt.IndexChar(val, 0)
	nb := len(val)
	b := (*rt.GoSlice)(unsafe.Pointer(&buf))

	
	for nb > 0 {
		
		dp := unsafe.Pointer(uintptr(b.Ptr) + uintptr(b.Len))
		dn := b.Cap - b.Len
		
		opts := uint64(0)
		if double {
			opts = types.F_DOUBLE_UNQUOTE
		}
		ret := native.Quote(sp, nb, dp, &dn, opts)
		
		b.Len += dn

		
		if ret >= 0 {
			break
		}

		
		*b = rt.GrowSlice(typeByte, *b, b.Cap*2)
		
		ret = ^ret
		
		nb -= ret
		sp = unsafe.Pointer(uintptr(sp) + uintptr(ret))
	}

	runtime.KeepAlive(buf)
	runtime.KeepAlive(sp)
	if double {
		buf = append(buf, `\""`...)
	} else {
		buf = append(buf, `"`...)
	}

	return buf
}

func HtmlEscape(dst []byte, src []byte) []byte {
	var sidx int

	dst = append(dst, src[:0]...) 
	sbuf := (*rt.GoSlice)(unsafe.Pointer(&src))
	dbuf := (*rt.GoSlice)(unsafe.Pointer(&dst))

	
	if cap(dst)-len(dst) < len(src)+types.BufPaddingSize {
		cap := len(src)*3/2 + types.BufPaddingSize
		*dbuf = rt.GrowSlice(typeByte, *dbuf, cap)
	}

	for sidx < sbuf.Len {
		sp := rt.Add(sbuf.Ptr, uintptr(sidx))
		dp := rt.Add(dbuf.Ptr, uintptr(dbuf.Len))

		sn := sbuf.Len - sidx
		dn := dbuf.Cap - dbuf.Len
		nb := native.HTMLEscape(sp, sn, dp, &dn)

		
		if dbuf.Len += dn; nb >= 0 {
			break
		}

		
		sidx += ^nb
		*dbuf = rt.GrowSlice(typeByte, *dbuf, dbuf.Cap*2)
	}
	return dst
}

func F64toa(buf []byte, v float64) ([]byte) {
	if v == 0 {
		return append(buf, '0')
	}
	buf = rt.GuardSlice2(buf, 64)
	ret := native.F64toa((*byte)(rt.IndexByte(buf, len(buf))), v)
	if ret > 0 {
		return buf[:len(buf)+ret]
	} else {
		return buf
	}
}

func F32toa(buf []byte, v float32) ([]byte) {
	if v == 0 {
		return append(buf, '0')
	}
	buf = rt.GuardSlice2(buf, 64)
	ret := native.F32toa((*byte)(rt.IndexByte(buf, len(buf))), v)
	if ret > 0 {
		return buf[:len(buf)+ret]
	} else {
		return buf
	}
}

func I64toa(buf []byte, v int64) ([]byte) {
	buf = rt.GuardSlice2(buf, 32)
	ret := native.I64toa((*byte)(rt.IndexByte(buf, len(buf))), v)
	if ret > 0 {
		return buf[:len(buf)+ret]
	} else {
		return buf
	}
}

func U64toa(buf []byte, v uint64) ([]byte) {
	buf = rt.GuardSlice2(buf, 32)
	ret := native.U64toa((*byte)(rt.IndexByte(buf, len(buf))), v)
	if ret > 0 {
		return buf[:len(buf)+ret]
	} else {
		return buf
	}
}

