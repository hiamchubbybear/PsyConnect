package kafka

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"time"

	"github.com/segmentio/kafka-go/protocol/offsetfetch"
)



type OffsetFetchRequest struct {
	
	Addr net.Addr

	
	GroupID string

	
	Topics map[string][]int
}



type OffsetFetchResponse struct {
	
	Throttle time.Duration

	
	Topics map[string][]OffsetFetchPartition

	
	
	
	
	
	
	Error error
}



type OffsetFetchPartition struct {
	
	Partition int

	
	
	CommittedOffset int64

	
	Metadata string

	
	
	
	
	
	
	Error error
}



func (c *Client) OffsetFetch(ctx context.Context, req *OffsetFetchRequest) (*OffsetFetchResponse, error) {

	
	
	
	
	var topics []offsetfetch.RequestTopic

	if len(req.Topics) > 0 {
		topics = make([]offsetfetch.RequestTopic, 0, len(req.Topics))

		for topicName, partitions := range req.Topics {
			indexes := make([]int32, len(partitions))

			for i, p := range partitions {
				indexes[i] = int32(p)
			}

			topics = append(topics, offsetfetch.RequestTopic{
				Name:             topicName,
				PartitionIndexes: indexes,
			})
		}
	}

	m, err := c.roundTrip(ctx, req.Addr, &offsetfetch.Request{
		GroupID: req.GroupID,
		Topics:  topics,
	})

	if err != nil {
		return nil, fmt.Errorf("kafka.(*Client).OffsetFetch: %w", err)
	}

	res := m.(*offsetfetch.Response)
	ret := &OffsetFetchResponse{
		Throttle: makeDuration(res.ThrottleTimeMs),
		Topics:   make(map[string][]OffsetFetchPartition, len(res.Topics)),
		Error:    makeError(res.ErrorCode, ""),
	}

	for _, t := range res.Topics {
		partitions := make([]OffsetFetchPartition, len(t.Partitions))

		for i, p := range t.Partitions {
			partitions[i] = OffsetFetchPartition{
				Partition:       int(p.PartitionIndex),
				CommittedOffset: p.CommittedOffset,
				Metadata:        p.Metadata,
				Error:           makeError(p.ErrorCode, ""),
			}
		}

		ret.Topics[t.Name] = partitions
	}

	return ret, nil
}

type offsetFetchRequestV1Topic struct {
	
	Topic string

	
	Partitions []int32
}

func (t offsetFetchRequestV1Topic) size() int32 {
	return sizeofString(t.Topic) +
		sizeofInt32Array(t.Partitions)
}

func (t offsetFetchRequestV1Topic) writeTo(wb *writeBuffer) {
	wb.writeString(t.Topic)
	wb.writeInt32Array(t.Partitions)
}

type offsetFetchRequestV1 struct {
	
	GroupID string

	
	Topics []offsetFetchRequestV1Topic
}

func (t offsetFetchRequestV1) size() int32 {
	return sizeofString(t.GroupID) +
		sizeofArray(len(t.Topics), func(i int) int32 { return t.Topics[i].size() })
}

func (t offsetFetchRequestV1) writeTo(wb *writeBuffer) {
	wb.writeString(t.GroupID)
	wb.writeArray(len(t.Topics), func(i int) { t.Topics[i].writeTo(wb) })
}

type offsetFetchResponseV1PartitionResponse struct {
	
	Partition int32

	
	Offset int64

	
	Metadata string

	
	ErrorCode int16
}

func (t offsetFetchResponseV1PartitionResponse) size() int32 {
	return sizeofInt32(t.Partition) +
		sizeofInt64(t.Offset) +
		sizeofString(t.Metadata) +
		sizeofInt16(t.ErrorCode)
}

func (t offsetFetchResponseV1PartitionResponse) writeTo(wb *writeBuffer) {
	wb.writeInt32(t.Partition)
	wb.writeInt64(t.Offset)
	wb.writeString(t.Metadata)
	wb.writeInt16(t.ErrorCode)
}

func (t *offsetFetchResponseV1PartitionResponse) readFrom(r *bufio.Reader, size int) (remain int, err error) {
	if remain, err = readInt32(r, size, &t.Partition); err != nil {
		return
	}
	if remain, err = readInt64(r, remain, &t.Offset); err != nil {
		return
	}
	if remain, err = readString(r, remain, &t.Metadata); err != nil {
		return
	}
	if remain, err = readInt16(r, remain, &t.ErrorCode); err != nil {
		return
	}
	return
}

type offsetFetchResponseV1Response struct {
	
	Topic string

	
	PartitionResponses []offsetFetchResponseV1PartitionResponse
}

func (t offsetFetchResponseV1Response) size() int32 {
	return sizeofString(t.Topic) +
		sizeofArray(len(t.PartitionResponses), func(i int) int32 { return t.PartitionResponses[i].size() })
}

func (t offsetFetchResponseV1Response) writeTo(wb *writeBuffer) {
	wb.writeString(t.Topic)
	wb.writeArray(len(t.PartitionResponses), func(i int) { t.PartitionResponses[i].writeTo(wb) })
}

func (t *offsetFetchResponseV1Response) readFrom(r *bufio.Reader, size int) (remain int, err error) {
	if remain, err = readString(r, size, &t.Topic); err != nil {
		return
	}

	fn := func(r *bufio.Reader, size int) (fnRemain int, fnErr error) {
		item := offsetFetchResponseV1PartitionResponse{}
		if fnRemain, fnErr = (&item).readFrom(r, size); err != nil {
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

type offsetFetchResponseV1 struct {
	
	Responses []offsetFetchResponseV1Response
}

func (t offsetFetchResponseV1) size() int32 {
	return sizeofArray(len(t.Responses), func(i int) int32 { return t.Responses[i].size() })
}

func (t offsetFetchResponseV1) writeTo(wb *writeBuffer) {
	wb.writeArray(len(t.Responses), func(i int) { t.Responses[i].writeTo(wb) })
}

func (t *offsetFetchResponseV1) readFrom(r *bufio.Reader, size int) (remain int, err error) {
	fn := func(r *bufio.Reader, withSize int) (fnRemain int, fnErr error) {
		item := offsetFetchResponseV1Response{}
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
