package mapstructure

import (
	"encoding"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)



func typedDecodeHook(h DecodeHookFunc) DecodeHookFunc {
	
	var f1 DecodeHookFuncType
	var f2 DecodeHookFuncKind
	var f3 DecodeHookFuncValue

	
	
	potential := []interface{}{f1, f2, f3}

	v := reflect.ValueOf(h)
	vt := v.Type()
	for _, raw := range potential {
		pt := reflect.ValueOf(raw).Type()
		if vt.ConvertibleTo(pt) {
			return v.Convert(pt).Interface()
		}
	}

	return nil
}




func cachedDecodeHook(raw DecodeHookFunc) func(from reflect.Value, to reflect.Value) (interface{}, error) {
	switch f := typedDecodeHook(raw).(type) {
	case DecodeHookFuncType:
		return func(from reflect.Value, to reflect.Value) (interface{}, error) {
			return f(from.Type(), to.Type(), from.Interface())
		}
	case DecodeHookFuncKind:
		return func(from reflect.Value, to reflect.Value) (interface{}, error) {
			return f(from.Kind(), to.Kind(), from.Interface())
		}
	case DecodeHookFuncValue:
		return func(from reflect.Value, to reflect.Value) (interface{}, error) {
			return f(from, to)
		}
	default:
		return func(from reflect.Value, to reflect.Value) (interface{}, error) {
			return nil, errors.New("invalid decode hook signature")
		}
	}
}




func DecodeHookExec(
	raw DecodeHookFunc,
	from reflect.Value, to reflect.Value,
) (interface{}, error) {
	switch f := typedDecodeHook(raw).(type) {
	case DecodeHookFuncType:
		return f(from.Type(), to.Type(), from.Interface())
	case DecodeHookFuncKind:
		return f(from.Kind(), to.Kind(), from.Interface())
	case DecodeHookFuncValue:
		return f(from, to)
	default:
		return nil, errors.New("invalid decode hook signature")
	}
}






func ComposeDecodeHookFunc(fs ...DecodeHookFunc) DecodeHookFunc {
	cached := make([]func(from reflect.Value, to reflect.Value) (interface{}, error), 0, len(fs))
	for _, f := range fs {
		cached = append(cached, cachedDecodeHook(f))
	}
	return func(f reflect.Value, t reflect.Value) (interface{}, error) {
		var err error
		data := f.Interface()

		newFrom := f
		for _, c := range cached {
			data, err = c(newFrom, t)
			if err != nil {
				return nil, err
			}
			newFrom = reflect.ValueOf(data)
		}

		return data, nil
	}
}



func OrComposeDecodeHookFunc(ff ...DecodeHookFunc) DecodeHookFunc {
	cached := make([]func(from reflect.Value, to reflect.Value) (interface{}, error), 0, len(ff))
	for _, f := range ff {
		cached = append(cached, cachedDecodeHook(f))
	}
	return func(a, b reflect.Value) (interface{}, error) {
		var allErrs string
		var out interface{}
		var err error

		for _, c := range cached {
			out, err = c(a, b)
			if err != nil {
				allErrs += err.Error() + "\n"
				continue
			}

			return out, nil
		}

		return nil, errors.New(allErrs)
	}
}



func StringToSliceHookFunc(sep string) DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.SliceOf(f) {
			return data, nil
		}

		raw := data.(string)
		if raw == "" {
			return []string{}, nil
		}

		return strings.Split(raw, sep), nil
	}
}



func StringToTimeDurationHookFunc() DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(time.Duration(5)) {
			return data, nil
		}

		
		return time.ParseDuration(data.(string))
	}
}



func StringToURLHookFunc() DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(&url.URL{}) {
			return data, nil
		}

		
		return url.Parse(data.(string))
	}
}



func StringToIPHookFunc() DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(net.IP{}) {
			return data, nil
		}

		
		ip := net.ParseIP(data.(string))
		if ip == nil {
			return net.IP{}, fmt.Errorf("failed parsing ip %v", data)
		}

		return ip, nil
	}
}



func StringToIPNetHookFunc() DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(net.IPNet{}) {
			return data, nil
		}

		
		_, net, err := net.ParseCIDR(data.(string))
		return net, err
	}
}



func StringToTimeHookFunc(layout string) DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(time.Time{}) {
			return data, nil
		}

		
		return time.Parse(layout, data.(string))
	}
}






func WeaklyTypedHook(
	f reflect.Kind,
	t reflect.Kind,
	data interface{},
) (interface{}, error) {
	dataVal := reflect.ValueOf(data)
	switch t {
	case reflect.String:
		switch f {
		case reflect.Bool:
			if dataVal.Bool() {
				return "1", nil
			}
			return "0", nil
		case reflect.Float32:
			return strconv.FormatFloat(dataVal.Float(), 'f', -1, 64), nil
		case reflect.Int:
			return strconv.FormatInt(dataVal.Int(), 10), nil
		case reflect.Slice:
			dataType := dataVal.Type()
			elemKind := dataType.Elem().Kind()
			if elemKind == reflect.Uint8 {
				return string(dataVal.Interface().([]uint8)), nil
			}
		case reflect.Uint:
			return strconv.FormatUint(dataVal.Uint(), 10), nil
		}
	}

	return data, nil
}

func RecursiveStructToMapHookFunc() DecodeHookFunc {
	return func(f reflect.Value, t reflect.Value) (interface{}, error) {
		if f.Kind() != reflect.Struct {
			return f.Interface(), nil
		}

		var i interface{} = struct{}{}
		if t.Type() != reflect.TypeOf(&i).Elem() {
			return f.Interface(), nil
		}

		m := make(map[string]interface{})
		t.Set(reflect.ValueOf(m))

		return f.Interface(), nil
	}
}




func TextUnmarshallerHookFunc() DecodeHookFuncType {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		result := reflect.New(t).Interface()
		unmarshaller, ok := result.(encoding.TextUnmarshaler)
		if !ok {
			return data, nil
		}
		str, ok := data.(string)
		if !ok {
			str = reflect.Indirect(reflect.ValueOf(&data)).Elem().String()
		}
		if err := unmarshaller.UnmarshalText([]byte(str)); err != nil {
			return nil, err
		}
		return result, nil
	}
}



func StringToNetIPAddrHookFunc() DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(netip.Addr{}) {
			return data, nil
		}

		
		return netip.ParseAddr(data.(string))
	}
}



func StringToNetIPAddrPortHookFunc() DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(netip.AddrPort{}) {
			return data, nil
		}

		
		return netip.ParseAddrPort(data.(string))
	}
}




func StringToBasicTypeHookFunc() DecodeHookFunc {
	return ComposeDecodeHookFunc(
		StringToInt8HookFunc(),
		StringToUint8HookFunc(),
		StringToInt16HookFunc(),
		StringToUint16HookFunc(),
		StringToInt32HookFunc(),
		StringToUint32HookFunc(),
		StringToInt64HookFunc(),
		StringToUint64HookFunc(),
		StringToIntHookFunc(),
		StringToUintHookFunc(),
		StringToFloat32HookFunc(),
		StringToFloat64HookFunc(),
		StringToBoolHookFunc(),
		
		
		
		StringToComplex64HookFunc(),
		StringToComplex128HookFunc(),
	)
}



func StringToInt8HookFunc() DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || t.Kind() != reflect.Int8 {
			return data, nil
		}

		
		i64, err := strconv.ParseInt(data.(string), 0, 8)
		return int8(i64), err
	}
}



func StringToUint8HookFunc() DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || t.Kind() != reflect.Uint8 {
			return data, nil
		}

		
		u64, err := strconv.ParseUint(data.(string), 0, 8)
		return uint8(u64), err
	}
}



func StringToInt16HookFunc() DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || t.Kind() != reflect.Int16 {
			return data, nil
		}

		
		i64, err := strconv.ParseInt(data.(string), 0, 16)
		return int16(i64), err
	}
}



func StringToUint16HookFunc() DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || t.Kind() != reflect.Uint16 {
			return data, nil
		}

		
		u64, err := strconv.ParseUint(data.(string), 0, 16)
		return uint16(u64), err
	}
}



func StringToInt32HookFunc() DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || t.Kind() != reflect.Int32 {
			return data, nil
		}

		
		i64, err := strconv.ParseInt(data.(string), 0, 32)
		return int32(i64), err
	}
}



func StringToUint32HookFunc() DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || t.Kind() != reflect.Uint32 {
			return data, nil
		}

		
		u64, err := strconv.ParseUint(data.(string), 0, 32)
		return uint32(u64), err
	}
}



func StringToInt64HookFunc() DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || t.Kind() != reflect.Int64 {
			return data, nil
		}

		
		return strconv.ParseInt(data.(string), 0, 64)
	}
}



func StringToUint64HookFunc() DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || t.Kind() != reflect.Uint64 {
			return data, nil
		}

		
		return strconv.ParseUint(data.(string), 0, 64)
	}
}



func StringToIntHookFunc() DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || t.Kind() != reflect.Int {
			return data, nil
		}

		
		i64, err := strconv.ParseInt(data.(string), 0, 0)
		return int(i64), err
	}
}



func StringToUintHookFunc() DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || t.Kind() != reflect.Uint {
			return data, nil
		}

		
		u64, err := strconv.ParseUint(data.(string), 0, 0)
		return uint(u64), err
	}
}



func StringToFloat32HookFunc() DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || t.Kind() != reflect.Float32 {
			return data, nil
		}

		
		f64, err := strconv.ParseFloat(data.(string), 32)
		return float32(f64), err
	}
}



func StringToFloat64HookFunc() DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || t.Kind() != reflect.Float64 {
			return data, nil
		}

		
		return strconv.ParseFloat(data.(string), 64)
	}
}



func StringToBoolHookFunc() DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || t.Kind() != reflect.Bool {
			return data, nil
		}

		
		return strconv.ParseBool(data.(string))
	}
}



func StringToByteHookFunc() DecodeHookFunc {
	return StringToUint8HookFunc()
}



func StringToRuneHookFunc() DecodeHookFunc {
	return StringToInt32HookFunc()
}



func StringToComplex64HookFunc() DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || t.Kind() != reflect.Complex64 {
			return data, nil
		}

		
		c128, err := strconv.ParseComplex(data.(string), 64)
		return complex64(c128), err
	}
}



func StringToComplex128HookFunc() DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String || t.Kind() != reflect.Complex128 {
			return data, nil
		}

		
		return strconv.ParseComplex(data.(string), 128)
	}
}
