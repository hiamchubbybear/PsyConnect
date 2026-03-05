package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)




func (v *validate) extractTypeInternal(current reflect.Value, nullable bool) (reflect.Value, reflect.Kind, bool) {

BEGIN:
	switch current.Kind() {
	case reflect.Ptr:

		nullable = true

		if current.IsNil() {
			return current, reflect.Ptr, nullable
		}

		current = current.Elem()
		goto BEGIN

	case reflect.Interface:

		nullable = true

		if current.IsNil() {
			return current, reflect.Interface, nullable
		}

		current = current.Elem()
		goto BEGIN

	case reflect.Invalid:
		return current, reflect.Invalid, nullable

	default:

		if v.v.hasCustomFuncs {

			if fn, ok := v.v.customFuncs[current.Type()]; ok {
				current = reflect.ValueOf(fn(current))
				goto BEGIN
			}
		}

		return current, current.Kind(), nullable
	}
}






func (v *validate) getStructFieldOKInternal(val reflect.Value, namespace string) (current reflect.Value, kind reflect.Kind, nullable bool, found bool) {

BEGIN:
	current, kind, nullable = v.ExtractType(val)
	if kind == reflect.Invalid {
		return
	}

	if namespace == "" {
		found = true
		return
	}

	switch kind {

	case reflect.Ptr, reflect.Interface:
		return

	case reflect.Struct:

		typ := current.Type()
		fld := namespace
		var ns string

		if !typ.ConvertibleTo(timeType) {

			idx := strings.Index(namespace, namespaceSeparator)

			if idx != -1 {
				fld = namespace[:idx]
				ns = namespace[idx+1:]
			} else {
				ns = ""
			}

			bracketIdx := strings.Index(fld, leftBracket)
			if bracketIdx != -1 {
				fld = fld[:bracketIdx]

				ns = namespace[bracketIdx:]
			}

			val = current.FieldByName(fld)
			namespace = ns
			goto BEGIN
		}

	case reflect.Array, reflect.Slice:
		idx := strings.Index(namespace, leftBracket)
		idx2 := strings.Index(namespace, rightBracket)

		arrIdx, _ := strconv.Atoi(namespace[idx+1 : idx2])

		if arrIdx >= current.Len() {
			return
		}

		startIdx := idx2 + 1

		if startIdx < len(namespace) {
			if namespace[startIdx:startIdx+1] == namespaceSeparator {
				startIdx++
			}
		}

		val = current.Index(arrIdx)
		namespace = namespace[startIdx:]
		goto BEGIN

	case reflect.Map:
		idx := strings.Index(namespace, leftBracket) + 1
		idx2 := strings.Index(namespace, rightBracket)

		endIdx := idx2

		if endIdx+1 < len(namespace) {
			if namespace[endIdx+1:endIdx+2] == namespaceSeparator {
				endIdx++
			}
		}

		key := namespace[idx:idx2]

		switch current.Type().Key().Kind() {
		case reflect.Int:
			i, _ := strconv.Atoi(key)
			val = current.MapIndex(reflect.ValueOf(i))
			namespace = namespace[endIdx+1:]

		case reflect.Int8:
			i, _ := strconv.ParseInt(key, 10, 8)
			val = current.MapIndex(reflect.ValueOf(int8(i)))
			namespace = namespace[endIdx+1:]

		case reflect.Int16:
			i, _ := strconv.ParseInt(key, 10, 16)
			val = current.MapIndex(reflect.ValueOf(int16(i)))
			namespace = namespace[endIdx+1:]

		case reflect.Int32:
			i, _ := strconv.ParseInt(key, 10, 32)
			val = current.MapIndex(reflect.ValueOf(int32(i)))
			namespace = namespace[endIdx+1:]

		case reflect.Int64:
			i, _ := strconv.ParseInt(key, 10, 64)
			val = current.MapIndex(reflect.ValueOf(i))
			namespace = namespace[endIdx+1:]

		case reflect.Uint:
			i, _ := strconv.ParseUint(key, 10, 0)
			val = current.MapIndex(reflect.ValueOf(uint(i)))
			namespace = namespace[endIdx+1:]

		case reflect.Uint8:
			i, _ := strconv.ParseUint(key, 10, 8)
			val = current.MapIndex(reflect.ValueOf(uint8(i)))
			namespace = namespace[endIdx+1:]

		case reflect.Uint16:
			i, _ := strconv.ParseUint(key, 10, 16)
			val = current.MapIndex(reflect.ValueOf(uint16(i)))
			namespace = namespace[endIdx+1:]

		case reflect.Uint32:
			i, _ := strconv.ParseUint(key, 10, 32)
			val = current.MapIndex(reflect.ValueOf(uint32(i)))
			namespace = namespace[endIdx+1:]

		case reflect.Uint64:
			i, _ := strconv.ParseUint(key, 10, 64)
			val = current.MapIndex(reflect.ValueOf(i))
			namespace = namespace[endIdx+1:]

		case reflect.Float32:
			f, _ := strconv.ParseFloat(key, 32)
			val = current.MapIndex(reflect.ValueOf(float32(f)))
			namespace = namespace[endIdx+1:]

		case reflect.Float64:
			f, _ := strconv.ParseFloat(key, 64)
			val = current.MapIndex(reflect.ValueOf(f))
			namespace = namespace[endIdx+1:]

		case reflect.Bool:
			b, _ := strconv.ParseBool(key)
			val = current.MapIndex(reflect.ValueOf(b))
			namespace = namespace[endIdx+1:]

		
		default:
			val = current.MapIndex(reflect.ValueOf(key))
			namespace = namespace[endIdx+1:]
		}

		goto BEGIN
	}

	
	panic("Invalid field namespace")
}



func asInt(param string) int64 {
	i, err := strconv.ParseInt(param, 0, 64)
	panicIf(err)

	return i
}



func asIntFromTimeDuration(param string) int64 {
	d, err := time.ParseDuration(param)
	if err != nil {
		
		return asInt(param)
	}
	return int64(d)
}



func asIntFromType(t reflect.Type, param string) int64 {
	switch t {
	case timeDurationType:
		return asIntFromTimeDuration(param)
	default:
		return asInt(param)
	}
}



func asUint(param string) uint64 {

	i, err := strconv.ParseUint(param, 0, 64)
	panicIf(err)

	return i
}



func asFloat64(param string) float64 {
	i, err := strconv.ParseFloat(param, 64)
	panicIf(err)
	return i
}



func asFloat32(param string) float64 {
	i, err := strconv.ParseFloat(param, 32)
	panicIf(err)
	return i
}



func asBool(param string) bool {

	i, err := strconv.ParseBool(param)
	panicIf(err)

	return i
}

func panicIf(err error) {
	if err != nil {
		panic(err.Error())
	}
}



func fieldMatchesRegexByStringerValOrString(regexFn func() *regexp.Regexp, fl FieldLevel) bool {
	regex := regexFn()
	switch fl.Field().Kind() {
	case reflect.String:
		return regex.MatchString(fl.Field().String())
	default:
		if stringer, ok := fl.Field().Interface().(fmt.Stringer); ok {
			return regex.MatchString(stringer.String())
		} else {
			return regex.MatchString(fl.Field().String())
		}
	}
}
