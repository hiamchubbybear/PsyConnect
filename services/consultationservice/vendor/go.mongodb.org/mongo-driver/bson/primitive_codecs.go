





package bson

import (
	"errors"
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
)

var tRawValue = reflect.TypeOf(RawValue{})
var tRaw = reflect.TypeOf(Raw(nil))

var primitiveCodecs PrimitiveCodecs






type PrimitiveCodecs struct{}






func (pc PrimitiveCodecs) RegisterPrimitiveCodecs(rb *bsoncodec.RegistryBuilder) {
	if rb == nil {
		panic(errors.New("argument to RegisterPrimitiveCodecs must not be nil"))
	}

	rb.
		RegisterTypeEncoder(tRawValue, bsoncodec.ValueEncoderFunc(pc.RawValueEncodeValue)).
		RegisterTypeEncoder(tRaw, bsoncodec.ValueEncoderFunc(pc.RawEncodeValue)).
		RegisterTypeDecoder(tRawValue, bsoncodec.ValueDecoderFunc(pc.RawValueDecodeValue)).
		RegisterTypeDecoder(tRaw, bsoncodec.ValueDecoderFunc(pc.RawDecodeValue))
}








func (PrimitiveCodecs) RawValueEncodeValue(_ bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.Type() != tRawValue {
		return bsoncodec.ValueEncoderError{
			Name:     "RawValueEncodeValue",
			Types:    []reflect.Type{tRawValue},
			Received: val,
		}
	}

	rawvalue := val.Interface().(RawValue)

	if !rawvalue.Type.IsValid() {
		return fmt.Errorf("the RawValue Type specifies an invalid BSON type: %#x", byte(rawvalue.Type))
	}

	return bsonrw.Copier{}.CopyValueFromBytes(vw, rawvalue.Type, rawvalue.Value)
}





func (PrimitiveCodecs) RawValueDecodeValue(_ bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() || val.Type() != tRawValue {
		return bsoncodec.ValueDecoderError{Name: "RawValueDecodeValue", Types: []reflect.Type{tRawValue}, Received: val}
	}

	t, value, err := bsonrw.Copier{}.CopyValueToBytes(vr)
	if err != nil {
		return err
	}

	val.Set(reflect.ValueOf(RawValue{Type: t, Value: value}))
	return nil
}





func (PrimitiveCodecs) RawEncodeValue(_ bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.Type() != tRaw {
		return bsoncodec.ValueEncoderError{Name: "RawEncodeValue", Types: []reflect.Type{tRaw}, Received: val}
	}

	rdr := val.Interface().(Raw)

	return bsonrw.Copier{}.CopyDocumentFromBytes(vw, rdr)
}





func (PrimitiveCodecs) RawDecodeValue(_ bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() || val.Type() != tRaw {
		return bsoncodec.ValueDecoderError{Name: "RawDecodeValue", Types: []reflect.Type{tRaw}, Received: val}
	}

	if val.IsNil() {
		val.Set(reflect.MakeSlice(val.Type(), 0, 0))
	}

	val.SetLen(0)

	rdr, err := bsonrw.Copier{}.AppendDocumentBytes(val.Interface().(Raw), vr)
	val.Set(reflect.ValueOf(rdr))
	return err
}
