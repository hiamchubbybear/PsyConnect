package kafka

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/segmentio/kafka-go/protocol/endtxn"
)


type EndTxnRequest struct {
	
	Addr net.Addr

	
	TransactionalID string

	
	ProducerID int

	
	ProducerEpoch int

	
	Committed bool
}


type EndTxnResponse struct {
	
	Throttle time.Duration

	
	
	
	Error error
}


func (c *Client) EndTxn(ctx context.Context, req *EndTxnRequest) (*EndTxnResponse, error) {
	m, err := c.roundTrip(ctx, req.Addr, &endtxn.Request{
		TransactionalID: req.TransactionalID,
		ProducerID:      int64(req.ProducerID),
		ProducerEpoch:   int16(req.ProducerEpoch),
		Committed:       req.Committed,
	})
	if err != nil {
		return nil, fmt.Errorf("kafka.(*Client).EndTxn: %w", err)
	}

	r := m.(*endtxn.Response)

	res := &EndTxnResponse{
		Throttle: makeDuration(r.ThrottleTimeMs),
		Error:    makeError(r.ErrorCode, ""),
	}

	return res, nil
}
