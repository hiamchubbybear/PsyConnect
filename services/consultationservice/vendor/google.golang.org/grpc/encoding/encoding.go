








package encoding

import (
	"io"
	"strings"

	"google.golang.org/grpc/internal/grpcutil"
)



const Identity = "identity"







type Compressor interface {
	
	
	
	Compress(w io.Writer) (io.WriteCloser, error)
	
	
	
	Decompress(r io.Reader) (io.Reader, error)
	
	
	
	Name() string
}

var registeredCompressor = make(map[string]Compressor)










func RegisterCompressor(c Compressor) {
	registeredCompressor[c.Name()] = c
	if !grpcutil.IsCompressorNameRegistered(c.Name()) {
		grpcutil.RegisteredCompressorNames = append(grpcutil.RegisteredCompressorNames, c.Name())
	}
}


func GetCompressor(name string) Compressor {
	return registeredCompressor[name]
}




type Codec interface {
	
	Marshal(v any) ([]byte, error)
	
	Unmarshal(data []byte, v any) error
	
	
	
	Name() string
}

var registeredCodecs = make(map[string]any)















func RegisterCodec(codec Codec) {
	if codec == nil {
		panic("cannot register a nil Codec")
	}
	if codec.Name() == "" {
		panic("cannot register Codec with empty string result for Name()")
	}
	contentSubtype := strings.ToLower(codec.Name())
	registeredCodecs[contentSubtype] = codec
}





func GetCodec(contentSubtype string) Codec {
	c, _ := registeredCodecs[contentSubtype].(Codec)
	return c
}
