



package impl

import (
	"sync"
	"sync/atomic"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/internal/errors"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type extensionFieldInfo struct {
	wiretag             uint64
	tagsize             int
	unmarshalNeedsValue bool
	funcs               valueCoderFuncs
	validation          validationInfo
}

func getExtensionFieldInfo(xt protoreflect.ExtensionType) *extensionFieldInfo {
	if xi, ok := xt.(*ExtensionInfo); ok {
		xi.lazyInit()
		return xi.info
	}
	
	
	
	
	
	
	return makeExtensionFieldInfo(xt.TypeDescriptor())
}

func makeExtensionFieldInfo(xd protoreflect.ExtensionDescriptor) *extensionFieldInfo {
	var wiretag uint64
	if !xd.IsPacked() {
		wiretag = protowire.EncodeTag(xd.Number(), wireTypes[xd.Kind()])
	} else {
		wiretag = protowire.EncodeTag(xd.Number(), protowire.BytesType)
	}
	e := &extensionFieldInfo{
		wiretag: wiretag,
		tagsize: protowire.SizeVarint(wiretag),
		funcs:   encoderFuncsForValue(xd),
	}
	
	
	
	switch xd.Kind() {
	case protoreflect.MessageKind, protoreflect.GroupKind, protoreflect.EnumKind:
		e.unmarshalNeedsValue = true
	default:
		if xd.Cardinality() == protoreflect.Repeated {
			e.unmarshalNeedsValue = true
		}
	}
	return e
}

type lazyExtensionValue struct {
	atomicOnce uint32 
	mu         sync.Mutex
	xi         *extensionFieldInfo
	value      protoreflect.Value
	b          []byte
}

type ExtensionField struct {
	typ protoreflect.ExtensionType

	
	
	value protoreflect.Value
	lazy  *lazyExtensionValue
}

func (f *ExtensionField) appendLazyBytes(xt protoreflect.ExtensionType, xi *extensionFieldInfo, num protowire.Number, wtyp protowire.Type, b []byte) {
	if f.lazy == nil {
		f.lazy = &lazyExtensionValue{xi: xi}
	}
	f.typ = xt
	f.lazy.xi = xi
	f.lazy.b = protowire.AppendTag(f.lazy.b, num, wtyp)
	f.lazy.b = append(f.lazy.b, b...)
}

func (f *ExtensionField) canLazy(xt protoreflect.ExtensionType) bool {
	if f.typ == nil {
		return true
	}
	if f.typ == xt && f.lazy != nil && atomic.LoadUint32(&f.lazy.atomicOnce) == 0 {
		return true
	}
	return false
}




func (f *ExtensionField) isUnexpandedLazy() bool {
	return f.lazy != nil && atomic.LoadUint32(&f.lazy.atomicOnce) == 0
}





func (f *ExtensionField) lazyBuffer() []byte {
	
	
	if !f.isUnexpandedLazy() {
		return nil
	}
	f.lazy.mu.Lock()
	defer f.lazy.mu.Unlock()
	return f.lazy.b
}

func (f *ExtensionField) lazyInit() {
	f.lazy.mu.Lock()
	defer f.lazy.mu.Unlock()
	if atomic.LoadUint32(&f.lazy.atomicOnce) == 1 {
		return
	}
	if f.lazy.xi != nil {
		b := f.lazy.b
		val := f.typ.New()
		for len(b) > 0 {
			var tag uint64
			if b[0] < 0x80 {
				tag = uint64(b[0])
				b = b[1:]
			} else if len(b) >= 2 && b[1] < 128 {
				tag = uint64(b[0]&0x7f) + uint64(b[1])<<7
				b = b[2:]
			} else {
				var n int
				tag, n = protowire.ConsumeVarint(b)
				if n < 0 {
					panic(errors.New("bad tag in lazy extension decoding"))
				}
				b = b[n:]
			}
			num := protowire.Number(tag >> 3)
			wtyp := protowire.Type(tag & 7)
			var out unmarshalOutput
			var err error
			val, out, err = f.lazy.xi.funcs.unmarshal(b, val, num, wtyp, lazyUnmarshalOptions)
			if err != nil {
				panic(errors.New("decode failure in lazy extension decoding: %v", err))
			}
			b = b[out.n:]
		}
		f.lazy.value = val
	} else {
		panic("No support for lazy fns for ExtensionField")
	}
	f.lazy.xi = nil
	f.lazy.b = nil
	atomic.StoreUint32(&f.lazy.atomicOnce, 1)
}



func (f *ExtensionField) Set(t protoreflect.ExtensionType, v protoreflect.Value) {
	f.typ = t
	f.value = v
	f.lazy = nil
}



func (f *ExtensionField) Value() protoreflect.Value {
	if f.lazy != nil {
		if atomic.LoadUint32(&f.lazy.atomicOnce) == 0 {
			f.lazyInit()
		}
		return f.lazy.value
	}
	return f.value
}



func (f ExtensionField) Type() protoreflect.ExtensionType {
	return f.typ
}



func (f ExtensionField) IsSet() bool {
	return f.typ != nil
}



func IsLazy(m protoreflect.Message, fd protoreflect.FieldDescriptor) bool {
	var mi *MessageInfo
	var p pointer
	switch m := m.(type) {
	case *messageState:
		mi = m.messageInfo()
		p = m.pointer()
	case *messageReflectWrapper:
		mi = m.messageInfo()
		p = m.pointer()
	default:
		return false
	}
	xd, ok := fd.(protoreflect.ExtensionTypeDescriptor)
	if !ok {
		return false
	}
	xt := xd.Type()
	ext := mi.extensionMap(p)
	if ext == nil {
		return false
	}
	f, ok := (*ext)[int32(fd.Number())]
	if !ok {
		return false
	}
	return f.typ == xt && f.lazy != nil && atomic.LoadUint32(&f.lazy.atomicOnce) == 0
}
