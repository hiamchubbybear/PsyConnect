


package telemetry

import (
	"encoding/hex"
	"errors"
	"fmt"
)

const (
	traceIDSize = 16
	spanIDSize  = 8
)


type TraceID [traceIDSize]byte


func (tid TraceID) String() string {
	return hex.EncodeToString(tid[:])
}


func (tid TraceID) IsEmpty() bool {
	return tid == [traceIDSize]byte{}
}


func (tid TraceID) MarshalJSON() ([]byte, error) {
	if tid.IsEmpty() {
		return []byte(`""`), nil
	}
	return marshalJSON(tid[:])
}



func (tid *TraceID) UnmarshalJSON(data []byte) error {
	*tid = [traceIDSize]byte{}
	return unmarshalJSON(tid[:], data)
}


type SpanID [spanIDSize]byte


func (sid SpanID) String() string {
	return hex.EncodeToString(sid[:])
}


func (sid SpanID) IsEmpty() bool {
	return sid == [spanIDSize]byte{}
}


func (sid SpanID) MarshalJSON() ([]byte, error) {
	if sid.IsEmpty() {
		return []byte(`""`), nil
	}
	return marshalJSON(sid[:])
}


func (sid *SpanID) UnmarshalJSON(data []byte) error {
	*sid = [spanIDSize]byte{}
	return unmarshalJSON(sid[:], data)
}


func marshalJSON(id []byte) ([]byte, error) {
	
	hexLen := hex.EncodedLen(len(id)) + 2

	b := make([]byte, hexLen)
	hex.Encode(b[1:hexLen-1], id)
	b[0], b[hexLen-1] = '"', '"'

	return b, nil
}


func unmarshalJSON(dst []byte, src []byte) error {
	if l := len(src); l >= 2 && src[0] == '"' && src[l-1] == '"' {
		src = src[1 : l-1]
	}
	nLen := len(src)
	if nLen == 0 {
		return nil
	}

	if len(dst) != hex.DecodedLen(nLen) {
		return errors.New("invalid length for ID")
	}

	_, err := hex.Decode(dst, src)
	if err != nil {
		return fmt.Errorf("cannot unmarshal ID from string '%s': %w", string(src), err)
	}
	return nil
}
