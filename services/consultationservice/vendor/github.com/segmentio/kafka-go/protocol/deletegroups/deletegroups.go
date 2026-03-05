package deletegroups

import "github.com/segmentio/kafka-go/protocol"

func init() {
	protocol.Register(&Request{}, &Response{})
}

type Request struct {
	
	
	_ struct{} `kafka:"min=v2,max=v2,tag"`

	GroupIDs []string `kafka:"min=v0,max=v2"`
}

func (r *Request) Group() string {
	
	if len(r.GroupIDs) > 0 {
		return r.GroupIDs[0]
	}
	return ""
}

func (r *Request) ApiKey() protocol.ApiKey { return protocol.DeleteGroups }

var (
	_ protocol.GroupMessage = (*Request)(nil)
)

type Response struct {
	
	
	_ struct{} `kafka:"min=v2,max=v2,tag"`

	ThrottleTimeMs int32           `kafka:"min=v0,max=v2"`
	Responses      []ResponseGroup `kafka:"min=v0,max=v2"`
}

func (r *Response) ApiKey() protocol.ApiKey { return protocol.DeleteGroups }

type ResponseGroup struct {
	GroupID   string `kafka:"min=v0,max=v2"`
	ErrorCode int16  `kafka:"min=v0,max=v2"`
}
