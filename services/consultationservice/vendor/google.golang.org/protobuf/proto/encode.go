



package proto

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/internal/encoding/messageset"
	"google.golang.org/protobuf/internal/order"
	"google.golang.org/protobuf/internal/pragma"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoiface"

	protoerrors "google.golang.org/protobuf/internal/errors"
)






type MarshalOptions struct {
	pragma.NoUnkeyedLiterals

	
	
	
	AllowPartial bool

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	Deterministic bool

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	UseCachedSize bool
}





func (o MarshalOptions) flags() protoiface.MarshalInputFlags {
	var flags protoiface.MarshalInputFlags

	
	

	if o.Deterministic {
		flags |= protoiface.MarshalDeterministic
	}

	if o.UseCachedSize {
		flags |= protoiface.MarshalUseCachedSize
	}

	return flags
}






func Marshal(m Message) ([]byte, error) {
	
	if m == nil {
		return nil, nil
	}

	out, err := MarshalOptions{}.marshal(nil, m.ProtoReflect())
	if len(out.Buf) == 0 && err == nil {
		out.Buf = emptyBytesForMessage(m)
	}
	return out.Buf, err
}


func (o MarshalOptions) Marshal(m Message) ([]byte, error) {
	
	if m == nil {
		return nil, nil
	}

	out, err := o.marshal(nil, m.ProtoReflect())
	if len(out.Buf) == 0 && err == nil {
		out.Buf = emptyBytesForMessage(m)
	}
	return out.Buf, err
}










func emptyBytesForMessage(m Message) []byte {
	if m == nil || !m.ProtoReflect().IsValid() {
		return nil
	}
	return emptyBuf[:]
}






func (o MarshalOptions) MarshalAppend(b []byte, m Message) ([]byte, error) {
	
	if m == nil {
		return b, nil
	}

	out, err := o.marshal(b, m.ProtoReflect())
	return out.Buf, err
}





func (o MarshalOptions) MarshalState(in protoiface.MarshalInput) (protoiface.MarshalOutput, error) {
	return o.marshal(in.Buf, in.Message)
}




func (o MarshalOptions) marshal(b []byte, m protoreflect.Message) (out protoiface.MarshalOutput, err error) {
	allowPartial := o.AllowPartial
	o.AllowPartial = true
	if methods := protoMethods(m); methods != nil && methods.Marshal != nil &&
		!(o.Deterministic && methods.Flags&protoiface.SupportMarshalDeterministic == 0) {
		in := protoiface.MarshalInput{
			Message: m,
			Buf:     b,
			Flags:   o.flags(),
		}
		if methods.Size != nil {
			sout := methods.Size(protoiface.SizeInput{
				Message: m,
				Flags:   in.Flags,
			})
			if cap(b) < len(b)+sout.Size {
				in.Buf = make([]byte, len(b), growcap(cap(b), len(b)+sout.Size))
				copy(in.Buf, b)
			}
			in.Flags |= protoiface.MarshalUseCachedSize
		}
		out, err = methods.Marshal(in)
	} else {
		out.Buf, err = o.marshalMessageSlow(b, m)
	}
	if err != nil {
		var mismatch *protoerrors.SizeMismatchError
		if errors.As(err, &mismatch) {
			return out, fmt.Errorf("marshaling %s: %v", string(m.Descriptor().FullName()), err)
		}
		return out, err
	}
	if allowPartial {
		return out, nil
	}
	return out, checkInitialized(m)
}

func (o MarshalOptions) marshalMessage(b []byte, m protoreflect.Message) ([]byte, error) {
	out, err := o.marshal(b, m)
	return out.Buf, err
}







func growcap(oldcap, wantcap int) (newcap int) {
	if wantcap > oldcap*2 {
		newcap = wantcap
	} else if oldcap < 1024 {
		
		
		
		newcap = oldcap * 2
	} else {
		newcap = oldcap
		for 0 < newcap && newcap < wantcap {
			newcap += newcap / 4
		}
		if newcap <= 0 {
			newcap = wantcap
		}
	}
	return newcap
}

func (o MarshalOptions) marshalMessageSlow(b []byte, m protoreflect.Message) ([]byte, error) {
	if messageset.IsMessageSet(m.Descriptor()) {
		return o.marshalMessageSet(b, m)
	}
	fieldOrder := order.AnyFieldOrder
	if o.Deterministic {
		
		
		
		fieldOrder = order.LegacyFieldOrder
	}
	var err error
	order.RangeFields(m, fieldOrder, func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		b, err = o.marshalField(b, fd, v)
		return err == nil
	})
	if err != nil {
		return b, err
	}
	b = append(b, m.GetUnknown()...)
	return b, nil
}

func (o MarshalOptions) marshalField(b []byte, fd protoreflect.FieldDescriptor, value protoreflect.Value) ([]byte, error) {
	switch {
	case fd.IsList():
		return o.marshalList(b, fd, value.List())
	case fd.IsMap():
		return o.marshalMap(b, fd, value.Map())
	default:
		b = protowire.AppendTag(b, fd.Number(), wireTypes[fd.Kind()])
		return o.marshalSingular(b, fd, value)
	}
}

func (o MarshalOptions) marshalList(b []byte, fd protoreflect.FieldDescriptor, list protoreflect.List) ([]byte, error) {
	if fd.IsPacked() && list.Len() > 0 {
		b = protowire.AppendTag(b, fd.Number(), protowire.BytesType)
		b, pos := appendSpeculativeLength(b)
		for i, llen := 0, list.Len(); i < llen; i++ {
			var err error
			b, err = o.marshalSingular(b, fd, list.Get(i))
			if err != nil {
				return b, err
			}
		}
		b = finishSpeculativeLength(b, pos)
		return b, nil
	}

	kind := fd.Kind()
	for i, llen := 0, list.Len(); i < llen; i++ {
		var err error
		b = protowire.AppendTag(b, fd.Number(), wireTypes[kind])
		b, err = o.marshalSingular(b, fd, list.Get(i))
		if err != nil {
			return b, err
		}
	}
	return b, nil
}

func (o MarshalOptions) marshalMap(b []byte, fd protoreflect.FieldDescriptor, mapv protoreflect.Map) ([]byte, error) {
	keyf := fd.MapKey()
	valf := fd.MapValue()
	keyOrder := order.AnyKeyOrder
	if o.Deterministic {
		keyOrder = order.GenericKeyOrder
	}
	var err error
	order.RangeEntries(mapv, keyOrder, func(key protoreflect.MapKey, value protoreflect.Value) bool {
		b = protowire.AppendTag(b, fd.Number(), protowire.BytesType)
		var pos int
		b, pos = appendSpeculativeLength(b)

		b, err = o.marshalField(b, keyf, key.Value())
		if err != nil {
			return false
		}
		b, err = o.marshalField(b, valf, value)
		if err != nil {
			return false
		}
		b = finishSpeculativeLength(b, pos)
		return true
	})
	return b, err
}




const speculativeLength = 1

func appendSpeculativeLength(b []byte) ([]byte, int) {
	pos := len(b)
	b = append(b, "\x00\x00\x00\x00"[:speculativeLength]...)
	return b, pos
}

func finishSpeculativeLength(b []byte, pos int) []byte {
	mlen := len(b) - pos - speculativeLength
	msiz := protowire.SizeVarint(uint64(mlen))
	if msiz != speculativeLength {
		for i := 0; i < msiz-speculativeLength; i++ {
			b = append(b, 0)
		}
		copy(b[pos+msiz:], b[pos+speculativeLength:])
		b = b[:pos+msiz+mlen]
	}
	protowire.AppendVarint(b[:pos], uint64(mlen))
	return b
}
