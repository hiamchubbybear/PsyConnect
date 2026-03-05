package locales

import (
	"strconv"
	"time"

	"github.com/go-playground/locales/currency"
)
















type PluralRule int


const (
	PluralRuleUnknown PluralRule = iota
	PluralRuleZero               
	PluralRuleOne                
	PluralRuleTwo                
	PluralRuleFew                
	PluralRuleMany               
	PluralRuleOther              
)

const (
	pluralsString = "UnknownZeroOneTwoFewManyOther"
)




type Translator interface {

	
	

	
	Locale() string

	
	
	PluralsCardinal() []PluralRule

	
	
	PluralsOrdinal() []PluralRule

	
	
	PluralsRange() []PluralRule

	
	CardinalPluralRule(num float64, v uint64) PluralRule

	
	OrdinalPluralRule(num float64, v uint64) PluralRule

	
	RangePluralRule(num1 float64, v1 uint64, num2 float64, v2 uint64) PluralRule

	
	MonthAbbreviated(month time.Month) string

	
	MonthsAbbreviated() []string

	
	MonthNarrow(month time.Month) string

	
	MonthsNarrow() []string

	
	MonthWide(month time.Month) string

	
	MonthsWide() []string

	
	WeekdayAbbreviated(weekday time.Weekday) string

	
	WeekdaysAbbreviated() []string

	
	WeekdayNarrow(weekday time.Weekday) string

	
	WeekdaysNarrow() []string

	
	WeekdayShort(weekday time.Weekday) string

	
	WeekdaysShort() []string

	
	WeekdayWide(weekday time.Weekday) string

	
	WeekdaysWide() []string

	

	
	FmtNumber(num float64, v uint64) string

	
	
	FmtPercent(num float64, v uint64) string

	
	FmtCurrency(num float64, v uint64, currency currency.Type) string

	
	
	FmtAccounting(num float64, v uint64, currency currency.Type) string

	
	FmtDateShort(t time.Time) string

	
	FmtDateMedium(t time.Time) string

	
	FmtDateLong(t time.Time) string

	
	FmtDateFull(t time.Time) string

	
	FmtTimeShort(t time.Time) string

	
	FmtTimeMedium(t time.Time) string

	
	FmtTimeLong(t time.Time) string

	
	FmtTimeFull(t time.Time) string
}


func (p PluralRule) String() string {

	switch p {
	case PluralRuleZero:
		return pluralsString[7:11]
	case PluralRuleOne:
		return pluralsString[11:14]
	case PluralRuleTwo:
		return pluralsString[14:17]
	case PluralRuleFew:
		return pluralsString[17:20]
	case PluralRuleMany:
		return pluralsString[20:24]
	case PluralRuleOther:
		return pluralsString[24:]
	default:
		return pluralsString[:7]
	}
}


















































func W(n float64, v uint64) (w int64) {

	s := strconv.FormatFloat(n-float64(int64(n)), 'f', int(v), 64)

	
	
	if len(s) != 1 {

		s = s[2:]
		end := len(s) + 1

		for i := end; i >= 0; i-- {
			if s[i] != '0' {
				end = i + 1
				break
			}
		}

		w = int64(len(s[:end]))
	}

	return
}


func F(n float64, v uint64) (f int64) {

	s := strconv.FormatFloat(n-float64(int64(n)), 'f', int(v), 64)

	
	
	if len(s) != 1 {

		
		
		f, _ = strconv.ParseInt(s[2:], 10, 64)
	}

	return
}


func T(n float64, v uint64) (t int64) {

	s := strconv.FormatFloat(n-float64(int64(n)), 'f', int(v), 64)

	
	
	if len(s) != 1 {

		s = s[2:]
		end := len(s) + 1

		for i := end; i >= 0; i-- {
			if s[i] != '0' {
				end = i + 1
				break
			}
		}

		
		
		t, _ = strconv.ParseInt(s[:end], 10, 64)
	}

	return
}
