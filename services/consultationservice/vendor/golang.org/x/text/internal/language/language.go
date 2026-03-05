



//go:generate go run gen.go gen_common.go -output tables.go

package language 




import (
	"errors"
	"fmt"
	"strings"
)

const (
	
	
	maxCoreSize = 12

	
	
	max99thPercentileSize = 32

	
	
	maxSimpleUExtensionSize = 14
)




type Tag struct {
	
	
	
	

	LangID   Language
	RegionID Region
	
	
	
	
	
	
	
	ScriptID Script
	pVariant byte   
	pExt     uint16 

	
	
	str string
}



func Make(s string) Tag {
	t, _ := Parse(s)
	return t
}




func (t Tag) Raw() (b Language, s Script, r Region) {
	return t.LangID, t.ScriptID, t.RegionID
}


func (t Tag) equalTags(a Tag) bool {
	return t.LangID == a.LangID && t.ScriptID == a.ScriptID && t.RegionID == a.RegionID
}


func (t Tag) IsRoot() bool {
	if int(t.pVariant) < len(t.str) {
		return false
	}
	return t.equalTags(Und)
}



func (t Tag) IsPrivateUse() bool {
	return t.str != "" && t.pVariant == 0
}




func (t *Tag) RemakeString() {
	if t.str == "" {
		return
	}
	extra := t.str[t.pVariant:]
	if t.pVariant > 0 {
		extra = extra[1:]
	}
	if t.equalTags(Und) && strings.HasPrefix(extra, "x-") {
		t.str = extra
		t.pVariant = 0
		t.pExt = 0
		return
	}
	var buf [max99thPercentileSize]byte 
	b := buf[:t.genCoreBytes(buf[:])]
	if extra != "" {
		diff := len(b) - int(t.pVariant)
		b = append(b, '-')
		b = append(b, extra...)
		t.pVariant = uint8(int(t.pVariant) + diff)
		t.pExt = uint16(int(t.pExt) + diff)
	} else {
		t.pVariant = uint8(len(b))
		t.pExt = uint16(len(b))
	}
	t.str = string(b)
}




func (t *Tag) genCoreBytes(buf []byte) int {
	n := t.LangID.StringToBuf(buf[:])
	if t.ScriptID != 0 {
		n += copy(buf[n:], "-")
		n += copy(buf[n:], t.ScriptID.String())
	}
	if t.RegionID != 0 {
		n += copy(buf[n:], "-")
		n += copy(buf[n:], t.RegionID.String())
	}
	return n
}


func (t Tag) String() string {
	if t.str != "" {
		return t.str
	}
	if t.ScriptID == 0 && t.RegionID == 0 {
		return t.LangID.String()
	}
	buf := [maxCoreSize]byte{}
	return string(buf[:t.genCoreBytes(buf[:])])
}


func (t Tag) MarshalText() (text []byte, err error) {
	if t.str != "" {
		text = append(text, t.str...)
	} else if t.ScriptID == 0 && t.RegionID == 0 {
		text = append(text, t.LangID.String()...)
	} else {
		buf := [maxCoreSize]byte{}
		text = buf[:t.genCoreBytes(buf[:])]
	}
	return text, nil
}


func (t *Tag) UnmarshalText(text []byte) error {
	tag, err := Parse(string(text))
	*t = tag
	return err
}



func (t Tag) Variants() string {
	if t.pVariant == 0 {
		return ""
	}
	return t.str[t.pVariant:t.pExt]
}


func (t Tag) VariantOrPrivateUseTags() string {
	if t.pExt > 0 {
		return t.str[t.pVariant:t.pExt]
	}
	return t.str[t.pVariant:]
}



func (t Tag) HasString() bool {
	return t.str != ""
}




func (t Tag) Parent() Tag {
	if t.str != "" {
		
		b, s, r := t.Raw()
		t = Tag{LangID: b, ScriptID: s, RegionID: r}
		if t.RegionID == 0 && t.ScriptID != 0 && t.LangID != 0 {
			base, _ := addTags(Tag{LangID: t.LangID})
			if base.ScriptID == t.ScriptID {
				return Tag{LangID: t.LangID}
			}
		}
		return t
	}
	if t.LangID != 0 {
		if t.RegionID != 0 {
			maxScript := t.ScriptID
			if maxScript == 0 {
				max, _ := addTags(t)
				maxScript = max.ScriptID
			}

			for i := range parents {
				if Language(parents[i].lang) == t.LangID && Script(parents[i].maxScript) == maxScript {
					for _, r := range parents[i].fromRegion {
						if Region(r) == t.RegionID {
							return Tag{
								LangID:   t.LangID,
								ScriptID: Script(parents[i].script),
								RegionID: Region(parents[i].toRegion),
							}
						}
					}
				}
			}

			
			base, _ := addTags(Tag{LangID: t.LangID})
			if base.ScriptID != maxScript {
				return Tag{LangID: t.LangID, ScriptID: maxScript}
			}
			return Tag{LangID: t.LangID}
		} else if t.ScriptID != 0 {
			
			
			base, _ := addTags(Tag{LangID: t.LangID})
			if base.ScriptID != t.ScriptID {
				return Und
			}
			return Tag{LangID: t.LangID}
		}
	}
	return Und
}


func ParseExtension(s string) (ext string, err error) {
	defer func() {
		if recover() != nil {
			ext = ""
			err = ErrSyntax
		}
	}()

	scan := makeScannerString(s)
	var end int
	if n := len(scan.token); n != 1 {
		return "", ErrSyntax
	}
	scan.toLower(0, len(scan.b))
	end = parseExtension(&scan)
	if end != len(s) {
		return "", ErrSyntax
	}
	return string(scan.b), nil
}


func (t Tag) HasVariants() bool {
	return uint16(t.pVariant) < t.pExt
}


func (t Tag) HasExtensions() bool {
	return int(t.pExt) < len(t.str)
}




func (t Tag) Extension(x byte) (ext string, ok bool) {
	for i := int(t.pExt); i < len(t.str)-1; {
		var ext string
		i, ext = getExtension(t.str, i)
		if ext[0] == x {
			return ext, true
		}
	}
	return "", false
}


func (t Tag) Extensions() []string {
	e := []string{}
	for i := int(t.pExt); i < len(t.str)-1; {
		var ext string
		i, ext = getExtension(t.str, i)
		e = append(e, ext)
	}
	return e
}









func (t Tag) TypeForKey(key string) string {
	if _, start, end, _ := t.findTypeForKey(key); end != start {
		s := t.str[start:end]
		if p := strings.IndexByte(s, '-'); p >= 0 {
			s = s[:p]
		}
		return s
	}
	return ""
}

var (
	errPrivateUse       = errors.New("cannot set a key on a private use tag")
	errInvalidArguments = errors.New("invalid key or type")
)





func (t Tag) SetTypeForKey(key, value string) (Tag, error) {
	if t.IsPrivateUse() {
		return t, errPrivateUse
	}
	if len(key) != 2 {
		return t, errInvalidArguments
	}

	
	if value == "" {
		start, sep, end, _ := t.findTypeForKey(key)
		if start != sep {
			
			switch {
			case t.str[start-2] != '-': 
			case end == len(t.str), 
				end+2 < len(t.str) && t.str[end+2] == '-': 
				start -= 2
			}
			if start == int(t.pVariant) && end == len(t.str) {
				t.str = ""
				t.pVariant, t.pExt = 0, 0
			} else {
				t.str = fmt.Sprintf("%s%s", t.str[:start], t.str[end:])
			}
		}
		return t, nil
	}

	if len(value) < 3 || len(value) > 8 {
		return t, errInvalidArguments
	}

	var (
		buf    [maxCoreSize + maxSimpleUExtensionSize]byte
		uStart int 
	)

	
	if t.str == "" {
		uStart = t.genCoreBytes(buf[:])
		buf[uStart] = '-'
		uStart++
	}

	
	b := buf[uStart:]
	copy(b, "u-")
	copy(b[2:], key)
	b[4] = '-'
	b = b[:5+copy(b[5:], value)]
	scan := makeScanner(b)
	if parseExtensions(&scan); scan.err != nil {
		return t, scan.err
	}

	
	if t.str == "" {
		t.pVariant, t.pExt = byte(uStart-1), uint16(uStart-1)
		t.str = string(buf[:uStart+len(b)])
	} else {
		s := t.str
		start, sep, end, hasExt := t.findTypeForKey(key)
		if start == sep {
			if hasExt {
				b = b[2:]
			}
			t.str = fmt.Sprintf("%s-%s%s", s[:sep], b, s[end:])
		} else {
			t.str = fmt.Sprintf("%s-%s%s", s[:start+3], value, s[end:])
		}
	}
	return t, nil
}






func (t Tag) findTypeForKey(key string) (start, sep, end int, hasExt bool) {
	p := int(t.pExt)
	if len(key) != 2 || p == len(t.str) || p == 0 {
		return p, p, p, false
	}
	s := t.str

	
	for p++; s[p] != 'u'; p++ {
		if s[p] > 'u' {
			p--
			return p, p, p, false
		}
		if p = nextExtension(s, p); p == len(s) {
			return len(s), len(s), len(s), false
		}
	}
	
	p++

	
	curKey := ""

	
	for {
		end = p
		for p++; p < len(s) && s[p] != '-'; p++ {
		}
		n := p - end - 1
		if n <= 2 && curKey == key {
			if sep < end {
				sep++
			}
			return start, sep, end, true
		}
		switch n {
		case 0, 
			1: 
			return end, end, end, true
		case 2:
			
			curKey = s[end+1 : p]
			if curKey > key {
				return end, end, end, true
			}
			start = end
			sep = p
		}
	}
}




func ParseBase(s string) (l Language, err error) {
	defer func() {
		if recover() != nil {
			l = 0
			err = ErrSyntax
		}
	}()

	if n := len(s); n < 2 || 3 < n {
		return 0, ErrSyntax
	}
	var buf [3]byte
	return getLangID(buf[:copy(buf[:], s)])
}




func ParseScript(s string) (scr Script, err error) {
	defer func() {
		if recover() != nil {
			scr = 0
			err = ErrSyntax
		}
	}()

	if len(s) != 4 {
		return 0, ErrSyntax
	}
	var buf [4]byte
	return getScriptID(script, buf[:copy(buf[:], s)])
}



func EncodeM49(r int) (Region, error) {
	return getRegionM49(r)
}




func ParseRegion(s string) (r Region, err error) {
	defer func() {
		if recover() != nil {
			r = 0
			err = ErrSyntax
		}
	}()

	if n := len(s); n < 2 || 3 < n {
		return 0, ErrSyntax
	}
	var buf [3]byte
	return getRegionID(buf[:copy(buf[:], s)])
}



func (r Region) IsCountry() bool {
	if r == 0 || r.IsGroup() || r.IsPrivateUse() && r != _XK {
		return false
	}
	return true
}



func (r Region) IsGroup() bool {
	if r == 0 {
		return false
	}
	return int(regionInclusion[r]) < len(regionContainment)
}



func (r Region) Contains(c Region) bool {
	if r == c {
		return true
	}
	g := regionInclusion[r]
	if g >= nRegionGroups {
		return false
	}
	m := regionContainment[g]

	d := regionInclusion[c]
	b := regionInclusionBits[d]

	
	
	
	if d >= nRegionGroups {
		return b&m != 0
	}
	return b&^m == 0
}

var errNoTLD = errors.New("language: region is not a valid ccTLD")








func (r Region) TLD() (Region, error) {
	
	
	if r == _GB {
		r = _UK
	}
	if (r.typ() & ccTLD) == 0 {
		return 0, errNoTLD
	}
	return r, nil
}




func (r Region) Canonicalize() Region {
	if cr := normRegion(r); cr != 0 {
		return cr
	}
	return r
}


type Variant struct {
	ID  uint8
	str string
}



func ParseVariant(s string) (v Variant, err error) {
	defer func() {
		if recover() != nil {
			v = Variant{}
			err = ErrSyntax
		}
	}()

	s = strings.ToLower(s)
	if id, ok := variantIndex[s]; ok {
		return Variant{id, s}, nil
	}
	return Variant{}, NewValueError([]byte(s))
}


func (v Variant) String() string {
	return v.str
}
