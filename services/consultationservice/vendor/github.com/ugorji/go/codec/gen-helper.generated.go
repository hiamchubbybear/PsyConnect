






package codec

import (
	"encoding"
	"reflect"
)


const GenVersion = 28








func GenHelper() (g genHelper) { return }

type genHelper struct{}

func (genHelper) Encoder(e *Encoder) (ge genHelperEncoder, ee genHelperEncDriver) {
	ge = genHelperEncoder{e: e}
	ee = genHelperEncDriver{encDriver: e.e}
	return
}

func (genHelper) Decoder(d *Decoder) (gd genHelperDecoder, dd genHelperDecDriver) {
	gd = genHelperDecoder{d: d}
	dd = genHelperDecDriver{decDriver: d.d}
	return
}

type genHelperEncDriver struct {
	encDriver
}

type genHelperDecDriver struct {
	decDriver
}


type genHelperEncoder struct {
	M mustHdl
	F fastpathT
	e *Encoder
}


type genHelperDecoder struct {
	C checkOverflow
	F fastpathT
	d *Decoder
}


func (f genHelperEncoder) EncBasicHandle() *BasicHandle {
	return f.e.h
}


func (f genHelperEncoder) EncWr() *encWr {
	return f.e.w()
}


func (f genHelperEncoder) EncBinary() bool {
	return f.e.be 
}


func (f genHelperEncoder) IsJSONHandle() bool {
	return f.e.js
}


func (f genHelperEncoder) EncFallback(iv interface{}) {
	
	f.e.encodeValue(reflect.ValueOf(iv), nil)
}


func (f genHelperEncoder) EncTextMarshal(iv encoding.TextMarshaler) {
	bs, fnerr := iv.MarshalText()
	f.e.marshalUtf8(bs, fnerr)
}


func (f genHelperEncoder) EncJSONMarshal(iv jsonMarshaler) {
	bs, fnerr := iv.MarshalJSON()
	f.e.marshalAsis(bs, fnerr)
}


func (f genHelperEncoder) EncBinaryMarshal(iv encoding.BinaryMarshaler) {
	bs, fnerr := iv.MarshalBinary()
	f.e.marshalRaw(bs, fnerr)
}


func (f genHelperEncoder) EncRaw(iv Raw) { f.e.rawBytes(iv) }


func (f genHelperEncoder) Extension(v interface{}) (xfn *extTypeTagFn) {
	return f.e.h.getExtForI(v)
}


func (f genHelperEncoder) EncExtension(v interface{}, xfFn *extTypeTagFn) {
	f.e.e.EncodeExt(v, xfFn.rt, xfFn.tag, xfFn.ext)
}


func (f genHelperEncoder) EncWriteMapStart(length int) { f.e.mapStart(length) }


func (f genHelperEncoder) EncWriteMapEnd() { f.e.mapEnd() }


func (f genHelperEncoder) EncWriteArrayStart(length int) { f.e.arrayStart(length) }


func (f genHelperEncoder) EncWriteArrayEnd() { f.e.arrayEnd() }


func (f genHelperEncoder) EncWriteArrayElem() { f.e.arrayElem() }


func (f genHelperEncoder) EncWriteMapElemKey() { f.e.mapElemKey() }


func (f genHelperEncoder) EncWriteMapElemValue() { f.e.mapElemValue() }


func (f genHelperEncoder) EncEncodeComplex64(v complex64) { f.e.encodeComplex64(v) }


func (f genHelperEncoder) EncEncodeComplex128(v complex128) { f.e.encodeComplex128(v) }


func (f genHelperEncoder) EncEncode(v interface{}) { f.e.encode(v) }


func (f genHelperEncoder) EncFnGivenAddr(v interface{}) *codecFn {
	return f.e.h.fn(reflect.TypeOf(v).Elem())
}


func (f genHelperEncoder) EncEncodeNumBoolStrKindGivenAddr(v interface{}, encFn *codecFn) {
	f.e.encodeValueNonNil(reflect.ValueOf(v).Elem(), encFn)
}


func (f genHelperEncoder) EncEncodeMapNonNil(v interface{}) {
	if skipFastpathTypeSwitchInDirectCall || !fastpathEncodeTypeSwitch(v, f.e) {
		f.e.encodeValueNonNil(reflect.ValueOf(v), nil)
	}
}




func (f genHelperDecoder) DecBasicHandle() *BasicHandle {
	return f.d.h
}


func (f genHelperDecoder) DecBinary() bool {
	return f.d.be 
}


func (f genHelperDecoder) DecSwallow() { f.d.swallow() }







func (f genHelperDecoder) DecScratchArrayBuffer() *[decScratchByteArrayLen]byte {
	return &f.d.b
}


func (f genHelperDecoder) DecFallback(iv interface{}, chkPtr bool) {
	rv := reflect.ValueOf(iv)
	if chkPtr {
		if x, _ := isDecodeable(rv); !x {
			f.d.haltAsNotDecodeable(rv)
		}
	}
	f.d.decodeValue(rv, nil)
}


func (f genHelperDecoder) DecSliceHelperStart() (decSliceHelper, int) {
	return f.d.decSliceHelperStart()
}


func (f genHelperDecoder) DecStructFieldNotFound(index int, name string) {
	f.d.structFieldNotFound(index, name)
}


func (f genHelperDecoder) DecArrayCannotExpand(sliceLen, streamLen int) {
	f.d.arrayCannotExpand(sliceLen, streamLen)
}


func (f genHelperDecoder) DecTextUnmarshal(tm encoding.TextUnmarshaler) {
	halt.onerror(tm.UnmarshalText(f.d.d.DecodeStringAsBytes()))
}


func (f genHelperDecoder) DecJSONUnmarshal(tm jsonUnmarshaler) {
	f.d.jsonUnmarshalV(tm)
}


func (f genHelperDecoder) DecBinaryUnmarshal(bm encoding.BinaryUnmarshaler) {
	halt.onerror(bm.UnmarshalBinary(f.d.d.DecodeBytes(nil)))
}


func (f genHelperDecoder) DecRaw() []byte { return f.d.rawBytes() }


func (f genHelperDecoder) IsJSONHandle() bool {
	return f.d.js
}


func (f genHelperDecoder) Extension(v interface{}) (xfn *extTypeTagFn) {
	return f.d.h.getExtForI(v)
}


func (f genHelperDecoder) DecExtension(v interface{}, xfFn *extTypeTagFn) {
	f.d.d.DecodeExt(v, xfFn.rt, xfFn.tag, xfFn.ext)
}


func (f genHelperDecoder) DecInferLen(clen, maxlen, unit int) (rvlen int) {
	return decInferLen(clen, maxlen, unit)
}


func (f genHelperDecoder) DecReadMapStart() int { return f.d.mapStart(f.d.d.ReadMapStart()) }


func (f genHelperDecoder) DecReadMapEnd() { f.d.mapEnd() }


func (f genHelperDecoder) DecReadArrayStart() int { return f.d.arrayStart(f.d.d.ReadArrayStart()) }


func (f genHelperDecoder) DecReadArrayEnd() { f.d.arrayEnd() }


func (f genHelperDecoder) DecReadArrayElem() { f.d.arrayElem() }


func (f genHelperDecoder) DecReadMapElemKey() { f.d.mapElemKey() }


func (f genHelperDecoder) DecReadMapElemValue() { f.d.mapElemValue() }


func (f genHelperDecoder) DecDecodeFloat32() float32 { return f.d.decodeFloat32() }


func (f genHelperDecoder) DecStringZC(v []byte) string { return f.d.stringZC(v) }


func (f genHelperDecoder) DecodeBytesInto(v []byte) []byte { return f.d.decodeBytesInto(v) }


func (f genHelperDecoder) DecContainerNext(j, containerLen int, hasLen bool) bool {
	
	
	if hasLen {
		return j < containerLen
	}
	return !f.d.checkBreak()
}
