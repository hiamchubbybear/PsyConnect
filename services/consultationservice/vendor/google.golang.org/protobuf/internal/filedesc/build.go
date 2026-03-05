







package filedesc

import (
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/internal/genid"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)


type Builder struct {
	
	GoPackagePath string

	
	
	RawDescriptor []byte

	
	NumEnums int32
	
	
	NumMessages int32
	
	NumExtensions int32
	
	NumServices int32

	
	
	TypeResolver interface {
		protoregistry.ExtensionTypeResolver
	}

	
	
	
	FileRegistry interface {
		FindFileByPath(string) (protoreflect.FileDescriptor, error)
		FindDescriptorByName(protoreflect.FullName) (protoreflect.Descriptor, error)
		RegisterFile(protoreflect.FileDescriptor) error
	}
}




type resolverByIndex interface {
	FindEnumByIndex(int32, int32, []Enum, []Message) protoreflect.EnumDescriptor
	FindMessageByIndex(int32, int32, []Enum, []Message) protoreflect.MessageDescriptor
}


const (
	listFieldDeps int32 = iota
	listExtTargets
	listExtDeps
	listMethInDeps
	listMethOutDeps
)


type Out struct {
	File protoreflect.FileDescriptor

	
	Enums []Enum
	
	
	Messages []Message
	
	Extensions []Extension
	
	Services []Service
}







func (db Builder) Build() (out Out) {
	
	if db.NumEnums+db.NumMessages+db.NumExtensions+db.NumServices == 0 {
		db.unmarshalCounts(db.RawDescriptor, true)
	}

	
	if db.TypeResolver == nil {
		db.TypeResolver = protoregistry.GlobalTypes
	}
	if db.FileRegistry == nil {
		db.FileRegistry = protoregistry.GlobalFiles
	}

	fd := newRawFile(db)
	out.File = fd
	out.Enums = fd.allEnums
	out.Messages = fd.allMessages
	out.Extensions = fd.allExtensions
	out.Services = fd.allServices

	if err := db.FileRegistry.RegisterFile(fd); err != nil {
		panic(err)
	}
	return out
}




func (db *Builder) unmarshalCounts(b []byte, isFile bool) {
	for len(b) > 0 {
		num, typ, n := protowire.ConsumeTag(b)
		b = b[n:]
		switch typ {
		case protowire.BytesType:
			v, m := protowire.ConsumeBytes(b)
			b = b[m:]
			if isFile {
				switch num {
				case genid.FileDescriptorProto_EnumType_field_number:
					db.NumEnums++
				case genid.FileDescriptorProto_MessageType_field_number:
					db.unmarshalCounts(v, false)
					db.NumMessages++
				case genid.FileDescriptorProto_Extension_field_number:
					db.NumExtensions++
				case genid.FileDescriptorProto_Service_field_number:
					db.NumServices++
				}
			} else {
				switch num {
				case genid.DescriptorProto_EnumType_field_number:
					db.NumEnums++
				case genid.DescriptorProto_NestedType_field_number:
					db.unmarshalCounts(v, false)
					db.NumMessages++
				case genid.DescriptorProto_Extension_field_number:
					db.NumExtensions++
				}
			}
		default:
			m := protowire.ConsumeFieldValue(num, typ, b)
			b = b[m:]
		}
	}
}
