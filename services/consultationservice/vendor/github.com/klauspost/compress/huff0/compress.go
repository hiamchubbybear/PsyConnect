package huff0

import (
	"fmt"
	"math"
	"runtime"
	"sync"
)





func Compress1X(in []byte, s *Scratch) (out []byte, reUsed bool, err error) {
	s, err = s.prepare(in)
	if err != nil {
		return nil, false, err
	}
	return compress(in, s, s.compress1X)
}






func Compress4X(in []byte, s *Scratch) (out []byte, reUsed bool, err error) {
	s, err = s.prepare(in)
	if err != nil {
		return nil, false, err
	}
	if false {
		
		const parallelThreshold = 8 << 10
		if len(in) < parallelThreshold || runtime.GOMAXPROCS(0) == 1 {
			return compress(in, s, s.compress4X)
		}
		return compress(in, s, s.compress4Xp)
	}
	return compress(in, s, s.compress4X)
}

func compress(in []byte, s *Scratch, compressor func(src []byte) ([]byte, error)) (out []byte, reUsed bool, err error) {
	
	if s.Reuse == ReusePolicyNone {
		s.prevTable = s.prevTable[:0]
	}

	
	maxCount := s.maxCount
	var canReuse = false
	if maxCount == 0 {
		maxCount, canReuse = s.countSimple(in)
	} else {
		canReuse = s.canUseTable(s.prevTable)
	}

	
	wantSize := len(in)
	if s.WantLogLess > 0 {
		wantSize -= wantSize >> s.WantLogLess
	}

	
	s.clearCount = true
	s.maxCount = 0
	if maxCount >= len(in) {
		if maxCount > len(in) {
			return nil, false, fmt.Errorf("maxCount (%d) > length (%d)", maxCount, len(in))
		}
		if len(in) == 1 {
			return nil, false, ErrIncompressible
		}
		
		return nil, false, ErrUseRLE
	}
	if maxCount == 1 || maxCount < (len(in)>>7) {
		
		return nil, false, ErrIncompressible
	}
	if s.Reuse == ReusePolicyMust && !canReuse {
		
		return nil, false, ErrIncompressible
	}
	if (s.Reuse == ReusePolicyPrefer || s.Reuse == ReusePolicyMust) && canReuse {
		keepTable := s.cTable
		keepTL := s.actualTableLog
		s.cTable = s.prevTable
		s.actualTableLog = s.prevTableLog
		s.Out, err = compressor(in)
		s.cTable = keepTable
		s.actualTableLog = keepTL
		if err == nil && len(s.Out) < wantSize {
			s.OutData = s.Out
			return s.Out, true, nil
		}
		if s.Reuse == ReusePolicyMust {
			return nil, false, ErrIncompressible
		}
		
		s.prevTable = s.prevTable[:0]
	}

	
	err = s.buildCTable()
	if err != nil {
		return nil, false, err
	}

	if false && !s.canUseTable(s.cTable) {
		panic("invalid table generated")
	}

	if s.Reuse == ReusePolicyAllow && canReuse {
		hSize := len(s.Out)
		oldSize := s.prevTable.estimateSize(s.count[:s.symbolLen])
		newSize := s.cTable.estimateSize(s.count[:s.symbolLen])
		if oldSize <= hSize+newSize || hSize+12 >= wantSize {
			
			keepTable := s.cTable
			keepTL := s.actualTableLog

			s.cTable = s.prevTable
			s.actualTableLog = s.prevTableLog
			s.Out, err = compressor(in)

			
			s.cTable = keepTable
			s.actualTableLog = keepTL
			if err != nil {
				return nil, false, err
			}
			if len(s.Out) >= wantSize {
				return nil, false, ErrIncompressible
			}
			s.OutData = s.Out
			return s.Out, true, nil
		}
	}

	
	err = s.cTable.write(s)
	if err != nil {
		s.OutTable = nil
		return nil, false, err
	}
	s.OutTable = s.Out

	
	s.Out, err = compressor(in)
	if err != nil {
		s.OutTable = nil
		return nil, false, err
	}
	if len(s.Out) >= wantSize {
		s.OutTable = nil
		return nil, false, ErrIncompressible
	}
	
	s.prevTable, s.prevTableLog, s.cTable = s.cTable, s.actualTableLog, s.prevTable[:0]
	s.OutData = s.Out[len(s.OutTable):]
	return s.Out, false, nil
}


func EstimateSizes(in []byte, s *Scratch) (tableSz, dataSz, reuseSz int, err error) {
	s, err = s.prepare(in)
	if err != nil {
		return 0, 0, 0, err
	}

	
	tableSz, dataSz, reuseSz = -1, -1, -1
	maxCount := s.maxCount
	var canReuse = false
	if maxCount == 0 {
		maxCount, canReuse = s.countSimple(in)
	} else {
		canReuse = s.canUseTable(s.prevTable)
	}

	
	wantSize := len(in)
	if s.WantLogLess > 0 {
		wantSize -= wantSize >> s.WantLogLess
	}

	
	s.clearCount = true
	s.maxCount = 0
	if maxCount >= len(in) {
		if maxCount > len(in) {
			return 0, 0, 0, fmt.Errorf("maxCount (%d) > length (%d)", maxCount, len(in))
		}
		if len(in) == 1 {
			return 0, 0, 0, ErrIncompressible
		}
		
		return 0, 0, 0, ErrUseRLE
	}
	if maxCount == 1 || maxCount < (len(in)>>7) {
		
		return 0, 0, 0, ErrIncompressible
	}

	
	err = s.buildCTable()
	if err != nil {
		return 0, 0, 0, err
	}

	if false && !s.canUseTable(s.cTable) {
		panic("invalid table generated")
	}

	tableSz, err = s.cTable.estTableSize(s)
	if err != nil {
		return 0, 0, 0, err
	}
	if canReuse {
		reuseSz = s.prevTable.estimateSize(s.count[:s.symbolLen])
	}
	dataSz = s.cTable.estimateSize(s.count[:s.symbolLen])

	
	return tableSz, dataSz, reuseSz, nil
}

func (s *Scratch) compress1X(src []byte) ([]byte, error) {
	return s.compress1xDo(s.Out, src)
}

func (s *Scratch) compress1xDo(dst, src []byte) ([]byte, error) {
	var bw = bitWriter{out: dst}

	
	n := len(src)
	n -= n & 3
	cTable := s.cTable[:256]

	
	for i := len(src) & 3; i > 0; i-- {
		bw.encSymbol(cTable, src[n+i-1])
	}
	n -= 4
	if s.actualTableLog <= 8 {
		for ; n >= 0; n -= 4 {
			tmp := src[n : n+4]
			
			bw.flush32()
			bw.encFourSymbols(cTable[tmp[3]], cTable[tmp[2]], cTable[tmp[1]], cTable[tmp[0]])
		}
	} else {
		for ; n >= 0; n -= 4 {
			tmp := src[n : n+4]
			
			bw.flush32()
			bw.encTwoSymbols(cTable, tmp[3], tmp[2])
			bw.flush32()
			bw.encTwoSymbols(cTable, tmp[1], tmp[0])
		}
	}
	err := bw.close()
	return bw.out, err
}

var sixZeros [6]byte

func (s *Scratch) compress4X(src []byte) ([]byte, error) {
	if len(src) < 12 {
		return nil, ErrIncompressible
	}
	segmentSize := (len(src) + 3) / 4

	
	offsetIdx := len(s.Out)
	s.Out = append(s.Out, sixZeros[:]...)

	for i := 0; i < 4; i++ {
		toDo := src
		if len(toDo) > segmentSize {
			toDo = toDo[:segmentSize]
		}
		src = src[len(toDo):]

		var err error
		idx := len(s.Out)
		s.Out, err = s.compress1xDo(s.Out, toDo)
		if err != nil {
			return nil, err
		}
		if len(s.Out)-idx > math.MaxUint16 {
			
			return nil, ErrIncompressible
		}
		
		if i < 3 {
			
			length := len(s.Out) - idx
			s.Out[i*2+offsetIdx] = byte(length)
			s.Out[i*2+offsetIdx+1] = byte(length >> 8)
		}
	}

	return s.Out, nil
}


func (s *Scratch) compress4Xp(src []byte) ([]byte, error) {
	if len(src) < 12 {
		return nil, ErrIncompressible
	}
	
	s.Out = s.Out[:6]

	segmentSize := (len(src) + 3) / 4
	var wg sync.WaitGroup
	var errs [4]error
	wg.Add(4)
	for i := 0; i < 4; i++ {
		toDo := src
		if len(toDo) > segmentSize {
			toDo = toDo[:segmentSize]
		}
		src = src[len(toDo):]

		
		go func(i int) {
			s.tmpOut[i], errs[i] = s.compress1xDo(s.tmpOut[i][:0], toDo)
			wg.Done()
		}(i)
	}
	wg.Wait()
	for i := 0; i < 4; i++ {
		if errs[i] != nil {
			return nil, errs[i]
		}
		o := s.tmpOut[i]
		if len(o) > math.MaxUint16 {
			
			return nil, ErrIncompressible
		}
		
		if i < 3 {
			
			s.Out[i*2] = byte(len(o))
			s.Out[i*2+1] = byte(len(o) >> 8)
		}

		
		s.Out = append(s.Out, o...)
	}
	return s.Out, nil
}




func (s *Scratch) countSimple(in []byte) (max int, reuse bool) {
	reuse = true
	for _, v := range in {
		s.count[v]++
	}
	m := uint32(0)
	if len(s.prevTable) > 0 {
		for i, v := range s.count[:] {
			if v == 0 {
				continue
			}
			if v > m {
				m = v
			}
			s.symbolLen = uint16(i) + 1
			if i >= len(s.prevTable) {
				reuse = false
			} else if s.prevTable[i].nBits == 0 {
				reuse = false
			}
		}
		return int(m), reuse
	}
	for i, v := range s.count[:] {
		if v == 0 {
			continue
		}
		if v > m {
			m = v
		}
		s.symbolLen = uint16(i) + 1
	}
	return int(m), false
}

func (s *Scratch) canUseTable(c cTable) bool {
	if len(c) < int(s.symbolLen) {
		return false
	}
	for i, v := range s.count[:s.symbolLen] {
		if v != 0 && c[i].nBits == 0 {
			return false
		}
	}
	return true
}


func (s *Scratch) validateTable(c cTable) bool {
	if len(c) < int(s.symbolLen) {
		return false
	}
	for i, v := range s.count[:s.symbolLen] {
		if v != 0 {
			if c[i].nBits == 0 {
				return false
			}
			if c[i].nBits > s.actualTableLog {
				return false
			}
		}
	}
	return true
}


func (s *Scratch) minTableLog() uint8 {
	minBitsSrc := highBit32(uint32(s.br.remain())) + 1
	minBitsSymbols := highBit32(uint32(s.symbolLen-1)) + 2
	if minBitsSrc < minBitsSymbols {
		return uint8(minBitsSrc)
	}
	return uint8(minBitsSymbols)
}


func (s *Scratch) optimalTableLog() {
	tableLog := s.TableLog
	minBits := s.minTableLog()
	maxBitsSrc := uint8(highBit32(uint32(s.br.remain()-1))) - 1
	if maxBitsSrc < tableLog {
		
		tableLog = maxBitsSrc
	}
	if minBits > tableLog {
		tableLog = minBits
	}
	
	if tableLog < minTablelog {
		tableLog = minTablelog
	}
	if tableLog > tableLogMax {
		tableLog = tableLogMax
	}
	s.actualTableLog = tableLog
}

type cTableEntry struct {
	val   uint16
	nBits uint8
	
}

const huffNodesMask = huffNodesLen - 1

func (s *Scratch) buildCTable() error {
	s.optimalTableLog()
	s.huffSort()
	if cap(s.cTable) < maxSymbolValue+1 {
		s.cTable = make([]cTableEntry, s.symbolLen, maxSymbolValue+1)
	} else {
		s.cTable = s.cTable[:s.symbolLen]
		for i := range s.cTable {
			s.cTable[i] = cTableEntry{}
		}
	}

	var startNode = int16(s.symbolLen)
	nonNullRank := s.symbolLen - 1

	nodeNb := startNode
	huffNode := s.nodes[1 : huffNodesLen+1]

	
	
	huffNode0 := s.nodes[0 : huffNodesLen+1]

	for huffNode[nonNullRank].count() == 0 {
		nonNullRank--
	}

	lowS := int16(nonNullRank)
	nodeRoot := nodeNb + lowS - 1
	lowN := nodeNb
	huffNode[nodeNb].setCount(huffNode[lowS].count() + huffNode[lowS-1].count())
	huffNode[lowS].setParent(nodeNb)
	huffNode[lowS-1].setParent(nodeNb)
	nodeNb++
	lowS -= 2
	for n := nodeNb; n <= nodeRoot; n++ {
		huffNode[n].setCount(1 << 30)
	}
	
	huffNode0[0].setCount(1 << 31)

	
	for nodeNb <= nodeRoot {
		var n1, n2 int16
		if huffNode0[lowS+1].count() < huffNode0[lowN+1].count() {
			n1 = lowS
			lowS--
		} else {
			n1 = lowN
			lowN++
		}
		if huffNode0[lowS+1].count() < huffNode0[lowN+1].count() {
			n2 = lowS
			lowS--
		} else {
			n2 = lowN
			lowN++
		}

		huffNode[nodeNb].setCount(huffNode0[n1+1].count() + huffNode0[n2+1].count())
		huffNode0[n1+1].setParent(nodeNb)
		huffNode0[n2+1].setParent(nodeNb)
		nodeNb++
	}

	
	huffNode[nodeRoot].setNbBits(0)
	for n := nodeRoot - 1; n >= startNode; n-- {
		huffNode[n].setNbBits(huffNode[huffNode[n].parent()].nbBits() + 1)
	}
	for n := uint16(0); n <= nonNullRank; n++ {
		huffNode[n].setNbBits(huffNode[huffNode[n].parent()].nbBits() + 1)
	}
	s.actualTableLog = s.setMaxHeight(int(nonNullRank))
	maxNbBits := s.actualTableLog

	
	if maxNbBits > tableLogMax {
		return fmt.Errorf("internal error: maxNbBits (%d) > tableLogMax (%d)", maxNbBits, tableLogMax)
	}
	var nbPerRank [tableLogMax + 1]uint16
	var valPerRank [16]uint16
	for _, v := range huffNode[:nonNullRank+1] {
		nbPerRank[v.nbBits()]++
	}
	
	{
		min := uint16(0)
		for n := maxNbBits; n > 0; n-- {
			
			valPerRank[n] = min
			min += nbPerRank[n]
			min >>= 1
		}
	}

	
	for _, v := range huffNode[:nonNullRank+1] {
		s.cTable[v.symbol()].nBits = v.nbBits()
	}

	
	t := s.cTable[:s.symbolLen]
	for n, val := range t {
		nbits := val.nBits & 15
		v := valPerRank[nbits]
		t[n].val = v
		valPerRank[nbits] = v + 1
	}

	return nil
}


func (s *Scratch) huffSort() {
	type rankPos struct {
		base    uint32
		current uint32
	}

	
	nodes := s.nodes[:huffNodesLen+1]
	s.nodes = nodes
	nodes = nodes[1 : huffNodesLen+1]

	
	var rank [32]rankPos
	for _, v := range s.count[:s.symbolLen] {
		r := highBit32(v+1) & 31
		rank[r].base++
	}
	
	const maxBitLength = 18 + 1
	for n := maxBitLength; n > 0; n-- {
		rank[n-1].base += rank[n].base
	}
	for n := range rank[:maxBitLength] {
		rank[n].current = rank[n].base
	}
	for n, c := range s.count[:s.symbolLen] {
		r := (highBit32(c+1) + 1) & 31
		pos := rank[r].current
		rank[r].current++
		prev := nodes[(pos-1)&huffNodesMask]
		for pos > rank[r].base && c > prev.count() {
			nodes[pos&huffNodesMask] = prev
			pos--
			prev = nodes[(pos-1)&huffNodesMask]
		}
		nodes[pos&huffNodesMask] = makeNodeElt(c, byte(n))
	}
}

func (s *Scratch) setMaxHeight(lastNonNull int) uint8 {
	maxNbBits := s.actualTableLog
	huffNode := s.nodes[1 : huffNodesLen+1]
	

	largestBits := huffNode[lastNonNull].nbBits()

	
	if largestBits <= maxNbBits {
		return largestBits
	}
	totalCost := int(0)
	baseCost := int(1) << (largestBits - maxNbBits)
	n := uint32(lastNonNull)

	for huffNode[n].nbBits() > maxNbBits {
		totalCost += baseCost - (1 << (largestBits - huffNode[n].nbBits()))
		huffNode[n].setNbBits(maxNbBits)
		n--
	}
	

	for huffNode[n].nbBits() == maxNbBits {
		n--
	}
	

	
	totalCost >>= largestBits - maxNbBits 

	
	{
		const noSymbol = 0xF0F0F0F0
		var rankLast [tableLogMax + 2]uint32

		for i := range rankLast[:] {
			rankLast[i] = noSymbol
		}

		
		{
			currentNbBits := maxNbBits
			for pos := int(n); pos >= 0; pos-- {
				if huffNode[pos].nbBits() >= currentNbBits {
					continue
				}
				currentNbBits = huffNode[pos].nbBits() 
				rankLast[maxNbBits-currentNbBits] = uint32(pos)
			}
		}

		for totalCost > 0 {
			nBitsToDecrease := uint8(highBit32(uint32(totalCost))) + 1

			for ; nBitsToDecrease > 1; nBitsToDecrease-- {
				highPos := rankLast[nBitsToDecrease]
				lowPos := rankLast[nBitsToDecrease-1]
				if highPos == noSymbol {
					continue
				}
				if lowPos == noSymbol {
					break
				}
				highTotal := huffNode[highPos].count()
				lowTotal := 2 * huffNode[lowPos].count()
				if highTotal <= lowTotal {
					break
				}
			}
			
			
			
			for (nBitsToDecrease <= tableLogMax) && (rankLast[nBitsToDecrease] == noSymbol) {
				nBitsToDecrease++
			}
			totalCost -= 1 << (nBitsToDecrease - 1)
			if rankLast[nBitsToDecrease-1] == noSymbol {
				
				rankLast[nBitsToDecrease-1] = rankLast[nBitsToDecrease]
			}
			huffNode[rankLast[nBitsToDecrease]].setNbBits(1 +
				huffNode[rankLast[nBitsToDecrease]].nbBits())
			if rankLast[nBitsToDecrease] == 0 {
				
				rankLast[nBitsToDecrease] = noSymbol
			} else {
				rankLast[nBitsToDecrease]--
				if huffNode[rankLast[nBitsToDecrease]].nbBits() != maxNbBits-nBitsToDecrease {
					rankLast[nBitsToDecrease] = noSymbol 
				}
			}
		}

		for totalCost < 0 { 
			if rankLast[1] == noSymbol { 
				for huffNode[n].nbBits() == maxNbBits {
					n--
				}
				huffNode[n+1].setNbBits(huffNode[n+1].nbBits() - 1)
				rankLast[1] = n + 1
				totalCost++
				continue
			}
			huffNode[rankLast[1]+1].setNbBits(huffNode[rankLast[1]+1].nbBits() - 1)
			rankLast[1]++
			totalCost++
		}
	}
	return maxNbBits
}










type nodeElt uint64

func makeNodeElt(count uint32, symbol byte) nodeElt {
	return nodeElt(count) | nodeElt(symbol)<<48
}

func (e *nodeElt) count() uint32  { return uint32(*e) }
func (e *nodeElt) parent() uint16 { return uint16(*e >> 32) }
func (e *nodeElt) symbol() byte   { return byte(*e >> 48) }
func (e *nodeElt) nbBits() uint8  { return uint8(*e >> 56) }

func (e *nodeElt) setCount(c uint32) { *e = (*e)&0xffffffff00000000 | nodeElt(c) }
func (e *nodeElt) setParent(p int16) { *e = (*e)&0xffff0000ffffffff | nodeElt(uint16(p))<<32 }
func (e *nodeElt) setNbBits(n uint8) { *e = (*e)&0x00ffffffffffffff | nodeElt(n)<<56 }
