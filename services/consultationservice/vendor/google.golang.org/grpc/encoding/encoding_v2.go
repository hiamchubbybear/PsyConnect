

package encoding

import (
	"strings"

	"google.golang.org/grpc/mem"
)




type CodecV2 interface {
	
	
	
	Marshal(v any) (out mem.BufferSlice, err error)
	
	
	
	
	Unmarshal(data mem.BufferSlice, v any) error
	
	
	
	Name() string
}


















func RegisterCodecV2(codec CodecV2) {
	if codec == nil {
		panic("cannot register a nil CodecV2")
	}
	if codec.Name() == "" {
		panic("cannot register CodecV2 with empty string result for Name()")
	}
	contentSubtype := strings.ToLower(codec.Name())
	registeredCodecs[contentSubtype] = codec
}





func GetCodecV2(contentSubtype string) CodecV2 {
	c, _ := registeredCodecs[contentSubtype].(CodecV2)
	return c
}
