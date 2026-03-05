



package flate

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

const (
	
	offsetCodeCount = 30

	
	endBlockMarker = 256

	
	lengthCodesStart = 257

	
	codegenCodeCount = 19
	badCode          = 255

	
	
	maxPredefinedTokens = 250

	
	
	
	
	bufferFlushSize = 246
)


const lengthExtraBitsMinCode = 8


var lengthExtraBits = [32]uint8{
	 0, 0, 0,
	 0, 0, 0, 0, 0, 1, 1, 1, 1, 2,
	 2, 2, 2, 3, 3, 3, 3, 4, 4, 4,
	 4, 5, 5, 5, 5, 0,
}


var lengthBase = [32]uint8{
	0, 1, 2, 3, 4, 5, 6, 7, 8, 10,
	12, 14, 16, 20, 24, 28, 32, 40, 48, 56,
	64, 80, 96, 112, 128, 160, 192, 224, 255,
}


const offsetExtraBitsMinCode = 4


var offsetExtraBits = [32]int8{
	0, 0, 0, 0, 1, 1, 2, 2, 3, 3,
	4, 4, 5, 5, 6, 6, 7, 7, 8, 8,
	9, 9, 10, 10, 11, 11, 12, 12, 13, 13,
	
	14, 14,
}

var offsetCombined = [32]uint32{}

func init() {
	var offsetBase = [32]uint32{
		
		0x000000, 0x000001, 0x000002, 0x000003, 0x000004,
		0x000006, 0x000008, 0x00000c, 0x000010, 0x000018,
		0x000020, 0x000030, 0x000040, 0x000060, 0x000080,
		0x0000c0, 0x000100, 0x000180, 0x000200, 0x000300,
		0x000400, 0x000600, 0x000800, 0x000c00, 0x001000,
		0x001800, 0x002000, 0x003000, 0x004000, 0x006000,

		
		0x008000, 0x00c000,
	}

	for i := range offsetCombined[:] {
		
		if offsetExtraBits[i] == 0 || offsetBase[i] > 0x006000 {
			continue
		}
		offsetCombined[i] = uint32(offsetExtraBits[i]) | (offsetBase[i] << 8)
	}
}


var codegenOrder = []uint32{16, 17, 18, 0, 8, 7, 9, 6, 10, 5, 11, 4, 12, 3, 13, 2, 14, 1, 15}

type huffmanBitWriter struct {
	
	
	
	writer io.Writer

	
	
	bits            uint64
	nbits           uint8
	nbytes          uint8
	lastHuffMan     bool
	literalEncoding *huffmanEncoder
	tmpLitEncoding  *huffmanEncoder
	offsetEncoding  *huffmanEncoder
	codegenEncoding *huffmanEncoder
	err             error
	lastHeader      int
	
	logNewTablePenalty uint
	bytes              [256 + 8]byte
	literalFreq        [lengthCodesStart + 32]uint16
	offsetFreq         [32]uint16
	codegenFreq        [codegenCodeCount]uint16

	
	codegen [literalCount + offsetCodeCount + 1]uint8
}




















func newHuffmanBitWriter(w io.Writer) *huffmanBitWriter {
	return &huffmanBitWriter{
		writer:          w,
		literalEncoding: newHuffmanEncoder(literalCount),
		tmpLitEncoding:  newHuffmanEncoder(literalCount),
		codegenEncoding: newHuffmanEncoder(codegenCodeCount),
		offsetEncoding:  newHuffmanEncoder(offsetCodeCount),
	}
}

func (w *huffmanBitWriter) reset(writer io.Writer) {
	w.writer = writer
	w.bits, w.nbits, w.nbytes, w.err = 0, 0, 0, nil
	w.lastHeader = 0
	w.lastHuffMan = false
}

func (w *huffmanBitWriter) canReuse(t *tokens) (ok bool) {
	a := t.offHist[:offsetCodeCount]
	b := w.offsetEncoding.codes
	b = b[:len(a)]
	for i, v := range a {
		if v != 0 && b[i].zero() {
			return false
		}
	}

	a = t.extraHist[:literalCount-256]
	b = w.literalEncoding.codes[256:literalCount]
	b = b[:len(a)]
	for i, v := range a {
		if v != 0 && b[i].zero() {
			return false
		}
	}

	a = t.litHist[:256]
	b = w.literalEncoding.codes[:len(a)]
	for i, v := range a {
		if v != 0 && b[i].zero() {
			return false
		}
	}
	return true
}

func (w *huffmanBitWriter) flush() {
	if w.err != nil {
		w.nbits = 0
		return
	}
	if w.lastHeader > 0 {
		
		w.writeCode(w.literalEncoding.codes[endBlockMarker])
		w.lastHeader = 0
	}
	n := w.nbytes
	for w.nbits != 0 {
		w.bytes[n] = byte(w.bits)
		w.bits >>= 8
		if w.nbits > 8 { 
			w.nbits -= 8
		} else {
			w.nbits = 0
		}
		n++
	}
	w.bits = 0
	w.write(w.bytes[:n])
	w.nbytes = 0
}

func (w *huffmanBitWriter) write(b []byte) {
	if w.err != nil {
		return
	}
	_, w.err = w.writer.Write(b)
}

func (w *huffmanBitWriter) writeBits(b int32, nb uint8) {
	w.bits |= uint64(b) << (w.nbits & 63)
	w.nbits += nb
	if w.nbits >= 48 {
		w.writeOutBits()
	}
}

func (w *huffmanBitWriter) writeBytes(bytes []byte) {
	if w.err != nil {
		return
	}
	n := w.nbytes
	if w.nbits&7 != 0 {
		w.err = InternalError("writeBytes with unfinished bits")
		return
	}
	for w.nbits != 0 {
		w.bytes[n] = byte(w.bits)
		w.bits >>= 8
		w.nbits -= 8
		n++
	}
	if n != 0 {
		w.write(w.bytes[:n])
	}
	w.nbytes = 0
	w.write(bytes)
}













func (w *huffmanBitWriter) generateCodegen(numLiterals int, numOffsets int, litEnc, offEnc *huffmanEncoder) {
	for i := range w.codegenFreq {
		w.codegenFreq[i] = 0
	}
	
	
	
	
	codegen := w.codegen[:] 
	
	cgnl := codegen[:numLiterals]
	for i := range cgnl {
		cgnl[i] = litEnc.codes[i].len()
	}

	cgnl = codegen[numLiterals : numLiterals+numOffsets]
	for i := range cgnl {
		cgnl[i] = offEnc.codes[i].len()
	}
	codegen[numLiterals+numOffsets] = badCode

	size := codegen[0]
	count := 1
	outIndex := 0
	for inIndex := 1; size != badCode; inIndex++ {
		
		
		nextSize := codegen[inIndex]
		if nextSize == size {
			count++
			continue
		}
		
		if size != 0 {
			codegen[outIndex] = size
			outIndex++
			w.codegenFreq[size]++
			count--
			for count >= 3 {
				n := 6
				if n > count {
					n = count
				}
				codegen[outIndex] = 16
				outIndex++
				codegen[outIndex] = uint8(n - 3)
				outIndex++
				w.codegenFreq[16]++
				count -= n
			}
		} else {
			for count >= 11 {
				n := 138
				if n > count {
					n = count
				}
				codegen[outIndex] = 18
				outIndex++
				codegen[outIndex] = uint8(n - 11)
				outIndex++
				w.codegenFreq[18]++
				count -= n
			}
			if count >= 3 {
				
				codegen[outIndex] = 17
				outIndex++
				codegen[outIndex] = uint8(count - 3)
				outIndex++
				w.codegenFreq[17]++
				count = 0
			}
		}
		count--
		for ; count >= 0; count-- {
			codegen[outIndex] = size
			outIndex++
			w.codegenFreq[size]++
		}
		
		size = nextSize
		count = 1
	}
	
	codegen[outIndex] = badCode
}

func (w *huffmanBitWriter) codegens() int {
	numCodegens := len(w.codegenFreq)
	for numCodegens > 4 && w.codegenFreq[codegenOrder[numCodegens-1]] == 0 {
		numCodegens--
	}
	return numCodegens
}

func (w *huffmanBitWriter) headerSize() (size, numCodegens int) {
	numCodegens = len(w.codegenFreq)
	for numCodegens > 4 && w.codegenFreq[codegenOrder[numCodegens-1]] == 0 {
		numCodegens--
	}
	return 3 + 5 + 5 + 4 + (3 * numCodegens) +
		w.codegenEncoding.bitLength(w.codegenFreq[:]) +
		int(w.codegenFreq[16])*2 +
		int(w.codegenFreq[17])*3 +
		int(w.codegenFreq[18])*7, numCodegens
}


func (w *huffmanBitWriter) dynamicReuseSize(litEnc, offEnc *huffmanEncoder) (size int) {
	size = litEnc.bitLength(w.literalFreq[:]) +
		offEnc.bitLength(w.offsetFreq[:])
	return size
}


func (w *huffmanBitWriter) dynamicSize(litEnc, offEnc *huffmanEncoder, extraBits int) (size, numCodegens int) {
	header, numCodegens := w.headerSize()
	size = header +
		litEnc.bitLength(w.literalFreq[:]) +
		offEnc.bitLength(w.offsetFreq[:]) +
		extraBits
	return size, numCodegens
}



func (w *huffmanBitWriter) extraBitSize() int {
	total := 0
	for i, n := range w.literalFreq[257:literalCount] {
		total += int(n) * int(lengthExtraBits[i&31])
	}
	for i, n := range w.offsetFreq[:offsetCodeCount] {
		total += int(n) * int(offsetExtraBits[i&31])
	}
	return total
}


func (w *huffmanBitWriter) fixedSize(extraBits int) int {
	return 3 +
		fixedLiteralEncoding.bitLength(w.literalFreq[:]) +
		fixedOffsetEncoding.bitLength(w.offsetFreq[:]) +
		extraBits
}




func (w *huffmanBitWriter) storedSize(in []byte) (int, bool) {
	if in == nil {
		return 0, false
	}
	if len(in) <= maxStoreBlockSize {
		return (len(in) + 5) * 8, true
	}
	return 0, false
}

func (w *huffmanBitWriter) writeCode(c hcode) {
	
	w.bits |= c.code64() << (w.nbits & 63)
	w.nbits += c.len()
	if w.nbits >= 48 {
		w.writeOutBits()
	}
}


func (w *huffmanBitWriter) writeOutBits() {
	bits := w.bits
	w.bits >>= 48
	w.nbits -= 48
	n := w.nbytes

	
	binary.LittleEndian.PutUint64(w.bytes[n:], bits)
	n += 6

	if n >= bufferFlushSize {
		if w.err != nil {
			n = 0
			return
		}
		w.write(w.bytes[:n])
		n = 0
	}

	w.nbytes = n
}






func (w *huffmanBitWriter) writeDynamicHeader(numLiterals int, numOffsets int, numCodegens int, isEof bool) {
	if w.err != nil {
		return
	}
	var firstBits int32 = 4
	if isEof {
		firstBits = 5
	}
	w.writeBits(firstBits, 3)
	w.writeBits(int32(numLiterals-257), 5)
	w.writeBits(int32(numOffsets-1), 5)
	w.writeBits(int32(numCodegens-4), 4)

	for i := 0; i < numCodegens; i++ {
		value := uint(w.codegenEncoding.codes[codegenOrder[i]].len())
		w.writeBits(int32(value), 3)
	}

	i := 0
	for {
		var codeWord = uint32(w.codegen[i])
		i++
		if codeWord == badCode {
			break
		}
		w.writeCode(w.codegenEncoding.codes[codeWord])

		switch codeWord {
		case 16:
			w.writeBits(int32(w.codegen[i]), 2)
			i++
		case 17:
			w.writeBits(int32(w.codegen[i]), 3)
			i++
		case 18:
			w.writeBits(int32(w.codegen[i]), 7)
			i++
		}
	}
}




func (w *huffmanBitWriter) writeStoredHeader(length int, isEof bool) {
	if w.err != nil {
		return
	}
	if w.lastHeader > 0 {
		
		w.writeCode(w.literalEncoding.codes[endBlockMarker])
		w.lastHeader = 0
	}

	
	if length == 0 && isEof {
		w.writeFixedHeader(isEof)
		
		w.writeBits(0, 7)
		w.flush()
		return
	}

	var flag int32
	if isEof {
		flag = 1
	}
	w.writeBits(flag, 3)
	w.flush()
	w.writeBits(int32(length), 16)
	w.writeBits(int32(^uint16(length)), 16)
}

func (w *huffmanBitWriter) writeFixedHeader(isEof bool) {
	if w.err != nil {
		return
	}
	if w.lastHeader > 0 {
		
		w.writeCode(w.literalEncoding.codes[endBlockMarker])
		w.lastHeader = 0
	}

	
	var value int32 = 2
	if isEof {
		value = 3
	}
	w.writeBits(value, 3)
}






func (w *huffmanBitWriter) writeBlock(tokens *tokens, eof bool, input []byte) {
	if w.err != nil {
		return
	}

	tokens.AddEOB()
	if w.lastHeader > 0 {
		
		w.writeCode(w.literalEncoding.codes[endBlockMarker])
		w.lastHeader = 0
	}
	numLiterals, numOffsets := w.indexTokens(tokens, false)
	w.generate()
	var extraBits int
	storedSize, storable := w.storedSize(input)
	if storable {
		extraBits = w.extraBitSize()
	}

	
	
	var literalEncoding = fixedLiteralEncoding
	var offsetEncoding = fixedOffsetEncoding
	var size = math.MaxInt32
	if tokens.n < maxPredefinedTokens {
		size = w.fixedSize(extraBits)
	}

	
	var numCodegens int

	
	
	w.generateCodegen(numLiterals, numOffsets, w.literalEncoding, w.offsetEncoding)
	w.codegenEncoding.generate(w.codegenFreq[:], 7)
	dynamicSize, numCodegens := w.dynamicSize(w.literalEncoding, w.offsetEncoding, extraBits)

	if dynamicSize < size {
		size = dynamicSize
		literalEncoding = w.literalEncoding
		offsetEncoding = w.offsetEncoding
	}

	
	if storable && storedSize <= size {
		w.writeStoredHeader(len(input), eof)
		w.writeBytes(input)
		return
	}

	
	if literalEncoding == fixedLiteralEncoding {
		w.writeFixedHeader(eof)
	} else {
		w.writeDynamicHeader(numLiterals, numOffsets, numCodegens, eof)
	}

	
	w.writeTokens(tokens.Slice(), literalEncoding.codes, offsetEncoding.codes)
}






func (w *huffmanBitWriter) writeBlockDynamic(tokens *tokens, eof bool, input []byte, sync bool) {
	if w.err != nil {
		return
	}

	sync = sync || eof
	if sync {
		tokens.AddEOB()
	}

	
	if (w.lastHuffMan || eof) && w.lastHeader > 0 {
		
		w.writeCode(w.literalEncoding.codes[endBlockMarker])
		w.lastHeader = 0
		w.lastHuffMan = false
	}

	
	
	
	const fillReuse = false

	
	if !fillReuse && w.lastHeader > 0 && !w.canReuse(tokens) {
		w.writeCode(w.literalEncoding.codes[endBlockMarker])
		w.lastHeader = 0
	}

	numLiterals, numOffsets := w.indexTokens(tokens, !sync)
	extraBits := 0
	ssize, storable := w.storedSize(input)

	const usePrefs = true
	if storable || w.lastHeader > 0 {
		extraBits = w.extraBitSize()
	}

	var size int

	
	if w.lastHeader > 0 {
		
		
		newSize := w.lastHeader + tokens.EstimatedBits()
		newSize += int(w.literalEncoding.codes[endBlockMarker].len()) + newSize>>w.logNewTablePenalty

		
		
		reuseSize := w.dynamicReuseSize(w.literalEncoding, w.offsetEncoding) + extraBits

		
		if newSize < reuseSize {
			
			w.writeCode(w.literalEncoding.codes[endBlockMarker])
			size = newSize
			w.lastHeader = 0
		} else {
			size = reuseSize
		}

		if tokens.n < maxPredefinedTokens {
			if preSize := w.fixedSize(extraBits) + 7; usePrefs && preSize < size {
				
				if storable && ssize <= size {
					w.writeStoredHeader(len(input), eof)
					w.writeBytes(input)
					return
				}
				w.writeFixedHeader(eof)
				if !sync {
					tokens.AddEOB()
				}
				w.writeTokens(tokens.Slice(), fixedLiteralEncoding.codes, fixedOffsetEncoding.codes)
				return
			}
		}
		
		if storable && ssize <= size {
			w.writeStoredHeader(len(input), eof)
			w.writeBytes(input)
			return
		}
	}

	
	if w.lastHeader == 0 {
		if fillReuse && !sync {
			w.fillTokens()
			numLiterals, numOffsets = maxNumLit, maxNumDist
		} else {
			w.literalFreq[endBlockMarker] = 1
		}

		w.generate()
		
		
		w.generateCodegen(numLiterals, numOffsets, w.literalEncoding, w.offsetEncoding)
		w.codegenEncoding.generate(w.codegenFreq[:], 7)

		var numCodegens int
		if fillReuse && !sync {
			
			w.indexTokens(tokens, true)
		}
		size, numCodegens = w.dynamicSize(w.literalEncoding, w.offsetEncoding, extraBits)

		
		if tokens.n < maxPredefinedTokens {
			if preSize := w.fixedSize(extraBits); usePrefs && preSize <= size {
				
				if storable && ssize <= preSize {
					w.writeStoredHeader(len(input), eof)
					w.writeBytes(input)
					return
				}
				w.writeFixedHeader(eof)
				if !sync {
					tokens.AddEOB()
				}
				w.writeTokens(tokens.Slice(), fixedLiteralEncoding.codes, fixedOffsetEncoding.codes)
				return
			}
		}

		if storable && ssize <= size {
			
			w.writeStoredHeader(len(input), eof)
			w.writeBytes(input)
			return
		}

		
		w.writeDynamicHeader(numLiterals, numOffsets, numCodegens, eof)
		if !sync {
			w.lastHeader, _ = w.headerSize()
		}
		w.lastHuffMan = false
	}

	if sync {
		w.lastHeader = 0
	}
	
	w.writeTokens(tokens.Slice(), w.literalEncoding.codes, w.offsetEncoding.codes)
}

func (w *huffmanBitWriter) fillTokens() {
	for i, v := range w.literalFreq[:literalCount] {
		if v == 0 {
			w.literalFreq[i] = 1
		}
	}
	for i, v := range w.offsetFreq[:offsetCodeCount] {
		if v == 0 {
			w.offsetFreq[i] = 1
		}
	}
}





func (w *huffmanBitWriter) indexTokens(t *tokens, filled bool) (numLiterals, numOffsets int) {
	
	*(*[256]uint16)(w.literalFreq[:]) = t.litHist
	
	*(*[32]uint16)(w.literalFreq[256:]) = t.extraHist
	w.offsetFreq = t.offHist

	if t.n == 0 {
		return
	}
	if filled {
		return maxNumLit, maxNumDist
	}
	
	numLiterals = len(w.literalFreq)
	for w.literalFreq[numLiterals-1] == 0 {
		numLiterals--
	}
	
	numOffsets = len(w.offsetFreq)
	for numOffsets > 0 && w.offsetFreq[numOffsets-1] == 0 {
		numOffsets--
	}
	if numOffsets == 0 {
		
		
		w.offsetFreq[0] = 1
		numOffsets = 1
	}
	return
}

func (w *huffmanBitWriter) generate() {
	w.literalEncoding.generate(w.literalFreq[:literalCount], 15)
	w.offsetEncoding.generate(w.offsetFreq[:offsetCodeCount], 15)
}



func (w *huffmanBitWriter) writeTokens(tokens []token, leCodes, oeCodes []hcode) {
	if w.err != nil {
		return
	}
	if len(tokens) == 0 {
		return
	}

	
	var deferEOB bool
	if tokens[len(tokens)-1] == endBlockMarker {
		tokens = tokens[:len(tokens)-1]
		deferEOB = true
	}

	
	lits := leCodes[:256]
	offs := oeCodes[:32]
	lengths := leCodes[lengthCodesStart:]
	lengths = lengths[:32]

	
	bits, nbits, nbytes := w.bits, w.nbits, w.nbytes

	for _, t := range tokens {
		if t < 256 {
			
			c := lits[t]
			bits |= c.code64() << (nbits & 63)
			nbits += c.len()
			if nbits >= 48 {
				binary.LittleEndian.PutUint64(w.bytes[nbytes:], bits)
				
				bits >>= 48
				nbits -= 48
				nbytes += 6
				if nbytes >= bufferFlushSize {
					if w.err != nil {
						nbytes = 0
						return
					}
					_, w.err = w.writer.Write(w.bytes[:nbytes])
					nbytes = 0
				}
			}
			continue
		}

		
		length := t.length()
		lengthCode := lengthCode(length) & 31
		if false {
			w.writeCode(lengths[lengthCode])
		} else {
			
			c := lengths[lengthCode]
			bits |= c.code64() << (nbits & 63)
			nbits += c.len()
			if nbits >= 48 {
				binary.LittleEndian.PutUint64(w.bytes[nbytes:], bits)
				
				bits >>= 48
				nbits -= 48
				nbytes += 6
				if nbytes >= bufferFlushSize {
					if w.err != nil {
						nbytes = 0
						return
					}
					_, w.err = w.writer.Write(w.bytes[:nbytes])
					nbytes = 0
				}
			}
		}

		if lengthCode >= lengthExtraBitsMinCode {
			extraLengthBits := lengthExtraBits[lengthCode]
			
			extraLength := int32(length - lengthBase[lengthCode])
			bits |= uint64(extraLength) << (nbits & 63)
			nbits += extraLengthBits
			if nbits >= 48 {
				binary.LittleEndian.PutUint64(w.bytes[nbytes:], bits)
				
				bits >>= 48
				nbits -= 48
				nbytes += 6
				if nbytes >= bufferFlushSize {
					if w.err != nil {
						nbytes = 0
						return
					}
					_, w.err = w.writer.Write(w.bytes[:nbytes])
					nbytes = 0
				}
			}
		}
		
		offset := t.offset()
		offsetCode := (offset >> 16) & 31
		if false {
			w.writeCode(offs[offsetCode])
		} else {
			
			c := offs[offsetCode]
			bits |= c.code64() << (nbits & 63)
			nbits += c.len()
			if nbits >= 48 {
				binary.LittleEndian.PutUint64(w.bytes[nbytes:], bits)
				
				bits >>= 48
				nbits -= 48
				nbytes += 6
				if nbytes >= bufferFlushSize {
					if w.err != nil {
						nbytes = 0
						return
					}
					_, w.err = w.writer.Write(w.bytes[:nbytes])
					nbytes = 0
				}
			}
		}

		if offsetCode >= offsetExtraBitsMinCode {
			offsetComb := offsetCombined[offsetCode]
			
			bits |= uint64((offset-(offsetComb>>8))&matchOffsetOnlyMask) << (nbits & 63)
			nbits += uint8(offsetComb)
			if nbits >= 48 {
				binary.LittleEndian.PutUint64(w.bytes[nbytes:], bits)
				
				bits >>= 48
				nbits -= 48
				nbytes += 6
				if nbytes >= bufferFlushSize {
					if w.err != nil {
						nbytes = 0
						return
					}
					_, w.err = w.writer.Write(w.bytes[:nbytes])
					nbytes = 0
				}
			}
		}
	}
	
	w.bits, w.nbits, w.nbytes = bits, nbits, nbytes

	if deferEOB {
		w.writeCode(leCodes[endBlockMarker])
	}
}



var huffOffset *huffmanEncoder

func init() {
	w := newHuffmanBitWriter(nil)
	w.offsetFreq[0] = 1
	huffOffset = newHuffmanEncoder(offsetCodeCount)
	huffOffset.generate(w.offsetFreq[:offsetCodeCount], 15)
}




func (w *huffmanBitWriter) writeBlockHuff(eof bool, input []byte, sync bool) {
	if w.err != nil {
		return
	}

	
	for i := range w.literalFreq[:] {
		w.literalFreq[i] = 0
	}
	if !w.lastHuffMan {
		for i := range w.offsetFreq[:] {
			w.offsetFreq[i] = 0
		}
	}

	const numLiterals = endBlockMarker + 1
	const numOffsets = 1

	
	
	
	
	const guessHeaderSizeBits = 70 * 8
	histogram(input, w.literalFreq[:numLiterals])
	ssize, storable := w.storedSize(input)
	if storable && len(input) > 1024 {
		
		abs := float64(0)
		avg := float64(len(input)) / 256
		max := float64(len(input) * 2)
		for _, v := range w.literalFreq[:256] {
			diff := float64(v) - avg
			abs += diff * diff
			if abs > max {
				break
			}
		}
		if abs < max {
			if debugDeflate {
				fmt.Println("stored", abs, "<", max)
			}
			
			w.writeStoredHeader(len(input), eof)
			w.writeBytes(input)
			return
		}
	}
	w.literalFreq[endBlockMarker] = 1
	w.tmpLitEncoding.generate(w.literalFreq[:numLiterals], 15)
	estBits := w.tmpLitEncoding.canReuseBits(w.literalFreq[:numLiterals])
	if estBits < math.MaxInt32 {
		estBits += w.lastHeader
		if w.lastHeader == 0 {
			estBits += guessHeaderSizeBits
		}
		estBits += estBits >> w.logNewTablePenalty
	}

	
	if storable && ssize <= estBits {
		if debugDeflate {
			fmt.Println("stored,", ssize, "<=", estBits)
		}
		w.writeStoredHeader(len(input), eof)
		w.writeBytes(input)
		return
	}

	if w.lastHeader > 0 {
		reuseSize := w.literalEncoding.canReuseBits(w.literalFreq[:256])

		if estBits < reuseSize {
			if debugDeflate {
				fmt.Println("NOT reusing, reuse:", reuseSize/8, "> new:", estBits/8, "header est:", w.lastHeader/8, "bytes")
			}
			
			w.writeCode(w.literalEncoding.codes[endBlockMarker])
			w.lastHeader = 0
		} else if debugDeflate {
			fmt.Println("reusing, reuse:", reuseSize/8, "> new:", estBits/8, "- header est:", w.lastHeader/8)
		}
	}

	count := 0
	if w.lastHeader == 0 {
		
		w.literalEncoding, w.tmpLitEncoding = w.tmpLitEncoding, w.literalEncoding
		
		
		w.generateCodegen(numLiterals, numOffsets, w.literalEncoding, huffOffset)
		w.codegenEncoding.generate(w.codegenFreq[:], 7)
		numCodegens := w.codegens()

		
		w.writeDynamicHeader(numLiterals, numOffsets, numCodegens, eof)
		w.lastHuffMan = true
		w.lastHeader, _ = w.headerSize()
		if debugDeflate {
			count += w.lastHeader
			fmt.Println("header:", count/8)
		}
	}

	encoding := w.literalEncoding.codes[:256]
	
	bits, nbits, nbytes := w.bits, w.nbits, w.nbytes

	if debugDeflate {
		count -= int(nbytes)*8 + int(nbits)
	}
	
	
	for len(input) > 3 {
		
		if nbits >= 8 {
			n := nbits >> 3
			binary.LittleEndian.PutUint64(w.bytes[nbytes:], bits)
			bits >>= (n * 8) & 63
			nbits -= n * 8
			nbytes += n
		}
		if nbytes >= bufferFlushSize {
			if w.err != nil {
				nbytes = 0
				return
			}
			if debugDeflate {
				count += int(nbytes) * 8
			}
			_, w.err = w.writer.Write(w.bytes[:nbytes])
			nbytes = 0
		}
		a, b := encoding[input[0]], encoding[input[1]]
		bits |= a.code64() << (nbits & 63)
		bits |= b.code64() << ((nbits + a.len()) & 63)
		c := encoding[input[2]]
		nbits += b.len() + a.len()
		bits |= c.code64() << (nbits & 63)
		nbits += c.len()
		input = input[3:]
	}

	
	for _, t := range input {
		if nbits >= 48 {
			binary.LittleEndian.PutUint64(w.bytes[nbytes:], bits)
			
			bits >>= 48
			nbits -= 48
			nbytes += 6
			if nbytes >= bufferFlushSize {
				if w.err != nil {
					nbytes = 0
					return
				}
				if debugDeflate {
					count += int(nbytes) * 8
				}
				_, w.err = w.writer.Write(w.bytes[:nbytes])
				nbytes = 0
			}
		}
		
		c := encoding[t]
		bits |= c.code64() << (nbits & 63)

		nbits += c.len()
		if debugDeflate {
			count += int(c.len())
		}
	}
	
	w.bits, w.nbits, w.nbytes = bits, nbits, nbytes

	if debugDeflate {
		nb := count + int(nbytes)*8 + int(nbits)
		fmt.Println("wrote", nb, "bits,", nb/8, "bytes.")
	}
	
	if w.nbits >= 48 {
		w.writeOutBits()
	}

	if eof || sync {
		w.writeCode(w.literalEncoding.codes[endBlockMarker])
		w.lastHeader = 0
		w.lastHuffMan = false
	}
}
