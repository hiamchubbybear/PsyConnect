



package proto

import (
	"google.golang.org/protobuf/internal/errors"
	"google.golang.org/protobuf/reflect/protoreflect"
)













type Message = protoreflect.ProtoMessage







var Error error

func init() {
	Error = errors.Error
}



func MessageName(m Message) protoreflect.FullName {
	if m == nil {
		return ""
	}
	return m.ProtoReflect().Descriptor().FullName()
}
