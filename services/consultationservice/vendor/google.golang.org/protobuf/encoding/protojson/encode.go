



package protojson

import (
	"encoding/base64"
	"fmt"

	"google.golang.org/protobuf/internal/encoding/json"
	"google.golang.org/protobuf/internal/encoding/messageset"
	"google.golang.org/protobuf/internal/errors"
	"google.golang.org/protobuf/internal/filedesc"
	"google.golang.org/protobuf/internal/flags"
	"google.golang.org/protobuf/internal/genid"
	"google.golang.org/protobuf/internal/order"
	"google.golang.org/protobuf/internal/pragma"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

const defaultIndent = "  "






func Format(m proto.Message) string {
	return MarshalOptions{Multiline: true}.Format(m)
}





func Marshal(m proto.Message) ([]byte, error) {
	return MarshalOptions{}.Marshal(m)
}


type MarshalOptions struct {
	pragma.NoUnkeyedLiterals

	
	
	
	Multiline bool

	
	
	
	
	Indent string

	
	
	
	AllowPartial bool

	
	
	UseProtoNames bool

	
	UseEnumNumbers bool

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	EmitUnpopulated bool

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	EmitDefaultValues bool

	
	
	Resolver interface {
		protoregistry.ExtensionTypeResolver
		protoregistry.MessageTypeResolver
	}
}






func (o MarshalOptions) Format(m proto.Message) string {
	if m == nil || !m.ProtoReflect().IsValid() {
		return "<nil>" 
	}
	o.AllowPartial = true
	b, _ := o.Marshal(m)
	return string(b)
}





func (o MarshalOptions) Marshal(m proto.Message) ([]byte, error) {
	return o.marshal(nil, m)
}



func (o MarshalOptions) MarshalAppend(b []byte, m proto.Message) ([]byte, error) {
	return o.marshal(b, m)
}




func (o MarshalOptions) marshal(b []byte, m proto.Message) ([]byte, error) {
	if o.Multiline && o.Indent == "" {
		o.Indent = defaultIndent
	}
	if o.Resolver == nil {
		o.Resolver = protoregistry.GlobalTypes
	}

	internalEnc, err := json.NewEncoder(b, o.Indent)
	if err != nil {
		return nil, err
	}

	
	
	if m == nil {
		return append(b, '{', '}'), nil
	}

	enc := encoder{internalEnc, o}
	if err := enc.marshalMessage(m.ProtoReflect(), ""); err != nil {
		return nil, err
	}
	if o.AllowPartial {
		return enc.Bytes(), nil
	}
	return enc.Bytes(), proto.CheckInitialized(m)
}

type encoder struct {
	*json.Encoder
	opts MarshalOptions
}


var typeFieldDesc = func() protoreflect.FieldDescriptor {
	var fd filedesc.Field
	fd.L0.FullName = "@type"
	fd.L0.Index = -1
	fd.L1.Cardinality = protoreflect.Optional
	fd.L1.Kind = protoreflect.StringKind
	return &fd
}()



type typeURLFieldRanger struct {
	order.FieldRanger
	typeURL string
}

func (m typeURLFieldRanger) Range(f func(protoreflect.FieldDescriptor, protoreflect.Value) bool) {
	if !f(typeFieldDesc, protoreflect.ValueOfString(m.typeURL)) {
		return
	}
	m.FieldRanger.Range(f)
}



type unpopulatedFieldRanger struct {
	protoreflect.Message

	skipNull bool
}

func (m unpopulatedFieldRanger) Range(f func(protoreflect.FieldDescriptor, protoreflect.Value) bool) {
	fds := m.Descriptor().Fields()
	for i := 0; i < fds.Len(); i++ {
		fd := fds.Get(i)
		if m.Has(fd) || fd.ContainingOneof() != nil {
			continue 
		}

		v := m.Get(fd)
		if fd.HasPresence() {
			if m.skipNull {
				continue
			}
			v = protoreflect.Value{} 
		}
		if !f(fd, v) {
			return
		}
	}
	m.Message.Range(f)
}




func (e encoder) marshalMessage(m protoreflect.Message, typeURL string) error {
	if !flags.ProtoLegacy && messageset.IsMessageSet(m.Descriptor()) {
		return errors.New("no support for proto1 MessageSets")
	}

	if marshal := wellKnownTypeMarshaler(m.Descriptor().FullName()); marshal != nil {
		return marshal(e, m)
	}

	e.StartObject()
	defer e.EndObject()

	var fields order.FieldRanger = m
	switch {
	case e.opts.EmitUnpopulated:
		fields = unpopulatedFieldRanger{Message: m, skipNull: false}
	case e.opts.EmitDefaultValues:
		fields = unpopulatedFieldRanger{Message: m, skipNull: true}
	}
	if typeURL != "" {
		fields = typeURLFieldRanger{fields, typeURL}
	}

	var err error
	order.RangeFields(fields, order.IndexNameFieldOrder, func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		name := fd.JSONName()
		if e.opts.UseProtoNames {
			name = fd.TextName()
		}

		if err = e.WriteName(name); err != nil {
			return false
		}
		if err = e.marshalValue(v, fd); err != nil {
			return false
		}
		return true
	})
	return err
}


func (e encoder) marshalValue(val protoreflect.Value, fd protoreflect.FieldDescriptor) error {
	switch {
	case fd.IsList():
		return e.marshalList(val.List(), fd)
	case fd.IsMap():
		return e.marshalMap(val.Map(), fd)
	default:
		return e.marshalSingular(val, fd)
	}
}



func (e encoder) marshalSingular(val protoreflect.Value, fd protoreflect.FieldDescriptor) error {
	if !val.IsValid() {
		e.WriteNull()
		return nil
	}

	switch kind := fd.Kind(); kind {
	case protoreflect.BoolKind:
		e.WriteBool(val.Bool())

	case protoreflect.StringKind:
		if e.WriteString(val.String()) != nil {
			return errors.InvalidUTF8(string(fd.FullName()))
		}

	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		e.WriteInt(val.Int())

	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		e.WriteUint(val.Uint())

	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Uint64Kind,
		protoreflect.Sfixed64Kind, protoreflect.Fixed64Kind:
		
		e.WriteString(val.String())

	case protoreflect.FloatKind:
		
		e.WriteFloat(val.Float(), 32)

	case protoreflect.DoubleKind:
		
		e.WriteFloat(val.Float(), 64)

	case protoreflect.BytesKind:
		e.WriteString(base64.StdEncoding.EncodeToString(val.Bytes()))

	case protoreflect.EnumKind:
		if fd.Enum().FullName() == genid.NullValue_enum_fullname {
			e.WriteNull()
		} else {
			desc := fd.Enum().Values().ByNumber(val.Enum())
			if e.opts.UseEnumNumbers || desc == nil {
				e.WriteInt(int64(val.Enum()))
			} else {
				e.WriteString(string(desc.Name()))
			}
		}

	case protoreflect.MessageKind, protoreflect.GroupKind:
		if err := e.marshalMessage(val.Message(), ""); err != nil {
			return err
		}

	default:
		panic(fmt.Sprintf("%v has unknown kind: %v", fd.FullName(), kind))
	}
	return nil
}


func (e encoder) marshalList(list protoreflect.List, fd protoreflect.FieldDescriptor) error {
	e.StartArray()
	defer e.EndArray()

	for i := 0; i < list.Len(); i++ {
		item := list.Get(i)
		if err := e.marshalSingular(item, fd); err != nil {
			return err
		}
	}
	return nil
}


func (e encoder) marshalMap(mmap protoreflect.Map, fd protoreflect.FieldDescriptor) error {
	e.StartObject()
	defer e.EndObject()

	var err error
	order.RangeEntries(mmap, order.GenericKeyOrder, func(k protoreflect.MapKey, v protoreflect.Value) bool {
		if err = e.WriteName(k.String()); err != nil {
			return false
		}
		if err = e.marshalSingular(v, fd.MapValue()); err != nil {
			return false
		}
		return true
	})
	return err
}
