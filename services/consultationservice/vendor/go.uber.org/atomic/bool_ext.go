



















package atomic

import (
	"strconv"
)

//go:generate bin/gen-atomicwrapper -name=Bool -type=bool -wrapped=Uint32 -pack=boolToInt -unpack=truthy -cas -swap -json -file=bool.go

func truthy(n uint32) bool {
	return n == 1
}

func boolToInt(b bool) uint32 {
	if b {
		return 1
	}
	return 0
}


func (b *Bool) Toggle() (old bool) {
	for {
		old := b.Load()
		if b.CAS(old, !old) {
			return old
		}
	}
}


func (b *Bool) String() string {
	return strconv.FormatBool(b.Load())
}
