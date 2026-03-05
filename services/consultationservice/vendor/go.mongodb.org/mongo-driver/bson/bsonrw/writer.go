





package bsonrw

import (
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)




type ArrayWriter interface {
	WriteArrayElement() (ValueWriter, error)
	WriteArrayEnd() error
}




type DocumentWriter interface {
	WriteDocumentElement(string) (ValueWriter, error)
	WriteDocumentEnd() error
}




type ValueWriter interface {
	WriteArray() (ArrayWriter, error)
	WriteBinary(b []byte) error
	WriteBinaryWithSubtype(b []byte, btype byte) error
	WriteBoolean(bool) error
	WriteCodeWithScope(code string) (DocumentWriter, error)
	WriteDBPointer(ns string, oid primitive.ObjectID) error
	WriteDateTime(dt int64) error
	WriteDecimal128(primitive.Decimal128) error
	WriteDouble(float64) error
	WriteInt32(int32) error
	WriteInt64(int64) error
	WriteJavascript(code string) error
	WriteMaxKey() error
	WriteMinKey() error
	WriteNull() error
	WriteObjectID(primitive.ObjectID) error
	WriteRegex(pattern, options string) error
	WriteString(string) error
	WriteDocument() (DocumentWriter, error)
	WriteSymbol(symbol string) error
	WriteTimestamp(t, i uint32) error
	WriteUndefined() error
}




type ValueWriterFlusher interface {
	ValueWriter
	Flush() error
}






type BytesWriter interface {
	WriteValueBytes(t bsontype.Type, b []byte) error
}




type SliceWriter []byte




func (sw *SliceWriter) Write(p []byte) (int, error) {
	written := len(p)
	*sw = append(*sw, p...)
	return written, nil
}
