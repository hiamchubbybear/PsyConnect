



package protoreflect

import (
	"fmt"
	"math"
)



























































type Value value

















func ValueOf(v any) Value {
	switch v := v.(type) {
	case nil:
		return Value{}
	case bool:
		return ValueOfBool(v)
	case int32:
		return ValueOfInt32(v)
	case int64:
		return ValueOfInt64(v)
	case uint32:
		return ValueOfUint32(v)
	case uint64:
		return ValueOfUint64(v)
	case float32:
		return ValueOfFloat32(v)
	case float64:
		return ValueOfFloat64(v)
	case string:
		return ValueOfString(v)
	case []byte:
		return ValueOfBytes(v)
	case EnumNumber:
		return ValueOfEnum(v)
	case Message, List, Map:
		return valueOfIface(v)
	case ProtoMessage:
		panic(fmt.Sprintf("invalid proto.Message(%T) type, expected a protoreflect.Message type", v))
	default:
		panic(fmt.Sprintf("invalid type: %T", v))
	}
}


func ValueOfBool(v bool) Value {
	if v {
		return Value{typ: boolType, num: 1}
	} else {
		return Value{typ: boolType, num: 0}
	}
}


func ValueOfInt32(v int32) Value {
	return Value{typ: int32Type, num: uint64(v)}
}


func ValueOfInt64(v int64) Value {
	return Value{typ: int64Type, num: uint64(v)}
}


func ValueOfUint32(v uint32) Value {
	return Value{typ: uint32Type, num: uint64(v)}
}


func ValueOfUint64(v uint64) Value {
	return Value{typ: uint64Type, num: v}
}


func ValueOfFloat32(v float32) Value {
	return Value{typ: float32Type, num: uint64(math.Float64bits(float64(v)))}
}


func ValueOfFloat64(v float64) Value {
	return Value{typ: float64Type, num: uint64(math.Float64bits(float64(v)))}
}


func ValueOfString(v string) Value {
	return valueOfString(v)
}


func ValueOfBytes(v []byte) Value {
	return valueOfBytes(v[:len(v):len(v)])
}


func ValueOfEnum(v EnumNumber) Value {
	return Value{typ: enumType, num: uint64(v)}
}


func ValueOfMessage(v Message) Value {
	return valueOfIface(v)
}


func ValueOfList(v List) Value {
	return valueOfIface(v)
}


func ValueOfMap(v Map) Value {
	return valueOfIface(v)
}


func (v Value) IsValid() bool {
	return v.typ != nilType
}




func (v Value) Interface() any {
	switch v.typ {
	case nilType:
		return nil
	case boolType:
		return v.Bool()
	case int32Type:
		return int32(v.Int())
	case int64Type:
		return int64(v.Int())
	case uint32Type:
		return uint32(v.Uint())
	case uint64Type:
		return uint64(v.Uint())
	case float32Type:
		return float32(v.Float())
	case float64Type:
		return float64(v.Float())
	case stringType:
		return v.String()
	case bytesType:
		return v.Bytes()
	case enumType:
		return v.Enum()
	default:
		return v.getIface()
	}
}

func (v Value) typeName() string {
	switch v.typ {
	case nilType:
		return "nil"
	case boolType:
		return "bool"
	case int32Type:
		return "int32"
	case int64Type:
		return "int64"
	case uint32Type:
		return "uint32"
	case uint64Type:
		return "uint64"
	case float32Type:
		return "float32"
	case float64Type:
		return "float64"
	case stringType:
		return "string"
	case bytesType:
		return "bytes"
	case enumType:
		return "enum"
	default:
		switch v := v.getIface().(type) {
		case Message:
			return "message"
		case List:
			return "list"
		case Map:
			return "map"
		default:
			return fmt.Sprintf("<unknown: %T>", v)
		}
	}
}

func (v Value) panicMessage(what string) string {
	return fmt.Sprintf("type mismatch: cannot convert %v to %s", v.typeName(), what)
}


func (v Value) Bool() bool {
	switch v.typ {
	case boolType:
		return v.num > 0
	default:
		panic(v.panicMessage("bool"))
	}
}


func (v Value) Int() int64 {
	switch v.typ {
	case int32Type, int64Type:
		return int64(v.num)
	default:
		panic(v.panicMessage("int"))
	}
}


func (v Value) Uint() uint64 {
	switch v.typ {
	case uint32Type, uint64Type:
		return uint64(v.num)
	default:
		panic(v.panicMessage("uint"))
	}
}


func (v Value) Float() float64 {
	switch v.typ {
	case float32Type, float64Type:
		return math.Float64frombits(uint64(v.num))
	default:
		panic(v.panicMessage("float"))
	}
}



func (v Value) String() string {
	switch v.typ {
	case stringType:
		return v.getString()
	default:
		return fmt.Sprint(v.Interface())
	}
}


func (v Value) Bytes() []byte {
	switch v.typ {
	case bytesType:
		return v.getBytes()
	default:
		panic(v.panicMessage("bytes"))
	}
}


func (v Value) Enum() EnumNumber {
	switch v.typ {
	case enumType:
		return EnumNumber(v.num)
	default:
		panic(v.panicMessage("enum"))
	}
}


func (v Value) Message() Message {
	switch vi := v.getIface().(type) {
	case Message:
		return vi
	default:
		panic(v.panicMessage("message"))
	}
}


func (v Value) List() List {
	switch vi := v.getIface().(type) {
	case List:
		return vi
	default:
		panic(v.panicMessage("list"))
	}
}


func (v Value) Map() Map {
	switch vi := v.getIface().(type) {
	case Map:
		return vi
	default:
		panic(v.panicMessage("map"))
	}
}


func (v Value) MapKey() MapKey {
	switch v.typ {
	case boolType, int32Type, int64Type, uint32Type, uint64Type, stringType:
		return MapKey(v)
	default:
		panic(v.panicMessage("map key"))
	}
}























type MapKey value


func (k MapKey) IsValid() bool {
	return Value(k).IsValid()
}


func (k MapKey) Interface() any {
	return Value(k).Interface()
}


func (k MapKey) Bool() bool {
	return Value(k).Bool()
}


func (k MapKey) Int() int64 {
	return Value(k).Int()
}


func (k MapKey) Uint() uint64 {
	return Value(k).Uint()
}



func (k MapKey) String() string {
	return Value(k).String()
}


func (k MapKey) Value() Value {
	return Value(k)
}
