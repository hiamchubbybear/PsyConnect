// +build go1.16,!go1.20



package rt

import (
    _ `unsafe`
)

//go:linkname GrowSlice runtime.growslice

func GrowSlice(et *GoType, old GoSlice, cap int) GoSlice
