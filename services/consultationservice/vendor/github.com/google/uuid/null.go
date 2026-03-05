



package uuid

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

var jsonNull = []byte("null")














type NullUUID struct {
	UUID  UUID
	Valid bool 
}


func (nu *NullUUID) Scan(value interface{}) error {
	if value == nil {
		nu.UUID, nu.Valid = Nil, false
		return nil
	}

	err := nu.UUID.Scan(value)
	if err != nil {
		nu.Valid = false
		return err
	}

	nu.Valid = true
	return nil
}


func (nu NullUUID) Value() (driver.Value, error) {
	if !nu.Valid {
		return nil, nil
	}
	
	return nu.UUID.Value()
}


func (nu NullUUID) MarshalBinary() ([]byte, error) {
	if nu.Valid {
		return nu.UUID[:], nil
	}

	return []byte(nil), nil
}


func (nu *NullUUID) UnmarshalBinary(data []byte) error {
	if len(data) != 16 {
		return fmt.Errorf("invalid UUID (got %d bytes)", len(data))
	}
	copy(nu.UUID[:], data)
	nu.Valid = true
	return nil
}


func (nu NullUUID) MarshalText() ([]byte, error) {
	if nu.Valid {
		return nu.UUID.MarshalText()
	}

	return jsonNull, nil
}


func (nu *NullUUID) UnmarshalText(data []byte) error {
	id, err := ParseBytes(data)
	if err != nil {
		nu.Valid = false
		return err
	}
	nu.UUID = id
	nu.Valid = true
	return nil
}


func (nu NullUUID) MarshalJSON() ([]byte, error) {
	if nu.Valid {
		return json.Marshal(nu.UUID)
	}

	return jsonNull, nil
}


func (nu *NullUUID) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, jsonNull) {
		*nu = NullUUID{}
		return nil 
	}
	err := json.Unmarshal(data, &nu.UUID)
	nu.Valid = err == nil
	return err
}
