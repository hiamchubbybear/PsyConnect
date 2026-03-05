package kafka

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/segmentio/kafka-go/protocol/addoffsetstotxn"
)


type AddOffsetsToTxnRequest struct {
	
	Addr net.Addr

	
	TransactionalID string

	
	
	ProducerID int

	
	ProducerEpoch int

	
	GroupID string
}


type AddOffsetsToTxnResponse struct {
	
	Throttle time.Duration

	
	
	
	
	
	Error error
}


func (c *Client) AddOffsetsToTxn(
	ctx context.Context,
	req *AddOffsetsToTxnRequest,
) (*AddOffsetsToTxnResponse, error) {
	m, err := c.roundTrip(ctx, req.Addr, &addoffsetstotxn.Request{
		TransactionalID: req.TransactionalID,
		ProducerID:      int64(req.ProducerID),
		ProducerEpoch:   int16(req.ProducerEpoch),
		GroupID:         req.GroupID,
	})
	if err != nil {
		return nil, fmt.Errorf("kafka.(*Client).AddOffsetsToTxn: %w", err)
	}

	r := m.(*addoffsetstotxn.Response)

	res := &AddOffsetsToTxnResponse{
		Throttle: makeDuration(r.ThrottleTimeMs),
		Error:    makeError(r.ErrorCode, ""),
	}

	return res, nil
}
