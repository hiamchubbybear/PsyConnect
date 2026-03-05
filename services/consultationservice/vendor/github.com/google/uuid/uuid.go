



package uuid

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
)



type UUID [16]byte


type Version byte


type Variant byte


const (
	Invalid   = Variant(iota) 
	RFC4122                   
	Reserved                  
	Microsoft                 
	Future                    
)

const randPoolSize = 16 * 16

var (
	rander      = rand.Reader 
	poolEnabled = false
	poolMu      sync.Mutex
	poolPos     = randPoolSize     
	pool        [randPoolSize]byte 
)

type invalidLengthError struct{ len int }

func (err invalidLengthError) Error() string {
	return fmt.Sprintf("invalid UUID length: %d", err.len)
}


func IsInvalidLengthError(err error) bool {
	_, ok := err.(invalidLengthError)
	return ok
}










func Parse(s string) (UUID, error) {
	var uuid UUID
	switch len(s) {
	
	case 36:

	
	case 36 + 9:
		if !strings.EqualFold(s[:9], "urn:uuid:") {
			return uuid, fmt.Errorf("invalid urn prefix: %q", s[:9])
		}
		s = s[9:]

	
	case 36 + 2:
		s = s[1:]

	
	case 32:
		var ok bool
		for i := range uuid {
			uuid[i], ok = xtob(s[i*2], s[i*2+1])
			if !ok {
				return uuid, errors.New("invalid UUID format")
			}
		}
		return uuid, nil
	default:
		return uuid, invalidLengthError{len(s)}
	}
	
	
	if s[8] != '-' || s[13] != '-' || s[18] != '-' || s[23] != '-' {
		return uuid, errors.New("invalid UUID format")
	}
	for i, x := range [16]int{
		0, 2, 4, 6,
		9, 11,
		14, 16,
		19, 21,
		24, 26, 28, 30, 32, 34,
	} {
		v, ok := xtob(s[x], s[x+1])
		if !ok {
			return uuid, errors.New("invalid UUID format")
		}
		uuid[i] = v
	}
	return uuid, nil
}


func ParseBytes(b []byte) (UUID, error) {
	var uuid UUID
	switch len(b) {
	case 36: 
	case 36 + 9: 
		if !bytes.EqualFold(b[:9], []byte("urn:uuid:")) {
			return uuid, fmt.Errorf("invalid urn prefix: %q", b[:9])
		}
		b = b[9:]
	case 36 + 2: 
		b = b[1:]
	case 32: 
		var ok bool
		for i := 0; i < 32; i += 2 {
			uuid[i/2], ok = xtob(b[i], b[i+1])
			if !ok {
				return uuid, errors.New("invalid UUID format")
			}
		}
		return uuid, nil
	default:
		return uuid, invalidLengthError{len(b)}
	}
	
	
	if b[8] != '-' || b[13] != '-' || b[18] != '-' || b[23] != '-' {
		return uuid, errors.New("invalid UUID format")
	}
	for i, x := range [16]int{
		0, 2, 4, 6,
		9, 11,
		14, 16,
		19, 21,
		24, 26, 28, 30, 32, 34,
	} {
		v, ok := xtob(b[x], b[x+1])
		if !ok {
			return uuid, errors.New("invalid UUID format")
		}
		uuid[i] = v
	}
	return uuid, nil
}



func MustParse(s string) UUID {
	uuid, err := Parse(s)
	if err != nil {
		panic(`uuid: Parse(` + s + `): ` + err.Error())
	}
	return uuid
}



func FromBytes(b []byte) (uuid UUID, err error) {
	err = uuid.UnmarshalBinary(b)
	return uuid, err
}


func Must(uuid UUID, err error) UUID {
	if err != nil {
		panic(err)
	}
	return uuid
}







func Validate(s string) error {
	switch len(s) {
	
	case 36:

	
	case 36 + 9:
		if !strings.EqualFold(s[:9], "urn:uuid:") {
			return fmt.Errorf("invalid urn prefix: %q", s[:9])
		}
		s = s[9:]

	
	case 36 + 2:
		if s[0] != '{' || s[len(s)-1] != '}' {
			return fmt.Errorf("invalid bracketed UUID format")
		}
		s = s[1 : len(s)-1]

	
	case 32:
		for i := 0; i < len(s); i += 2 {
			_, ok := xtob(s[i], s[i+1])
			if !ok {
				return errors.New("invalid UUID format")
			}
		}

	default:
		return invalidLengthError{len(s)}
	}

	
	if len(s) == 36 {
		if s[8] != '-' || s[13] != '-' || s[18] != '-' || s[23] != '-' {
			return errors.New("invalid UUID format")
		}
		for _, x := range []int{0, 2, 4, 6, 9, 11, 14, 16, 19, 21, 24, 26, 28, 30, 32, 34} {
			if _, ok := xtob(s[x], s[x+1]); !ok {
				return errors.New("invalid UUID format")
			}
		}
	}

	return nil
}



func (uuid UUID) String() string {
	var buf [36]byte
	encodeHex(buf[:], uuid)
	return string(buf[:])
}



func (uuid UUID) URN() string {
	var buf [36 + 9]byte
	copy(buf[:], "urn:uuid:")
	encodeHex(buf[9:], uuid)
	return string(buf[:])
}

func encodeHex(dst []byte, uuid UUID) {
	hex.Encode(dst, uuid[:4])
	dst[8] = '-'
	hex.Encode(dst[9:13], uuid[4:6])
	dst[13] = '-'
	hex.Encode(dst[14:18], uuid[6:8])
	dst[18] = '-'
	hex.Encode(dst[19:23], uuid[8:10])
	dst[23] = '-'
	hex.Encode(dst[24:], uuid[10:])
}


func (uuid UUID) Variant() Variant {
	switch {
	case (uuid[8] & 0xc0) == 0x80:
		return RFC4122
	case (uuid[8] & 0xe0) == 0xc0:
		return Microsoft
	case (uuid[8] & 0xe0) == 0xe0:
		return Future
	default:
		return Reserved
	}
}


func (uuid UUID) Version() Version {
	return Version(uuid[6] >> 4)
}

func (v Version) String() string {
	if v > 15 {
		return fmt.Sprintf("BAD_VERSION_%d", v)
	}
	return fmt.Sprintf("VERSION_%d", v)
}

func (v Variant) String() string {
	switch v {
	case RFC4122:
		return "RFC4122"
	case Reserved:
		return "Reserved"
	case Microsoft:
		return "Microsoft"
	case Future:
		return "Future"
	case Invalid:
		return "Invalid"
	}
	return fmt.Sprintf("BadVariant%d", int(v))
}







func SetRand(r io.Reader) {
	if r == nil {
		rander = rand.Reader
		return
	}
	rander = r
}












func EnableRandPool() {
	poolEnabled = true
}







func DisableRandPool() {
	poolEnabled = false
	defer poolMu.Unlock()
	poolMu.Lock()
	poolPos = randPoolSize
}


type UUIDs []UUID


func (uuids UUIDs) Strings() []string {
	var uuidStrs = make([]string, len(uuids))
	for i, uuid := range uuids {
		uuidStrs[i] = uuid.String()
	}
	return uuidStrs
}
