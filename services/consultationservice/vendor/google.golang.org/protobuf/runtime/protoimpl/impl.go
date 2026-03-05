










package protoimpl

import (
	"google.golang.org/protobuf/internal/filedesc"
	"google.golang.org/protobuf/internal/filetype"
	"google.golang.org/protobuf/internal/impl"
	"google.golang.org/protobuf/internal/protolazy"
)


const UnsafeEnabled = impl.UnsafeEnabled

type (
	
	DescBuilder = filedesc.Builder
	TypeBuilder = filetype.Builder

	
	EnumInfo      = impl.EnumInfo
	MessageInfo   = impl.MessageInfo
	ExtensionInfo = impl.ExtensionInfo

	
	MessageState     = impl.MessageState
	SizeCache        = impl.SizeCache
	WeakFields       = impl.WeakFields
	UnknownFields    = impl.UnknownFields
	ExtensionFields  = impl.ExtensionFields
	ExtensionFieldV1 = impl.ExtensionField

	Pointer = impl.Pointer

	LazyUnmarshalInfo  = *protolazy.XXX_lazyUnmarshalInfo
	RaceDetectHookData = impl.RaceDetectHookData
)

var X impl.Export
