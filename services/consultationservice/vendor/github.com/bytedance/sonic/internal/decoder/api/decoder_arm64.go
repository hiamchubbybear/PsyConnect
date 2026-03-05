// +build go1.17,!go1.25



package api

import (
	`github.com/bytedance/sonic/internal/decoder/optdec`
	`github.com/bytedance/sonic/internal/envs`
)

var (
	pretouchImpl = optdec.Pretouch
	decodeImpl = optdec.Decode
)


func init() {
    
	envs.EnableOptDec()
	envs.EnableFastMap()
}


