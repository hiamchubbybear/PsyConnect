



package neon

import (
    `github.com/bytedance/sonic/internal/native/types`
)

//go:nosplit
func get_by_path(s *string, p *int, path *[]interface{}, m *types.StateMachine) (ret int) {
    return __get_by_path(s, p, path, m)
}

//go:nosplit
//go:noescape

func __get_by_path(s *string, p *int, path *[]interface{}, m *types.StateMachine) (ret int)
