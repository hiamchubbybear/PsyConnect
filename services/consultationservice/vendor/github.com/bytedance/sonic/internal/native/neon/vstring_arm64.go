


package neon

import (
    `github.com/bytedance/sonic/internal/native/types`
)

//go:nosplit
func vstring(s *string, p *int, v *types.JsonState, flags uint64) {
    __vstring(s, p, v, flags)
}

//go:nosplit
//go:noescape

func __vstring(s *string, p *int, v *types.JsonState, flags uint64)
