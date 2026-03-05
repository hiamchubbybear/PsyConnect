



package prototext

import (
	"fmt"
	"strconv"
	"unicode/utf8"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/internal/encoding/messageset"
	"google.golang.org/protobuf/internal/encoding/text"
	"google.golang.org/protobuf/internal/errors"
	"google.golang.org/protobuf/internal/flags"
	"google.golang.org/protobuf/internal/genid"
	"google.golang.org/protobuf/internal/order"
	"google.golang.org/protobuf/internal/pragma"
	"google.golang.org/protobuf/internal/strs"
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

	
	
	EmitASCII bool

	
	
	
	allowInvalidUTF8 bool

	
	
	
	AllowPartial bool

	
	
	
	EmitUnknown bool

	
	
	Resolver interface {
		protoregistry.ExtensionTypeResolver
		protoregistry.MessageTypeResolver
	}
}






func (o MarshalOptions) Format(m proto.Message) string {
	if m == nil || !m.ProtoReflect().IsValid() {
		return "<nil>" 
	}
	o.allowInvalidUTF8 = true
	o.AllowPartial = true
	o.EmitUnknown = true
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
	var delims = [2]byte{'{', '}'}

	if o.Multiline && o.Indent == "" {
		o.Indent = defaultIndent
	}
	if o.Resolver == nil {
		o.Resolver = protoregistry.GlobalTypes
	}

	internalEnc, err := text.NewEncoder(b, o.Indent, delims, o.EmitASCII)
	if err != nil {
		return nil, err
	}

	
	
	if m == nil {
		return b, nil
	}

	enc := encoder{internalEnc, o}
	err = enc.marshalMessage(m.ProtoReflect(), false)
	if err != nil {
		return nil, err
	}
	out := enc.Bytes()
	if len(o.Indent) > 0 && len(out) > 0 {
		out = append(out, '\n')
	}
	if o.AllowPartial {
		return out, nil
	}
	return out, proto.CheckInitialized(m)
}

type encoder struct {
	*text.Encoder
	opts MarshalOptions
}


func (e encoder) marshalMessage(m protoreflect.Message, inclDelims bool) error {
	messageDesc := m.Descriptor()
	if !flags.ProtoLegacy && messageset.IsMessageSet(messageDesc) {
		return errors.New("no support for proto1 MessageSets")
	}

	if inclDelims {
		e.StartMessage()
		defer e.EndMessage()
	}

	
	if messageDesc.FullName() == genid.Any_message_fullname {
		if e.marshalAny(m) {
			return nil
		}
		
	}

	
	var err error
	order.RangeFields(m, order.IndexNameFieldOrder, func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		if err = e.marshalField(fd.TextName(), v, fd); err != nil {
			return false
		}
		return true
	})
	if err != nil {
		return err
	}

	
	if e.opts.EmitUnknown {
		e.marshalUnknown(m.GetUnknown())
	}

	return nil
}


func (e encoder) marshalField(name string, val protoreflect.Value, fd protoreflect.FieldDescriptor) error {
	switch {
	case fd.IsList():
		return e.marshalList(name, val.List(), fd)
	case fd.IsMap():
		return e.marshalMap(name, val.Map(), fd)
	default:
		e.WriteName(name)
		return e.marshalSingular(val, fd)
	}
}



func (e encoder) marshalSingular(val protoreflect.Value, fd protoreflect.FieldDescriptor) error {
	kind := fd.Kind()
	switch kind {
	case protoreflect.BoolKind:
		e.WriteBool(val.Bool())

	case protoreflect.StringKind:
		s := val.String()
		if !e.opts.allowInvalidUTF8 && strs.EnforceUTF8(fd) && !utf8.ValidString(s) {
			return errors.InvalidUTF8(string(fd.FullName()))
		}
		e.WriteString(s)

	case protoreflect.Int32Kind, protoreflect.Int64Kind,
		protoreflect.Sint32Kind, protoreflect.Sint64Kind,
		protoreflect.Sfixed32Kind, protoreflect.Sfixed64Kind:
		e.WriteInt(val.Int())

	case protoreflect.Uint32Kind, protoreflect.Uint64Kind,
		protoreflect.Fixed32Kind, protoreflect.Fixed64Kind:
		e.WriteUint(val.Uint())

	case protoreflect.FloatKind:
		
		e.WriteFloat(val.Float(), 32)

	case protoreflect.DoubleKind:
		
		e.WriteFloat(val.Float(), 64)

	case protoreflect.BytesKind:
		e.WriteString(string(val.Bytes()))

	case protoreflect.EnumKind:
		num := val.Enum()
		if desc := fd.Enum().Values().ByNumber(num); desc != nil {
			e.WriteLiteral(string(desc.Name()))
		} else {
			
			e.WriteInt(int64(num))
		}

	case protoreflect.MessageKind, protoreflect.GroupKind:
		return e.marshalMessage(val.Message(), true)

	default:
		panic(fmt.Sprintf("%v has unknown kind: %v", fd.FullName(), kind))
	}
	return nil
}


func (e encoder) marshalList(name string, list protoreflect.List, fd protoreflect.FieldDescriptor) error {
	size := list.Len()
	for i := 0; i < size; i++ {
		e.WriteName(name)
		if err := e.marshalSingular(list.Get(i), fd); err != nil {
			return err
		}
	}
	return nil
}


func (e encoder) marshalMap(name string, mmap protoreflect.Map, fd protoreflect.FieldDescriptor) error {
	var err error
	order.RangeEntries(mmap, order.GenericKeyOrder, func(key protoreflect.MapKey, val protoreflect.Value) bool {
		e.WriteName(name)
		e.StartMessage()
		defer e.EndMessage()

		e.WriteName(string(genid.MapEntry_Key_field_name))
		err = e.marshalSingular(key.Value(), fd.MapKey())
		if err != nil {
			return false
		}

		e.WriteName(string(genid.MapEntry_Value_field_name))
		err = e.marshalSingular(val, fd.MapValue())
		if err != nil {
			return false
		}
		return true
	})
	return err
}



func (e encoder) marshalUnknown(b []byte) {
	const dec = 10
	const hex = 16
	for len(b) > 0 {
		num, wtype, n := protowire.ConsumeTag(b)
		b = b[n:]
		e.WriteName(strconv.FormatInt(int64(num), dec))

		switch wtype {
		case protowire.VarintType:
			var v uint64
			v, n = protowire.ConsumeVarint(b)
			e.WriteUint(v)
		case protowire.Fixed32Type:
			var v uint32
			v, n = protowire.ConsumeFixed32(b)
			e.WriteLiteral("0x" + strconv.FormatUint(uint64(v), hex))
		case protowire.Fixed64Type:
			var v uint64
			v, n = protowire.ConsumeFixed64(b)
			e.WriteLiteral("0x" + strconv.FormatUint(v, hex))
		case protowire.BytesType:
			var v []byte
			v, n = protowire.ConsumeBytes(b)
			e.WriteString(string(v))
		case protowire.StartGroupType:
			e.StartMessage()
			var v []byte
			v, n = protowire.ConsumeGroup(num, b)
			e.marshalUnknown(v)
			e.EndMessage()
		default:
			panic(fmt.Sprintf("prototext: error parsing unknown field wire type: %v", wtype))
		}

		b = b[n:]
	}
}



func (e encoder) marshalAny(any protoreflect.Message) bool {
	
	fds := any.Descriptor().Fields()
	fdType := fds.ByNumber(genid.Any_TypeUrl_field_number)
	typeURL := any.Get(fdType).String()
	mt, err := e.opts.Resolver.FindMessageByURL(typeURL)
	if err != nil {
		return false
	}
	m := mt.New().Interface()

	
	fdValue := fds.ByNumber(genid.Any_Value_field_number)
	value := any.Get(fdValue)
	err = proto.UnmarshalOptions{
		AllowPartial: true,
		Resolver:     e.opts.Resolver,
	}.Unmarshal(value.Bytes(), m)
	if err != nil {
		return false
	}

	
	
	pos := e.Snapshot()

	
	e.WriteName("[" + typeURL + "]")
	err = e.marshalMessage(m.ProtoReflect(), true)
	if err != nil {
		e.Reset(pos)
		return false
	}
	return true
}
