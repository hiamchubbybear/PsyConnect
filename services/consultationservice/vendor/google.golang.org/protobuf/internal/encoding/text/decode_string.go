



package text

import (
	"bytes"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"

	"google.golang.org/protobuf/internal/strs"
)







func (d *Decoder) parseStringValue() (Token, error) {
	
	
	
	
	
	in0 := d.in
	var ss []string
	for len(d.in) > 0 && (d.in[0] == '"' || d.in[0] == '\'') {
		s, err := d.parseString()
		if err != nil {
			return Token{}, err
		}
		ss = append(ss, s)
	}
	
	return Token{
		kind:  Scalar,
		attrs: stringValue,
		pos:   len(d.orig) - len(in0),
		raw:   in0[:len(in0)-len(d.in)],
		str:   strings.Join(ss, ""),
	}, nil
}


func (d *Decoder) parseString() (string, error) {
	in := d.in
	if len(in) == 0 {
		return "", ErrUnexpectedEOF
	}
	quote := in[0]
	in = in[1:]
	i := indexNeedEscapeInBytes(in)
	in, out := in[i:], in[:i:i] 
	for len(in) > 0 {
		switch r, n := utf8.DecodeRune(in); {
		case r == utf8.RuneError && n == 1:
			return "", d.newSyntaxError("invalid UTF-8 detected")
		case r == 0 || r == '\n':
			return "", d.newSyntaxError("invalid character %q in string", r)
		case r == rune(quote):
			in = in[1:]
			d.consume(len(d.in) - len(in))
			return string(out), nil
		case r == '\\':
			if len(in) < 2 {
				return "", ErrUnexpectedEOF
			}
			switch r := in[1]; r {
			case '"', '\'', '\\', '?':
				in, out = in[2:], append(out, r)
			case 'a':
				in, out = in[2:], append(out, '\a')
			case 'b':
				in, out = in[2:], append(out, '\b')
			case 'n':
				in, out = in[2:], append(out, '\n')
			case 'r':
				in, out = in[2:], append(out, '\r')
			case 't':
				in, out = in[2:], append(out, '\t')
			case 'v':
				in, out = in[2:], append(out, '\v')
			case 'f':
				in, out = in[2:], append(out, '\f')
			case '0', '1', '2', '3', '4', '5', '6', '7':
				
				n := len(in[1:]) - len(bytes.TrimLeft(in[1:], "01234567"))
				if n > 3 {
					n = 3
				}
				v, err := strconv.ParseUint(string(in[1:1+n]), 8, 8)
				if err != nil {
					return "", d.newSyntaxError("invalid octal escape code %q in string", in[:1+n])
				}
				in, out = in[1+n:], append(out, byte(v))
			case 'x':
				
				n := len(in[2:]) - len(bytes.TrimLeft(in[2:], "0123456789abcdefABCDEF"))
				if n > 2 {
					n = 2
				}
				v, err := strconv.ParseUint(string(in[2:2+n]), 16, 8)
				if err != nil {
					return "", d.newSyntaxError("invalid hex escape code %q in string", in[:2+n])
				}
				in, out = in[2+n:], append(out, byte(v))
			case 'u', 'U':
				
				n := 6
				if r == 'U' {
					n = 10
				}
				if len(in) < n {
					return "", ErrUnexpectedEOF
				}
				v, err := strconv.ParseUint(string(in[2:n]), 16, 32)
				if utf8.MaxRune < v || err != nil {
					return "", d.newSyntaxError("invalid Unicode escape code %q in string", in[:n])
				}
				in = in[n:]

				r := rune(v)
				if utf16.IsSurrogate(r) {
					if len(in) < 6 {
						return "", ErrUnexpectedEOF
					}
					v, err := strconv.ParseUint(string(in[2:6]), 16, 16)
					r = utf16.DecodeRune(r, rune(v))
					if in[0] != '\\' || in[1] != 'u' || r == unicode.ReplacementChar || err != nil {
						return "", d.newSyntaxError("invalid Unicode escape code %q in string", in[:6])
					}
					in = in[6:]
				}
				out = append(out, string(r)...)
			default:
				return "", d.newSyntaxError("invalid escape code %q in string", in[:2])
			}
		default:
			i := indexNeedEscapeInBytes(in[n:])
			in, out = in[n+i:], append(out, in[:n+i]...)
		}
	}
	return "", ErrUnexpectedEOF
}



func indexNeedEscapeInBytes(b []byte) int { return indexNeedEscapeInString(strs.UnsafeString(b)) }




func UnmarshalString(s string) (string, error) {
	d := NewDecoder([]byte(s))
	return d.parseString()
}
