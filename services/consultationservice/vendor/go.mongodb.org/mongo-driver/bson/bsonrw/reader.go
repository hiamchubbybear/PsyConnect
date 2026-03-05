





package bsonrw

import (
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)



type ArrayReader interface {
	ReadValue() (ValueReader, error)
}



type DocumentReader interface {
	ReadElement() (string, ValueReader, error)
}




type ValueReader interface {
	Type() bsontype.Type
	Skip() error

	ReadArray() (ArrayReader, error)
	ReadBinary() (b []byte, btype byte, err error)
	ReadBoolean() (bool, error)
	ReadDocument() (DocumentReader, error)
	ReadCodeWithScope() (code string, dr DocumentReader, err error)
	ReadDBPointer() (ns string, oid primitive.ObjectID, err error)
	ReadDateTime() (int64, error)
	ReadDecimal128() (primitive.Decimal128, error)
	ReadDouble() (float64, error)
	ReadInt32() (int32, error)
	ReadInt64() (int64, error)
	ReadJavascript() (code string, err error)
	ReadMaxKey() error
	ReadMinKey() error
	ReadNull() error
	ReadObjectID() (primitive.ObjectID, error)
	ReadRegex() (pattern, options string, err error)
	ReadString() (string, error)
	ReadSymbol() (symbol string, err error)
	ReadTimestamp() (t, i uint32, err error)
	ReadUndefined() error
}








type BytesReader interface {
	ReadValueBytes(dst []byte) (bsontype.Type, []byte, error)
}
