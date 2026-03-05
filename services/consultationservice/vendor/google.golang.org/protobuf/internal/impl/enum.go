



package impl

import (
	"reflect"

	"google.golang.org/protobuf/reflect/protoreflect"
)

type EnumInfo struct {
	GoReflectType reflect.Type 
	Desc          protoreflect.EnumDescriptor
}

func (t *EnumInfo) New(n protoreflect.EnumNumber) protoreflect.Enum {
	return reflect.ValueOf(n).Convert(t.GoReflectType).Interface().(protoreflect.Enum)
}
func (t *EnumInfo) Descriptor() protoreflect.EnumDescriptor { return t.Desc }
