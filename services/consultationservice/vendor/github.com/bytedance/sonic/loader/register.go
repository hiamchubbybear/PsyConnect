// +build !bytedance_tango



package loader

import (
	"sync/atomic"
	"unsafe"
)

func registerModule(mod *moduledata) {
    registerModuleLockFree(&lastmoduledatap, mod)
}

func registerModuleLockFree(tail **moduledata, mod *moduledata) {
    for {
        oldTail := loadModule(tail)
        if casModule(tail, oldTail, mod) {
            storeModule(&oldTail.next, mod)
            break
        }
    }
}

func loadModule(p **moduledata) *moduledata {
    return (*moduledata)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(p))))
}

func storeModule(p **moduledata, value *moduledata) {
    atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(p)), unsafe.Pointer(value))
}

func casModule(p **moduledata, oldValue *moduledata, newValue *moduledata) bool {
    return atomic.CompareAndSwapPointer(
        (*unsafe.Pointer)(unsafe.Pointer(p)),
        unsafe.Pointer(oldValue),
        unsafe.Pointer(newValue),
    )
}

