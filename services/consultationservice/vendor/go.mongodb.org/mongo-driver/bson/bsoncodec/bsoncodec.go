





package bsoncodec 

import (
	"fmt"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	emptyValue = reflect.Value{}
)






type Marshaler interface {
	MarshalBSON() ([]byte, error)
}







type ValueMarshaler interface {
	MarshalBSONValue() (bsontype.Type, []byte, error)
}







type Unmarshaler interface {
	UnmarshalBSON([]byte) error
}







type ValueUnmarshaler interface {
	UnmarshalBSONValue(bsontype.Type, []byte) error
}



type ValueEncoderError struct {
	Name     string
	Types    []reflect.Type
	Kinds    []reflect.Kind
	Received reflect.Value
}

func (vee ValueEncoderError) Error() string {
	typeKinds := make([]string, 0, len(vee.Types)+len(vee.Kinds))
	for _, t := range vee.Types {
		typeKinds = append(typeKinds, t.String())
	}
	for _, k := range vee.Kinds {
		if k == reflect.Map {
			typeKinds = append(typeKinds, "map[string]*")
			continue
		}
		typeKinds = append(typeKinds, k.String())
	}
	received := vee.Received.Kind().String()
	if vee.Received.IsValid() {
		received = vee.Received.Type().String()
	}
	return fmt.Sprintf("%s can only encode valid %s, but got %s", vee.Name, strings.Join(typeKinds, ", "), received)
}



type ValueDecoderError struct {
	Name     string
	Types    []reflect.Type
	Kinds    []reflect.Kind
	Received reflect.Value
}

func (vde ValueDecoderError) Error() string {
	typeKinds := make([]string, 0, len(vde.Types)+len(vde.Kinds))
	for _, t := range vde.Types {
		typeKinds = append(typeKinds, t.String())
	}
	for _, k := range vde.Kinds {
		if k == reflect.Map {
			typeKinds = append(typeKinds, "map[string]*")
			continue
		}
		typeKinds = append(typeKinds, k.String())
	}
	received := vde.Received.Kind().String()
	if vde.Received.IsValid() {
		received = vde.Received.Type().String()
	}
	return fmt.Sprintf("%s can only decode valid and settable %s, but got %s", vde.Name, strings.Join(typeKinds, ", "), received)
}



type EncodeContext struct {
	*Registry

	
	
	
	
	
	MinSize bool

	errorOnInlineDuplicates bool
	stringifyMapKeysWithFmt bool
	nilMapAsEmpty           bool
	nilSliceAsEmpty         bool
	nilByteSliceAsEmpty     bool
	omitZeroStruct          bool
	useJSONStructTags       bool
}





func (ec *EncodeContext) ErrorOnInlineDuplicates() {
	ec.errorOnInlineDuplicates = true
}





func (ec *EncodeContext) StringifyMapKeysWithFmt() {
	ec.stringifyMapKeysWithFmt = true
}





func (ec *EncodeContext) NilMapAsEmpty() {
	ec.nilMapAsEmpty = true
}





func (ec *EncodeContext) NilSliceAsEmpty() {
	ec.nilSliceAsEmpty = true
}





func (ec *EncodeContext) NilByteSliceAsEmpty() {
	ec.nilByteSliceAsEmpty = true
}








func (ec *EncodeContext) OmitZeroStruct() {
	ec.omitZeroStruct = true
}





func (ec *EncodeContext) UseJSONStructTags() {
	ec.useJSONStructTags = true
}



type DecodeContext struct {
	*Registry

	
	
	
	
	
	
	Truncate bool

	
	
	
	
	
	
	Ancestor reflect.Type

	
	
	
	
	defaultDocumentType reflect.Type

	binaryAsSlice     bool
	useJSONStructTags bool
	useLocalTimeZone  bool
	zeroMaps          bool
	zeroStructs       bool
}





func (dc *DecodeContext) BinaryAsSlice() {
	dc.binaryAsSlice = true
}





func (dc *DecodeContext) UseJSONStructTags() {
	dc.useJSONStructTags = true
}





func (dc *DecodeContext) UseLocalTimeZone() {
	dc.useLocalTimeZone = true
}





func (dc *DecodeContext) ZeroMaps() {
	dc.zeroMaps = true
}





func (dc *DecodeContext) ZeroStructs() {
	dc.zeroStructs = true
}





func (dc *DecodeContext) DefaultDocumentM() {
	dc.defaultDocumentType = reflect.TypeOf(primitive.M{})
}





func (dc *DecodeContext) DefaultDocumentD() {
	dc.defaultDocumentType = reflect.TypeOf(primitive.D{})
}





type ValueCodec interface {
	ValueEncoder
	ValueDecoder
}


type ValueEncoder interface {
	EncodeValue(EncodeContext, bsonrw.ValueWriter, reflect.Value) error
}



type ValueEncoderFunc func(EncodeContext, bsonrw.ValueWriter, reflect.Value) error


func (fn ValueEncoderFunc) EncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	return fn(ec, vw, val)
}


type ValueDecoder interface {
	DecodeValue(DecodeContext, bsonrw.ValueReader, reflect.Value) error
}



type ValueDecoderFunc func(DecodeContext, bsonrw.ValueReader, reflect.Value) error


func (fn ValueDecoderFunc) DecodeValue(dc DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	return fn(dc, vr, val)
}


type typeDecoder interface {
	decodeType(DecodeContext, bsonrw.ValueReader, reflect.Type) (reflect.Value, error)
}


type typeDecoderFunc func(DecodeContext, bsonrw.ValueReader, reflect.Type) (reflect.Value, error)

func (fn typeDecoderFunc) decodeType(dc DecodeContext, vr bsonrw.ValueReader, t reflect.Type) (reflect.Value, error) {
	return fn(dc, vr, t)
}


type decodeAdapter struct {
	ValueDecoderFunc
	typeDecoderFunc
}

var _ ValueDecoder = decodeAdapter{}
var _ typeDecoder = decodeAdapter{}



func decodeTypeOrValue(decoder ValueDecoder, dc DecodeContext, vr bsonrw.ValueReader, t reflect.Type) (reflect.Value, error) {
	td, _ := decoder.(typeDecoder)
	return decodeTypeOrValueWithInfo(decoder, td, dc, vr, t, true)
}

func decodeTypeOrValueWithInfo(vd ValueDecoder, td typeDecoder, dc DecodeContext, vr bsonrw.ValueReader, t reflect.Type, convert bool) (reflect.Value, error) {
	if td != nil {
		val, err := td.decodeType(dc, vr, t)
		if err == nil && convert && val.Type() != t {
			
			
			
			
			
			
			
			val = val.Convert(t)
		}
		return val, err
	}

	val := reflect.New(t).Elem()
	err := vd.DecodeValue(dc, vr, val)
	return val, err
}







type CodecZeroer interface {
	IsTypeZero(interface{}) bool
}
