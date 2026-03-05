





package bsonrw

import (
	"errors"
	"fmt"
	"io"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)






type Copier struct{}






func NewCopier() Copier {
	return Copier{}
}





func CopyDocument(dst ValueWriter, src ValueReader) error {
	return Copier{}.CopyDocument(dst, src)
}





func (c Copier) CopyDocument(dst ValueWriter, src ValueReader) error {
	dr, err := src.ReadDocument()
	if err != nil {
		return err
	}

	dw, err := dst.WriteDocument()
	if err != nil {
		return err
	}

	return c.copyDocumentCore(dw, dr)
}






func (c Copier) CopyArrayFromBytes(dst ValueWriter, src []byte) error {
	aw, err := dst.WriteArray()
	if err != nil {
		return err
	}

	err = c.CopyBytesToArrayWriter(aw, src)
	if err != nil {
		return err
	}

	return aw.WriteArrayEnd()
}






func (c Copier) CopyDocumentFromBytes(dst ValueWriter, src []byte) error {
	dw, err := dst.WriteDocument()
	if err != nil {
		return err
	}

	err = c.CopyBytesToDocumentWriter(dw, src)
	if err != nil {
		return err
	}

	return dw.WriteDocumentEnd()
}

type writeElementFn func(key string) (ValueWriter, error)






func (c Copier) CopyBytesToArrayWriter(dst ArrayWriter, src []byte) error {
	wef := func(_ string) (ValueWriter, error) {
		return dst.WriteArrayElement()
	}

	return c.copyBytesToValueWriter(src, wef)
}






func (c Copier) CopyBytesToDocumentWriter(dst DocumentWriter, src []byte) error {
	wef := func(key string) (ValueWriter, error) {
		return dst.WriteDocumentElement(key)
	}

	return c.copyBytesToValueWriter(src, wef)
}

func (c Copier) copyBytesToValueWriter(src []byte, wef writeElementFn) error {
	
	length, rem, ok := bsoncore.ReadLength(src)
	if !ok {
		return fmt.Errorf("couldn't read length from src, not enough bytes. length=%d", len(src))
	}
	if len(src) < int(length) {
		return fmt.Errorf("length read exceeds number of bytes available. length=%d bytes=%d", len(src), length)
	}
	rem = rem[:length-4]

	var t bsontype.Type
	var key string
	var val bsoncore.Value
	for {
		t, rem, ok = bsoncore.ReadType(rem)
		if !ok {
			return io.EOF
		}
		if t == bsontype.Type(0) {
			if len(rem) != 0 {
				return fmt.Errorf("document end byte found before end of document. remaining bytes=%v", rem)
			}
			break
		}

		key, rem, ok = bsoncore.ReadKey(rem)
		if !ok {
			return fmt.Errorf("invalid key found. remaining bytes=%v", rem)
		}

		
		vw, err := wef(key)
		if err != nil {
			return err
		}

		val, rem, ok = bsoncore.ReadValue(rem, t)
		if !ok {
			return fmt.Errorf("not enough bytes available to read type. bytes=%d type=%s", len(rem), t)
		}
		err = c.CopyValueFromBytes(vw, t, val.Data)
		if err != nil {
			return err
		}
	}
	return nil
}






func (c Copier) CopyDocumentToBytes(src ValueReader) ([]byte, error) {
	return c.AppendDocumentBytes(nil, src)
}






func (c Copier) AppendDocumentBytes(dst []byte, src ValueReader) ([]byte, error) {
	if br, ok := src.(BytesReader); ok {
		_, dst, err := br.ReadValueBytes(dst)
		return dst, err
	}

	vw := vwPool.Get().(*valueWriter)
	defer putValueWriter(vw)

	vw.reset(dst)

	err := c.CopyDocument(vw, src)
	dst = vw.buf
	return dst, err
}





func (c Copier) AppendArrayBytes(dst []byte, src ValueReader) ([]byte, error) {
	if br, ok := src.(BytesReader); ok {
		_, dst, err := br.ReadValueBytes(dst)
		return dst, err
	}

	vw := vwPool.Get().(*valueWriter)
	defer putValueWriter(vw)

	vw.reset(dst)

	err := c.copyArray(vw, src)
	dst = vw.buf
	return dst, err
}




func (c Copier) CopyValueFromBytes(dst ValueWriter, t bsontype.Type, src []byte) error {
	if wvb, ok := dst.(BytesWriter); ok {
		return wvb.WriteValueBytes(t, src)
	}

	vr := vrPool.Get().(*valueReader)
	defer vrPool.Put(vr)

	vr.reset(src)
	vr.pushElement(t)

	return c.CopyValue(dst, vr)
}





func (c Copier) CopyValueToBytes(src ValueReader) (bsontype.Type, []byte, error) {
	return c.AppendValueBytes(nil, src)
}






func (c Copier) AppendValueBytes(dst []byte, src ValueReader) (bsontype.Type, []byte, error) {
	if br, ok := src.(BytesReader); ok {
		return br.ReadValueBytes(dst)
	}

	vw := vwPool.Get().(*valueWriter)
	defer putValueWriter(vw)

	start := len(dst)

	vw.reset(dst)
	vw.push(mElement)

	err := c.CopyValue(vw, src)
	if err != nil {
		return 0, dst, err
	}

	return bsontype.Type(vw.buf[start]), vw.buf[start+2:], nil
}





func (c Copier) CopyValue(dst ValueWriter, src ValueReader) error {
	var err error
	switch src.Type() {
	case bsontype.Double:
		var f64 float64
		f64, err = src.ReadDouble()
		if err != nil {
			break
		}
		err = dst.WriteDouble(f64)
	case bsontype.String:
		var str string
		str, err = src.ReadString()
		if err != nil {
			return err
		}
		err = dst.WriteString(str)
	case bsontype.EmbeddedDocument:
		err = c.CopyDocument(dst, src)
	case bsontype.Array:
		err = c.copyArray(dst, src)
	case bsontype.Binary:
		var data []byte
		var subtype byte
		data, subtype, err = src.ReadBinary()
		if err != nil {
			break
		}
		err = dst.WriteBinaryWithSubtype(data, subtype)
	case bsontype.Undefined:
		err = src.ReadUndefined()
		if err != nil {
			break
		}
		err = dst.WriteUndefined()
	case bsontype.ObjectID:
		var oid primitive.ObjectID
		oid, err = src.ReadObjectID()
		if err != nil {
			break
		}
		err = dst.WriteObjectID(oid)
	case bsontype.Boolean:
		var b bool
		b, err = src.ReadBoolean()
		if err != nil {
			break
		}
		err = dst.WriteBoolean(b)
	case bsontype.DateTime:
		var dt int64
		dt, err = src.ReadDateTime()
		if err != nil {
			break
		}
		err = dst.WriteDateTime(dt)
	case bsontype.Null:
		err = src.ReadNull()
		if err != nil {
			break
		}
		err = dst.WriteNull()
	case bsontype.Regex:
		var pattern, options string
		pattern, options, err = src.ReadRegex()
		if err != nil {
			break
		}
		err = dst.WriteRegex(pattern, options)
	case bsontype.DBPointer:
		var ns string
		var pointer primitive.ObjectID
		ns, pointer, err = src.ReadDBPointer()
		if err != nil {
			break
		}
		err = dst.WriteDBPointer(ns, pointer)
	case bsontype.JavaScript:
		var js string
		js, err = src.ReadJavascript()
		if err != nil {
			break
		}
		err = dst.WriteJavascript(js)
	case bsontype.Symbol:
		var symbol string
		symbol, err = src.ReadSymbol()
		if err != nil {
			break
		}
		err = dst.WriteSymbol(symbol)
	case bsontype.CodeWithScope:
		var code string
		var srcScope DocumentReader
		code, srcScope, err = src.ReadCodeWithScope()
		if err != nil {
			break
		}

		var dstScope DocumentWriter
		dstScope, err = dst.WriteCodeWithScope(code)
		if err != nil {
			break
		}
		err = c.copyDocumentCore(dstScope, srcScope)
	case bsontype.Int32:
		var i32 int32
		i32, err = src.ReadInt32()
		if err != nil {
			break
		}
		err = dst.WriteInt32(i32)
	case bsontype.Timestamp:
		var t, i uint32
		t, i, err = src.ReadTimestamp()
		if err != nil {
			break
		}
		err = dst.WriteTimestamp(t, i)
	case bsontype.Int64:
		var i64 int64
		i64, err = src.ReadInt64()
		if err != nil {
			break
		}
		err = dst.WriteInt64(i64)
	case bsontype.Decimal128:
		var d128 primitive.Decimal128
		d128, err = src.ReadDecimal128()
		if err != nil {
			break
		}
		err = dst.WriteDecimal128(d128)
	case bsontype.MinKey:
		err = src.ReadMinKey()
		if err != nil {
			break
		}
		err = dst.WriteMinKey()
	case bsontype.MaxKey:
		err = src.ReadMaxKey()
		if err != nil {
			break
		}
		err = dst.WriteMaxKey()
	default:
		err = fmt.Errorf("Cannot copy unknown BSON type %s", src.Type())
	}

	return err
}

func (c Copier) copyArray(dst ValueWriter, src ValueReader) error {
	ar, err := src.ReadArray()
	if err != nil {
		return err
	}

	aw, err := dst.WriteArray()
	if err != nil {
		return err
	}

	for {
		vr, err := ar.ReadValue()
		if errors.Is(err, ErrEOA) {
			break
		}
		if err != nil {
			return err
		}

		vw, err := aw.WriteArrayElement()
		if err != nil {
			return err
		}

		err = c.CopyValue(vw, vr)
		if err != nil {
			return err
		}
	}

	return aw.WriteArrayEnd()
}

func (c Copier) copyDocumentCore(dw DocumentWriter, dr DocumentReader) error {
	for {
		key, vr, err := dr.ReadElement()
		if errors.Is(err, ErrEOD) {
			break
		}
		if err != nil {
			return err
		}

		vw, err := dw.WriteDocumentElement(key)
		if err != nil {
			return err
		}

		err = c.CopyValue(vw, vr)
		if err != nil {
			return err
		}
	}

	return dw.WriteDocumentEnd()
}
