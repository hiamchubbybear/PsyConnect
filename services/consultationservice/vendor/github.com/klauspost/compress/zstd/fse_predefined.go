



package zstd

import (
	"fmt"
	"math"
	"sync"
)

var (
	
	
	
	fsePredef [3]fseDecoder

	
	
	
	fsePredefEnc [3]fseEncoder

	
	
	symbolTableX [3][]baseOffset

	
	
	maxTableSymbol = [3]uint8{tableLiteralLengths: maxLiteralLengthSymbol, tableOffsets: maxOffsetLengthSymbol, tableMatchLengths: maxMatchLengthSymbol}

	
	bitTables = [3][]byte{tableLiteralLengths: llBitsTable[:], tableOffsets: nil, tableMatchLengths: mlBitsTable[:]}
)

type tableIndex uint8

const (
	
	tableLiteralLengths tableIndex = 0
	tableOffsets        tableIndex = 1
	tableMatchLengths   tableIndex = 2

	maxLiteralLengthSymbol = 35
	maxOffsetLengthSymbol  = 30
	maxMatchLengthSymbol   = 52
)


type baseOffset struct {
	baseLine uint32
	addBits  uint8
}


func fillBase(dst []baseOffset, base uint32, bits ...uint8) {
	if len(bits) != len(dst) {
		panic(fmt.Sprintf("len(dst) (%d) != len(bits) (%d)", len(dst), len(bits)))
	}
	for i, bit := range bits {
		if base > math.MaxInt32 {
			panic("invalid decoding table, base overflows int32")
		}

		dst[i] = baseOffset{
			baseLine: base,
			addBits:  bit,
		}
		base += 1 << bit
	}
}

var predef sync.Once

func initPredefined() {
	predef.Do(func() {
		
		tmp := make([]baseOffset, 36)
		for i := range tmp[:16] {
			tmp[i] = baseOffset{
				baseLine: uint32(i),
				addBits:  0,
			}
		}
		fillBase(tmp[16:], 16, 1, 1, 1, 1, 2, 2, 3, 3, 4, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16)
		symbolTableX[tableLiteralLengths] = tmp

		
		tmp = make([]baseOffset, 53)
		for i := range tmp[:32] {
			tmp[i] = baseOffset{
				
				baseLine: uint32(i) + 3,
				addBits:  0,
			}
		}
		fillBase(tmp[32:], 35, 1, 1, 1, 1, 2, 2, 3, 3, 4, 4, 5, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16)
		symbolTableX[tableMatchLengths] = tmp

		
		tmp = make([]baseOffset, maxOffsetBits+1)
		tmp[1] = baseOffset{
			baseLine: 1,
			addBits:  1,
		}
		fillBase(tmp[2:], 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30)
		symbolTableX[tableOffsets] = tmp

		
		
		for i := range fsePredef[:] {
			f := &fsePredef[i]
			switch tableIndex(i) {
			case tableLiteralLengths:
				
				f.actualTableLog = 6
				copy(f.norm[:], []int16{4, 3, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 1, 1, 1,
					2, 2, 2, 2, 2, 2, 2, 2, 2, 3, 2, 1, 1, 1, 1, 1,
					-1, -1, -1, -1})
				f.symbolLen = 36
			case tableOffsets:
				
				f.actualTableLog = 5
				copy(f.norm[:], []int16{
					1, 1, 1, 1, 1, 1, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1,
					1, 1, 1, 1, 1, 1, 1, 1, -1, -1, -1, -1, -1})
				f.symbolLen = 29
			case tableMatchLengths:
				
				f.actualTableLog = 6
				copy(f.norm[:], []int16{
					1, 4, 3, 2, 2, 2, 2, 2, 2, 1, 1, 1, 1, 1, 1, 1,
					1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
					1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, -1, -1,
					-1, -1, -1, -1, -1})
				f.symbolLen = 53
			}
			if err := f.buildDtable(); err != nil {
				panic(fmt.Errorf("building table %v: %v", tableIndex(i), err))
			}
			if err := f.transform(symbolTableX[i]); err != nil {
				panic(fmt.Errorf("building table %v: %v", tableIndex(i), err))
			}
			f.preDefined = true

			
			enc := &fsePredefEnc[i]
			copy(enc.norm[:], f.norm[:])
			enc.symbolLen = f.symbolLen
			enc.actualTableLog = f.actualTableLog
			if err := enc.buildCTable(); err != nil {
				panic(fmt.Errorf("building encoding table %v: %v", tableIndex(i), err))
			}
			enc.setBits(bitTables[i])
			enc.preDefined = true
		}
	})
}
