





package bsoncore

import (
	"strconv"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type ArrayBuilder struct {
	arr     []byte
	indexes []int32
	keys    []int
}


func NewArrayBuilder() *ArrayBuilder {
	return (&ArrayBuilder{}).startArray()
}


func (a *ArrayBuilder) startArray() *ArrayBuilder {
	var index int32
	index, a.arr = AppendArrayStart(a.arr)
	a.indexes = append(a.indexes, index)
	a.keys = append(a.keys, 0)
	return a
}



func (a *ArrayBuilder) Build() Array {
	lastIndex := len(a.indexes) - 1
	lastKey := len(a.keys) - 1
	a.arr, _ = AppendArrayEnd(a.arr, a.indexes[lastIndex])
	a.indexes = a.indexes[:lastIndex]
	a.keys = a.keys[:lastKey]
	return a.arr
}


func (a *ArrayBuilder) incrementKey() string {
	idx := len(a.keys) - 1
	key := strconv.Itoa(a.keys[idx])
	a.keys[idx]++
	return key
}


func (a *ArrayBuilder) AppendInt32(i32 int32) *ArrayBuilder {
	a.arr = AppendInt32Element(a.arr, a.incrementKey(), i32)
	return a
}


func (a *ArrayBuilder) AppendDocument(doc []byte) *ArrayBuilder {
	a.arr = AppendDocumentElement(a.arr, a.incrementKey(), doc)
	return a
}


func (a *ArrayBuilder) AppendArray(arr []byte) *ArrayBuilder {
	a.arr = AppendArrayElement(a.arr, a.incrementKey(), arr)
	return a
}


func (a *ArrayBuilder) AppendDouble(f float64) *ArrayBuilder {
	a.arr = AppendDoubleElement(a.arr, a.incrementKey(), f)
	return a
}


func (a *ArrayBuilder) AppendString(str string) *ArrayBuilder {
	a.arr = AppendStringElement(a.arr, a.incrementKey(), str)
	return a
}


func (a *ArrayBuilder) AppendObjectID(oid primitive.ObjectID) *ArrayBuilder {
	a.arr = AppendObjectIDElement(a.arr, a.incrementKey(), oid)
	return a
}



func (a *ArrayBuilder) AppendBinary(subtype byte, b []byte) *ArrayBuilder {
	a.arr = AppendBinaryElement(a.arr, a.incrementKey(), subtype, b)
	return a
}


func (a *ArrayBuilder) AppendUndefined() *ArrayBuilder {
	a.arr = AppendUndefinedElement(a.arr, a.incrementKey())
	return a
}


func (a *ArrayBuilder) AppendBoolean(b bool) *ArrayBuilder {
	a.arr = AppendBooleanElement(a.arr, a.incrementKey(), b)
	return a
}


func (a *ArrayBuilder) AppendDateTime(dt int64) *ArrayBuilder {
	a.arr = AppendDateTimeElement(a.arr, a.incrementKey(), dt)
	return a
}


func (a *ArrayBuilder) AppendNull() *ArrayBuilder {
	a.arr = AppendNullElement(a.arr, a.incrementKey())
	return a
}


func (a *ArrayBuilder) AppendRegex(pattern, options string) *ArrayBuilder {
	a.arr = AppendRegexElement(a.arr, a.incrementKey(), pattern, options)
	return a
}


func (a *ArrayBuilder) AppendDBPointer(ns string, oid primitive.ObjectID) *ArrayBuilder {
	a.arr = AppendDBPointerElement(a.arr, a.incrementKey(), ns, oid)
	return a
}


func (a *ArrayBuilder) AppendJavaScript(js string) *ArrayBuilder {
	a.arr = AppendJavaScriptElement(a.arr, a.incrementKey(), js)
	return a
}


func (a *ArrayBuilder) AppendSymbol(symbol string) *ArrayBuilder {
	a.arr = AppendSymbolElement(a.arr, a.incrementKey(), symbol)
	return a
}


func (a *ArrayBuilder) AppendCodeWithScope(code string, scope Document) *ArrayBuilder {
	a.arr = AppendCodeWithScopeElement(a.arr, a.incrementKey(), code, scope)
	return a
}


func (a *ArrayBuilder) AppendTimestamp(t, i uint32) *ArrayBuilder {
	a.arr = AppendTimestampElement(a.arr, a.incrementKey(), t, i)
	return a
}


func (a *ArrayBuilder) AppendInt64(i64 int64) *ArrayBuilder {
	a.arr = AppendInt64Element(a.arr, a.incrementKey(), i64)
	return a
}


func (a *ArrayBuilder) AppendDecimal128(d128 primitive.Decimal128) *ArrayBuilder {
	a.arr = AppendDecimal128Element(a.arr, a.incrementKey(), d128)
	return a
}


func (a *ArrayBuilder) AppendMaxKey() *ArrayBuilder {
	a.arr = AppendMaxKeyElement(a.arr, a.incrementKey())
	return a
}


func (a *ArrayBuilder) AppendMinKey() *ArrayBuilder {
	a.arr = AppendMinKeyElement(a.arr, a.incrementKey())
	return a
}


func (a *ArrayBuilder) AppendValue(val Value) *ArrayBuilder {
	a.arr = AppendValueElement(a.arr, a.incrementKey(), val)
	return a
}



func (a *ArrayBuilder) StartArray() *ArrayBuilder {
	a.arr = AppendHeader(a.arr, bsontype.Array, a.incrementKey())
	a.startArray()
	return a
}


func (a *ArrayBuilder) FinishArray() *ArrayBuilder {
	a.arr = a.Build()
	return a
}
