


package zstd

import (
	"encoding/binary"
	"errors"
	"io"
)



const HeaderMaxSize = 14 + 3


type Header struct {
	
	
	
	SingleSegment bool

	
	
	WindowSize uint64

	
	
	DictionaryID uint32

	
	HasFCS bool

	
	FrameContentSize uint64

	
	
	Skippable bool

	
	
	SkippableID int

	
	
	SkippableSize uint32

	
	
	
	
	
	
	
	
	
	
	
	HeaderSize int

	
	FirstBlock struct {
		
		OK bool

		
		Last bool

		
		
		
		
		Compressed bool

		
		
		DecompressedSize int

		
		
		
		CompressedSize int
	}

	
	
	HasCheckSum bool
}







func (h *Header) Decode(in []byte) error {
	*h = Header{}
	if len(in) < 4 {
		return io.ErrUnexpectedEOF
	}
	h.HeaderSize += 4
	b, in := in[:4], in[4:]
	if string(b) != frameMagic {
		if string(b[1:4]) != skippableFrameMagic || b[0]&0xf0 != 0x50 {
			return ErrMagicMismatch
		}
		if len(in) < 4 {
			return io.ErrUnexpectedEOF
		}
		h.HeaderSize += 4
		h.Skippable = true
		h.SkippableID = int(b[0] & 0xf)
		h.SkippableSize = binary.LittleEndian.Uint32(in)
		return nil
	}

	
	
	if len(in) < 1 {
		return io.ErrUnexpectedEOF
	}
	fhd, in := in[0], in[1:]
	h.HeaderSize++
	h.SingleSegment = fhd&(1<<5) != 0
	h.HasCheckSum = fhd&(1<<2) != 0
	if fhd&(1<<3) != 0 {
		return errors.New("reserved bit set on frame header")
	}

	if !h.SingleSegment {
		if len(in) < 1 {
			return io.ErrUnexpectedEOF
		}
		var wd byte
		wd, in = in[0], in[1:]
		h.HeaderSize++
		windowLog := 10 + (wd >> 3)
		windowBase := uint64(1) << windowLog
		windowAdd := (windowBase / 8) * uint64(wd&0x7)
		h.WindowSize = windowBase + windowAdd
	}

	
	
	if size := fhd & 3; size != 0 {
		if size == 3 {
			size = 4
		}
		if len(in) < int(size) {
			return io.ErrUnexpectedEOF
		}
		b, in = in[:size], in[size:]
		h.HeaderSize += int(size)
		switch len(b) {
		case 1:
			h.DictionaryID = uint32(b[0])
		case 2:
			h.DictionaryID = uint32(b[0]) | (uint32(b[1]) << 8)
		case 4:
			h.DictionaryID = uint32(b[0]) | (uint32(b[1]) << 8) | (uint32(b[2]) << 16) | (uint32(b[3]) << 24)
		}
	}

	
	
	var fcsSize int
	v := fhd >> 6
	switch v {
	case 0:
		if h.SingleSegment {
			fcsSize = 1
		}
	default:
		fcsSize = 1 << v
	}

	if fcsSize > 0 {
		h.HasFCS = true
		if len(in) < fcsSize {
			return io.ErrUnexpectedEOF
		}
		b, in = in[:fcsSize], in[fcsSize:]
		h.HeaderSize += int(fcsSize)
		switch len(b) {
		case 1:
			h.FrameContentSize = uint64(b[0])
		case 2:
			
			h.FrameContentSize = uint64(b[0]) | (uint64(b[1]) << 8) + 256
		case 4:
			h.FrameContentSize = uint64(b[0]) | (uint64(b[1]) << 8) | (uint64(b[2]) << 16) | (uint64(b[3]) << 24)
		case 8:
			d1 := uint32(b[0]) | (uint32(b[1]) << 8) | (uint32(b[2]) << 16) | (uint32(b[3]) << 24)
			d2 := uint32(b[4]) | (uint32(b[5]) << 8) | (uint32(b[6]) << 16) | (uint32(b[7]) << 24)
			h.FrameContentSize = uint64(d1) | (uint64(d2) << 32)
		}
	}

	
	if len(in) < 3 {
		return nil
	}
	tmp := in[:3]
	bh := uint32(tmp[0]) | (uint32(tmp[1]) << 8) | (uint32(tmp[2]) << 16)
	h.FirstBlock.Last = bh&1 != 0
	blockType := blockType((bh >> 1) & 3)
	
	cSize := int(bh >> 3)
	switch blockType {
	case blockTypeReserved:
		return nil
	case blockTypeRLE:
		h.FirstBlock.Compressed = true
		h.FirstBlock.DecompressedSize = cSize
		h.FirstBlock.CompressedSize = 1
	case blockTypeCompressed:
		h.FirstBlock.Compressed = true
		h.FirstBlock.CompressedSize = cSize
	case blockTypeRaw:
		h.FirstBlock.DecompressedSize = cSize
		h.FirstBlock.CompressedSize = cSize
	default:
		panic("Invalid block type")
	}

	h.FirstBlock.OK = true
	return nil
}
