





package bsoncodec

import (
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson/bsonoptions"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)



















type ByteSliceCodec struct {
	
	
	
	
	
	EncodeNilAsEmpty bool
}

var (
	defaultByteSliceCodec = NewByteSliceCodec()

	
	
	
	_ typeDecoder = defaultByteSliceCodec
)





func NewByteSliceCodec(opts ...*bsonoptions.ByteSliceCodecOptions) *ByteSliceCodec {
	byteSliceOpt := bsonoptions.MergeByteSliceCodecOptions(opts...)
	codec := ByteSliceCodec{}
	if byteSliceOpt.EncodeNilAsEmpty != nil {
		codec.EncodeNilAsEmpty = *byteSliceOpt.EncodeNilAsEmpty
	}
	return &codec
}


func (bsc *ByteSliceCodec) EncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.Type() != tByteSlice {
		return ValueEncoderError{Name: "ByteSliceEncodeValue", Types: []reflect.Type{tByteSlice}, Received: val}
	}
	if val.IsNil() && !bsc.EncodeNilAsEmpty && !ec.nilByteSliceAsEmpty {
		return vw.WriteNull()
	}
	return vw.WriteBinary(val.Interface().([]byte))
}

func (bsc *ByteSliceCodec) decodeType(_ DecodeContext, vr bsonrw.ValueReader, t reflect.Type) (reflect.Value, error) {
	if t != tByteSlice {
		return emptyValue, ValueDecoderError{
			Name:     "ByteSliceDecodeValue",
			Types:    []reflect.Type{tByteSlice},
			Received: reflect.Zero(t),
		}
	}

	var data []byte
	var err error
	switch vrType := vr.Type(); vrType {
	case bsontype.String:
		str, err := vr.ReadString()
		if err != nil {
			return emptyValue, err
		}
		data = []byte(str)
	case bsontype.Symbol:
		sym, err := vr.ReadSymbol()
		if err != nil {
			return emptyValue, err
		}
		data = []byte(sym)
	case bsontype.Binary:
		var subtype byte
		data, subtype, err = vr.ReadBinary()
		if err != nil {
			return emptyValue, err
		}
		if subtype != bsontype.BinaryGeneric && subtype != bsontype.BinaryBinaryOld {
			return emptyValue, decodeBinaryError{subtype: subtype, typeName: "[]byte"}
		}
	case bsontype.Null:
		err = vr.ReadNull()
	case bsontype.Undefined:
		err = vr.ReadUndefined()
	default:
		return emptyValue, fmt.Errorf("cannot decode %v into a []byte", vrType)
	}
	if err != nil {
		return emptyValue, err
	}

	return reflect.ValueOf(data), nil
}


func (bsc *ByteSliceCodec) DecodeValue(dc DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() || val.Type() != tByteSlice {
		return ValueDecoderError{Name: "ByteSliceDecodeValue", Types: []reflect.Type{tByteSlice}, Received: val}
	}

	elem, err := bsc.decodeType(dc, vr, tByteSlice)
	if err != nil {
		return err
	}

	val.Set(elem)
	return nil
}
