








package protoiface

import (
	"google.golang.org/protobuf/internal/pragma"
	"google.golang.org/protobuf/reflect/protoreflect"
)


type Methods = struct {
	pragma.NoUnkeyedLiterals

	
	Flags SupportFlags

	
	
	Size func(SizeInput) SizeOutput

	
	
	
	Marshal func(MarshalInput) (MarshalOutput, error)

	
	
	Unmarshal func(UnmarshalInput) (UnmarshalOutput, error)

	
	Merge func(MergeInput) MergeOutput

	
	CheckInitialized func(CheckInitializedInput) (CheckInitializedOutput, error)

	
	Equal func(EqualInput) EqualOutput
}


type SupportFlags = uint64

const (
	
	SupportMarshalDeterministic SupportFlags = 1 << iota

	
	SupportUnmarshalDiscardUnknown
)


type SizeInput = struct {
	pragma.NoUnkeyedLiterals

	Message protoreflect.Message
	Flags   MarshalInputFlags
}


type SizeOutput = struct {
	pragma.NoUnkeyedLiterals

	Size int
}


type MarshalInput = struct {
	pragma.NoUnkeyedLiterals

	Message protoreflect.Message
	Buf     []byte 
	Flags   MarshalInputFlags
}


type MarshalOutput = struct {
	pragma.NoUnkeyedLiterals

	Buf []byte 
}



type MarshalInputFlags = uint8

const (
	MarshalDeterministic MarshalInputFlags = 1 << iota
	MarshalUseCachedSize
)


type UnmarshalInput = struct {
	pragma.NoUnkeyedLiterals

	Message  protoreflect.Message
	Buf      []byte 
	Flags    UnmarshalInputFlags
	Resolver interface {
		FindExtensionByName(field protoreflect.FullName) (protoreflect.ExtensionType, error)
		FindExtensionByNumber(message protoreflect.FullName, field protoreflect.FieldNumber) (protoreflect.ExtensionType, error)
	}
	Depth int
}


type UnmarshalOutput = struct {
	pragma.NoUnkeyedLiterals

	Flags UnmarshalOutputFlags
}



type UnmarshalInputFlags = uint8

const (
	UnmarshalDiscardUnknown UnmarshalInputFlags = 1 << iota

	
	
	UnmarshalAliasBuffer

	
	
	UnmarshalValidated

	
	
	UnmarshalCheckRequired

	
	
	UnmarshalNoLazyDecoding
)


type UnmarshalOutputFlags = uint8

const (
	
	
	
	UnmarshalInitialized UnmarshalOutputFlags = 1 << iota
)


type MergeInput = struct {
	pragma.NoUnkeyedLiterals

	Source      protoreflect.Message
	Destination protoreflect.Message
}


type MergeOutput = struct {
	pragma.NoUnkeyedLiterals

	Flags MergeOutputFlags
}


type MergeOutputFlags = uint8

const (
	
	
	MergeComplete MergeOutputFlags = 1 << iota
)


type CheckInitializedInput = struct {
	pragma.NoUnkeyedLiterals

	Message protoreflect.Message
}


type CheckInitializedOutput = struct {
	pragma.NoUnkeyedLiterals
}


type EqualInput = struct {
	pragma.NoUnkeyedLiterals

	MessageA protoreflect.Message
	MessageB protoreflect.Message
}


type EqualOutput = struct {
	pragma.NoUnkeyedLiterals

	Equal bool
}
