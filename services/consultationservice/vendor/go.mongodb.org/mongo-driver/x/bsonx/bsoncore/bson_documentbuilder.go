





package bsoncore

import (
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type DocumentBuilder struct {
	doc     []byte
	indexes []int32
}


func (db *DocumentBuilder) startDocument() *DocumentBuilder {
	var index int32
	index, db.doc = AppendDocumentStart(db.doc)
	db.indexes = append(db.indexes, index)
	return db
}


func NewDocumentBuilder() *DocumentBuilder {
	return (&DocumentBuilder{}).startDocument()
}



func (db *DocumentBuilder) Build() Document {
	last := len(db.indexes) - 1
	db.doc, _ = AppendDocumentEnd(db.doc, db.indexes[last])
	db.indexes = db.indexes[:last]
	return db.doc
}


func (db *DocumentBuilder) AppendInt32(key string, i32 int32) *DocumentBuilder {
	db.doc = AppendInt32Element(db.doc, key, i32)
	return db
}



func (db *DocumentBuilder) AppendDocument(key string, doc []byte) *DocumentBuilder {
	db.doc = AppendDocumentElement(db.doc, key, doc)
	return db
}


func (db *DocumentBuilder) AppendArray(key string, arr []byte) *DocumentBuilder {
	db.doc = AppendHeader(db.doc, bsontype.Array, key)
	db.doc = AppendArray(db.doc, arr)
	return db
}


func (db *DocumentBuilder) AppendDouble(key string, f float64) *DocumentBuilder {
	db.doc = AppendDoubleElement(db.doc, key, f)
	return db
}


func (db *DocumentBuilder) AppendString(key string, str string) *DocumentBuilder {
	db.doc = AppendStringElement(db.doc, key, str)
	return db
}


func (db *DocumentBuilder) AppendObjectID(key string, oid primitive.ObjectID) *DocumentBuilder {
	db.doc = AppendObjectIDElement(db.doc, key, oid)
	return db
}



func (db *DocumentBuilder) AppendBinary(key string, subtype byte, b []byte) *DocumentBuilder {
	db.doc = AppendBinaryElement(db.doc, key, subtype, b)
	return db
}


func (db *DocumentBuilder) AppendUndefined(key string) *DocumentBuilder {
	db.doc = AppendUndefinedElement(db.doc, key)
	return db
}


func (db *DocumentBuilder) AppendBoolean(key string, b bool) *DocumentBuilder {
	db.doc = AppendBooleanElement(db.doc, key, b)
	return db
}


func (db *DocumentBuilder) AppendDateTime(key string, dt int64) *DocumentBuilder {
	db.doc = AppendDateTimeElement(db.doc, key, dt)
	return db
}


func (db *DocumentBuilder) AppendNull(key string) *DocumentBuilder {
	db.doc = AppendNullElement(db.doc, key)
	return db
}


func (db *DocumentBuilder) AppendRegex(key, pattern, options string) *DocumentBuilder {
	db.doc = AppendRegexElement(db.doc, key, pattern, options)
	return db
}


func (db *DocumentBuilder) AppendDBPointer(key string, ns string, oid primitive.ObjectID) *DocumentBuilder {
	db.doc = AppendDBPointerElement(db.doc, key, ns, oid)
	return db
}


func (db *DocumentBuilder) AppendJavaScript(key, js string) *DocumentBuilder {
	db.doc = AppendJavaScriptElement(db.doc, key, js)
	return db
}


func (db *DocumentBuilder) AppendSymbol(key, symbol string) *DocumentBuilder {
	db.doc = AppendSymbolElement(db.doc, key, symbol)
	return db
}


func (db *DocumentBuilder) AppendCodeWithScope(key string, code string, scope Document) *DocumentBuilder {
	db.doc = AppendCodeWithScopeElement(db.doc, key, code, scope)
	return db
}


func (db *DocumentBuilder) AppendTimestamp(key string, t, i uint32) *DocumentBuilder {
	db.doc = AppendTimestampElement(db.doc, key, t, i)
	return db
}


func (db *DocumentBuilder) AppendInt64(key string, i64 int64) *DocumentBuilder {
	db.doc = AppendInt64Element(db.doc, key, i64)
	return db
}


func (db *DocumentBuilder) AppendDecimal128(key string, d128 primitive.Decimal128) *DocumentBuilder {
	db.doc = AppendDecimal128Element(db.doc, key, d128)
	return db
}


func (db *DocumentBuilder) AppendMaxKey(key string) *DocumentBuilder {
	db.doc = AppendMaxKeyElement(db.doc, key)
	return db
}


func (db *DocumentBuilder) AppendMinKey(key string) *DocumentBuilder {
	db.doc = AppendMinKeyElement(db.doc, key)
	return db
}


func (db *DocumentBuilder) AppendValue(key string, val Value) *DocumentBuilder {
	db.doc = AppendValueElement(db.doc, key, val)
	return db
}



func (db *DocumentBuilder) StartDocument(key string) *DocumentBuilder {
	db.doc = AppendHeader(db.doc, bsontype.EmbeddedDocument, key)
	db = db.startDocument()
	return db
}


func (db *DocumentBuilder) FinishDocument() *DocumentBuilder {
	db.doc = db.Build()
	return db
}
