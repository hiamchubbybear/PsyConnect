





package bsoncodec

import (
	"encoding"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/bsonoptions"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

var defaultMapCodec = NewMapCodec()



















type MapCodec struct {
	
	
	
	
	DecodeZerosMap bool

	
	
	
	
	EncodeNilAsEmpty bool

	
	
	
	
	
	EncodeKeysWithStringer bool
}



type KeyMarshaler interface {
	MarshalKey() (key string, err error)
}







type KeyUnmarshaler interface {
	UnmarshalKey(key string) error
}





func NewMapCodec(opts ...*bsonoptions.MapCodecOptions) *MapCodec {
	mapOpt := bsonoptions.MergeMapCodecOptions(opts...)

	codec := MapCodec{}
	if mapOpt.DecodeZerosMap != nil {
		codec.DecodeZerosMap = *mapOpt.DecodeZerosMap
	}
	if mapOpt.EncodeNilAsEmpty != nil {
		codec.EncodeNilAsEmpty = *mapOpt.EncodeNilAsEmpty
	}
	if mapOpt.EncodeKeysWithStringer != nil {
		codec.EncodeKeysWithStringer = *mapOpt.EncodeKeysWithStringer
	}
	return &codec
}


func (mc *MapCodec) EncodeValue(ec EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.Kind() != reflect.Map {
		return ValueEncoderError{Name: "MapEncodeValue", Kinds: []reflect.Kind{reflect.Map}, Received: val}
	}

	if val.IsNil() && !mc.EncodeNilAsEmpty && !ec.nilMapAsEmpty {
		
		
		
		
		
		err := vw.WriteNull()
		if err == nil {
			return nil
		}
	}

	dw, err := vw.WriteDocument()
	if err != nil {
		return err
	}

	return mc.mapEncodeValue(ec, dw, val, nil)
}




func (mc *MapCodec) mapEncodeValue(ec EncodeContext, dw bsonrw.DocumentWriter, val reflect.Value, collisionFn func(string) bool) error {

	elemType := val.Type().Elem()
	encoder, err := ec.LookupEncoder(elemType)
	if err != nil && elemType.Kind() != reflect.Interface {
		return err
	}

	keys := val.MapKeys()
	for _, key := range keys {
		keyStr, err := mc.encodeKey(key, ec.stringifyMapKeysWithFmt)
		if err != nil {
			return err
		}

		if collisionFn != nil && collisionFn(keyStr) {
			return fmt.Errorf("Key %s of inlined map conflicts with a struct field name", key)
		}

		currEncoder, currVal, lookupErr := defaultValueEncoders.lookupElementEncoder(ec, encoder, val.MapIndex(key))
		if lookupErr != nil && !errors.Is(lookupErr, errInvalidValue) {
			return lookupErr
		}

		vw, err := dw.WriteDocumentElement(keyStr)
		if err != nil {
			return err
		}

		if errors.Is(lookupErr, errInvalidValue) {
			err = vw.WriteNull()
			if err != nil {
				return err
			}
			continue
		}

		err = currEncoder.EncodeValue(ec, vw, currVal)
		if err != nil {
			return err
		}
	}

	return dw.WriteDocumentEnd()
}


func (mc *MapCodec) DecodeValue(dc DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if val.Kind() != reflect.Map || (!val.CanSet() && val.IsNil()) {
		return ValueDecoderError{Name: "MapDecodeValue", Kinds: []reflect.Kind{reflect.Map}, Received: val}
	}

	switch vrType := vr.Type(); vrType {
	case bsontype.Type(0), bsontype.EmbeddedDocument:
	case bsontype.Null:
		val.Set(reflect.Zero(val.Type()))
		return vr.ReadNull()
	case bsontype.Undefined:
		val.Set(reflect.Zero(val.Type()))
		return vr.ReadUndefined()
	default:
		return fmt.Errorf("cannot decode %v into a %s", vrType, val.Type())
	}

	dr, err := vr.ReadDocument()
	if err != nil {
		return err
	}

	if val.IsNil() {
		val.Set(reflect.MakeMap(val.Type()))
	}

	if val.Len() > 0 && (mc.DecodeZerosMap || dc.zeroMaps) {
		clearMap(val)
	}

	eType := val.Type().Elem()
	decoder, err := dc.LookupDecoder(eType)
	if err != nil {
		return err
	}
	eTypeDecoder, _ := decoder.(typeDecoder)

	if eType == tEmpty {
		dc.Ancestor = val.Type()
	}

	keyType := val.Type().Key()

	for {
		key, vr, err := dr.ReadElement()
		if errors.Is(err, bsonrw.ErrEOD) {
			break
		}
		if err != nil {
			return err
		}

		k, err := mc.decodeKey(key, keyType)
		if err != nil {
			return err
		}

		elem, err := decodeTypeOrValueWithInfo(decoder, eTypeDecoder, dc, vr, eType, true)
		if err != nil {
			return newDecodeError(key, err)
		}

		val.SetMapIndex(k, elem)
	}
	return nil
}

func clearMap(m reflect.Value) {
	var none reflect.Value
	for _, k := range m.MapKeys() {
		m.SetMapIndex(k, none)
	}
}

func (mc *MapCodec) encodeKey(val reflect.Value, encodeKeysWithStringer bool) (string, error) {
	if mc.EncodeKeysWithStringer || encodeKeysWithStringer {
		return fmt.Sprint(val), nil
	}

	
	if val.Kind() == reflect.String {
		return val.String(), nil
	}
	
	if km, ok := val.Interface().(KeyMarshaler); ok {
		if val.Kind() == reflect.Ptr && val.IsNil() {
			return "", nil
		}
		buf, err := km.MarshalKey()
		if err == nil {
			return buf, nil
		}
		return "", err
	}
	
	if km, ok := val.Interface().(encoding.TextMarshaler); ok {
		if val.Kind() == reflect.Ptr && val.IsNil() {
			return "", nil
		}

		buf, err := km.MarshalText()
		if err != nil {
			return "", err
		}

		return string(buf), nil
	}

	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(val.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(val.Uint(), 10), nil
	}
	return "", fmt.Errorf("unsupported key type: %v", val.Type())
}

var keyUnmarshalerType = reflect.TypeOf((*KeyUnmarshaler)(nil)).Elem()
var textUnmarshalerType = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()

func (mc *MapCodec) decodeKey(key string, keyType reflect.Type) (reflect.Value, error) {
	keyVal := reflect.ValueOf(key)
	var err error
	switch {
	
	case !mc.EncodeKeysWithStringer && reflect.PtrTo(keyType).Implements(keyUnmarshalerType):
		keyVal = reflect.New(keyType)
		v := keyVal.Interface().(KeyUnmarshaler)
		err = v.UnmarshalKey(key)
		keyVal = keyVal.Elem()
	
	case reflect.PtrTo(keyType).Implements(textUnmarshalerType):
		keyVal = reflect.New(keyType)
		v := keyVal.Interface().(encoding.TextUnmarshaler)
		err = v.UnmarshalText([]byte(key))
		keyVal = keyVal.Elem()
	
	default:
		switch keyType.Kind() {
		case reflect.String:
			keyVal = reflect.ValueOf(key).Convert(keyType)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			n, parseErr := strconv.ParseInt(key, 10, 64)
			if parseErr != nil || reflect.Zero(keyType).OverflowInt(n) {
				err = fmt.Errorf("failed to unmarshal number key %v", key)
			}
			keyVal = reflect.ValueOf(n).Convert(keyType)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			n, parseErr := strconv.ParseUint(key, 10, 64)
			if parseErr != nil || reflect.Zero(keyType).OverflowUint(n) {
				err = fmt.Errorf("failed to unmarshal number key %v", key)
				break
			}
			keyVal = reflect.ValueOf(n).Convert(keyType)
		case reflect.Float32, reflect.Float64:
			if mc.EncodeKeysWithStringer {
				parsed, err := strconv.ParseFloat(key, 64)
				if err != nil {
					return keyVal, fmt.Errorf("Map key is defined to be a decimal type (%v) but got error %w", keyType.Kind(), err)
				}
				keyVal = reflect.ValueOf(parsed)
				break
			}
			fallthrough
		default:
			return keyVal, fmt.Errorf("unsupported key type: %v", keyType)
		}
	}
	return keyVal, err
}
