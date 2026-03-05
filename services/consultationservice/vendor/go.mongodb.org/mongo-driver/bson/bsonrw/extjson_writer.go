





package bsonrw

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"math"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"go.mongodb.org/mongo-driver/bson/primitive"
)




type ExtJSONValueWriterPool struct {
	pool sync.Pool
}




func NewExtJSONValueWriterPool() *ExtJSONValueWriterPool {
	return &ExtJSONValueWriterPool{
		pool: sync.Pool{
			New: func() interface{} {
				return new(extJSONValueWriter)
			},
		},
	}
}




func (bvwp *ExtJSONValueWriterPool) Get(w io.Writer, canonical, escapeHTML bool) ValueWriter {
	vw := bvwp.pool.Get().(*extJSONValueWriter)
	if writer, ok := w.(*SliceWriter); ok {
		vw.reset(*writer, canonical, escapeHTML)
		vw.w = writer
		return vw
	}
	vw.buf = vw.buf[:0]
	vw.w = w
	return vw
}





func (bvwp *ExtJSONValueWriterPool) Put(vw ValueWriter) (ok bool) {
	bvw, ok := vw.(*extJSONValueWriter)
	if !ok {
		return false
	}

	if _, ok := bvw.w.(*SliceWriter); ok {
		bvw.buf = nil
	}
	bvw.w = nil

	bvwp.pool.Put(bvw)
	return true
}

type ejvwState struct {
	mode mode
}

type extJSONValueWriter struct {
	w   io.Writer
	buf []byte

	stack      []ejvwState
	frame      int64
	canonical  bool
	escapeHTML bool
	newlines   bool
}


func NewExtJSONValueWriter(w io.Writer, canonical, escapeHTML bool) (ValueWriter, error) {
	if w == nil {
		return nil, errNilWriter
	}

	
	
	
	return newExtJSONWriter(w, canonical, escapeHTML, true), nil
}

func newExtJSONWriter(w io.Writer, canonical, escapeHTML, newlines bool) *extJSONValueWriter {
	stack := make([]ejvwState, 1, 5)
	stack[0] = ejvwState{mode: mTopLevel}

	return &extJSONValueWriter{
		w:          w,
		buf:        []byte{},
		stack:      stack,
		canonical:  canonical,
		escapeHTML: escapeHTML,
		newlines:   newlines,
	}
}

func newExtJSONWriterFromSlice(buf []byte, canonical, escapeHTML bool) *extJSONValueWriter {
	stack := make([]ejvwState, 1, 5)
	stack[0] = ejvwState{mode: mTopLevel}

	return &extJSONValueWriter{
		buf:        buf,
		stack:      stack,
		canonical:  canonical,
		escapeHTML: escapeHTML,
	}
}

func (ejvw *extJSONValueWriter) reset(buf []byte, canonical, escapeHTML bool) {
	if ejvw.stack == nil {
		ejvw.stack = make([]ejvwState, 1, 5)
	}

	ejvw.stack = ejvw.stack[:1]
	ejvw.stack[0] = ejvwState{mode: mTopLevel}
	ejvw.canonical = canonical
	ejvw.escapeHTML = escapeHTML
	ejvw.frame = 0
	ejvw.buf = buf
	ejvw.w = nil
}

func (ejvw *extJSONValueWriter) advanceFrame() {
	if ejvw.frame+1 >= int64(len(ejvw.stack)) { 
		length := len(ejvw.stack)
		if length+1 >= cap(ejvw.stack) {
			
			buf := make([]ejvwState, 2*cap(ejvw.stack)+1)
			copy(buf, ejvw.stack)
			ejvw.stack = buf
		}
		ejvw.stack = ejvw.stack[:length+1]
	}
	ejvw.frame++
}

func (ejvw *extJSONValueWriter) push(m mode) {
	ejvw.advanceFrame()

	ejvw.stack[ejvw.frame].mode = m
}

func (ejvw *extJSONValueWriter) pop() {
	switch ejvw.stack[ejvw.frame].mode {
	case mElement, mValue:
		ejvw.frame--
	case mDocument, mArray, mCodeWithScope:
		ejvw.frame -= 2 
	}
}

func (ejvw *extJSONValueWriter) invalidTransitionErr(destination mode, name string, modes []mode) error {
	te := TransitionError{
		name:        name,
		current:     ejvw.stack[ejvw.frame].mode,
		destination: destination,
		modes:       modes,
		action:      "write",
	}
	if ejvw.frame != 0 {
		te.parent = ejvw.stack[ejvw.frame-1].mode
	}
	return te
}

func (ejvw *extJSONValueWriter) ensureElementValue(destination mode, callerName string, addmodes ...mode) error {
	switch ejvw.stack[ejvw.frame].mode {
	case mElement, mValue:
	default:
		modes := []mode{mElement, mValue}
		if addmodes != nil {
			modes = append(modes, addmodes...)
		}
		return ejvw.invalidTransitionErr(destination, callerName, modes)
	}

	return nil
}

func (ejvw *extJSONValueWriter) writeExtendedSingleValue(key string, value string, quotes bool) {
	var s string
	if quotes {
		s = fmt.Sprintf(`{"$%s":"%s"}`, key, value)
	} else {
		s = fmt.Sprintf(`{"$%s":%s}`, key, value)
	}

	ejvw.buf = append(ejvw.buf, []byte(s)...)
}

func (ejvw *extJSONValueWriter) WriteArray() (ArrayWriter, error) {
	if err := ejvw.ensureElementValue(mArray, "WriteArray"); err != nil {
		return nil, err
	}

	ejvw.buf = append(ejvw.buf, '[')

	ejvw.push(mArray)
	return ejvw, nil
}

func (ejvw *extJSONValueWriter) WriteBinary(b []byte) error {
	return ejvw.WriteBinaryWithSubtype(b, 0x00)
}

func (ejvw *extJSONValueWriter) WriteBinaryWithSubtype(b []byte, btype byte) error {
	if err := ejvw.ensureElementValue(mode(0), "WriteBinaryWithSubtype"); err != nil {
		return err
	}

	var buf bytes.Buffer
	buf.WriteString(`{"$binary":{"base64":"`)
	buf.WriteString(base64.StdEncoding.EncodeToString(b))
	buf.WriteString(fmt.Sprintf(`","subType":"%02x"}},`, btype))

	ejvw.buf = append(ejvw.buf, buf.Bytes()...)

	ejvw.pop()
	return nil
}

func (ejvw *extJSONValueWriter) WriteBoolean(b bool) error {
	if err := ejvw.ensureElementValue(mode(0), "WriteBoolean"); err != nil {
		return err
	}

	ejvw.buf = append(ejvw.buf, []byte(strconv.FormatBool(b))...)
	ejvw.buf = append(ejvw.buf, ',')

	ejvw.pop()
	return nil
}

func (ejvw *extJSONValueWriter) WriteCodeWithScope(code string) (DocumentWriter, error) {
	if err := ejvw.ensureElementValue(mCodeWithScope, "WriteCodeWithScope"); err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.WriteString(`{"$code":`)
	writeStringWithEscapes(code, &buf, ejvw.escapeHTML)
	buf.WriteString(`,"$scope":{`)

	ejvw.buf = append(ejvw.buf, buf.Bytes()...)

	ejvw.push(mCodeWithScope)
	return ejvw, nil
}

func (ejvw *extJSONValueWriter) WriteDBPointer(ns string, oid primitive.ObjectID) error {
	if err := ejvw.ensureElementValue(mode(0), "WriteDBPointer"); err != nil {
		return err
	}

	var buf bytes.Buffer
	buf.WriteString(`{"$dbPointer":{"$ref":"`)
	buf.WriteString(ns)
	buf.WriteString(`","$id":{"$oid":"`)
	buf.WriteString(oid.Hex())
	buf.WriteString(`"}}},`)

	ejvw.buf = append(ejvw.buf, buf.Bytes()...)

	ejvw.pop()
	return nil
}

func (ejvw *extJSONValueWriter) WriteDateTime(dt int64) error {
	if err := ejvw.ensureElementValue(mode(0), "WriteDateTime"); err != nil {
		return err
	}

	t := time.Unix(dt/1e3, dt%1e3*1e6).UTC()

	if ejvw.canonical || t.Year() < 1970 || t.Year() > 9999 {
		s := fmt.Sprintf(`{"$numberLong":"%d"}`, dt)
		ejvw.writeExtendedSingleValue("date", s, false)
	} else {
		ejvw.writeExtendedSingleValue("date", t.Format(rfc3339Milli), true)
	}

	ejvw.buf = append(ejvw.buf, ',')

	ejvw.pop()
	return nil
}

func (ejvw *extJSONValueWriter) WriteDecimal128(d primitive.Decimal128) error {
	if err := ejvw.ensureElementValue(mode(0), "WriteDecimal128"); err != nil {
		return err
	}

	ejvw.writeExtendedSingleValue("numberDecimal", d.String(), true)
	ejvw.buf = append(ejvw.buf, ',')

	ejvw.pop()
	return nil
}

func (ejvw *extJSONValueWriter) WriteDocument() (DocumentWriter, error) {
	if ejvw.stack[ejvw.frame].mode == mTopLevel {
		ejvw.buf = append(ejvw.buf, '{')
		return ejvw, nil
	}

	if err := ejvw.ensureElementValue(mDocument, "WriteDocument", mTopLevel); err != nil {
		return nil, err
	}

	ejvw.buf = append(ejvw.buf, '{')
	ejvw.push(mDocument)
	return ejvw, nil
}

func (ejvw *extJSONValueWriter) WriteDouble(f float64) error {
	if err := ejvw.ensureElementValue(mode(0), "WriteDouble"); err != nil {
		return err
	}

	s := formatDouble(f)

	if ejvw.canonical {
		ejvw.writeExtendedSingleValue("numberDouble", s, true)
	} else {
		switch s {
		case "Infinity":
			fallthrough
		case "-Infinity":
			fallthrough
		case "NaN":
			s = fmt.Sprintf(`{"$numberDouble":"%s"}`, s)
		}
		ejvw.buf = append(ejvw.buf, []byte(s)...)
	}

	ejvw.buf = append(ejvw.buf, ',')

	ejvw.pop()
	return nil
}

func (ejvw *extJSONValueWriter) WriteInt32(i int32) error {
	if err := ejvw.ensureElementValue(mode(0), "WriteInt32"); err != nil {
		return err
	}

	s := strconv.FormatInt(int64(i), 10)

	if ejvw.canonical {
		ejvw.writeExtendedSingleValue("numberInt", s, true)
	} else {
		ejvw.buf = append(ejvw.buf, []byte(s)...)
	}

	ejvw.buf = append(ejvw.buf, ',')

	ejvw.pop()
	return nil
}

func (ejvw *extJSONValueWriter) WriteInt64(i int64) error {
	if err := ejvw.ensureElementValue(mode(0), "WriteInt64"); err != nil {
		return err
	}

	s := strconv.FormatInt(i, 10)

	if ejvw.canonical {
		ejvw.writeExtendedSingleValue("numberLong", s, true)
	} else {
		ejvw.buf = append(ejvw.buf, []byte(s)...)
	}

	ejvw.buf = append(ejvw.buf, ',')

	ejvw.pop()
	return nil
}

func (ejvw *extJSONValueWriter) WriteJavascript(code string) error {
	if err := ejvw.ensureElementValue(mode(0), "WriteJavascript"); err != nil {
		return err
	}

	var buf bytes.Buffer
	writeStringWithEscapes(code, &buf, ejvw.escapeHTML)

	ejvw.writeExtendedSingleValue("code", buf.String(), false)
	ejvw.buf = append(ejvw.buf, ',')

	ejvw.pop()
	return nil
}

func (ejvw *extJSONValueWriter) WriteMaxKey() error {
	if err := ejvw.ensureElementValue(mode(0), "WriteMaxKey"); err != nil {
		return err
	}

	ejvw.writeExtendedSingleValue("maxKey", "1", false)
	ejvw.buf = append(ejvw.buf, ',')

	ejvw.pop()
	return nil
}

func (ejvw *extJSONValueWriter) WriteMinKey() error {
	if err := ejvw.ensureElementValue(mode(0), "WriteMinKey"); err != nil {
		return err
	}

	ejvw.writeExtendedSingleValue("minKey", "1", false)
	ejvw.buf = append(ejvw.buf, ',')

	ejvw.pop()
	return nil
}

func (ejvw *extJSONValueWriter) WriteNull() error {
	if err := ejvw.ensureElementValue(mode(0), "WriteNull"); err != nil {
		return err
	}

	ejvw.buf = append(ejvw.buf, []byte("null")...)
	ejvw.buf = append(ejvw.buf, ',')

	ejvw.pop()
	return nil
}

func (ejvw *extJSONValueWriter) WriteObjectID(oid primitive.ObjectID) error {
	if err := ejvw.ensureElementValue(mode(0), "WriteObjectID"); err != nil {
		return err
	}

	ejvw.writeExtendedSingleValue("oid", oid.Hex(), true)
	ejvw.buf = append(ejvw.buf, ',')

	ejvw.pop()
	return nil
}

func (ejvw *extJSONValueWriter) WriteRegex(pattern string, options string) error {
	if err := ejvw.ensureElementValue(mode(0), "WriteRegex"); err != nil {
		return err
	}

	options = sortStringAlphebeticAscending(options)
	var buf bytes.Buffer
	buf.WriteString(`{"$regularExpression":{"pattern":`)
	writeStringWithEscapes(pattern, &buf, ejvw.escapeHTML)
	buf.WriteString(`,"options":`)
	writeStringWithEscapes(options, &buf, ejvw.escapeHTML)
	buf.WriteString(`}},`)

	ejvw.buf = append(ejvw.buf, buf.Bytes()...)

	ejvw.pop()
	return nil
}

func (ejvw *extJSONValueWriter) WriteString(s string) error {
	if err := ejvw.ensureElementValue(mode(0), "WriteString"); err != nil {
		return err
	}

	var buf bytes.Buffer
	writeStringWithEscapes(s, &buf, ejvw.escapeHTML)

	ejvw.buf = append(ejvw.buf, buf.Bytes()...)
	ejvw.buf = append(ejvw.buf, ',')

	ejvw.pop()
	return nil
}

func (ejvw *extJSONValueWriter) WriteSymbol(symbol string) error {
	if err := ejvw.ensureElementValue(mode(0), "WriteSymbol"); err != nil {
		return err
	}

	var buf bytes.Buffer
	writeStringWithEscapes(symbol, &buf, ejvw.escapeHTML)

	ejvw.writeExtendedSingleValue("symbol", buf.String(), false)
	ejvw.buf = append(ejvw.buf, ',')

	ejvw.pop()
	return nil
}

func (ejvw *extJSONValueWriter) WriteTimestamp(t uint32, i uint32) error {
	if err := ejvw.ensureElementValue(mode(0), "WriteTimestamp"); err != nil {
		return err
	}

	var buf bytes.Buffer
	buf.WriteString(`{"$timestamp":{"t":`)
	buf.WriteString(strconv.FormatUint(uint64(t), 10))
	buf.WriteString(`,"i":`)
	buf.WriteString(strconv.FormatUint(uint64(i), 10))
	buf.WriteString(`}},`)

	ejvw.buf = append(ejvw.buf, buf.Bytes()...)

	ejvw.pop()
	return nil
}

func (ejvw *extJSONValueWriter) WriteUndefined() error {
	if err := ejvw.ensureElementValue(mode(0), "WriteUndefined"); err != nil {
		return err
	}

	ejvw.writeExtendedSingleValue("undefined", "true", false)
	ejvw.buf = append(ejvw.buf, ',')

	ejvw.pop()
	return nil
}

func (ejvw *extJSONValueWriter) WriteDocumentElement(key string) (ValueWriter, error) {
	switch ejvw.stack[ejvw.frame].mode {
	case mDocument, mTopLevel, mCodeWithScope:
		var buf bytes.Buffer
		writeStringWithEscapes(key, &buf, ejvw.escapeHTML)

		ejvw.buf = append(ejvw.buf, []byte(fmt.Sprintf(`%s:`, buf.String()))...)
		ejvw.push(mElement)
	default:
		return nil, ejvw.invalidTransitionErr(mElement, "WriteDocumentElement", []mode{mDocument, mTopLevel, mCodeWithScope})
	}

	return ejvw, nil
}

func (ejvw *extJSONValueWriter) WriteDocumentEnd() error {
	switch ejvw.stack[ejvw.frame].mode {
	case mDocument, mTopLevel, mCodeWithScope:
	default:
		return fmt.Errorf("incorrect mode to end document: %s", ejvw.stack[ejvw.frame].mode)
	}

	
	if ejvw.buf[len(ejvw.buf)-1] == ',' {
		ejvw.buf[len(ejvw.buf)-1] = '}'
	} else {
		ejvw.buf = append(ejvw.buf, '}')
	}

	switch ejvw.stack[ejvw.frame].mode {
	case mCodeWithScope:
		ejvw.buf = append(ejvw.buf, '}')
		fallthrough
	case mDocument:
		ejvw.buf = append(ejvw.buf, ',')
	case mTopLevel:
		
		
		
		if ejvw.newlines {
			ejvw.buf = append(ejvw.buf, '\n')
		}
		if ejvw.w != nil {
			if _, err := ejvw.w.Write(ejvw.buf); err != nil {
				return err
			}
			ejvw.buf = ejvw.buf[:0]
		}
	}

	ejvw.pop()
	return nil
}

func (ejvw *extJSONValueWriter) WriteArrayElement() (ValueWriter, error) {
	switch ejvw.stack[ejvw.frame].mode {
	case mArray:
		ejvw.push(mValue)
	default:
		return nil, ejvw.invalidTransitionErr(mValue, "WriteArrayElement", []mode{mArray})
	}

	return ejvw, nil
}

func (ejvw *extJSONValueWriter) WriteArrayEnd() error {
	switch ejvw.stack[ejvw.frame].mode {
	case mArray:
		
		if ejvw.buf[len(ejvw.buf)-1] == ',' {
			ejvw.buf[len(ejvw.buf)-1] = ']'
		} else {
			ejvw.buf = append(ejvw.buf, ']')
		}

		ejvw.buf = append(ejvw.buf, ',')

		ejvw.pop()
	default:
		return fmt.Errorf("incorrect mode to end array: %s", ejvw.stack[ejvw.frame].mode)
	}

	return nil
}

func formatDouble(f float64) string {
	var s string
	switch {
	case math.IsInf(f, 1):
		s = "Infinity"
	case math.IsInf(f, -1):
		s = "-Infinity"
	case math.IsNaN(f):
		s = "NaN"
	default:
		
		
		s = strconv.FormatFloat(f, 'G', -1, 64)
		if !strings.ContainsRune(s, 'E') && !strings.ContainsRune(s, '.') {
			s += ".0"
		}
	}

	return s
}

var hexChars = "0123456789abcdef"

func writeStringWithEscapes(s string, buf *bytes.Buffer, escapeHTML bool) {
	buf.WriteByte('"')
	start := 0
	for i := 0; i < len(s); {
		if b := s[i]; b < utf8.RuneSelf {
			if htmlSafeSet[b] || (!escapeHTML && safeSet[b]) {
				i++
				continue
			}
			if start < i {
				buf.WriteString(s[start:i])
			}
			switch b {
			case '\\', '"':
				buf.WriteByte('\\')
				buf.WriteByte(b)
			case '\n':
				buf.WriteByte('\\')
				buf.WriteByte('n')
			case '\r':
				buf.WriteByte('\\')
				buf.WriteByte('r')
			case '\t':
				buf.WriteByte('\\')
				buf.WriteByte('t')
			case '\b':
				buf.WriteByte('\\')
				buf.WriteByte('b')
			case '\f':
				buf.WriteByte('\\')
				buf.WriteByte('f')
			default:
				
				
				
				
				
				buf.WriteString(`\u00`)
				buf.WriteByte(hexChars[b>>4])
				buf.WriteByte(hexChars[b&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				buf.WriteString(s[start:i])
			}
			buf.WriteString(`\ufffd`)
			i += size
			start = i
			continue
		}
		
		
		
		
		
		
		
		if c == '\u2028' || c == '\u2029' {
			if start < i {
				buf.WriteString(s[start:i])
			}
			buf.WriteString(`\u202`)
			buf.WriteByte(hexChars[c&0xF])
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(s) {
		buf.WriteString(s[start:])
	}
	buf.WriteByte('"')
}

type sortableString []rune

func (ss sortableString) Len() int {
	return len(ss)
}

func (ss sortableString) Less(i, j int) bool {
	return ss[i] < ss[j]
}

func (ss sortableString) Swap(i, j int) {
	ss[i], ss[j] = ss[j], ss[i]
}

func sortStringAlphebeticAscending(s string) string {
	ss := sortableString([]rune(s))
	sort.Sort(ss)
	return string([]rune(ss))
}
