



package proto

import (
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/internal/encoding/messageset"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoiface"
)






func Size(m Message) int {
	return MarshalOptions{}.Size(m)
}






func (o MarshalOptions) Size(m Message) int {
	
	if m == nil {
		return 0
	}

	return o.size(m.ProtoReflect())
}




func (o MarshalOptions) size(m protoreflect.Message) (size int) {
	methods := protoMethods(m)
	if methods != nil && methods.Size != nil {
		out := methods.Size(protoiface.SizeInput{
			Message: m,
			Flags:   o.flags(),
		})
		return out.Size
	}
	if methods != nil && methods.Marshal != nil {
		
		
		out, _ := methods.Marshal(protoiface.MarshalInput{
			Message: m,
			Flags:   o.flags(),
		})
		return len(out.Buf)
	}
	return o.sizeMessageSlow(m)
}

func (o MarshalOptions) sizeMessageSlow(m protoreflect.Message) (size int) {
	if messageset.IsMessageSet(m.Descriptor()) {
		return o.sizeMessageSet(m)
	}
	m.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		size += o.sizeField(fd, v)
		return true
	})
	size += len(m.GetUnknown())
	return size
}

func (o MarshalOptions) sizeField(fd protoreflect.FieldDescriptor, value protoreflect.Value) (size int) {
	num := fd.Number()
	switch {
	case fd.IsList():
		return o.sizeList(num, fd, value.List())
	case fd.IsMap():
		return o.sizeMap(num, fd, value.Map())
	default:
		return protowire.SizeTag(num) + o.sizeSingular(num, fd.Kind(), value)
	}
}

func (o MarshalOptions) sizeList(num protowire.Number, fd protoreflect.FieldDescriptor, list protoreflect.List) (size int) {
	sizeTag := protowire.SizeTag(num)

	if fd.IsPacked() && list.Len() > 0 {
		content := 0
		for i, llen := 0, list.Len(); i < llen; i++ {
			content += o.sizeSingular(num, fd.Kind(), list.Get(i))
		}
		return sizeTag + protowire.SizeBytes(content)
	}

	for i, llen := 0, list.Len(); i < llen; i++ {
		size += sizeTag + o.sizeSingular(num, fd.Kind(), list.Get(i))
	}
	return size
}

func (o MarshalOptions) sizeMap(num protowire.Number, fd protoreflect.FieldDescriptor, mapv protoreflect.Map) (size int) {
	sizeTag := protowire.SizeTag(num)

	mapv.Range(func(key protoreflect.MapKey, value protoreflect.Value) bool {
		size += sizeTag
		size += protowire.SizeBytes(o.sizeField(fd.MapKey(), key.Value()) + o.sizeField(fd.MapValue(), value))
		return true
	})
	return size
}
