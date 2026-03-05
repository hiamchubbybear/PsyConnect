


package baggage 

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"

	"go.opentelemetry.io/otel/internal/baggage"
)

const (
	maxMembers               = 180
	maxBytesPerMembers       = 4096
	maxBytesPerBaggageString = 8192

	listDelimiter     = ","
	keyValueDelimiter = "="
	propertyDelimiter = ";"
)

var (
	errInvalidKey      = errors.New("invalid key")
	errInvalidValue    = errors.New("invalid value")
	errInvalidProperty = errors.New("invalid baggage list-member property")
	errInvalidMember   = errors.New("invalid baggage list-member")
	errMemberNumber    = errors.New("too many list-members in baggage-string")
	errMemberBytes     = errors.New("list-member too large")
	errBaggageBytes    = errors.New("baggage-string too large")
)


type Property struct {
	key, value string

	
	
	hasValue bool
}










func NewKeyProperty(key string) (Property, error) {
	if !validateBaggageName(key) {
		return newInvalidProperty(), fmt.Errorf("%w: %q", errInvalidKey, key)
	}

	p := Property{key: key}
	return p, nil
}








func NewKeyValueProperty(key, value string) (Property, error) {
	if !validateKey(key) {
		return newInvalidProperty(), fmt.Errorf("%w: %q", errInvalidKey, key)
	}

	if !validateValue(value) {
		return newInvalidProperty(), fmt.Errorf("%w: %q", errInvalidValue, value)
	}
	decodedValue, err := url.PathUnescape(value)
	if err != nil {
		return newInvalidProperty(), fmt.Errorf("%w: %q", errInvalidValue, value)
	}
	return NewKeyValuePropertyRaw(key, decodedValue)
}










func NewKeyValuePropertyRaw(key, value string) (Property, error) {
	if !validateBaggageName(key) {
		return newInvalidProperty(), fmt.Errorf("%w: %q", errInvalidKey, key)
	}
	if !validateBaggageValue(value) {
		return newInvalidProperty(), fmt.Errorf("%w: %q", errInvalidValue, value)
	}

	p := Property{
		key:      key,
		value:    value,
		hasValue: true,
	}
	return p, nil
}

func newInvalidProperty() Property {
	return Property{}
}




func parseProperty(property string) (Property, error) {
	if property == "" {
		return newInvalidProperty(), nil
	}

	p, ok := parsePropertyInternal(property)
	if !ok {
		return newInvalidProperty(), fmt.Errorf("%w: %q", errInvalidProperty, property)
	}

	return p, nil
}



func (p Property) validate() error {
	errFunc := func(err error) error {
		return fmt.Errorf("invalid property: %w", err)
	}

	if !validateBaggageName(p.key) {
		return errFunc(fmt.Errorf("%w: %q", errInvalidKey, p.key))
	}
	if !p.hasValue && p.value != "" {
		return errFunc(errors.New("inconsistent value"))
	}
	if p.hasValue && !validateBaggageValue(p.value) {
		return errFunc(fmt.Errorf("%w: %q", errInvalidValue, p.value))
	}
	return nil
}


func (p Property) Key() string {
	return p.key
}




func (p Property) Value() (string, bool) {
	return p.value, p.hasValue
}






func (p Property) String() string {
	
	if !validateKey(p.key) {
		return ""
	}

	if p.hasValue {
		return fmt.Sprintf("%s%s%v", p.key, keyValueDelimiter, valueEscape(p.value))
	}
	return p.key
}

type properties []Property

func fromInternalProperties(iProps []baggage.Property) properties {
	if len(iProps) == 0 {
		return nil
	}

	props := make(properties, len(iProps))
	for i, p := range iProps {
		props[i] = Property{
			key:      p.Key,
			value:    p.Value,
			hasValue: p.HasValue,
		}
	}
	return props
}

func (p properties) asInternal() []baggage.Property {
	if len(p) == 0 {
		return nil
	}

	iProps := make([]baggage.Property, len(p))
	for i, prop := range p {
		iProps[i] = baggage.Property{
			Key:      prop.key,
			Value:    prop.value,
			HasValue: prop.hasValue,
		}
	}
	return iProps
}

func (p properties) Copy() properties {
	if len(p) == 0 {
		return nil
	}

	props := make(properties, len(p))
	copy(props, p)
	return props
}



func (p properties) validate() error {
	for _, prop := range p {
		if err := prop.validate(); err != nil {
			return err
		}
	}
	return nil
}



func (p properties) String() string {
	props := make([]string, 0, len(p))
	for _, prop := range p {
		s := prop.String()

		
		if s != "" {
			props = append(props, s)
		}
	}
	return strings.Join(props, propertyDelimiter)
}



type Member struct {
	key, value string
	properties properties

	
	
	
	hasData bool
}








func NewMember(key, value string, props ...Property) (Member, error) {
	if !validateKey(key) {
		return newInvalidMember(), fmt.Errorf("%w: %q", errInvalidKey, key)
	}

	if !validateValue(value) {
		return newInvalidMember(), fmt.Errorf("%w: %q", errInvalidValue, value)
	}
	decodedValue, err := url.PathUnescape(value)
	if err != nil {
		return newInvalidMember(), fmt.Errorf("%w: %q", errInvalidValue, value)
	}
	return NewMemberRaw(key, decodedValue, props...)
}










func NewMemberRaw(key, value string, props ...Property) (Member, error) {
	m := Member{
		key:        key,
		value:      value,
		properties: properties(props).Copy(),
		hasData:    true,
	}
	if err := m.validate(); err != nil {
		return newInvalidMember(), err
	}
	return m, nil
}

func newInvalidMember() Member {
	return Member{}
}




func parseMember(member string) (Member, error) {
	if n := len(member); n > maxBytesPerMembers {
		return newInvalidMember(), fmt.Errorf("%w: %d", errMemberBytes, n)
	}

	var props properties
	keyValue, properties, found := strings.Cut(member, propertyDelimiter)
	if found {
		
		for _, pStr := range strings.Split(properties, propertyDelimiter) {
			p, err := parseProperty(pStr)
			if err != nil {
				return newInvalidMember(), err
			}
			props = append(props, p)
		}
	}
	

	
	k, v, found := strings.Cut(keyValue, keyValueDelimiter)
	if !found {
		return newInvalidMember(), fmt.Errorf("%w: %q", errInvalidMember, member)
	}
	
	
	key := strings.TrimSpace(k)
	if !validateKey(key) {
		return newInvalidMember(), fmt.Errorf("%w: %q", errInvalidKey, key)
	}

	rawVal := strings.TrimSpace(v)
	if !validateValue(rawVal) {
		return newInvalidMember(), fmt.Errorf("%w: %q", errInvalidValue, v)
	}

	
	unescapeVal, err := url.PathUnescape(rawVal)
	if err != nil {
		return newInvalidMember(), fmt.Errorf("%w: %w", errInvalidValue, err)
	}

	value := replaceInvalidUTF8Sequences(len(rawVal), unescapeVal)
	return Member{key: key, value: value, properties: props, hasData: true}, nil
}


func replaceInvalidUTF8Sequences(c int, unescapeVal string) string {
	if utf8.ValidString(unescapeVal) {
		return unescapeVal
	}
	
	

	var b strings.Builder
	b.Grow(c)
	for i := 0; i < len(unescapeVal); {
		r, size := utf8.DecodeRuneInString(unescapeVal[i:])
		if r == utf8.RuneError && size == 1 {
			
			_, _ = b.WriteString("�")
		} else {
			_, _ = b.WriteRune(r)
		}
		i += size
	}

	return b.String()
}



func (m Member) validate() error {
	if !m.hasData {
		return fmt.Errorf("%w: %q", errInvalidMember, m)
	}

	if !validateBaggageName(m.key) {
		return fmt.Errorf("%w: %q", errInvalidKey, m.key)
	}
	if !validateBaggageValue(m.value) {
		return fmt.Errorf("%w: %q", errInvalidValue, m.value)
	}
	return m.properties.validate()
}


func (m Member) Key() string { return m.key }


func (m Member) Value() string { return m.value }


func (m Member) Properties() []Property { return m.properties.Copy() }






func (m Member) String() string {
	
	if !validateKey(m.key) {
		return ""
	}

	s := m.key + keyValueDelimiter + valueEscape(m.value)
	if len(m.properties) > 0 {
		s += propertyDelimiter + m.properties.String()
	}
	return s
}



type Baggage struct { 
	list baggage.List
}





func New(members ...Member) (Baggage, error) {
	if len(members) == 0 {
		return Baggage{}, nil
	}

	b := make(baggage.List)
	for _, m := range members {
		if !m.hasData {
			return Baggage{}, errInvalidMember
		}

		
		b[m.key] = baggage.Item{
			Value:      m.value,
			Properties: m.properties.asInternal(),
		}
	}

	
	if len(b) > maxMembers {
		return Baggage{}, errMemberNumber
	}

	bag := Baggage{b}
	if n := len(bag.String()); n > maxBytesPerBaggageString {
		return Baggage{}, fmt.Errorf("%w: %d", errBaggageBytes, n)
	}

	return bag, nil
}









func Parse(bStr string) (Baggage, error) {
	if bStr == "" {
		return Baggage{}, nil
	}

	if n := len(bStr); n > maxBytesPerBaggageString {
		return Baggage{}, fmt.Errorf("%w: %d", errBaggageBytes, n)
	}

	b := make(baggage.List)
	for _, memberStr := range strings.Split(bStr, listDelimiter) {
		m, err := parseMember(memberStr)
		if err != nil {
			return Baggage{}, err
		}
		
		b[m.key] = baggage.Item{
			Value:      m.value,
			Properties: m.properties.asInternal(),
		}
	}

	
	
	
	if len(b) > maxMembers {
		return Baggage{}, errMemberNumber
	}

	return Baggage{b}, nil
}







func (b Baggage) Member(key string) Member {
	v, ok := b.list[key]
	if !ok {
		
		
		
		
		return newInvalidMember()
	}

	return Member{
		key:        key,
		value:      v.Value,
		properties: fromInternalProperties(v.Properties),
		hasData:    true,
	}
}






func (b Baggage) Members() []Member {
	if len(b.list) == 0 {
		return nil
	}

	members := make([]Member, 0, len(b.list))
	for k, v := range b.list {
		members = append(members, Member{
			key:        k,
			value:      v.Value,
			properties: fromInternalProperties(v.Properties),
			hasData:    true,
		})
	}
	return members
}







func (b Baggage) SetMember(member Member) (Baggage, error) {
	if !member.hasData {
		return b, errInvalidMember
	}

	n := len(b.list)
	if _, ok := b.list[member.key]; !ok {
		n++
	}
	list := make(baggage.List, n)

	for k, v := range b.list {
		
		if k == member.key {
			continue
		}
		list[k] = v
	}

	list[member.key] = baggage.Item{
		Value:      member.value,
		Properties: member.properties.asInternal(),
	}

	return Baggage{list: list}, nil
}



func (b Baggage) DeleteMember(key string) Baggage {
	n := len(b.list)
	if _, ok := b.list[key]; ok {
		n--
	}
	list := make(baggage.List, n)

	for k, v := range b.list {
		if k == key {
			continue
		}
		list[k] = v
	}

	return Baggage{list: list}
}


func (b Baggage) Len() int {
	return len(b.list)
}






func (b Baggage) String() string {
	members := make([]string, 0, len(b.list))
	for k, v := range b.list {
		s := Member{
			key:        k,
			value:      v.Value,
			properties: fromInternalProperties(v.Properties),
		}.String()

		
		if s != "" {
			members = append(members, s)
		}
	}
	return strings.Join(members, listDelimiter)
}



func parsePropertyInternal(s string) (p Property, ok bool) {
	
	
	
	index := skipSpace(s, 0)

	
	keyStart := index
	keyEnd := index
	for _, c := range s[keyStart:] {
		if !validateKeyChar(c) {
			break
		}
		keyEnd++
	}

	
	
	if keyStart == keyEnd {
		return
	}

	
	index = skipSpace(s, keyEnd)

	if index == len(s) {
		
		ok = true
		p.key = s[keyStart:keyEnd]
		return
	}

	
	
	if s[index] != keyValueDelimiter[0] {
		return
	}

	
	
	index = skipSpace(s, index+1)

	
	
	
	valueStart := index
	valueEnd := index
	for _, c := range s[valueStart:] {
		if !validateValueChar(c) {
			break
		}
		valueEnd++
	}

	
	index = skipSpace(s, valueEnd)

	
	
	
	if index != len(s) {
		return
	}

	
	rawVal := s[valueStart:valueEnd]
	unescapeVal, err := url.PathUnescape(rawVal)
	if err != nil {
		return
	}
	value := replaceInvalidUTF8Sequences(len(rawVal), unescapeVal)

	ok = true
	p.key = s[keyStart:keyEnd]
	p.hasValue = true

	p.value = value
	return
}

func skipSpace(s string, offset int) int {
	i := offset
	for ; i < len(s); i++ {
		c := s[i]
		if c != ' ' && c != '\t' {
			break
		}
	}
	return i
}

var safeKeyCharset = [utf8.RuneSelf]bool{
	
	'#':  true,
	'$':  true,
	'%':  true,
	'&':  true,
	'\'': true,

	
	'0': true,
	'1': true,
	'2': true,
	'3': true,
	'4': true,
	'5': true,
	'6': true,
	'7': true,
	'8': true,
	'9': true,

	
	'A': true,
	'B': true,
	'C': true,
	'D': true,
	'E': true,
	'F': true,
	'G': true,
	'H': true,
	'I': true,
	'J': true,
	'K': true,
	'L': true,
	'M': true,
	'N': true,
	'O': true,
	'P': true,
	'Q': true,
	'R': true,
	'S': true,
	'T': true,
	'U': true,
	'V': true,
	'W': true,
	'X': true,
	'Y': true,
	'Z': true,

	
	'^': true,
	'_': true,
	'`': true,
	'a': true,
	'b': true,
	'c': true,
	'd': true,
	'e': true,
	'f': true,
	'g': true,
	'h': true,
	'i': true,
	'j': true,
	'k': true,
	'l': true,
	'm': true,
	'n': true,
	'o': true,
	'p': true,
	'q': true,
	'r': true,
	's': true,
	't': true,
	'u': true,
	'v': true,
	'w': true,
	'x': true,
	'y': true,
	'z': true,

	
	'!': true,
	'*': true,
	'+': true,
	'-': true,
	'.': true,
	'|': true,
	'~': true,
}



func validateBaggageName(s string) bool {
	if len(s) == 0 {
		return false
	}

	return utf8.ValidString(s)
}




func validateBaggageValue(s string) bool {
	return utf8.ValidString(s)
}


func validateKey(s string) bool {
	if len(s) == 0 {
		return false
	}

	for _, c := range s {
		if !validateKeyChar(c) {
			return false
		}
	}

	return true
}

func validateKeyChar(c int32) bool {
	return c >= 0 && c < int32(utf8.RuneSelf) && safeKeyCharset[c]
}


func validateValue(s string) bool {
	for _, c := range s {
		if !validateValueChar(c) {
			return false
		}
	}

	return true
}

var safeValueCharset = [utf8.RuneSelf]bool{
	'!': true, 

	
	'#':  true,
	'$':  true,
	'%':  true,
	'&':  true,
	'\'': true,
	'(':  true,
	')':  true,
	'*':  true,
	'+':  true,

	
	'-': true,
	'.': true,
	'/': true,
	'0': true,
	'1': true,
	'2': true,
	'3': true,
	'4': true,
	'5': true,
	'6': true,
	'7': true,
	'8': true,
	'9': true,
	':': true,

	
	'<': true, 
	'=': true, 
	'>': true, 
	'?': true, 
	'@': true, 
	'A': true, 
	'B': true, 
	'C': true, 
	'D': true, 
	'E': true, 
	'F': true, 
	'G': true, 
	'H': true, 
	'I': true, 
	'J': true, 
	'K': true, 
	'L': true, 
	'M': true, 
	'N': true, 
	'O': true, 
	'P': true, 
	'Q': true, 
	'R': true, 
	'S': true, 
	'T': true, 
	'U': true, 
	'V': true, 
	'W': true, 
	'X': true, 
	'Y': true, 
	'Z': true, 
	'[': true, 

	
	']': true, 
	'^': true, 
	'_': true, 
	'`': true, 
	'a': true, 
	'b': true, 
	'c': true, 
	'd': true, 
	'e': true, 
	'f': true, 
	'g': true, 
	'h': true, 
	'i': true, 
	'j': true, 
	'k': true, 
	'l': true, 
	'm': true, 
	'n': true, 
	'o': true, 
	'p': true, 
	'q': true, 
	'r': true, 
	's': true, 
	't': true, 
	'u': true, 
	'v': true, 
	'w': true, 
	'x': true, 
	'y': true, 
	'z': true, 
	'{': true, 
	'|': true, 
	'}': true, 
	'~': true, 
}

func validateValueChar(c int32) bool {
	return c >= 0 && c < int32(utf8.RuneSelf) && safeValueCharset[c]
}






func valueEscape(s string) string {
	hexCount := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if shouldEscape(c) {
			hexCount++
		}
	}

	if hexCount == 0 {
		return s
	}

	var buf [64]byte
	var t []byte

	required := len(s) + 2*hexCount
	if required <= len(buf) {
		t = buf[:required]
	} else {
		t = make([]byte, required)
	}

	j := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if shouldEscape(s[i]) {
			const upperhex = "0123456789ABCDEF"
			t[j] = '%'
			t[j+1] = upperhex[c>>4]
			t[j+2] = upperhex[c&15]
			j += 3
		} else {
			t[j] = c
			j++
		}
	}

	return string(t)
}



func shouldEscape(c byte) bool {
	if c == '%' {
		
		return true
	}
	return !validateValueChar(int32(c))
}
