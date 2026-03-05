



package neon

import (
    

    `github.com/bytedance/sonic/internal/native/types`
)

//go:nosplit
func skip_object(s *string, p *int, m *types.StateMachine, flags uint64) (ret int) {
    return __skip_object(s, p, m, flags)
}

//go:nosplit
//go:noescape

func __skip_object(s *string, p *int, m *types.StateMachine, flags uint64) (ret int)
