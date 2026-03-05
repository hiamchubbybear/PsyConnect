


package telemetry

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)


type Scope struct {
	Name         string `json:"name,omitempty"`
	Version      string `json:"version,omitempty"`
	Attrs        []Attr `json:"attributes,omitempty"`
	DroppedAttrs uint32 `json:"droppedAttributesCount,omitempty"`
}


func (s *Scope) UnmarshalJSON(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))

	t, err := decoder.Token()
	if err != nil {
		return err
	}
	if t != json.Delim('{') {
		return errors.New("invalid Scope type")
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
			return fmt.Errorf("invalid Scope field: %#v", keyIface)
		}

		switch key {
		case "name":
			err = decoder.Decode(&s.Name)
		case "version":
			err = decoder.Decode(&s.Version)
		case "attributes":
			err = decoder.Decode(&s.Attrs)
		case "droppedAttributesCount", "dropped_attributes_count":
			err = decoder.Decode(&s.DroppedAttrs)
		default:
			
		}

		if err != nil {
			return err
		}
	}
	return nil
}
