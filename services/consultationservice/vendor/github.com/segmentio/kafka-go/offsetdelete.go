package kafka

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/segmentio/kafka-go/protocol/offsetdelete"
)



type OffsetDelete struct {
	Topic     string
	Partition int
}



type OffsetDeleteRequest struct {
	
	Addr net.Addr

	
	GroupID string

	
	Topics map[string][]int
}



type OffsetDeleteResponse struct {
	
	Error error

	
	Throttle time.Duration

	
	
	Topics map[string][]OffsetDeletePartition
}



type OffsetDeletePartition struct {
	
	Partition int

	
	
	Error error
}



func (c *Client) OffsetDelete(ctx context.Context, req *OffsetDeleteRequest) (*OffsetDeleteResponse, error) {
	topics := make([]offsetdelete.RequestTopic, 0, len(req.Topics))

	for topicName, partitionIndexes := range req.Topics {
		partitions := make([]offsetdelete.RequestPartition, len(partitionIndexes))

		for i, c := range partitionIndexes {
			partitions[i] = offsetdelete.RequestPartition{
				PartitionIndex: int32(c),
			}
		}

		topics = append(topics, offsetdelete.RequestTopic{
			Name:       topicName,
			Partitions: partitions,
		})
	}

	m, err := c.roundTrip(ctx, req.Addr, &offsetdelete.Request{
		GroupID: req.GroupID,
		Topics:  topics,
	})
	if err != nil {
		return nil, fmt.Errorf("kafka.(*Client).OffsetDelete: %w", err)
	}
	r := m.(*offsetdelete.Response)

	res := &OffsetDeleteResponse{
		Error:    makeError(r.ErrorCode, ""),
		Throttle: makeDuration(r.ThrottleTimeMs),
		Topics:   make(map[string][]OffsetDeletePartition, len(r.Topics)),
	}

	for _, topic := range r.Topics {
		partitions := make([]OffsetDeletePartition, len(topic.Partitions))

		for i, p := range topic.Partitions {
			partitions[i] = OffsetDeletePartition{
				Partition: int(p.PartitionIndex),
				Error:     makeError(p.ErrorCode, ""),
			}
		}

		res.Topics[topic.Name] = partitions
	}

	return res, nil
}
