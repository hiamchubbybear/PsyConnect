//go:build go1.17 && !go1.25
// +build go1.17,!go1.25



package api

import (
	"github.com/bytedance/sonic/internal/envs"
	"github.com/bytedance/sonic/internal/decoder/jitdec"
	"github.com/bytedance/sonic/internal/decoder/optdec"
)

var (
	pretouchImpl = jitdec.Pretouch
	decodeImpl = jitdec.Decode
) 

 func init() {
	if envs.UseOptDec {
		pretouchImpl = optdec.Pretouch
		decodeImpl = optdec.Decode
	}
 }
