package ut

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-playground/locales"
)

const (
	paramZero          = "{0}"
	paramOne           = "{1}"
	unknownTranslation = ""
)





type Translator interface {
	locales.Translator

	
	
	
	Add(key interface{}, text string, override bool) error

	
	
	
	
	
	AddCardinal(key interface{}, text string, rule locales.PluralRule, override bool) error

	
	
	
	
	
	
	AddOrdinal(key interface{}, text string, rule locales.PluralRule, override bool) error

	
	
	
	AddRange(key interface{}, text string, rule locales.PluralRule, override bool) error

	
	T(key interface{}, params ...string) (string, error)

	
	
	C(key interface{}, num float64, digits uint64, param string) (string, error)

	
	
	O(key interface{}, num float64, digits uint64, param string) (string, error)

	
	
	R(key interface{}, num1 float64, digits1 uint64, num2 float64, digits2 uint64, param1, param2 string) (string, error)

	
	
	VerifyTranslations() error
}

var _ Translator = new(translator)
var _ locales.Translator = new(translator)

type translator struct {
	locales.Translator
	translations        map[interface{}]*transText
	cardinalTanslations map[interface{}][]*transText 
	ordinalTanslations  map[interface{}][]*transText
	rangeTanslations    map[interface{}][]*transText
}

type transText struct {
	text    string
	indexes []int
}

func newTranslator(trans locales.Translator) Translator {
	return &translator{
		Translator:          trans,
		translations:        make(map[interface{}]*transText), 
		cardinalTanslations: make(map[interface{}][]*transText),
		ordinalTanslations:  make(map[interface{}][]*transText),
		rangeTanslations:    make(map[interface{}][]*transText),
	}
}




func (t *translator) Add(key interface{}, text string, override bool) error {

	if _, ok := t.translations[key]; ok && !override {
		return &ErrConflictingTranslation{locale: t.Locale(), key: key, text: text}
	}

	lb := strings.Count(text, "{")
	rb := strings.Count(text, "}")

	if lb != rb {
		return &ErrMissingBracket{locale: t.Locale(), key: key, text: text}
	}

	trans := &transText{
		text: text,
	}

	var idx int

	for i := 0; i < lb; i++ {
		s := "{" + strconv.Itoa(i) + "}"
		idx = strings.Index(text, s)
		if idx == -1 {
			return &ErrBadParamSyntax{locale: t.Locale(), param: s, key: key, text: text}
		}

		trans.indexes = append(trans.indexes, idx)
		trans.indexes = append(trans.indexes, idx+len(s))
	}

	t.translations[key] = trans

	return nil
}






func (t *translator) AddCardinal(key interface{}, text string, rule locales.PluralRule, override bool) error {

	var verified bool

	
	for _, pr := range t.PluralsCardinal() {
		if pr == rule {
			verified = true
			break
		}
	}

	if !verified {
		return &ErrCardinalTranslation{text: fmt.Sprintf("error: cardinal plural rule '%s' does not exist for locale '%s' key: '%v' text: '%s'", rule, t.Locale(), key, text)}
	}

	tarr, ok := t.cardinalTanslations[key]
	if ok {
		
		if len(tarr) > 0 && tarr[rule] != nil && !override {
			return &ErrConflictingTranslation{locale: t.Locale(), key: key, rule: rule, text: text}
		}

	} else {
		tarr = make([]*transText, 7)
		t.cardinalTanslations[key] = tarr
	}

	trans := &transText{
		text:    text,
		indexes: make([]int, 2),
	}

	tarr[rule] = trans

	idx := strings.Index(text, paramZero)
	if idx == -1 {
		tarr[rule] = nil
		return &ErrCardinalTranslation{text: fmt.Sprintf("error: parameter '%s' not found, may want to use 'Add' instead of 'AddCardinal'. locale: '%s' key: '%v' text: '%s'", paramZero, t.Locale(), key, text)}
	}

	trans.indexes[0] = idx
	trans.indexes[1] = idx + len(paramZero)

	return nil
}






func (t *translator) AddOrdinal(key interface{}, text string, rule locales.PluralRule, override bool) error {

	var verified bool

	
	for _, pr := range t.PluralsOrdinal() {
		if pr == rule {
			verified = true
			break
		}
	}

	if !verified {
		return &ErrOrdinalTranslation{text: fmt.Sprintf("error: ordinal plural rule '%s' does not exist for locale '%s' key: '%v' text: '%s'", rule, t.Locale(), key, text)}
	}

	tarr, ok := t.ordinalTanslations[key]
	if ok {
		
		if len(tarr) > 0 && tarr[rule] != nil && !override {
			return &ErrConflictingTranslation{locale: t.Locale(), key: key, rule: rule, text: text}
		}

	} else {
		tarr = make([]*transText, 7)
		t.ordinalTanslations[key] = tarr
	}

	trans := &transText{
		text:    text,
		indexes: make([]int, 2),
	}

	tarr[rule] = trans

	idx := strings.Index(text, paramZero)
	if idx == -1 {
		tarr[rule] = nil
		return &ErrOrdinalTranslation{text: fmt.Sprintf("error: parameter '%s' not found, may want to use 'Add' instead of 'AddOrdinal'. locale: '%s' key: '%v' text: '%s'", paramZero, t.Locale(), key, text)}
	}

	trans.indexes[0] = idx
	trans.indexes[1] = idx + len(paramZero)

	return nil
}




func (t *translator) AddRange(key interface{}, text string, rule locales.PluralRule, override bool) error {

	var verified bool

	
	for _, pr := range t.PluralsRange() {
		if pr == rule {
			verified = true
			break
		}
	}

	if !verified {
		return &ErrRangeTranslation{text: fmt.Sprintf("error: range plural rule '%s' does not exist for locale '%s' key: '%v' text: '%s'", rule, t.Locale(), key, text)}
	}

	tarr, ok := t.rangeTanslations[key]
	if ok {
		
		if len(tarr) > 0 && tarr[rule] != nil && !override {
			return &ErrConflictingTranslation{locale: t.Locale(), key: key, rule: rule, text: text}
		}

	} else {
		tarr = make([]*transText, 7)
		t.rangeTanslations[key] = tarr
	}

	trans := &transText{
		text:    text,
		indexes: make([]int, 4),
	}

	tarr[rule] = trans

	idx := strings.Index(text, paramZero)
	if idx == -1 {
		tarr[rule] = nil
		return &ErrRangeTranslation{text: fmt.Sprintf("error: parameter '%s' not found, are you sure you're adding a Range Translation? locale: '%s' key: '%v' text: '%s'", paramZero, t.Locale(), key, text)}
	}

	trans.indexes[0] = idx
	trans.indexes[1] = idx + len(paramZero)

	idx = strings.Index(text, paramOne)
	if idx == -1 {
		tarr[rule] = nil
		return &ErrRangeTranslation{text: fmt.Sprintf("error: parameter '%s' not found, a Range Translation requires two parameters. locale: '%s' key: '%v' text: '%s'", paramOne, t.Locale(), key, text)}
	}

	trans.indexes[2] = idx
	trans.indexes[3] = idx + len(paramOne)

	return nil
}


func (t *translator) T(key interface{}, params ...string) (string, error) {

	trans, ok := t.translations[key]
	if !ok {
		return unknownTranslation, ErrUnknowTranslation
	}

	b := make([]byte, 0, 64)

	var start, end, count int

	for i := 0; i < len(trans.indexes); i++ {
		end = trans.indexes[i]
		b = append(b, trans.text[start:end]...)
		b = append(b, params[count]...)
		i++
		start = trans.indexes[i]
		count++
	}

	b = append(b, trans.text[start:]...)

	return string(b), nil
}


func (t *translator) C(key interface{}, num float64, digits uint64, param string) (string, error) {

	tarr, ok := t.cardinalTanslations[key]
	if !ok {
		return unknownTranslation, ErrUnknowTranslation
	}

	rule := t.CardinalPluralRule(num, digits)

	trans := tarr[rule]

	b := make([]byte, 0, 64)
	b = append(b, trans.text[:trans.indexes[0]]...)
	b = append(b, param...)
	b = append(b, trans.text[trans.indexes[1]:]...)

	return string(b), nil
}


func (t *translator) O(key interface{}, num float64, digits uint64, param string) (string, error) {

	tarr, ok := t.ordinalTanslations[key]
	if !ok {
		return unknownTranslation, ErrUnknowTranslation
	}

	rule := t.OrdinalPluralRule(num, digits)

	trans := tarr[rule]

	b := make([]byte, 0, 64)
	b = append(b, trans.text[:trans.indexes[0]]...)
	b = append(b, param...)
	b = append(b, trans.text[trans.indexes[1]:]...)

	return string(b), nil
}



func (t *translator) R(key interface{}, num1 float64, digits1 uint64, num2 float64, digits2 uint64, param1, param2 string) (string, error) {

	tarr, ok := t.rangeTanslations[key]
	if !ok {
		return unknownTranslation, ErrUnknowTranslation
	}

	rule := t.RangePluralRule(num1, digits1, num2, digits2)

	trans := tarr[rule]

	b := make([]byte, 0, 64)
	b = append(b, trans.text[:trans.indexes[0]]...)
	b = append(b, param1...)
	b = append(b, trans.text[trans.indexes[1]:trans.indexes[2]]...)
	b = append(b, param2...)
	b = append(b, trans.text[trans.indexes[3]:]...)

	return string(b), nil
}



func (t *translator) VerifyTranslations() error {

	for k, v := range t.cardinalTanslations {

		for _, rule := range t.PluralsCardinal() {

			if v[rule] == nil {
				return &ErrMissingPluralTranslation{locale: t.Locale(), translationType: "plural", rule: rule, key: k}
			}
		}
	}

	for k, v := range t.ordinalTanslations {

		for _, rule := range t.PluralsOrdinal() {

			if v[rule] == nil {
				return &ErrMissingPluralTranslation{locale: t.Locale(), translationType: "ordinal", rule: rule, key: k}
			}
		}
	}

	for k, v := range t.rangeTanslations {

		for _, rule := range t.PluralsRange() {

			if v[rule] == nil {
				return &ErrMissingPluralTranslation{locale: t.Locale(), translationType: "range", rule: rule, key: k}
			}
		}
	}

	return nil
}
