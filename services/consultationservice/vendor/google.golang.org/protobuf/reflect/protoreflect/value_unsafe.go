



package protoreflect

import (
	"unsafe"

	"google.golang.org/protobuf/internal/pragma"
)

type (
	ifaceHeader struct {
		_    [0]any 
		Type unsafe.Pointer
		Data unsafe.Pointer
	}
)

var (
	nilType     = typeOf(nil)
	boolType    = typeOf(*new(bool))
	int32Type   = typeOf(*new(int32))
	int64Type   = typeOf(*new(int64))
	uint32Type  = typeOf(*new(uint32))
	uint64Type  = typeOf(*new(uint64))
	float32Type = typeOf(*new(float32))
	float64Type = typeOf(*new(float64))
	stringType  = typeOf(*new(string))
	bytesType   = typeOf(*new([]byte))
	enumType    = typeOf(*new(EnumNumber))
)



func typeOf(t any) unsafe.Pointer {
	return (*ifaceHeader)(unsafe.Pointer(&t)).Type
}







type value struct {
	pragma.DoNotCompare 

	
	typ unsafe.Pointer 

	
	ptr unsafe.Pointer 

	
	
	
	
	
	num uint64 
}

func valueOfString(v string) Value {
	return Value{typ: stringType, ptr: unsafe.Pointer(unsafe.StringData(v)), num: uint64(len(v))}
}
func valueOfBytes(v []byte) Value {
	return Value{typ: bytesType, ptr: unsafe.Pointer(unsafe.SliceData(v)), num: uint64(len(v))}
}
func valueOfIface(v any) Value {
	p := (*ifaceHeader)(unsafe.Pointer(&v))
	return Value{typ: p.Type, ptr: p.Data}
}

func (v Value) getString() string {
	return unsafe.String((*byte)(v.ptr), v.num)
}
func (v Value) getBytes() []byte {
	return unsafe.Slice((*byte)(v.ptr), v.num)
}
func (v Value) getIface() (x any) {
	*(*ifaceHeader)(unsafe.Pointer(&x)) = ifaceHeader{Type: v.typ, Data: v.ptr}
	return x
}
