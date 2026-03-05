package jsoniter

import (
	"io"
)



type Stream struct {
	cfg        *frozenConfig
	out        io.Writer
	buf        []byte
	Error      error
	indention  int
	Attachment interface{} 
}





func NewStream(cfg API, out io.Writer, bufSize int) *Stream {
	return &Stream{
		cfg:       cfg.(*frozenConfig),
		out:       out,
		buf:       make([]byte, 0, bufSize),
		Error:     nil,
		indention: 0,
	}
}


func (stream *Stream) Pool() StreamPool {
	return stream.cfg
}


func (stream *Stream) Reset(out io.Writer) {
	stream.out = out
	stream.buf = stream.buf[:0]
}


func (stream *Stream) Available() int {
	return cap(stream.buf) - len(stream.buf)
}


func (stream *Stream) Buffered() int {
	return len(stream.buf)
}


func (stream *Stream) Buffer() []byte {
	return stream.buf
}


func (stream *Stream) SetBuffer(buf []byte) {
	stream.buf = buf
}





func (stream *Stream) Write(p []byte) (nn int, err error) {
	stream.buf = append(stream.buf, p...)
	if stream.out != nil {
		nn, err = stream.out.Write(stream.buf)
		stream.buf = stream.buf[nn:]
		return
	}
	return len(p), nil
}


func (stream *Stream) writeByte(c byte) {
	stream.buf = append(stream.buf, c)
}

func (stream *Stream) writeTwoBytes(c1 byte, c2 byte) {
	stream.buf = append(stream.buf, c1, c2)
}

func (stream *Stream) writeThreeBytes(c1 byte, c2 byte, c3 byte) {
	stream.buf = append(stream.buf, c1, c2, c3)
}

func (stream *Stream) writeFourBytes(c1 byte, c2 byte, c3 byte, c4 byte) {
	stream.buf = append(stream.buf, c1, c2, c3, c4)
}

func (stream *Stream) writeFiveBytes(c1 byte, c2 byte, c3 byte, c4 byte, c5 byte) {
	stream.buf = append(stream.buf, c1, c2, c3, c4, c5)
}


func (stream *Stream) Flush() error {
	if stream.out == nil {
		return nil
	}
	if stream.Error != nil {
		return stream.Error
	}
	_, err := stream.out.Write(stream.buf)
	if err != nil {
		if stream.Error == nil {
			stream.Error = err
		}
		return err
	}
	stream.buf = stream.buf[:0]
	return nil
}


func (stream *Stream) WriteRaw(s string) {
	stream.buf = append(stream.buf, s...)
}


func (stream *Stream) WriteNil() {
	stream.writeFourBytes('n', 'u', 'l', 'l')
}


func (stream *Stream) WriteTrue() {
	stream.writeFourBytes('t', 'r', 'u', 'e')
}


func (stream *Stream) WriteFalse() {
	stream.writeFiveBytes('f', 'a', 'l', 's', 'e')
}


func (stream *Stream) WriteBool(val bool) {
	if val {
		stream.WriteTrue()
	} else {
		stream.WriteFalse()
	}
}


func (stream *Stream) WriteObjectStart() {
	stream.indention += stream.cfg.indentionStep
	stream.writeByte('{')
	stream.writeIndention(0)
}


func (stream *Stream) WriteObjectField(field string) {
	stream.WriteString(field)
	if stream.indention > 0 {
		stream.writeTwoBytes(':', ' ')
	} else {
		stream.writeByte(':')
	}
}


func (stream *Stream) WriteObjectEnd() {
	stream.writeIndention(stream.cfg.indentionStep)
	stream.indention -= stream.cfg.indentionStep
	stream.writeByte('}')
}


func (stream *Stream) WriteEmptyObject() {
	stream.writeByte('{')
	stream.writeByte('}')
}


func (stream *Stream) WriteMore() {
	stream.writeByte(',')
	stream.writeIndention(0)
}


func (stream *Stream) WriteArrayStart() {
	stream.indention += stream.cfg.indentionStep
	stream.writeByte('[')
	stream.writeIndention(0)
}


func (stream *Stream) WriteEmptyArray() {
	stream.writeTwoBytes('[', ']')
}


func (stream *Stream) WriteArrayEnd() {
	stream.writeIndention(stream.cfg.indentionStep)
	stream.indention -= stream.cfg.indentionStep
	stream.writeByte(']')
}

func (stream *Stream) writeIndention(delta int) {
	if stream.indention == 0 {
		return
	}
	stream.writeByte('\n')
	toWrite := stream.indention - delta
	for i := 0; i < toWrite; i++ {
		stream.buf = append(stream.buf, ' ')
	}
}
