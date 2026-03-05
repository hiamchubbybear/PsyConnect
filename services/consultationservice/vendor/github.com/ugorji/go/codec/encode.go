


package codec

import (
	"encoding"
	"errors"
	"io"
	"reflect"
	"sort"
	"strconv"
	"time"
)



const defEncByteBufSize = 1 << 10 

var errEncoderNotInitialized = errors.New("Encoder not initialized")


type encDriver interface {
	EncodeNil()
	EncodeInt(i int64)
	EncodeUint(i uint64)
	EncodeBool(b bool)
	EncodeFloat32(f float32)
	EncodeFloat64(f float64)
	EncodeRawExt(re *RawExt)
	EncodeExt(v interface{}, basetype reflect.Type, xtag uint64, ext Ext)
	
	EncodeString(v string)
	EncodeStringBytesRaw(v []byte)
	EncodeTime(time.Time)
	WriteArrayStart(length int)
	WriteArrayEnd()
	WriteMapStart(length int)
	WriteMapEnd()

	
	reset()

	encoder() *Encoder

	driverStateManager
}

type encDriverContainerTracker interface {
	WriteArrayElem()
	WriteMapElemKey()
	WriteMapElemValue()
}

type encDriverNoState struct{}

func (encDriverNoState) captureState() interface{}  { return nil }
func (encDriverNoState) reset()                     {}
func (encDriverNoState) resetState()                {}
func (encDriverNoState) restoreState(v interface{}) {}

type encDriverNoopContainerWriter struct{}

func (encDriverNoopContainerWriter) WriteArrayStart(length int) {}
func (encDriverNoopContainerWriter) WriteArrayEnd()             {}
func (encDriverNoopContainerWriter) WriteMapStart(length int)   {}
func (encDriverNoopContainerWriter) WriteMapEnd()               {}


type encStructFieldObj struct {
	key   string
	rv    reflect.Value
	intf  interface{}
	ascii bool
	isRv  bool
}

type encStructFieldObjSlice []encStructFieldObj

func (p encStructFieldObjSlice) Len() int      { return len(p) }
func (p encStructFieldObjSlice) Swap(i, j int) { p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] }
func (p encStructFieldObjSlice) Less(i, j int) bool {
	return p[uint(i)].key < p[uint(j)].key
}


type EncodeOptions struct {
	
	
	
	WriterBufferSize int

	
	
	
	
	
	
	ChanRecvTimeout time.Duration

	
	StructToArray bool

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	Canonical bool

	
	
	
	
	
	
	
	CheckCircularRef bool

	
	
	
	
	
	
	
	
	
	
	RecursiveEmptyCheck bool

	
	
	
	
	
	Raw bool

	
	
	
	
	
	
	
	
	
	
	
	StringToRaw bool

	
	
	
	
	
	
	
	
	OptimumSize bool

	
	
	
	
	
	
	NoAddressableReadonly bool
}



func (e *Encoder) rawExt(f *codecFnInfo, rv reflect.Value) {
	e.e.EncodeRawExt(rv2i(rv).(*RawExt))
}

func (e *Encoder) ext(f *codecFnInfo, rv reflect.Value) {
	e.e.EncodeExt(rv2i(rv), f.ti.rt, f.xfTag, f.xfFn)
}

func (e *Encoder) selferMarshal(f *codecFnInfo, rv reflect.Value) {
	rv2i(rv).(Selfer).CodecEncodeSelf(e)
}

func (e *Encoder) binaryMarshal(f *codecFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(encoding.BinaryMarshaler).MarshalBinary()
	e.marshalRaw(bs, fnerr)
}

func (e *Encoder) textMarshal(f *codecFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(encoding.TextMarshaler).MarshalText()
	e.marshalUtf8(bs, fnerr)
}

func (e *Encoder) jsonMarshal(f *codecFnInfo, rv reflect.Value) {
	bs, fnerr := rv2i(rv).(jsonMarshaler).MarshalJSON()
	e.marshalAsis(bs, fnerr)
}

func (e *Encoder) raw(f *codecFnInfo, rv reflect.Value) {
	e.rawBytes(rv2i(rv).(Raw))
}

func (e *Encoder) encodeComplex64(v complex64) {
	if imag(v) != 0 {
		e.errorf("cannot encode complex number: %v, with imaginary values: %v", v, imag(v))
	}
	e.e.EncodeFloat32(real(v))
}

func (e *Encoder) encodeComplex128(v complex128) {
	if imag(v) != 0 {
		e.errorf("cannot encode complex number: %v, with imaginary values: %v", v, imag(v))
	}
	e.e.EncodeFloat64(real(v))
}

func (e *Encoder) kBool(f *codecFnInfo, rv reflect.Value) {
	e.e.EncodeBool(rvGetBool(rv))
}

func (e *Encoder) kTime(f *codecFnInfo, rv reflect.Value) {
	e.e.EncodeTime(rvGetTime(rv))
}

func (e *Encoder) kString(f *codecFnInfo, rv reflect.Value) {
	e.e.EncodeString(rvGetString(rv))
}

func (e *Encoder) kFloat32(f *codecFnInfo, rv reflect.Value) {
	e.e.EncodeFloat32(rvGetFloat32(rv))
}

func (e *Encoder) kFloat64(f *codecFnInfo, rv reflect.Value) {
	e.e.EncodeFloat64(rvGetFloat64(rv))
}

func (e *Encoder) kComplex64(f *codecFnInfo, rv reflect.Value) {
	e.encodeComplex64(rvGetComplex64(rv))
}

func (e *Encoder) kComplex128(f *codecFnInfo, rv reflect.Value) {
	e.encodeComplex128(rvGetComplex128(rv))
}

func (e *Encoder) kInt(f *codecFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt(rv)))
}

func (e *Encoder) kInt8(f *codecFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt8(rv)))
}

func (e *Encoder) kInt16(f *codecFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt16(rv)))
}

func (e *Encoder) kInt32(f *codecFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt32(rv)))
}

func (e *Encoder) kInt64(f *codecFnInfo, rv reflect.Value) {
	e.e.EncodeInt(int64(rvGetInt64(rv)))
}

func (e *Encoder) kUint(f *codecFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint(rv)))
}

func (e *Encoder) kUint8(f *codecFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint8(rv)))
}

func (e *Encoder) kUint16(f *codecFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint16(rv)))
}

func (e *Encoder) kUint32(f *codecFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint32(rv)))
}

func (e *Encoder) kUint64(f *codecFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUint64(rv)))
}

func (e *Encoder) kUintptr(f *codecFnInfo, rv reflect.Value) {
	e.e.EncodeUint(uint64(rvGetUintptr(rv)))
}

func (e *Encoder) kErr(f *codecFnInfo, rv reflect.Value) {
	e.errorf("unsupported kind %s, for %#v", rv.Kind(), rv)
}

func chanToSlice(rv reflect.Value, rtslice reflect.Type, timeout time.Duration) (rvcs reflect.Value) {
	rvcs = rvZeroK(rtslice, reflect.Slice)
	if timeout < 0 { 
		for {
			recv, recvOk := rv.Recv()
			if !recvOk {
				break
			}
			rvcs = reflect.Append(rvcs, recv)
		}
	} else {
		cases := make([]reflect.SelectCase, 2)
		cases[0] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: rv}
		if timeout == 0 {
			cases[1] = reflect.SelectCase{Dir: reflect.SelectDefault}
		} else {
			tt := time.NewTimer(timeout)
			cases[1] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(tt.C)}
		}
		for {
			chosen, recv, recvOk := reflect.Select(cases)
			if chosen == 1 || !recvOk {
				break
			}
			rvcs = reflect.Append(rvcs, recv)
		}
	}
	return
}

func (e *Encoder) kSeqFn(rtelem reflect.Type) (fn *codecFn) {
	for rtelem.Kind() == reflect.Ptr {
		rtelem = rtelem.Elem()
	}
	
	
	if rtelem.Kind() != reflect.Interface {
		fn = e.h.fn(rtelem)
	}
	return
}

func (e *Encoder) kSliceWMbs(rv reflect.Value, ti *typeInfo) {
	var l = rvLenSlice(rv)
	if l == 0 {
		e.mapStart(0)
	} else {
		e.haltOnMbsOddLen(l)
		e.mapStart(l >> 1) 
		fn := e.kSeqFn(ti.elem)
		for j := 0; j < l; j++ {
			if j&1 == 0 { 
				e.mapElemKey()
			} else {
				e.mapElemValue()
			}
			e.encodeValue(rvSliceIndex(rv, j, ti), fn)
		}
	}
	e.mapEnd()
}

func (e *Encoder) kSliceW(rv reflect.Value, ti *typeInfo) {
	var l = rvLenSlice(rv)
	e.arrayStart(l)
	if l > 0 {
		fn := e.kSeqFn(ti.elem)
		for j := 0; j < l; j++ {
			e.arrayElem()
			e.encodeValue(rvSliceIndex(rv, j, ti), fn)
		}
	}
	e.arrayEnd()
}

func (e *Encoder) kArrayWMbs(rv reflect.Value, ti *typeInfo) {
	var l = rv.Len()
	if l == 0 {
		e.mapStart(0)
	} else {
		e.haltOnMbsOddLen(l)
		e.mapStart(l >> 1) 
		fn := e.kSeqFn(ti.elem)
		for j := 0; j < l; j++ {
			if j&1 == 0 { 
				e.mapElemKey()
			} else {
				e.mapElemValue()
			}
			e.encodeValue(rv.Index(j), fn)
		}
	}
	e.mapEnd()
}

func (e *Encoder) kArrayW(rv reflect.Value, ti *typeInfo) {
	var l = rv.Len()
	e.arrayStart(l)
	if l > 0 {
		fn := e.kSeqFn(ti.elem)
		for j := 0; j < l; j++ {
			e.arrayElem()
			e.encodeValue(rv.Index(j), fn)
		}
	}
	e.arrayEnd()
}

func (e *Encoder) kChan(f *codecFnInfo, rv reflect.Value) {
	if f.ti.chandir&uint8(reflect.RecvDir) == 0 {
		e.errorf("send-only channel cannot be encoded")
	}
	if !f.ti.mbs && uint8TypId == rt2id(f.ti.elem) {
		e.kSliceBytesChan(rv)
		return
	}
	rtslice := reflect.SliceOf(f.ti.elem)
	rv = chanToSlice(rv, rtslice, e.h.ChanRecvTimeout)
	ti := e.h.getTypeInfo(rt2id(rtslice), rtslice)
	if f.ti.mbs {
		e.kSliceWMbs(rv, ti)
	} else {
		e.kSliceW(rv, ti)
	}
}

func (e *Encoder) kSlice(f *codecFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		e.kSliceWMbs(rv, f.ti)
	} else if f.ti.rtid == uint8SliceTypId || uint8TypId == rt2id(f.ti.elem) {
		e.e.EncodeStringBytesRaw(rvGetBytes(rv))
	} else {
		e.kSliceW(rv, f.ti)
	}
}

func (e *Encoder) kArray(f *codecFnInfo, rv reflect.Value) {
	if f.ti.mbs {
		e.kArrayWMbs(rv, f.ti)
	} else if handleBytesWithinKArray && uint8TypId == rt2id(f.ti.elem) {
		e.e.EncodeStringBytesRaw(rvGetArrayBytes(rv, []byte{}))
	} else {
		e.kArrayW(rv, f.ti)
	}
}

func (e *Encoder) kSliceBytesChan(rv reflect.Value) {
	
	

	bs0 := e.blist.peek(32, true)
	bs := bs0

	irv := rv2i(rv)
	ch, ok := irv.(<-chan byte)
	if !ok {
		ch = irv.(chan byte)
	}

L1:
	switch timeout := e.h.ChanRecvTimeout; {
	case timeout == 0: 
		for {
			select {
			case b := <-ch:
				bs = append(bs, b)
			default:
				break L1
			}
		}
	case timeout > 0: 
		tt := time.NewTimer(timeout)
		for {
			select {
			case b := <-ch:
				bs = append(bs, b)
			case <-tt.C:
				
				break L1
			}
		}
	default: 
		for b := range ch {
			bs = append(bs, b)
		}
	}

	e.e.EncodeStringBytesRaw(bs)
	e.blist.put(bs)
	if !byteSliceSameData(bs0, bs) {
		e.blist.put(bs0)
	}
}

func (e *Encoder) kStructSfi(f *codecFnInfo) []*structFieldInfo {
	if e.h.Canonical {
		return f.ti.sfi.sorted()
	}
	return f.ti.sfi.source()
}

func (e *Encoder) kStructNoOmitempty(f *codecFnInfo, rv reflect.Value) {
	var tisfi []*structFieldInfo
	if f.ti.toArray || e.h.StructToArray { 
		tisfi = f.ti.sfi.source()
		e.arrayStart(len(tisfi))
		for _, si := range tisfi {
			e.arrayElem()
			e.encodeValue(si.path.field(rv), nil)
		}
		e.arrayEnd()
	} else {
		tisfi = e.kStructSfi(f)
		e.mapStart(len(tisfi))
		keytyp := f.ti.keyType
		for _, si := range tisfi {
			e.mapElemKey()
			e.kStructFieldKey(keytyp, si.path.encNameAsciiAlphaNum, si.encName)
			e.mapElemValue()
			e.encodeValue(si.path.field(rv), nil)
		}
		e.mapEnd()
	}
}

func (e *Encoder) kStructFieldKey(keyType valueType, encNameAsciiAlphaNum bool, encName string) {
	encStructFieldKey(encName, e.e, e.w(), keyType, encNameAsciiAlphaNum, e.js)
}

func (e *Encoder) kStruct(f *codecFnInfo, rv reflect.Value) {
	var newlen int
	ti := f.ti
	toMap := !(ti.toArray || e.h.StructToArray)
	var mf map[string]interface{}
	if ti.flagMissingFielder {
		mf = rv2i(rv).(MissingFielder).CodecMissingFields()
		toMap = true
		newlen += len(mf)
	} else if ti.flagMissingFielderPtr {
		rv2 := e.addrRV(rv, ti.rt, ti.ptr)
		mf = rv2i(rv2).(MissingFielder).CodecMissingFields()
		toMap = true
		newlen += len(mf)
	}
	tisfi := ti.sfi.source()
	newlen += len(tisfi)

	var fkvs = e.slist.get(newlen)[:newlen]

	recur := e.h.RecursiveEmptyCheck

	var kv sfiRv
	var j int
	if toMap {
		newlen = 0
		for _, si := range e.kStructSfi(f) {
			kv.r = si.path.field(rv)
			if si.path.omitEmpty && isEmptyValue(kv.r, e.h.TypeInfos, recur) {
				continue
			}
			kv.v = si
			fkvs[newlen] = kv
			newlen++
		}

		var mf2s []stringIntf
		if len(mf) > 0 {
			mf2s = make([]stringIntf, 0, len(mf))
			for k, v := range mf {
				if k == "" {
					continue
				}
				if ti.infoFieldOmitempty && isEmptyValue(reflect.ValueOf(v), e.h.TypeInfos, recur) {
					continue
				}
				mf2s = append(mf2s, stringIntf{k, v})
			}
		}

		e.mapStart(newlen + len(mf2s))

		
		
		

		if len(mf2s) > 0 && e.h.Canonical {
			mf2w := make([]encStructFieldObj, newlen+len(mf2s))
			for j = 0; j < newlen; j++ {
				kv = fkvs[j]
				mf2w[j] = encStructFieldObj{kv.v.encName, kv.r, nil, kv.v.path.encNameAsciiAlphaNum, true}
			}
			for _, v := range mf2s {
				mf2w[j] = encStructFieldObj{v.v, reflect.Value{}, v.i, false, false}
				j++
			}
			sort.Sort((encStructFieldObjSlice)(mf2w))
			for _, v := range mf2w {
				e.mapElemKey()
				e.kStructFieldKey(ti.keyType, v.ascii, v.key)
				e.mapElemValue()
				if v.isRv {
					e.encodeValue(v.rv, nil)
				} else {
					e.encode(v.intf)
				}
			}
		} else {
			keytyp := ti.keyType
			for j = 0; j < newlen; j++ {
				kv = fkvs[j]
				e.mapElemKey()
				e.kStructFieldKey(keytyp, kv.v.path.encNameAsciiAlphaNum, kv.v.encName)
				e.mapElemValue()
				e.encodeValue(kv.r, nil)
			}
			for _, v := range mf2s {
				e.mapElemKey()
				e.kStructFieldKey(keytyp, false, v.v)
				e.mapElemValue()
				e.encode(v.i)
			}
		}

		e.mapEnd()
	} else {
		newlen = len(tisfi)
		for i, si := range tisfi { 
			kv.r = si.path.field(rv)
			
			
			if si.path.omitEmpty && isEmptyValue(kv.r, e.h.TypeInfos, recur) {
				switch kv.r.Kind() {
				case reflect.Struct, reflect.Interface, reflect.Ptr, reflect.Array, reflect.Map, reflect.Slice:
					kv.r = reflect.Value{} 
				}
			}
			fkvs[i] = kv
		}
		
		e.arrayStart(newlen)
		for j = 0; j < newlen; j++ {
			e.arrayElem()
			e.encodeValue(fkvs[j].r, nil)
		}
		e.arrayEnd()
	}

	
	
	
	e.slist.put(fkvs)
}

func (e *Encoder) kMap(f *codecFnInfo, rv reflect.Value) {
	l := rvLenMap(rv)
	e.mapStart(l)
	if l == 0 {
		e.mapEnd()
		return
	}

	
	
	
	
	
	
	

	var keyFn, valFn *codecFn

	ktypeKind := reflect.Kind(f.ti.keykind)
	vtypeKind := reflect.Kind(f.ti.elemkind)

	rtval := f.ti.elem
	rtvalkind := vtypeKind
	for rtvalkind == reflect.Ptr {
		rtval = rtval.Elem()
		rtvalkind = rtval.Kind()
	}
	if rtvalkind != reflect.Interface {
		valFn = e.h.fn(rtval)
	}

	var rvv = mapAddrLoopvarRV(f.ti.elem, vtypeKind)

	rtkey := f.ti.key
	var keyTypeIsString = stringTypId == rt2id(rtkey) 
	if keyTypeIsString {
		keyFn = e.h.fn(rtkey)
	} else {
		for rtkey.Kind() == reflect.Ptr {
			rtkey = rtkey.Elem()
		}
		if rtkey.Kind() != reflect.Interface {
			keyFn = e.h.fn(rtkey)
		}
	}

	if e.h.Canonical {
		e.kMapCanonical(f.ti, rv, rvv, keyFn, valFn)
		e.mapEnd()
		return
	}

	var rvk = mapAddrLoopvarRV(f.ti.key, ktypeKind)

	var it mapIter
	mapRange(&it, rv, rvk, rvv, true)

	for it.Next() {
		e.mapElemKey()
		if keyTypeIsString {
			e.e.EncodeString(it.Key().String())
		} else {
			e.encodeValue(it.Key(), keyFn)
		}
		e.mapElemValue()
		e.encodeValue(it.Value(), valFn)
	}
	it.Done()

	e.mapEnd()
}

func (e *Encoder) kMapCanonical(ti *typeInfo, rv, rvv reflect.Value, keyFn, valFn *codecFn) {
	
	
	
	
	rtkey := ti.key
	rtkeydecl := rtkey.PkgPath() == "" && rtkey.Name() != "" 

	mks := rv.MapKeys()
	rtkeyKind := rtkey.Kind()
	kfast := mapKeyFastKindFor(rtkeyKind)
	visindirect := mapStoresElemIndirect(uintptr(ti.elemsize))
	visref := refBitset.isset(ti.elemkind)

	switch rtkeyKind {
	case reflect.Bool:
		
		
		

		
		
		if len(mks) == 2 && mks[0].Bool() {
			mks[0], mks[1] = mks[1], mks[0]
		}
		for i := range mks {
			e.mapElemKey()
			if rtkeydecl {
				e.e.EncodeBool(mks[i].Bool())
			} else {
				e.encodeValueNonNil(mks[i], keyFn)
			}
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mks[i], rvv, kfast, visindirect, visref), valFn)
		}
	case reflect.String:
		mksv := make([]stringRv, len(mks))
		for i, k := range mks {
			v := &mksv[i]
			v.r = k
			v.v = k.String()
		}
		sort.Sort(stringRvSlice(mksv))
		for i := range mksv {
			e.mapElemKey()
			if rtkeydecl {
				e.e.EncodeString(mksv[i].v)
			} else {
				e.encodeValueNonNil(mksv[i].r, keyFn)
			}
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mksv[i].r, rvv, kfast, visindirect, visref), valFn)
		}
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint, reflect.Uintptr:
		mksv := make([]uint64Rv, len(mks))
		for i, k := range mks {
			v := &mksv[i]
			v.r = k
			v.v = k.Uint()
		}
		sort.Sort(uint64RvSlice(mksv))
		for i := range mksv {
			e.mapElemKey()
			if rtkeydecl {
				e.e.EncodeUint(mksv[i].v)
			} else {
				e.encodeValueNonNil(mksv[i].r, keyFn)
			}
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mksv[i].r, rvv, kfast, visindirect, visref), valFn)
		}
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		mksv := make([]int64Rv, len(mks))
		for i, k := range mks {
			v := &mksv[i]
			v.r = k
			v.v = k.Int()
		}
		sort.Sort(int64RvSlice(mksv))
		for i := range mksv {
			e.mapElemKey()
			if rtkeydecl {
				e.e.EncodeInt(mksv[i].v)
			} else {
				e.encodeValueNonNil(mksv[i].r, keyFn)
			}
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mksv[i].r, rvv, kfast, visindirect, visref), valFn)
		}
	case reflect.Float32:
		mksv := make([]float64Rv, len(mks))
		for i, k := range mks {
			v := &mksv[i]
			v.r = k
			v.v = k.Float()
		}
		sort.Sort(float64RvSlice(mksv))
		for i := range mksv {
			e.mapElemKey()
			if rtkeydecl {
				e.e.EncodeFloat32(float32(mksv[i].v))
			} else {
				e.encodeValueNonNil(mksv[i].r, keyFn)
			}
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mksv[i].r, rvv, kfast, visindirect, visref), valFn)
		}
	case reflect.Float64:
		mksv := make([]float64Rv, len(mks))
		for i, k := range mks {
			v := &mksv[i]
			v.r = k
			v.v = k.Float()
		}
		sort.Sort(float64RvSlice(mksv))
		for i := range mksv {
			e.mapElemKey()
			if rtkeydecl {
				e.e.EncodeFloat64(mksv[i].v)
			} else {
				e.encodeValueNonNil(mksv[i].r, keyFn)
			}
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mksv[i].r, rvv, kfast, visindirect, visref), valFn)
		}
	default:
		if rtkey == timeTyp {
			mksv := make([]timeRv, len(mks))
			for i, k := range mks {
				v := &mksv[i]
				v.r = k
				v.v = rv2i(k).(time.Time)
			}
			sort.Sort(timeRvSlice(mksv))
			for i := range mksv {
				e.mapElemKey()
				e.e.EncodeTime(mksv[i].v)
				e.mapElemValue()
				e.encodeValue(mapGet(rv, mksv[i].r, rvv, kfast, visindirect, visref), valFn)
			}
			break
		}

		
		
		bs0 := e.blist.get(len(mks) * 16)
		mksv := bs0
		mksbv := make([]bytesRv, len(mks))

		func() {
			
			defer func(wb bytesEncAppender, bytes bool, c containerState, state interface{}) {
				e.wb = wb
				e.bytes = bytes
				e.c = c
				e.e.restoreState(state)
			}(e.wb, e.bytes, e.c, e.e.captureState())

			
			e.wb = bytesEncAppender{mksv[:0], &mksv}
			e.bytes = true
			e.c = 0
			e.e.resetState()

			for i, k := range mks {
				v := &mksbv[i]
				l := len(mksv)

				e.c = containerMapKey
				e.encodeValue(k, nil)
				e.atEndOfEncode()
				e.w().end()

				v.r = k
				v.v = mksv[l:]
			}
		}()

		sort.Sort(bytesRvSlice(mksbv))
		for j := range mksbv {
			e.mapElemKey()
			e.encWr.writeb(mksbv[j].v)
			e.mapElemValue()
			e.encodeValue(mapGet(rv, mksbv[j].r, rvv, kfast, visindirect, visref), valFn)
		}
		e.blist.put(mksv)
		if !byteSliceSameData(bs0, mksv) {
			e.blist.put(bs0)
		}
	}
}









type Encoder struct {
	panicHdl

	e encDriver

	h *BasicHandle

	
	encWr

	
	hh Handle

	blist bytesFreelist
	err   error

	

	

	
	
	
	
	
	
	
	ci []interface{} 

	perType encPerType

	slist sfiRvFreelist
}





func NewEncoder(w io.Writer, h Handle) *Encoder {
	e := h.newEncDriver().encoder()
	if w != nil {
		e.Reset(w)
	}
	return e
}






func NewEncoderBytes(out *[]byte, h Handle) *Encoder {
	e := h.newEncDriver().encoder()
	if out != nil {
		e.ResetBytes(out)
	}
	return e
}

func (e *Encoder) HandleName() string {
	return e.hh.Name()
}

func (e *Encoder) init(h Handle) {
	initHandle(h)
	e.err = errEncoderNotInitialized
	e.bytes = true
	e.hh = h
	e.h = h.getBasicHandle()
	e.be = e.hh.isBinary()
}

func (e *Encoder) w() *encWr {
	return &e.encWr
}

func (e *Encoder) resetCommon() {
	e.e.reset()
	if e.ci != nil {
		e.ci = e.ci[:0]
	}
	e.c = 0
	e.calls = 0
	e.seq = 0
	e.err = nil
}





func (e *Encoder) Reset(w io.Writer) {
	e.bytes = false
	if e.wf == nil {
		e.wf = new(bufioEncWriter)
	}
	e.wf.reset(w, e.h.WriterBufferSize, &e.blist)
	e.resetCommon()
}


func (e *Encoder) ResetBytes(out *[]byte) {
	e.bytes = true
	e.wb.reset(encInBytes(out), out)
	e.resetCommon()
}






















































































func (e *Encoder) Encode(v interface{}) (err error) {
	
	
	
	if !debugging {
		defer func() {
			
			
			if x := recover(); x != nil {
				panicValToErr(e, x, &e.err)
				err = e.err
			}
		}()
	}

	e.MustEncode(v)
	return
}




func (e *Encoder) MustEncode(v interface{}) {
	halt.onerror(e.err)
	if e.hh == nil {
		halt.onerror(errNoFormatHandle)
	}

	e.calls++
	e.encode(v)
	e.calls--
	if e.calls == 0 {
		e.atEndOfEncode()
		e.w().end()
	}
}





func (e *Encoder) Release() {
}

func (e *Encoder) encode(iv interface{}) {
	
	

	if iv == nil {
		e.e.EncodeNil()
		return
	}

	rv, ok := isNil(iv)
	if ok {
		e.e.EncodeNil()
		return
	}

	switch v := iv.(type) {
	
	
	case Raw:
		e.rawBytes(v)
	case reflect.Value:
		e.encodeValue(v, nil)

	case string:
		e.e.EncodeString(v)
	case bool:
		e.e.EncodeBool(v)
	case int:
		e.e.EncodeInt(int64(v))
	case int8:
		e.e.EncodeInt(int64(v))
	case int16:
		e.e.EncodeInt(int64(v))
	case int32:
		e.e.EncodeInt(int64(v))
	case int64:
		e.e.EncodeInt(v)
	case uint:
		e.e.EncodeUint(uint64(v))
	case uint8:
		e.e.EncodeUint(uint64(v))
	case uint16:
		e.e.EncodeUint(uint64(v))
	case uint32:
		e.e.EncodeUint(uint64(v))
	case uint64:
		e.e.EncodeUint(v)
	case uintptr:
		e.e.EncodeUint(uint64(v))
	case float32:
		e.e.EncodeFloat32(v)
	case float64:
		e.e.EncodeFloat64(v)
	case complex64:
		e.encodeComplex64(v)
	case complex128:
		e.encodeComplex128(v)
	case time.Time:
		e.e.EncodeTime(v)
	case []byte:
		e.e.EncodeStringBytesRaw(v)
	case *Raw:
		e.rawBytes(*v)
	case *string:
		e.e.EncodeString(*v)
	case *bool:
		e.e.EncodeBool(*v)
	case *int:
		e.e.EncodeInt(int64(*v))
	case *int8:
		e.e.EncodeInt(int64(*v))
	case *int16:
		e.e.EncodeInt(int64(*v))
	case *int32:
		e.e.EncodeInt(int64(*v))
	case *int64:
		e.e.EncodeInt(*v)
	case *uint:
		e.e.EncodeUint(uint64(*v))
	case *uint8:
		e.e.EncodeUint(uint64(*v))
	case *uint16:
		e.e.EncodeUint(uint64(*v))
	case *uint32:
		e.e.EncodeUint(uint64(*v))
	case *uint64:
		e.e.EncodeUint(*v)
	case *uintptr:
		e.e.EncodeUint(uint64(*v))
	case *float32:
		e.e.EncodeFloat32(*v)
	case *float64:
		e.e.EncodeFloat64(*v)
	case *complex64:
		e.encodeComplex64(*v)
	case *complex128:
		e.encodeComplex128(*v)
	case *time.Time:
		e.e.EncodeTime(*v)
	case *[]byte:
		if *v == nil {
			e.e.EncodeNil()
		} else {
			e.e.EncodeStringBytesRaw(*v)
		}
	default:
		
		if skipFastpathTypeSwitchInDirectCall || !fastpathEncodeTypeSwitch(iv, e) {
			e.encodeValue(rv, nil)
		}
	}
}





func (e *Encoder) encodeValue(rv reflect.Value, fn *codecFn) {
	

	

	var sptr interface{}
	var rvp reflect.Value
	var rvpValid bool
TOP:
	switch rv.Kind() {
	case reflect.Ptr:
		if rvIsNil(rv) {
			e.e.EncodeNil()
			return
		}
		rvpValid = true
		rvp = rv
		rv = rv.Elem()
		goto TOP
	case reflect.Interface:
		if rvIsNil(rv) {
			e.e.EncodeNil()
			return
		}
		rvpValid = false
		rvp = reflect.Value{}
		rv = rv.Elem()
		goto TOP
	case reflect.Struct:
		if rvpValid && e.h.CheckCircularRef {
			sptr = rv2i(rvp)
			for _, vv := range e.ci {
				if eq4i(sptr, vv) { 
					e.errorf("circular reference found: %p, %T", sptr, sptr)
				}
			}
			e.ci = append(e.ci, sptr)
		}
	case reflect.Slice, reflect.Map, reflect.Chan:
		if rvIsNil(rv) {
			e.e.EncodeNil()
			return
		}
	case reflect.Invalid, reflect.Func:
		e.e.EncodeNil()
		return
	}

	if fn == nil {
		fn = e.h.fn(rv.Type())
	}

	if !fn.i.addrE { 
		
	} else if rvpValid {
		rv = rvp
	} else {
		rv = e.addrRV(rv, fn.i.ti.rt, fn.i.ti.ptr)
	}
	fn.fe(e, &fn.i, rv)

	if sptr != nil { 
		e.ci = e.ci[:len(e.ci)-1]
	}
}



func (e *Encoder) encodeValueNonNil(rv reflect.Value, fn *codecFn) {
	if fn == nil {
		fn = e.h.fn(rv.Type())
	}

	if fn.i.addrE { 
		rv = e.addrRV(rv, fn.i.ti.rt, fn.i.ti.ptr)
	}
	fn.fe(e, &fn.i, rv)
}


func (e *Encoder) addrRV(rv reflect.Value, typ, ptrType reflect.Type) (rva reflect.Value) {
	if rv.CanAddr() {
		return rvAddr(rv, ptrType)
	}
	if e.h.NoAddressableReadonly {
		rva = reflect.New(typ)
		rvSetDirect(rva.Elem(), rv)
		return
	}
	return rvAddr(e.perType.AddressableRO(rv), ptrType)
}

func (e *Encoder) marshalUtf8(bs []byte, fnerr error) {
	e.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.EncodeString(stringView(bs))
	}
}

func (e *Encoder) marshalAsis(bs []byte, fnerr error) {
	e.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.encWr.writeb(bs) 
	}
}

func (e *Encoder) marshalRaw(bs []byte, fnerr error) {
	e.onerror(fnerr)
	if bs == nil {
		e.e.EncodeNil()
	} else {
		e.e.EncodeStringBytesRaw(bs)
	}
}

func (e *Encoder) rawBytes(vv Raw) {
	v := []byte(vv)
	if !e.h.Raw {
		e.errorf("Raw values cannot be encoded: %v", v)
	}
	e.encWr.writeb(v)
}

func (e *Encoder) wrapErr(v error, err *error) {
	*err = wrapCodecErr(v, e.hh.Name(), 0, true)
}





func (e *Encoder) mapStart(length int) {
	e.e.WriteMapStart(length)
	e.c = containerMapStart
}

func (e *Encoder) mapElemKey() {
	if e.js {
		e.jsondriver().WriteMapElemKey()
	}
	e.c = containerMapKey
}

func (e *Encoder) mapElemValue() {
	if e.js {
		e.jsondriver().WriteMapElemValue()
	}
	e.c = containerMapValue
}

func (e *Encoder) mapEnd() {
	e.e.WriteMapEnd()
	e.c = 0
}

func (e *Encoder) arrayStart(length int) {
	e.e.WriteArrayStart(length)
	e.c = containerArrayStart
}

func (e *Encoder) arrayElem() {
	if e.js {
		e.jsondriver().WriteArrayElem()
	}
	e.c = containerArrayElem
}

func (e *Encoder) arrayEnd() {
	e.e.WriteArrayEnd()
	e.c = 0
}



func (e *Encoder) haltOnMbsOddLen(length int) {
	if length&1 != 0 { 
		e.errorf("mapBySlice requires even slice length, but got %v", length)
	}
}

func (e *Encoder) atEndOfEncode() {
	
	if e.js {
		e.jsondriver().atEndOfEncode()
	}
}

func (e *Encoder) sideEncode(v interface{}, basetype reflect.Type, bs *[]byte) {
	
	
	
	
	

	defer func(wb bytesEncAppender, bytes bool, c containerState, state interface{}) {
		e.wb = wb
		e.bytes = bytes
		e.c = c
		e.e.restoreState(state)
	}(e.wb, e.bytes, e.c, e.e.captureState())

	e.wb = bytesEncAppender{encInBytes(bs)[:0], bs}
	e.bytes = true
	e.c = 0
	e.e.resetState()

	
	rv := baseRV(v)
	e.encodeValue(rv, e.h.fnNoExt(basetype))
	e.atEndOfEncode()
	e.w().end()
}

func encInBytes(out *[]byte) (in []byte) {
	in = *out
	if in == nil {
		in = make([]byte, defEncByteBufSize)
	}
	return
}

func encStructFieldKey(encName string, ee encDriver, w *encWr,
	keyType valueType, encNameAsciiAlphaNum bool, js bool) {
	
	

	if keyType == valueTypeString {
		if js && encNameAsciiAlphaNum { 
			w.writeqstr(encName)
		} else { 
			ee.EncodeString(encName)
		}
	} else if keyType == valueTypeInt {
		ee.EncodeInt(must.Int(strconv.ParseInt(encName, 10, 64)))
	} else if keyType == valueTypeUint {
		ee.EncodeUint(must.Uint(strconv.ParseUint(encName, 10, 64)))
	} else if keyType == valueTypeFloat {
		ee.EncodeFloat64(must.Float(strconv.ParseFloat(encName, 64)))
	} else {
		halt.errorf("invalid struct key type: %v", keyType)
	}
}
