


package telemetry

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)


type Resource struct {
	
	
	
	Attrs []Attr `json:"attributes,omitempty"`
	
	
	DroppedAttrs uint32 `json:"droppedAttributesCount,omitempty"`
}


func (r *Resource) UnmarshalJSON(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))

	t, err := decoder.Token()
	if err != nil {
		return err
	}
	if t != json.Delim('{') {
		return errors.New("invalid Resource type")
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
			return fmt.Errorf("invalid Resource field: %#v", keyIface)
		}

		switch key {
		case "attributes":
			err = decoder.Decode(&r.Attrs)
		case "droppedAttributesCount", "dropped_attributes_count":
			err = decoder.Decode(&r.DroppedAttrs)
		default:
			
		}

		if err != nil {
			return err
		}
	}
	return nil
}
