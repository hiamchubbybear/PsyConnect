



package impl

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"google.golang.org/protobuf/internal/genid"
	"google.golang.org/protobuf/reflect/protoreflect"
)







type MessageInfo struct {
	
	GoReflectType reflect.Type 

	
	Desc protoreflect.MessageDescriptor

	
	
	Exporter exporter

	
	OneofWrappers []any

	initMu   sync.Mutex 
	initDone uint32

	reflectMessageInfo 
	coderMessageInfo   
}




type exporter func(v any, i int) any




func getMessageInfo(mt reflect.Type) *MessageInfo {
	m, ok := reflect.Zero(mt).Interface().(protoreflect.ProtoMessage)
	if !ok {
		return nil
	}
	mr, ok := m.ProtoReflect().(interface{ ProtoMessageInfo() *MessageInfo })
	if !ok {
		return nil
	}
	return mr.ProtoMessageInfo()
}

func (mi *MessageInfo) init() {
	
	
	
	if atomic.LoadUint32(&mi.initDone) == 0 {
		mi.initOnce()
	}
}

func (mi *MessageInfo) initOnce() {
	mi.initMu.Lock()
	defer mi.initMu.Unlock()
	if mi.initDone == 1 {
		return
	}
	if opaqueInitHook(mi) {
		return
	}

	t := mi.GoReflectType
	if t.Kind() != reflect.Ptr && t.Elem().Kind() != reflect.Struct {
		panic(fmt.Sprintf("got %v, want *struct kind", t))
	}
	t = t.Elem()

	si := mi.makeStructInfo(t)
	mi.makeReflectFuncs(t, si)
	mi.makeCoderMethods(t, si)

	atomic.StoreUint32(&mi.initDone, 1)
}




func (mi *MessageInfo) getPointer(m protoreflect.Message) (p pointer, ok bool) {
	switch m := m.(type) {
	case *messageState:
		return m.pointer(), m.messageInfo() == mi
	case *messageReflectWrapper:
		return m.pointer(), m.messageInfo() == mi
	}
	return pointer{}, false
}

type (
	SizeCache       = int32
	WeakFields      = map[int32]protoreflect.ProtoMessage
	UnknownFields   = unknownFieldsA 
	unknownFieldsA  = []byte
	unknownFieldsB  = *[]byte
	ExtensionFields = map[int32]ExtensionField
)

var (
	sizecacheType       = reflect.TypeOf(SizeCache(0))
	unknownFieldsAType  = reflect.TypeOf(unknownFieldsA(nil))
	unknownFieldsBType  = reflect.TypeOf(unknownFieldsB(nil))
	extensionFieldsType = reflect.TypeOf(ExtensionFields(nil))
)

type structInfo struct {
	sizecacheOffset offset
	sizecacheType   reflect.Type
	unknownOffset   offset
	unknownType     reflect.Type
	extensionOffset offset
	extensionType   reflect.Type

	lazyOffset     offset
	presenceOffset offset

	fieldsByNumber        map[protoreflect.FieldNumber]reflect.StructField
	oneofsByName          map[protoreflect.Name]reflect.StructField
	oneofWrappersByType   map[reflect.Type]protoreflect.FieldNumber
	oneofWrappersByNumber map[protoreflect.FieldNumber]reflect.Type
}

func (mi *MessageInfo) makeStructInfo(t reflect.Type) structInfo {
	si := structInfo{
		sizecacheOffset: invalidOffset,
		unknownOffset:   invalidOffset,
		extensionOffset: invalidOffset,
		lazyOffset:      invalidOffset,
		presenceOffset:  invalidOffset,

		fieldsByNumber:        map[protoreflect.FieldNumber]reflect.StructField{},
		oneofsByName:          map[protoreflect.Name]reflect.StructField{},
		oneofWrappersByType:   map[reflect.Type]protoreflect.FieldNumber{},
		oneofWrappersByNumber: map[protoreflect.FieldNumber]reflect.Type{},
	}

fieldLoop:
	for i := 0; i < t.NumField(); i++ {
		switch f := t.Field(i); f.Name {
		case genid.SizeCache_goname, genid.SizeCacheA_goname:
			if f.Type == sizecacheType {
				si.sizecacheOffset = offsetOf(f)
				si.sizecacheType = f.Type
			}
		case genid.UnknownFields_goname, genid.UnknownFieldsA_goname:
			if f.Type == unknownFieldsAType || f.Type == unknownFieldsBType {
				si.unknownOffset = offsetOf(f)
				si.unknownType = f.Type
			}
		case genid.ExtensionFields_goname, genid.ExtensionFieldsA_goname, genid.ExtensionFieldsB_goname:
			if f.Type == extensionFieldsType {
				si.extensionOffset = offsetOf(f)
				si.extensionType = f.Type
			}
		case "lazyFields", "XXX_lazyUnmarshalInfo":
			si.lazyOffset = offsetOf(f)
		case "XXX_presence":
			si.presenceOffset = offsetOf(f)
		default:
			for _, s := range strings.Split(f.Tag.Get("protobuf"), ",") {
				if len(s) > 0 && strings.Trim(s, "0123456789") == "" {
					n, _ := strconv.ParseUint(s, 10, 64)
					si.fieldsByNumber[protoreflect.FieldNumber(n)] = f
					continue fieldLoop
				}
			}
			if s := f.Tag.Get("protobuf_oneof"); len(s) > 0 {
				si.oneofsByName[protoreflect.Name(s)] = f
				continue fieldLoop
			}
		}
	}

	
	oneofWrappers := mi.OneofWrappers
	methods := make([]reflect.Method, 0, 2)
	if m, ok := reflect.PtrTo(t).MethodByName("XXX_OneofFuncs"); ok {
		methods = append(methods, m)
	}
	if m, ok := reflect.PtrTo(t).MethodByName("XXX_OneofWrappers"); ok {
		methods = append(methods, m)
	}
	for _, fn := range methods {
		for _, v := range fn.Func.Call([]reflect.Value{reflect.Zero(fn.Type.In(0))}) {
			if vs, ok := v.Interface().([]any); ok {
				oneofWrappers = vs
			}
		}
	}
	for _, v := range oneofWrappers {
		tf := reflect.TypeOf(v).Elem()
		f := tf.Field(0)
		for _, s := range strings.Split(f.Tag.Get("protobuf"), ",") {
			if len(s) > 0 && strings.Trim(s, "0123456789") == "" {
				n, _ := strconv.ParseUint(s, 10, 64)
				si.oneofWrappersByType[tf] = protoreflect.FieldNumber(n)
				si.oneofWrappersByNumber[protoreflect.FieldNumber(n)] = tf
				break
			}
		}
	}

	return si
}

func (mi *MessageInfo) New() protoreflect.Message {
	m := reflect.New(mi.GoReflectType.Elem()).Interface()
	if r, ok := m.(protoreflect.ProtoMessage); ok {
		return r.ProtoReflect()
	}
	return mi.MessageOf(m)
}
func (mi *MessageInfo) Zero() protoreflect.Message {
	return mi.MessageOf(reflect.Zero(mi.GoReflectType).Interface())
}
func (mi *MessageInfo) Descriptor() protoreflect.MessageDescriptor {
	return mi.Desc
}
func (mi *MessageInfo) Enum(i int) protoreflect.EnumType {
	mi.init()
	fd := mi.Desc.Fields().Get(i)
	return Export{}.EnumTypeOf(mi.fieldTypes[fd.Number()])
}
func (mi *MessageInfo) Message(i int) protoreflect.MessageType {
	mi.init()
	fd := mi.Desc.Fields().Get(i)
	switch {
	case fd.IsMap():
		return mapEntryType{fd.Message(), mi.fieldTypes[fd.Number()]}
	default:
		return Export{}.MessageTypeOf(mi.fieldTypes[fd.Number()])
	}
}

type mapEntryType struct {
	desc    protoreflect.MessageDescriptor
	valType any 
}

func (mt mapEntryType) New() protoreflect.Message {
	return nil
}
func (mt mapEntryType) Zero() protoreflect.Message {
	return nil
}
func (mt mapEntryType) Descriptor() protoreflect.MessageDescriptor {
	return mt.desc
}
func (mt mapEntryType) Enum(i int) protoreflect.EnumType {
	fd := mt.desc.Fields().Get(i)
	if fd.Enum() == nil {
		return nil
	}
	return Export{}.EnumTypeOf(mt.valType)
}
func (mt mapEntryType) Message(i int) protoreflect.MessageType {
	fd := mt.desc.Fields().Get(i)
	if fd.Message() == nil {
		return nil
	}
	return Export{}.MessageTypeOf(mt.valType)
}
