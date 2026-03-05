



package strs

import (
	"unsafe"

	"google.golang.org/protobuf/reflect/protoreflect"
)






func UnsafeString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}





func UnsafeBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}



type Builder struct {
	buf []byte
}



func (sb *Builder) AppendFullName(prefix protoreflect.FullName, name protoreflect.Name) protoreflect.FullName {
	n := len(prefix) + len(".") + len(name)
	if len(prefix) == 0 {
		n -= len(".")
	}
	sb.grow(n)
	sb.buf = append(sb.buf, prefix...)
	sb.buf = append(sb.buf, '.')
	sb.buf = append(sb.buf, name...)
	return protoreflect.FullName(sb.last(n))
}



func (sb *Builder) MakeString(b []byte) string {
	sb.grow(len(b))
	sb.buf = append(sb.buf, b...)
	return sb.last(len(b))
}

func (sb *Builder) grow(n int) {
	if cap(sb.buf)-len(sb.buf) >= n {
		return
	}

	
	
	
	sb.buf = make([]byte, 0, 2*(cap(sb.buf)+n))
}

func (sb *Builder) last(n int) string {
	return UnsafeString(sb.buf[len(sb.buf)-n:])
}
