





//go:build cse
// +build cse

package mongocrypt


import "C"
import (
	"unsafe"
)


type binary struct {
	p       *C.uint8_t
	wrapped *C.mongocrypt_binary_t
}


func newBinary() *binary {
	return &binary{
		wrapped: C.mongocrypt_binary_new(),
	}
}


func newBinaryFromBytes(data []byte) *binary {
	if len(data) == 0 {
		return newBinary()
	}

	
	addr := (*C.uint8_t)(C.CBytes(data)) 
	dataLen := C.uint32_t(len(data))     
	return &binary{
		p:       addr,
		wrapped: C.mongocrypt_binary_new_from_data(addr, dataLen),
	}
}


func (b *binary) toBytes() []byte {
	dataPtr := C.mongocrypt_binary_data(b.wrapped) 
	dataLen := C.mongocrypt_binary_len(b.wrapped)  

	return C.GoBytes(unsafe.Pointer(dataPtr), C.int(dataLen))
}


func (b *binary) close() {
	if b.p != nil {
		C.free(unsafe.Pointer(b.p))
	}
	C.mongocrypt_binary_destroy(b.wrapped)
}
