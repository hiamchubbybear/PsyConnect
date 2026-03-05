





package bson

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
)


var ErrDecodeToNil = errors.New("cannot Decode to nil value")




var decPool = sync.Pool{
	New: func() interface{} {
		return new(Decoder)
	},
}



type Decoder struct {
	dc bsoncodec.DecodeContext
	vr bsonrw.ValueReader

	
	
	defaultDocumentM bool
	defaultDocumentD bool

	binaryAsSlice     bool
	useJSONStructTags bool
	useLocalTimeZone  bool
	zeroMaps          bool
	zeroStructs       bool
}


func NewDecoder(vr bsonrw.ValueReader) (*Decoder, error) {
	if vr == nil {
		return nil, errors.New("cannot create a new Decoder with a nil ValueReader")
	}

	return &Decoder{
		dc: bsoncodec.DecodeContext{Registry: DefaultRegistry},
		vr: vr,
	}, nil
}





func NewDecoderWithContext(dc bsoncodec.DecodeContext, vr bsonrw.ValueReader) (*Decoder, error) {
	if dc.Registry == nil {
		dc.Registry = DefaultRegistry
	}
	if vr == nil {
		return nil, errors.New("cannot create a new Decoder with a nil ValueReader")
	}

	return &Decoder{
		dc: dc,
		vr: vr,
	}, nil
}





func (d *Decoder) Decode(val interface{}) error {
	if unmarshaler, ok := val.(Unmarshaler); ok {
		
		buf, err := bsonrw.Copier{}.CopyDocumentToBytes(d.vr)
		if err != nil {
			return err
		}
		return unmarshaler.UnmarshalBSON(buf)
	}

	rval := reflect.ValueOf(val)
	switch rval.Kind() {
	case reflect.Ptr:
		if rval.IsNil() {
			return ErrDecodeToNil
		}
		rval = rval.Elem()
	case reflect.Map:
		if rval.IsNil() {
			return ErrDecodeToNil
		}
	default:
		return fmt.Errorf("argument to Decode must be a pointer or a map, but got %v", rval)
	}
	decoder, err := d.dc.LookupDecoder(rval.Type())
	if err != nil {
		return err
	}

	if d.defaultDocumentM {
		d.dc.DefaultDocumentM()
	}
	if d.defaultDocumentD {
		d.dc.DefaultDocumentD()
	}
	if d.binaryAsSlice {
		d.dc.BinaryAsSlice()
	}
	if d.useJSONStructTags {
		d.dc.UseJSONStructTags()
	}
	if d.useLocalTimeZone {
		d.dc.UseLocalTimeZone()
	}
	if d.zeroMaps {
		d.dc.ZeroMaps()
	}
	if d.zeroStructs {
		d.dc.ZeroStructs()
	}

	return decoder.DecodeValue(d.dc, d.vr, rval)
}



func (d *Decoder) Reset(vr bsonrw.ValueReader) error {
	
	d.vr = vr
	return nil
}


func (d *Decoder) SetRegistry(r *bsoncodec.Registry) error {
	
	d.dc.Registry = r
	return nil
}




func (d *Decoder) SetContext(dc bsoncodec.DecodeContext) error {
	
	d.dc = dc
	return nil
}



func (d *Decoder) DefaultDocumentM() {
	d.defaultDocumentM = true
}



func (d *Decoder) DefaultDocumentD() {
	d.defaultDocumentD = true
}




func (d *Decoder) AllowTruncatingDoubles() {
	d.dc.Truncate = true
}



func (d *Decoder) BinaryAsSlice() {
	d.binaryAsSlice = true
}



func (d *Decoder) UseJSONStructTags() {
	d.useJSONStructTags = true
}



func (d *Decoder) UseLocalTimeZone() {
	d.useLocalTimeZone = true
}



func (d *Decoder) ZeroMaps() {
	d.zeroMaps = true
}



func (d *Decoder) ZeroStructs() {
	d.zeroStructs = true
}
