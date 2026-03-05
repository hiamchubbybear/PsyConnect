



package cpu

import (
	"runtime"
)


type byteOrder interface {
	Uint32([]byte) uint32
	Uint64([]byte) uint64
}

type littleEndian struct{}
type bigEndian struct{}

func (littleEndian) Uint32(b []byte) uint32 {
	_ = b[3] 
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
}

func (littleEndian) Uint64(b []byte) uint64 {
	_ = b[7] 
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48 | uint64(b[7])<<56
}

func (bigEndian) Uint32(b []byte) uint32 {
	_ = b[3] 
	return uint32(b[3]) | uint32(b[2])<<8 | uint32(b[1])<<16 | uint32(b[0])<<24
}

func (bigEndian) Uint64(b []byte) uint64 {
	_ = b[7] 
	return uint64(b[7]) | uint64(b[6])<<8 | uint64(b[5])<<16 | uint64(b[4])<<24 |
		uint64(b[3])<<32 | uint64(b[2])<<40 | uint64(b[1])<<48 | uint64(b[0])<<56
}



func hostByteOrder() byteOrder {
	switch runtime.GOARCH {
	case "386", "amd64", "amd64p32",
		"alpha",
		"arm", "arm64",
		"loong64",
		"mipsle", "mips64le", "mips64p32le",
		"nios2",
		"ppc64le",
		"riscv", "riscv64",
		"sh":
		return littleEndian{}
	case "armbe", "arm64be",
		"m68k",
		"mips", "mips64", "mips64p32",
		"ppc", "ppc64",
		"s390", "s390x",
		"shbe",
		"sparc", "sparc64":
		return bigEndian{}
	}
	panic("unknown architecture")
}
