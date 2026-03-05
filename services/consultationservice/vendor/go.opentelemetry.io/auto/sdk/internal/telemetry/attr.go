


package telemetry


type Attr struct {
	Key   string `json:"key,omitempty"`
	Value Value  `json:"value,omitempty"`
}


func String(key, value string) Attr {
	return Attr{key, StringValue(value)}
}


func Int64(key string, value int64) Attr {
	return Attr{key, Int64Value(value)}
}


func Int(key string, value int) Attr {
	return Int64(key, int64(value))
}


func Float64(key string, value float64) Attr {
	return Attr{key, Float64Value(value)}
}


func Bool(key string, value bool) Attr {
	return Attr{key, BoolValue(value)}
}



func Bytes(key string, value []byte) Attr {
	return Attr{key, BytesValue(value)}
}



func Slice(key string, value ...Value) Attr {
	return Attr{key, SliceValue(value...)}
}



func Map(key string, value ...Attr) Attr {
	return Attr{key, MapValue(value...)}
}


func (a Attr) Equal(b Attr) bool {
	return a.Key == b.Key && a.Value.Equal(b.Value)
}
