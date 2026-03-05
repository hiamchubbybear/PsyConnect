



package language

import (
	"errors"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/text/internal/language"
)




type ValueError interface {
	error

	
	Subtag() string
}









func Parse(s string) (t Tag, err error) {
	return Default.Parse(s)
}









func (c CanonType) Parse(s string) (t Tag, err error) {
	defer func() {
		if recover() != nil {
			t = Tag{}
			err = language.ErrSyntax
		}
	}()

	tt, err := language.Parse(s)
	if err != nil {
		return makeTag(tt), err
	}
	tt, changed := canonicalize(c, tt)
	if changed {
		tt.RemakeString()
	}
	return makeTag(tt), nil
}











func Compose(part ...interface{}) (t Tag, err error) {
	return Default.Compose(part...)
}











func (c CanonType) Compose(part ...interface{}) (t Tag, err error) {
	defer func() {
		if recover() != nil {
			t = Tag{}
			err = language.ErrSyntax
		}
	}()

	var b language.Builder
	if err = update(&b, part...); err != nil {
		return und, err
	}
	b.Tag, _ = canonicalize(c, b.Tag)
	return makeTag(b.Make()), err
}

var errInvalidArgument = errors.New("invalid Extension or Variant")

func update(b *language.Builder, part ...interface{}) (err error) {
	for _, x := range part {
		switch v := x.(type) {
		case Tag:
			b.SetTag(v.tag())
		case Base:
			b.Tag.LangID = v.langID
		case Script:
			b.Tag.ScriptID = v.scriptID
		case Region:
			b.Tag.RegionID = v.regionID
		case Variant:
			if v.variant == "" {
				err = errInvalidArgument
				break
			}
			b.AddVariant(v.variant)
		case Extension:
			if v.s == "" {
				err = errInvalidArgument
				break
			}
			b.SetExt(v.s)
		case []Variant:
			b.ClearVariants()
			for _, v := range v {
				b.AddVariant(v.variant)
			}
		case []Extension:
			b.ClearExtensions()
			for _, e := range v {
				b.SetExt(e.s)
			}
		
		case error:
			if v != nil {
				err = v
			}
		}
	}
	return
}

var errInvalidWeight = errors.New("ParseAcceptLanguage: invalid weight")
var errTagListTooLarge = errors.New("tag list exceeds max length")








func ParseAcceptLanguage(s string) (tag []Tag, q []float32, err error) {
	defer func() {
		if recover() != nil {
			tag = nil
			q = nil
			err = language.ErrSyntax
		}
	}()

	if strings.Count(s, "-") > 1000 {
		return nil, nil, errTagListTooLarge
	}

	var entry string
	for s != "" {
		if entry, s = split(s, ','); entry == "" {
			continue
		}

		entry, weight := split(entry, ';')

		
		t, err := Parse(entry)
		if err != nil {
			id, ok := acceptFallback[entry]
			if !ok {
				return nil, nil, err
			}
			t = makeTag(language.Tag{LangID: id})
		}

		
		w := 1.0
		if weight != "" {
			weight = consume(weight, 'q')
			weight = consume(weight, '=')
			
			
			if w, err = strconv.ParseFloat(weight, 32); err != nil {
				return nil, nil, errInvalidWeight
			}
			
			if w <= 0 {
				continue
			}
		}

		tag = append(tag, t)
		q = append(q, float32(w))
	}
	sort.Stable(&tagSort{tag, q})
	return tag, q, nil
}



func consume(s string, c byte) string {
	if s == "" || s[0] != c {
		return ""
	}
	return strings.TrimSpace(s[1:])
}

func split(s string, c byte) (head, tail string) {
	if i := strings.IndexByte(s, c); i >= 0 {
		return strings.TrimSpace(s[:i]), strings.TrimSpace(s[i+1:])
	}
	return strings.TrimSpace(s), ""
}



var acceptFallback = map[string]language.Language{
	"english": _en,
	"deutsch": _de,
	"italian": _it,
	"french":  _fr,
	"*":       _mul, 
}

type tagSort struct {
	tag []Tag
	q   []float32
}

func (s *tagSort) Len() int {
	return len(s.q)
}

func (s *tagSort) Less(i, j int) bool {
	return s.q[i] > s.q[j]
}

func (s *tagSort) Swap(i, j int) {
	s.tag[i], s.tag[j] = s.tag[j], s.tag[i]
	s.q[i], s.q[j] = s.q[j], s.q[i]
}
