





//go:build cse
// +build cse

package mongocrypt


import "C"
import (
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)


type Context struct {
	wrapped *C.mongocrypt_ctx_t
}


func newContext(wrapped *C.mongocrypt_ctx_t) *Context {
	return &Context{
		wrapped: wrapped,
	}
}


func (c *Context) State() State {
	return State(int(C.mongocrypt_ctx_state(c.wrapped)))
}


func (c *Context) NextOperation() (bsoncore.Document, error) {
	opDocBinary := newBinary() 
	defer opDocBinary.close()

	if ok := C.mongocrypt_ctx_mongo_op(c.wrapped, opDocBinary.wrapped); !ok {
		return nil, c.createErrorFromStatus()
	}
	return opDocBinary.toBytes(), nil
}


func (c *Context) AddOperationResult(result bsoncore.Document) error {
	resultBinary := newBinaryFromBytes(result)
	defer resultBinary.close()

	if ok := C.mongocrypt_ctx_mongo_feed(c.wrapped, resultBinary.wrapped); !ok {
		return c.createErrorFromStatus()
	}
	return nil
}


func (c *Context) CompleteOperation() error {
	if ok := C.mongocrypt_ctx_mongo_done(c.wrapped); !ok {
		return c.createErrorFromStatus()
	}
	return nil
}


func (c *Context) NextKmsContext() *KmsContext {
	ctx := C.mongocrypt_ctx_next_kms_ctx(c.wrapped)
	if ctx == nil {
		return nil
	}
	return newKmsContext(ctx)
}


func (c *Context) FinishKmsContexts() error {
	if ok := C.mongocrypt_ctx_kms_done(c.wrapped); !ok {
		return c.createErrorFromStatus()
	}
	return nil
}


func (c *Context) Finish() (bsoncore.Document, error) {
	docBinary := newBinary() 
	defer docBinary.close()

	if ok := C.mongocrypt_ctx_finalize(c.wrapped, docBinary.wrapped); !ok {
		return nil, c.createErrorFromStatus()
	}
	return docBinary.toBytes(), nil
}


func (c *Context) Close() {
	C.mongocrypt_ctx_destroy(c.wrapped)
}


func (c *Context) createErrorFromStatus() error {
	status := C.mongocrypt_status_new()
	defer C.mongocrypt_status_destroy(status)
	C.mongocrypt_ctx_status(c.wrapped, status)
	return errorFromStatus(status)
}


func (c *Context) ProvideKmsProviders(kmsProviders bsoncore.Document) error {
	kmsProvidersBinary := newBinaryFromBytes(kmsProviders)
	defer kmsProvidersBinary.close()

	if ok := C.mongocrypt_ctx_provide_kms_providers(c.wrapped, kmsProvidersBinary.wrapped); !ok {
		return c.createErrorFromStatus()
	}
	return nil
}
