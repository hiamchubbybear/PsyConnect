package kafka

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/segmentio/kafka-go/protocol"
	produceAPI "github.com/segmentio/kafka-go/protocol/produce"
	"github.com/segmentio/kafka-go/protocol/rawproduce"
)



type RawProduceRequest struct {
	
	Addr net.Addr

	
	Topic string

	
	Partition int

	
	RequiredAcks RequiredAcks

	
	
	
	
	MessageVersion int

	
	
	TransactionalID string

	
	RawRecords protocol.RawRecordSet
}








func (c *Client) RawProduce(ctx context.Context, req *RawProduceRequest) (*ProduceResponse, error) {
	m, err := c.roundTrip(ctx, req.Addr, &rawproduce.Request{
		TransactionalID: req.TransactionalID,
		Acks:            int16(req.RequiredAcks),
		Timeout:         c.timeoutMs(ctx, defaultProduceTimeout),
		Topics: []rawproduce.RequestTopic{{
			Topic: req.Topic,
			Partitions: []rawproduce.RequestPartition{{
				Partition: int32(req.Partition),
				RecordSet: req.RawRecords,
			}},
		}},
	})

	switch {
	case err == nil:
	case errors.Is(err, protocol.ErrNoRecord):
		return new(ProduceResponse), nil
	default:
		return nil, fmt.Errorf("kafka.(*Client).RawProduce: %w", err)
	}

	if req.RequiredAcks == RequireNone {
		return nil, nil
	}

	res := m.(*produceAPI.Response)
	if len(res.Topics) == 0 {
		return nil, fmt.Errorf("kafka.(*Client).RawProduce: %w", protocol.ErrNoTopic)
	}
	topic := &res.Topics[0]
	if len(topic.Partitions) == 0 {
		return nil, fmt.Errorf("kafka.(*Client).RawProduce: %w", protocol.ErrNoPartition)
	}
	partition := &topic.Partitions[0]

	ret := &ProduceResponse{
		Throttle:       makeDuration(res.ThrottleTimeMs),
		Error:          makeError(partition.ErrorCode, partition.ErrorMessage),
		BaseOffset:     partition.BaseOffset,
		LogAppendTime:  makeTime(partition.LogAppendTime),
		LogStartOffset: partition.LogStartOffset,
	}

	if len(partition.RecordErrors) != 0 {
		ret.RecordErrors = make(map[int]error, len(partition.RecordErrors))

		for _, recErr := range partition.RecordErrors {
			ret.RecordErrors[int(recErr.BatchIndex)] = errors.New(recErr.BatchIndexErrorMessage)
		}
	}

	return ret, nil
}
