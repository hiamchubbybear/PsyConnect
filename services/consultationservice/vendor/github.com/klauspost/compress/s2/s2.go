

































package s2

import (
	"bytes"
	"hash/crc32"
)


const (
	tagLiteral = 0x00
	tagCopy1   = 0x01
	tagCopy2   = 0x02
	tagCopy4   = 0x03
)

const (
	checksumSize     = 4
	chunkHeaderSize  = 4
	magicChunk       = "\xff\x06\x00\x00" + magicBody
	magicChunkSnappy = "\xff\x06\x00\x00" + magicBodySnappy
	magicBodySnappy  = "sNaPpY"
	magicBody        = "S2sTwO"

	
	
	
	
	maxBlockSize = 4 << 20

	
	minBlockSize = 4 << 10

	skippableFrameHeader = 4
	maxChunkSize         = 1<<24 - 1 

	
	defaultBlockSize = 1 << 20

	
	maxSnappyBlockSize = 1 << 16

	obufHeaderLen = checksumSize + chunkHeaderSize
)

const (
	chunkTypeCompressedData   = 0x00
	chunkTypeUncompressedData = 0x01
	ChunkTypeIndex            = 0x99
	chunkTypePadding          = 0xfe
	chunkTypeStreamIdentifier = 0xff
)

var crcTable = crc32.MakeTable(crc32.Castagnoli)



func crc(b []byte) uint32 {
	c := crc32.Update(0, crcTable, b)
	return c>>15 | c<<17 + 0xa282ead8
}



func literalExtraSize(n int64) int64 {
	if n == 0 {
		return 0
	}
	switch {
	case n < 60:
		return 1
	case n < 1<<8:
		return 2
	case n < 1<<16:
		return 3
	case n < 1<<24:
		return 4
	default:
		return 5
	}
}

type byter interface {
	Bytes() []byte
}

var _ byter = &bytes.Buffer{}
