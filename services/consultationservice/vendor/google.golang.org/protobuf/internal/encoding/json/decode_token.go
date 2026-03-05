



package json

import (
	"bytes"
	"fmt"
	"strconv"
)


type Kind uint16

const (
	Invalid Kind = (1 << iota) / 2
	EOF
	Null
	Bool
	Number
	String
	Name
	ObjectOpen
	ObjectClose
	ArrayOpen
	ArrayClose

	
	
	comma
)

func (k Kind) String() string {
	switch k {
	case EOF:
		return "eof"
	case Null:
		return "null"
	case Bool:
		return "bool"
	case Number:
		return "number"
	case String:
		return "string"
	case ObjectOpen:
		return "{"
	case ObjectClose:
		return "}"
	case Name:
		return "name"
	case ArrayOpen:
		return "["
	case ArrayClose:
		return "]"
	case comma:
		return ","
	}
	return "<invalid>"
}







type Token struct {
	
	kind Kind
	
	pos int
	
	
	raw []byte
	
	boo bool
	
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


func (t Token) Name() string {
	if t.kind == Name {
		return t.str
	}
	panic(fmt.Sprintf("Token is not a Name: %v", t.RawString()))
}


func (t Token) Bool() bool {
	if t.kind == Bool {
		return t.boo
	}
	panic(fmt.Sprintf("Token is not a Bool: %v", t.RawString()))
}



func (t Token) ParsedString() string {
	if t.kind == String {
		return t.str
	}
	panic(fmt.Sprintf("Token is not a String: %v", t.RawString()))
}








func (t Token) Float(bitSize int) (float64, bool) {
	if t.kind != Number {
		return 0, false
	}
	f, err := strconv.ParseFloat(t.RawString(), bitSize)
	if err != nil {
		return 0, false
	}
	return f, true
}






func (t Token) Int(bitSize int) (int64, bool) {
	s, ok := t.getIntStr()
	if !ok {
		return 0, false
	}
	n, err := strconv.ParseInt(s, 10, bitSize)
	if err != nil {
		return 0, false
	}
	return n, true
}






func (t Token) Uint(bitSize int) (uint64, bool) {
	s, ok := t.getIntStr()
	if !ok {
		return 0, false
	}
	n, err := strconv.ParseUint(s, 10, bitSize)
	if err != nil {
		return 0, false
	}
	return n, true
}

func (t Token) getIntStr() (string, bool) {
	if t.kind != Number {
		return "", false
	}
	parts, ok := parseNumberParts(t.raw)
	if !ok {
		return "", false
	}
	return normalizeToIntString(parts)
}


func TokenEquals(x, y Token) bool {
	return x.kind == y.kind &&
		x.pos == y.pos &&
		bytes.Equal(x.raw, y.raw) &&
		x.boo == y.boo &&
		x.str == y.str
}
