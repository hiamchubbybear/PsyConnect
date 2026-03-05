

package lz4stream

import "github.com/pierrec/lz4/v4/internal/lz4block"













type DescriptorFlags uint16


func (x DescriptorFlags) ContentChecksum() bool   { return x>>2&1 != 0 }
func (x DescriptorFlags) Size() bool              { return x>>3&1 != 0 }
func (x DescriptorFlags) BlockChecksum() bool     { return x>>4&1 != 0 }
func (x DescriptorFlags) BlockIndependence() bool { return x>>5&1 != 0 }
func (x DescriptorFlags) Version() uint16         { return uint16(x >> 6 & 0x3) }
func (x DescriptorFlags) BlockSizeIndex() lz4block.BlockSizeIndex {
	return lz4block.BlockSizeIndex(x >> 12 & 0x7)
}


func (x *DescriptorFlags) ContentChecksumSet(v bool) *DescriptorFlags {
	const b = 1 << 2
	if v {
		*x = *x&^b | b
	} else {
		*x &^= b
	}
	return x
}
func (x *DescriptorFlags) SizeSet(v bool) *DescriptorFlags {
	const b = 1 << 3
	if v {
		*x = *x&^b | b
	} else {
		*x &^= b
	}
	return x
}
func (x *DescriptorFlags) BlockChecksumSet(v bool) *DescriptorFlags {
	const b = 1 << 4
	if v {
		*x = *x&^b | b
	} else {
		*x &^= b
	}
	return x
}
func (x *DescriptorFlags) BlockIndependenceSet(v bool) *DescriptorFlags {
	const b = 1 << 5
	if v {
		*x = *x&^b | b
	} else {
		*x &^= b
	}
	return x
}
func (x *DescriptorFlags) VersionSet(v uint16) *DescriptorFlags {
	*x = *x&^(0x3<<6) | (DescriptorFlags(v) & 0x3 << 6)
	return x
}
func (x *DescriptorFlags) BlockSizeIndexSet(v lz4block.BlockSizeIndex) *DescriptorFlags {
	*x = *x&^(0x7<<12) | (DescriptorFlags(v) & 0x7 << 12)
	return x
}








type DataBlockSize uint32


func (x DataBlockSize) size() int          { return int(x & 0x7FFFFFFF) }
func (x DataBlockSize) Uncompressed() bool { return x>>31&1 != 0 }


func (x *DataBlockSize) sizeSet(v int) *DataBlockSize {
	*x = *x&^0x7FFFFFFF | DataBlockSize(v)&0x7FFFFFFF
	return x
}
func (x *DataBlockSize) UncompressedSet(v bool) *DataBlockSize {
	const b = 1 << 31
	if v {
		*x = *x&^b | b
	} else {
		*x &^= b
	}
	return x
}
