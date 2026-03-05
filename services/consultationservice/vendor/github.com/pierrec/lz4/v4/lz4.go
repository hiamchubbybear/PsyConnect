







package lz4

import (
	"github.com/pierrec/lz4/v4/internal/lz4block"
	"github.com/pierrec/lz4/v4/internal/lz4errors"
)

func _() {
	
	var x [1]struct{}
	_ = x[lz4block.CompressionLevel(Fast)-lz4block.Fast]
	_ = x[Block64Kb-BlockSize(lz4block.Block64Kb)]
	_ = x[Block256Kb-BlockSize(lz4block.Block256Kb)]
	_ = x[Block1Mb-BlockSize(lz4block.Block1Mb)]
	_ = x[Block4Mb-BlockSize(lz4block.Block4Mb)]
}


func CompressBlockBound(n int) int {
	return lz4block.CompressBlockBound(n)
}







func UncompressBlock(src, dst []byte) (int, error) {
	return lz4block.UncompressBlock(src, dst, nil)
}







func UncompressBlockWithDict(src, dst, dict []byte) (int, error) {
	return lz4block.UncompressBlock(src, dst, dict)
}







type Compressor struct{ c lz4block.Compressor }












func (c *Compressor) CompressBlock(src, dst []byte) (int, error) {
	return c.c.CompressBlock(src, dst)
}



















func CompressBlock(src, dst []byte, _ []int) (int, error) {
	return lz4block.CompressBlock(src, dst)
}








type CompressorHC struct {
	
	
	Level CompressionLevel
	c     lz4block.CompressorHC
}












func (c *CompressorHC) CompressBlock(src, dst []byte) (int, error) {
	return c.c.CompressBlock(src, dst, lz4block.CompressionLevel(c.Level))
}





func CompressBlockHC(src, dst []byte, depth CompressionLevel, _, _ []int) (int, error) {
	return lz4block.CompressBlockHC(src, dst, lz4block.CompressionLevel(depth))
}

const (
	
	
	ErrInvalidSourceShortBuffer = lz4errors.ErrInvalidSourceShortBuffer
	
	ErrInvalidFrame = lz4errors.ErrInvalidFrame
	
	ErrInternalUnhandledState = lz4errors.ErrInternalUnhandledState
	
	ErrInvalidHeaderChecksum = lz4errors.ErrInvalidHeaderChecksum
	
	ErrInvalidBlockChecksum = lz4errors.ErrInvalidBlockChecksum
	
	ErrInvalidFrameChecksum = lz4errors.ErrInvalidFrameChecksum
	
	ErrOptionInvalidCompressionLevel = lz4errors.ErrOptionInvalidCompressionLevel
	
	ErrOptionClosedOrError = lz4errors.ErrOptionClosedOrError
	
	ErrOptionInvalidBlockSize = lz4errors.ErrOptionInvalidBlockSize
	
	ErrOptionNotApplicable = lz4errors.ErrOptionNotApplicable
	
	ErrWriterNotClosed = lz4errors.ErrWriterNotClosed
)
