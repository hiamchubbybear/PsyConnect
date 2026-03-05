



package proto

















func ValueOrNil[T any](has bool, getter func() T) *T {
	if !has {
		return nil
	}
	v := getter()
	return &v
}

































func ValueOrDefault[T interface {
	*P
	Message
}, P any](val T) T {
	if val == nil {
		return T(new(P))
	}
	return val
}



func ValueOrDefaultBytes(val []byte) []byte {
	if val == nil {
		return []byte{}
	}
	return val
}
