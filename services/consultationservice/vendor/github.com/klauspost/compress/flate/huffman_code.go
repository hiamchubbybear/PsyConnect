



package flate

import (
	"math"
	"math/bits"
)

const (
	maxBitsLimit = 16
	
	literalCount = 286
)


type hcode uint32

func (h hcode) len() uint8 {
	return uint8(h)
}

func (h hcode) code64() uint64 {
	return uint64(h >> 8)
}

func (h hcode) zero() bool {
	return h == 0
}

type huffmanEncoder struct {
	codes    []hcode
	bitCount [17]int32

	
	
	
	freqcache [literalCount + 1]literalNode
}

type literalNode struct {
	literal uint16
	freq    uint16
}


type levelInfo struct {
	
	level int32

	
	lastFreq int32

	
	nextCharFreq int32

	
	
	nextPairFreq int32

	
	
	needed int32
}


func (h *hcode) set(code uint16, length uint8) {
	*h = hcode(length) | (hcode(code) << 8)
}

func newhcode(code uint16, length uint8) hcode {
	return hcode(length) | (hcode(code) << 8)
}

func reverseBits(number uint16, bitLength byte) uint16 {
	return bits.Reverse16(number << ((16 - bitLength) & 15))
}

func maxNode() literalNode { return literalNode{math.MaxUint16, math.MaxUint16} }

func newHuffmanEncoder(size int) *huffmanEncoder {
	
	c := uint(bits.Len32(uint32(size - 1)))
	return &huffmanEncoder{codes: make([]hcode, size, 1<<c)}
}


func generateFixedLiteralEncoding() *huffmanEncoder {
	h := newHuffmanEncoder(literalCount)
	codes := h.codes
	var ch uint16
	for ch = 0; ch < literalCount; ch++ {
		var bits uint16
		var size uint8
		switch {
		case ch < 144:
			
			bits = ch + 48
			size = 8
		case ch < 256:
			
			bits = ch + 400 - 144
			size = 9
		case ch < 280:
			
			bits = ch - 256
			size = 7
		default:
			
			bits = ch + 192 - 280
			size = 8
		}
		codes[ch] = newhcode(reverseBits(bits, size), size)
	}
	return h
}

func generateFixedOffsetEncoding() *huffmanEncoder {
	h := newHuffmanEncoder(30)
	codes := h.codes
	for ch := range codes {
		codes[ch] = newhcode(reverseBits(uint16(ch), 5), 5)
	}
	return h
}

var fixedLiteralEncoding = generateFixedLiteralEncoding()
var fixedOffsetEncoding = generateFixedOffsetEncoding()

func (h *huffmanEncoder) bitLength(freq []uint16) int {
	var total int
	for i, f := range freq {
		if f != 0 {
			total += int(f) * int(h.codes[i].len())
		}
	}
	return total
}

func (h *huffmanEncoder) bitLengthRaw(b []byte) int {
	var total int
	for _, f := range b {
		total += int(h.codes[f].len())
	}
	return total
}


func (h *huffmanEncoder) canReuseBits(freq []uint16) int {
	var total int
	for i, f := range freq {
		if f != 0 {
			code := h.codes[i]
			if code.zero() {
				return math.MaxInt32
			}
			total += int(f) * int(code.len())
		}
	}
	return total
}



















func (h *huffmanEncoder) bitCounts(list []literalNode, maxBits int32) []int32 {
	if maxBits >= maxBitsLimit {
		panic("flate: maxBits too large")
	}
	n := int32(len(list))
	list = list[0 : n+1]
	list[n] = maxNode()

	
	
	if maxBits > n-1 {
		maxBits = n - 1
	}

	
	
	
	
	var levels [maxBitsLimit]levelInfo
	
	
	
	
	var leafCounts [maxBitsLimit][maxBitsLimit]int32

	
	l2f := int32(list[2].freq)
	l1f := int32(list[1].freq)
	l0f := int32(list[0].freq) + int32(list[1].freq)

	for level := int32(1); level <= maxBits; level++ {
		
		
		levels[level] = levelInfo{
			level:        level,
			lastFreq:     l1f,
			nextCharFreq: l2f,
			nextPairFreq: l0f,
		}
		leafCounts[level][level] = 2
		if level == 1 {
			levels[level].nextPairFreq = math.MaxInt32
		}
	}

	
	levels[maxBits].needed = 2*n - 4

	level := uint32(maxBits)
	for level < 16 {
		l := &levels[level]
		if l.nextPairFreq == math.MaxInt32 && l.nextCharFreq == math.MaxInt32 {
			
			
			
			
			l.needed = 0
			levels[level+1].nextPairFreq = math.MaxInt32
			level++
			continue
		}

		prevFreq := l.lastFreq
		if l.nextCharFreq < l.nextPairFreq {
			
			n := leafCounts[level][level] + 1
			l.lastFreq = l.nextCharFreq
			
			leafCounts[level][level] = n
			e := list[n]
			if e.literal < math.MaxUint16 {
				l.nextCharFreq = int32(e.freq)
			} else {
				l.nextCharFreq = math.MaxInt32
			}
		} else {
			
			
			
			l.lastFreq = l.nextPairFreq
			
			if true {
				save := leafCounts[level][level]
				leafCounts[level] = leafCounts[level-1]
				leafCounts[level][level] = save
			} else {
				copy(leafCounts[level][:level], leafCounts[level-1][:level])
			}
			levels[l.level-1].needed = 2
		}

		if l.needed--; l.needed == 0 {
			
			
			
			
			if l.level == maxBits {
				
				break
			}
			levels[l.level+1].nextPairFreq = prevFreq + l.lastFreq
			level++
		} else {
			
			for levels[level-1].needed > 0 {
				level--
			}
		}
	}

	
	
	if leafCounts[maxBits][maxBits] != n {
		panic("leafCounts[maxBits][maxBits] != n")
	}

	bitCount := h.bitCount[:maxBits+1]
	bits := 1
	counts := &leafCounts[maxBits]
	for level := maxBits; level > 0; level-- {
		
		
		bitCount[bits] = counts[level] - counts[level-1]
		bits++
	}
	return bitCount
}



func (h *huffmanEncoder) assignEncodingAndSize(bitCount []int32, list []literalNode) {
	code := uint16(0)
	for n, bits := range bitCount {
		code <<= 1
		if n == 0 || bits == 0 {
			continue
		}
		
		
		
		
		chunk := list[len(list)-int(bits):]

		sortByLiteral(chunk)
		for _, node := range chunk {
			h.codes[node.literal] = newhcode(reverseBits(code, uint8(n)), uint8(n))
			code++
		}
		list = list[0 : len(list)-int(bits)]
	}
}





func (h *huffmanEncoder) generate(freq []uint16, maxBits int32) {
	list := h.freqcache[:len(freq)+1]
	codes := h.codes[:len(freq)]
	
	count := 0
	
	for i, f := range freq {
		if f != 0 {
			list[count] = literalNode{uint16(i), f}
			count++
		} else {
			codes[i] = 0
		}
	}
	list[count] = literalNode{}

	list = list[:count]
	if count <= 2 {
		
		
		for i, node := range list {
			
			h.codes[node.literal].set(uint16(i), 1)
		}
		return
	}
	sortByFreq(list)

	
	bitCount := h.bitCounts(list, maxBits)
	
	h.assignEncodingAndSize(bitCount, list)
}


func atLeastOne(v float32) float32 {
	if v < 1 {
		return 1
	}
	if v > 15 {
		return 15
	}
	return v
}

func histogram(b []byte, h []uint16) {
	if true && len(b) >= 8<<10 {
		
		histogramSplit(b, h)
	} else {
		h = h[:256]
		for _, t := range b {
			h[t]++
		}
	}
}

func histogramSplit(b []byte, h []uint16) {
	
	
	h = h[:256]
	for len(b)&3 != 0 {
		h[b[0]]++
		b = b[1:]
	}
	n := len(b) / 4
	x, y, z, w := b[:n], b[n:], b[n+n:], b[n+n+n:]
	y, z, w = y[:len(x)], z[:len(x)], w[:len(x)]
	for i, t := range x {
		v0 := &h[t]
		v1 := &h[y[i]]
		v3 := &h[w[i]]
		v2 := &h[z[i]]
		*v0++
		*v1++
		*v2++
		*v3++
	}
}
