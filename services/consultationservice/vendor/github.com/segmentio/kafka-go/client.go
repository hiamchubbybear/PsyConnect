package kafka

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/segmentio/kafka-go/protocol"
)

const (
	defaultCreateTopicsTimeout     = 2 * time.Second
	defaultDeleteTopicsTimeout     = 2 * time.Second
	defaultCreatePartitionsTimeout = 2 * time.Second
	defaultProduceTimeout          = 500 * time.Millisecond
	defaultMaxWait                 = 500 * time.Millisecond
)








type Client struct {
	
	
	
	
	
	Addr net.Addr

	
	
	
	Timeout time.Duration

	
	
	
	Transport RoundTripper
}









type TopicAndGroup struct {
	Topic   string
	GroupId string
}






func (c *Client) ConsumerOffsets(ctx context.Context, tg TopicAndGroup) (map[int]int64, error) {
	metadata, err := c.Metadata(ctx, &MetadataRequest{
		Topics: []string{tg.Topic},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get topic metadata :%w", err)
	}

	topic := metadata.Topics[0]
	partitions := make([]int, len(topic.Partitions))

	for i := range topic.Partitions {
		partitions[i] = topic.Partitions[i].ID
	}

	offsets, err := c.OffsetFetch(ctx, &OffsetFetchRequest{
		GroupID: tg.GroupId,
		Topics: map[string][]int{
			tg.Topic: partitions,
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get offsets: %w", err)
	}

	topicOffsets := offsets.Topics[topic.Name]
	partitionOffsets := make(map[int]int64, len(topicOffsets))

	for _, off := range topicOffsets {
		partitionOffsets[off.Partition] = off.CommittedOffset
	}

	return partitionOffsets, nil
}

func (c *Client) roundTrip(ctx context.Context, addr net.Addr, msg protocol.Message) (protocol.Message, error) {
	if c.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.Timeout)
		defer cancel()
	}

	if addr == nil {
		if addr = c.Addr; addr == nil {
			return nil, errors.New("no address was given for the kafka cluster in the request or on the client")
		}
	}

	return c.transport().RoundTrip(ctx, addr, msg)
}

func (c *Client) transport() RoundTripper {
	if c.Transport != nil {
		return c.Transport
	}
	return DefaultTransport
}

func (c *Client) timeout(ctx context.Context, defaultTimeout time.Duration) time.Duration {
	timeout := c.Timeout

	if deadline, ok := ctx.Deadline(); ok {
		if remain := time.Until(deadline); remain < timeout {
			timeout = remain
		}
	}

	if timeout > 0 {
		
		
		
		return timeout / 2
	}

	return defaultTimeout
}

func (c *Client) timeoutMs(ctx context.Context, defaultTimeout time.Duration) int32 {
	return milliseconds(c.timeout(ctx, defaultTimeout))
}
