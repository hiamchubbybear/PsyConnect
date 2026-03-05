


package neon

import (
    `github.com/bytedance/sonic/internal/native/types`
)

//go:nosplit
func vnumber(s *string, p *int, v *types.JsonState) {
    __vnumber(s, p, v)
}

//go:nosplit
//go:noescape

func __vnumber(s *string, p *int, v *types.JsonState)


