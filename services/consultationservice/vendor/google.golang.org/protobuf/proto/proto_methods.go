




//go:build !protoreflect
// +build !protoreflect

package proto

import (
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoiface"
)

const hasProtoMethods = true

func protoMethods(m protoreflect.Message) *protoiface.Methods {
	return m.ProtoMethods()
}
