





package bsoncodec

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson/bsonoptions"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)



















type EmptyInterfaceCodec struct {
	
	
	
	
	DecodeBinaryAsSlice bool
}

var (
	defaultEmptyInterfaceCodec = NewEmptyInterfaceCodec()

	
	
	
	_ typeDecoder = defaultEmptyInterfaceCodec
)





func NewEmptyInterfaceCodec(opts ...*bsonoptions.EmptyInterfaceCodecOptions) *EmptyInterfaceCodec {
	interfaceOpt := bsonoptions.MergeEmptyInterfaceCodecOptions(opts...)

	codec := EmptyInterfaceCodec{}
	if interfaceOpt.DecodeBinaryAsSlice != nil {
		codec.DecodeBinaryAsSlice = *interfaceOpt.DecodeBinaryAsSlice
	}
	return &codec
}


func (eic EmptyInterfaceCodec) EncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.Type() != tEmpty {
		return ValueEncoderError{Name: "EmptyInterfaceEncodeValue", Types: []reflect.Type{tEmpty}, Received: val}
	}

	if val.IsNil() {
		return vw.WriteNull()
	}
	encoder, err := ec.LookupEncoder(val.Elem().Type())
	if err != nil {
		return err
	}

	return encoder.EncodeValue(ec, vw, val.Elem())
}

func (eic EmptyInterfaceCodec) getEmptyInterfaceDecodeType(dc DecodeContext, valueType bsontype.Type) (reflect.Type, error) {
	isDocument := valueType == bsontype.Type(0) || valueType == bsontype.EmbeddedDocument
	if isDocument {
		if dc.defaultDocumentType != nil {
			
			
			return dc.defaultDocumentType, nil
		}
		if dc.Ancestor != nil {
			
			
			
			return dc.Ancestor, nil
		}
	}

	rtype, err := dc.LookupTypeMapEntry(valueType)
	if err == nil {
		return rtype, nil
	}

	if isDocument {
		
		
		var lookupType bsontype.Type
		switch valueType {
		case bsontype.Type(0):
			lookupType = bsontype.EmbeddedDocument
		case bsontype.EmbeddedDocument:
			lookupType = bsontype.Type(0)
		}

		rtype, err = dc.LookupTypeMapEntry(lookupType)
		if err == nil {
			return rtype, nil
		}
	}

	return nil, err
}

func (eic EmptyInterfaceCodec) decodeType(dc DecodeContext, vr bsonrw.ValueReader, t reflect.Type) (reflect.Value, error) {
	if t != tEmpty {
		return emptyValue, ValueDecoderError{Name: "EmptyInterfaceDecodeValue", Types: []reflect.Type{tEmpty}, Received: reflect.Zero(t)}
	}

	rtype, err := eic.getEmptyInterfaceDecodeType(dc, vr.Type())
	if err != nil {
		switch vr.Type() {
		case bsontype.Null:
			return reflect.Zero(t), vr.ReadNull()
		default:
			return emptyValue, err
		}
	}

	decoder, err := dc.LookupDecoder(rtype)
	if err != nil {
		return emptyValue, err
	}

	elem, err := decodeTypeOrValue(decoder, dc, vr, rtype)
	if err != nil {
		return emptyValue, err
	}

	if (eic.DecodeBinaryAsSlice || dc.binaryAsSlice) && rtype == tBinary {
		binElem := elem.Interface().(primitive.Binary)
		if binElem.Subtype == bsontype.BinaryGeneric || binElem.Subtype == bsontype.BinaryBinaryOld {
			elem = reflect.ValueOf(binElem.Data)
		}
	}

	return elem, nil
}


func (eic EmptyInterfaceCodec) DecodeValue(dc DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() || val.Type() != tEmpty {
		return ValueDecoderError{Name: "EmptyInterfaceDecodeValue", Types: []reflect.Type{tEmpty}, Received: val}
	}

	elem, err := eic.decodeType(dc, vr, val.Type())
	if err != nil {
		return err
	}

	val.Set(elem)
	return nil
}
