














package protoregistry

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"google.golang.org/protobuf/internal/encoding/messageset"
	"google.golang.org/protobuf/internal/errors"
	"google.golang.org/protobuf/internal/flags"
	"google.golang.org/protobuf/reflect/protoreflect"
)













var conflictPolicy = "panic" 




var ignoreConflict = func(d protoreflect.Descriptor, err error) bool {
	const env = "GOLANG_PROTOBUF_REGISTRATION_CONFLICT"
	const faq = "https://protobuf.dev/reference/go/faq#namespace-conflict"
	policy := conflictPolicy
	if v := os.Getenv(env); v != "" {
		policy = v
	}
	switch policy {
	case "panic":
		panic(fmt.Sprintf("%v\nSee %v\n", err, faq))
	case "warn":
		fmt.Fprintf(os.Stderr, "WARNING: %v\nSee %v\n\n", err, faq)
		return true
	case "ignore":
		return true
	default:
		panic("invalid " + env + " value: " + os.Getenv(env))
	}
}

var globalMutex sync.RWMutex


var GlobalFiles *Files = new(Files)



var GlobalTypes *Types = new(Types)





var NotFound = errors.New("not found")




type Files struct {
	
	
	
	
	
	
	
	
	
	
	
	
	descsByName map[protoreflect.FullName]any
	filesByPath map[string][]protoreflect.FileDescriptor
	numFiles    int
}

type packageDescriptor struct {
	files []protoreflect.FileDescriptor
}








func (r *Files) RegisterFile(file protoreflect.FileDescriptor) error {
	if r == GlobalFiles {
		globalMutex.Lock()
		defer globalMutex.Unlock()
	}
	if r.descsByName == nil {
		r.descsByName = map[protoreflect.FullName]any{
			"": &packageDescriptor{},
		}
		r.filesByPath = make(map[string][]protoreflect.FileDescriptor)
	}
	path := file.Path()
	if prev := r.filesByPath[path]; len(prev) > 0 {
		r.checkGenProtoConflict(path)
		err := errors.New("file %q is already registered", file.Path())
		err = amendErrorWithCaller(err, prev[0], file)
		if !(r == GlobalFiles && ignoreConflict(file, err)) {
			return err
		}
	}

	for name := file.Package(); name != ""; name = name.Parent() {
		switch prev := r.descsByName[name]; prev.(type) {
		case nil, *packageDescriptor:
		default:
			err := errors.New("file %q has a package name conflict over %v", file.Path(), name)
			err = amendErrorWithCaller(err, prev, file)
			if r == GlobalFiles && ignoreConflict(file, err) {
				err = nil
			}
			return err
		}
	}
	var err error
	var hasConflict bool
	rangeTopLevelDescriptors(file, func(d protoreflect.Descriptor) {
		if prev := r.descsByName[d.FullName()]; prev != nil {
			hasConflict = true
			err = errors.New("file %q has a name conflict over %v", file.Path(), d.FullName())
			err = amendErrorWithCaller(err, prev, file)
			if r == GlobalFiles && ignoreConflict(d, err) {
				err = nil
			}
		}
	})
	if hasConflict {
		return err
	}

	for name := file.Package(); name != ""; name = name.Parent() {
		if r.descsByName[name] == nil {
			r.descsByName[name] = &packageDescriptor{}
		}
	}
	p := r.descsByName[file.Package()].(*packageDescriptor)
	p.files = append(p.files, file)
	rangeTopLevelDescriptors(file, func(d protoreflect.Descriptor) {
		r.descsByName[d.FullName()] = d
	})
	r.filesByPath[path] = append(r.filesByPath[path], file)
	r.numFiles++
	return nil
}






func (r *Files) checkGenProtoConflict(path string) {
	if r != GlobalFiles {
		return
	}
	var prevPath string
	const prevModule = "google.golang.org/genproto"
	const prevVersion = "cb27e3aa (May 26th, 2020)"
	switch path {
	case "google/protobuf/field_mask.proto":
		prevPath = prevModule + "/protobuf/field_mask"
	case "google/protobuf/api.proto":
		prevPath = prevModule + "/protobuf/api"
	case "google/protobuf/type.proto":
		prevPath = prevModule + "/protobuf/ptype"
	case "google/protobuf/source_context.proto":
		prevPath = prevModule + "/protobuf/source_context"
	default:
		return
	}
	pkgName := strings.TrimSuffix(strings.TrimPrefix(path, "google/protobuf/"), ".proto")
	pkgName = strings.Replace(pkgName, "_", "", -1) + "pb" 
	currPath := "google.golang.org/protobuf/types/known/" + pkgName
	panic(fmt.Sprintf(""+
		"duplicate registration of %q\n"+
		"\n"+
		"The generated definition for this file has moved:\n"+
		"\tfrom: %q\n"+
		"\tto:   %q\n"+
		"A dependency on the %q module must\n"+
		"be at version %v or higher.\n"+
		"\n"+
		"Upgrade the dependency by running:\n"+
		"\tgo get -u %v\n",
		path, prevPath, currPath, prevModule, prevVersion, prevPath))
}




func (r *Files) FindDescriptorByName(name protoreflect.FullName) (protoreflect.Descriptor, error) {
	if r == nil {
		return nil, NotFound
	}
	if r == GlobalFiles {
		globalMutex.RLock()
		defer globalMutex.RUnlock()
	}
	prefix := name
	suffix := nameSuffix("")
	for prefix != "" {
		if d, ok := r.descsByName[prefix]; ok {
			switch d := d.(type) {
			case protoreflect.EnumDescriptor:
				if d.FullName() == name {
					return d, nil
				}
			case protoreflect.EnumValueDescriptor:
				if d.FullName() == name {
					return d, nil
				}
			case protoreflect.MessageDescriptor:
				if d.FullName() == name {
					return d, nil
				}
				if d := findDescriptorInMessage(d, suffix); d != nil && d.FullName() == name {
					return d, nil
				}
			case protoreflect.ExtensionDescriptor:
				if d.FullName() == name {
					return d, nil
				}
			case protoreflect.ServiceDescriptor:
				if d.FullName() == name {
					return d, nil
				}
				if d := d.Methods().ByName(suffix.Pop()); d != nil && d.FullName() == name {
					return d, nil
				}
			}
			return nil, NotFound
		}
		prefix = prefix.Parent()
		suffix = nameSuffix(name[len(prefix)+len("."):])
	}
	return nil, NotFound
}

func findDescriptorInMessage(md protoreflect.MessageDescriptor, suffix nameSuffix) protoreflect.Descriptor {
	name := suffix.Pop()
	if suffix == "" {
		if ed := md.Enums().ByName(name); ed != nil {
			return ed
		}
		for i := md.Enums().Len() - 1; i >= 0; i-- {
			if vd := md.Enums().Get(i).Values().ByName(name); vd != nil {
				return vd
			}
		}
		if xd := md.Extensions().ByName(name); xd != nil {
			return xd
		}
		if fd := md.Fields().ByName(name); fd != nil {
			return fd
		}
		if od := md.Oneofs().ByName(name); od != nil {
			return od
		}
	}
	if md := md.Messages().ByName(name); md != nil {
		if suffix == "" {
			return md
		}
		return findDescriptorInMessage(md, suffix)
	}
	return nil
}

type nameSuffix string

func (s *nameSuffix) Pop() (name protoreflect.Name) {
	if i := strings.IndexByte(string(*s), '.'); i >= 0 {
		name, *s = protoreflect.Name((*s)[:i]), (*s)[i+1:]
	} else {
		name, *s = protoreflect.Name((*s)), ""
	}
	return name
}





func (r *Files) FindFileByPath(path string) (protoreflect.FileDescriptor, error) {
	if r == nil {
		return nil, NotFound
	}
	if r == GlobalFiles {
		globalMutex.RLock()
		defer globalMutex.RUnlock()
	}
	fds := r.filesByPath[path]
	switch len(fds) {
	case 0:
		return nil, NotFound
	case 1:
		return fds[0], nil
	default:
		return nil, errors.New("multiple files named %q", path)
	}
}



func (r *Files) NumFiles() int {
	if r == nil {
		return 0
	}
	if r == GlobalFiles {
		globalMutex.RLock()
		defer globalMutex.RUnlock()
	}
	return r.numFiles
}




func (r *Files) RangeFiles(f func(protoreflect.FileDescriptor) bool) {
	if r == nil {
		return
	}
	if r == GlobalFiles {
		globalMutex.RLock()
		defer globalMutex.RUnlock()
	}
	for _, files := range r.filesByPath {
		for _, file := range files {
			if !f(file) {
				return
			}
		}
	}
}


func (r *Files) NumFilesByPackage(name protoreflect.FullName) int {
	if r == nil {
		return 0
	}
	if r == GlobalFiles {
		globalMutex.RLock()
		defer globalMutex.RUnlock()
	}
	p, ok := r.descsByName[name].(*packageDescriptor)
	if !ok {
		return 0
	}
	return len(p.files)
}



func (r *Files) RangeFilesByPackage(name protoreflect.FullName, f func(protoreflect.FileDescriptor) bool) {
	if r == nil {
		return
	}
	if r == GlobalFiles {
		globalMutex.RLock()
		defer globalMutex.RUnlock()
	}
	p, ok := r.descsByName[name].(*packageDescriptor)
	if !ok {
		return
	}
	for _, file := range p.files {
		if !f(file) {
			return
		}
	}
}



func rangeTopLevelDescriptors(fd protoreflect.FileDescriptor, f func(protoreflect.Descriptor)) {
	eds := fd.Enums()
	for i := eds.Len() - 1; i >= 0; i-- {
		f(eds.Get(i))
		vds := eds.Get(i).Values()
		for i := vds.Len() - 1; i >= 0; i-- {
			f(vds.Get(i))
		}
	}
	mds := fd.Messages()
	for i := mds.Len() - 1; i >= 0; i-- {
		f(mds.Get(i))
	}
	xds := fd.Extensions()
	for i := xds.Len() - 1; i >= 0; i-- {
		f(xds.Get(i))
	}
	sds := fd.Services()
	for i := sds.Len() - 1; i >= 0; i-- {
		f(sds.Get(i))
	}
}







type MessageTypeResolver interface {
	
	
	
	
	FindMessageByName(message protoreflect.FullName) (protoreflect.MessageType, error)

	
	
	
	
	FindMessageByURL(url string) (protoreflect.MessageType, error)
}







type ExtensionTypeResolver interface {
	
	
	
	
	
	
	FindExtensionByName(field protoreflect.FullName) (protoreflect.ExtensionType, error)

	
	
	
	
	FindExtensionByNumber(message protoreflect.FullName, field protoreflect.FieldNumber) (protoreflect.ExtensionType, error)
}

var (
	_ MessageTypeResolver   = (*Types)(nil)
	_ ExtensionTypeResolver = (*Types)(nil)
)



type Types struct {
	typesByName         typesByName
	extensionsByMessage extensionsByMessage

	numEnums      int
	numMessages   int
	numExtensions int
}

type (
	typesByName         map[protoreflect.FullName]any
	extensionsByMessage map[protoreflect.FullName]extensionsByNumber
	extensionsByNumber  map[protoreflect.FieldNumber]protoreflect.ExtensionType
)




func (r *Types) RegisterMessage(mt protoreflect.MessageType) error {
	
	
	md := mt.Descriptor()

	if r == GlobalTypes {
		globalMutex.Lock()
		defer globalMutex.Unlock()
	}

	if err := r.register("message", md, mt); err != nil {
		return err
	}
	r.numMessages++
	return nil
}




func (r *Types) RegisterEnum(et protoreflect.EnumType) error {
	
	
	ed := et.Descriptor()

	if r == GlobalTypes {
		globalMutex.Lock()
		defer globalMutex.Unlock()
	}

	if err := r.register("enum", ed, et); err != nil {
		return err
	}
	r.numEnums++
	return nil
}




func (r *Types) RegisterExtension(xt protoreflect.ExtensionType) error {
	
	
	
	
	
	xd := xt.TypeDescriptor()

	if r == GlobalTypes {
		globalMutex.Lock()
		defer globalMutex.Unlock()
	}

	field := xd.Number()
	message := xd.ContainingMessage().FullName()
	if prev := r.extensionsByMessage[message][field]; prev != nil {
		err := errors.New("extension number %d is already registered on message %v", field, message)
		err = amendErrorWithCaller(err, prev, xt)
		if !(r == GlobalTypes && ignoreConflict(xd, err)) {
			return err
		}
	}

	if err := r.register("extension", xd, xt); err != nil {
		return err
	}
	if r.extensionsByMessage == nil {
		r.extensionsByMessage = make(extensionsByMessage)
	}
	if r.extensionsByMessage[message] == nil {
		r.extensionsByMessage[message] = make(extensionsByNumber)
	}
	r.extensionsByMessage[message][field] = xt
	r.numExtensions++
	return nil
}

func (r *Types) register(kind string, desc protoreflect.Descriptor, typ any) error {
	name := desc.FullName()
	prev := r.typesByName[name]
	if prev != nil {
		err := errors.New("%v %v is already registered", kind, name)
		err = amendErrorWithCaller(err, prev, typ)
		if !(r == GlobalTypes && ignoreConflict(desc, err)) {
			return err
		}
	}
	if r.typesByName == nil {
		r.typesByName = make(typesByName)
	}
	r.typesByName[name] = typ
	return nil
}





func (r *Types) FindEnumByName(enum protoreflect.FullName) (protoreflect.EnumType, error) {
	if r == nil {
		return nil, NotFound
	}
	if r == GlobalTypes {
		globalMutex.RLock()
		defer globalMutex.RUnlock()
	}
	if v := r.typesByName[enum]; v != nil {
		if et, _ := v.(protoreflect.EnumType); et != nil {
			return et, nil
		}
		return nil, errors.New("found wrong type: got %v, want enum", typeName(v))
	}
	return nil, NotFound
}





func (r *Types) FindMessageByName(message protoreflect.FullName) (protoreflect.MessageType, error) {
	if r == nil {
		return nil, NotFound
	}
	if r == GlobalTypes {
		globalMutex.RLock()
		defer globalMutex.RUnlock()
	}
	if v := r.typesByName[message]; v != nil {
		if mt, _ := v.(protoreflect.MessageType); mt != nil {
			return mt, nil
		}
		return nil, errors.New("found wrong type: got %v, want message", typeName(v))
	}
	return nil, NotFound
}





func (r *Types) FindMessageByURL(url string) (protoreflect.MessageType, error) {
	
	
	if r == nil {
		return nil, NotFound
	}
	if r == GlobalTypes {
		globalMutex.RLock()
		defer globalMutex.RUnlock()
	}
	message := protoreflect.FullName(url)
	if i := strings.LastIndexByte(url, '/'); i >= 0 {
		message = message[i+len("/"):]
	}

	if v := r.typesByName[message]; v != nil {
		if mt, _ := v.(protoreflect.MessageType); mt != nil {
			return mt, nil
		}
		return nil, errors.New("found wrong type: got %v, want message", typeName(v))
	}
	return nil, NotFound
}







func (r *Types) FindExtensionByName(field protoreflect.FullName) (protoreflect.ExtensionType, error) {
	if r == nil {
		return nil, NotFound
	}
	if r == GlobalTypes {
		globalMutex.RLock()
		defer globalMutex.RUnlock()
	}
	if v := r.typesByName[field]; v != nil {
		if xt, _ := v.(protoreflect.ExtensionType); xt != nil {
			return xt, nil
		}

		
		
		
		
		
		
		if flags.ProtoLegacy {
			if _, ok := v.(protoreflect.MessageType); ok {
				field := field.Append(messageset.ExtensionName)
				if v := r.typesByName[field]; v != nil {
					if xt, _ := v.(protoreflect.ExtensionType); xt != nil {
						if messageset.IsMessageSetExtension(xt.TypeDescriptor()) {
							return xt, nil
						}
					}
				}
			}
		}

		return nil, errors.New("found wrong type: got %v, want extension", typeName(v))
	}
	return nil, NotFound
}





func (r *Types) FindExtensionByNumber(message protoreflect.FullName, field protoreflect.FieldNumber) (protoreflect.ExtensionType, error) {
	if r == nil {
		return nil, NotFound
	}
	if r == GlobalTypes {
		globalMutex.RLock()
		defer globalMutex.RUnlock()
	}
	if xt, ok := r.extensionsByMessage[message][field]; ok {
		return xt, nil
	}
	return nil, NotFound
}


func (r *Types) NumEnums() int {
	if r == nil {
		return 0
	}
	if r == GlobalTypes {
		globalMutex.RLock()
		defer globalMutex.RUnlock()
	}
	return r.numEnums
}



func (r *Types) RangeEnums(f func(protoreflect.EnumType) bool) {
	if r == nil {
		return
	}
	if r == GlobalTypes {
		globalMutex.RLock()
		defer globalMutex.RUnlock()
	}
	for _, typ := range r.typesByName {
		if et, ok := typ.(protoreflect.EnumType); ok {
			if !f(et) {
				return
			}
		}
	}
}


func (r *Types) NumMessages() int {
	if r == nil {
		return 0
	}
	if r == GlobalTypes {
		globalMutex.RLock()
		defer globalMutex.RUnlock()
	}
	return r.numMessages
}



func (r *Types) RangeMessages(f func(protoreflect.MessageType) bool) {
	if r == nil {
		return
	}
	if r == GlobalTypes {
		globalMutex.RLock()
		defer globalMutex.RUnlock()
	}
	for _, typ := range r.typesByName {
		if mt, ok := typ.(protoreflect.MessageType); ok {
			if !f(mt) {
				return
			}
		}
	}
}


func (r *Types) NumExtensions() int {
	if r == nil {
		return 0
	}
	if r == GlobalTypes {
		globalMutex.RLock()
		defer globalMutex.RUnlock()
	}
	return r.numExtensions
}



func (r *Types) RangeExtensions(f func(protoreflect.ExtensionType) bool) {
	if r == nil {
		return
	}
	if r == GlobalTypes {
		globalMutex.RLock()
		defer globalMutex.RUnlock()
	}
	for _, typ := range r.typesByName {
		if xt, ok := typ.(protoreflect.ExtensionType); ok {
			if !f(xt) {
				return
			}
		}
	}
}



func (r *Types) NumExtensionsByMessage(message protoreflect.FullName) int {
	if r == nil {
		return 0
	}
	if r == GlobalTypes {
		globalMutex.RLock()
		defer globalMutex.RUnlock()
	}
	return len(r.extensionsByMessage[message])
}



func (r *Types) RangeExtensionsByMessage(message protoreflect.FullName, f func(protoreflect.ExtensionType) bool) {
	if r == nil {
		return
	}
	if r == GlobalTypes {
		globalMutex.RLock()
		defer globalMutex.RUnlock()
	}
	for _, xt := range r.extensionsByMessage[message] {
		if !f(xt) {
			return
		}
	}
}

func typeName(t any) string {
	switch t.(type) {
	case protoreflect.EnumType:
		return "enum"
	case protoreflect.MessageType:
		return "message"
	case protoreflect.ExtensionType:
		return "extension"
	default:
		return fmt.Sprintf("%T", t)
	}
}

func amendErrorWithCaller(err error, prev, curr any) error {
	prevPkg := goPackage(prev)
	currPkg := goPackage(curr)
	if prevPkg == "" || currPkg == "" || prevPkg == currPkg {
		return err
	}
	return errors.New("%s\n\tpreviously from: %q\n\tcurrently from:  %q", err, prevPkg, currPkg)
}

func goPackage(v any) string {
	switch d := v.(type) {
	case protoreflect.EnumType:
		v = d.Descriptor()
	case protoreflect.MessageType:
		v = d.Descriptor()
	case protoreflect.ExtensionType:
		v = d.TypeDescriptor()
	}
	if d, ok := v.(protoreflect.Descriptor); ok {
		v = d.ParentFile()
	}
	if d, ok := v.(interface{ GoPackagePath() string }); ok {
		return d.GoPackagePath()
	}
	return ""
}
