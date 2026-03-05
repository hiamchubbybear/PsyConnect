



package protoiface

type MessageV1 interface {
	Reset()
	String() string
	ProtoMessage()
}

type ExtensionRangeV1 struct {
	Start, End int32 
}
