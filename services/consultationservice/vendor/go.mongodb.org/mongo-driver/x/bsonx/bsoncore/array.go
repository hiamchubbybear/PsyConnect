





package bsoncore

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)



func NewArrayLengthError(length, rem int) error {
	return lengthError("array", length, rem)
}


type Array []byte



func NewArrayFromReader(r io.Reader) (Array, error) {
	return newBufferFromReader(r)
}



func (a Array) Index(index uint) Value {
	value, err := a.IndexErr(index)
	if err != nil {
		panic(err)
	}
	return value
}


func (a Array) IndexErr(index uint) (Value, error) {
	elem, err := indexErr(a, index)
	if err != nil {
		return Value{}, err
	}
	return elem.Value(), err
}



func (a Array) DebugString() string {
	if len(a) < 5 {
		return "<malformed>"
	}
	var buf strings.Builder
	buf.WriteString("Array")
	length, rem, _ := ReadLength(a) 
	buf.WriteByte('(')
	buf.WriteString(strconv.Itoa(int(length)))
	length -= 4
	buf.WriteString(")[")
	var elem Element
	var ok bool
	for length > 1 {
		elem, rem, ok = ReadElement(rem)
		length -= int32(len(elem))
		if !ok {
			buf.WriteString(fmt.Sprintf("<malformed (%d)>", length))
			break
		}
		buf.WriteString(elem.Value().DebugString())
		if length != 1 {
			buf.WriteByte(',')
		}
	}
	buf.WriteByte(']')

	return buf.String()
}



func (a Array) String() string {
	if len(a) < 5 {
		return ""
	}
	var buf strings.Builder
	buf.WriteByte('[')

	length, rem, _ := ReadLength(a) 

	length -= 4

	var elem Element
	var ok bool
	for length > 1 {
		elem, rem, ok = ReadElement(rem)
		length -= int32(len(elem))
		if !ok {
			return ""
		}
		buf.WriteString(elem.Value().String())
		if length > 1 {
			buf.WriteByte(',')
		}
	}
	if length != 1 { 
		return ""
	}

	buf.WriteByte(']')
	return buf.String()
}




func (a Array) Values() ([]Value, error) {
	return values(a)
}


func (a Array) Validate() error {
	length, rem, ok := ReadLength(a)
	if !ok {
		return NewInsufficientBytesError(a, rem)
	}
	if int(length) > len(a) {
		return NewArrayLengthError(int(length), len(a))
	}
	if a[length-1] != 0x00 {
		return ErrMissingNull
	}

	length -= 4
	var elem Element

	var keyNum int64
	for length > 1 {
		elem, rem, ok = ReadElement(rem)
		length -= int32(len(elem))
		if !ok {
			return NewInsufficientBytesError(a, rem)
		}

		
		err := elem.Validate()
		if err != nil {
			return err
		}

		
		if fmt.Sprint(keyNum) != elem.Key() {
			return fmt.Errorf("array key %q is out of order or invalid", elem.Key())
		}
		keyNum++
	}

	if len(rem) < 1 || rem[0] != 0x00 {
		return ErrMissingNull
	}
	return nil
}
