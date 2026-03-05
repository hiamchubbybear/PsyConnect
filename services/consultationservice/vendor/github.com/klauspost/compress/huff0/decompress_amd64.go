//go:build amd64 && !appengine && !noasm && gc
// +build amd64,!appengine,!noasm,gc



package huff0

import (
	"errors"
	"fmt"

	"github.com/klauspost/compress/internal/cpuinfo"
)




//go:noescape
func decompress4x_main_loop_amd64(ctx *decompress4xContext)





//go:noescape
func decompress4x_8b_main_loop_amd64(ctx *decompress4xContext)


const fallback8BitSize = 800

type decompress4xContext struct {
	pbr      *[4]bitReaderShifted
	peekBits uint8
	out      *byte
	dstEvery int
	tbl      *dEntrySingle
	decoded  int
	limit    *byte
}





func (d *Decoder) Decompress4X(dst, src []byte) ([]byte, error) {
	if len(d.dt.single) == 0 {
		return nil, errors.New("no table loaded")
	}
	if len(src) < 6+(4*1) {
		return nil, errors.New("input too small")
	}

	use8BitTables := d.actualTableLog <= 8
	if cap(dst) < fallback8BitSize && use8BitTables {
		return d.decompress4X8bit(dst, src)
	}

	var br [4]bitReaderShifted
	
	start := 6
	for i := 0; i < 3; i++ {
		length := int(src[i*2]) | (int(src[i*2+1]) << 8)
		if start+length >= len(src) {
			return nil, errors.New("truncated input (or invalid offset)")
		}
		err := br[i].init(src[start : start+length])
		if err != nil {
			return nil, err
		}
		start += length
	}
	err := br[3].init(src[start:])
	if err != nil {
		return nil, err
	}

	
	dstSize := cap(dst)
	dst = dst[:dstSize]
	out := dst
	dstEvery := (dstSize + 3) / 4

	const tlSize = 1 << tableLogMax
	const tlMask = tlSize - 1
	single := d.dt.single[:tlSize]

	var decoded int

	if len(out) > 4*4 && !(br[0].off < 4 || br[1].off < 4 || br[2].off < 4 || br[3].off < 4) {
		ctx := decompress4xContext{
			pbr:      &br,
			peekBits: uint8((64 - d.actualTableLog) & 63), 
			out:      &out[0],
			dstEvery: dstEvery,
			tbl:      &single[0],
			limit:    &out[dstEvery-4], 
		}
		if use8BitTables {
			decompress4x_8b_main_loop_amd64(&ctx)
		} else {
			decompress4x_main_loop_amd64(&ctx)
		}

		decoded = ctx.decoded
		out = out[decoded/4:]
	}

	
	remainBytes := dstEvery - (decoded / 4)
	for i := range br {
		offset := dstEvery * i
		endsAt := offset + remainBytes
		if endsAt > len(out) {
			endsAt = len(out)
		}
		br := &br[i]
		bitsLeft := br.remaining()
		for bitsLeft > 0 {
			br.fill()
			if offset >= endsAt {
				return nil, errors.New("corruption detected: stream overrun 4")
			}

			
			val := br.peekBitsFast(d.actualTableLog)
			v := single[val&tlMask].entry
			nBits := uint8(v)
			br.advance(nBits)
			bitsLeft -= uint(nBits)
			out[offset] = uint8(v >> 8)
			offset++
		}
		if offset != endsAt {
			return nil, fmt.Errorf("corruption detected: short output block %d, end %d != %d", i, offset, endsAt)
		}
		decoded += offset - dstEvery*i
		err = br.close()
		if err != nil {
			return nil, err
		}
	}
	if dstSize != decoded {
		return nil, errors.New("corruption detected: short output block")
	}
	return dst, nil
}




//go:noescape
func decompress1x_main_loop_amd64(ctx *decompress1xContext)




//go:noescape
func decompress1x_main_loop_bmi2(ctx *decompress1xContext)

type decompress1xContext struct {
	pbr      *bitReaderShifted
	peekBits uint8
	out      *byte
	outCap   int
	tbl      *dEntrySingle
	decoded  int
}


const error_max_decoded_size_exeeded = -1




func (d *Decoder) Decompress1X(dst, src []byte) ([]byte, error) {
	if len(d.dt.single) == 0 {
		return nil, errors.New("no table loaded")
	}
	var br bitReaderShifted
	err := br.init(src)
	if err != nil {
		return dst, err
	}
	maxDecodedSize := cap(dst)
	dst = dst[:maxDecodedSize]

	const tlSize = 1 << tableLogMax
	const tlMask = tlSize - 1

	if maxDecodedSize >= 4 {
		ctx := decompress1xContext{
			pbr:      &br,
			out:      &dst[0],
			outCap:   maxDecodedSize,
			peekBits: uint8((64 - d.actualTableLog) & 63), 
			tbl:      &d.dt.single[0],
		}

		if cpuinfo.HasBMI2() {
			decompress1x_main_loop_bmi2(&ctx)
		} else {
			decompress1x_main_loop_amd64(&ctx)
		}
		if ctx.decoded == error_max_decoded_size_exeeded {
			return nil, ErrMaxDecodedSizeExceeded
		}

		dst = dst[:ctx.decoded]
	}

	
	bitsLeft := uint8(br.off)*8 + 64 - br.bitsRead
	for bitsLeft > 0 {
		br.fill()
		if len(dst) >= maxDecodedSize {
			br.close()
			return nil, ErrMaxDecodedSizeExceeded
		}
		v := d.dt.single[br.peekBitsFast(d.actualTableLog)&tlMask]
		nBits := uint8(v.entry)
		br.advance(nBits)
		bitsLeft -= nBits
		dst = append(dst, uint8(v.entry>>8))
	}
	return dst, br.close()
}
