



package snappy

import (
	"io"

	"github.com/klauspost/compress/s2"
)

var (
	
	ErrCorrupt = s2.ErrCorrupt
	
	ErrTooLarge = s2.ErrTooLarge
	
	ErrUnsupported = s2.ErrUnsupported
)

const (
	
	
	
	
	
	
	
	
	maxBlockSize = 65536
)


func DecodedLen(src []byte) (int, error) {
	return s2.DecodedLen(src)
}








func Decode(dst, src []byte) ([]byte, error) {
	return s2.Decode(dst, src)
}




func NewReader(r io.Reader) *Reader {
	return s2.NewReader(r, s2.ReaderMaxBlockSize(maxBlockSize))
}




type Reader = s2.Reader
