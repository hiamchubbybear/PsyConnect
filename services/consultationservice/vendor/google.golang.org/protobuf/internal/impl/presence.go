



package impl

import (
	"sync/atomic"
	"unsafe"
)


type presenceSize uint32


type presence struct {
	
	P unsafe.Pointer
}

func (p presence) toElem(num uint32) (ret *uint32) {
	const (
		bitsPerByte = 8
		siz         = unsafe.Sizeof(*ret)
	)
	
	
	
	offset := uintptr(num) / (siz * bitsPerByte) * siz
	return (*uint32)(unsafe.Pointer(uintptr(p.P) + offset))
}


func (p presence) Present(num uint32) bool {
	if p.P == nil {
		return false
	}
	return Export{}.Present(p.toElem(num), num)
}


func (p presence) SetPresent(num uint32, size presenceSize) {
	Export{}.SetPresent(p.toElem(num), num, uint32(size))
}



func (p presence) SetPresentUnatomic(num uint32, size presenceSize) {
	Export{}.SetPresentNonAtomic(p.toElem(num), num, uint32(size))
}


func (p presence) ClearPresent(num uint32) {
	Export{}.ClearPresent(p.toElem(num), num)
}






func (p presence) LoadPresenceCache() (current uint32) {
	if p.P == nil {
		return 0
	}
	return atomic.LoadUint32((*uint32)(p.P))
}






func (p presence) PresentInCache(num uint32, cachedElement *uint32, current *uint32) bool {
	if num/32 != *cachedElement {
		o := uintptr(num/32) * unsafe.Sizeof(uint32(0))
		q := (*uint32)(unsafe.Pointer(uintptr(p.P) + o))
		*current = atomic.LoadUint32(q)
		*cachedElement = num / 32
	}
	return (*current & (1 << (num % 32))) > 0
}


func (p presence) AnyPresent(size presenceSize) bool {
	n := uintptr((size + 31) / 32)
	for j := uintptr(0); j < n; j++ {
		o := j * unsafe.Sizeof(uint32(0))
		q := (*uint32)(unsafe.Pointer(uintptr(p.P) + o))
		b := atomic.LoadUint32(q)
		if b > 0 {
			return true
		}
	}
	return false
}











func (p presence) toRaceDetectData() *RaceDetectHookData {
	var template struct {
		d RaceDetectHookData
		a [1]uint32
	}
	o := (uintptr(unsafe.Pointer(&template.a)) - uintptr(unsafe.Pointer(&template.d)))
	return (*RaceDetectHookData)(unsafe.Pointer(uintptr(p.P) - o))
}

func atomicLoadShadowPresence(p **[]byte) *[]byte {
	return (*[]byte)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(p))))
}
func atomicStoreShadowPresence(p **[]byte, v *[]byte) {
	atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(p)), nil, unsafe.Pointer(v))
}












func findPointerToRaceDetectData(ptr *uint32, num uint32) *RaceDetectHookData {
	var template struct {
		d RaceDetectHookData
		a [1]uint32
	}
	o := (uintptr(unsafe.Pointer(&template.a)) - uintptr(unsafe.Pointer(&template.d))) + uintptr(num/32)*unsafe.Sizeof(uint32(0))
	return (*RaceDetectHookData)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) - o))
}
