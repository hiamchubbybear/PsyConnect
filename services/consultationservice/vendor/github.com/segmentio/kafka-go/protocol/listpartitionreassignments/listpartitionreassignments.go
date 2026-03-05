package listpartitionreassignments

import "github.com/segmentio/kafka-go/protocol"

func init() {
	protocol.Register(&Request{}, &Response{})
}



type Request struct {
	
	
	_ struct{} `kafka:"min=v0,max=v0,tag"`

	TimeoutMs int32          `kafka:"min=v0,max=v0"`
	Topics    []RequestTopic `kafka:"min=v0,max=v0,nullable"`
}

type RequestTopic struct {
	
	
	_ struct{} `kafka:"min=v0,max=v0,tag"`

	Name             string  `kafka:"min=v0,max=v0"`
	PartitionIndexes []int32 `kafka:"min=v0,max=v0"`
}

func (r *Request) ApiKey() protocol.ApiKey {
	return protocol.ListPartitionReassignments
}

func (r *Request) Broker(cluster protocol.Cluster) (protocol.Broker, error) {
	return cluster.Brokers[cluster.Controller], nil
}

type Response struct {
	
	
	_ struct{} `kafka:"min=v0,max=v0,tag"`

	ThrottleTimeMs int32           `kafka:"min=v0,max=v0"`
	ErrorCode      int16           `kafka:"min=v0,max=v0"`
	ErrorMessage   string          `kafka:"min=v0,max=v0,nullable"`
	Topics         []ResponseTopic `kafka:"min=v0,max=v0"`
}

type ResponseTopic struct {
	
	
	_ struct{} `kafka:"min=v0,max=v0,tag"`

	Name       string              `kafka:"min=v0,max=v0"`
	Partitions []ResponsePartition `kafka:"min=v0,max=v0"`
}

type ResponsePartition struct {
	
	
	_ struct{} `kafka:"min=v0,max=v0,tag"`

	PartitionIndex   int32   `kafka:"min=v0,max=v0"`
	Replicas         []int32 `kafka:"min=v0,max=v0"`
	AddingReplicas   []int32 `kafka:"min=v0,max=v0"`
	RemovingReplicas []int32 `kafka:"min=v0,max=v0"`
}

func (r *Response) ApiKey() protocol.ApiKey {
	return protocol.ListPartitionReassignments
}
