





package bson

import (
	"errors"
	"reflect"
	"sync"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
)




var encPool = sync.Pool{
	New: func() interface{} {
		return new(Encoder)
	},
}



type Encoder struct {
	ec bsoncodec.EncodeContext
	vw bsonrw.ValueWriter

	errorOnInlineDuplicates bool
	intMinSize              bool
	stringifyMapKeysWithFmt bool
	nilMapAsEmpty           bool
	nilSliceAsEmpty         bool
	nilByteSliceAsEmpty     bool
	omitZeroStruct          bool
	useJSONStructTags       bool
}


func NewEncoder(vw bsonrw.ValueWriter) (*Encoder, error) {
	
	if vw == nil {
		return nil, errors.New("cannot create a new Encoder with a nil ValueWriter")
	}

	return &Encoder{
		ec: bsoncodec.EncodeContext{Registry: DefaultRegistry},
		vw: vw,
	}, nil
}





func NewEncoderWithContext(ec bsoncodec.EncodeContext, vw bsonrw.ValueWriter) (*Encoder, error) {
	if ec.Registry == nil {
		ec = bsoncodec.EncodeContext{Registry: DefaultRegistry}
	}
	if vw == nil {
		return nil, errors.New("cannot create a new Encoder with a nil ValueWriter")
	}

	return &Encoder{
		ec: ec,
		vw: vw,
	}, nil
}




func (e *Encoder) Encode(val interface{}) error {
	if marshaler, ok := val.(Marshaler); ok {
		
		buf, err := marshaler.MarshalBSON()
		if err != nil {
			return err
		}
		return bsonrw.Copier{}.CopyDocumentFromBytes(e.vw, buf)
	}

	encoder, err := e.ec.LookupEncoder(reflect.TypeOf(val))
	if err != nil {
		return err
	}

	
	
	if e.errorOnInlineDuplicates {
		e.ec.ErrorOnInlineDuplicates()
	}
	if e.intMinSize {
		e.ec.MinSize = true
	}
	if e.stringifyMapKeysWithFmt {
		e.ec.StringifyMapKeysWithFmt()
	}
	if e.nilMapAsEmpty {
		e.ec.NilMapAsEmpty()
	}
	if e.nilSliceAsEmpty {
		e.ec.NilSliceAsEmpty()
	}
	if e.nilByteSliceAsEmpty {
		e.ec.NilByteSliceAsEmpty()
	}
	if e.omitZeroStruct {
		e.ec.OmitZeroStruct()
	}
	if e.useJSONStructTags {
		e.ec.UseJSONStructTags()
	}

	return encoder.EncodeValue(e.ec, e.vw, reflect.ValueOf(val))
}



func (e *Encoder) Reset(vw bsonrw.ValueWriter) error {
	
	e.vw = vw
	return nil
}


func (e *Encoder) SetRegistry(r *bsoncodec.Registry) error {
	
	e.ec.Registry = r
	return nil
}




func (e *Encoder) SetContext(ec bsoncodec.EncodeContext) error {
	
	e.ec = ec
	return nil
}



func (e *Encoder) ErrorOnInlineDuplicates() {
	e.errorOnInlineDuplicates = true
}




func (e *Encoder) IntMinSize() {
	e.intMinSize = true
}



func (e *Encoder) StringifyMapKeysWithFmt() {
	e.stringifyMapKeysWithFmt = true
}



func (e *Encoder) NilMapAsEmpty() {
	e.nilMapAsEmpty = true
}



func (e *Encoder) NilSliceAsEmpty() {
	e.nilSliceAsEmpty = true
}



func (e *Encoder) NilByteSliceAsEmpty() {
	e.nilByteSliceAsEmpty = true
}









func (e *Encoder) OmitZeroStruct() {
	e.omitZeroStruct = true
}



func (e *Encoder) UseJSONStructTags() {
	e.useJSONStructTags = true
}
