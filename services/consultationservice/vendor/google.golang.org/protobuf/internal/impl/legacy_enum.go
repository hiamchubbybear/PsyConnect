



package impl

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"google.golang.org/protobuf/internal/filedesc"
	"google.golang.org/protobuf/internal/strs"
	"google.golang.org/protobuf/reflect/protoreflect"
)




func legacyEnumName(ed protoreflect.EnumDescriptor) string {
	var protoPkg string
	enumName := string(ed.FullName())
	if fd := ed.ParentFile(); fd != nil {
		protoPkg = string(fd.Package())
		enumName = strings.TrimPrefix(enumName, protoPkg+".")
	}
	if protoPkg == "" {
		return strs.GoCamelCase(enumName)
	}
	return protoPkg + "." + strs.GoCamelCase(enumName)
}



func legacyWrapEnum(v reflect.Value) protoreflect.Enum {
	et := legacyLoadEnumType(v.Type())
	return et.New(protoreflect.EnumNumber(v.Int()))
}

var legacyEnumTypeCache sync.Map 



func legacyLoadEnumType(t reflect.Type) protoreflect.EnumType {
	
	if et, ok := legacyEnumTypeCache.Load(t); ok {
		return et.(protoreflect.EnumType)
	}

	
	var et protoreflect.EnumType
	ed := LegacyLoadEnumDesc(t)
	et = &legacyEnumType{
		desc:   ed,
		goType: t,
	}
	if et, ok := legacyEnumTypeCache.LoadOrStore(t, et); ok {
		return et.(protoreflect.EnumType)
	}
	return et
}

type legacyEnumType struct {
	desc   protoreflect.EnumDescriptor
	goType reflect.Type
	m      sync.Map 
}

func (t *legacyEnumType) New(n protoreflect.EnumNumber) protoreflect.Enum {
	if e, ok := t.m.Load(n); ok {
		return e.(protoreflect.Enum)
	}
	e := &legacyEnumWrapper{num: n, pbTyp: t, goTyp: t.goType}
	t.m.Store(n, e)
	return e
}
func (t *legacyEnumType) Descriptor() protoreflect.EnumDescriptor {
	return t.desc
}

type legacyEnumWrapper struct {
	num   protoreflect.EnumNumber
	pbTyp protoreflect.EnumType
	goTyp reflect.Type
}

func (e *legacyEnumWrapper) Descriptor() protoreflect.EnumDescriptor {
	return e.pbTyp.Descriptor()
}
func (e *legacyEnumWrapper) Type() protoreflect.EnumType {
	return e.pbTyp
}
func (e *legacyEnumWrapper) Number() protoreflect.EnumNumber {
	return e.num
}
func (e *legacyEnumWrapper) ProtoReflect() protoreflect.Enum {
	return e
}
func (e *legacyEnumWrapper) protoUnwrap() any {
	v := reflect.New(e.goTyp).Elem()
	v.SetInt(int64(e.num))
	return v.Interface()
}

var (
	_ protoreflect.Enum = (*legacyEnumWrapper)(nil)
	_ unwrapper         = (*legacyEnumWrapper)(nil)
)

var legacyEnumDescCache sync.Map 





func LegacyLoadEnumDesc(t reflect.Type) protoreflect.EnumDescriptor {
	
	if ed, ok := legacyEnumDescCache.Load(t); ok {
		return ed.(protoreflect.EnumDescriptor)
	}

	
	ev := reflect.Zero(t).Interface()
	if _, ok := ev.(protoreflect.Enum); ok {
		panic(fmt.Sprintf("%v already implements proto.Enum", t))
	}
	edV1, ok := ev.(enumV1)
	if !ok {
		return aberrantLoadEnumDesc(t)
	}
	b, idxs := edV1.EnumDescriptor()

	var ed protoreflect.EnumDescriptor
	if len(idxs) == 1 {
		ed = legacyLoadFileDesc(b).Enums().Get(idxs[0])
	} else {
		md := legacyLoadFileDesc(b).Messages().Get(idxs[0])
		for _, i := range idxs[1 : len(idxs)-1] {
			md = md.Messages().Get(i)
		}
		ed = md.Enums().Get(idxs[len(idxs)-1])
	}
	if ed, ok := legacyEnumDescCache.LoadOrStore(t, ed); ok {
		return ed.(protoreflect.EnumDescriptor)
	}
	return ed
}

var aberrantEnumDescCache sync.Map 









func aberrantLoadEnumDesc(t reflect.Type) protoreflect.EnumDescriptor {
	
	if ed, ok := aberrantEnumDescCache.Load(t); ok {
		return ed.(protoreflect.EnumDescriptor)
	}

	
	ed := &filedesc.Enum{L2: new(filedesc.EnumL2)}
	ed.L0.FullName = AberrantDeriveFullName(t) 
	ed.L0.ParentFile = filedesc.SurrogateProto3
	ed.L1.EditionFeatures = ed.L0.ParentFile.L1.EditionFeatures
	ed.L2.Values.List = append(ed.L2.Values.List, filedesc.EnumValue{})

	

	vd := &ed.L2.Values.List[0]
	vd.L0.FullName = ed.L0.FullName + "_UNKNOWN" 
	vd.L0.ParentFile = ed.L0.ParentFile
	vd.L0.Parent = ed

	
	
	

	if ed, ok := aberrantEnumDescCache.LoadOrStore(t, ed); ok {
		return ed.(protoreflect.EnumDescriptor)
	}
	return ed
}






func AberrantDeriveFullName(t reflect.Type) protoreflect.FullName {
	sanitize := func(r rune) rune {
		switch {
		case r == '/':
			return '.'
		case 'a' <= r && r <= 'z', 'A' <= r && r <= 'Z', '0' <= r && r <= '9':
			return r
		default:
			return '_'
		}
	}
	prefix := strings.Map(sanitize, t.PkgPath())
	suffix := strings.Map(sanitize, t.Name())
	if suffix == "" {
		suffix = fmt.Sprintf("UnknownX%X", reflect.ValueOf(t).Pointer())
	}

	ss := append(strings.Split(prefix, "."), suffix)
	for i, s := range ss {
		if s == "" || ('0' <= s[0] && s[0] <= '9') {
			ss[i] = "x" + s
		}
	}
	return protoreflect.FullName(strings.Join(ss, "."))
}
