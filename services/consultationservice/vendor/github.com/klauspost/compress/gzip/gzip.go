



package gzip

import (
	"errors"
	"fmt"
	"hash/crc32"
	"io"

	"github.com/klauspost/compress/flate"
)



const (
	NoCompression       = flate.NoCompression
	BestSpeed           = flate.BestSpeed
	BestCompression     = flate.BestCompression
	DefaultCompression  = flate.DefaultCompression
	ConstantCompression = flate.ConstantCompression
	HuffmanOnly         = flate.HuffmanOnly

	
	
	
	
	
	StatelessCompression = -3
)



type Writer struct {
	Header      
	w           io.Writer
	level       int
	err         error
	compressor  *flate.Writer
	digest      uint32 
	size        uint32 
	wroteHeader bool
	closed      bool
	buf         [10]byte
}









func NewWriter(w io.Writer) *Writer {
	z, _ := NewWriterLevel(w, DefaultCompression)
	return z
}







func NewWriterLevel(w io.Writer, level int) (*Writer, error) {
	if level < StatelessCompression || level > BestCompression {
		return nil, fmt.Errorf("gzip: invalid compression level: %d", level)
	}
	z := new(Writer)
	z.init(w, level)
	return z, nil
}

func (z *Writer) init(w io.Writer, level int) {
	compressor := z.compressor
	if level != StatelessCompression {
		if compressor != nil {
			compressor.Reset(w)
		}
	}

	*z = Writer{
		Header: Header{
			OS: 255, 
		},
		w:          w,
		level:      level,
		compressor: compressor,
	}
}





func (z *Writer) Reset(w io.Writer) {
	z.init(w, z.level)
}


func (z *Writer) writeBytes(b []byte) error {
	if len(b) > 0xffff {
		return errors.New("gzip.Write: Extra data is too large")
	}
	le.PutUint16(z.buf[:2], uint16(len(b)))
	_, err := z.w.Write(z.buf[:2])
	if err != nil {
		return err
	}
	_, err = z.w.Write(b)
	return err
}



func (z *Writer) writeString(s string) (err error) {
	
	needconv := false
	for _, v := range s {
		if v == 0 || v > 0xff {
			return errors.New("gzip.Write: non-Latin-1 header string")
		}
		if v > 0x7f {
			needconv = true
		}
	}
	if needconv {
		b := make([]byte, 0, len(s))
		for _, v := range s {
			b = append(b, byte(v))
		}
		_, err = z.w.Write(b)
	} else {
		_, err = io.WriteString(z.w, s)
	}
	if err != nil {
		return err
	}
	
	z.buf[0] = 0
	_, err = z.w.Write(z.buf[:1])
	return err
}



func (z *Writer) Write(p []byte) (int, error) {
	if z.err != nil {
		return 0, z.err
	}
	var n int
	
	if !z.wroteHeader {
		z.wroteHeader = true
		z.buf[0] = gzipID1
		z.buf[1] = gzipID2
		z.buf[2] = gzipDeflate
		z.buf[3] = 0
		if z.Extra != nil {
			z.buf[3] |= 0x04
		}
		if z.Name != "" {
			z.buf[3] |= 0x08
		}
		if z.Comment != "" {
			z.buf[3] |= 0x10
		}
		le.PutUint32(z.buf[4:8], uint32(z.ModTime.Unix()))
		if z.level == BestCompression {
			z.buf[8] = 2
		} else if z.level == BestSpeed {
			z.buf[8] = 4
		} else {
			z.buf[8] = 0
		}
		z.buf[9] = z.OS
		n, z.err = z.w.Write(z.buf[:10])
		if z.err != nil {
			return n, z.err
		}
		if z.Extra != nil {
			z.err = z.writeBytes(z.Extra)
			if z.err != nil {
				return n, z.err
			}
		}
		if z.Name != "" {
			z.err = z.writeString(z.Name)
			if z.err != nil {
				return n, z.err
			}
		}
		if z.Comment != "" {
			z.err = z.writeString(z.Comment)
			if z.err != nil {
				return n, z.err
			}
		}

		if z.compressor == nil && z.level != StatelessCompression {
			z.compressor, _ = flate.NewWriter(z.w, z.level)
		}
	}
	z.size += uint32(len(p))
	z.digest = crc32.Update(z.digest, crc32.IEEETable, p)
	if z.level == StatelessCompression {
		return len(p), flate.StatelessDeflate(z.w, p, false, nil)
	}
	n, z.err = z.compressor.Write(p)
	return n, z.err
}









func (z *Writer) Flush() error {
	if z.err != nil {
		return z.err
	}
	if z.closed || z.level == StatelessCompression {
		return nil
	}
	if !z.wroteHeader {
		z.Write(nil)
		if z.err != nil {
			return z.err
		}
	}
	z.err = z.compressor.Flush()
	return z.err
}



func (z *Writer) Close() error {
	if z.err != nil {
		return z.err
	}
	if z.closed {
		return nil
	}
	z.closed = true
	if !z.wroteHeader {
		z.Write(nil)
		if z.err != nil {
			return z.err
		}
	}
	if z.level == StatelessCompression {
		z.err = flate.StatelessDeflate(z.w, nil, true, nil)
	} else {
		z.err = z.compressor.Close()
	}
	if z.err != nil {
		return z.err
	}
	le.PutUint32(z.buf[:4], z.digest)
	le.PutUint32(z.buf[4:8], z.size)
	_, z.err = z.w.Write(z.buf[:8])
	return z.err
}
