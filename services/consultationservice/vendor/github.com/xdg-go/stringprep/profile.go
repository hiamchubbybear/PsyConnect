package stringprep

import (
	"golang.org/x/text/unicode/norm"
)


type Profile struct {
	Mappings  []Mapping
	Normalize bool
	Prohibits []Set
	CheckBiDi bool
}

var errProhibited = "prohibited character"



func (p Profile) Prepare(s string) (string, error) {
	
	temp := make([]rune, 0, len(s))

	
	for _, r := range s {
		rs, ok := p.applyMaps(r)
		if ok {
			temp = append(temp, rs...)
		} else {
			temp = append(temp, r)
		}
	}

	
	var out string
	if p.Normalize {
		out = norm.NFKC.String(string(temp))
	} else {
		out = string(temp)
	}

	
	for _, r := range out {
		if p.runeIsProhibited(r) {
			return "", Error{Msg: errProhibited, Rune: r}
		}
	}

	
	if p.CheckBiDi {
		if err := passesBiDiRules(out); err != nil {
			return "", err
		}
	}

	return out, nil
}

func (p Profile) applyMaps(r rune) ([]rune, bool) {
	for _, m := range p.Mappings {
		rs, ok := m.Map(r)
		if ok {
			return rs, true
		}
	}
	return nil, false
}

func (p Profile) runeIsProhibited(r rune) bool {
	for _, s := range p.Prohibits {
		if s.Contains(r) {
			return true
		}
	}
	return false
}
