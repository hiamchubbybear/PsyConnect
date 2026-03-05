





package options

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)


type DataKeyOptions struct {
	KeyAltNames []string
	KeyMaterial []byte
	MasterKey   bsoncore.Document
}


func DataKey() *DataKeyOptions {
	return &DataKeyOptions{}
}


func (dko *DataKeyOptions) SetKeyAltNames(names []string) *DataKeyOptions {
	dko.KeyAltNames = names
	return dko
}


func (dko *DataKeyOptions) SetMasterKey(key bsoncore.Document) *DataKeyOptions {
	dko.MasterKey = key
	return dko
}


func (dko *DataKeyOptions) SetKeyMaterial(keyMaterial []byte) *DataKeyOptions {
	dko.KeyMaterial = keyMaterial
	return dko
}


type QueryType int


const (
	QueryTypeEquality QueryType = 1
)


type ExplicitEncryptionOptions struct {
	KeyID            *primitive.Binary
	KeyAltName       *string
	Algorithm        string
	QueryType        string
	ContentionFactor *int64
	RangeOptions     *ExplicitRangeOptions
}


type ExplicitRangeOptions struct {
	Min        *bsoncore.Value
	Max        *bsoncore.Value
	Sparsity   *int64
	TrimFactor *int32
	Precision  *int32
}


func ExplicitEncryption() *ExplicitEncryptionOptions {
	return &ExplicitEncryptionOptions{}
}


func (eeo *ExplicitEncryptionOptions) SetKeyID(keyID primitive.Binary) *ExplicitEncryptionOptions {
	eeo.KeyID = &keyID
	return eeo
}


func (eeo *ExplicitEncryptionOptions) SetKeyAltName(keyAltName string) *ExplicitEncryptionOptions {
	eeo.KeyAltName = &keyAltName
	return eeo
}


func (eeo *ExplicitEncryptionOptions) SetAlgorithm(algorithm string) *ExplicitEncryptionOptions {
	eeo.Algorithm = algorithm
	return eeo
}


func (eeo *ExplicitEncryptionOptions) SetQueryType(queryType string) *ExplicitEncryptionOptions {
	eeo.QueryType = queryType
	return eeo
}


func (eeo *ExplicitEncryptionOptions) SetContentionFactor(contentionFactor int64) *ExplicitEncryptionOptions {
	eeo.ContentionFactor = &contentionFactor
	return eeo
}


func (eeo *ExplicitEncryptionOptions) SetRangeOptions(ro ExplicitRangeOptions) *ExplicitEncryptionOptions {
	eeo.RangeOptions = &ro
	return eeo
}



type RewrapManyDataKeyOptions struct {
	
	Provider *string

	
	MasterKey bsoncore.Document
}


func RewrapManyDataKey() *RewrapManyDataKeyOptions {
	return new(RewrapManyDataKeyOptions)
}


func (rmdko *RewrapManyDataKeyOptions) SetProvider(provider string) *RewrapManyDataKeyOptions {
	rmdko.Provider = &provider
	return rmdko
}


func (rmdko *RewrapManyDataKeyOptions) SetMasterKey(masterKey bsoncore.Document) *RewrapManyDataKeyOptions {
	rmdko.MasterKey = masterKey
	return rmdko
}






func MergeRewrapManyDataKeyOptions(opts ...*RewrapManyDataKeyOptions) *RewrapManyDataKeyOptions {
	rmdkOpts := RewrapManyDataKey()
	for _, rmdko := range opts {
		if rmdko == nil {
			continue
		}
		if provider := rmdko.Provider; provider != nil {
			rmdkOpts.Provider = provider
		}
		if masterKey := rmdko.MasterKey; masterKey != nil {
			rmdkOpts.MasterKey = masterKey
		}
	}
	return rmdkOpts
}
