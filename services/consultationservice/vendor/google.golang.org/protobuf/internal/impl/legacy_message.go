



package impl

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"google.golang.org/protobuf/internal/descopts"
	ptag "google.golang.org/protobuf/internal/encoding/tag"
	"google.golang.org/protobuf/internal/errors"
	"google.golang.org/protobuf/internal/filedesc"
	"google.golang.org/protobuf/internal/strs"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoiface"
)



func legacyWrapMessage(v reflect.Value) protoreflect.Message {
	t := v.Type()
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return aberrantMessage{v: v}
	}
	mt := legacyLoadMessageInfo(t, "")
	return mt.MessageOf(v.Interface())
}




func legacyLoadMessageType(t reflect.Type, name protoreflect.FullName) protoreflect.MessageType {
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return aberrantMessageType{t}
	}
	return legacyLoadMessageInfo(t, name)
}

var legacyMessageTypeCache sync.Map 




func legacyLoadMessageInfo(t reflect.Type, name protoreflect.FullName) *MessageInfo {
	
	if mt, ok := legacyMessageTypeCache.Load(t); ok {
		return mt.(*MessageInfo)
	}

	
	mi := &MessageInfo{
		Desc:          legacyLoadMessageDesc(t, name),
		GoReflectType: t,
	}

	var hasMarshal, hasUnmarshal bool
	v := reflect.Zero(t).Interface()
	if _, hasMarshal = v.(legacyMarshaler); hasMarshal {
		mi.methods.Marshal = legacyMarshal

		
		
		
		
		mi.methods.Flags |= protoiface.SupportMarshalDeterministic
	}
	if _, hasUnmarshal = v.(legacyUnmarshaler); hasUnmarshal {
		mi.methods.Unmarshal = legacyUnmarshal
	}
	if _, hasMerge := v.(legacyMerger); hasMerge || (hasMarshal && hasUnmarshal) {
		mi.methods.Merge = legacyMerge
	}

	if mi, ok := legacyMessageTypeCache.LoadOrStore(t, mi); ok {
		return mi.(*MessageInfo)
	}
	return mi
}

var legacyMessageDescCache sync.Map 





func LegacyLoadMessageDesc(t reflect.Type) protoreflect.MessageDescriptor {
	return legacyLoadMessageDesc(t, "")
}
func legacyLoadMessageDesc(t reflect.Type, name protoreflect.FullName) protoreflect.MessageDescriptor {
	
	if mi, ok := legacyMessageDescCache.Load(t); ok {
		return mi.(protoreflect.MessageDescriptor)
	}

	
	mv := reflect.Zero(t).Interface()
	if _, ok := mv.(protoreflect.ProtoMessage); ok {
		panic(fmt.Sprintf("%v already implements proto.Message", t))
	}
	mdV1, ok := mv.(messageV1)
	if !ok {
		return aberrantLoadMessageDesc(t, name)
	}

	
	
	
	
	b, idxs := func() ([]byte, []int) {
		defer func() {
			recover()
		}()
		return mdV1.Descriptor()
	}()
	if b == nil {
		return aberrantLoadMessageDesc(t, name)
	}

	
	
	
	if t.Elem().Kind() == reflect.Struct {
		if nfield := t.Elem().NumField(); nfield > 0 {
			hasProtoField := false
			for i := 0; i < nfield; i++ {
				f := t.Elem().Field(i)
				if f.Tag.Get("protobuf") != "" || f.Tag.Get("protobuf_oneof") != "" || strings.HasPrefix(f.Name, "XXX_") {
					hasProtoField = true
					break
				}
			}
			if !hasProtoField {
				return aberrantLoadMessageDesc(t, name)
			}
		}
	}

	md := legacyLoadFileDesc(b).Messages().Get(idxs[0])
	for _, i := range idxs[1:] {
		md = md.Messages().Get(i)
	}
	if name != "" && md.FullName() != name {
		panic(fmt.Sprintf("mismatching message name: got %v, want %v", md.FullName(), name))
	}
	if md, ok := legacyMessageDescCache.LoadOrStore(t, md); ok {
		return md.(protoreflect.MessageDescriptor)
	}
	return md
}

var (
	aberrantMessageDescLock  sync.Mutex
	aberrantMessageDescCache map[reflect.Type]protoreflect.MessageDescriptor
)






func aberrantLoadMessageDesc(t reflect.Type, name protoreflect.FullName) protoreflect.MessageDescriptor {
	aberrantMessageDescLock.Lock()
	defer aberrantMessageDescLock.Unlock()
	if aberrantMessageDescCache == nil {
		aberrantMessageDescCache = make(map[reflect.Type]protoreflect.MessageDescriptor)
	}
	return aberrantLoadMessageDescReentrant(t, name)
}
func aberrantLoadMessageDescReentrant(t reflect.Type, name protoreflect.FullName) protoreflect.MessageDescriptor {
	
	if md, ok := aberrantMessageDescCache[t]; ok {
		return md
	}

	
	
	
	md := &filedesc.Message{L2: new(filedesc.MessageL2)}
	md.L0.FullName = aberrantDeriveMessageName(t, name)
	md.L0.ParentFile = filedesc.SurrogateProto2
	aberrantMessageDescCache[t] = md

	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return md
	}

	
	for i := 0; i < t.Elem().NumField(); i++ {
		f := t.Elem().Field(i)
		if tag := f.Tag.Get("protobuf"); tag != "" {
			switch f.Type.Kind() {
			case reflect.Bool, reflect.Int32, reflect.Int64, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64, reflect.String:
				md.L0.ParentFile = filedesc.SurrogateProto3
			}
			for _, s := range strings.Split(tag, ",") {
				if s == "proto3" {
					md.L0.ParentFile = filedesc.SurrogateProto3
				}
			}
		}
	}

	md.L1.EditionFeatures = md.L0.ParentFile.L1.EditionFeatures
	
	var oneofWrappers []reflect.Type
	methods := make([]reflect.Method, 0, 2)
	if m, ok := t.MethodByName("XXX_OneofFuncs"); ok {
		methods = append(methods, m)
	}
	if m, ok := t.MethodByName("XXX_OneofWrappers"); ok {
		methods = append(methods, m)
	}
	for _, fn := range methods {
		for _, v := range fn.Func.Call([]reflect.Value{reflect.Zero(fn.Type.In(0))}) {
			if vs, ok := v.Interface().([]any); ok {
				for _, v := range vs {
					oneofWrappers = append(oneofWrappers, reflect.TypeOf(v))
				}
			}
		}
	}

	
	if fn, ok := t.MethodByName("ExtensionRangeArray"); ok {
		vs := fn.Func.Call([]reflect.Value{reflect.Zero(fn.Type.In(0))})[0]
		for i := 0; i < vs.Len(); i++ {
			v := vs.Index(i)
			md.L2.ExtensionRanges.List = append(md.L2.ExtensionRanges.List, [2]protoreflect.FieldNumber{
				protoreflect.FieldNumber(v.FieldByName("Start").Int()),
				protoreflect.FieldNumber(v.FieldByName("End").Int() + 1),
			})
			md.L2.ExtensionRangeOptions = append(md.L2.ExtensionRangeOptions, nil)
		}
	}

	
	for i := 0; i < t.Elem().NumField(); i++ {
		f := t.Elem().Field(i)
		if tag := f.Tag.Get("protobuf"); tag != "" {
			tagKey := f.Tag.Get("protobuf_key")
			tagVal := f.Tag.Get("protobuf_val")
			aberrantAppendField(md, f.Type, tag, tagKey, tagVal)
		}
		if tag := f.Tag.Get("protobuf_oneof"); tag != "" {
			n := len(md.L2.Oneofs.List)
			md.L2.Oneofs.List = append(md.L2.Oneofs.List, filedesc.Oneof{})
			od := &md.L2.Oneofs.List[n]
			od.L0.FullName = md.FullName().Append(protoreflect.Name(tag))
			od.L0.ParentFile = md.L0.ParentFile
			od.L1.EditionFeatures = md.L1.EditionFeatures
			od.L0.Parent = md
			od.L0.Index = n

			for _, t := range oneofWrappers {
				if t.Implements(f.Type) {
					f := t.Elem().Field(0)
					if tag := f.Tag.Get("protobuf"); tag != "" {
						aberrantAppendField(md, f.Type, tag, "", "")
						fd := &md.L2.Fields.List[len(md.L2.Fields.List)-1]
						fd.L1.ContainingOneof = od
						fd.L1.EditionFeatures = od.L1.EditionFeatures
						od.L1.Fields.List = append(od.L1.Fields.List, fd)
					}
				}
			}
		}
	}

	return md
}

func aberrantDeriveMessageName(t reflect.Type, name protoreflect.FullName) protoreflect.FullName {
	if name.IsValid() {
		return name
	}
	func() {
		defer func() { recover() }() 
		if m, ok := reflect.Zero(t).Interface().(interface{ XXX_MessageName() string }); ok {
			name = protoreflect.FullName(m.XXX_MessageName())
		}
	}()
	if name.IsValid() {
		return name
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return AberrantDeriveFullName(t)
}

func aberrantAppendField(md *filedesc.Message, goType reflect.Type, tag, tagKey, tagVal string) {
	t := goType
	isOptional := t.Kind() == reflect.Ptr && t.Elem().Kind() != reflect.Struct
	isRepeated := t.Kind() == reflect.Slice && t.Elem().Kind() != reflect.Uint8
	if isOptional || isRepeated {
		t = t.Elem()
	}
	fd := ptag.Unmarshal(tag, t, placeholderEnumValues{}).(*filedesc.Field)

	
	n := len(md.L2.Fields.List)
	md.L2.Fields.List = append(md.L2.Fields.List, *fd)
	fd = &md.L2.Fields.List[n]
	fd.L0.FullName = md.FullName().Append(fd.Name())
	fd.L0.ParentFile = md.L0.ParentFile
	fd.L0.Parent = md
	fd.L0.Index = n

	if fd.L1.EditionFeatures.IsPacked {
		fd.L1.Options = func() protoreflect.ProtoMessage {
			opts := descopts.Field.ProtoReflect().New()
			if fd.L1.EditionFeatures.IsPacked {
				opts.Set(opts.Descriptor().Fields().ByName("packed"), protoreflect.ValueOfBool(fd.L1.EditionFeatures.IsPacked))
			}
			return opts.Interface()
		}
	}

	
	if fd.Enum() == nil && fd.Kind() == protoreflect.EnumKind {
		switch v := reflect.Zero(t).Interface().(type) {
		case protoreflect.Enum:
			fd.L1.Enum = v.Descriptor()
		default:
			fd.L1.Enum = LegacyLoadEnumDesc(t)
		}
	}
	if fd.Message() == nil && (fd.Kind() == protoreflect.MessageKind || fd.Kind() == protoreflect.GroupKind) {
		switch v := reflect.Zero(t).Interface().(type) {
		case protoreflect.ProtoMessage:
			fd.L1.Message = v.ProtoReflect().Descriptor()
		case messageV1:
			fd.L1.Message = LegacyLoadMessageDesc(t)
		default:
			if t.Kind() == reflect.Map {
				n := len(md.L1.Messages.List)
				md.L1.Messages.List = append(md.L1.Messages.List, filedesc.Message{L2: new(filedesc.MessageL2)})
				md2 := &md.L1.Messages.List[n]
				md2.L0.FullName = md.FullName().Append(protoreflect.Name(strs.MapEntryName(string(fd.Name()))))
				md2.L0.ParentFile = md.L0.ParentFile
				md2.L0.Parent = md
				md2.L0.Index = n
				md2.L1.EditionFeatures = md.L1.EditionFeatures

				md2.L1.IsMapEntry = true
				md2.L2.Options = func() protoreflect.ProtoMessage {
					opts := descopts.Message.ProtoReflect().New()
					opts.Set(opts.Descriptor().Fields().ByName("map_entry"), protoreflect.ValueOfBool(true))
					return opts.Interface()
				}

				aberrantAppendField(md2, t.Key(), tagKey, "", "")
				aberrantAppendField(md2, t.Elem(), tagVal, "", "")

				fd.L1.Message = md2
				break
			}
			fd.L1.Message = aberrantLoadMessageDescReentrant(t, "")
		}
	}
}

type placeholderEnumValues struct {
	protoreflect.EnumValueDescriptors
}

func (placeholderEnumValues) ByNumber(n protoreflect.EnumNumber) protoreflect.EnumValueDescriptor {
	return filedesc.PlaceholderEnumValue(protoreflect.FullName(fmt.Sprintf("UNKNOWN_%d", n)))
}


type legacyMarshaler interface {
	Marshal() ([]byte, error)
}


type legacyUnmarshaler interface {
	Unmarshal([]byte) error
}


type legacyMerger interface {
	Merge(protoiface.MessageV1)
}

var aberrantProtoMethods = &protoiface.Methods{
	Marshal:   legacyMarshal,
	Unmarshal: legacyUnmarshal,
	Merge:     legacyMerge,

	
	
	
	
	Flags: protoiface.SupportMarshalDeterministic,
}

func legacyMarshal(in protoiface.MarshalInput) (protoiface.MarshalOutput, error) {
	v := in.Message.(unwrapper).protoUnwrap()
	marshaler, ok := v.(legacyMarshaler)
	if !ok {
		return protoiface.MarshalOutput{}, errors.New("%T does not implement Marshal", v)
	}
	out, err := marshaler.Marshal()
	if in.Buf != nil {
		out = append(in.Buf, out...)
	}
	return protoiface.MarshalOutput{
		Buf: out,
	}, err
}

func legacyUnmarshal(in protoiface.UnmarshalInput) (protoiface.UnmarshalOutput, error) {
	v := in.Message.(unwrapper).protoUnwrap()
	unmarshaler, ok := v.(legacyUnmarshaler)
	if !ok {
		return protoiface.UnmarshalOutput{}, errors.New("%T does not implement Unmarshal", v)
	}
	return protoiface.UnmarshalOutput{}, unmarshaler.Unmarshal(in.Buf)
}

func legacyMerge(in protoiface.MergeInput) protoiface.MergeOutput {
	
	dstv := in.Destination.(unwrapper).protoUnwrap()
	merger, ok := dstv.(legacyMerger)
	if ok {
		merger.Merge(Export{}.ProtoMessageV1Of(in.Source))
		return protoiface.MergeOutput{Flags: protoiface.MergeComplete}
	}

	
	
	srcv := in.Source.(unwrapper).protoUnwrap()
	marshaler, ok := srcv.(legacyMarshaler)
	if !ok {
		return protoiface.MergeOutput{}
	}
	dstv = in.Destination.(unwrapper).protoUnwrap()
	unmarshaler, ok := dstv.(legacyUnmarshaler)
	if !ok {
		return protoiface.MergeOutput{}
	}
	if !in.Source.IsValid() {
		
		
		
		
		return protoiface.MergeOutput{Flags: protoiface.MergeComplete}
	}
	b, err := marshaler.Marshal()
	if err != nil {
		return protoiface.MergeOutput{}
	}
	err = unmarshaler.Unmarshal(b)
	if err != nil {
		return protoiface.MergeOutput{}
	}
	return protoiface.MergeOutput{Flags: protoiface.MergeComplete}
}


type aberrantMessageType struct {
	t reflect.Type
}

func (mt aberrantMessageType) New() protoreflect.Message {
	if mt.t.Kind() == reflect.Ptr {
		return aberrantMessage{reflect.New(mt.t.Elem())}
	}
	return aberrantMessage{reflect.Zero(mt.t)}
}
func (mt aberrantMessageType) Zero() protoreflect.Message {
	return aberrantMessage{reflect.Zero(mt.t)}
}
func (mt aberrantMessageType) GoType() reflect.Type {
	return mt.t
}
func (mt aberrantMessageType) Descriptor() protoreflect.MessageDescriptor {
	return LegacyLoadMessageDesc(mt.t)
}






type aberrantMessage struct {
	v reflect.Value
}


func (m aberrantMessage) Reset() {
	if mr, ok := m.v.Interface().(interface{ Reset() }); ok {
		mr.Reset()
		return
	}
	if m.v.Kind() == reflect.Ptr && !m.v.IsNil() {
		m.v.Elem().Set(reflect.Zero(m.v.Type().Elem()))
	}
}

func (m aberrantMessage) ProtoReflect() protoreflect.Message {
	return m
}

func (m aberrantMessage) Descriptor() protoreflect.MessageDescriptor {
	return LegacyLoadMessageDesc(m.v.Type())
}
func (m aberrantMessage) Type() protoreflect.MessageType {
	return aberrantMessageType{m.v.Type()}
}
func (m aberrantMessage) New() protoreflect.Message {
	if m.v.Type().Kind() == reflect.Ptr {
		return aberrantMessage{reflect.New(m.v.Type().Elem())}
	}
	return aberrantMessage{reflect.Zero(m.v.Type())}
}
func (m aberrantMessage) Interface() protoreflect.ProtoMessage {
	return m
}
func (m aberrantMessage) Range(f func(protoreflect.FieldDescriptor, protoreflect.Value) bool) {
	return
}
func (m aberrantMessage) Has(protoreflect.FieldDescriptor) bool {
	return false
}
func (m aberrantMessage) Clear(protoreflect.FieldDescriptor) {
	panic("invalid Message.Clear on " + string(m.Descriptor().FullName()))
}
func (m aberrantMessage) Get(fd protoreflect.FieldDescriptor) protoreflect.Value {
	if fd.Default().IsValid() {
		return fd.Default()
	}
	panic("invalid Message.Get on " + string(m.Descriptor().FullName()))
}
func (m aberrantMessage) Set(protoreflect.FieldDescriptor, protoreflect.Value) {
	panic("invalid Message.Set on " + string(m.Descriptor().FullName()))
}
func (m aberrantMessage) Mutable(protoreflect.FieldDescriptor) protoreflect.Value {
	panic("invalid Message.Mutable on " + string(m.Descriptor().FullName()))
}
func (m aberrantMessage) NewField(protoreflect.FieldDescriptor) protoreflect.Value {
	panic("invalid Message.NewField on " + string(m.Descriptor().FullName()))
}
func (m aberrantMessage) WhichOneof(protoreflect.OneofDescriptor) protoreflect.FieldDescriptor {
	panic("invalid Message.WhichOneof descriptor on " + string(m.Descriptor().FullName()))
}
func (m aberrantMessage) GetUnknown() protoreflect.RawFields {
	return nil
}
func (m aberrantMessage) SetUnknown(protoreflect.RawFields) {
	
}
func (m aberrantMessage) IsValid() bool {
	if m.v.Kind() == reflect.Ptr {
		return !m.v.IsNil()
	}
	return false
}
func (m aberrantMessage) ProtoMethods() *protoiface.Methods {
	return aberrantProtoMethods
}
func (m aberrantMessage) protoUnwrap() any {
	return m.v.Interface()
}
