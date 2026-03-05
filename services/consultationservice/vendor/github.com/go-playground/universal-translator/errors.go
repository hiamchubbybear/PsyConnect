package ut

import (
	"errors"
	"fmt"

	"github.com/go-playground/locales"
)

var (
	
	ErrUnknowTranslation = errors.New("Unknown Translation")
)

var _ error = new(ErrConflictingTranslation)
var _ error = new(ErrRangeTranslation)
var _ error = new(ErrOrdinalTranslation)
var _ error = new(ErrCardinalTranslation)
var _ error = new(ErrMissingPluralTranslation)
var _ error = new(ErrExistingTranslator)


type ErrExistingTranslator struct {
	locale string
}


func (e *ErrExistingTranslator) Error() string {
	return fmt.Sprintf("error: conflicting translator for locale '%s'", e.locale)
}


type ErrConflictingTranslation struct {
	locale string
	key    interface{}
	rule   locales.PluralRule
	text   string
}


func (e *ErrConflictingTranslation) Error() string {

	if _, ok := e.key.(string); !ok {
		return fmt.Sprintf("error: conflicting key '%#v' rule '%s' with text '%s' for locale '%s', value being ignored", e.key, e.rule, e.text, e.locale)
	}

	return fmt.Sprintf("error: conflicting key '%s' rule '%s' with text '%s' for locale '%s', value being ignored", e.key, e.rule, e.text, e.locale)
}


type ErrRangeTranslation struct {
	text string
}


func (e *ErrRangeTranslation) Error() string {
	return e.text
}


type ErrOrdinalTranslation struct {
	text string
}


func (e *ErrOrdinalTranslation) Error() string {
	return e.text
}


type ErrCardinalTranslation struct {
	text string
}


func (e *ErrCardinalTranslation) Error() string {
	return e.text
}



type ErrMissingPluralTranslation struct {
	locale          string
	key             interface{}
	rule            locales.PluralRule
	translationType string
}


func (e *ErrMissingPluralTranslation) Error() string {

	if _, ok := e.key.(string); !ok {
		return fmt.Sprintf("error: missing '%s' plural rule '%s' for translation with key '%#v' and locale '%s'", e.translationType, e.rule, e.key, e.locale)
	}

	return fmt.Sprintf("error: missing '%s' plural rule '%s' for translation with key '%s' and locale '%s'", e.translationType, e.rule, e.key, e.locale)
}



type ErrMissingBracket struct {
	locale string
	key    interface{}
	text   string
}


func (e *ErrMissingBracket) Error() string {
	return fmt.Sprintf("error: missing bracket '{}', in translation. locale: '%s' key: '%v' text: '%s'", e.locale, e.key, e.text)
}



type ErrBadParamSyntax struct {
	locale string
	param  string
	key    interface{}
	text   string
}


func (e *ErrBadParamSyntax) Error() string {
	return fmt.Sprintf("error: bad parameter syntax, missing parameter '%s' in translation. locale: '%s' key: '%v' text: '%s'", e.param, e.locale, e.key, e.text)
}





type ErrMissingLocale struct {
	locale string
}


func (e *ErrMissingLocale) Error() string {
	return fmt.Sprintf("error: locale '%s' not registered.", e.locale)
}



type ErrBadPluralDefinition struct {
	tl translation
}


func (e *ErrBadPluralDefinition) Error() string {
	return fmt.Sprintf("error: bad plural definition '%#v'", e.tl)
}
