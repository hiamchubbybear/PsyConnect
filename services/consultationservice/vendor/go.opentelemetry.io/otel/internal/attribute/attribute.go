



package attribute 

import (
	"reflect"
)


func BoolSliceValue(v []bool) interface{} {
	var zero bool
	cp := reflect.New(reflect.ArrayOf(len(v), reflect.TypeOf(zero))).Elem()
	reflect.Copy(cp, reflect.ValueOf(v))
	return cp.Interface()
}


func Int64SliceValue(v []int64) interface{} {
	var zero int64
	cp := reflect.New(reflect.ArrayOf(len(v), reflect.TypeOf(zero))).Elem()
	reflect.Copy(cp, reflect.ValueOf(v))
	return cp.Interface()
}


func Float64SliceValue(v []float64) interface{} {
	var zero float64
	cp := reflect.New(reflect.ArrayOf(len(v), reflect.TypeOf(zero))).Elem()
	reflect.Copy(cp, reflect.ValueOf(v))
	return cp.Interface()
}


func StringSliceValue(v []string) interface{} {
	var zero string
	cp := reflect.New(reflect.ArrayOf(len(v), reflect.TypeOf(zero))).Elem()
	reflect.Copy(cp, reflect.ValueOf(v))
	return cp.Interface()
}


func AsBoolSlice(v interface{}) []bool {
	rv := reflect.ValueOf(v)
	if rv.Type().Kind() != reflect.Array {
		return nil
	}
	cpy := make([]bool, rv.Len())
	if len(cpy) > 0 {
		_ = reflect.Copy(reflect.ValueOf(cpy), rv)
	}
	return cpy
}


func AsInt64Slice(v interface{}) []int64 {
	rv := reflect.ValueOf(v)
	if rv.Type().Kind() != reflect.Array {
		return nil
	}
	cpy := make([]int64, rv.Len())
	if len(cpy) > 0 {
		_ = reflect.Copy(reflect.ValueOf(cpy), rv)
	}
	return cpy
}


func AsFloat64Slice(v interface{}) []float64 {
	rv := reflect.ValueOf(v)
	if rv.Type().Kind() != reflect.Array {
		return nil
	}
	cpy := make([]float64, rv.Len())
	if len(cpy) > 0 {
		_ = reflect.Copy(reflect.ValueOf(cpy), rv)
	}
	return cpy
}


func AsStringSlice(v interface{}) []string {
	rv := reflect.ValueOf(v)
	if rv.Type().Kind() != reflect.Array {
		return nil
	}
	cpy := make([]string, rv.Len())
	if len(cpy) > 0 {
		_ = reflect.Copy(reflect.ValueOf(cpy), rv)
	}
	return cpy
}
