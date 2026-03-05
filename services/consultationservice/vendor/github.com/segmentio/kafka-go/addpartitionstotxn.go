package kafka

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/segmentio/kafka-go/protocol/addpartitionstotxn"
)



type AddPartitionToTxn struct {
	
	Partition int
}


type AddPartitionsToTxnRequest struct {
	
	Addr net.Addr

	
	TransactionalID string

	
	
	ProducerID int

	
	ProducerEpoch int

	
	Topics map[string][]AddPartitionToTxn
}


type AddPartitionsToTxnResponse struct {
	
	Throttle time.Duration

	
	Topics map[string][]AddPartitionToTxnPartition
}



type AddPartitionToTxnPartition struct {
	
	Partition int

	
	
	
	
	
	Error error
}


func (c *Client) AddPartitionsToTxn(
	ctx context.Context,
	req *AddPartitionsToTxnRequest,
) (*AddPartitionsToTxnResponse, error) {
	protoReq := &addpartitionstotxn.Request{
		TransactionalID: req.TransactionalID,
		ProducerID:      int64(req.ProducerID),
		ProducerEpoch:   int16(req.ProducerEpoch),
	}
	protoReq.Topics = make([]addpartitionstotxn.RequestTopic, 0, len(req.Topics))

	for topic, partitions := range req.Topics {
		reqTopic := addpartitionstotxn.RequestTopic{
			Name:       topic,
			Partitions: make([]int32, len(partitions)),
		}
		for i, partition := range partitions {
			reqTopic.Partitions[i] = int32(partition.Partition)
		}
		protoReq.Topics = append(protoReq.Topics, reqTopic)
	}

	m, err := c.roundTrip(ctx, req.Addr, protoReq)
	if err != nil {
		return nil, fmt.Errorf("kafka.(*Client).AddPartitionsToTxn: %w", err)
	}

	r := m.(*addpartitionstotxn.Response)

	res := &AddPartitionsToTxnResponse{
		Throttle: makeDuration(r.ThrottleTimeMs),
		Topics:   make(map[string][]AddPartitionToTxnPartition, len(r.Results)),
	}

	for _, result := range r.Results {
		partitions := make([]AddPartitionToTxnPartition, 0, len(result.Results))
		for _, rp := range result.Results {
			partitions = append(partitions, AddPartitionToTxnPartition{
				Partition: int(rp.PartitionIndex),
				Error:     makeError(rp.ErrorCode, ""),
			})
		}
		res.Topics[result.Name] = partitions
	}

	return res, nil
}
