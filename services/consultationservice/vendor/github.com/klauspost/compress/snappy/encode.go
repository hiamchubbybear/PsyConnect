



package snappy

import (
	"io"

	"github.com/klauspost/compress/s2"
)








func Encode(dst, src []byte) []byte {
	return s2.EncodeSnappyBetter(dst, src)
}





func MaxEncodedLen(srcLen int) int {
	return s2.MaxEncodedLen(srcLen)
}










func NewWriter(w io.Writer) *Writer {
	return s2.NewWriter(w, s2.WriterSnappyCompat(), s2.WriterBetterCompression(), s2.WriterFlushOnWrite(), s2.WriterConcurrency(1))
}








func NewBufferedWriter(w io.Writer) *Writer {
	return s2.NewWriter(w, s2.WriterSnappyCompat(), s2.WriterBetterCompression())
}




type Writer = s2.Writer
