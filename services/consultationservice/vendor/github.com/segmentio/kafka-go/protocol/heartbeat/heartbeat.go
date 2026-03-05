package heartbeat

import "github.com/segmentio/kafka-go/protocol"

func init() {
	protocol.Register(&Request{}, &Response{})
}


type Request struct {
	
	
	_ struct{} `kafka:"min=v4,max=v4,tag"`

	GroupID         string `kafka:"min=v0,max=v4"`
	GenerationID    int32  `kafka:"min=v0,max=v4"`
	MemberID        string `kafka:"min=v0,max=v4"`
	GroupInstanceID string `kafka:"min=v3,max=v4,nullable"`
}

func (r *Request) ApiKey() protocol.ApiKey {
	return protocol.Heartbeat
}

type Response struct {
	
	
	_ struct{} `kafka:"min=v4,max=v4,tag"`

	ErrorCode      int16 `kafka:"min=v0,max=v4"`
	ThrottleTimeMs int32 `kafka:"min=v1,max=v4"`
}

func (r *Response) ApiKey() protocol.ApiKey {
	return protocol.Heartbeat
}
