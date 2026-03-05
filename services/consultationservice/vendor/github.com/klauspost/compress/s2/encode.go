




package s2

import (
	"encoding/binary"
	"math"
	"math/bits"
)













func Encode(dst, src []byte) []byte {
	if n := MaxEncodedLen(len(src)); n < 0 {
		panic(ErrTooLarge)
	} else if cap(dst) < n {
		dst = make([]byte, n)
	} else {
		dst = dst[:n]
	}

	
	d := binary.PutUvarint(dst, uint64(len(src)))

	if len(src) == 0 {
		return dst[:d]
	}
	if len(src) < minNonLiteralBlockSize {
		d += emitLiteral(dst[d:], src)
		return dst[:d]
	}
	n := encodeBlock(dst[d:], src)
	if n > 0 {
		d += n
		return dst[:d]
	}
	
	d += emitLiteral(dst[d:], src)
	return dst[:d]
}





func EstimateBlockSize(src []byte) (d int) {
	if len(src) < 6 || int64(len(src)) > 0xffffffff {
		return -1
	}
	if len(src) <= 1024 {
		d = calcBlockSizeSmall(src)
	} else {
		d = calcBlockSize(src)
	}

	if d == 0 {
		return -1
	}
	
	d += (bits.Len64(uint64(len(src))) + 7) / 7

	if d >= len(src) {
		return -1
	}
	return d
}
















func EncodeBetter(dst, src []byte) []byte {
	if n := MaxEncodedLen(len(src)); n < 0 {
		panic(ErrTooLarge)
	} else if len(dst) < n {
		dst = make([]byte, n)
	}

	
	d := binary.PutUvarint(dst, uint64(len(src)))

	if len(src) == 0 {
		return dst[:d]
	}
	if len(src) < minNonLiteralBlockSize {
		d += emitLiteral(dst[d:], src)
		return dst[:d]
	}
	n := encodeBlockBetter(dst[d:], src)
	if n > 0 {
		d += n
		return dst[:d]
	}
	
	d += emitLiteral(dst[d:], src)
	return dst[:d]
}
















func EncodeBest(dst, src []byte) []byte {
	if n := MaxEncodedLen(len(src)); n < 0 {
		panic(ErrTooLarge)
	} else if len(dst) < n {
		dst = make([]byte, n)
	}

	
	d := binary.PutUvarint(dst, uint64(len(src)))

	if len(src) == 0 {
		return dst[:d]
	}
	if len(src) < minNonLiteralBlockSize {
		d += emitLiteral(dst[d:], src)
		return dst[:d]
	}
	n := encodeBlockBest(dst[d:], src, nil)
	if n > 0 {
		d += n
		return dst[:d]
	}
	
	d += emitLiteral(dst[d:], src)
	return dst[:d]
}















func EncodeSnappy(dst, src []byte) []byte {
	if n := MaxEncodedLen(len(src)); n < 0 {
		panic(ErrTooLarge)
	} else if cap(dst) < n {
		dst = make([]byte, n)
	} else {
		dst = dst[:n]
	}

	
	d := binary.PutUvarint(dst, uint64(len(src)))

	if len(src) == 0 {
		return dst[:d]
	}
	if len(src) < minNonLiteralBlockSize {
		d += emitLiteral(dst[d:], src)
		return dst[:d]
	}

	n := encodeBlockSnappy(dst[d:], src)
	if n > 0 {
		d += n
		return dst[:d]
	}
	
	d += emitLiteral(dst[d:], src)
	return dst[:d]
}















func EncodeSnappyBetter(dst, src []byte) []byte {
	if n := MaxEncodedLen(len(src)); n < 0 {
		panic(ErrTooLarge)
	} else if cap(dst) < n {
		dst = make([]byte, n)
	} else {
		dst = dst[:n]
	}

	
	d := binary.PutUvarint(dst, uint64(len(src)))

	if len(src) == 0 {
		return dst[:d]
	}
	if len(src) < minNonLiteralBlockSize {
		d += emitLiteral(dst[d:], src)
		return dst[:d]
	}

	n := encodeBlockBetterSnappy(dst[d:], src)
	if n > 0 {
		d += n
		return dst[:d]
	}
	
	d += emitLiteral(dst[d:], src)
	return dst[:d]
}















func EncodeSnappyBest(dst, src []byte) []byte {
	if n := MaxEncodedLen(len(src)); n < 0 {
		panic(ErrTooLarge)
	} else if cap(dst) < n {
		dst = make([]byte, n)
	} else {
		dst = dst[:n]
	}

	
	d := binary.PutUvarint(dst, uint64(len(src)))

	if len(src) == 0 {
		return dst[:d]
	}
	if len(src) < minNonLiteralBlockSize {
		d += emitLiteral(dst[d:], src)
		return dst[:d]
	}

	n := encodeBlockBestSnappy(dst[d:], src)
	if n > 0 {
		d += n
		return dst[:d]
	}
	
	d += emitLiteral(dst[d:], src)
	return dst[:d]
}






func ConcatBlocks(dst []byte, blocks ...[]byte) ([]byte, error) {
	totalSize := uint64(0)
	compSize := 0
	for _, b := range blocks {
		l, hdr, err := decodedLen(b)
		if err != nil {
			return nil, err
		}
		totalSize += uint64(l)
		compSize += len(b) - hdr
	}
	if totalSize == 0 {
		dst = append(dst, 0)
		return dst, nil
	}
	if totalSize > math.MaxUint32 {
		return nil, ErrTooLarge
	}
	var tmp [binary.MaxVarintLen32]byte
	hdrSize := binary.PutUvarint(tmp[:], totalSize)
	wantSize := hdrSize + compSize

	if cap(dst)-len(dst) < wantSize {
		dst = append(make([]byte, 0, wantSize+len(dst)), dst...)
	}
	dst = append(dst, tmp[:hdrSize]...)
	for _, b := range blocks {
		_, hdr, err := decodedLen(b)
		if err != nil {
			return nil, err
		}
		dst = append(dst, b[hdr:]...)
	}
	return dst, nil
}









const inputMargin = 8



const minNonLiteralBlockSize = 32

const intReduction = 2 - (1 << (^uint(0) >> 63)) 




const MaxBlockSize = (1<<(32-intReduction) - 1) - binary.MaxVarintLen32 - 5






func MaxEncodedLen(srcLen int) int {
	n := uint64(srcLen)
	if intReduction == 1 {
		
		if n > math.MaxInt32 {
			
			return -1
		}
	} else if n > 0xffffffff {
		
		
		return -1
	}
	
	n = n + uint64((bits.Len64(n)+7)/7)

	
	n += uint64(literalExtraSize(int64(srcLen)))
	if intReduction == 1 {
		
		if n > math.MaxInt32 {
			return -1
		}
	} else if n > 0xffffffff {
		
		
		return -1
	}
	return int(n)
}
