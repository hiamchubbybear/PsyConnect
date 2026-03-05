package kafka

import (
	"bufio"
)



type saslHandshakeRequestV0 struct {
	
	Mechanism string
}

func (t saslHandshakeRequestV0) size() int32 {
	return sizeofString(t.Mechanism)
}

func (t *saslHandshakeRequestV0) readFrom(r *bufio.Reader, sz int) (remain int, err error) {
	return readString(r, sz, &t.Mechanism)
}

func (t saslHandshakeRequestV0) writeTo(wb *writeBuffer) {
	wb.writeString(t.Mechanism)
}



type saslHandshakeResponseV0 struct {
	
	ErrorCode int16

	
	EnabledMechanisms []string
}

func (t saslHandshakeResponseV0) size() int32 {
	return sizeofInt16(t.ErrorCode) + sizeofStringArray(t.EnabledMechanisms)
}

func (t saslHandshakeResponseV0) writeTo(wb *writeBuffer) {
	wb.writeInt16(t.ErrorCode)
	wb.writeStringArray(t.EnabledMechanisms)
}

func (t *saslHandshakeResponseV0) readFrom(r *bufio.Reader, sz int) (remain int, err error) {
	if remain, err = readInt16(r, sz, &t.ErrorCode); err != nil {
		return
	}
	if remain, err = readStringArray(r, remain, &t.EnabledMechanisms); err != nil {
		return
	}
	return
}
