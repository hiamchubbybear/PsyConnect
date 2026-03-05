



package json

import (
	"math"
	"math/bits"
	"strconv"
	"strings"
	"unicode/utf8"

	"google.golang.org/protobuf/internal/detrand"
	"google.golang.org/protobuf/internal/errors"
)


type kind uint8

const (
	_ kind = (1 << iota) / 2
	name
	scalar
	objectOpen
	objectClose
	arrayOpen
	arrayClose
)



type Encoder struct {
	indent   string
	lastKind kind
	indents  []byte
	out      []byte
}





func NewEncoder(buf []byte, indent string) (*Encoder, error) {
	e := &Encoder{
		out: buf,
	}
	if len(indent) > 0 {
		if strings.Trim(indent, " \t") != "" {
			return nil, errors.New("indent may only be composed of space or tab characters")
		}
		e.indent = indent
	}
	return e, nil
}


func (e *Encoder) Bytes() []byte {
	return e.out
}


func (e *Encoder) WriteNull() {
	e.prepareNext(scalar)
	e.out = append(e.out, "null"...)
}


func (e *Encoder) WriteBool(b bool) {
	e.prepareNext(scalar)
	if b {
		e.out = append(e.out, "true"...)
	} else {
		e.out = append(e.out, "false"...)
	}
}



func (e *Encoder) WriteString(s string) error {
	e.prepareNext(scalar)
	var err error
	if e.out, err = appendString(e.out, s); err != nil {
		return err
	}
	return nil
}


var errInvalidUTF8 = errors.New("invalid UTF-8")

func appendString(out []byte, in string) ([]byte, error) {
	out = append(out, '"')
	i := indexNeedEscapeInString(in)
	in, out = in[i:], append(out, in[:i]...)
	for len(in) > 0 {
		switch r, n := utf8.DecodeRuneInString(in); {
		case r == utf8.RuneError && n == 1:
			return out, errInvalidUTF8
		case r < ' ' || r == '"' || r == '\\':
			out = append(out, '\\')
			switch r {
			case '"', '\\':
				out = append(out, byte(r))
			case '\b':
				out = append(out, 'b')
			case '\f':
				out = append(out, 'f')
			case '\n':
				out = append(out, 'n')
			case '\r':
				out = append(out, 'r')
			case '\t':
				out = append(out, 't')
			default:
				out = append(out, 'u')
				out = append(out, "0000"[1+(bits.Len32(uint32(r))-1)/4:]...)
				out = strconv.AppendUint(out, uint64(r), 16)
			}
			in = in[n:]
		default:
			i := indexNeedEscapeInString(in[n:])
			in, out = in[n+i:], append(out, in[:n+i]...)
		}
	}
	out = append(out, '"')
	return out, nil
}



func indexNeedEscapeInString(s string) int {
	for i, r := range s {
		if r < ' ' || r == '\\' || r == '"' || r == utf8.RuneError {
			return i
		}
	}
	return len(s)
}


func (e *Encoder) WriteFloat(n float64, bitSize int) {
	e.prepareNext(scalar)
	e.out = appendFloat(e.out, n, bitSize)
}


func appendFloat(out []byte, n float64, bitSize int) []byte {
	switch {
	case math.IsNaN(n):
		return append(out, `"NaN"`...)
	case math.IsInf(n, +1):
		return append(out, `"Infinity"`...)
	case math.IsInf(n, -1):
		return append(out, `"-Infinity"`...)
	}

	
	
	fmt := byte('f')
	if abs := math.Abs(n); abs != 0 {
		if bitSize == 64 && (abs < 1e-6 || abs >= 1e21) ||
			bitSize == 32 && (float32(abs) < 1e-6 || float32(abs) >= 1e21) {
			fmt = 'e'
		}
	}
	out = strconv.AppendFloat(out, n, fmt, -1, bitSize)
	if fmt == 'e' {
		n := len(out)
		if n >= 4 && out[n-4] == 'e' && out[n-3] == '-' && out[n-2] == '0' {
			out[n-2] = out[n-1]
			out = out[:n-1]
		}
	}
	return out
}


func (e *Encoder) WriteInt(n int64) {
	e.prepareNext(scalar)
	e.out = strconv.AppendInt(e.out, n, 10)
}


func (e *Encoder) WriteUint(n uint64) {
	e.prepareNext(scalar)
	e.out = strconv.AppendUint(e.out, n, 10)
}


func (e *Encoder) StartObject() {
	e.prepareNext(objectOpen)
	e.out = append(e.out, '{')
}


func (e *Encoder) EndObject() {
	e.prepareNext(objectClose)
	e.out = append(e.out, '}')
}




func (e *Encoder) WriteName(s string) error {
	e.prepareNext(name)
	var err error
	
	e.out, err = appendString(e.out, s)
	e.out = append(e.out, ':')
	return err
}


func (e *Encoder) StartArray() {
	e.prepareNext(arrayOpen)
	e.out = append(e.out, '[')
}


func (e *Encoder) EndArray() {
	e.prepareNext(arrayClose)
	e.out = append(e.out, ']')
}



func (e *Encoder) prepareNext(next kind) {
	defer func() {
		
		e.lastKind = next
	}()

	if len(e.indent) == 0 {
		
		if e.lastKind&(scalar|objectClose|arrayClose) != 0 &&
			next&(name|scalar|objectOpen|arrayOpen) != 0 {
			e.out = append(e.out, ',')
			
			
			if detrand.Bool() {
				e.out = append(e.out, ' ')
			}
		}
		return
	}

	switch {
	case e.lastKind&(objectOpen|arrayOpen) != 0:
		
		if next&(objectClose|arrayClose) == 0 {
			e.indents = append(e.indents, e.indent...)
			e.out = append(e.out, '\n')
			e.out = append(e.out, e.indents...)
		}

	case e.lastKind&(scalar|objectClose|arrayClose) != 0:
		switch {
		
		case next&(name|scalar|objectOpen|arrayOpen) != 0:
			e.out = append(e.out, ',', '\n')

		
		case next&(objectClose|arrayClose) != 0:
			e.indents = e.indents[:len(e.indents)-len(e.indent)]
			e.out = append(e.out, '\n')
		}
		e.out = append(e.out, e.indents...)

	case e.lastKind&name != 0:
		e.out = append(e.out, ' ')
		
		
		if detrand.Bool() {
			e.out = append(e.out, ' ')
		}
	}
}
