








package primitive

import (
	"crypto/rand"
	"encoding"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync/atomic"
	"time"
)


var ErrInvalidHex = errors.New("the provided hex string is not a valid ObjectID")


type ObjectID [12]byte


var NilObjectID ObjectID

var objectIDCounter = readRandomUint32()
var processUnique = processUniqueBytes()

var _ encoding.TextMarshaler = ObjectID{}
var _ encoding.TextUnmarshaler = &ObjectID{}


func NewObjectID() ObjectID {
	return NewObjectIDFromTimestamp(time.Now())
}


func NewObjectIDFromTimestamp(timestamp time.Time) ObjectID {
	var b [12]byte

	binary.BigEndian.PutUint32(b[0:4], uint32(timestamp.Unix()))
	copy(b[4:9], processUnique[:])
	putUint24(b[9:12], atomic.AddUint32(&objectIDCounter, 1))

	return b
}


func (id ObjectID) Timestamp() time.Time {
	unixSecs := binary.BigEndian.Uint32(id[0:4])
	return time.Unix(int64(unixSecs), 0).UTC()
}


func (id ObjectID) Hex() string {
	var buf [24]byte
	hex.Encode(buf[:], id[:])
	return string(buf[:])
}

func (id ObjectID) String() string {
	return fmt.Sprintf("ObjectID(%q)", id.Hex())
}


func (id ObjectID) IsZero() bool {
	return id == NilObjectID
}



func ObjectIDFromHex(s string) (ObjectID, error) {
	if len(s) != 24 {
		return NilObjectID, ErrInvalidHex
	}

	var oid [12]byte
	_, err := hex.Decode(oid[:], []byte(s))
	if err != nil {
		return NilObjectID, err
	}

	return oid, nil
}




func IsValidObjectID(s string) bool {
	_, err := ObjectIDFromHex(s)
	return err == nil
}



func (id ObjectID) MarshalText() ([]byte, error) {
	return []byte(id.Hex()), nil
}



func (id *ObjectID) UnmarshalText(b []byte) error {
	oid, err := ObjectIDFromHex(string(b))
	if err != nil {
		return err
	}
	*id = oid
	return nil
}


func (id ObjectID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.Hex())
}





func (id *ObjectID) UnmarshalJSON(b []byte) error {
	
	
	
	if string(b) == "null" {
		return nil
	}

	var err error
	switch len(b) {
	case 12:
		copy(id[:], b)
	default:
		
		var res interface{}
		err := json.Unmarshal(b, &res)
		if err != nil {
			return err
		}
		str, ok := res.(string)
		if !ok {
			m, ok := res.(map[string]interface{})
			if !ok {
				return errors.New("not an extended JSON ObjectID")
			}
			oid, ok := m["$oid"]
			if !ok {
				return errors.New("not an extended JSON ObjectID")
			}
			str, ok = oid.(string)
			if !ok {
				return errors.New("not an extended JSON ObjectID")
			}
		}

		
		if len(str) == 0 {
			copy(id[:], NilObjectID[:])
			return nil
		}

		if len(str) != 24 {
			return fmt.Errorf("cannot unmarshal into an ObjectID, the length must be 24 but it is %d", len(str))
		}

		_, err = hex.Decode(id[:], []byte(str))
		if err != nil {
			return err
		}
	}

	return err
}

func processUniqueBytes() [5]byte {
	var b [5]byte
	_, err := io.ReadFull(rand.Reader, b[:])
	if err != nil {
		panic(fmt.Errorf("cannot initialize objectid package with crypto.rand.Reader: %w", err))
	}

	return b
}

func readRandomUint32() uint32 {
	var b [4]byte
	_, err := io.ReadFull(rand.Reader, b[:])
	if err != nil {
		panic(fmt.Errorf("cannot initialize objectid package with crypto.rand.Reader: %w", err))
	}

	return (uint32(b[0]) << 0) | (uint32(b[1]) << 8) | (uint32(b[2]) << 16) | (uint32(b[3]) << 24)
}

func putUint24(b []byte, v uint32) {
	b[0] = byte(v >> 16)
	b[1] = byte(v >> 8)
	b[2] = byte(v)
}
