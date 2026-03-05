


package neon

import (
    `github.com/bytedance/sonic/internal/native/types`
)

//go:nosplit
func vsigned(s *string, p *int, v *types.JsonState) {
    __vsigned(s, p, v)
}

//go:nosplit
//go:noescape

func __vsigned(s *string, p *int, v *types.JsonState)
