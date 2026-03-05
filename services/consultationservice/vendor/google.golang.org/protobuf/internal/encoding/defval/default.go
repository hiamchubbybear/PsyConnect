








package defval

import (
	"fmt"
	"math"
	"strconv"

	ptext "google.golang.org/protobuf/internal/encoding/text"
	"google.golang.org/protobuf/internal/errors"
	"google.golang.org/protobuf/reflect/protoreflect"
)


type Format int

const (
	_ Format = iota

	
	
	Descriptor

	
	GoTag
)



func Unmarshal(s string, k protoreflect.Kind, evs protoreflect.EnumValueDescriptors, f Format) (protoreflect.Value, protoreflect.EnumValueDescriptor, error) {
	switch k {
	case protoreflect.BoolKind:
		if f == GoTag {
			switch s {
			case "1":
				return protoreflect.ValueOfBool(true), nil, nil
			case "0":
				return protoreflect.ValueOfBool(false), nil, nil
			}
		} else {
			switch s {
			case "true":
				return protoreflect.ValueOfBool(true), nil, nil
			case "false":
				return protoreflect.ValueOfBool(false), nil, nil
			}
		}
	case protoreflect.EnumKind:
		if f == GoTag {
			
			if n, err := strconv.ParseInt(s, 10, 32); err == nil {
				if ev := evs.ByNumber(protoreflect.EnumNumber(n)); ev != nil {
					return protoreflect.ValueOfEnum(ev.Number()), ev, nil
				}
			}
		} else {
			
			ev := evs.ByName(protoreflect.Name(s))
			if ev != nil {
				return protoreflect.ValueOfEnum(ev.Number()), ev, nil
			}
		}
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		if v, err := strconv.ParseInt(s, 10, 32); err == nil {
			return protoreflect.ValueOfInt32(int32(v)), nil, nil
		}
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		if v, err := strconv.ParseInt(s, 10, 64); err == nil {
			return protoreflect.ValueOfInt64(int64(v)), nil, nil
		}
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		if v, err := strconv.ParseUint(s, 10, 32); err == nil {
			return protoreflect.ValueOfUint32(uint32(v)), nil, nil
		}
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		if v, err := strconv.ParseUint(s, 10, 64); err == nil {
			return protoreflect.ValueOfUint64(uint64(v)), nil, nil
		}
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		var v float64
		var err error
		switch s {
		case "-inf":
			v = math.Inf(-1)
		case "inf":
			v = math.Inf(+1)
		case "nan":
			v = math.NaN()
		default:
			v, err = strconv.ParseFloat(s, 64)
		}
		if err == nil {
			if k == protoreflect.FloatKind {
				return protoreflect.ValueOfFloat32(float32(v)), nil, nil
			} else {
				return protoreflect.ValueOfFloat64(float64(v)), nil, nil
			}
		}
	case protoreflect.StringKind:
		
		return protoreflect.ValueOfString(s), nil, nil
	case protoreflect.BytesKind:
		if b, ok := unmarshalBytes(s); ok {
			return protoreflect.ValueOfBytes(b), nil, nil
		}
	}
	return protoreflect.Value{}, nil, errors.New("could not parse value for %v: %q", k, s)
}




func Marshal(v protoreflect.Value, ev protoreflect.EnumValueDescriptor, k protoreflect.Kind, f Format) (string, error) {
	switch k {
	case protoreflect.BoolKind:
		if f == GoTag {
			if v.Bool() {
				return "1", nil
			} else {
				return "0", nil
			}
		} else {
			if v.Bool() {
				return "true", nil
			} else {
				return "false", nil
			}
		}
	case protoreflect.EnumKind:
		if f == GoTag {
			return strconv.FormatInt(int64(v.Enum()), 10), nil
		} else {
			return string(ev.Name()), nil
		}
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind, protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return strconv.FormatInt(v.Int(), 10), nil
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind, protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return strconv.FormatUint(v.Uint(), 10), nil
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		f := v.Float()
		switch {
		case math.IsInf(f, -1):
			return "-inf", nil
		case math.IsInf(f, +1):
			return "inf", nil
		case math.IsNaN(f):
			return "nan", nil
		default:
			if k == protoreflect.FloatKind {
				return strconv.FormatFloat(f, 'g', -1, 32), nil
			} else {
				return strconv.FormatFloat(f, 'g', -1, 64), nil
			}
		}
	case protoreflect.StringKind:
		
		return v.String(), nil
	case protoreflect.BytesKind:
		if s, ok := marshalBytes(v.Bytes()); ok {
			return s, nil
		}
	}
	return "", errors.New("could not format value for %v: %v", k, v)
}


func unmarshalBytes(s string) ([]byte, bool) {
	
	
	v, err := ptext.UnmarshalString(`"` + s + `"`)
	if err != nil {
		return nil, false
	}
	return []byte(v), true
}




func marshalBytes(b []byte) (string, bool) {
	var s []byte
	for _, c := range b {
		switch c {
		case '\n':
			s = append(s, `\n`...)
		case '\r':
			s = append(s, `\r`...)
		case '\t':
			s = append(s, `\t`...)
		case '"':
			s = append(s, `\"`...)
		case '\'':
			s = append(s, `\'`...)
		case '\\':
			s = append(s, `\\`...)
		default:
			if printableASCII := c >= 0x20 && c <= 0x7e; printableASCII {
				s = append(s, c)
			} else {
				s = append(s, fmt.Sprintf(`\%03o`, c)...)
			}
		}
	}
	return string(s), true
}
