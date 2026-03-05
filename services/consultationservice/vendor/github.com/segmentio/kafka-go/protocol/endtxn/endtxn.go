package endtxn

import "github.com/segmentio/kafka-go/protocol"

func init() {
	protocol.Register(&Request{}, &Response{})
}

type Request struct {
	
	
	_ struct{} `kafka:"min=v3,max=v3,tag"`

	TransactionalID string `kafka:"min=v0,max=v2|min=v3,max=v3,compact"`
	ProducerID      int64  `kafka:"min=v0,max=v3"`
	ProducerEpoch   int16  `kafka:"min=v0,max=v3"`
	Committed       bool   `kafka:"min=v0,max=v3"`
}

func (r *Request) ApiKey() protocol.ApiKey { return protocol.EndTxn }

func (r *Request) Transaction() string { return r.TransactionalID }

var _ protocol.TransactionalMessage = (*Request)(nil)

type Response struct {
	
	
	_ struct{} `kafka:"min=v3,max=v3,tag"`

	ThrottleTimeMs int32 `kafka:"min=v0,max=v3"`
	ErrorCode      int16 `kafka:"min=v0,max=v3"`
}

func (r *Response) ApiKey() protocol.ApiKey { return protocol.EndTxn }
