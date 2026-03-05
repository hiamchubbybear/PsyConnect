



package impl

import (
	"bytes"
	"compress/gzip"
	"io"
	"sync"

	"google.golang.org/protobuf/internal/filedesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)









type (
	enumV1 interface {
		EnumDescriptor() ([]byte, []int)
	}
	messageV1 interface {
		Descriptor() ([]byte, []int)
	}
)

var legacyFileDescCache sync.Map 







func legacyLoadFileDesc(b []byte) protoreflect.FileDescriptor {
	
	if fd, ok := legacyFileDescCache.Load(&b[0]); ok {
		return fd.(protoreflect.FileDescriptor)
	}

	
	zr, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		panic(err)
	}
	b2, err := io.ReadAll(zr)
	if err != nil {
		panic(err)
	}

	fd := filedesc.Builder{
		RawDescriptor: b2,
		FileRegistry:  resolverOnly{protoregistry.GlobalFiles}, 
	}.Build().File
	if fd, ok := legacyFileDescCache.LoadOrStore(&b[0], fd); ok {
		return fd.(protoreflect.FileDescriptor)
	}
	return fd
}

type resolverOnly struct {
	reg *protoregistry.Files
}

func (r resolverOnly) FindFileByPath(path string) (protoreflect.FileDescriptor, error) {
	return r.reg.FindFileByPath(path)
}
func (r resolverOnly) FindDescriptorByName(name protoreflect.FullName) (protoreflect.Descriptor, error) {
	return r.reg.FindDescriptorByName(name)
}
func (resolverOnly) RegisterFile(protoreflect.FileDescriptor) error {
	return nil
}
