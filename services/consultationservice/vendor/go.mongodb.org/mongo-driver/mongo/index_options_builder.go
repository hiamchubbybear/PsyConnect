





package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
)




type IndexOptionsBuilder struct {
	document bson.D
}




func NewIndexOptionsBuilder() *IndexOptionsBuilder {
	return &IndexOptionsBuilder{}
}




func (iob *IndexOptionsBuilder) Background(background bool) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"background", background})
	return iob
}




func (iob *IndexOptionsBuilder) ExpireAfterSeconds(expireAfterSeconds int32) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"expireAfterSeconds", expireAfterSeconds})
	return iob
}




func (iob *IndexOptionsBuilder) Name(name string) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"name", name})
	return iob
}




func (iob *IndexOptionsBuilder) Sparse(sparse bool) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"sparse", sparse})
	return iob
}




func (iob *IndexOptionsBuilder) StorageEngine(storageEngine interface{}) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"storageEngine", storageEngine})
	return iob
}




func (iob *IndexOptionsBuilder) Unique(unique bool) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"unique", unique})
	return iob
}




func (iob *IndexOptionsBuilder) Version(version int32) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"v", version})
	return iob
}




func (iob *IndexOptionsBuilder) DefaultLanguage(defaultLanguage string) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"default_language", defaultLanguage})
	return iob
}




func (iob *IndexOptionsBuilder) LanguageOverride(languageOverride string) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"language_override", languageOverride})
	return iob
}




func (iob *IndexOptionsBuilder) TextVersion(textVersion int32) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"textIndexVersion", textVersion})
	return iob
}




func (iob *IndexOptionsBuilder) Weights(weights interface{}) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"weights", weights})
	return iob
}




func (iob *IndexOptionsBuilder) SphereVersion(sphereVersion int32) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"2dsphereIndexVersion", sphereVersion})
	return iob
}




func (iob *IndexOptionsBuilder) Bits(bits int32) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"bits", bits})
	return iob
}




func (iob *IndexOptionsBuilder) Max(max float64) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"max", max})
	return iob
}




func (iob *IndexOptionsBuilder) Min(min float64) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"min", min})
	return iob
}




func (iob *IndexOptionsBuilder) BucketSize(bucketSize int32) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"bucketSize", bucketSize})
	return iob
}




func (iob *IndexOptionsBuilder) PartialFilterExpression(partialFilterExpression interface{}) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"partialFilterExpression", partialFilterExpression})
	return iob
}




func (iob *IndexOptionsBuilder) Collation(collation interface{}) *IndexOptionsBuilder {
	iob.document = append(iob.document, bson.E{"collation", collation})
	return iob
}




func (iob *IndexOptionsBuilder) Build() bson.D {
	return iob.document
}
