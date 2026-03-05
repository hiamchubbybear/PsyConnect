



package sha3





import (
	"crypto"
	"hash"
)




func New224() hash.Hash {
	return new224()
}




func New256() hash.Hash {
	return new256()
}




func New384() hash.Hash {
	return new384()
}




func New512() hash.Hash {
	return new512()
}

func init() {
	crypto.RegisterHash(crypto.SHA3_224, New224)
	crypto.RegisterHash(crypto.SHA3_256, New256)
	crypto.RegisterHash(crypto.SHA3_384, New384)
	crypto.RegisterHash(crypto.SHA3_512, New512)
}

const (
	dsbyteSHA3   = 0b00000110
	dsbyteKeccak = 0b00000001
	dsbyteShake  = 0b00011111
	dsbyteCShake = 0b00000100

	
	
	rateK256  = (1600 - 256) / 8
	rateK448  = (1600 - 448) / 8
	rateK512  = (1600 - 512) / 8
	rateK768  = (1600 - 768) / 8
	rateK1024 = (1600 - 1024) / 8
)

func new224Generic() *state {
	return &state{rate: rateK448, outputLen: 28, dsbyte: dsbyteSHA3}
}

func new256Generic() *state {
	return &state{rate: rateK512, outputLen: 32, dsbyte: dsbyteSHA3}
}

func new384Generic() *state {
	return &state{rate: rateK768, outputLen: 48, dsbyte: dsbyteSHA3}
}

func new512Generic() *state {
	return &state{rate: rateK1024, outputLen: 64, dsbyte: dsbyteSHA3}
}





func NewLegacyKeccak256() hash.Hash {
	return &state{rate: rateK512, outputLen: 32, dsbyte: dsbyteKeccak}
}





func NewLegacyKeccak512() hash.Hash {
	return &state{rate: rateK1024, outputLen: 64, dsbyte: dsbyteKeccak}
}


func Sum224(data []byte) (digest [28]byte) {
	h := New224()
	h.Write(data)
	h.Sum(digest[:0])
	return
}


func Sum256(data []byte) (digest [32]byte) {
	h := New256()
	h.Write(data)
	h.Sum(digest[:0])
	return
}


func Sum384(data []byte) (digest [48]byte) {
	h := New384()
	h.Write(data)
	h.Sum(digest[:0])
	return
}


func Sum512(data []byte) (digest [64]byte) {
	h := New512()
	h.Write(data)
	h.Sum(digest[:0])
	return
}
