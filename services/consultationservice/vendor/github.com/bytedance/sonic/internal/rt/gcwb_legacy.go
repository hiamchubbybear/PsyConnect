// +build go1.17,!go1.21



package rt

import (
    _ `unsafe`
)

//go:linkname GcWriteBarrierAX runtime.gcWriteBarrier
func GcWriteBarrierAX()

//go:linkname RuntimeWriteBarrier runtime.writeBarrier
var RuntimeWriteBarrier uintptr
