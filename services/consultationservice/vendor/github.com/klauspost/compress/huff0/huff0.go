


package huff0

import (
	"errors"
	"fmt"
	"math"
	"math/bits"
	"sync"

	"github.com/klauspost/compress/fse"
)

const (
	maxSymbolValue = 255

	
	
	tableLogMax     = 11
	tableLogDefault = 11
	minTablelog     = 5
	huffNodesLen    = 512

	
	BlockSizeMax = 1<<18 - 1
)

var (
	
	ErrIncompressible = errors.New("input is not compressible")

	
	ErrUseRLE = errors.New("input is single value repeated")

	
	ErrTooBig = errors.New("input too big")

	
	ErrMaxDecodedSizeExceeded = errors.New("maximum output size exceeded")
)

type ReusePolicy uint8

const (
	
	ReusePolicyAllow ReusePolicy = iota

	
	
	
	
	ReusePolicyPrefer

	
	
	ReusePolicyNone

	
	ReusePolicyMust
)

type Scratch struct {
	count [maxSymbolValue + 1]uint32

	
	
	

	
	
	
	
	
	Out []byte

	
	
	OutTable []byte

	
	
	OutData []byte

	
	
	
	MaxDecodedSize int

	br byteReader

	
	MaxSymbolValue uint8

	
	
	TableLog uint8

	
	Reuse ReusePolicy

	
	
	
	
	WantLogLess uint8

	symbolLen      uint16 
	maxCount       int    
	clearCount     bool   
	actualTableLog uint8  
	prevTableLog   uint8  
	prevTable      cTable 
	cTable         cTable 
	dt             dTable 
	nodes          []nodeElt
	tmpOut         [4][]byte
	fse            *fse.Scratch
	decPool        sync.Pool 
	huffWeight     [maxSymbolValue + 1]byte
}


func (s *Scratch) TransferCTable(src *Scratch) {
	if cap(s.prevTable) < len(src.prevTable) {
		s.prevTable = make(cTable, 0, maxSymbolValue+1)
	}
	s.prevTable = s.prevTable[:len(src.prevTable)]
	copy(s.prevTable, src.prevTable)
	s.prevTableLog = src.prevTableLog
}

func (s *Scratch) prepare(in []byte) (*Scratch, error) {
	if len(in) > BlockSizeMax {
		return nil, ErrTooBig
	}
	if s == nil {
		s = &Scratch{}
	}
	if s.MaxSymbolValue == 0 {
		s.MaxSymbolValue = maxSymbolValue
	}
	if s.TableLog == 0 {
		s.TableLog = tableLogDefault
	}
	if s.TableLog > tableLogMax || s.TableLog < minTablelog {
		return nil, fmt.Errorf(" invalid tableLog %d (%d -> %d)", s.TableLog, minTablelog, tableLogMax)
	}
	if s.MaxDecodedSize <= 0 || s.MaxDecodedSize > BlockSizeMax {
		s.MaxDecodedSize = BlockSizeMax
	}
	if s.clearCount && s.maxCount == 0 {
		for i := range s.count {
			s.count[i] = 0
		}
		s.clearCount = false
	}
	if cap(s.Out) == 0 {
		s.Out = make([]byte, 0, len(in))
	}
	s.Out = s.Out[:0]

	s.OutTable = nil
	s.OutData = nil
	if cap(s.nodes) < huffNodesLen+1 {
		s.nodes = make([]nodeElt, 0, huffNodesLen+1)
	}
	s.nodes = s.nodes[:0]
	if s.fse == nil {
		s.fse = &fse.Scratch{}
	}
	s.br.init(in)

	return s, nil
}

type cTable []cTableEntry

func (c cTable) write(s *Scratch) error {
	var (
		
		bitsToWeight [tableLogMax + 1]byte
		huffLog      = s.actualTableLog
		
		maxSymbolValue = uint8(s.symbolLen - 1)
		huffWeight     = s.huffWeight[:256]
	)
	const (
		maxFSETableLog = 6
	)
	
	bitsToWeight[0] = 0
	for n := uint8(1); n < huffLog+1; n++ {
		bitsToWeight[n] = huffLog + 1 - n
	}

	
	hist := s.fse.Histogram()
	hist = hist[:256]
	for i := range hist[:16] {
		hist[i] = 0
	}
	for n := uint8(0); n < maxSymbolValue; n++ {
		v := bitsToWeight[c[n].nBits] & 15
		huffWeight[n] = v
		hist[v]++
	}

	
	if maxSymbolValue >= 2 {
		huffMaxCnt := uint32(0)
		huffMax := uint8(0)
		for i, v := range hist[:16] {
			if v == 0 {
				continue
			}
			huffMax = byte(i)
			if v > huffMaxCnt {
				huffMaxCnt = v
			}
		}
		s.fse.HistogramFinished(huffMax, int(huffMaxCnt))
		s.fse.TableLog = maxFSETableLog
		b, err := fse.Compress(huffWeight[:maxSymbolValue], s.fse)
		if err == nil && len(b) < int(s.symbolLen>>1) {
			s.Out = append(s.Out, uint8(len(b)))
			s.Out = append(s.Out, b...)
			return nil
		}
		
	}
	
	if maxSymbolValue > (256 - 128) {
		
		return ErrIncompressible
	}
	op := s.Out
	
	op = append(op, 128|(maxSymbolValue-1))
	
	huffWeight[maxSymbolValue] = 0
	for n := uint16(0); n < uint16(maxSymbolValue); n += 2 {
		op = append(op, (huffWeight[n]<<4)|huffWeight[n+1])
	}
	s.Out = op
	return nil
}

func (c cTable) estTableSize(s *Scratch) (sz int, err error) {
	var (
		
		bitsToWeight [tableLogMax + 1]byte
		huffLog      = s.actualTableLog
		
		maxSymbolValue = uint8(s.symbolLen - 1)
		huffWeight     = s.huffWeight[:256]
	)
	const (
		maxFSETableLog = 6
	)
	
	bitsToWeight[0] = 0
	for n := uint8(1); n < huffLog+1; n++ {
		bitsToWeight[n] = huffLog + 1 - n
	}

	
	hist := s.fse.Histogram()
	hist = hist[:256]
	for i := range hist[:16] {
		hist[i] = 0
	}
	for n := uint8(0); n < maxSymbolValue; n++ {
		v := bitsToWeight[c[n].nBits] & 15
		huffWeight[n] = v
		hist[v]++
	}

	
	if maxSymbolValue >= 2 {
		huffMaxCnt := uint32(0)
		huffMax := uint8(0)
		for i, v := range hist[:16] {
			if v == 0 {
				continue
			}
			huffMax = byte(i)
			if v > huffMaxCnt {
				huffMaxCnt = v
			}
		}
		s.fse.HistogramFinished(huffMax, int(huffMaxCnt))
		s.fse.TableLog = maxFSETableLog
		b, err := fse.Compress(huffWeight[:maxSymbolValue], s.fse)
		if err == nil && len(b) < int(s.symbolLen>>1) {
			sz += 1 + len(b)
			return sz, nil
		}
		
	}
	
	if maxSymbolValue > (256 - 128) {
		
		return 0, ErrIncompressible
	}
	
	sz += 1 + int(maxSymbolValue/2)
	return sz, nil
}



func (c cTable) estimateSize(hist []uint32) int {
	nbBits := uint32(7)
	for i, v := range c[:len(hist)] {
		nbBits += uint32(v.nBits) * hist[i]
	}
	return int(nbBits >> 3)
}


func (s *Scratch) minSize(total int) int {
	nbBits := float64(7)
	fTotal := float64(total)
	for _, v := range s.count[:s.symbolLen] {
		n := float64(v)
		if n > 0 {
			nbBits += math.Log2(fTotal/n) * n
		}
	}
	return int(nbBits) >> 3
}

func highBit32(val uint32) (n uint32) {
	return uint32(bits.Len32(val) - 1)
}
