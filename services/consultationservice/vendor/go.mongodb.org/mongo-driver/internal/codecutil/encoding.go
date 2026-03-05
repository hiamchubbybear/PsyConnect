





package codecutil

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

var ErrNilValue = errors.New("value is nil")



type MarshalError struct {
	Value interface{}
	Err   error
}


func (e MarshalError) Error() string {
	return fmt.Sprintf("cannot transform type %s to a BSON Document: %v",
		reflect.TypeOf(e.Value), e.Err)
}


type EncoderFn func(io.Writer) (*bson.Encoder, error)



func MarshalValue(val interface{}, encFn EncoderFn) (bsoncore.Value, error) {
	
	if bval, ok := val.(bsoncore.Value); ok {
		return bval, nil
	}

	if val == nil {
		return bsoncore.Value{}, ErrNilValue
	}

	buf := new(bytes.Buffer)

	enc, err := encFn(buf)
	if err != nil {
		return bsoncore.Value{}, err
	}

	
	
	err = enc.Encode(bson.D{{Key: "", Value: val}})
	if err != nil {
		return bsoncore.Value{}, MarshalError{Value: val, Err: err}
	}

	return bsoncore.Document(buf.Bytes()).Index(0).Value(), nil
}
