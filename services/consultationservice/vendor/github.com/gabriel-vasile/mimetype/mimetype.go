




package mimetype

import (
	"io"
	"mime"
	"os"
	"sync/atomic"
)

var defaultLimit uint32 = 3072


var readLimit uint32 = defaultLimit





func Detect(in []byte) *MIME {
	
	l := atomic.LoadUint32(&readLimit)
	if l > 0 && len(in) > int(l) {
		in = in[:l]
	}
	mu.RLock()
	defer mu.RUnlock()
	return root.match(in, l)
}











func DetectReader(r io.Reader) (*MIME, error) {
	var in []byte
	var err error

	
	l := atomic.LoadUint32(&readLimit)
	if l == 0 {
		in, err = io.ReadAll(r)
		if err != nil {
			return errMIME, err
		}
	} else {
		var n int
		in = make([]byte, l)
		
		
		n, err = io.ReadFull(r, in)
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			return errMIME, err
		}
		in = in[:n]
	}

	mu.RLock()
	defer mu.RUnlock()
	return root.match(in, l), nil
}






func DetectFile(path string) (*MIME, error) {
	f, err := os.Open(path)
	if err != nil {
		return errMIME, err
	}
	defer f.Close()

	return DetectReader(f)
}





func EqualsAny(s string, mimes ...string) bool {
	s, _, _ = mime.ParseMediaType(s)
	for _, m := range mimes {
		m, _, _ = mime.ParseMediaType(m)
		if s == m {
			return true
		}
	}

	return false
}






func SetLimit(limit uint32) {
	
	atomic.StoreUint32(&readLimit, limit)
}



func Extend(detector func(raw []byte, limit uint32) bool, mime, extension string, aliases ...string) {
	root.Extend(detector, mime, extension, aliases...)
}



func Lookup(mime string) *MIME {
	mu.RLock()
	defer mu.RUnlock()
	return root.lookup(mime)
}
