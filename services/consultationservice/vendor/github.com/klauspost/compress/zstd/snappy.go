



package zstd

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"

	"github.com/klauspost/compress/huff0"
	snappy "github.com/klauspost/compress/internal/snapref"
)

const (
	snappyTagLiteral = 0x00
	snappyTagCopy1   = 0x01
	snappyTagCopy2   = 0x02
	snappyTagCopy4   = 0x03
)

const (
	snappyChecksumSize = 4
	snappyMagicBody    = "sNaPpY"

	
	
	
	
	
	
	
	
	snappyMaxBlockSize = 65536

	
	
	
	
	snappyMaxEncodedLenOfMaxBlockSize = 76490
)

const (
	chunkTypeCompressedData   = 0x00
	chunkTypeUncompressedData = 0x01
	chunkTypePadding          = 0xfe
	chunkTypeStreamIdentifier = 0xff
)

var (
	
	ErrSnappyCorrupt = errors.New("snappy: corrupt input")
	
	ErrSnappyTooLarge = errors.New("snappy: decoded block is too large")
	
	ErrSnappyUnsupported = errors.New("snappy: unsupported input")

	errUnsupportedLiteralLength = errors.New("snappy: unsupported literal length")
)










type SnappyConverter struct {
	r     io.Reader
	err   error
	buf   []byte
	block *blockEnc
}




func (r *SnappyConverter) Convert(in io.Reader, w io.Writer) (int64, error) {
	initPredefined()
	r.err = nil
	r.r = in
	if r.block == nil {
		r.block = &blockEnc{}
		r.block.init()
	}
	r.block.initNewEncode()
	if len(r.buf) != snappyMaxEncodedLenOfMaxBlockSize+snappyChecksumSize {
		r.buf = make([]byte, snappyMaxEncodedLenOfMaxBlockSize+snappyChecksumSize)
	}
	r.block.litEnc.Reuse = huff0.ReusePolicyNone
	var written int64
	var readHeader bool
	{
		var header []byte
		var n int
		header, r.err = frameHeader{WindowSize: snappyMaxBlockSize}.appendTo(r.buf[:0])

		n, r.err = w.Write(header)
		if r.err != nil {
			return written, r.err
		}
		written += int64(n)
	}

	for {
		if !r.readFull(r.buf[:4], true) {
			
			r.block.reset(nil)
			r.block.last = true
			err := r.block.encodeLits(r.block.literals, false)
			if err != nil {
				return written, err
			}
			n, err := w.Write(r.block.output)
			if err != nil {
				return written, err
			}
			written += int64(n)

			return written, r.err
		}
		chunkType := r.buf[0]
		if !readHeader {
			if chunkType != chunkTypeStreamIdentifier {
				println("chunkType != chunkTypeStreamIdentifier", chunkType)
				r.err = ErrSnappyCorrupt
				return written, r.err
			}
			readHeader = true
		}
		chunkLen := int(r.buf[1]) | int(r.buf[2])<<8 | int(r.buf[3])<<16
		if chunkLen > len(r.buf) {
			println("chunkLen > len(r.buf)", chunkType)
			r.err = ErrSnappyUnsupported
			return written, r.err
		}

		
		
		switch chunkType {
		case chunkTypeCompressedData:
			
			if chunkLen < snappyChecksumSize {
				println("chunkLen < snappyChecksumSize", chunkLen, snappyChecksumSize)
				r.err = ErrSnappyCorrupt
				return written, r.err
			}
			buf := r.buf[:chunkLen]
			if !r.readFull(buf, false) {
				return written, r.err
			}
			
			buf = buf[snappyChecksumSize:]

			n, hdr, err := snappyDecodedLen(buf)
			if err != nil {
				r.err = err
				return written, r.err
			}
			buf = buf[hdr:]
			if n > snappyMaxBlockSize {
				println("n > snappyMaxBlockSize", n, snappyMaxBlockSize)
				r.err = ErrSnappyCorrupt
				return written, r.err
			}
			r.block.reset(nil)
			r.block.pushOffsets()
			if err := decodeSnappy(r.block, buf); err != nil {
				r.err = err
				return written, r.err
			}
			if r.block.size+r.block.extraLits != n {
				printf("invalid size, want %d, got %d\n", n, r.block.size+r.block.extraLits)
				r.err = ErrSnappyCorrupt
				return written, r.err
			}
			err = r.block.encode(nil, false, false)
			switch err {
			case errIncompressible:
				r.block.popOffsets()
				r.block.reset(nil)
				r.block.literals, err = snappy.Decode(r.block.literals[:n], r.buf[snappyChecksumSize:chunkLen])
				if err != nil {
					return written, err
				}
				err = r.block.encodeLits(r.block.literals, false)
				if err != nil {
					return written, err
				}
			case nil:
			default:
				return written, err
			}

			n, r.err = w.Write(r.block.output)
			if r.err != nil {
				return written, err
			}
			written += int64(n)
			continue
		case chunkTypeUncompressedData:
			if debugEncoder {
				println("Uncompressed, chunklen", chunkLen)
			}
			
			if chunkLen < snappyChecksumSize {
				println("chunkLen < snappyChecksumSize", chunkLen, snappyChecksumSize)
				r.err = ErrSnappyCorrupt
				return written, r.err
			}
			r.block.reset(nil)
			buf := r.buf[:snappyChecksumSize]
			if !r.readFull(buf, false) {
				return written, r.err
			}
			checksum := uint32(buf[0]) | uint32(buf[1])<<8 | uint32(buf[2])<<16 | uint32(buf[3])<<24
			
			n := chunkLen - snappyChecksumSize
			if n > snappyMaxBlockSize {
				println("n > snappyMaxBlockSize", n, snappyMaxBlockSize)
				r.err = ErrSnappyCorrupt
				return written, r.err
			}
			r.block.literals = r.block.literals[:n]
			if !r.readFull(r.block.literals, false) {
				return written, r.err
			}
			if snappyCRC(r.block.literals) != checksum {
				println("literals crc mismatch")
				r.err = ErrSnappyCorrupt
				return written, r.err
			}
			err := r.block.encodeLits(r.block.literals, false)
			if err != nil {
				return written, err
			}
			n, r.err = w.Write(r.block.output)
			if r.err != nil {
				return written, err
			}
			written += int64(n)
			continue

		case chunkTypeStreamIdentifier:
			if debugEncoder {
				println("stream id", chunkLen, len(snappyMagicBody))
			}
			
			if chunkLen != len(snappyMagicBody) {
				println("chunkLen != len(snappyMagicBody)", chunkLen, len(snappyMagicBody))
				r.err = ErrSnappyCorrupt
				return written, r.err
			}
			if !r.readFull(r.buf[:len(snappyMagicBody)], false) {
				return written, r.err
			}
			for i := 0; i < len(snappyMagicBody); i++ {
				if r.buf[i] != snappyMagicBody[i] {
					println("r.buf[i] != snappyMagicBody[i]", r.buf[i], snappyMagicBody[i], i)
					r.err = ErrSnappyCorrupt
					return written, r.err
				}
			}
			continue
		}

		if chunkType <= 0x7f {
			
			println("chunkType <= 0x7f")
			r.err = ErrSnappyUnsupported
			return written, r.err
		}
		
		
		if !r.readFull(r.buf[:chunkLen], false) {
			return written, r.err
		}
	}
}



func decodeSnappy(blk *blockEnc, src []byte) error {
	
	var s, length int
	lits := blk.extraLits
	var offset uint32
	for s < len(src) {
		switch src[s] & 0x03 {
		case snappyTagLiteral:
			x := uint32(src[s] >> 2)
			switch {
			case x < 60:
				s++
			case x == 60:
				s += 2
				if uint(s) > uint(len(src)) { 
					println("uint(s) > uint(len(src)", s, src)
					return ErrSnappyCorrupt
				}
				x = uint32(src[s-1])
			case x == 61:
				s += 3
				if uint(s) > uint(len(src)) { 
					println("uint(s) > uint(len(src)", s, src)
					return ErrSnappyCorrupt
				}
				x = uint32(src[s-2]) | uint32(src[s-1])<<8
			case x == 62:
				s += 4
				if uint(s) > uint(len(src)) { 
					println("uint(s) > uint(len(src)", s, src)
					return ErrSnappyCorrupt
				}
				x = uint32(src[s-3]) | uint32(src[s-2])<<8 | uint32(src[s-1])<<16
			case x == 63:
				s += 5
				if uint(s) > uint(len(src)) { 
					println("uint(s) > uint(len(src)", s, src)
					return ErrSnappyCorrupt
				}
				x = uint32(src[s-4]) | uint32(src[s-3])<<8 | uint32(src[s-2])<<16 | uint32(src[s-1])<<24
			}
			if x > snappyMaxBlockSize {
				println("x > snappyMaxBlockSize", x, snappyMaxBlockSize)
				return ErrSnappyCorrupt
			}
			length = int(x) + 1
			if length <= 0 {
				println("length <= 0 ", length)

				return errUnsupportedLiteralLength
			}
			
			
			

			blk.literals = append(blk.literals, src[s:s+length]...)
			
			lits += length
			s += length
			continue

		case snappyTagCopy1:
			s += 2
			if uint(s) > uint(len(src)) { 
				println("uint(s) > uint(len(src)", s, len(src))
				return ErrSnappyCorrupt
			}
			length = 4 + int(src[s-2])>>2&0x7
			offset = uint32(src[s-2])&0xe0<<3 | uint32(src[s-1])

		case snappyTagCopy2:
			s += 3
			if uint(s) > uint(len(src)) { 
				println("uint(s) > uint(len(src)", s, len(src))
				return ErrSnappyCorrupt
			}
			length = 1 + int(src[s-3])>>2
			offset = uint32(src[s-2]) | uint32(src[s-1])<<8

		case snappyTagCopy4:
			s += 5
			if uint(s) > uint(len(src)) { 
				println("uint(s) > uint(len(src)", s, len(src))
				return ErrSnappyCorrupt
			}
			length = 1 + int(src[s-5])>>2
			offset = uint32(src[s-4]) | uint32(src[s-3])<<8 | uint32(src[s-2])<<16 | uint32(src[s-1])<<24
		}

		if offset <= 0 || blk.size+lits < int(offset)  {
			println("offset <= 0 || blk.size+lits < int(offset)", offset, blk.size+lits, int(offset), blk.size, lits)

			return ErrSnappyCorrupt
		}

		
		
		
		if false {
			offset = blk.matchOffset(offset, uint32(lits))
		} else {
			offset += 3
		}

		blk.sequences = append(blk.sequences, seq{
			litLen:   uint32(lits),
			offset:   offset,
			matchLen: uint32(length) - zstdMinMatch,
		})
		blk.size += length + lits
		lits = 0
	}
	blk.extraLits = lits
	return nil
}

func (r *SnappyConverter) readFull(p []byte, allowEOF bool) (ok bool) {
	if _, r.err = io.ReadFull(r.r, p); r.err != nil {
		if r.err == io.ErrUnexpectedEOF || (r.err == io.EOF && !allowEOF) {
			r.err = ErrSnappyCorrupt
		}
		return false
	}
	return true
}

var crcTable = crc32.MakeTable(crc32.Castagnoli)



func snappyCRC(b []byte) uint32 {
	c := crc32.Update(0, crcTable, b)
	return c>>15 | c<<17 + 0xa282ead8
}



func snappyDecodedLen(src []byte) (blockLen, headerLen int, err error) {
	v, n := binary.Uvarint(src)
	if n <= 0 || v > 0xffffffff {
		return 0, 0, ErrSnappyCorrupt
	}

	const wordSize = 32 << (^uint(0) >> 32 & 1)
	if wordSize == 32 && v > 0x7fffffff {
		return 0, 0, ErrSnappyTooLarge
	}
	return int(v), n, nil
}
