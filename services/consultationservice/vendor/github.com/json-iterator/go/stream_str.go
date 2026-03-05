package jsoniter

import (
	"unicode/utf8"
)








var htmlSafeSet = [utf8.RuneSelf]bool{
	' ':      true,
	'!':      true,
	'"':      false,
	'#':      true,
	'$':      true,
	'%':      true,
	'&':      false,
	'\'':     true,
	'(':      true,
	')':      true,
	'*':      true,
	'+':      true,
	',':      true,
	'-':      true,
	'.':      true,
	'/':      true,
	'0':      true,
	'1':      true,
	'2':      true,
	'3':      true,
	'4':      true,
	'5':      true,
	'6':      true,
	'7':      true,
	'8':      true,
	'9':      true,
	':':      true,
	';':      true,
	'<':      false,
	'=':      true,
	'>':      false,
	'?':      true,
	'@':      true,
	'A':      true,
	'B':      true,
	'C':      true,
	'D':      true,
	'E':      true,
	'F':      true,
	'G':      true,
	'H':      true,
	'I':      true,
	'J':      true,
	'K':      true,
	'L':      true,
	'M':      true,
	'N':      true,
	'O':      true,
	'P':      true,
	'Q':      true,
	'R':      true,
	'S':      true,
	'T':      true,
	'U':      true,
	'V':      true,
	'W':      true,
	'X':      true,
	'Y':      true,
	'Z':      true,
	'[':      true,
	'\\':     false,
	']':      true,
	'^':      true,
	'_':      true,
	'`':      true,
	'a':      true,
	'b':      true,
	'c':      true,
	'd':      true,
	'e':      true,
	'f':      true,
	'g':      true,
	'h':      true,
	'i':      true,
	'j':      true,
	'k':      true,
	'l':      true,
	'm':      true,
	'n':      true,
	'o':      true,
	'p':      true,
	'q':      true,
	'r':      true,
	's':      true,
	't':      true,
	'u':      true,
	'v':      true,
	'w':      true,
	'x':      true,
	'y':      true,
	'z':      true,
	'{':      true,
	'|':      true,
	'}':      true,
	'~':      true,
	'\u007f': true,
}







var safeSet = [utf8.RuneSelf]bool{
	' ':      true,
	'!':      true,
	'"':      false,
	'#':      true,
	'$':      true,
	'%':      true,
	'&':      true,
	'\'':     true,
	'(':      true,
	')':      true,
	'*':      true,
	'+':      true,
	',':      true,
	'-':      true,
	'.':      true,
	'/':      true,
	'0':      true,
	'1':      true,
	'2':      true,
	'3':      true,
	'4':      true,
	'5':      true,
	'6':      true,
	'7':      true,
	'8':      true,
	'9':      true,
	':':      true,
	';':      true,
	'<':      true,
	'=':      true,
	'>':      true,
	'?':      true,
	'@':      true,
	'A':      true,
	'B':      true,
	'C':      true,
	'D':      true,
	'E':      true,
	'F':      true,
	'G':      true,
	'H':      true,
	'I':      true,
	'J':      true,
	'K':      true,
	'L':      true,
	'M':      true,
	'N':      true,
	'O':      true,
	'P':      true,
	'Q':      true,
	'R':      true,
	'S':      true,
	'T':      true,
	'U':      true,
	'V':      true,
	'W':      true,
	'X':      true,
	'Y':      true,
	'Z':      true,
	'[':      true,
	'\\':     false,
	']':      true,
	'^':      true,
	'_':      true,
	'`':      true,
	'a':      true,
	'b':      true,
	'c':      true,
	'd':      true,
	'e':      true,
	'f':      true,
	'g':      true,
	'h':      true,
	'i':      true,
	'j':      true,
	'k':      true,
	'l':      true,
	'm':      true,
	'n':      true,
	'o':      true,
	'p':      true,
	'q':      true,
	'r':      true,
	's':      true,
	't':      true,
	'u':      true,
	'v':      true,
	'w':      true,
	'x':      true,
	'y':      true,
	'z':      true,
	'{':      true,
	'|':      true,
	'}':      true,
	'~':      true,
	'\u007f': true,
}

var hex = "0123456789abcdef"


func (stream *Stream) WriteStringWithHTMLEscaped(s string) {
	valLen := len(s)
	stream.buf = append(stream.buf, '"')
	
	i := 0
	for ; i < valLen; i++ {
		c := s[i]
		if c < utf8.RuneSelf && htmlSafeSet[c] {
			stream.buf = append(stream.buf, c)
		} else {
			break
		}
	}
	if i == valLen {
		stream.buf = append(stream.buf, '"')
		return
	}
	writeStringSlowPathWithHTMLEscaped(stream, i, s, valLen)
}

func writeStringSlowPathWithHTMLEscaped(stream *Stream, i int, s string, valLen int) {
	start := i
	
	for i < valLen {
		if b := s[i]; b < utf8.RuneSelf {
			if htmlSafeSet[b] {
				i++
				continue
			}
			if start < i {
				stream.WriteRaw(s[start:i])
			}
			switch b {
			case '\\', '"':
				stream.writeTwoBytes('\\', b)
			case '\n':
				stream.writeTwoBytes('\\', 'n')
			case '\r':
				stream.writeTwoBytes('\\', 'r')
			case '\t':
				stream.writeTwoBytes('\\', 't')
			default:
				
				
				
				
				
				stream.WriteRaw(`\u00`)
				stream.writeTwoBytes(hex[b>>4], hex[b&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				stream.WriteRaw(s[start:i])
			}
			stream.WriteRaw(`\ufffd`)
			i++
			start = i
			continue
		}
		
		
		
		
		
		
		
		if c == '\u2028' || c == '\u2029' {
			if start < i {
				stream.WriteRaw(s[start:i])
			}
			stream.WriteRaw(`\u202`)
			stream.writeByte(hex[c&0xF])
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(s) {
		stream.WriteRaw(s[start:])
	}
	stream.writeByte('"')
}


func (stream *Stream) WriteString(s string) {
	valLen := len(s)
	stream.buf = append(stream.buf, '"')
	
	i := 0
	for ; i < valLen; i++ {
		c := s[i]
		if c > 31 && c != '"' && c != '\\' {
			stream.buf = append(stream.buf, c)
		} else {
			break
		}
	}
	if i == valLen {
		stream.buf = append(stream.buf, '"')
		return
	}
	writeStringSlowPath(stream, i, s, valLen)
}

func writeStringSlowPath(stream *Stream, i int, s string, valLen int) {
	start := i
	
	for i < valLen {
		if b := s[i]; b < utf8.RuneSelf {
			if safeSet[b] {
				i++
				continue
			}
			if start < i {
				stream.WriteRaw(s[start:i])
			}
			switch b {
			case '\\', '"':
				stream.writeTwoBytes('\\', b)
			case '\n':
				stream.writeTwoBytes('\\', 'n')
			case '\r':
				stream.writeTwoBytes('\\', 'r')
			case '\t':
				stream.writeTwoBytes('\\', 't')
			default:
				
				
				
				
				
				stream.WriteRaw(`\u00`)
				stream.writeTwoBytes(hex[b>>4], hex[b&0xF])
			}
			i++
			start = i
			continue
		}
		i++
		continue
	}
	if start < len(s) {
		stream.WriteRaw(s[start:])
	}
	stream.writeByte('"')
}
