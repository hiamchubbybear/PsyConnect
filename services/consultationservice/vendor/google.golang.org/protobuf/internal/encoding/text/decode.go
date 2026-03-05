



package text

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"unicode/utf8"

	"google.golang.org/protobuf/internal/errors"
)


type Decoder struct {
	
	
	lastCall call

	
	lastToken Token

	
	lastErr error

	
	
	
	
	
	openStack []byte

	
	orig []byte
	
	in []byte
}


func NewDecoder(b []byte) *Decoder {
	return &Decoder{orig: b, in: b}
}


var ErrUnexpectedEOF = errors.New("%v", io.ErrUnexpectedEOF)


type call uint8

const (
	readCall call = iota
	peekCall
)


func (d *Decoder) Peek() (Token, error) {
	defer func() { d.lastCall = peekCall }()
	if d.lastCall == readCall {
		d.lastToken, d.lastErr = d.Read()
	}
	return d.lastToken, d.lastErr
}



func (d *Decoder) Read() (Token, error) {
	defer func() { d.lastCall = readCall }()
	if d.lastCall == peekCall {
		return d.lastToken, d.lastErr
	}

	tok, err := d.parseNext(d.lastToken.kind)
	if err != nil {
		return Token{}, err
	}

	switch tok.kind {
	case comma, semicolon:
		tok, err = d.parseNext(tok.kind)
		if err != nil {
			return Token{}, err
		}
	}
	d.lastToken = tok
	return tok, nil
}

const (
	mismatchedFmt = "mismatched close character %q"
	unexpectedFmt = "unexpected character %q"
)


func (d *Decoder) parseNext(lastKind Kind) (Token, error) {
	
	d.consume(0)
	isEOF := false
	if len(d.in) == 0 {
		isEOF = true
	}

	switch lastKind {
	case EOF:
		return d.consumeToken(EOF, 0, 0), nil

	case bof:
		
		if isEOF {
			return d.consumeToken(EOF, 0, 0), nil
		}
		return d.parseFieldName()

	case Name:
		
		if isEOF {
			return Token{}, ErrUnexpectedEOF
		}
		switch ch := d.in[0]; ch {
		case '{', '<':
			d.pushOpenStack(ch)
			return d.consumeToken(MessageOpen, 1, 0), nil
		case '[':
			d.pushOpenStack(ch)
			return d.consumeToken(ListOpen, 1, 0), nil
		default:
			return d.parseScalar()
		}

	case Scalar:
		openKind, closeCh := d.currentOpenKind()
		switch openKind {
		case bof:
			
			
			if isEOF {
				return d.consumeToken(EOF, 0, 0), nil
			}
			switch d.in[0] {
			case ',':
				return d.consumeToken(comma, 1, 0), nil
			case ';':
				return d.consumeToken(semicolon, 1, 0), nil
			default:
				return d.parseFieldName()
			}

		case MessageOpen:
			
			if isEOF {
				return Token{}, ErrUnexpectedEOF
			}
			switch ch := d.in[0]; ch {
			case closeCh:
				d.popOpenStack()
				return d.consumeToken(MessageClose, 1, 0), nil
			case otherCloseChar[closeCh]:
				return Token{}, d.newSyntaxError(mismatchedFmt, ch)
			case ',':
				return d.consumeToken(comma, 1, 0), nil
			case ';':
				return d.consumeToken(semicolon, 1, 0), nil
			default:
				return d.parseFieldName()
			}

		case ListOpen:
			
			if isEOF {
				return Token{}, ErrUnexpectedEOF
			}
			switch ch := d.in[0]; ch {
			case ']':
				d.popOpenStack()
				return d.consumeToken(ListClose, 1, 0), nil
			case ',':
				return d.consumeToken(comma, 1, 0), nil
			default:
				return Token{}, d.newSyntaxError(unexpectedFmt, ch)
			}
		}

	case MessageOpen:
		
		if isEOF {
			return Token{}, ErrUnexpectedEOF
		}
		_, closeCh := d.currentOpenKind()
		switch ch := d.in[0]; ch {
		case closeCh:
			d.popOpenStack()
			return d.consumeToken(MessageClose, 1, 0), nil
		case otherCloseChar[closeCh]:
			return Token{}, d.newSyntaxError(mismatchedFmt, ch)
		default:
			return d.parseFieldName()
		}

	case MessageClose:
		openKind, closeCh := d.currentOpenKind()
		switch openKind {
		case bof:
			
			
			if isEOF {
				return d.consumeToken(EOF, 0, 0), nil
			}
			switch ch := d.in[0]; ch {
			case ',':
				return d.consumeToken(comma, 1, 0), nil
			case ';':
				return d.consumeToken(semicolon, 1, 0), nil
			default:
				return d.parseFieldName()
			}

		case MessageOpen:
			
			if isEOF {
				return Token{}, ErrUnexpectedEOF
			}
			switch ch := d.in[0]; ch {
			case closeCh:
				d.popOpenStack()
				return d.consumeToken(MessageClose, 1, 0), nil
			case otherCloseChar[closeCh]:
				return Token{}, d.newSyntaxError(mismatchedFmt, ch)
			case ',':
				return d.consumeToken(comma, 1, 0), nil
			case ';':
				return d.consumeToken(semicolon, 1, 0), nil
			default:
				return d.parseFieldName()
			}

		case ListOpen:
			
			if isEOF {
				return Token{}, ErrUnexpectedEOF
			}
			switch ch := d.in[0]; ch {
			case closeCh:
				d.popOpenStack()
				return d.consumeToken(ListClose, 1, 0), nil
			case ',':
				return d.consumeToken(comma, 1, 0), nil
			default:
				return Token{}, d.newSyntaxError(unexpectedFmt, ch)
			}
		}

	case ListOpen:
		
		if isEOF {
			return Token{}, ErrUnexpectedEOF
		}
		switch ch := d.in[0]; ch {
		case ']':
			d.popOpenStack()
			return d.consumeToken(ListClose, 1, 0), nil
		case '{', '<':
			d.pushOpenStack(ch)
			return d.consumeToken(MessageOpen, 1, 0), nil
		default:
			return d.parseScalar()
		}

	case ListClose:
		openKind, closeCh := d.currentOpenKind()
		switch openKind {
		case bof:
			
			
			if isEOF {
				return d.consumeToken(EOF, 0, 0), nil
			}
			switch ch := d.in[0]; ch {
			case ',':
				return d.consumeToken(comma, 1, 0), nil
			case ';':
				return d.consumeToken(semicolon, 1, 0), nil
			default:
				return d.parseFieldName()
			}

		case MessageOpen:
			
			if isEOF {
				return Token{}, ErrUnexpectedEOF
			}
			switch ch := d.in[0]; ch {
			case closeCh:
				d.popOpenStack()
				return d.consumeToken(MessageClose, 1, 0), nil
			case otherCloseChar[closeCh]:
				return Token{}, d.newSyntaxError(mismatchedFmt, ch)
			case ',':
				return d.consumeToken(comma, 1, 0), nil
			case ';':
				return d.consumeToken(semicolon, 1, 0), nil
			default:
				return d.parseFieldName()
			}

		default:
			
		}

	case comma, semicolon:
		openKind, closeCh := d.currentOpenKind()
		switch openKind {
		case bof:
			
			if isEOF {
				return d.consumeToken(EOF, 0, 0), nil
			}
			return d.parseFieldName()

		case MessageOpen:
			
			if isEOF {
				return Token{}, ErrUnexpectedEOF
			}
			switch ch := d.in[0]; ch {
			case closeCh:
				d.popOpenStack()
				return d.consumeToken(MessageClose, 1, 0), nil
			case otherCloseChar[closeCh]:
				return Token{}, d.newSyntaxError(mismatchedFmt, ch)
			default:
				return d.parseFieldName()
			}

		case ListOpen:
			if lastKind == semicolon {
				
				
				
				break
			}
			
			if isEOF {
				return Token{}, ErrUnexpectedEOF
			}
			switch ch := d.in[0]; ch {
			case '{', '<':
				d.pushOpenStack(ch)
				return d.consumeToken(MessageOpen, 1, 0), nil
			default:
				return d.parseScalar()
			}
		}
	}

	line, column := d.Position(len(d.orig) - len(d.in))
	panic(fmt.Sprintf("Decoder.parseNext: bug at handling line %d:%d with lastKind=%v", line, column, lastKind))
}

var otherCloseChar = map[byte]byte{
	'}': '>',
	'>': '}',
}





func (d *Decoder) currentOpenKind() (Kind, byte) {
	if len(d.openStack) == 0 {
		return bof, 0
	}
	openCh := d.openStack[len(d.openStack)-1]
	switch openCh {
	case '{':
		return MessageOpen, '}'
	case '<':
		return MessageOpen, '>'
	case '[':
		return ListOpen, ']'
	}
	panic(fmt.Sprintf("Decoder: openStack contains invalid byte %c", openCh))
}

func (d *Decoder) pushOpenStack(ch byte) {
	d.openStack = append(d.openStack, ch)
}

func (d *Decoder) popOpenStack() {
	d.openStack = d.openStack[:len(d.openStack)-1]
}


func (d *Decoder) parseFieldName() (tok Token, err error) {
	defer func() {
		if err == nil && d.tryConsumeChar(':') {
			tok.attrs |= hasSeparator
		}
	}()

	
	if d.in[0] == '[' {
		return d.parseTypeName()
	}

	
	if size := parseIdent(d.in, false); size > 0 {
		return d.consumeToken(Name, size, uint8(IdentName)), nil
	}

	
	
	if num := parseNumber(d.in); num.size > 0 {
		str := num.string(d.in)
		if !num.neg && num.kind == numDec {
			if _, err := strconv.ParseInt(str, 10, 32); err == nil {
				return d.consumeToken(Name, num.size, uint8(FieldNumber)), nil
			}
		}
		return Token{}, d.newSyntaxError("invalid field number: %s", str)
	}

	return Token{}, d.newSyntaxError("invalid field name: %s", errId(d.in))
}






func (d *Decoder) parseTypeName() (Token, error) {
	startPos := len(d.orig) - len(d.in)
	
	
	s := consume(d.in[1:], 0)
	if len(s) == 0 {
		return Token{}, ErrUnexpectedEOF
	}

	var name []byte
	for len(s) > 0 && isTypeNameChar(s[0]) {
		name = append(name, s[0])
		s = s[1:]
	}
	s = consume(s, 0)

	var closed bool
	for len(s) > 0 && !closed {
		switch {
		case s[0] == ']':
			s = s[1:]
			closed = true

		case s[0] == '/', s[0] == '.':
			if len(name) > 0 && (name[len(name)-1] == '/' || name[len(name)-1] == '.') {
				return Token{}, d.newSyntaxError("invalid type URL/extension field name: %s",
					d.orig[startPos:len(d.orig)-len(s)+1])
			}
			name = append(name, s[0])
			s = s[1:]
			s = consume(s, 0)
			for len(s) > 0 && isTypeNameChar(s[0]) {
				name = append(name, s[0])
				s = s[1:]
			}
			s = consume(s, 0)

		default:
			return Token{}, d.newSyntaxError(
				"invalid type URL/extension field name: %s", d.orig[startPos:len(d.orig)-len(s)+1])
		}
	}

	if !closed {
		return Token{}, ErrUnexpectedEOF
	}

	
	size := len(name)
	if size == 0 || name[0] == '.' || name[size-1] == '.' || name[size-1] == '/' {
		return Token{}, d.newSyntaxError("invalid type URL/extension field name: %s",
			d.orig[startPos:len(d.orig)-len(s)])
	}

	d.in = s
	endPos := len(d.orig) - len(d.in)
	d.consume(0)

	return Token{
		kind:  Name,
		attrs: uint8(TypeName),
		pos:   startPos,
		raw:   d.orig[startPos:endPos],
		str:   string(name),
	}, nil
}

func isTypeNameChar(b byte) bool {
	return (b == '-' || b == '_' ||
		('0' <= b && b <= '9') ||
		('a' <= b && b <= 'z') ||
		('A' <= b && b <= 'Z'))
}

func isWhiteSpace(b byte) bool {
	switch b {
	case ' ', '\n', '\r', '\t':
		return true
	default:
		return false
	}
}





func parseIdent(input []byte, allowNeg bool) int {
	var size int

	s := input
	if len(s) == 0 {
		return 0
	}

	if allowNeg && s[0] == '-' {
		s = s[1:]
		size++
		if len(s) == 0 {
			return 0
		}
	}

	switch {
	case s[0] == '_',
		'a' <= s[0] && s[0] <= 'z',
		'A' <= s[0] && s[0] <= 'Z':
		s = s[1:]
		size++
	default:
		return 0
	}

	for len(s) > 0 && (s[0] == '_' ||
		'a' <= s[0] && s[0] <= 'z' ||
		'A' <= s[0] && s[0] <= 'Z' ||
		'0' <= s[0] && s[0] <= '9') {
		s = s[1:]
		size++
	}

	if len(s) > 0 && !isDelim(s[0]) {
		return 0
	}

	return size
}


func (d *Decoder) parseScalar() (Token, error) {
	if d.in[0] == '"' || d.in[0] == '\'' {
		return d.parseStringValue()
	}

	if tok, ok := d.parseLiteralValue(); ok {
		return tok, nil
	}

	if tok, ok := d.parseNumberValue(); ok {
		return tok, nil
	}

	return Token{}, d.newSyntaxError("invalid scalar value: %s", errId(d.in))
}




func (d *Decoder) parseLiteralValue() (Token, bool) {
	size := parseIdent(d.in, true)
	if size == 0 {
		return Token{}, false
	}
	return d.consumeToken(Scalar, size, literalValue), true
}



func (d *Decoder) consumeToken(kind Kind, size int, attrs uint8) Token {
	
	tok := Token{
		kind:  kind,
		attrs: attrs,
		pos:   len(d.orig) - len(d.in),
		raw:   d.in[:size],
	}
	d.consume(size)
	return tok
}



func (d *Decoder) newSyntaxError(f string, x ...any) error {
	e := errors.New(f, x...)
	line, column := d.Position(len(d.orig) - len(d.in))
	return errors.New("syntax error (line %d:%d): %v", line, column, e)
}



func (d *Decoder) Position(idx int) (line int, column int) {
	b := d.orig[:idx]
	line = bytes.Count(b, []byte("\n")) + 1
	if i := bytes.LastIndexByte(b, '\n'); i >= 0 {
		b = b[i+1:]
	}
	column = utf8.RuneCount(b) + 1 
	return line, column
}

func (d *Decoder) tryConsumeChar(c byte) bool {
	if len(d.in) > 0 && d.in[0] == c {
		d.consume(1)
		return true
	}
	return false
}


func (d *Decoder) consume(n int) {
	d.in = consume(d.in, n)
	return
}


func consume(b []byte, n int) []byte {
	b = b[n:]
	for len(b) > 0 {
		switch b[0] {
		case ' ', '\n', '\r', '\t':
			b = b[1:]
		case '#':
			if i := bytes.IndexByte(b, '\n'); i >= 0 {
				b = b[i+len("\n"):]
			} else {
				b = nil
			}
		default:
			return b
		}
	}
	return b
}



func errId(seq []byte) []byte {
	const maxLen = 32
	for i := 0; i < len(seq); {
		if i > maxLen {
			return append(seq[:i:i], "…"...)
		}
		r, size := utf8.DecodeRune(seq[i:])
		if r > utf8.RuneSelf || (r != '/' && isDelim(byte(r))) {
			if i == 0 {
				
				
				
				i = size
			}
			return seq[:i:i]
		}
		i += size
	}
	
	return seq
}


func isDelim(c byte) bool {
	return !(c == '-' || c == '+' || c == '.' || c == '_' ||
		('a' <= c && c <= 'z') ||
		('A' <= c && c <= 'Z') ||
		('0' <= c && c <= '9'))
}
