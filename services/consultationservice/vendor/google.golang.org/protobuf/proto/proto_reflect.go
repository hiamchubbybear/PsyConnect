




//go:build protoreflect
// +build protoreflect

package proto

import (
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoiface"
)

const hasProtoMethods = false

func protoMethods(m protoreflect.Message) *protoiface.Methods {
	return nil
}
