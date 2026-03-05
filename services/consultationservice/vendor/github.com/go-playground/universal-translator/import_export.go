package ut

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"io"

	"github.com/go-playground/locales"
)

type translation struct {
	Locale           string      `json:"locale"`
	Key              interface{} `json:"key"` 
	Translation      string      `json:"trans"`
	PluralType       string      `json:"type,omitempty"`
	PluralRule       string      `json:"rule,omitempty"`
	OverrideExisting bool        `json:"override,omitempty"`
}

const (
	cardinalType = "Cardinal"
	ordinalType  = "Ordinal"
	rangeType    = "Range"
)


type ImportExportFormat uint8


const (
	FormatJSON ImportExportFormat = iota
)




func (t *UniversalTranslator) Export(format ImportExportFormat, dirname string) error {

	_, err := os.Stat(dirname)
	if err != nil {

		if !os.IsNotExist(err) {
			return err
		}

		if err = os.MkdirAll(dirname, 0744); err != nil {
			return err
		}
	}

	
	var trans []translation
	var b []byte
	var ext string

	for _, locale := range t.translators {

		for k, v := range locale.(*translator).translations {
			trans = append(trans, translation{
				Locale:      locale.Locale(),
				Key:         k,
				Translation: v.text,
			})
		}

		for k, pluralTrans := range locale.(*translator).cardinalTanslations {

			for i, plural := range pluralTrans {

				
				
				if plural == nil {
					continue
				}

				trans = append(trans, translation{
					Locale:      locale.Locale(),
					Key:         k.(string),
					Translation: plural.text,
					PluralType:  cardinalType,
					PluralRule:  locales.PluralRule(i).String(),
				})
			}
		}

		for k, pluralTrans := range locale.(*translator).ordinalTanslations {

			for i, plural := range pluralTrans {

				
				
				if plural == nil {
					continue
				}

				trans = append(trans, translation{
					Locale:      locale.Locale(),
					Key:         k.(string),
					Translation: plural.text,
					PluralType:  ordinalType,
					PluralRule:  locales.PluralRule(i).String(),
				})
			}
		}

		for k, pluralTrans := range locale.(*translator).rangeTanslations {

			for i, plural := range pluralTrans {

				
				
				if plural == nil {
					continue
				}

				trans = append(trans, translation{
					Locale:      locale.Locale(),
					Key:         k.(string),
					Translation: plural.text,
					PluralType:  rangeType,
					PluralRule:  locales.PluralRule(i).String(),
				})
			}
		}

		switch format {
		case FormatJSON:
			b, err = json.MarshalIndent(trans, "", "    ")
			ext = ".json"
		}

		if err != nil {
			return err
		}

		err = os.WriteFile(filepath.Join(dirname, fmt.Sprintf("%s%s", locale.Locale(), ext)), b, 0644)
		if err != nil {
			return err
		}

		trans = trans[0:0]
	}

	return nil
}




func (t *UniversalTranslator) Import(format ImportExportFormat, dirnameOrFilename string) error {

	fi, err := os.Stat(dirnameOrFilename)
	if err != nil {
		return err
	}

	processFn := func(filename string) error {

		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer f.Close()

		return t.ImportByReader(format, f)
	}

	if !fi.IsDir() {
		return processFn(dirnameOrFilename)
	}

	
	walker := func(path string, info os.FileInfo, err error) error {

		if info.IsDir() {
			return nil
		}

		switch format {
		case FormatJSON:
			
			if filepath.Ext(info.Name()) != ".json" {
				return nil
			}
		}

		return processFn(path)
	}

	return filepath.Walk(dirnameOrFilename, walker)
}




func (t *UniversalTranslator) ImportByReader(format ImportExportFormat, reader io.Reader) error {

	b, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	var trans []translation

	switch format {
	case FormatJSON:
		err = json.Unmarshal(b, &trans)
	}

	if err != nil {
		return err
	}

	for _, tl := range trans {

		locale, found := t.FindTranslator(tl.Locale)
		if !found {
			return &ErrMissingLocale{locale: tl.Locale}
		}

		pr := stringToPR(tl.PluralRule)

		if pr == locales.PluralRuleUnknown {

			err = locale.Add(tl.Key, tl.Translation, tl.OverrideExisting)
			if err != nil {
				return err
			}

			continue
		}

		switch tl.PluralType {
		case cardinalType:
			err = locale.AddCardinal(tl.Key, tl.Translation, pr, tl.OverrideExisting)
		case ordinalType:
			err = locale.AddOrdinal(tl.Key, tl.Translation, pr, tl.OverrideExisting)
		case rangeType:
			err = locale.AddRange(tl.Key, tl.Translation, pr, tl.OverrideExisting)
		default:
			return &ErrBadPluralDefinition{tl: tl}
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func stringToPR(s string) locales.PluralRule {

	switch s {
	case "Zero":
		return locales.PluralRuleZero
	case "One":
		return locales.PluralRuleOne
	case "Two":
		return locales.PluralRuleTwo
	case "Few":
		return locales.PluralRuleFew
	case "Many":
		return locales.PluralRuleMany
	case "Other":
		return locales.PluralRuleOther
	default:
		return locales.PluralRuleUnknown
	}

}
