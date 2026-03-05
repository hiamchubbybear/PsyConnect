





package bsoncore

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)


type ValidationError string

func (ve ValidationError) Error() string { return string(ve) }



func NewDocumentLengthError(length, rem int) error {
	return lengthError("document", length, rem)
}

func lengthError(bufferType string, length, rem int) error {
	return ValidationError(fmt.Sprintf("%v length exceeds available bytes. length=%d remainingBytes=%d",
		bufferType, length, rem))
}


type InsufficientBytesError struct {
	Source    []byte
	Remaining []byte
}



func NewInsufficientBytesError(src, rem []byte) InsufficientBytesError {
	return InsufficientBytesError{Source: src, Remaining: rem}
}


func (ibe InsufficientBytesError) Error() string {
	return "too few bytes to read next component"
}


func (ibe InsufficientBytesError) Equal(err2 error) bool {
	switch err2.(type) {
	case InsufficientBytesError:
		return true
	default:
		return false
	}
}



type InvalidDepthTraversalError struct {
	Key  string
	Type bsontype.Type
}

func (idte InvalidDepthTraversalError) Error() string {
	return fmt.Sprintf(
		"attempt to traverse into %s, but it's type is %s, not %s nor %s",
		idte.Key, idte.Type, bsontype.EmbeddedDocument, bsontype.Array,
	)
}


const ErrMissingNull ValidationError = "document or array end is missing null byte"



const ErrInvalidLength ValidationError = "document or array length is invalid"


var ErrNilReader = errors.New("nil reader")


var ErrEmptyKey = errors.New("empty key provided")


var ErrElementNotFound = errors.New("element not found")


var ErrOutOfBounds = errors.New("out of bounds")


type Document []byte



func NewDocumentFromReader(r io.Reader) (Document, error) {
	return newBufferFromReader(r)
}

func newBufferFromReader(r io.Reader) ([]byte, error) {
	if r == nil {
		return nil, ErrNilReader
	}

	var lengthBytes [4]byte

	
	_, err := io.ReadFull(r, lengthBytes[:])
	if err != nil {
		return nil, err
	}

	length, _, _ := readi32(lengthBytes[:]) 
	if length < 0 {
		return nil, ErrInvalidLength
	}
	buffer := make([]byte, length)

	copy(buffer, lengthBytes[:])

	_, err = io.ReadFull(r, buffer[4:])
	if err != nil {
		return nil, err
	}

	if buffer[length-1] != 0x00 {
		return nil, ErrMissingNull
	}

	return buffer, nil
}





func (d Document) Lookup(key ...string) Value {
	val, _ := d.LookupErr(key...)
	return val
}


func (d Document) LookupErr(key ...string) (Value, error) {
	if len(key) < 1 {
		return Value{}, ErrEmptyKey
	}
	length, rem, ok := ReadLength(d)
	if !ok {
		return Value{}, NewInsufficientBytesError(d, rem)
	}

	length -= 4

	var elem Element
	for length > 1 {
		elem, rem, ok = ReadElement(rem)
		length -= int32(len(elem))
		if !ok {
			return Value{}, NewInsufficientBytesError(d, rem)
		}
		
		if string(elem.KeyBytes()) != key[0] {
			continue
		}
		if len(key) > 1 {
			tt := bsontype.Type(elem[0])
			switch tt {
			case bsontype.EmbeddedDocument:
				val, err := elem.Value().Document().LookupErr(key[1:]...)
				if err != nil {
					return Value{}, err
				}
				return val, nil
			case bsontype.Array:
				
				val, err := Document(elem.Value().Array()).LookupErr(key[1:]...)
				if err != nil {
					return Value{}, err
				}
				return val, nil
			default:
				return Value{}, InvalidDepthTraversalError{Key: elem.Key(), Type: tt}
			}
		}
		return elem.ValueErr()
	}
	return Value{}, ErrElementNotFound
}



func (d Document) Index(index uint) Element {
	elem, err := d.IndexErr(index)
	if err != nil {
		panic(err)
	}
	return elem
}


func (d Document) IndexErr(index uint) (Element, error) {
	return indexErr(d, index)
}

func indexErr(b []byte, index uint) (Element, error) {
	length, rem, ok := ReadLength(b)
	if !ok {
		return nil, NewInsufficientBytesError(b, rem)
	}

	length -= 4

	var current uint
	var elem Element
	for length > 1 {
		elem, rem, ok = ReadElement(rem)
		length -= int32(len(elem))
		if !ok {
			return nil, NewInsufficientBytesError(b, rem)
		}
		if current != index {
			current++
			continue
		}
		return elem, nil
	}
	return nil, ErrOutOfBounds
}



func (d Document) DebugString() string {
	if len(d) < 5 {
		return "<malformed>"
	}
	var buf strings.Builder
	buf.WriteString("Document")
	length, rem, _ := ReadLength(d) 
	buf.WriteByte('(')
	buf.WriteString(strconv.Itoa(int(length)))
	length -= 4
	buf.WriteString("){")
	var elem Element
	var ok bool
	for length > 1 {
		elem, rem, ok = ReadElement(rem)
		length -= int32(len(elem))
		if !ok {
			buf.WriteString(fmt.Sprintf("<malformed (%d)>", length))
			break
		}
		buf.WriteString(elem.DebugString())
	}
	buf.WriteByte('}')

	return buf.String()
}



func (d Document) String() string {
	if len(d) < 5 {
		return ""
	}
	var buf strings.Builder
	buf.WriteByte('{')

	length, rem, _ := ReadLength(d) 

	length -= 4

	var elem Element
	var ok bool
	first := true
	for length > 1 {
		if !first {
			buf.WriteByte(',')
		}
		elem, rem, ok = ReadElement(rem)
		length -= int32(len(elem))
		if !ok {
			return ""
		}
		buf.WriteString(elem.String())
		first = false
	}
	buf.WriteByte('}')

	return buf.String()
}




func (d Document) Elements() ([]Element, error) {
	length, rem, ok := ReadLength(d)
	if !ok {
		return nil, NewInsufficientBytesError(d, rem)
	}

	length -= 4

	var elem Element
	var elems []Element
	for length > 1 {
		elem, rem, ok = ReadElement(rem)
		length -= int32(len(elem))
		if !ok {
			return elems, NewInsufficientBytesError(d, rem)
		}
		if err := elem.Validate(); err != nil {
			return elems, err
		}
		elems = append(elems, elem)
	}
	return elems, nil
}




func (d Document) Values() ([]Value, error) {
	return values(d)
}

func values(b []byte) ([]Value, error) {
	length, rem, ok := ReadLength(b)
	if !ok {
		return nil, NewInsufficientBytesError(b, rem)
	}

	length -= 4

	var elem Element
	var vals []Value
	for length > 1 {
		elem, rem, ok = ReadElement(rem)
		length -= int32(len(elem))
		if !ok {
			return vals, NewInsufficientBytesError(b, rem)
		}
		if err := elem.Value().Validate(); err != nil {
			return vals, err
		}
		vals = append(vals, elem.Value())
	}
	return vals, nil
}


func (d Document) Validate() error {
	length, rem, ok := ReadLength(d)
	if !ok {
		return NewInsufficientBytesError(d, rem)
	}
	if int(length) > len(d) {
		return NewDocumentLengthError(int(length), len(d))
	}
	if d[length-1] != 0x00 {
		return ErrMissingNull
	}

	length -= 4
	var elem Element

	for length > 1 {
		elem, rem, ok = ReadElement(rem)
		length -= int32(len(elem))
		if !ok {
			return NewInsufficientBytesError(d, rem)
		}
		err := elem.Validate()
		if err != nil {
			return err
		}
	}

	if len(rem) < 1 || rem[0] != 0x00 {
		return ErrMissingNull
	}
	return nil
}
