


package attribute 

import (
	"fmt"
)


type KeyValue struct {
	Key   Key
	Value Value
}


func (kv KeyValue) Valid() bool {
	return kv.Key.Defined() && kv.Value.Type() != INVALID
}


func Bool(k string, v bool) KeyValue {
	return Key(k).Bool(v)
}


func BoolSlice(k string, v []bool) KeyValue {
	return Key(k).BoolSlice(v)
}


func Int(k string, v int) KeyValue {
	return Key(k).Int(v)
}


func IntSlice(k string, v []int) KeyValue {
	return Key(k).IntSlice(v)
}


func Int64(k string, v int64) KeyValue {
	return Key(k).Int64(v)
}


func Int64Slice(k string, v []int64) KeyValue {
	return Key(k).Int64Slice(v)
}


func Float64(k string, v float64) KeyValue {
	return Key(k).Float64(v)
}


func Float64Slice(k string, v []float64) KeyValue {
	return Key(k).Float64Slice(v)
}


func String(k, v string) KeyValue {
	return Key(k).String(v)
}


func StringSlice(k string, v []string) KeyValue {
	return Key(k).StringSlice(v)
}



func Stringer(k string, v fmt.Stringer) KeyValue {
	return Key(k).String(v.String())
}
