package json

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/goccy/go-json/internal/encoder"
)



type Marshaler interface {
	MarshalJSON() ([]byte, error)
}



type MarshalerContext interface {
	MarshalJSON(context.Context) ([]byte, error)
}









type Unmarshaler interface {
	UnmarshalJSON([]byte) error
}



type UnmarshalerContext interface {
	UnmarshalJSON(context.Context, []byte) error
}

































































































































func Marshal(v interface{}) ([]byte, error) {
	return MarshalWithOption(v)
}


func MarshalNoEscape(v interface{}) ([]byte, error) {
	return marshalNoEscape(v)
}


func MarshalContext(ctx context.Context, v interface{}, optFuncs ...EncodeOptionFunc) ([]byte, error) {
	return marshalContext(ctx, v, optFuncs...)
}


func MarshalWithOption(v interface{}, optFuncs ...EncodeOptionFunc) ([]byte, error) {
	return marshal(v, optFuncs...)
}




func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return MarshalIndentWithOption(v, prefix, indent)
}


func MarshalIndentWithOption(v interface{}, prefix, indent string, optFuncs ...EncodeOptionFunc) ([]byte, error) {
	return marshalIndent(v, prefix, indent, optFuncs...)
}










































































func Unmarshal(data []byte, v interface{}) error {
	return unmarshal(data, v)
}




func UnmarshalContext(ctx context.Context, data []byte, v interface{}, optFuncs ...DecodeOptionFunc) error {
	return unmarshalContext(ctx, data, v)
}

func UnmarshalWithOption(data []byte, v interface{}, optFuncs ...DecodeOptionFunc) error {
	return unmarshal(data, v, optFuncs...)
}

func UnmarshalNoEscape(data []byte, v interface{}, optFuncs ...DecodeOptionFunc) error {
	return unmarshalNoEscape(data, v, optFuncs...)
}









type Token = json.Token


type Number = json.Number




type RawMessage = json.RawMessage


type Delim = json.Delim



func Compact(dst *bytes.Buffer, src []byte) error {
	return encoder.Compact(dst, src, false)
}












func Indent(dst *bytes.Buffer, src []byte, prefix, indent string) error {
	return encoder.Indent(dst, src, prefix, indent)
}







func HTMLEscape(dst *bytes.Buffer, src []byte) {
	var v interface{}
	dec := NewDecoder(bytes.NewBuffer(src))
	dec.UseNumber()
	if err := dec.Decode(&v); err != nil {
		return
	}
	buf, _ := marshal(v)
	dst.Write(buf)
}


func Valid(data []byte) bool {
	var v interface{}
	decoder := NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&v)
	if err != nil {
		return false
	}
	if !decoder.More() {
		return true
	}
	return decoder.InputOffset() >= int64(len(data))
}

func init() {
	encoder.Marshal = Marshal
	encoder.Unmarshal = Unmarshal
}
