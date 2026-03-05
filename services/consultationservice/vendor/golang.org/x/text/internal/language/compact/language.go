



//go:generate go run gen.go gen_index.go -output tables.go
//go:generate go run gen_parents.go

package compact




import (
	"strings"

	"golang.org/x/text/internal/language"
)




type Tag struct {
	
	language ID
	locale   ID
	full     fullTag 
}

const _und = 0

type fullTag interface {
	IsRoot() bool
	Parent() language.Tag
}


func Make(t language.Tag) (tag Tag) {
	if region := t.TypeForKey("rg"); len(region) == 6 && region[2:] == "zzzz" {
		if r, err := language.ParseRegion(region[:2]); err == nil {
			tFull := t
			t, _ = t.SetTypeForKey("rg", "")
			
			var exact1, exact2 bool
			tag.language, exact1 = FromTag(t)
			t.RegionID = r
			tag.locale, exact2 = FromTag(t)
			if !exact1 || !exact2 {
				tag.full = tFull
			}
			return tag
		}
	}
	lang, ok := FromTag(t)
	tag.language = lang
	tag.locale = lang
	if !ok {
		tag.full = t
	}
	return tag
}


func (t Tag) Tag() language.Tag {
	if t.full != nil {
		return t.full.(language.Tag)
	}
	tag := t.language.Tag()
	if t.language != t.locale {
		loc := t.locale.Tag()
		tag, _ = tag.SetTypeForKey("rg", strings.ToLower(loc.RegionID.String())+"zzzz")
	}
	return tag
}


func (t *Tag) IsCompact() bool {
	return t.full == nil
}



func (t Tag) MayHaveVariants() bool {
	return t.full != nil || int(t.language) >= len(coreTags)
}



func (t Tag) MayHaveExtensions() bool {
	return t.full != nil ||
		int(t.language) >= len(coreTags) ||
		t.language != t.locale
}


func (t Tag) IsRoot() bool {
	if t.full != nil {
		return t.full.IsRoot()
	}
	return t.language == _und
}




func (t Tag) Parent() Tag {
	if t.full != nil {
		return Make(t.full.Parent())
	}
	if t.language != t.locale {
		
		return Tag{language: t.language, locale: t.language}
	}
	
	
	
	
	lang, _ := FromTag(t.language.Tag().Parent())
	return Tag{language: lang, locale: lang}
}


func nextToken(s string) (t, tail string) {
	p := strings.Index(s[1:], "-")
	if p == -1 {
		return s[1:], ""
	}
	p++
	return s[1:p], s[p:]
}






func LanguageID(t Tag) (id ID, exact bool) {
	return t.language, t.full == nil
}










func RegionalID(t Tag) (id ID, exact bool) {
	return t.locale, t.full == nil
}





func (t Tag) LanguageTag() Tag {
	if t.full == nil {
		return Tag{language: t.language, locale: t.language}
	}
	tt := t.Tag()
	tt.SetTypeForKey("rg", "")
	tt.SetTypeForKey("va", "")
	return Make(tt)
}





func (t Tag) RegionalTag() Tag {
	rt := Tag{language: t.locale, locale: t.locale}
	if t.full == nil {
		return rt
	}
	b := language.Builder{}
	tag := t.Tag()
	
	b.SetTag(t.locale.Tag())
	if v := tag.Variants(); v != "" {
		for _, v := range strings.Split(v, "-") {
			b.AddVariant(v)
		}
	}
	for _, e := range tag.Extensions() {
		b.AddExt(e)
	}
	return t
}


func FromTag(t language.Tag) (id ID, exact bool) {
	
	
	
	exact = true

	b, s, r := t.Raw()
	if t.HasString() {
		if t.IsPrivateUse() {
			
			return 0, false
		}
		hasExtra := false
		if t.HasVariants() {
			if t.HasExtensions() {
				build := language.Builder{}
				build.SetTag(language.Tag{LangID: b, ScriptID: s, RegionID: r})
				build.AddVariant(t.Variants())
				exact = false
				t = build.Make()
			}
			hasExtra = true
		} else if _, ok := t.Extension('u'); ok {
			
			
			old := t
			variant := t.TypeForKey("va")
			t = language.Tag{LangID: b, ScriptID: s, RegionID: r}
			if variant != "" {
				t, _ = t.SetTypeForKey("va", variant)
				hasExtra = true
			}
			exact = old == t
		} else {
			exact = false
		}
		if hasExtra {
			
			for i, s := range specialTags {
				if s == t {
					return ID(i + len(coreTags)), exact
				}
			}
			exact = false
		}
	}
	if x, ok := getCoreIndex(t); ok {
		return x, exact
	}
	exact = false
	if r != 0 && s == 0 {
		
		t, _ := t.Maximize()
		if x, ok := getCoreIndex(t); ok {
			return x, exact
		}
	}
	for t = t.Parent(); t != root; t = t.Parent() {
		
		
		
		if x, ok := getCoreIndex(t); ok {
			return x, exact
		}
	}
	return 0, exact
}

var root = language.Tag{}
