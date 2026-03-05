


package zstd

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"math"
)


const debug = false


const debugEncoder = debug


const debugDecoder = debug


const debugAsserts = debug || false


const debugSequences = false


const debugMatches = false


const forcePreDef = false


const zstdMinMatch = 3


const fcsUnknown = math.MaxUint64

var (
	
	
	ErrReservedBlockType = errors.New("invalid input: reserved block type encountered")

	
	
	ErrCompressedSizeTooBig = errors.New("invalid input: compressed size too big")

	
	
	ErrBlockTooSmall = errors.New("block too small")

	
	
	ErrUnexpectedBlockSize = errors.New("unexpected block size")

	
	
	ErrMagicMismatch = errors.New("invalid input: magic number mismatch")

	
	
	ErrWindowSizeExceeded = errors.New("window size exceeded")

	
	
	ErrWindowSizeTooSmall = errors.New("invalid input: window size was too small")

	
	ErrDecoderSizeExceeded = errors.New("decompressed size exceeds configured limit")

	
	ErrUnknownDictionary = errors.New("unknown dictionary")

	
	
	ErrFrameSizeExceeded = errors.New("frame size exceeded")

	
	
	ErrFrameSizeMismatch = errors.New("frame size does not match size on stream")

	
	ErrCRCMismatch = errors.New("CRC check failed")

	
	
	ErrDecoderClosed = errors.New("decoder used after Close")

	
	
	ErrDecoderNilInput = errors.New("nil input provided as reader")
)

func println(a ...interface{}) {
	if debug || debugDecoder || debugEncoder {
		log.Println(a...)
	}
}

func printf(format string, a ...interface{}) {
	if debug || debugDecoder || debugEncoder {
		log.Printf(format, a...)
	}
}

func load3232(b []byte, i int32) uint32 {
	return binary.LittleEndian.Uint32(b[:len(b):len(b)][i:])
}

func load6432(b []byte, i int32) uint64 {
	return binary.LittleEndian.Uint64(b[:len(b):len(b)][i:])
}

type byter interface {
	Bytes() []byte
	Len() int
}

var _ byter = &bytes.Buffer{}
