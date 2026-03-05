





package stringprep

var errHasLCat = "BiDi string can't have runes from category L"
var errFirstRune = "BiDi string first rune must have category R or AL"
var errLastRune = "BiDi string last rune must have category R or AL"


func checkBiDiProhibitedRune(s string) error {
	for _, r := range s {
		if TableC8.Contains(r) {
			return Error{Msg: errProhibited, Rune: r}
		}
	}
	return nil
}


func checkBiDiLCat(s string) error {
	for _, r := range s {
		if TableD2.Contains(r) {
			return Error{Msg: errHasLCat, Rune: r}
		}
	}
	return nil
}


func checkBadFirstAndLastRandALCat(s string) error {
	rs := []rune(s)
	if !TableD1.Contains(rs[0]) {
		return Error{Msg: errFirstRune, Rune: rs[0]}
	}
	n := len(rs) - 1
	if !TableD1.Contains(rs[n]) {
		return Error{Msg: errLastRune, Rune: rs[n]}
	}
	return nil
}


func hasBiDiRandALCat(s string) bool {
	for _, r := range s {
		if TableD1.Contains(r) {
			return true
		}
	}
	return false
}


func passesBiDiRules(s string) error {
	if len(s) == 0 {
		return nil
	}
	if err := checkBiDiProhibitedRune(s); err != nil {
		return err
	}
	if hasBiDiRandALCat(s) {
		if err := checkBiDiLCat(s); err != nil {
			return err
		}
		if err := checkBadFirstAndLastRandALCat(s); err != nil {
			return err
		}
	}
	return nil
}
