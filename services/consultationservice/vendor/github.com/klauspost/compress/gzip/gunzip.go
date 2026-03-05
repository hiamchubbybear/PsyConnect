





package gzip

import (
	"bufio"
	"compress/gzip"
	"encoding/binary"
	"hash/crc32"
	"io"
	"time"

	"github.com/klauspost/compress/flate"
)

const (
	gzipID1     = 0x1f
	gzipID2     = 0x8b
	gzipDeflate = 8
	flagText    = 1 << 0
	flagHdrCrc  = 1 << 1
	flagExtra   = 1 << 2
	flagName    = 1 << 3
	flagComment = 1 << 4
)

var (
	
	ErrChecksum = gzip.ErrChecksum
	
	ErrHeader = gzip.ErrHeader
)

var le = binary.LittleEndian


func noEOF(err error) error {
	if err == io.EOF {
		return io.ErrUnexpectedEOF
	}
	return err
}






type Header struct {
	Comment string    
	Extra   []byte    
	ModTime time.Time 
	Name    string    
	OS      byte      
}















type Reader struct {
	Header       
	r            flate.Reader
	br           *bufio.Reader
	decompressor io.ReadCloser
	digest       uint32 
	size         uint32 
	buf          [512]byte
	err          error
	multistream  bool
}








func NewReader(r io.Reader) (*Reader, error) {
	z := new(Reader)
	if err := z.Reset(r); err != nil {
		return nil, err
	}
	return z, nil
}




func (z *Reader) Reset(r io.Reader) error {
	*z = Reader{
		decompressor: z.decompressor,
		multistream:  true,
	}
	if rr, ok := r.(flate.Reader); ok {
		z.r = rr
	} else {
		
		if z.br != nil {
			z.br.Reset(r)
		} else {
			z.br = bufio.NewReader(r)
		}
		z.r = z.br
	}
	z.Header, z.err = z.readHeader()
	return z.err
}

















func (z *Reader) Multistream(ok bool) {
	z.multistream = ok
}





func (z *Reader) readString() (string, error) {
	var err error
	needConv := false
	for i := 0; ; i++ {
		if i >= len(z.buf) {
			return "", ErrHeader
		}
		z.buf[i], err = z.r.ReadByte()
		if err != nil {
			return "", err
		}
		if z.buf[i] > 0x7f {
			needConv = true
		}
		if z.buf[i] == 0 {
			
			z.digest = crc32.Update(z.digest, crc32.IEEETable, z.buf[:i+1])

			
			if needConv {
				s := make([]rune, 0, i)
				for _, v := range z.buf[:i] {
					s = append(s, rune(v))
				}
				return string(s), nil
			}
			return string(z.buf[:i]), nil
		}
	}
}



func (z *Reader) readHeader() (hdr Header, err error) {
	if _, err = io.ReadFull(z.r, z.buf[:10]); err != nil {
		
		
		
		
		
		
		
		return hdr, err
	}
	if z.buf[0] != gzipID1 || z.buf[1] != gzipID2 || z.buf[2] != gzipDeflate {
		return hdr, ErrHeader
	}
	flg := z.buf[3]
	hdr.ModTime = time.Unix(int64(le.Uint32(z.buf[4:8])), 0)
	
	hdr.OS = z.buf[9]
	z.digest = crc32.ChecksumIEEE(z.buf[:10])

	if flg&flagExtra != 0 {
		if _, err = io.ReadFull(z.r, z.buf[:2]); err != nil {
			return hdr, noEOF(err)
		}
		z.digest = crc32.Update(z.digest, crc32.IEEETable, z.buf[:2])
		data := make([]byte, le.Uint16(z.buf[:2]))
		if _, err = io.ReadFull(z.r, data); err != nil {
			return hdr, noEOF(err)
		}
		z.digest = crc32.Update(z.digest, crc32.IEEETable, data)
		hdr.Extra = data
	}

	var s string
	if flg&flagName != 0 {
		if s, err = z.readString(); err != nil {
			return hdr, err
		}
		hdr.Name = s
	}

	if flg&flagComment != 0 {
		if s, err = z.readString(); err != nil {
			return hdr, err
		}
		hdr.Comment = s
	}

	if flg&flagHdrCrc != 0 {
		if _, err = io.ReadFull(z.r, z.buf[:2]); err != nil {
			return hdr, noEOF(err)
		}
		digest := le.Uint16(z.buf[:2])
		if digest != uint16(z.digest) {
			return hdr, ErrHeader
		}
	}

	z.digest = 0
	if z.decompressor == nil {
		z.decompressor = flate.NewReader(z.r)
	} else {
		z.decompressor.(flate.Resetter).Reset(z.r, nil)
	}
	return hdr, nil
}


func (z *Reader) Read(p []byte) (n int, err error) {
	if z.err != nil {
		return 0, z.err
	}

	for n == 0 {
		n, z.err = z.decompressor.Read(p)
		z.digest = crc32.Update(z.digest, crc32.IEEETable, p[:n])
		z.size += uint32(n)
		if z.err != io.EOF {
			
			return n, z.err
		}

		
		if _, err := io.ReadFull(z.r, z.buf[:8]); err != nil {
			z.err = noEOF(err)
			return n, z.err
		}
		digest := le.Uint32(z.buf[:4])
		size := le.Uint32(z.buf[4:8])
		if digest != z.digest || size != z.size {
			z.err = ErrChecksum
			return n, z.err
		}
		z.digest, z.size = 0, 0

		
		if !z.multistream {
			return n, io.EOF
		}
		z.err = nil 

		if _, z.err = z.readHeader(); z.err != nil {
			return n, z.err
		}
	}

	return n, nil
}

type crcer interface {
	io.Writer
	Sum32() uint32
	Reset()
}
type crcUpdater struct {
	z *Reader
}

func (c *crcUpdater) Write(p []byte) (int, error) {
	c.z.digest = crc32.Update(c.z.digest, crc32.IEEETable, p)
	return len(p), nil
}

func (c *crcUpdater) Sum32() uint32 {
	return c.z.digest
}

func (c *crcUpdater) Reset() {
	c.z.digest = 0
}


func (z *Reader) WriteTo(w io.Writer) (int64, error) {
	total := int64(0)
	crcWriter := crcer(crc32.NewIEEE())
	if z.digest != 0 {
		crcWriter = &crcUpdater{z: z}
	}
	for {
		if z.err != nil {
			if z.err == io.EOF {
				return total, nil
			}
			return total, z.err
		}

		
		mw := io.MultiWriter(w, crcWriter)
		n, err := z.decompressor.(io.WriterTo).WriteTo(mw)
		total += n
		z.size += uint32(n)
		if err != nil {
			z.err = err
			return total, z.err
		}

		
		if _, err := io.ReadFull(z.r, z.buf[0:8]); err != nil {
			if err == io.EOF {
				err = io.ErrUnexpectedEOF
			}
			z.err = err
			return total, err
		}
		z.digest = crcWriter.Sum32()
		digest := le.Uint32(z.buf[:4])
		size := le.Uint32(z.buf[4:8])
		if digest != z.digest || size != z.size {
			z.err = ErrChecksum
			return total, z.err
		}
		z.digest, z.size = 0, 0

		
		if !z.multistream {
			return total, nil
		}
		crcWriter.Reset()
		z.err = nil 

		if _, z.err = z.readHeader(); z.err != nil {
			if z.err == io.EOF {
				return total, nil
			}
			return total, z.err
		}
	}
}




func (z *Reader) Close() error { return z.decompressor.Close() }
