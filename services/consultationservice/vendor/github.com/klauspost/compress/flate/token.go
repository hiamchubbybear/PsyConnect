



package flate

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

const (
	
	
	
	
	lengthShift         = 22
	offsetMask          = 1<<lengthShift - 1
	typeMask            = 3 << 30
	literalType         = 0 << 30
	matchType           = 1 << 30
	matchOffsetOnlyMask = 0xffff
)



var lengthCodes = [256]uint8{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 8,
	9, 9, 10, 10, 11, 11, 12, 12, 12, 12,
	13, 13, 13, 13, 14, 14, 14, 14, 15, 15,
	15, 15, 16, 16, 16, 16, 16, 16, 16, 16,
	17, 17, 17, 17, 17, 17, 17, 17, 18, 18,
	18, 18, 18, 18, 18, 18, 19, 19, 19, 19,
	19, 19, 19, 19, 20, 20, 20, 20, 20, 20,
	20, 20, 20, 20, 20, 20, 20, 20, 20, 20,
	21, 21, 21, 21, 21, 21, 21, 21, 21, 21,
	21, 21, 21, 21, 21, 21, 22, 22, 22, 22,
	22, 22, 22, 22, 22, 22, 22, 22, 22, 22,
	22, 22, 23, 23, 23, 23, 23, 23, 23, 23,
	23, 23, 23, 23, 23, 23, 23, 23, 24, 24,
	24, 24, 24, 24, 24, 24, 24, 24, 24, 24,
	24, 24, 24, 24, 24, 24, 24, 24, 24, 24,
	24, 24, 24, 24, 24, 24, 24, 24, 24, 24,
	25, 25, 25, 25, 25, 25, 25, 25, 25, 25,
	25, 25, 25, 25, 25, 25, 25, 25, 25, 25,
	25, 25, 25, 25, 25, 25, 25, 25, 25, 25,
	25, 25, 26, 26, 26, 26, 26, 26, 26, 26,
	26, 26, 26, 26, 26, 26, 26, 26, 26, 26,
	26, 26, 26, 26, 26, 26, 26, 26, 26, 26,
	26, 26, 26, 26, 27, 27, 27, 27, 27, 27,
	27, 27, 27, 27, 27, 27, 27, 27, 27, 27,
	27, 27, 27, 27, 27, 27, 27, 27, 27, 27,
	27, 27, 27, 27, 27, 28,
}


var lengthCodes1 = [256]uint8{
	1, 2, 3, 4, 5, 6, 7, 8, 9, 9,
	10, 10, 11, 11, 12, 12, 13, 13, 13, 13,
	14, 14, 14, 14, 15, 15, 15, 15, 16, 16,
	16, 16, 17, 17, 17, 17, 17, 17, 17, 17,
	18, 18, 18, 18, 18, 18, 18, 18, 19, 19,
	19, 19, 19, 19, 19, 19, 20, 20, 20, 20,
	20, 20, 20, 20, 21, 21, 21, 21, 21, 21,
	21, 21, 21, 21, 21, 21, 21, 21, 21, 21,
	22, 22, 22, 22, 22, 22, 22, 22, 22, 22,
	22, 22, 22, 22, 22, 22, 23, 23, 23, 23,
	23, 23, 23, 23, 23, 23, 23, 23, 23, 23,
	23, 23, 24, 24, 24, 24, 24, 24, 24, 24,
	24, 24, 24, 24, 24, 24, 24, 24, 25, 25,
	25, 25, 25, 25, 25, 25, 25, 25, 25, 25,
	25, 25, 25, 25, 25, 25, 25, 25, 25, 25,
	25, 25, 25, 25, 25, 25, 25, 25, 25, 25,
	26, 26, 26, 26, 26, 26, 26, 26, 26, 26,
	26, 26, 26, 26, 26, 26, 26, 26, 26, 26,
	26, 26, 26, 26, 26, 26, 26, 26, 26, 26,
	26, 26, 27, 27, 27, 27, 27, 27, 27, 27,
	27, 27, 27, 27, 27, 27, 27, 27, 27, 27,
	27, 27, 27, 27, 27, 27, 27, 27, 27, 27,
	27, 27, 27, 27, 28, 28, 28, 28, 28, 28,
	28, 28, 28, 28, 28, 28, 28, 28, 28, 28,
	28, 28, 28, 28, 28, 28, 28, 28, 28, 28,
	28, 28, 28, 28, 28, 29,
}

var offsetCodes = [256]uint32{
	0, 1, 2, 3, 4, 4, 5, 5, 6, 6, 6, 6, 7, 7, 7, 7,
	8, 8, 8, 8, 8, 8, 8, 8, 9, 9, 9, 9, 9, 9, 9, 9,
	10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10,
	11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11, 11,
	12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12,
	12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12,
	13, 13, 13, 13, 13, 13, 13, 13, 13, 13, 13, 13, 13, 13, 13, 13,
	13, 13, 13, 13, 13, 13, 13, 13, 13, 13, 13, 13, 13, 13, 13, 13,
	14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14,
	14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14,
	14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14,
	14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14, 14,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
}


var offsetCodes14 = [256]uint32{
	14, 15, 16, 17, 18, 18, 19, 19, 20, 20, 20, 20, 21, 21, 21, 21,
	22, 22, 22, 22, 22, 22, 22, 22, 23, 23, 23, 23, 23, 23, 23, 23,
	24, 24, 24, 24, 24, 24, 24, 24, 24, 24, 24, 24, 24, 24, 24, 24,
	25, 25, 25, 25, 25, 25, 25, 25, 25, 25, 25, 25, 25, 25, 25, 25,
	26, 26, 26, 26, 26, 26, 26, 26, 26, 26, 26, 26, 26, 26, 26, 26,
	26, 26, 26, 26, 26, 26, 26, 26, 26, 26, 26, 26, 26, 26, 26, 26,
	27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27,
	27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27,
	28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28,
	28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28,
	28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28,
	28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28, 28,
	29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29,
	29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29,
	29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29,
	29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29,
}

type token uint32

type tokens struct {
	extraHist [32]uint16  
	offHist   [32]uint16  
	litHist   [256]uint16 
	nFilled   int
	n         uint16 
	tokens    [maxStoreBlockSize + 1]token
}

func (t *tokens) Reset() {
	if t.n == 0 {
		return
	}
	t.n = 0
	t.nFilled = 0
	for i := range t.litHist[:] {
		t.litHist[i] = 0
	}
	for i := range t.extraHist[:] {
		t.extraHist[i] = 0
	}
	for i := range t.offHist[:] {
		t.offHist[i] = 0
	}
}

func (t *tokens) Fill() {
	if t.n == 0 {
		return
	}
	for i, v := range t.litHist[:] {
		if v == 0 {
			t.litHist[i] = 1
			t.nFilled++
		}
	}
	for i, v := range t.extraHist[:literalCount-256] {
		if v == 0 {
			t.nFilled++
			t.extraHist[i] = 1
		}
	}
	for i, v := range t.offHist[:offsetCodeCount] {
		if v == 0 {
			t.offHist[i] = 1
		}
	}
}

func indexTokens(in []token) tokens {
	var t tokens
	t.indexTokens(in)
	return t
}

func (t *tokens) indexTokens(in []token) {
	t.Reset()
	for _, tok := range in {
		if tok < matchType {
			t.AddLiteral(tok.literal())
			continue
		}
		t.AddMatch(uint32(tok.length()), tok.offset()&matchOffsetOnlyMask)
	}
}


func emitLiteral(dst *tokens, lit []byte) {
	for _, v := range lit {
		dst.tokens[dst.n] = token(v)
		dst.litHist[v]++
		dst.n++
	}
}

func (t *tokens) AddLiteral(lit byte) {
	t.tokens[t.n] = token(lit)
	t.litHist[lit]++
	t.n++
}


func mFastLog2(val float32) float32 {
	ux := int32(math.Float32bits(val))
	log2 := (float32)(((ux >> 23) & 255) - 128)
	ux &= -0x7f800001
	ux += 127 << 23
	uval := math.Float32frombits(uint32(ux))
	log2 += ((-0.34484843)*uval+2.02466578)*uval - 0.67487759
	return log2
}




func (t *tokens) EstimatedBits() int {
	shannon := float32(0)
	bits := int(0)
	nMatches := 0
	total := int(t.n) + t.nFilled
	if total > 0 {
		invTotal := 1.0 / float32(total)
		for _, v := range t.litHist[:] {
			if v > 0 {
				n := float32(v)
				shannon += atLeastOne(-mFastLog2(n*invTotal)) * n
			}
		}
		
		shannon += 15
		for i, v := range t.extraHist[1 : literalCount-256] {
			if v > 0 {
				n := float32(v)
				shannon += atLeastOne(-mFastLog2(n*invTotal)) * n
				bits += int(lengthExtraBits[i&31]) * int(v)
				nMatches += int(v)
			}
		}
	}
	if nMatches > 0 {
		invTotal := 1.0 / float32(nMatches)
		for i, v := range t.offHist[:offsetCodeCount] {
			if v > 0 {
				n := float32(v)
				shannon += atLeastOne(-mFastLog2(n*invTotal)) * n
				bits += int(offsetExtraBits[i&31]) * int(v)
			}
		}
	}
	return int(shannon) + bits
}



func (t *tokens) AddMatch(xlength uint32, xoffset uint32) {
	if debugDeflate {
		if xlength >= maxMatchLength+baseMatchLength {
			panic(fmt.Errorf("invalid length: %v", xlength))
		}
		if xoffset >= maxMatchOffset+baseMatchOffset {
			panic(fmt.Errorf("invalid offset: %v", xoffset))
		}
	}
	oCode := offsetCode(xoffset)
	xoffset |= oCode << 16

	t.extraHist[lengthCodes1[uint8(xlength)]]++
	t.offHist[oCode&31]++
	t.tokens[t.n] = token(matchType | xlength<<lengthShift | xoffset)
	t.n++
}



func (t *tokens) AddMatchLong(xlength int32, xoffset uint32) {
	if debugDeflate {
		if xoffset >= maxMatchOffset+baseMatchOffset {
			panic(fmt.Errorf("invalid offset: %v", xoffset))
		}
	}
	oc := offsetCode(xoffset)
	xoffset |= oc << 16
	for xlength > 0 {
		xl := xlength
		if xl > 258 {
			
			if xl > 258+baseMatchLength {
				xl = 258
			} else {
				xl = 258 - baseMatchLength
			}
		}
		xlength -= xl
		xl -= baseMatchLength
		t.extraHist[lengthCodes1[uint8(xl)]]++
		t.offHist[oc&31]++
		t.tokens[t.n] = token(matchType | uint32(xl)<<lengthShift | xoffset)
		t.n++
	}
}

func (t *tokens) AddEOB() {
	t.tokens[t.n] = token(endBlockMarker)
	t.extraHist[0]++
	t.n++
}

func (t *tokens) Slice() []token {
	return t.tokens[:t.n]
}


func (t *tokens) VarInt() []byte {
	var b = make([]byte, binary.MaxVarintLen32*int(t.n))
	var off int
	for _, v := range t.tokens[:t.n] {
		off += binary.PutUvarint(b[off:], uint64(v))
	}
	return b[:off]
}



func (t *tokens) FromVarInt(b []byte) error {
	var buf = bytes.NewReader(b)
	var toks []token
	for {
		r, err := binary.ReadUvarint(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		toks = append(toks, token(r))
	}
	t.indexTokens(toks)
	return nil
}


func (t token) typ() uint32 { return uint32(t) & typeMask }


func (t token) literal() uint8 { return uint8(t) }


func (t token) offset() uint32 { return uint32(t) & offsetMask }

func (t token) length() uint8 { return uint8(t >> lengthShift) }


func lengthCode(len uint8) uint8 { return lengthCodes[len] }


func offsetCode(off uint32) uint32 {
	if false {
		if off < uint32(len(offsetCodes)) {
			return offsetCodes[off&255]
		} else if off>>7 < uint32(len(offsetCodes)) {
			return offsetCodes[(off>>7)&255] + 14
		} else {
			return offsetCodes[(off>>14)&255] + 28
		}
	}
	if off < uint32(len(offsetCodes)) {
		return offsetCodes[uint8(off)]
	}
	return offsetCodes14[uint8(off>>7)]
}
