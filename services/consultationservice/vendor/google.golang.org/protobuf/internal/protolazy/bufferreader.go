





package protolazy

import (
	"errors"
	"fmt"
	"io"

	"google.golang.org/protobuf/encoding/protowire"
)


type BufferReader struct {
	Buf []byte
	Pos int
}


func NewBufferReader(buf []byte) BufferReader {
	return BufferReader{Buf: buf, Pos: 0}
}

var errOutOfBounds = errors.New("protobuf decoding: out of bounds")
var errOverflow = errors.New("proto: integer overflow")

func (b *BufferReader) DecodeVarintSlow() (x uint64, err error) {
	i := b.Pos
	l := len(b.Buf)

	for shift := uint(0); shift < 64; shift += 7 {
		if i >= l {
			err = io.ErrUnexpectedEOF
			return
		}
		v := b.Buf[i]
		i++
		x |= (uint64(v) & 0x7F) << shift
		if v < 0x80 {
			b.Pos = i
			return
		}
	}

	
	err = errOverflow
	return
}


func (b *BufferReader) DecodeVarint() (x uint64, err error) {
	i := b.Pos
	buf := b.Buf

	if i >= len(buf) {
		return 0, io.ErrUnexpectedEOF
	} else if buf[i] < 0x80 {
		b.Pos++
		return uint64(buf[i]), nil
	} else if len(buf)-i < 10 {
		return b.DecodeVarintSlow()
	}

	var v uint64
	
	x = uint64(buf[i]) & 127
	i++

	v = uint64(buf[i])
	i++
	x |= (v & 127) << 7
	if v < 128 {
		goto done
	}

	v = uint64(buf[i])
	i++
	x |= (v & 127) << 14
	if v < 128 {
		goto done
	}

	v = uint64(buf[i])
	i++
	x |= (v & 127) << 21
	if v < 128 {
		goto done
	}

	v = uint64(buf[i])
	i++
	x |= (v & 127) << 28
	if v < 128 {
		goto done
	}

	v = uint64(buf[i])
	i++
	x |= (v & 127) << 35
	if v < 128 {
		goto done
	}

	v = uint64(buf[i])
	i++
	x |= (v & 127) << 42
	if v < 128 {
		goto done
	}

	v = uint64(buf[i])
	i++
	x |= (v & 127) << 49
	if v < 128 {
		goto done
	}

	v = uint64(buf[i])
	i++
	x |= (v & 127) << 56
	if v < 128 {
		goto done
	}

	v = uint64(buf[i])
	i++
	x |= (v & 127) << 63
	if v < 128 {
		goto done
	}

	return 0, errOverflow

done:
	b.Pos = i
	return
}


func (b *BufferReader) DecodeVarint32() (x uint32, err error) {
	i := b.Pos
	buf := b.Buf

	if i >= len(buf) {
		return 0, io.ErrUnexpectedEOF
	} else if buf[i] < 0x80 {
		b.Pos++
		return uint32(buf[i]), nil
	} else if len(buf)-i < 5 {
		v, err := b.DecodeVarintSlow()
		return uint32(v), err
	}

	var v uint32
	
	x = uint32(buf[i]) & 127
	i++

	v = uint32(buf[i])
	i++
	x |= (v & 127) << 7
	if v < 128 {
		goto done
	}

	v = uint32(buf[i])
	i++
	x |= (v & 127) << 14
	if v < 128 {
		goto done
	}

	v = uint32(buf[i])
	i++
	x |= (v & 127) << 21
	if v < 128 {
		goto done
	}

	v = uint32(buf[i])
	i++
	x |= (v & 127) << 28
	if v < 128 {
		goto done
	}

	return 0, errOverflow

done:
	b.Pos = i
	return
}


func (b *BufferReader) SkipValue(tag uint32) (err error) {
	wireType := tag & 0x7
	switch protowire.Type(wireType) {
	case protowire.VarintType:
		err = b.SkipVarint()
	case protowire.Fixed64Type:
		err = b.SkipFixed64()
	case protowire.BytesType:
		var n uint32
		n, err = b.DecodeVarint32()
		if err == nil {
			err = b.Skip(int(n))
		}
	case protowire.StartGroupType:
		err = b.SkipGroup(tag)
	case protowire.Fixed32Type:
		err = b.SkipFixed32()
	default:
		err = fmt.Errorf("Unexpected wire type (%d)", wireType)
	}
	return
}


func (b *BufferReader) SkipGroup(tag uint32) (err error) {
	tagStack := make([]uint32, 0, 16)
	tagStack = append(tagStack, tag)
	var n uint32
	for len(tagStack) > 0 {
		tag, err = b.DecodeVarint32()
		if err != nil {
			return err
		}
		switch protowire.Type(tag & 0x7) {
		case protowire.VarintType:
			err = b.SkipVarint()
		case protowire.Fixed64Type:
			err = b.Skip(8)
		case protowire.BytesType:
			n, err = b.DecodeVarint32()
			if err == nil {
				err = b.Skip(int(n))
			}
		case protowire.StartGroupType:
			tagStack = append(tagStack, tag)
		case protowire.Fixed32Type:
			err = b.SkipFixed32()
		case protowire.EndGroupType:
			if protoFieldNumber(tagStack[len(tagStack)-1]) == protoFieldNumber(tag) {
				tagStack = tagStack[:len(tagStack)-1]
			} else {
				err = fmt.Errorf("end group tag %d does not match begin group tag %d at pos %d",
					protoFieldNumber(tag), protoFieldNumber(tagStack[len(tagStack)-1]), b.Pos)
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}


func (b *BufferReader) SkipVarint() (err error) {
	i := b.Pos

	if len(b.Buf)-i < 10 {
		
		if _, err := b.DecodeVarintSlow(); err != nil {
			return err
		}
		return nil
	}

	if b.Buf[i] < 0x80 {
		goto out
	}
	i++

	if b.Buf[i] < 0x80 {
		goto out
	}
	i++

	if b.Buf[i] < 0x80 {
		goto out
	}
	i++

	if b.Buf[i] < 0x80 {
		goto out
	}
	i++

	if b.Buf[i] < 0x80 {
		goto out
	}
	i++

	if b.Buf[i] < 0x80 {
		goto out
	}
	i++

	if b.Buf[i] < 0x80 {
		goto out
	}
	i++

	if b.Buf[i] < 0x80 {
		goto out
	}
	i++

	if b.Buf[i] < 0x80 {
		goto out
	}
	i++

	if b.Buf[i] < 0x80 {
		goto out
	}
	return errOverflow

out:
	b.Pos = i + 1
	return nil
}


func (b *BufferReader) Skip(n int) (err error) {
	if len(b.Buf) < b.Pos+n {
		return io.ErrUnexpectedEOF
	}
	b.Pos += n
	return
}


func (b *BufferReader) SkipFixed64() (err error) {
	return b.Skip(8)
}


func (b *BufferReader) SkipFixed32() (err error) {
	return b.Skip(4)
}


func (b *BufferReader) SkipBytes() (err error) {
	n, err := b.DecodeVarint32()
	if err != nil {
		return err
	}
	return b.Skip(int(n))
}


func (b *BufferReader) Done() bool {
	return b.Pos == len(b.Buf)
}


func (b *BufferReader) Remaining() int {
	return len(b.Buf) - b.Pos
}
