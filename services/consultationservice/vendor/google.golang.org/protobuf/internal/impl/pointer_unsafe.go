



package impl

import (
	"reflect"
	"sync/atomic"
	"unsafe"

	"google.golang.org/protobuf/internal/protolazy"
)

const UnsafeEnabled = true


type Pointer unsafe.Pointer



type offset uintptr


func offsetOf(f reflect.StructField) offset {
	return offset(f.Offset)
}


func (f offset) IsValid() bool { return f != invalidOffset }


var invalidOffset = ^offset(0)


var zeroOffset = offset(0)


type pointer struct{ p unsafe.Pointer }


func pointerOf(p Pointer) pointer {
	return pointer{p: unsafe.Pointer(p)}
}


func pointerOfValue(v reflect.Value) pointer {
	return pointer{p: unsafe.Pointer(v.Pointer())}
}


func pointerOfIface(v any) pointer {
	type ifaceHeader struct {
		Type unsafe.Pointer
		Data unsafe.Pointer
	}
	return pointer{p: (*ifaceHeader)(unsafe.Pointer(&v)).Data}
}


func (p pointer) IsNil() bool {
	return p.p == nil
}



func (p pointer) Apply(f offset) pointer {
	if p.IsNil() {
		panic("invalid nil pointer")
	}
	return pointer{p: unsafe.Pointer(uintptr(p.p) + uintptr(f))}
}



func (p pointer) AsValueOf(t reflect.Type) reflect.Value {
	return reflect.NewAt(t, p.p)
}



func (p pointer) AsIfaceOf(t reflect.Type) any {
	
	return p.AsValueOf(t).Interface()
}

func (p pointer) Bool() *bool                           { return (*bool)(p.p) }
func (p pointer) BoolPtr() **bool                       { return (**bool)(p.p) }
func (p pointer) BoolSlice() *[]bool                    { return (*[]bool)(p.p) }
func (p pointer) Int32() *int32                         { return (*int32)(p.p) }
func (p pointer) Int32Ptr() **int32                     { return (**int32)(p.p) }
func (p pointer) Int32Slice() *[]int32                  { return (*[]int32)(p.p) }
func (p pointer) Int64() *int64                         { return (*int64)(p.p) }
func (p pointer) Int64Ptr() **int64                     { return (**int64)(p.p) }
func (p pointer) Int64Slice() *[]int64                  { return (*[]int64)(p.p) }
func (p pointer) Uint32() *uint32                       { return (*uint32)(p.p) }
func (p pointer) Uint32Ptr() **uint32                   { return (**uint32)(p.p) }
func (p pointer) Uint32Slice() *[]uint32                { return (*[]uint32)(p.p) }
func (p pointer) Uint64() *uint64                       { return (*uint64)(p.p) }
func (p pointer) Uint64Ptr() **uint64                   { return (**uint64)(p.p) }
func (p pointer) Uint64Slice() *[]uint64                { return (*[]uint64)(p.p) }
func (p pointer) Float32() *float32                     { return (*float32)(p.p) }
func (p pointer) Float32Ptr() **float32                 { return (**float32)(p.p) }
func (p pointer) Float32Slice() *[]float32              { return (*[]float32)(p.p) }
func (p pointer) Float64() *float64                     { return (*float64)(p.p) }
func (p pointer) Float64Ptr() **float64                 { return (**float64)(p.p) }
func (p pointer) Float64Slice() *[]float64              { return (*[]float64)(p.p) }
func (p pointer) String() *string                       { return (*string)(p.p) }
func (p pointer) StringPtr() **string                   { return (**string)(p.p) }
func (p pointer) StringSlice() *[]string                { return (*[]string)(p.p) }
func (p pointer) Bytes() *[]byte                        { return (*[]byte)(p.p) }
func (p pointer) BytesPtr() **[]byte                    { return (**[]byte)(p.p) }
func (p pointer) BytesSlice() *[][]byte                 { return (*[][]byte)(p.p) }
func (p pointer) Extensions() *map[int32]ExtensionField { return (*map[int32]ExtensionField)(p.p) }
func (p pointer) LazyInfoPtr() **protolazy.XXX_lazyUnmarshalInfo {
	return (**protolazy.XXX_lazyUnmarshalInfo)(p.p)
}

func (p pointer) PresenceInfo() presence {
	return presence{P: p.p}
}

func (p pointer) Elem() pointer {
	return pointer{p: *(*unsafe.Pointer)(p.p)}
}




func (p pointer) PointerSlice() []pointer {
	
	
	return *(*[]pointer)(p.p)
}


func (p pointer) AppendPointerSlice(v pointer) {
	*(*[]pointer)(p.p) = append(*(*[]pointer)(p.p), v)
}


func (p pointer) SetPointer(v pointer) {
	*(*unsafe.Pointer)(p.p) = (unsafe.Pointer)(v.p)
}

func (p pointer) growBoolSlice(addCap int) {
	sp := p.BoolSlice()
	s := make([]bool, 0, addCap+len(*sp))
	s = s[:len(*sp)]
	copy(s, *sp)
	*sp = s
}

func (p pointer) growInt32Slice(addCap int) {
	sp := p.Int32Slice()
	s := make([]int32, 0, addCap+len(*sp))
	s = s[:len(*sp)]
	copy(s, *sp)
	*sp = s
}

func (p pointer) growUint32Slice(addCap int) {
	p.growInt32Slice(addCap)
}

func (p pointer) growFloat32Slice(addCap int) {
	p.growInt32Slice(addCap)
}

func (p pointer) growInt64Slice(addCap int) {
	sp := p.Int64Slice()
	s := make([]int64, 0, addCap+len(*sp))
	s = s[:len(*sp)]
	copy(s, *sp)
	*sp = s
}

func (p pointer) growUint64Slice(addCap int) {
	p.growInt64Slice(addCap)
}

func (p pointer) growFloat64Slice(addCap int) {
	p.growInt64Slice(addCap)
}


const _ = uint(unsafe.Sizeof(unsafe.Pointer(nil)) - unsafe.Sizeof(MessageState{}))

func (Export) MessageStateOf(p Pointer) *messageState {
	
	return (*messageState)(unsafe.Pointer(p))
}
func (ms *messageState) pointer() pointer {
	
	return pointer{p: unsafe.Pointer(ms)}
}
func (ms *messageState) messageInfo() *MessageInfo {
	mi := ms.LoadMessageInfo()
	if mi == nil {
		panic("invalid nil message info; this suggests memory corruption due to a race or shallow copy on the message struct")
	}
	return mi
}
func (ms *messageState) LoadMessageInfo() *MessageInfo {
	return (*MessageInfo)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&ms.atomicMessageInfo))))
}
func (ms *messageState) StoreMessageInfo(mi *MessageInfo) {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&ms.atomicMessageInfo)), unsafe.Pointer(mi))
}

type atomicNilMessage struct{ p unsafe.Pointer } 

func (m *atomicNilMessage) Init(mi *MessageInfo) *messageReflectWrapper {
	if p := atomic.LoadPointer(&m.p); p != nil {
		return (*messageReflectWrapper)(p)
	}
	w := &messageReflectWrapper{mi: mi}
	atomic.CompareAndSwapPointer(&m.p, nil, (unsafe.Pointer)(w))
	return (*messageReflectWrapper)(atomic.LoadPointer(&m.p))
}
