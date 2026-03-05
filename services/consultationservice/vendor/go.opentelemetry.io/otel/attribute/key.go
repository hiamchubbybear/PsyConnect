


package attribute 



type Key string





func (k Key) Bool(v bool) KeyValue {
	return KeyValue{
		Key:   k,
		Value: BoolValue(v),
	}
}





func (k Key) BoolSlice(v []bool) KeyValue {
	return KeyValue{
		Key:   k,
		Value: BoolSliceValue(v),
	}
}





func (k Key) Int(v int) KeyValue {
	return KeyValue{
		Key:   k,
		Value: IntValue(v),
	}
}





func (k Key) IntSlice(v []int) KeyValue {
	return KeyValue{
		Key:   k,
		Value: IntSliceValue(v),
	}
}





func (k Key) Int64(v int64) KeyValue {
	return KeyValue{
		Key:   k,
		Value: Int64Value(v),
	}
}





func (k Key) Int64Slice(v []int64) KeyValue {
	return KeyValue{
		Key:   k,
		Value: Int64SliceValue(v),
	}
}





func (k Key) Float64(v float64) KeyValue {
	return KeyValue{
		Key:   k,
		Value: Float64Value(v),
	}
}





func (k Key) Float64Slice(v []float64) KeyValue {
	return KeyValue{
		Key:   k,
		Value: Float64SliceValue(v),
	}
}





func (k Key) String(v string) KeyValue {
	return KeyValue{
		Key:   k,
		Value: StringValue(v),
	}
}





func (k Key) StringSlice(v []string) KeyValue {
	return KeyValue{
		Key:   k,
		Value: StringSliceValue(v),
	}
}


func (k Key) Defined() bool {
	return len(k) != 0
}
