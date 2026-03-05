



//go:build gc && !purego

package sha3





import (
	"hash"

	"golang.org/x/sys/cpu"
)



type code uint64

const (
	
	sha3_224  code = 32
	sha3_256       = 33
	sha3_384       = 34
	sha3_512       = 35
	shake_128      = 36
	shake_256      = 37
	nopad          = 0x100
)




//go:noescape
func kimd(function code, chain *[200]byte, src []byte)




//go:noescape
func klmd(function code, chain *[200]byte, dst, src []byte)

type asmState struct {
	a         [200]byte       
	buf       []byte          
	rate      int             
	storage   [3072]byte      
	outputLen int             
	function  code            
	state     spongeDirection 
}

func newAsmState(function code) *asmState {
	var s asmState
	s.function = function
	switch function {
	case sha3_224:
		s.rate = 144
		s.outputLen = 28
	case sha3_256:
		s.rate = 136
		s.outputLen = 32
	case sha3_384:
		s.rate = 104
		s.outputLen = 48
	case sha3_512:
		s.rate = 72
		s.outputLen = 64
	case shake_128:
		s.rate = 168
		s.outputLen = 32
	case shake_256:
		s.rate = 136
		s.outputLen = 64
	default:
		panic("sha3: unrecognized function code")
	}

	
	s.resetBuf()
	return &s
}

func (s *asmState) clone() *asmState {
	c := *s
	c.buf = c.storage[:len(s.buf):cap(s.buf)]
	return &c
}



func (s *asmState) copyIntoBuf(b []byte) {
	bufLen := len(s.buf)
	s.buf = s.buf[:len(s.buf)+len(b)]
	copy(s.buf[bufLen:], b)
}



func (s *asmState) resetBuf() {
	max := (cap(s.storage) / s.rate) * s.rate
	s.buf = s.storage[:0:max]
}



func (s *asmState) Write(b []byte) (int, error) {
	if s.state != spongeAbsorbing {
		panic("sha3: Write after Read")
	}
	length := len(b)
	for len(b) > 0 {
		if len(s.buf) == 0 && len(b) >= cap(s.buf) {
			
			
			remainder := len(b) % s.rate
			kimd(s.function, &s.a, b[:len(b)-remainder])
			if remainder != 0 {
				s.copyIntoBuf(b[len(b)-remainder:])
			}
			return length, nil
		}

		if len(s.buf) == cap(s.buf) {
			
			kimd(s.function, &s.a, s.buf)
			s.buf = s.buf[:0]
		}

		
		n := len(b)
		if len(b) > cap(s.buf)-len(s.buf) {
			n = cap(s.buf) - len(s.buf)
		}
		s.copyIntoBuf(b[:n])
		b = b[n:]
	}
	return length, nil
}


func (s *asmState) Read(out []byte) (n int, err error) {
	
	
	if s.function != shake_128 && s.function != shake_256 {
		panic("sha3: can only call Read for SHAKE functions")
	}

	n = len(out)

	
	if s.state == spongeAbsorbing {
		s.state = spongeSqueezing

		
		if len(out)%s.rate == 0 {
			klmd(s.function, &s.a, out, s.buf) 
			s.buf = s.buf[:0]
			return
		}

		
		max := cap(s.buf)
		if max > len(out) {
			max = (len(out)/s.rate)*s.rate + s.rate
		}
		klmd(s.function, &s.a, s.buf[:max], s.buf)
		s.buf = s.buf[:max]
	}

	for len(out) > 0 {
		
		if len(s.buf) != 0 {
			c := copy(out, s.buf)
			out = out[c:]
			s.buf = s.buf[c:]
			continue
		}

		
		if len(out)%s.rate == 0 {
			klmd(s.function|nopad, &s.a, out, nil)
			return
		}

		
		s.resetBuf()
		if cap(s.buf) > len(out) {
			s.buf = s.buf[:(len(out)/s.rate)*s.rate+s.rate]
		}
		klmd(s.function|nopad, &s.a, s.buf, nil)
	}
	return
}



func (s *asmState) Sum(b []byte) []byte {
	if s.state != spongeAbsorbing {
		panic("sha3: Sum after Read")
	}

	
	a := s.a

	
	
	switch s.function {
	case sha3_224, sha3_256, sha3_384, sha3_512:
		klmd(s.function, &a, nil, s.buf)
		return append(b, a[:s.outputLen]...)
	case shake_128, shake_256:
		d := make([]byte, s.outputLen, 64)
		klmd(s.function, &a, d, s.buf)
		return append(b, d[:s.outputLen]...)
	default:
		panic("sha3: unknown function")
	}
}


func (s *asmState) Reset() {
	for i := range s.a {
		s.a[i] = 0
	}
	s.resetBuf()
	s.state = spongeAbsorbing
}


func (s *asmState) Size() int {
	return s.outputLen
}





func (s *asmState) BlockSize() int {
	return s.rate
}


func (s *asmState) Clone() ShakeHash {
	return s.clone()
}



func new224() hash.Hash {
	if cpu.S390X.HasSHA3 {
		return newAsmState(sha3_224)
	}
	return new224Generic()
}



func new256() hash.Hash {
	if cpu.S390X.HasSHA3 {
		return newAsmState(sha3_256)
	}
	return new256Generic()
}



func new384() hash.Hash {
	if cpu.S390X.HasSHA3 {
		return newAsmState(sha3_384)
	}
	return new384Generic()
}



func new512() hash.Hash {
	if cpu.S390X.HasSHA3 {
		return newAsmState(sha3_512)
	}
	return new512Generic()
}



func newShake128() ShakeHash {
	if cpu.S390X.HasSHA3 {
		return newAsmState(shake_128)
	}
	return newShake128Generic()
}



func newShake256() ShakeHash {
	if cpu.S390X.HasSHA3 {
		return newAsmState(shake_256)
	}
	return newShake256Generic()
}
