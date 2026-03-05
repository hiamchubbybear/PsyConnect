



package neon

import (
    

    `github.com/bytedance/sonic/internal/native/types`
)

//go:nosplit
//go:noescape

func __skip_one(s *string, p *int, m *types.StateMachine, flags uint64) (ret int)

//go:nosplit
func skip_one(s *string, p *int, m *types.StateMachine, flags uint64) (ret int) {
    return __skip_one(s, p, m, flags)
}
