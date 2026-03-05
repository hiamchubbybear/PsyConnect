





package bsoncodec

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

var _ ValueEncoder = &PointerCodec{}
var _ ValueDecoder = &PointerCodec{}













type PointerCodec struct {
	ecache typeEncoderCache
	dcache typeDecoderCache
}





func NewPointerCodec() *PointerCodec {
	return &PointerCodec{}
}



func (pc *PointerCodec) EncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if val.Kind() != reflect.Ptr {
		if !val.IsValid() {
			return vw.WriteNull()
		}
		return ValueEncoderError{Name: "PointerCodec.EncodeValue", Kinds: []reflect.Kind{reflect.Ptr}, Received: val}
	}

	if val.IsNil() {
		return vw.WriteNull()
	}

	typ := val.Type()
	if v, ok := pc.ecache.Load(typ); ok {
		if v == nil {
			return ErrNoEncoder{Type: typ}
		}
		return v.EncodeValue(ec, vw, val.Elem())
	}
	
	enc, err := ec.LookupEncoder(typ.Elem())
	enc = pc.ecache.LoadOrStore(typ, enc)
	if err != nil {
		return err
	}
	return enc.EncodeValue(ec, vw, val.Elem())
}



func (pc *PointerCodec) DecodeValue(dc DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() || val.Kind() != reflect.Ptr {
		return ValueDecoderError{Name: "PointerCodec.DecodeValue", Kinds: []reflect.Kind{reflect.Ptr}, Received: val}
	}

	typ := val.Type()
	if vr.Type() == bsontype.Null {
		val.Set(reflect.Zero(typ))
		return vr.ReadNull()
	}
	if vr.Type() == bsontype.Undefined {
		val.Set(reflect.Zero(typ))
		return vr.ReadUndefined()
	}

	if val.IsNil() {
		val.Set(reflect.New(typ.Elem()))
	}

	if v, ok := pc.dcache.Load(typ); ok {
		if v == nil {
			return ErrNoDecoder{Type: typ}
		}
		return v.DecodeValue(dc, vr, val.Elem())
	}
	
	dec, err := dc.LookupDecoder(typ.Elem())
	dec = pc.dcache.LoadOrStore(typ, dec)
	if err != nil {
		return err
	}
	return dec.DecodeValue(dc, vr, val.Elem())
}
