



package text

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"

	"google.golang.org/protobuf/internal/flags"
)


type Kind uint8


const (
	Invalid Kind = iota
	EOF
	Name   
	Scalar 
	MessageOpen
	MessageClose
	ListOpen
	ListClose

	
	comma
	semicolon

	
	
	bof = Invalid
)

func (t Kind) String() string {
	switch t {
	case Invalid:
		return "<invalid>"
	case EOF:
		return "eof"
	case Scalar:
		return "scalar"
	case Name:
		return "name"
	case MessageOpen:
		return "{"
	case MessageClose:
		return "}"
	case ListOpen:
		return "["
	case ListClose:
		return "]"
	case comma:
		return ","
	case semicolon:
		return ";"
	default:
		return fmt.Sprintf("<invalid:%v>", uint8(t))
	}
}


type NameKind uint8


const (
	IdentName NameKind = iota + 1
	TypeName
	FieldNumber
)

func (t NameKind) String() string {
	switch t {
	case IdentName:
		return "IdentName"
	case TypeName:
		return "TypeName"
	case FieldNumber:
		return "FieldNumber"
	default:
		return fmt.Sprintf("<invalid:%v>", uint8(t))
	}
}






const hasSeparator = 1 << 7


const (
	numberValue = iota + 1
	stringValue
	literalValue
)


const isNegative = 1 << 7



type Token struct {
	
	kind Kind
	
	
	
	attrs uint8
	
	
	
	numAttrs uint8
	
	pos int
	
	
	raw []byte
	
	
	
	
	str string
}


func (t Token) Kind() Kind {
	return t.kind
}


func (t Token) RawString() string {
	return string(t.raw)
}


func (t Token) Pos() int {
	return t.pos
}



func (t Token) NameKind() NameKind {
	if t.kind == Name {
		return NameKind(t.attrs &^ hasSeparator)
	}
	panic(fmt.Sprintf("Token is not a Name type: %s", t.kind))
}



func (t Token) HasSeparator() bool {
	if t.kind == Name {
		return t.attrs&hasSeparator != 0
	}
	panic(fmt.Sprintf("Token is not a Name type: %s", t.kind))
}


func (t Token) IdentName() string {
	if t.kind == Name && t.attrs&uint8(IdentName) != 0 {
		return string(t.raw)
	}
	panic(fmt.Sprintf("Token is not an IdentName: %s:%s", t.kind, NameKind(t.attrs&^hasSeparator)))
}


func (t Token) TypeName() string {
	if t.kind == Name && t.attrs&uint8(TypeName) != 0 {
		return t.str
	}
	panic(fmt.Sprintf("Token is not a TypeName: %s:%s", t.kind, NameKind(t.attrs&^hasSeparator)))
}




func (t Token) FieldNumber() int32 {
	if t.kind != Name || t.attrs&uint8(FieldNumber) == 0 {
		panic(fmt.Sprintf("Token is not a FieldNumber: %s:%s", t.kind, NameKind(t.attrs&^hasSeparator)))
	}
	
	
	num, _ := strconv.ParseInt(string(t.raw), 10, 32)
	return int32(num)
}


func (t Token) String() (string, bool) {
	if t.kind != Scalar || t.attrs != stringValue {
		return "", false
	}
	return t.str, true
}


func (t Token) Enum() (string, bool) {
	if t.kind != Scalar || t.attrs != literalValue || (len(t.raw) > 0 && t.raw[0] == '-') {
		return "", false
	}
	return string(t.raw), true
}


func (t Token) Bool() (bool, bool) {
	if t.kind != Scalar {
		return false, false
	}
	switch t.attrs {
	case literalValue:
		if b, ok := boolLits[string(t.raw)]; ok {
			return b, true
		}
	case numberValue:
		
		
		n, err := strconv.ParseUint(t.str, 0, 64)
		if err == nil {
			switch n {
			case 0:
				return false, true
			case 1:
				return true, true
			}
		}
	}
	return false, false
}


var boolLits = map[string]bool{
	"t":     true,
	"true":  true,
	"True":  true,
	"f":     false,
	"false": false,
	"False": false,
}


func (t Token) Uint64() (uint64, bool) {
	if t.kind != Scalar || t.attrs != numberValue ||
		t.numAttrs&isNegative > 0 || t.numAttrs&numFloat > 0 {
		return 0, false
	}
	n, err := strconv.ParseUint(t.str, 0, 64)
	if err != nil {
		return 0, false
	}
	return n, true
}


func (t Token) Uint32() (uint32, bool) {
	if t.kind != Scalar || t.attrs != numberValue ||
		t.numAttrs&isNegative > 0 || t.numAttrs&numFloat > 0 {
		return 0, false
	}
	n, err := strconv.ParseUint(t.str, 0, 32)
	if err != nil {
		return 0, false
	}
	return uint32(n), true
}


func (t Token) Int64() (int64, bool) {
	if t.kind != Scalar || t.attrs != numberValue || t.numAttrs&numFloat > 0 {
		return 0, false
	}
	if n, err := strconv.ParseInt(t.str, 0, 64); err == nil {
		return n, true
	}
	
	
	if flags.ProtoLegacy && (t.numAttrs == numHex) {
		if n, err := strconv.ParseUint(t.str, 0, 64); err == nil {
			return int64(n), true
		}
	}
	return 0, false
}


func (t Token) Int32() (int32, bool) {
	if t.kind != Scalar || t.attrs != numberValue || t.numAttrs&numFloat > 0 {
		return 0, false
	}
	if n, err := strconv.ParseInt(t.str, 0, 32); err == nil {
		return int32(n), true
	}
	
	
	if flags.ProtoLegacy && (t.numAttrs == numHex) {
		if n, err := strconv.ParseUint(t.str, 0, 32); err == nil {
			return int32(n), true
		}
	}
	return 0, false
}


func (t Token) Float64() (float64, bool) {
	if t.kind != Scalar {
		return 0, false
	}
	switch t.attrs {
	case literalValue:
		if f, ok := floatLits[strings.ToLower(string(t.raw))]; ok {
			return f, true
		}
	case numberValue:
		n, err := strconv.ParseFloat(t.str, 64)
		if err == nil {
			return n, true
		}
		nerr := err.(*strconv.NumError)
		if nerr.Err == strconv.ErrRange {
			return n, true
		}
	}
	return 0, false
}


func (t Token) Float32() (float32, bool) {
	if t.kind != Scalar {
		return 0, false
	}
	switch t.attrs {
	case literalValue:
		if f, ok := floatLits[strings.ToLower(string(t.raw))]; ok {
			return float32(f), true
		}
	case numberValue:
		n, err := strconv.ParseFloat(t.str, 64)
		if err == nil {
			
			return float32(n), true
		}
		nerr := err.(*strconv.NumError)
		if nerr.Err == strconv.ErrRange {
			return float32(n), true
		}
	}
	return 0, false
}



var floatLits = map[string]float64{
	"nan":       math.NaN(),
	"inf":       math.Inf(1),
	"infinity":  math.Inf(1),
	"-inf":      math.Inf(-1),
	"-infinity": math.Inf(-1),
}


func TokenEquals(x, y Token) bool {
	return x.kind == y.kind &&
		x.attrs == y.attrs &&
		x.numAttrs == y.numAttrs &&
		x.pos == y.pos &&
		bytes.Equal(x.raw, y.raw) &&
		x.str == y.str
}
