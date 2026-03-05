



package language

import (
	"errors"
	"strings"

	"golang.org/x/text/internal/language"
)


type MatchOption func(*matcher)




func PreferSameScript(preferSame bool) MatchOption {
	return func(m *matcher) { m.preferSameScript = preferSame }
}










func MatchStrings(m Matcher, lang ...string) (tag Tag, index int) {
	for _, accept := range lang {
		desired, _, err := ParseAcceptLanguage(accept)
		if err != nil {
			continue
		}
		if tag, index, conf := m.Match(desired...); conf != No {
			return tag, index
		}
	}
	tag, index, _ = m.Match()
	return
}






type Matcher interface {
	Match(t ...Tag) (tag Tag, index int, c Confidence)
}



func Comprehends(speaker, alternative Tag) Confidence {
	_, _, c := NewMatcher([]Tag{alternative}).Match(speaker)
	return c
}















func NewMatcher(t []Tag, options ...MatchOption) Matcher {
	return newMatcher(t, options)
}

func (m *matcher) Match(want ...Tag) (t Tag, index int, c Confidence) {
	var tt language.Tag
	match, w, c := m.getBest(want...)
	if match != nil {
		tt, index = match.tag, match.index
	} else {
		
		tt = m.default_.tag
		if m.preferSameScript {
		outer:
			for _, w := range want {
				script, _ := w.Script()
				if script.scriptID == 0 {
					
					
					continue
				}
				for i, h := range m.supported {
					if script.scriptID == h.maxScript {
						tt, index = h.tag, i
						break outer
					}
				}
			}
		}
		
	}
	if w.RegionID != tt.RegionID && w.RegionID != 0 {
		if w.RegionID != 0 && tt.RegionID != 0 && tt.RegionID.Contains(w.RegionID) {
			tt.RegionID = w.RegionID
			tt.RemakeString()
		} else if r := w.RegionID.String(); len(r) == 2 {
			
			tt, _ = tt.SetTypeForKey("rg", strings.ToLower(r)+"zzzz")
		}
	}
	
	
	
	
	if e := w.Extensions(); len(e) > 0 {
		b := language.Builder{}
		b.SetTag(tt)
		for _, e := range e {
			b.AddExt(e)
		}
		tt = b.Make()
	}
	return makeTag(tt), index, c
}



var ErrMissingLikelyTagsData = errors.New("missing likely tags data")























































































































type matcher struct {
	default_         *haveTag
	supported        []*haveTag
	index            map[language.Language]*matchHeader
	passSettings     bool
	preferSameScript bool
}



type matchHeader struct {
	haveTags []*haveTag
	original bool
}



type haveTag struct {
	tag language.Tag

	
	index int

	
	
	conf Confidence

	
	maxRegion language.Region
	maxScript language.Script

	
	
	
	altScript language.Script

	
	nextMax uint16
}

func makeHaveTag(tag language.Tag, index int) (haveTag, language.Language) {
	max := tag
	if tag.LangID != 0 || tag.RegionID != 0 || tag.ScriptID != 0 {
		max, _ = canonicalize(All, max)
		max, _ = max.Maximize()
		max.RemakeString()
	}
	return haveTag{tag, index, Exact, max.RegionID, max.ScriptID, altScript(max.LangID, max.ScriptID), 0}, max.LangID
}




func altScript(l language.Language, s language.Script) language.Script {
	for _, alt := range matchScript {
		
		if (language.Language(alt.wantLang) == l || language.Language(alt.haveLang) == l) &&
			language.Script(alt.haveScript) == s {
			return language.Script(alt.wantScript)
		}
	}
	return 0
}



func (h *matchHeader) addIfNew(n haveTag, exact bool) {
	h.original = h.original || exact
	
	for _, v := range h.haveTags {
		if equalsRest(v.tag, n.tag) {
			return
		}
	}
	
	
	for i, v := range h.haveTags {
		if v.maxScript == n.maxScript &&
			v.maxRegion == n.maxRegion &&
			v.tag.VariantOrPrivateUseTags() == n.tag.VariantOrPrivateUseTags() {
			for h.haveTags[i].nextMax != 0 {
				i = int(h.haveTags[i].nextMax)
			}
			h.haveTags[i].nextMax = uint16(len(h.haveTags))
			break
		}
	}
	h.haveTags = append(h.haveTags, &n)
}



func (m *matcher) header(l language.Language) *matchHeader {
	if h := m.index[l]; h != nil {
		return h
	}
	h := &matchHeader{}
	m.index[l] = h
	return h
}

func toConf(d uint8) Confidence {
	if d <= 10 {
		return High
	}
	if d < 30 {
		return Low
	}
	return No
}




func newMatcher(supported []Tag, options []MatchOption) *matcher {
	m := &matcher{
		index:            make(map[language.Language]*matchHeader),
		preferSameScript: true,
	}
	for _, o := range options {
		o(m)
	}
	if len(supported) == 0 {
		m.default_ = &haveTag{}
		return m
	}
	
	
	for i, tag := range supported {
		tt := tag.tag()
		pair, _ := makeHaveTag(tt, i)
		m.header(tt.LangID).addIfNew(pair, true)
		m.supported = append(m.supported, &pair)
	}
	m.default_ = m.header(supported[0].lang()).haveTags[0]
	
	
	for i, tag := range supported {
		tt := tag.tag()
		pair, max := makeHaveTag(tt, i)
		if max != tt.LangID {
			m.header(max).addIfNew(pair, true)
		}
	}

	
	
	
	update := func(want, have uint16, conf Confidence) {
		if hh := m.index[language.Language(have)]; hh != nil {
			if !hh.original {
				return
			}
			hw := m.header(language.Language(want))
			for _, ht := range hh.haveTags {
				v := *ht
				if conf < v.conf {
					v.conf = conf
				}
				v.nextMax = 0 
				if v.altScript != 0 {
					v.altScript = altScript(language.Language(want), v.maxScript)
				}
				hw.addIfNew(v, conf == Exact && hh.original)
			}
		}
	}

	
	
	for _, ml := range matchLang {
		update(ml.want, ml.have, toConf(ml.distance))
		if !ml.oneway {
			update(ml.have, ml.want, toConf(ml.distance))
		}
	}

	
	
	
	
	
	for i, lm := range language.AliasMap {
		
		
		conf := Exact
		if language.AliasTypes[i] != language.Macro {
			if !isExactEquivalent(language.Language(lm.From)) {
				conf = High
			}
			update(lm.To, lm.From, conf)
		}
		update(lm.From, lm.To, conf)
	}
	return m
}



func (m *matcher) getBest(want ...Tag) (got *haveTag, orig language.Tag, c Confidence) {
	best := bestMatch{}
	for i, ww := range want {
		w := ww.tag()
		var max language.Tag
		
		h := m.index[w.LangID]
		if w.LangID != 0 {
			if h == nil {
				continue
			}
			
			max, _ = canonicalize(Legacy|Deprecated|Macro, w)
			
			
			if w.RegionID != max.RegionID {
				w.RegionID = max.RegionID
			}
			
			
			max, _ = max.Maximize()
		} else {
			
			if h != nil {
				for i := range h.haveTags {
					have := h.haveTags[i]
					if equalsRest(have.tag, w) {
						return have, w, Exact
					}
				}
			}
			if w.ScriptID == 0 && w.RegionID == 0 {
				
				
				continue
			}
			max, _ = w.Maximize()
			if h = m.index[max.LangID]; h == nil {
				continue
			}
		}
		pin := true
		for _, t := range want[i+1:] {
			if w.LangID == t.lang() {
				pin = false
				break
			}
		}
		
		for i := range h.haveTags {
			have := h.haveTags[i]
			best.update(have, w, max.ScriptID, max.RegionID, pin)
			if best.conf == Exact {
				for have.nextMax != 0 {
					have = h.haveTags[have.nextMax]
					best.update(have, w, max.ScriptID, max.RegionID, pin)
				}
				return best.have, best.want, best.conf
			}
		}
	}
	if best.conf <= No {
		if len(want) != 0 {
			return nil, want[0].tag(), No
		}
		return nil, language.Tag{}, No
	}
	return best.have, best.want, best.conf
}


type bestMatch struct {
	have            *haveTag
	want            language.Tag
	conf            Confidence
	pinnedRegion    language.Region
	pinLanguage     bool
	sameRegionGroup bool
	
	origLang     bool
	origReg      bool
	paradigmReg  bool
	regGroupDist uint8
	origScript   bool
}















func (m *bestMatch) update(have *haveTag, tag language.Tag, maxScript language.Script, maxRegion language.Region, pin bool) {
	
	c := have.conf
	if c < m.conf {
		return
	}
	
	if m.pinLanguage && tag.LangID != m.want.LangID {
		return
	}
	
	if tag.LangID == m.want.LangID && m.sameRegionGroup {
		_, sameGroup := regionGroupDist(m.pinnedRegion, have.maxRegion, have.maxScript, m.want.LangID)
		if !sameGroup {
			return
		}
	}
	if c == Exact && have.maxScript == maxScript {
		
		
		m.pinLanguage = pin
	}
	if equalsRest(have.tag, tag) {
	} else if have.maxScript != maxScript {
		
		
		
		if Low < m.conf || have.altScript != maxScript {
			return
		}
		c = Low
	} else if have.maxRegion != maxRegion {
		if High < c {
			
			c = High
		}
	}

	
	
	
	
	beaten := false 
	if c != m.conf {
		if c < m.conf {
			return
		}
		beaten = true
	}

	
	
	origLang := have.tag.LangID == tag.LangID && tag.LangID != 0
	if !beaten && m.origLang != origLang {
		if m.origLang {
			return
		}
		beaten = true
	}

	
	origReg := have.tag.RegionID == tag.RegionID && tag.RegionID != 0
	if !beaten && m.origReg != origReg {
		if m.origReg {
			return
		}
		beaten = true
	}

	regGroupDist, sameGroup := regionGroupDist(have.maxRegion, maxRegion, maxScript, tag.LangID)
	if !beaten && m.regGroupDist != regGroupDist {
		if regGroupDist > m.regGroupDist {
			return
		}
		beaten = true
	}

	paradigmReg := isParadigmLocale(tag.LangID, have.maxRegion)
	if !beaten && m.paradigmReg != paradigmReg {
		if !paradigmReg {
			return
		}
		beaten = true
	}

	
	origScript := have.tag.ScriptID == tag.ScriptID && tag.ScriptID != 0
	if !beaten && m.origScript != origScript {
		if m.origScript {
			return
		}
		beaten = true
	}

	
	if beaten {
		m.have = have
		m.want = tag
		m.conf = c
		m.pinnedRegion = maxRegion
		m.sameRegionGroup = sameGroup
		m.origLang = origLang
		m.origReg = origReg
		m.paradigmReg = paradigmReg
		m.origScript = origScript
		m.regGroupDist = regGroupDist
	}
}

func isParadigmLocale(lang language.Language, r language.Region) bool {
	for _, e := range paradigmLocales {
		if language.Language(e[0]) == lang && (r == language.Region(e[1]) || r == language.Region(e[2])) {
			return true
		}
	}
	return false
}



func regionGroupDist(a, b language.Region, script language.Script, lang language.Language) (dist uint8, same bool) {
	const defaultDistance = 4

	aGroup := uint(regionToGroups[a]) << 1
	bGroup := uint(regionToGroups[b]) << 1
	for _, ri := range matchRegion {
		if language.Language(ri.lang) == lang && (ri.script == 0 || language.Script(ri.script) == script) {
			group := uint(1 << (ri.group &^ 0x80))
			if 0x80&ri.group == 0 {
				if aGroup&bGroup&group != 0 { 
					return ri.distance, ri.distance == defaultDistance
				}
			} else {
				if (aGroup|bGroup)&group == 0 { 
					return ri.distance, ri.distance == defaultDistance
				}
			}
		}
	}
	return defaultDistance, true
}


func equalsRest(a, b language.Tag) bool {
	
	
	return a.ScriptID == b.ScriptID && a.RegionID == b.RegionID && a.VariantOrPrivateUseTags() == b.VariantOrPrivateUseTags()
}



func isExactEquivalent(l language.Language) bool {
	for _, o := range notEquivalent {
		if o == l {
			return false
		}
	}
	return true
}

var notEquivalent []language.Language

func init() {
	
	
	for _, lm := range language.AliasMap {
		tag := language.Tag{LangID: language.Language(lm.From)}
		if tag, _ = canonicalize(All, tag); tag.ScriptID != 0 || tag.RegionID != 0 {
			notEquivalent = append(notEquivalent, language.Language(lm.From))
		}
	}
	
	for i, v := range paradigmLocales {
		t := language.Tag{LangID: language.Language(v[0])}
		max, _ := t.Maximize()
		if v[1] == 0 {
			paradigmLocales[i][1] = uint16(max.RegionID)
		}
		if v[2] == 0 {
			paradigmLocales[i][2] = uint16(max.RegionID)
		}
	}
}
