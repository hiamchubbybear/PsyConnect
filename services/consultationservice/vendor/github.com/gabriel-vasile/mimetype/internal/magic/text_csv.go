package magic

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"errors"
	"io"
	"sync"
)


var readerPool = sync.Pool{
	New: func() any {
		
		return bufio.NewReader(nil)
	},
}

func newReader(r io.Reader) *bufio.Reader {
	br := readerPool.Get().(*bufio.Reader)
	br.Reset(r)
	return br
}


func Csv(raw []byte, limit uint32) bool {
	return sv(raw, ',', limit)
}


func Tsv(raw []byte, limit uint32) bool {
	return sv(raw, '\t', limit)
}

func sv(in []byte, comma rune, limit uint32) bool {
	in = dropLastLine(in, limit)

	br := newReader(bytes.NewReader(in))
	defer readerPool.Put(br)
	r := csv.NewReader(br)
	r.Comma = comma
	r.ReuseRecord = true
	r.LazyQuotes = true
	r.Comment = '#'

	lines := 0
	for {
		_, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return false
		}
		lines++
	}

	return r.FieldsPerRecord > 1 && lines > 1
}






func dropLastLine(b []byte, readLimit uint32) []byte {
	if readLimit == 0 || uint32(len(b)) < readLimit {
		return b
	}
	for i := len(b) - 1; i > 0; i-- {
		if b[i] == '\n' {
			return b[:i]
		}
	}
	return b
}
