package kafka

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"time"

	"github.com/segmentio/kafka-go/protocol/offsetcommit"
)






type OffsetCommit struct {
	Partition int
	Offset    int64
	Metadata  string
}



type OffsetCommitRequest struct {
	
	Addr net.Addr

	
	GroupID string

	
	GenerationID int

	
	MemberID string

	
	InstanceID string

	
	
	
	
	Topics map[string][]OffsetCommit
}



type OffsetCommitResponse struct {
	
	Throttle time.Duration

	
	
	Topics map[string][]OffsetCommitPartition
}



type OffsetCommitPartition struct {
	
	Partition int

	
	
	
	
	
	
	Error error
}



func (c *Client) OffsetCommit(ctx context.Context, req *OffsetCommitRequest) (*OffsetCommitResponse, error) {
	now := time.Now().UnixNano() / int64(time.Millisecond)
	topics := make([]offsetcommit.RequestTopic, 0, len(req.Topics))

	for topicName, commits := range req.Topics {
		partitions := make([]offsetcommit.RequestPartition, len(commits))

		for i, c := range commits {
			partitions[i] = offsetcommit.RequestPartition{
				PartitionIndex:    int32(c.Partition),
				CommittedOffset:   c.Offset,
				CommittedMetadata: c.Metadata,
				
				
				
				CommitTimestamp: now,
			}
		}

		topics = append(topics, offsetcommit.RequestTopic{
			Name:       topicName,
			Partitions: partitions,
		})
	}

	m, err := c.roundTrip(ctx, req.Addr, &offsetcommit.Request{
		GroupID:         req.GroupID,
		GenerationID:    int32(req.GenerationID),
		MemberID:        req.MemberID,
		GroupInstanceID: req.InstanceID,
		Topics:          topics,
		
		
		
		
		RetentionTimeMs: int64((24 * time.Hour) / time.Millisecond),
	})
	if err != nil {
		return nil, fmt.Errorf("kafka.(*Client).OffsetCommit: %w", err)
	}
	r := m.(*offsetcommit.Response)

	res := &OffsetCommitResponse{
		Throttle: makeDuration(r.ThrottleTimeMs),
		Topics:   make(map[string][]OffsetCommitPartition, len(r.Topics)),
	}

	for _, topic := range r.Topics {
		partitions := make([]OffsetCommitPartition, len(topic.Partitions))

		for i, p := range topic.Partitions {
			partitions[i] = OffsetCommitPartition{
				Partition: int(p.PartitionIndex),
				Error:     makeError(p.ErrorCode, ""),
			}
		}

		res.Topics[topic.Name] = partitions
	}

	return res, nil
}

type offsetCommitRequestV2Partition struct {
	
	Partition int32

	
	Offset int64

	
	Metadata string
}

func (t offsetCommitRequestV2Partition) size() int32 {
	return sizeofInt32(t.Partition) +
		sizeofInt64(t.Offset) +
		sizeofString(t.Metadata)
}

func (t offsetCommitRequestV2Partition) writeTo(wb *writeBuffer) {
	wb.writeInt32(t.Partition)
	wb.writeInt64(t.Offset)
	wb.writeString(t.Metadata)
}

type offsetCommitRequestV2Topic struct {
	
	Topic string

	
	Partitions []offsetCommitRequestV2Partition
}

func (t offsetCommitRequestV2Topic) size() int32 {
	return sizeofString(t.Topic) +
		sizeofArray(len(t.Partitions), func(i int) int32 { return t.Partitions[i].size() })
}

func (t offsetCommitRequestV2Topic) writeTo(wb *writeBuffer) {
	wb.writeString(t.Topic)
	wb.writeArray(len(t.Partitions), func(i int) { t.Partitions[i].writeTo(wb) })
}

type offsetCommitRequestV2 struct {
	
	GroupID string

	
	GenerationID int32

	
	MemberID string

	
	RetentionTime int64

	
	Topics []offsetCommitRequestV2Topic
}

func (t offsetCommitRequestV2) size() int32 {
	return sizeofString(t.GroupID) +
		sizeofInt32(t.GenerationID) +
		sizeofString(t.MemberID) +
		sizeofInt64(t.RetentionTime) +
		sizeofArray(len(t.Topics), func(i int) int32 { return t.Topics[i].size() })
}

func (t offsetCommitRequestV2) writeTo(wb *writeBuffer) {
	wb.writeString(t.GroupID)
	wb.writeInt32(t.GenerationID)
	wb.writeString(t.MemberID)
	wb.writeInt64(t.RetentionTime)
	wb.writeArray(len(t.Topics), func(i int) { t.Topics[i].writeTo(wb) })
}

type offsetCommitResponseV2PartitionResponse struct {
	Partition int32

	
	ErrorCode int16
}

func (t offsetCommitResponseV2PartitionResponse) size() int32 {
	return sizeofInt32(t.Partition) +
		sizeofInt16(t.ErrorCode)
}

func (t offsetCommitResponseV2PartitionResponse) writeTo(wb *writeBuffer) {
	wb.writeInt32(t.Partition)
	wb.writeInt16(t.ErrorCode)
}

func (t *offsetCommitResponseV2PartitionResponse) readFrom(r *bufio.Reader, size int) (remain int, err error) {
	if remain, err = readInt32(r, size, &t.Partition); err != nil {
		return
	}
	if remain, err = readInt16(r, remain, &t.ErrorCode); err != nil {
		return
	}
	return
}

type offsetCommitResponseV2Response struct {
	Topic              string
	PartitionResponses []offsetCommitResponseV2PartitionResponse
}

func (t offsetCommitResponseV2Response) size() int32 {
	return sizeofString(t.Topic) +
		sizeofArray(len(t.PartitionResponses), func(i int) int32 { return t.PartitionResponses[i].size() })
}

func (t offsetCommitResponseV2Response) writeTo(wb *writeBuffer) {
	wb.writeString(t.Topic)
	wb.writeArray(len(t.PartitionResponses), func(i int) { t.PartitionResponses[i].writeTo(wb) })
}

func (t *offsetCommitResponseV2Response) readFrom(r *bufio.Reader, size int) (remain int, err error) {
	if remain, err = readString(r, size, &t.Topic); err != nil {
		return
	}

	fn := func(r *bufio.Reader, withSize int) (fnRemain int, fnErr error) {
		item := offsetCommitResponseV2PartitionResponse{}
		if fnRemain, fnErr = (&item).readFrom(r, withSize); fnErr != nil {
			return
		}
		t.PartitionResponses = append(t.PartitionResponses, item)
		return
	}
	if remain, err = readArrayWith(r, remain, fn); err != nil {
		return
	}

	return
}

type offsetCommitResponseV2 struct {
	Responses []offsetCommitResponseV2Response
}

func (t offsetCommitResponseV2) size() int32 {
	return sizeofArray(len(t.Responses), func(i int) int32 { return t.Responses[i].size() })
}

func (t offsetCommitResponseV2) writeTo(wb *writeBuffer) {
	wb.writeArray(len(t.Responses), func(i int) { t.Responses[i].writeTo(wb) })
}

func (t *offsetCommitResponseV2) readFrom(r *bufio.Reader, size int) (remain int, err error) {
	fn := func(r *bufio.Reader, withSize int) (fnRemain int, fnErr error) {
		item := offsetCommitResponseV2Response{}
		if fnRemain, fnErr = (&item).readFrom(r, withSize); fnErr != nil {
			return
		}
		t.Responses = append(t.Responses, item)
		return
	}
	if remain, err = readArrayWith(r, size, fn); err != nil {
		return
	}

	return
}
