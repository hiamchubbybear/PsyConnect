





package rand

import (
	"encoding/binary"
	"io"
	"math/bits"
)













type PCGSource struct {
	low  uint64
	high uint64
}

const (
	maxUint32 = (1 << 32) - 1

	multiplier = 47026247687942121848144207491837523525
	mulHigh    = multiplier >> 64
	mulLow     = multiplier & maxUint64

	increment = 117397592171526113268558934119004209487
	incHigh   = increment >> 64
	incLow    = increment & maxUint64

	
	initializer = 245720598905631564143578724636268694099
	initHigh    = initializer >> 64
	initLow     = initializer & maxUint64
)


func (pcg *PCGSource) Seed(seed uint64) {
	pcg.low = seed
	pcg.high = seed 
}


func (pcg *PCGSource) Uint64() uint64 {
	pcg.multiply()
	pcg.add()
	
	return bits.RotateLeft64(pcg.high^pcg.low, -int(pcg.high>>58))
}

func (pcg *PCGSource) add() {
	var carry uint64
	pcg.low, carry = Add64(pcg.low, incLow, 0)
	pcg.high, _ = Add64(pcg.high, incHigh, carry)
}

func (pcg *PCGSource) multiply() {
	hi, lo := Mul64(pcg.low, mulLow)
	hi += pcg.high * mulLow
	hi += pcg.low * mulHigh
	pcg.low = lo
	pcg.high = hi
}


func (pcg *PCGSource) MarshalBinary() ([]byte, error) {
	var buf [16]byte
	binary.BigEndian.PutUint64(buf[:8], pcg.high)
	binary.BigEndian.PutUint64(buf[8:], pcg.low)
	return buf[:], nil
}


func (pcg *PCGSource) UnmarshalBinary(data []byte) error {
	if len(data) < 16 {
		return io.ErrUnexpectedEOF
	}
	pcg.low = binary.BigEndian.Uint64(data[8:])
	pcg.high = binary.BigEndian.Uint64(data[:8])
	return nil
}
