





























































































































































package mapstructure

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/go-viper/mapstructure/v2/internal/errors"
)















type DecodeHookFunc interface{}



type DecodeHookFuncType func(reflect.Type, reflect.Type, interface{}) (interface{}, error)



type DecodeHookFuncKind func(reflect.Kind, reflect.Kind, interface{}) (interface{}, error)



type DecodeHookFuncValue func(from reflect.Value, to reflect.Value) (interface{}, error)



type DecoderConfig struct {
	
	
	
	
	
	
	
	
	
	DecodeHook DecodeHookFunc

	
	
	
	ErrorUnused bool

	
	
	
	
	ErrorUnset bool

	
	
	
	ZeroFields bool

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	WeaklyTypedInput bool

	
	
	
	
	
	
	Squash bool

	
	
	Metadata *Metadata

	
	
	Result interface{}

	
	
	TagName string

	
	
	SquashTagOption string

	
	
	IgnoreUntaggedFields bool

	
	
	
	MatchName func(mapKey, fieldName string) bool

	
	
	DecodeNil bool
}







type Decoder struct {
	config           *DecoderConfig
	cachedDecodeHook func(from reflect.Value, to reflect.Value) (interface{}, error)
}



type Metadata struct {
	
	Keys []string

	
	
	Unused []string

	
	
	
	Unset []string
}



func Decode(input interface{}, output interface{}) error {
	config := &DecoderConfig{
		Metadata: nil,
		Result:   output,
	}

	decoder, err := NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}



func WeakDecode(input, output interface{}) error {
	config := &DecoderConfig{
		Metadata:         nil,
		Result:           output,
		WeaklyTypedInput: true,
	}

	decoder, err := NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}



func DecodeMetadata(input interface{}, output interface{}, metadata *Metadata) error {
	config := &DecoderConfig{
		Metadata: metadata,
		Result:   output,
	}

	decoder, err := NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}




func WeakDecodeMetadata(input interface{}, output interface{}, metadata *Metadata) error {
	config := &DecoderConfig{
		Metadata:         metadata,
		Result:           output,
		WeaklyTypedInput: true,
	}

	decoder, err := NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}




func NewDecoder(config *DecoderConfig) (*Decoder, error) {
	val := reflect.ValueOf(config.Result)
	if val.Kind() != reflect.Ptr {
		return nil, errors.New("result must be a pointer")
	}

	val = val.Elem()
	if !val.CanAddr() {
		return nil, errors.New("result must be addressable (a pointer)")
	}

	if config.Metadata != nil {
		if config.Metadata.Keys == nil {
			config.Metadata.Keys = make([]string, 0)
		}

		if config.Metadata.Unused == nil {
			config.Metadata.Unused = make([]string, 0)
		}

		if config.Metadata.Unset == nil {
			config.Metadata.Unset = make([]string, 0)
		}
	}

	if config.TagName == "" {
		config.TagName = "mapstructure"
	}

	if config.SquashTagOption == "" {
		config.SquashTagOption = "squash"
	}

	if config.MatchName == nil {
		config.MatchName = strings.EqualFold
	}

	result := &Decoder{
		config: config,
	}
	if config.DecodeHook != nil {
		result.cachedDecodeHook = cachedDecodeHook(config.DecodeHook)
	}

	return result, nil
}



func (d *Decoder) Decode(input interface{}) error {
	err := d.decode("", input, reflect.ValueOf(d.config.Result).Elem())

	
	var joinedErr interface{ Unwrap() []error }
	if errors.As(err, &joinedErr) {
		return fmt.Errorf("decoding failed due to the following error(s):\n\n%w", err)
	}

	return err
}


func isNil(input interface{}) bool {
	if input == nil {
		return true
	}
	val := reflect.ValueOf(input)
	return val.Kind() == reflect.Ptr && val.IsNil()
}


func (d *Decoder) decode(name string, input interface{}, outVal reflect.Value) error {
	var (
		inputVal   = reflect.ValueOf(input)
		outputKind = getKind(outVal)
		decodeNil  = d.config.DecodeNil && d.cachedDecodeHook != nil
	)
	if isNil(input) {
		
		input = nil
	}
	if input == nil {
		
		
		if d.config.ZeroFields {
			outVal.Set(reflect.Zero(outVal.Type()))

			if d.config.Metadata != nil && name != "" {
				d.config.Metadata.Keys = append(d.config.Metadata.Keys, name)
			}
		}
		if !decodeNil {
			return nil
		}
	}
	if !inputVal.IsValid() {
		if !decodeNil {
			
			
			outVal.Set(reflect.Zero(outVal.Type()))
			if d.config.Metadata != nil && name != "" {
				d.config.Metadata.Keys = append(d.config.Metadata.Keys, name)
			}
			return nil
		}
		
		switch outputKind {
		case reflect.Struct, reflect.Map:
			var mapVal map[string]interface{}
			inputVal = reflect.ValueOf(mapVal) 
		case reflect.Slice, reflect.Array:
			var sliceVal []interface{}
			inputVal = reflect.ValueOf(sliceVal) 
		default:
			inputVal = reflect.Zero(outVal.Type())
		}
	}

	if d.cachedDecodeHook != nil {
		
		var err error
		input, err = d.cachedDecodeHook(inputVal, outVal)
		if err != nil {
			return fmt.Errorf("error decoding '%s': %w", name, err)
		}
	}
	if isNil(input) {
		return nil
	}

	var err error
	addMetaKey := true
	switch outputKind {
	case reflect.Bool:
		err = d.decodeBool(name, input, outVal)
	case reflect.Interface:
		err = d.decodeBasic(name, input, outVal)
	case reflect.String:
		err = d.decodeString(name, input, outVal)
	case reflect.Int:
		err = d.decodeInt(name, input, outVal)
	case reflect.Uint:
		err = d.decodeUint(name, input, outVal)
	case reflect.Float32:
		err = d.decodeFloat(name, input, outVal)
	case reflect.Complex64:
		err = d.decodeComplex(name, input, outVal)
	case reflect.Struct:
		err = d.decodeStruct(name, input, outVal)
	case reflect.Map:
		err = d.decodeMap(name, input, outVal)
	case reflect.Ptr:
		addMetaKey, err = d.decodePtr(name, input, outVal)
	case reflect.Slice:
		err = d.decodeSlice(name, input, outVal)
	case reflect.Array:
		err = d.decodeArray(name, input, outVal)
	case reflect.Func:
		err = d.decodeFunc(name, input, outVal)
	default:
		
		return fmt.Errorf("%s: unsupported type: %s", name, outputKind)
	}

	
	
	if addMetaKey && d.config.Metadata != nil && name != "" {
		d.config.Metadata.Keys = append(d.config.Metadata.Keys, name)
	}

	return err
}



func (d *Decoder) decodeBasic(name string, data interface{}, val reflect.Value) error {
	if val.IsValid() && val.Elem().IsValid() {
		elem := val.Elem()

		
		
		
		copied := false
		if !elem.CanAddr() {
			copied = true

			
			copy := reflect.New(elem.Type())

			
			copy.Elem().Set(elem)

			
			elem = copy
		}

		
		
		if err := d.decode(name, data, elem); err != nil || !copied {
			return err
		}

		
		val.Set(elem.Elem())
		return nil
	}

	dataVal := reflect.ValueOf(data)

	
	
	
	if dataVal.Kind() == reflect.Ptr && dataVal.Type().Elem() == val.Type() {
		dataVal = reflect.Indirect(dataVal)
	}

	if !dataVal.IsValid() {
		dataVal = reflect.Zero(val.Type())
	}

	dataValType := dataVal.Type()
	if !dataValType.AssignableTo(val.Type()) {
		return fmt.Errorf(
			"'%s' expected type '%s', got '%s'",
			name, val.Type(), dataValType)
	}

	val.Set(dataVal)
	return nil
}

func (d *Decoder) decodeString(name string, data interface{}, val reflect.Value) error {
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataKind := getKind(dataVal)

	converted := true
	switch {
	case dataKind == reflect.String:
		val.SetString(dataVal.String())
	case dataKind == reflect.Bool && d.config.WeaklyTypedInput:
		if dataVal.Bool() {
			val.SetString("1")
		} else {
			val.SetString("0")
		}
	case dataKind == reflect.Int && d.config.WeaklyTypedInput:
		val.SetString(strconv.FormatInt(dataVal.Int(), 10))
	case dataKind == reflect.Uint && d.config.WeaklyTypedInput:
		val.SetString(strconv.FormatUint(dataVal.Uint(), 10))
	case dataKind == reflect.Float32 && d.config.WeaklyTypedInput:
		val.SetString(strconv.FormatFloat(dataVal.Float(), 'f', -1, 64))
	case dataKind == reflect.Slice && d.config.WeaklyTypedInput,
		dataKind == reflect.Array && d.config.WeaklyTypedInput:
		dataType := dataVal.Type()
		elemKind := dataType.Elem().Kind()
		switch elemKind {
		case reflect.Uint8:
			var uints []uint8
			if dataKind == reflect.Array {
				uints = make([]uint8, dataVal.Len(), dataVal.Len())
				for i := range uints {
					uints[i] = dataVal.Index(i).Interface().(uint8)
				}
			} else {
				uints = dataVal.Interface().([]uint8)
			}
			val.SetString(string(uints))
		default:
			converted = false
		}
	default:
		converted = false
	}

	if !converted {
		return fmt.Errorf(
			"'%s' expected type '%s', got unconvertible type '%s', value: '%v'",
			name, val.Type(), dataVal.Type(), data)
	}

	return nil
}

func (d *Decoder) decodeInt(name string, data interface{}, val reflect.Value) error {
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataKind := getKind(dataVal)
	dataType := dataVal.Type()

	switch {
	case dataKind == reflect.Int:
		val.SetInt(dataVal.Int())
	case dataKind == reflect.Uint:
		val.SetInt(int64(dataVal.Uint()))
	case dataKind == reflect.Float32:
		val.SetInt(int64(dataVal.Float()))
	case dataKind == reflect.Bool && d.config.WeaklyTypedInput:
		if dataVal.Bool() {
			val.SetInt(1)
		} else {
			val.SetInt(0)
		}
	case dataKind == reflect.String && d.config.WeaklyTypedInput:
		str := dataVal.String()
		if str == "" {
			str = "0"
		}

		i, err := strconv.ParseInt(str, 0, val.Type().Bits())
		if err == nil {
			val.SetInt(i)
		} else {
			return fmt.Errorf("cannot parse '%s' as int: %s", name, err)
		}
	case dataType.PkgPath() == "encoding/json" && dataType.Name() == "Number":
		jn := data.(json.Number)
		i, err := jn.Int64()
		if err != nil {
			return fmt.Errorf(
				"error decoding json.Number into %s: %s", name, err)
		}
		val.SetInt(i)
	default:
		return fmt.Errorf(
			"'%s' expected type '%s', got unconvertible type '%s', value: '%v'",
			name, val.Type(), dataVal.Type(), data)
	}

	return nil
}

func (d *Decoder) decodeUint(name string, data interface{}, val reflect.Value) error {
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataKind := getKind(dataVal)
	dataType := dataVal.Type()

	switch {
	case dataKind == reflect.Int:
		i := dataVal.Int()
		if i < 0 && !d.config.WeaklyTypedInput {
			return fmt.Errorf("cannot parse '%s', %d overflows uint",
				name, i)
		}
		val.SetUint(uint64(i))
	case dataKind == reflect.Uint:
		val.SetUint(dataVal.Uint())
	case dataKind == reflect.Float32:
		f := dataVal.Float()
		if f < 0 && !d.config.WeaklyTypedInput {
			return fmt.Errorf("cannot parse '%s', %f overflows uint",
				name, f)
		}
		val.SetUint(uint64(f))
	case dataKind == reflect.Bool && d.config.WeaklyTypedInput:
		if dataVal.Bool() {
			val.SetUint(1)
		} else {
			val.SetUint(0)
		}
	case dataKind == reflect.String && d.config.WeaklyTypedInput:
		str := dataVal.String()
		if str == "" {
			str = "0"
		}

		i, err := strconv.ParseUint(str, 0, val.Type().Bits())
		if err == nil {
			val.SetUint(i)
		} else {
			return fmt.Errorf("cannot parse '%s' as uint: %s", name, err)
		}
	case dataType.PkgPath() == "encoding/json" && dataType.Name() == "Number":
		jn := data.(json.Number)
		i, err := strconv.ParseUint(string(jn), 0, 64)
		if err != nil {
			return fmt.Errorf(
				"error decoding json.Number into %s: %s", name, err)
		}
		val.SetUint(i)
	default:
		return fmt.Errorf(
			"'%s' expected type '%s', got unconvertible type '%s', value: '%v'",
			name, val.Type(), dataVal.Type(), data)
	}

	return nil
}

func (d *Decoder) decodeBool(name string, data interface{}, val reflect.Value) error {
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataKind := getKind(dataVal)

	switch {
	case dataKind == reflect.Bool:
		val.SetBool(dataVal.Bool())
	case dataKind == reflect.Int && d.config.WeaklyTypedInput:
		val.SetBool(dataVal.Int() != 0)
	case dataKind == reflect.Uint && d.config.WeaklyTypedInput:
		val.SetBool(dataVal.Uint() != 0)
	case dataKind == reflect.Float32 && d.config.WeaklyTypedInput:
		val.SetBool(dataVal.Float() != 0)
	case dataKind == reflect.String && d.config.WeaklyTypedInput:
		b, err := strconv.ParseBool(dataVal.String())
		if err == nil {
			val.SetBool(b)
		} else if dataVal.String() == "" {
			val.SetBool(false)
		} else {
			return fmt.Errorf("cannot parse '%s' as bool: %s", name, err)
		}
	default:
		return fmt.Errorf(
			"'%s' expected type '%s', got unconvertible type '%#v', value: '%#v'",
			name, val, dataVal, data)
	}

	return nil
}

func (d *Decoder) decodeFloat(name string, data interface{}, val reflect.Value) error {
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataKind := getKind(dataVal)
	dataType := dataVal.Type()

	switch {
	case dataKind == reflect.Int:
		val.SetFloat(float64(dataVal.Int()))
	case dataKind == reflect.Uint:
		val.SetFloat(float64(dataVal.Uint()))
	case dataKind == reflect.Float32:
		val.SetFloat(dataVal.Float())
	case dataKind == reflect.Bool && d.config.WeaklyTypedInput:
		if dataVal.Bool() {
			val.SetFloat(1)
		} else {
			val.SetFloat(0)
		}
	case dataKind == reflect.String && d.config.WeaklyTypedInput:
		str := dataVal.String()
		if str == "" {
			str = "0"
		}

		f, err := strconv.ParseFloat(str, val.Type().Bits())
		if err == nil {
			val.SetFloat(f)
		} else {
			return fmt.Errorf("cannot parse '%s' as float: %s", name, err)
		}
	case dataType.PkgPath() == "encoding/json" && dataType.Name() == "Number":
		jn := data.(json.Number)
		i, err := jn.Float64()
		if err != nil {
			return fmt.Errorf(
				"error decoding json.Number into %s: %s", name, err)
		}
		val.SetFloat(i)
	default:
		return fmt.Errorf(
			"'%s' expected type '%s', got unconvertible type '%s', value: '%v'",
			name, val.Type(), dataVal.Type(), data)
	}

	return nil
}

func (d *Decoder) decodeComplex(name string, data interface{}, val reflect.Value) error {
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataKind := getKind(dataVal)

	switch {
	case dataKind == reflect.Complex64:
		val.SetComplex(dataVal.Complex())
	default:
		return fmt.Errorf(
			"'%s' expected type '%s', got unconvertible type '%s', value: '%v'",
			name, val.Type(), dataVal.Type(), data)
	}

	return nil
}

func (d *Decoder) decodeMap(name string, data interface{}, val reflect.Value) error {
	valType := val.Type()
	valKeyType := valType.Key()
	valElemType := valType.Elem()

	
	valMap := val

	
	if valMap.IsNil() || d.config.ZeroFields {
		
		mapType := reflect.MapOf(valKeyType, valElemType)
		valMap = reflect.MakeMap(mapType)
	}

	dataVal := reflect.ValueOf(data)

	
	for dataVal.Kind() == reflect.Pointer {
		dataVal = reflect.Indirect(dataVal)
	}

	
	switch dataVal.Kind() {
	case reflect.Map:
		return d.decodeMapFromMap(name, dataVal, val, valMap)

	case reflect.Struct:
		return d.decodeMapFromStruct(name, dataVal, val, valMap)

	case reflect.Array, reflect.Slice:
		if d.config.WeaklyTypedInput {
			return d.decodeMapFromSlice(name, dataVal, val, valMap)
		}

		fallthrough

	default:
		return fmt.Errorf("'%s' expected a map, got '%s'", name, dataVal.Kind())
	}
}

func (d *Decoder) decodeMapFromSlice(name string, dataVal reflect.Value, val reflect.Value, valMap reflect.Value) error {
	
	if dataVal.Len() == 0 {
		val.Set(valMap)
		return nil
	}

	for i := 0; i < dataVal.Len(); i++ {
		err := d.decode(
			name+"["+strconv.Itoa(i)+"]",
			dataVal.Index(i).Interface(), val)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Decoder) decodeMapFromMap(name string, dataVal reflect.Value, val reflect.Value, valMap reflect.Value) error {
	valType := val.Type()
	valKeyType := valType.Key()
	valElemType := valType.Elem()

	
	var errs []error

	
	if dataVal.Len() == 0 {
		if dataVal.IsNil() {
			if !val.IsNil() {
				val.Set(dataVal)
			}
		} else {
			
			val.Set(valMap)
		}

		return nil
	}

	for _, k := range dataVal.MapKeys() {
		fieldName := name + "[" + k.String() + "]"

		
		currentKey := reflect.Indirect(reflect.New(valKeyType))
		if err := d.decode(fieldName, k.Interface(), currentKey); err != nil {
			errs = append(errs, err)
			continue
		}

		
		v := dataVal.MapIndex(k).Interface()
		currentVal := reflect.Indirect(reflect.New(valElemType))
		if err := d.decode(fieldName, v, currentVal); err != nil {
			errs = append(errs, err)
			continue
		}

		valMap.SetMapIndex(currentKey, currentVal)
	}

	
	val.Set(valMap)

	return errors.Join(errs...)
}

func (d *Decoder) decodeMapFromStruct(name string, dataVal reflect.Value, val reflect.Value, valMap reflect.Value) error {
	typ := dataVal.Type()
	for i := 0; i < typ.NumField(); i++ {
		
		
		f := typ.Field(i)
		if f.PkgPath != "" {
			continue
		}

		
		
		v := dataVal.Field(i)
		if !v.Type().AssignableTo(valMap.Type().Elem()) {
			return fmt.Errorf("cannot assign type '%s' to map value field of type '%s'", v.Type(), valMap.Type().Elem())
		}

		tagValue := f.Tag.Get(d.config.TagName)
		keyName := f.Name

		if tagValue == "" && d.config.IgnoreUntaggedFields {
			continue
		}

		
		squash := d.config.Squash && v.Kind() == reflect.Struct && f.Anonymous

		v = dereferencePtrToStructIfNeeded(v, d.config.TagName)

		
		if index := strings.Index(tagValue, ","); index != -1 {
			if tagValue[:index] == "-" {
				continue
			}
			
			if strings.Index(tagValue[index+1:], "omitempty") != -1 && isEmptyValue(v) {
				continue
			}

			
			squash = squash || strings.Contains(tagValue[index+1:], d.config.SquashTagOption)
			if squash {
				
				if v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Struct {
					v = v.Elem()
				}

				
				if v.Kind() != reflect.Struct {
					return fmt.Errorf("cannot squash non-struct type '%s'", v.Type())
				}
			} else {
				if strings.Index(tagValue[index+1:], "remain") != -1 {
					if v.Kind() != reflect.Map {
						return fmt.Errorf("error remain-tag field with invalid type: '%s'", v.Type())
					}

					ptr := v.MapRange()
					for ptr.Next() {
						valMap.SetMapIndex(ptr.Key(), ptr.Value())
					}
					continue
				}
			}
			if keyNameTagValue := tagValue[:index]; keyNameTagValue != "" {
				keyName = keyNameTagValue
			}
		} else if len(tagValue) > 0 {
			if tagValue == "-" {
				continue
			}
			keyName = tagValue
		}

		switch v.Kind() {
		
		case reflect.Struct:
			x := reflect.New(v.Type())
			x.Elem().Set(v)

			vType := valMap.Type()
			vKeyType := vType.Key()
			vElemType := vType.Elem()
			mType := reflect.MapOf(vKeyType, vElemType)
			vMap := reflect.MakeMap(mType)

			
			
			
			
			addrVal := reflect.New(vMap.Type())
			reflect.Indirect(addrVal).Set(vMap)

			err := d.decode(keyName, x.Interface(), reflect.Indirect(addrVal))
			if err != nil {
				return err
			}

			
			
			vMap = reflect.Indirect(addrVal)

			if squash {
				for _, k := range vMap.MapKeys() {
					valMap.SetMapIndex(k, vMap.MapIndex(k))
				}
			} else {
				valMap.SetMapIndex(reflect.ValueOf(keyName), vMap)
			}

		default:
			valMap.SetMapIndex(reflect.ValueOf(keyName), v)
		}
	}

	if val.CanAddr() {
		val.Set(valMap)
	}

	return nil
}

func (d *Decoder) decodePtr(name string, data interface{}, val reflect.Value) (bool, error) {
	
	
	isNil := data == nil
	if !isNil {
		switch v := reflect.Indirect(reflect.ValueOf(data)); v.Kind() {
		case reflect.Chan,
			reflect.Func,
			reflect.Interface,
			reflect.Map,
			reflect.Ptr,
			reflect.Slice:
			isNil = v.IsNil()
		}
	}
	if isNil {
		if !val.IsNil() && val.CanSet() {
			nilValue := reflect.New(val.Type()).Elem()
			val.Set(nilValue)
		}

		return true, nil
	}

	
	
	valType := val.Type()
	valElemType := valType.Elem()
	if val.CanSet() {
		realVal := val
		if realVal.IsNil() || d.config.ZeroFields {
			realVal = reflect.New(valElemType)
		}

		if err := d.decode(name, data, reflect.Indirect(realVal)); err != nil {
			return false, err
		}

		val.Set(realVal)
	} else {
		if err := d.decode(name, data, reflect.Indirect(val)); err != nil {
			return false, err
		}
	}
	return false, nil
}

func (d *Decoder) decodeFunc(name string, data interface{}, val reflect.Value) error {
	
	
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	if val.Type() != dataVal.Type() {
		return fmt.Errorf(
			"'%s' expected type '%s', got unconvertible type '%s', value: '%v'",
			name, val.Type(), dataVal.Type(), data)
	}
	val.Set(dataVal)
	return nil
}

func (d *Decoder) decodeSlice(name string, data interface{}, val reflect.Value) error {
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataValKind := dataVal.Kind()
	valType := val.Type()
	valElemType := valType.Elem()
	sliceType := reflect.SliceOf(valElemType)

	
	if dataValKind != reflect.Array && dataValKind != reflect.Slice {
		if d.config.WeaklyTypedInput {
			switch {
			
			case dataValKind == reflect.Slice, dataValKind == reflect.Array:
				break

			
			case dataValKind == reflect.Map:
				if dataVal.Len() == 0 {
					val.Set(reflect.MakeSlice(sliceType, 0, 0))
					return nil
				}
				
				return d.decodeSlice(name, []interface{}{data}, val)

			case dataValKind == reflect.String && valElemType.Kind() == reflect.Uint8:
				return d.decodeSlice(name, []byte(dataVal.String()), val)

			
			
			default:
				
				return d.decodeSlice(name, []interface{}{data}, val)
			}
		}

		return fmt.Errorf(
			"'%s': source data must be an array or slice, got %s", name, dataValKind)
	}

	
	if dataValKind != reflect.Array && dataVal.IsNil() {
		return nil
	}

	valSlice := val
	if valSlice.IsNil() || d.config.ZeroFields {
		
		valSlice = reflect.MakeSlice(sliceType, dataVal.Len(), dataVal.Len())
	} else if valSlice.Len() > dataVal.Len() {
		valSlice = valSlice.Slice(0, dataVal.Len())
	}

	
	var errs []error

	for i := 0; i < dataVal.Len(); i++ {
		currentData := dataVal.Index(i).Interface()
		for valSlice.Len() <= i {
			valSlice = reflect.Append(valSlice, reflect.Zero(valElemType))
		}
		currentField := valSlice.Index(i)

		fieldName := name + "[" + strconv.Itoa(i) + "]"
		if err := d.decode(fieldName, currentData, currentField); err != nil {
			errs = append(errs, err)
		}
	}

	
	val.Set(valSlice)

	return errors.Join(errs...)
}

func (d *Decoder) decodeArray(name string, data interface{}, val reflect.Value) error {
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataValKind := dataVal.Kind()
	valType := val.Type()
	valElemType := valType.Elem()
	arrayType := reflect.ArrayOf(valType.Len(), valElemType)

	valArray := val

	if isComparable(valArray) && valArray.Interface() == reflect.Zero(valArray.Type()).Interface() || d.config.ZeroFields {
		
		if dataValKind != reflect.Array && dataValKind != reflect.Slice {
			if d.config.WeaklyTypedInput {
				switch {
				
				case dataValKind == reflect.Map:
					if dataVal.Len() == 0 {
						val.Set(reflect.Zero(arrayType))
						return nil
					}

				
				
				default:
					
					return d.decodeArray(name, []interface{}{data}, val)
				}
			}

			return fmt.Errorf(
				"'%s': source data must be an array or slice, got %s", name, dataValKind)

		}
		if dataVal.Len() > arrayType.Len() {
			return fmt.Errorf(
				"'%s': expected source data to have length less or equal to %d, got %d", name, arrayType.Len(), dataVal.Len())
		}

		
		valArray = reflect.New(arrayType).Elem()
	}

	
	var errs []error

	for i := 0; i < dataVal.Len(); i++ {
		currentData := dataVal.Index(i).Interface()
		currentField := valArray.Index(i)

		fieldName := name + "[" + strconv.Itoa(i) + "]"
		if err := d.decode(fieldName, currentData, currentField); err != nil {
			errs = append(errs, err)
		}
	}

	
	val.Set(valArray)

	return errors.Join(errs...)
}

func (d *Decoder) decodeStruct(name string, data interface{}, val reflect.Value) error {
	dataVal := reflect.Indirect(reflect.ValueOf(data))

	
	
	if dataVal.Type() == val.Type() {
		val.Set(dataVal)
		return nil
	}

	dataValKind := dataVal.Kind()
	switch dataValKind {
	case reflect.Map:
		return d.decodeStructFromMap(name, dataVal, val)

	case reflect.Struct:
		
		
		

		
		mapType := reflect.TypeOf((map[string]interface{})(nil))
		mval := reflect.MakeMap(mapType)

		
		
		
		
		addrVal := reflect.New(mval.Type())

		reflect.Indirect(addrVal).Set(mval)
		if err := d.decodeMapFromStruct(name, dataVal, reflect.Indirect(addrVal), mval); err != nil {
			return err
		}

		result := d.decodeStructFromMap(name, reflect.Indirect(addrVal), val)
		return result

	default:
		return fmt.Errorf("'%s' expected a map, got '%s'", name, dataVal.Kind())
	}
}

func (d *Decoder) decodeStructFromMap(name string, dataVal, val reflect.Value) error {
	dataValType := dataVal.Type()
	if kind := dataValType.Key().Kind(); kind != reflect.String && kind != reflect.Interface {
		return fmt.Errorf(
			"'%s' needs a map with string keys, has '%s' keys",
			name, dataValType.Key().Kind())
	}

	dataValKeys := make(map[reflect.Value]struct{})
	dataValKeysUnused := make(map[interface{}]struct{})
	for _, dataValKey := range dataVal.MapKeys() {
		dataValKeys[dataValKey] = struct{}{}
		dataValKeysUnused[dataValKey.Interface()] = struct{}{}
	}

	targetValKeysUnused := make(map[interface{}]struct{})

	var errs []error

	
	
	
	structs := make([]reflect.Value, 1, 5)
	structs[0] = val

	
	
	type field struct {
		field reflect.StructField
		val   reflect.Value
	}

	
	
	var remainField *field

	fields := []field{}
	for len(structs) > 0 {
		structVal := structs[0]
		structs = structs[1:]

		structType := structVal.Type()

		for i := 0; i < structType.NumField(); i++ {
			fieldType := structType.Field(i)
			fieldVal := structVal.Field(i)
			if fieldVal.Kind() == reflect.Ptr && fieldVal.Elem().Kind() == reflect.Struct {
				
				fieldVal = fieldVal.Elem()
			}

			
			squash := d.config.Squash && fieldVal.Kind() == reflect.Struct && fieldType.Anonymous
			remain := false

			
			tagParts := strings.Split(fieldType.Tag.Get(d.config.TagName), ",")
			for _, tag := range tagParts[1:] {
				if tag == d.config.SquashTagOption {
					squash = true
					break
				}

				if tag == "remain" {
					remain = true
					break
				}
			}

			if squash {
				switch fieldVal.Kind() {
				case reflect.Struct:
					structs = append(structs, fieldVal)
				case reflect.Interface:
					if !fieldVal.IsNil() {
						structs = append(structs, fieldVal.Elem().Elem())
					}
				default:
					errs = append(errs, fmt.Errorf("%s: unsupported type for squash: %s", fieldType.Name, fieldVal.Kind()))
				}
				continue
			}

			
			if remain {
				remainField = &field{fieldType, fieldVal}
			} else {
				
				fields = append(fields, field{fieldType, fieldVal})
			}
		}
	}

	
	for _, f := range fields {
		field, fieldValue := f.field, f.val
		fieldName := field.Name

		tagValue := field.Tag.Get(d.config.TagName)
		if tagValue == "" && d.config.IgnoreUntaggedFields {
			continue
		}
		tagValue = strings.SplitN(tagValue, ",", 2)[0]
		if tagValue != "" {
			fieldName = tagValue
		}

		rawMapKey := reflect.ValueOf(fieldName)
		rawMapVal := dataVal.MapIndex(rawMapKey)
		if !rawMapVal.IsValid() {
			
			
			for dataValKey := range dataValKeys {
				mK, ok := dataValKey.Interface().(string)
				if !ok {
					
					continue
				}

				if d.config.MatchName(mK, fieldName) {
					rawMapKey = dataValKey
					rawMapVal = dataVal.MapIndex(dataValKey)
					break
				}
			}

			if !rawMapVal.IsValid() {
				
				
				targetValKeysUnused[fieldName] = struct{}{}
				continue
			}
		}

		if !fieldValue.IsValid() {
			
			panic("field is not valid")
		}

		
		
		if !fieldValue.CanSet() {
			continue
		}

		
		delete(dataValKeysUnused, rawMapKey.Interface())

		
		
		if name != "" {
			fieldName = name + "." + fieldName
		}

		if err := d.decode(fieldName, rawMapVal.Interface(), fieldValue); err != nil {
			errs = append(errs, err)
		}
	}

	
	
	if remainField != nil && len(dataValKeysUnused) > 0 {
		
		remain := map[interface{}]interface{}{}
		for key := range dataValKeysUnused {
			remain[key] = dataVal.MapIndex(reflect.ValueOf(key)).Interface()
		}

		
		if err := d.decodeMap(name, remain, remainField.val); err != nil {
			errs = append(errs, err)
		}

		
		
		dataValKeysUnused = nil
	}

	if d.config.ErrorUnused && len(dataValKeysUnused) > 0 {
		keys := make([]string, 0, len(dataValKeysUnused))
		for rawKey := range dataValKeysUnused {
			keys = append(keys, rawKey.(string))
		}
		sort.Strings(keys)

		err := fmt.Errorf("'%s' has invalid keys: %s", name, strings.Join(keys, ", "))
		errs = append(errs, err)
	}

	if d.config.ErrorUnset && len(targetValKeysUnused) > 0 {
		keys := make([]string, 0, len(targetValKeysUnused))
		for rawKey := range targetValKeysUnused {
			keys = append(keys, rawKey.(string))
		}
		sort.Strings(keys)

		err := fmt.Errorf("'%s' has unset fields: %s", name, strings.Join(keys, ", "))
		errs = append(errs, err)
	}

	if err := errors.Join(errs...); err != nil {
		return err
	}

	
	if d.config.Metadata != nil {
		for rawKey := range dataValKeysUnused {
			key := rawKey.(string)
			if name != "" {
				key = name + "." + key
			}

			d.config.Metadata.Unused = append(d.config.Metadata.Unused, key)
		}
		for rawKey := range targetValKeysUnused {
			key := rawKey.(string)
			if name != "" {
				key = name + "." + key
			}

			d.config.Metadata.Unset = append(d.config.Metadata.Unset, key)
		}
	}

	return nil
}

func isEmptyValue(v reflect.Value) bool {
	switch getKind(v) {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

func getKind(val reflect.Value) reflect.Kind {
	kind := val.Kind()

	switch {
	case kind >= reflect.Int && kind <= reflect.Int64:
		return reflect.Int
	case kind >= reflect.Uint && kind <= reflect.Uint64:
		return reflect.Uint
	case kind >= reflect.Float32 && kind <= reflect.Float64:
		return reflect.Float32
	case kind >= reflect.Complex64 && kind <= reflect.Complex128:
		return reflect.Complex64
	default:
		return kind
	}
}

func isStructTypeConvertibleToMap(typ reflect.Type, checkMapstructureTags bool, tagName string) bool {
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if f.PkgPath == "" && !checkMapstructureTags { 
			return true
		}
		if checkMapstructureTags && f.Tag.Get(tagName) != "" { 
			return true
		}
	}
	return false
}

func dereferencePtrToStructIfNeeded(v reflect.Value, tagName string) reflect.Value {
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return v
	}
	deref := v.Elem()
	derefT := deref.Type()
	if isStructTypeConvertibleToMap(derefT, true, tagName) {
		return deref
	}
	return v
}
