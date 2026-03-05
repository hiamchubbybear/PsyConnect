

package unquote

import (
    `unsafe`
    `runtime`

    `github.com/bytedance/sonic/internal/native`
    `github.com/bytedance/sonic/internal/native/types`
    `github.com/bytedance/sonic/internal/rt`
)



func String(s string) (ret string, err types.ParsingError) {
    mm := make([]byte, 0, len(s))
    err = intoBytesUnsafe(s, &mm, true)
    ret = rt.Mem2Str(mm)
    return
}


func IntoBytes(s string, m *[]byte) types.ParsingError {
    if cap(*m) < len(s) {
        return types.ERR_EOF
    } else {
        return intoBytesUnsafe(s, m, true)
    }
}



func _String(s string, replace bool) (ret string, err error) {
    mm := make([]byte, 0, len(s))
    err = intoBytesUnsafe(s, &mm, replace)
    ret = rt.Mem2Str(mm)
    return
}

func intoBytesUnsafe(s string, m *[]byte, replace bool) types.ParsingError {
    pos := -1
    slv := (*rt.GoSlice)(unsafe.Pointer(m))
    str := (*rt.GoString)(unsafe.Pointer(&s))

    flags := uint64(0)
    if replace {
        
        flags |= types.F_UNICODE_REPLACE
    }

    ret := native.Unquote(str.Ptr, str.Len, slv.Ptr, &pos, flags)

    
    if ret < 0 {
        return types.ParsingError(-ret)
    }

    
    slv.Len = ret
    runtime.KeepAlive(s)
    return 0
}



