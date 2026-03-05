





package bsoncore

import (
	"bytes"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)


type MalformedElementError string

func (mee MalformedElementError) Error() string { return string(mee) }


const ErrElementMissingKey MalformedElementError = "element is missing key"


const ErrElementMissingType MalformedElementError = "element is missing type"


type Element []byte



func (e Element) Key() string {
	key, _ := e.KeyErr()
	return key
}




func (e Element) KeyBytes() []byte {
	key, _ := e.KeyBytesErr()
	return key
}


func (e Element) KeyErr() (string, error) {
	key, err := e.KeyBytesErr()
	return string(key), err
}



func (e Element) KeyBytesErr() ([]byte, error) {
	if len(e) == 0 {
		return nil, ErrElementMissingType
	}
	idx := bytes.IndexByte(e[1:], 0x00)
	if idx == -1 {
		return nil, ErrElementMissingKey
	}
	return e[1 : idx+1], nil
}


func (e Element) Validate() error {
	if len(e) < 1 {
		return ErrElementMissingType
	}
	idx := bytes.IndexByte(e[1:], 0x00)
	if idx == -1 {
		return ErrElementMissingKey
	}
	return Value{Type: bsontype.Type(e[0]), Data: e[idx+2:]}.Validate()
}




func (e Element) CompareKey(key []byte) bool {
	if len(e) < 2 {
		return false
	}
	idx := bytes.IndexByte(e[1:], 0x00)
	if idx == -1 {
		return false
	}
	if index := bytes.IndexByte(key, 0x00); index > -1 {
		key = key[:index]
	}
	return bytes.Equal(e[1:idx+1], key)
}



func (e Element) Value() Value {
	val, _ := e.ValueErr()
	return val
}


func (e Element) ValueErr() (Value, error) {
	if len(e) == 0 {
		return Value{}, ErrElementMissingType
	}
	idx := bytes.IndexByte(e[1:], 0x00)
	if idx == -1 {
		return Value{}, ErrElementMissingKey
	}

	val, rem, exists := ReadValue(e[idx+2:], bsontype.Type(e[0]))
	if !exists {
		return Value{}, NewInsufficientBytesError(e, rem)
	}
	return val, nil
}


func (e Element) String() string {
	if len(e) == 0 {
		return ""
	}
	t := bsontype.Type(e[0])
	idx := bytes.IndexByte(e[1:], 0x00)
	if idx == -1 {
		return ""
	}
	key, valBytes := []byte(e[1:idx+1]), []byte(e[idx+2:])
	val, _, valid := ReadValue(valBytes, t)
	if !valid {
		return ""
	}
	return "\"" + string(key) + "\": " + val.String()
}



func (e Element) DebugString() string {
	if len(e) == 0 {
		return "<malformed>"
	}
	t := bsontype.Type(e[0])
	idx := bytes.IndexByte(e[1:], 0x00)
	if idx == -1 {
		return fmt.Sprintf(`bson.Element{[%s]<malformed>}`, t)
	}
	key, valBytes := []byte(e[1:idx+1]), []byte(e[idx+2:])
	val, _, valid := ReadValue(valBytes, t)
	if !valid {
		return fmt.Sprintf(`bson.Element{[%s]"%s": <malformed>}`, t, key)
	}
	return fmt.Sprintf(`bson.Element{[%s]"%s": %v}`, t, key, val)
}
