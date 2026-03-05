


package telemetry

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)











type Traces struct {
	
	
	
	
	
	ResourceSpans []*ResourceSpans `json:"resourceSpans,omitempty"`
}


func (td *Traces) UnmarshalJSON(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))

	t, err := decoder.Token()
	if err != nil {
		return err
	}
	if t != json.Delim('{') {
		return errors.New("invalid TracesData type")
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
			return fmt.Errorf("invalid TracesData field: %#v", keyIface)
		}

		switch key {
		case "resourceSpans", "resource_spans":
			err = decoder.Decode(&td.ResourceSpans)
		default:
			
		}

		if err != nil {
			return err
		}
	}
	return nil
}


type ResourceSpans struct {
	
	
	Resource Resource `json:"resource"`
	
	ScopeSpans []*ScopeSpans `json:"scopeSpans,omitempty"`
	
	
	SchemaURL string `json:"schemaUrl,omitempty"`
}


func (rs *ResourceSpans) UnmarshalJSON(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))

	t, err := decoder.Token()
	if err != nil {
		return err
	}
	if t != json.Delim('{') {
		return errors.New("invalid ResourceSpans type")
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
			return fmt.Errorf("invalid ResourceSpans field: %#v", keyIface)
		}

		switch key {
		case "resource":
			err = decoder.Decode(&rs.Resource)
		case "scopeSpans", "scope_spans":
			err = decoder.Decode(&rs.ScopeSpans)
		case "schemaUrl", "schema_url":
			err = decoder.Decode(&rs.SchemaURL)
		default:
			
		}

		if err != nil {
			return err
		}
	}
	return nil
}


type ScopeSpans struct {
	
	
	
	Scope *Scope `json:"scope"`
	
	Spans []*Span `json:"spans,omitempty"`
	
	
	
	
	SchemaURL string `json:"schemaUrl,omitempty"`
}


func (ss *ScopeSpans) UnmarshalJSON(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))

	t, err := decoder.Token()
	if err != nil {
		return err
	}
	if t != json.Delim('{') {
		return errors.New("invalid ScopeSpans type")
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
			return fmt.Errorf("invalid ScopeSpans field: %#v", keyIface)
		}

		switch key {
		case "scope":
			err = decoder.Decode(&ss.Scope)
		case "spans":
			err = decoder.Decode(&ss.Spans)
		case "schemaUrl", "schema_url":
			err = decoder.Decode(&ss.SchemaURL)
		default:
			
		}

		if err != nil {
			return err
		}
	}
	return nil
}
