

package grpc

import (
	"google.golang.org/grpc/encoding"
	_ "google.golang.org/grpc/encoding/proto" 
	"google.golang.org/grpc/mem"
)





type baseCodec interface {
	Marshal(v any) (mem.BufferSlice, error)
	Unmarshal(data mem.BufferSlice, v any) error
}






func getCodec(name string) encoding.CodecV2 {
	if codecV1 := encoding.GetCodec(name); codecV1 != nil {
		return newCodecV1Bridge(codecV1)
	}

	return encoding.GetCodecV2(name)
}

func newCodecV0Bridge(c Codec) baseCodec {
	return codecV0Bridge{codec: c}
}

func newCodecV1Bridge(c encoding.Codec) encoding.CodecV2 {
	return codecV1Bridge{
		codecV0Bridge: codecV0Bridge{codec: c},
		name:          c.Name(),
	}
}

var _ baseCodec = codecV0Bridge{}

type codecV0Bridge struct {
	codec interface {
		Marshal(v any) ([]byte, error)
		Unmarshal(data []byte, v any) error
	}
}

func (c codecV0Bridge) Marshal(v any) (mem.BufferSlice, error) {
	data, err := c.codec.Marshal(v)
	if err != nil {
		return nil, err
	}
	return mem.BufferSlice{mem.SliceBuffer(data)}, nil
}

func (c codecV0Bridge) Unmarshal(data mem.BufferSlice, v any) (err error) {
	return c.codec.Unmarshal(data.Materialize(), v)
}

var _ encoding.CodecV2 = codecV1Bridge{}

type codecV1Bridge struct {
	codecV0Bridge
	name string
}

func (c codecV1Bridge) Name() string {
	return c.name
}






type Codec interface {
	
	Marshal(v any) ([]byte, error)
	
	Unmarshal(data []byte, v any) error
	
	
	String() string
}
