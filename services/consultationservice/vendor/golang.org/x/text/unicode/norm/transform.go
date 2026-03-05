



package norm

import (
	"unicode/utf8"

	"golang.org/x/text/transform"
)


func (Form) Reset() {}





func (f Form) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	
	b := src
	eof := atEOF
	if ns := len(dst); ns < len(b) {
		err = transform.ErrShortDst
		eof = false
		b = b[:ns]
	}
	i, ok := formTable[f].quickSpan(inputBytes(b), 0, len(b), eof)
	n := copy(dst, b[:i])
	if !ok {
		nDst, nSrc, err = f.transform(dst[n:], src[n:], atEOF)
		return nDst + n, nSrc + n, err
	}

	if err == nil && n < len(src) && !atEOF {
		err = transform.ErrShortSrc
	}
	return n, n, err
}

func flushTransform(rb *reorderBuffer) bool {
	
	if len(rb.out) < rb.nrune*utf8.UTFMax {
		return false
	}
	rb.out = rb.out[rb.flushCopy(rb.out):]
	return true
}

var errs = []error{nil, transform.ErrShortDst, transform.ErrShortSrc}



func (f Form) transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	
	rb := reorderBuffer{}
	rb.init(f, src)
	for {
		
		rb.setFlusher(dst[nDst:], flushTransform)
		end := decomposeSegment(&rb, nSrc, atEOF)
		if end < 0 {
			return nDst, nSrc, errs[-end]
		}
		nDst = len(dst) - len(rb.out)
		nSrc = end

		
		end = rb.nsrc
		eof := atEOF
		if n := nSrc + len(dst) - nDst; n < end {
			err = transform.ErrShortDst
			end = n
			eof = false
		}
		end, ok := rb.f.quickSpan(rb.src, nSrc, end, eof)
		n := copy(dst[nDst:], rb.src.bytes[nSrc:end])
		nSrc += n
		nDst += n
		if ok {
			if err == nil && n < rb.nsrc && !atEOF {
				err = transform.ErrShortSrc
			}
			return nDst, nSrc, err
		}
	}
}
