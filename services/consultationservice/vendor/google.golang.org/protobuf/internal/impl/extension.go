



package impl

import (
	"reflect"
	"sync"
	"sync/atomic"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoiface"
)






type ExtensionInfo struct {
	
	
	
	
	
	
	
	
	
	
	
	
	
	init uint32
	mu   sync.Mutex

	goType reflect.Type
	desc   extensionTypeDescriptor
	conv   Converter
	info   *extensionFieldInfo 

	
	
	
	
	
	ExtendedType protoiface.MessageV1

	
	
	
	
	
	
	ExtensionType any

	
	
	
	Field int32

	
	
	
	Name string

	
	
	
	Tag string

	
	
	
	Filename string
}


const (
	extensionInfoUninitialized = 0
	extensionInfoDescInit      = 1
	extensionInfoFullInit      = 2
)

func InitExtensionInfo(xi *ExtensionInfo, xd protoreflect.ExtensionDescriptor, goType reflect.Type) {
	xi.goType = goType
	xi.desc = extensionTypeDescriptor{xd, xi}
	xi.init = extensionInfoDescInit
}

func (xi *ExtensionInfo) New() protoreflect.Value {
	return xi.lazyInit().New()
}
func (xi *ExtensionInfo) Zero() protoreflect.Value {
	return xi.lazyInit().Zero()
}
func (xi *ExtensionInfo) ValueOf(v any) protoreflect.Value {
	return xi.lazyInit().PBValueOf(reflect.ValueOf(v))
}
func (xi *ExtensionInfo) InterfaceOf(v protoreflect.Value) any {
	return xi.lazyInit().GoValueOf(v).Interface()
}
func (xi *ExtensionInfo) IsValidValue(v protoreflect.Value) bool {
	return xi.lazyInit().IsValidPB(v)
}
func (xi *ExtensionInfo) IsValidInterface(v any) bool {
	return xi.lazyInit().IsValidGo(reflect.ValueOf(v))
}
func (xi *ExtensionInfo) TypeDescriptor() protoreflect.ExtensionTypeDescriptor {
	if atomic.LoadUint32(&xi.init) < extensionInfoDescInit {
		xi.lazyInitSlow()
	}
	return &xi.desc
}

func (xi *ExtensionInfo) lazyInit() Converter {
	if atomic.LoadUint32(&xi.init) < extensionInfoFullInit {
		xi.lazyInitSlow()
	}
	return xi.conv
}

func (xi *ExtensionInfo) lazyInitSlow() {
	xi.mu.Lock()
	defer xi.mu.Unlock()

	if xi.init == extensionInfoFullInit {
		return
	}
	defer atomic.StoreUint32(&xi.init, extensionInfoFullInit)

	if xi.desc.ExtensionDescriptor == nil {
		xi.initFromLegacy()
	}
	if !xi.desc.ExtensionDescriptor.IsPlaceholder() {
		if xi.ExtensionType == nil {
			xi.initToLegacy()
		}
		xi.conv = NewConverter(xi.goType, xi.desc.ExtensionDescriptor)
		xi.info = makeExtensionFieldInfo(xi.desc.ExtensionDescriptor)
		xi.info.validation = newValidationInfo(xi.desc.ExtensionDescriptor, xi.goType)
	}
}

type extensionTypeDescriptor struct {
	protoreflect.ExtensionDescriptor
	xi *ExtensionInfo
}

func (xtd *extensionTypeDescriptor) Type() protoreflect.ExtensionType {
	return xtd.xi
}
func (xtd *extensionTypeDescriptor) Descriptor() protoreflect.ExtensionDescriptor {
	return xtd.ExtensionDescriptor
}
