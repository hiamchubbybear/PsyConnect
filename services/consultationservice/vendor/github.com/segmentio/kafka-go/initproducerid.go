package kafka

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/segmentio/kafka-go/protocol/initproducerid"
)


type InitProducerIDRequest struct {
	
	Addr net.Addr

	
	TransactionalID string

	
	TransactionTimeoutMs int

	
	
	
	ProducerID int

	
	
	
	
	ProducerEpoch int
}


type ProducerSession struct {
	
	ProducerID int

	
	ProducerEpoch int
}


type InitProducerIDResponse struct {
	
	Producer *ProducerSession

	
	Throttle time.Duration

	
	
	
	
	Error error
}



func (c *Client) InitProducerID(ctx context.Context, req *InitProducerIDRequest) (*InitProducerIDResponse, error) {
	m, err := c.roundTrip(ctx, req.Addr, &initproducerid.Request{
		TransactionalID:      req.TransactionalID,
		TransactionTimeoutMs: int32(req.TransactionTimeoutMs),
		ProducerID:           int64(req.ProducerID),
		ProducerEpoch:        int16(req.ProducerEpoch),
	})
	if err != nil {
		return nil, fmt.Errorf("kafka.(*Client).InitProducerId: %w", err)
	}

	res := m.(*initproducerid.Response)

	return &InitProducerIDResponse{
		Producer: &ProducerSession{
			ProducerID:    int(res.ProducerID),
			ProducerEpoch: int(res.ProducerEpoch),
		},
		Throttle: makeDuration(res.ThrottleTimeMs),
		Error:    makeError(res.ErrorCode, ""),
	}, nil
}
