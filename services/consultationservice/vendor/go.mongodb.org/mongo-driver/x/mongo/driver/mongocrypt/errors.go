





//go:build cse
// +build cse

package mongocrypt


import "C"
import (
	"fmt"
)


type Error struct {
	Code    int32
	Message string
}


func (e Error) Error() string {
	return fmt.Sprintf("mongocrypt error %d: %v", e.Code, e.Message)
}


func errorFromStatus(status *C.mongocrypt_status_t) error {
	cCode := C.mongocrypt_status_code(status) 
	
	
	cMsg := C.mongocrypt_status_message(status, nil) 
	var msg string
	if cMsg != nil {
		msg = C.GoString(cMsg)
	}

	return Error{
		Code:    int32(cCode),
		Message: msg,
	}
}
