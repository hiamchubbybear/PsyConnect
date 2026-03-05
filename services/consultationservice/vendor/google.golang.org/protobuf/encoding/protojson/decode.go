



package protojson

import (
	"encoding/base64"
	"fmt"
	"math"
	"strconv"
	"strings"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/internal/encoding/json"
	"google.golang.org/protobuf/internal/encoding/messageset"
	"google.golang.org/protobuf/internal/errors"
	"google.golang.org/protobuf/internal/flags"
	"google.golang.org/protobuf/internal/genid"
	"google.golang.org/protobuf/internal/pragma"
	"google.golang.org/protobuf/internal/set"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)



func Unmarshal(b []byte, m proto.Message) error {
	return UnmarshalOptions{}.Unmarshal(b, m)
}


type UnmarshalOptions struct {
	pragma.NoUnkeyedLiterals

	
	
	AllowPartial bool

	
	DiscardUnknown bool

	
	
	
	Resolver interface {
		protoregistry.MessageTypeResolver
		protoregistry.ExtensionTypeResolver
	}

	
	
	RecursionLimit int
}






func (o UnmarshalOptions) Unmarshal(b []byte, m proto.Message) error {
	return o.unmarshal(b, m)
}




func (o UnmarshalOptions) unmarshal(b []byte, m proto.Message) error {
	proto.Reset(m)

	if o.Resolver == nil {
		o.Resolver = protoregistry.GlobalTypes
	}
	if o.RecursionLimit == 0 {
		o.RecursionLimit = protowire.DefaultRecursionLimit
	}

	dec := decoder{json.NewDecoder(b), o}
	if err := dec.unmarshalMessage(m.ProtoReflect(), false); err != nil {
		return err
	}

	
	tok, err := dec.Read()
	if err != nil {
		return err
	}
	if tok.Kind() != json.EOF {
		return dec.unexpectedTokenError(tok)
	}

	if o.AllowPartial {
		return nil
	}
	return proto.CheckInitialized(m)
}

type decoder struct {
	*json.Decoder
	opts UnmarshalOptions
}


func (d decoder) newError(pos int, f string, x ...any) error {
	line, column := d.Position(pos)
	head := fmt.Sprintf("(line %d:%d): ", line, column)
	return errors.New(head+f, x...)
}


func (d decoder) unexpectedTokenError(tok json.Token) error {
	return d.syntaxError(tok.Pos(), "unexpected token %s", tok.RawString())
}


func (d decoder) syntaxError(pos int, f string, x ...any) error {
	line, column := d.Position(pos)
	head := fmt.Sprintf("syntax error (line %d:%d): ", line, column)
	return errors.New(head+f, x...)
}


func (d decoder) unmarshalMessage(m protoreflect.Message, skipTypeURL bool) error {
	d.opts.RecursionLimit--
	if d.opts.RecursionLimit < 0 {
		return errors.New("exceeded max recursion depth")
	}
	if unmarshal := wellKnownTypeUnmarshaler(m.Descriptor().FullName()); unmarshal != nil {
		return unmarshal(d, m)
	}

	tok, err := d.Read()
	if err != nil {
		return err
	}
	if tok.Kind() != json.ObjectOpen {
		return d.unexpectedTokenError(tok)
	}

	messageDesc := m.Descriptor()
	if !flags.ProtoLegacy && messageset.IsMessageSet(messageDesc) {
		return errors.New("no support for proto1 MessageSets")
	}

	var seenNums set.Ints
	var seenOneofs set.Ints
	fieldDescs := messageDesc.Fields()
	for {
		
		tok, err := d.Read()
		if err != nil {
			return err
		}
		switch tok.Kind() {
		default:
			return d.unexpectedTokenError(tok)
		case json.ObjectClose:
			return nil
		case json.Name:
			
		}

		name := tok.Name()
		
		
		
		if skipTypeURL && name == "@type" {
			d.Read()
			continue
		}

		
		var fd protoreflect.FieldDescriptor
		if strings.HasPrefix(name, "[") && strings.HasSuffix(name, "]") {
			
			extName := protoreflect.FullName(name[1 : len(name)-1])
			extType, err := d.opts.Resolver.FindExtensionByName(extName)
			if err != nil && err != protoregistry.NotFound {
				return d.newError(tok.Pos(), "unable to resolve %s: %v", tok.RawString(), err)
			}
			if extType != nil {
				fd = extType.TypeDescriptor()
				if !messageDesc.ExtensionRanges().Has(fd.Number()) || fd.ContainingMessage().FullName() != messageDesc.FullName() {
					return d.newError(tok.Pos(), "message %v cannot be extended by %v", messageDesc.FullName(), fd.FullName())
				}
			}
		} else {
			
			fd = fieldDescs.ByJSONName(name)
			if fd == nil {
				fd = fieldDescs.ByTextName(name)
			}
		}

		if fd == nil {
			
			if d.opts.DiscardUnknown {
				if err := d.skipJSONValue(); err != nil {
					return err
				}
				continue
			}
			return d.newError(tok.Pos(), "unknown field %v", tok.RawString())
		}

		
		num := uint64(fd.Number())
		if seenNums.Has(num) {
			return d.newError(tok.Pos(), "duplicate field %v", tok.RawString())
		}
		seenNums.Set(num)

		
		
		if tok, _ := d.Peek(); tok.Kind() == json.Null && !isKnownValue(fd) && !isNullValue(fd) {
			d.Read()
			continue
		}

		switch {
		case fd.IsList():
			list := m.Mutable(fd).List()
			if err := d.unmarshalList(list, fd); err != nil {
				return err
			}
		case fd.IsMap():
			mmap := m.Mutable(fd).Map()
			if err := d.unmarshalMap(mmap, fd); err != nil {
				return err
			}
		default:
			
			if od := fd.ContainingOneof(); od != nil {
				idx := uint64(od.Index())
				if seenOneofs.Has(idx) {
					return d.newError(tok.Pos(), "error parsing %s, oneof %v is already set", tok.RawString(), od.FullName())
				}
				seenOneofs.Set(idx)
			}

			
			if err := d.unmarshalSingular(m, fd); err != nil {
				return err
			}
		}
	}
}

func isKnownValue(fd protoreflect.FieldDescriptor) bool {
	md := fd.Message()
	return md != nil && md.FullName() == genid.Value_message_fullname
}

func isNullValue(fd protoreflect.FieldDescriptor) bool {
	ed := fd.Enum()
	return ed != nil && ed.FullName() == genid.NullValue_enum_fullname
}



func (d decoder) unmarshalSingular(m protoreflect.Message, fd protoreflect.FieldDescriptor) error {
	var val protoreflect.Value
	var err error
	switch fd.Kind() {
	case protoreflect.MessageKind, protoreflect.GroupKind:
		val = m.NewField(fd)
		err = d.unmarshalMessage(val.Message(), false)
	default:
		val, err = d.unmarshalScalar(fd)
	}

	if err != nil {
		return err
	}
	if val.IsValid() {
		m.Set(fd, val)
	}
	return nil
}



func (d decoder) unmarshalScalar(fd protoreflect.FieldDescriptor) (protoreflect.Value, error) {
	const b32 int = 32
	const b64 int = 64

	tok, err := d.Read()
	if err != nil {
		return protoreflect.Value{}, err
	}

	kind := fd.Kind()
	switch kind {
	case protoreflect.BoolKind:
		if tok.Kind() == json.Bool {
			return protoreflect.ValueOfBool(tok.Bool()), nil
		}

	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		if v, ok := unmarshalInt(tok, b32); ok {
			return v, nil
		}

	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		if v, ok := unmarshalInt(tok, b64); ok {
			return v, nil
		}

	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		if v, ok := unmarshalUint(tok, b32); ok {
			return v, nil
		}

	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		if v, ok := unmarshalUint(tok, b64); ok {
			return v, nil
		}

	case protoreflect.FloatKind:
		if v, ok := unmarshalFloat(tok, b32); ok {
			return v, nil
		}

	case protoreflect.DoubleKind:
		if v, ok := unmarshalFloat(tok, b64); ok {
			return v, nil
		}

	case protoreflect.StringKind:
		if tok.Kind() == json.String {
			return protoreflect.ValueOfString(tok.ParsedString()), nil
		}

	case protoreflect.BytesKind:
		if v, ok := unmarshalBytes(tok); ok {
			return v, nil
		}

	case protoreflect.EnumKind:
		if v, ok := unmarshalEnum(tok, fd, d.opts.DiscardUnknown); ok {
			return v, nil
		}

	default:
		panic(fmt.Sprintf("unmarshalScalar: invalid scalar kind %v", kind))
	}

	return protoreflect.Value{}, d.newError(tok.Pos(), "invalid value for %v field %v: %v", kind, fd.JSONName(), tok.RawString())
}

func unmarshalInt(tok json.Token, bitSize int) (protoreflect.Value, bool) {
	switch tok.Kind() {
	case json.Number:
		return getInt(tok, bitSize)

	case json.String:
		
		s := strings.TrimSpace(tok.ParsedString())
		if len(s) != len(tok.ParsedString()) {
			return protoreflect.Value{}, false
		}
		dec := json.NewDecoder([]byte(s))
		tok, err := dec.Read()
		if err != nil {
			return protoreflect.Value{}, false
		}
		return getInt(tok, bitSize)
	}
	return protoreflect.Value{}, false
}

func getInt(tok json.Token, bitSize int) (protoreflect.Value, bool) {
	n, ok := tok.Int(bitSize)
	if !ok {
		return protoreflect.Value{}, false
	}
	if bitSize == 32 {
		return protoreflect.ValueOfInt32(int32(n)), true
	}
	return protoreflect.ValueOfInt64(n), true
}

func unmarshalUint(tok json.Token, bitSize int) (protoreflect.Value, bool) {
	switch tok.Kind() {
	case json.Number:
		return getUint(tok, bitSize)

	case json.String:
		
		s := strings.TrimSpace(tok.ParsedString())
		if len(s) != len(tok.ParsedString()) {
			return protoreflect.Value{}, false
		}
		dec := json.NewDecoder([]byte(s))
		tok, err := dec.Read()
		if err != nil {
			return protoreflect.Value{}, false
		}
		return getUint(tok, bitSize)
	}
	return protoreflect.Value{}, false
}

func getUint(tok json.Token, bitSize int) (protoreflect.Value, bool) {
	n, ok := tok.Uint(bitSize)
	if !ok {
		return protoreflect.Value{}, false
	}
	if bitSize == 32 {
		return protoreflect.ValueOfUint32(uint32(n)), true
	}
	return protoreflect.ValueOfUint64(n), true
}

func unmarshalFloat(tok json.Token, bitSize int) (protoreflect.Value, bool) {
	switch tok.Kind() {
	case json.Number:
		return getFloat(tok, bitSize)

	case json.String:
		s := tok.ParsedString()
		switch s {
		case "NaN":
			if bitSize == 32 {
				return protoreflect.ValueOfFloat32(float32(math.NaN())), true
			}
			return protoreflect.ValueOfFloat64(math.NaN()), true
		case "Infinity":
			if bitSize == 32 {
				return protoreflect.ValueOfFloat32(float32(math.Inf(+1))), true
			}
			return protoreflect.ValueOfFloat64(math.Inf(+1)), true
		case "-Infinity":
			if bitSize == 32 {
				return protoreflect.ValueOfFloat32(float32(math.Inf(-1))), true
			}
			return protoreflect.ValueOfFloat64(math.Inf(-1)), true
		}

		
		if len(s) != len(strings.TrimSpace(s)) {
			return protoreflect.Value{}, false
		}
		dec := json.NewDecoder([]byte(s))
		tok, err := dec.Read()
		if err != nil {
			return protoreflect.Value{}, false
		}
		return getFloat(tok, bitSize)
	}
	return protoreflect.Value{}, false
}

func getFloat(tok json.Token, bitSize int) (protoreflect.Value, bool) {
	n, ok := tok.Float(bitSize)
	if !ok {
		return protoreflect.Value{}, false
	}
	if bitSize == 32 {
		return protoreflect.ValueOfFloat32(float32(n)), true
	}
	return protoreflect.ValueOfFloat64(n), true
}

func unmarshalBytes(tok json.Token) (protoreflect.Value, bool) {
	if tok.Kind() != json.String {
		return protoreflect.Value{}, false
	}

	s := tok.ParsedString()
	enc := base64.StdEncoding
	if strings.ContainsAny(s, "-_") {
		enc = base64.URLEncoding
	}
	if len(s)%4 != 0 {
		enc = enc.WithPadding(base64.NoPadding)
	}
	b, err := enc.DecodeString(s)
	if err != nil {
		return protoreflect.Value{}, false
	}
	return protoreflect.ValueOfBytes(b), true
}

func unmarshalEnum(tok json.Token, fd protoreflect.FieldDescriptor, discardUnknown bool) (protoreflect.Value, bool) {
	switch tok.Kind() {
	case json.String:
		
		s := tok.ParsedString()
		if enumVal := fd.Enum().Values().ByName(protoreflect.Name(s)); enumVal != nil {
			return protoreflect.ValueOfEnum(enumVal.Number()), true
		}
		if discardUnknown {
			return protoreflect.Value{}, true
		}

	case json.Number:
		if n, ok := tok.Int(32); ok {
			return protoreflect.ValueOfEnum(protoreflect.EnumNumber(n)), true
		}

	case json.Null:
		
		if isNullValue(fd) {
			return protoreflect.ValueOfEnum(0), true
		}
	}

	return protoreflect.Value{}, false
}

func (d decoder) unmarshalList(list protoreflect.List, fd protoreflect.FieldDescriptor) error {
	tok, err := d.Read()
	if err != nil {
		return err
	}
	if tok.Kind() != json.ArrayOpen {
		return d.unexpectedTokenError(tok)
	}

	switch fd.Kind() {
	case protoreflect.MessageKind, protoreflect.GroupKind:
		for {
			tok, err := d.Peek()
			if err != nil {
				return err
			}

			if tok.Kind() == json.ArrayClose {
				d.Read()
				return nil
			}

			val := list.NewElement()
			if err := d.unmarshalMessage(val.Message(), false); err != nil {
				return err
			}
			list.Append(val)
		}
	default:
		for {
			tok, err := d.Peek()
			if err != nil {
				return err
			}

			if tok.Kind() == json.ArrayClose {
				d.Read()
				return nil
			}

			val, err := d.unmarshalScalar(fd)
			if err != nil {
				return err
			}
			if val.IsValid() {
				list.Append(val)
			}
		}
	}

	return nil
}

func (d decoder) unmarshalMap(mmap protoreflect.Map, fd protoreflect.FieldDescriptor) error {
	tok, err := d.Read()
	if err != nil {
		return err
	}
	if tok.Kind() != json.ObjectOpen {
		return d.unexpectedTokenError(tok)
	}

	
	
	
	var unmarshalMapValue func() (protoreflect.Value, error)
	switch fd.MapValue().Kind() {
	case protoreflect.MessageKind, protoreflect.GroupKind:
		unmarshalMapValue = func() (protoreflect.Value, error) {
			val := mmap.NewValue()
			if err := d.unmarshalMessage(val.Message(), false); err != nil {
				return protoreflect.Value{}, err
			}
			return val, nil
		}
	default:
		unmarshalMapValue = func() (protoreflect.Value, error) {
			return d.unmarshalScalar(fd.MapValue())
		}
	}

Loop:
	for {
		
		tok, err := d.Read()
		if err != nil {
			return err
		}
		switch tok.Kind() {
		default:
			return d.unexpectedTokenError(tok)
		case json.ObjectClose:
			break Loop
		case json.Name:
			
		}

		
		pkey, err := d.unmarshalMapKey(tok, fd.MapKey())
		if err != nil {
			return err
		}

		
		if mmap.Has(pkey) {
			return d.newError(tok.Pos(), "duplicate map key %v", tok.RawString())
		}

		
		pval, err := unmarshalMapValue()
		if err != nil {
			return err
		}
		if pval.IsValid() {
			mmap.Set(pkey, pval)
		}
	}

	return nil
}



func (d decoder) unmarshalMapKey(tok json.Token, fd protoreflect.FieldDescriptor) (protoreflect.MapKey, error) {
	const b32 = 32
	const b64 = 64
	const base10 = 10

	name := tok.Name()
	kind := fd.Kind()
	switch kind {
	case protoreflect.StringKind:
		return protoreflect.ValueOfString(name).MapKey(), nil

	case protoreflect.BoolKind:
		switch name {
		case "true":
			return protoreflect.ValueOfBool(true).MapKey(), nil
		case "false":
			return protoreflect.ValueOfBool(false).MapKey(), nil
		}

	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		if n, err := strconv.ParseInt(name, base10, b32); err == nil {
			return protoreflect.ValueOfInt32(int32(n)).MapKey(), nil
		}

	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		if n, err := strconv.ParseInt(name, base10, b64); err == nil {
			return protoreflect.ValueOfInt64(int64(n)).MapKey(), nil
		}

	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		if n, err := strconv.ParseUint(name, base10, b32); err == nil {
			return protoreflect.ValueOfUint32(uint32(n)).MapKey(), nil
		}

	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		if n, err := strconv.ParseUint(name, base10, b64); err == nil {
			return protoreflect.ValueOfUint64(uint64(n)).MapKey(), nil
		}

	default:
		panic(fmt.Sprintf("invalid kind for map key: %v", kind))
	}

	return protoreflect.MapKey{}, d.newError(tok.Pos(), "invalid value for %v key: %s", kind, tok.RawString())
}
