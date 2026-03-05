



package protojson

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"google.golang.org/protobuf/internal/encoding/json"
	"google.golang.org/protobuf/internal/errors"
	"google.golang.org/protobuf/internal/genid"
	"google.golang.org/protobuf/internal/strs"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type marshalFunc func(encoder, protoreflect.Message) error



func wellKnownTypeMarshaler(name protoreflect.FullName) marshalFunc {
	if name.Parent() == genid.GoogleProtobuf_package {
		switch name.Name() {
		case genid.Any_message_name:
			return encoder.marshalAny
		case genid.Timestamp_message_name:
			return encoder.marshalTimestamp
		case genid.Duration_message_name:
			return encoder.marshalDuration
		case genid.BoolValue_message_name,
			genid.Int32Value_message_name,
			genid.Int64Value_message_name,
			genid.UInt32Value_message_name,
			genid.UInt64Value_message_name,
			genid.FloatValue_message_name,
			genid.DoubleValue_message_name,
			genid.StringValue_message_name,
			genid.BytesValue_message_name:
			return encoder.marshalWrapperType
		case genid.Struct_message_name:
			return encoder.marshalStruct
		case genid.ListValue_message_name:
			return encoder.marshalListValue
		case genid.Value_message_name:
			return encoder.marshalKnownValue
		case genid.FieldMask_message_name:
			return encoder.marshalFieldMask
		case genid.Empty_message_name:
			return encoder.marshalEmpty
		}
	}
	return nil
}

type unmarshalFunc func(decoder, protoreflect.Message) error



func wellKnownTypeUnmarshaler(name protoreflect.FullName) unmarshalFunc {
	if name.Parent() == genid.GoogleProtobuf_package {
		switch name.Name() {
		case genid.Any_message_name:
			return decoder.unmarshalAny
		case genid.Timestamp_message_name:
			return decoder.unmarshalTimestamp
		case genid.Duration_message_name:
			return decoder.unmarshalDuration
		case genid.BoolValue_message_name,
			genid.Int32Value_message_name,
			genid.Int64Value_message_name,
			genid.UInt32Value_message_name,
			genid.UInt64Value_message_name,
			genid.FloatValue_message_name,
			genid.DoubleValue_message_name,
			genid.StringValue_message_name,
			genid.BytesValue_message_name:
			return decoder.unmarshalWrapperType
		case genid.Struct_message_name:
			return decoder.unmarshalStruct
		case genid.ListValue_message_name:
			return decoder.unmarshalListValue
		case genid.Value_message_name:
			return decoder.unmarshalKnownValue
		case genid.FieldMask_message_name:
			return decoder.unmarshalFieldMask
		case genid.Empty_message_name:
			return decoder.unmarshalEmpty
		}
	}
	return nil
}







func (e encoder) marshalAny(m protoreflect.Message) error {
	fds := m.Descriptor().Fields()
	fdType := fds.ByNumber(genid.Any_TypeUrl_field_number)
	fdValue := fds.ByNumber(genid.Any_Value_field_number)

	if !m.Has(fdType) {
		if !m.Has(fdValue) {
			
			e.StartObject()
			e.EndObject()
			return nil
		} else {
			
			return errors.New("%s: %v is not set", genid.Any_message_fullname, genid.Any_TypeUrl_field_name)
		}
	}

	typeVal := m.Get(fdType)
	valueVal := m.Get(fdValue)

	
	typeURL := typeVal.String()
	emt, err := e.opts.Resolver.FindMessageByURL(typeURL)
	if err != nil {
		return errors.New("%s: unable to resolve %q: %v", genid.Any_message_fullname, typeURL, err)
	}

	em := emt.New()
	err = proto.UnmarshalOptions{
		AllowPartial: true, 
		Resolver:     e.opts.Resolver,
	}.Unmarshal(valueVal.Bytes(), em.Interface())
	if err != nil {
		return errors.New("%s: unable to unmarshal %q: %v", genid.Any_message_fullname, typeURL, err)
	}

	
	
	
	if marshal := wellKnownTypeMarshaler(emt.Descriptor().FullName()); marshal != nil {
		e.StartObject()
		defer e.EndObject()

		
		e.WriteName("@type")
		if err := e.WriteString(typeURL); err != nil {
			return err
		}

		e.WriteName("value")
		return marshal(e, em)
	}

	
	if err := e.marshalMessage(em, typeURL); err != nil {
		return err
	}

	return nil
}

func (d decoder) unmarshalAny(m protoreflect.Message) error {
	
	start, err := d.Peek()
	if err != nil {
		return err
	}
	if start.Kind() != json.ObjectOpen {
		return d.unexpectedTokenError(start)
	}

	
	
	
	dec := decoder{d.Clone(), UnmarshalOptions{RecursionLimit: d.opts.RecursionLimit}}
	tok, err := findTypeURL(dec)
	switch err {
	case errEmptyObject:
		
		d.Read() 
		d.Read() 
		return nil

	case errMissingType:
		if d.opts.DiscardUnknown {
			
			return d.skipJSONValue()
		}
		
		return d.newError(start.Pos(), err.Error())

	default:
		if err != nil {
			return err
		}
	}

	typeURL := tok.ParsedString()
	emt, err := d.opts.Resolver.FindMessageByURL(typeURL)
	if err != nil {
		return d.newError(tok.Pos(), "unable to resolve %v: %q", tok.RawString(), err)
	}

	
	em := emt.New()
	if unmarshal := wellKnownTypeUnmarshaler(emt.Descriptor().FullName()); unmarshal != nil {
		
		
		if err := d.unmarshalAnyValue(unmarshal, em); err != nil {
			return err
		}
	} else {
		
		if err := d.unmarshalMessage(em, true); err != nil {
			return err
		}
	}
	
	
	b, err := proto.MarshalOptions{
		AllowPartial:  true, 
		Deterministic: true,
	}.Marshal(em.Interface())
	if err != nil {
		return d.newError(start.Pos(), "error in marshaling Any.value field: %v", err)
	}

	fds := m.Descriptor().Fields()
	fdType := fds.ByNumber(genid.Any_TypeUrl_field_number)
	fdValue := fds.ByNumber(genid.Any_Value_field_number)

	m.Set(fdType, protoreflect.ValueOfString(typeURL))
	m.Set(fdValue, protoreflect.ValueOfBytes(b))
	return nil
}

var errEmptyObject = fmt.Errorf(`empty object`)
var errMissingType = fmt.Errorf(`missing "@type" field`)






func findTypeURL(d decoder) (json.Token, error) {
	var typeURL string
	var typeTok json.Token
	numFields := 0
	
	d.Read()

Loop:
	for {
		tok, err := d.Read()
		if err != nil {
			return json.Token{}, err
		}

		switch tok.Kind() {
		case json.ObjectClose:
			if typeURL == "" {
				
				if numFields > 0 {
					return json.Token{}, errMissingType
				}
				return json.Token{}, errEmptyObject
			}
			break Loop

		case json.Name:
			numFields++
			if tok.Name() != "@type" {
				
				if err := d.skipJSONValue(); err != nil {
					return json.Token{}, err
				}
				continue
			}

			
			if typeURL != "" {
				return json.Token{}, d.newError(tok.Pos(), `duplicate "@type" field`)
			}
			
			tok, err := d.Read()
			if err != nil {
				return json.Token{}, err
			}
			if tok.Kind() != json.String {
				return json.Token{}, d.newError(tok.Pos(), `@type field value is not a string: %v`, tok.RawString())
			}
			typeURL = tok.ParsedString()
			if typeURL == "" {
				return json.Token{}, d.newError(tok.Pos(), `@type field contains empty value`)
			}
			typeTok = tok
		}
	}

	return typeTok, nil
}




func (d decoder) skipJSONValue() error {
	var open int
	for {
		tok, err := d.Read()
		if err != nil {
			return err
		}
		switch tok.Kind() {
		case json.ObjectClose, json.ArrayClose:
			open--
		case json.ObjectOpen, json.ArrayOpen:
			open++
			if open > d.opts.RecursionLimit {
				return errors.New("exceeded max recursion depth")
			}
		case json.EOF:
			
			
			return errors.New("unexpected EOF")
		}
		if open == 0 {
			return nil
		}
	}
}



func (d decoder) unmarshalAnyValue(unmarshal unmarshalFunc, m protoreflect.Message) error {
	
	d.Read()

	var found bool 
	for {
		tok, err := d.Read()
		if err != nil {
			return err
		}
		switch tok.Kind() {
		case json.ObjectClose:
			if !found {
				
				
				if m.Descriptor().FullName() != genid.Empty_message_fullname {
					return d.newError(tok.Pos(), `missing "value" field`)
				}
			}
			return nil

		case json.Name:
			switch tok.Name() {
			case "@type":
				
				d.Read()

			case "value":
				if found {
					return d.newError(tok.Pos(), `duplicate "value" field`)
				}
				
				if err := unmarshal(d, m); err != nil {
					return err
				}
				found = true

			default:
				if d.opts.DiscardUnknown {
					if err := d.skipJSONValue(); err != nil {
						return err
					}
					continue
				}
				return d.newError(tok.Pos(), "unknown field %v", tok.RawString())
			}
		}
	}
}



func (e encoder) marshalWrapperType(m protoreflect.Message) error {
	fd := m.Descriptor().Fields().ByNumber(genid.WrapperValue_Value_field_number)
	val := m.Get(fd)
	return e.marshalSingular(val, fd)
}

func (d decoder) unmarshalWrapperType(m protoreflect.Message) error {
	fd := m.Descriptor().Fields().ByNumber(genid.WrapperValue_Value_field_number)
	val, err := d.unmarshalScalar(fd)
	if err != nil {
		return err
	}
	m.Set(fd, val)
	return nil
}



func (e encoder) marshalEmpty(protoreflect.Message) error {
	e.StartObject()
	e.EndObject()
	return nil
}

func (d decoder) unmarshalEmpty(protoreflect.Message) error {
	tok, err := d.Read()
	if err != nil {
		return err
	}
	if tok.Kind() != json.ObjectOpen {
		return d.unexpectedTokenError(tok)
	}

	for {
		tok, err := d.Read()
		if err != nil {
			return err
		}
		switch tok.Kind() {
		case json.ObjectClose:
			return nil

		case json.Name:
			if d.opts.DiscardUnknown {
				if err := d.skipJSONValue(); err != nil {
					return err
				}
				continue
			}
			return d.newError(tok.Pos(), "unknown field %v", tok.RawString())

		default:
			return d.unexpectedTokenError(tok)
		}
	}
}




func (e encoder) marshalStruct(m protoreflect.Message) error {
	fd := m.Descriptor().Fields().ByNumber(genid.Struct_Fields_field_number)
	return e.marshalMap(m.Get(fd).Map(), fd)
}

func (d decoder) unmarshalStruct(m protoreflect.Message) error {
	fd := m.Descriptor().Fields().ByNumber(genid.Struct_Fields_field_number)
	return d.unmarshalMap(m.Mutable(fd).Map(), fd)
}





func (e encoder) marshalListValue(m protoreflect.Message) error {
	fd := m.Descriptor().Fields().ByNumber(genid.ListValue_Values_field_number)
	return e.marshalList(m.Get(fd).List(), fd)
}

func (d decoder) unmarshalListValue(m protoreflect.Message) error {
	fd := m.Descriptor().Fields().ByNumber(genid.ListValue_Values_field_number)
	return d.unmarshalList(m.Mutable(fd).List(), fd)
}





func (e encoder) marshalKnownValue(m protoreflect.Message) error {
	od := m.Descriptor().Oneofs().ByName(genid.Value_Kind_oneof_name)
	fd := m.WhichOneof(od)
	if fd == nil {
		return errors.New("%s: none of the oneof fields is set", genid.Value_message_fullname)
	}
	if fd.Number() == genid.Value_NumberValue_field_number {
		if v := m.Get(fd).Float(); math.IsNaN(v) || math.IsInf(v, 0) {
			return errors.New("%s: invalid %v value", genid.Value_NumberValue_field_fullname, v)
		}
	}
	return e.marshalSingular(m.Get(fd), fd)
}

func (d decoder) unmarshalKnownValue(m protoreflect.Message) error {
	tok, err := d.Peek()
	if err != nil {
		return err
	}

	var fd protoreflect.FieldDescriptor
	var val protoreflect.Value
	switch tok.Kind() {
	case json.Null:
		d.Read()
		fd = m.Descriptor().Fields().ByNumber(genid.Value_NullValue_field_number)
		val = protoreflect.ValueOfEnum(0)

	case json.Bool:
		tok, err := d.Read()
		if err != nil {
			return err
		}
		fd = m.Descriptor().Fields().ByNumber(genid.Value_BoolValue_field_number)
		val = protoreflect.ValueOfBool(tok.Bool())

	case json.Number:
		tok, err := d.Read()
		if err != nil {
			return err
		}
		fd = m.Descriptor().Fields().ByNumber(genid.Value_NumberValue_field_number)
		var ok bool
		val, ok = unmarshalFloat(tok, 64)
		if !ok {
			return d.newError(tok.Pos(), "invalid %v: %v", genid.Value_message_fullname, tok.RawString())
		}

	case json.String:
		
		
		
		
		
		
		tok, err := d.Read()
		if err != nil {
			return err
		}
		fd = m.Descriptor().Fields().ByNumber(genid.Value_StringValue_field_number)
		val = protoreflect.ValueOfString(tok.ParsedString())

	case json.ObjectOpen:
		fd = m.Descriptor().Fields().ByNumber(genid.Value_StructValue_field_number)
		val = m.NewField(fd)
		if err := d.unmarshalStruct(val.Message()); err != nil {
			return err
		}

	case json.ArrayOpen:
		fd = m.Descriptor().Fields().ByNumber(genid.Value_ListValue_field_number)
		val = m.NewField(fd)
		if err := d.unmarshalListValue(val.Message()); err != nil {
			return err
		}

	default:
		return d.newError(tok.Pos(), "invalid %v: %v", genid.Value_message_fullname, tok.RawString())
	}

	m.Set(fd, val)
	return nil
}













const (
	secondsInNanos       = 999999999
	maxSecondsInDuration = 315576000000
)

func (e encoder) marshalDuration(m protoreflect.Message) error {
	fds := m.Descriptor().Fields()
	fdSeconds := fds.ByNumber(genid.Duration_Seconds_field_number)
	fdNanos := fds.ByNumber(genid.Duration_Nanos_field_number)

	secsVal := m.Get(fdSeconds)
	nanosVal := m.Get(fdNanos)
	secs := secsVal.Int()
	nanos := nanosVal.Int()
	if secs < -maxSecondsInDuration || secs > maxSecondsInDuration {
		return errors.New("%s: seconds out of range %v", genid.Duration_message_fullname, secs)
	}
	if nanos < -secondsInNanos || nanos > secondsInNanos {
		return errors.New("%s: nanos out of range %v", genid.Duration_message_fullname, nanos)
	}
	if (secs > 0 && nanos < 0) || (secs < 0 && nanos > 0) {
		return errors.New("%s: signs of seconds and nanos do not match", genid.Duration_message_fullname)
	}
	
	
	var sign string
	if secs < 0 || nanos < 0 {
		sign, secs, nanos = "-", -1*secs, -1*nanos
	}
	x := fmt.Sprintf("%s%d.%09d", sign, secs, nanos)
	x = strings.TrimSuffix(x, "000")
	x = strings.TrimSuffix(x, "000")
	x = strings.TrimSuffix(x, ".000")
	e.WriteString(x + "s")
	return nil
}

func (d decoder) unmarshalDuration(m protoreflect.Message) error {
	tok, err := d.Read()
	if err != nil {
		return err
	}
	if tok.Kind() != json.String {
		return d.unexpectedTokenError(tok)
	}

	secs, nanos, ok := parseDuration(tok.ParsedString())
	if !ok {
		return d.newError(tok.Pos(), "invalid %v value %v", genid.Duration_message_fullname, tok.RawString())
	}
	
	
	if secs < -maxSecondsInDuration || secs > maxSecondsInDuration {
		return d.newError(tok.Pos(), "%v value out of range: %v", genid.Duration_message_fullname, tok.RawString())
	}

	fds := m.Descriptor().Fields()
	fdSeconds := fds.ByNumber(genid.Duration_Seconds_field_number)
	fdNanos := fds.ByNumber(genid.Duration_Nanos_field_number)

	m.Set(fdSeconds, protoreflect.ValueOfInt64(secs))
	m.Set(fdNanos, protoreflect.ValueOfInt32(nanos))
	return nil
}







func parseDuration(input string) (int64, int32, bool) {
	b := []byte(input)
	size := len(b)
	if size < 2 {
		return 0, 0, false
	}
	if b[size-1] != 's' {
		return 0, 0, false
	}
	b = b[:size-1]

	
	var neg bool
	switch b[0] {
	case '-':
		neg = true
		b = b[1:]
	case '+':
		b = b[1:]
	}
	if len(b) == 0 {
		return 0, 0, false
	}

	
	var intp []byte
	switch {
	case b[0] == '0':
		b = b[1:]

	case '1' <= b[0] && b[0] <= '9':
		intp = b[0:]
		b = b[1:]
		n := 1
		for len(b) > 0 && '0' <= b[0] && b[0] <= '9' {
			n++
			b = b[1:]
		}
		intp = intp[:n]

	case b[0] == '.':
		

	default:
		return 0, 0, false
	}

	hasFrac := false
	var frac [9]byte
	if len(b) > 0 {
		if b[0] != '.' {
			return 0, 0, false
		}
		
		b = b[1:]
		n := 0
		for len(b) > 0 && n < 9 && '0' <= b[0] && b[0] <= '9' {
			frac[n] = b[0]
			n++
			b = b[1:]
		}
		
		if len(b) > 0 {
			return 0, 0, false
		}
		
		for i := n; i < 9; i++ {
			frac[i] = '0'
		}
		hasFrac = true
	}

	var secs int64
	if len(intp) > 0 {
		var err error
		secs, err = strconv.ParseInt(string(intp), 10, 64)
		if err != nil {
			return 0, 0, false
		}
	}

	var nanos int64
	if hasFrac {
		nanob := bytes.TrimLeft(frac[:], "0")
		if len(nanob) > 0 {
			var err error
			nanos, err = strconv.ParseInt(string(nanob), 10, 32)
			if err != nil {
				return 0, 0, false
			}
		}
	}

	if neg {
		if secs > 0 {
			secs = -secs
		}
		if nanos > 0 {
			nanos = -nanos
		}
	}
	return secs, int32(nanos), true
}














const (
	maxTimestampSeconds = 253402300799
	minTimestampSeconds = -62135596800
)

func (e encoder) marshalTimestamp(m protoreflect.Message) error {
	fds := m.Descriptor().Fields()
	fdSeconds := fds.ByNumber(genid.Timestamp_Seconds_field_number)
	fdNanos := fds.ByNumber(genid.Timestamp_Nanos_field_number)

	secsVal := m.Get(fdSeconds)
	nanosVal := m.Get(fdNanos)
	secs := secsVal.Int()
	nanos := nanosVal.Int()
	if secs < minTimestampSeconds || secs > maxTimestampSeconds {
		return errors.New("%s: seconds out of range %v", genid.Timestamp_message_fullname, secs)
	}
	if nanos < 0 || nanos > secondsInNanos {
		return errors.New("%s: nanos out of range %v", genid.Timestamp_message_fullname, nanos)
	}
	
	
	t := time.Unix(secs, nanos).UTC()
	x := t.Format("2006-01-02T15:04:05.000000000")
	x = strings.TrimSuffix(x, "000")
	x = strings.TrimSuffix(x, "000")
	x = strings.TrimSuffix(x, ".000")
	e.WriteString(x + "Z")
	return nil
}

func (d decoder) unmarshalTimestamp(m protoreflect.Message) error {
	tok, err := d.Read()
	if err != nil {
		return err
	}
	if tok.Kind() != json.String {
		return d.unexpectedTokenError(tok)
	}

	s := tok.ParsedString()
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return d.newError(tok.Pos(), "invalid %v value %v", genid.Timestamp_message_fullname, tok.RawString())
	}
	
	secs := t.Unix()
	if secs < minTimestampSeconds || secs > maxTimestampSeconds {
		return d.newError(tok.Pos(), "%v value out of range: %v", genid.Timestamp_message_fullname, tok.RawString())
	}
	
	i := strings.LastIndexByte(s, '.')  
	j := strings.LastIndexAny(s, "Z-+") 
	if i >= 0 && j >= i && j-i > len(".999999999") {
		return d.newError(tok.Pos(), "invalid %v value %v", genid.Timestamp_message_fullname, tok.RawString())
	}

	fds := m.Descriptor().Fields()
	fdSeconds := fds.ByNumber(genid.Timestamp_Seconds_field_number)
	fdNanos := fds.ByNumber(genid.Timestamp_Nanos_field_number)

	m.Set(fdSeconds, protoreflect.ValueOfInt64(secs))
	m.Set(fdNanos, protoreflect.ValueOfInt32(int32(t.Nanosecond())))
	return nil
}






func (e encoder) marshalFieldMask(m protoreflect.Message) error {
	fd := m.Descriptor().Fields().ByNumber(genid.FieldMask_Paths_field_number)
	list := m.Get(fd).List()
	paths := make([]string, 0, list.Len())

	for i := 0; i < list.Len(); i++ {
		s := list.Get(i).String()
		if !protoreflect.FullName(s).IsValid() {
			return errors.New("%s contains invalid path: %q", genid.FieldMask_Paths_field_fullname, s)
		}
		
		cc := strs.JSONCamelCase(s)
		if s != strs.JSONSnakeCase(cc) {
			return errors.New("%s contains irreversible value %q", genid.FieldMask_Paths_field_fullname, s)
		}
		paths = append(paths, cc)
	}

	e.WriteString(strings.Join(paths, ","))
	return nil
}

func (d decoder) unmarshalFieldMask(m protoreflect.Message) error {
	tok, err := d.Read()
	if err != nil {
		return err
	}
	if tok.Kind() != json.String {
		return d.unexpectedTokenError(tok)
	}
	str := strings.TrimSpace(tok.ParsedString())
	if str == "" {
		return nil
	}
	paths := strings.Split(str, ",")

	fd := m.Descriptor().Fields().ByNumber(genid.FieldMask_Paths_field_number)
	list := m.Mutable(fd).List()

	for _, s0 := range paths {
		s := strs.JSONSnakeCase(s0)
		if strings.Contains(s0, "_") || !protoreflect.FullName(s).IsValid() {
			return d.newError(tok.Pos(), "%v contains invalid path: %q", genid.FieldMask_Paths_field_fullname, s0)
		}
		list.Append(protoreflect.ValueOfString(s))
	}
	return nil
}
