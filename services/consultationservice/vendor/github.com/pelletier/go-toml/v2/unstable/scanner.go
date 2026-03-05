package unstable

import "github.com/pelletier/go-toml/v2/internal/characters"

func scanFollows(b []byte, pattern string) bool {
	n := len(pattern)

	return len(b) >= n && string(b[:n]) == pattern
}

func scanFollowsMultilineBasicStringDelimiter(b []byte) bool {
	return scanFollows(b, `"""`)
}

func scanFollowsMultilineLiteralStringDelimiter(b []byte) bool {
	return scanFollows(b, `'''`)
}

func scanFollowsTrue(b []byte) bool {
	return scanFollows(b, `true`)
}

func scanFollowsFalse(b []byte) bool {
	return scanFollows(b, `false`)
}

func scanFollowsInf(b []byte) bool {
	return scanFollows(b, `inf`)
}

func scanFollowsNan(b []byte) bool {
	return scanFollows(b, `nan`)
}

func scanUnquotedKey(b []byte) ([]byte, []byte) {
	
	for i := 0; i < len(b); i++ {
		if !isUnquotedKeyChar(b[i]) {
			return b[:i], b[i:]
		}
	}

	return b, b[len(b):]
}

func isUnquotedKeyChar(r byte) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_'
}

func scanLiteralString(b []byte) ([]byte, []byte, error) {
	
	
	
	for i := 1; i < len(b); {
		switch b[i] {
		case '\'':
			return b[:i+1], b[i+1:], nil
		case '\n', '\r':
			return nil, nil, NewParserError(b[i:i+1], "literal strings cannot have new lines")
		}
		size := characters.Utf8ValidNext(b[i:])
		if size == 0 {
			return nil, nil, NewParserError(b[i:i+1], "invalid character")
		}
		i += size
	}

	return nil, nil, NewParserError(b[len(b):], "unterminated literal string")
}

func scanMultilineLiteralString(b []byte) ([]byte, []byte, error) {
	
	
	
	
	
	
	
	
	for i := 3; i < len(b); {
		switch b[i] {
		case '\'':
			if scanFollowsMultilineLiteralStringDelimiter(b[i:]) {
				i += 3

				
				
				
				
				

				if i >= len(b) || b[i] != '\'' {
					return b[:i], b[i:], nil
				}
				i++

				if i >= len(b) || b[i] != '\'' {
					return b[:i], b[i:], nil
				}
				i++

				if i < len(b) && b[i] == '\'' {
					return nil, nil, NewParserError(b[i-3:i+1], "''' not allowed in multiline literal string")
				}

				return b[:i], b[i:], nil
			}
		case '\r':
			if len(b) < i+2 {
				return nil, nil, NewParserError(b[len(b):], `need a \n after \r`)
			}
			if b[i+1] != '\n' {
				return nil, nil, NewParserError(b[i:i+2], `need a \n after \r`)
			}
			i += 2 
			continue
		}
		size := characters.Utf8ValidNext(b[i:])
		if size == 0 {
			return nil, nil, NewParserError(b[i:i+1], "invalid character")
		}
		i += size
	}

	return nil, nil, NewParserError(b[len(b):], `multiline literal string not terminated by '''`)
}

func scanWindowsNewline(b []byte) ([]byte, []byte, error) {
	const lenCRLF = 2
	if len(b) < lenCRLF {
		return nil, nil, NewParserError(b, "windows new line expected")
	}

	if b[1] != '\n' {
		return nil, nil, NewParserError(b, `windows new line should be \r\n`)
	}

	return b[:lenCRLF], b[lenCRLF:], nil
}

func scanWhitespace(b []byte) ([]byte, []byte) {
	for i := 0; i < len(b); i++ {
		switch b[i] {
		case ' ', '\t':
			continue
		default:
			return b[:i], b[i:]
		}
	}

	return b, b[len(b):]
}

func scanComment(b []byte) ([]byte, []byte, error) {
	
	
	
	
	

	for i := 1; i < len(b); {
		if b[i] == '\n' {
			return b[:i], b[i:], nil
		}
		if b[i] == '\r' {
			if i+1 < len(b) && b[i+1] == '\n' {
				return b[:i+1], b[i+1:], nil
			}
			return nil, nil, NewParserError(b[i:i+1], "invalid character in comment")
		}
		size := characters.Utf8ValidNext(b[i:])
		if size == 0 {
			return nil, nil, NewParserError(b[i:i+1], "invalid character in comment")
		}

		i += size
	}

	return b, b[len(b):], nil
}

func scanBasicString(b []byte) ([]byte, bool, []byte, error) {
	
	
	
	
	
	escaped := false
	i := 1

	for ; i < len(b); i++ {
		switch b[i] {
		case '"':
			return b[:i+1], escaped, b[i+1:], nil
		case '\n', '\r':
			return nil, escaped, nil, NewParserError(b[i:i+1], "basic strings cannot have new lines")
		case '\\':
			if len(b) < i+2 {
				return nil, escaped, nil, NewParserError(b[i:i+1], "need a character after \\")
			}
			escaped = true
			i++ 
		}
	}

	return nil, escaped, nil, NewParserError(b[len(b):], `basic string not terminated by "`)
}

func scanMultilineBasicString(b []byte) ([]byte, bool, []byte, error) {
	
	
	
	
	
	
	
	
	
	

	escaped := false
	i := 3

	for ; i < len(b); i++ {
		switch b[i] {
		case '"':
			if scanFollowsMultilineBasicStringDelimiter(b[i:]) {
				i += 3

				
				
				
				
				

				if i >= len(b) || b[i] != '"' {
					return b[:i], escaped, b[i:], nil
				}
				i++

				if i >= len(b) || b[i] != '"' {
					return b[:i], escaped, b[i:], nil
				}
				i++

				if i < len(b) && b[i] == '"' {
					return nil, escaped, nil, NewParserError(b[i-3:i+1], `""" not allowed in multiline basic string`)
				}

				return b[:i], escaped, b[i:], nil
			}
		case '\\':
			if len(b) < i+2 {
				return nil, escaped, nil, NewParserError(b[len(b):], "need a character after \\")
			}
			escaped = true
			i++ 
		case '\r':
			if len(b) < i+2 {
				return nil, escaped, nil, NewParserError(b[len(b):], `need a \n after \r`)
			}
			if b[i+1] != '\n' {
				return nil, escaped, nil, NewParserError(b[i:i+2], `need a \n after \r`)
			}
			i++ 
		}
	}

	return nil, escaped, nil, NewParserError(b[len(b):], `multiline basic string not terminated by """`)
}
