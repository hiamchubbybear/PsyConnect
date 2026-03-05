package kafka

import (
	"github.com/segmentio/kafka-go/protocol"
)


type Header = protocol.Header







type Bytes = protocol.Bytes




func NewBytes(b []byte) Bytes { return protocol.NewBytes(b) }


func ReadAll(b Bytes) ([]byte, error) { return protocol.ReadAll(b) }




type Record = protocol.Record






type RecordReader = protocol.RecordReader



func NewRecordReader(records ...Record) RecordReader {
	return protocol.NewRecordReader(records...)
}
