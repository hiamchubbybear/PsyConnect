package danger

import (
	"reflect"
	"unsafe"
)









type TypeID unsafe.Pointer

func MakeTypeID(t reflect.Type) TypeID {
	
	
	
	return TypeID((*[2]unsafe.Pointer)(unsafe.Pointer(&t))[1])
}
