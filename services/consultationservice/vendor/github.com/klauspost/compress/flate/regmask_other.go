//go:build !amd64
// +build !amd64

package flate

const (
	
	
	

	
	reg8SizeMask8  = 0xff
	reg8SizeMask16 = 0xff
	reg8SizeMask32 = 0xff
	reg8SizeMask64 = 0xff

	
	reg16SizeMask8  = 0xffff
	reg16SizeMask16 = 0xffff
	reg16SizeMask32 = 0xffff
	reg16SizeMask64 = 0xffff

	
	reg32SizeMask8  = 0xffffffff
	reg32SizeMask16 = 0xffffffff
	reg32SizeMask32 = 0xffffffff
	reg32SizeMask64 = 0xffffffff

	
	reg64SizeMask8  = 0xffffffffffffffff
	reg64SizeMask16 = 0xffffffffffffffff
	reg64SizeMask32 = 0xffffffffffffffff
	reg64SizeMask64 = 0xffffffffffffffff

	
	regSizeMaskUint8  = ^uint(0)
	regSizeMaskUint16 = ^uint(0)
	regSizeMaskUint32 = ^uint(0)
	regSizeMaskUint64 = ^uint(0)
)
