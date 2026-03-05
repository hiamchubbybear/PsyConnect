

package jitdec

import (
    `encoding`
    `encoding/json`
    `unsafe`

    `github.com/bytedance/sonic/internal/native`
    `github.com/bytedance/sonic/internal/rt`
)

func decodeTypedPointer(s string, i int, vt *rt.GoType, vp unsafe.Pointer, sb *_Stack, fv uint64) (int, error) {
    if fn, err := findOrCompile(vt); err != nil {
        return 0, err
    } else {
        rt.MoreStack(_FP_size + _VD_size + native.MaxFrameSize)
        ret, err := fn(s, i, vp, sb, fv, "", nil)
        return ret, err
    }
}

func decodeJsonUnmarshaler(vv interface{}, s string) error {
    return vv.(json.Unmarshaler).UnmarshalJSON(rt.Str2Mem(s))
}


type MismatchQuotedError struct {}

func (*MismatchQuotedError) Error() string {
    return "mismatch quoted"
}

func decodeJsonUnmarshalerQuoted(vv interface{}, s string) error {
    if len(s) < 2 || s[0] != '"' || s[len(s)-1] != '"' {
        return &MismatchQuotedError{}
    }
    return vv.(json.Unmarshaler).UnmarshalJSON(rt.Str2Mem(s[1:len(s)-1]))
}

func decodeTextUnmarshaler(vv interface{}, s string) error {
    return vv.(encoding.TextUnmarshaler).UnmarshalText(rt.Str2Mem(s))
}
