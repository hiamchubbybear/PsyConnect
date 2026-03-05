

package sonic

import (
	"errors"
)




type NoCopyRawMessage []byte


func (m NoCopyRawMessage) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return m, nil
}


func (m *NoCopyRawMessage) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("sonic.NoCopyRawMessage: UnmarshalJSON on nil pointer")
	}
	*m = data
	return nil
}
