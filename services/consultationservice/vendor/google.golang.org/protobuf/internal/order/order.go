



package order

import (
	"google.golang.org/protobuf/reflect/protoreflect"
)



type FieldOrder func(x, y protoreflect.FieldDescriptor) bool

var (
	
	AnyFieldOrder FieldOrder = nil

	
	
	LegacyFieldOrder FieldOrder = func(x, y protoreflect.FieldDescriptor) bool {
		ox, oy := x.ContainingOneof(), y.ContainingOneof()
		inOneof := func(od protoreflect.OneofDescriptor) bool {
			return od != nil && !od.IsSynthetic()
		}

		
		if x.IsExtension() != y.IsExtension() {
			return x.IsExtension() && !y.IsExtension()
		}
		
		if inOneof(ox) != inOneof(oy) {
			return !inOneof(ox) && inOneof(oy)
		}
		
		if inOneof(ox) && inOneof(oy) && ox != oy {
			return ox.Index() < oy.Index()
		}
		
		return x.Number() < y.Number()
	}

	
	NumberFieldOrder FieldOrder = func(x, y protoreflect.FieldDescriptor) bool {
		return x.Number() < y.Number()
	}

	
	
	
	IndexNameFieldOrder FieldOrder = func(x, y protoreflect.FieldDescriptor) bool {
		
		if x.IsExtension() != y.IsExtension() {
			return !x.IsExtension() && y.IsExtension()
		}
		
		if x.IsExtension() && y.IsExtension() {
			return x.FullName() < y.FullName()
		}
		
		return x.Index() < y.Index()
	}
)



type KeyOrder func(x, y protoreflect.MapKey) bool

var (
	
	AnyKeyOrder KeyOrder = nil

	
	
	GenericKeyOrder KeyOrder = func(x, y protoreflect.MapKey) bool {
		switch x.Interface().(type) {
		case bool:
			return !x.Bool() && y.Bool()
		case int32, int64:
			return x.Int() < y.Int()
		case uint32, uint64:
			return x.Uint() < y.Uint()
		case string:
			return x.String() < y.String()
		default:
			panic("invalid map key type")
		}
	}
)
