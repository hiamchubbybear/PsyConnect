



package impl

import (
	"strconv"
	"sync/atomic"
	"unsafe"

	"google.golang.org/protobuf/reflect/protoreflect"
)

func (Export) UnmarshalField(msg any, fieldNum int32) {
	UnmarshalField(msg.(protoreflect.ProtoMessage).ProtoReflect(), protoreflect.FieldNumber(fieldNum))
}






func (Export) Present(part *uint32, num uint32) bool {
	
	
	raceDetectHookPresent(part, num)
	return atomic.LoadUint32(part)&(1<<(num%32)) > 0
}





func (Export) SetPresent(part *uint32, num uint32, size uint32) {
	
	
	raceDetectHookSetPresent(part, num, presenceSize(size))
	for {
		old := atomic.LoadUint32(part)
		if atomic.CompareAndSwapUint32(part, old, old|(1<<(num%32))) {
			return
		}
	}
}




func (Export) SetPresentNonAtomic(part *uint32, num uint32, size uint32) {
	
	
	raceDetectHookSetPresent(part, num, presenceSize(size))
	*part |= 1 << (num % 32)
}




func (Export) ClearPresent(part *uint32, num uint32) {
	
	
	raceDetectHookClearPresent(part, num)
	for {
		old := atomic.LoadUint32(part)
		if atomic.CompareAndSwapUint32(part, old, old&^(1<<(num%32))) {
			return
		}
	}
}




func interfaceToPointer(i *any) pointer {
	return pointer{p: (*[2]unsafe.Pointer)(unsafe.Pointer(i))[1]}
}

func (p pointer) atomicGetPointer() pointer {
	return pointer{p: atomic.LoadPointer((*unsafe.Pointer)(p.p))}
}

func (p pointer) atomicSetPointer(q pointer) {
	atomic.StorePointer((*unsafe.Pointer)(p.p), q.p)
}





func (Export) AtomicCheckPointerIsNil(ptr any) bool {
	return interfaceToPointer(&ptr).atomicGetPointer().IsNil()
}





func (Export) AtomicSetPointer(dstPtr, valPtr any) {
	interfaceToPointer(&dstPtr).atomicSetPointer(interfaceToPointer(&valPtr))
}



func (Export) AtomicLoadPointer(ptr Pointer, dst Pointer) {
	*(*unsafe.Pointer)(unsafe.Pointer(dst)) = atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(ptr)))
}






func (Export) AtomicInitializePointer(ptr Pointer, dst Pointer) {
	if !atomic.CompareAndSwapPointer((*unsafe.Pointer)(ptr), unsafe.Pointer(nil), *(*unsafe.Pointer)(dst)) {
		*(*unsafe.Pointer)(unsafe.Pointer(dst)) = atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(ptr)))
	}
}



func (Export) MessageFieldStringOf(md protoreflect.MessageDescriptor, n protoreflect.FieldNumber) string {
	fd := md.Fields().ByNumber(n)
	if fd != nil {
		return string(fd.Name())
	}
	return strconv.Itoa(int(n))
}
