





package tag

import (
	"reflect"
	"strconv"
	"strings"

	"google.golang.org/protobuf/internal/encoding/defval"
	"google.golang.org/protobuf/internal/filedesc"
	"google.golang.org/protobuf/internal/strs"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var byteType = reflect.TypeOf(byte(0))











func Unmarshal(tag string, goType reflect.Type, evs protoreflect.EnumValueDescriptors) protoreflect.FieldDescriptor {
	f := new(filedesc.Field)
	f.L0.ParentFile = filedesc.SurrogateProto2
	f.L1.EditionFeatures = f.L0.ParentFile.L1.EditionFeatures
	for len(tag) > 0 {
		i := strings.IndexByte(tag, ',')
		if i < 0 {
			i = len(tag)
		}
		switch s := tag[:i]; {
		case strings.HasPrefix(s, "name="):
			f.L0.FullName = protoreflect.FullName(s[len("name="):])
		case strings.Trim(s, "0123456789") == "":
			n, _ := strconv.ParseUint(s, 10, 32)
			f.L1.Number = protoreflect.FieldNumber(n)
		case s == "opt":
			f.L1.Cardinality = protoreflect.Optional
		case s == "req":
			f.L1.Cardinality = protoreflect.Required
		case s == "rep":
			f.L1.Cardinality = protoreflect.Repeated
		case s == "varint":
			switch goType.Kind() {
			case reflect.Bool:
				f.L1.Kind = protoreflect.BoolKind
			case reflect.Int32:
				f.L1.Kind = protoreflect.Int32Kind
			case reflect.Int64:
				f.L1.Kind = protoreflect.Int64Kind
			case reflect.Uint32:
				f.L1.Kind = protoreflect.Uint32Kind
			case reflect.Uint64:
				f.L1.Kind = protoreflect.Uint64Kind
			}
		case s == "zigzag32":
			if goType.Kind() == reflect.Int32 {
				f.L1.Kind = protoreflect.Sint32Kind
			}
		case s == "zigzag64":
			if goType.Kind() == reflect.Int64 {
				f.L1.Kind = protoreflect.Sint64Kind
			}
		case s == "fixed32":
			switch goType.Kind() {
			case reflect.Int32:
				f.L1.Kind = protoreflect.Sfixed32Kind
			case reflect.Uint32:
				f.L1.Kind = protoreflect.Fixed32Kind
			case reflect.Float32:
				f.L1.Kind = protoreflect.FloatKind
			}
		case s == "fixed64":
			switch goType.Kind() {
			case reflect.Int64:
				f.L1.Kind = protoreflect.Sfixed64Kind
			case reflect.Uint64:
				f.L1.Kind = protoreflect.Fixed64Kind
			case reflect.Float64:
				f.L1.Kind = protoreflect.DoubleKind
			}
		case s == "bytes":
			switch {
			case goType.Kind() == reflect.String:
				f.L1.Kind = protoreflect.StringKind
			case goType.Kind() == reflect.Slice && goType.Elem() == byteType:
				f.L1.Kind = protoreflect.BytesKind
			default:
				f.L1.Kind = protoreflect.MessageKind
			}
		case s == "group":
			f.L1.Kind = protoreflect.GroupKind
		case strings.HasPrefix(s, "enum="):
			f.L1.Kind = protoreflect.EnumKind
		case strings.HasPrefix(s, "json="):
			jsonName := s[len("json="):]
			if jsonName != strs.JSONCamelCase(string(f.L0.FullName.Name())) {
				f.L1.StringName.InitJSON(jsonName)
			}
		case s == "packed":
			f.L1.EditionFeatures.IsPacked = true
		case strings.HasPrefix(s, "def="):
			
			
			s, i = tag[len("def="):], len(tag)
			v, ev, _ := defval.Unmarshal(s, f.L1.Kind, evs, defval.GoTag)
			f.L1.Default = filedesc.DefaultValue(v, ev)
		case s == "proto3":
			f.L0.ParentFile = filedesc.SurrogateProto3
		}
		tag = strings.TrimPrefix(tag[i:], ",")
	}

	
	
	if f.L1.Kind == protoreflect.GroupKind {
		f.L0.FullName = protoreflect.FullName(strings.ToLower(string(f.L0.FullName)))
	}
	return f
}









func Marshal(fd protoreflect.FieldDescriptor, enumName string) string {
	var tag []string
	switch fd.Kind() {
	case protoreflect.BoolKind, protoreflect.EnumKind, protoreflect.Int32Kind, protoreflect.Uint32Kind, protoreflect.Int64Kind, protoreflect.Uint64Kind:
		tag = append(tag, "varint")
	case protoreflect.Sint32Kind:
		tag = append(tag, "zigzag32")
	case protoreflect.Sint64Kind:
		tag = append(tag, "zigzag64")
	case protoreflect.Sfixed32Kind, protoreflect.Fixed32Kind, protoreflect.FloatKind:
		tag = append(tag, "fixed32")
	case protoreflect.Sfixed64Kind, protoreflect.Fixed64Kind, protoreflect.DoubleKind:
		tag = append(tag, "fixed64")
	case protoreflect.StringKind, protoreflect.BytesKind, protoreflect.MessageKind:
		tag = append(tag, "bytes")
	case protoreflect.GroupKind:
		tag = append(tag, "group")
	}
	tag = append(tag, strconv.Itoa(int(fd.Number())))
	switch fd.Cardinality() {
	case protoreflect.Optional:
		tag = append(tag, "opt")
	case protoreflect.Required:
		tag = append(tag, "req")
	case protoreflect.Repeated:
		tag = append(tag, "rep")
	}
	if fd.IsPacked() {
		tag = append(tag, "packed")
	}
	name := string(fd.Name())
	if fd.Kind() == protoreflect.GroupKind {
		
		
		
		name = string(fd.Message().Name())
	}
	tag = append(tag, "name="+name)
	if jsonName := fd.JSONName(); jsonName != "" && jsonName != name && !fd.IsExtension() {
		
		
		tag = append(tag, "json="+jsonName)
	}
	
	
	
	if fd.Syntax() == protoreflect.Proto3 && !fd.IsExtension() {
		tag = append(tag, "proto3")
	}
	if fd.Kind() == protoreflect.EnumKind && enumName != "" {
		tag = append(tag, "enum="+enumName)
	}
	if fd.ContainingOneof() != nil {
		tag = append(tag, "oneof")
	}
	
	if fd.HasDefault() {
		def, _ := defval.Marshal(fd.Default(), fd.DefaultEnumValue(), fd.Kind(), defval.GoTag)
		tag = append(tag, "def="+def)
	}
	return strings.Join(tag, ",")
}
