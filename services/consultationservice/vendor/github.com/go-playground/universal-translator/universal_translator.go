package ut

import (
	"strings"

	"github.com/go-playground/locales"
)


type UniversalTranslator struct {
	translators map[string]Translator
	fallback    Translator
}



func New(fallback locales.Translator, supportedLocales ...locales.Translator) *UniversalTranslator {

	t := &UniversalTranslator{
		translators: make(map[string]Translator),
	}

	for _, v := range supportedLocales {

		trans := newTranslator(v)
		t.translators[strings.ToLower(trans.Locale())] = trans

		if fallback.Locale() == v.Locale() {
			t.fallback = trans
		}
	}

	if t.fallback == nil && fallback != nil {
		t.fallback = newTranslator(fallback)
	}

	return t
}




func (t *UniversalTranslator) FindTranslator(locales ...string) (trans Translator, found bool) {

	for _, locale := range locales {

		if trans, found = t.translators[strings.ToLower(locale)]; found {
			return
		}
	}

	return t.fallback, false
}



func (t *UniversalTranslator) GetTranslator(locale string) (trans Translator, found bool) {

	if trans, found = t.translators[strings.ToLower(locale)]; found {
		return
	}

	return t.fallback, false
}


func (t *UniversalTranslator) GetFallback() Translator {
	return t.fallback
}





func (t *UniversalTranslator) AddTranslator(translator locales.Translator, override bool) error {

	lc := strings.ToLower(translator.Locale())
	_, ok := t.translators[lc]
	if ok && !override {
		return &ErrExistingTranslator{locale: translator.Locale()}
	}

	trans := newTranslator(translator)

	if t.fallback.Locale() == translator.Locale() {

		
		
		if !override {
			return &ErrExistingTranslator{locale: translator.Locale()}
		}

		t.fallback = trans
	}

	t.translators[lc] = trans

	return nil
}



func (t *UniversalTranslator) VerifyTranslations() (err error) {

	for _, trans := range t.translators {
		err = trans.VerifyTranslations()
		if err != nil {
			return
		}
	}

	return
}
