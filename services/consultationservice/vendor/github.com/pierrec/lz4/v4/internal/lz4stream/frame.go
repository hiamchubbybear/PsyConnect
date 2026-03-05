
package lz4stream

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/pierrec/lz4/v4/internal/lz4block"
	"github.com/pierrec/lz4/v4/internal/lz4errors"
	"github.com/pierrec/lz4/v4/internal/xxh32"
)

//go:generate go run gen.go

const (
	frameMagic       uint32 = 0x184D2204
	frameSkipMagic   uint32 = 0x184D2A50
	frameMagicLegacy uint32 = 0x184C2102
)

func NewFrame() *Frame {
	return &Frame{}
}

type Frame struct {
	buf        [15]byte 
	Magic      uint32
	Descriptor FrameDescriptor
	Blocks     Blocks
	Checksum   uint32
	checksum   xxh32.XXHZero
}



func (f *Frame) Reset(num int) {
	f.Magic = 0
	f.Descriptor.Checksum = 0
	f.Descriptor.ContentSize = 0
	_ = f.Blocks.close(f, num)
	f.Checksum = 0
}

func (f *Frame) InitW(dst io.Writer, num int, legacy bool) {
	if legacy {
		f.Magic = frameMagicLegacy
		idx := lz4block.Index(lz4block.Block8Mb)
		f.Descriptor.Flags.BlockSizeIndexSet(idx)
	} else {
		f.Magic = frameMagic
		f.Descriptor.initW()
	}
	f.Blocks.initW(f, dst, num)
	f.checksum.Reset()
}

func (f *Frame) CloseW(dst io.Writer, num int) error {
	if err := f.Blocks.close(f, num); err != nil {
		return err
	}
	if f.isLegacy() {
		return nil
	}
	buf := f.buf[:0]
	
	buf = append(buf, 0, 0, 0, 0)
	if f.Descriptor.Flags.ContentChecksum() {
		buf = f.checksum.Sum(buf)
	}
	_, err := dst.Write(buf)
	return err
}

func (f *Frame) isLegacy() bool {
	return f.Magic == frameMagicLegacy
}

func (f *Frame) ParseHeaders(src io.Reader) error {
	if f.Magic > 0 {
		
		return nil
	}

newFrame:
	var err error
	if f.Magic, err = f.readUint32(src); err != nil {
		return err
	}
	switch m := f.Magic; {
	case m == frameMagic || m == frameMagicLegacy:
	
	case m>>8 == frameSkipMagic>>8:
		skip, err := f.readUint32(src)
		if err != nil {
			return err
		}
		if _, err := io.CopyN(ioutil.Discard, src, int64(skip)); err != nil {
			return err
		}
		goto newFrame
	default:
		return lz4errors.ErrInvalidFrame
	}
	if err := f.Descriptor.initR(f, src); err != nil {
		return err
	}
	f.checksum.Reset()
	return nil
}

func (f *Frame) InitR(src io.Reader, num int) (chan []byte, error) {
	return f.Blocks.initR(f, num, src)
}

func (f *Frame) CloseR(src io.Reader) (err error) {
	if f.isLegacy() {
		return nil
	}
	if !f.Descriptor.Flags.ContentChecksum() {
		return nil
	}
	if f.Checksum, err = f.readUint32(src); err != nil {
		return err
	}
	if c := f.checksum.Sum32(); c != f.Checksum {
		return fmt.Errorf("%w: got %x; expected %x", lz4errors.ErrInvalidFrameChecksum, c, f.Checksum)
	}
	return nil
}

type FrameDescriptor struct {
	Flags       DescriptorFlags
	ContentSize uint64
	Checksum    uint8
}

func (fd *FrameDescriptor) initW() {
	fd.Flags.VersionSet(1)
	fd.Flags.BlockIndependenceSet(true)
}

func (fd *FrameDescriptor) Write(f *Frame, dst io.Writer) error {
	if fd.Checksum > 0 {
		
		return nil
	}

	buf := f.buf[:4]
	
	binary.LittleEndian.PutUint32(buf, f.Magic)
	if !f.isLegacy() {
		buf = buf[:4+2]
		binary.LittleEndian.PutUint16(buf[4:], uint16(fd.Flags))

		if fd.Flags.Size() {
			buf = buf[:4+2+8]
			binary.LittleEndian.PutUint64(buf[4+2:], fd.ContentSize)
		}
		fd.Checksum = descriptorChecksum(buf[4:])
		buf = append(buf, fd.Checksum)
	}

	_, err := dst.Write(buf)
	return err
}

func (fd *FrameDescriptor) initR(f *Frame, src io.Reader) error {
	if f.isLegacy() {
		idx := lz4block.Index(lz4block.Block8Mb)
		f.Descriptor.Flags.BlockSizeIndexSet(idx)
		return nil
	}
	
	buf := f.buf[:3]
	if _, err := io.ReadFull(src, buf); err != nil {
		return err
	}
	descr := binary.LittleEndian.Uint16(buf)
	fd.Flags = DescriptorFlags(descr)
	if fd.Flags.Size() {
		
		buf = buf[:3+8]
		if _, err := io.ReadFull(src, buf[3:]); err != nil {
			return err
		}
		fd.ContentSize = binary.LittleEndian.Uint64(buf[2:])
	}
	fd.Checksum = buf[len(buf)-1] 
	buf = buf[:len(buf)-1]        
	if c := descriptorChecksum(buf); fd.Checksum != c {
		return fmt.Errorf("%w: got %x; expected %x", lz4errors.ErrInvalidHeaderChecksum, c, fd.Checksum)
	}
	
	if idx := fd.Flags.BlockSizeIndex(); !idx.IsValid() {
		return lz4errors.ErrOptionInvalidBlockSize
	}
	return nil
}

func descriptorChecksum(buf []byte) byte {
	return byte(xxh32.ChecksumZero(buf) >> 8)
}
