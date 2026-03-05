


package trace 

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
)

const (
	
	
	FlagsSampled = TraceFlags(0x01)

	errInvalidHexID errorConst = "trace-id and span-id can only contain [0-9a-f] characters, all lowercase"

	errInvalidTraceIDLength errorConst = "hex encoded trace-id must have length equals to 32"
	errNilTraceID           errorConst = "trace-id can't be all zero"

	errInvalidSpanIDLength errorConst = "hex encoded span-id must have length equals to 16"
	errNilSpanID           errorConst = "span-id can't be all zero"
)

type errorConst string

func (e errorConst) Error() string {
	return string(e)
}



type TraceID [16]byte

var (
	nilTraceID TraceID
	_          json.Marshaler = nilTraceID
)



func (t TraceID) IsValid() bool {
	return !bytes.Equal(t[:], nilTraceID[:])
}



func (t TraceID) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}


func (t TraceID) String() string {
	return hex.EncodeToString(t[:])
}


type SpanID [8]byte

var (
	nilSpanID SpanID
	_         json.Marshaler = nilSpanID
)



func (s SpanID) IsValid() bool {
	return !bytes.Equal(s[:], nilSpanID[:])
}



func (s SpanID) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}


func (s SpanID) String() string {
	return hex.EncodeToString(s[:])
}





func TraceIDFromHex(h string) (TraceID, error) {
	t := TraceID{}
	if len(h) != 32 {
		return t, errInvalidTraceIDLength
	}

	if err := decodeHex(h, t[:]); err != nil {
		return t, err
	}

	if !t.IsValid() {
		return t, errNilTraceID
	}
	return t, nil
}




func SpanIDFromHex(h string) (SpanID, error) {
	s := SpanID{}
	if len(h) != 16 {
		return s, errInvalidSpanIDLength
	}

	if err := decodeHex(h, s[:]); err != nil {
		return s, err
	}

	if !s.IsValid() {
		return s, errNilSpanID
	}
	return s, nil
}

func decodeHex(h string, b []byte) error {
	for _, r := range h {
		switch {
		case 'a' <= r && r <= 'f':
			continue
		case '0' <= r && r <= '9':
			continue
		default:
			return errInvalidHexID
		}
	}

	decoded, err := hex.DecodeString(h)
	if err != nil {
		return err
	}

	copy(b, decoded)
	return nil
}


type TraceFlags byte 


func (tf TraceFlags) IsSampled() bool {
	return tf&FlagsSampled == FlagsSampled
}


func (tf TraceFlags) WithSampled(sampled bool) TraceFlags { 
	if sampled {
		return tf | FlagsSampled
	}

	return tf &^ FlagsSampled
}



func (tf TraceFlags) MarshalJSON() ([]byte, error) {
	return json.Marshal(tf.String())
}


func (tf TraceFlags) String() string {
	return hex.EncodeToString([]byte{byte(tf)}[:])
}



type SpanContextConfig struct {
	TraceID    TraceID
	SpanID     SpanID
	TraceFlags TraceFlags
	TraceState TraceState
	Remote     bool
}



func NewSpanContext(config SpanContextConfig) SpanContext {
	return SpanContext{
		traceID:    config.TraceID,
		spanID:     config.SpanID,
		traceFlags: config.TraceFlags,
		traceState: config.TraceState,
		remote:     config.Remote,
	}
}


type SpanContext struct {
	traceID    TraceID
	spanID     SpanID
	traceFlags TraceFlags
	traceState TraceState
	remote     bool
}

var _ json.Marshaler = SpanContext{}



func (sc SpanContext) IsValid() bool {
	return sc.HasTraceID() && sc.HasSpanID()
}


func (sc SpanContext) IsRemote() bool {
	return sc.remote
}


func (sc SpanContext) WithRemote(remote bool) SpanContext {
	return SpanContext{
		traceID:    sc.traceID,
		spanID:     sc.spanID,
		traceFlags: sc.traceFlags,
		traceState: sc.traceState,
		remote:     remote,
	}
}


func (sc SpanContext) TraceID() TraceID {
	return sc.traceID
}


func (sc SpanContext) HasTraceID() bool {
	return sc.traceID.IsValid()
}


func (sc SpanContext) WithTraceID(traceID TraceID) SpanContext {
	return SpanContext{
		traceID:    traceID,
		spanID:     sc.spanID,
		traceFlags: sc.traceFlags,
		traceState: sc.traceState,
		remote:     sc.remote,
	}
}


func (sc SpanContext) SpanID() SpanID {
	return sc.spanID
}


func (sc SpanContext) HasSpanID() bool {
	return sc.spanID.IsValid()
}


func (sc SpanContext) WithSpanID(spanID SpanID) SpanContext {
	return SpanContext{
		traceID:    sc.traceID,
		spanID:     spanID,
		traceFlags: sc.traceFlags,
		traceState: sc.traceState,
		remote:     sc.remote,
	}
}


func (sc SpanContext) TraceFlags() TraceFlags {
	return sc.traceFlags
}


func (sc SpanContext) IsSampled() bool {
	return sc.traceFlags.IsSampled()
}


func (sc SpanContext) WithTraceFlags(flags TraceFlags) SpanContext {
	return SpanContext{
		traceID:    sc.traceID,
		spanID:     sc.spanID,
		traceFlags: flags,
		traceState: sc.traceState,
		remote:     sc.remote,
	}
}


func (sc SpanContext) TraceState() TraceState {
	return sc.traceState
}


func (sc SpanContext) WithTraceState(state TraceState) SpanContext {
	return SpanContext{
		traceID:    sc.traceID,
		spanID:     sc.spanID,
		traceFlags: sc.traceFlags,
		traceState: state,
		remote:     sc.remote,
	}
}


func (sc SpanContext) Equal(other SpanContext) bool {
	return sc.traceID == other.traceID &&
		sc.spanID == other.spanID &&
		sc.traceFlags == other.traceFlags &&
		sc.traceState.String() == other.traceState.String() &&
		sc.remote == other.remote
}


func (sc SpanContext) MarshalJSON() ([]byte, error) {
	return json.Marshal(SpanContextConfig{
		TraceID:    sc.traceID,
		SpanID:     sc.spanID,
		TraceFlags: sc.traceFlags,
		TraceState: sc.traceState,
		Remote:     sc.remote,
	})
}
