



package sha3

import (
	"crypto/subtle"
	"encoding/binary"
	"errors"
	"unsafe"

	"golang.org/x/sys/cpu"
)


type spongeDirection int

const (
	
	spongeAbsorbing spongeDirection = iota
	
	spongeSqueezing
)

type state struct {
	a [1600 / 8]byte 

	
	
	
	n, rate int

	
	
	
	
	
	
	
	
	
	
	
	
	dsbyte byte

	outputLen int             
	state     spongeDirection 
}


func (d *state) BlockSize() int { return d.rate }


func (d *state) Size() int { return d.outputLen }



func (d *state) Reset() {
	
	for i := range d.a {
		d.a[i] = 0
	}
	d.state = spongeAbsorbing
	d.n = 0
}

func (d *state) clone() *state {
	ret := *d
	return &ret
}


func (d *state) permute() {
	var a *[25]uint64
	if cpu.IsBigEndian {
		a = new([25]uint64)
		for i := range a {
			a[i] = binary.LittleEndian.Uint64(d.a[i*8:])
		}
	} else {
		a = (*[25]uint64)(unsafe.Pointer(&d.a))
	}

	keccakF1600(a)
	d.n = 0

	if cpu.IsBigEndian {
		for i := range a {
			binary.LittleEndian.PutUint64(d.a[i*8:], a[i])
		}
	}
}



func (d *state) padAndPermute() {
	
	
	
	
	d.a[d.n] ^= d.dsbyte
	
	
	
	d.a[d.rate-1] ^= 0x80
	
	d.permute()
	d.state = spongeSqueezing
}



func (d *state) Write(p []byte) (n int, err error) {
	if d.state != spongeAbsorbing {
		panic("sha3: Write after Read")
	}

	n = len(p)

	for len(p) > 0 {
		x := subtle.XORBytes(d.a[d.n:d.rate], d.a[d.n:d.rate], p)
		d.n += x
		p = p[x:]

		
		if d.n == d.rate {
			d.permute()
		}
	}

	return
}


func (d *state) Read(out []byte) (n int, err error) {
	
	if d.state == spongeAbsorbing {
		d.padAndPermute()
	}

	n = len(out)

	
	for len(out) > 0 {
		
		if d.n == d.rate {
			d.permute()
		}

		x := copy(out, d.a[d.n:d.rate])
		d.n += x
		out = out[x:]
	}

	return
}



func (d *state) Sum(in []byte) []byte {
	if d.state != spongeAbsorbing {
		panic("sha3: Sum after Read")
	}

	
	
	dup := d.clone()
	hash := make([]byte, dup.outputLen, 64) 
	dup.Read(hash)
	return append(in, hash...)
}

const (
	magicSHA3   = "sha\x08"
	magicShake  = "sha\x09"
	magicCShake = "sha\x0a"
	magicKeccak = "sha\x0b"
	
	marshaledSize = len(magicSHA3) + 1 + 200 + 1 + 1
)

func (d *state) MarshalBinary() ([]byte, error) {
	return d.AppendBinary(make([]byte, 0, marshaledSize))
}

func (d *state) AppendBinary(b []byte) ([]byte, error) {
	switch d.dsbyte {
	case dsbyteSHA3:
		b = append(b, magicSHA3...)
	case dsbyteShake:
		b = append(b, magicShake...)
	case dsbyteCShake:
		b = append(b, magicCShake...)
	case dsbyteKeccak:
		b = append(b, magicKeccak...)
	default:
		panic("unknown dsbyte")
	}
	
	b = append(b, byte(d.rate))
	b = append(b, d.a[:]...)
	b = append(b, byte(d.n), byte(d.state))
	return b, nil
}

func (d *state) UnmarshalBinary(b []byte) error {
	if len(b) != marshaledSize {
		return errors.New("sha3: invalid hash state")
	}

	magic := string(b[:len(magicSHA3)])
	b = b[len(magicSHA3):]
	switch {
	case magic == magicSHA3 && d.dsbyte == dsbyteSHA3:
	case magic == magicShake && d.dsbyte == dsbyteShake:
	case magic == magicCShake && d.dsbyte == dsbyteCShake:
	case magic == magicKeccak && d.dsbyte == dsbyteKeccak:
	default:
		return errors.New("sha3: invalid hash state identifier")
	}

	rate := int(b[0])
	b = b[1:]
	if rate != d.rate {
		return errors.New("sha3: invalid hash state function")
	}

	copy(d.a[:], b)
	b = b[len(d.a):]

	n, state := int(b[0]), spongeDirection(b[1])
	if n > d.rate {
		return errors.New("sha3: invalid hash state")
	}
	d.n = n
	if state != spongeAbsorbing && state != spongeSqueezing {
		return errors.New("sha3: invalid hash state")
	}
	d.state = state

	return nil
}
