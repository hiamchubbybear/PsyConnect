package kafka

import "github.com/segmentio/kafka-go/protocol"


type Broker struct {
	Host string
	Port int
	ID   int
	Rack string
}


type Topic struct {
	
	Name string

	
	Internal bool

	
	Partitions []Partition

	
	
	
	
	
	
	Error error
}


type Partition struct {
	
	
	Topic string
	ID    int

	
	
	
	
	
	
	Leader   Broker
	Replicas []Broker
	Isr      []Broker

	
	OfflineReplicas []Broker

	
	
	
	
	
	
	Error error
}













func Marshal(v interface{}) ([]byte, error) {
	return protocol.Marshal(-1, v)
}




func Unmarshal(b []byte, v interface{}) error {
	return protocol.Unmarshal(b, -1, v)
}


type Version int16




func (n Version) Marshal(v interface{}) ([]byte, error) {
	return protocol.Marshal(int16(n), v)
}




func (n Version) Unmarshal(b []byte, v interface{}) error {
	return protocol.Unmarshal(b, int16(n), v)
}
