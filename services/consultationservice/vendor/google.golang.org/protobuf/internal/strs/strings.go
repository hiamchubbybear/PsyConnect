




package strs

import (
	"go/token"
	"strings"
	"unicode"
	"unicode/utf8"

	"google.golang.org/protobuf/internal/flags"
	"google.golang.org/protobuf/reflect/protoreflect"
)


func EnforceUTF8(fd protoreflect.FieldDescriptor) bool {
	if flags.ProtoLegacy || fd.Syntax() == protoreflect.Editions {
		if fd, ok := fd.(interface{ EnforceUTF8() bool }); ok {
			return fd.EnforceUTF8()
		}
	}
	return fd.Syntax() == protoreflect.Proto3
}





func GoCamelCase(s string) string {
	
	
	
	
	var b []byte
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '.' && i+1 < len(s) && isASCIILower(s[i+1]):
			
		case c == '.':
			b = append(b, '_') 
		case c == '_' && (i == 0 || s[i-1] == '.'):
			
			
			b = append(b, 'X') 
		case c == '_' && i+1 < len(s) && isASCIILower(s[i+1]):
			
		case isASCIIDigit(c):
			b = append(b, c)
		default:
			
			
			if isASCIILower(c) {
				c -= 'a' - 'A' 
			}
			b = append(b, c)

			
			for ; i+1 < len(s) && isASCIILower(s[i+1]); i++ {
				b = append(b, s[i+1])
			}
		}
	}
	return string(b)
}


func GoSanitized(s string) string {
	
	
	s = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return r
		}
		return '_'
	}, s)

	
	
	r, _ := utf8.DecodeRuneInString(s)
	if token.Lookup(s).IsKeyword() || !unicode.IsLetter(r) {
		return "_" + s
	}
	return s
}



func JSONCamelCase(s string) string {
	var b []byte
	var wasUnderscore bool
	for i := 0; i < len(s); i++ { 
		c := s[i]
		if c != '_' {
			if wasUnderscore && isASCIILower(c) {
				c -= 'a' - 'A' 
			}
			b = append(b, c)
		}
		wasUnderscore = c == '_'
	}
	return string(b)
}



func JSONSnakeCase(s string) string {
	var b []byte
	for i := 0; i < len(s); i++ { 
		c := s[i]
		if isASCIIUpper(c) {
			b = append(b, '_')
			c += 'a' - 'A' 
		}
		b = append(b, c)
	}
	return string(b)
}



func MapEntryName(s string) string {
	var b []byte
	upperNext := true
	for _, c := range s {
		switch {
		case c == '_':
			upperNext = true
		case upperNext:
			b = append(b, byte(unicode.ToUpper(c)))
			upperNext = false
		default:
			b = append(b, byte(c))
		}
	}
	b = append(b, "Entry"...)
	return string(b)
}



func EnumValueName(s string) string {
	var b []byte
	upperNext := true
	for _, c := range s {
		switch {
		case c == '_':
			upperNext = true
		case upperNext:
			b = append(b, byte(unicode.ToUpper(c)))
			upperNext = false
		default:
			b = append(b, byte(unicode.ToLower(c)))
			upperNext = false
		}
	}
	return string(b)
}




func TrimEnumPrefix(s, prefix string) string {
	s0 := s 
	for len(s) > 0 && len(prefix) > 0 {
		if s[0] == '_' {
			s = s[1:]
			continue
		}
		if unicode.ToLower(rune(s[0])) != rune(prefix[0]) {
			return s0 
		}
		s, prefix = s[1:], prefix[1:]
	}
	if len(prefix) > 0 {
		return s0 
	}
	s = strings.TrimLeft(s, "_")
	if len(s) == 0 {
		return s0 
	}
	return s
}

func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}
func isASCIIUpper(c byte) bool {
	return 'A' <= c && c <= 'Z'
}
func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}
