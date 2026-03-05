



package neon

import (
    

    `github.com/bytedance/sonic/internal/native/types`
)

//go:nosplit
func validate_one(s *string, p *int, m *types.StateMachine) (ret int) {
    return __validate_one(s, p, m)
}

//go:nosplit
//go:noescape

func __validate_one(s *string, p *int, m *types.StateMachine) (ret int)
