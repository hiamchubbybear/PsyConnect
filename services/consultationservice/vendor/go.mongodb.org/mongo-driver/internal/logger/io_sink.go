





package logger

import (
	"encoding/json"
	"io"
	"math"
	"sync"
	"time"
)


type IOSink struct {
	enc *json.Encoder

	
	
	
	encMu sync.Mutex
}


var _ LogSink = &IOSink{}



func NewIOSink(out io.Writer) *IOSink {
	return &IOSink{
		enc: json.NewEncoder(out),
	}
}


func (sink *IOSink) Info(_ int, msg string, keysAndValues ...interface{}) {
	mapSize := len(keysAndValues) / 2
	if math.MaxInt-mapSize >= 2 {
		mapSize += 2
	}
	kvMap := make(map[string]interface{}, mapSize)

	kvMap[KeyTimestamp] = time.Now().UnixNano()
	kvMap[KeyMessage] = msg

	for i := 0; i < len(keysAndValues); i += 2 {
		kvMap[keysAndValues[i].(string)] = keysAndValues[i+1]
	}

	sink.encMu.Lock()
	defer sink.encMu.Unlock()

	_ = sink.enc.Encode(kvMap)
}


func (sink *IOSink) Error(err error, msg string, kv ...interface{}) {
	kv = append(kv, KeyError, err.Error())
	sink.Info(0, msg, kv...)
}
