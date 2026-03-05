





package bson

import (
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
)















var DefaultRegistry = NewRegistry()






func NewRegistryBuilder() *bsoncodec.RegistryBuilder {
	rb := bsoncodec.NewRegistryBuilder()
	bsoncodec.DefaultValueEncoders{}.RegisterDefaultEncoders(rb)
	bsoncodec.DefaultValueDecoders{}.RegisterDefaultDecoders(rb)
	primitiveCodecs.RegisterPrimitiveCodecs(rb)
	return rb
}




func NewRegistry() *bsoncodec.Registry {
	return NewRegistryBuilder().Build()
}
