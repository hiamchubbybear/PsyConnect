package zstd

import (
	"errors"
	"fmt"
	"math"
	"math/bits"
	"runtime"
	"strings"
)


type EOption func(*encoderOptions) error


type encoderOptions struct {
	concurrent      int
	level           EncoderLevel
	single          *bool
	pad             int
	blockSize       int
	windowSize      int
	crc             bool
	fullZero        bool
	noEntropy       bool
	allLitEntropy   bool
	customWindow    bool
	customALEntropy bool
	customBlockSize bool
	lowMem          bool
	dict            *dict
}

func (o *encoderOptions) setDefault() {
	*o = encoderOptions{
		concurrent:    runtime.GOMAXPROCS(0),
		crc:           true,
		single:        nil,
		blockSize:     maxCompressedBlockSize,
		windowSize:    8 << 20,
		level:         SpeedDefault,
		allLitEntropy: false,
		lowMem:        false,
	}
}


func (o encoderOptions) encoder() encoder {
	switch o.level {
	case SpeedFastest:
		if o.dict != nil {
			return &fastEncoderDict{fastEncoder: fastEncoder{fastBase: fastBase{maxMatchOff: int32(o.windowSize), bufferReset: math.MaxInt32 - int32(o.windowSize*2), lowMem: o.lowMem}}}
		}
		return &fastEncoder{fastBase: fastBase{maxMatchOff: int32(o.windowSize), bufferReset: math.MaxInt32 - int32(o.windowSize*2), lowMem: o.lowMem}}

	case SpeedDefault:
		if o.dict != nil {
			return &doubleFastEncoderDict{fastEncoderDict: fastEncoderDict{fastEncoder: fastEncoder{fastBase: fastBase{maxMatchOff: int32(o.windowSize), bufferReset: math.MaxInt32 - int32(o.windowSize*2), lowMem: o.lowMem}}}}
		}
		return &doubleFastEncoder{fastEncoder: fastEncoder{fastBase: fastBase{maxMatchOff: int32(o.windowSize), bufferReset: math.MaxInt32 - int32(o.windowSize*2), lowMem: o.lowMem}}}
	case SpeedBetterCompression:
		if o.dict != nil {
			return &betterFastEncoderDict{betterFastEncoder: betterFastEncoder{fastBase: fastBase{maxMatchOff: int32(o.windowSize), bufferReset: math.MaxInt32 - int32(o.windowSize*2), lowMem: o.lowMem}}}
		}
		return &betterFastEncoder{fastBase: fastBase{maxMatchOff: int32(o.windowSize), bufferReset: math.MaxInt32 - int32(o.windowSize*2), lowMem: o.lowMem}}
	case SpeedBestCompression:
		return &bestFastEncoder{fastBase: fastBase{maxMatchOff: int32(o.windowSize), bufferReset: math.MaxInt32 - int32(o.windowSize*2), lowMem: o.lowMem}}
	}
	panic("unknown compression level")
}



func WithEncoderCRC(b bool) EOption {
	return func(o *encoderOptions) error { o.crc = b; return nil }
}






func WithEncoderConcurrency(n int) EOption {
	return func(o *encoderOptions) error {
		if n <= 0 {
			return fmt.Errorf("concurrency must be at least 1")
		}
		o.concurrent = n
		return nil
	}
}






func WithWindowSize(n int) EOption {
	return func(o *encoderOptions) error {
		switch {
		case n < MinWindowSize:
			return fmt.Errorf("window size must be at least %d", MinWindowSize)
		case n > MaxWindowSize:
			return fmt.Errorf("window size must be at most %d", MaxWindowSize)
		case (n & (n - 1)) != 0:
			return errors.New("window size must be a power of 2")
		}

		o.windowSize = n
		o.customWindow = true
		if o.blockSize > o.windowSize {
			o.blockSize = o.windowSize
			o.customBlockSize = true
		}
		return nil
	}
}







func WithEncoderPadding(n int) EOption {
	return func(o *encoderOptions) error {
		if n <= 0 {
			return fmt.Errorf("padding must be at least 1")
		}
		
		if n == 1 {
			n = 0
		}
		if n > 1<<30 {
			return fmt.Errorf("padding must less than 1GB (1<<30 bytes) ")
		}
		o.pad = n
		return nil
	}
}





type EncoderLevel int

const (
	speedNotSet EncoderLevel = iota

	
	
	SpeedFastest

	
	
	SpeedDefault

	
	
	
	SpeedBetterCompression

	
	
	SpeedBestCompression

	
	
	speedLast
)




func EncoderLevelFromString(s string) (bool, EncoderLevel) {
	for l := speedNotSet + 1; l < speedLast; l++ {
		if strings.EqualFold(s, l.String()) {
			return true, l
		}
	}
	return false, SpeedDefault
}




func EncoderLevelFromZstd(level int) EncoderLevel {
	switch {
	case level < 3:
		return SpeedFastest
	case level >= 3 && level < 6:
		return SpeedDefault
	case level >= 6 && level < 10:
		return SpeedBetterCompression
	default:
		return SpeedBestCompression
	}
}


func (e EncoderLevel) String() string {
	switch e {
	case SpeedFastest:
		return "fastest"
	case SpeedDefault:
		return "default"
	case SpeedBetterCompression:
		return "better"
	case SpeedBestCompression:
		return "best"
	default:
		return "invalid"
	}
}


func WithEncoderLevel(l EncoderLevel) EOption {
	return func(o *encoderOptions) error {
		switch {
		case l <= speedNotSet || l >= speedLast:
			return fmt.Errorf("unknown encoder level")
		}
		o.level = l
		if !o.customWindow {
			switch o.level {
			case SpeedFastest:
				o.windowSize = 4 << 20
				if !o.customBlockSize {
					o.blockSize = 1 << 16
				}
			case SpeedDefault:
				o.windowSize = 8 << 20
			case SpeedBetterCompression:
				o.windowSize = 16 << 20
			case SpeedBestCompression:
				o.windowSize = 32 << 20
			}
		}
		if !o.customALEntropy {
			o.allLitEntropy = l > SpeedDefault
		}

		return nil
	}
}




func WithZeroFrames(b bool) EOption {
	return func(o *encoderOptions) error {
		o.fullZero = b
		return nil
	}
}





func WithAllLitEntropyCompression(b bool) EOption {
	return func(o *encoderOptions) error {
		o.customALEntropy = true
		o.allLitEntropy = b
		return nil
	}
}




func WithNoEntropyCompression(b bool) EOption {
	return func(o *encoderOptions) error {
		o.noEntropy = b
		return nil
	}
}











func WithSingleSegment(b bool) EOption {
	return func(o *encoderOptions) error {
		o.single = &b
		return nil
	}
}





func WithLowerEncoderMem(b bool) EOption {
	return func(o *encoderOptions) error {
		o.lowMem = b
		return nil
	}
}









func WithEncoderDict(dict []byte) EOption {
	return func(o *encoderOptions) error {
		d, err := loadDict(dict)
		if err != nil {
			return err
		}
		o.dict = d
		return nil
	}
}





func WithEncoderDictRaw(id uint32, content []byte) EOption {
	return func(o *encoderOptions) error {
		if bits.UintSize > 32 && uint(len(content)) > dictMaxLength {
			return fmt.Errorf("dictionary of size %d > 2GiB too large", len(content))
		}
		o.dict = &dict{id: id, content: content, offsets: [3]int{1, 4, 8}}
		return nil
	}
}
