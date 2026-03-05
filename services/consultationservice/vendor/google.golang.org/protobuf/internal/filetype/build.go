





package filetype

import (
	"reflect"

	"google.golang.org/protobuf/internal/descopts"
	"google.golang.org/protobuf/internal/filedesc"
	pimpl "google.golang.org/protobuf/internal/impl"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)


































type Builder struct {
	
	File filedesc.Builder

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	GoTypes []any

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	DependencyIndexes []int32

	
	EnumInfos []pimpl.EnumInfo

	
	
	
	
	MessageInfos []pimpl.MessageInfo

	
	
	
	
	ExtensionInfos []pimpl.ExtensionInfo

	
	
	TypeRegistry interface {
		RegisterMessage(protoreflect.MessageType) error
		RegisterEnum(protoreflect.EnumType) error
		RegisterExtension(protoreflect.ExtensionType) error
	}
}


type Out struct {
	File protoreflect.FileDescriptor
}

func (tb Builder) Build() (out Out) {
	
	
	if tb.File.FileRegistry == nil {
		tb.File.FileRegistry = protoregistry.GlobalFiles
	}
	tb.File.FileRegistry = &resolverByIndex{
		goTypes:      tb.GoTypes,
		depIdxs:      tb.DependencyIndexes,
		fileRegistry: tb.File.FileRegistry,
	}

	
	if tb.TypeRegistry == nil {
		tb.TypeRegistry = protoregistry.GlobalTypes
	}

	fbOut := tb.File.Build()
	out.File = fbOut.File

	
	enumGoTypes := tb.GoTypes[:len(fbOut.Enums)]
	if len(tb.EnumInfos) != len(fbOut.Enums) {
		panic("mismatching enum lengths")
	}
	if len(fbOut.Enums) > 0 {
		for i := range fbOut.Enums {
			tb.EnumInfos[i] = pimpl.EnumInfo{
				GoReflectType: reflect.TypeOf(enumGoTypes[i]),
				Desc:          &fbOut.Enums[i],
			}
			
			if err := tb.TypeRegistry.RegisterEnum(&tb.EnumInfos[i]); err != nil {
				panic(err)
			}
		}
	}

	
	messageGoTypes := tb.GoTypes[len(fbOut.Enums):][:len(fbOut.Messages)]
	if len(tb.MessageInfos) != len(fbOut.Messages) {
		panic("mismatching message lengths")
	}
	if len(fbOut.Messages) > 0 {
		for i := range fbOut.Messages {
			if messageGoTypes[i] == nil {
				continue 
			}

			tb.MessageInfos[i].GoReflectType = reflect.TypeOf(messageGoTypes[i])
			tb.MessageInfos[i].Desc = &fbOut.Messages[i]

			
			if err := tb.TypeRegistry.RegisterMessage(&tb.MessageInfos[i]); err != nil {
				panic(err)
			}
		}

		
		
		if out.File.Path() == "google/protobuf/descriptor.proto" && out.File.Package() == "google.protobuf" {
			for i := range fbOut.Messages {
				switch fbOut.Messages[i].Name() {
				case "FileOptions":
					descopts.File = messageGoTypes[i].(protoreflect.ProtoMessage)
				case "EnumOptions":
					descopts.Enum = messageGoTypes[i].(protoreflect.ProtoMessage)
				case "EnumValueOptions":
					descopts.EnumValue = messageGoTypes[i].(protoreflect.ProtoMessage)
				case "MessageOptions":
					descopts.Message = messageGoTypes[i].(protoreflect.ProtoMessage)
				case "FieldOptions":
					descopts.Field = messageGoTypes[i].(protoreflect.ProtoMessage)
				case "OneofOptions":
					descopts.Oneof = messageGoTypes[i].(protoreflect.ProtoMessage)
				case "ExtensionRangeOptions":
					descopts.ExtensionRange = messageGoTypes[i].(protoreflect.ProtoMessage)
				case "ServiceOptions":
					descopts.Service = messageGoTypes[i].(protoreflect.ProtoMessage)
				case "MethodOptions":
					descopts.Method = messageGoTypes[i].(protoreflect.ProtoMessage)
				}
			}
		}
	}

	
	if len(tb.ExtensionInfos) != len(fbOut.Extensions) {
		panic("mismatching extension lengths")
	}
	var depIdx int32
	for i := range fbOut.Extensions {
		
		
		const listExtDeps = 2
		var goType reflect.Type
		switch fbOut.Extensions[i].L1.Kind {
		case protoreflect.EnumKind:
			j := depIdxs.Get(tb.DependencyIndexes, listExtDeps, depIdx)
			goType = reflect.TypeOf(tb.GoTypes[j])
			depIdx++
		case protoreflect.MessageKind, protoreflect.GroupKind:
			j := depIdxs.Get(tb.DependencyIndexes, listExtDeps, depIdx)
			goType = reflect.TypeOf(tb.GoTypes[j])
			depIdx++
		default:
			goType = goTypeForPBKind[fbOut.Extensions[i].L1.Kind]
		}
		if fbOut.Extensions[i].IsList() {
			goType = reflect.SliceOf(goType)
		}

		pimpl.InitExtensionInfo(&tb.ExtensionInfos[i], &fbOut.Extensions[i], goType)

		
		if err := tb.TypeRegistry.RegisterExtension(&tb.ExtensionInfos[i]); err != nil {
			panic(err)
		}
	}

	return out
}

var goTypeForPBKind = map[protoreflect.Kind]reflect.Type{
	protoreflect.BoolKind:     reflect.TypeOf(bool(false)),
	protoreflect.Int32Kind:    reflect.TypeOf(int32(0)),
	protoreflect.Sint32Kind:   reflect.TypeOf(int32(0)),
	protoreflect.Sfixed32Kind: reflect.TypeOf(int32(0)),
	protoreflect.Int64Kind:    reflect.TypeOf(int64(0)),
	protoreflect.Sint64Kind:   reflect.TypeOf(int64(0)),
	protoreflect.Sfixed64Kind: reflect.TypeOf(int64(0)),
	protoreflect.Uint32Kind:   reflect.TypeOf(uint32(0)),
	protoreflect.Fixed32Kind:  reflect.TypeOf(uint32(0)),
	protoreflect.Uint64Kind:   reflect.TypeOf(uint64(0)),
	protoreflect.Fixed64Kind:  reflect.TypeOf(uint64(0)),
	protoreflect.FloatKind:    reflect.TypeOf(float32(0)),
	protoreflect.DoubleKind:   reflect.TypeOf(float64(0)),
	protoreflect.StringKind:   reflect.TypeOf(string("")),
	protoreflect.BytesKind:    reflect.TypeOf([]byte(nil)),
}

type depIdxs []int32


func (x depIdxs) Get(i, j int32) int32 {
	return x[x[int32(len(x))-i-1]+j]
}

type (
	resolverByIndex struct {
		goTypes []any
		depIdxs depIdxs
		fileRegistry
	}
	fileRegistry interface {
		FindFileByPath(string) (protoreflect.FileDescriptor, error)
		FindDescriptorByName(protoreflect.FullName) (protoreflect.Descriptor, error)
		RegisterFile(protoreflect.FileDescriptor) error
	}
)

func (r *resolverByIndex) FindEnumByIndex(i, j int32, es []filedesc.Enum, ms []filedesc.Message) protoreflect.EnumDescriptor {
	if depIdx := int(r.depIdxs.Get(i, j)); int(depIdx) < len(es)+len(ms) {
		return &es[depIdx]
	} else {
		return pimpl.Export{}.EnumDescriptorOf(r.goTypes[depIdx])
	}
}

func (r *resolverByIndex) FindMessageByIndex(i, j int32, es []filedesc.Enum, ms []filedesc.Message) protoreflect.MessageDescriptor {
	if depIdx := int(r.depIdxs.Get(i, j)); depIdx < len(es)+len(ms) {
		return &ms[depIdx-len(es)]
	} else {
		return pimpl.Export{}.MessageDescriptorOf(r.goTypes[depIdx])
	}
}
