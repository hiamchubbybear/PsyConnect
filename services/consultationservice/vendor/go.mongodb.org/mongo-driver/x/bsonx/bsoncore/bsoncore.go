





package bsoncore 

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	
	EmptyDocumentLength = 5
	
	nullTerminator       = string(byte(0))
	invalidKeyPanicMsg   = "BSON element keys cannot contain null bytes"
	invalidRegexPanicMsg = "BSON regex values cannot contain null bytes"
)


func AppendType(dst []byte, t bsontype.Type) []byte { return append(dst, byte(t)) }


func AppendKey(dst []byte, key string) []byte { return append(dst, key+nullTerminator...) }



func AppendHeader(dst []byte, t bsontype.Type, key string) []byte {
	if !isValidCString(key) {
		panic(invalidKeyPanicMsg)
	}

	dst = AppendType(dst, t)
	dst = append(dst, key...)
	return append(dst, 0x00)
	
}





func ReadType(src []byte) (bsontype.Type, []byte, bool) {
	if len(src) < 1 {
		return 0, src, false
	}
	return bsontype.Type(src[0]), src[1:], true
}




func ReadKey(src []byte) (string, []byte, bool) { return readcstring(src) }




func ReadKeyBytes(src []byte) ([]byte, []byte, bool) { return readcstringbytes(src) }



func ReadHeader(src []byte) (t bsontype.Type, key string, rem []byte, ok bool) {
	t, rem, ok = ReadType(src)
	if !ok {
		return 0, "", src, false
	}
	key, rem, ok = ReadKey(rem)
	if !ok {
		return 0, "", src, false
	}

	return t, key, rem, true
}



func ReadHeaderBytes(src []byte) (header []byte, rem []byte, ok bool) {
	if len(src) < 1 {
		return nil, src, false
	}
	idx := bytes.IndexByte(src[1:], 0x00)
	if idx == -1 {
		return nil, src, false
	}
	return src[:idx], src[idx+1:], true
}



func ReadElement(src []byte) (Element, []byte, bool) {
	if len(src) < 1 {
		return nil, src, false
	}
	t := bsontype.Type(src[0])
	idx := bytes.IndexByte(src[1:], 0x00)
	if idx == -1 {
		return nil, src, false
	}
	length, ok := valueLength(src[idx+2:], t) 
	if !ok {
		return nil, src, false
	}
	elemLength := 1 + idx + 1 + int(length)
	if elemLength > len(src) {
		return nil, src, false
	}
	if elemLength < 0 {
		return nil, src, false
	}
	return src[:elemLength], src[elemLength:], true
}


func AppendValueElement(dst []byte, key string, value Value) []byte {
	dst = AppendHeader(dst, value.Type, key)
	dst = append(dst, value.Data...)
	return dst
}



func ReadValue(src []byte, t bsontype.Type) (Value, []byte, bool) {
	data, rem, ok := readValue(src, t)
	if !ok {
		return Value{}, src, false
	}
	return Value{Type: t, Data: data}, rem, true
}


func AppendDouble(dst []byte, f float64) []byte {
	return appendu64(dst, math.Float64bits(f))
}



func AppendDoubleElement(dst []byte, key string, f float64) []byte {
	return AppendDouble(AppendHeader(dst, bsontype.Double, key), f)
}



func ReadDouble(src []byte) (float64, []byte, bool) {
	bits, src, ok := readu64(src)
	if !ok {
		return 0, src, false
	}
	return math.Float64frombits(bits), src, true
}


func AppendString(dst []byte, s string) []byte {
	return appendstring(dst, s)
}



func AppendStringElement(dst []byte, key, val string) []byte {
	return AppendString(AppendHeader(dst, bsontype.String, key), val)
}



func ReadString(src []byte) (string, []byte, bool) {
	return readstring(src)
}



func AppendDocumentStart(dst []byte) (index int32, b []byte) {
	
	
	
	
	return ReserveLength(dst)
}



func AppendDocumentStartInline(dst []byte, index *int32) []byte {
	idx, doc := AppendDocumentStart(dst)
	*index = idx
	return doc
}


func AppendDocumentElementStart(dst []byte, key string) (index int32, b []byte) {
	return AppendDocumentStart(AppendHeader(dst, bsontype.EmbeddedDocument, key))
}



func AppendDocumentEnd(dst []byte, index int32) ([]byte, error) {
	if int(index) > len(dst)-4 {
		return dst, fmt.Errorf("not enough bytes available after index to write length")
	}
	dst = append(dst, 0x00)
	dst = UpdateLength(dst, index, int32(len(dst[index:])))
	return dst, nil
}


func AppendDocument(dst []byte, doc []byte) []byte { return append(dst, doc...) }



func AppendDocumentElement(dst []byte, key string, doc []byte) []byte {
	return AppendDocument(AppendHeader(dst, bsontype.EmbeddedDocument, key), doc)
}



func BuildDocument(dst []byte, elems ...[]byte) []byte {
	idx, dst := ReserveLength(dst)
	for _, elem := range elems {
		dst = append(dst, elem...)
	}
	dst = append(dst, 0x00)
	dst = UpdateLength(dst, idx, int32(len(dst[idx:])))
	return dst
}


func BuildDocumentValue(elems ...[]byte) Value {
	return Value{Type: bsontype.EmbeddedDocument, Data: BuildDocument(nil, elems...)}
}



func BuildDocumentElement(dst []byte, key string, elems ...[]byte) []byte {
	return BuildDocument(AppendHeader(dst, bsontype.EmbeddedDocument, key), elems...)
}


var BuildDocumentFromElements = BuildDocument



func ReadDocument(src []byte) (doc Document, rem []byte, ok bool) { return readLengthBytes(src) }



func AppendArrayStart(dst []byte) (index int32, b []byte) { return ReserveLength(dst) }



func AppendArrayElementStart(dst []byte, key string) (index int32, b []byte) {
	return AppendArrayStart(AppendHeader(dst, bsontype.Array, key))
}



func AppendArrayEnd(dst []byte, index int32) ([]byte, error) { return AppendDocumentEnd(dst, index) }


func AppendArray(dst []byte, arr []byte) []byte { return append(dst, arr...) }



func AppendArrayElement(dst []byte, key string, arr []byte) []byte {
	return AppendArray(AppendHeader(dst, bsontype.Array, key), arr)
}


func BuildArray(dst []byte, values ...Value) []byte {
	idx, dst := ReserveLength(dst)
	for pos, val := range values {
		dst = AppendValueElement(dst, strconv.Itoa(pos), val)
	}
	dst = append(dst, 0x00)
	dst = UpdateLength(dst, idx, int32(len(dst[idx:])))
	return dst
}


func BuildArrayElement(dst []byte, key string, values ...Value) []byte {
	return BuildArray(AppendHeader(dst, bsontype.Array, key), values...)
}



func ReadArray(src []byte) (arr Array, rem []byte, ok bool) { return readLengthBytes(src) }


func AppendBinary(dst []byte, subtype byte, b []byte) []byte {
	if subtype == 0x02 {
		return appendBinarySubtype2(dst, subtype, b)
	}
	dst = append(appendLength(dst, int32(len(b))), subtype)
	return append(dst, b...)
}



func AppendBinaryElement(dst []byte, key string, subtype byte, b []byte) []byte {
	return AppendBinary(AppendHeader(dst, bsontype.Binary, key), subtype, b)
}



func ReadBinary(src []byte) (subtype byte, bin []byte, rem []byte, ok bool) {
	length, rem, ok := ReadLength(src)
	if !ok {
		return 0x00, nil, src, false
	}
	if len(rem) < 1 { 
		return 0x00, nil, src, false
	}
	subtype, rem = rem[0], rem[1:]

	if len(rem) < int(length) {
		return 0x00, nil, src, false
	}

	if subtype == 0x02 {
		length, rem, ok = ReadLength(rem)
		if !ok || len(rem) < int(length) {
			return 0x00, nil, src, false
		}
	}

	return subtype, rem[:length], rem[length:], true
}



func AppendUndefinedElement(dst []byte, key string) []byte {
	return AppendHeader(dst, bsontype.Undefined, key)
}


func AppendObjectID(dst []byte, oid primitive.ObjectID) []byte { return append(dst, oid[:]...) }



func AppendObjectIDElement(dst []byte, key string, oid primitive.ObjectID) []byte {
	return AppendObjectID(AppendHeader(dst, bsontype.ObjectID, key), oid)
}



func ReadObjectID(src []byte) (primitive.ObjectID, []byte, bool) {
	if len(src) < 12 {
		return primitive.ObjectID{}, src, false
	}
	var oid primitive.ObjectID
	copy(oid[:], src[0:12])
	return oid, src[12:], true
}


func AppendBoolean(dst []byte, b bool) []byte {
	if b {
		return append(dst, 0x01)
	}
	return append(dst, 0x00)
}



func AppendBooleanElement(dst []byte, key string, b bool) []byte {
	return AppendBoolean(AppendHeader(dst, bsontype.Boolean, key), b)
}



func ReadBoolean(src []byte) (bool, []byte, bool) {
	if len(src) < 1 {
		return false, src, false
	}

	return src[0] == 0x01, src[1:], true
}


func AppendDateTime(dst []byte, dt int64) []byte { return appendi64(dst, dt) }



func AppendDateTimeElement(dst []byte, key string, dt int64) []byte {
	return AppendDateTime(AppendHeader(dst, bsontype.DateTime, key), dt)
}



func ReadDateTime(src []byte) (int64, []byte, bool) { return readi64(src) }


func AppendTime(dst []byte, t time.Time) []byte {
	return AppendDateTime(dst, t.Unix()*1000+int64(t.Nanosecond()/1e6))
}



func AppendTimeElement(dst []byte, key string, t time.Time) []byte {
	return AppendTime(AppendHeader(dst, bsontype.DateTime, key), t)
}



func ReadTime(src []byte) (time.Time, []byte, bool) {
	dt, rem, ok := readi64(src)
	return time.Unix(dt/1e3, dt%1e3*1e6), rem, ok
}



func AppendNullElement(dst []byte, key string) []byte { return AppendHeader(dst, bsontype.Null, key) }


func AppendRegex(dst []byte, pattern, options string) []byte {
	if !isValidCString(pattern) || !isValidCString(options) {
		panic(invalidRegexPanicMsg)
	}

	return append(dst, pattern+nullTerminator+options+nullTerminator...)
}



func AppendRegexElement(dst []byte, key, pattern, options string) []byte {
	return AppendRegex(AppendHeader(dst, bsontype.Regex, key), pattern, options)
}



func ReadRegex(src []byte) (pattern, options string, rem []byte, ok bool) {
	pattern, rem, ok = readcstring(src)
	if !ok {
		return "", "", src, false
	}
	options, rem, ok = readcstring(rem)
	if !ok {
		return "", "", src, false
	}
	return pattern, options, rem, true
}


func AppendDBPointer(dst []byte, ns string, oid primitive.ObjectID) []byte {
	return append(appendstring(dst, ns), oid[:]...)
}



func AppendDBPointerElement(dst []byte, key, ns string, oid primitive.ObjectID) []byte {
	return AppendDBPointer(AppendHeader(dst, bsontype.DBPointer, key), ns, oid)
}



func ReadDBPointer(src []byte) (ns string, oid primitive.ObjectID, rem []byte, ok bool) {
	ns, rem, ok = readstring(src)
	if !ok {
		return "", primitive.ObjectID{}, src, false
	}
	oid, rem, ok = ReadObjectID(rem)
	if !ok {
		return "", primitive.ObjectID{}, src, false
	}
	return ns, oid, rem, true
}


func AppendJavaScript(dst []byte, js string) []byte { return appendstring(dst, js) }



func AppendJavaScriptElement(dst []byte, key, js string) []byte {
	return AppendJavaScript(AppendHeader(dst, bsontype.JavaScript, key), js)
}



func ReadJavaScript(src []byte) (js string, rem []byte, ok bool) { return readstring(src) }


func AppendSymbol(dst []byte, symbol string) []byte { return appendstring(dst, symbol) }



func AppendSymbolElement(dst []byte, key, symbol string) []byte {
	return AppendSymbol(AppendHeader(dst, bsontype.Symbol, key), symbol)
}



func ReadSymbol(src []byte) (symbol string, rem []byte, ok bool) { return readstring(src) }


func AppendCodeWithScope(dst []byte, code string, scope []byte) []byte {
	length := int32(4 + 4 + len(code) + 1 + len(scope)) 
	dst = appendLength(dst, length)

	return append(appendstring(dst, code), scope...)
}




func AppendCodeWithScopeElement(dst []byte, key, code string, scope []byte) []byte {
	return AppendCodeWithScope(AppendHeader(dst, bsontype.CodeWithScope, key), code, scope)
}



func ReadCodeWithScope(src []byte) (code string, scope []byte, rem []byte, ok bool) {
	length, rem, ok := ReadLength(src)
	if !ok || len(src) < int(length) {
		return "", nil, src, false
	}

	code, rem, ok = readstring(rem)
	if !ok {
		return "", nil, src, false
	}

	scope, rem, ok = ReadDocument(rem)
	if !ok {
		return "", nil, src, false
	}
	return code, scope, rem, true
}


func AppendInt32(dst []byte, i32 int32) []byte { return appendi32(dst, i32) }



func AppendInt32Element(dst []byte, key string, i32 int32) []byte {
	return AppendInt32(AppendHeader(dst, bsontype.Int32, key), i32)
}



func ReadInt32(src []byte) (int32, []byte, bool) { return readi32(src) }


func AppendTimestamp(dst []byte, t, i uint32) []byte {
	return appendu32(appendu32(dst, i), t) 
}



func AppendTimestampElement(dst []byte, key string, t, i uint32) []byte {
	return AppendTimestamp(AppendHeader(dst, bsontype.Timestamp, key), t, i)
}



func ReadTimestamp(src []byte) (t, i uint32, rem []byte, ok bool) {
	i, rem, ok = readu32(src)
	if !ok {
		return 0, 0, src, false
	}
	t, rem, ok = readu32(rem)
	if !ok {
		return 0, 0, src, false
	}
	return t, i, rem, true
}


func AppendInt64(dst []byte, i64 int64) []byte { return appendi64(dst, i64) }



func AppendInt64Element(dst []byte, key string, i64 int64) []byte {
	return AppendInt64(AppendHeader(dst, bsontype.Int64, key), i64)
}



func ReadInt64(src []byte) (int64, []byte, bool) { return readi64(src) }


func AppendDecimal128(dst []byte, d128 primitive.Decimal128) []byte {
	high, low := d128.GetBytes()
	return appendu64(appendu64(dst, low), high)
}



func AppendDecimal128Element(dst []byte, key string, d128 primitive.Decimal128) []byte {
	return AppendDecimal128(AppendHeader(dst, bsontype.Decimal128, key), d128)
}



func ReadDecimal128(src []byte) (primitive.Decimal128, []byte, bool) {
	l, rem, ok := readu64(src)
	if !ok {
		return primitive.Decimal128{}, src, false
	}

	h, rem, ok := readu64(rem)
	if !ok {
		return primitive.Decimal128{}, src, false
	}

	return primitive.NewDecimal128(h, l), rem, true
}



func AppendMaxKeyElement(dst []byte, key string) []byte {
	return AppendHeader(dst, bsontype.MaxKey, key)
}



func AppendMinKeyElement(dst []byte, key string) []byte {
	return AppendHeader(dst, bsontype.MinKey, key)
}


func EqualValue(t1, t2 bsontype.Type, v1, v2 []byte) bool {
	if t1 != t2 {
		return false
	}
	v1, _, ok := readValue(v1, t1)
	if !ok {
		return false
	}
	v2, _, ok = readValue(v2, t2)
	if !ok {
		return false
	}
	return bytes.Equal(v1, v2)
}




func valueLength(src []byte, t bsontype.Type) (int32, bool) {
	var length int32
	ok := true
	switch t {
	case bsontype.Array, bsontype.EmbeddedDocument, bsontype.CodeWithScope:
		length, _, ok = ReadLength(src)
	case bsontype.Binary:
		length, _, ok = ReadLength(src)
		length += 4 + 1 
	case bsontype.Boolean:
		length = 1
	case bsontype.DBPointer:
		length, _, ok = ReadLength(src)
		length += 4 + 12 
	case bsontype.DateTime, bsontype.Double, bsontype.Int64, bsontype.Timestamp:
		length = 8
	case bsontype.Decimal128:
		length = 16
	case bsontype.Int32:
		length = 4
	case bsontype.JavaScript, bsontype.String, bsontype.Symbol:
		length, _, ok = ReadLength(src)
		length += 4
	case bsontype.MaxKey, bsontype.MinKey, bsontype.Null, bsontype.Undefined:
		length = 0
	case bsontype.ObjectID:
		length = 12
	case bsontype.Regex:
		regex := bytes.IndexByte(src, 0x00)
		if regex < 0 {
			ok = false
			break
		}
		pattern := bytes.IndexByte(src[regex+1:], 0x00)
		if pattern < 0 {
			ok = false
			break
		}
		length = int32(int64(regex) + 1 + int64(pattern) + 1)
	default:
		ok = false
	}

	return length, ok
}

func readValue(src []byte, t bsontype.Type) ([]byte, []byte, bool) {
	length, ok := valueLength(src, t)
	if !ok || int(length) > len(src) {
		return nil, src, false
	}

	return src[:length], src[length:], true
}



func ReserveLength(dst []byte) (int32, []byte) {
	index := len(dst)
	return int32(index), append(dst, 0x00, 0x00, 0x00, 0x00)
}


func UpdateLength(dst []byte, index, length int32) []byte {
	binary.LittleEndian.PutUint32(dst[index:], uint32(length))
	return dst
}

func appendLength(dst []byte, l int32) []byte { return appendi32(dst, l) }

func appendi32(dst []byte, i32 int32) []byte {
	b := []byte{0, 0, 0, 0}
	binary.LittleEndian.PutUint32(b, uint32(i32))
	return append(dst, b...)
}




func ReadLength(src []byte) (int32, []byte, bool) {
	ln, src, ok := readi32(src)
	if ln < 0 {
		return ln, src, false
	}
	return ln, src, ok
}

func readi32(src []byte) (int32, []byte, bool) {
	if len(src) < 4 {
		return 0, src, false
	}
	return int32(binary.LittleEndian.Uint32(src)), src[4:], true
}

func appendi64(dst []byte, i64 int64) []byte {
	b := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	binary.LittleEndian.PutUint64(b, uint64(i64))
	return append(dst, b...)
}

func readi64(src []byte) (int64, []byte, bool) {
	if len(src) < 8 {
		return 0, src, false
	}
	return int64(binary.LittleEndian.Uint64(src)), src[8:], true
}

func appendu32(dst []byte, u32 uint32) []byte {
	b := []byte{0, 0, 0, 0}
	binary.LittleEndian.PutUint32(b, u32)
	return append(dst, b...)
}

func readu32(src []byte) (uint32, []byte, bool) {
	if len(src) < 4 {
		return 0, src, false
	}

	return binary.LittleEndian.Uint32(src), src[4:], true
}

func appendu64(dst []byte, u64 uint64) []byte {
	b := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	binary.LittleEndian.PutUint64(b, u64)
	return append(dst, b...)
}

func readu64(src []byte) (uint64, []byte, bool) {
	if len(src) < 8 {
		return 0, src, false
	}
	return binary.LittleEndian.Uint64(src), src[8:], true
}


func readcstring(src []byte) (string, []byte, bool) {
	idx := bytes.IndexByte(src, 0x00)
	if idx < 0 {
		return "", src, false
	}
	return string(src[:idx]), src[idx+1:], true
}


func readcstringbytes(src []byte) ([]byte, []byte, bool) {
	idx := bytes.IndexByte(src, 0x00)
	if idx < 0 {
		return nil, src, false
	}
	return src[:idx], src[idx+1:], true
}

func appendstring(dst []byte, s string) []byte {
	l := int32(len(s) + 1)
	dst = appendLength(dst, l)
	dst = append(dst, s...)
	return append(dst, 0x00)
}

func readstring(src []byte) (string, []byte, bool) {
	l, rem, ok := ReadLength(src)
	if !ok {
		return "", src, false
	}
	if len(src[4:]) < int(l) || l == 0 {
		return "", src, false
	}

	return string(rem[:l-1]), rem[l:], true
}



func readLengthBytes(src []byte) ([]byte, []byte, bool) {
	l, _, ok := ReadLength(src)
	if !ok {
		return nil, src, false
	}
	if l < 4 {
		return nil, src, false
	}
	if len(src) < int(l) {
		return nil, src, false
	}
	return src[:l], src[l:], true
}

func appendBinarySubtype2(dst []byte, subtype byte, b []byte) []byte {
	dst = appendLength(dst, int32(len(b)+4)) 
	dst = append(dst, subtype)
	dst = appendLength(dst, int32(len(b)))
	return append(dst, b...)
}

func isValidCString(cs string) bool {
	return !strings.ContainsRune(cs, '\x00')
}
