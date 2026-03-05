





//go:build !cse
// +build !cse








package mongocrypt

import (
	"context"

	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver/mongocrypt/options"
)

const cseNotSupportedMsg = "client-side encryption not enabled. add the cse build tag to support"


type MongoCrypt struct{}



func Version() string {
	return ""
}


func NewMongoCrypt(*options.MongoCryptOptions) (*MongoCrypt, error) {
	panic(cseNotSupportedMsg)
}


func (m *MongoCrypt) CreateEncryptionContext(string, bsoncore.Document) (*Context, error) {
	panic(cseNotSupportedMsg)
}


func (m *MongoCrypt) CreateExplicitEncryptionExpressionContext(bsoncore.Document, *options.ExplicitEncryptionOptions) (*Context, error) {
	panic(cseNotSupportedMsg)
}


func (m *MongoCrypt) CreateDecryptionContext(bsoncore.Document) (*Context, error) {
	panic(cseNotSupportedMsg)
}


func (m *MongoCrypt) CreateDataKeyContext(string, *options.DataKeyOptions) (*Context, error) {
	panic(cseNotSupportedMsg)
}


func (m *MongoCrypt) CreateExplicitEncryptionContext(bsoncore.Document, *options.ExplicitEncryptionOptions) (*Context, error) {
	panic(cseNotSupportedMsg)
}


func (m *MongoCrypt) RewrapDataKeyContext([]byte, *options.RewrapManyDataKeyOptions) (*Context, error) {
	panic(cseNotSupportedMsg)
}


func (m *MongoCrypt) CreateExplicitDecryptionContext(bsoncore.Document) (*Context, error) {
	panic(cseNotSupportedMsg)
}



func (m *MongoCrypt) CryptSharedLibVersion() uint64 {
	panic(cseNotSupportedMsg)
}



func (m *MongoCrypt) CryptSharedLibVersionString() string {
	panic(cseNotSupportedMsg)
}


func (m *MongoCrypt) Close() {
	panic(cseNotSupportedMsg)
}


func (m *MongoCrypt) GetKmsProviders(context.Context) (bsoncore.Document, error) {
	panic(cseNotSupportedMsg)
}
