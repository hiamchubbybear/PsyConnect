package compress

import (
	"encoding"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/segmentio/kafka-go/compress/gzip"
	"github.com/segmentio/kafka-go/compress/lz4"
	"github.com/segmentio/kafka-go/compress/snappy"
	"github.com/segmentio/kafka-go/compress/zstd"
)


type Compression int8

const (
	None   Compression = 0
	Gzip   Compression = 1
	Snappy Compression = 2
	Lz4    Compression = 3
	Zstd   Compression = 4
)

func (c Compression) Codec() Codec {
	if i := int(c); i >= 0 && i < len(Codecs) {
		return Codecs[i]
	}
	return nil
}

func (c Compression) String() string {
	if codec := c.Codec(); codec != nil {
		return codec.Name()
	}
	return "uncompressed"
}

func (c Compression) MarshalText() ([]byte, error) {
	return []byte(c.String()), nil
}

func (c *Compression) UnmarshalText(b []byte) error {
	switch string(b) {
	case "none", "uncompressed":
		*c = None
		return nil
	}

	for _, codec := range Codecs[None+1:] {
		if codec.Name() == string(b) {
			*c = Compression(codec.Code())
			return nil
		}
	}

	i, err := strconv.ParseInt(string(b), 10, 64)
	if err == nil && i >= 0 && i < int64(len(Codecs)) {
		*c = Compression(i)
		return nil
	}

	s := &strings.Builder{}
	s.WriteString("none, uncompressed")

	for i, codec := range Codecs[None+1:] {
		if i < (len(Codecs) - 1) {
			s.WriteString(", ")
		} else {
			s.WriteString(", or ")
		}
		s.WriteString(codec.Name())
	}

	return fmt.Errorf("compression format must be one of %s, not %q", s, b)
}

var (
	_ encoding.TextMarshaler   = Compression(0)
	_ encoding.TextUnmarshaler = (*Compression)(nil)
)





type Codec interface {
	
	Code() int8

	
	Name() string

	
	NewReader(r io.Reader) io.ReadCloser

	
	NewWriter(w io.Writer) io.WriteCloser
}

var (
	
	GzipCodec gzip.Codec

	
	SnappyCodec snappy.Codec

	
	Lz4Codec lz4.Codec

	
	ZstdCodec zstd.Codec

	
	Codecs = [...]Codec{
		None:   nil,
		Gzip:   &GzipCodec,
		Snappy: &SnappyCodec,
		Lz4:    &Lz4Codec,
		Zstd:   &ZstdCodec,
	}
)
