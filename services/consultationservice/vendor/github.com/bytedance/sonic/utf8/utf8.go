

package utf8

import (
	`runtime`

    `github.com/bytedance/sonic/internal/rt`
    `github.com/bytedance/sonic/internal/native/types`
    `github.com/bytedance/sonic/internal/native`
)


func CorrectWith(dst []byte, src []byte, repl string) []byte {
    sstr := rt.Mem2Str(src)
    sidx := 0

    
    m := types.NewStateMachine()
    m.Sp = 0 

    for sidx < len(sstr) {
        scur  := sidx
        ecode := native.ValidateUTF8(&sstr, &sidx, m)

        if m.Sp != 0 {
            if m.Sp > len(sstr) {
                panic("numbers of invalid utf8 exceed the string len!")
            }
        }
        
        for i := 0; i < m.Sp; i++ {
            ipos := m.Vt[i] 
            dst  = append(dst, sstr[scur:ipos]...)
            dst  = append(dst, repl...)
            scur = m.Vt[i] + 1
        }
        
        dst = append(dst, sstr[scur:sidx]...)

        
        if ecode != 0 {
            m.Sp = 0
        }
    }

    types.FreeStateMachine(m)
    return dst
}


func Validate(src []byte) bool {
	if src == nil {
		return true
	}
    return ValidateString(rt.Mem2Str(src))
}


func ValidateString(src string) bool {
	if src == "" {
		return true
	}
    ret := native.ValidateUTF8Fast(&src) == 0
	runtime.KeepAlive(src)
	return ret
}
