



//go:generate go run gen.go -output tables.go

package language




import (
	"strings"

	"golang.org/x/text/internal/language"
	"golang.org/x/text/internal/language/compact"
)




type Tag compact.Tag

func makeTag(t language.Tag) (tag Tag) {
	return Tag(compact.Make(t))
}

func (t *Tag) tag() language.Tag {
	return (*compact.Tag)(t).Tag()
}

func (t *Tag) isCompact() bool {
	return (*compact.Tag)(t).IsCompact()
}


func (t *Tag) lang() language.Language { return t.tag().LangID }
func (t *Tag) region() language.Region { return t.tag().RegionID }
func (t *Tag) script() language.Script { return t.tag().ScriptID }



func Make(s string) Tag {
	return Default.Make(s)
}



func (c CanonType) Make(s string) Tag {
	t, _ := c.Parse(s)
	return t
}



func (t Tag) Raw() (b Base, s Script, r Region) {
	tt := t.tag()
	return Base{tt.LangID}, Script{tt.ScriptID}, Region{tt.RegionID}
}


func (t Tag) IsRoot() bool {
	return compact.Tag(t).IsRoot()
}


type CanonType int

const (
	
	DeprecatedBase CanonType = 1 << iota
	
	DeprecatedScript
	
	DeprecatedRegion
	
	SuppressScript
	
	
	Legacy
	
	
	Macro
	
	
	
	CLDR

	
	Raw CanonType = 0

	
	Deprecated = DeprecatedBase | DeprecatedScript | DeprecatedRegion

	
	BCP47 = Deprecated | SuppressScript

	
	All = BCP47 | Legacy | Macro

	
	
	
	
	
	Default = Deprecated | Legacy

	canonLang = DeprecatedBase | Legacy | Macro

	
)



func canonicalize(c CanonType, t language.Tag) (language.Tag, bool) {
	if c == Raw {
		return t, false
	}
	changed := false
	if c&SuppressScript != 0 {
		if t.LangID.SuppressScript() == t.ScriptID {
			t.ScriptID = 0
			changed = true
		}
	}
	if c&canonLang != 0 {
		for {
			if l, aliasType := t.LangID.Canonicalize(); l != t.LangID {
				switch aliasType {
				case language.Legacy:
					if c&Legacy != 0 {
						if t.LangID == _sh && t.ScriptID == 0 {
							t.ScriptID = _Latn
						}
						t.LangID = l
						changed = true
					}
				case language.Macro:
					if c&Macro != 0 {
						
						
						
						
						
						
						
						
						
						if c&CLDR == 0 || t.LangID != _nb {
							changed = true
							t.LangID = l
						}
					}
				case language.Deprecated:
					if c&DeprecatedBase != 0 {
						if t.LangID == _mo && t.RegionID == 0 {
							t.RegionID = _MD
						}
						t.LangID = l
						changed = true
						
						continue
					}
				}
			} else if c&Legacy != 0 && t.LangID == _no && c&CLDR != 0 {
				t.LangID = _nb
				changed = true
			}
			break
		}
	}
	if c&DeprecatedScript != 0 {
		if t.ScriptID == _Qaai {
			changed = true
			t.ScriptID = _Zinh
		}
	}
	if c&DeprecatedRegion != 0 {
		if r := t.RegionID.Canonicalize(); r != t.RegionID {
			changed = true
			t.RegionID = r
		}
	}
	return t, changed
}


func (c CanonType) Canonicalize(t Tag) (Tag, error) {
	
	if t.isCompact() {
		if _, changed := canonicalize(c, compact.Tag(t).Tag()); !changed {
			return t, nil
		}
	}
	
	
	if tag, changed := canonicalize(c, t.tag()); changed {
		tag.RemakeString()
		return makeTag(tag), nil
	}
	return t, nil

}






type Confidence int

const (
	No    Confidence = iota 
	Low                     
	High                    
	Exact                   
)

var confName = []string{"No", "Low", "High", "Exact"}

func (c Confidence) String() string {
	return confName[c]
}


func (t Tag) String() string {
	return t.tag().String()
}


func (t Tag) MarshalText() (text []byte, err error) {
	return t.tag().MarshalText()
}


func (t *Tag) UnmarshalText(text []byte) error {
	var tag language.Tag
	err := tag.UnmarshalText(text)
	*t = makeTag(tag)
	return err
}




func (t Tag) Base() (Base, Confidence) {
	if b := t.lang(); b != 0 {
		return Base{b}, Exact
	}
	tt := t.tag()
	c := High
	if tt.ScriptID == 0 && !tt.RegionID.IsCountry() {
		c = Low
	}
	if tag, err := tt.Maximize(); err == nil && tag.LangID != 0 {
		return Base{tag.LangID}, c
	}
	return Base{0}, No
}















func (t Tag) Script() (Script, Confidence) {
	if scr := t.script(); scr != 0 {
		return Script{scr}, Exact
	}
	tt := t.tag()
	sc, c := language.Script(_Zzzz), No
	if scr := tt.LangID.SuppressScript(); scr != 0 {
		
		
		if tt.RegionID == 0 {
			return Script{scr}, High
		}
		sc, c = scr, High
	}
	if tag, err := tt.Maximize(); err == nil {
		if tag.ScriptID != sc {
			sc, c = tag.ScriptID, Low
		}
	} else {
		tt, _ = canonicalize(Deprecated|Macro, tt)
		if tag, err := tt.Maximize(); err == nil && tag.ScriptID != sc {
			sc, c = tag.ScriptID, Low
		}
	}
	return Script{sc}, c
}




func (t Tag) Region() (Region, Confidence) {
	if r := t.region(); r != 0 {
		return Region{r}, Exact
	}
	tt := t.tag()
	if tt, err := tt.Maximize(); err == nil {
		return Region{tt.RegionID}, Low 
	}
	tt, _ = canonicalize(Deprecated|Macro, tt)
	if tag, err := tt.Maximize(); err == nil {
		return Region{tag.RegionID}, Low
	}
	return Region{_ZZ}, No 
}



func (t Tag) Variants() []Variant {
	if !compact.Tag(t).MayHaveVariants() {
		return nil
	}
	v := []Variant{}
	x, str := "", t.tag().Variants()
	for str != "" {
		x, str = nextToken(str)
		v = append(v, Variant{x})
	}
	return v
}









func (t Tag) Parent() Tag {
	return Tag(compact.Tag(t).Parent())
}


func nextToken(s string) (t, tail string) {
	p := strings.Index(s[1:], "-")
	if p == -1 {
		return s[1:], ""
	}
	p++
	return s[1:p], s[p:]
}


type Extension struct {
	s string
}



func (e Extension) String() string {
	return e.s
}


func ParseExtension(s string) (e Extension, err error) {
	ext, err := language.ParseExtension(s)
	return Extension{ext}, err
}



func (e Extension) Type() byte {
	if e.s == "" {
		return 0
	}
	return e.s[0]
}


func (e Extension) Tokens() []string {
	return strings.Split(e.s, "-")
}




func (t Tag) Extension(x byte) (ext Extension, ok bool) {
	if !compact.Tag(t).MayHaveExtensions() {
		return Extension{}, false
	}
	e, ok := t.tag().Extension(x)
	return Extension{e}, ok
}


func (t Tag) Extensions() []Extension {
	if !compact.Tag(t).MayHaveExtensions() {
		return nil
	}
	e := []Extension{}
	for _, ext := range t.tag().Extensions() {
		e = append(e, Extension{ext})
	}
	return e
}









func (t Tag) TypeForKey(key string) string {
	if !compact.Tag(t).MayHaveExtensions() {
		if key != "rg" && key != "va" {
			return ""
		}
	}
	return t.tag().TypeForKey(key)
}





func (t Tag) SetTypeForKey(key, value string) (Tag, error) {
	tt, err := t.tag().SetTypeForKey(key, value)
	return makeTag(tt), err
}



const NumCompactTags = compact.NumCompactTags






func CompactIndex(t Tag) (index int, exact bool) {
	id, exact := compact.LanguageID(compact.Tag(t))
	return int(id), exact
}

var root = language.Tag{}



type Base struct {
	langID language.Language
}




func ParseBase(s string) (Base, error) {
	l, err := language.ParseBase(s)
	return Base{l}, err
}


func (b Base) String() string {
	return b.langID.String()
}


func (b Base) ISO3() string {
	return b.langID.ISO3()
}


func (b Base) IsPrivateUse() bool {
	return b.langID.IsPrivateUse()
}



type Script struct {
	scriptID language.Script
}




func ParseScript(s string) (Script, error) {
	sc, err := language.ParseScript(s)
	return Script{sc}, err
}



func (s Script) String() string {
	return s.scriptID.String()
}


func (s Script) IsPrivateUse() bool {
	return s.scriptID.IsPrivateUse()
}


type Region struct {
	regionID language.Region
}



func EncodeM49(r int) (Region, error) {
	rid, err := language.EncodeM49(r)
	return Region{rid}, err
}




func ParseRegion(s string) (Region, error) {
	r, err := language.ParseRegion(s)
	return Region{r}, err
}



func (r Region) String() string {
	return r.regionID.String()
}




func (r Region) ISO3() string {
	return r.regionID.ISO3()
}



func (r Region) M49() int {
	return r.regionID.M49()
}




func (r Region) IsPrivateUse() bool {
	return r.regionID.IsPrivateUse()
}



func (r Region) IsCountry() bool {
	return r.regionID.IsCountry()
}



func (r Region) IsGroup() bool {
	return r.regionID.IsGroup()
}



func (r Region) Contains(c Region) bool {
	return r.regionID.Contains(c.regionID)
}








func (r Region) TLD() (Region, error) {
	tld, err := r.regionID.TLD()
	return Region{tld}, err
}




func (r Region) Canonicalize() Region {
	return Region{r.regionID.Canonicalize()}
}


type Variant struct {
	variant string
}



func ParseVariant(s string) (Variant, error) {
	v, err := language.ParseVariant(s)
	return Variant{v.String()}, err
}


func (v Variant) String() string {
	return v.variant
}
