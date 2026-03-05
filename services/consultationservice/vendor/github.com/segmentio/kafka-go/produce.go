package kafka

import (
	"bufio"
	"context"
	"encoding"
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go/protocol"
	produceAPI "github.com/segmentio/kafka-go/protocol/produce"
)

type RequiredAcks int

const (
	RequireNone RequiredAcks = 0
	RequireOne  RequiredAcks = 1
	RequireAll  RequiredAcks = -1
)

func (acks RequiredAcks) String() string {
	switch acks {
	case RequireNone:
		return "none"
	case RequireOne:
		return "one"
	case RequireAll:
		return "all"
	default:
		return "unknown"
	}
}

func (acks RequiredAcks) MarshalText() ([]byte, error) {
	return []byte(acks.String()), nil
}

func (acks *RequiredAcks) UnmarshalText(b []byte) error {
	switch string(b) {
	case "none":
		*acks = RequireNone
	case "one":
		*acks = RequireOne
	case "all":
		*acks = RequireAll
	default:
		x, err := strconv.ParseInt(string(b), 10, 64)
		parsed := RequiredAcks(x)
		if err != nil || (parsed != RequireNone && parsed != RequireOne && parsed != RequireAll) {
			return fmt.Errorf("required acks must be one of none, one, or all, not %q", b)
		}
		*acks = parsed
	}
	return nil
}

var (
	_ encoding.TextMarshaler   = RequiredAcks(0)
	_ encoding.TextUnmarshaler = (*RequiredAcks)(nil)
)



type ProduceRequest struct {
	
	Addr net.Addr

	
	Topic string

	
	Partition int

	
	RequiredAcks RequiredAcks

	
	
	
	
	MessageVersion int

	
	
	TransactionalID string

	
	Records RecordReader

	
	
	Compression Compression
}



type ProduceResponse struct {
	
	Throttle time.Duration

	
	
	
	
	
	Error error

	
	
	
	
	BaseOffset int64

	
	
	
	
	LogAppendTime time.Time

	
	
	
	
	LogStartOffset int64

	
	
	
	
	
	RecordErrors map[int]error
}








func (c *Client) Produce(ctx context.Context, req *ProduceRequest) (*ProduceResponse, error) {
	attributes := protocol.Attributes(req.Compression) & 0x7

	m, err := c.roundTrip(ctx, req.Addr, &produceAPI.Request{
		TransactionalID: req.TransactionalID,
		Acks:            int16(req.RequiredAcks),
		Timeout:         c.timeoutMs(ctx, defaultProduceTimeout),
		Topics: []produceAPI.RequestTopic{{
			Topic: req.Topic,
			Partitions: []produceAPI.RequestPartition{{
				Partition: int32(req.Partition),
				RecordSet: protocol.RecordSet{
					Attributes: attributes,
					Records:    req.Records,
				},
			}},
		}},
	})

	switch {
	case err == nil:
	case errors.Is(err, protocol.ErrNoRecord):
		return new(ProduceResponse), nil
	default:
		return nil, fmt.Errorf("kafka.(*Client).Produce: %w", err)
	}

	if req.RequiredAcks == RequireNone {
		return nil, nil
	}

	res := m.(*produceAPI.Response)
	if len(res.Topics) == 0 {
		return nil, fmt.Errorf("kafka.(*Client).Produce: %w", protocol.ErrNoTopic)
	}
	topic := &res.Topics[0]
	if len(topic.Partitions) == 0 {
		return nil, fmt.Errorf("kafka.(*Client).Produce: %w", protocol.ErrNoPartition)
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

type produceRequestV2 struct {
	RequiredAcks int16
	Timeout      int32
	Topics       []produceRequestTopicV2
}

func (r produceRequestV2) size() int32 {
	return 2 + 4 + sizeofArray(len(r.Topics), func(i int) int32 { return r.Topics[i].size() })
}

func (r produceRequestV2) writeTo(wb *writeBuffer) {
	wb.writeInt16(r.RequiredAcks)
	wb.writeInt32(r.Timeout)
	wb.writeArray(len(r.Topics), func(i int) { r.Topics[i].writeTo(wb) })
}

type produceRequestTopicV2 struct {
	TopicName  string
	Partitions []produceRequestPartitionV2
}

func (t produceRequestTopicV2) size() int32 {
	return sizeofString(t.TopicName) +
		sizeofArray(len(t.Partitions), func(i int) int32 { return t.Partitions[i].size() })
}

func (t produceRequestTopicV2) writeTo(wb *writeBuffer) {
	wb.writeString(t.TopicName)
	wb.writeArray(len(t.Partitions), func(i int) { t.Partitions[i].writeTo(wb) })
}

type produceRequestPartitionV2 struct {
	Partition      int32
	MessageSetSize int32
	MessageSet     messageSet
}

func (p produceRequestPartitionV2) size() int32 {
	return 4 + 4 + p.MessageSet.size()
}

func (p produceRequestPartitionV2) writeTo(wb *writeBuffer) {
	wb.writeInt32(p.Partition)
	wb.writeInt32(p.MessageSetSize)
	p.MessageSet.writeTo(wb)
}

type produceResponsePartitionV2 struct {
	Partition int32
	ErrorCode int16
	Offset    int64
	Timestamp int64
}

func (p produceResponsePartitionV2) size() int32 {
	return 4 + 2 + 8 + 8
}

func (p produceResponsePartitionV2) writeTo(wb *writeBuffer) {
	wb.writeInt32(p.Partition)
	wb.writeInt16(p.ErrorCode)
	wb.writeInt64(p.Offset)
	wb.writeInt64(p.Timestamp)
}

func (p *produceResponsePartitionV2) readFrom(r *bufio.Reader, sz int) (remain int, err error) {
	if remain, err = readInt32(r, sz, &p.Partition); err != nil {
		return
	}
	if remain, err = readInt16(r, remain, &p.ErrorCode); err != nil {
		return
	}
	if remain, err = readInt64(r, remain, &p.Offset); err != nil {
		return
	}
	if remain, err = readInt64(r, remain, &p.Timestamp); err != nil {
		return
	}
	return
}

type produceResponsePartitionV7 struct {
	Partition   int32
	ErrorCode   int16
	Offset      int64
	Timestamp   int64
	StartOffset int64
}

func (p produceResponsePartitionV7) size() int32 {
	return 4 + 2 + 8 + 8 + 8
}

func (p produceResponsePartitionV7) writeTo(wb *writeBuffer) {
	wb.writeInt32(p.Partition)
	wb.writeInt16(p.ErrorCode)
	wb.writeInt64(p.Offset)
	wb.writeInt64(p.Timestamp)
	wb.writeInt64(p.StartOffset)
}

func (p *produceResponsePartitionV7) readFrom(r *bufio.Reader, sz int) (remain int, err error) {
	if remain, err = readInt32(r, sz, &p.Partition); err != nil {
		return
	}
	if remain, err = readInt16(r, remain, &p.ErrorCode); err != nil {
		return
	}
	if remain, err = readInt64(r, remain, &p.Offset); err != nil {
		return
	}
	if remain, err = readInt64(r, remain, &p.Timestamp); err != nil {
		return
	}
	if remain, err = readInt64(r, remain, &p.StartOffset); err != nil {
		return
	}
	return
}
