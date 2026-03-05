package lz4stream

import (
	"encoding/binary"
	"fmt"
	"io"
	"sync"

	"github.com/pierrec/lz4/v4/internal/lz4block"
	"github.com/pierrec/lz4/v4/internal/lz4errors"
	"github.com/pierrec/lz4/v4/internal/xxh32"
)

type Blocks struct {
	Block  *FrameDataBlock
	Blocks chan chan *FrameDataBlock
	mu     sync.Mutex
	err    error
}

func (b *Blocks) initW(f *Frame, dst io.Writer, num int) {
	if num == 1 {
		b.Blocks = nil
		b.Block = NewFrameDataBlock(f)
		return
	}
	b.Block = nil
	if cap(b.Blocks) != num {
		b.Blocks = make(chan chan *FrameDataBlock, num)
	}
	
	go func() {
		
		for c := range b.Blocks {
			
			
			
			
			block := <-c
			if block == nil {
				
				
				close(c)
				return
			}
			
			if b.err == nil {
				
				if err := block.Write(f, dst); err != nil {
					
					b.err = err
					
				}
			}
			close(c)
		}
	}()
}

func (b *Blocks) close(f *Frame, num int) error {
	if num == 1 {
		if b.Block != nil {
			b.Block.Close(f)
		}
		err := b.err
		b.err = nil
		return err
	}
	if b.Blocks == nil {
		err := b.err
		b.err = nil
		return err
	}
	c := make(chan *FrameDataBlock)
	b.Blocks <- c
	c <- nil
	<-c
	err := b.err
	b.err = nil
	return err
}


func (b *Blocks) ErrorR() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.err
}






func (b *Blocks) initR(f *Frame, num int, src io.Reader) (chan []byte, error) {
	size := f.Descriptor.Flags.BlockSizeIndex()
	if num == 1 {
		b.Blocks = nil
		b.Block = NewFrameDataBlock(f)
		return nil, nil
	}
	b.Block = nil
	blocks := make(chan chan []byte, num)
	
	data := make(chan []byte)
	
	

	
	var cum uint32
	go func() {
		var cumx uint32
		var err error
		for b.ErrorR() == nil {
			block := NewFrameDataBlock(f)
			cumx, err = block.Read(f, src, 0)
			if err != nil {
				block.Close(f)
				break
			}
			
			if b.ErrorR() != nil {
				block.Close(f)
				break
			}
			c := make(chan []byte)
			blocks <- c
			go func() {
				defer block.Close(f)
				data, err := block.Uncompress(f, size.Get(), nil, false)
				if err != nil {
					b.closeR(err)
					
					close(c)
				} else {
					c <- data
				}
			}()
		}
		
		c := make(chan []byte)
		blocks <- c
		c <- nil 
		<-c      
		if f.isLegacy() && cum == cumx {
			err = io.EOF
		}
		b.closeR(err)
		close(data)
	}()
	
	
	go func(leg bool) {
		defer close(blocks)
		skipBlocks := false
		for c := range blocks {
			buf, ok := <-c
			if !ok {
				
				
				skipBlocks = true
				continue
			}
			if buf == nil {
				
				close(c)
				return
			}
			if skipBlocks {
				
				continue
			}
			
			if f.Descriptor.Flags.ContentChecksum() {
				_, _ = f.checksum.Write(buf)
			}
			if leg {
				cum += uint32(len(buf))
			}
			data <- buf
			close(c)
		}
	}(f.isLegacy())
	return data, nil
}


func (b *Blocks) closeR(err error) {
	b.mu.Lock()
	if b.err == nil {
		b.err = err
	}
	b.mu.Unlock()
}

func NewFrameDataBlock(f *Frame) *FrameDataBlock {
	buf := f.Descriptor.Flags.BlockSizeIndex().Get()
	return &FrameDataBlock{Data: buf, data: buf}
}

type FrameDataBlock struct {
	Size     DataBlockSize
	Data     []byte 
	Checksum uint32
	data     []byte 
	src      []byte 
	err      error  
}

func (b *FrameDataBlock) Close(f *Frame) {
	b.Size = 0
	b.Checksum = 0
	b.err = nil
	if b.data != nil {
		
		lz4block.Put(b.data)
		b.Data = nil
		b.data = nil
		b.src = nil
	}
}


func (b *FrameDataBlock) Compress(f *Frame, src []byte, level lz4block.CompressionLevel) *FrameDataBlock {
	data := b.data
	if f.isLegacy() {
		
		
		src = src[:8<<20]
	} else {
		data = data[:len(src)] 
	}
	var n int
	switch level {
	case lz4block.Fast:
		n, _ = lz4block.CompressBlock(src, data)
	default:
		n, _ = lz4block.CompressBlockHC(src, data, level)
	}
	if n == 0 {
		b.Size.UncompressedSet(true)
		b.Data = src
	} else {
		b.Size.UncompressedSet(false)
		b.Data = data[:n]
	}
	b.Size.sizeSet(len(b.Data))
	b.src = src 

	if f.Descriptor.Flags.BlockChecksum() {
		b.Checksum = xxh32.ChecksumZero(src)
	}
	return b
}

func (b *FrameDataBlock) Write(f *Frame, dst io.Writer) error {
	
	
	if f.Descriptor.Flags.ContentChecksum() {
		_, _ = f.checksum.Write(b.src)
	}
	buf := f.buf[:]
	binary.LittleEndian.PutUint32(buf, uint32(b.Size))
	if _, err := dst.Write(buf[:4]); err != nil {
		return err
	}

	if _, err := dst.Write(b.Data); err != nil {
		return err
	}

	if b.Checksum == 0 {
		return nil
	}
	binary.LittleEndian.PutUint32(buf, b.Checksum)
	_, err := dst.Write(buf[:4])
	return err
}


func (b *FrameDataBlock) Read(f *Frame, src io.Reader, cum uint32) (uint32, error) {
	x, err := f.readUint32(src)
	if err != nil {
		return 0, err
	}
	if f.isLegacy() {
		switch x {
		case frameMagicLegacy:
			
			return b.Read(f, src, cum)
		case cum:
			
			
			
			return 0, io.EOF
		}
	} else if x == 0 {
		
		return 0, io.EOF
	}
	b.Size = DataBlockSize(x)

	size := b.Size.size()
	if size > cap(b.data) {
		return x, lz4errors.ErrOptionInvalidBlockSize
	}
	b.data = b.data[:size]
	if _, err := io.ReadFull(src, b.data); err != nil {
		return x, err
	}
	if f.Descriptor.Flags.BlockChecksum() {
		sum, err := f.readUint32(src)
		if err != nil {
			return 0, err
		}
		b.Checksum = sum
	}
	return x, nil
}

func (b *FrameDataBlock) Uncompress(f *Frame, dst, dict []byte, sum bool) ([]byte, error) {
	if b.Size.Uncompressed() {
		n := copy(dst, b.data)
		dst = dst[:n]
	} else {
		n, err := lz4block.UncompressBlock(b.data, dst, dict)
		if err != nil {
			return nil, err
		}
		dst = dst[:n]
	}
	if f.Descriptor.Flags.BlockChecksum() {
		if c := xxh32.ChecksumZero(dst); c != b.Checksum {
			err := fmt.Errorf("%w: got %x; expected %x", lz4errors.ErrInvalidBlockChecksum, c, b.Checksum)
			return nil, err
		}
	}
	if sum && f.Descriptor.Flags.ContentChecksum() {
		_, _ = f.checksum.Write(dst)
	}
	return dst, nil
}

func (f *Frame) readUint32(r io.Reader) (x uint32, err error) {
	if _, err = io.ReadFull(r, f.buf[:4]); err != nil {
		return
	}
	x = binary.LittleEndian.Uint32(f.buf[:4])
	return
}
