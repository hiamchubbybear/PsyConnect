



package filedesc

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	"google.golang.org/protobuf/internal/descfmt"
	"google.golang.org/protobuf/internal/descopts"
	"google.golang.org/protobuf/internal/encoding/defval"
	"google.golang.org/protobuf/internal/encoding/messageset"
	"google.golang.org/protobuf/internal/genid"
	"google.golang.org/protobuf/internal/pragma"
	"google.golang.org/protobuf/internal/strs"
	"google.golang.org/protobuf/reflect/protoreflect"
)


type Edition int32



const (
	EditionUnknown     Edition = 0
	EditionProto2      Edition = 998
	EditionProto3      Edition = 999
	Edition2023        Edition = 1000
	Edition2024        Edition = 1001
	EditionUnsupported Edition = 100000
)















type (
	File struct {
		fileRaw
		L1 FileL1

		once uint32     
		mu   sync.Mutex 
		L2   *FileL2
	}
	FileL1 struct {
		Syntax  protoreflect.Syntax
		Edition Edition 
		Path    string
		Package protoreflect.FullName

		Enums      Enums
		Messages   Messages
		Extensions Extensions
		Services   Services

		EditionFeatures EditionFeatures
	}
	FileL2 struct {
		Options   func() protoreflect.ProtoMessage
		Imports   FileImports
		Locations SourceLocations
	}

	
	
	
	EditionFeatures struct {
		
		
		StripEnumPrefix int

		
		
		IsFieldPresence bool

		
		
		IsLegacyRequired bool

		
		
		IsOpenEnum bool

		
		
		IsPacked bool

		
		
		IsUTF8Validated bool

		
		
		IsDelimitedEncoded bool

		
		
		IsJSONCompliant bool

		
		
		GenerateLegacyUnmarshalJSON bool
		
		
		APILevel int
	}
)

func (fd *File) ParentFile() protoreflect.FileDescriptor { return fd }
func (fd *File) Parent() protoreflect.Descriptor         { return nil }
func (fd *File) Index() int                              { return 0 }
func (fd *File) Syntax() protoreflect.Syntax             { return fd.L1.Syntax }


func (fd *File) Edition() int32                  { return int32(fd.L1.Edition) }
func (fd *File) Name() protoreflect.Name         { return fd.L1.Package.Name() }
func (fd *File) FullName() protoreflect.FullName { return fd.L1.Package }
func (fd *File) IsPlaceholder() bool             { return false }
func (fd *File) Options() protoreflect.ProtoMessage {
	if f := fd.lazyInit().Options; f != nil {
		return f()
	}
	return descopts.File
}
func (fd *File) Path() string                                  { return fd.L1.Path }
func (fd *File) Package() protoreflect.FullName                { return fd.L1.Package }
func (fd *File) Imports() protoreflect.FileImports             { return &fd.lazyInit().Imports }
func (fd *File) Enums() protoreflect.EnumDescriptors           { return &fd.L1.Enums }
func (fd *File) Messages() protoreflect.MessageDescriptors     { return &fd.L1.Messages }
func (fd *File) Extensions() protoreflect.ExtensionDescriptors { return &fd.L1.Extensions }
func (fd *File) Services() protoreflect.ServiceDescriptors     { return &fd.L1.Services }
func (fd *File) SourceLocations() protoreflect.SourceLocations { return &fd.lazyInit().Locations }
func (fd *File) Format(s fmt.State, r rune)                    { descfmt.FormatDesc(s, r, fd) }
func (fd *File) ProtoType(protoreflect.FileDescriptor)         {}
func (fd *File) ProtoInternal(pragma.DoNotImplement)           {}

func (fd *File) lazyInit() *FileL2 {
	if atomic.LoadUint32(&fd.once) == 0 {
		fd.lazyInitOnce()
	}
	return fd.L2
}

func (fd *File) lazyInitOnce() {
	fd.mu.Lock()
	if fd.L2 == nil {
		fd.lazyRawInit() 
	}
	atomic.StoreUint32(&fd.once, 1)
	fd.mu.Unlock()
}






func (fd *File) GoPackagePath() string {
	return fd.builder.GoPackagePath
}

type (
	Enum struct {
		Base
		L1 EnumL1
		L2 *EnumL2 
	}
	EnumL1 struct {
		eagerValues bool 

		EditionFeatures EditionFeatures
	}
	EnumL2 struct {
		Options        func() protoreflect.ProtoMessage
		Values         EnumValues
		ReservedNames  Names
		ReservedRanges EnumRanges
	}

	EnumValue struct {
		Base
		L1 EnumValueL1
	}
	EnumValueL1 struct {
		Options func() protoreflect.ProtoMessage
		Number  protoreflect.EnumNumber
	}
)

func (ed *Enum) Options() protoreflect.ProtoMessage {
	if f := ed.lazyInit().Options; f != nil {
		return f()
	}
	return descopts.Enum
}
func (ed *Enum) Values() protoreflect.EnumValueDescriptors {
	if ed.L1.eagerValues {
		return &ed.L2.Values
	}
	return &ed.lazyInit().Values
}
func (ed *Enum) ReservedNames() protoreflect.Names       { return &ed.lazyInit().ReservedNames }
func (ed *Enum) ReservedRanges() protoreflect.EnumRanges { return &ed.lazyInit().ReservedRanges }
func (ed *Enum) Format(s fmt.State, r rune)              { descfmt.FormatDesc(s, r, ed) }
func (ed *Enum) ProtoType(protoreflect.EnumDescriptor)   {}
func (ed *Enum) lazyInit() *EnumL2 {
	ed.L0.ParentFile.lazyInit() 
	return ed.L2
}
func (ed *Enum) IsClosed() bool {
	return !ed.L1.EditionFeatures.IsOpenEnum
}

func (ed *EnumValue) Options() protoreflect.ProtoMessage {
	if f := ed.L1.Options; f != nil {
		return f()
	}
	return descopts.EnumValue
}
func (ed *EnumValue) Number() protoreflect.EnumNumber            { return ed.L1.Number }
func (ed *EnumValue) Format(s fmt.State, r rune)                 { descfmt.FormatDesc(s, r, ed) }
func (ed *EnumValue) ProtoType(protoreflect.EnumValueDescriptor) {}

type (
	Message struct {
		Base
		L1 MessageL1
		L2 *MessageL2 
	}
	MessageL1 struct {
		Enums        Enums
		Messages     Messages
		Extensions   Extensions
		IsMapEntry   bool 
		IsMessageSet bool 

		EditionFeatures EditionFeatures
	}
	MessageL2 struct {
		Options               func() protoreflect.ProtoMessage
		Fields                Fields
		Oneofs                Oneofs
		ReservedNames         Names
		ReservedRanges        FieldRanges
		RequiredNumbers       FieldNumbers 
		ExtensionRanges       FieldRanges
		ExtensionRangeOptions []func() protoreflect.ProtoMessage 
	}

	Field struct {
		Base
		L1 FieldL1
	}
	FieldL1 struct {
		Options          func() protoreflect.ProtoMessage
		Number           protoreflect.FieldNumber
		Cardinality      protoreflect.Cardinality 
		Kind             protoreflect.Kind
		StringName       stringName
		IsProto3Optional bool 
		IsLazy           bool 
		Default          defaultValue
		ContainingOneof  protoreflect.OneofDescriptor 
		Enum             protoreflect.EnumDescriptor
		Message          protoreflect.MessageDescriptor

		EditionFeatures EditionFeatures
	}

	Oneof struct {
		Base
		L1 OneofL1
	}
	OneofL1 struct {
		Options func() protoreflect.ProtoMessage
		Fields  OneofFields 

		EditionFeatures EditionFeatures
	}
)

func (md *Message) Options() protoreflect.ProtoMessage {
	if f := md.lazyInit().Options; f != nil {
		return f()
	}
	return descopts.Message
}
func (md *Message) IsMapEntry() bool                           { return md.L1.IsMapEntry }
func (md *Message) Fields() protoreflect.FieldDescriptors      { return &md.lazyInit().Fields }
func (md *Message) Oneofs() protoreflect.OneofDescriptors      { return &md.lazyInit().Oneofs }
func (md *Message) ReservedNames() protoreflect.Names          { return &md.lazyInit().ReservedNames }
func (md *Message) ReservedRanges() protoreflect.FieldRanges   { return &md.lazyInit().ReservedRanges }
func (md *Message) RequiredNumbers() protoreflect.FieldNumbers { return &md.lazyInit().RequiredNumbers }
func (md *Message) ExtensionRanges() protoreflect.FieldRanges  { return &md.lazyInit().ExtensionRanges }
func (md *Message) ExtensionRangeOptions(i int) protoreflect.ProtoMessage {
	if f := md.lazyInit().ExtensionRangeOptions[i]; f != nil {
		return f()
	}
	return descopts.ExtensionRange
}
func (md *Message) Enums() protoreflect.EnumDescriptors           { return &md.L1.Enums }
func (md *Message) Messages() protoreflect.MessageDescriptors     { return &md.L1.Messages }
func (md *Message) Extensions() protoreflect.ExtensionDescriptors { return &md.L1.Extensions }
func (md *Message) ProtoType(protoreflect.MessageDescriptor)      {}
func (md *Message) Format(s fmt.State, r rune)                    { descfmt.FormatDesc(s, r, md) }
func (md *Message) lazyInit() *MessageL2 {
	md.L0.ParentFile.lazyInit() 
	return md.L2
}






func (md *Message) IsMessageSet() bool {
	return md.L1.IsMessageSet
}

func (fd *Field) Options() protoreflect.ProtoMessage {
	if f := fd.L1.Options; f != nil {
		return f()
	}
	return descopts.Field
}
func (fd *Field) Number() protoreflect.FieldNumber      { return fd.L1.Number }
func (fd *Field) Cardinality() protoreflect.Cardinality { return fd.L1.Cardinality }
func (fd *Field) Kind() protoreflect.Kind {
	return fd.L1.Kind
}
func (fd *Field) HasJSONName() bool { return fd.L1.StringName.hasJSON }
func (fd *Field) JSONName() string  { return fd.L1.StringName.getJSON(fd) }
func (fd *Field) TextName() string  { return fd.L1.StringName.getText(fd) }
func (fd *Field) HasPresence() bool {
	if fd.L1.Cardinality == protoreflect.Repeated {
		return false
	}
	return fd.IsExtension() || fd.L1.EditionFeatures.IsFieldPresence || fd.L1.Message != nil || fd.L1.ContainingOneof != nil
}
func (fd *Field) HasOptionalKeyword() bool {
	return (fd.L0.ParentFile.L1.Syntax == protoreflect.Proto2 && fd.L1.Cardinality == protoreflect.Optional && fd.L1.ContainingOneof == nil) || fd.L1.IsProto3Optional
}
func (fd *Field) IsPacked() bool {
	if fd.L1.Cardinality != protoreflect.Repeated {
		return false
	}
	switch fd.L1.Kind {
	case protoreflect.StringKind, protoreflect.BytesKind, protoreflect.MessageKind, protoreflect.GroupKind:
		return false
	}
	return fd.L1.EditionFeatures.IsPacked
}
func (fd *Field) IsExtension() bool { return false }
func (fd *Field) IsWeak() bool      { return false }
func (fd *Field) IsLazy() bool      { return fd.L1.IsLazy }
func (fd *Field) IsList() bool      { return fd.Cardinality() == protoreflect.Repeated && !fd.IsMap() }
func (fd *Field) IsMap() bool       { return fd.Message() != nil && fd.Message().IsMapEntry() }
func (fd *Field) MapKey() protoreflect.FieldDescriptor {
	if !fd.IsMap() {
		return nil
	}
	return fd.Message().Fields().ByNumber(genid.MapEntry_Key_field_number)
}
func (fd *Field) MapValue() protoreflect.FieldDescriptor {
	if !fd.IsMap() {
		return nil
	}
	return fd.Message().Fields().ByNumber(genid.MapEntry_Value_field_number)
}
func (fd *Field) HasDefault() bool                                   { return fd.L1.Default.has }
func (fd *Field) Default() protoreflect.Value                        { return fd.L1.Default.get(fd) }
func (fd *Field) DefaultEnumValue() protoreflect.EnumValueDescriptor { return fd.L1.Default.enum }
func (fd *Field) ContainingOneof() protoreflect.OneofDescriptor      { return fd.L1.ContainingOneof }
func (fd *Field) ContainingMessage() protoreflect.MessageDescriptor {
	return fd.L0.Parent.(protoreflect.MessageDescriptor)
}
func (fd *Field) Enum() protoreflect.EnumDescriptor {
	return fd.L1.Enum
}
func (fd *Field) Message() protoreflect.MessageDescriptor {
	return fd.L1.Message
}
func (fd *Field) IsMapEntry() bool {
	parent, ok := fd.L0.Parent.(protoreflect.MessageDescriptor)
	return ok && parent.IsMapEntry()
}
func (fd *Field) Format(s fmt.State, r rune)             { descfmt.FormatDesc(s, r, fd) }
func (fd *Field) ProtoType(protoreflect.FieldDescriptor) {}








func (fd *Field) EnforceUTF8() bool {
	return fd.L1.EditionFeatures.IsUTF8Validated
}

func (od *Oneof) IsSynthetic() bool {
	return od.L0.ParentFile.L1.Syntax == protoreflect.Proto3 && len(od.L1.Fields.List) == 1 && od.L1.Fields.List[0].HasOptionalKeyword()
}
func (od *Oneof) Options() protoreflect.ProtoMessage {
	if f := od.L1.Options; f != nil {
		return f()
	}
	return descopts.Oneof
}
func (od *Oneof) Fields() protoreflect.FieldDescriptors  { return &od.L1.Fields }
func (od *Oneof) Format(s fmt.State, r rune)             { descfmt.FormatDesc(s, r, od) }
func (od *Oneof) ProtoType(protoreflect.OneofDescriptor) {}

type (
	Extension struct {
		Base
		L1 ExtensionL1
		L2 *ExtensionL2 
	}
	ExtensionL1 struct {
		Number          protoreflect.FieldNumber
		Extendee        protoreflect.MessageDescriptor
		Cardinality     protoreflect.Cardinality
		Kind            protoreflect.Kind
		IsLazy          bool
		EditionFeatures EditionFeatures
	}
	ExtensionL2 struct {
		Options          func() protoreflect.ProtoMessage
		StringName       stringName
		IsProto3Optional bool 
		Default          defaultValue
		Enum             protoreflect.EnumDescriptor
		Message          protoreflect.MessageDescriptor
	}
)

func (xd *Extension) Options() protoreflect.ProtoMessage {
	if f := xd.lazyInit().Options; f != nil {
		return f()
	}
	return descopts.Field
}
func (xd *Extension) Number() protoreflect.FieldNumber      { return xd.L1.Number }
func (xd *Extension) Cardinality() protoreflect.Cardinality { return xd.L1.Cardinality }
func (xd *Extension) Kind() protoreflect.Kind               { return xd.L1.Kind }
func (xd *Extension) HasJSONName() bool                     { return xd.lazyInit().StringName.hasJSON }
func (xd *Extension) JSONName() string                      { return xd.lazyInit().StringName.getJSON(xd) }
func (xd *Extension) TextName() string                      { return xd.lazyInit().StringName.getText(xd) }
func (xd *Extension) HasPresence() bool                     { return xd.L1.Cardinality != protoreflect.Repeated }
func (xd *Extension) HasOptionalKeyword() bool {
	return (xd.L0.ParentFile.L1.Syntax == protoreflect.Proto2 && xd.L1.Cardinality == protoreflect.Optional) || xd.lazyInit().IsProto3Optional
}
func (xd *Extension) IsPacked() bool {
	if xd.L1.Cardinality != protoreflect.Repeated {
		return false
	}
	switch xd.L1.Kind {
	case protoreflect.StringKind, protoreflect.BytesKind, protoreflect.MessageKind, protoreflect.GroupKind:
		return false
	}
	return xd.L1.EditionFeatures.IsPacked
}
func (xd *Extension) IsExtension() bool                      { return true }
func (xd *Extension) IsWeak() bool                           { return false }
func (xd *Extension) IsLazy() bool                           { return xd.L1.IsLazy }
func (xd *Extension) IsList() bool                           { return xd.Cardinality() == protoreflect.Repeated }
func (xd *Extension) IsMap() bool                            { return false }
func (xd *Extension) MapKey() protoreflect.FieldDescriptor   { return nil }
func (xd *Extension) MapValue() protoreflect.FieldDescriptor { return nil }
func (xd *Extension) HasDefault() bool                       { return xd.lazyInit().Default.has }
func (xd *Extension) Default() protoreflect.Value            { return xd.lazyInit().Default.get(xd) }
func (xd *Extension) DefaultEnumValue() protoreflect.EnumValueDescriptor {
	return xd.lazyInit().Default.enum
}
func (xd *Extension) ContainingOneof() protoreflect.OneofDescriptor     { return nil }
func (xd *Extension) ContainingMessage() protoreflect.MessageDescriptor { return xd.L1.Extendee }
func (xd *Extension) Enum() protoreflect.EnumDescriptor                 { return xd.lazyInit().Enum }
func (xd *Extension) Message() protoreflect.MessageDescriptor           { return xd.lazyInit().Message }
func (xd *Extension) Format(s fmt.State, r rune)                        { descfmt.FormatDesc(s, r, xd) }
func (xd *Extension) ProtoType(protoreflect.FieldDescriptor)            {}
func (xd *Extension) ProtoInternal(pragma.DoNotImplement)               {}
func (xd *Extension) lazyInit() *ExtensionL2 {
	xd.L0.ParentFile.lazyInit() 
	return xd.L2
}

type (
	Service struct {
		Base
		L1 ServiceL1
		L2 *ServiceL2 
	}
	ServiceL1 struct{}
	ServiceL2 struct {
		Options func() protoreflect.ProtoMessage
		Methods Methods
	}

	Method struct {
		Base
		L1 MethodL1
	}
	MethodL1 struct {
		Options           func() protoreflect.ProtoMessage
		Input             protoreflect.MessageDescriptor
		Output            protoreflect.MessageDescriptor
		IsStreamingClient bool
		IsStreamingServer bool
	}
)

func (sd *Service) Options() protoreflect.ProtoMessage {
	if f := sd.lazyInit().Options; f != nil {
		return f()
	}
	return descopts.Service
}
func (sd *Service) Methods() protoreflect.MethodDescriptors  { return &sd.lazyInit().Methods }
func (sd *Service) Format(s fmt.State, r rune)               { descfmt.FormatDesc(s, r, sd) }
func (sd *Service) ProtoType(protoreflect.ServiceDescriptor) {}
func (sd *Service) ProtoInternal(pragma.DoNotImplement)      {}
func (sd *Service) lazyInit() *ServiceL2 {
	sd.L0.ParentFile.lazyInit() 
	return sd.L2
}

func (md *Method) Options() protoreflect.ProtoMessage {
	if f := md.L1.Options; f != nil {
		return f()
	}
	return descopts.Method
}
func (md *Method) Input() protoreflect.MessageDescriptor   { return md.L1.Input }
func (md *Method) Output() protoreflect.MessageDescriptor  { return md.L1.Output }
func (md *Method) IsStreamingClient() bool                 { return md.L1.IsStreamingClient }
func (md *Method) IsStreamingServer() bool                 { return md.L1.IsStreamingServer }
func (md *Method) Format(s fmt.State, r rune)              { descfmt.FormatDesc(s, r, md) }
func (md *Method) ProtoType(protoreflect.MethodDescriptor) {}
func (md *Method) ProtoInternal(pragma.DoNotImplement)     {}



var (
	SurrogateProto2      = &File{L1: FileL1{Syntax: protoreflect.Proto2}, L2: &FileL2{}}
	SurrogateProto3      = &File{L1: FileL1{Syntax: protoreflect.Proto3}, L2: &FileL2{}}
	SurrogateEdition2023 = &File{L1: FileL1{Syntax: protoreflect.Editions, Edition: Edition2023}, L2: &FileL2{}}
)

type (
	Base struct {
		L0 BaseL0
	}
	BaseL0 struct {
		FullName   protoreflect.FullName 
		ParentFile *File                 
		Parent     protoreflect.Descriptor
		Index      int
	}
)

func (d *Base) Name() protoreflect.Name         { return d.L0.FullName.Name() }
func (d *Base) FullName() protoreflect.FullName { return d.L0.FullName }
func (d *Base) ParentFile() protoreflect.FileDescriptor {
	if d.L0.ParentFile == SurrogateProto2 || d.L0.ParentFile == SurrogateProto3 {
		return nil 
	}
	return d.L0.ParentFile
}
func (d *Base) Parent() protoreflect.Descriptor     { return d.L0.Parent }
func (d *Base) Index() int                          { return d.L0.Index }
func (d *Base) Syntax() protoreflect.Syntax         { return d.L0.ParentFile.Syntax() }
func (d *Base) IsPlaceholder() bool                 { return false }
func (d *Base) ProtoInternal(pragma.DoNotImplement) {}

type stringName struct {
	hasJSON  bool
	once     sync.Once
	nameJSON string
	nameText string
}


func (s *stringName) InitJSON(name string) {
	s.hasJSON = true
	s.nameJSON = name
}




func isGroupLike(fd protoreflect.FieldDescriptor) bool {
	
	if fd.Kind() != protoreflect.GroupKind {
		return false
	}

	
	if strings.ToLower(string(fd.Message().Name())) != string(fd.Name()) {
		return false
	}

	
	if fd.Message().ParentFile() != fd.ParentFile() {
		return false
	}

	
	
	
	if fd.IsExtension() {
		return fd.Parent() == fd.Message().Parent()
	}
	return fd.ContainingMessage() == fd.Message().Parent()
}

func (s *stringName) lazyInit(fd protoreflect.FieldDescriptor) *stringName {
	s.once.Do(func() {
		if fd.IsExtension() {
			
			var name string
			if messageset.IsMessageSetExtension(fd) {
				name = string("[" + fd.FullName().Parent() + "]")
			} else {
				name = string("[" + fd.FullName() + "]")
			}
			s.nameJSON = name
			s.nameText = name
		} else {
			
			if !s.hasJSON {
				s.nameJSON = strs.JSONCamelCase(string(fd.Name()))
			}

			
			s.nameText = string(fd.Name())
			if isGroupLike(fd) {
				s.nameText = string(fd.Message().Name())
			}
		}
	})
	return s
}

func (s *stringName) getJSON(fd protoreflect.FieldDescriptor) string { return s.lazyInit(fd).nameJSON }
func (s *stringName) getText(fd protoreflect.FieldDescriptor) string { return s.lazyInit(fd).nameText }

func DefaultValue(v protoreflect.Value, ev protoreflect.EnumValueDescriptor) defaultValue {
	dv := defaultValue{has: v.IsValid(), val: v, enum: ev}
	if b, ok := v.Interface().([]byte); ok {
		
		
		dv.bytes = append([]byte(nil), b...)
	}
	return dv
}

func unmarshalDefault(b []byte, k protoreflect.Kind, pf *File, ed protoreflect.EnumDescriptor) defaultValue {
	var evs protoreflect.EnumValueDescriptors
	if k == protoreflect.EnumKind {
		
		
		if e, ok := ed.(*Enum); ok && e.L0.ParentFile == pf {
			evs = &e.L2.Values
		} else {
			evs = ed.Values()
		}

		
		
		if ed.IsPlaceholder() && protoreflect.Name(b).IsValid() {
			v := protoreflect.ValueOfEnum(0)
			ev := PlaceholderEnumValue(ed.FullName().Parent().Append(protoreflect.Name(b)))
			return DefaultValue(v, ev)
		}
	}

	v, ev, err := defval.Unmarshal(string(b), k, evs, defval.Descriptor)
	if err != nil {
		panic(err)
	}
	return DefaultValue(v, ev)
}

type defaultValue struct {
	has   bool
	val   protoreflect.Value
	enum  protoreflect.EnumValueDescriptor
	bytes []byte
}

func (dv *defaultValue) get(fd protoreflect.FieldDescriptor) protoreflect.Value {
	
	if !dv.has {
		if fd.Cardinality() == protoreflect.Repeated {
			return protoreflect.Value{}
		}
		switch fd.Kind() {
		case protoreflect.BoolKind:
			return protoreflect.ValueOfBool(false)
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
			return protoreflect.ValueOfInt32(0)
		case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
			return protoreflect.ValueOfInt64(0)
		case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
			return protoreflect.ValueOfUint32(0)
		case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
			return protoreflect.ValueOfUint64(0)
		case protoreflect.FloatKind:
			return protoreflect.ValueOfFloat32(0)
		case protoreflect.DoubleKind:
			return protoreflect.ValueOfFloat64(0)
		case protoreflect.StringKind:
			return protoreflect.ValueOfString("")
		case protoreflect.BytesKind:
			return protoreflect.ValueOfBytes(nil)
		case protoreflect.EnumKind:
			if evs := fd.Enum().Values(); evs.Len() > 0 {
				return protoreflect.ValueOfEnum(evs.Get(0).Number())
			}
			return protoreflect.ValueOfEnum(0)
		}
	}

	if len(dv.bytes) > 0 && !bytes.Equal(dv.bytes, dv.val.Bytes()) {
		
		
		
		panic(fmt.Sprintf("detected mutation on the default bytes for %v", fd.FullName()))
	}
	return dv.val
}
