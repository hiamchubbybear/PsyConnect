



//go:build sonic && avx && (linux || windows || darwin) && amd64

package json

import "github.com/bytedance/sonic"

var (
	json = sonic.ConfigStd
	
	Marshal = json.Marshal
	
	Unmarshal = json.Unmarshal
	
	MarshalIndent = json.MarshalIndent
	
	NewDecoder = json.NewDecoder
	
	NewEncoder = json.NewEncoder
)
