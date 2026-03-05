package json

import (
	"encoding/json"
)


type Codec struct{}

func (Codec) Encode(v map[string]any) ([]byte, error) {
	
	return json.MarshalIndent(v, "", "  ")
}

func (Codec) Decode(b []byte, v map[string]any) error {
	return json.Unmarshal(b, &v)
}
