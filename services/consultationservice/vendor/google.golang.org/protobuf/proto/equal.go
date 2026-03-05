



package proto

import (
	"reflect"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoiface"
)





























func Equal(x, y Message) bool {
	if x == nil || y == nil {
		return x == nil && y == nil
	}
	if reflect.TypeOf(x).Kind() == reflect.Ptr && x == y {
		
		return true
	}
	mx := x.ProtoReflect()
	my := y.ProtoReflect()
	if mx.IsValid() != my.IsValid() {
		return false
	}

	
	pmx := protoMethods(mx)
	pmy := protoMethods(my)
	if pmx != nil && pmy != nil && pmx.Equal != nil && pmy.Equal != nil {
		return pmx.Equal(protoiface.EqualInput{MessageA: mx, MessageB: my}).Equal
	}

	vx := protoreflect.ValueOfMessage(mx)
	vy := protoreflect.ValueOfMessage(my)
	return vx.Equal(vy)
}
