package jsoniter

import (
	"bytes"
	"io"
)


type RawMessage []byte





func Unmarshal(data []byte, v interface{}) error {
	return ConfigDefault.Unmarshal(data, v)
}


func UnmarshalFromString(str string, v interface{}) error {
	return ConfigDefault.UnmarshalFromString(str, v)
}


func Get(data []byte, path ...interface{}) Any {
	return ConfigDefault.Get(data, path...)
}





func Marshal(v interface{}) ([]byte, error) {
	return ConfigDefault.Marshal(v)
}


func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return ConfigDefault.MarshalIndent(v, prefix, indent)
}


func MarshalToString(v interface{}) (string, error) {
	return ConfigDefault.MarshalToString(v)
}







func NewDecoder(reader io.Reader) *Decoder {
	return ConfigDefault.NewDecoder(reader)
}



type Decoder struct {
	iter *Iterator
}


func (adapter *Decoder) Decode(obj interface{}) error {
	if adapter.iter.head == adapter.iter.tail && adapter.iter.reader != nil {
		if !adapter.iter.loadMore() {
			return io.EOF
		}
	}
	adapter.iter.ReadVal(obj)
	err := adapter.iter.Error
	if err == io.EOF {
		return nil
	}
	return adapter.iter.Error
}


func (adapter *Decoder) More() bool {
	iter := adapter.iter
	if iter.Error != nil {
		return false
	}
	c := iter.nextToken()
	if c == 0 {
		return false
	}
	iter.unreadByte()
	return c != ']' && c != '}'
}


func (adapter *Decoder) Buffered() io.Reader {
	remaining := adapter.iter.buf[adapter.iter.head:adapter.iter.tail]
	return bytes.NewReader(remaining)
}



func (adapter *Decoder) UseNumber() {
	cfg := adapter.iter.cfg.configBeforeFrozen
	cfg.UseNumber = true
	adapter.iter.cfg = cfg.frozeWithCacheReuse(adapter.iter.cfg.extraExtensions)
}




func (adapter *Decoder) DisallowUnknownFields() {
	cfg := adapter.iter.cfg.configBeforeFrozen
	cfg.DisallowUnknownFields = true
	adapter.iter.cfg = cfg.frozeWithCacheReuse(adapter.iter.cfg.extraExtensions)
}


func NewEncoder(writer io.Writer) *Encoder {
	return ConfigDefault.NewEncoder(writer)
}


type Encoder struct {
	stream *Stream
}


func (adapter *Encoder) Encode(val interface{}) error {
	adapter.stream.WriteVal(val)
	adapter.stream.WriteRaw("\n")
	adapter.stream.Flush()
	return adapter.stream.Error
}


func (adapter *Encoder) SetIndent(prefix, indent string) {
	config := adapter.stream.cfg.configBeforeFrozen
	config.IndentionStep = len(indent)
	adapter.stream.cfg = config.frozeWithCacheReuse(adapter.stream.cfg.extraExtensions)
}


func (adapter *Encoder) SetEscapeHTML(escapeHTML bool) {
	config := adapter.stream.cfg.configBeforeFrozen
	config.EscapeHTML = escapeHTML
	adapter.stream.cfg = config.frozeWithCacheReuse(adapter.stream.cfg.extraExtensions)
}


func Valid(data []byte) bool {
	return ConfigDefault.Valid(data)
}
