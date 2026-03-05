





package bsoncore

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type ElementTypeError struct {
	Method string
	Type   bsontype.Type
}


func (ete ElementTypeError) Error() string {
	return "Call of " + ete.Method + " on " + ete.Type.String() + " type"
}


type Value struct {
	Type bsontype.Type
	Data []byte
}


func (v Value) Validate() error {
	_, _, valid := readValue(v.Data, v.Type)
	if !valid {
		return NewInsufficientBytesError(v.Data, v.Data)
	}
	return nil
}


func (v Value) IsNumber() bool {
	switch v.Type {
	case bsontype.Double, bsontype.Int32, bsontype.Int64, bsontype.Decimal128:
		return true
	default:
		return false
	}
}



func (v Value) AsInt32() int32 {
	if !v.IsNumber() {
		panic(ElementTypeError{"bsoncore.Value.AsInt32", v.Type})
	}
	var i32 int32
	switch v.Type {
	case bsontype.Double:
		f64, _, ok := ReadDouble(v.Data)
		if !ok {
			panic(NewInsufficientBytesError(v.Data, v.Data))
		}
		i32 = int32(f64)
	case bsontype.Int32:
		var ok bool
		i32, _, ok = ReadInt32(v.Data)
		if !ok {
			panic(NewInsufficientBytesError(v.Data, v.Data))
		}
	case bsontype.Int64:
		i64, _, ok := ReadInt64(v.Data)
		if !ok {
			panic(NewInsufficientBytesError(v.Data, v.Data))
		}
		i32 = int32(i64)
	case bsontype.Decimal128:
		panic(ElementTypeError{"bsoncore.Value.AsInt32", v.Type})
	}
	return i32
}



func (v Value) AsInt32OK() (int32, bool) {
	if !v.IsNumber() {
		return 0, false
	}
	var i32 int32
	switch v.Type {
	case bsontype.Double:
		f64, _, ok := ReadDouble(v.Data)
		if !ok {
			return 0, false
		}
		i32 = int32(f64)
	case bsontype.Int32:
		var ok bool
		i32, _, ok = ReadInt32(v.Data)
		if !ok {
			return 0, false
		}
	case bsontype.Int64:
		i64, _, ok := ReadInt64(v.Data)
		if !ok {
			return 0, false
		}
		i32 = int32(i64)
	case bsontype.Decimal128:
		return 0, false
	}
	return i32, true
}



func (v Value) AsInt64() int64 {
	if !v.IsNumber() {
		panic(ElementTypeError{"bsoncore.Value.AsInt64", v.Type})
	}
	var i64 int64
	switch v.Type {
	case bsontype.Double:
		f64, _, ok := ReadDouble(v.Data)
		if !ok {
			panic(NewInsufficientBytesError(v.Data, v.Data))
		}
		i64 = int64(f64)
	case bsontype.Int32:
		var ok bool
		i32, _, ok := ReadInt32(v.Data)
		if !ok {
			panic(NewInsufficientBytesError(v.Data, v.Data))
		}
		i64 = int64(i32)
	case bsontype.Int64:
		var ok bool
		i64, _, ok = ReadInt64(v.Data)
		if !ok {
			panic(NewInsufficientBytesError(v.Data, v.Data))
		}
	case bsontype.Decimal128:
		panic(ElementTypeError{"bsoncore.Value.AsInt64", v.Type})
	}
	return i64
}



func (v Value) AsInt64OK() (int64, bool) {
	if !v.IsNumber() {
		return 0, false
	}
	var i64 int64
	switch v.Type {
	case bsontype.Double:
		f64, _, ok := ReadDouble(v.Data)
		if !ok {
			return 0, false
		}
		i64 = int64(f64)
	case bsontype.Int32:
		var ok bool
		i32, _, ok := ReadInt32(v.Data)
		if !ok {
			return 0, false
		}
		i64 = int64(i32)
	case bsontype.Int64:
		var ok bool
		i64, _, ok = ReadInt64(v.Data)
		if !ok {
			return 0, false
		}
	case bsontype.Decimal128:
		return 0, false
	}
	return i64, true
}














func (v Value) Equal(v2 Value) bool {
	if v.Type != v2.Type {
		return false
	}

	return bytes.Equal(v.Data, v2.Data)
}



func (v Value) String() string {
	switch v.Type {
	case bsontype.Double:
		f64, ok := v.DoubleOK()
		if !ok {
			return ""
		}
		return fmt.Sprintf(`{"$numberDouble":"%s"}`, formatDouble(f64))
	case bsontype.String:
		str, ok := v.StringValueOK()
		if !ok {
			return ""
		}
		return escapeString(str)
	case bsontype.EmbeddedDocument:
		doc, ok := v.DocumentOK()
		if !ok {
			return ""
		}
		return doc.String()
	case bsontype.Array:
		arr, ok := v.ArrayOK()
		if !ok {
			return ""
		}
		return arr.String()
	case bsontype.Binary:
		subtype, data, ok := v.BinaryOK()
		if !ok {
			return ""
		}
		return fmt.Sprintf(`{"$binary":{"base64":"%s","subType":"%02x"}}`, base64.StdEncoding.EncodeToString(data), subtype)
	case bsontype.Undefined:
		return `{"$undefined":true}`
	case bsontype.ObjectID:
		oid, ok := v.ObjectIDOK()
		if !ok {
			return ""
		}
		return fmt.Sprintf(`{"$oid":"%s"}`, oid.Hex())
	case bsontype.Boolean:
		b, ok := v.BooleanOK()
		if !ok {
			return ""
		}
		return strconv.FormatBool(b)
	case bsontype.DateTime:
		dt, ok := v.DateTimeOK()
		if !ok {
			return ""
		}
		return fmt.Sprintf(`{"$date":{"$numberLong":"%d"}}`, dt)
	case bsontype.Null:
		return "null"
	case bsontype.Regex:
		pattern, options, ok := v.RegexOK()
		if !ok {
			return ""
		}
		return fmt.Sprintf(
			`{"$regularExpression":{"pattern":%s,"options":"%s"}}`,
			escapeString(pattern), sortStringAlphebeticAscending(options),
		)
	case bsontype.DBPointer:
		ns, pointer, ok := v.DBPointerOK()
		if !ok {
			return ""
		}
		return fmt.Sprintf(`{"$dbPointer":{"$ref":%s,"$id":{"$oid":"%s"}}}`, escapeString(ns), pointer.Hex())
	case bsontype.JavaScript:
		js, ok := v.JavaScriptOK()
		if !ok {
			return ""
		}
		return fmt.Sprintf(`{"$code":%s}`, escapeString(js))
	case bsontype.Symbol:
		symbol, ok := v.SymbolOK()
		if !ok {
			return ""
		}
		return fmt.Sprintf(`{"$symbol":%s}`, escapeString(symbol))
	case bsontype.CodeWithScope:
		code, scope, ok := v.CodeWithScopeOK()
		if !ok {
			return ""
		}
		return fmt.Sprintf(`{"$code":%s,"$scope":%s}`, code, scope)
	case bsontype.Int32:
		i32, ok := v.Int32OK()
		if !ok {
			return ""
		}
		return fmt.Sprintf(`{"$numberInt":"%d"}`, i32)
	case bsontype.Timestamp:
		t, i, ok := v.TimestampOK()
		if !ok {
			return ""
		}
		return fmt.Sprintf(`{"$timestamp":{"t":%v,"i":%v}}`, t, i)
	case bsontype.Int64:
		i64, ok := v.Int64OK()
		if !ok {
			return ""
		}
		return fmt.Sprintf(`{"$numberLong":"%d"}`, i64)
	case bsontype.Decimal128:
		d128, ok := v.Decimal128OK()
		if !ok {
			return ""
		}
		return fmt.Sprintf(`{"$numberDecimal":"%s"}`, d128.String())
	case bsontype.MinKey:
		return `{"$minKey":1}`
	case bsontype.MaxKey:
		return `{"$maxKey":1}`
	default:
		return ""
	}
}



func (v Value) DebugString() string {
	switch v.Type {
	case bsontype.String:
		str, ok := v.StringValueOK()
		if !ok {
			return "<malformed>"
		}
		return escapeString(str)
	case bsontype.EmbeddedDocument:
		doc, ok := v.DocumentOK()
		if !ok {
			return "<malformed>"
		}
		return doc.DebugString()
	case bsontype.Array:
		arr, ok := v.ArrayOK()
		if !ok {
			return "<malformed>"
		}
		return arr.DebugString()
	case bsontype.CodeWithScope:
		code, scope, ok := v.CodeWithScopeOK()
		if !ok {
			return ""
		}
		return fmt.Sprintf(`{"$code":%s,"$scope":%s}`, code, scope.DebugString())
	default:
		str := v.String()
		if str == "" {
			return "<malformed>"
		}
		return str
	}
}



func (v Value) Double() float64 {
	if v.Type != bsontype.Double {
		panic(ElementTypeError{"bsoncore.Value.Double", v.Type})
	}
	f64, _, ok := ReadDouble(v.Data)
	if !ok {
		panic(NewInsufficientBytesError(v.Data, v.Data))
	}
	return f64
}


func (v Value) DoubleOK() (float64, bool) {
	if v.Type != bsontype.Double {
		return 0, false
	}
	f64, _, ok := ReadDouble(v.Data)
	if !ok {
		return 0, false
	}
	return f64, true
}






func (v Value) StringValue() string {
	if v.Type != bsontype.String {
		panic(ElementTypeError{"bsoncore.Value.StringValue", v.Type})
	}
	str, _, ok := ReadString(v.Data)
	if !ok {
		panic(NewInsufficientBytesError(v.Data, v.Data))
	}
	return str
}



func (v Value) StringValueOK() (string, bool) {
	if v.Type != bsontype.String {
		return "", false
	}
	str, _, ok := ReadString(v.Data)
	if !ok {
		return "", false
	}
	return str, true
}



func (v Value) Document() Document {
	if v.Type != bsontype.EmbeddedDocument {
		panic(ElementTypeError{"bsoncore.Value.Document", v.Type})
	}
	doc, _, ok := ReadDocument(v.Data)
	if !ok {
		panic(NewInsufficientBytesError(v.Data, v.Data))
	}
	return doc
}



func (v Value) DocumentOK() (Document, bool) {
	if v.Type != bsontype.EmbeddedDocument {
		return nil, false
	}
	doc, _, ok := ReadDocument(v.Data)
	if !ok {
		return nil, false
	}
	return doc, true
}



func (v Value) Array() Array {
	if v.Type != bsontype.Array {
		panic(ElementTypeError{"bsoncore.Value.Array", v.Type})
	}
	arr, _, ok := ReadArray(v.Data)
	if !ok {
		panic(NewInsufficientBytesError(v.Data, v.Data))
	}
	return arr
}



func (v Value) ArrayOK() (Array, bool) {
	if v.Type != bsontype.Array {
		return nil, false
	}
	arr, _, ok := ReadArray(v.Data)
	if !ok {
		return nil, false
	}
	return arr, true
}



func (v Value) Binary() (subtype byte, data []byte) {
	if v.Type != bsontype.Binary {
		panic(ElementTypeError{"bsoncore.Value.Binary", v.Type})
	}
	subtype, data, _, ok := ReadBinary(v.Data)
	if !ok {
		panic(NewInsufficientBytesError(v.Data, v.Data))
	}
	return subtype, data
}



func (v Value) BinaryOK() (subtype byte, data []byte, ok bool) {
	if v.Type != bsontype.Binary {
		return 0x00, nil, false
	}
	subtype, data, _, ok = ReadBinary(v.Data)
	if !ok {
		return 0x00, nil, false
	}
	return subtype, data, true
}



func (v Value) ObjectID() primitive.ObjectID {
	if v.Type != bsontype.ObjectID {
		panic(ElementTypeError{"bsoncore.Value.ObjectID", v.Type})
	}
	oid, _, ok := ReadObjectID(v.Data)
	if !ok {
		panic(NewInsufficientBytesError(v.Data, v.Data))
	}
	return oid
}



func (v Value) ObjectIDOK() (primitive.ObjectID, bool) {
	if v.Type != bsontype.ObjectID {
		return primitive.ObjectID{}, false
	}
	oid, _, ok := ReadObjectID(v.Data)
	if !ok {
		return primitive.ObjectID{}, false
	}
	return oid, true
}



func (v Value) Boolean() bool {
	if v.Type != bsontype.Boolean {
		panic(ElementTypeError{"bsoncore.Value.Boolean", v.Type})
	}
	b, _, ok := ReadBoolean(v.Data)
	if !ok {
		panic(NewInsufficientBytesError(v.Data, v.Data))
	}
	return b
}



func (v Value) BooleanOK() (bool, bool) {
	if v.Type != bsontype.Boolean {
		return false, false
	}
	b, _, ok := ReadBoolean(v.Data)
	if !ok {
		return false, false
	}
	return b, true
}



func (v Value) DateTime() int64 {
	if v.Type != bsontype.DateTime {
		panic(ElementTypeError{"bsoncore.Value.DateTime", v.Type})
	}
	dt, _, ok := ReadDateTime(v.Data)
	if !ok {
		panic(NewInsufficientBytesError(v.Data, v.Data))
	}
	return dt
}



func (v Value) DateTimeOK() (int64, bool) {
	if v.Type != bsontype.DateTime {
		return 0, false
	}
	dt, _, ok := ReadDateTime(v.Data)
	if !ok {
		return 0, false
	}
	return dt, true
}



func (v Value) Time() time.Time {
	if v.Type != bsontype.DateTime {
		panic(ElementTypeError{"bsoncore.Value.Time", v.Type})
	}
	dt, _, ok := ReadDateTime(v.Data)
	if !ok {
		panic(NewInsufficientBytesError(v.Data, v.Data))
	}
	return time.Unix(dt/1000, dt%1000*1000000)
}



func (v Value) TimeOK() (time.Time, bool) {
	if v.Type != bsontype.DateTime {
		return time.Time{}, false
	}
	dt, _, ok := ReadDateTime(v.Data)
	if !ok {
		return time.Time{}, false
	}
	return time.Unix(dt/1000, dt%1000*1000000), true
}



func (v Value) Regex() (pattern, options string) {
	if v.Type != bsontype.Regex {
		panic(ElementTypeError{"bsoncore.Value.Regex", v.Type})
	}
	pattern, options, _, ok := ReadRegex(v.Data)
	if !ok {
		panic(NewInsufficientBytesError(v.Data, v.Data))
	}
	return pattern, options
}



func (v Value) RegexOK() (pattern, options string, ok bool) {
	if v.Type != bsontype.Regex {
		return "", "", false
	}
	pattern, options, _, ok = ReadRegex(v.Data)
	if !ok {
		return "", "", false
	}
	return pattern, options, true
}



func (v Value) DBPointer() (string, primitive.ObjectID) {
	if v.Type != bsontype.DBPointer {
		panic(ElementTypeError{"bsoncore.Value.DBPointer", v.Type})
	}
	ns, pointer, _, ok := ReadDBPointer(v.Data)
	if !ok {
		panic(NewInsufficientBytesError(v.Data, v.Data))
	}
	return ns, pointer
}



func (v Value) DBPointerOK() (string, primitive.ObjectID, bool) {
	if v.Type != bsontype.DBPointer {
		return "", primitive.ObjectID{}, false
	}
	ns, pointer, _, ok := ReadDBPointer(v.Data)
	if !ok {
		return "", primitive.ObjectID{}, false
	}
	return ns, pointer, true
}



func (v Value) JavaScript() string {
	if v.Type != bsontype.JavaScript {
		panic(ElementTypeError{"bsoncore.Value.JavaScript", v.Type})
	}
	js, _, ok := ReadJavaScript(v.Data)
	if !ok {
		panic(NewInsufficientBytesError(v.Data, v.Data))
	}
	return js
}



func (v Value) JavaScriptOK() (string, bool) {
	if v.Type != bsontype.JavaScript {
		return "", false
	}
	js, _, ok := ReadJavaScript(v.Data)
	if !ok {
		return "", false
	}
	return js, true
}



func (v Value) Symbol() string {
	if v.Type != bsontype.Symbol {
		panic(ElementTypeError{"bsoncore.Value.Symbol", v.Type})
	}
	symbol, _, ok := ReadSymbol(v.Data)
	if !ok {
		panic(NewInsufficientBytesError(v.Data, v.Data))
	}
	return symbol
}



func (v Value) SymbolOK() (string, bool) {
	if v.Type != bsontype.Symbol {
		return "", false
	}
	symbol, _, ok := ReadSymbol(v.Data)
	if !ok {
		return "", false
	}
	return symbol, true
}



func (v Value) CodeWithScope() (string, Document) {
	if v.Type != bsontype.CodeWithScope {
		panic(ElementTypeError{"bsoncore.Value.CodeWithScope", v.Type})
	}
	code, scope, _, ok := ReadCodeWithScope(v.Data)
	if !ok {
		panic(NewInsufficientBytesError(v.Data, v.Data))
	}
	return code, scope
}



func (v Value) CodeWithScopeOK() (string, Document, bool) {
	if v.Type != bsontype.CodeWithScope {
		return "", nil, false
	}
	code, scope, _, ok := ReadCodeWithScope(v.Data)
	if !ok {
		return "", nil, false
	}
	return code, scope, true
}



func (v Value) Int32() int32 {
	if v.Type != bsontype.Int32 {
		panic(ElementTypeError{"bsoncore.Value.Int32", v.Type})
	}
	i32, _, ok := ReadInt32(v.Data)
	if !ok {
		panic(NewInsufficientBytesError(v.Data, v.Data))
	}
	return i32
}



func (v Value) Int32OK() (int32, bool) {
	if v.Type != bsontype.Int32 {
		return 0, false
	}
	i32, _, ok := ReadInt32(v.Data)
	if !ok {
		return 0, false
	}
	return i32, true
}



func (v Value) Timestamp() (t, i uint32) {
	if v.Type != bsontype.Timestamp {
		panic(ElementTypeError{"bsoncore.Value.Timestamp", v.Type})
	}
	t, i, _, ok := ReadTimestamp(v.Data)
	if !ok {
		panic(NewInsufficientBytesError(v.Data, v.Data))
	}
	return t, i
}



func (v Value) TimestampOK() (t, i uint32, ok bool) {
	if v.Type != bsontype.Timestamp {
		return 0, 0, false
	}
	t, i, _, ok = ReadTimestamp(v.Data)
	if !ok {
		return 0, 0, false
	}
	return t, i, true
}



func (v Value) Int64() int64 {
	if v.Type != bsontype.Int64 {
		panic(ElementTypeError{"bsoncore.Value.Int64", v.Type})
	}
	i64, _, ok := ReadInt64(v.Data)
	if !ok {
		panic(NewInsufficientBytesError(v.Data, v.Data))
	}
	return i64
}



func (v Value) Int64OK() (int64, bool) {
	if v.Type != bsontype.Int64 {
		return 0, false
	}
	i64, _, ok := ReadInt64(v.Data)
	if !ok {
		return 0, false
	}
	return i64, true
}



func (v Value) Decimal128() primitive.Decimal128 {
	if v.Type != bsontype.Decimal128 {
		panic(ElementTypeError{"bsoncore.Value.Decimal128", v.Type})
	}
	d128, _, ok := ReadDecimal128(v.Data)
	if !ok {
		panic(NewInsufficientBytesError(v.Data, v.Data))
	}
	return d128
}



func (v Value) Decimal128OK() (primitive.Decimal128, bool) {
	if v.Type != bsontype.Decimal128 {
		return primitive.Decimal128{}, false
	}
	d128, _, ok := ReadDecimal128(v.Data)
	if !ok {
		return primitive.Decimal128{}, false
	}
	return d128, true
}

var hexChars = "0123456789abcdef"

func escapeString(s string) string {
	escapeHTML := true
	var buf bytes.Buffer
	buf.WriteByte('"')
	start := 0
	for i := 0; i < len(s); {
		if b := s[i]; b < utf8.RuneSelf {
			if htmlSafeSet[b] || (!escapeHTML && safeSet[b]) {
				i++
				continue
			}
			if start < i {
				buf.WriteString(s[start:i])
			}
			switch b {
			case '\\', '"':
				buf.WriteByte('\\')
				buf.WriteByte(b)
			case '\n':
				buf.WriteByte('\\')
				buf.WriteByte('n')
			case '\r':
				buf.WriteByte('\\')
				buf.WriteByte('r')
			case '\t':
				buf.WriteByte('\\')
				buf.WriteByte('t')
			case '\b':
				buf.WriteByte('\\')
				buf.WriteByte('b')
			case '\f':
				buf.WriteByte('\\')
				buf.WriteByte('f')
			default:
				
				
				
				
				
				buf.WriteString(`\u00`)
				buf.WriteByte(hexChars[b>>4])
				buf.WriteByte(hexChars[b&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				buf.WriteString(s[start:i])
			}
			buf.WriteString(`\ufffd`)
			i += size
			start = i
			continue
		}
		
		
		
		
		
		
		
		if c == '\u2028' || c == '\u2029' {
			if start < i {
				buf.WriteString(s[start:i])
			}
			buf.WriteString(`\u202`)
			buf.WriteByte(hexChars[c&0xF])
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(s) {
		buf.WriteString(s[start:])
	}
	buf.WriteByte('"')
	return buf.String()
}

func formatDouble(f float64) string {
	var s string
	switch {
	case math.IsInf(f, 1):
		s = "Infinity"
	case math.IsInf(f, -1):
		s = "-Infinity"
	case math.IsNaN(f):
		s = "NaN"
	default:
		
		
		s = strconv.FormatFloat(f, 'G', -1, 64)
		if !strings.ContainsRune(s, '.') {
			s += ".0"
		}
	}

	return s
}

type sortableString []rune

func (ss sortableString) Len() int {
	return len(ss)
}

func (ss sortableString) Less(i, j int) bool {
	return ss[i] < ss[j]
}

func (ss sortableString) Swap(i, j int) {
	ss[i], ss[j] = ss[j], ss[i]
}

func sortStringAlphebeticAscending(s string) string {
	ss := sortableString([]rune(s))
	sort.Sort(ss)
	return string([]rune(ss))
}
