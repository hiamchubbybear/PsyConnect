



package protoreflect















type Descriptor interface {
	
	
	
	
	ParentFile() FileDescriptor

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	Parent() Descriptor

	
	
	
	Index() int

	
	Syntax() Syntax 

	
	Name() Name 

	
	
	
	
	
	
	
	FullName() FullName 

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	IsPlaceholder() bool

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	Options() ProtoMessage

	doNotImplement
}






type FileDescriptor interface {
	Descriptor 

	
	Path() string 
	
	Package() FullName 

	
	Imports() FileImports

	
	Enums() EnumDescriptors
	
	Messages() MessageDescriptors
	
	Extensions() ExtensionDescriptors
	
	Services() ServiceDescriptors

	
	SourceLocations() SourceLocations

	isFileDescriptor
}
type isFileDescriptor interface{ ProtoType(FileDescriptor) }


type FileImports interface {
	
	Len() int
	
	Get(i int) FileImport

	doNotImplement
}


type FileImport struct {
	
	
	
	FileDescriptor

	
	
	
	
	
	
	IsPublic bool

	
	IsWeak bool
}







type MessageDescriptor interface {
	Descriptor

	
	
	
	
	
	
	
	
	
	
	IsMapEntry() bool

	
	Fields() FieldDescriptors
	
	Oneofs() OneofDescriptors

	
	ReservedNames() Names
	
	ReservedRanges() FieldRanges
	
	
	RequiredNumbers() FieldNumbers
	
	
	ExtensionRanges() FieldRanges
	
	
	
	
	
	
	ExtensionRangeOptions(i int) ProtoMessage

	
	Enums() EnumDescriptors
	
	Messages() MessageDescriptors
	
	Extensions() ExtensionDescriptors

	isMessageDescriptor
}
type isMessageDescriptor interface{ ProtoType(MessageDescriptor) }




type MessageType interface {
	
	
	New() Message

	
	
	Zero() Message

	
	
	
	Descriptor() MessageDescriptor
}



type MessageFieldTypes interface {
	MessageType

	
	
	
	
	
	Enum(i int) EnumType

	
	
	
	
	
	Message(i int) MessageType
}


type MessageDescriptors interface {
	
	Len() int
	
	Get(i int) MessageDescriptor
	
	
	ByName(s Name) MessageDescriptor

	doNotImplement
}







type FieldDescriptor interface {
	Descriptor

	
	Number() FieldNumber
	
	Cardinality() Cardinality
	
	Kind() Kind

	
	HasJSONName() bool

	
	
	
	JSONName() string

	
	
	
	
	TextName() string

	
	
	HasPresence() bool

	
	
	
	IsExtension() bool

	
	
	HasOptionalKeyword() bool

	
	IsWeak() bool

	
	
	
	IsPacked() bool

	
	
	
	
	IsList() bool

	
	
	
	
	IsMap() bool

	
	
	MapKey() FieldDescriptor

	
	
	MapValue() FieldDescriptor

	
	HasDefault() bool

	
	
	
	
	
	Default() Value

	
	
	DefaultEnumValue() EnumValueDescriptor

	
	
	ContainingOneof() OneofDescriptor

	
	
	
	ContainingMessage() MessageDescriptor

	
	
	Enum() EnumDescriptor

	
	
	Message() MessageDescriptor

	isFieldDescriptor
}
type isFieldDescriptor interface{ ProtoType(FieldDescriptor) }


type FieldDescriptors interface {
	
	Len() int
	
	Get(i int) FieldDescriptor
	
	
	ByName(s Name) FieldDescriptor
	
	
	ByJSONName(s string) FieldDescriptor
	
	
	ByTextName(s string) FieldDescriptor
	
	
	ByNumber(n FieldNumber) FieldDescriptor

	doNotImplement
}



type OneofDescriptor interface {
	Descriptor

	
	
	
	IsSynthetic() bool

	
	Fields() FieldDescriptors

	isOneofDescriptor
}
type isOneofDescriptor interface{ ProtoType(OneofDescriptor) }


type OneofDescriptors interface {
	
	Len() int
	
	Get(i int) OneofDescriptor
	
	
	ByName(s Name) OneofDescriptor

	doNotImplement
}


type ExtensionDescriptor = FieldDescriptor


type ExtensionTypeDescriptor interface {
	ExtensionDescriptor

	
	Type() ExtensionType

	
	
	Descriptor() ExtensionDescriptor
}


type ExtensionDescriptors interface {
	
	Len() int
	
	Get(i int) ExtensionDescriptor
	
	
	ByName(s Name) ExtensionDescriptor

	doNotImplement
}























type ExtensionType interface {
	
	
	New() Value

	
	
	
	Zero() Value

	
	TypeDescriptor() ExtensionTypeDescriptor

	
	
	
	
	
	ValueOf(any) Value

	
	
	
	
	
	
	
	InterfaceOf(Value) any

	
	IsValidValue(Value) bool

	
	IsValidInterface(any) bool
}






type EnumDescriptor interface {
	Descriptor

	
	Values() EnumValueDescriptors

	
	ReservedNames() Names
	
	ReservedRanges() EnumRanges

	
	
	
	
	IsClosed() bool

	isEnumDescriptor
}
type isEnumDescriptor interface{ ProtoType(EnumDescriptor) }


type EnumType interface {
	
	New(n EnumNumber) Enum

	
	
	
	Descriptor() EnumDescriptor
}


type EnumDescriptors interface {
	
	Len() int
	
	Get(i int) EnumDescriptor
	
	
	ByName(s Name) EnumDescriptor

	doNotImplement
}









type EnumValueDescriptor interface {
	Descriptor

	
	Number() EnumNumber

	isEnumValueDescriptor
}
type isEnumValueDescriptor interface{ ProtoType(EnumValueDescriptor) }


type EnumValueDescriptors interface {
	
	Len() int
	
	Get(i int) EnumValueDescriptor
	
	
	ByName(s Name) EnumValueDescriptor
	
	
	
	ByNumber(n EnumNumber) EnumValueDescriptor

	doNotImplement
}





type ServiceDescriptor interface {
	Descriptor

	
	Methods() MethodDescriptors

	isServiceDescriptor
}
type isServiceDescriptor interface{ ProtoType(ServiceDescriptor) }


type ServiceDescriptors interface {
	
	Len() int
	
	Get(i int) ServiceDescriptor
	
	
	ByName(s Name) ServiceDescriptor

	doNotImplement
}



type MethodDescriptor interface {
	Descriptor

	
	Input() MessageDescriptor
	
	Output() MessageDescriptor
	
	IsStreamingClient() bool
	
	IsStreamingServer() bool

	isMethodDescriptor
}
type isMethodDescriptor interface{ ProtoType(MethodDescriptor) }


type MethodDescriptors interface {
	
	Len() int
	
	Get(i int) MethodDescriptor
	
	
	ByName(s Name) MethodDescriptor

	doNotImplement
}
