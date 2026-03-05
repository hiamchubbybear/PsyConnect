





package options

import (
	"net/http"

	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)


type MongoCryptOptions struct {
	KmsProviders               bsoncore.Document
	LocalSchemaMap             map[string]bsoncore.Document
	BypassQueryAnalysis        bool
	EncryptedFieldsMap         map[string]bsoncore.Document
	CryptSharedLibDisabled     bool
	CryptSharedLibOverridePath string
	HTTPClient                 *http.Client
}


func MongoCrypt() *MongoCryptOptions {
	return &MongoCryptOptions{}
}


func (mo *MongoCryptOptions) SetKmsProviders(kmsProviders bsoncore.Document) *MongoCryptOptions {
	mo.KmsProviders = kmsProviders
	return mo
}


func (mo *MongoCryptOptions) SetLocalSchemaMap(localSchemaMap map[string]bsoncore.Document) *MongoCryptOptions {
	mo.LocalSchemaMap = localSchemaMap
	return mo
}


func (mo *MongoCryptOptions) SetBypassQueryAnalysis(bypassQueryAnalysis bool) *MongoCryptOptions {
	mo.BypassQueryAnalysis = bypassQueryAnalysis
	return mo
}


func (mo *MongoCryptOptions) SetEncryptedFieldsMap(efcMap map[string]bsoncore.Document) *MongoCryptOptions {
	mo.EncryptedFieldsMap = efcMap
	return mo
}


func (mo *MongoCryptOptions) SetCryptSharedLibDisabled(disabled bool) *MongoCryptOptions {
	mo.CryptSharedLibDisabled = disabled
	return mo
}



func (mo *MongoCryptOptions) SetCryptSharedLibOverridePath(path string) *MongoCryptOptions {
	mo.CryptSharedLibOverridePath = path
	return mo
}


func (mo *MongoCryptOptions) SetHTTPClient(httpClient *http.Client) *MongoCryptOptions {
	mo.HTTPClient = httpClient
	return mo
}
