








package unsafeheader

import (
	"unsafe"
)







type Slice struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}







type String struct {
	Data unsafe.Pointer
	Len  int
}
