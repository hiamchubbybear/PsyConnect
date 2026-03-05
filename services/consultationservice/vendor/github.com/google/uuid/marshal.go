



package uuid

import "fmt"


func (uuid UUID) MarshalText() ([]byte, error) {
	var js [36]byte
	encodeHex(js[:], uuid)
	return js[:], nil
}


func (uuid *UUID) UnmarshalText(data []byte) error {
	id, err := ParseBytes(data)
	if err != nil {
		return err
	}
	*uuid = id
	return nil
}


func (uuid UUID) MarshalBinary() ([]byte, error) {
	return uuid[:], nil
}


func (uuid *UUID) UnmarshalBinary(data []byte) error {
	if len(data) != 16 {
		return fmt.Errorf("invalid UUID (got %d bytes)", len(data))
	}
	copy(uuid[:], data)
	return nil
}
