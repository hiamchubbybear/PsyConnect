


package attribute 

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"go.opentelemetry.io/otel/internal"
	"go.opentelemetry.io/otel/internal/attribute"
)

//go:generate stringer -type=Type


type Type int 


type Value struct {
	vtype    Type
	numeric  uint64
	stringly string
	slice    interface{}
}

const (
	
	INVALID Type = iota
	
	BOOL
	
	INT64
	
	FLOAT64
	
	STRING
	
	BOOLSLICE
	
	INT64SLICE
	
	FLOAT64SLICE
	
	STRINGSLICE
)


func BoolValue(v bool) Value {
	return Value{
		vtype:   BOOL,
		numeric: internal.BoolToRaw(v),
	}
}


func BoolSliceValue(v []bool) Value {
	return Value{vtype: BOOLSLICE, slice: attribute.BoolSliceValue(v)}
}


func IntValue(v int) Value {
	return Int64Value(int64(v))
}


func IntSliceValue(v []int) Value {
	var int64Val int64
	cp := reflect.New(reflect.ArrayOf(len(v), reflect.TypeOf(int64Val)))
	for i, val := range v {
		cp.Elem().Index(i).SetInt(int64(val))
	}
	return Value{
		vtype: INT64SLICE,
		slice: cp.Elem().Interface(),
	}
}


func Int64Value(v int64) Value {
	return Value{
		vtype:   INT64,
		numeric: internal.Int64ToRaw(v),
	}
}


func Int64SliceValue(v []int64) Value {
	return Value{vtype: INT64SLICE, slice: attribute.Int64SliceValue(v)}
}


func Float64Value(v float64) Value {
	return Value{
		vtype:   FLOAT64,
		numeric: internal.Float64ToRaw(v),
	}
}


func Float64SliceValue(v []float64) Value {
	return Value{vtype: FLOAT64SLICE, slice: attribute.Float64SliceValue(v)}
}


func StringValue(v string) Value {
	return Value{
		vtype:    STRING,
		stringly: v,
	}
}


func StringSliceValue(v []string) Value {
	return Value{vtype: STRINGSLICE, slice: attribute.StringSliceValue(v)}
}


func (v Value) Type() Type {
	return v.vtype
}



func (v Value) AsBool() bool {
	return internal.RawToBool(v.numeric)
}



func (v Value) AsBoolSlice() []bool {
	if v.vtype != BOOLSLICE {
		return nil
	}
	return v.asBoolSlice()
}

func (v Value) asBoolSlice() []bool {
	return attribute.AsBoolSlice(v.slice)
}



func (v Value) AsInt64() int64 {
	return internal.RawToInt64(v.numeric)
}



func (v Value) AsInt64Slice() []int64 {
	if v.vtype != INT64SLICE {
		return nil
	}
	return v.asInt64Slice()
}

func (v Value) asInt64Slice() []int64 {
	return attribute.AsInt64Slice(v.slice)
}



func (v Value) AsFloat64() float64 {
	return internal.RawToFloat64(v.numeric)
}



func (v Value) AsFloat64Slice() []float64 {
	if v.vtype != FLOAT64SLICE {
		return nil
	}
	return v.asFloat64Slice()
}

func (v Value) asFloat64Slice() []float64 {
	return attribute.AsFloat64Slice(v.slice)
}



func (v Value) AsString() string {
	return v.stringly
}



func (v Value) AsStringSlice() []string {
	if v.vtype != STRINGSLICE {
		return nil
	}
	return v.asStringSlice()
}

func (v Value) asStringSlice() []string {
	return attribute.AsStringSlice(v.slice)
}

type unknownValueType struct{}


func (v Value) AsInterface() interface{} {
	switch v.Type() {
	case BOOL:
		return v.AsBool()
	case BOOLSLICE:
		return v.asBoolSlice()
	case INT64:
		return v.AsInt64()
	case INT64SLICE:
		return v.asInt64Slice()
	case FLOAT64:
		return v.AsFloat64()
	case FLOAT64SLICE:
		return v.asFloat64Slice()
	case STRING:
		return v.stringly
	case STRINGSLICE:
		return v.asStringSlice()
	}
	return unknownValueType{}
}


func (v Value) Emit() string {
	switch v.Type() {
	case BOOLSLICE:
		return fmt.Sprint(v.asBoolSlice())
	case BOOL:
		return strconv.FormatBool(v.AsBool())
	case INT64SLICE:
		j, err := json.Marshal(v.asInt64Slice())
		if err != nil {
			return fmt.Sprintf("invalid: %v", v.asInt64Slice())
		}
		return string(j)
	case INT64:
		return strconv.FormatInt(v.AsInt64(), 10)
	case FLOAT64SLICE:
		j, err := json.Marshal(v.asFloat64Slice())
		if err != nil {
			return fmt.Sprintf("invalid: %v", v.asFloat64Slice())
		}
		return string(j)
	case FLOAT64:
		return fmt.Sprint(v.AsFloat64())
	case STRINGSLICE:
		j, err := json.Marshal(v.asStringSlice())
		if err != nil {
			return fmt.Sprintf("invalid: %v", v.asStringSlice())
		}
		return string(j)
	case STRING:
		return v.stringly
	default:
		return "unknown"
	}
}


func (v Value) MarshalJSON() ([]byte, error) {
	var jsonVal struct {
		Type  string
		Value interface{}
	}
	jsonVal.Type = v.Type().String()
	jsonVal.Value = v.AsInterface()
	return json.Marshal(jsonVal)
}
