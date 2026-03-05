



package sha3












import (
	"bytes"
	"encoding/binary"
	"errors"
	"hash"
	"io"
	"math/bits"
)





type ShakeHash interface {
	hash.Hash

	
	
	
	
	io.Reader

	
	Clone() ShakeHash
}


type cshakeState struct {
	*state 

	
	
	
	
	
	initBlock []byte
}

func bytepad(data []byte, rate int) []byte {
	out := make([]byte, 0, 9+len(data)+rate-1)
	out = append(out, leftEncode(uint64(rate))...)
	out = append(out, data...)
	if padlen := rate - len(out)%rate; padlen < rate {
		out = append(out, make([]byte, padlen)...)
	}
	return out
}

func leftEncode(x uint64) []byte {
	
	n := (bits.Len64(x) + 7) / 8
	if n == 0 {
		n = 1
	}
	
	b := make([]byte, 9)
	binary.BigEndian.PutUint64(b[1:], x)
	b = b[9-n-1:]
	b[0] = byte(n)
	return b
}

func newCShake(N, S []byte, rate, outputLen int, dsbyte byte) ShakeHash {
	c := cshakeState{state: &state{rate: rate, outputLen: outputLen, dsbyte: dsbyte}}
	c.initBlock = make([]byte, 0, 9+len(N)+9+len(S)) 
	c.initBlock = append(c.initBlock, leftEncode(uint64(len(N))*8)...)
	c.initBlock = append(c.initBlock, N...)
	c.initBlock = append(c.initBlock, leftEncode(uint64(len(S))*8)...)
	c.initBlock = append(c.initBlock, S...)
	c.Write(bytepad(c.initBlock, c.rate))
	return &c
}


func (c *cshakeState) Reset() {
	c.state.Reset()
	c.Write(bytepad(c.initBlock, c.rate))
}


func (c *cshakeState) Clone() ShakeHash {
	b := make([]byte, len(c.initBlock))
	copy(b, c.initBlock)
	return &cshakeState{state: c.clone(), initBlock: b}
}


func (c *state) Clone() ShakeHash {
	return c.clone()
}

func (c *cshakeState) MarshalBinary() ([]byte, error) {
	return c.AppendBinary(make([]byte, 0, marshaledSize+len(c.initBlock)))
}

func (c *cshakeState) AppendBinary(b []byte) ([]byte, error) {
	b, err := c.state.AppendBinary(b)
	if err != nil {
		return nil, err
	}
	b = append(b, c.initBlock...)
	return b, nil
}

func (c *cshakeState) UnmarshalBinary(b []byte) error {
	if len(b) <= marshaledSize {
		return errors.New("sha3: invalid hash state")
	}
	if err := c.state.UnmarshalBinary(b[:marshaledSize]); err != nil {
		return err
	}
	c.initBlock = bytes.Clone(b[marshaledSize:])
	return nil
}




func NewShake128() ShakeHash {
	return newShake128()
}




func NewShake256() ShakeHash {
	return newShake256()
}

func newShake128Generic() *state {
	return &state{rate: rateK256, outputLen: 32, dsbyte: dsbyteShake}
}

func newShake256Generic() *state {
	return &state{rate: rateK512, outputLen: 64, dsbyte: dsbyteShake}
}







func NewCShake128(N, S []byte) ShakeHash {
	if len(N) == 0 && len(S) == 0 {
		return NewShake128()
	}
	return newCShake(N, S, rateK256, 32, dsbyteCShake)
}







func NewCShake256(N, S []byte) ShakeHash {
	if len(N) == 0 && len(S) == 0 {
		return NewShake256()
	}
	return newCShake(N, S, rateK512, 64, dsbyteCShake)
}


func ShakeSum128(hash, data []byte) {
	h := NewShake128()
	h.Write(data)
	h.Read(hash)
}


func ShakeSum256(hash, data []byte) {
	h := NewShake256()
	h.Write(data)
	h.Read(hash)
}
