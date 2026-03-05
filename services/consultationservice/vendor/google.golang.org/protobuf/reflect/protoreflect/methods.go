



package protoreflect

import (
	"google.golang.org/protobuf/internal/pragma"
)







type (
	methods = struct {
		pragma.NoUnkeyedLiterals
		Flags            supportFlags
		Size             func(sizeInput) sizeOutput
		Marshal          func(marshalInput) (marshalOutput, error)
		Unmarshal        func(unmarshalInput) (unmarshalOutput, error)
		Merge            func(mergeInput) mergeOutput
		CheckInitialized func(checkInitializedInput) (checkInitializedOutput, error)
		Equal            func(equalInput) equalOutput
	}
	supportFlags = uint64
	sizeInput    = struct {
		pragma.NoUnkeyedLiterals
		Message Message
		Flags   uint8
	}
	sizeOutput = struct {
		pragma.NoUnkeyedLiterals
		Size int
	}
	marshalInput = struct {
		pragma.NoUnkeyedLiterals
		Message Message
		Buf     []byte
		Flags   uint8
	}
	marshalOutput = struct {
		pragma.NoUnkeyedLiterals
		Buf []byte
	}
	unmarshalInput = struct {
		pragma.NoUnkeyedLiterals
		Message  Message
		Buf      []byte
		Flags    uint8
		Resolver interface {
			FindExtensionByName(field FullName) (ExtensionType, error)
			FindExtensionByNumber(message FullName, field FieldNumber) (ExtensionType, error)
		}
		Depth int
	}
	unmarshalOutput = struct {
		pragma.NoUnkeyedLiterals
		Flags uint8
	}
	mergeInput = struct {
		pragma.NoUnkeyedLiterals
		Source      Message
		Destination Message
	}
	mergeOutput = struct {
		pragma.NoUnkeyedLiterals
		Flags uint8
	}
	checkInitializedInput = struct {
		pragma.NoUnkeyedLiterals
		Message Message
	}
	checkInitializedOutput = struct {
		pragma.NoUnkeyedLiterals
	}
	equalInput = struct {
		pragma.NoUnkeyedLiterals
		MessageA Message
		MessageB Message
	}
	equalOutput = struct {
		pragma.NoUnkeyedLiterals
		Equal bool
	}
)
