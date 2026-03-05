



package proto

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoiface"
)












func Merge(dst, src Message) {
	
	

	dstMsg, srcMsg := dst.ProtoReflect(), src.ProtoReflect()
	if dstMsg.Descriptor() != srcMsg.Descriptor() {
		if got, want := dstMsg.Descriptor().FullName(), srcMsg.Descriptor().FullName(); got != want {
			panic(fmt.Sprintf("descriptor mismatch: %v != %v", got, want))
		}
		panic("descriptor mismatch")
	}
	mergeOptions{}.mergeMessage(dstMsg, srcMsg)
}



func Clone(m Message) Message {
	
	
	
	
	
	
	
	
	if m == nil {
		return nil
	}
	src := m.ProtoReflect()
	if !src.IsValid() {
		return src.Type().Zero().Interface()
	}
	dst := src.New()
	mergeOptions{}.mergeMessage(dst, src)
	return dst.Interface()
}



func CloneOf[M Message](m M) M {
	return Clone(m).(M)
}



type mergeOptions struct{}

func (o mergeOptions) mergeMessage(dst, src protoreflect.Message) {
	methods := protoMethods(dst)
	if methods != nil && methods.Merge != nil {
		in := protoiface.MergeInput{
			Destination: dst,
			Source:      src,
		}
		out := methods.Merge(in)
		if out.Flags&protoiface.MergeComplete != 0 {
			return
		}
	}

	if !dst.IsValid() {
		panic(fmt.Sprintf("cannot merge into invalid %v message", dst.Descriptor().FullName()))
	}

	src.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		switch {
		case fd.IsList():
			o.mergeList(dst.Mutable(fd).List(), v.List(), fd)
		case fd.IsMap():
			o.mergeMap(dst.Mutable(fd).Map(), v.Map(), fd.MapValue())
		case fd.Message() != nil:
			o.mergeMessage(dst.Mutable(fd).Message(), v.Message())
		case fd.Kind() == protoreflect.BytesKind:
			dst.Set(fd, o.cloneBytes(v))
		default:
			dst.Set(fd, v)
		}
		return true
	})

	if len(src.GetUnknown()) > 0 {
		dst.SetUnknown(append(dst.GetUnknown(), src.GetUnknown()...))
	}
}

func (o mergeOptions) mergeList(dst, src protoreflect.List, fd protoreflect.FieldDescriptor) {
	
	for i, n := 0, src.Len(); i < n; i++ {
		switch v := src.Get(i); {
		case fd.Message() != nil:
			dstv := dst.NewElement()
			o.mergeMessage(dstv.Message(), v.Message())
			dst.Append(dstv)
		case fd.Kind() == protoreflect.BytesKind:
			dst.Append(o.cloneBytes(v))
		default:
			dst.Append(v)
		}
	}
}

func (o mergeOptions) mergeMap(dst, src protoreflect.Map, fd protoreflect.FieldDescriptor) {
	
	src.Range(func(k protoreflect.MapKey, v protoreflect.Value) bool {
		switch {
		case fd.Message() != nil:
			dstv := dst.NewValue()
			o.mergeMessage(dstv.Message(), v.Message())
			dst.Set(k, dstv)
		case fd.Kind() == protoreflect.BytesKind:
			dst.Set(k, o.cloneBytes(v))
		default:
			dst.Set(k, v)
		}
		return true
	})
}

func (o mergeOptions) cloneBytes(v protoreflect.Value) protoreflect.Value {
	return protoreflect.ValueOfBytes(append([]byte{}, v.Bytes()...))
}
