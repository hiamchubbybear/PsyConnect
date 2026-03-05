package kafka

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/segmentio/kafka-go/protocol/createpartitions"
)



type CreatePartitionsRequest struct {
	
	Addr net.Addr

	
	Topics []TopicPartitionsConfig

	
	
	ValidateOnly bool
}



type CreatePartitionsResponse struct {
	
	Throttle time.Duration

	
	
	
	
	
	Errors map[string]error
}



func (c *Client) CreatePartitions(ctx context.Context, req *CreatePartitionsRequest) (*CreatePartitionsResponse, error) {
	topics := make([]createpartitions.RequestTopic, len(req.Topics))

	for i, t := range req.Topics {
		topics[i] = createpartitions.RequestTopic{
			Name:        t.Name,
			Count:       t.Count,
			Assignments: t.assignments(),
		}
	}

	m, err := c.roundTrip(ctx, req.Addr, &createpartitions.Request{
		Topics:       topics,
		TimeoutMs:    c.timeoutMs(ctx, defaultCreatePartitionsTimeout),
		ValidateOnly: req.ValidateOnly,
	})

	if err != nil {
		return nil, fmt.Errorf("kafka.(*Client).CreatePartitions: %w", err)
	}

	res := m.(*createpartitions.Response)
	ret := &CreatePartitionsResponse{
		Throttle: makeDuration(res.ThrottleTimeMs),
		Errors:   make(map[string]error, len(res.Results)),
	}

	for _, t := range res.Results {
		ret.Errors[t.Name] = makeError(t.ErrorCode, t.ErrorMessage)
	}

	return ret, nil
}

type TopicPartitionsConfig struct {
	
	Name string

	
	Count int32

	
	TopicPartitionAssignments []TopicPartitionAssignment
}

func (t *TopicPartitionsConfig) assignments() []createpartitions.RequestAssignment {
	if len(t.TopicPartitionAssignments) == 0 {
		return nil
	}
	assignments := make([]createpartitions.RequestAssignment, len(t.TopicPartitionAssignments))
	for i, a := range t.TopicPartitionAssignments {
		assignments[i] = createpartitions.RequestAssignment{
			BrokerIDs: a.BrokerIDs,
		}
	}
	return assignments
}

type TopicPartitionAssignment struct {
	
	BrokerIDs []int32
}
