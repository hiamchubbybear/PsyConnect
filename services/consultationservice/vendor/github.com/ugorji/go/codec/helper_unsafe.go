


//go:build !safe && !codec.safe && !appengine && go1.9
// +build !safe,!codec.safe,!appengine,go1.9






package codec

import (
	"reflect"
	_ "runtime" 
	"sync/atomic"
	"time"
	"unsafe"
)





































const safeMode = false













const helperUnsafeDirectAssignMapEntry = true


const (
	unsafeFlagStickyRO = 1 << 5
	unsafeFlagEmbedRO  = 1 << 6
	unsafeFlagIndir    = 1 << 7
	unsafeFlagAddr     = 1 << 8
	unsafeFlagRO       = unsafeFlagStickyRO | unsafeFlagEmbedRO
	
	
)




const transientSizeMax = 64


const transientValueHasStringSlice = false

type unsafeString struct {
	Data unsafe.Pointer
	Len  int
}

type unsafeSlice struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}

type unsafeIntf struct {
	typ unsafe.Pointer
	ptr unsafe.Pointer
}

type unsafeReflectValue struct {
	unsafeIntf
	flag uintptr
}


type unsafeRuntimeType struct {
	size uintptr
	
}




var (
	unsafeZeroAddr  = unsafe.Pointer(&unsafeZeroArr[0])
	unsafeZeroSlice = unsafeSlice{unsafeZeroAddr, 0, 0}
)










type unsafePerTypeElem struct {
	arr   [transientSizeMax]byte 
	slice unsafeSlice            
}

func (x *unsafePerTypeElem) addrFor(k reflect.Kind) unsafe.Pointer {
	if k == reflect.String || k == reflect.Slice {
		x.slice = unsafeSlice{} 
		return unsafe.Pointer(&x.slice)
	}
	x.arr = [transientSizeMax]byte{} 
	return unsafe.Pointer(&x.arr)
}

type perType struct {
	elems [2]unsafePerTypeElem
}

type decPerType struct {
	perType
}

type encPerType struct{}






func (x *perType) TransientAddrK(t reflect.Type, k reflect.Kind) reflect.Value {
	return rvZeroAddrTransientAnyK(t, k, x.elems[0].addrFor(k))
}

func (x *perType) TransientAddr2K(t reflect.Type, k reflect.Kind) reflect.Value {
	return rvZeroAddrTransientAnyK(t, k, x.elems[1].addrFor(k))
}

func (encPerType) AddressableRO(v reflect.Value) reflect.Value {
	return rvAddressableReadonly(v)
}




func byteAt(b []byte, index uint) byte {
	
	return *(*byte)(unsafe.Pointer(uintptr((*unsafeSlice)(unsafe.Pointer(&b)).Data) + uintptr(index)))
}

func byteSliceOf(b []byte, start, end uint) []byte {
	s := (*unsafeSlice)(unsafe.Pointer(&b))
	s.Data = unsafe.Pointer(uintptr(s.Data) + uintptr(start))
	s.Len = int(end - start)
	s.Cap -= int(start)
	return b
}






func setByteAt(b []byte, index uint, val byte) {
	
	*(*byte)(unsafe.Pointer(uintptr((*unsafeSlice)(unsafe.Pointer(&b)).Data) + uintptr(index))) = val
}




func stringView(v []byte) string {
	return *(*string)(unsafe.Pointer(&v))
}




func bytesView(v string) (b []byte) {
	sx := (*unsafeString)(unsafe.Pointer(&v))
	bx := (*unsafeSlice)(unsafe.Pointer(&b))
	bx.Data, bx.Len, bx.Cap = sx.Data, sx.Len, sx.Len
	return
}

func byteSliceSameData(v1 []byte, v2 []byte) bool {
	return (*unsafeSlice)(unsafe.Pointer(&v1)).Data == (*unsafeSlice)(unsafe.Pointer(&v2)).Data
}





func okBytes2(b []byte) [2]byte {
	return *((*[2]byte)(((*unsafeSlice)(unsafe.Pointer(&b))).Data))
}

func okBytes3(b []byte) [3]byte {
	return *((*[3]byte)(((*unsafeSlice)(unsafe.Pointer(&b))).Data))
}

func okBytes4(b []byte) [4]byte {
	return *((*[4]byte)(((*unsafeSlice)(unsafe.Pointer(&b))).Data))
}

func okBytes8(b []byte) [8]byte {
	return *((*[8]byte)(((*unsafeSlice)(unsafe.Pointer(&b))).Data))
}




func isNil(v interface{}) (rv reflect.Value, isnil bool) {
	var ui = (*unsafeIntf)(unsafe.Pointer(&v))
	isnil = ui.ptr == nil
	if !isnil {
		rv, isnil = unsafeIsNilIntfOrSlice(ui, v)
	}
	return
}

func unsafeIsNilIntfOrSlice(ui *unsafeIntf, v interface{}) (rv reflect.Value, isnil bool) {
	rv = reflect.ValueOf(v) 
	tk := rv.Kind()
	isnil = (tk == reflect.Interface || tk == reflect.Slice) && *(*unsafe.Pointer)(ui.ptr) == nil
	return
}





func rvRefPtr(v *unsafeReflectValue) unsafe.Pointer {
	if v.flag&unsafeFlagIndir != 0 {
		return *(*unsafe.Pointer)(v.ptr)
	}
	return v.ptr
}

func eq4i(i0, i1 interface{}) bool {
	v0 := (*unsafeIntf)(unsafe.Pointer(&i0))
	v1 := (*unsafeIntf)(unsafe.Pointer(&i1))
	return v0.typ == v1.typ && v0.ptr == v1.ptr
}

func rv4iptr(i interface{}) (v reflect.Value) {
	
	uv := (*unsafeReflectValue)(unsafe.Pointer(&v))
	uv.unsafeIntf = *(*unsafeIntf)(unsafe.Pointer(&i))
	uv.flag = uintptr(rkindPtr)
	return
}

func rv4istr(i interface{}) (v reflect.Value) {
	
	uv := (*unsafeReflectValue)(unsafe.Pointer(&v))
	uv.unsafeIntf = *(*unsafeIntf)(unsafe.Pointer(&i))
	uv.flag = uintptr(rkindString) | unsafeFlagIndir
	return
}

func rv2i(rv reflect.Value) (i interface{}) {
	
	
	
	
	
	

	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	if refBitset.isset(byte(rv.Kind())) && urv.flag&unsafeFlagIndir != 0 {
		urv.ptr = *(*unsafe.Pointer)(urv.ptr)
	}
	return *(*interface{})(unsafe.Pointer(&urv.unsafeIntf))
}

func rvAddr(rv reflect.Value, ptrType reflect.Type) reflect.Value {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	urv.flag = (urv.flag & unsafeFlagRO) | uintptr(reflect.Ptr)
	urv.typ = ((*unsafeIntf)(unsafe.Pointer(&ptrType))).ptr
	return rv
}

func rvIsNil(rv reflect.Value) bool {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	if urv.flag&unsafeFlagIndir != 0 {
		return *(*unsafe.Pointer)(urv.ptr) == nil
	}
	return urv.ptr == nil
}

func rvSetSliceLen(rv reflect.Value, length int) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	(*unsafeString)(urv.ptr).Len = length
}

func rvZeroAddrK(t reflect.Type, k reflect.Kind) (rv reflect.Value) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	urv.typ = ((*unsafeIntf)(unsafe.Pointer(&t))).ptr
	urv.flag = uintptr(k) | unsafeFlagIndir | unsafeFlagAddr
	urv.ptr = unsafeNew(urv.typ)
	return
}

func rvZeroAddrTransientAnyK(t reflect.Type, k reflect.Kind, addr unsafe.Pointer) (rv reflect.Value) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	urv.typ = ((*unsafeIntf)(unsafe.Pointer(&t))).ptr
	urv.flag = uintptr(k) | unsafeFlagIndir | unsafeFlagAddr
	urv.ptr = addr
	return
}

func rvZeroK(t reflect.Type, k reflect.Kind) (rv reflect.Value) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	urv.typ = ((*unsafeIntf)(unsafe.Pointer(&t))).ptr
	if refBitset.isset(byte(k)) {
		urv.flag = uintptr(k)
	} else if rtsize2(urv.typ) <= uintptr(len(unsafeZeroArr)) {
		urv.flag = uintptr(k) | unsafeFlagIndir
		urv.ptr = unsafeZeroAddr
	} else { 
		urv.flag = uintptr(k) | unsafeFlagIndir | unsafeFlagAddr
		urv.ptr = unsafeNew(urv.typ)
	}
	return
}



func rvConvert(v reflect.Value, t reflect.Type) reflect.Value {
	uv := (*unsafeReflectValue)(unsafe.Pointer(&v))
	uv.typ = ((*unsafeIntf)(unsafe.Pointer(&t))).ptr
	return v
}









func rvAddressableReadonly(v reflect.Value) reflect.Value {
	
	
	
	

	uv := (*unsafeReflectValue)(unsafe.Pointer(&v))
	uv.flag = uv.flag | unsafeFlagAddr 

	return v
}

func rtsize2(rt unsafe.Pointer) uintptr {
	return ((*unsafeRuntimeType)(rt)).size
}

func rt2id(rt reflect.Type) uintptr {
	return uintptr(((*unsafeIntf)(unsafe.Pointer(&rt))).ptr)
}

func i2rtid(i interface{}) uintptr {
	return uintptr(((*unsafeIntf)(unsafe.Pointer(&i))).typ)
}



func unsafeCmpZero(ptr unsafe.Pointer, size int) bool {
	
	var s1 = unsafeString{ptr, size}
	var s2 = unsafeString{unsafeZeroAddr, size}
	if size > len(unsafeZeroArr) {
		arr := make([]byte, size)
		s2.Data = unsafe.Pointer(&arr[0])
	}
	return *(*string)(unsafe.Pointer(&s1)) == *(*string)(unsafe.Pointer(&s2)) 
}

func isEmptyValue(v reflect.Value, tinfos *TypeInfos, recursive bool) bool {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&v))
	if urv.flag == 0 {
		return true
	}
	if recursive {
		return isEmptyValueFallbackRecur(urv, v, tinfos)
	}
	return unsafeCmpZero(urv.ptr, int(rtsize2(urv.typ)))
}

func isEmptyValueFallbackRecur(urv *unsafeReflectValue, v reflect.Value, tinfos *TypeInfos) bool {
	const recursive = true

	switch v.Kind() {
	case reflect.Invalid:
		return true
	case reflect.String:
		return (*unsafeString)(urv.ptr).Len == 0
	case reflect.Slice:
		return (*unsafeSlice)(urv.ptr).Len == 0
	case reflect.Bool:
		return !*(*bool)(urv.ptr)
	case reflect.Int:
		return *(*int)(urv.ptr) == 0
	case reflect.Int8:
		return *(*int8)(urv.ptr) == 0
	case reflect.Int16:
		return *(*int16)(urv.ptr) == 0
	case reflect.Int32:
		return *(*int32)(urv.ptr) == 0
	case reflect.Int64:
		return *(*int64)(urv.ptr) == 0
	case reflect.Uint:
		return *(*uint)(urv.ptr) == 0
	case reflect.Uint8:
		return *(*uint8)(urv.ptr) == 0
	case reflect.Uint16:
		return *(*uint16)(urv.ptr) == 0
	case reflect.Uint32:
		return *(*uint32)(urv.ptr) == 0
	case reflect.Uint64:
		return *(*uint64)(urv.ptr) == 0
	case reflect.Uintptr:
		return *(*uintptr)(urv.ptr) == 0
	case reflect.Float32:
		return *(*float32)(urv.ptr) == 0
	case reflect.Float64:
		return *(*float64)(urv.ptr) == 0
	case reflect.Complex64:
		return unsafeCmpZero(urv.ptr, 8)
	case reflect.Complex128:
		return unsafeCmpZero(urv.ptr, 16)
	case reflect.Struct:
		
		if tinfos == nil {
			tinfos = defTypeInfos
		}
		ti := tinfos.find(uintptr(urv.typ))
		if ti == nil {
			ti = tinfos.load(v.Type())
		}
		return unsafeCmpZero(urv.ptr, int(ti.size))
	case reflect.Interface, reflect.Ptr:
		
		isnil := urv.ptr == nil || *(*unsafe.Pointer)(urv.ptr) == nil
		if recursive && !isnil {
			return isEmptyValue(v.Elem(), tinfos, recursive)
		}
		return isnil
	case reflect.UnsafePointer:
		return urv.ptr == nil || *(*unsafe.Pointer)(urv.ptr) == nil
	case reflect.Chan:
		return urv.ptr == nil || len_chan(rvRefPtr(urv)) == 0
	case reflect.Map:
		return urv.ptr == nil || len_map(rvRefPtr(urv)) == 0
	case reflect.Array:
		return v.Len() == 0 ||
			urv.ptr == nil ||
			urv.typ == nil ||
			rtsize2(urv.typ) == 0 ||
			unsafeCmpZero(urv.ptr, int(rtsize2(urv.typ)))
	}
	return false
}



type structFieldInfos struct {
	c      unsafe.Pointer 
	s      unsafe.Pointer 
	length int
}

func (x *structFieldInfos) load(source, sorted []*structFieldInfo) {
	s := (*unsafeSlice)(unsafe.Pointer(&sorted))
	x.s = s.Data
	x.length = s.Len
	s = (*unsafeSlice)(unsafe.Pointer(&source))
	x.c = s.Data
}

func (x *structFieldInfos) sorted() (v []*structFieldInfo) {
	*(*unsafeSlice)(unsafe.Pointer(&v)) = unsafeSlice{x.s, x.length, x.length}
	
	
	
	
	return
}

func (x *structFieldInfos) source() (v []*structFieldInfo) {
	*(*unsafeSlice)(unsafe.Pointer(&v)) = unsafeSlice{x.c, x.length, x.length}
	return
}









type atomicTypeInfoSlice struct {
	v unsafe.Pointer 
}

func (x *atomicTypeInfoSlice) load() (s []rtid2ti) {
	x2 := atomic.LoadPointer(&x.v)
	if x2 != nil {
		s = *(*[]rtid2ti)(x2)
	}
	return
}

func (x *atomicTypeInfoSlice) store(p []rtid2ti) {
	atomic.StorePointer(&x.v, unsafe.Pointer(&p))
}






type atomicRtidFnSlice struct {
	v unsafe.Pointer 
}

func (x *atomicRtidFnSlice) load() (s []codecRtidFn) {
	x2 := atomic.LoadPointer(&x.v)
	if x2 != nil {
		s = *(*[]codecRtidFn)(x2)
	}
	return
}

func (x *atomicRtidFnSlice) store(p []codecRtidFn) {
	atomic.StorePointer(&x.v, unsafe.Pointer(&p))
}


type atomicClsErr struct {
	v unsafe.Pointer 
}

func (x *atomicClsErr) load() (e clsErr) {
	x2 := (*clsErr)(atomic.LoadPointer(&x.v))
	if x2 != nil {
		e = *x2
	}
	return
}

func (x *atomicClsErr) store(p clsErr) {
	atomic.StorePointer(&x.v, unsafe.Pointer(&p))
}










type unsafeDecNakedWrapper struct {
	fauxUnion
	ru, ri, rf, rl, rs, rb, rt reflect.Value 
}

func (n *unsafeDecNakedWrapper) init() {
	n.ru = rv4iptr(&n.u).Elem()
	n.ri = rv4iptr(&n.i).Elem()
	n.rf = rv4iptr(&n.f).Elem()
	n.rl = rv4iptr(&n.l).Elem()
	n.rs = rv4iptr(&n.s).Elem()
	n.rt = rv4iptr(&n.t).Elem()
	n.rb = rv4iptr(&n.b).Elem()
	
}

var defUnsafeDecNakedWrapper unsafeDecNakedWrapper

func init() {
	defUnsafeDecNakedWrapper.init()
}

func (n *fauxUnion) ru() (v reflect.Value) {
	v = defUnsafeDecNakedWrapper.ru
	((*unsafeReflectValue)(unsafe.Pointer(&v))).ptr = unsafe.Pointer(&n.u)
	return
}
func (n *fauxUnion) ri() (v reflect.Value) {
	v = defUnsafeDecNakedWrapper.ri
	((*unsafeReflectValue)(unsafe.Pointer(&v))).ptr = unsafe.Pointer(&n.i)
	return
}
func (n *fauxUnion) rf() (v reflect.Value) {
	v = defUnsafeDecNakedWrapper.rf
	((*unsafeReflectValue)(unsafe.Pointer(&v))).ptr = unsafe.Pointer(&n.f)
	return
}
func (n *fauxUnion) rl() (v reflect.Value) {
	v = defUnsafeDecNakedWrapper.rl
	((*unsafeReflectValue)(unsafe.Pointer(&v))).ptr = unsafe.Pointer(&n.l)
	return
}
func (n *fauxUnion) rs() (v reflect.Value) {
	v = defUnsafeDecNakedWrapper.rs
	((*unsafeReflectValue)(unsafe.Pointer(&v))).ptr = unsafe.Pointer(&n.s)
	return
}
func (n *fauxUnion) rt() (v reflect.Value) {
	v = defUnsafeDecNakedWrapper.rt
	((*unsafeReflectValue)(unsafe.Pointer(&v))).ptr = unsafe.Pointer(&n.t)
	return
}
func (n *fauxUnion) rb() (v reflect.Value) {
	v = defUnsafeDecNakedWrapper.rb
	((*unsafeReflectValue)(unsafe.Pointer(&v))).ptr = unsafe.Pointer(&n.b)
	return
}


func rvSetBytes(rv reflect.Value, v []byte) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*[]byte)(urv.ptr) = v
}

func rvSetString(rv reflect.Value, v string) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*string)(urv.ptr) = v
}

func rvSetBool(rv reflect.Value, v bool) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*bool)(urv.ptr) = v
}

func rvSetTime(rv reflect.Value, v time.Time) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*time.Time)(urv.ptr) = v
}

func rvSetFloat32(rv reflect.Value, v float32) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*float32)(urv.ptr) = v
}

func rvSetFloat64(rv reflect.Value, v float64) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*float64)(urv.ptr) = v
}

func rvSetComplex64(rv reflect.Value, v complex64) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*complex64)(urv.ptr) = v
}

func rvSetComplex128(rv reflect.Value, v complex128) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*complex128)(urv.ptr) = v
}

func rvSetInt(rv reflect.Value, v int) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*int)(urv.ptr) = v
}

func rvSetInt8(rv reflect.Value, v int8) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*int8)(urv.ptr) = v
}

func rvSetInt16(rv reflect.Value, v int16) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*int16)(urv.ptr) = v
}

func rvSetInt32(rv reflect.Value, v int32) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*int32)(urv.ptr) = v
}

func rvSetInt64(rv reflect.Value, v int64) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*int64)(urv.ptr) = v
}

func rvSetUint(rv reflect.Value, v uint) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*uint)(urv.ptr) = v
}

func rvSetUintptr(rv reflect.Value, v uintptr) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*uintptr)(urv.ptr) = v
}

func rvSetUint8(rv reflect.Value, v uint8) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*uint8)(urv.ptr) = v
}

func rvSetUint16(rv reflect.Value, v uint16) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*uint16)(urv.ptr) = v
}

func rvSetUint32(rv reflect.Value, v uint32) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*uint32)(urv.ptr) = v
}

func rvSetUint64(rv reflect.Value, v uint64) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	*(*uint64)(urv.ptr) = v
}




func rvSetZero(rv reflect.Value) {
	rvSetDirectZero(rv)
}

func rvSetIntf(rv reflect.Value, v reflect.Value) {
	rv.Set(v)
}




func rvSetDirect(rv reflect.Value, v reflect.Value) {
	
	
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	uv := (*unsafeReflectValue)(unsafe.Pointer(&v))
	if uv.flag&unsafeFlagIndir == 0 {
		*(*unsafe.Pointer)(urv.ptr) = uv.ptr
	} else if uv.ptr == unsafeZeroAddr {
		if urv.ptr != unsafeZeroAddr {
			typedmemclr(urv.typ, urv.ptr)
		}
	} else {
		typedmemmove(urv.typ, urv.ptr, uv.ptr)
	}
}


func rvSetDirectZero(rv reflect.Value) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	if urv.ptr != unsafeZeroAddr {
		typedmemclr(urv.typ, urv.ptr)
	}
}




func rvMakeSlice(rv reflect.Value, ti *typeInfo, xlen, xcap int) (_ reflect.Value, set bool) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	ux := (*unsafeSlice)(urv.ptr)
	t := ((*unsafeIntf)(unsafe.Pointer(&ti.elem))).ptr
	s := unsafeSlice{newarray(t, xcap), xlen, xcap}
	if ux.Len > 0 {
		typedslicecopy(t, s, *ux)
	}
	*ux = s
	return rv, true
}




func rvSlice(rv reflect.Value, length int) reflect.Value {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	var x []struct{}
	ux := (*unsafeSlice)(unsafe.Pointer(&x))
	*ux = *(*unsafeSlice)(urv.ptr)
	ux.Len = length
	urv.ptr = unsafe.Pointer(ux)
	return rv
}




func rvGrowSlice(rv reflect.Value, ti *typeInfo, cap, incr int) (v reflect.Value, newcap int, set bool) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	ux := (*unsafeSlice)(urv.ptr)
	t := ((*unsafeIntf)(unsafe.Pointer(&ti.elem))).ptr
	*ux = unsafeGrowslice(t, *ux, cap, incr)
	ux.Len = ux.Cap
	return rv, ux.Cap, true
}



func rvSliceIndex(rv reflect.Value, i int, ti *typeInfo) (v reflect.Value) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	uv := (*unsafeReflectValue)(unsafe.Pointer(&v))
	uv.ptr = unsafe.Pointer(uintptr(((*unsafeSlice)(urv.ptr)).Data) + uintptr(int(ti.elemsize)*i))
	uv.typ = ((*unsafeIntf)(unsafe.Pointer(&ti.elem))).ptr
	uv.flag = uintptr(ti.elemkind) | unsafeFlagIndir | unsafeFlagAddr
	return
}

func rvSliceZeroCap(t reflect.Type) (v reflect.Value) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&v))
	urv.typ = ((*unsafeIntf)(unsafe.Pointer(&t))).ptr
	urv.flag = uintptr(reflect.Slice) | unsafeFlagIndir
	urv.ptr = unsafe.Pointer(&unsafeZeroSlice)
	return
}

func rvLenSlice(rv reflect.Value) int {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	return (*unsafeSlice)(urv.ptr).Len
}

func rvCapSlice(rv reflect.Value) int {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	return (*unsafeSlice)(urv.ptr).Cap
}

func rvArrayIndex(rv reflect.Value, i int, ti *typeInfo) (v reflect.Value) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	uv := (*unsafeReflectValue)(unsafe.Pointer(&v))
	uv.ptr = unsafe.Pointer(uintptr(urv.ptr) + uintptr(int(ti.elemsize)*i))
	uv.typ = ((*unsafeIntf)(unsafe.Pointer(&ti.elem))).ptr
	uv.flag = uintptr(ti.elemkind) | unsafeFlagIndir | unsafeFlagAddr
	return
}


func rvGetArrayBytes(rv reflect.Value, scratch []byte) (bs []byte) {
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	bx := (*unsafeSlice)(unsafe.Pointer(&bs))
	bx.Data = urv.ptr
	bx.Len = rv.Len()
	bx.Cap = bx.Len
	return
}

func rvGetArray4Slice(rv reflect.Value) (v reflect.Value) {
	
	
	
	
	
	
	

	t := reflectArrayOf(rvLenSlice(rv), rv.Type().Elem())
	

	uv := (*unsafeReflectValue)(unsafe.Pointer(&v))
	uv.flag = uintptr(reflect.Array) | unsafeFlagIndir | unsafeFlagAddr
	uv.typ = ((*unsafeIntf)(unsafe.Pointer(&t))).ptr

	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	uv.ptr = *(*unsafe.Pointer)(urv.ptr) 

	return
}

func rvGetSlice4Array(rv reflect.Value, v interface{}) {
	
	uv := (*unsafeIntf)(unsafe.Pointer(&v))
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))

	s := (*unsafeSlice)(uv.ptr)
	s.Data = urv.ptr
	s.Len = rv.Len()
	s.Cap = s.Len
}

func rvCopySlice(dest, src reflect.Value, elemType reflect.Type) {
	typedslicecopy((*unsafeIntf)(unsafe.Pointer(&elemType)).ptr,
		*(*unsafeSlice)((*unsafeReflectValue)(unsafe.Pointer(&dest)).ptr),
		*(*unsafeSlice)((*unsafeReflectValue)(unsafe.Pointer(&src)).ptr))
}



func rvGetBool(rv reflect.Value) bool {
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	return *(*bool)(v.ptr)
}

func rvGetBytes(rv reflect.Value) []byte {
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	return *(*[]byte)(v.ptr)
}

func rvGetTime(rv reflect.Value) time.Time {
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	return *(*time.Time)(v.ptr)
}

func rvGetString(rv reflect.Value) string {
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	return *(*string)(v.ptr)
}

func rvGetFloat64(rv reflect.Value) float64 {
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	return *(*float64)(v.ptr)
}

func rvGetFloat32(rv reflect.Value) float32 {
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	return *(*float32)(v.ptr)
}

func rvGetComplex64(rv reflect.Value) complex64 {
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	return *(*complex64)(v.ptr)
}

func rvGetComplex128(rv reflect.Value) complex128 {
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	return *(*complex128)(v.ptr)
}

func rvGetInt(rv reflect.Value) int {
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	return *(*int)(v.ptr)
}

func rvGetInt8(rv reflect.Value) int8 {
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	return *(*int8)(v.ptr)
}

func rvGetInt16(rv reflect.Value) int16 {
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	return *(*int16)(v.ptr)
}

func rvGetInt32(rv reflect.Value) int32 {
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	return *(*int32)(v.ptr)
}

func rvGetInt64(rv reflect.Value) int64 {
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	return *(*int64)(v.ptr)
}

func rvGetUint(rv reflect.Value) uint {
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	return *(*uint)(v.ptr)
}

func rvGetUint8(rv reflect.Value) uint8 {
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	return *(*uint8)(v.ptr)
}

func rvGetUint16(rv reflect.Value) uint16 {
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	return *(*uint16)(v.ptr)
}

func rvGetUint32(rv reflect.Value) uint32 {
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	return *(*uint32)(v.ptr)
}

func rvGetUint64(rv reflect.Value) uint64 {
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	return *(*uint64)(v.ptr)
}

func rvGetUintptr(rv reflect.Value) uintptr {
	v := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	return *(*uintptr)(v.ptr)
}

func rvLenMap(rv reflect.Value) int {
	
	
	
	

	

	return len_map(rvRefPtr((*unsafeReflectValue)(unsafe.Pointer(&rv))))
}






















































type unsafeMapIter struct {
	mtyp, mptr unsafe.Pointer
	k, v       reflect.Value
	kisref     bool
	visref     bool
	mapvalues  bool
	done       bool
	started    bool
	_          [3]byte 
	it         struct {
		key   unsafe.Pointer
		value unsafe.Pointer
		_     [20]uintptr 
	}
}

func (t *unsafeMapIter) Next() (r bool) {
	if t == nil || t.done {
		return
	}
	if t.started {
		mapiternext((unsafe.Pointer)(&t.it))
	} else {
		t.started = true
	}

	t.done = t.it.key == nil
	if t.done {
		return
	}

	if helperUnsafeDirectAssignMapEntry || t.kisref {
		(*unsafeReflectValue)(unsafe.Pointer(&t.k)).ptr = t.it.key
	} else {
		k := (*unsafeReflectValue)(unsafe.Pointer(&t.k))
		typedmemmove(k.typ, k.ptr, t.it.key)
	}

	if t.mapvalues {
		if helperUnsafeDirectAssignMapEntry || t.visref {
			(*unsafeReflectValue)(unsafe.Pointer(&t.v)).ptr = t.it.value
		} else {
			v := (*unsafeReflectValue)(unsafe.Pointer(&t.v))
			typedmemmove(v.typ, v.ptr, t.it.value)
		}
	}

	return true
}

func (t *unsafeMapIter) Key() (r reflect.Value) {
	return t.k
}

func (t *unsafeMapIter) Value() (r reflect.Value) {
	return t.v
}

func (t *unsafeMapIter) Done() {}

type mapIter struct {
	unsafeMapIter
}

func mapRange(t *mapIter, m, k, v reflect.Value, mapvalues bool) {
	if rvIsNil(m) {
		t.done = true
		return
	}
	t.done = false
	t.started = false
	t.mapvalues = mapvalues

	

	urv := (*unsafeReflectValue)(unsafe.Pointer(&m))
	t.mtyp = urv.typ
	t.mptr = rvRefPtr(urv)

	
	mapiterinit(t.mtyp, t.mptr, unsafe.Pointer(&t.it))

	t.k = k
	t.kisref = refBitset.isset(byte(k.Kind()))

	if mapvalues {
		t.v = v
		t.visref = refBitset.isset(byte(v.Kind()))
	} else {
		t.v = reflect.Value{}
	}
}



func unsafeMapKVPtr(urv *unsafeReflectValue) unsafe.Pointer {
	if urv.flag&unsafeFlagIndir == 0 {
		return unsafe.Pointer(&urv.ptr)
	}
	return urv.ptr
}











func mapAddrLoopvarRV(t reflect.Type, k reflect.Kind) (rv reflect.Value) {
	
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	urv.flag = uintptr(k) | unsafeFlagIndir | unsafeFlagAddr
	urv.typ = ((*unsafeIntf)(unsafe.Pointer(&t))).ptr
	
	
	if !helperUnsafeDirectAssignMapEntry {
		urv.ptr = unsafeNew(urv.typ)
	}
	return
}



func (e *Encoder) jsondriver() *jsonEncDriver {
	return (*jsonEncDriver)((*unsafeIntf)(unsafe.Pointer(&e.e)).ptr)
}

func (d *Decoder) zerocopystate() bool {
	return d.decByteState == decByteStateZerocopy && d.h.ZeroCopy
}

func (d *Decoder) stringZC(v []byte) (s string) {
	

	
	if d.decByteState == decByteStateZerocopy && d.h.ZeroCopy {
		return stringView(v)
	}
	return d.string(v)
}

func (d *Decoder) mapKeyString(callFnRvk *bool, kstrbs, kstr2bs *[]byte) string {
	if !d.zerocopystate() {
		*callFnRvk = true
		if d.decByteState == decByteStateReuseBuf {
			*kstrbs = append((*kstrbs)[:0], (*kstr2bs)...)
			*kstr2bs = *kstrbs
		}
	}
	return stringView(*kstr2bs)
}



func (d *Decoder) jsondriver() *jsonDecDriver {
	return (*jsonDecDriver)((*unsafeIntf)(unsafe.Pointer(&d.d)).ptr)
}



func (n *structFieldInfoPathNode) rvField(v reflect.Value) (rv reflect.Value) {
	
	uv := (*unsafeReflectValue)(unsafe.Pointer(&v))
	urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))
	
	urv.flag = uv.flag&(unsafeFlagStickyRO|unsafeFlagIndir|unsafeFlagAddr) | uintptr(n.kind)
	urv.typ = ((*unsafeIntf)(unsafe.Pointer(&n.typ))).ptr
	urv.ptr = unsafe.Pointer(uintptr(uv.ptr) + uintptr(n.offset))
	return
}






func len_map_chan(m unsafe.Pointer) int {
	if m == nil {
		return 0
	}
	return *((*int)(m))
}

func len_map(m unsafe.Pointer) int {
	
	return len_map_chan(m)
}
func len_chan(m unsafe.Pointer) int {
	
	return len_map_chan(m)
}

func unsafeNew(typ unsafe.Pointer) unsafe.Pointer {
	return mallocgc(rtsize2(typ), typ, true)
}




















//go:linkname memmove runtime.memmove
//go:noescape
func memmove(to, from unsafe.Pointer, n uintptr)

//go:linkname mallocgc runtime.mallocgc
//go:noescape
func mallocgc(size uintptr, typ unsafe.Pointer, needzero bool) unsafe.Pointer

//go:linkname newarray runtime.newarray
//go:noescape
func newarray(typ unsafe.Pointer, n int) unsafe.Pointer

//go:linkname mapiterinit runtime.mapiterinit
//go:noescape
func mapiterinit(typ unsafe.Pointer, m unsafe.Pointer, it unsafe.Pointer)

//go:linkname mapiternext runtime.mapiternext
//go:noescape
func mapiternext(it unsafe.Pointer) (key unsafe.Pointer)

//go:linkname mapdelete runtime.mapdelete
//go:noescape
func mapdelete(typ unsafe.Pointer, m unsafe.Pointer, key unsafe.Pointer)

//go:linkname mapassign runtime.mapassign
//go:noescape
func mapassign(typ unsafe.Pointer, m unsafe.Pointer, key unsafe.Pointer) unsafe.Pointer

//go:linkname mapaccess2 runtime.mapaccess2
//go:noescape
func mapaccess2(typ unsafe.Pointer, m unsafe.Pointer, key unsafe.Pointer) (val unsafe.Pointer, ok bool)







//go:linkname typedslicecopy reflect.typedslicecopy
//go:noescape
func typedslicecopy(elemType unsafe.Pointer, dst, src unsafeSlice) int

//go:linkname typedmemmove reflect.typedmemmove
//go:noescape
func typedmemmove(typ unsafe.Pointer, dst, src unsafe.Pointer)

//go:linkname typedmemclr reflect.typedmemclr
//go:noescape
func typedmemclr(typ unsafe.Pointer, dst unsafe.Pointer)
