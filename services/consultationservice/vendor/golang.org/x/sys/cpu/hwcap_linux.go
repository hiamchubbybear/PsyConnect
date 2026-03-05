



package cpu

import (
	"os"
)

const (
	_AT_HWCAP  = 16
	_AT_HWCAP2 = 26

	procAuxv = "/proc/self/auxv"

	uintSize = int(32 << (^uint(0) >> 63))
)




var hwCap uint
var hwCap2 uint

func readHWCAP() error {
	
	if a := getAuxv(); len(a) > 0 {
		for len(a) >= 2 {
			tag, val := a[0], uint(a[1])
			a = a[2:]
			switch tag {
			case _AT_HWCAP:
				hwCap = val
			case _AT_HWCAP2:
				hwCap2 = val
			}
		}
		return nil
	}

	buf, err := os.ReadFile(procAuxv)
	if err != nil {
		
		
		
		
		return err
	}
	bo := hostByteOrder()
	for len(buf) >= 2*(uintSize/8) {
		var tag, val uint
		switch uintSize {
		case 32:
			tag = uint(bo.Uint32(buf[0:]))
			val = uint(bo.Uint32(buf[4:]))
			buf = buf[8:]
		case 64:
			tag = uint(bo.Uint64(buf[0:]))
			val = uint(bo.Uint64(buf[8:]))
			buf = buf[16:]
		}
		switch tag {
		case _AT_HWCAP:
			hwCap = val
		case _AT_HWCAP2:
			hwCap2 = val
		}
	}
	return nil
}
