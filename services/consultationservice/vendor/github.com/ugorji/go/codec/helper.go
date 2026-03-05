


package codec



























































































































































































import (
	"encoding"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unicode/utf8"
)










const debugging = false

const (
	
	
	
	containerLenUnknown = -1

	
	
	containerLenNil = math.MinInt32

	
	
	
	
	
	
	handleBytesWithinKArray = true

	
	
	supportMarshalInterfaces = true

	
	bytesFreeListNoCache = false

	
	
	cacheLineSize = 64

	wordSizeBits = 32 << (^uint(0) >> 63) 
	wordSize     = wordSizeBits / 8

	
	
	
	
	skipFastpathTypeSwitchInDirectCall = false
)

const cpu32Bit = ^uint(0)>>32 == 0

type rkind byte

const (
	rkindPtr    = rkind(reflect.Ptr)
	rkindString = rkind(reflect.String)
	rkindChan   = rkind(reflect.Chan)
)

type mapKeyFastKind uint8

const (
	mapKeyFastKind32 = iota + 1
	mapKeyFastKind32ptr
	mapKeyFastKind64
	mapKeyFastKind64ptr
	mapKeyFastKindStr
)

var (
	
	
	
	handleInitMu sync.Mutex

	must mustHdl
	halt panicHdl

	digitCharBitset      bitset256
	numCharBitset        bitset256
	whitespaceCharBitset bitset256
	asciiAlphaNumBitset  bitset256

	
	
	
	
	
	

	
	refBitset bitset32

	
	isnilBitset bitset32

	
	numBoolBitset bitset32

	
	numBoolStrSliceBitset bitset32

	
	scalarBitset bitset32

	mapKeyFastKindVals [32]mapKeyFastKind

	
	codecgen bool

	oneByteArr    [1]byte
	zeroByteSlice = oneByteArr[:0:0]

	eofReader devNullReader
)

var (
	errMapTypeNotMapKind     = errors.New("MapType MUST be of Map Kind")
	errSliceTypeNotSliceKind = errors.New("SliceType MUST be of Slice Kind")

	errExtFnWriteExtUnsupported   = errors.New("BytesExt.WriteExt is not supported")
	errExtFnReadExtUnsupported    = errors.New("BytesExt.ReadExt is not supported")
	errExtFnConvertExtUnsupported = errors.New("InterfaceExt.ConvertExt is not supported")
	errExtFnUpdateExtUnsupported  = errors.New("InterfaceExt.UpdateExt is not supported")

	errPanicUndefined = errors.New("panic: undefined error")

	errHandleInited = errors.New("cannot modify initialized Handle")

	errNoFormatHandle = errors.New("no handle (cannot identify format)")
)

var pool4tiload = sync.Pool{
	New: func() interface{} {
		return &typeInfoLoad{
			etypes:   make([]uintptr, 0, 4),
			sfis:     make([]structFieldInfo, 0, 4),
			sfiNames: make(map[string]uint16, 4),
		}
	},
}

func init() {
	xx := func(f mapKeyFastKind, k ...reflect.Kind) {
		for _, v := range k {
			mapKeyFastKindVals[byte(v)&31] = f 
		}
	}

	var f mapKeyFastKind

	f = mapKeyFastKind64
	if wordSizeBits == 32 {
		f = mapKeyFastKind32
	}
	xx(f, reflect.Int, reflect.Uint, reflect.Uintptr)

	f = mapKeyFastKind64ptr
	if wordSizeBits == 32 {
		f = mapKeyFastKind32ptr
	}
	xx(f, reflect.Ptr)

	xx(mapKeyFastKindStr, reflect.String)
	xx(mapKeyFastKind32, reflect.Uint32, reflect.Int32, reflect.Float32)
	xx(mapKeyFastKind64, reflect.Uint64, reflect.Int64, reflect.Float64)

	numBoolBitset.
		set(byte(reflect.Bool)).
		set(byte(reflect.Int)).
		set(byte(reflect.Int8)).
		set(byte(reflect.Int16)).
		set(byte(reflect.Int32)).
		set(byte(reflect.Int64)).
		set(byte(reflect.Uint)).
		set(byte(reflect.Uint8)).
		set(byte(reflect.Uint16)).
		set(byte(reflect.Uint32)).
		set(byte(reflect.Uint64)).
		set(byte(reflect.Uintptr)).
		set(byte(reflect.Float32)).
		set(byte(reflect.Float64)).
		set(byte(reflect.Complex64)).
		set(byte(reflect.Complex128))

	numBoolStrSliceBitset = numBoolBitset

	numBoolStrSliceBitset.
		set(byte(reflect.String)).
		set(byte(reflect.Slice))

	scalarBitset = numBoolBitset

	scalarBitset.
		set(byte(reflect.String))

	

	refBitset.
		set(byte(reflect.Map)).
		set(byte(reflect.Ptr)).
		set(byte(reflect.Func)).
		set(byte(reflect.Chan)).
		set(byte(reflect.UnsafePointer))

	isnilBitset = refBitset

	isnilBitset.
		set(byte(reflect.Interface)).
		set(byte(reflect.Slice))

	
	
	
	

	for i := byte(0); i <= utf8.RuneSelf; i++ {
		if (i >= '0' && i <= '9') || (i >= 'a' && i <= 'z') || (i >= 'A' && i <= 'Z') {
			asciiAlphaNumBitset.set(i)
		}
		switch i {
		case ' ', '\t', '\r', '\n':
			whitespaceCharBitset.set(i)
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			digitCharBitset.set(i)
			numCharBitset.set(i)
		case '.', '+', '-':
			numCharBitset.set(i)
		case 'e', 'E':
			numCharBitset.set(i)
		}
	}
}





type driverStateManager interface {
	resetState()
	captureState() interface{}
	restoreState(state interface{})
}

type bdAndBdread struct {
	bdRead bool
	bd     byte
}

func (x bdAndBdread) captureState() interface{}   { return x }
func (x *bdAndBdread) resetState()                { x.bd, x.bdRead = 0, false }
func (x *bdAndBdread) reset()                     { x.resetState() }
func (x *bdAndBdread) restoreState(v interface{}) { *x = v.(bdAndBdread) }

type clsErr struct {
	err    error 
	closed bool  
}

type charEncoding uint8

const (
	_ charEncoding = iota 
	cUTF8
	cUTF16LE
	cUTF16BE
	cUTF32LE
	cUTF32BE
	
	cRAW charEncoding = 255
)


type valueType uint8

const (
	valueTypeUnset valueType = iota
	valueTypeNil
	valueTypeInt
	valueTypeUint
	valueTypeFloat
	valueTypeBool
	valueTypeString
	valueTypeSymbol
	valueTypeBytes
	valueTypeMap
	valueTypeArray
	valueTypeTime
	valueTypeExt

	
)

var valueTypeStrings = [...]string{
	"Unset",
	"Nil",
	"Int",
	"Uint",
	"Float",
	"Bool",
	"String",
	"Symbol",
	"Bytes",
	"Map",
	"Array",
	"Timestamp",
	"Ext",
}

func (x valueType) String() string {
	if int(x) < len(valueTypeStrings) {
		return valueTypeStrings[x]
	}
	return strconv.FormatInt(int64(x), 10)
}



type containerState uint8

const (
	_ containerState = iota

	containerMapStart
	containerMapKey
	containerMapValue
	containerMapEnd
	containerArrayStart
	containerArrayElem
	containerArrayEnd
)





const rgetMaxRecursion = 2










type fauxUnion struct {
	

	
	u uint64
	i int64
	f float64
	l []byte
	s string

	
	t time.Time
	b bool

	
	v valueType
}


type typeInfoLoad struct {
	etypes   []uintptr
	sfis     []structFieldInfo
	sfiNames map[string]uint16
}

func (x *typeInfoLoad) reset() {
	x.etypes = x.etypes[:0]
	x.sfis = x.sfis[:0]
	for k := range x.sfiNames { 
		delete(x.sfiNames, k)
	}
}




type jsonMarshaler interface {
	MarshalJSON() ([]byte, error)
}
type jsonUnmarshaler interface {
	UnmarshalJSON([]byte) error
}

type isZeroer interface {
	IsZero() bool
}

type isCodecEmptyer interface {
	IsCodecEmpty() bool
}

type codecError struct {
	err    error
	name   string
	pos    int
	encode bool
}

func (e *codecError) Cause() error {
	return e.err
}

func (e *codecError) Unwrap() error {
	return e.err
}

func (e *codecError) Error() string {
	if e.encode {
		return fmt.Sprintf("%s encode error: %v", e.name, e.err)
	}
	return fmt.Sprintf("%s decode error [pos %d]: %v", e.name, e.pos, e.err)
}

func wrapCodecErr(in error, name string, numbytesread int, encode bool) (out error) {
	x, ok := in.(*codecError)
	if ok && x.pos == numbytesread && x.name == name && x.encode == encode {
		return in
	}
	return &codecError{in, name, numbytesread, encode}
}

var (
	bigen bigenHelper

	bigenstd = binary.BigEndian

	structInfoFieldName = "_struct"

	mapStrIntfTyp  = reflect.TypeOf(map[string]interface{}(nil))
	mapIntfIntfTyp = reflect.TypeOf(map[interface{}]interface{}(nil))
	intfSliceTyp   = reflect.TypeOf([]interface{}(nil))
	intfTyp        = intfSliceTyp.Elem()

	reflectValTyp = reflect.TypeOf((*reflect.Value)(nil)).Elem()

	stringTyp     = reflect.TypeOf("")
	timeTyp       = reflect.TypeOf(time.Time{})
	rawExtTyp     = reflect.TypeOf(RawExt{})
	rawTyp        = reflect.TypeOf(Raw{})
	uintptrTyp    = reflect.TypeOf(uintptr(0))
	uint8Typ      = reflect.TypeOf(uint8(0))
	uint8SliceTyp = reflect.TypeOf([]uint8(nil))
	uintTyp       = reflect.TypeOf(uint(0))
	intTyp        = reflect.TypeOf(int(0))

	mapBySliceTyp = reflect.TypeOf((*MapBySlice)(nil)).Elem()

	binaryMarshalerTyp   = reflect.TypeOf((*encoding.BinaryMarshaler)(nil)).Elem()
	binaryUnmarshalerTyp = reflect.TypeOf((*encoding.BinaryUnmarshaler)(nil)).Elem()

	textMarshalerTyp   = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
	textUnmarshalerTyp = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()

	jsonMarshalerTyp   = reflect.TypeOf((*jsonMarshaler)(nil)).Elem()
	jsonUnmarshalerTyp = reflect.TypeOf((*jsonUnmarshaler)(nil)).Elem()

	selferTyp                = reflect.TypeOf((*Selfer)(nil)).Elem()
	missingFielderTyp        = reflect.TypeOf((*MissingFielder)(nil)).Elem()
	iszeroTyp                = reflect.TypeOf((*isZeroer)(nil)).Elem()
	isCodecEmptyerTyp        = reflect.TypeOf((*isCodecEmptyer)(nil)).Elem()
	isSelferViaCodecgenerTyp = reflect.TypeOf((*isSelferViaCodecgener)(nil)).Elem()

	uint8TypId      = rt2id(uint8Typ)
	uint8SliceTypId = rt2id(uint8SliceTyp)
	rawExtTypId     = rt2id(rawExtTyp)
	rawTypId        = rt2id(rawTyp)
	intfTypId       = rt2id(intfTyp)
	timeTypId       = rt2id(timeTyp)
	stringTypId     = rt2id(stringTyp)

	mapStrIntfTypId  = rt2id(mapStrIntfTyp)
	mapIntfIntfTypId = rt2id(mapIntfIntfTyp)
	intfSliceTypId   = rt2id(intfSliceTyp)
	

	intBitsize  = uint8(intTyp.Bits())
	uintBitsize = uint8(uintTyp.Bits())

	
	bsAll0xff = []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

	chkOvf checkOverflow
)

var defTypeInfos = NewTypeInfos([]string{"codec", "json"})








var SelfExt = &extFailWrapper{}


















type Selfer interface {
	CodecEncodeSelf(*Encoder)
	CodecDecodeSelf(*Decoder)
}

type isSelferViaCodecgener interface {
	codecSelferViaCodecgen()
}










type MissingFielder interface {
	
	
	
	CodecMissingField(field []byte, value interface{}) bool

	
	
	
	CodecMissingFields() map[string]interface{}
}






















type MapBySlice interface {
	MapBySlice()
}








type basicHandleRuntimeState struct {
	
	
	rtidFns      atomicRtidFnSlice
	rtidFnsNoExt atomicRtidFnSlice

	
	
	

	extHandle

	intf2impls

	mu sync.Mutex

	jsonHandle   bool
	binaryHandle bool

	
	
	
	timeBuiltin bool
	_           bool 
}




type BasicHandle struct {
	
	

	
	
	
	TypeInfos *TypeInfos

	*basicHandleRuntimeState

	

	DecodeOptions

	

	EncodeOptions

	RPCOptions

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	TimeNotBuiltin bool

	
	
	
	
	ExplicitRelease bool

	
	inited uint32 

}




func initHandle(hh Handle) {
	x := hh.getBasicHandle()

	
	
	
	
	
	
	
	

	
	
	
	if atomic.LoadUint32(&x.inited) == 0 {
		x.initHandle(hh)
	}
}

func (x *BasicHandle) basicInit() {
	x.rtidFns.store(nil)
	x.rtidFnsNoExt.store(nil)
	x.timeBuiltin = !x.TimeNotBuiltin
}

func (x *BasicHandle) init() {}

func (x *BasicHandle) isInited() bool {
	return atomic.LoadUint32(&x.inited) != 0
}


func (x *BasicHandle) clearInited() {
	atomic.StoreUint32(&x.inited, 0)
}



func (x *basicHandleRuntimeState) TimeBuiltin() bool {
	return x.timeBuiltin
}

func (x *basicHandleRuntimeState) isJs() bool {
	return x.jsonHandle
}

func (x *basicHandleRuntimeState) isBe() bool {
	return x.binaryHandle
}

func (x *basicHandleRuntimeState) setExt(rt reflect.Type, tag uint64, ext Ext) (err error) {
	rk := rt.Kind()
	for rk == reflect.Ptr {
		rt = rt.Elem()
		rk = rt.Kind()
	}

	if rt.PkgPath() == "" || rk == reflect.Interface { 
		return fmt.Errorf("codec.Handle.SetExt: Takes named type, not a pointer or interface: %v", rt)
	}

	rtid := rt2id(rt)
	
	
	
	switch rtid {
	case rawTypId, rawExtTypId:
		return
	case timeTypId:
		if x.timeBuiltin {
			return
		}
	}

	for i := range x.extHandle {
		v := &x.extHandle[i]
		if v.rtid == rtid {
			v.tag, v.ext = tag, ext
			return
		}
	}
	rtidptr := rt2id(reflect.PtrTo(rt))
	x.extHandle = append(x.extHandle, extTypeTagFn{rtid, rtidptr, rt, tag, ext})
	return
}




//go:noinline
func (x *BasicHandle) initHandle(hh Handle) {
	handleInitMu.Lock()
	defer handleInitMu.Unlock() 
	if x.inited == 0 {
		if x.basicHandleRuntimeState == nil {
			x.basicHandleRuntimeState = new(basicHandleRuntimeState)
		}
		x.jsonHandle = hh.isJson()
		x.binaryHandle = hh.isBinary()
		
		if x.MapType != nil && x.MapType.Kind() != reflect.Map {
			halt.onerror(errMapTypeNotMapKind)
		}
		if x.SliceType != nil && x.SliceType.Kind() != reflect.Slice {
			halt.onerror(errSliceTypeNotSliceKind)
		}
		x.basicInit()
		hh.init()
		atomic.StoreUint32(&x.inited, 1)
	}
}

func (x *BasicHandle) getBasicHandle() *BasicHandle {
	return x
}

func (x *BasicHandle) typeInfos() *TypeInfos {
	if x.TypeInfos != nil {
		return x.TypeInfos
	}
	return defTypeInfos
}

func (x *BasicHandle) getTypeInfo(rtid uintptr, rt reflect.Type) (pti *typeInfo) {
	return x.typeInfos().get(rtid, rt)
}

func findRtidFn(s []codecRtidFn, rtid uintptr) (i uint, fn *codecFn) {
	
	

	
	var h uint 
	var j = uint(len(s))
LOOP:
	if i < j {
		h = (i + j) >> 1 
		if s[h].rtid < rtid {
			i = h + 1
		} else {
			j = h
		}
		goto LOOP
	}
	if i < uint(len(s)) && s[i].rtid == rtid {
		fn = s[i].fn
	}
	return
}

func (x *BasicHandle) fn(rt reflect.Type) (fn *codecFn) {
	return x.fnVia(rt, x.typeInfos(), &x.rtidFns, x.CheckCircularRef, true)
}

func (x *BasicHandle) fnNoExt(rt reflect.Type) (fn *codecFn) {
	return x.fnVia(rt, x.typeInfos(), &x.rtidFnsNoExt, x.CheckCircularRef, false)
}

func (x *basicHandleRuntimeState) fnVia(rt reflect.Type, tinfos *TypeInfos, fs *atomicRtidFnSlice, checkCircularRef, checkExt bool) (fn *codecFn) {
	rtid := rt2id(rt)
	sp := fs.load()
	if sp != nil {
		if _, fn = findRtidFn(sp, rtid); fn != nil {
			return
		}
	}

	fn = x.fnLoad(rt, rtid, tinfos, checkCircularRef, checkExt)
	x.mu.Lock()
	sp = fs.load()
	
	
	if sp == nil {
		sp = []codecRtidFn{{rtid, fn}}
		fs.store(sp)
	} else {
		idx, fn2 := findRtidFn(sp, rtid)
		if fn2 == nil {
			sp2 := make([]codecRtidFn, len(sp)+1)
			copy(sp2[idx+1:], sp[idx:])
			copy(sp2, sp[:idx])
			sp2[idx] = codecRtidFn{rtid, fn}
			fs.store(sp2)
		}
	}
	x.mu.Unlock()
	return
}

func fnloadFastpathUnderlying(ti *typeInfo) (f *fastpathE, u reflect.Type) {
	var rtid uintptr
	var idx int
	rtid = rt2id(ti.fastpathUnderlying)
	idx = fastpathAvIndex(rtid)
	if idx == -1 {
		return
	}
	f = &fastpathAv[idx]
	if uint8(reflect.Array) == ti.kind {
		u = reflectArrayOf(ti.rt.Len(), ti.elem)
	} else {
		u = f.rt
	}
	return
}

func (x *basicHandleRuntimeState) fnLoad(rt reflect.Type, rtid uintptr, tinfos *TypeInfos, checkCircularRef, checkExt bool) (fn *codecFn) {
	fn = new(codecFn)
	fi := &(fn.i)
	ti := tinfos.get(rtid, rt)
	fi.ti = ti
	rk := reflect.Kind(ti.kind)

	
	
	

	fi.addrDf = true
	

	if rtid == timeTypId && x.timeBuiltin {
		fn.fe = (*Encoder).kTime
		fn.fd = (*Decoder).kTime
	} else if rtid == rawTypId {
		fn.fe = (*Encoder).raw
		fn.fd = (*Decoder).raw
	} else if rtid == rawExtTypId {
		fn.fe = (*Encoder).rawExt
		fn.fd = (*Decoder).rawExt
		fi.addrD = true
		fi.addrE = true
	} else if xfFn := x.getExt(rtid, checkExt); xfFn != nil {
		fi.xfTag, fi.xfFn = xfFn.tag, xfFn.ext
		fn.fe = (*Encoder).ext
		fn.fd = (*Decoder).ext
		fi.addrD = true
		if rk == reflect.Struct || rk == reflect.Array {
			fi.addrE = true
		}
	} else if (ti.flagSelfer || ti.flagSelferPtr) &&
		!(checkCircularRef && ti.flagSelferViaCodecgen && ti.kind == byte(reflect.Struct)) {
		
		fn.fe = (*Encoder).selferMarshal
		fn.fd = (*Decoder).selferUnmarshal
		fi.addrD = ti.flagSelferPtr
		fi.addrE = ti.flagSelferPtr
	} else if supportMarshalInterfaces && x.isBe() &&
		(ti.flagBinaryMarshaler || ti.flagBinaryMarshalerPtr) &&
		(ti.flagBinaryUnmarshaler || ti.flagBinaryUnmarshalerPtr) {
		fn.fe = (*Encoder).binaryMarshal
		fn.fd = (*Decoder).binaryUnmarshal
		fi.addrD = ti.flagBinaryUnmarshalerPtr
		fi.addrE = ti.flagBinaryMarshalerPtr
	} else if supportMarshalInterfaces && !x.isBe() && x.isJs() &&
		(ti.flagJsonMarshaler || ti.flagJsonMarshalerPtr) &&
		(ti.flagJsonUnmarshaler || ti.flagJsonUnmarshalerPtr) {
		
		fn.fe = (*Encoder).jsonMarshal
		fn.fd = (*Decoder).jsonUnmarshal
		fi.addrD = ti.flagJsonUnmarshalerPtr
		fi.addrE = ti.flagJsonMarshalerPtr
	} else if supportMarshalInterfaces && !x.isBe() &&
		(ti.flagTextMarshaler || ti.flagTextMarshalerPtr) &&
		(ti.flagTextUnmarshaler || ti.flagTextUnmarshalerPtr) {
		fn.fe = (*Encoder).textMarshal
		fn.fd = (*Decoder).textUnmarshal
		fi.addrD = ti.flagTextUnmarshalerPtr
		fi.addrE = ti.flagTextMarshalerPtr
	} else {
		if fastpathEnabled && (rk == reflect.Map || rk == reflect.Slice || rk == reflect.Array) {
			
			
			
			
			
			
			
			
			
			
			
			var rtid2 uintptr
			if !ti.flagHasPkgPath { 
				rtid2 = rtid
				if rk == reflect.Array {
					rtid2 = rt2id(ti.key) 
				}
				if idx := fastpathAvIndex(rtid2); idx != -1 {
					fn.fe = fastpathAv[idx].encfn
					fn.fd = fastpathAv[idx].decfn
					fi.addrD = true
					fi.addrDf = false
					if rk == reflect.Array {
						fi.addrD = false 
					}
				}
			} else { 
				
				xfe, xrt := fnloadFastpathUnderlying(ti)
				if xfe != nil {
					xfnf := xfe.encfn
					xfnf2 := xfe.decfn
					if rk == reflect.Array {
						fi.addrD = false 
						fn.fd = func(d *Decoder, xf *codecFnInfo, xrv reflect.Value) {
							xfnf2(d, xf, rvConvert(xrv, xrt))
						}
					} else {
						fi.addrD = true
						fi.addrDf = false 
						xptr2rt := reflect.PtrTo(xrt)
						fn.fd = func(d *Decoder, xf *codecFnInfo, xrv reflect.Value) {
							if xrv.Kind() == reflect.Ptr {
								xfnf2(d, xf, rvConvert(xrv, xptr2rt))
							} else {
								xfnf2(d, xf, rvConvert(xrv, xrt))
							}
						}
					}
					fn.fe = func(e *Encoder, xf *codecFnInfo, xrv reflect.Value) {
						xfnf(e, xf, rvConvert(xrv, xrt))
					}
				}
			}
		}
		if fn.fe == nil && fn.fd == nil {
			switch rk {
			case reflect.Bool:
				fn.fe = (*Encoder).kBool
				fn.fd = (*Decoder).kBool
			case reflect.String:
				
				
				
				
				
				
				
				
				

				fn.fe = (*Encoder).kString
				fn.fd = (*Decoder).kString
			case reflect.Int:
				fn.fd = (*Decoder).kInt
				fn.fe = (*Encoder).kInt
			case reflect.Int8:
				fn.fe = (*Encoder).kInt8
				fn.fd = (*Decoder).kInt8
			case reflect.Int16:
				fn.fe = (*Encoder).kInt16
				fn.fd = (*Decoder).kInt16
			case reflect.Int32:
				fn.fe = (*Encoder).kInt32
				fn.fd = (*Decoder).kInt32
			case reflect.Int64:
				fn.fe = (*Encoder).kInt64
				fn.fd = (*Decoder).kInt64
			case reflect.Uint:
				fn.fd = (*Decoder).kUint
				fn.fe = (*Encoder).kUint
			case reflect.Uint8:
				fn.fe = (*Encoder).kUint8
				fn.fd = (*Decoder).kUint8
			case reflect.Uint16:
				fn.fe = (*Encoder).kUint16
				fn.fd = (*Decoder).kUint16
			case reflect.Uint32:
				fn.fe = (*Encoder).kUint32
				fn.fd = (*Decoder).kUint32
			case reflect.Uint64:
				fn.fe = (*Encoder).kUint64
				fn.fd = (*Decoder).kUint64
			case reflect.Uintptr:
				fn.fe = (*Encoder).kUintptr
				fn.fd = (*Decoder).kUintptr
			case reflect.Float32:
				fn.fe = (*Encoder).kFloat32
				fn.fd = (*Decoder).kFloat32
			case reflect.Float64:
				fn.fe = (*Encoder).kFloat64
				fn.fd = (*Decoder).kFloat64
			case reflect.Complex64:
				fn.fe = (*Encoder).kComplex64
				fn.fd = (*Decoder).kComplex64
			case reflect.Complex128:
				fn.fe = (*Encoder).kComplex128
				fn.fd = (*Decoder).kComplex128
			case reflect.Chan:
				fn.fe = (*Encoder).kChan
				fn.fd = (*Decoder).kChan
			case reflect.Slice:
				fn.fe = (*Encoder).kSlice
				fn.fd = (*Decoder).kSlice
			case reflect.Array:
				fi.addrD = false 
				fn.fe = (*Encoder).kArray
				fn.fd = (*Decoder).kArray
			case reflect.Struct:
				if ti.anyOmitEmpty ||
					ti.flagMissingFielder ||
					ti.flagMissingFielderPtr {
					fn.fe = (*Encoder).kStruct
				} else {
					fn.fe = (*Encoder).kStructNoOmitempty
				}
				fn.fd = (*Decoder).kStruct
			case reflect.Map:
				fn.fe = (*Encoder).kMap
				fn.fd = (*Decoder).kMap
			case reflect.Interface:
				
				fn.fd = (*Decoder).kInterface
				fn.fe = (*Encoder).kErr
			default:
				
				fn.fe = (*Encoder).kErr
				fn.fd = (*Decoder).kErr
			}
		}
	}
	return
}
















type Handle interface {
	Name() string
	getBasicHandle() *BasicHandle
	newEncDriver() encDriver
	newDecDriver() decDriver
	isBinary() bool
	isJson() bool 
	
	desc(bd byte) string
	
	init()
}





type Raw []byte







type RawExt struct {
	Tag uint64
	
	
	Data []byte
	
	
	
	Value interface{}
}

func (re *RawExt) setData(xbs []byte, zerocopy bool) {
	if zerocopy {
		re.Data = xbs
	} else {
		re.Data = append(re.Data[:0], xbs...)
	}
}



type BytesExt interface {
	
	
	
	WriteExt(v interface{}) []byte

	
	
	
	ReadExt(dst interface{}, src []byte)
}





type InterfaceExt interface {
	
	
	
	
	ConvertExt(v interface{}) interface{}

	
	
	
	
	UpdateExt(dst interface{}, src interface{})
}


type Ext interface {
	BytesExt
	InterfaceExt
}


type addExtWrapper struct {
	encFn func(reflect.Value) ([]byte, error)
	decFn func(reflect.Value, []byte) error
}

func (x addExtWrapper) WriteExt(v interface{}) []byte {
	bs, err := x.encFn(reflect.ValueOf(v))
	halt.onerror(err)
	return bs
}

func (x addExtWrapper) ReadExt(v interface{}, bs []byte) {
	halt.onerror(x.decFn(reflect.ValueOf(v), bs))
}

func (x addExtWrapper) ConvertExt(v interface{}) interface{} {
	return x.WriteExt(v)
}

func (x addExtWrapper) UpdateExt(dest interface{}, v interface{}) {
	x.ReadExt(dest, v.([]byte))
}

type bytesExtFailer struct{}

func (bytesExtFailer) WriteExt(v interface{}) []byte {
	halt.onerror(errExtFnWriteExtUnsupported)
	return nil
}
func (bytesExtFailer) ReadExt(v interface{}, bs []byte) {
	halt.onerror(errExtFnReadExtUnsupported)
}

type interfaceExtFailer struct{}

func (interfaceExtFailer) ConvertExt(v interface{}) interface{} {
	halt.onerror(errExtFnConvertExtUnsupported)
	return nil
}
func (interfaceExtFailer) UpdateExt(dest interface{}, v interface{}) {
	halt.onerror(errExtFnUpdateExtUnsupported)
}

type bytesExtWrapper struct {
	interfaceExtFailer
	BytesExt
}

type interfaceExtWrapper struct {
	bytesExtFailer
	InterfaceExt
}

type extFailWrapper struct {
	bytesExtFailer
	interfaceExtFailer
}

type binaryEncodingType struct{}

func (binaryEncodingType) isBinary() bool { return true }
func (binaryEncodingType) isJson() bool   { return false }

type textEncodingType struct{}

func (textEncodingType) isBinary() bool { return false }
func (textEncodingType) isJson() bool   { return false }

type notJsonType struct{}

func (notJsonType) isJson() bool { return false }




type noBuiltInTypes struct{}

func (noBuiltInTypes) EncodeBuiltin(rt uintptr, v interface{}) {}
func (noBuiltInTypes) DecodeBuiltin(rt uintptr, v interface{}) {}













type bigenHelper struct{}

func (z bigenHelper) PutUint16(v uint16) (b [2]byte) {
	return [...]byte{
		byte(v >> 8),
		byte(v),
	}
}

func (z bigenHelper) PutUint32(v uint32) (b [4]byte) {
	return [...]byte{
		byte(v >> 24),
		byte(v >> 16),
		byte(v >> 8),
		byte(v),
	}
}

func (z bigenHelper) PutUint64(v uint64) (b [8]byte) {
	return [...]byte{
		byte(v >> 56),
		byte(v >> 48),
		byte(v >> 40),
		byte(v >> 32),
		byte(v >> 24),
		byte(v >> 16),
		byte(v >> 8),
		byte(v),
	}
}

func (z bigenHelper) Uint16(b [2]byte) (v uint16) {
	return uint16(b[1]) |
		uint16(b[0])<<8
}

func (z bigenHelper) Uint32(b [4]byte) (v uint32) {
	return uint32(b[3]) |
		uint32(b[2])<<8 |
		uint32(b[1])<<16 |
		uint32(b[0])<<24
}

func (z bigenHelper) Uint64(b [8]byte) (v uint64) {
	return uint64(b[7]) |
		uint64(b[6])<<8 |
		uint64(b[5])<<16 |
		uint64(b[4])<<24 |
		uint64(b[3])<<32 |
		uint64(b[2])<<40 |
		uint64(b[1])<<48 |
		uint64(b[0])<<56
}

func (z bigenHelper) writeUint16(w *encWr, v uint16) {
	x := z.PutUint16(v)
	w.writen2(x[0], x[1])
}

func (z bigenHelper) writeUint32(w *encWr, v uint32) {
	
	
	
	
	w.writen4(z.PutUint32(v))
}

func (z bigenHelper) writeUint64(w *encWr, v uint64) {
	w.writen8(z.PutUint64(v))
}

type extTypeTagFn struct {
	rtid    uintptr
	rtidptr uintptr
	rt      reflect.Type
	tag     uint64
	ext     Ext
}

type extHandle []extTypeTagFn





func (x *BasicHandle) AddExt(rt reflect.Type, tag byte,
	encfn func(reflect.Value) ([]byte, error),
	decfn func(reflect.Value, []byte) error) (err error) {
	if encfn == nil || decfn == nil {
		return x.SetExt(rt, uint64(tag), nil)
	}
	return x.SetExt(rt, uint64(tag), addExtWrapper{encfn, decfn})
}







func (x *BasicHandle) SetExt(rt reflect.Type, tag uint64, ext Ext) (err error) {
	if x.isInited() {
		return errHandleInited
	}
	if x.basicHandleRuntimeState == nil {
		x.basicHandleRuntimeState = new(basicHandleRuntimeState)
	}
	return x.basicHandleRuntimeState.setExt(rt, tag, ext)
}

func (o extHandle) getExtForI(x interface{}) (v *extTypeTagFn) {
	if len(o) > 0 {
		v = o.getExt(i2rtid(x), true)
	}
	return
}

func (o extHandle) getExt(rtid uintptr, check bool) (v *extTypeTagFn) {
	if !check {
		return
	}
	for i := range o {
		v = &o[i]
		if v.rtid == rtid || v.rtidptr == rtid {
			return
		}
	}
	return nil
}

func (o extHandle) getExtForTag(tag uint64) (v *extTypeTagFn) {
	for i := range o {
		v = &o[i]
		if v.tag == tag {
			return
		}
	}
	return nil
}

type intf2impl struct {
	rtid uintptr 
	impl reflect.Type
}

type intf2impls []intf2impl







func (o *intf2impls) Intf2Impl(intf, impl reflect.Type) (err error) {
	if impl != nil && !impl.Implements(intf) {
		return fmt.Errorf("Intf2Impl: %v does not implement %v", impl, intf)
	}
	rtid := rt2id(intf)
	o2 := *o
	for i := range o2 {
		v := &o2[i]
		if v.rtid == rtid {
			v.impl = impl
			return
		}
	}
	*o = append(o2, intf2impl{rtid, impl})
	return
}

func (o intf2impls) intf2impl(rtid uintptr) (rv reflect.Value) {
	for i := range o {
		v := &o[i]
		if v.rtid == rtid {
			if v.impl == nil {
				return
			}
			vkind := v.impl.Kind()
			if vkind == reflect.Ptr {
				return reflect.New(v.impl.Elem())
			}
			return rvZeroAddrK(v.impl, vkind)
		}
	}
	return
}






type structFieldInfoPathNode struct {
	parent *structFieldInfoPathNode

	offset   uint16
	index    uint16
	kind     uint8
	numderef uint8

	
	

	encNameAsciiAlphaNum bool 
	omitEmpty            bool

	typ reflect.Type
}


func (path *structFieldInfoPathNode) depth() (d int) {
TOP:
	if path != nil {
		d++
		path = path.parent
		goto TOP
	}
	return
}


func (path *structFieldInfoPathNode) field(v reflect.Value) (rv2 reflect.Value) {
	if parent := path.parent; parent != nil {
		v = parent.field(v)
		for j, k := uint8(0), parent.numderef; j < k; j++ {
			if rvIsNil(v) {
				return
			}
			v = v.Elem()
		}
	}
	return path.rvField(v)
}



func (path *structFieldInfoPathNode) fieldAlloc(v reflect.Value) (rv2 reflect.Value) {
	if parent := path.parent; parent != nil {
		v = parent.fieldAlloc(v)
		for j, k := uint8(0), parent.numderef; j < k; j++ {
			if rvIsNil(v) {
				rvSetDirect(v, reflect.New(v.Type().Elem()))
			}
			v = v.Elem()
		}
	}
	return path.rvField(v)
}

type structFieldInfo struct {
	encName string 

	

	

	
	

	path structFieldInfoPathNode
}

func parseStructInfo(stag string) (toArray, omitEmpty bool, keytype valueType) {
	keytype = valueTypeString 
	if stag == "" {
		return
	}
	ss := strings.Split(stag, ",")
	if len(ss) < 2 {
		return
	}
	for _, s := range ss[1:] {
		switch s {
		case "omitempty":
			omitEmpty = true
		case "toarray":
			toArray = true
		case "int":
			keytype = valueTypeInt
		case "uint":
			keytype = valueTypeUint
		case "float":
			keytype = valueTypeFloat
			
			
		case "string":
			keytype = valueTypeString
		}
	}
	return
}

func (si *structFieldInfo) parseTag(stag string) {
	if stag == "" {
		return
	}
	for i, s := range strings.Split(stag, ",") {
		if i == 0 {
			if s != "" {
				si.encName = s
			}
		} else {
			switch s {
			case "omitempty":
				si.path.omitEmpty = true
			}
		}
	}
}

type sfiSortedByEncName []*structFieldInfo

func (p sfiSortedByEncName) Len() int           { return len(p) }
func (p sfiSortedByEncName) Swap(i, j int)      { p[uint(i)], p[uint(j)] = p[uint(j)], p[uint(i)] }
func (p sfiSortedByEncName) Less(i, j int) bool { return p[uint(i)].encName < p[uint(j)].encName }



type typeInfo4Container struct {
	elem reflect.Type
	
	
	
	
	key reflect.Type

	
	
	
	
	fastpathUnderlying reflect.Type

	tikey  *typeInfo
	tielem *typeInfo
}










type typeInfo struct {
	rt  reflect.Type
	ptr reflect.Type

	

	rtid uintptr

	numMeth uint16 
	kind    uint8
	chandir uint8

	anyOmitEmpty bool      
	toArray      bool      
	keyType      valueType 
	mbs          bool      

	sfi4Name map[string]*structFieldInfo 

	*typeInfo4Container

	

	size, keysize, elemsize uint32

	keykind, elemkind uint8

	flagHasPkgPath   bool 
	flagComparable   bool
	flagCanTransient bool

	flagMarshalInterface  bool 
	flagSelferViaCodecgen bool

	
	flagIsZeroer    bool
	flagIsZeroerPtr bool

	flagIsCodecEmptyer    bool
	flagIsCodecEmptyerPtr bool

	flagBinaryMarshaler    bool
	flagBinaryMarshalerPtr bool

	flagBinaryUnmarshaler    bool
	flagBinaryUnmarshalerPtr bool

	flagTextMarshaler    bool
	flagTextMarshalerPtr bool

	flagTextUnmarshaler    bool
	flagTextUnmarshalerPtr bool

	flagJsonMarshaler    bool
	flagJsonMarshalerPtr bool

	flagJsonUnmarshaler    bool
	flagJsonUnmarshalerPtr bool

	flagSelfer    bool
	flagSelferPtr bool

	flagMissingFielder    bool
	flagMissingFielderPtr bool

	infoFieldOmitempty bool

	sfi structFieldInfos
}

func (ti *typeInfo) siForEncName(name []byte) (si *structFieldInfo) {
	return ti.sfi4Name[string(name)]
}

func (ti *typeInfo) resolve(x []structFieldInfo, ss map[string]uint16) (n int) {
	n = len(x)

	for i := range x {
		ui := uint16(i)
		xn := x[i].encName
		j, ok := ss[xn]
		if ok {
			i2clear := ui                              
			if x[i].path.depth() < x[j].path.depth() { 
				ss[xn] = ui
				i2clear = j
			}
			if x[i2clear].encName != "" {
				x[i2clear].encName = ""
				n--
			}
		} else {
			ss[xn] = ui
		}
	}

	return
}

func (ti *typeInfo) init(x []structFieldInfo, n int) {
	var anyOmitEmpty bool

	
	m := make(map[string]*structFieldInfo, n)
	w := make([]structFieldInfo, n)
	y := make([]*structFieldInfo, n+n)
	z := y[n:]
	y = y[:n]
	n = 0
	for i := range x {
		if x[i].encName == "" {
			continue
		}
		if !anyOmitEmpty && x[i].path.omitEmpty {
			anyOmitEmpty = true
		}
		w[n] = x[i]
		y[n] = &w[n]
		m[x[i].encName] = &w[n]
		n++
	}
	if n != len(y) {
		halt.errorf("failure reading struct %v - expecting %d of %d valid fields, got %d", ti.rt, len(y), len(x), n)
	}

	copy(z, y)
	sort.Sort(sfiSortedByEncName(z))

	ti.anyOmitEmpty = anyOmitEmpty
	ti.sfi.load(y, z)
	ti.sfi4Name = m
}














func transientBitsetFlags() *bitset32 {
	if transientValueHasStringSlice {
		return &numBoolStrSliceBitset
	}
	return &numBoolBitset
}

func isCanTransient(t reflect.Type, k reflect.Kind) (v bool) {
	var bs = transientBitsetFlags()
	if bs.isset(byte(k)) {
		v = true
	} else if k == reflect.Slice {
		elem := t.Elem()
		v = numBoolBitset.isset(byte(elem.Kind()))
	} else if k == reflect.Array {
		elem := t.Elem()
		v = isCanTransient(elem, elem.Kind())
	} else if k == reflect.Struct {
		v = true
		for j, jlen := 0, t.NumField(); j < jlen; j++ {
			f := t.Field(j)
			if !isCanTransient(f.Type, f.Type.Kind()) {
				v = false
				return
			}
		}
	} else {
		v = false
	}
	return
}

func (ti *typeInfo) doSetFlagCanTransient() {
	if transientSizeMax > 0 {
		ti.flagCanTransient = ti.size <= transientSizeMax
	} else {
		ti.flagCanTransient = true
	}
	if ti.flagCanTransient {
		if !transientBitsetFlags().isset(ti.kind) {
			ti.flagCanTransient = isCanTransient(ti.rt, reflect.Kind(ti.kind))
		}
	}
}

type rtid2ti struct {
	rtid uintptr
	ti   *typeInfo
}





type TypeInfos struct {
	infos atomicTypeInfoSlice
	mu    sync.Mutex
	_     uint64 
	tags  []string
	_     uint64 
}





func NewTypeInfos(tags []string) *TypeInfos {
	return &TypeInfos{tags: tags}
}

func (x *TypeInfos) structTag(t reflect.StructTag) (s string) {
	
	
	for _, x := range x.tags {
		s = t.Get(x)
		if s != "" {
			return s
		}
	}
	return
}

func findTypeInfo(s []rtid2ti, rtid uintptr) (i uint, ti *typeInfo) {
	
	

	var h uint
	var j = uint(len(s))
LOOP:
	if i < j {
		h = (i + j) >> 1 
		if s[h].rtid < rtid {
			i = h + 1
		} else {
			j = h
		}
		goto LOOP
	}
	if i < uint(len(s)) && s[i].rtid == rtid {
		ti = s[i].ti
	}
	return
}

func (x *TypeInfos) get(rtid uintptr, rt reflect.Type) (pti *typeInfo) {
	if pti = x.find(rtid); pti == nil {
		pti = x.load(rt)
	}
	return
}

func (x *TypeInfos) find(rtid uintptr) (pti *typeInfo) {
	sp := x.infos.load()
	if sp != nil {
		_, pti = findTypeInfo(sp, rtid)
	}
	return
}

func (x *TypeInfos) load(rt reflect.Type) (pti *typeInfo) {
	rk := rt.Kind()

	if rk == reflect.Ptr { 
		halt.errorf("invalid kind passed to TypeInfos.get: %v - %v", rk, rt)
	}

	rtid := rt2id(rt)

	
	
	ti := typeInfo{
		rt:      rt,
		ptr:     reflect.PtrTo(rt),
		rtid:    rtid,
		kind:    uint8(rk),
		size:    uint32(rt.Size()),
		numMeth: uint16(rt.NumMethod()),
		keyType: valueTypeString, 

		
		flagHasPkgPath: rt.PkgPath() != "",
	}

	
	bset := func(when bool, b *bool) {
		if when {
			*b = true
		}
	}

	var b1, b2 bool

	b1, b2 = implIntf(rt, binaryMarshalerTyp)
	bset(b1, &ti.flagBinaryMarshaler)
	bset(b2, &ti.flagBinaryMarshalerPtr)
	b1, b2 = implIntf(rt, binaryUnmarshalerTyp)
	bset(b1, &ti.flagBinaryUnmarshaler)
	bset(b2, &ti.flagBinaryUnmarshalerPtr)
	b1, b2 = implIntf(rt, textMarshalerTyp)
	bset(b1, &ti.flagTextMarshaler)
	bset(b2, &ti.flagTextMarshalerPtr)
	b1, b2 = implIntf(rt, textUnmarshalerTyp)
	bset(b1, &ti.flagTextUnmarshaler)
	bset(b2, &ti.flagTextUnmarshalerPtr)
	b1, b2 = implIntf(rt, jsonMarshalerTyp)
	bset(b1, &ti.flagJsonMarshaler)
	bset(b2, &ti.flagJsonMarshalerPtr)
	b1, b2 = implIntf(rt, jsonUnmarshalerTyp)
	bset(b1, &ti.flagJsonUnmarshaler)
	bset(b2, &ti.flagJsonUnmarshalerPtr)
	b1, b2 = implIntf(rt, selferTyp)
	bset(b1, &ti.flagSelfer)
	bset(b2, &ti.flagSelferPtr)
	b1, b2 = implIntf(rt, missingFielderTyp)
	bset(b1, &ti.flagMissingFielder)
	bset(b2, &ti.flagMissingFielderPtr)
	b1, b2 = implIntf(rt, iszeroTyp)
	bset(b1, &ti.flagIsZeroer)
	bset(b2, &ti.flagIsZeroerPtr)
	b1, b2 = implIntf(rt, isCodecEmptyerTyp)
	bset(b1, &ti.flagIsCodecEmptyer)
	bset(b2, &ti.flagIsCodecEmptyerPtr)

	b1, b2 = implIntf(rt, isSelferViaCodecgenerTyp)
	ti.flagSelferViaCodecgen = b1 || b2

	ti.flagMarshalInterface = ti.flagSelfer || ti.flagSelferPtr ||
		ti.flagSelferViaCodecgen ||
		ti.flagBinaryMarshaler || ti.flagBinaryMarshalerPtr ||
		ti.flagBinaryUnmarshaler || ti.flagBinaryUnmarshalerPtr ||
		ti.flagTextMarshaler || ti.flagTextMarshalerPtr ||
		ti.flagTextUnmarshaler || ti.flagTextUnmarshalerPtr ||
		ti.flagJsonMarshaler || ti.flagJsonMarshalerPtr ||
		ti.flagJsonUnmarshaler || ti.flagJsonUnmarshalerPtr

	b1 = rt.Comparable()
	
	ti.flagComparable = b1

	ti.doSetFlagCanTransient()

	var tt reflect.Type
	switch rk {
	case reflect.Struct:
		var omitEmpty bool
		if f, ok := rt.FieldByName(structInfoFieldName); ok {
			ti.toArray, omitEmpty, ti.keyType = parseStructInfo(x.structTag(f.Tag))
			ti.infoFieldOmitempty = omitEmpty
		} else {
			ti.keyType = valueTypeString
		}
		pp, pi := &pool4tiload, pool4tiload.Get()
		pv := pi.(*typeInfoLoad)
		pv.reset()
		pv.etypes = append(pv.etypes, ti.rtid)
		x.rget(rt, rtid, nil, pv, omitEmpty)
		n := ti.resolve(pv.sfis, pv.sfiNames)
		ti.init(pv.sfis, n)
		pp.Put(pi)
	case reflect.Map:
		ti.typeInfo4Container = new(typeInfo4Container)
		ti.elem = rt.Elem()
		for tt = ti.elem; tt.Kind() == reflect.Ptr; tt = tt.Elem() {
		}
		ti.tielem = x.get(rt2id(tt), tt)
		ti.elemkind = uint8(ti.elem.Kind())
		ti.elemsize = uint32(ti.elem.Size())
		ti.key = rt.Key()
		for tt = ti.key; tt.Kind() == reflect.Ptr; tt = tt.Elem() {
		}
		ti.tikey = x.get(rt2id(tt), tt)
		ti.keykind = uint8(ti.key.Kind())
		ti.keysize = uint32(ti.key.Size())
		if ti.flagHasPkgPath {
			ti.fastpathUnderlying = reflect.MapOf(ti.key, ti.elem)
		}
	case reflect.Slice:
		ti.typeInfo4Container = new(typeInfo4Container)
		ti.mbs, b2 = implIntf(rt, mapBySliceTyp)
		if !ti.mbs && b2 {
			ti.mbs = b2
		}
		ti.elem = rt.Elem()
		for tt = ti.elem; tt.Kind() == reflect.Ptr; tt = tt.Elem() {
		}
		ti.tielem = x.get(rt2id(tt), tt)
		ti.elemkind = uint8(ti.elem.Kind())
		ti.elemsize = uint32(ti.elem.Size())
		if ti.flagHasPkgPath {
			ti.fastpathUnderlying = reflect.SliceOf(ti.elem)
		}
	case reflect.Chan:
		ti.typeInfo4Container = new(typeInfo4Container)
		ti.elem = rt.Elem()
		for tt = ti.elem; tt.Kind() == reflect.Ptr; tt = tt.Elem() {
		}
		ti.tielem = x.get(rt2id(tt), tt)
		ti.elemkind = uint8(ti.elem.Kind())
		ti.elemsize = uint32(ti.elem.Size())
		ti.chandir = uint8(rt.ChanDir())
		ti.key = reflect.SliceOf(ti.elem)
		ti.keykind = uint8(reflect.Slice)
	case reflect.Array:
		ti.typeInfo4Container = new(typeInfo4Container)
		ti.mbs, b2 = implIntf(rt, mapBySliceTyp)
		if !ti.mbs && b2 {
			ti.mbs = b2
		}
		ti.elem = rt.Elem()
		ti.elemkind = uint8(ti.elem.Kind())
		ti.elemsize = uint32(ti.elem.Size())
		for tt = ti.elem; tt.Kind() == reflect.Ptr; tt = tt.Elem() {
		}
		ti.tielem = x.get(rt2id(tt), tt)
		ti.key = reflect.SliceOf(ti.elem)
		ti.keykind = uint8(reflect.Slice)
		ti.keysize = uint32(ti.key.Size())
		if ti.flagHasPkgPath {
			ti.fastpathUnderlying = ti.key
		}

		
		
		
		
		
	}

	x.mu.Lock()
	sp := x.infos.load()
	
	
	if sp == nil {
		pti = &ti
		sp = []rtid2ti{{rtid, pti}}
		x.infos.store(sp)
	} else {
		var idx uint
		idx, pti = findTypeInfo(sp, rtid)
		if pti == nil {
			pti = &ti
			sp2 := make([]rtid2ti, len(sp)+1)
			copy(sp2[idx+1:], sp[idx:])
			copy(sp2, sp[:idx])
			sp2[idx] = rtid2ti{rtid, pti}
			x.infos.store(sp2)
		}
	}
	x.mu.Unlock()
	return
}

func (x *TypeInfos) rget(rt reflect.Type, rtid uintptr,
	path *structFieldInfoPathNode, pv *typeInfoLoad, omitEmpty bool) {
	
	
	
	
	
	
	
	
	flen := rt.NumField()
LOOP:
	for j, jlen := uint16(0), uint16(flen); j < jlen; j++ {
		f := rt.Field(int(j))
		fkind := f.Type.Kind()

		
		switch fkind {
		case reflect.Func, reflect.UnsafePointer:
			continue LOOP
		}

		isUnexported := f.PkgPath != ""
		if isUnexported && !f.Anonymous {
			continue
		}
		stag := x.structTag(f.Tag)
		if stag == "-" {
			continue
		}
		var si structFieldInfo

		var numderef uint8 = 0
		for xft := f.Type; xft.Kind() == reflect.Ptr; xft = xft.Elem() {
			numderef++
		}

		var parsed bool
		
		
		if f.Anonymous && fkind != reflect.Interface {
			
			ft := f.Type
			isPtr := ft.Kind() == reflect.Ptr
			for ft.Kind() == reflect.Ptr {
				ft = ft.Elem()
			}
			isStruct := ft.Kind() == reflect.Struct

			
			
			
			
			if (isUnexported && !isStruct) || (!allowSetUnexportedEmbeddedPtr && isUnexported && isPtr) {
				continue
			}
			doInline := stag == ""
			if !doInline {
				si.parseTag(stag)
				parsed = true
				doInline = si.encName == "" 
			}
			if doInline && isStruct {
				
				ftid := rt2id(ft)
				
				
				
				
				processIt := true
				numk := 0
				for _, k := range pv.etypes {
					if k == ftid {
						numk++
						if numk == rgetMaxRecursion {
							processIt = false
							break
						}
					}
				}
				if processIt {
					pv.etypes = append(pv.etypes, ftid)
					path2 := &structFieldInfoPathNode{
						parent:   path,
						typ:      f.Type,
						offset:   uint16(f.Offset),
						index:    j,
						kind:     uint8(fkind),
						numderef: numderef,
					}
					x.rget(ft, ftid, path2, pv, omitEmpty)
				}
				continue
			}
		}

		
		if isUnexported || f.Name == "" { 
			continue
		}

		si.path = structFieldInfoPathNode{
			parent:   path,
			typ:      f.Type,
			offset:   uint16(f.Offset),
			index:    j,
			kind:     uint8(fkind),
			numderef: numderef,
			
			encNameAsciiAlphaNum: true,
			
			omitEmpty: si.path.omitEmpty,
		}

		if !parsed {
			si.encName = f.Name
			si.parseTag(stag)
			parsed = true
		} else if si.encName == "" {
			si.encName = f.Name
		}

		

		if omitEmpty {
			si.path.omitEmpty = true
		}

		for i := len(si.encName) - 1; i >= 0; i-- { 
			if !asciiAlphaNumBitset.isset(si.encName[i]) {
				si.path.encNameAsciiAlphaNum = false
				break
			}
		}

		pv.sfis = append(pv.sfis, si)
	}
}

func implIntf(rt, iTyp reflect.Type) (base bool, indir bool) {
	

	
	

	
	
	

	base = rt.Implements(iTyp)
	if base {
		indir = true
	} else {
		indir = reflect.PtrTo(rt).Implements(iTyp)
	}
	return
}

func bool2int(b bool) (v uint8) {
	
	if b {
		v = 1
	}
	return
}

func isSliceBoundsError(s string) bool {
	return strings.Contains(s, "index out of range") ||
		strings.Contains(s, "slice bounds out of range")
}

func sprintf(format string, v ...interface{}) string {
	return fmt.Sprintf(format, v...)
}

func panicValToErr(h errDecorator, v interface{}, err *error) {
	if v == *err {
		return
	}
	switch xerr := v.(type) {
	case nil:
	case runtime.Error:
		d, dok := h.(*Decoder)
		if dok && d.bytes && isSliceBoundsError(xerr.Error()) {
			*err = io.ErrUnexpectedEOF
		} else {
			h.wrapErr(xerr, err)
		}
	case error:
		switch xerr {
		case nil:
		case io.EOF, io.ErrUnexpectedEOF, errEncoderNotInitialized, errDecoderNotInitialized:
			
			*err = xerr
		default:
			h.wrapErr(xerr, err)
		}
	default:
		
		h.wrapErr(fmt.Errorf("%v", v), err)
	}
}

func usableByteSlice(bs []byte, slen int) (out []byte, changed bool) {
	const maxCap = 1024 * 1024 * 64 
	const skipMaxCap = false        
	if slen <= 0 {
		return []byte{}, true
	}
	if slen <= cap(bs) {
		return bs[:slen], false
	}
	
	if skipMaxCap || slen <= maxCap {
		return make([]byte, slen), true
	}
	return make([]byte, maxCap), true
}

func mapKeyFastKindFor(k reflect.Kind) mapKeyFastKind {
	return mapKeyFastKindVals[k&31]
}



type codecFnInfo struct {
	ti     *typeInfo
	xfFn   Ext
	xfTag  uint64
	addrD  bool
	addrDf bool 
	addrE  bool
	
}





type codecFn struct {
	i  codecFnInfo
	fe func(*Encoder, *codecFnInfo, reflect.Value)
	fd func(*Decoder, *codecFnInfo, reflect.Value)
	
}

type codecRtidFn struct {
	rtid uintptr
	fn   *codecFn
}

func makeExt(ext interface{}) Ext {
	switch t := ext.(type) {
	case Ext:
		return t
	case BytesExt:
		return &bytesExtWrapper{BytesExt: t}
	case InterfaceExt:
		return &interfaceExtWrapper{InterfaceExt: t}
	}
	return &extFailWrapper{}
}

func baseRV(v interface{}) (rv reflect.Value) {
	
	for rv = reflect.ValueOf(v); rv.Kind() == reflect.Ptr; rv = rv.Elem() {
	}
	return
}








type checkOverflow struct{}

func (checkOverflow) Float32(v float64) (overflow bool) {
	if v < 0 {
		v = -v
	}
	return math.MaxFloat32 < v && v <= math.MaxFloat64
}
func (checkOverflow) Uint(v uint64, bitsize uint8) (overflow bool) {
	if v != 0 && v != (v<<(64-bitsize))>>(64-bitsize) {
		overflow = true
	}
	return
}
func (checkOverflow) Int(v int64, bitsize uint8) (overflow bool) {
	if v != 0 && v != (v<<(64-bitsize))>>(64-bitsize) {
		overflow = true
	}
	return
}

func (checkOverflow) Uint2Int(v uint64, neg bool) (overflow bool) {
	return (neg && v > 1<<63) || (!neg && v >= 1<<63)
}

func (checkOverflow) SignedInt(v uint64) (overflow bool) {
	
	
	
	
	
	
	
	
	
	
	
	

	
	
	overflow = (v>>63) != 0 && v&0x7fffffffffffffff > math.MaxInt64-1

	return
}

func (x checkOverflow) Float32V(v float64) float64 {
	if x.Float32(v) {
		halt.errorf("float32 overflow: %v", v)
	}
	return v
}
func (x checkOverflow) UintV(v uint64, bitsize uint8) uint64 {
	if x.Uint(v, bitsize) {
		halt.errorf("uint64 overflow: %v", v)
	}
	return v
}
func (x checkOverflow) IntV(v int64, bitsize uint8) int64 {
	if x.Int(v, bitsize) {
		halt.errorf("int64 overflow: %v", v)
	}
	return v
}
func (x checkOverflow) SignedIntV(v uint64) int64 {
	if x.SignedInt(v) {
		halt.errorf("uint64 to int64 overflow: %v", v)
	}
	return int64(v)
}



func isNaN64(f float64) bool { return f != f }

func isWhitespaceChar(v byte) bool {
	

	return v < 33
	
	
	
	
}

func isNumberChar(v byte) bool {
	

	return numCharBitset.isset(v)
	
	
}



type ioFlusher interface {
	Flush() error
}

type ioBuffered interface {
	Buffered() int
}



type sfiRv struct {
	v *structFieldInfo
	r reflect.Value
}























type bitset32 [32]bool

func (x *bitset32) set(pos byte) *bitset32 {
	x[pos&31] = true 
	return x
}
func (x *bitset32) isset(pos byte) bool {
	return x[pos&31] 
}

type bitset256 [256]bool

func (x *bitset256) set(pos byte) *bitset256 {
	x[pos] = true
	return x
}
func (x *bitset256) isset(pos byte) bool {
	return x[pos]
}



type panicHdl struct{}


func (panicHdl) onerror(err error) {
	if err != nil {
		panic(err)
	}
}






//go:noinline
func (panicHdl) errorf(format string, params ...interface{}) {
	if format == "" {
		panic(errPanicUndefined)
	}
	if len(params) == 0 {
		panic(errors.New(format))
	}
	panic(fmt.Errorf(format, params...))
}



type errDecorator interface {
	wrapErr(in error, out *error)
}

type errDecoratorDef struct{}

func (errDecoratorDef) wrapErr(v error, e *error) { *e = v }



type mustHdl struct{}

func (mustHdl) String(s string, err error) string {
	halt.onerror(err)
	return s
}
func (mustHdl) Int(s int64, err error) int64 {
	halt.onerror(err)
	return s
}
func (mustHdl) Uint(s uint64, err error) uint64 {
	halt.onerror(err)
	return s
}
func (mustHdl) Float(s float64, err error) float64 {
	halt.onerror(err)
	return s
}



func freelistCapacity(length int) (capacity int) {
	for capacity = 8; capacity <= length; capacity *= 2 {
	}
	return
}

























type bytesFreelist [][]byte



func (x *bytesFreelist) peek(length int, pop bool) (out []byte) {
	if bytesFreeListNoCache {
		return make([]byte, 0, freelistCapacity(length))
	}
	y := *x
	if len(y) > 0 {
		out = y[len(y)-1]
	}
	
	const minLenBytes = 64
	if length < minLenBytes {
		length = minLenBytes
	}
	if cap(out) < length {
		out = make([]byte, 0, freelistCapacity(length))
		y = append(y, out)
		*x = y
	}
	if pop && len(y) > 0 {
		y = y[:len(y)-1]
		*x = y
	}
	return
}



func (x *bytesFreelist) get(length int) (out []byte) {
	if bytesFreeListNoCache {
		return make([]byte, 0, freelistCapacity(length))
	}
	y := *x
	
	
	for i := 0; i < len(y); i++ {
		v := y[i]
		if cap(v) >= length {
			
			copy(y[i:], y[i+1:])
			*x = y[:len(y)-1]
			return v
		}
	}
	return make([]byte, 0, freelistCapacity(length))
}

func (x *bytesFreelist) put(v []byte) {
	if bytesFreeListNoCache || cap(v) == 0 {
		return
	}
	if len(v) != 0 {
		v = v[:0]
	}
	
	y := append(*x, v)
	*x = y
	
	
	for i := 0; i < len(y)-1; i++ {
		z := y[i]
		if cap(z) > cap(v) {
			copy(y[i+1:], y[i:])
			y[i] = v
			return
		}
	}
}

func (x *bytesFreelist) check(v []byte, length int) (out []byte) {
	
	if cap(v) >= length {
		return v[:0]
	}
	return x.checkPutGet(v, length)
}

func (x *bytesFreelist) checkPutGet(v []byte, length int) []byte {
	
	const useSeparateCalls = false

	if useSeparateCalls {
		x.put(v)
		return x.get(length)
	}

	if bytesFreeListNoCache {
		return make([]byte, 0, freelistCapacity(length))
	}

	
	y := *x
	var put = cap(v) == 0 
	if !put {
		y = append(y, v)
		*x = y
	}
	for i := 0; i < len(y); i++ {
		z := y[i]
		if put {
			if cap(z) >= length {
				copy(y[i:], y[i+1:])
				y = y[:len(y)-1]
				*x = y
				return z
			}
		} else {
			if cap(z) > cap(v) {
				copy(y[i+1:], y[i:])
				y[i] = v
				put = true
			}
		}
	}
	return make([]byte, 0, freelistCapacity(length))
}












type sfiRvFreelist [][]sfiRv

func (x *sfiRvFreelist) get(length int) (out []sfiRv) {
	y := *x

	
	
	for i := 0; i < len(y); i++ {
		v := y[i]
		if cap(v) >= length {
			
			copy(y[i:], y[i+1:])
			*x = y[:len(y)-1]
			return v
		}
	}
	return make([]sfiRv, 0, freelistCapacity(length))
}

func (x *sfiRvFreelist) put(v []sfiRv) {
	if len(v) != 0 {
		v = v[:0]
	}
	
	y := append(*x, v)
	*x = y
	
	
	for i := 0; i < len(y)-1; i++ {
		z := y[i]
		if cap(z) > cap(v) {
			copy(y[i+1:], y[i:])
			y[i] = v
			return
		}
	}
}








const (
	internMaxStrLen = 16     
	internCap       = 64 * 2 
)

type internerMap map[string]string

func (x *internerMap) init() {
	*x = make(map[string]string, internCap)
}

func (x internerMap) string(v []byte) (s string) {
	s, ok := x[string(v)] 
	if !ok {
		s = string(v) 
		x[s] = s
	}
	return
}
