



package impl

import (
	"fmt"
	"math/bits"
	"os"
	"reflect"
	"sort"
	"sync/atomic"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/internal/errors"
	"google.golang.org/protobuf/internal/protolazy"
	"google.golang.org/protobuf/reflect/protoreflect"
	preg "google.golang.org/protobuf/reflect/protoregistry"
	piface "google.golang.org/protobuf/runtime/protoiface"
)

var enableLazy int32 = func() int32 {
	if os.Getenv("GOPROTODEBUG") == "nolazy" {
		return 0
	}
	return 1
}()


func EnableLazyUnmarshal(enable bool) {
	if enable {
		atomic.StoreInt32(&enableLazy, 1)
		return
	}
	atomic.StoreInt32(&enableLazy, 0)
}


func LazyEnabled() bool {
	return atomic.LoadInt32(&enableLazy) != 0
}


func UnmarshalField(m interface{}, num protowire.Number) {
	switch m := m.(type) {
	case *messageState:
		m.messageInfo().lazyUnmarshal(m.pointer(), num)
	case *messageReflectWrapper:
		m.messageInfo().lazyUnmarshal(m.pointer(), num)
	default:
		panic(fmt.Sprintf("unsupported wrapper type %T", m))
	}
}

func (mi *MessageInfo) lazyUnmarshal(p pointer, num protoreflect.FieldNumber) {
	var f *coderFieldInfo
	if int(num) < len(mi.denseCoderFields) {
		f = mi.denseCoderFields[num]
	} else {
		f = mi.coderFields[num]
	}
	if f == nil {
		panic(fmt.Sprintf("lazyUnmarshal: field info for %v.%v", mi.Desc.FullName(), num))
	}
	lazy := *p.Apply(mi.lazyOffset).LazyInfoPtr()
	start, end, found, _, multipleEntries := lazy.FindFieldInProto(uint32(num))
	if !found && multipleEntries == nil {
		panic(fmt.Sprintf("lazyUnmarshal: can't find field data for %v.%v", mi.Desc.FullName(), num))
	}
	
	
	fp := pointerOfValue(reflect.New(f.ft))
	if multipleEntries != nil {
		for _, entry := range multipleEntries {
			mi.unmarshalField(lazy.Buffer()[entry.Start:entry.End], fp, f, lazy, lazy.UnmarshalFlags())
		}
	} else {
		mi.unmarshalField(lazy.Buffer()[start:end], fp, f, lazy, lazy.UnmarshalFlags())
	}
	p.Apply(f.offset).AtomicSetPointerIfNil(fp.Elem())
}

func (mi *MessageInfo) unmarshalField(b []byte, p pointer, f *coderFieldInfo, lazyInfo *protolazy.XXX_lazyUnmarshalInfo, flags piface.UnmarshalInputFlags) error {
	opts := lazyUnmarshalOptions
	opts.flags |= flags
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
				return errors.New("invalid wire data")
			}
			b = b[n:]
		}
		var num protowire.Number
		if n := tag >> 3; n < uint64(protowire.MinValidNumber) || n > uint64(protowire.MaxValidNumber) {
			return errors.New("invalid wire data")
		} else {
			num = protowire.Number(n)
		}
		wtyp := protowire.Type(tag & 7)
		if num == f.num {
			o, err := f.funcs.unmarshal(b, p, wtyp, f, opts)
			if err == nil {
				b = b[o.n:]
				continue
			}
			if err != errUnknown {
				return err
			}
		}
		n := protowire.ConsumeFieldValue(num, wtyp, b)
		if n < 0 {
			return errors.New("invalid wire data")
		}
		b = b[n:]
	}
	return nil
}

func (mi *MessageInfo) skipField(b []byte, f *coderFieldInfo, wtyp protowire.Type, opts unmarshalOptions) (out unmarshalOutput, _ ValidationStatus) {
	fmi := f.validation.mi
	if fmi == nil {
		fd := mi.Desc.Fields().ByNumber(f.num)
		if fd == nil {
			return out, ValidationUnknown
		}
		messageName := fd.Message().FullName()
		messageType, err := preg.GlobalTypes.FindMessageByName(messageName)
		if err != nil {
			return out, ValidationUnknown
		}
		var ok bool
		fmi, ok = messageType.(*MessageInfo)
		if !ok {
			return out, ValidationUnknown
		}
	}
	fmi.init()
	switch f.validation.typ {
	case validationTypeMessage:
		if wtyp != protowire.BytesType {
			return out, ValidationWrongWireType
		}
		v, n := protowire.ConsumeBytes(b)
		if n < 0 {
			return out, ValidationInvalid
		}
		out, st := fmi.validate(v, 0, opts)
		out.n = n
		return out, st
	case validationTypeGroup:
		if wtyp != protowire.StartGroupType {
			return out, ValidationWrongWireType
		}
		out, st := fmi.validate(b, f.num, opts)
		return out, st
	default:
		return out, ValidationUnknown
	}
}




func (mi *MessageInfo) unmarshalPointerLazy(b []byte, p pointer, groupTag protowire.Number, opts unmarshalOptions) (out unmarshalOutput, err error) {
	initialized := true
	var requiredMask uint64
	var lazy **protolazy.XXX_lazyUnmarshalInfo
	var presence presence
	var lazyIndex []protolazy.IndexEntry
	var lastNum protowire.Number
	outOfOrder := false
	lazyDecode := false
	presence = p.Apply(mi.presenceOffset).PresenceInfo()
	lazy = p.Apply(mi.lazyOffset).LazyInfoPtr()
	if !presence.AnyPresent(mi.presenceSize) {
		if opts.CanBeLazy() {
			
			
			
			lazyDecode = true
			if *lazy == nil {
				*lazy = &protolazy.XXX_lazyUnmarshalInfo{}
			}
			(*lazy).SetUnmarshalFlags(opts.flags)
			if !opts.AliasBuffer() {
				
				
				
				b = append([]byte{}, b...)
				opts.flags |= piface.UnmarshalAliasBuffer
			}
			(*lazy).SetBuffer(b)
		}
	}
	
	
	
	
	type lazyAction uint8
	const (
		lazyValidateOnly   lazyAction = iota 
		lazyUnmarshalNow                     
		lazyUnmarshalLater                   
	)
	var lazyFields map[*coderFieldInfo]lazyAction
	var exts *map[int32]ExtensionField
	start := len(b)
	pos := 0
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
				return out, errDecode
			}
			b = b[n:]
		}
		var num protowire.Number
		if n := tag >> 3; n < uint64(protowire.MinValidNumber) || n > uint64(protowire.MaxValidNumber) {
			return out, errors.New("invalid field number")
		} else {
			num = protowire.Number(n)
		}
		wtyp := protowire.Type(tag & 7)

		if wtyp == protowire.EndGroupType {
			if num != groupTag {
				return out, errors.New("mismatching end group marker")
			}
			groupTag = 0
			break
		}

		var f *coderFieldInfo
		if int(num) < len(mi.denseCoderFields) {
			f = mi.denseCoderFields[num]
		} else {
			f = mi.coderFields[num]
		}
		var n int
		err := errUnknown
		discardUnknown := false
	Field:
		switch {
		case f != nil:
			if f.funcs.unmarshal == nil {
				break
			}
			if f.isLazy && lazyDecode {
				switch {
				case lazyFields == nil || lazyFields[f] == lazyValidateOnly:
					
					o, valid := mi.skipField(b, f, wtyp, opts)
					switch valid {
					case ValidationValid:
						
						err = nil
						presence.SetPresentUnatomic(f.presenceIndex, mi.presenceSize)
						requiredMask |= f.validation.requiredBit
						if !o.initialized {
							initialized = false
						}
						n = o.n
						break Field
					case ValidationInvalid:
						return out, errors.New("invalid proto wire format")
					case ValidationWrongWireType:
						break Field
					case ValidationUnknown:
						if lazyFields == nil {
							lazyFields = make(map[*coderFieldInfo]lazyAction)
						}
						if presence.Present(f.presenceIndex) {
							
							
							
							
							
							presence.ClearPresent(f.presenceIndex)
							lazyFields[f] = lazyUnmarshalLater
							discardUnknown = true
							break Field
						} else {
							
							
							
							lazyFields[f] = lazyUnmarshalNow
						}
					}
				case lazyFields[f] == lazyUnmarshalLater:
					
					
					discardUnknown = true
					break Field
				default:
					
				}
			}
			if f.isLazy && !lazyDecode && presence.Present(f.presenceIndex) {
				if p.Apply(f.offset).AtomicGetPointer().IsNil() {
					mi.lazyUnmarshal(p, f.num)
				}
			}
			var o unmarshalOutput
			o, err = f.funcs.unmarshal(b, p.Apply(f.offset), wtyp, f, opts)
			n = o.n
			if err != nil {
				break
			}
			requiredMask |= f.validation.requiredBit
			if f.funcs.isInit != nil && !o.initialized {
				initialized = false
			}
			if f.presenceIndex != noPresence {
				presence.SetPresentUnatomic(f.presenceIndex, mi.presenceSize)
			}
		default:
			
			if exts == nil && mi.extensionOffset.IsValid() {
				exts = p.Apply(mi.extensionOffset).Extensions()
				if *exts == nil {
					*exts = make(map[int32]ExtensionField)
				}
			}
			if exts == nil {
				break
			}
			var o unmarshalOutput
			o, err = mi.unmarshalExtension(b, num, wtyp, *exts, opts)
			if err != nil {
				break
			}
			n = o.n
			if !o.initialized {
				initialized = false
			}
		}
		if err != nil {
			if err != errUnknown {
				return out, err
			}
			n = protowire.ConsumeFieldValue(num, wtyp, b)
			if n < 0 {
				return out, errDecode
			}
			if !discardUnknown && !opts.DiscardUnknown() && mi.unknownOffset.IsValid() {
				u := mi.mutableUnknownBytes(p)
				*u = protowire.AppendTag(*u, num, wtyp)
				*u = append(*u, b[:n]...)
			}
		}
		b = b[n:]
		end := start - len(b)
		if lazyDecode && f != nil && f.isLazy {
			if num != lastNum {
				lazyIndex = append(lazyIndex, protolazy.IndexEntry{
					FieldNum: uint32(num),
					Start:    uint32(pos),
					End:      uint32(end),
				})
			} else {
				i := len(lazyIndex) - 1
				lazyIndex[i].End = uint32(end)
				lazyIndex[i].MultipleContiguous = true
			}
		}
		if num < lastNum {
			outOfOrder = true
		}
		pos = end
		lastNum = num
	}
	if groupTag != 0 {
		return out, errors.New("missing end group marker")
	}
	if lazyFields != nil {
		
		for f, action := range lazyFields {
			if action != lazyUnmarshalLater {
				continue
			}
			initialized = false
			if *lazy == nil {
				*lazy = &protolazy.XXX_lazyUnmarshalInfo{}
			}
			if err := mi.unmarshalField((*lazy).Buffer(), p.Apply(f.offset), f, *lazy, opts.flags); err != nil {
				return out, err
			}
			presence.SetPresentUnatomic(f.presenceIndex, mi.presenceSize)
		}
	}
	if lazyDecode {
		if outOfOrder {
			sort.Slice(lazyIndex, func(i, j int) bool {
				return lazyIndex[i].FieldNum < lazyIndex[j].FieldNum ||
					(lazyIndex[i].FieldNum == lazyIndex[j].FieldNum &&
						lazyIndex[i].Start < lazyIndex[j].Start)
			})
		}
		if *lazy == nil {
			*lazy = &protolazy.XXX_lazyUnmarshalInfo{}
		}

		(*lazy).SetIndex(lazyIndex)
	}
	if mi.numRequiredFields > 0 && bits.OnesCount64(requiredMask) != int(mi.numRequiredFields) {
		initialized = false
	}
	if initialized {
		out.initialized = true
	}
	out.n = start - len(b)
	return out, nil
}
