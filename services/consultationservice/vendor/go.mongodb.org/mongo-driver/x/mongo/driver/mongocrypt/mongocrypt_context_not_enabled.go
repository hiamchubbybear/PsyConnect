





//go:build !cse
// +build !cse

package mongocrypt

import (
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)


type Context struct{}


func (c *Context) State() State {
	panic(cseNotSupportedMsg)
}


func (c *Context) NextOperation() (bsoncore.Document, error) {
	panic(cseNotSupportedMsg)
}


func (c *Context) AddOperationResult(bsoncore.Document) error {
	panic(cseNotSupportedMsg)
}


func (c *Context) CompleteOperation() error {
	panic(cseNotSupportedMsg)
}


func (c *Context) NextKmsContext() *KmsContext {
	panic(cseNotSupportedMsg)
}


func (c *Context) FinishKmsContexts() error {
	panic(cseNotSupportedMsg)
}


func (c *Context) Finish() (bsoncore.Document, error) {
	panic(cseNotSupportedMsg)
}


func (c *Context) Close() {
	panic(cseNotSupportedMsg)
}


func (c *Context) ProvideKmsProviders(bsoncore.Document) error {
	panic(cseNotSupportedMsg)
}
