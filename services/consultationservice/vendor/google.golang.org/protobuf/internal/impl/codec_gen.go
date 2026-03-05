





package impl

import (
	"math"
	"unicode/utf8"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/reflect/protoreflect"
)


func sizeBool(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	v := *p.Bool()
	return f.tagsize + protowire.SizeVarint(protowire.EncodeBool(v))
}


func appendBool(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := *p.Bool()
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendVarint(b, protowire.EncodeBool(v))
	return b, nil
}


func consumeBool(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	if wtyp != protowire.VarintType {
		return out, errUnknown
	}
	var v uint64
	var n int
	if len(b) >= 1 && b[0] < 0x80 {
		v = uint64(b[0])
		n = 1
	} else if len(b) >= 2 && b[1] < 128 {
		v = uint64(b[0]&0x7f) + uint64(b[1])<<7
		n = 2
	} else {
		v, n = protowire.ConsumeVarint(b)
	}
	if n < 0 {
		return out, errDecode
	}
	*p.Bool() = protowire.DecodeBool(v)
	out.n = n
	return out, nil
}

var coderBool = pointerCoderFuncs{
	size:      sizeBool,
	marshal:   appendBool,
	unmarshal: consumeBool,
	merge:     mergeBool,
}



func sizeBoolNoZero(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	v := *p.Bool()
	if v == false {
		return 0
	}
	return f.tagsize + protowire.SizeVarint(protowire.EncodeBool(v))
}



func appendBoolNoZero(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := *p.Bool()
	if v == false {
		return b, nil
	}
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendVarint(b, protowire.EncodeBool(v))
	return b, nil
}

var coderBoolNoZero = pointerCoderFuncs{
	size:      sizeBoolNoZero,
	marshal:   appendBoolNoZero,
	unmarshal: consumeBool,
	merge:     mergeBoolNoZero,
}



func sizeBoolPtr(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	v := **p.BoolPtr()
	return f.tagsize + protowire.SizeVarint(protowire.EncodeBool(v))
}



func appendBoolPtr(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := **p.BoolPtr()
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendVarint(b, protowire.EncodeBool(v))
	return b, nil
}


func consumeBoolPtr(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	if wtyp != protowire.VarintType {
		return out, errUnknown
	}
	var v uint64
	var n int
	if len(b) >= 1 && b[0] < 0x80 {
		v = uint64(b[0])
		n = 1
	} else if len(b) >= 2 && b[1] < 128 {
		v = uint64(b[0]&0x7f) + uint64(b[1])<<7
		n = 2
	} else {
		v, n = protowire.ConsumeVarint(b)
	}
	if n < 0 {
		return out, errDecode
	}
	vp := p.BoolPtr()
	if *vp == nil {
		*vp = new(bool)
	}
	**vp = protowire.DecodeBool(v)
	out.n = n
	return out, nil
}

var coderBoolPtr = pointerCoderFuncs{
	size:      sizeBoolPtr,
	marshal:   appendBoolPtr,
	unmarshal: consumeBoolPtr,
	merge:     mergeBoolPtr,
}


func sizeBoolSlice(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	s := *p.BoolSlice()
	for _, v := range s {
		size += f.tagsize + protowire.SizeVarint(protowire.EncodeBool(v))
	}
	return size
}


func appendBoolSlice(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	s := *p.BoolSlice()
	for _, v := range s {
		b = protowire.AppendVarint(b, f.wiretag)
		b = protowire.AppendVarint(b, protowire.EncodeBool(v))
	}
	return b, nil
}


func consumeBoolSlice(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	sp := p.BoolSlice()
	if wtyp == protowire.BytesType {
		b, n := protowire.ConsumeBytes(b)
		if n < 0 {
			return out, errDecode
		}
		count := 0
		for _, v := range b {
			if v < 0x80 {
				count++
			}
		}
		if count > 0 {
			p.growBoolSlice(count)
		}
		s := *sp
		for len(b) > 0 {
			var v uint64
			var n int
			if len(b) >= 1 && b[0] < 0x80 {
				v = uint64(b[0])
				n = 1
			} else if len(b) >= 2 && b[1] < 128 {
				v = uint64(b[0]&0x7f) + uint64(b[1])<<7
				n = 2
			} else {
				v, n = protowire.ConsumeVarint(b)
			}
			if n < 0 {
				return out, errDecode
			}
			s = append(s, protowire.DecodeBool(v))
			b = b[n:]
		}
		*sp = s
		out.n = n
		return out, nil
	}
	if wtyp != protowire.VarintType {
		return out, errUnknown
	}
	var v uint64
	var n int
	if len(b) >= 1 && b[0] < 0x80 {
		v = uint64(b[0])
		n = 1
	} else if len(b) >= 2 && b[1] < 128 {
		v = uint64(b[0]&0x7f) + uint64(b[1])<<7
		n = 2
	} else {
		v, n = protowire.ConsumeVarint(b)
	}
	if n < 0 {
		return out, errDecode
	}
	*sp = append(*sp, protowire.DecodeBool(v))
	out.n = n
	return out, nil
}

var coderBoolSlice = pointerCoderFuncs{
	size:      sizeBoolSlice,
	marshal:   appendBoolSlice,
	unmarshal: consumeBoolSlice,
	merge:     mergeBoolSlice,
}


func sizeBoolPackedSlice(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	s := *p.BoolSlice()
	if len(s) == 0 {
		return 0
	}
	n := 0
	for _, v := range s {
		n += protowire.SizeVarint(protowire.EncodeBool(v))
	}
	return f.tagsize + protowire.SizeBytes(n)
}


func appendBoolPackedSlice(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	s := *p.BoolSlice()
	if len(s) == 0 {
		return b, nil
	}
	b = protowire.AppendVarint(b, f.wiretag)
	n := 0
	for _, v := range s {
		n += protowire.SizeVarint(protowire.EncodeBool(v))
	}
	b = protowire.AppendVarint(b, uint64(n))
	for _, v := range s {
		b = protowire.AppendVarint(b, protowire.EncodeBool(v))
	}
	return b, nil
}

var coderBoolPackedSlice = pointerCoderFuncs{
	size:      sizeBoolPackedSlice,
	marshal:   appendBoolPackedSlice,
	unmarshal: consumeBoolSlice,
	merge:     mergeBoolSlice,
}


func sizeBoolValue(v protoreflect.Value, tagsize int, opts marshalOptions) int {
	return tagsize + protowire.SizeVarint(protowire.EncodeBool(v.Bool()))
}


func appendBoolValue(b []byte, v protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	b = protowire.AppendVarint(b, wiretag)
	b = protowire.AppendVarint(b, protowire.EncodeBool(v.Bool()))
	return b, nil
}


func consumeBoolValue(b []byte, _ protoreflect.Value, _ protowire.Number, wtyp protowire.Type, opts unmarshalOptions) (_ protoreflect.Value, out unmarshalOutput, err error) {
	if wtyp != protowire.VarintType {
		return protoreflect.Value{}, out, errUnknown
	}
	var v uint64
	var n int
	if len(b) >= 1 && b[0] < 0x80 {
		v = uint64(b[0])
		n = 1
	} else if len(b) >= 2 && b[1] < 128 {
		v = uint64(b[0]&0x7f) + uint64(b[1])<<7
		n = 2
	} else {
		v, n = protowire.ConsumeVarint(b)
	}
	if n < 0 {
		return protoreflect.Value{}, out, errDecode
	}
	out.n = n
	return protoreflect.ValueOfBool(protowire.DecodeBool(v)), out, nil
}

var coderBoolValue = valueCoderFuncs{
	size:      sizeBoolValue,
	marshal:   appendBoolValue,
	unmarshal: consumeBoolValue,
	merge:     mergeScalarValue,
}


func sizeBoolSliceValue(listv protoreflect.Value, tagsize int, opts marshalOptions) (size int) {
	list := listv.List()
	for i, llen := 0, list.Len(); i < llen; i++ {
		v := list.Get(i)
		size += tagsize + protowire.SizeVarint(protowire.EncodeBool(v.Bool()))
	}
	return size
}


func appendBoolSliceValue(b []byte, listv protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	list := listv.List()
	for i, llen := 0, list.Len(); i < llen; i++ {
		v := list.Get(i)
		b = protowire.AppendVarint(b, wiretag)
		b = protowire.AppendVarint(b, protowire.EncodeBool(v.Bool()))
	}
	return b, nil
}


func consumeBoolSliceValue(b []byte, listv protoreflect.Value, _ protowire.Number, wtyp protowire.Type, opts unmarshalOptions) (_ protoreflect.Value, out unmarshalOutput, err error) {
	list := listv.List()
	if wtyp == protowire.BytesType {
		b, n := protowire.ConsumeBytes(b)
		if n < 0 {
			return protoreflect.Value{}, out, errDecode
		}
		for len(b) > 0 {
			var v uint64
			var n int
			if len(b) >= 1 && b[0] < 0x80 {
				v = uint64(b[0])
				n = 1
			} else if len(b) >= 2 && b[1] < 128 {
				v = uint64(b[0]&0x7f) + uint64(b[1])<<7
				n = 2
			} else {
				v, n = protowire.ConsumeVarint(b)
			}
			if n < 0 {
				return protoreflect.Value{}, out, errDecode
			}
			list.Append(protoreflect.ValueOfBool(protowire.DecodeBool(v)))
			b = b[n:]
		}
		out.n = n
		return listv, out, nil
	}
	if wtyp != protowire.VarintType {
		return protoreflect.Value{}, out, errUnknown
	}
	var v uint64
	var n int
	if len(b) >= 1 && b[0] < 0x80 {
		v = uint64(b[0])
		n = 1
	} else if len(b) >= 2 && b[1] < 128 {
		v = uint64(b[0]&0x7f) + uint64(b[1])<<7
		n = 2
	} else {
		v, n = protowire.ConsumeVarint(b)
	}
	if n < 0 {
		return protoreflect.Value{}, out, errDecode
	}
	list.Append(protoreflect.ValueOfBool(protowire.DecodeBool(v)))
	out.n = n
	return listv, out, nil
}

var coderBoolSliceValue = valueCoderFuncs{
	size:      sizeBoolSliceValue,
	marshal:   appendBoolSliceValue,
	unmarshal: consumeBoolSliceValue,
	merge:     mergeListValue,
}


func sizeBoolPackedSliceValue(listv protoreflect.Value, tagsize int, opts marshalOptions) (size int) {
	list := listv.List()
	llen := list.Len()
	if llen == 0 {
		return 0
	}
	n := 0
	for i, llen := 0, llen; i < llen; i++ {
		v := list.Get(i)
		n += protowire.SizeVarint(protowire.EncodeBool(v.Bool()))
	}
	return tagsize + protowire.SizeBytes(n)
}


func appendBoolPackedSliceValue(b []byte, listv protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	list := listv.List()
	llen := list.Len()
	if llen == 0 {
		return b, nil
	}
	b = protowire.AppendVarint(b, wiretag)
	n := 0
	for i := 0; i < llen; i++ {
		v := list.Get(i)
		n += protowire.SizeVarint(protowire.EncodeBool(v.Bool()))
	}
	b = protowire.AppendVarint(b, uint64(n))
	for i := 0; i < llen; i++ {
		v := list.Get(i)
		b = protowire.AppendVarint(b, protowire.EncodeBool(v.Bool()))
	}
	return b, nil
}

var coderBoolPackedSliceValue = valueCoderFuncs{
	size:      sizeBoolPackedSliceValue,
	marshal:   appendBoolPackedSliceValue,
	unmarshal: consumeBoolSliceValue,
	merge:     mergeListValue,
}


func sizeEnumValue(v protoreflect.Value, tagsize int, opts marshalOptions) int {
	return tagsize + protowire.SizeVarint(uint64(v.Enum()))
}


func appendEnumValue(b []byte, v protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	b = protowire.AppendVarint(b, wiretag)
	b = protowire.AppendVarint(b, uint64(v.Enum()))
	return b, nil
}


func consumeEnumValue(b []byte, _ protoreflect.Value, _ protowire.Number, wtyp protowire.Type, opts unmarshalOptions) (_ protoreflect.Value, out unmarshalOutput, err error) {
	if wtyp != protowire.VarintType {
		return protoreflect.Value{}, out, errUnknown
	}
	var v uint64
	var n int
	if len(b) >= 1 && b[0] < 0x80 {
		v = uint64(b[0])
		n = 1
	} else if len(b) >= 2 && b[1] < 128 {
		v = uint64(b[0]&0x7f) + uint64(b[1])<<7
		n = 2
	} else {
		v, n = protowire.ConsumeVarint(b)
	}
	if n < 0 {
		return protoreflect.Value{}, out, errDecode
	}
	out.n = n
	return protoreflect.ValueOfEnum(protoreflect.EnumNumber(v)), out, nil
}

var coderEnumValue = valueCoderFuncs{
	size:      sizeEnumValue,
	marshal:   appendEnumValue,
	unmarshal: consumeEnumValue,
	merge:     mergeScalarValue,
}


func sizeEnumSliceValue(listv protoreflect.Value, tagsize int, opts marshalOptions) (size int) {
	list := listv.List()
	for i, llen := 0, list.Len(); i < llen; i++ {
		v := list.Get(i)
		size += tagsize + protowire.SizeVarint(uint64(v.Enum()))
	}
	return size
}


func appendEnumSliceValue(b []byte, listv protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	list := listv.List()
	for i, llen := 0, list.Len(); i < llen; i++ {
		v := list.Get(i)
		b = protowire.AppendVarint(b, wiretag)
		b = protowire.AppendVarint(b, uint64(v.Enum()))
	}
	return b, nil
}


func consumeEnumSliceValue(b []byte, listv protoreflect.Value, _ protowire.Number, wtyp protowire.Type, opts unmarshalOptions) (_ protoreflect.Value, out unmarshalOutput, err error) {
	list := listv.List()
	if wtyp == protowire.BytesType {
		b, n := protowire.ConsumeBytes(b)
		if n < 0 {
			return protoreflect.Value{}, out, errDecode
		}
		for len(b) > 0 {
			var v uint64
			var n int
			if len(b) >= 1 && b[0] < 0x80 {
				v = uint64(b[0])
				n = 1
			} else if len(b) >= 2 && b[1] < 128 {
				v = uint64(b[0]&0x7f) + uint64(b[1])<<7
				n = 2
			} else {
				v, n = protowire.ConsumeVarint(b)
			}
			if n < 0 {
				return protoreflect.Value{}, out, errDecode
			}
			list.Append(protoreflect.ValueOfEnum(protoreflect.EnumNumber(v)))
			b = b[n:]
		}
		out.n = n
		return listv, out, nil
	}
	if wtyp != protowire.VarintType {
		return protoreflect.Value{}, out, errUnknown
	}
	var v uint64
	var n int
	if len(b) >= 1 && b[0] < 0x80 {
		v = uint64(b[0])
		n = 1
	} else if len(b) >= 2 && b[1] < 128 {
		v = uint64(b[0]&0x7f) + uint64(b[1])<<7
		n = 2
	} else {
		v, n = protowire.ConsumeVarint(b)
	}
	if n < 0 {
		return protoreflect.Value{}, out, errDecode
	}
	list.Append(protoreflect.ValueOfEnum(protoreflect.EnumNumber(v)))
	out.n = n
	return listv, out, nil
}

var coderEnumSliceValue = valueCoderFuncs{
	size:      sizeEnumSliceValue,
	marshal:   appendEnumSliceValue,
	unmarshal: consumeEnumSliceValue,
	merge:     mergeListValue,
}


func sizeEnumPackedSliceValue(listv protoreflect.Value, tagsize int, opts marshalOptions) (size int) {
	list := listv.List()
	llen := list.Len()
	if llen == 0 {
		return 0
	}
	n := 0
	for i, llen := 0, llen; i < llen; i++ {
		v := list.Get(i)
		n += protowire.SizeVarint(uint64(v.Enum()))
	}
	return tagsize + protowire.SizeBytes(n)
}


func appendEnumPackedSliceValue(b []byte, listv protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	list := listv.List()
	llen := list.Len()
	if llen == 0 {
		return b, nil
	}
	b = protowire.AppendVarint(b, wiretag)
	n := 0
	for i := 0; i < llen; i++ {
		v := list.Get(i)
		n += protowire.SizeVarint(uint64(v.Enum()))
	}
	b = protowire.AppendVarint(b, uint64(n))
	for i := 0; i < llen; i++ {
		v := list.Get(i)
		b = protowire.AppendVarint(b, uint64(v.Enum()))
	}
	return b, nil
}

var coderEnumPackedSliceValue = valueCoderFuncs{
	size:      sizeEnumPackedSliceValue,
	marshal:   appendEnumPackedSliceValue,
	unmarshal: consumeEnumSliceValue,
	merge:     mergeListValue,
}


func sizeInt32(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	v := *p.Int32()
	return f.tagsize + protowire.SizeVarint(uint64(v))
}


func appendInt32(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := *p.Int32()
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendVarint(b, uint64(v))
	return b, nil
}


func consumeInt32(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	if wtyp != protowire.VarintType {
		return out, errUnknown
	}
	var v uint64
	var n int
	if len(b) >= 1 && b[0] < 0x80 {
		v = uint64(b[0])
		n = 1
	} else if len(b) >= 2 && b[1] < 128 {
		v = uint64(b[0]&0x7f) + uint64(b[1])<<7
		n = 2
	} else {
		v, n = protowire.ConsumeVarint(b)
	}
	if n < 0 {
		return out, errDecode
	}
	*p.Int32() = int32(v)
	out.n = n
	return out, nil
}

var coderInt32 = pointerCoderFuncs{
	size:      sizeInt32,
	marshal:   appendInt32,
	unmarshal: consumeInt32,
	merge:     mergeInt32,
}



func sizeInt32NoZero(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	v := *p.Int32()
	if v == 0 {
		return 0
	}
	return f.tagsize + protowire.SizeVarint(uint64(v))
}



func appendInt32NoZero(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := *p.Int32()
	if v == 0 {
		return b, nil
	}
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendVarint(b, uint64(v))
	return b, nil
}

var coderInt32NoZero = pointerCoderFuncs{
	size:      sizeInt32NoZero,
	marshal:   appendInt32NoZero,
	unmarshal: consumeInt32,
	merge:     mergeInt32NoZero,
}



func sizeInt32Ptr(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	v := **p.Int32Ptr()
	return f.tagsize + protowire.SizeVarint(uint64(v))
}



func appendInt32Ptr(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := **p.Int32Ptr()
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendVarint(b, uint64(v))
	return b, nil
}


func consumeInt32Ptr(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	if wtyp != protowire.VarintType {
		return out, errUnknown
	}
	var v uint64
	var n int
	if len(b) >= 1 && b[0] < 0x80 {
		v = uint64(b[0])
		n = 1
	} else if len(b) >= 2 && b[1] < 128 {
		v = uint64(b[0]&0x7f) + uint64(b[1])<<7
		n = 2
	} else {
		v, n = protowire.ConsumeVarint(b)
	}
	if n < 0 {
		return out, errDecode
	}
	vp := p.Int32Ptr()
	if *vp == nil {
		*vp = new(int32)
	}
	**vp = int32(v)
	out.n = n
	return out, nil
}

var coderInt32Ptr = pointerCoderFuncs{
	size:      sizeInt32Ptr,
	marshal:   appendInt32Ptr,
	unmarshal: consumeInt32Ptr,
	merge:     mergeInt32Ptr,
}


func sizeInt32Slice(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	s := *p.Int32Slice()
	for _, v := range s {
		size += f.tagsize + protowire.SizeVarint(uint64(v))
	}
	return size
}


func appendInt32Slice(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	s := *p.Int32Slice()
	for _, v := range s {
		b = protowire.AppendVarint(b, f.wiretag)
		b = protowire.AppendVarint(b, uint64(v))
	}
	return b, nil
}


func consumeInt32Slice(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	sp := p.Int32Slice()
	if wtyp == protowire.BytesType {
		b, n := protowire.ConsumeBytes(b)
		if n < 0 {
			return out, errDecode
		}
		count := 0
		for _, v := range b {
			if v < 0x80 {
				count++
			}
		}
		if count > 0 {
			p.growInt32Slice(count)
		}
		s := *sp
		for len(b) > 0 {
			var v uint64
			var n int
			if len(b) >= 1 && b[0] < 0x80 {
				v = uint64(b[0])
				n = 1
			} else if len(b) >= 2 && b[1] < 128 {
				v = uint64(b[0]&0x7f) + uint64(b[1])<<7
				n = 2
			} else {
				v, n = protowire.ConsumeVarint(b)
			}
			if n < 0 {
				return out, errDecode
			}
			s = append(s, int32(v))
			b = b[n:]
		}
		*sp = s
		out.n = n
		return out, nil
	}
	if wtyp != protowire.VarintType {
		return out, errUnknown
	}
	var v uint64
	var n int
	if len(b) >= 1 && b[0] < 0x80 {
		v = uint64(b[0])
		n = 1
	} else if len(b) >= 2 && b[1] < 128 {
		v = uint64(b[0]&0x7f) + uint64(b[1])<<7
		n = 2
	} else {
		v, n = protowire.ConsumeVarint(b)
	}
	if n < 0 {
		return out, errDecode
	}
	*sp = append(*sp, int32(v))
	out.n = n
	return out, nil
}

var coderInt32Slice = pointerCoderFuncs{
	size:      sizeInt32Slice,
	marshal:   appendInt32Slice,
	unmarshal: consumeInt32Slice,
	merge:     mergeInt32Slice,
}


func sizeInt32PackedSlice(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	s := *p.Int32Slice()
	if len(s) == 0 {
		return 0
	}
	n := 0
	for _, v := range s {
		n += protowire.SizeVarint(uint64(v))
	}
	return f.tagsize + protowire.SizeBytes(n)
}


func appendInt32PackedSlice(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	s := *p.Int32Slice()
	if len(s) == 0 {
		return b, nil
	}
	b = protowire.AppendVarint(b, f.wiretag)
	n := 0
	for _, v := range s {
		n += protowire.SizeVarint(uint64(v))
	}
	b = protowire.AppendVarint(b, uint64(n))
	for _, v := range s {
		b = protowire.AppendVarint(b, uint64(v))
	}
	return b, nil
}

var coderInt32PackedSlice = pointerCoderFuncs{
	size:      sizeInt32PackedSlice,
	marshal:   appendInt32PackedSlice,
	unmarshal: consumeInt32Slice,
	merge:     mergeInt32Slice,
}


func sizeInt32Value(v protoreflect.Value, tagsize int, opts marshalOptions) int {
	return tagsize + protowire.SizeVarint(uint64(int32(v.Int())))
}


func appendInt32Value(b []byte, v protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	b = protowire.AppendVarint(b, wiretag)
	b = protowire.AppendVarint(b, uint64(int32(v.Int())))
	return b, nil
}


func consumeInt32Value(b []byte, _ protoreflect.Value, _ protowire.Number, wtyp protowire.Type, opts unmarshalOptions) (_ protoreflect.Value, out unmarshalOutput, err error) {
	if wtyp != protowire.VarintType {
		return protoreflect.Value{}, out, errUnknown
	}
	var v uint64
	var n int
	if len(b) >= 1 && b[0] < 0x80 {
		v = uint64(b[0])
		n = 1
	} else if len(b) >= 2 && b[1] < 128 {
		v = uint64(b[0]&0x7f) + uint64(b[1])<<7
		n = 2
	} else {
		v, n = protowire.ConsumeVarint(b)
	}
	if n < 0 {
		return protoreflect.Value{}, out, errDecode
	}
	out.n = n
	return protoreflect.ValueOfInt32(int32(v)), out, nil
}

var coderInt32Value = valueCoderFuncs{
	size:      sizeInt32Value,
	marshal:   appendInt32Value,
	unmarshal: consumeInt32Value,
	merge:     mergeScalarValue,
}


func sizeInt32SliceValue(listv protoreflect.Value, tagsize int, opts marshalOptions) (size int) {
	list := listv.List()
	for i, llen := 0, list.Len(); i < llen; i++ {
		v := list.Get(i)
		size += tagsize + protowire.SizeVarint(uint64(int32(v.Int())))
	}
	return size
}


func appendInt32SliceValue(b []byte, listv protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	list := listv.List()
	for i, llen := 0, list.Len(); i < llen; i++ {
		v := list.Get(i)
		b = protowire.AppendVarint(b, wiretag)
		b = protowire.AppendVarint(b, uint64(int32(v.Int())))
	}
	return b, nil
}


func consumeInt32SliceValue(b []byte, listv protoreflect.Value, _ protowire.Number, wtyp protowire.Type, opts unmarshalOptions) (_ protoreflect.Value, out unmarshalOutput, err error) {
	list := listv.List()
	if wtyp == protowire.BytesType {
		b, n := protowire.ConsumeBytes(b)
		if n < 0 {
			return protoreflect.Value{}, out, errDecode
		}
		for len(b) > 0 {
			var v uint64
			var n int
			if len(b) >= 1 && b[0] < 0x80 {
				v = uint64(b[0])
				n = 1
			} else if len(b) >= 2 && b[1] < 128 {
				v = uint64(b[0]&0x7f) + uint64(b[1])<<7
				n = 2
			} else {
				v, n = protowire.ConsumeVarint(b)
			}
			if n < 0 {
				return protoreflect.Value{}, out, errDecode
			}
			list.Append(protoreflect.ValueOfInt32(int32(v)))
			b = b[n:]
		}
		out.n = n
		return listv, out, nil
	}
	if wtyp != protowire.VarintType {
		return protoreflect.Value{}, out, errUnknown
	}
	var v uint64
	var n int
	if len(b) >= 1 && b[0] < 0x80 {
		v = uint64(b[0])
		n = 1
	} else if len(b) >= 2 && b[1] < 128 {
		v = uint64(b[0]&0x7f) + uint64(b[1])<<7
		n = 2
	} else {
		v, n = protowire.ConsumeVarint(b)
	}
	if n < 0 {
		return protoreflect.Value{}, out, errDecode
	}
	list.Append(protoreflect.ValueOfInt32(int32(v)))
	out.n = n
	return listv, out, nil
}

var coderInt32SliceValue = valueCoderFuncs{
	size:      sizeInt32SliceValue,
	marshal:   appendInt32SliceValue,
	unmarshal: consumeInt32SliceValue,
	merge:     mergeListValue,
}


func sizeInt32PackedSliceValue(listv protoreflect.Value, tagsize int, opts marshalOptions) (size int) {
	list := listv.List()
	llen := list.Len()
	if llen == 0 {
		return 0
	}
	n := 0
	for i, llen := 0, llen; i < llen; i++ {
		v := list.Get(i)
		n += protowire.SizeVarint(uint64(int32(v.Int())))
	}
	return tagsize + protowire.SizeBytes(n)
}


func appendInt32PackedSliceValue(b []byte, listv protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	list := listv.List()
	llen := list.Len()
	if llen == 0 {
		return b, nil
	}
	b = protowire.AppendVarint(b, wiretag)
	n := 0
	for i := 0; i < llen; i++ {
		v := list.Get(i)
		n += protowire.SizeVarint(uint64(int32(v.Int())))
	}
	b = protowire.AppendVarint(b, uint64(n))
	for i := 0; i < llen; i++ {
		v := list.Get(i)
		b = protowire.AppendVarint(b, uint64(int32(v.Int())))
	}
	return b, nil
}

var coderInt32PackedSliceValue = valueCoderFuncs{
	size:      sizeInt32PackedSliceValue,
	marshal:   appendInt32PackedSliceValue,
	unmarshal: consumeInt32SliceValue,
	merge:     mergeListValue,
}


func sizeSint32(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	v := *p.Int32()
	return f.tagsize + protowire.SizeVarint(protowire.EncodeZigZag(int64(v)))
}


func appendSint32(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := *p.Int32()
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendVarint(b, protowire.EncodeZigZag(int64(v)))
	return b, nil
}


func consumeSint32(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	if wtyp != protowire.VarintType {
		return out, errUnknown
	}
	var v uint64
	var n int
	if len(b) >= 1 && b[0] < 0x80 {
		v = uint64(b[0])
		n = 1
	} else if len(b) >= 2 && b[1] < 128 {
		v = uint64(b[0]&0x7f) + uint64(b[1])<<7
		n = 2
	} else {
		v, n = protowire.ConsumeVarint(b)
	}
	if n < 0 {
		return out, errDecode
	}
	*p.Int32() = int32(protowire.DecodeZigZag(v & math.MaxUint32))
	out.n = n
	return out, nil
}

var coderSint32 = pointerCoderFuncs{
	size:      sizeSint32,
	marshal:   appendSint32,
	unmarshal: consumeSint32,
	merge:     mergeInt32,
}



func sizeSint32NoZero(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	v := *p.Int32()
	if v == 0 {
		return 0
	}
	return f.tagsize + protowire.SizeVarint(protowire.EncodeZigZag(int64(v)))
}



func appendSint32NoZero(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := *p.Int32()
	if v == 0 {
		return b, nil
	}
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendVarint(b, protowire.EncodeZigZag(int64(v)))
	return b, nil
}

var coderSint32NoZero = pointerCoderFuncs{
	size:      sizeSint32NoZero,
	marshal:   appendSint32NoZero,
	unmarshal: consumeSint32,
	merge:     mergeInt32NoZero,
}



func sizeSint32Ptr(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	v := **p.Int32Ptr()
	return f.tagsize + protowire.SizeVarint(protowire.EncodeZigZag(int64(v)))
}



func appendSint32Ptr(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := **p.Int32Ptr()
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendVarint(b, protowire.EncodeZigZag(int64(v)))
	return b, nil
}


func consumeSint32Ptr(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	if wtyp != protowire.VarintType {
		return out, errUnknown
	}
	var v uint64
	var n int
	if len(b) >= 1 && b[0] < 0x80 {
		v = uint64(b[0])
		n = 1
	} else if len(b) >= 2 && b[1] < 128 {
		v = uint64(b[0]&0x7f) + uint64(b[1])<<7
		n = 2
	} else {
		v, n = protowire.ConsumeVarint(b)
	}
	if n < 0 {
		return out, errDecode
	}
	vp := p.Int32Ptr()
	if *vp == nil {
		*vp = new(int32)
	}
	**vp = int32(protowire.DecodeZigZag(v & math.MaxUint32))
	out.n = n
	return out, nil
}

var coderSint32Ptr = pointerCoderFuncs{
	size:      sizeSint32Ptr,
	marshal:   appendSint32Ptr,
	unmarshal: consumeSint32Ptr,
	merge:     mergeInt32Ptr,
}


func sizeSint32Slice(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	s := *p.Int32Slice()
	for _, v := range s {
		size += f.tagsize + protowire.SizeVarint(protowire.EncodeZigZag(int64(v)))
	}
	return size
}


func appendSint32Slice(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	s := *p.Int32Slice()
	for _, v := range s {
		b = protowire.AppendVarint(b, f.wiretag)
		b = protowire.AppendVarint(b, protowire.EncodeZigZag(int64(v)))
	}
	return b, nil
}


func consumeSint32Slice(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	sp := p.Int32Slice()
	if wtyp == protowire.BytesType {
		b, n := protowire.ConsumeBytes(b)
		if n < 0 {
			return out, errDecode
		}
		count := 0
		for _, v := range b {
			if v < 0x80 {
				count++
			}
		}
		if count > 0 {
			p.growInt32Slice(count)
		}
		s := *sp
		for len(b) > 0 {
			var v uint64
			var n int
			if len(b) >= 1 && b[0] < 0x80 {
				v = uint64(b[0])
				n = 1
			} else if len(b) >= 2 && b[1] < 128 {
				v = uint64(b[0]&0x7f) + uint64(b[1])<<7
				n = 2
			} else {
				v, n = protowire.ConsumeVarint(b)
			}
			if n < 0 {
				return out, errDecode
			}
			s = append(s, int32(protowire.DecodeZigZag(v&math.MaxUint32)))
			b = b[n:]
		}
		*sp = s
		out.n = n
		return out, nil
	}
	if wtyp != protowire.VarintType {
		return out, errUnknown
	}
	var v uint64
	var n int
	if len(b) >= 1 && b[0] < 0x80 {
		v = uint64(b[0])
		n = 1
	} else if len(b) >= 2 && b[1] < 128 {
		v = uint64(b[0]&0x7f) + uint64(b[1])<<7
		n = 2
	} else {
		v, n = protowire.ConsumeVarint(b)
	}
	if n < 0 {
		return out, errDecode
	}
	*sp = append(*sp, int32(protowire.DecodeZigZag(v&math.MaxUint32)))
	out.n = n
	return out, nil
}

var coderSint32Slice = pointerCoderFuncs{
	size:      sizeSint32Slice,
	marshal:   appendSint32Slice,
	unmarshal: consumeSint32Slice,
	merge:     mergeInt32Slice,
}


func sizeSint32PackedSlice(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	s := *p.Int32Slice()
	if len(s) == 0 {
		return 0
	}
	n := 0
	for _, v := range s {
		n += protowire.SizeVarint(protowire.EncodeZigZag(int64(v)))
	}
	return f.tagsize + protowire.SizeBytes(n)
}


func appendSint32PackedSlice(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	s := *p.Int32Slice()
	if len(s) == 0 {
		return b, nil
	}
	b = protowire.AppendVarint(b, f.wiretag)
	n := 0
	for _, v := range s {
		n += protowire.SizeVarint(protowire.EncodeZigZag(int64(v)))
	}
	b = protowire.AppendVarint(b, uint64(n))
	for _, v := range s {
		b = protowire.AppendVarint(b, protowire.EncodeZigZag(int64(v)))
	}
	return b, nil
}

var coderSint32PackedSlice = pointerCoderFuncs{
	size:      sizeSint32PackedSlice,
	marshal:   appendSint32PackedSlice,
	unmarshal: consumeSint32Slice,
	merge:     mergeInt32Slice,
}


func sizeSint32Value(v protoreflect.Value, tagsize int, opts marshalOptions) int {
	return tagsize + protowire.SizeVarint(protowire.EncodeZigZag(int64(int32(v.Int()))))
}


func appendSint32Value(b []byte, v protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	b = protowire.AppendVarint(b, wiretag)
	b = protowire.AppendVarint(b, protowire.EncodeZigZag(int64(int32(v.Int()))))
	return b, nil
}


func consumeSint32Value(b []byte, _ protoreflect.Value, _ protowire.Number, wtyp protowire.Type, opts unmarshalOptions) (_ protoreflect.Value, out unmarshalOutput, err error) {
	if wtyp != protowire.VarintType {
		return protoreflect.Value{}, out, errUnknown
	}
	var v uint64
	var n int
	if len(b) >= 1 && b[0] < 0x80 {
		v = uint64(b[0])
		n = 1
	} else if len(b) >= 2 && b[1] < 128 {
		v = uint64(b[0]&0x7f) + uint64(b[1])<<7
		n = 2
	} else {
		v, n = protowire.ConsumeVarint(b)
	}
	if n < 0 {
		return protoreflect.Value{}, out, errDecode
	}
	out.n = n
	return protoreflect.ValueOfInt32(int32(protowire.DecodeZigZag(v & math.MaxUint32))), out, nil
}

var coderSint32Value = valueCoderFuncs{
	size:      sizeSint32Value,
	marshal:   appendSint32Value,
	unmarshal: consumeSint32Value,
	merge:     mergeScalarValue,
}


func sizeSint32SliceValue(listv protoreflect.Value, tagsize int, opts marshalOptions) (size int) {
	list := listv.List()
	for i, llen := 0, list.Len(); i < llen; i++ {
		v := list.Get(i)
		size += tagsize + protowire.SizeVarint(protowire.EncodeZigZag(int64(int32(v.Int()))))
	}
	return size
}


func appendSint32SliceValue(b []byte, listv protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	list := listv.List()
	for i, llen := 0, list.Len(); i < llen; i++ {
		v := list.Get(i)
		b = protowire.AppendVarint(b, wiretag)
		b = protowire.AppendVarint(b, protowire.EncodeZigZag(int64(int32(v.Int()))))
	}
	return b, nil
}


func consumeSint32SliceValue(b []byte, listv protoreflect.Value, _ protowire.Number, wtyp protowire.Type, opts unmarshalOptions) (_ protoreflect.Value, out unmarshalOutput, err error) {
	list := listv.List()
	if wtyp == protowire.BytesType {
		b, n := protowire.ConsumeBytes(b)
		if n < 0 {
			return protoreflect.Value{}, out, errDecode
		}
		for len(b) > 0 {
			var v uint64
			var n int
			if len(b) >= 1 && b[0] < 0x80 {
				v = uint64(b[0])
				n = 1
			} else if len(b) >= 2 && b[1] < 128 {
				v = uint64(b[0]&0x7f) + uint64(b[1])<<7
				n = 2
			} else {
				v, n = protowire.ConsumeVarint(b)
			}
			if n < 0 {
				return protoreflect.Value{}, out, errDecode
			}
			list.Append(protoreflect.ValueOfInt32(int32(protowire.DecodeZigZag(v & math.MaxUint32))))
			b = b[n:]
		}
		out.n = n
		return listv, out, nil
	}
	if wtyp != protowire.VarintType {
		return protoreflect.Value{}, out, errUnknown
	}
	var v uint64
	var n int
	if len(b) >= 1 && b[0] < 0x80 {
		v = uint64(b[0])
		n = 1
	} else if len(b) >= 2 && b[1] < 128 {
		v = uint64(b[0]&0x7f) + uint64(b[1])<<7
		n = 2
	} else {
		v, n = protowire.ConsumeVarint(b)
	}
	if n < 0 {
		return protoreflect.Value{}, out, errDecode
	}
	list.Append(protoreflect.ValueOfInt32(int32(protowire.DecodeZigZag(v & math.MaxUint32))))
	out.n = n
	return listv, out, nil
}

var coderSint32SliceValue = valueCoderFuncs{
	size:      sizeSint32SliceValue,
	marshal:   appendSint32SliceValue,
	unmarshal: consumeSint32SliceValue,
	merge:     mergeListValue,
}


func sizeSint32PackedSliceValue(listv protoreflect.Value, tagsize int, opts marshalOptions) (size int) {
	list := listv.List()
	llen := list.Len()
	if llen == 0 {
		return 0
	}
	n := 0
	for i, llen := 0, llen; i < llen; i++ {
		v := list.Get(i)
		n += protowire.SizeVarint(protowire.EncodeZigZag(int64(int32(v.Int()))))
	}
	return tagsize + protowire.SizeBytes(n)
}


func appendSint32PackedSliceValue(b []byte, listv protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	list := listv.List()
	llen := list.Len()
	if llen == 0 {
		return b, nil
	}
	b = protowire.AppendVarint(b, wiretag)
	n := 0
	for i := 0; i < llen; i++ {
		v := list.Get(i)
		n += protowire.SizeVarint(protowire.EncodeZigZag(int64(int32(v.Int()))))
	}
	b = protowire.AppendVarint(b, uint64(n))
	for i := 0; i < llen; i++ {
		v := list.Get(i)
		b = protowire.AppendVarint(b, protowire.EncodeZigZag(int64(int32(v.Int()))))
	}
	return b, nil
}

var coderSint32PackedSliceValue = valueCoderFuncs{
	size:      sizeSint32PackedSliceValue,
	marshal:   appendSint32PackedSliceValue,
	unmarshal: consumeSint32SliceValue,
	merge:     mergeListValue,
}


func sizeUint32(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	v := *p.Uint32()
	return f.tagsize + protowire.SizeVarint(uint64(v))
}


func appendUint32(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := *p.Uint32()
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendVarint(b, uint64(v))
	return b, nil
}


func consumeUint32(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	if wtyp != protowire.VarintType {
		return out, errUnknown
	}
	var v uint64
	var n int
	if len(b) >= 1 && b[0] < 0x80 {
		v = uint64(b[0])
		n = 1
	} else if len(b) >= 2 && b[1] < 128 {
		v = uint64(b[0]&0x7f) + uint64(b[1])<<7
		n = 2
	} else {
		v, n = protowire.ConsumeVarint(b)
	}
	if n < 0 {
		return out, errDecode
	}
	*p.Uint32() = uint32(v)
	out.n = n
	return out, nil
}

var coderUint32 = pointerCoderFuncs{
	size:      sizeUint32,
	marshal:   appendUint32,
	unmarshal: consumeUint32,
	merge:     mergeUint32,
}



func sizeUint32NoZero(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	v := *p.Uint32()
	if v == 0 {
		return 0
	}
	return f.tagsize + protowire.SizeVarint(uint64(v))
}



func appendUint32NoZero(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := *p.Uint32()
	if v == 0 {
		return b, nil
	}
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendVarint(b, uint64(v))
	return b, nil
}

var coderUint32NoZero = pointerCoderFuncs{
	size:      sizeUint32NoZero,
	marshal:   appendUint32NoZero,
	unmarshal: consumeUint32,
	merge:     mergeUint32NoZero,
}



func sizeUint32Ptr(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	v := **p.Uint32Ptr()
	return f.tagsize + protowire.SizeVarint(uint64(v))
}



func appendUint32Ptr(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := **p.Uint32Ptr()
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendVarint(b, uint64(v))
	return b, nil
}


func consumeUint32Ptr(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	if wtyp != protowire.VarintType {
		return out, errUnknown
	}
	var v uint64
	var n int
	if len(b) >= 1 && b[0] < 0x80 {
		v = uint64(b[0])
		n = 1
	} else if len(b) >= 2 && b[1] < 128 {
		v = uint64(b[0]&0x7f) + uint64(b[1])<<7
		n = 2
	} else {
		v, n = protowire.ConsumeVarint(b)
	}
	if n < 0 {
		return out, errDecode
	}
	vp := p.Uint32Ptr()
	if *vp == nil {
		*vp = new(uint32)
	}
	**vp = uint32(v)
	out.n = n
	return out, nil
}

var coderUint32Ptr = pointerCoderFuncs{
	size:      sizeUint32Ptr,
	marshal:   appendUint32Ptr,
	unmarshal: consumeUint32Ptr,
	merge:     mergeUint32Ptr,
}


func sizeUint32Slice(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	s := *p.Uint32Slice()
	for _, v := range s {
		size += f.tagsize + protowire.SizeVarint(uint64(v))
	}
	return size
}


func appendUint32Slice(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	s := *p.Uint32Slice()
	for _, v := range s {
		b = protowire.AppendVarint(b, f.wiretag)
		b = protowire.AppendVarint(b, uint64(v))
	}
	return b, nil
}


func consumeUint32Slice(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	sp := p.Uint32Slice()
	if wtyp == protowire.BytesType {
		b, n := protowire.ConsumeBytes(b)
		if n < 0 {
			return out, errDecode
		}
		count := 0
		for _, v := range b {
			if v < 0x80 {
				count++
			}
		}
		if count > 0 {
			p.growUint32Slice(count)
		}
		s := *sp
		for len(b) > 0 {
			var v uint64
			var n int
			if len(b) >= 1 && b[0] < 0x80 {
				v = uint64(b[0])
				n = 1
			} else if len(b) >= 2 && b[1] < 128 {
				v = uint64(b[0]&0x7f) + uint64(b[1])<<7
				n = 2
			} else {
				v, n = protowire.ConsumeVarint(b)
			}
			if n < 0 {
				return out, errDecode
			}
			s = append(s, uint32(v))
			b = b[n:]
		}
		*sp = s
		out.n = n
		return out, nil
	}
	if wtyp != protowire.VarintType {
		return out, errUnknown
	}
	var v uint64
	var n int
	if len(b) >= 1 && b[0] < 0x80 {
		v = uint64(b[0])
		n = 1
	} else if len(b) >= 2 && b[1] < 128 {
		v = uint64(b[0]&0x7f) + uint64(b[1])<<7
		n = 2
	} else {
		v, n = protowire.ConsumeVarint(b)
	}
	if n < 0 {
		return out, errDecode
	}
	*sp = append(*sp, uint32(v))
	out.n = n
	return out, nil
}

var coderUint32Slice = pointerCoderFuncs{
	size:      sizeUint32Slice,
	marshal:   appendUint32Slice,
	unmarshal: consumeUint32Slice,
	merge:     mergeUint32Slice,
}


func sizeUint32PackedSlice(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	s := *p.Uint32Slice()
	if len(s) == 0 {
		return 0
	}
	n := 0
	for _, v := range s {
		n += protowire.SizeVarint(uint64(v))
	}
	return f.tagsize + protowire.SizeBytes(n)
}


func appendUint32PackedSlice(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	s := *p.Uint32Slice()
	if len(s) == 0 {
		return b, nil
	}
	b = protowire.AppendVarint(b, f.wiretag)
	n := 0
	for _, v := range s {
		n += protowire.SizeVarint(uint64(v))
	}
	b = protowire.AppendVarint(b, uint64(n))
	for _, v := range s {
		b = protowire.AppendVarint(b, uint64(v))
	}
	return b, nil
}

var coderUint32PackedSlice = pointerCoderFuncs{
	size:      sizeUint32PackedSlice,
	marshal:   appendUint32PackedSlice,
	unmarshal: consumeUint32Slice,
	merge:     mergeUint32Slice,
}


func sizeUint32Value(v protoreflect.Value, tagsize int, opts marshalOptions) int {
	return tagsize + protowire.SizeVarint(uint64(uint32(v.Uint())))
}


func appendUint32Value(b []byte, v protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	b = protowire.AppendVarint(b, wiretag)
	b = protowire.AppendVarint(b, uint64(uint32(v.Uint())))
	return b, nil
}


func consumeUint32Value(b []byte, _ protoreflect.Value, _ protowire.Number, wtyp protowire.Type, opts unmarshalOptions) (_ protoreflect.Value, out unmarshalOutput, err error) {
	if wtyp != protowire.VarintType {
		return protoreflect.Value{}, out, errUnknown
	}
	var v uint64
	var n int
	if len(b) >= 1 && b[0] < 0x80 {
		v = uint64(b[0])
		n = 1
	} else if len(b) >= 2 && b[1] < 128 {
		v = uint64(b[0]&0x7f) + uint64(b[1])<<7
		n = 2
	} else {
		v, n = protowire.ConsumeVarint(b)
	}
	if n < 0 {
		return protoreflect.Value{}, out, errDecode
	}
	out.n = n
	return protoreflect.ValueOfUint32(uint32(v)), out, nil
}

var coderUint32Value = valueCoderFuncs{
	size:      sizeUint32Value,
	marshal:   appendUint32Value,
	unmarshal: consumeUint32Value,
	merge:     mergeScalarValue,
}


func sizeUint32SliceValue(listv protoreflect.Value, tagsize int, opts marshalOptions) (size int) {
	list := listv.List()
	for i, llen := 0, list.Len(); i < llen; i++ {
		v := list.Get(i)
		size += tagsize + protowire.SizeVarint(uint64(uint32(v.Uint())))
	}
	return size
}


func appendUint32SliceValue(b []byte, listv protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	list := listv.List()
	for i, llen := 0, list.Len(); i < llen; i++ {
		v := list.Get(i)
		b = protowire.AppendVarint(b, wiretag)
		b = protowire.AppendVarint(b, uint64(uint32(v.Uint())))
	}
	return b, nil
}


func consumeUint32SliceValue(b []byte, listv protoreflect.Value, _ protowire.Number, wtyp protowire.Type, opts unmarshalOptions) (_ protoreflect.Value, out unmarshalOutput, err error) {
	list := listv.List()
	if wtyp == protowire.BytesType {
		b, n := protowire.ConsumeBytes(b)
		if n < 0 {
			return protoreflect.Value{}, out, errDecode
		}
		for len(b) > 0 {
			var v uint64
			var n int
			if len(b) >= 1 && b[0] < 0x80 {
				v = uint64(b[0])
				n = 1
			} else if len(b) >= 2 && b[1] < 128 {
				v = uint64(b[0]&0x7f) + uint64(b[1])<<7
				n = 2
			} else {
				v, n = protowire.ConsumeVarint(b)
			}
			if n < 0 {
				return protoreflect.Value{}, out, errDecode
			}
			list.Append(protoreflect.ValueOfUint32(uint32(v)))
			b = b[n:]
		}
		out.n = n
		return listv, out, nil
	}
	if wtyp != protowire.VarintType {
		return protoreflect.Value{}, out, errUnknown
	}
	var v uint64
	var n int
	if len(b) >= 1 && b[0] < 0x80 {
		v = uint64(b[0])
		n = 1
	} else if len(b) >= 2 && b[1] < 128 {
		v = uint64(b[0]&0x7f) + uint64(b[1])<<7
		n = 2
	} else {
		v, n = protowire.ConsumeVarint(b)
	}
	if n < 0 {
		return protoreflect.Value{}, out, errDecode
	}
	list.Append(protoreflect.ValueOfUint32(uint32(v)))
	out.n = n
	return listv, out, nil
}

var coderUint32SliceValue = valueCoderFuncs{
	size:      sizeUint32SliceValue,
	marshal:   appendUint32SliceValue,
	unmarshal: consumeUint32SliceValue,
	merge:     mergeListValue,
}


func sizeUint32PackedSliceValue(listv protoreflect.Value, tagsize int, opts marshalOptions) (size int) {
	list := listv.List()
	llen := list.Len()
	if llen == 0 {
		return 0
	}
	n := 0
	for i, llen := 0, llen; i < llen; i++ {
		v := list.Get(i)
		n += protowire.SizeVarint(uint64(uint32(v.Uint())))
	}
	return tagsize + protowire.SizeBytes(n)
}


func appendUint32PackedSliceValue(b []byte, listv protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	list := listv.List()
	llen := list.Len()
	if llen == 0 {
		return b, nil
	}
	b = protowire.AppendVarint(b, wiretag)
	n := 0
	for i := 0; i < llen; i++ {
		v := list.Get(i)
		n += protowire.SizeVarint(uint64(uint32(v.Uint())))
	}
	b = protowire.AppendVarint(b, uint64(n))
	for i := 0; i < llen; i++ {
		v := list.Get(i)
		b = protowire.AppendVarint(b, uint64(uint32(v.Uint())))
	}
	return b, nil
}

var coderUint32PackedSliceValue = valueCoderFuncs{
	size:      sizeUint32PackedSliceValue,
	marshal:   appendUint32PackedSliceValue,
	unmarshal: consumeUint32SliceValue,
	merge:     mergeListValue,
}


func sizeInt64(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	v := *p.Int64()
	return f.tagsize + protowire.SizeVarint(uint64(v))
}


func appendInt64(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := *p.Int64()
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendVarint(b, uint64(v))
	return b, nil
}


func consumeInt64(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	if wtyp != protowire.VarintType {
		return out, errUnknown
	}
	var v uint64
	var n int
	if len(b) >= 1 && b[0] < 0x80 {
		v = uint64(b[0])
		n = 1
	} else if len(b) >= 2 && b[1] < 128 {
		v = uint64(b[0]&0x7f) + uint64(b[1])<<7
		n = 2
	} else {
		v, n = protowire.ConsumeVarint(b)
	}
	if n < 0 {
		return out, errDecode
	}
	*p.Int64() = int64(v)
	out.n = n
	return out, nil
}

var coderInt64 = pointerCoderFuncs{
	size:      sizeInt64,
	marshal:   appendInt64,
	unmarshal: consumeInt64,
	merge:     mergeInt64,
}



func sizeInt64NoZero(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	v := *p.Int64()
	if v == 0 {
		return 0
	}
	return f.tagsize + protowire.SizeVarint(uint64(v))
}



func appendInt64NoZero(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := *p.Int64()
	if v == 0 {
		return b, nil
	}
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendVarint(b, uint64(v))
	return b, nil
}

var coderInt64NoZero = pointerCoderFuncs{
	size:      sizeInt64NoZero,
	marshal:   appendInt64NoZero,
	unmarshal: consumeInt64,
	merge:     mergeInt64NoZero,
}



func sizeInt64Ptr(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	v := **p.Int64Ptr()
	return f.tagsize + protowire.SizeVarint(uint64(v))
}



func appendInt64Ptr(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := **p.Int64Ptr()
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendVarint(b, uint64(v))
	return b, nil
}


func consumeInt64Ptr(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	if wtyp != protowire.VarintType {
		return out, errUnknown
	}
	var v uint64
	var n int
	if len(b) >= 1 && b[0] < 0x80 {
		v = uint64(b[0])
		n = 1
	} else if len(b) >= 2 && b[1] < 128 {
		v = uint64(b[0]&0x7f) + uint64(b[1])<<7
		n = 2
	} else {
		v, n = protowire.ConsumeVarint(b)
	}
	if n < 0 {
		return out, errDecode
	}
	vp := p.Int64Ptr()
	if *vp == nil {
		*vp = new(int64)
	}
	**vp = int64(v)
	out.n = n
	return out, nil
}

var coderInt64Ptr = pointerCoderFuncs{
	size:      sizeInt64Ptr,
	marshal:   appendInt64Ptr,
	unmarshal: consumeInt64Ptr,
	merge:     mergeInt64Ptr,
}


func sizeInt64Slice(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	s := *p.Int64Slice()
	for _, v := range s {
		size += f.tagsize + protowire.SizeVarint(uint64(v))
	}
	return size
}


func appendInt64Slice(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	s := *p.Int64Slice()
	for _, v := range s {
		b = protowire.AppendVarint(b, f.wiretag)
		b = protowire.AppendVarint(b, uint64(v))
	}
	return b, nil
}


func consumeInt64Slice(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	sp := p.Int64Slice()
	if wtyp == protowire.BytesType {
		b, n := protowire.ConsumeBytes(b)
		if n < 0 {
			return out, errDecode
		}
		count := 0
		for _, v := range b {
			if v < 0x80 {
				count++
			}
		}
		if count > 0 {
			p.growInt64Slice(count)
		}
		s := *sp
		for len(b) > 0 {
			var v uint64
			var n int
			if len(b) >= 1 && b[0] < 0x80 {
				v = uint64(b[0])
				n = 1
			} else if len(b) >= 2 && b[1] < 128 {
				v = uint64(b[0]&0x7f) + uint64(b[1])<<7
				n = 2
			} else {
				v, n = protowire.ConsumeVarint(b)
			}
			if n < 0 {
				return out, errDecode
			}
			s = append(s, int64(v))
			b = b[n:]
		}
		*sp = s
		out.n = n
		return out, nil
	}
	if wtyp != protowire.VarintType {
		return out, errUnknown
	}
	var v uint64
	var n int
	if len(b) >= 1 && b[0] < 0x80 {
		v = uint64(b[0])
		n = 1
	} else if len(b) >= 2 && b[1] < 128 {
		v = uint64(b[0]&0x7f) + uint64(b[1])<<7
		n = 2
	} else {
		v, n = protowire.ConsumeVarint(b)
	}
	if n < 0 {
		return out, errDecode
	}
	*sp = append(*sp, int64(v))
	out.n = n
	return out, nil
}

var coderInt64Slice = pointerCoderFuncs{
	size:      sizeInt64Slice,
	marshal:   appendInt64Slice,
	unmarshal: consumeInt64Slice,
	merge:     mergeInt64Slice,
}


func sizeInt64PackedSlice(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	s := *p.Int64Slice()
	if len(s) == 0 {
		return 0
	}
	n := 0
	for _, v := range s {
		n += protowire.SizeVarint(uint64(v))
	}
	return f.tagsize + protowire.SizeBytes(n)
}


func appendInt64PackedSlice(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	s := *p.Int64Slice()
	if len(s) == 0 {
		return b, nil
	}
	b = protowire.AppendVarint(b, f.wiretag)
	n := 0
	for _, v := range s {
		n += protowire.SizeVarint(uint64(v))
	}
	b = protowire.AppendVarint(b, uint64(n))
	for _, v := range s {
		b = protowire.AppendVarint(b, uint64(v))
	}
	return b, nil
}

var coderInt64PackedSlice = pointerCoderFuncs{
	size:      sizeInt64PackedSlice,
	marshal:   appendInt64PackedSlice,
	unmarshal: consumeInt64Slice,
	merge:     mergeInt64Slice,
}


func sizeInt64Value(v protoreflect.Value, tagsize int, opts marshalOptions) int {
	return tagsize + protowire.SizeVarint(uint64(v.Int()))
}


func appendInt64Value(b []byte, v protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	b = protowire.AppendVarint(b, wiretag)
	b = protowire.AppendVarint(b, uint64(v.Int()))
	return b, nil
}


func consumeInt64Value(b []byte, _ protoreflect.Value, _ protowire.Number, wtyp protowire.Type, opts unmarshalOptions) (_ protoreflect.Value, out unmarshalOutput, err error) {
	if wtyp != protowire.VarintType {
		return protoreflect.Value{}, out, errUnknown
	}
	var v uint64
	var n int
	if len(b) >= 1 && b[0] < 0x80 {
		v = uint64(b[0])
		n = 1
	} else if len(b) >= 2 && b[1] < 128 {
		v = uint64(b[0]&0x7f) + uint64(b[1])<<7
		n = 2
	} else {
		v, n = protowire.ConsumeVarint(b)
	}
	if n < 0 {
		return protoreflect.Value{}, out, errDecode
	}
	out.n = n
	return protoreflect.ValueOfInt64(int64(v)), out, nil
}

var coderInt64Value = valueCoderFuncs{
	size:      sizeInt64Value,
	marshal:   appendInt64Value,
	unmarshal: consumeInt64Value,
	merge:     mergeScalarValue,
}


func sizeInt64SliceValue(listv protoreflect.Value, tagsize int, opts marshalOptions) (size int) {
	list := listv.List()
	for i, llen := 0, list.Len(); i < llen; i++ {
		v := list.Get(i)
		size += tagsize + protowire.SizeVarint(uint64(v.Int()))
	}
	return size
}


func appendInt64SliceValue(b []byte, listv protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	list := listv.List()
	for i, llen := 0, list.Len(); i < llen; i++ {
		v := list.Get(i)
		b = protowire.AppendVarint(b, wiretag)
		b = protowire.AppendVarint(b, uint64(v.Int()))
	}
	return b, nil
}


func consumeInt64SliceValue(b []byte, listv protoreflect.Value, _ protowire.Number, wtyp protowire.Type, opts unmarshalOptions) (_ protoreflect.Value, out unmarshalOutput, err error) {
	list := listv.List()
	if wtyp == protowire.BytesType {
		b, n := protowire.ConsumeBytes(b)
		if n < 0 {
			return protoreflect.Value{}, out, errDecode
		}
		for len(b) > 0 {
			var v uint64
			var n int
			if len(b) >= 1 && b[0] < 0x80 {
				v = uint64(b[0])
				n = 1
			} else if len(b) >= 2 && b[1] < 128 {
				v = uint64(b[0]&0x7f) + uint64(b[1])<<7
				n = 2
			} else {
				v, n = protowire.ConsumeVarint(b)
			}
			if n < 0 {
				return protoreflect.Value{}, out, errDecode
			}
			list.Append(protoreflect.ValueOfInt64(int64(v)))
			b = b[n:]
		}
		out.n = n
		return listv, out, nil
	}
	if wtyp != protowire.VarintType {
		return protoreflect.Value{}, out, errUnknown
	}
	var v uint64
	var n int
	if len(b) >= 1 && b[0] < 0x80 {
		v = uint64(b[0])
		n = 1
	} else if len(b) >= 2 && b[1] < 128 {
		v = uint64(b[0]&0x7f) + uint64(b[1])<<7
		n = 2
	} else {
		v, n = protowire.ConsumeVarint(b)
	}
	if n < 0 {
		return protoreflect.Value{}, out, errDecode
	}
	list.Append(protoreflect.ValueOfInt64(int64(v)))
	out.n = n
	return listv, out, nil
}

var coderInt64SliceValue = valueCoderFuncs{
	size:      sizeInt64SliceValue,
	marshal:   appendInt64SliceValue,
	unmarshal: consumeInt64SliceValue,
	merge:     mergeListValue,
}


func sizeInt64PackedSliceValue(listv protoreflect.Value, tagsize int, opts marshalOptions) (size int) {
	list := listv.List()
	llen := list.Len()
	if llen == 0 {
		return 0
	}
	n := 0
	for i, llen := 0, llen; i < llen; i++ {
		v := list.Get(i)
		n += protowire.SizeVarint(uint64(v.Int()))
	}
	return tagsize + protowire.SizeBytes(n)
}


func appendInt64PackedSliceValue(b []byte, listv protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	list := listv.List()
	llen := list.Len()
	if llen == 0 {
		return b, nil
	}
	b = protowire.AppendVarint(b, wiretag)
	n := 0
	for i := 0; i < llen; i++ {
		v := list.Get(i)
		n += protowire.SizeVarint(uint64(v.Int()))
	}
	b = protowire.AppendVarint(b, uint64(n))
	for i := 0; i < llen; i++ {
		v := list.Get(i)
		b = protowire.AppendVarint(b, uint64(v.Int()))
	}
	return b, nil
}

var coderInt64PackedSliceValue = valueCoderFuncs{
	size:      sizeInt64PackedSliceValue,
	marshal:   appendInt64PackedSliceValue,
	unmarshal: consumeInt64SliceValue,
	merge:     mergeListValue,
}


func sizeSint64(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	v := *p.Int64()
	return f.tagsize + protowire.SizeVarint(protowire.EncodeZigZag(v))
}


func appendSint64(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := *p.Int64()
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendVarint(b, protowire.EncodeZigZag(v))
	return b, nil
}


func consumeSint64(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	if wtyp != protowire.VarintType {
		return out, errUnknown
	}
	var v uint64
	var n int
	if len(b) >= 1 && b[0] < 0x80 {
		v = uint64(b[0])
		n = 1
	} else if len(b) >= 2 && b[1] < 128 {
		v = uint64(b[0]&0x7f) + uint64(b[1])<<7
		n = 2
	} else {
		v, n = protowire.ConsumeVarint(b)
	}
	if n < 0 {
		return out, errDecode
	}
	*p.Int64() = protowire.DecodeZigZag(v)
	out.n = n
	return out, nil
}

var coderSint64 = pointerCoderFuncs{
	size:      sizeSint64,
	marshal:   appendSint64,
	unmarshal: consumeSint64,
	merge:     mergeInt64,
}



func sizeSint64NoZero(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	v := *p.Int64()
	if v == 0 {
		return 0
	}
	return f.tagsize + protowire.SizeVarint(protowire.EncodeZigZag(v))
}



func appendSint64NoZero(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := *p.Int64()
	if v == 0 {
		return b, nil
	}
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendVarint(b, protowire.EncodeZigZag(v))
	return b, nil
}

var coderSint64NoZero = pointerCoderFuncs{
	size:      sizeSint64NoZero,
	marshal:   appendSint64NoZero,
	unmarshal: consumeSint64,
	merge:     mergeInt64NoZero,
}



func sizeSint64Ptr(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	v := **p.Int64Ptr()
	return f.tagsize + protowire.SizeVarint(protowire.EncodeZigZag(v))
}



func appendSint64Ptr(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := **p.Int64Ptr()
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendVarint(b, protowire.EncodeZigZag(v))
	return b, nil
}


func consumeSint64Ptr(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	if wtyp != protowire.VarintType {
		return out, errUnknown
	}
	var v uint64
	var n int
	if len(b) >= 1 && b[0] < 0x80 {
		v = uint64(b[0])
		n = 1
	} else if len(b) >= 2 && b[1] < 128 {
		v = uint64(b[0]&0x7f) + uint64(b[1])<<7
		n = 2
	} else {
		v, n = protowire.ConsumeVarint(b)
	}
	if n < 0 {
		return out, errDecode
	}
	vp := p.Int64Ptr()
	if *vp == nil {
		*vp = new(int64)
	}
	**vp = protowire.DecodeZigZag(v)
	out.n = n
	return out, nil
}

var coderSint64Ptr = pointerCoderFuncs{
	size:      sizeSint64Ptr,
	marshal:   appendSint64Ptr,
	unmarshal: consumeSint64Ptr,
	merge:     mergeInt64Ptr,
}


func sizeSint64Slice(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	s := *p.Int64Slice()
	for _, v := range s {
		size += f.tagsize + protowire.SizeVarint(protowire.EncodeZigZag(v))
	}
	return size
}


func appendSint64Slice(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	s := *p.Int64Slice()
	for _, v := range s {
		b = protowire.AppendVarint(b, f.wiretag)
		b = protowire.AppendVarint(b, protowire.EncodeZigZag(v))
	}
	return b, nil
}


func consumeSint64Slice(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	sp := p.Int64Slice()
	if wtyp == protowire.BytesType {
		b, n := protowire.ConsumeBytes(b)
		if n < 0 {
			return out, errDecode
		}
		count := 0
		for _, v := range b {
			if v < 0x80 {
				count++
			}
		}
		if count > 0 {
			p.growInt64Slice(count)
		}
		s := *sp
		for len(b) > 0 {
			var v uint64
			var n int
			if len(b) >= 1 && b[0] < 0x80 {
				v = uint64(b[0])
				n = 1
			} else if len(b) >= 2 && b[1] < 128 {
				v = uint64(b[0]&0x7f) + uint64(b[1])<<7
				n = 2
			} else {
				v, n = protowire.ConsumeVarint(b)
			}
			if n < 0 {
				return out, errDecode
			}
			s = append(s, protowire.DecodeZigZag(v))
			b = b[n:]
		}
		*sp = s
		out.n = n
		return out, nil
	}
	if wtyp != protowire.VarintType {
		return out, errUnknown
	}
	var v uint64
	var n int
	if len(b) >= 1 && b[0] < 0x80 {
		v = uint64(b[0])
		n = 1
	} else if len(b) >= 2 && b[1] < 128 {
		v = uint64(b[0]&0x7f) + uint64(b[1])<<7
		n = 2
	} else {
		v, n = protowire.ConsumeVarint(b)
	}
	if n < 0 {
		return out, errDecode
	}
	*sp = append(*sp, protowire.DecodeZigZag(v))
	out.n = n
	return out, nil
}

var coderSint64Slice = pointerCoderFuncs{
	size:      sizeSint64Slice,
	marshal:   appendSint64Slice,
	unmarshal: consumeSint64Slice,
	merge:     mergeInt64Slice,
}


func sizeSint64PackedSlice(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	s := *p.Int64Slice()
	if len(s) == 0 {
		return 0
	}
	n := 0
	for _, v := range s {
		n += protowire.SizeVarint(protowire.EncodeZigZag(v))
	}
	return f.tagsize + protowire.SizeBytes(n)
}


func appendSint64PackedSlice(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	s := *p.Int64Slice()
	if len(s) == 0 {
		return b, nil
	}
	b = protowire.AppendVarint(b, f.wiretag)
	n := 0
	for _, v := range s {
		n += protowire.SizeVarint(protowire.EncodeZigZag(v))
	}
	b = protowire.AppendVarint(b, uint64(n))
	for _, v := range s {
		b = protowire.AppendVarint(b, protowire.EncodeZigZag(v))
	}
	return b, nil
}

var coderSint64PackedSlice = pointerCoderFuncs{
	size:      sizeSint64PackedSlice,
	marshal:   appendSint64PackedSlice,
	unmarshal: consumeSint64Slice,
	merge:     mergeInt64Slice,
}


func sizeSint64Value(v protoreflect.Value, tagsize int, opts marshalOptions) int {
	return tagsize + protowire.SizeVarint(protowire.EncodeZigZag(v.Int()))
}


func appendSint64Value(b []byte, v protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	b = protowire.AppendVarint(b, wiretag)
	b = protowire.AppendVarint(b, protowire.EncodeZigZag(v.Int()))
	return b, nil
}


func consumeSint64Value(b []byte, _ protoreflect.Value, _ protowire.Number, wtyp protowire.Type, opts unmarshalOptions) (_ protoreflect.Value, out unmarshalOutput, err error) {
	if wtyp != protowire.VarintType {
		return protoreflect.Value{}, out, errUnknown
	}
	var v uint64
	var n int
	if len(b) >= 1 && b[0] < 0x80 {
		v = uint64(b[0])
		n = 1
	} else if len(b) >= 2 && b[1] < 128 {
		v = uint64(b[0]&0x7f) + uint64(b[1])<<7
		n = 2
	} else {
		v, n = protowire.ConsumeVarint(b)
	}
	if n < 0 {
		return protoreflect.Value{}, out, errDecode
	}
	out.n = n
	return protoreflect.ValueOfInt64(protowire.DecodeZigZag(v)), out, nil
}

var coderSint64Value = valueCoderFuncs{
	size:      sizeSint64Value,
	marshal:   appendSint64Value,
	unmarshal: consumeSint64Value,
	merge:     mergeScalarValue,
}


func sizeSint64SliceValue(listv protoreflect.Value, tagsize int, opts marshalOptions) (size int) {
	list := listv.List()
	for i, llen := 0, list.Len(); i < llen; i++ {
		v := list.Get(i)
		size += tagsize + protowire.SizeVarint(protowire.EncodeZigZag(v.Int()))
	}
	return size
}


func appendSint64SliceValue(b []byte, listv protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	list := listv.List()
	for i, llen := 0, list.Len(); i < llen; i++ {
		v := list.Get(i)
		b = protowire.AppendVarint(b, wiretag)
		b = protowire.AppendVarint(b, protowire.EncodeZigZag(v.Int()))
	}
	return b, nil
}


func consumeSint64SliceValue(b []byte, listv protoreflect.Value, _ protowire.Number, wtyp protowire.Type, opts unmarshalOptions) (_ protoreflect.Value, out unmarshalOutput, err error) {
	list := listv.List()
	if wtyp == protowire.BytesType {
		b, n := protowire.ConsumeBytes(b)
		if n < 0 {
			return protoreflect.Value{}, out, errDecode
		}
		for len(b) > 0 {
			var v uint64
			var n int
			if len(b) >= 1 && b[0] < 0x80 {
				v = uint64(b[0])
				n = 1
			} else if len(b) >= 2 && b[1] < 128 {
				v = uint64(b[0]&0x7f) + uint64(b[1])<<7
				n = 2
			} else {
				v, n = protowire.ConsumeVarint(b)
			}
			if n < 0 {
				return protoreflect.Value{}, out, errDecode
			}
			list.Append(protoreflect.ValueOfInt64(protowire.DecodeZigZag(v)))
			b = b[n:]
		}
		out.n = n
		return listv, out, nil
	}
	if wtyp != protowire.VarintType {
		return protoreflect.Value{}, out, errUnknown
	}
	var v uint64
	var n int
	if len(b) >= 1 && b[0] < 0x80 {
		v = uint64(b[0])
		n = 1
	} else if len(b) >= 2 && b[1] < 128 {
		v = uint64(b[0]&0x7f) + uint64(b[1])<<7
		n = 2
	} else {
		v, n = protowire.ConsumeVarint(b)
	}
	if n < 0 {
		return protoreflect.Value{}, out, errDecode
	}
	list.Append(protoreflect.ValueOfInt64(protowire.DecodeZigZag(v)))
	out.n = n
	return listv, out, nil
}

var coderSint64SliceValue = valueCoderFuncs{
	size:      sizeSint64SliceValue,
	marshal:   appendSint64SliceValue,
	unmarshal: consumeSint64SliceValue,
	merge:     mergeListValue,
}


func sizeSint64PackedSliceValue(listv protoreflect.Value, tagsize int, opts marshalOptions) (size int) {
	list := listv.List()
	llen := list.Len()
	if llen == 0 {
		return 0
	}
	n := 0
	for i, llen := 0, llen; i < llen; i++ {
		v := list.Get(i)
		n += protowire.SizeVarint(protowire.EncodeZigZag(v.Int()))
	}
	return tagsize + protowire.SizeBytes(n)
}


func appendSint64PackedSliceValue(b []byte, listv protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	list := listv.List()
	llen := list.Len()
	if llen == 0 {
		return b, nil
	}
	b = protowire.AppendVarint(b, wiretag)
	n := 0
	for i := 0; i < llen; i++ {
		v := list.Get(i)
		n += protowire.SizeVarint(protowire.EncodeZigZag(v.Int()))
	}
	b = protowire.AppendVarint(b, uint64(n))
	for i := 0; i < llen; i++ {
		v := list.Get(i)
		b = protowire.AppendVarint(b, protowire.EncodeZigZag(v.Int()))
	}
	return b, nil
}

var coderSint64PackedSliceValue = valueCoderFuncs{
	size:      sizeSint64PackedSliceValue,
	marshal:   appendSint64PackedSliceValue,
	unmarshal: consumeSint64SliceValue,
	merge:     mergeListValue,
}


func sizeUint64(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	v := *p.Uint64()
	return f.tagsize + protowire.SizeVarint(v)
}


func appendUint64(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := *p.Uint64()
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendVarint(b, v)
	return b, nil
}


func consumeUint64(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	if wtyp != protowire.VarintType {
		return out, errUnknown
	}
	var v uint64
	var n int
	if len(b) >= 1 && b[0] < 0x80 {
		v = uint64(b[0])
		n = 1
	} else if len(b) >= 2 && b[1] < 128 {
		v = uint64(b[0]&0x7f) + uint64(b[1])<<7
		n = 2
	} else {
		v, n = protowire.ConsumeVarint(b)
	}
	if n < 0 {
		return out, errDecode
	}
	*p.Uint64() = v
	out.n = n
	return out, nil
}

var coderUint64 = pointerCoderFuncs{
	size:      sizeUint64,
	marshal:   appendUint64,
	unmarshal: consumeUint64,
	merge:     mergeUint64,
}



func sizeUint64NoZero(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	v := *p.Uint64()
	if v == 0 {
		return 0
	}
	return f.tagsize + protowire.SizeVarint(v)
}



func appendUint64NoZero(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := *p.Uint64()
	if v == 0 {
		return b, nil
	}
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendVarint(b, v)
	return b, nil
}

var coderUint64NoZero = pointerCoderFuncs{
	size:      sizeUint64NoZero,
	marshal:   appendUint64NoZero,
	unmarshal: consumeUint64,
	merge:     mergeUint64NoZero,
}



func sizeUint64Ptr(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	v := **p.Uint64Ptr()
	return f.tagsize + protowire.SizeVarint(v)
}



func appendUint64Ptr(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := **p.Uint64Ptr()
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendVarint(b, v)
	return b, nil
}


func consumeUint64Ptr(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	if wtyp != protowire.VarintType {
		return out, errUnknown
	}
	var v uint64
	var n int
	if len(b) >= 1 && b[0] < 0x80 {
		v = uint64(b[0])
		n = 1
	} else if len(b) >= 2 && b[1] < 128 {
		v = uint64(b[0]&0x7f) + uint64(b[1])<<7
		n = 2
	} else {
		v, n = protowire.ConsumeVarint(b)
	}
	if n < 0 {
		return out, errDecode
	}
	vp := p.Uint64Ptr()
	if *vp == nil {
		*vp = new(uint64)
	}
	**vp = v
	out.n = n
	return out, nil
}

var coderUint64Ptr = pointerCoderFuncs{
	size:      sizeUint64Ptr,
	marshal:   appendUint64Ptr,
	unmarshal: consumeUint64Ptr,
	merge:     mergeUint64Ptr,
}


func sizeUint64Slice(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	s := *p.Uint64Slice()
	for _, v := range s {
		size += f.tagsize + protowire.SizeVarint(v)
	}
	return size
}


func appendUint64Slice(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	s := *p.Uint64Slice()
	for _, v := range s {
		b = protowire.AppendVarint(b, f.wiretag)
		b = protowire.AppendVarint(b, v)
	}
	return b, nil
}


func consumeUint64Slice(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	sp := p.Uint64Slice()
	if wtyp == protowire.BytesType {
		b, n := protowire.ConsumeBytes(b)
		if n < 0 {
			return out, errDecode
		}
		count := 0
		for _, v := range b {
			if v < 0x80 {
				count++
			}
		}
		if count > 0 {
			p.growUint64Slice(count)
		}
		s := *sp
		for len(b) > 0 {
			var v uint64
			var n int
			if len(b) >= 1 && b[0] < 0x80 {
				v = uint64(b[0])
				n = 1
			} else if len(b) >= 2 && b[1] < 128 {
				v = uint64(b[0]&0x7f) + uint64(b[1])<<7
				n = 2
			} else {
				v, n = protowire.ConsumeVarint(b)
			}
			if n < 0 {
				return out, errDecode
			}
			s = append(s, v)
			b = b[n:]
		}
		*sp = s
		out.n = n
		return out, nil
	}
	if wtyp != protowire.VarintType {
		return out, errUnknown
	}
	var v uint64
	var n int
	if len(b) >= 1 && b[0] < 0x80 {
		v = uint64(b[0])
		n = 1
	} else if len(b) >= 2 && b[1] < 128 {
		v = uint64(b[0]&0x7f) + uint64(b[1])<<7
		n = 2
	} else {
		v, n = protowire.ConsumeVarint(b)
	}
	if n < 0 {
		return out, errDecode
	}
	*sp = append(*sp, v)
	out.n = n
	return out, nil
}

var coderUint64Slice = pointerCoderFuncs{
	size:      sizeUint64Slice,
	marshal:   appendUint64Slice,
	unmarshal: consumeUint64Slice,
	merge:     mergeUint64Slice,
}


func sizeUint64PackedSlice(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	s := *p.Uint64Slice()
	if len(s) == 0 {
		return 0
	}
	n := 0
	for _, v := range s {
		n += protowire.SizeVarint(v)
	}
	return f.tagsize + protowire.SizeBytes(n)
}


func appendUint64PackedSlice(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	s := *p.Uint64Slice()
	if len(s) == 0 {
		return b, nil
	}
	b = protowire.AppendVarint(b, f.wiretag)
	n := 0
	for _, v := range s {
		n += protowire.SizeVarint(v)
	}
	b = protowire.AppendVarint(b, uint64(n))
	for _, v := range s {
		b = protowire.AppendVarint(b, v)
	}
	return b, nil
}

var coderUint64PackedSlice = pointerCoderFuncs{
	size:      sizeUint64PackedSlice,
	marshal:   appendUint64PackedSlice,
	unmarshal: consumeUint64Slice,
	merge:     mergeUint64Slice,
}


func sizeUint64Value(v protoreflect.Value, tagsize int, opts marshalOptions) int {
	return tagsize + protowire.SizeVarint(v.Uint())
}


func appendUint64Value(b []byte, v protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	b = protowire.AppendVarint(b, wiretag)
	b = protowire.AppendVarint(b, v.Uint())
	return b, nil
}


func consumeUint64Value(b []byte, _ protoreflect.Value, _ protowire.Number, wtyp protowire.Type, opts unmarshalOptions) (_ protoreflect.Value, out unmarshalOutput, err error) {
	if wtyp != protowire.VarintType {
		return protoreflect.Value{}, out, errUnknown
	}
	var v uint64
	var n int
	if len(b) >= 1 && b[0] < 0x80 {
		v = uint64(b[0])
		n = 1
	} else if len(b) >= 2 && b[1] < 128 {
		v = uint64(b[0]&0x7f) + uint64(b[1])<<7
		n = 2
	} else {
		v, n = protowire.ConsumeVarint(b)
	}
	if n < 0 {
		return protoreflect.Value{}, out, errDecode
	}
	out.n = n
	return protoreflect.ValueOfUint64(v), out, nil
}

var coderUint64Value = valueCoderFuncs{
	size:      sizeUint64Value,
	marshal:   appendUint64Value,
	unmarshal: consumeUint64Value,
	merge:     mergeScalarValue,
}


func sizeUint64SliceValue(listv protoreflect.Value, tagsize int, opts marshalOptions) (size int) {
	list := listv.List()
	for i, llen := 0, list.Len(); i < llen; i++ {
		v := list.Get(i)
		size += tagsize + protowire.SizeVarint(v.Uint())
	}
	return size
}


func appendUint64SliceValue(b []byte, listv protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	list := listv.List()
	for i, llen := 0, list.Len(); i < llen; i++ {
		v := list.Get(i)
		b = protowire.AppendVarint(b, wiretag)
		b = protowire.AppendVarint(b, v.Uint())
	}
	return b, nil
}


func consumeUint64SliceValue(b []byte, listv protoreflect.Value, _ protowire.Number, wtyp protowire.Type, opts unmarshalOptions) (_ protoreflect.Value, out unmarshalOutput, err error) {
	list := listv.List()
	if wtyp == protowire.BytesType {
		b, n := protowire.ConsumeBytes(b)
		if n < 0 {
			return protoreflect.Value{}, out, errDecode
		}
		for len(b) > 0 {
			var v uint64
			var n int
			if len(b) >= 1 && b[0] < 0x80 {
				v = uint64(b[0])
				n = 1
			} else if len(b) >= 2 && b[1] < 128 {
				v = uint64(b[0]&0x7f) + uint64(b[1])<<7
				n = 2
			} else {
				v, n = protowire.ConsumeVarint(b)
			}
			if n < 0 {
				return protoreflect.Value{}, out, errDecode
			}
			list.Append(protoreflect.ValueOfUint64(v))
			b = b[n:]
		}
		out.n = n
		return listv, out, nil
	}
	if wtyp != protowire.VarintType {
		return protoreflect.Value{}, out, errUnknown
	}
	var v uint64
	var n int
	if len(b) >= 1 && b[0] < 0x80 {
		v = uint64(b[0])
		n = 1
	} else if len(b) >= 2 && b[1] < 128 {
		v = uint64(b[0]&0x7f) + uint64(b[1])<<7
		n = 2
	} else {
		v, n = protowire.ConsumeVarint(b)
	}
	if n < 0 {
		return protoreflect.Value{}, out, errDecode
	}
	list.Append(protoreflect.ValueOfUint64(v))
	out.n = n
	return listv, out, nil
}

var coderUint64SliceValue = valueCoderFuncs{
	size:      sizeUint64SliceValue,
	marshal:   appendUint64SliceValue,
	unmarshal: consumeUint64SliceValue,
	merge:     mergeListValue,
}


func sizeUint64PackedSliceValue(listv protoreflect.Value, tagsize int, opts marshalOptions) (size int) {
	list := listv.List()
	llen := list.Len()
	if llen == 0 {
		return 0
	}
	n := 0
	for i, llen := 0, llen; i < llen; i++ {
		v := list.Get(i)
		n += protowire.SizeVarint(v.Uint())
	}
	return tagsize + protowire.SizeBytes(n)
}


func appendUint64PackedSliceValue(b []byte, listv protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	list := listv.List()
	llen := list.Len()
	if llen == 0 {
		return b, nil
	}
	b = protowire.AppendVarint(b, wiretag)
	n := 0
	for i := 0; i < llen; i++ {
		v := list.Get(i)
		n += protowire.SizeVarint(v.Uint())
	}
	b = protowire.AppendVarint(b, uint64(n))
	for i := 0; i < llen; i++ {
		v := list.Get(i)
		b = protowire.AppendVarint(b, v.Uint())
	}
	return b, nil
}

var coderUint64PackedSliceValue = valueCoderFuncs{
	size:      sizeUint64PackedSliceValue,
	marshal:   appendUint64PackedSliceValue,
	unmarshal: consumeUint64SliceValue,
	merge:     mergeListValue,
}


func sizeSfixed32(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {

	return f.tagsize + protowire.SizeFixed32()
}


func appendSfixed32(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := *p.Int32()
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendFixed32(b, uint32(v))
	return b, nil
}


func consumeSfixed32(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	if wtyp != protowire.Fixed32Type {
		return out, errUnknown
	}
	v, n := protowire.ConsumeFixed32(b)
	if n < 0 {
		return out, errDecode
	}
	*p.Int32() = int32(v)
	out.n = n
	return out, nil
}

var coderSfixed32 = pointerCoderFuncs{
	size:      sizeSfixed32,
	marshal:   appendSfixed32,
	unmarshal: consumeSfixed32,
	merge:     mergeInt32,
}



func sizeSfixed32NoZero(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	v := *p.Int32()
	if v == 0 {
		return 0
	}
	return f.tagsize + protowire.SizeFixed32()
}



func appendSfixed32NoZero(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := *p.Int32()
	if v == 0 {
		return b, nil
	}
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendFixed32(b, uint32(v))
	return b, nil
}

var coderSfixed32NoZero = pointerCoderFuncs{
	size:      sizeSfixed32NoZero,
	marshal:   appendSfixed32NoZero,
	unmarshal: consumeSfixed32,
	merge:     mergeInt32NoZero,
}



func sizeSfixed32Ptr(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	return f.tagsize + protowire.SizeFixed32()
}



func appendSfixed32Ptr(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := **p.Int32Ptr()
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendFixed32(b, uint32(v))
	return b, nil
}


func consumeSfixed32Ptr(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	if wtyp != protowire.Fixed32Type {
		return out, errUnknown
	}
	v, n := protowire.ConsumeFixed32(b)
	if n < 0 {
		return out, errDecode
	}
	vp := p.Int32Ptr()
	if *vp == nil {
		*vp = new(int32)
	}
	**vp = int32(v)
	out.n = n
	return out, nil
}

var coderSfixed32Ptr = pointerCoderFuncs{
	size:      sizeSfixed32Ptr,
	marshal:   appendSfixed32Ptr,
	unmarshal: consumeSfixed32Ptr,
	merge:     mergeInt32Ptr,
}


func sizeSfixed32Slice(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	s := *p.Int32Slice()
	size = len(s) * (f.tagsize + protowire.SizeFixed32())
	return size
}


func appendSfixed32Slice(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	s := *p.Int32Slice()
	for _, v := range s {
		b = protowire.AppendVarint(b, f.wiretag)
		b = protowire.AppendFixed32(b, uint32(v))
	}
	return b, nil
}


func consumeSfixed32Slice(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	sp := p.Int32Slice()
	if wtyp == protowire.BytesType {
		b, n := protowire.ConsumeBytes(b)
		if n < 0 {
			return out, errDecode
		}
		count := len(b) / protowire.SizeFixed32()
		if count > 0 {
			p.growInt32Slice(count)
		}
		s := *sp
		for len(b) > 0 {
			v, n := protowire.ConsumeFixed32(b)
			if n < 0 {
				return out, errDecode
			}
			s = append(s, int32(v))
			b = b[n:]
		}
		*sp = s
		out.n = n
		return out, nil
	}
	if wtyp != protowire.Fixed32Type {
		return out, errUnknown
	}
	v, n := protowire.ConsumeFixed32(b)
	if n < 0 {
		return out, errDecode
	}
	*sp = append(*sp, int32(v))
	out.n = n
	return out, nil
}

var coderSfixed32Slice = pointerCoderFuncs{
	size:      sizeSfixed32Slice,
	marshal:   appendSfixed32Slice,
	unmarshal: consumeSfixed32Slice,
	merge:     mergeInt32Slice,
}


func sizeSfixed32PackedSlice(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	s := *p.Int32Slice()
	if len(s) == 0 {
		return 0
	}
	n := len(s) * protowire.SizeFixed32()
	return f.tagsize + protowire.SizeBytes(n)
}


func appendSfixed32PackedSlice(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	s := *p.Int32Slice()
	if len(s) == 0 {
		return b, nil
	}
	b = protowire.AppendVarint(b, f.wiretag)
	n := len(s) * protowire.SizeFixed32()
	b = protowire.AppendVarint(b, uint64(n))
	for _, v := range s {
		b = protowire.AppendFixed32(b, uint32(v))
	}
	return b, nil
}

var coderSfixed32PackedSlice = pointerCoderFuncs{
	size:      sizeSfixed32PackedSlice,
	marshal:   appendSfixed32PackedSlice,
	unmarshal: consumeSfixed32Slice,
	merge:     mergeInt32Slice,
}


func sizeSfixed32Value(v protoreflect.Value, tagsize int, opts marshalOptions) int {
	return tagsize + protowire.SizeFixed32()
}


func appendSfixed32Value(b []byte, v protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	b = protowire.AppendVarint(b, wiretag)
	b = protowire.AppendFixed32(b, uint32(v.Int()))
	return b, nil
}


func consumeSfixed32Value(b []byte, _ protoreflect.Value, _ protowire.Number, wtyp protowire.Type, opts unmarshalOptions) (_ protoreflect.Value, out unmarshalOutput, err error) {
	if wtyp != protowire.Fixed32Type {
		return protoreflect.Value{}, out, errUnknown
	}
	v, n := protowire.ConsumeFixed32(b)
	if n < 0 {
		return protoreflect.Value{}, out, errDecode
	}
	out.n = n
	return protoreflect.ValueOfInt32(int32(v)), out, nil
}

var coderSfixed32Value = valueCoderFuncs{
	size:      sizeSfixed32Value,
	marshal:   appendSfixed32Value,
	unmarshal: consumeSfixed32Value,
	merge:     mergeScalarValue,
}


func sizeSfixed32SliceValue(listv protoreflect.Value, tagsize int, opts marshalOptions) (size int) {
	list := listv.List()
	size = list.Len() * (tagsize + protowire.SizeFixed32())
	return size
}


func appendSfixed32SliceValue(b []byte, listv protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	list := listv.List()
	for i, llen := 0, list.Len(); i < llen; i++ {
		v := list.Get(i)
		b = protowire.AppendVarint(b, wiretag)
		b = protowire.AppendFixed32(b, uint32(v.Int()))
	}
	return b, nil
}


func consumeSfixed32SliceValue(b []byte, listv protoreflect.Value, _ protowire.Number, wtyp protowire.Type, opts unmarshalOptions) (_ protoreflect.Value, out unmarshalOutput, err error) {
	list := listv.List()
	if wtyp == protowire.BytesType {
		b, n := protowire.ConsumeBytes(b)
		if n < 0 {
			return protoreflect.Value{}, out, errDecode
		}
		for len(b) > 0 {
			v, n := protowire.ConsumeFixed32(b)
			if n < 0 {
				return protoreflect.Value{}, out, errDecode
			}
			list.Append(protoreflect.ValueOfInt32(int32(v)))
			b = b[n:]
		}
		out.n = n
		return listv, out, nil
	}
	if wtyp != protowire.Fixed32Type {
		return protoreflect.Value{}, out, errUnknown
	}
	v, n := protowire.ConsumeFixed32(b)
	if n < 0 {
		return protoreflect.Value{}, out, errDecode
	}
	list.Append(protoreflect.ValueOfInt32(int32(v)))
	out.n = n
	return listv, out, nil
}

var coderSfixed32SliceValue = valueCoderFuncs{
	size:      sizeSfixed32SliceValue,
	marshal:   appendSfixed32SliceValue,
	unmarshal: consumeSfixed32SliceValue,
	merge:     mergeListValue,
}


func sizeSfixed32PackedSliceValue(listv protoreflect.Value, tagsize int, opts marshalOptions) (size int) {
	list := listv.List()
	llen := list.Len()
	if llen == 0 {
		return 0
	}
	n := llen * protowire.SizeFixed32()
	return tagsize + protowire.SizeBytes(n)
}


func appendSfixed32PackedSliceValue(b []byte, listv protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	list := listv.List()
	llen := list.Len()
	if llen == 0 {
		return b, nil
	}
	b = protowire.AppendVarint(b, wiretag)
	n := llen * protowire.SizeFixed32()
	b = protowire.AppendVarint(b, uint64(n))
	for i := 0; i < llen; i++ {
		v := list.Get(i)
		b = protowire.AppendFixed32(b, uint32(v.Int()))
	}
	return b, nil
}

var coderSfixed32PackedSliceValue = valueCoderFuncs{
	size:      sizeSfixed32PackedSliceValue,
	marshal:   appendSfixed32PackedSliceValue,
	unmarshal: consumeSfixed32SliceValue,
	merge:     mergeListValue,
}


func sizeFixed32(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {

	return f.tagsize + protowire.SizeFixed32()
}


func appendFixed32(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := *p.Uint32()
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendFixed32(b, v)
	return b, nil
}


func consumeFixed32(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	if wtyp != protowire.Fixed32Type {
		return out, errUnknown
	}
	v, n := protowire.ConsumeFixed32(b)
	if n < 0 {
		return out, errDecode
	}
	*p.Uint32() = v
	out.n = n
	return out, nil
}

var coderFixed32 = pointerCoderFuncs{
	size:      sizeFixed32,
	marshal:   appendFixed32,
	unmarshal: consumeFixed32,
	merge:     mergeUint32,
}



func sizeFixed32NoZero(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	v := *p.Uint32()
	if v == 0 {
		return 0
	}
	return f.tagsize + protowire.SizeFixed32()
}



func appendFixed32NoZero(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := *p.Uint32()
	if v == 0 {
		return b, nil
	}
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendFixed32(b, v)
	return b, nil
}

var coderFixed32NoZero = pointerCoderFuncs{
	size:      sizeFixed32NoZero,
	marshal:   appendFixed32NoZero,
	unmarshal: consumeFixed32,
	merge:     mergeUint32NoZero,
}



func sizeFixed32Ptr(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	return f.tagsize + protowire.SizeFixed32()
}



func appendFixed32Ptr(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := **p.Uint32Ptr()
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendFixed32(b, v)
	return b, nil
}


func consumeFixed32Ptr(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	if wtyp != protowire.Fixed32Type {
		return out, errUnknown
	}
	v, n := protowire.ConsumeFixed32(b)
	if n < 0 {
		return out, errDecode
	}
	vp := p.Uint32Ptr()
	if *vp == nil {
		*vp = new(uint32)
	}
	**vp = v
	out.n = n
	return out, nil
}

var coderFixed32Ptr = pointerCoderFuncs{
	size:      sizeFixed32Ptr,
	marshal:   appendFixed32Ptr,
	unmarshal: consumeFixed32Ptr,
	merge:     mergeUint32Ptr,
}


func sizeFixed32Slice(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	s := *p.Uint32Slice()
	size = len(s) * (f.tagsize + protowire.SizeFixed32())
	return size
}


func appendFixed32Slice(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	s := *p.Uint32Slice()
	for _, v := range s {
		b = protowire.AppendVarint(b, f.wiretag)
		b = protowire.AppendFixed32(b, v)
	}
	return b, nil
}


func consumeFixed32Slice(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	sp := p.Uint32Slice()
	if wtyp == protowire.BytesType {
		b, n := protowire.ConsumeBytes(b)
		if n < 0 {
			return out, errDecode
		}
		count := len(b) / protowire.SizeFixed32()
		if count > 0 {
			p.growUint32Slice(count)
		}
		s := *sp
		for len(b) > 0 {
			v, n := protowire.ConsumeFixed32(b)
			if n < 0 {
				return out, errDecode
			}
			s = append(s, v)
			b = b[n:]
		}
		*sp = s
		out.n = n
		return out, nil
	}
	if wtyp != protowire.Fixed32Type {
		return out, errUnknown
	}
	v, n := protowire.ConsumeFixed32(b)
	if n < 0 {
		return out, errDecode
	}
	*sp = append(*sp, v)
	out.n = n
	return out, nil
}

var coderFixed32Slice = pointerCoderFuncs{
	size:      sizeFixed32Slice,
	marshal:   appendFixed32Slice,
	unmarshal: consumeFixed32Slice,
	merge:     mergeUint32Slice,
}


func sizeFixed32PackedSlice(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	s := *p.Uint32Slice()
	if len(s) == 0 {
		return 0
	}
	n := len(s) * protowire.SizeFixed32()
	return f.tagsize + protowire.SizeBytes(n)
}


func appendFixed32PackedSlice(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	s := *p.Uint32Slice()
	if len(s) == 0 {
		return b, nil
	}
	b = protowire.AppendVarint(b, f.wiretag)
	n := len(s) * protowire.SizeFixed32()
	b = protowire.AppendVarint(b, uint64(n))
	for _, v := range s {
		b = protowire.AppendFixed32(b, v)
	}
	return b, nil
}

var coderFixed32PackedSlice = pointerCoderFuncs{
	size:      sizeFixed32PackedSlice,
	marshal:   appendFixed32PackedSlice,
	unmarshal: consumeFixed32Slice,
	merge:     mergeUint32Slice,
}


func sizeFixed32Value(v protoreflect.Value, tagsize int, opts marshalOptions) int {
	return tagsize + protowire.SizeFixed32()
}


func appendFixed32Value(b []byte, v protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	b = protowire.AppendVarint(b, wiretag)
	b = protowire.AppendFixed32(b, uint32(v.Uint()))
	return b, nil
}


func consumeFixed32Value(b []byte, _ protoreflect.Value, _ protowire.Number, wtyp protowire.Type, opts unmarshalOptions) (_ protoreflect.Value, out unmarshalOutput, err error) {
	if wtyp != protowire.Fixed32Type {
		return protoreflect.Value{}, out, errUnknown
	}
	v, n := protowire.ConsumeFixed32(b)
	if n < 0 {
		return protoreflect.Value{}, out, errDecode
	}
	out.n = n
	return protoreflect.ValueOfUint32(uint32(v)), out, nil
}

var coderFixed32Value = valueCoderFuncs{
	size:      sizeFixed32Value,
	marshal:   appendFixed32Value,
	unmarshal: consumeFixed32Value,
	merge:     mergeScalarValue,
}


func sizeFixed32SliceValue(listv protoreflect.Value, tagsize int, opts marshalOptions) (size int) {
	list := listv.List()
	size = list.Len() * (tagsize + protowire.SizeFixed32())
	return size
}


func appendFixed32SliceValue(b []byte, listv protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	list := listv.List()
	for i, llen := 0, list.Len(); i < llen; i++ {
		v := list.Get(i)
		b = protowire.AppendVarint(b, wiretag)
		b = protowire.AppendFixed32(b, uint32(v.Uint()))
	}
	return b, nil
}


func consumeFixed32SliceValue(b []byte, listv protoreflect.Value, _ protowire.Number, wtyp protowire.Type, opts unmarshalOptions) (_ protoreflect.Value, out unmarshalOutput, err error) {
	list := listv.List()
	if wtyp == protowire.BytesType {
		b, n := protowire.ConsumeBytes(b)
		if n < 0 {
			return protoreflect.Value{}, out, errDecode
		}
		for len(b) > 0 {
			v, n := protowire.ConsumeFixed32(b)
			if n < 0 {
				return protoreflect.Value{}, out, errDecode
			}
			list.Append(protoreflect.ValueOfUint32(uint32(v)))
			b = b[n:]
		}
		out.n = n
		return listv, out, nil
	}
	if wtyp != protowire.Fixed32Type {
		return protoreflect.Value{}, out, errUnknown
	}
	v, n := protowire.ConsumeFixed32(b)
	if n < 0 {
		return protoreflect.Value{}, out, errDecode
	}
	list.Append(protoreflect.ValueOfUint32(uint32(v)))
	out.n = n
	return listv, out, nil
}

var coderFixed32SliceValue = valueCoderFuncs{
	size:      sizeFixed32SliceValue,
	marshal:   appendFixed32SliceValue,
	unmarshal: consumeFixed32SliceValue,
	merge:     mergeListValue,
}


func sizeFixed32PackedSliceValue(listv protoreflect.Value, tagsize int, opts marshalOptions) (size int) {
	list := listv.List()
	llen := list.Len()
	if llen == 0 {
		return 0
	}
	n := llen * protowire.SizeFixed32()
	return tagsize + protowire.SizeBytes(n)
}


func appendFixed32PackedSliceValue(b []byte, listv protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	list := listv.List()
	llen := list.Len()
	if llen == 0 {
		return b, nil
	}
	b = protowire.AppendVarint(b, wiretag)
	n := llen * protowire.SizeFixed32()
	b = protowire.AppendVarint(b, uint64(n))
	for i := 0; i < llen; i++ {
		v := list.Get(i)
		b = protowire.AppendFixed32(b, uint32(v.Uint()))
	}
	return b, nil
}

var coderFixed32PackedSliceValue = valueCoderFuncs{
	size:      sizeFixed32PackedSliceValue,
	marshal:   appendFixed32PackedSliceValue,
	unmarshal: consumeFixed32SliceValue,
	merge:     mergeListValue,
}


func sizeFloat(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {

	return f.tagsize + protowire.SizeFixed32()
}


func appendFloat(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := *p.Float32()
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendFixed32(b, math.Float32bits(v))
	return b, nil
}


func consumeFloat(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	if wtyp != protowire.Fixed32Type {
		return out, errUnknown
	}
	v, n := protowire.ConsumeFixed32(b)
	if n < 0 {
		return out, errDecode
	}
	*p.Float32() = math.Float32frombits(v)
	out.n = n
	return out, nil
}

var coderFloat = pointerCoderFuncs{
	size:      sizeFloat,
	marshal:   appendFloat,
	unmarshal: consumeFloat,
	merge:     mergeFloat32,
}



func sizeFloatNoZero(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	v := *p.Float32()
	if v == 0 && !math.Signbit(float64(v)) {
		return 0
	}
	return f.tagsize + protowire.SizeFixed32()
}



func appendFloatNoZero(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := *p.Float32()
	if v == 0 && !math.Signbit(float64(v)) {
		return b, nil
	}
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendFixed32(b, math.Float32bits(v))
	return b, nil
}

var coderFloatNoZero = pointerCoderFuncs{
	size:      sizeFloatNoZero,
	marshal:   appendFloatNoZero,
	unmarshal: consumeFloat,
	merge:     mergeFloat32NoZero,
}



func sizeFloatPtr(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	return f.tagsize + protowire.SizeFixed32()
}



func appendFloatPtr(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := **p.Float32Ptr()
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendFixed32(b, math.Float32bits(v))
	return b, nil
}


func consumeFloatPtr(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	if wtyp != protowire.Fixed32Type {
		return out, errUnknown
	}
	v, n := protowire.ConsumeFixed32(b)
	if n < 0 {
		return out, errDecode
	}
	vp := p.Float32Ptr()
	if *vp == nil {
		*vp = new(float32)
	}
	**vp = math.Float32frombits(v)
	out.n = n
	return out, nil
}

var coderFloatPtr = pointerCoderFuncs{
	size:      sizeFloatPtr,
	marshal:   appendFloatPtr,
	unmarshal: consumeFloatPtr,
	merge:     mergeFloat32Ptr,
}


func sizeFloatSlice(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	s := *p.Float32Slice()
	size = len(s) * (f.tagsize + protowire.SizeFixed32())
	return size
}


func appendFloatSlice(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	s := *p.Float32Slice()
	for _, v := range s {
		b = protowire.AppendVarint(b, f.wiretag)
		b = protowire.AppendFixed32(b, math.Float32bits(v))
	}
	return b, nil
}


func consumeFloatSlice(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	sp := p.Float32Slice()
	if wtyp == protowire.BytesType {
		b, n := protowire.ConsumeBytes(b)
		if n < 0 {
			return out, errDecode
		}
		count := len(b) / protowire.SizeFixed32()
		if count > 0 {
			p.growFloat32Slice(count)
		}
		s := *sp
		for len(b) > 0 {
			v, n := protowire.ConsumeFixed32(b)
			if n < 0 {
				return out, errDecode
			}
			s = append(s, math.Float32frombits(v))
			b = b[n:]
		}
		*sp = s
		out.n = n
		return out, nil
	}
	if wtyp != protowire.Fixed32Type {
		return out, errUnknown
	}
	v, n := protowire.ConsumeFixed32(b)
	if n < 0 {
		return out, errDecode
	}
	*sp = append(*sp, math.Float32frombits(v))
	out.n = n
	return out, nil
}

var coderFloatSlice = pointerCoderFuncs{
	size:      sizeFloatSlice,
	marshal:   appendFloatSlice,
	unmarshal: consumeFloatSlice,
	merge:     mergeFloat32Slice,
}


func sizeFloatPackedSlice(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	s := *p.Float32Slice()
	if len(s) == 0 {
		return 0
	}
	n := len(s) * protowire.SizeFixed32()
	return f.tagsize + protowire.SizeBytes(n)
}


func appendFloatPackedSlice(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	s := *p.Float32Slice()
	if len(s) == 0 {
		return b, nil
	}
	b = protowire.AppendVarint(b, f.wiretag)
	n := len(s) * protowire.SizeFixed32()
	b = protowire.AppendVarint(b, uint64(n))
	for _, v := range s {
		b = protowire.AppendFixed32(b, math.Float32bits(v))
	}
	return b, nil
}

var coderFloatPackedSlice = pointerCoderFuncs{
	size:      sizeFloatPackedSlice,
	marshal:   appendFloatPackedSlice,
	unmarshal: consumeFloatSlice,
	merge:     mergeFloat32Slice,
}


func sizeFloatValue(v protoreflect.Value, tagsize int, opts marshalOptions) int {
	return tagsize + protowire.SizeFixed32()
}


func appendFloatValue(b []byte, v protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	b = protowire.AppendVarint(b, wiretag)
	b = protowire.AppendFixed32(b, math.Float32bits(float32(v.Float())))
	return b, nil
}


func consumeFloatValue(b []byte, _ protoreflect.Value, _ protowire.Number, wtyp protowire.Type, opts unmarshalOptions) (_ protoreflect.Value, out unmarshalOutput, err error) {
	if wtyp != protowire.Fixed32Type {
		return protoreflect.Value{}, out, errUnknown
	}
	v, n := protowire.ConsumeFixed32(b)
	if n < 0 {
		return protoreflect.Value{}, out, errDecode
	}
	out.n = n
	return protoreflect.ValueOfFloat32(math.Float32frombits(uint32(v))), out, nil
}

var coderFloatValue = valueCoderFuncs{
	size:      sizeFloatValue,
	marshal:   appendFloatValue,
	unmarshal: consumeFloatValue,
	merge:     mergeScalarValue,
}


func sizeFloatSliceValue(listv protoreflect.Value, tagsize int, opts marshalOptions) (size int) {
	list := listv.List()
	size = list.Len() * (tagsize + protowire.SizeFixed32())
	return size
}


func appendFloatSliceValue(b []byte, listv protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	list := listv.List()
	for i, llen := 0, list.Len(); i < llen; i++ {
		v := list.Get(i)
		b = protowire.AppendVarint(b, wiretag)
		b = protowire.AppendFixed32(b, math.Float32bits(float32(v.Float())))
	}
	return b, nil
}


func consumeFloatSliceValue(b []byte, listv protoreflect.Value, _ protowire.Number, wtyp protowire.Type, opts unmarshalOptions) (_ protoreflect.Value, out unmarshalOutput, err error) {
	list := listv.List()
	if wtyp == protowire.BytesType {
		b, n := protowire.ConsumeBytes(b)
		if n < 0 {
			return protoreflect.Value{}, out, errDecode
		}
		for len(b) > 0 {
			v, n := protowire.ConsumeFixed32(b)
			if n < 0 {
				return protoreflect.Value{}, out, errDecode
			}
			list.Append(protoreflect.ValueOfFloat32(math.Float32frombits(uint32(v))))
			b = b[n:]
		}
		out.n = n
		return listv, out, nil
	}
	if wtyp != protowire.Fixed32Type {
		return protoreflect.Value{}, out, errUnknown
	}
	v, n := protowire.ConsumeFixed32(b)
	if n < 0 {
		return protoreflect.Value{}, out, errDecode
	}
	list.Append(protoreflect.ValueOfFloat32(math.Float32frombits(uint32(v))))
	out.n = n
	return listv, out, nil
}

var coderFloatSliceValue = valueCoderFuncs{
	size:      sizeFloatSliceValue,
	marshal:   appendFloatSliceValue,
	unmarshal: consumeFloatSliceValue,
	merge:     mergeListValue,
}


func sizeFloatPackedSliceValue(listv protoreflect.Value, tagsize int, opts marshalOptions) (size int) {
	list := listv.List()
	llen := list.Len()
	if llen == 0 {
		return 0
	}
	n := llen * protowire.SizeFixed32()
	return tagsize + protowire.SizeBytes(n)
}


func appendFloatPackedSliceValue(b []byte, listv protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	list := listv.List()
	llen := list.Len()
	if llen == 0 {
		return b, nil
	}
	b = protowire.AppendVarint(b, wiretag)
	n := llen * protowire.SizeFixed32()
	b = protowire.AppendVarint(b, uint64(n))
	for i := 0; i < llen; i++ {
		v := list.Get(i)
		b = protowire.AppendFixed32(b, math.Float32bits(float32(v.Float())))
	}
	return b, nil
}

var coderFloatPackedSliceValue = valueCoderFuncs{
	size:      sizeFloatPackedSliceValue,
	marshal:   appendFloatPackedSliceValue,
	unmarshal: consumeFloatSliceValue,
	merge:     mergeListValue,
}


func sizeSfixed64(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {

	return f.tagsize + protowire.SizeFixed64()
}


func appendSfixed64(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := *p.Int64()
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendFixed64(b, uint64(v))
	return b, nil
}


func consumeSfixed64(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	if wtyp != protowire.Fixed64Type {
		return out, errUnknown
	}
	v, n := protowire.ConsumeFixed64(b)
	if n < 0 {
		return out, errDecode
	}
	*p.Int64() = int64(v)
	out.n = n
	return out, nil
}

var coderSfixed64 = pointerCoderFuncs{
	size:      sizeSfixed64,
	marshal:   appendSfixed64,
	unmarshal: consumeSfixed64,
	merge:     mergeInt64,
}



func sizeSfixed64NoZero(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	v := *p.Int64()
	if v == 0 {
		return 0
	}
	return f.tagsize + protowire.SizeFixed64()
}



func appendSfixed64NoZero(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := *p.Int64()
	if v == 0 {
		return b, nil
	}
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendFixed64(b, uint64(v))
	return b, nil
}

var coderSfixed64NoZero = pointerCoderFuncs{
	size:      sizeSfixed64NoZero,
	marshal:   appendSfixed64NoZero,
	unmarshal: consumeSfixed64,
	merge:     mergeInt64NoZero,
}



func sizeSfixed64Ptr(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	return f.tagsize + protowire.SizeFixed64()
}



func appendSfixed64Ptr(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := **p.Int64Ptr()
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendFixed64(b, uint64(v))
	return b, nil
}


func consumeSfixed64Ptr(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	if wtyp != protowire.Fixed64Type {
		return out, errUnknown
	}
	v, n := protowire.ConsumeFixed64(b)
	if n < 0 {
		return out, errDecode
	}
	vp := p.Int64Ptr()
	if *vp == nil {
		*vp = new(int64)
	}
	**vp = int64(v)
	out.n = n
	return out, nil
}

var coderSfixed64Ptr = pointerCoderFuncs{
	size:      sizeSfixed64Ptr,
	marshal:   appendSfixed64Ptr,
	unmarshal: consumeSfixed64Ptr,
	merge:     mergeInt64Ptr,
}


func sizeSfixed64Slice(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	s := *p.Int64Slice()
	size = len(s) * (f.tagsize + protowire.SizeFixed64())
	return size
}


func appendSfixed64Slice(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	s := *p.Int64Slice()
	for _, v := range s {
		b = protowire.AppendVarint(b, f.wiretag)
		b = protowire.AppendFixed64(b, uint64(v))
	}
	return b, nil
}


func consumeSfixed64Slice(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	sp := p.Int64Slice()
	if wtyp == protowire.BytesType {
		b, n := protowire.ConsumeBytes(b)
		if n < 0 {
			return out, errDecode
		}
		count := len(b) / protowire.SizeFixed64()
		if count > 0 {
			p.growInt64Slice(count)
		}
		s := *sp
		for len(b) > 0 {
			v, n := protowire.ConsumeFixed64(b)
			if n < 0 {
				return out, errDecode
			}
			s = append(s, int64(v))
			b = b[n:]
		}
		*sp = s
		out.n = n
		return out, nil
	}
	if wtyp != protowire.Fixed64Type {
		return out, errUnknown
	}
	v, n := protowire.ConsumeFixed64(b)
	if n < 0 {
		return out, errDecode
	}
	*sp = append(*sp, int64(v))
	out.n = n
	return out, nil
}

var coderSfixed64Slice = pointerCoderFuncs{
	size:      sizeSfixed64Slice,
	marshal:   appendSfixed64Slice,
	unmarshal: consumeSfixed64Slice,
	merge:     mergeInt64Slice,
}


func sizeSfixed64PackedSlice(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	s := *p.Int64Slice()
	if len(s) == 0 {
		return 0
	}
	n := len(s) * protowire.SizeFixed64()
	return f.tagsize + protowire.SizeBytes(n)
}


func appendSfixed64PackedSlice(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	s := *p.Int64Slice()
	if len(s) == 0 {
		return b, nil
	}
	b = protowire.AppendVarint(b, f.wiretag)
	n := len(s) * protowire.SizeFixed64()
	b = protowire.AppendVarint(b, uint64(n))
	for _, v := range s {
		b = protowire.AppendFixed64(b, uint64(v))
	}
	return b, nil
}

var coderSfixed64PackedSlice = pointerCoderFuncs{
	size:      sizeSfixed64PackedSlice,
	marshal:   appendSfixed64PackedSlice,
	unmarshal: consumeSfixed64Slice,
	merge:     mergeInt64Slice,
}


func sizeSfixed64Value(v protoreflect.Value, tagsize int, opts marshalOptions) int {
	return tagsize + protowire.SizeFixed64()
}


func appendSfixed64Value(b []byte, v protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	b = protowire.AppendVarint(b, wiretag)
	b = protowire.AppendFixed64(b, uint64(v.Int()))
	return b, nil
}


func consumeSfixed64Value(b []byte, _ protoreflect.Value, _ protowire.Number, wtyp protowire.Type, opts unmarshalOptions) (_ protoreflect.Value, out unmarshalOutput, err error) {
	if wtyp != protowire.Fixed64Type {
		return protoreflect.Value{}, out, errUnknown
	}
	v, n := protowire.ConsumeFixed64(b)
	if n < 0 {
		return protoreflect.Value{}, out, errDecode
	}
	out.n = n
	return protoreflect.ValueOfInt64(int64(v)), out, nil
}

var coderSfixed64Value = valueCoderFuncs{
	size:      sizeSfixed64Value,
	marshal:   appendSfixed64Value,
	unmarshal: consumeSfixed64Value,
	merge:     mergeScalarValue,
}


func sizeSfixed64SliceValue(listv protoreflect.Value, tagsize int, opts marshalOptions) (size int) {
	list := listv.List()
	size = list.Len() * (tagsize + protowire.SizeFixed64())
	return size
}


func appendSfixed64SliceValue(b []byte, listv protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	list := listv.List()
	for i, llen := 0, list.Len(); i < llen; i++ {
		v := list.Get(i)
		b = protowire.AppendVarint(b, wiretag)
		b = protowire.AppendFixed64(b, uint64(v.Int()))
	}
	return b, nil
}


func consumeSfixed64SliceValue(b []byte, listv protoreflect.Value, _ protowire.Number, wtyp protowire.Type, opts unmarshalOptions) (_ protoreflect.Value, out unmarshalOutput, err error) {
	list := listv.List()
	if wtyp == protowire.BytesType {
		b, n := protowire.ConsumeBytes(b)
		if n < 0 {
			return protoreflect.Value{}, out, errDecode
		}
		for len(b) > 0 {
			v, n := protowire.ConsumeFixed64(b)
			if n < 0 {
				return protoreflect.Value{}, out, errDecode
			}
			list.Append(protoreflect.ValueOfInt64(int64(v)))
			b = b[n:]
		}
		out.n = n
		return listv, out, nil
	}
	if wtyp != protowire.Fixed64Type {
		return protoreflect.Value{}, out, errUnknown
	}
	v, n := protowire.ConsumeFixed64(b)
	if n < 0 {
		return protoreflect.Value{}, out, errDecode
	}
	list.Append(protoreflect.ValueOfInt64(int64(v)))
	out.n = n
	return listv, out, nil
}

var coderSfixed64SliceValue = valueCoderFuncs{
	size:      sizeSfixed64SliceValue,
	marshal:   appendSfixed64SliceValue,
	unmarshal: consumeSfixed64SliceValue,
	merge:     mergeListValue,
}


func sizeSfixed64PackedSliceValue(listv protoreflect.Value, tagsize int, opts marshalOptions) (size int) {
	list := listv.List()
	llen := list.Len()
	if llen == 0 {
		return 0
	}
	n := llen * protowire.SizeFixed64()
	return tagsize + protowire.SizeBytes(n)
}


func appendSfixed64PackedSliceValue(b []byte, listv protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	list := listv.List()
	llen := list.Len()
	if llen == 0 {
		return b, nil
	}
	b = protowire.AppendVarint(b, wiretag)
	n := llen * protowire.SizeFixed64()
	b = protowire.AppendVarint(b, uint64(n))
	for i := 0; i < llen; i++ {
		v := list.Get(i)
		b = protowire.AppendFixed64(b, uint64(v.Int()))
	}
	return b, nil
}

var coderSfixed64PackedSliceValue = valueCoderFuncs{
	size:      sizeSfixed64PackedSliceValue,
	marshal:   appendSfixed64PackedSliceValue,
	unmarshal: consumeSfixed64SliceValue,
	merge:     mergeListValue,
}


func sizeFixed64(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {

	return f.tagsize + protowire.SizeFixed64()
}


func appendFixed64(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := *p.Uint64()
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendFixed64(b, v)
	return b, nil
}


func consumeFixed64(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	if wtyp != protowire.Fixed64Type {
		return out, errUnknown
	}
	v, n := protowire.ConsumeFixed64(b)
	if n < 0 {
		return out, errDecode
	}
	*p.Uint64() = v
	out.n = n
	return out, nil
}

var coderFixed64 = pointerCoderFuncs{
	size:      sizeFixed64,
	marshal:   appendFixed64,
	unmarshal: consumeFixed64,
	merge:     mergeUint64,
}



func sizeFixed64NoZero(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	v := *p.Uint64()
	if v == 0 {
		return 0
	}
	return f.tagsize + protowire.SizeFixed64()
}



func appendFixed64NoZero(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := *p.Uint64()
	if v == 0 {
		return b, nil
	}
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendFixed64(b, v)
	return b, nil
}

var coderFixed64NoZero = pointerCoderFuncs{
	size:      sizeFixed64NoZero,
	marshal:   appendFixed64NoZero,
	unmarshal: consumeFixed64,
	merge:     mergeUint64NoZero,
}



func sizeFixed64Ptr(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	return f.tagsize + protowire.SizeFixed64()
}



func appendFixed64Ptr(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := **p.Uint64Ptr()
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendFixed64(b, v)
	return b, nil
}


func consumeFixed64Ptr(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	if wtyp != protowire.Fixed64Type {
		return out, errUnknown
	}
	v, n := protowire.ConsumeFixed64(b)
	if n < 0 {
		return out, errDecode
	}
	vp := p.Uint64Ptr()
	if *vp == nil {
		*vp = new(uint64)
	}
	**vp = v
	out.n = n
	return out, nil
}

var coderFixed64Ptr = pointerCoderFuncs{
	size:      sizeFixed64Ptr,
	marshal:   appendFixed64Ptr,
	unmarshal: consumeFixed64Ptr,
	merge:     mergeUint64Ptr,
}


func sizeFixed64Slice(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	s := *p.Uint64Slice()
	size = len(s) * (f.tagsize + protowire.SizeFixed64())
	return size
}


func appendFixed64Slice(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	s := *p.Uint64Slice()
	for _, v := range s {
		b = protowire.AppendVarint(b, f.wiretag)
		b = protowire.AppendFixed64(b, v)
	}
	return b, nil
}


func consumeFixed64Slice(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	sp := p.Uint64Slice()
	if wtyp == protowire.BytesType {
		b, n := protowire.ConsumeBytes(b)
		if n < 0 {
			return out, errDecode
		}
		count := len(b) / protowire.SizeFixed64()
		if count > 0 {
			p.growUint64Slice(count)
		}
		s := *sp
		for len(b) > 0 {
			v, n := protowire.ConsumeFixed64(b)
			if n < 0 {
				return out, errDecode
			}
			s = append(s, v)
			b = b[n:]
		}
		*sp = s
		out.n = n
		return out, nil
	}
	if wtyp != protowire.Fixed64Type {
		return out, errUnknown
	}
	v, n := protowire.ConsumeFixed64(b)
	if n < 0 {
		return out, errDecode
	}
	*sp = append(*sp, v)
	out.n = n
	return out, nil
}

var coderFixed64Slice = pointerCoderFuncs{
	size:      sizeFixed64Slice,
	marshal:   appendFixed64Slice,
	unmarshal: consumeFixed64Slice,
	merge:     mergeUint64Slice,
}


func sizeFixed64PackedSlice(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	s := *p.Uint64Slice()
	if len(s) == 0 {
		return 0
	}
	n := len(s) * protowire.SizeFixed64()
	return f.tagsize + protowire.SizeBytes(n)
}


func appendFixed64PackedSlice(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	s := *p.Uint64Slice()
	if len(s) == 0 {
		return b, nil
	}
	b = protowire.AppendVarint(b, f.wiretag)
	n := len(s) * protowire.SizeFixed64()
	b = protowire.AppendVarint(b, uint64(n))
	for _, v := range s {
		b = protowire.AppendFixed64(b, v)
	}
	return b, nil
}

var coderFixed64PackedSlice = pointerCoderFuncs{
	size:      sizeFixed64PackedSlice,
	marshal:   appendFixed64PackedSlice,
	unmarshal: consumeFixed64Slice,
	merge:     mergeUint64Slice,
}


func sizeFixed64Value(v protoreflect.Value, tagsize int, opts marshalOptions) int {
	return tagsize + protowire.SizeFixed64()
}


func appendFixed64Value(b []byte, v protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	b = protowire.AppendVarint(b, wiretag)
	b = protowire.AppendFixed64(b, v.Uint())
	return b, nil
}


func consumeFixed64Value(b []byte, _ protoreflect.Value, _ protowire.Number, wtyp protowire.Type, opts unmarshalOptions) (_ protoreflect.Value, out unmarshalOutput, err error) {
	if wtyp != protowire.Fixed64Type {
		return protoreflect.Value{}, out, errUnknown
	}
	v, n := protowire.ConsumeFixed64(b)
	if n < 0 {
		return protoreflect.Value{}, out, errDecode
	}
	out.n = n
	return protoreflect.ValueOfUint64(v), out, nil
}

var coderFixed64Value = valueCoderFuncs{
	size:      sizeFixed64Value,
	marshal:   appendFixed64Value,
	unmarshal: consumeFixed64Value,
	merge:     mergeScalarValue,
}


func sizeFixed64SliceValue(listv protoreflect.Value, tagsize int, opts marshalOptions) (size int) {
	list := listv.List()
	size = list.Len() * (tagsize + protowire.SizeFixed64())
	return size
}


func appendFixed64SliceValue(b []byte, listv protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	list := listv.List()
	for i, llen := 0, list.Len(); i < llen; i++ {
		v := list.Get(i)
		b = protowire.AppendVarint(b, wiretag)
		b = protowire.AppendFixed64(b, v.Uint())
	}
	return b, nil
}


func consumeFixed64SliceValue(b []byte, listv protoreflect.Value, _ protowire.Number, wtyp protowire.Type, opts unmarshalOptions) (_ protoreflect.Value, out unmarshalOutput, err error) {
	list := listv.List()
	if wtyp == protowire.BytesType {
		b, n := protowire.ConsumeBytes(b)
		if n < 0 {
			return protoreflect.Value{}, out, errDecode
		}
		for len(b) > 0 {
			v, n := protowire.ConsumeFixed64(b)
			if n < 0 {
				return protoreflect.Value{}, out, errDecode
			}
			list.Append(protoreflect.ValueOfUint64(v))
			b = b[n:]
		}
		out.n = n
		return listv, out, nil
	}
	if wtyp != protowire.Fixed64Type {
		return protoreflect.Value{}, out, errUnknown
	}
	v, n := protowire.ConsumeFixed64(b)
	if n < 0 {
		return protoreflect.Value{}, out, errDecode
	}
	list.Append(protoreflect.ValueOfUint64(v))
	out.n = n
	return listv, out, nil
}

var coderFixed64SliceValue = valueCoderFuncs{
	size:      sizeFixed64SliceValue,
	marshal:   appendFixed64SliceValue,
	unmarshal: consumeFixed64SliceValue,
	merge:     mergeListValue,
}


func sizeFixed64PackedSliceValue(listv protoreflect.Value, tagsize int, opts marshalOptions) (size int) {
	list := listv.List()
	llen := list.Len()
	if llen == 0 {
		return 0
	}
	n := llen * protowire.SizeFixed64()
	return tagsize + protowire.SizeBytes(n)
}


func appendFixed64PackedSliceValue(b []byte, listv protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	list := listv.List()
	llen := list.Len()
	if llen == 0 {
		return b, nil
	}
	b = protowire.AppendVarint(b, wiretag)
	n := llen * protowire.SizeFixed64()
	b = protowire.AppendVarint(b, uint64(n))
	for i := 0; i < llen; i++ {
		v := list.Get(i)
		b = protowire.AppendFixed64(b, v.Uint())
	}
	return b, nil
}

var coderFixed64PackedSliceValue = valueCoderFuncs{
	size:      sizeFixed64PackedSliceValue,
	marshal:   appendFixed64PackedSliceValue,
	unmarshal: consumeFixed64SliceValue,
	merge:     mergeListValue,
}


func sizeDouble(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {

	return f.tagsize + protowire.SizeFixed64()
}


func appendDouble(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := *p.Float64()
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendFixed64(b, math.Float64bits(v))
	return b, nil
}


func consumeDouble(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	if wtyp != protowire.Fixed64Type {
		return out, errUnknown
	}
	v, n := protowire.ConsumeFixed64(b)
	if n < 0 {
		return out, errDecode
	}
	*p.Float64() = math.Float64frombits(v)
	out.n = n
	return out, nil
}

var coderDouble = pointerCoderFuncs{
	size:      sizeDouble,
	marshal:   appendDouble,
	unmarshal: consumeDouble,
	merge:     mergeFloat64,
}



func sizeDoubleNoZero(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	v := *p.Float64()
	if v == 0 && !math.Signbit(float64(v)) {
		return 0
	}
	return f.tagsize + protowire.SizeFixed64()
}



func appendDoubleNoZero(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := *p.Float64()
	if v == 0 && !math.Signbit(float64(v)) {
		return b, nil
	}
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendFixed64(b, math.Float64bits(v))
	return b, nil
}

var coderDoubleNoZero = pointerCoderFuncs{
	size:      sizeDoubleNoZero,
	marshal:   appendDoubleNoZero,
	unmarshal: consumeDouble,
	merge:     mergeFloat64NoZero,
}



func sizeDoublePtr(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	return f.tagsize + protowire.SizeFixed64()
}



func appendDoublePtr(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := **p.Float64Ptr()
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendFixed64(b, math.Float64bits(v))
	return b, nil
}


func consumeDoublePtr(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	if wtyp != protowire.Fixed64Type {
		return out, errUnknown
	}
	v, n := protowire.ConsumeFixed64(b)
	if n < 0 {
		return out, errDecode
	}
	vp := p.Float64Ptr()
	if *vp == nil {
		*vp = new(float64)
	}
	**vp = math.Float64frombits(v)
	out.n = n
	return out, nil
}

var coderDoublePtr = pointerCoderFuncs{
	size:      sizeDoublePtr,
	marshal:   appendDoublePtr,
	unmarshal: consumeDoublePtr,
	merge:     mergeFloat64Ptr,
}


func sizeDoubleSlice(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	s := *p.Float64Slice()
	size = len(s) * (f.tagsize + protowire.SizeFixed64())
	return size
}


func appendDoubleSlice(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	s := *p.Float64Slice()
	for _, v := range s {
		b = protowire.AppendVarint(b, f.wiretag)
		b = protowire.AppendFixed64(b, math.Float64bits(v))
	}
	return b, nil
}


func consumeDoubleSlice(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	sp := p.Float64Slice()
	if wtyp == protowire.BytesType {
		b, n := protowire.ConsumeBytes(b)
		if n < 0 {
			return out, errDecode
		}
		count := len(b) / protowire.SizeFixed64()
		if count > 0 {
			p.growFloat64Slice(count)
		}
		s := *sp
		for len(b) > 0 {
			v, n := protowire.ConsumeFixed64(b)
			if n < 0 {
				return out, errDecode
			}
			s = append(s, math.Float64frombits(v))
			b = b[n:]
		}
		*sp = s
		out.n = n
		return out, nil
	}
	if wtyp != protowire.Fixed64Type {
		return out, errUnknown
	}
	v, n := protowire.ConsumeFixed64(b)
	if n < 0 {
		return out, errDecode
	}
	*sp = append(*sp, math.Float64frombits(v))
	out.n = n
	return out, nil
}

var coderDoubleSlice = pointerCoderFuncs{
	size:      sizeDoubleSlice,
	marshal:   appendDoubleSlice,
	unmarshal: consumeDoubleSlice,
	merge:     mergeFloat64Slice,
}


func sizeDoublePackedSlice(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	s := *p.Float64Slice()
	if len(s) == 0 {
		return 0
	}
	n := len(s) * protowire.SizeFixed64()
	return f.tagsize + protowire.SizeBytes(n)
}


func appendDoublePackedSlice(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	s := *p.Float64Slice()
	if len(s) == 0 {
		return b, nil
	}
	b = protowire.AppendVarint(b, f.wiretag)
	n := len(s) * protowire.SizeFixed64()
	b = protowire.AppendVarint(b, uint64(n))
	for _, v := range s {
		b = protowire.AppendFixed64(b, math.Float64bits(v))
	}
	return b, nil
}

var coderDoublePackedSlice = pointerCoderFuncs{
	size:      sizeDoublePackedSlice,
	marshal:   appendDoublePackedSlice,
	unmarshal: consumeDoubleSlice,
	merge:     mergeFloat64Slice,
}


func sizeDoubleValue(v protoreflect.Value, tagsize int, opts marshalOptions) int {
	return tagsize + protowire.SizeFixed64()
}


func appendDoubleValue(b []byte, v protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	b = protowire.AppendVarint(b, wiretag)
	b = protowire.AppendFixed64(b, math.Float64bits(v.Float()))
	return b, nil
}


func consumeDoubleValue(b []byte, _ protoreflect.Value, _ protowire.Number, wtyp protowire.Type, opts unmarshalOptions) (_ protoreflect.Value, out unmarshalOutput, err error) {
	if wtyp != protowire.Fixed64Type {
		return protoreflect.Value{}, out, errUnknown
	}
	v, n := protowire.ConsumeFixed64(b)
	if n < 0 {
		return protoreflect.Value{}, out, errDecode
	}
	out.n = n
	return protoreflect.ValueOfFloat64(math.Float64frombits(v)), out, nil
}

var coderDoubleValue = valueCoderFuncs{
	size:      sizeDoubleValue,
	marshal:   appendDoubleValue,
	unmarshal: consumeDoubleValue,
	merge:     mergeScalarValue,
}


func sizeDoubleSliceValue(listv protoreflect.Value, tagsize int, opts marshalOptions) (size int) {
	list := listv.List()
	size = list.Len() * (tagsize + protowire.SizeFixed64())
	return size
}


func appendDoubleSliceValue(b []byte, listv protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	list := listv.List()
	for i, llen := 0, list.Len(); i < llen; i++ {
		v := list.Get(i)
		b = protowire.AppendVarint(b, wiretag)
		b = protowire.AppendFixed64(b, math.Float64bits(v.Float()))
	}
	return b, nil
}


func consumeDoubleSliceValue(b []byte, listv protoreflect.Value, _ protowire.Number, wtyp protowire.Type, opts unmarshalOptions) (_ protoreflect.Value, out unmarshalOutput, err error) {
	list := listv.List()
	if wtyp == protowire.BytesType {
		b, n := protowire.ConsumeBytes(b)
		if n < 0 {
			return protoreflect.Value{}, out, errDecode
		}
		for len(b) > 0 {
			v, n := protowire.ConsumeFixed64(b)
			if n < 0 {
				return protoreflect.Value{}, out, errDecode
			}
			list.Append(protoreflect.ValueOfFloat64(math.Float64frombits(v)))
			b = b[n:]
		}
		out.n = n
		return listv, out, nil
	}
	if wtyp != protowire.Fixed64Type {
		return protoreflect.Value{}, out, errUnknown
	}
	v, n := protowire.ConsumeFixed64(b)
	if n < 0 {
		return protoreflect.Value{}, out, errDecode
	}
	list.Append(protoreflect.ValueOfFloat64(math.Float64frombits(v)))
	out.n = n
	return listv, out, nil
}

var coderDoubleSliceValue = valueCoderFuncs{
	size:      sizeDoubleSliceValue,
	marshal:   appendDoubleSliceValue,
	unmarshal: consumeDoubleSliceValue,
	merge:     mergeListValue,
}


func sizeDoublePackedSliceValue(listv protoreflect.Value, tagsize int, opts marshalOptions) (size int) {
	list := listv.List()
	llen := list.Len()
	if llen == 0 {
		return 0
	}
	n := llen * protowire.SizeFixed64()
	return tagsize + protowire.SizeBytes(n)
}


func appendDoublePackedSliceValue(b []byte, listv protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	list := listv.List()
	llen := list.Len()
	if llen == 0 {
		return b, nil
	}
	b = protowire.AppendVarint(b, wiretag)
	n := llen * protowire.SizeFixed64()
	b = protowire.AppendVarint(b, uint64(n))
	for i := 0; i < llen; i++ {
		v := list.Get(i)
		b = protowire.AppendFixed64(b, math.Float64bits(v.Float()))
	}
	return b, nil
}

var coderDoublePackedSliceValue = valueCoderFuncs{
	size:      sizeDoublePackedSliceValue,
	marshal:   appendDoublePackedSliceValue,
	unmarshal: consumeDoubleSliceValue,
	merge:     mergeListValue,
}


func sizeString(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	v := *p.String()
	return f.tagsize + protowire.SizeBytes(len(v))
}


func appendString(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := *p.String()
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendString(b, v)
	return b, nil
}


func consumeString(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	if wtyp != protowire.BytesType {
		return out, errUnknown
	}
	v, n := protowire.ConsumeBytes(b)
	if n < 0 {
		return out, errDecode
	}
	*p.String() = string(v)
	out.n = n
	return out, nil
}

var coderString = pointerCoderFuncs{
	size:      sizeString,
	marshal:   appendString,
	unmarshal: consumeString,
	merge:     mergeString,
}


func appendStringValidateUTF8(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := *p.String()
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendString(b, v)
	if !utf8.ValidString(v) {
		return b, errInvalidUTF8{}
	}
	return b, nil
}


func consumeStringValidateUTF8(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	if wtyp != protowire.BytesType {
		return out, errUnknown
	}
	v, n := protowire.ConsumeBytes(b)
	if n < 0 {
		return out, errDecode
	}
	if !utf8.Valid(v) {
		return out, errInvalidUTF8{}
	}
	*p.String() = string(v)
	out.n = n
	return out, nil
}

var coderStringValidateUTF8 = pointerCoderFuncs{
	size:      sizeString,
	marshal:   appendStringValidateUTF8,
	unmarshal: consumeStringValidateUTF8,
	merge:     mergeString,
}



func sizeStringNoZero(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	v := *p.String()
	if len(v) == 0 {
		return 0
	}
	return f.tagsize + protowire.SizeBytes(len(v))
}



func appendStringNoZero(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := *p.String()
	if len(v) == 0 {
		return b, nil
	}
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendString(b, v)
	return b, nil
}

var coderStringNoZero = pointerCoderFuncs{
	size:      sizeStringNoZero,
	marshal:   appendStringNoZero,
	unmarshal: consumeString,
	merge:     mergeStringNoZero,
}



func appendStringNoZeroValidateUTF8(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := *p.String()
	if len(v) == 0 {
		return b, nil
	}
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendString(b, v)
	if !utf8.ValidString(v) {
		return b, errInvalidUTF8{}
	}
	return b, nil
}

var coderStringNoZeroValidateUTF8 = pointerCoderFuncs{
	size:      sizeStringNoZero,
	marshal:   appendStringNoZeroValidateUTF8,
	unmarshal: consumeStringValidateUTF8,
	merge:     mergeStringNoZero,
}



func sizeStringPtr(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	v := **p.StringPtr()
	return f.tagsize + protowire.SizeBytes(len(v))
}



func appendStringPtr(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := **p.StringPtr()
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendString(b, v)
	return b, nil
}


func consumeStringPtr(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	if wtyp != protowire.BytesType {
		return out, errUnknown
	}
	v, n := protowire.ConsumeBytes(b)
	if n < 0 {
		return out, errDecode
	}
	vp := p.StringPtr()
	if *vp == nil {
		*vp = new(string)
	}
	**vp = string(v)
	out.n = n
	return out, nil
}

var coderStringPtr = pointerCoderFuncs{
	size:      sizeStringPtr,
	marshal:   appendStringPtr,
	unmarshal: consumeStringPtr,
	merge:     mergeStringPtr,
}



func appendStringPtrValidateUTF8(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := **p.StringPtr()
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendString(b, v)
	if !utf8.ValidString(v) {
		return b, errInvalidUTF8{}
	}
	return b, nil
}


func consumeStringPtrValidateUTF8(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	if wtyp != protowire.BytesType {
		return out, errUnknown
	}
	v, n := protowire.ConsumeBytes(b)
	if n < 0 {
		return out, errDecode
	}
	if !utf8.Valid(v) {
		return out, errInvalidUTF8{}
	}
	vp := p.StringPtr()
	if *vp == nil {
		*vp = new(string)
	}
	**vp = string(v)
	out.n = n
	return out, nil
}

var coderStringPtrValidateUTF8 = pointerCoderFuncs{
	size:      sizeStringPtr,
	marshal:   appendStringPtrValidateUTF8,
	unmarshal: consumeStringPtrValidateUTF8,
	merge:     mergeStringPtr,
}


func sizeStringSlice(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	s := *p.StringSlice()
	for _, v := range s {
		size += f.tagsize + protowire.SizeBytes(len(v))
	}
	return size
}


func appendStringSlice(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	s := *p.StringSlice()
	for _, v := range s {
		b = protowire.AppendVarint(b, f.wiretag)
		b = protowire.AppendString(b, v)
	}
	return b, nil
}


func consumeStringSlice(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	sp := p.StringSlice()
	if wtyp != protowire.BytesType {
		return out, errUnknown
	}
	v, n := protowire.ConsumeBytes(b)
	if n < 0 {
		return out, errDecode
	}
	*sp = append(*sp, string(v))
	out.n = n
	return out, nil
}

var coderStringSlice = pointerCoderFuncs{
	size:      sizeStringSlice,
	marshal:   appendStringSlice,
	unmarshal: consumeStringSlice,
	merge:     mergeStringSlice,
}


func appendStringSliceValidateUTF8(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	s := *p.StringSlice()
	for _, v := range s {
		b = protowire.AppendVarint(b, f.wiretag)
		b = protowire.AppendString(b, v)
		if !utf8.ValidString(v) {
			return b, errInvalidUTF8{}
		}
	}
	return b, nil
}


func consumeStringSliceValidateUTF8(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	if wtyp != protowire.BytesType {
		return out, errUnknown
	}
	v, n := protowire.ConsumeBytes(b)
	if n < 0 {
		return out, errDecode
	}
	if !utf8.Valid(v) {
		return out, errInvalidUTF8{}
	}
	sp := p.StringSlice()
	*sp = append(*sp, string(v))
	out.n = n
	return out, nil
}

var coderStringSliceValidateUTF8 = pointerCoderFuncs{
	size:      sizeStringSlice,
	marshal:   appendStringSliceValidateUTF8,
	unmarshal: consumeStringSliceValidateUTF8,
	merge:     mergeStringSlice,
}


func sizeStringValue(v protoreflect.Value, tagsize int, opts marshalOptions) int {
	return tagsize + protowire.SizeBytes(len(v.String()))
}


func appendStringValue(b []byte, v protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	b = protowire.AppendVarint(b, wiretag)
	b = protowire.AppendString(b, v.String())
	return b, nil
}


func consumeStringValue(b []byte, _ protoreflect.Value, _ protowire.Number, wtyp protowire.Type, opts unmarshalOptions) (_ protoreflect.Value, out unmarshalOutput, err error) {
	if wtyp != protowire.BytesType {
		return protoreflect.Value{}, out, errUnknown
	}
	v, n := protowire.ConsumeBytes(b)
	if n < 0 {
		return protoreflect.Value{}, out, errDecode
	}
	out.n = n
	return protoreflect.ValueOfString(string(v)), out, nil
}

var coderStringValue = valueCoderFuncs{
	size:      sizeStringValue,
	marshal:   appendStringValue,
	unmarshal: consumeStringValue,
	merge:     mergeScalarValue,
}


func appendStringValueValidateUTF8(b []byte, v protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	b = protowire.AppendVarint(b, wiretag)
	b = protowire.AppendString(b, v.String())
	if !utf8.ValidString(v.String()) {
		return b, errInvalidUTF8{}
	}
	return b, nil
}


func consumeStringValueValidateUTF8(b []byte, _ protoreflect.Value, _ protowire.Number, wtyp protowire.Type, opts unmarshalOptions) (_ protoreflect.Value, out unmarshalOutput, err error) {
	if wtyp != protowire.BytesType {
		return protoreflect.Value{}, out, errUnknown
	}
	v, n := protowire.ConsumeBytes(b)
	if n < 0 {
		return protoreflect.Value{}, out, errDecode
	}
	if !utf8.Valid(v) {
		return protoreflect.Value{}, out, errInvalidUTF8{}
	}
	out.n = n
	return protoreflect.ValueOfString(string(v)), out, nil
}

var coderStringValueValidateUTF8 = valueCoderFuncs{
	size:      sizeStringValue,
	marshal:   appendStringValueValidateUTF8,
	unmarshal: consumeStringValueValidateUTF8,
	merge:     mergeScalarValue,
}


func sizeStringSliceValue(listv protoreflect.Value, tagsize int, opts marshalOptions) (size int) {
	list := listv.List()
	for i, llen := 0, list.Len(); i < llen; i++ {
		v := list.Get(i)
		size += tagsize + protowire.SizeBytes(len(v.String()))
	}
	return size
}


func appendStringSliceValue(b []byte, listv protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	list := listv.List()
	for i, llen := 0, list.Len(); i < llen; i++ {
		v := list.Get(i)
		b = protowire.AppendVarint(b, wiretag)
		b = protowire.AppendString(b, v.String())
	}
	return b, nil
}


func consumeStringSliceValue(b []byte, listv protoreflect.Value, _ protowire.Number, wtyp protowire.Type, opts unmarshalOptions) (_ protoreflect.Value, out unmarshalOutput, err error) {
	list := listv.List()
	if wtyp != protowire.BytesType {
		return protoreflect.Value{}, out, errUnknown
	}
	v, n := protowire.ConsumeBytes(b)
	if n < 0 {
		return protoreflect.Value{}, out, errDecode
	}
	list.Append(protoreflect.ValueOfString(string(v)))
	out.n = n
	return listv, out, nil
}

var coderStringSliceValue = valueCoderFuncs{
	size:      sizeStringSliceValue,
	marshal:   appendStringSliceValue,
	unmarshal: consumeStringSliceValue,
	merge:     mergeListValue,
}


func sizeBytes(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	v := *p.Bytes()
	return f.tagsize + protowire.SizeBytes(len(v))
}


func appendBytes(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := *p.Bytes()
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendBytes(b, v)
	return b, nil
}


func consumeBytes(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	if wtyp != protowire.BytesType {
		return out, errUnknown
	}
	v, n := protowire.ConsumeBytes(b)
	if n < 0 {
		return out, errDecode
	}
	*p.Bytes() = append(emptyBuf[:], v...)
	out.n = n
	return out, nil
}

var coderBytes = pointerCoderFuncs{
	size:      sizeBytes,
	marshal:   appendBytes,
	unmarshal: consumeBytes,
	merge:     mergeBytes,
}


func appendBytesValidateUTF8(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := *p.Bytes()
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendBytes(b, v)
	if !utf8.Valid(v) {
		return b, errInvalidUTF8{}
	}
	return b, nil
}


func consumeBytesValidateUTF8(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	if wtyp != protowire.BytesType {
		return out, errUnknown
	}
	v, n := protowire.ConsumeBytes(b)
	if n < 0 {
		return out, errDecode
	}
	if !utf8.Valid(v) {
		return out, errInvalidUTF8{}
	}
	*p.Bytes() = append(emptyBuf[:], v...)
	out.n = n
	return out, nil
}

var coderBytesValidateUTF8 = pointerCoderFuncs{
	size:      sizeBytes,
	marshal:   appendBytesValidateUTF8,
	unmarshal: consumeBytesValidateUTF8,
	merge:     mergeBytes,
}



func sizeBytesNoZero(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	v := *p.Bytes()
	if len(v) == 0 {
		return 0
	}
	return f.tagsize + protowire.SizeBytes(len(v))
}



func appendBytesNoZero(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := *p.Bytes()
	if len(v) == 0 {
		return b, nil
	}
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendBytes(b, v)
	return b, nil
}



func consumeBytesNoZero(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	if wtyp != protowire.BytesType {
		return out, errUnknown
	}
	v, n := protowire.ConsumeBytes(b)
	if n < 0 {
		return out, errDecode
	}
	*p.Bytes() = append(([]byte)(nil), v...)
	out.n = n
	return out, nil
}

var coderBytesNoZero = pointerCoderFuncs{
	size:      sizeBytesNoZero,
	marshal:   appendBytesNoZero,
	unmarshal: consumeBytesNoZero,
	merge:     mergeBytesNoZero,
}



func appendBytesNoZeroValidateUTF8(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	v := *p.Bytes()
	if len(v) == 0 {
		return b, nil
	}
	b = protowire.AppendVarint(b, f.wiretag)
	b = protowire.AppendBytes(b, v)
	if !utf8.Valid(v) {
		return b, errInvalidUTF8{}
	}
	return b, nil
}


func consumeBytesNoZeroValidateUTF8(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	if wtyp != protowire.BytesType {
		return out, errUnknown
	}
	v, n := protowire.ConsumeBytes(b)
	if n < 0 {
		return out, errDecode
	}
	if !utf8.Valid(v) {
		return out, errInvalidUTF8{}
	}
	*p.Bytes() = append(([]byte)(nil), v...)
	out.n = n
	return out, nil
}

var coderBytesNoZeroValidateUTF8 = pointerCoderFuncs{
	size:      sizeBytesNoZero,
	marshal:   appendBytesNoZeroValidateUTF8,
	unmarshal: consumeBytesNoZeroValidateUTF8,
	merge:     mergeBytesNoZero,
}


func sizeBytesSlice(p pointer, f *coderFieldInfo, opts marshalOptions) (size int) {
	s := *p.BytesSlice()
	for _, v := range s {
		size += f.tagsize + protowire.SizeBytes(len(v))
	}
	return size
}


func appendBytesSlice(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	s := *p.BytesSlice()
	for _, v := range s {
		b = protowire.AppendVarint(b, f.wiretag)
		b = protowire.AppendBytes(b, v)
	}
	return b, nil
}


func consumeBytesSlice(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	sp := p.BytesSlice()
	if wtyp != protowire.BytesType {
		return out, errUnknown
	}
	v, n := protowire.ConsumeBytes(b)
	if n < 0 {
		return out, errDecode
	}
	*sp = append(*sp, append(emptyBuf[:], v...))
	out.n = n
	return out, nil
}

var coderBytesSlice = pointerCoderFuncs{
	size:      sizeBytesSlice,
	marshal:   appendBytesSlice,
	unmarshal: consumeBytesSlice,
	merge:     mergeBytesSlice,
}


func appendBytesSliceValidateUTF8(b []byte, p pointer, f *coderFieldInfo, opts marshalOptions) ([]byte, error) {
	s := *p.BytesSlice()
	for _, v := range s {
		b = protowire.AppendVarint(b, f.wiretag)
		b = protowire.AppendBytes(b, v)
		if !utf8.Valid(v) {
			return b, errInvalidUTF8{}
		}
	}
	return b, nil
}


func consumeBytesSliceValidateUTF8(b []byte, p pointer, wtyp protowire.Type, f *coderFieldInfo, opts unmarshalOptions) (out unmarshalOutput, err error) {
	if wtyp != protowire.BytesType {
		return out, errUnknown
	}
	v, n := protowire.ConsumeBytes(b)
	if n < 0 {
		return out, errDecode
	}
	if !utf8.Valid(v) {
		return out, errInvalidUTF8{}
	}
	sp := p.BytesSlice()
	*sp = append(*sp, append(emptyBuf[:], v...))
	out.n = n
	return out, nil
}

var coderBytesSliceValidateUTF8 = pointerCoderFuncs{
	size:      sizeBytesSlice,
	marshal:   appendBytesSliceValidateUTF8,
	unmarshal: consumeBytesSliceValidateUTF8,
	merge:     mergeBytesSlice,
}


func sizeBytesValue(v protoreflect.Value, tagsize int, opts marshalOptions) int {
	return tagsize + protowire.SizeBytes(len(v.Bytes()))
}


func appendBytesValue(b []byte, v protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	b = protowire.AppendVarint(b, wiretag)
	b = protowire.AppendBytes(b, v.Bytes())
	return b, nil
}


func consumeBytesValue(b []byte, _ protoreflect.Value, _ protowire.Number, wtyp protowire.Type, opts unmarshalOptions) (_ protoreflect.Value, out unmarshalOutput, err error) {
	if wtyp != protowire.BytesType {
		return protoreflect.Value{}, out, errUnknown
	}
	v, n := protowire.ConsumeBytes(b)
	if n < 0 {
		return protoreflect.Value{}, out, errDecode
	}
	out.n = n
	return protoreflect.ValueOfBytes(append(emptyBuf[:], v...)), out, nil
}

var coderBytesValue = valueCoderFuncs{
	size:      sizeBytesValue,
	marshal:   appendBytesValue,
	unmarshal: consumeBytesValue,
	merge:     mergeBytesValue,
}


func sizeBytesSliceValue(listv protoreflect.Value, tagsize int, opts marshalOptions) (size int) {
	list := listv.List()
	for i, llen := 0, list.Len(); i < llen; i++ {
		v := list.Get(i)
		size += tagsize + protowire.SizeBytes(len(v.Bytes()))
	}
	return size
}


func appendBytesSliceValue(b []byte, listv protoreflect.Value, wiretag uint64, opts marshalOptions) ([]byte, error) {
	list := listv.List()
	for i, llen := 0, list.Len(); i < llen; i++ {
		v := list.Get(i)
		b = protowire.AppendVarint(b, wiretag)
		b = protowire.AppendBytes(b, v.Bytes())
	}
	return b, nil
}


func consumeBytesSliceValue(b []byte, listv protoreflect.Value, _ protowire.Number, wtyp protowire.Type, opts unmarshalOptions) (_ protoreflect.Value, out unmarshalOutput, err error) {
	list := listv.List()
	if wtyp != protowire.BytesType {
		return protoreflect.Value{}, out, errUnknown
	}
	v, n := protowire.ConsumeBytes(b)
	if n < 0 {
		return protoreflect.Value{}, out, errDecode
	}
	list.Append(protoreflect.ValueOfBytes(append(emptyBuf[:], v...)))
	out.n = n
	return listv, out, nil
}

var coderBytesSliceValue = valueCoderFuncs{
	size:      sizeBytesSliceValue,
	marshal:   appendBytesSliceValue,
	unmarshal: consumeBytesSliceValue,
	merge:     mergeBytesListValue,
}


var emptyBuf [0]byte

var wireTypes = map[protoreflect.Kind]protowire.Type{
	protoreflect.BoolKind:     protowire.VarintType,
	protoreflect.EnumKind:     protowire.VarintType,
	protoreflect.Int32Kind:    protowire.VarintType,
	protoreflect.Sint32Kind:   protowire.VarintType,
	protoreflect.Uint32Kind:   protowire.VarintType,
	protoreflect.Int64Kind:    protowire.VarintType,
	protoreflect.Sint64Kind:   protowire.VarintType,
	protoreflect.Uint64Kind:   protowire.VarintType,
	protoreflect.Sfixed32Kind: protowire.Fixed32Type,
	protoreflect.Fixed32Kind:  protowire.Fixed32Type,
	protoreflect.FloatKind:    protowire.Fixed32Type,
	protoreflect.Sfixed64Kind: protowire.Fixed64Type,
	protoreflect.Fixed64Kind:  protowire.Fixed64Type,
	protoreflect.DoubleKind:   protowire.Fixed64Type,
	protoreflect.StringKind:   protowire.BytesType,
	protoreflect.BytesKind:    protowire.BytesType,
	protoreflect.MessageKind:  protowire.BytesType,
	protoreflect.GroupKind:    protowire.StartGroupType,
}
