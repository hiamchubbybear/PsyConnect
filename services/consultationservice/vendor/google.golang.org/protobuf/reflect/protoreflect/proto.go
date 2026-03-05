



































































































































package protoreflect

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/internal/pragma"
)

type doNotImplement pragma.DoNotImplement




type ProtoMessage interface{ ProtoReflect() Message }


type Syntax syntax

type syntax int8 

const (
	Proto2   Syntax = 2
	Proto3   Syntax = 3
	Editions Syntax = 4
)


func (s Syntax) IsValid() bool {
	switch s {
	case Proto2, Proto3, Editions:
		return true
	default:
		return false
	}
}


func (s Syntax) String() string {
	switch s {
	case Proto2:
		return "proto2"
	case Proto3:
		return "proto3"
	case Editions:
		return "editions"
	default:
		return fmt.Sprintf("<unknown:%d>", s)
	}
}


func (s Syntax) GoString() string {
	switch s {
	case Proto2:
		return "Proto2"
	case Proto3:
		return "Proto3"
	default:
		return fmt.Sprintf("Syntax(%d)", s)
	}
}


type Cardinality cardinality

type cardinality int8 


const (
	Optional Cardinality = 1 
	Required Cardinality = 2 
	Repeated Cardinality = 3 
)


func (c Cardinality) IsValid() bool {
	switch c {
	case Optional, Required, Repeated:
		return true
	default:
		return false
	}
}


func (c Cardinality) String() string {
	switch c {
	case Optional:
		return "optional"
	case Required:
		return "required"
	case Repeated:
		return "repeated"
	default:
		return fmt.Sprintf("<unknown:%d>", c)
	}
}


func (c Cardinality) GoString() string {
	switch c {
	case Optional:
		return "Optional"
	case Required:
		return "Required"
	case Repeated:
		return "Repeated"
	default:
		return fmt.Sprintf("Cardinality(%d)", c)
	}
}


type Kind kind

type kind int8 


const (
	BoolKind     Kind = 8
	EnumKind     Kind = 14
	Int32Kind    Kind = 5
	Sint32Kind   Kind = 17
	Uint32Kind   Kind = 13
	Int64Kind    Kind = 3
	Sint64Kind   Kind = 18
	Uint64Kind   Kind = 4
	Sfixed32Kind Kind = 15
	Fixed32Kind  Kind = 7
	FloatKind    Kind = 2
	Sfixed64Kind Kind = 16
	Fixed64Kind  Kind = 6
	DoubleKind   Kind = 1
	StringKind   Kind = 9
	BytesKind    Kind = 12
	MessageKind  Kind = 11
	GroupKind    Kind = 10
)


func (k Kind) IsValid() bool {
	switch k {
	case BoolKind, EnumKind,
		Int32Kind, Sint32Kind, Uint32Kind,
		Int64Kind, Sint64Kind, Uint64Kind,
		Sfixed32Kind, Fixed32Kind, FloatKind,
		Sfixed64Kind, Fixed64Kind, DoubleKind,
		StringKind, BytesKind, MessageKind, GroupKind:
		return true
	default:
		return false
	}
}


func (k Kind) String() string {
	switch k {
	case BoolKind:
		return "bool"
	case EnumKind:
		return "enum"
	case Int32Kind:
		return "int32"
	case Sint32Kind:
		return "sint32"
	case Uint32Kind:
		return "uint32"
	case Int64Kind:
		return "int64"
	case Sint64Kind:
		return "sint64"
	case Uint64Kind:
		return "uint64"
	case Sfixed32Kind:
		return "sfixed32"
	case Fixed32Kind:
		return "fixed32"
	case FloatKind:
		return "float"
	case Sfixed64Kind:
		return "sfixed64"
	case Fixed64Kind:
		return "fixed64"
	case DoubleKind:
		return "double"
	case StringKind:
		return "string"
	case BytesKind:
		return "bytes"
	case MessageKind:
		return "message"
	case GroupKind:
		return "group"
	default:
		return fmt.Sprintf("<unknown:%d>", k)
	}
}


func (k Kind) GoString() string {
	switch k {
	case BoolKind:
		return "BoolKind"
	case EnumKind:
		return "EnumKind"
	case Int32Kind:
		return "Int32Kind"
	case Sint32Kind:
		return "Sint32Kind"
	case Uint32Kind:
		return "Uint32Kind"
	case Int64Kind:
		return "Int64Kind"
	case Sint64Kind:
		return "Sint64Kind"
	case Uint64Kind:
		return "Uint64Kind"
	case Sfixed32Kind:
		return "Sfixed32Kind"
	case Fixed32Kind:
		return "Fixed32Kind"
	case FloatKind:
		return "FloatKind"
	case Sfixed64Kind:
		return "Sfixed64Kind"
	case Fixed64Kind:
		return "Fixed64Kind"
	case DoubleKind:
		return "DoubleKind"
	case StringKind:
		return "StringKind"
	case BytesKind:
		return "BytesKind"
	case MessageKind:
		return "MessageKind"
	case GroupKind:
		return "GroupKind"
	default:
		return fmt.Sprintf("Kind(%d)", k)
	}
}


type FieldNumber = protowire.Number


type FieldNumbers interface {
	
	Len() int
	
	Get(i int) FieldNumber
	
	Has(n FieldNumber) bool

	doNotImplement
}


type FieldRanges interface {
	
	Len() int
	
	Get(i int) [2]FieldNumber 
	
	Has(n FieldNumber) bool

	doNotImplement
}


type EnumNumber int32


type EnumRanges interface {
	
	Len() int
	
	Get(i int) [2]EnumNumber 
	
	Has(n EnumNumber) bool

	doNotImplement
}



type Name string 



func (s Name) IsValid() bool {
	return consumeIdent(string(s)) == len(s)
}


type Names interface {
	
	Len() int
	
	Get(i int) Name
	
	Has(s Name) bool

	doNotImplement
}







type FullName string 



func (s FullName) IsValid() bool {
	i := consumeIdent(string(s))
	if i < 0 {
		return false
	}
	for len(s) > i {
		if s[i] != '.' {
			return false
		}
		i++
		n := consumeIdent(string(s[i:]))
		if n < 0 {
			return false
		}
		i += n
	}
	return true
}

func consumeIdent(s string) (i int) {
	if len(s) == 0 || !isLetter(s[i]) {
		return -1
	}
	i++
	for len(s) > i && isLetterDigit(s[i]) {
		i++
	}
	return i
}
func isLetter(c byte) bool {
	return c == '_' || ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}
func isLetterDigit(c byte) bool {
	return isLetter(c) || ('0' <= c && c <= '9')
}



func (n FullName) Name() Name {
	if i := strings.LastIndexByte(string(n), '.'); i >= 0 {
		return Name(n[i+1:])
	}
	return Name(n)
}



func (n FullName) Parent() FullName {
	if i := strings.LastIndexByte(string(n), '.'); i >= 0 {
		return n[:i]
	}
	return ""
}




func (n FullName) Append(s Name) FullName {
	if n == "" {
		return FullName(s)
	}
	return n + "." + FullName(s)
}
