package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"time"

	"github.com/segmentio/kafka-go/compress"
)


type Attributes int16

const (
	Gzip          Attributes = Attributes(compress.Gzip)   
	Snappy        Attributes = Attributes(compress.Snappy) 
	Lz4           Attributes = Attributes(compress.Lz4)    
	Zstd          Attributes = Attributes(compress.Zstd)   
	Transactional Attributes = 1 << 4
	Control       Attributes = 1 << 5
)

func (a Attributes) Compression() compress.Compression {
	return compress.Compression(a & 7)
}

func (a Attributes) Transactional() bool {
	return (a & Transactional) != 0
}

func (a Attributes) Control() bool {
	return (a & Control) != 0
}

func (a Attributes) String() string {
	s := a.Compression().String()
	if a.Transactional() {
		s += "+transactional"
	}
	if a.Control() {
		s += "+control"
	}
	return s
}


type Header struct {
	Key   string
	Value []byte
}




type Record struct {
	
	
	Offset int64

	
	
	Time time.Time

	
	
	
	
	Key Bytes

	
	
	
	
	Value Bytes

	
	
	
	Headers []Header
}



type RecordSet struct {
	
	
	
	
	
	
	
	
	Version int8

	
	
	
	
	
	
	
	Attributes Attributes

	
	
	
	
	
	
	Records RecordReader
}




type bufferedReader interface {
	Discard(int) (int, error)
	Peek(int) ([]byte, error)
}




type bytesBuffer interface {
	Bytes() []byte
}



const magicByteOffset = 16




func (rs *RecordSet) ReadFrom(r io.Reader) (int64, error) {
	d, _ := r.(*decoder)
	if d == nil {
		d = &decoder{
			reader: r,
			remain: 4,
		}
	}

	*rs = RecordSet{}
	limit := d.remain
	size := d.readInt32()

	if d.err != nil {
		return int64(limit - d.remain), d.err
	}

	if size <= 0 {
		return 4, nil
	}

	stream := &RecordStream{
		Records: make([]RecordReader, 0, 4),
	}

	var err error
	d.remain = int(size)

	for d.remain > 0 && err == nil {
		var version byte

		if d.remain < (magicByteOffset + 1) {
			if len(stream.Records) != 0 {
				break
			}
			return 4, fmt.Errorf("impossible record set shorter than %d bytes", magicByteOffset+1)
		}

		switch r := d.reader.(type) {
		case bufferedReader:
			b, err := r.Peek(magicByteOffset + 1)
			if err != nil {
				n, _ := r.Discard(len(b))
				return 4 + int64(n), dontExpectEOF(err)
			}
			version = b[magicByteOffset]
		case bytesBuffer:
			version = r.Bytes()[magicByteOffset]
		default:
			b := make([]byte, magicByteOffset+1)
			if n, err := io.ReadFull(d.reader, b); err != nil {
				return 4 + int64(n), dontExpectEOF(err)
			}
			version = b[magicByteOffset]
			
			
			
			
			
			
			
			d.reader = io.MultiReader(bytes.NewReader(b), d.reader)
		}

		var tmp RecordSet
		switch version {
		case 0, 1:
			err = tmp.readFromVersion1(d)
		case 2:
			err = tmp.readFromVersion2(d)
		default:
			err = fmt.Errorf("unsupported message version %d for message of size %d", version, size)
		}

		if tmp.Version > rs.Version {
			rs.Version = tmp.Version
		}

		rs.Attributes |= tmp.Attributes

		if tmp.Records != nil {
			stream.Records = append(stream.Records, tmp.Records)
		}
	}

	if len(stream.Records) != 0 {
		rs.Records = stream
		
		
		err = nil
	}

	d.discardAll()
	rn := 4 + (int(size) - d.remain)
	d.remain = limit - rn
	return int64(rn), err
}









func (rs *RecordSet) WriteTo(w io.Writer) (int64, error) {
	if rs.Records == nil {
		return 0, ErrNoRecord
	}

	
	
	
	buffer, _ := w.(*pageBuffer)
	bufferOffset := int64(0)

	if buffer != nil {
		bufferOffset = buffer.Size()
	} else {
		buffer = newPageBuffer()
		defer buffer.unref()
	}

	size := packUint32(0)
	buffer.Write(size[:]) 

	var err error
	switch rs.Version {
	case 0, 1:
		err = rs.writeToVersion1(buffer, bufferOffset+4)
	case 2:
		err = rs.writeToVersion2(buffer, bufferOffset+4)
	default:
		err = fmt.Errorf("unsupported record set version %d", rs.Version)
	}
	if err != nil {
		return 0, err
	}

	n := buffer.Size() - bufferOffset
	if n == 0 {
		size = packUint32(^uint32(0))
	} else {
		size = packUint32(uint32(n) - 4)
	}
	buffer.WriteAt(size[:], bufferOffset)

	
	
	
	if buffer != w {
		return buffer.WriteTo(w)
	}

	return n, nil
}



type RawRecordSet struct {
	
	Reader io.Reader
}








func (rrs *RawRecordSet) ReadFrom(r io.Reader) (int64, error) {
	rs := &RecordSet{}
	n, err := rs.ReadFrom(r)
	if err != nil {
		return 0, err
	}

	buf := &bytes.Buffer{}
	rs.WriteTo(buf)
	*rrs = RawRecordSet{
		Reader: buf,
	}

	return n, nil
}



func (rrs *RawRecordSet) WriteTo(w io.Writer) (int64, error) {
	if rrs.Reader == nil {
		return 0, ErrNoRecord
	}

	return io.Copy(w, rrs.Reader)
}

func makeTime(t int64) time.Time {
	return time.Unix(t/1000, (t%1000)*int64(time.Millisecond))
}

func timestamp(t time.Time) int64 {
	if t.IsZero() {
		return 0
	}
	return t.UnixNano() / int64(time.Millisecond)
}

func packUint32(u uint32) (b [4]byte) {
	binary.BigEndian.PutUint32(b[:], u)
	return
}

func packUint64(u uint64) (b [8]byte) {
	binary.BigEndian.PutUint64(b[:], u)
	return
}
