package listoffsets

import (
	"sort"

	"github.com/segmentio/kafka-go/protocol"
)

func init() {
	protocol.Register(&Request{}, &Response{})
}

type Request struct {
	ReplicaID      int32          `kafka:"min=v1,max=v5"`
	IsolationLevel int8           `kafka:"min=v2,max=v5"`
	Topics         []RequestTopic `kafka:"min=v1,max=v5"`
}

type RequestTopic struct {
	Topic      string             `kafka:"min=v1,max=v5"`
	Partitions []RequestPartition `kafka:"min=v1,max=v5"`
}

type RequestPartition struct {
	Partition          int32 `kafka:"min=v1,max=v5"`
	CurrentLeaderEpoch int32 `kafka:"min=v4,max=v5"`
	Timestamp          int64 `kafka:"min=v1,max=v5"`
	
	
	
	
	
}

func (r *Request) ApiKey() protocol.ApiKey { return protocol.ListOffsets }

func (r *Request) Broker(cluster protocol.Cluster) (protocol.Broker, error) {
	
	
	partition := r.Topics[0].Partitions[0].Partition
	topic := r.Topics[0].Topic

	for _, p := range cluster.Topics[topic].Partitions {
		if p.ID == partition {
			return cluster.Brokers[p.Leader], nil
		}
	}

	return protocol.Broker{ID: -1}, nil
}

func (r *Request) Split(cluster protocol.Cluster) ([]protocol.Message, protocol.Merger, error) {
	
	
	
	
	
	
	
	
	
	
	
	
	
	requests := make([]Request, 0, 2*len(r.Topics))

	for _, t := range r.Topics {
		for _, p := range t.Partitions {
			requests = append(requests, Request{
				ReplicaID:      r.ReplicaID,
				IsolationLevel: r.IsolationLevel,
				Topics: []RequestTopic{{
					Topic: t.Topic,
					Partitions: []RequestPartition{{
						Partition:          p.Partition,
						CurrentLeaderEpoch: p.CurrentLeaderEpoch,
						Timestamp:          p.Timestamp,
					}},
				}},
			})
		}
	}

	messages := make([]protocol.Message, len(requests))

	for i := range requests {
		messages[i] = &requests[i]
	}

	return messages, new(Response), nil
}

type Response struct {
	ThrottleTimeMs int32           `kafka:"min=v2,max=v5"`
	Topics         []ResponseTopic `kafka:"min=v1,max=v5"`
}

type ResponseTopic struct {
	Topic      string              `kafka:"min=v1,max=v5"`
	Partitions []ResponsePartition `kafka:"min=v1,max=v5"`
}

type ResponsePartition struct {
	Partition   int32 `kafka:"min=v1,max=v5"`
	ErrorCode   int16 `kafka:"min=v1,max=v5"`
	Timestamp   int64 `kafka:"min=v1,max=v5"`
	Offset      int64 `kafka:"min=v1,max=v5"`
	LeaderEpoch int32 `kafka:"min=v4,max=v5"`
}

func (r *Response) ApiKey() protocol.ApiKey { return protocol.ListOffsets }

func (r *Response) Merge(requests []protocol.Message, results []interface{}) (protocol.Message, error) {
	type topicPartition struct {
		topic     string
		partition int32
	}

	
	
	
	
	
	
	
	
	
	timestamps := make([]map[topicPartition]int64, len(requests))

	for i, m := range requests {
		req := m.(*Request)
		ts := make(map[topicPartition]int64, len(req.Topics))

		for _, t := range req.Topics {
			for _, p := range t.Partitions {
				ts[topicPartition{
					topic:     t.Topic,
					partition: p.Partition,
				}] = p.Timestamp
			}
		}

		timestamps[i] = ts
	}

	topics := make(map[string][]ResponsePartition)
	errors := 0

	for i, res := range results {
		m, err := protocol.Result(res)
		if err != nil {
			for _, t := range requests[i].(*Request).Topics {
				partitions := topics[t.Topic]

				for _, p := range t.Partitions {
					partitions = append(partitions, ResponsePartition{
						Partition:   p.Partition,
						ErrorCode:   -1, 
						Timestamp:   -1,
						Offset:      -1,
						LeaderEpoch: -1,
					})
				}

				topics[t.Topic] = partitions
			}
			errors++
			continue
		}

		response := m.(*Response)

		if r.ThrottleTimeMs < response.ThrottleTimeMs {
			r.ThrottleTimeMs = response.ThrottleTimeMs
		}

		for _, t := range response.Topics {
			for _, p := range t.Partitions {
				if timestamp, ok := timestamps[i][topicPartition{
					topic:     t.Topic,
					partition: p.Partition,
				}]; ok {
					p.Timestamp = timestamp
				}
				topics[t.Topic] = append(topics[t.Topic], p)
			}
		}

	}

	if errors > 0 && errors == len(results) {
		_, err := protocol.Result(results[0])
		return nil, err
	}

	r.Topics = make([]ResponseTopic, 0, len(topics))

	for topicName, partitions := range topics {
		r.Topics = append(r.Topics, ResponseTopic{
			Topic:      topicName,
			Partitions: partitions,
		})
	}

	sort.Slice(r.Topics, func(i, j int) bool {
		return r.Topics[i].Topic < r.Topics[j].Topic
	})

	for _, t := range r.Topics {
		sort.Slice(t.Partitions, func(i, j int) bool {
			p1 := &t.Partitions[i]
			p2 := &t.Partitions[j]

			if p1.Partition != p2.Partition {
				return p1.Partition < p2.Partition
			}

			return p1.Offset < p2.Offset
		})
	}

	return r, nil
}

var (
	_ protocol.BrokerMessage = (*Request)(nil)
	_ protocol.Splitter      = (*Request)(nil)
	_ protocol.Merger        = (*Response)(nil)
)
