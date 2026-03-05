



package proto

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"
)




func Reset(m Message) {
	if mr, ok := m.(interface{ Reset() }); ok && hasProtoMethods {
		mr.Reset()
		return
	}
	resetMessage(m.ProtoReflect())
}

func resetMessage(m protoreflect.Message) {
	if !m.IsValid() {
		panic(fmt.Sprintf("cannot reset invalid %v message", m.Descriptor().FullName()))
	}

	
	fds := m.Descriptor().Fields()
	for i := 0; i < fds.Len(); i++ {
		m.Clear(fds.Get(i))
	}

	
	m.Range(func(fd protoreflect.FieldDescriptor, _ protoreflect.Value) bool {
		m.Clear(fd)
		return true
	})

	
	m.SetUnknown(nil)
}
