



package impl

import (
	"fmt"
	"reflect"
	"strconv"

	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/internal/errors"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoiface"
)



type Export struct{}



func (Export) NewError(f string, x ...any) error {
	return errors.New(f, x...)
}



type enum = any



func (Export) EnumOf(e enum) protoreflect.Enum {
	switch e := e.(type) {
	case nil:
		return nil
	case protoreflect.Enum:
		return e
	default:
		return legacyWrapEnum(reflect.ValueOf(e))
	}
}



func (Export) EnumDescriptorOf(e enum) protoreflect.EnumDescriptor {
	switch e := e.(type) {
	case nil:
		return nil
	case protoreflect.Enum:
		return e.Descriptor()
	default:
		return LegacyLoadEnumDesc(reflect.TypeOf(e))
	}
}



func (Export) EnumTypeOf(e enum) protoreflect.EnumType {
	switch e := e.(type) {
	case nil:
		return nil
	case protoreflect.Enum:
		return e.Type()
	default:
		return legacyLoadEnumType(reflect.TypeOf(e))
	}
}



func (Export) EnumStringOf(ed protoreflect.EnumDescriptor, n protoreflect.EnumNumber) string {
	ev := ed.Values().ByNumber(n)
	if ev != nil {
		return string(ev.Name())
	}
	return strconv.Itoa(int(n))
}



type message = any


type legacyMessageWrapper struct{ m protoreflect.ProtoMessage }

func (m legacyMessageWrapper) Reset()         { proto.Reset(m.m) }
func (m legacyMessageWrapper) String() string { return Export{}.MessageStringOf(m.m) }
func (m legacyMessageWrapper) ProtoMessage()  {}



func (Export) ProtoMessageV1Of(m message) protoiface.MessageV1 {
	switch mv := m.(type) {
	case nil:
		return nil
	case protoiface.MessageV1:
		return mv
	case unwrapper:
		return Export{}.ProtoMessageV1Of(mv.protoUnwrap())
	case protoreflect.ProtoMessage:
		return legacyMessageWrapper{mv}
	default:
		panic(fmt.Sprintf("message %T is neither a v1 or v2 Message", m))
	}
}

func (Export) protoMessageV2Of(m message) protoreflect.ProtoMessage {
	switch mv := m.(type) {
	case nil:
		return nil
	case protoreflect.ProtoMessage:
		return mv
	case legacyMessageWrapper:
		return mv.m
	case protoiface.MessageV1:
		return nil
	default:
		panic(fmt.Sprintf("message %T is neither a v1 or v2 Message", m))
	}
}



func (Export) ProtoMessageV2Of(m message) protoreflect.ProtoMessage {
	if m == nil {
		return nil
	}
	if mv := (Export{}).protoMessageV2Of(m); mv != nil {
		return mv
	}
	return legacyWrapMessage(reflect.ValueOf(m)).Interface()
}



func (Export) MessageOf(m message) protoreflect.Message {
	if m == nil {
		return nil
	}
	if mv := (Export{}).protoMessageV2Of(m); mv != nil {
		return mv.ProtoReflect()
	}
	return legacyWrapMessage(reflect.ValueOf(m))
}



func (Export) MessageDescriptorOf(m message) protoreflect.MessageDescriptor {
	if m == nil {
		return nil
	}
	if mv := (Export{}).protoMessageV2Of(m); mv != nil {
		return mv.ProtoReflect().Descriptor()
	}
	return LegacyLoadMessageDesc(reflect.TypeOf(m))
}



func (Export) MessageTypeOf(m message) protoreflect.MessageType {
	if m == nil {
		return nil
	}
	if mv := (Export{}).protoMessageV2Of(m); mv != nil {
		return mv.ProtoReflect().Type()
	}
	return legacyLoadMessageType(reflect.TypeOf(m), "")
}



func (Export) MessageStringOf(m protoreflect.ProtoMessage) string {
	return prototext.MarshalOptions{Multiline: false}.Format(m)
}
