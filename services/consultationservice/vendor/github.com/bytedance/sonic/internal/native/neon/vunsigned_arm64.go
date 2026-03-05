


package neon

import (
    `github.com/bytedance/sonic/internal/native/types`
)

//go:nosplit
func vunsigned(s *string, p *int, v *types.JsonState) {
    __vunsigned(s, p, v)
}

//go:nosplit
//go:noescape

func __vunsigned(s *string, p *int, v *types.JsonState)
