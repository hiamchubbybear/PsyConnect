






package randutil

import (
	crand "crypto/rand"
	"fmt"
	"io"

	xrand "go.mongodb.org/mongo-driver/internal/rand"
)




func NewLockedRand() *xrand.Rand {
	var randSrc = new(xrand.LockedSource)
	randSrc.Seed(cryptoSeed())
	return xrand.New(randSrc)
}




func cryptoSeed() uint64 {
	var b [8]byte
	_, err := io.ReadFull(crand.Reader, b[:])
	if err != nil {
		panic(fmt.Errorf("failed to read 8 bytes from a \"crypto/rand\".Reader: %v", err))
	}

	return (uint64(b[0]) << 0) | (uint64(b[1]) << 8) | (uint64(b[2]) << 16) | (uint64(b[3]) << 24) |
		(uint64(b[4]) << 32) | (uint64(b[5]) << 40) | (uint64(b[6]) << 48) | (uint64(b[7]) << 56)
}
