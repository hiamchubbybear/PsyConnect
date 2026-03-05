package yaml

import "gopkg.in/yaml.v3"


type Codec struct{}

func (Codec) Encode(v map[string]any) ([]byte, error) {
	return yaml.Marshal(v)
}

func (Codec) Decode(b []byte, v map[string]any) error {
	return yaml.Unmarshal(b, &v)
}
