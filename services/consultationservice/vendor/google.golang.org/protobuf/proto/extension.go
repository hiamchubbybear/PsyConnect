



package proto

import (
	"google.golang.org/protobuf/reflect/protoreflect"
)



func HasExtension(m Message, xt protoreflect.ExtensionType) bool {
	
	
	if m == nil || xt == nil {
		return false
	}

	
	
	mr := m.ProtoReflect()
	xd := xt.TypeDescriptor()
	if mr.Descriptor() != xd.ContainingMessage() {
		return false
	}

	return mr.Has(xd)
}




func ClearExtension(m Message, xt protoreflect.ExtensionType) {
	m.ProtoReflect().Clear(xt.TypeDescriptor())
}















































func GetExtension(m Message, xt protoreflect.ExtensionType) any {
	
	if m == nil {
		return xt.InterfaceOf(xt.Zero())
	}

	return xt.InterfaceOf(m.ProtoReflect().Get(xt.TypeDescriptor()))
}

































func SetExtension(m Message, xt protoreflect.ExtensionType, v any) {
	xd := xt.TypeDescriptor()
	pv := xt.ValueOf(v)

	
	isValid := true
	switch {
	case xd.IsList():
		isValid = pv.List().IsValid()
	case xd.IsMap():
		isValid = pv.Map().IsValid()
	case xd.Message() != nil:
		isValid = pv.Message().IsValid()
	}
	if !isValid {
		m.ProtoReflect().Clear(xd)
		return
	}

	m.ProtoReflect().Set(xd, pv)
}






func RangeExtensions(m Message, f func(protoreflect.ExtensionType, any) bool) {
	
	if m == nil {
		return
	}

	m.ProtoReflect().Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		if fd.IsExtension() {
			xt := fd.(protoreflect.ExtensionTypeDescriptor).Type()
			vi := xt.InterfaceOf(v)
			return f(xt, vi)
		}
		return true
	})
}
