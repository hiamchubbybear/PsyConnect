



package protoreflect

import "google.golang.org/protobuf/encoding/protowire"





type Enum interface {
	
	
	Descriptor() EnumDescriptor

	
	
	
	Type() EnumType

	
	Number() EnumNumber
}














type Message interface {
	
	
	Descriptor() MessageDescriptor

	
	
	
	Type() MessageType

	
	New() Message

	
	
	Interface() ProtoMessage

	
	
	
	
	
	Range(f func(FieldDescriptor, Value) bool)

	
	
	
	
	
	
	
	
	
	
	
	Has(FieldDescriptor) bool

	
	
	
	
	
	
	Clear(FieldDescriptor)

	
	
	
	
	
	
	Get(FieldDescriptor) Value

	
	
	
	
	
	
	
	
	
	
	Set(FieldDescriptor, Value)

	
	
	
	
	
	
	
	
	
	
	Mutable(FieldDescriptor) Value

	
	
	
	NewField(FieldDescriptor) Value

	
	
	
	WhichOneof(OneofDescriptor) FieldDescriptor

	
	
	
	GetUnknown() RawFields

	
	
	
	
	
	
	
	SetUnknown(RawFields)

	
	
	
	
	
	
	
	
	IsValid() bool

	
	
	
	
	
	
	ProtoMethods() *methods
}




type RawFields []byte


func (b RawFields) IsValid() bool {
	for len(b) > 0 {
		_, _, n := protowire.ConsumeField(b)
		if n < 0 {
			return false
		}
		b = b[n:]
	}
	return true
}




type List interface {
	
	
	Len() int

	
	
	Get(int) Value

	
	
	
	
	
	Set(int, Value)

	
	
	
	
	
	Append(Value)

	
	
	
	AppendMutable() Value

	
	
	
	Truncate(int)

	
	
	
	
	NewElement() Value

	
	
	
	
	
	
	IsValid() bool
}





type Map interface {
	
	Len() int

	
	
	
	
	
	Range(f func(MapKey, Value) bool)

	
	Has(MapKey) bool

	
	
	
	
	Clear(MapKey)

	
	
	Get(MapKey) Value

	
	
	
	
	
	
	Set(MapKey, Value)

	
	
	
	
	Mutable(MapKey) Value

	
	
	
	
	NewValue() Value

	
	
	
	
	
	
	
	
	IsValid() bool
}
