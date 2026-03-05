


//go:generate stringer -type=ValueKind -trimprefix=ValueKind

package telemetry

import (
	"bytes"
	"cmp"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"slices"
	"strconv"
	"unsafe"
)



type Value struct {
	
	noCmp [0]func() 

	
	
	num uint64
	
	
	
	
	any any
}

type (
	
	stringptr *byte
	
	bytesptr *byte
	
	sliceptr *Value
	
	mapptr *Attr
)


type ValueKind int


const (
	ValueKindEmpty ValueKind = iota
	ValueKindBool
	ValueKindFloat64
	ValueKindInt64
	ValueKindString
	ValueKindBytes
	ValueKindSlice
	ValueKindMap
)

var valueKindStrings = []string{
	"Empty",
	"Bool",
	"Float64",
	"Int64",
	"String",
	"Bytes",
	"Slice",
	"Map",
}

func (k ValueKind) String() string {
	if k >= 0 && int(k) < len(valueKindStrings) {
		return valueKindStrings[k]
	}
	return "<unknown telemetry.ValueKind>"
}


func StringValue(v string) Value {
	return Value{
		num: uint64(len(v)),
		any: stringptr(unsafe.StringData(v)),
	}
}


func IntValue(v int) Value { return Int64Value(int64(v)) }


func Int64Value(v int64) Value {
	return Value{num: uint64(v), any: ValueKindInt64}
}


func Float64Value(v float64) Value {
	return Value{num: math.Float64bits(v), any: ValueKindFloat64}
}


func BoolValue(v bool) Value { 
	var n uint64
	if v {
		n = 1
	}
	return Value{num: n, any: ValueKindBool}
}



func BytesValue(v []byte) Value {
	return Value{
		num: uint64(len(v)),
		any: bytesptr(unsafe.SliceData(v)),
	}
}



func SliceValue(vs ...Value) Value {
	return Value{
		num: uint64(len(vs)),
		any: sliceptr(unsafe.SliceData(vs)),
	}
}



func MapValue(kvs ...Attr) Value {
	return Value{
		num: uint64(len(kvs)),
		any: mapptr(unsafe.SliceData(kvs)),
	}
}


func (v Value) AsString() string {
	if sp, ok := v.any.(stringptr); ok {
		return unsafe.String(sp, v.num)
	}
	
	return ""
}



func (v Value) asString() string {
	return unsafe.String(v.any.(stringptr), v.num)
}


func (v Value) AsInt64() int64 {
	if v.Kind() != ValueKindInt64 {
		
		return 0
	}
	return v.asInt64()
}



func (v Value) asInt64() int64 {
	
	return int64(v.num) 
}


func (v Value) AsBool() bool {
	if v.Kind() != ValueKindBool {
		
		return false
	}
	return v.asBool()
}



func (v Value) asBool() bool { return v.num == 1 }


func (v Value) AsFloat64() float64 {
	if v.Kind() != ValueKindFloat64 {
		
		return 0
	}
	return v.asFloat64()
}



func (v Value) asFloat64() float64 { return math.Float64frombits(v.num) }


func (v Value) AsBytes() []byte {
	if sp, ok := v.any.(bytesptr); ok {
		return unsafe.Slice((*byte)(sp), v.num)
	}
	
	return nil
}



func (v Value) asBytes() []byte {
	return unsafe.Slice((*byte)(v.any.(bytesptr)), v.num)
}


func (v Value) AsSlice() []Value {
	if sp, ok := v.any.(sliceptr); ok {
		return unsafe.Slice((*Value)(sp), v.num)
	}
	
	return nil
}



func (v Value) asSlice() []Value {
	return unsafe.Slice((*Value)(v.any.(sliceptr)), v.num)
}


func (v Value) AsMap() []Attr {
	if sp, ok := v.any.(mapptr); ok {
		return unsafe.Slice((*Attr)(sp), v.num)
	}
	
	return nil
}



func (v Value) asMap() []Attr {
	return unsafe.Slice((*Attr)(v.any.(mapptr)), v.num)
}


func (v Value) Kind() ValueKind {
	switch x := v.any.(type) {
	case ValueKind:
		return x
	case stringptr:
		return ValueKindString
	case bytesptr:
		return ValueKindBytes
	case sliceptr:
		return ValueKindSlice
	case mapptr:
		return ValueKindMap
	default:
		return ValueKindEmpty
	}
}


func (v Value) Empty() bool { return v.Kind() == ValueKindEmpty }


func (v Value) Equal(w Value) bool {
	k1 := v.Kind()
	k2 := w.Kind()
	if k1 != k2 {
		return false
	}
	switch k1 {
	case ValueKindInt64, ValueKindBool:
		return v.num == w.num
	case ValueKindString:
		return v.asString() == w.asString()
	case ValueKindFloat64:
		return v.asFloat64() == w.asFloat64()
	case ValueKindSlice:
		return slices.EqualFunc(v.asSlice(), w.asSlice(), Value.Equal)
	case ValueKindMap:
		sv := sortMap(v.asMap())
		sw := sortMap(w.asMap())
		return slices.EqualFunc(sv, sw, Attr.Equal)
	case ValueKindBytes:
		return bytes.Equal(v.asBytes(), w.asBytes())
	case ValueKindEmpty:
		return true
	default:
		
		return false
	}
}

func sortMap(m []Attr) []Attr {
	sm := make([]Attr, len(m))
	copy(sm, m)
	slices.SortFunc(sm, func(a, b Attr) int {
		return cmp.Compare(a.Key, b.Key)
	})

	return sm
}





func (v Value) String() string {
	switch v.Kind() {
	case ValueKindString:
		return v.asString()
	case ValueKindInt64:
		
		return strconv.FormatInt(int64(v.num), 10) 
	case ValueKindFloat64:
		return strconv.FormatFloat(v.asFloat64(), 'g', -1, 64)
	case ValueKindBool:
		return strconv.FormatBool(v.asBool())
	case ValueKindBytes:
		return fmt.Sprint(v.asBytes())
	case ValueKindMap:
		return fmt.Sprint(v.asMap())
	case ValueKindSlice:
		return fmt.Sprint(v.asSlice())
	case ValueKindEmpty:
		return "<nil>"
	default:
		
		
		
		
		
		
		return fmt.Sprintf("<unhandled telemetry.ValueKind: %s>", v.Kind())
	}
}


func (v *Value) MarshalJSON() ([]byte, error) {
	switch v.Kind() {
	case ValueKindString:
		return json.Marshal(struct {
			Value string `json:"stringValue"`
		}{v.asString()})
	case ValueKindInt64:
		return json.Marshal(struct {
			Value string `json:"intValue"`
		}{strconv.FormatInt(int64(v.num), 10)})
	case ValueKindFloat64:
		return json.Marshal(struct {
			Value float64 `json:"doubleValue"`
		}{v.asFloat64()})
	case ValueKindBool:
		return json.Marshal(struct {
			Value bool `json:"boolValue"`
		}{v.asBool()})
	case ValueKindBytes:
		return json.Marshal(struct {
			Value []byte `json:"bytesValue"`
		}{v.asBytes()})
	case ValueKindMap:
		return json.Marshal(struct {
			Value struct {
				Values []Attr `json:"values"`
			} `json:"kvlistValue"`
		}{struct {
			Values []Attr `json:"values"`
		}{v.asMap()}})
	case ValueKindSlice:
		return json.Marshal(struct {
			Value struct {
				Values []Value `json:"values"`
			} `json:"arrayValue"`
		}{struct {
			Values []Value `json:"values"`
		}{v.asSlice()}})
	case ValueKindEmpty:
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown Value kind: %s", v.Kind().String())
	}
}


func (v *Value) UnmarshalJSON(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))

	t, err := decoder.Token()
	if err != nil {
		return err
	}
	if t != json.Delim('{') {
		return errors.New("invalid Value type")
	}

	for decoder.More() {
		keyIface, err := decoder.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				
				return nil
			}
			return err
		}

		key, ok := keyIface.(string)
		if !ok {
			return fmt.Errorf("invalid Value key: %#v", keyIface)
		}

		switch key {
		case "stringValue", "string_value":
			var val string
			err = decoder.Decode(&val)
			*v = StringValue(val)
		case "boolValue", "bool_value":
			var val bool
			err = decoder.Decode(&val)
			*v = BoolValue(val)
		case "intValue", "int_value":
			var val protoInt64
			err = decoder.Decode(&val)
			*v = Int64Value(val.Int64())
		case "doubleValue", "double_value":
			var val float64
			err = decoder.Decode(&val)
			*v = Float64Value(val)
		case "bytesValue", "bytes_value":
			var val64 string
			if err := decoder.Decode(&val64); err != nil {
				return err
			}
			var val []byte
			val, err = base64.StdEncoding.DecodeString(val64)
			*v = BytesValue(val)
		case "arrayValue", "array_value":
			var val struct{ Values []Value }
			err = decoder.Decode(&val)
			*v = SliceValue(val.Values...)
		case "kvlistValue", "kvlist_value":
			var val struct{ Values []Attr }
			err = decoder.Decode(&val)
			*v = MapValue(val.Values...)
		default:
			
			continue
		}
		
		return err
	}

	
	return nil
}
