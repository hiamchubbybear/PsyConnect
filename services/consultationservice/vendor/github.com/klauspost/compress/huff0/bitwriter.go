




package huff0



type bitWriter struct {
	bitContainer uint64
	nBits        uint8
	out          []byte
}



func (b *bitWriter) addBits16Clean(value uint16, bits uint8) {
	b.bitContainer |= uint64(value) << (b.nBits & 63)
	b.nBits += bits
}



func (b *bitWriter) encSymbol(ct cTable, symbol byte) {
	enc := ct[symbol]
	b.bitContainer |= uint64(enc.val) << (b.nBits & 63)
	if false {
		if enc.nBits == 0 {
			panic("nbits 0")
		}
	}
	b.nBits += enc.nBits
}



func (b *bitWriter) encTwoSymbols(ct cTable, av, bv byte) {
	encA := ct[av]
	encB := ct[bv]
	sh := b.nBits & 63
	combined := uint64(encA.val) | (uint64(encB.val) << (encA.nBits & 63))
	b.bitContainer |= combined << sh
	if false {
		if encA.nBits == 0 {
			panic("nbitsA 0")
		}
		if encB.nBits == 0 {
			panic("nbitsB 0")
		}
	}
	b.nBits += encA.nBits + encB.nBits
}




func (b *bitWriter) encFourSymbols(encA, encB, encC, encD cTableEntry) {
	bitsA := encA.nBits
	bitsB := bitsA + encB.nBits
	bitsC := bitsB + encC.nBits
	bitsD := bitsC + encD.nBits
	combined := uint64(encA.val) |
		(uint64(encB.val) << (bitsA & 63)) |
		(uint64(encC.val) << (bitsB & 63)) |
		(uint64(encD.val) << (bitsC & 63))
	b.bitContainer |= combined << (b.nBits & 63)
	b.nBits += bitsD
}


func (b *bitWriter) flush32() {
	if b.nBits < 32 {
		return
	}
	b.out = append(b.out,
		byte(b.bitContainer),
		byte(b.bitContainer>>8),
		byte(b.bitContainer>>16),
		byte(b.bitContainer>>24))
	b.nBits -= 32
	b.bitContainer >>= 32
}


func (b *bitWriter) flushAlign() {
	nbBytes := (b.nBits + 7) >> 3
	for i := uint8(0); i < nbBytes; i++ {
		b.out = append(b.out, byte(b.bitContainer>>(i*8)))
	}
	b.nBits = 0
	b.bitContainer = 0
}



func (b *bitWriter) close() error {
	
	b.addBits16Clean(1, 1)
	
	b.flushAlign()
	return nil
}
