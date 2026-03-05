



















package pbkdf2

import (
	"crypto/hmac"
	"encoding/binary"
	"hash"
)







func Key(password, salt []byte, iterCount, keyLen int, h func() hash.Hash) []byte {
	prf := hmac.New(h, password)
	hLen := prf.Size()
	numBlocks := keyLen / hLen
	
	if keyLen%hLen > 0 {
		numBlocks++
	}

	Ti := make([]byte, hLen)
	Uj := make([]byte, hLen)
	dk := make([]byte, 0, hLen*numBlocks)
	buf := make([]byte, 4)

	for i := uint32(1); i <= uint32(numBlocks); i++ {
		
		
		binary.BigEndian.PutUint32(buf, i)
		prf.Reset()
		prf.Write(salt)
		prf.Write(buf)
		Uj = Uj[:0]
		Uj = prf.Sum(Uj)

		
		copy(Ti, Uj)
		for j := 2; j <= iterCount; j++ {
			prf.Reset()
			prf.Write(Uj)
			Uj = Uj[:0]
			Uj = prf.Sum(Uj)
			for k := range Uj {
				Ti[k] ^= Uj[k]
			}
		}

		
		dk = append(dk, Ti...)
	}

	return dk[0:keyLen]
}
