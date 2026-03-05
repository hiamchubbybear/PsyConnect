


package telemetry

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"
)



type Span struct {
	
	
	
	
	
	
	TraceID TraceID `json:"traceId,omitempty"`
	
	
	
	
	
	
	SpanID SpanID `json:"spanId,omitempty"`
	
	
	
	TraceState string `json:"traceState,omitempty"`
	
	
	ParentSpanID SpanID `json:"parentSpanId,omitempty"`
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	Flags uint32 `json:"flags,omitempty"`
	
	
	
	
	
	
	
	
	
	
	
	Name string `json:"name"`
	
	
	
	Kind SpanKind `json:"kind,omitempty"`
	
	
	
	
	
	
	StartTime time.Time `json:"startTimeUnixNano,omitempty"`
	
	
	
	
	
	
	EndTime time.Time `json:"endTimeUnixNano,omitempty"`
	
	
	
	
	
	
	
	
	
	
	
	
	Attrs []Attr `json:"attributes,omitempty"`
	
	
	
	DroppedAttrs uint32 `json:"droppedAttributesCount,omitempty"`
	
	Events []*SpanEvent `json:"events,omitempty"`
	
	
	DroppedEvents uint32 `json:"droppedEventsCount,omitempty"`
	
	
	Links []*SpanLink `json:"links,omitempty"`
	
	
	DroppedLinks uint32 `json:"droppedLinksCount,omitempty"`
	
	
	Status *Status `json:"status,omitempty"`
}


func (s Span) MarshalJSON() ([]byte, error) {
	startT := s.StartTime.UnixNano()
	if s.StartTime.IsZero() || startT < 0 {
		startT = 0
	}

	endT := s.EndTime.UnixNano()
	if s.EndTime.IsZero() || endT < 0 {
		endT = 0
	}

	
	var parentSpanId string
	if !s.ParentSpanID.IsEmpty() {
		b := make([]byte, hex.EncodedLen(spanIDSize))
		hex.Encode(b, s.ParentSpanID[:])
		parentSpanId = string(b)
	}

	type Alias Span
	return json.Marshal(struct {
		Alias
		ParentSpanID string `json:"parentSpanId,omitempty"`
		StartTime    uint64 `json:"startTimeUnixNano,omitempty"`
		EndTime      uint64 `json:"endTimeUnixNano,omitempty"`
	}{
		Alias:        Alias(s),
		ParentSpanID: parentSpanId,
		StartTime:    uint64(startT),
		EndTime:      uint64(endT),
	})
}


func (s *Span) UnmarshalJSON(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))

	t, err := decoder.Token()
	if err != nil {
		return err
	}
	if t != json.Delim('{') {
		return errors.New("invalid Span type")
	}

	for decoder.More() {
		keyIface, err := decoder.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				
				return nil
			}
			return err
		}

		key, ok := keyIface.(string)
		if !ok {
			return fmt.Errorf("invalid Span field: %#v", keyIface)
		}

		switch key {
		case "traceId", "trace_id":
			err = decoder.Decode(&s.TraceID)
		case "spanId", "span_id":
			err = decoder.Decode(&s.SpanID)
		case "traceState", "trace_state":
			err = decoder.Decode(&s.TraceState)
		case "parentSpanId", "parent_span_id":
			err = decoder.Decode(&s.ParentSpanID)
		case "flags":
			err = decoder.Decode(&s.Flags)
		case "name":
			err = decoder.Decode(&s.Name)
		case "kind":
			err = decoder.Decode(&s.Kind)
		case "startTimeUnixNano", "start_time_unix_nano":
			var val protoUint64
			err = decoder.Decode(&val)
			s.StartTime = time.Unix(0, int64(val.Uint64()))
		case "endTimeUnixNano", "end_time_unix_nano":
			var val protoUint64
			err = decoder.Decode(&val)
			s.EndTime = time.Unix(0, int64(val.Uint64()))
		case "attributes":
			err = decoder.Decode(&s.Attrs)
		case "droppedAttributesCount", "dropped_attributes_count":
			err = decoder.Decode(&s.DroppedAttrs)
		case "events":
			err = decoder.Decode(&s.Events)
		case "droppedEventsCount", "dropped_events_count":
			err = decoder.Decode(&s.DroppedEvents)
		case "links":
			err = decoder.Decode(&s.Links)
		case "droppedLinksCount", "dropped_links_count":
			err = decoder.Decode(&s.DroppedLinks)
		case "status":
			err = decoder.Decode(&s.Status)
		default:
			
		}

		if err != nil {
			return err
		}
	}
	return nil
}















type SpanFlags int32

const (
	
	SpanFlagsTraceFlagsMask SpanFlags = 255
	
	
	
	SpanFlagsContextHasIsRemoteMask SpanFlags = 256
	
	SpanFlagsContextIsRemoteMask SpanFlags = 512
)



type SpanKind int32

const (
	
	
	SpanKindInternal SpanKind = 1
	
	
	SpanKindServer SpanKind = 2
	
	SpanKindClient SpanKind = 3
	
	
	
	
	SpanKindProducer SpanKind = 4
	
	
	
	SpanKindConsumer SpanKind = 5
)



type SpanEvent struct {
	
	Time time.Time `json:"timeUnixNano,omitempty"`
	
	
	Name string `json:"name,omitempty"`
	
	
	
	Attrs []Attr `json:"attributes,omitempty"`
	
	
	DroppedAttrs uint32 `json:"droppedAttributesCount,omitempty"`
}


func (e SpanEvent) MarshalJSON() ([]byte, error) {
	t := e.Time.UnixNano()
	if e.Time.IsZero() || t < 0 {
		t = 0
	}

	type Alias SpanEvent
	return json.Marshal(struct {
		Alias
		Time uint64 `json:"timeUnixNano,omitempty"`
	}{
		Alias: Alias(e),
		Time:  uint64(t),
	})
}


func (se *SpanEvent) UnmarshalJSON(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))

	t, err := decoder.Token()
	if err != nil {
		return err
	}
	if t != json.Delim('{') {
		return errors.New("invalid SpanEvent type")
	}

	for decoder.More() {
		keyIface, err := decoder.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				
				return nil
			}
			return err
		}

		key, ok := keyIface.(string)
		if !ok {
			return fmt.Errorf("invalid SpanEvent field: %#v", keyIface)
		}

		switch key {
		case "timeUnixNano", "time_unix_nano":
			var val protoUint64
			err = decoder.Decode(&val)
			se.Time = time.Unix(0, int64(val.Uint64()))
		case "name":
			err = decoder.Decode(&se.Name)
		case "attributes":
			err = decoder.Decode(&se.Attrs)
		case "droppedAttributesCount", "dropped_attributes_count":
			err = decoder.Decode(&se.DroppedAttrs)
		default:
			
		}

		if err != nil {
			return err
		}
	}
	return nil
}





type SpanLink struct {
	
	
	TraceID TraceID `json:"traceId,omitempty"`
	
	SpanID SpanID `json:"spanId,omitempty"`
	
	TraceState string `json:"traceState,omitempty"`
	
	
	
	Attrs []Attr `json:"attributes,omitempty"`
	
	
	DroppedAttrs uint32 `json:"droppedAttributesCount,omitempty"`
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	Flags uint32 `json:"flags,omitempty"`
}


func (sl *SpanLink) UnmarshalJSON(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))

	t, err := decoder.Token()
	if err != nil {
		return err
	}
	if t != json.Delim('{') {
		return errors.New("invalid SpanLink type")
	}

	for decoder.More() {
		keyIface, err := decoder.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				
				return nil
			}
			return err
		}

		key, ok := keyIface.(string)
		if !ok {
			return fmt.Errorf("invalid SpanLink field: %#v", keyIface)
		}

		switch key {
		case "traceId", "trace_id":
			err = decoder.Decode(&sl.TraceID)
		case "spanId", "span_id":
			err = decoder.Decode(&sl.SpanID)
		case "traceState", "trace_state":
			err = decoder.Decode(&sl.TraceState)
		case "attributes":
			err = decoder.Decode(&sl.Attrs)
		case "droppedAttributesCount", "dropped_attributes_count":
			err = decoder.Decode(&sl.DroppedAttrs)
		case "flags":
			err = decoder.Decode(&sl.Flags)
		default:
			
		}

		if err != nil {
			return err
		}
	}
	return nil
}
