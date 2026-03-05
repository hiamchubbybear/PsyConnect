



package zstd

import (
	"errors"
	"fmt"
	"math/bits"
	"runtime"
)


type DOption func(*decoderOptions) error


type decoderOptions struct {
	lowMem          bool
	concurrent      int
	maxDecodedSize  uint64
	maxWindowSize   uint64
	dicts           []*dict
	ignoreChecksum  bool
	limitToCap      bool
	decodeBufsBelow int
}

func (o *decoderOptions) setDefault() {
	*o = decoderOptions{
		
		lowMem:          true,
		concurrent:      runtime.GOMAXPROCS(0),
		maxWindowSize:   MaxWindowSize,
		decodeBufsBelow: 128 << 10,
	}
	if o.concurrent > 4 {
		o.concurrent = 4
	}
	o.maxDecodedSize = 64 << 30
}



func WithDecoderLowmem(b bool) DOption {
	return func(o *decoderOptions) error { o.lowMem = b; return nil }
}










func WithDecoderConcurrency(n int) DOption {
	return func(o *decoderOptions) error {
		if n < 0 {
			return errors.New("concurrency must be at least 1")
		}
		if n == 0 {
			o.concurrent = runtime.GOMAXPROCS(0)
		} else {
			o.concurrent = n
		}
		return nil
	}
}





func WithDecoderMaxMemory(n uint64) DOption {
	return func(o *decoderOptions) error {
		if n == 0 {
			return errors.New("WithDecoderMaxMemory must be at least 1")
		}
		if n > 1<<63 {
			return errors.New("WithDecoderMaxmemory must be less than 1 << 63")
		}
		o.maxDecodedSize = n
		return nil
	}
}









func WithDecoderDicts(dicts ...[]byte) DOption {
	return func(o *decoderOptions) error {
		for _, b := range dicts {
			d, err := loadDict(b)
			if err != nil {
				return err
			}
			o.dicts = append(o.dicts, d)
		}
		return nil
	}
}



func WithDecoderDictRaw(id uint32, content []byte) DOption {
	return func(o *decoderOptions) error {
		if bits.UintSize > 32 && uint(len(content)) > dictMaxLength {
			return fmt.Errorf("dictionary of size %d > 2GiB too large", len(content))
		}
		o.dicts = append(o.dicts, &dict{id: id, content: content, offsets: [3]int{1, 4, 8}})
		return nil
	}
}






func WithDecoderMaxWindow(size uint64) DOption {
	return func(o *decoderOptions) error {
		if size < MinWindowSize {
			return errors.New("WithMaxWindowSize must be at least 1KB, 1024 bytes")
		}
		if size > (1<<41)+7*(1<<38) {
			return errors.New("WithMaxWindowSize must be less than (1<<41) + 7*(1<<38) ~ 3.75TB")
		}
		o.maxWindowSize = size
		return nil
	}
}





func WithDecodeAllCapLimit(b bool) DOption {
	return func(o *decoderOptions) error {
		o.limitToCap = b
		return nil
	}
}






func WithDecodeBuffersBelow(size int) DOption {
	return func(o *decoderOptions) error {
		o.decodeBufsBelow = size
		return nil
	}
}


func IgnoreChecksum(b bool) DOption {
	return func(o *decoderOptions) error {
		o.ignoreChecksum = b
		return nil
	}
}
