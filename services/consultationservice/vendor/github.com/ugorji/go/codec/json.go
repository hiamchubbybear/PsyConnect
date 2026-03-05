


package codec














import (
	"encoding/base64"
	"math"
	"reflect"
	"strconv"
	"time"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
)








const jsonLits = `"true"false"null"`

var jsonLitb = []byte(jsonLits)

const (
	jsonLitT = 1
	jsonLitF = 6
	jsonLitN = 12
)

const jsonEncodeUintSmallsString = "" +
	"00010203040506070809" +
	"10111213141516171819" +
	"20212223242526272829" +
	"30313233343536373839" +
	"40414243444546474849" +
	"50515253545556575859" +
	"60616263646566676869" +
	"70717273747576777879" +
	"80818283848586878889" +
	"90919293949596979899"

var jsonEncodeUintSmallsStringBytes = []byte(jsonEncodeUintSmallsString)

const (
	jsonU4Chk2 = '0'
	jsonU4Chk1 = 'a' - 10
	jsonU4Chk0 = 'A' - 10
)

const (
	
	
	
	
	
	
	jsonValidateSymbols = true

	
	
	
	
	
	jsonEscapeMultiByteUnicodeSep = true

	
	
	
	jsonNakedBoolNullInQuotedStr = true

	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	jsonManualInlineDecRdInHotZones = true

	jsonSpacesOrTabsLen = 128

	
)

var (
	
	jsonTabs, jsonSpaces [jsonSpacesOrTabsLen]byte

	jsonCharHtmlSafeSet bitset256
	jsonCharSafeSet     bitset256
)

func init() {
	var i byte
	for i = 0; i < jsonSpacesOrTabsLen; i++ {
		jsonSpaces[i] = ' '
		jsonTabs[i] = '\t'
	}

	
	
	
	for i = 32; i < utf8.RuneSelf; i++ {
		switch i {
		case '"', '\\':
		case '<', '>', '&':
			jsonCharSafeSet.set(i) 
		default:
			jsonCharSafeSet.set(i)
			jsonCharHtmlSafeSet.set(i)
		}
	}
}



type jsonEncState struct {
	di int8   
	d  bool   
	dl uint16 
}

func (x jsonEncState) captureState() interface{}   { return x }
func (x *jsonEncState) restoreState(v interface{}) { *x = v.(jsonEncState) }

type jsonEncDriver struct {
	noBuiltInTypes
	h *JsonHandle

	

	
	jsonEncState

	ks bool 
	is byte 

	typical bool
	rawext  bool 

	s *bitset256 

	

	
	
	
	
	
	
	
	
	
	
	
	b [48]byte

	e Encoder
}

func (e *jsonEncDriver) encoder() *Encoder { return &e.e }

func (e *jsonEncDriver) writeIndent() {
	e.e.encWr.writen1('\n')
	x := int(e.di) * int(e.dl)
	if e.di < 0 {
		x = -x
		for x > jsonSpacesOrTabsLen {
			e.e.encWr.writeb(jsonTabs[:])
			x -= jsonSpacesOrTabsLen
		}
		e.e.encWr.writeb(jsonTabs[:x])
	} else {
		for x > jsonSpacesOrTabsLen {
			e.e.encWr.writeb(jsonSpaces[:])
			x -= jsonSpacesOrTabsLen
		}
		e.e.encWr.writeb(jsonSpaces[:x])
	}
}

func (e *jsonEncDriver) WriteArrayElem() {
	if e.e.c != containerArrayStart {
		e.e.encWr.writen1(',')
	}
	if e.d {
		e.writeIndent()
	}
}

func (e *jsonEncDriver) WriteMapElemKey() {
	if e.e.c != containerMapStart {
		e.e.encWr.writen1(',')
	}
	if e.d {
		e.writeIndent()
	}
}

func (e *jsonEncDriver) WriteMapElemValue() {
	if e.d {
		e.e.encWr.writen2(':', ' ')
	} else {
		e.e.encWr.writen1(':')
	}
}

func (e *jsonEncDriver) EncodeNil() {
	
	

	e.e.encWr.writestr(jsonLits[jsonLitN : jsonLitN+4])
}

func (e *jsonEncDriver) EncodeTime(t time.Time) {
	
	

	if t.IsZero() {
		e.EncodeNil()
	} else {
		e.b[0] = '"'
		b := fmtTime(t, time.RFC3339Nano, e.b[1:1])
		e.b[len(b)+1] = '"'
		e.e.encWr.writeb(e.b[:len(b)+2])
	}
}

func (e *jsonEncDriver) EncodeExt(rv interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
	if ext == SelfExt {
		e.e.encodeValue(baseRV(rv), e.h.fnNoExt(basetype))
	} else if v := ext.ConvertExt(rv); v == nil {
		e.EncodeNil()
	} else {
		e.e.encode(v)
	}
}

func (e *jsonEncDriver) EncodeRawExt(re *RawExt) {
	
	if re.Value == nil {
		e.EncodeNil()
	} else {
		e.e.encode(re.Value)
	}
}

var jsonEncBoolStrs = [2][2]string{
	{jsonLits[jsonLitF : jsonLitF+5], jsonLits[jsonLitT : jsonLitT+4]},
	{jsonLits[jsonLitF-1 : jsonLitF+6], jsonLits[jsonLitT-1 : jsonLitT+5]},
}

func (e *jsonEncDriver) EncodeBool(b bool) {
	e.e.encWr.writestr(
		jsonEncBoolStrs[bool2int(e.ks && e.e.c == containerMapKey)%2][bool2int(b)%2])
}

















func (e *jsonEncDriver) encodeFloat(f float64, bitsize, fmt byte, prec int8) {
	var blen uint
	if e.ks && e.e.c == containerMapKey {
		blen = 2 + uint(len(strconv.AppendFloat(e.b[1:1], f, fmt, int(prec), int(bitsize))))
		
		e.b[0] = '"'
		e.b[blen-1] = '"'
		e.e.encWr.writeb(e.b[:blen])
	} else {
		e.e.encWr.writeb(strconv.AppendFloat(e.b[:0], f, fmt, int(prec), int(bitsize)))
	}
}

func (e *jsonEncDriver) EncodeFloat64(f float64) {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		e.EncodeNil()
		return
	}
	fmt, prec := jsonFloatStrconvFmtPrec64(f)
	e.encodeFloat(f, 64, fmt, prec)
}

func (e *jsonEncDriver) EncodeFloat32(f float32) {
	if math.IsNaN(float64(f)) || math.IsInf(float64(f), 0) {
		e.EncodeNil()
		return
	}
	fmt, prec := jsonFloatStrconvFmtPrec32(f)
	e.encodeFloat(float64(f), 32, fmt, prec)
}

func (e *jsonEncDriver) encodeUint(neg bool, quotes bool, u uint64) {
	
	

	
	var ss = jsonEncodeUintSmallsStringBytes

	
	
	var a = e.b[0:24]
	var i = uint(len(a))

	if quotes {
		i--
		setByteAt(a, i, '"')
		
	}
	
	var is uint
	var us = uint(u)
	for us >= 100 {
		is = us % 100 * 2
		us /= 100
		i -= 2
		setByteAt(a, i+1, byteAt(ss, is+1))
		setByteAt(a, i, byteAt(ss, is))
		
		
	}

	
	is = us * 2
	i--
	setByteAt(a, i, byteAt(ss, is+1))
	
	if us >= 10 {
		i--
		setByteAt(a, i, byteAt(ss, is))
		
	}
	if neg {
		i--
		setByteAt(a, i, '-')
		
	}
	if quotes {
		i--
		setByteAt(a, i, '"')
		
	}
	e.e.encWr.writeb(a[i:])
}

func (e *jsonEncDriver) EncodeInt(v int64) {
	quotes := e.is == 'A' || e.is == 'L' && (v > 1<<53 || v < -(1<<53)) ||
		(e.ks && e.e.c == containerMapKey)

	if cpu32Bit {
		if quotes {
			blen := 2 + len(strconv.AppendInt(e.b[1:1], v, 10))
			e.b[0] = '"'
			e.b[blen-1] = '"'
			e.e.encWr.writeb(e.b[:blen])
		} else {
			e.e.encWr.writeb(strconv.AppendInt(e.b[:0], v, 10))
		}
		return
	}

	if v < 0 {
		e.encodeUint(true, quotes, uint64(-v))
	} else {
		e.encodeUint(false, quotes, uint64(v))
	}
}

func (e *jsonEncDriver) EncodeUint(v uint64) {
	quotes := e.is == 'A' || e.is == 'L' && v > 1<<53 ||
		(e.ks && e.e.c == containerMapKey)

	if cpu32Bit {
		
		if quotes {
			blen := 2 + len(strconv.AppendUint(e.b[1:1], v, 10))
			e.b[0] = '"'
			e.b[blen-1] = '"'
			e.e.encWr.writeb(e.b[:blen])
		} else {
			e.e.encWr.writeb(strconv.AppendUint(e.b[:0], v, 10))
		}
		return
	}

	e.encodeUint(false, quotes, v)
}

func (e *jsonEncDriver) EncodeString(v string) {
	if e.h.StringToRaw {
		e.EncodeStringBytesRaw(bytesView(v))
		return
	}
	e.quoteStr(v)
}

func (e *jsonEncDriver) EncodeStringBytesRaw(v []byte) {
	
	if v == nil {
		e.EncodeNil()
		return
	}

	if e.rawext {
		iv := e.h.RawBytesExt.ConvertExt(v)
		if iv == nil {
			e.EncodeNil()
		} else {
			e.e.encode(iv)
		}
		return
	}

	slen := base64.StdEncoding.EncodedLen(len(v)) + 2

	
	

	bs := e.e.blist.peek(slen, false)
	bs = bs[:slen]

	base64.StdEncoding.Encode(bs[1:], v)
	bs[len(bs)-1] = '"'
	bs[0] = '"'
	e.e.encWr.writeb(bs)
}






func (e *jsonEncDriver) WriteArrayStart(length int) {
	if e.d {
		e.dl++
	}
	e.e.encWr.writen1('[')
}

func (e *jsonEncDriver) WriteArrayEnd() {
	if e.d {
		e.dl--
		e.writeIndent()
	}
	e.e.encWr.writen1(']')
}

func (e *jsonEncDriver) WriteMapStart(length int) {
	if e.d {
		e.dl++
	}
	e.e.encWr.writen1('{')
}

func (e *jsonEncDriver) WriteMapEnd() {
	if e.d {
		e.dl--
		if e.e.c != containerMapStart {
			e.writeIndent()
		}
	}
	e.e.encWr.writen1('}')
}

func (e *jsonEncDriver) quoteStr(s string) {
	
	const hex = "0123456789abcdef"
	w := e.e.w()
	w.writen1('"')
	var i, start uint
	for i < uint(len(s)) {
		
		

		
		

		
		
		if e.s.isset(s[i]) {
			i++
			continue
		}
		
		if s[i] < utf8.RuneSelf {
			if start < i {
				w.writestr(s[start:i])
			}
			switch s[i] {
			case '\\', '"':
				w.writen2('\\', s[i])
			case '\n':
				w.writen2('\\', 'n')
			case '\r':
				w.writen2('\\', 'r')
			case '\b':
				w.writen2('\\', 'b')
			case '\f':
				w.writen2('\\', 'f')
			case '\t':
				w.writen2('\\', 't')
			default:
				w.writestr(`\u00`)
				w.writen2(hex[s[i]>>4], hex[s[i]&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError && size == 1 { 
			if start < i {
				w.writestr(s[start:i])
			}
			w.writestr(`\uFFFD`)
			i++
			start = i
			continue
		}
		
		
		if jsonEscapeMultiByteUnicodeSep && (c == '\u2028' || c == '\u2029') {
			if start < i {
				w.writestr(s[start:i])
			}
			w.writestr(`\u202`)
			w.writen1(hex[c&0xF])
			i += uint(size)
			start = i
			continue
		}
		i += uint(size)
	}
	if start < uint(len(s)) {
		w.writestr(s[start:])
	}
	w.writen1('"')
}

func (e *jsonEncDriver) atEndOfEncode() {
	if e.h.TermWhitespace {
		var c byte = ' ' 
		if e.e.c != 0 {
			c = '\n' 
		}
		e.e.encWr.writen1(c)
	}
}



type jsonDecState struct {
	rawext bool 

	tok  uint8   
	_    bool    
	_    byte    
	bstr [4]byte 

	
	
	buf *[]byte
}

func (x jsonDecState) captureState() interface{}   { return x }
func (x *jsonDecState) restoreState(v interface{}) { *x = v.(jsonDecState) }

type jsonDecDriver struct {
	noBuiltInTypes
	decDriverNoopNumberHelper
	h *JsonHandle

	jsonDecState

	

	

	d Decoder
}

func (d *jsonDecDriver) descBd() (s string) { panic("descBd unsupported") }

func (d *jsonDecDriver) decoder() *Decoder {
	return &d.d
}

func (d *jsonDecDriver) ReadMapStart() int {
	d.advance()
	if d.tok == 'n' {
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.d.decRd.readn3())
		return containerLenNil
	}
	if d.tok != '{' {
		d.d.errorf("read map - expect char '%c' but got char '%c'", '{', d.tok)
	}
	d.tok = 0
	return containerLenUnknown
}

func (d *jsonDecDriver) ReadArrayStart() int {
	d.advance()
	if d.tok == 'n' {
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.d.decRd.readn3())
		return containerLenNil
	}
	if d.tok != '[' {
		d.d.errorf("read array - expect char '%c' but got char '%c'", '[', d.tok)
	}
	d.tok = 0
	return containerLenUnknown
}







func (d *jsonDecDriver) CheckBreak() bool {
	d.advance()
	return d.tok == '}' || d.tok == ']'
}

func (d *jsonDecDriver) ReadArrayElem() {
	const xc uint8 = ','
	if d.d.c != containerArrayStart {
		d.advance()
		if d.tok != xc {
			d.readDelimError(xc)
		}
		d.tok = 0
	}
}

func (d *jsonDecDriver) ReadArrayEnd() {
	const xc uint8 = ']'
	d.advance()
	if d.tok != xc {
		d.readDelimError(xc)
	}
	d.tok = 0
}

func (d *jsonDecDriver) ReadMapElemKey() {
	const xc uint8 = ','
	if d.d.c != containerMapStart {
		d.advance()
		if d.tok != xc {
			d.readDelimError(xc)
		}
		d.tok = 0
	}
}

func (d *jsonDecDriver) ReadMapElemValue() {
	const xc uint8 = ':'
	d.advance()
	if d.tok != xc {
		d.readDelimError(xc)
	}
	d.tok = 0
}

func (d *jsonDecDriver) ReadMapEnd() {
	const xc uint8 = '}'
	d.advance()
	if d.tok != xc {
		d.readDelimError(xc)
	}
	d.tok = 0
}

func (d *jsonDecDriver) readDelimError(xc uint8) {
	d.d.errorf("read json delimiter - expect char '%c' but got char '%c'", xc, d.tok)
}





func (d *jsonDecDriver) checkLit3(got, expect [3]byte) {
	d.tok = 0
	if jsonValidateSymbols && got != expect {
		d.d.errorf("expecting %s: got %s", expect, got)
	}
}

func (d *jsonDecDriver) checkLit4(got, expect [4]byte) {
	d.tok = 0
	if jsonValidateSymbols && got != expect {
		d.d.errorf("expecting %s: got %s", expect, got)
	}
}

func (d *jsonDecDriver) skipWhitespace() {
	d.tok = d.d.decRd.skipWhitespace()
}

func (d *jsonDecDriver) advance() {
	if d.tok == 0 {
		d.skipWhitespace()
	}
}

func (d *jsonDecDriver) nextValueBytes(v []byte) []byte {
	v, cursor := d.nextValueBytesR(v)
	decNextValueBytesHelper{d: &d.d}.bytesRdV(&v, cursor)
	return v
}

func (d *jsonDecDriver) nextValueBytesR(v0 []byte) (v []byte, cursor uint) {
	v = v0
	var h = decNextValueBytesHelper{d: &d.d}
	dr := &d.d.decRd

	consumeString := func() {
	TOP:
		bs := dr.jsonReadAsisChars()
		h.appendN(&v, bs...)
		if bs[len(bs)-1] != '"' {
			
			h.append1(&v, dr.readn1())
			goto TOP
		}
	}

	d.advance()           
	cursor = d.d.rb.c - 1 

	switch d.tok {
	default:
		h.appendN(&v, dr.jsonReadNum()...)
	case 'n':
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.d.decRd.readn3())
		h.appendS(&v, jsonLits[jsonLitN:jsonLitN+4])
	case 'f':
		d.checkLit4([4]byte{'a', 'l', 's', 'e'}, d.d.decRd.readn4())
		h.appendS(&v, jsonLits[jsonLitF:jsonLitF+5])
	case 't':
		d.checkLit3([3]byte{'r', 'u', 'e'}, d.d.decRd.readn3())
		h.appendS(&v, jsonLits[jsonLitT:jsonLitT+4])
	case '"':
		h.append1(&v, '"')
		consumeString()
	case '{', '[':
		var elem struct{}
		var stack []struct{}

		stack = append(stack, elem)

		h.append1(&v, d.tok)

		for len(stack) != 0 {
			c := dr.readn1()
			h.append1(&v, c)
			switch c {
			case '"':
				consumeString()
			case '{', '[':
				stack = append(stack, elem)
			case '}', ']':
				stack = stack[:len(stack)-1]
			}
		}
	}
	d.tok = 0
	return
}

func (d *jsonDecDriver) TryNil() bool {
	d.advance()
	
	
	if d.tok == 'n' {
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.d.decRd.readn3())
		return true
	}
	return false
}

func (d *jsonDecDriver) DecodeBool() (v bool) {
	d.advance()
	
	fquot := d.d.c == containerMapKey && d.tok == '"'
	if fquot {
		d.tok = d.d.decRd.readn1()
	}
	switch d.tok {
	case 'f':
		d.checkLit4([4]byte{'a', 'l', 's', 'e'}, d.d.decRd.readn4())
		
	case 't':
		d.checkLit3([3]byte{'r', 'u', 'e'}, d.d.decRd.readn3())
		v = true
	case 'n':
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.d.decRd.readn3())
		
	default:
		d.d.errorf("decode bool: got first char %c", d.tok)
		
	}
	if fquot {
		d.d.decRd.readn1()
	}
	return
}

func (d *jsonDecDriver) DecodeTime() (t time.Time) {
	
	d.advance()
	if d.tok == 'n' {
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.d.decRd.readn3())
		return
	}
	d.ensureReadingString()
	bs := d.readUnescapedString()
	t, err := time.Parse(time.RFC3339, stringView(bs))
	d.d.onerror(err)
	return
}

func (d *jsonDecDriver) ContainerType() (vt valueType) {
	
	d.advance()

	
	

	
	
	if d.tok == '{' {
		return valueTypeMap
	} else if d.tok == '[' {
		return valueTypeArray
	} else if d.tok == 'n' {
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.d.decRd.readn3())
		return valueTypeNil
	} else if d.tok == '"' {
		return valueTypeString
	}
	return valueTypeUnset
}

func (d *jsonDecDriver) decNumBytes() (bs []byte) {
	d.advance()
	dr := &d.d.decRd
	if d.tok == '"' {
		bs = dr.readUntil('"')
	} else if d.tok == 'n' {
		d.checkLit3([3]byte{'u', 'l', 'l'}, dr.readn3())
	} else {
		if jsonManualInlineDecRdInHotZones {
			if dr.bytes {
				bs = dr.rb.jsonReadNum()
			} else {
				bs = dr.ri.jsonReadNum()
			}
		} else {
			bs = dr.jsonReadNum()
		}
	}
	d.tok = 0
	return
}

func (d *jsonDecDriver) DecodeUint64() (u uint64) {
	b := d.decNumBytes()
	u, neg, ok := parseInteger_bytes(b)
	if neg {
		d.d.errorf("negative number cannot be decoded as uint64")
	}
	if !ok {
		d.d.onerror(strconvParseErr(b, "ParseUint"))
	}
	return
}

func (d *jsonDecDriver) DecodeInt64() (v int64) {
	b := d.decNumBytes()
	u, neg, ok := parseInteger_bytes(b)
	if !ok {
		d.d.onerror(strconvParseErr(b, "ParseInt"))
	}
	if chkOvf.Uint2Int(u, neg) {
		d.d.errorf("overflow decoding number from %s", b)
	}
	if neg {
		v = -int64(u)
	} else {
		v = int64(u)
	}
	return
}

func (d *jsonDecDriver) DecodeFloat64() (f float64) {
	var err error
	bs := d.decNumBytes()
	if len(bs) == 0 {
		return
	}
	f, err = parseFloat64(bs)
	d.d.onerror(err)
	return
}

func (d *jsonDecDriver) DecodeFloat32() (f float32) {
	var err error
	bs := d.decNumBytes()
	if len(bs) == 0 {
		return
	}
	f, err = parseFloat32(bs)
	d.d.onerror(err)
	return
}

func (d *jsonDecDriver) DecodeExt(rv interface{}, basetype reflect.Type, xtag uint64, ext Ext) {
	d.advance()
	if d.tok == 'n' {
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.d.decRd.readn3())
		return
	}
	if ext == nil {
		re := rv.(*RawExt)
		re.Tag = xtag
		d.d.decode(&re.Value)
	} else if ext == SelfExt {
		d.d.decodeValue(baseRV(rv), d.h.fnNoExt(basetype))
	} else {
		d.d.interfaceExtConvertAndDecode(rv, ext)
	}
}

func (d *jsonDecDriver) decBytesFromArray(bs []byte) []byte {
	if bs != nil {
		bs = bs[:0]
	}
	d.tok = 0
	bs = append(bs, uint8(d.DecodeUint64()))
	d.tok = d.d.decRd.skipWhitespace() 
	for d.tok != ']' {
		if d.tok != ',' {
			d.d.errorf("read array element - expect char '%c' but got char '%c'", ',', d.tok)
		}
		d.tok = 0
		bs = append(bs, uint8(chkOvf.UintV(d.DecodeUint64(), 8)))
		d.tok = d.d.decRd.skipWhitespace() 
	}
	d.tok = 0
	return bs
}

func (d *jsonDecDriver) DecodeBytes(bs []byte) (bsOut []byte) {
	d.d.decByteState = decByteStateNone
	d.advance()
	if d.tok == 'n' {
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.d.decRd.readn3())
		return nil
	}
	
	if d.rawext {
		bsOut = bs
		d.d.interfaceExtConvertAndDecode(&bsOut, d.h.RawBytesExt)
		return
	}
	
	if d.tok == '[' {
		
		if bs == nil {
			d.d.decByteState = decByteStateReuseBuf
			bs = d.d.b[:]
		}
		return d.decBytesFromArray(bs)
	}

	
	

	d.ensureReadingString()
	bs1 := d.readUnescapedString()
	slen := base64.StdEncoding.DecodedLen(len(bs1))
	if slen == 0 {
		bsOut = []byte{}
	} else if slen <= cap(bs) {
		bsOut = bs[:slen]
	} else if bs == nil {
		d.d.decByteState = decByteStateReuseBuf
		bsOut = d.d.blist.check(*d.buf, slen)
		bsOut = bsOut[:slen]
		*d.buf = bsOut
	} else {
		bsOut = make([]byte, slen)
	}
	slen2, err := base64.StdEncoding.Decode(bsOut, bs1)
	if err != nil {
		d.d.errorf("error decoding base64 binary '%s': %v", bs1, err)
	}
	if slen != slen2 {
		bsOut = bsOut[:slen2]
	}
	return
}

func (d *jsonDecDriver) DecodeStringAsBytes() (s []byte) {
	d.d.decByteState = decByteStateNone
	d.advance()

	
	if d.tok == '"' {
		return d.dblQuoteStringAsBytes()
	}

	
	switch d.tok {
	case 'n':
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.d.decRd.readn3())
		return nil 
	case 'f':
		d.checkLit4([4]byte{'a', 'l', 's', 'e'}, d.d.decRd.readn4())
		return jsonLitb[jsonLitF : jsonLitF+5]
	case 't':
		d.checkLit3([3]byte{'r', 'u', 'e'}, d.d.decRd.readn3())
		return jsonLitb[jsonLitT : jsonLitT+4]
	default:
		
		d.tok = 0
		return d.d.decRd.jsonReadNum()
	}
}

func (d *jsonDecDriver) ensureReadingString() {
	if d.tok != '"' {
		d.d.errorf("expecting string starting with '\"'; got '%c'", d.tok)
	}
}

func (d *jsonDecDriver) readUnescapedString() (bs []byte) {
	
	bs = d.d.decRd.readUntil('"')
	d.tok = 0
	return
}

func (d *jsonDecDriver) dblQuoteStringAsBytes() (buf []byte) {
	checkUtf8 := d.h.ValidateUnicode
	d.d.decByteState = decByteStateNone
	
	buf = (*d.buf)[:0]
	dr := &d.d.decRd
	d.tok = 0

	var bs []byte
	var c byte
	var firstTime bool = true

	for {
		if firstTime {
			firstTime = false
			if dr.bytes {
				bs = dr.rb.jsonReadAsisChars()
				if bs[len(bs)-1] == '"' {
					d.d.decByteState = decByteStateZerocopy
					return bs[:len(bs)-1]
				}
				goto APPEND
			}
		}

		if jsonManualInlineDecRdInHotZones {
			if dr.bytes {
				bs = dr.rb.jsonReadAsisChars()
			} else {
				bs = dr.ri.jsonReadAsisChars()
			}
		} else {
			bs = dr.jsonReadAsisChars()
		}

	APPEND:
		_ = bs[0] 
		buf = append(buf, bs[:len(bs)-1]...)
		c = bs[len(bs)-1]

		if c == '"' {
			break
		}

		
		c = dr.readn1()

		switch c {
		case '"', '\\', '/', '\'':
			buf = append(buf, c)
		case 'b':
			buf = append(buf, '\b')
		case 'f':
			buf = append(buf, '\f')
		case 'n':
			buf = append(buf, '\n')
		case 'r':
			buf = append(buf, '\r')
		case 't':
			buf = append(buf, '\t')
		case 'u':
			rr := d.appendStringAsBytesSlashU()
			if checkUtf8 && rr == unicode.ReplacementChar {
				d.d.errorf("invalid UTF-8 character found after: %s", buf)
			}
			buf = append(buf, d.bstr[:utf8.EncodeRune(d.bstr[:], rr)]...)
		default:
			*d.buf = buf
			d.d.errorf("unsupported escaped value: %c", c)
		}
	}
	*d.buf = buf
	d.d.decByteState = decByteStateReuseBuf
	return
}

func (d *jsonDecDriver) appendStringAsBytesSlashU() (r rune) {
	var rr uint32
	var csu [2]byte
	var cs [4]byte = d.d.decRd.readn4()
	if rr = jsonSlashURune(cs); rr == unicode.ReplacementChar {
		return unicode.ReplacementChar
	}
	r = rune(rr)
	if utf16.IsSurrogate(r) {
		csu = d.d.decRd.readn2()
		cs = d.d.decRd.readn4()
		if csu[0] == '\\' && csu[1] == 'u' {
			if rr = jsonSlashURune(cs); rr == unicode.ReplacementChar {
				return unicode.ReplacementChar
			}
			return utf16.DecodeRune(r, rune(rr))
		}
		return unicode.ReplacementChar
	}
	return
}

func jsonSlashURune(cs [4]byte) (rr uint32) {
	for _, c := range cs {
		
		
		if c >= '0' && c <= '9' {
			rr = rr*16 + uint32(c-jsonU4Chk2)
		} else if c >= 'a' && c <= 'f' {
			rr = rr*16 + uint32(c-jsonU4Chk1)
		} else if c >= 'A' && c <= 'F' {
			rr = rr*16 + uint32(c-jsonU4Chk0)
		} else {
			return unicode.ReplacementChar
		}
	}
	return
}

func (d *jsonDecDriver) nakedNum(z *fauxUnion, bs []byte) (err error) {
	
	if d.h.PreferFloat {
		z.v = valueTypeFloat
		z.f, err = parseFloat64(bs)
	} else {
		err = parseNumber(bs, z, d.h.SignedInteger)
	}
	return
}

func (d *jsonDecDriver) DecodeNaked() {
	z := d.d.naked()

	d.advance()
	var bs []byte
	switch d.tok {
	case 'n':
		d.checkLit3([3]byte{'u', 'l', 'l'}, d.d.decRd.readn3())
		z.v = valueTypeNil
	case 'f':
		d.checkLit4([4]byte{'a', 'l', 's', 'e'}, d.d.decRd.readn4())
		z.v = valueTypeBool
		z.b = false
	case 't':
		d.checkLit3([3]byte{'r', 'u', 'e'}, d.d.decRd.readn3())
		z.v = valueTypeBool
		z.b = true
	case '{':
		z.v = valueTypeMap 
	case '[':
		z.v = valueTypeArray 
	case '"':
		
		bs = d.dblQuoteStringAsBytes()
		if jsonNakedBoolNullInQuotedStr &&
			d.h.MapKeyAsString && len(bs) > 0 && d.d.c == containerMapKey {
			switch string(bs) {
			
			
			case "true":
				z.v = valueTypeBool
				z.b = true
			case "false":
				z.v = valueTypeBool
				z.b = false
			default:
				
				if err := d.nakedNum(z, bs); err != nil {
					z.v = valueTypeString
					z.s = d.d.stringZC(bs)
				}
			}
		} else {
			z.v = valueTypeString
			z.s = d.d.stringZC(bs)
		}
	default: 
		bs = d.d.decRd.jsonReadNum()
		d.tok = 0
		if len(bs) == 0 {
			d.d.errorf("decode number from empty string")
		}
		if err := d.nakedNum(z, bs); err != nil {
			d.d.errorf("decode number from %s: %v", bs, err)
		}
	}
}































type JsonHandle struct {
	textEncodingType
	BasicHandle

	
	
	
	Indent int8

	
	
	
	
	
	
	
	
	
	
	
	IntegerAsString byte

	
	
	
	
	HTMLCharsAsIs bool

	
	
	
	PreferFloat bool

	
	
	
	
	
	TermWhitespace bool

	
	
	
	
	MapKeyAsString bool

	

	
	
	

	
	
	RawBytesExt InterfaceExt
}

func (h *JsonHandle) isJson() bool { return true }


func (h *JsonHandle) Name() string { return "json" }

func (h *JsonHandle) desc(bd byte) string { return string(bd) }

func (h *JsonHandle) typical() bool {
	return h.Indent == 0 && !h.MapKeyAsString && h.IntegerAsString != 'A' && h.IntegerAsString != 'L'
}

func (h *JsonHandle) newEncDriver() encDriver {
	var e = &jsonEncDriver{h: h}
	
	
	e.e.e = e
	e.e.js = true
	e.e.init(h)
	e.reset()
	return e
}

func (h *JsonHandle) newDecDriver() decDriver {
	var d = &jsonDecDriver{h: h}
	var x []byte
	d.buf = &x
	d.d.d = d
	d.d.js = true
	d.d.jsms = h.MapKeyAsString
	d.d.init(h)
	d.reset()
	return d
}

func (e *jsonEncDriver) resetState() {
	e.dl = 0
}

func (e *jsonEncDriver) reset() {
	e.resetState()
	
	
	e.typical = e.h.typical()
	if e.h.HTMLCharsAsIs {
		e.s = &jsonCharSafeSet
	} else {
		e.s = &jsonCharHtmlSafeSet
	}
	e.rawext = e.h.RawBytesExt != nil
	e.di = int8(e.h.Indent)
	e.d = e.h.Indent != 0
	e.ks = e.h.MapKeyAsString
	e.is = e.h.IntegerAsString
}

func (d *jsonDecDriver) resetState() {
	*d.buf = d.d.blist.check(*d.buf, 256)
	d.tok = 0
}

func (d *jsonDecDriver) reset() {
	d.resetState()
	d.rawext = d.h.RawBytesExt != nil
}

func jsonFloatStrconvFmtPrec64(f float64) (fmt byte, prec int8) {
	fmt = 'f'
	prec = -1
	fbits := math.Float64bits(f)
	abs := math.Float64frombits(fbits &^ (1 << 63))
	if abs == 0 || abs == 1 {
		prec = 1
	} else if abs < 1e-6 || abs >= 1e21 {
		fmt = 'e'
	} else if noFrac64(fbits) {
		prec = 1
	}
	return
}

func jsonFloatStrconvFmtPrec32(f float32) (fmt byte, prec int8) {
	fmt = 'f'
	prec = -1
	
	fbits := math.Float32bits(f)
	abs := math.Float32frombits(fbits &^ (1 << 31))
	if abs == 0 || abs == 1 {
		prec = 1
	} else if abs < 1e-6 || abs >= 1e21 {
		fmt = 'e'
	} else if noFrac32(fbits) {
		prec = 1
	}
	return
}

var _ decDriverContainerTracker = (*jsonDecDriver)(nil)
var _ encDriverContainerTracker = (*jsonEncDriver)(nil)
var _ decDriver = (*jsonDecDriver)(nil)
var _ encDriver = (*jsonEncDriver)(nil)
