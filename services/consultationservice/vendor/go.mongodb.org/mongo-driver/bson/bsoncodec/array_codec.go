





package bsoncodec

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)




type ArrayCodec struct{}

var defaultArrayCodec = NewArrayCodec()





func NewArrayCodec() *ArrayCodec {
	return &ArrayCodec{}
}


func (ac *ArrayCodec) EncodeValue(_ EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.Type() != tCoreArray {
		return ValueEncoderError{Name: "CoreArrayEncodeValue", Types: []reflect.Type{tCoreArray}, Received: val}
	}

	arr := val.Interface().(bsoncore.Array)
	return bsonrw.Copier{}.CopyArrayFromBytes(vw, arr)
}


func (ac *ArrayCodec) DecodeValue(_ DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() || val.Type() != tCoreArray {
		return ValueDecoderError{Name: "CoreArrayDecodeValue", Types: []reflect.Type{tCoreArray}, Received: val}
	}

	if val.IsNil() {
		val.Set(reflect.MakeSlice(val.Type(), 0, 0))
	}

	val.SetLen(0)
	arr, err := bsonrw.Copier{}.AppendArrayBytes(val.Interface().(bsoncore.Array), vr)
	val.Set(reflect.ValueOf(arr))
	return err
}
