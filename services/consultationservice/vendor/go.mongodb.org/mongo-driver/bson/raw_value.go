





package bson

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)


var ErrNilContext = errors.New("DecodeContext cannot be nil")


var ErrNilRegistry = errors.New("Registry cannot be nil")





type RawValue struct {
	Type  bsontype.Type
	Value []byte

	r *bsoncodec.Registry
}



func (rv RawValue) IsZero() bool {
	return rv.Type == 0x00 && len(rv.Value) == 0
}





func (rv RawValue) Unmarshal(val interface{}) error {
	reg := rv.r
	if reg == nil {
		reg = DefaultRegistry
	}
	return rv.UnmarshalWithRegistry(reg, val)
}


func (rv RawValue) Equal(rv2 RawValue) bool {
	if rv.Type != rv2.Type {
		return false
	}

	if !bytes.Equal(rv.Value, rv2.Value) {
		return false
	}

	return true
}



func (rv RawValue) UnmarshalWithRegistry(r *bsoncodec.Registry, val interface{}) error {
	if r == nil {
		return ErrNilRegistry
	}

	vr := bsonrw.NewBSONValueReader(rv.Type, rv.Value)
	rval := reflect.ValueOf(val)
	if rval.Kind() != reflect.Ptr {
		return fmt.Errorf("argument to Unmarshal* must be a pointer to a type, but got %v", rval)
	}
	rval = rval.Elem()
	dec, err := r.LookupDecoder(rval.Type())
	if err != nil {
		return err
	}
	return dec.DecodeValue(bsoncodec.DecodeContext{Registry: r}, vr, rval)
}







func (rv RawValue) UnmarshalWithContext(dc *bsoncodec.DecodeContext, val interface{}) error {
	if dc == nil {
		return ErrNilContext
	}

	vr := bsonrw.NewBSONValueReader(rv.Type, rv.Value)
	rval := reflect.ValueOf(val)
	if rval.Kind() != reflect.Ptr {
		return fmt.Errorf("argument to Unmarshal* must be a pointer to a type, but got %v", rval)
	}
	rval = rval.Elem()
	dec, err := dc.LookupDecoder(rval.Type())
	if err != nil {
		return err
	}
	return dec.DecodeValue(*dc, vr, rval)
}

func convertFromCoreValue(v bsoncore.Value) RawValue { return RawValue{Type: v.Type, Value: v.Data} }
func convertToCoreValue(v RawValue) bsoncore.Value {
	return bsoncore.Value{Type: v.Type, Data: v.Value}
}


func (rv RawValue) Validate() error { return convertToCoreValue(rv).Validate() }


func (rv RawValue) IsNumber() bool { return convertToCoreValue(rv).IsNumber() }



func (rv RawValue) String() string { return convertToCoreValue(rv).String() }



func (rv RawValue) DebugString() string { return convertToCoreValue(rv).DebugString() }



func (rv RawValue) Double() float64 { return convertToCoreValue(rv).Double() }


func (rv RawValue) DoubleOK() (float64, bool) { return convertToCoreValue(rv).DoubleOK() }






func (rv RawValue) StringValue() string { return convertToCoreValue(rv).StringValue() }



func (rv RawValue) StringValueOK() (string, bool) { return convertToCoreValue(rv).StringValueOK() }



func (rv RawValue) Document() Raw { return Raw(convertToCoreValue(rv).Document()) }



func (rv RawValue) DocumentOK() (Raw, bool) {
	doc, ok := convertToCoreValue(rv).DocumentOK()
	return Raw(doc), ok
}



func (rv RawValue) Array() Raw { return Raw(convertToCoreValue(rv).Array()) }



func (rv RawValue) ArrayOK() (Raw, bool) {
	doc, ok := convertToCoreValue(rv).ArrayOK()
	return Raw(doc), ok
}



func (rv RawValue) Binary() (subtype byte, data []byte) { return convertToCoreValue(rv).Binary() }



func (rv RawValue) BinaryOK() (subtype byte, data []byte, ok bool) {
	return convertToCoreValue(rv).BinaryOK()
}



func (rv RawValue) ObjectID() primitive.ObjectID { return convertToCoreValue(rv).ObjectID() }



func (rv RawValue) ObjectIDOK() (primitive.ObjectID, bool) {
	return convertToCoreValue(rv).ObjectIDOK()
}



func (rv RawValue) Boolean() bool { return convertToCoreValue(rv).Boolean() }



func (rv RawValue) BooleanOK() (bool, bool) { return convertToCoreValue(rv).BooleanOK() }



func (rv RawValue) DateTime() int64 { return convertToCoreValue(rv).DateTime() }



func (rv RawValue) DateTimeOK() (int64, bool) { return convertToCoreValue(rv).DateTimeOK() }



func (rv RawValue) Time() time.Time { return convertToCoreValue(rv).Time() }



func (rv RawValue) TimeOK() (time.Time, bool) { return convertToCoreValue(rv).TimeOK() }



func (rv RawValue) Regex() (pattern, options string) { return convertToCoreValue(rv).Regex() }



func (rv RawValue) RegexOK() (pattern, options string, ok bool) {
	return convertToCoreValue(rv).RegexOK()
}



func (rv RawValue) DBPointer() (string, primitive.ObjectID) {
	return convertToCoreValue(rv).DBPointer()
}



func (rv RawValue) DBPointerOK() (string, primitive.ObjectID, bool) {
	return convertToCoreValue(rv).DBPointerOK()
}



func (rv RawValue) JavaScript() string { return convertToCoreValue(rv).JavaScript() }



func (rv RawValue) JavaScriptOK() (string, bool) { return convertToCoreValue(rv).JavaScriptOK() }



func (rv RawValue) Symbol() string { return convertToCoreValue(rv).Symbol() }



func (rv RawValue) SymbolOK() (string, bool) { return convertToCoreValue(rv).SymbolOK() }



func (rv RawValue) CodeWithScope() (string, Raw) {
	code, scope := convertToCoreValue(rv).CodeWithScope()
	return code, Raw(scope)
}



func (rv RawValue) CodeWithScopeOK() (string, Raw, bool) {
	code, scope, ok := convertToCoreValue(rv).CodeWithScopeOK()
	return code, Raw(scope), ok
}



func (rv RawValue) Int32() int32 { return convertToCoreValue(rv).Int32() }



func (rv RawValue) Int32OK() (int32, bool) { return convertToCoreValue(rv).Int32OK() }






func (rv RawValue) AsInt32() int32 { return convertToCoreValue(rv).AsInt32() }






func (rv RawValue) AsInt32OK() (int32, bool) { return convertToCoreValue(rv).AsInt32OK() }



func (rv RawValue) Timestamp() (t, i uint32) { return convertToCoreValue(rv).Timestamp() }



func (rv RawValue) TimestampOK() (t, i uint32, ok bool) { return convertToCoreValue(rv).TimestampOK() }



func (rv RawValue) Int64() int64 { return convertToCoreValue(rv).Int64() }



func (rv RawValue) Int64OK() (int64, bool) { return convertToCoreValue(rv).Int64OK() }



func (rv RawValue) AsInt64() int64 { return convertToCoreValue(rv).AsInt64() }



func (rv RawValue) AsInt64OK() (int64, bool) { return convertToCoreValue(rv).AsInt64OK() }



func (rv RawValue) Decimal128() primitive.Decimal128 { return convertToCoreValue(rv).Decimal128() }



func (rv RawValue) Decimal128OK() (primitive.Decimal128, bool) {
	return convertToCoreValue(rv).Decimal128OK()
}
