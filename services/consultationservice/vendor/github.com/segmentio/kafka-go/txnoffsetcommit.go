package kafka

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/segmentio/kafka-go/protocol/txnoffsetcommit"
)



type TxnOffsetCommitRequest struct {
	
	Addr net.Addr

	
	TransactionalID string

	
	GroupID string

	
	
	ProducerID int

	
	ProducerEpoch int

	
	GenerationID int

	
	MemberID string

	
	GroupInstanceID string

	
	
	
	
	Topics map[string][]TxnOffsetCommit
}






type TxnOffsetCommit struct {
	Partition int
	Offset    int64
	Metadata  string
}



type TxnOffsetCommitResponse struct {
	
	Throttle time.Duration

	
	
	Topics map[string][]TxnOffsetCommitPartition
}



type TxnOffsetCommitPartition struct {
	
	Partition int

	
	
	
	
	
	
	Error error
}



func (c *Client) TxnOffsetCommit(
	ctx context.Context,
	req *TxnOffsetCommitRequest,
) (*TxnOffsetCommitResponse, error) {
	protoReq := &txnoffsetcommit.Request{
		TransactionalID: req.TransactionalID,
		GroupID:         req.GroupID,
		ProducerID:      int64(req.ProducerID),
		ProducerEpoch:   int16(req.ProducerEpoch),
		GenerationID:    int32(req.GenerationID),
		MemberID:        req.MemberID,
		GroupInstanceID: req.GroupInstanceID,
		Topics:          make([]txnoffsetcommit.RequestTopic, 0, len(req.Topics)),
	}

	for topic, partitions := range req.Topics {
		parts := make([]txnoffsetcommit.RequestPartition, len(partitions))
		for i, partition := range partitions {
			parts[i] = txnoffsetcommit.RequestPartition{
				Partition:         int32(partition.Partition),
				CommittedOffset:   int64(partition.Offset),
				CommittedMetadata: partition.Metadata,
			}
		}
		t := txnoffsetcommit.RequestTopic{
			Name:       topic,
			Partitions: parts,
		}

		protoReq.Topics = append(protoReq.Topics, t)
	}

	m, err := c.roundTrip(ctx, req.Addr, protoReq)
	if err != nil {
		return nil, fmt.Errorf("kafka.(*Client).TxnOffsetCommit: %w", err)
	}

	r := m.(*txnoffsetcommit.Response)

	res := &TxnOffsetCommitResponse{
		Throttle: makeDuration(r.ThrottleTimeMs),
		Topics:   make(map[string][]TxnOffsetCommitPartition, len(r.Topics)),
	}

	for _, topic := range r.Topics {
		partitions := make([]TxnOffsetCommitPartition, 0, len(topic.Partitions))
		for _, partition := range topic.Partitions {
			partitions = append(partitions, TxnOffsetCommitPartition{
				Partition: int(partition.Partition),
				Error:     makeError(partition.ErrorCode, ""),
			})
		}
		res.Topics[topic.Name] = partitions
	}

	return res, nil
}
