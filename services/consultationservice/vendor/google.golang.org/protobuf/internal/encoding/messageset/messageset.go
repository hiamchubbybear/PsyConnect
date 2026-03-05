




package messageset

import (
	"math"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/internal/errors"
	"google.golang.org/protobuf/reflect/protoreflect"
)












const (
	FieldItem    = protowire.Number(1)
	FieldTypeID  = protowire.Number(2)
	FieldMessage = protowire.Number(3)
)











const ExtensionName = "message_set_extension"


func IsMessageSet(md protoreflect.MessageDescriptor) bool {
	xmd, ok := md.(interface{ IsMessageSet() bool })
	return ok && xmd.IsMessageSet()
}


func IsMessageSetExtension(fd protoreflect.FieldDescriptor) bool {
	switch {
	case fd.Name() != ExtensionName:
		return false
	case !IsMessageSet(fd.ContainingMessage()):
		return false
	case fd.FullName().Parent() != fd.Message().FullName():
		return false
	}
	return true
}



func SizeField(num protowire.Number) int {
	return 2*protowire.SizeTag(FieldItem) + protowire.SizeTag(FieldTypeID) + protowire.SizeVarint(uint64(num))
}








func Unmarshal(b []byte, wantLen bool, fn func(typeID protowire.Number, value []byte) error) error {
	for len(b) > 0 {
		num, wtyp, n := protowire.ConsumeTag(b)
		if n < 0 {
			return protowire.ParseError(n)
		}
		b = b[n:]
		if num != FieldItem || wtyp != protowire.StartGroupType {
			n := protowire.ConsumeFieldValue(num, wtyp, b)
			if n < 0 {
				return protowire.ParseError(n)
			}
			b = b[n:]
			continue
		}
		typeID, value, n, err := ConsumeFieldValue(b, wantLen)
		if err != nil {
			return err
		}
		b = b[n:]
		if typeID == 0 {
			continue
		}
		if err := fn(typeID, value); err != nil {
			return err
		}
	}
	return nil
}







func ConsumeFieldValue(b []byte, wantLen bool) (typeid protowire.Number, message []byte, n int, err error) {
	ilen := len(b)
	for {
		num, wtyp, n := protowire.ConsumeTag(b)
		if n < 0 {
			return 0, nil, 0, protowire.ParseError(n)
		}
		b = b[n:]
		switch {
		case num == FieldItem && wtyp == protowire.EndGroupType:
			if wantLen && len(message) == 0 {
				
				
				message = protowire.AppendVarint(message, 0)
			}
			return typeid, message, ilen - len(b), nil
		case num == FieldTypeID && wtyp == protowire.VarintType:
			v, n := protowire.ConsumeVarint(b)
			if n < 0 {
				return 0, nil, 0, protowire.ParseError(n)
			}
			b = b[n:]
			if v < 1 || v > math.MaxInt32 {
				return 0, nil, 0, errors.New("invalid type_id in message set")
			}
			typeid = protowire.Number(v)
		case num == FieldMessage && wtyp == protowire.BytesType:
			m, n := protowire.ConsumeBytes(b)
			if n < 0 {
				return 0, nil, 0, protowire.ParseError(n)
			}
			if message == nil {
				if wantLen {
					message = b[:n:n]
				} else {
					message = m[:len(m):len(m)]
				}
			} else {
				
				
				
				
				
				
				
				if wantLen {
					_, nn := protowire.ConsumeVarint(message)
					m0 := message[nn:]
					message = nil
					message = protowire.AppendVarint(message, uint64(len(m0)+len(m)))
					message = append(message, m0...)
					message = append(message, m...)
				} else {
					message = append(message, m...)
				}
			}
			b = b[n:]
		default:
			
			n := protowire.ConsumeFieldValue(num, wtyp, b)
			if n < 0 {
				return 0, nil, 0, protowire.ParseError(n)
			}
			b = b[n:]
		}
	}
}




func AppendFieldStart(b []byte, num protowire.Number) []byte {
	b = protowire.AppendTag(b, FieldItem, protowire.StartGroupType)
	b = protowire.AppendTag(b, FieldTypeID, protowire.VarintType)
	b = protowire.AppendVarint(b, uint64(num))
	return b
}


func AppendFieldEnd(b []byte) []byte {
	return protowire.AppendTag(b, FieldItem, protowire.EndGroupType)
}




func SizeUnknown(unknown []byte) (size int) {
	for len(unknown) > 0 {
		num, typ, n := protowire.ConsumeTag(unknown)
		if n < 0 || typ != protowire.BytesType {
			return 0
		}
		unknown = unknown[n:]
		_, n = protowire.ConsumeBytes(unknown)
		if n < 0 {
			return 0
		}
		unknown = unknown[n:]
		size += SizeField(num) + protowire.SizeTag(FieldMessage) + n
	}
	return size
}









func AppendUnknown(b, unknown []byte) ([]byte, error) {
	for len(unknown) > 0 {
		num, typ, n := protowire.ConsumeTag(unknown)
		if n < 0 || typ != protowire.BytesType {
			return nil, errors.New("invalid data in message set unknown fields")
		}
		unknown = unknown[n:]
		_, n = protowire.ConsumeBytes(unknown)
		if n < 0 {
			return nil, errors.New("invalid data in message set unknown fields")
		}
		b = AppendFieldStart(b, num)
		b = protowire.AppendTag(b, FieldMessage, protowire.BytesType)
		b = append(b, unknown[:n]...)
		b = AppendFieldEnd(b)
		unknown = unknown[n:]
	}
	return b, nil
}
