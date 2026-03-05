


package codec

import (
	"encoding"
	"errors"
	"io"
	"math"
	"reflect"
	"strconv"
	"time"
)

const msgBadDesc = "unrecognized descriptor byte"

const (
	decDefMaxDepth         = 1024                
	decDefChanCap          = 64                  
	decScratchByteArrayLen = (8 + 2 + 2 + 1) * 8 

	

	
	
	
	
	
	
	
	
	
	
	
	decFailNonEmptyIntf = false

	
	
	
	
	
	
	decUseTransient = true
)

var (
	errNeedMapOrArrayDecodeToStruct = errors.New("only encoded map or array can decode into struct")
	errCannotDecodeIntoNil          = errors.New("cannot decode into nil")

	errExpandSliceCannotChange = errors.New("expand slice: cannot change")

	errDecoderNotInitialized = errors.New("Decoder not initialized")

	errDecUnreadByteNothingToRead   = errors.New("cannot unread - nothing has been read")
	errDecUnreadByteLastByteNotRead = errors.New("cannot unread - last byte has not been read")
	errDecUnreadByteUnknown         = errors.New("cannot unread - reason unknown")
	errMaxDepthExceeded             = errors.New("maximum decoding depth exceeded")
)



type decByteState uint8

const (
	decByteStateNone     decByteState = iota
	decByteStateZerocopy              
	decByteStateReuseBuf              
	
)

type decNotDecodeableReason uint8

const (
	decNotDecodeableReasonUnknown decNotDecodeableReason = iota
	decNotDecodeableReasonBadKind
	decNotDecodeableReasonNonAddrValue
	decNotDecodeableReasonNilReference
)

type decDriver interface {
	
	CheckBreak() bool

	
	
	
	
	TryNil() bool

	
	
	
	
	
	ContainerType() (vt valueType)

	
	
	
	
	
	
	
	
	
	
	
	DecodeNaked()

	DecodeInt64() (i int64)
	DecodeUint64() (ui uint64)

	DecodeFloat64() (f float64)
	DecodeBool() (b bool)

	
	
	
	
	
	
	DecodeStringAsBytes() (v []byte)

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	DecodeBytes(in []byte) (out []byte)
	

	
	DecodeExt(v interface{}, basetype reflect.Type, xtag uint64, ext Ext)
	

	DecodeTime() (t time.Time)

	
	
	
	ReadArrayStart() int

	
	
	
	ReadMapStart() int

	reset()

	

	
	
	
	
	
	
	
	
	
	nextValueBytes(start []byte) []byte

	
	descBd() string

	decoder() *Decoder

	driverStateManager
	decNegintPosintFloatNumber
}

type decDriverContainerTracker interface {
	ReadArrayElem()
	ReadMapElemKey()
	ReadMapElemValue()
	ReadArrayEnd()
	ReadMapEnd()
}

type decNegintPosintFloatNumber interface {
	decInteger() (ui uint64, neg, ok bool)
	decFloat() (f float64, ok bool)
}

type decDriverNoopNumberHelper struct{}

func (x decDriverNoopNumberHelper) decInteger() (ui uint64, neg, ok bool) {
	panic("decInteger unsupported")
}
func (x decDriverNoopNumberHelper) decFloat() (f float64, ok bool) { panic("decFloat unsupported") }

type decDriverNoopContainerReader struct{}



func (x decDriverNoopContainerReader) ReadArrayEnd()        {}
func (x decDriverNoopContainerReader) ReadMapEnd()          {}
func (x decDriverNoopContainerReader) CheckBreak() (v bool) { return }


type DecodeOptions struct {
	
	
	
	MapType reflect.Type

	
	
	SliceType reflect.Type

	
	
	
	
	
	
	
	MaxInitLen int

	
	
	
	ReaderBufferSize int

	
	
	MaxDepth int16

	
	
	ErrorIfNoField bool

	
	
	
	ErrorIfNoArrayExpand bool

	
	SignedInteger bool

	
	
	
	
	
	
	
	
	
	
	
	
	MapValueReset bool

	
	
	
	SliceElementReset bool

	
	
	
	
	
	
	
	
	
	
	
	InterfaceReset bool

	
	
	
	
	
	
	
	
	
	
	InternString bool

	
	
	
	
	
	
	
	
	PreferArrayOverSlice bool

	
	
	
	
	
	
	
	
	DeleteOnNilMapValue bool

	
	
	RawToString bool

	
	
	
	
	
	
	
	
	
	
	
	ZeroCopy bool

	
	
	
	
	PreferPointerForStructOrArray bool

	
	
	
	
	ValidateUnicode bool
}



func (d *Decoder) rawExt(f *codecFnInfo, rv reflect.Value) {
	d.d.DecodeExt(rv2i(rv), f.ti.rt, 0, nil)
}

func (d *Decoder) ext(f *codecFnInfo, rv reflect.Value) {
	d.d.DecodeExt(rv2i(rv), f.ti.rt, f.xfTag, f.xfFn)
}

func (d *Decoder) selferUnmarshal(f *codecFnInfo, rv reflect.Value) {
	rv2i(rv).(Selfer).CodecDecodeSelf(d)
}

func (d *Decoder) binaryUnmarshal(f *codecFnInfo, rv reflect.Value) {
	bm := rv2i(rv).(encoding.BinaryUnmarshaler)
	xbs := d.d.DecodeBytes(nil)
	fnerr := bm.UnmarshalBinary(xbs)
	d.onerror(fnerr)
}

func (d *Decoder) textUnmarshal(f *codecFnInfo, rv reflect.Value) {
	tm := rv2i(rv).(encoding.TextUnmarshaler)
	fnerr := tm.UnmarshalText(d.d.DecodeStringAsBytes())
	d.onerror(fnerr)
}

func (d *Decoder) jsonUnmarshal(f *codecFnInfo, rv reflect.Value) {
	d.jsonUnmarshalV(rv2i(rv).(jsonUnmarshaler))
}

func (d *Decoder) jsonUnmarshalV(tm jsonUnmarshaler) {
	
	var bs0 = []byte{}
	if !d.bytes {
		bs0 = d.blist.get(256)
	}
	bs := d.d.nextValueBytes(bs0)
	fnerr := tm.UnmarshalJSON(bs)
	if !d.bytes {
		d.blist.put(bs)
		if !byteSliceSameData(bs0, bs) {
			d.blist.put(bs0)
		}
	}
	d.onerror(fnerr)
}

func (d *Decoder) kErr(f *codecFnInfo, rv reflect.Value) {
	d.errorf("no decoding function defined for kind %v", rv.Kind())
}

func (d *Decoder) raw(f *codecFnInfo, rv reflect.Value) {
	rvSetBytes(rv, d.rawBytes())
}

func (d *Decoder) kString(f *codecFnInfo, rv reflect.Value) {
	rvSetString(rv, d.stringZC(d.d.DecodeStringAsBytes()))
}

func (d *Decoder) kBool(f *codecFnInfo, rv reflect.Value) {
	rvSetBool(rv, d.d.DecodeBool())
}

func (d *Decoder) kTime(f *codecFnInfo, rv reflect.Value) {
	rvSetTime(rv, d.d.DecodeTime())
}

func (d *Decoder) kFloat32(f *codecFnInfo, rv reflect.Value) {
	rvSetFloat32(rv, d.decodeFloat32())
}

func (d *Decoder) kFloat64(f *codecFnInfo, rv reflect.Value) {
	rvSetFloat64(rv, d.d.DecodeFloat64())
}

func (d *Decoder) kComplex64(f *codecFnInfo, rv reflect.Value) {
	rvSetComplex64(rv, complex(d.decodeFloat32(), 0))
}

func (d *Decoder) kComplex128(f *codecFnInfo, rv reflect.Value) {
	rvSetComplex128(rv, complex(d.d.DecodeFloat64(), 0))
}

func (d *Decoder) kInt(f *codecFnInfo, rv reflect.Value) {
	rvSetInt(rv, int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize)))
}

func (d *Decoder) kInt8(f *codecFnInfo, rv reflect.Value) {
	rvSetInt8(rv, int8(chkOvf.IntV(d.d.DecodeInt64(), 8)))
}

func (d *Decoder) kInt16(f *codecFnInfo, rv reflect.Value) {
	rvSetInt16(rv, int16(chkOvf.IntV(d.d.DecodeInt64(), 16)))
}

func (d *Decoder) kInt32(f *codecFnInfo, rv reflect.Value) {
	rvSetInt32(rv, int32(chkOvf.IntV(d.d.DecodeInt64(), 32)))
}

func (d *Decoder) kInt64(f *codecFnInfo, rv reflect.Value) {
	rvSetInt64(rv, d.d.DecodeInt64())
}

func (d *Decoder) kUint(f *codecFnInfo, rv reflect.Value) {
	rvSetUint(rv, uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize)))
}

func (d *Decoder) kUintptr(f *codecFnInfo, rv reflect.Value) {
	rvSetUintptr(rv, uintptr(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize)))
}

func (d *Decoder) kUint8(f *codecFnInfo, rv reflect.Value) {
	rvSetUint8(rv, uint8(chkOvf.UintV(d.d.DecodeUint64(), 8)))
}

func (d *Decoder) kUint16(f *codecFnInfo, rv reflect.Value) {
	rvSetUint16(rv, uint16(chkOvf.UintV(d.d.DecodeUint64(), 16)))
}

func (d *Decoder) kUint32(f *codecFnInfo, rv reflect.Value) {
	rvSetUint32(rv, uint32(chkOvf.UintV(d.d.DecodeUint64(), 32)))
}

func (d *Decoder) kUint64(f *codecFnInfo, rv reflect.Value) {
	rvSetUint64(rv, d.d.DecodeUint64())
}

func (d *Decoder) kInterfaceNaked(f *codecFnInfo) (rvn reflect.Value) {
	
	
	
	n := d.naked()
	d.d.DecodeNaked()

	
	
	
	
	
	
	
	if decFailNonEmptyIntf && f.ti.numMeth > 0 {
		d.errorf("cannot decode non-nil codec value into nil %v (%v methods)", f.ti.rt, f.ti.numMeth)
	}
	switch n.v {
	case valueTypeMap:
		mtid := d.mtid
		if mtid == 0 {
			if d.jsms { 
				mtid = mapStrIntfTypId 
			} else {
				mtid = mapIntfIntfTypId
			}
		}
		if mtid == mapStrIntfTypId {
			var v2 map[string]interface{}
			d.decode(&v2)
			rvn = rv4iptr(&v2).Elem()
		} else if mtid == mapIntfIntfTypId {
			var v2 map[interface{}]interface{}
			d.decode(&v2)
			rvn = rv4iptr(&v2).Elem()
		} else if d.mtr {
			rvn = reflect.New(d.h.MapType)
			d.decode(rv2i(rvn))
			rvn = rvn.Elem()
		} else {
			rvn = rvZeroAddrK(d.h.MapType, reflect.Map)
			d.decodeValue(rvn, nil)
		}
	case valueTypeArray:
		if d.stid == 0 || d.stid == intfSliceTypId {
			var v2 []interface{}
			d.decode(&v2)
			rvn = rv4iptr(&v2).Elem()
		} else if d.str {
			rvn = reflect.New(d.h.SliceType)
			d.decode(rv2i(rvn))
			rvn = rvn.Elem()
		} else {
			rvn = rvZeroAddrK(d.h.SliceType, reflect.Slice)
			d.decodeValue(rvn, nil)
		}
		if reflectArrayOfSupported && d.h.PreferArrayOverSlice {
			rvn = rvGetArray4Slice(rvn)
		}
	case valueTypeExt:
		tag, bytes := n.u, n.l 
		bfn := d.h.getExtForTag(tag)
		var re = RawExt{Tag: tag}
		if bytes == nil {
			
			
			if bfn == nil {
				d.decode(&re.Value)
				rvn = rv4iptr(&re).Elem()
			} else {
				if bfn.ext == SelfExt {
					rvn = rvZeroAddrK(bfn.rt, bfn.rt.Kind())
					d.decodeValue(rvn, d.h.fnNoExt(bfn.rt))
				} else {
					rvn = reflect.New(bfn.rt)
					d.interfaceExtConvertAndDecode(rv2i(rvn), bfn.ext)
					rvn = rvn.Elem()
				}
			}
		} else {
			
			if bfn == nil {
				re.setData(bytes, false)
				rvn = rv4iptr(&re).Elem()
			} else {
				rvn = reflect.New(bfn.rt)
				if bfn.ext == SelfExt {
					d.sideDecode(rv2i(rvn), bfn.rt, bytes)
				} else {
					bfn.ext.ReadExt(rv2i(rvn), bytes)
				}
				rvn = rvn.Elem()
			}
		}
		
		if d.h.PreferPointerForStructOrArray && rvn.CanAddr() {
			if rk := rvn.Kind(); rk == reflect.Array || rk == reflect.Struct {
				rvn = rvn.Addr()
			}
		}
	case valueTypeNil:
		
		
	case valueTypeInt:
		rvn = n.ri()
	case valueTypeUint:
		rvn = n.ru()
	case valueTypeFloat:
		rvn = n.rf()
	case valueTypeBool:
		rvn = n.rb()
	case valueTypeString, valueTypeSymbol:
		rvn = n.rs()
	case valueTypeBytes:
		rvn = n.rl()
	case valueTypeTime:
		rvn = n.rt()
	default:
		halt.errorf("kInterfaceNaked: unexpected valueType: %d", n.v)
	}
	return
}

func (d *Decoder) kInterface(f *codecFnInfo, rv reflect.Value) {
	
	
	
	
	
	
	
	

	isnilrv := rvIsNil(rv)

	var rvn reflect.Value

	if d.h.InterfaceReset {
		
		rvn = d.h.intf2impl(f.ti.rtid)
		if !rvn.IsValid() {
			rvn = d.kInterfaceNaked(f)
			if rvn.IsValid() {
				rvSetIntf(rv, rvn)
			} else if !isnilrv {
				decSetNonNilRV2Zero4Intf(rv)
			}
			return
		}
	} else if isnilrv {
		
		rvn = d.h.intf2impl(f.ti.rtid)
		if !rvn.IsValid() {
			rvn = d.kInterfaceNaked(f)
			if rvn.IsValid() {
				rvSetIntf(rv, rvn)
			}
			return
		}
	} else {
		
		rvn = rv.Elem()
	}

	

	canDecode, _ := isDecodeable(rvn)

	
	
	

	if !canDecode {
		rvn2 := d.oneShotAddrRV(rvn.Type(), rvn.Kind())
		rvSetDirect(rvn2, rvn)
		rvn = rvn2
	}

	d.decodeValue(rvn, nil)
	rvSetIntf(rv, rvn)
}

func decStructFieldKeyNotString(dd decDriver, keyType valueType, b *[decScratchByteArrayLen]byte) (rvkencname []byte) {
	if keyType == valueTypeInt {
		rvkencname = strconv.AppendInt(b[:0], dd.DecodeInt64(), 10)
	} else if keyType == valueTypeUint {
		rvkencname = strconv.AppendUint(b[:0], dd.DecodeUint64(), 10)
	} else if keyType == valueTypeFloat {
		rvkencname = strconv.AppendFloat(b[:0], dd.DecodeFloat64(), 'f', -1, 64)
	} else {
		halt.errorf("invalid struct key type: %v", keyType)
	}
	return
}

func (d *Decoder) kStructField(si *structFieldInfo, rv reflect.Value) {
	if d.d.TryNil() {
		if rv = si.path.field(rv); rv.IsValid() {
			decSetNonNilRV2Zero(rv)
		}
		return
	}
	d.decodeValueNoCheckNil(si.path.fieldAlloc(rv), nil)
}

func (d *Decoder) kStruct(f *codecFnInfo, rv reflect.Value) {
	ctyp := d.d.ContainerType()
	ti := f.ti
	var mf MissingFielder
	if ti.flagMissingFielder {
		mf = rv2i(rv).(MissingFielder)
	} else if ti.flagMissingFielderPtr {
		mf = rv2i(rvAddr(rv, ti.ptr)).(MissingFielder)
	}
	if ctyp == valueTypeMap {
		containerLen := d.mapStart(d.d.ReadMapStart())
		if containerLen == 0 {
			d.mapEnd()
			return
		}
		hasLen := containerLen >= 0
		var name2 []byte
		if mf != nil {
			var namearr2 [16]byte
			name2 = namearr2[:0]
		}
		var rvkencname []byte
		for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
			d.mapElemKey()
			if ti.keyType == valueTypeString {
				rvkencname = d.d.DecodeStringAsBytes()
			} else {
				rvkencname = decStructFieldKeyNotString(d.d, ti.keyType, &d.b)
			}
			d.mapElemValue()
			if si := ti.siForEncName(rvkencname); si != nil {
				d.kStructField(si, rv)
			} else if mf != nil {
				
				name2 = append(name2[:0], rvkencname...)
				var f interface{}
				d.decode(&f)
				if !mf.CodecMissingField(name2, f) && d.h.ErrorIfNoField {
					d.errorf("no matching struct field when decoding stream map with key: %s ", stringView(name2))
				}
			} else {
				d.structFieldNotFound(-1, stringView(rvkencname))
			}
		}
		d.mapEnd()
	} else if ctyp == valueTypeArray {
		containerLen := d.arrayStart(d.d.ReadArrayStart())
		if containerLen == 0 {
			d.arrayEnd()
			return
		}
		
		
		tisfi := ti.sfi.source()
		hasLen := containerLen >= 0

		
		
		
		for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
			d.arrayElem()
			if j < len(tisfi) {
				d.kStructField(tisfi[j], rv)
			} else {
				d.structFieldNotFound(j, "")
			}
		}

		d.arrayEnd()
	} else {
		d.onerror(errNeedMapOrArrayDecodeToStruct)
	}
}

func (d *Decoder) kSlice(f *codecFnInfo, rv reflect.Value) {
	
	

	

	ti := f.ti
	rvCanset := rv.CanSet()

	ctyp := d.d.ContainerType()
	if ctyp == valueTypeBytes || ctyp == valueTypeString {
		
		if !(ti.rtid == uint8SliceTypId || ti.elemkind == uint8(reflect.Uint8)) {
			d.errorf("bytes/string in stream must decode into slice/array of bytes, not %v", ti.rt)
		}
		rvbs := rvGetBytes(rv)
		if !rvCanset {
			
			rvbs = rvbs[:len(rvbs):len(rvbs)]
		}
		bs2 := d.decodeBytesInto(rvbs)
		
		if !(len(bs2) > 0 && len(bs2) == len(rvbs) && &bs2[0] == &rvbs[0]) {
			if rvCanset {
				rvSetBytes(rv, bs2)
			} else if len(rvbs) > 0 && len(bs2) > 0 {
				copy(rvbs, bs2)
			}
		}
		return
	}

	slh, containerLenS := d.decSliceHelperStart() 

	
	if containerLenS == 0 {
		if rvCanset {
			if rvIsNil(rv) {
				rvSetDirect(rv, rvSliceZeroCap(ti.rt))
			} else {
				rvSetSliceLen(rv, 0)
			}
		}
		slh.End()
		return
	}

	rtelem0Mut := !scalarBitset.isset(ti.elemkind)
	rtelem := ti.elem

	for k := reflect.Kind(ti.elemkind); k == reflect.Ptr; k = rtelem.Kind() {
		rtelem = rtelem.Elem()
	}

	var fn *codecFn

	var rvChanged bool

	var rv0 = rv
	var rv9 reflect.Value

	rvlen := rvLenSlice(rv)
	rvcap := rvCapSlice(rv)
	hasLen := containerLenS > 0
	if hasLen {
		if containerLenS > rvcap {
			oldRvlenGtZero := rvlen > 0
			rvlen1 := decInferLen(containerLenS, d.h.MaxInitLen, int(ti.elemsize))
			if rvlen1 == rvlen {
			} else if rvlen1 <= rvcap {
				if rvCanset {
					rvlen = rvlen1
					rvSetSliceLen(rv, rvlen)
				}
			} else if rvCanset { 
				rvlen = rvlen1
				rv, rvCanset = rvMakeSlice(rv, f.ti, rvlen, rvlen)
				rvcap = rvlen
				rvChanged = !rvCanset
			} else { 
				d.errorf("cannot decode into non-settable slice")
			}
			if rvChanged && oldRvlenGtZero && rtelem0Mut {
				rvCopySlice(rv, rv0, rtelem) 
			}
		} else if containerLenS != rvlen {
			if rvCanset {
				rvlen = containerLenS
				rvSetSliceLen(rv, rvlen)
			}
		}
	}

	
	var elemReset = d.h.SliceElementReset

	var j int

	for ; d.containerNext(j, containerLenS, hasLen); j++ {
		if j == 0 {
			if rvIsNil(rv) { 
				if rvCanset {
					rvlen = decInferLen(containerLenS, d.h.MaxInitLen, int(ti.elemsize))
					rv, rvCanset = rvMakeSlice(rv, f.ti, rvlen, rvlen)
					rvcap = rvlen
					rvChanged = !rvCanset
				} else {
					d.errorf("cannot decode into non-settable slice")
				}
			}
			if fn == nil {
				fn = d.h.fn(rtelem)
			}
		}
		
		if j >= rvlen {
			slh.ElemContainerState(j)

			
			

			if rvlen < rvcap {
				rvlen = rvcap
				if rvCanset {
					rvSetSliceLen(rv, rvlen)
				} else if rvChanged {
					rv = rvSlice(rv, rvlen)
				} else {
					d.onerror(errExpandSliceCannotChange)
				}
			} else {
				if !(rvCanset || rvChanged) {
					d.onerror(errExpandSliceCannotChange)
				}
				rv, rvcap, rvCanset = rvGrowSlice(rv, f.ti, rvcap, 1)
				rvlen = rvcap
				rvChanged = !rvCanset
			}
		} else {
			slh.ElemContainerState(j)
		}
		rv9 = rvSliceIndex(rv, j, f.ti)
		if elemReset {
			rvSetZero(rv9)
		}
		d.decodeValue(rv9, fn)
	}
	if j < rvlen {
		if rvCanset {
			rvSetSliceLen(rv, j)
		} else if rvChanged {
			rv = rvSlice(rv, j)
		}
		
	} else if j == 0 && rvIsNil(rv) {
		if rvCanset {
			rv = rvSliceZeroCap(ti.rt)
			rvCanset = false
			rvChanged = true
		}
	}
	slh.End()

	if rvChanged { 
		rvSetDirect(rv0, rv)
	}
}

func (d *Decoder) kArray(f *codecFnInfo, rv reflect.Value) {
	

	ctyp := d.d.ContainerType()
	if handleBytesWithinKArray && (ctyp == valueTypeBytes || ctyp == valueTypeString) {
		
		if f.ti.elemkind != uint8(reflect.Uint8) {
			d.errorf("bytes/string in stream can decode into array of bytes, but not %v", f.ti.rt)
		}
		rvbs := rvGetArrayBytes(rv, nil)
		bs2 := d.decodeBytesInto(rvbs)
		if !byteSliceSameData(rvbs, bs2) && len(rvbs) > 0 && len(bs2) > 0 {
			copy(rvbs, bs2)
		}
		return
	}

	slh, containerLenS := d.decSliceHelperStart() 

	
	if containerLenS == 0 {
		slh.End()
		return
	}

	rtelem := f.ti.elem
	for k := reflect.Kind(f.ti.elemkind); k == reflect.Ptr; k = rtelem.Kind() {
		rtelem = rtelem.Elem()
	}

	var fn *codecFn

	var rv9 reflect.Value

	rvlen := rv.Len() 
	hasLen := containerLenS > 0
	if hasLen && containerLenS > rvlen {
		d.errorf("cannot decode into array with length: %v, less than container length: %v", rvlen, containerLenS)
	}

	
	var elemReset = d.h.SliceElementReset

	for j := 0; d.containerNext(j, containerLenS, hasLen); j++ {
		
		if j >= rvlen {
			slh.arrayCannotExpand(hasLen, rvlen, j, containerLenS)
			return
		}

		slh.ElemContainerState(j)
		rv9 = rvArrayIndex(rv, j, f.ti)
		if elemReset {
			rvSetZero(rv9)
		}

		if fn == nil {
			fn = d.h.fn(rtelem)
		}
		d.decodeValue(rv9, fn)
	}
	slh.End()
}

func (d *Decoder) kChan(f *codecFnInfo, rv reflect.Value) {
	
	

	ti := f.ti
	if ti.chandir&uint8(reflect.SendDir) == 0 {
		d.errorf("receive-only channel cannot be decoded")
	}
	ctyp := d.d.ContainerType()
	if ctyp == valueTypeBytes || ctyp == valueTypeString {
		
		if !(ti.rtid == uint8SliceTypId || ti.elemkind == uint8(reflect.Uint8)) {
			d.errorf("bytes/string in stream must decode into slice/array of bytes, not %v", ti.rt)
		}
		bs2 := d.d.DecodeBytes(nil)
		irv := rv2i(rv)
		ch, ok := irv.(chan<- byte)
		if !ok {
			ch = irv.(chan byte)
		}
		for _, b := range bs2 {
			ch <- b
		}
		return
	}

	var rvCanset = rv.CanSet()

	
	slh, containerLenS := d.decSliceHelperStart()

	
	if containerLenS == 0 {
		if rvCanset && rvIsNil(rv) {
			rvSetDirect(rv, reflect.MakeChan(ti.rt, 0))
		}
		slh.End()
		return
	}

	rtelem := ti.elem
	useTransient := decUseTransient && ti.elemkind != byte(reflect.Ptr) && ti.tielem.flagCanTransient

	for k := reflect.Kind(ti.elemkind); k == reflect.Ptr; k = rtelem.Kind() {
		rtelem = rtelem.Elem()
	}

	var fn *codecFn

	var rvChanged bool
	var rv0 = rv
	var rv9 reflect.Value

	var rvlen int 
	hasLen := containerLenS > 0

	for j := 0; d.containerNext(j, containerLenS, hasLen); j++ {
		if j == 0 {
			if rvIsNil(rv) {
				if hasLen {
					rvlen = decInferLen(containerLenS, d.h.MaxInitLen, int(ti.elemsize))
				} else {
					rvlen = decDefChanCap
				}
				if rvCanset {
					rv = reflect.MakeChan(ti.rt, rvlen)
					rvChanged = true
				} else {
					d.errorf("cannot decode into non-settable chan")
				}
			}
			if fn == nil {
				fn = d.h.fn(rtelem)
			}
		}
		slh.ElemContainerState(j)
		if rv9.IsValid() {
			rvSetZero(rv9)
		} else if decUseTransient && useTransient {
			rv9 = d.perType.TransientAddrK(ti.elem, reflect.Kind(ti.elemkind))
		} else {
			rv9 = rvZeroAddrK(ti.elem, reflect.Kind(ti.elemkind))
		}
		if !d.d.TryNil() {
			d.decodeValueNoCheckNil(rv9, fn)
		}
		rv.Send(rv9)
	}
	slh.End()

	if rvChanged { 
		rvSetDirect(rv0, rv)
	}

}

func (d *Decoder) kMap(f *codecFnInfo, rv reflect.Value) {
	containerLen := d.mapStart(d.d.ReadMapStart())
	ti := f.ti
	if rvIsNil(rv) {
		rvlen := decInferLen(containerLen, d.h.MaxInitLen, int(ti.keysize+ti.elemsize))
		rvSetDirect(rv, makeMapReflect(ti.rt, rvlen))
	}

	if containerLen == 0 {
		d.mapEnd()
		return
	}

	ktype, vtype := ti.key, ti.elem
	ktypeId := rt2id(ktype)
	vtypeKind := reflect.Kind(ti.elemkind)
	ktypeKind := reflect.Kind(ti.keykind)
	kfast := mapKeyFastKindFor(ktypeKind)
	visindirect := mapStoresElemIndirect(uintptr(ti.elemsize))
	visref := refBitset.isset(ti.elemkind)

	vtypePtr := vtypeKind == reflect.Ptr
	ktypePtr := ktypeKind == reflect.Ptr

	vTransient := decUseTransient && !vtypePtr && ti.tielem.flagCanTransient
	kTransient := decUseTransient && !ktypePtr && ti.tikey.flagCanTransient

	var vtypeElem reflect.Type

	var keyFn, valFn *codecFn
	var ktypeLo, vtypeLo = ktype, vtype

	if ktypeKind == reflect.Ptr {
		for ktypeLo = ktype.Elem(); ktypeLo.Kind() == reflect.Ptr; ktypeLo = ktypeLo.Elem() {
		}
	}

	if vtypePtr {
		vtypeElem = vtype.Elem()
		for vtypeLo = vtypeElem; vtypeLo.Kind() == reflect.Ptr; vtypeLo = vtypeLo.Elem() {
		}
	}

	rvkMut := !scalarBitset.isset(ti.keykind) 
	rvvMut := !scalarBitset.isset(ti.elemkind)
	rvvCanNil := isnilBitset.isset(ti.elemkind)

	
	
	
	
	
	
	
	var rvk, rvkn, rvv, rvvn, rvva, rvvz reflect.Value

	
	var doMapGet, doMapSet bool

	if !d.h.MapValueReset {
		if rvvMut && (vtypeKind != reflect.Interface || !d.h.InterfaceReset) {
			doMapGet = true
			rvva = mapAddrLoopvarRV(vtype, vtypeKind)
		}
	}

	ktypeIsString := ktypeId == stringTypId
	ktypeIsIntf := ktypeId == intfTypId

	hasLen := containerLen > 0

	
	
	

	
	
	
	
	var kstrbs []byte
	var kstr2bs []byte
	var s string

	var callFnRvk bool

	fnRvk2 := func() (s string) {
		callFnRvk = false
		if len(kstr2bs) < 2 {
			return string(kstr2bs)
		}
		return d.mapKeyString(&callFnRvk, &kstrbs, &kstr2bs)
	}

	

	for j := 0; d.containerNext(j, containerLen, hasLen); j++ {
		callFnRvk = false
		if j == 0 {
			
			
			if decUseTransient && vTransient && kTransient {
				rvk = d.perType.TransientAddr2K(ktype, ktypeKind)
			} else {
				rvk = rvZeroAddrK(ktype, ktypeKind)
			}
			if !rvkMut {
				rvkn = rvk
			}
			if !rvvMut {
				if decUseTransient && vTransient {
					rvvn = d.perType.TransientAddrK(vtype, vtypeKind)
				} else {
					rvvn = rvZeroAddrK(vtype, vtypeKind)
				}
			}
			if !ktypeIsString && keyFn == nil {
				keyFn = d.h.fn(ktypeLo)
			}
			if valFn == nil {
				valFn = d.h.fn(vtypeLo)
			}
		} else if rvkMut {
			rvSetZero(rvk)
		} else {
			rvk = rvkn
		}

		d.mapElemKey()
		if ktypeIsString {
			kstr2bs = d.d.DecodeStringAsBytes()
			rvSetString(rvk, fnRvk2())
		} else {
			d.decByteState = decByteStateNone
			d.decodeValue(rvk, keyFn)
			
			if ktypeIsIntf {
				if rvk2 := rvk.Elem(); rvk2.IsValid() && rvk2.Type() == uint8SliceTyp {
					kstr2bs = rvGetBytes(rvk2)
					rvSetIntf(rvk, rv4istr(fnRvk2()))
				}
				
			}
		}

		d.mapElemValue()

		if d.d.TryNil() {
			
			if !rvvz.IsValid() {
				rvvz = rvZeroK(vtype, vtypeKind)
			}
			if callFnRvk {
				s = d.string(kstr2bs)
				if ktypeIsString {
					rvSetString(rvk, s)
				} else { 
					rvSetIntf(rvk, rv4istr(s))
				}
			}
			mapSet(rv, rvk, rvvz, kfast, visindirect, visref)
			continue
		}

		
		

		
		doMapSet = true

		if !rvvMut {
			rvv = rvvn
		} else if !doMapGet {
			goto NEW_RVV
		} else {
			rvv = mapGet(rv, rvk, rvva, kfast, visindirect, visref)
			if !rvv.IsValid() || (rvvCanNil && rvIsNil(rvv)) {
				goto NEW_RVV
			}
			switch vtypeKind {
			case reflect.Ptr, reflect.Map: 
				doMapSet = false
			case reflect.Interface:
				
				rvvn = rvv.Elem()
				if k := rvvn.Kind(); (k == reflect.Ptr || k == reflect.Map) && !rvIsNil(rvvn) {
					d.decodeValueNoCheckNil(rvvn, nil) 
					continue
				}
				
				rvvn = rvZeroAddrK(vtype, vtypeKind)
				rvSetIntf(rvvn, rvv)
				rvv = rvvn
			default:
				
				if decUseTransient && vTransient {
					rvvn = d.perType.TransientAddrK(vtype, vtypeKind)
				} else {
					rvvn = rvZeroAddrK(vtype, vtypeKind)
				}
				rvSetDirect(rvvn, rvv)
				rvv = rvvn
			}
		}
		goto DECODE_VALUE_NO_CHECK_NIL

	NEW_RVV:
		if vtypePtr {
			rvv = reflect.New(vtypeElem) 
		} else if decUseTransient && vTransient {
			rvv = d.perType.TransientAddrK(vtype, vtypeKind)
		} else {
			rvv = rvZeroAddrK(vtype, vtypeKind)
		}

	DECODE_VALUE_NO_CHECK_NIL:
		d.decodeValueNoCheckNil(rvv, valFn)

		if doMapSet {
			if callFnRvk {
				s = d.string(kstr2bs)
				if ktypeIsString {
					rvSetString(rvk, s)
				} else { 
					rvSetIntf(rvk, rv4istr(s))
				}
			}
			mapSet(rv, rvk, rvv, kfast, visindirect, visref)
		}
	}

	d.mapEnd()
}









type Decoder struct {
	panicHdl

	d decDriver

	
	mtid uintptr
	stid uintptr

	h *BasicHandle

	blist bytesFreelist

	
	decRd

	
	n fauxUnion

	hh  Handle
	err error

	perType decPerType

	
	is internerMap

	
	
	maxdepth int16
	depth    int16

	
	
	
	calls uint16 

	c containerState

	decByteState

	
	
	
	b [decScratchByteArrayLen]byte
}





func NewDecoder(r io.Reader, h Handle) *Decoder {
	d := h.newDecDriver().decoder()
	if r != nil {
		d.Reset(r)
	}
	return d
}



func NewDecoderBytes(in []byte, h Handle) *Decoder {
	d := h.newDecDriver().decoder()
	if in != nil {
		d.ResetBytes(in)
	}
	return d
}








func NewDecoderString(s string, h Handle) *Decoder {
	return NewDecoderBytes(bytesView(s), h)
}

func (d *Decoder) HandleName() string {
	return d.hh.Name()
}

func (d *Decoder) r() *decRd {
	return &d.decRd
}

func (d *Decoder) init(h Handle) {
	initHandle(h)
	d.cbreak = d.js || d.cbor
	d.bytes = true
	d.err = errDecoderNotInitialized
	d.h = h.getBasicHandle()
	d.hh = h
	d.be = h.isBinary()
	if d.h.InternString && d.is == nil {
		d.is.init()
	}
	
}

func (d *Decoder) resetCommon() {
	d.d.reset()
	d.err = nil
	d.c = 0
	d.decByteState = decByteStateNone
	d.depth = 0
	d.calls = 0
	
	d.maxdepth = decDefMaxDepth
	if d.h.MaxDepth > 0 {
		d.maxdepth = d.h.MaxDepth
	}
	d.mtid = 0
	d.stid = 0
	d.mtr = false
	d.str = false
	if d.h.MapType != nil {
		d.mtid = rt2id(d.h.MapType)
		d.mtr = fastpathAvIndex(d.mtid) != -1
	}
	if d.h.SliceType != nil {
		d.stid = rt2id(d.h.SliceType)
		d.str = fastpathAvIndex(d.stid) != -1
	}
}



func (d *Decoder) Reset(r io.Reader) {
	if r == nil {
		r = &eofReader
	}
	d.bytes = false
	if d.ri == nil {
		d.ri = new(ioDecReader)
	}
	d.ri.reset(r, d.h.ReaderBufferSize, &d.blist)
	d.decReader = d.ri
	d.resetCommon()
}



func (d *Decoder) ResetBytes(in []byte) {
	if in == nil {
		in = []byte{}
	}
	d.bytes = true
	d.decReader = &d.rb
	d.rb.reset(in)
	d.resetCommon()
}








func (d *Decoder) ResetString(s string) {
	d.ResetBytes(bytesView(s))
}

func (d *Decoder) naked() *fauxUnion {
	return &d.n
}
































































func (d *Decoder) Decode(v interface{}) (err error) {
	
	
	
	if !debugging {
		defer func() {
			if x := recover(); x != nil {
				panicValToErr(d, x, &d.err)
				err = d.err
			}
		}()
	}

	d.MustDecode(v)
	return
}




func (d *Decoder) MustDecode(v interface{}) {
	halt.onerror(d.err)
	if d.hh == nil {
		halt.onerror(errNoFormatHandle)
	}

	
	d.calls++
	d.decode(v)
	d.calls--
}





func (d *Decoder) Release() {
}

func (d *Decoder) swallow() {
	d.d.nextValueBytes(nil)
}

func (d *Decoder) swallowErr() (err error) {
	if !debugging {
		defer func() {
			if x := recover(); x != nil {
				panicValToErr(d, x, &err)
			}
		}()
	}
	d.swallow()
	return
}

func setZero(iv interface{}) {
	if iv == nil {
		return
	}
	rv, ok := isNil(iv)
	if ok {
		return
	}
	
	switch v := iv.(type) {
	case *string:
		*v = ""
	case *bool:
		*v = false
	case *int:
		*v = 0
	case *int8:
		*v = 0
	case *int16:
		*v = 0
	case *int32:
		*v = 0
	case *int64:
		*v = 0
	case *uint:
		*v = 0
	case *uint8:
		*v = 0
	case *uint16:
		*v = 0
	case *uint32:
		*v = 0
	case *uint64:
		*v = 0
	case *float32:
		*v = 0
	case *float64:
		*v = 0
	case *complex64:
		*v = 0
	case *complex128:
		*v = 0
	case *[]byte:
		*v = nil
	case *Raw:
		*v = nil
	case *time.Time:
		*v = time.Time{}
	case reflect.Value:
		decSetNonNilRV2Zero(v)
	default:
		if !fastpathDecodeSetZeroTypeSwitch(iv) {
			decSetNonNilRV2Zero(rv)
		}
	}
}


func decSetNonNilRV2Zero(v reflect.Value) {
	
	
	
	
	
	
	
	

	k := v.Kind()
	if k == reflect.Interface {
		decSetNonNilRV2Zero4Intf(v)
	} else if k == reflect.Ptr {
		decSetNonNilRV2Zero4Ptr(v)
	} else if v.CanSet() {
		rvSetDirectZero(v)
	}
}

func decSetNonNilRV2Zero4Ptr(v reflect.Value) {
	ve := v.Elem()
	if ve.CanSet() {
		rvSetZero(ve) 
	} else if v.CanSet() {
		rvSetZero(v)
	}
}

func decSetNonNilRV2Zero4Intf(v reflect.Value) {
	ve := v.Elem()
	if ve.CanSet() {
		rvSetDirectZero(ve) 
	} else if v.CanSet() {
		rvSetZero(v)
	}
}

func (d *Decoder) decode(iv interface{}) {
	
	

	if iv == nil {
		d.onerror(errCannotDecodeIntoNil)
	}

	switch v := iv.(type) {
	
	
	case reflect.Value:
		if x, _ := isDecodeable(v); !x {
			d.haltAsNotDecodeable(v)
		}
		d.decodeValue(v, nil)
	case *string:
		*v = d.stringZC(d.d.DecodeStringAsBytes())
	case *bool:
		*v = d.d.DecodeBool()
	case *int:
		*v = int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))
	case *int8:
		*v = int8(chkOvf.IntV(d.d.DecodeInt64(), 8))
	case *int16:
		*v = int16(chkOvf.IntV(d.d.DecodeInt64(), 16))
	case *int32:
		*v = int32(chkOvf.IntV(d.d.DecodeInt64(), 32))
	case *int64:
		*v = d.d.DecodeInt64()
	case *uint:
		*v = uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))
	case *uint8:
		*v = uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))
	case *uint16:
		*v = uint16(chkOvf.UintV(d.d.DecodeUint64(), 16))
	case *uint32:
		*v = uint32(chkOvf.UintV(d.d.DecodeUint64(), 32))
	case *uint64:
		*v = d.d.DecodeUint64()
	case *float32:
		*v = d.decodeFloat32()
	case *float64:
		*v = d.d.DecodeFloat64()
	case *complex64:
		*v = complex(d.decodeFloat32(), 0)
	case *complex128:
		*v = complex(d.d.DecodeFloat64(), 0)
	case *[]byte:
		*v = d.decodeBytesInto(*v)
	case []byte:
		
		b := d.decodeBytesInto(v[:len(v):len(v)])
		if !(len(b) > 0 && len(b) == len(v) && &b[0] == &v[0]) { 
			copy(v, b)
		}
	case *time.Time:
		*v = d.d.DecodeTime()
	case *Raw:
		*v = d.rawBytes()

	case *interface{}:
		d.decodeValue(rv4iptr(v), nil)

	default:
		
		if skipFastpathTypeSwitchInDirectCall || !fastpathDecodeTypeSwitch(iv, d) {
			v := reflect.ValueOf(iv)
			if x, _ := isDecodeable(v); !x {
				d.haltAsNotDecodeable(v)
			}
			d.decodeValue(v, nil)
		}
	}
}









func (d *Decoder) decodeValue(rv reflect.Value, fn *codecFn) {
	if d.d.TryNil() {
		decSetNonNilRV2Zero(rv)
		return
	}
	d.decodeValueNoCheckNil(rv, fn)
}

func (d *Decoder) decodeValueNoCheckNil(rv reflect.Value, fn *codecFn) {
	
	
	var rvp reflect.Value
	var rvpValid bool
PTR:
	if rv.Kind() == reflect.Ptr {
		rvpValid = true
		if rvIsNil(rv) {
			rvSetDirect(rv, reflect.New(rv.Type().Elem()))
		}
		rvp = rv
		rv = rv.Elem()
		goto PTR
	}

	if fn == nil {
		fn = d.h.fn(rv.Type())
	}
	if fn.i.addrD {
		if rvpValid {
			rv = rvp
		} else if rv.CanAddr() {
			rv = rvAddr(rv, fn.i.ti.ptr)
		} else if fn.i.addrDf {
			d.errorf("cannot decode into a non-pointer value")
		}
	}
	fn.fd(d, &fn.i, rv)
}

func (d *Decoder) structFieldNotFound(index int, rvkencname string) {
	
	
	
	if d.h.ErrorIfNoField {
		if index >= 0 {
			d.errorf("no matching struct field found when decoding stream array at index %v", index)
		} else if rvkencname != "" {
			d.errorf("no matching struct field found when decoding stream map with key " + rvkencname)
		}
	}
	d.swallow()
}

func (d *Decoder) arrayCannotExpand(sliceLen, streamLen int) {
	if d.h.ErrorIfNoArrayExpand {
		d.errorf("cannot expand array len during decode from %v to %v", sliceLen, streamLen)
	}
}

func (d *Decoder) haltAsNotDecodeable(rv reflect.Value) {
	if !rv.IsValid() {
		d.onerror(errCannotDecodeIntoNil)
	}
	
	if !rv.CanInterface() {
		d.errorf("cannot decode into a value without an interface: %v", rv)
	}
	d.errorf("cannot decode into value of kind: %v, %#v", rv.Kind(), rv2i(rv))
}

func (d *Decoder) depthIncr() {
	d.depth++
	if d.depth >= d.maxdepth {
		d.onerror(errMaxDepthExceeded)
	}
}

func (d *Decoder) depthDecr() {
	d.depth--
}





func (d *Decoder) string(v []byte) (s string) {
	if d.is == nil || d.c != containerMapKey || len(v) < 2 || len(v) > internMaxStrLen {
		return string(v)
	}
	return d.is.string(v)
}

func (d *Decoder) zerocopy() bool {
	return d.bytes && d.h.ZeroCopy
}




func (d *Decoder) decodeBytesInto(in []byte) (v []byte) {
	if in == nil {
		in = []byte{}
	}
	return d.d.DecodeBytes(in)
}

func (d *Decoder) rawBytes() (v []byte) {
	
	
	v = d.d.nextValueBytes([]byte{})
	if d.bytes && !d.h.ZeroCopy {
		vv := make([]byte, len(v))
		copy(vv, v) 
		v = vv
	}
	return
}

func (d *Decoder) wrapErr(v error, err *error) {
	*err = wrapCodecErr(v, d.hh.Name(), d.NumBytesRead(), false)
}


func (d *Decoder) NumBytesRead() int {
	return int(d.r().numread())
}





func (d *Decoder) decodeFloat32() float32 {
	if d.js {
		return d.jsondriver().DecodeFloat32() 
	}
	return float32(chkOvf.Float32V(d.d.DecodeFloat64()))
}














func (d *Decoder) checkBreak() (v bool) {
	
	
	

	
	
	
	
	
	

	if d.cbreak {
		v = d.d.CheckBreak()
	}
	return
}

func (d *Decoder) containerNext(j, containerLen int, hasLen bool) bool {
	

	
	if hasLen {
		return j < containerLen
	}
	return !d.checkBreak()
}

func (d *Decoder) mapStart(v int) int {
	if v != containerLenNil {
		d.depthIncr()
		d.c = containerMapStart
	}
	return v
}

func (d *Decoder) mapElemKey() {
	if d.js {
		d.jsondriver().ReadMapElemKey()
	}
	d.c = containerMapKey
}

func (d *Decoder) mapElemValue() {
	if d.js {
		d.jsondriver().ReadMapElemValue()
	}
	d.c = containerMapValue
}

func (d *Decoder) mapEnd() {
	if d.js {
		d.jsondriver().ReadMapEnd()
	}
	
	d.depthDecr()
	d.c = 0
}

func (d *Decoder) arrayStart(v int) int {
	if v != containerLenNil {
		d.depthIncr()
		d.c = containerArrayStart
	}
	return v
}

func (d *Decoder) arrayElem() {
	if d.js {
		d.jsondriver().ReadArrayElem()
	}
	d.c = containerArrayElem
}

func (d *Decoder) arrayEnd() {
	if d.js {
		d.jsondriver().ReadArrayEnd()
	}
	
	d.depthDecr()
	d.c = 0
}

func (d *Decoder) interfaceExtConvertAndDecode(v interface{}, ext InterfaceExt) {
	
	
	

	
	
	
	
	
	
	

	var s interface{}
	rv := reflect.ValueOf(v)
	rv2 := rv.Elem()
	rvk := rv2.Kind()
	if rvk == reflect.Struct || rvk == reflect.Array {
		s = ext.ConvertExt(v)
	} else {
		s = ext.ConvertExt(rv2i(rv2))
	}
	rv = reflect.ValueOf(s)

	
	
	

	if !rv.CanAddr() {
		rvk = rv.Kind()
		rv2 = d.oneShotAddrRV(rv.Type(), rvk)
		if rvk == reflect.Interface {
			rvSetIntf(rv2, rv)
		} else {
			rvSetDirect(rv2, rv)
		}
		rv = rv2
	}

	d.decodeValue(rv, nil)
	ext.UpdateExt(v, rv2i(rv))
}

func (d *Decoder) sideDecode(v interface{}, basetype reflect.Type, bs []byte) {
	

	defer func(rb bytesDecReader, bytes bool,
		c containerState, dbs decByteState, depth int16, r decReader, state interface{}) {
		d.rb = rb
		d.bytes = bytes
		d.c = c
		d.decByteState = dbs
		d.depth = depth
		d.decReader = r
		d.d.restoreState(state)
	}(d.rb, d.bytes, d.c, d.decByteState, d.depth, d.decReader, d.d.captureState())

	
	d.rb = bytesDecReader{bs[:len(bs):len(bs)], 0}
	d.bytes = true
	d.decReader = &d.rb
	d.d.resetState()
	d.c = 0
	d.decByteState = decByteStateNone
	d.depth = 0

	
	d.decodeValue(baseRV(v), d.h.fnNoExt(basetype))
}

func (d *Decoder) fauxUnionReadRawBytes(asString bool) {
	if asString || d.h.RawToString {
		d.n.v = valueTypeString
		
		d.n.s = d.stringZC(d.d.DecodeBytes(nil))
	} else {
		d.n.v = valueTypeBytes
		d.n.l = d.d.DecodeBytes([]byte{})
	}
}

func (d *Decoder) oneShotAddrRV(rvt reflect.Type, rvk reflect.Kind) reflect.Value {
	if decUseTransient &&
		(numBoolStrSliceBitset.isset(byte(rvk)) ||
			((rvk == reflect.Struct || rvk == reflect.Array) &&
				d.h.getTypeInfo(rt2id(rvt), rvt).flagCanTransient)) {
		return d.perType.TransientAddrK(rvt, rvk)
	}
	return rvZeroAddrK(rvt, rvk)
}







type decSliceHelper struct {
	d     *Decoder
	ct    valueType
	Array bool
	IsNil bool
}

func (d *Decoder) decSliceHelperStart() (x decSliceHelper, clen int) {
	x.ct = d.d.ContainerType()
	x.d = d
	switch x.ct {
	case valueTypeNil:
		x.IsNil = true
	case valueTypeArray:
		x.Array = true
		clen = d.arrayStart(d.d.ReadArrayStart())
	case valueTypeMap:
		clen = d.mapStart(d.d.ReadMapStart())
		clen += clen
	default:
		d.errorf("only encoded map or array can be decoded into a slice (%d)", x.ct)
	}
	return
}

func (x decSliceHelper) End() {
	if x.IsNil {
	} else if x.Array {
		x.d.arrayEnd()
	} else {
		x.d.mapEnd()
	}
}

func (x decSliceHelper) ElemContainerState(index int) {
	

	if x.Array {
		x.d.arrayElem()
	} else if index&1 == 0 { 
		x.d.mapElemKey()
	} else {
		x.d.mapElemValue()
	}
}

func (x decSliceHelper) arrayCannotExpand(hasLen bool, lenv, j, containerLenS int) {
	x.d.arrayCannotExpand(lenv, j+1)
	
	x.ElemContainerState(j)
	x.d.swallow()
	j++
	for ; x.d.containerNext(j, containerLenS, hasLen); j++ {
		x.ElemContainerState(j)
		x.d.swallow()
	}
	x.End()
}









type decNextValueBytesHelper struct {
	d *Decoder
}

func (x decNextValueBytesHelper) append1(v *[]byte, b byte) {
	if *v != nil && !x.d.bytes {
		*v = append(*v, b)
	}
}

func (x decNextValueBytesHelper) appendN(v *[]byte, b ...byte) {
	if *v != nil && !x.d.bytes {
		*v = append(*v, b...)
	}
}

func (x decNextValueBytesHelper) appendS(v *[]byte, b string) {
	if *v != nil && !x.d.bytes {
		*v = append(*v, b...)
	}
}

func (x decNextValueBytesHelper) bytesRdV(v *[]byte, startpos uint) {
	if x.d.bytes {
		*v = x.d.rb.b[startpos:x.d.rb.c]
	}
}






type decNegintPosintFloatNumberHelper struct {
	d *Decoder
}

func (x decNegintPosintFloatNumberHelper) uint64(ui uint64, neg, ok bool) uint64 {
	if ok && !neg {
		return ui
	}
	return x.uint64TryFloat(ok)
}

func (x decNegintPosintFloatNumberHelper) uint64TryFloat(ok bool) (ui uint64) {
	if ok { 
		x.d.errorf("assigning negative signed value to unsigned type")
	}
	f, ok := x.d.d.decFloat()
	if ok && f >= 0 && noFrac64(math.Float64bits(f)) {
		ui = uint64(f)
	} else {
		x.d.errorf("invalid number loading uint64, with descriptor: %v", x.d.d.descBd())
	}
	return ui
}

func decNegintPosintFloatNumberHelperInt64v(ui uint64, neg, incrIfNeg bool) (i int64) {
	if neg && incrIfNeg {
		ui++
	}
	i = chkOvf.SignedIntV(ui)
	if neg {
		i = -i
	}
	return
}

func (x decNegintPosintFloatNumberHelper) int64(ui uint64, neg, ok bool) (i int64) {
	if ok {
		return decNegintPosintFloatNumberHelperInt64v(ui, neg, x.d.cbor)
	}
	
	
	
	f, ok := x.d.d.decFloat()
	if ok && noFrac64(math.Float64bits(f)) {
		i = int64(f)
	} else {
		x.d.errorf("invalid number loading uint64, with descriptor: %v", x.d.d.descBd())
	}
	return
}

func (x decNegintPosintFloatNumberHelper) float64(f float64, ok bool) float64 {
	if ok {
		return f
	}
	return x.float64TryInteger()
}

func (x decNegintPosintFloatNumberHelper) float64TryInteger() float64 {
	ui, neg, ok := x.d.d.decInteger()
	if !ok {
		x.d.errorf("invalid descriptor for float: %v", x.d.d.descBd())
	}
	return float64(decNegintPosintFloatNumberHelperInt64v(ui, neg, x.d.cbor))
}












func isDecodeable(rv reflect.Value) (canDecode bool, reason decNotDecodeableReason) {
	switch rv.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Chan, reflect.Map:
		canDecode = !rvIsNil(rv)
		reason = decNotDecodeableReasonNilReference
	case reflect.Func, reflect.Interface, reflect.Invalid, reflect.UnsafePointer:
		reason = decNotDecodeableReasonBadKind
	default:
		canDecode = rv.CanAddr()
		reason = decNotDecodeableReasonNonAddrValue
	}
	return
}

func decByteSlice(r *decRd, clen, maxInitLen int, bs []byte) (bsOut []byte) {
	if clen <= 0 {
		bsOut = zeroByteSlice
	} else if cap(bs) >= clen {
		bsOut = bs[:clen]
		r.readb(bsOut)
	} else {
		var len2 int
		for len2 < clen {
			len3 := decInferLen(clen-len2, maxInitLen, 1)
			bs3 := bsOut
			bsOut = make([]byte, len2+len3)
			copy(bsOut, bs3)
			r.readb(bsOut[len2:])
			len2 += len3
		}
	}
	return
}






func decInferLen(clen, maxlen, unit int) int {
	
	
	
	const (
		minLenIfUnset = 8
		maxMem        = 256 * 1024 
	)

	

	
	
	
	
	

	if clen == 0 || clen == containerLenNil {
		return 0
	}
	if clen < 0 {
		
		clen = 64 / unit
		if clen > minLenIfUnset {
			return clen
		}
		return minLenIfUnset
	}
	if unit <= 0 {
		return clen
	}
	if maxlen <= 0 {
		maxlen = maxMem / unit
	}
	if clen < maxlen {
		return clen
	}
	return maxlen
}
