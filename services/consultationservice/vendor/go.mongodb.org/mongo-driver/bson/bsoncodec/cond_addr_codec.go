





package bsoncodec

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson/bsonrw"
)


type condAddrEncoder struct {
	canAddrEnc ValueEncoder
	elseEnc    ValueEncoder
}

var _ ValueEncoder = (*condAddrEncoder)(nil)


func newCondAddrEncoder(canAddrEnc, elseEnc ValueEncoder) *condAddrEncoder {
	encoder := condAddrEncoder{canAddrEnc: canAddrEnc, elseEnc: elseEnc}
	return &encoder
}


func (cae *condAddrEncoder) EncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if val.CanAddr() {
		return cae.canAddrEnc.EncodeValue(ec, vw, val)
	}
	if cae.elseEnc != nil {
		return cae.elseEnc.EncodeValue(ec, vw, val)
	}
	return ErrNoEncoder{Type: val.Type()}
}


type condAddrDecoder struct {
	canAddrDec ValueDecoder
	elseDec    ValueDecoder
}

var _ ValueDecoder = (*condAddrDecoder)(nil)


func newCondAddrDecoder(canAddrDec, elseDec ValueDecoder) *condAddrDecoder {
	decoder := condAddrDecoder{canAddrDec: canAddrDec, elseDec: elseDec}
	return &decoder
}


func (cad *condAddrDecoder) DecodeValue(dc DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if val.CanAddr() {
		return cad.canAddrDec.DecodeValue(dc, vr, val)
	}
	if cad.elseDec != nil {
		return cad.elseDec.DecodeValue(dc, vr, val)
	}
	return ErrNoDecoder{Type: val.Type()}
}
