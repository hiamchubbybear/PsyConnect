





//go:build cse
// +build cse

package mongocrypt


import "C"


type KmsContext struct {
	wrapped *C.mongocrypt_kms_ctx_t
}


func newKmsContext(wrapped *C.mongocrypt_kms_ctx_t) *KmsContext {
	return &KmsContext{
		wrapped: wrapped,
	}
}


func (kc *KmsContext) HostName() (string, error) {
	var hostname *C.char 
	if ok := C.mongocrypt_kms_ctx_endpoint(kc.wrapped, &hostname); !ok {
		return "", kc.createErrorFromStatus()
	}
	return C.GoString(hostname), nil
}


func (kc *KmsContext) KMSProvider() string {
	kmsProvider := C.mongocrypt_kms_ctx_get_kms_provider(kc.wrapped, nil)
	return C.GoString(kmsProvider)
}


func (kc *KmsContext) Message() ([]byte, error) {
	msgBinary := newBinary()
	defer msgBinary.close()

	if ok := C.mongocrypt_kms_ctx_message(kc.wrapped, msgBinary.wrapped); !ok {
		return nil, kc.createErrorFromStatus()
	}
	return msgBinary.toBytes(), nil
}



func (kc *KmsContext) BytesNeeded() int32 {
	return int32(C.mongocrypt_kms_ctx_bytes_needed(kc.wrapped))
}


func (kc *KmsContext) FeedResponse(response []byte) error {
	responseBinary := newBinaryFromBytes(response)
	defer responseBinary.close()

	if ok := C.mongocrypt_kms_ctx_feed(kc.wrapped, responseBinary.wrapped); !ok {
		return kc.createErrorFromStatus()
	}
	return nil
}


func (kc *KmsContext) createErrorFromStatus() error {
	status := C.mongocrypt_status_new()
	defer C.mongocrypt_status_destroy(status)
	C.mongocrypt_kms_ctx_status(kc.wrapped, status)
	return errorFromStatus(status)
}
