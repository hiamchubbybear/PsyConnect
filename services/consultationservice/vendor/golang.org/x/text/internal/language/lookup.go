



package language

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"

	"golang.org/x/text/internal/tag"
)



func findIndex(idx tag.Index, key []byte, form string) (index int, err error) {
	if !tag.FixCase(form, key) {
		return 0, ErrSyntax
	}
	i := idx.Index(key)
	if i == -1 {
		return 0, NewValueError(key)
	}
	return i, nil
}

func searchUint(imap []uint16, key uint16) int {
	return sort.Search(len(imap), func(i int) bool {
		return imap[i] >= key
	})
}

type Language uint16



func getLangID(s []byte) (Language, error) {
	if len(s) == 2 {
		return getLangISO2(s)
	}
	return getLangISO3(s)
}




func (id Language) Canonicalize() (Language, AliasType) {
	return normLang(id)
}


func normLang(id Language) (Language, AliasType) {
	k := sort.Search(len(AliasMap), func(i int) bool {
		return AliasMap[i].From >= uint16(id)
	})
	if k < len(AliasMap) && AliasMap[k].From == uint16(id) {
		return Language(AliasMap[k].To), AliasTypes[k]
	}
	return id, AliasTypeUnknown
}



func getLangISO2(s []byte) (Language, error) {
	if !tag.FixCase("zz", s) {
		return 0, ErrSyntax
	}
	if i := lang.Index(s); i != -1 && lang.Elem(i)[3] != 0 {
		return Language(i), nil
	}
	return 0, NewValueError(s)
}

const base = 'z' - 'a' + 1

func strToInt(s []byte) uint {
	v := uint(0)
	for i := 0; i < len(s); i++ {
		v *= base
		v += uint(s[i] - 'a')
	}
	return v
}



func intToStr(v uint, s []byte) {
	for i := len(s) - 1; i >= 0; i-- {
		s[i] = byte(v%base) + 'a'
		v /= base
	}
}



func getLangISO3(s []byte) (Language, error) {
	if tag.FixCase("und", s) {
		
		for i := lang.Index(s[:2]); i != -1; i = lang.Next(s[:2], i) {
			if e := lang.Elem(i); e[3] == 0 && e[2] == s[2] {
				
				
				
				id := Language(i)
				if id == nonCanonicalUnd {
					return 0, nil
				}
				return id, nil
			}
		}
		if i := altLangISO3.Index(s); i != -1 {
			return Language(altLangIndex[altLangISO3.Elem(i)[3]]), nil
		}
		n := strToInt(s)
		if langNoIndex[n/8]&(1<<(n%8)) != 0 {
			return Language(n) + langNoIndexOffset, nil
		}
		
		for i := lang.Index(s[:1]); i != -1; i = lang.Next(s[:1], i) {
			if e := lang.Elem(i); e[2] == s[1] && e[3] == s[2] {
				return Language(i), nil
			}
		}
		return 0, NewValueError(s)
	}
	return 0, ErrSyntax
}



func (id Language) StringToBuf(b []byte) int {
	if id >= langNoIndexOffset {
		intToStr(uint(id)-langNoIndexOffset, b[:3])
		return 3
	} else if id == 0 {
		return copy(b, "und")
	}
	l := lang[id<<2:]
	if l[3] == 0 {
		return copy(b, l[:3])
	}
	return copy(b, l[:2])
}




func (b Language) String() string {
	if b == 0 {
		return "und"
	} else if b >= langNoIndexOffset {
		b -= langNoIndexOffset
		buf := [3]byte{}
		intToStr(uint(b), buf[:])
		return string(buf[:])
	}
	l := lang.Elem(int(b))
	if l[3] == 0 {
		return l[:3]
	}
	return l[:2]
}


func (b Language) ISO3() string {
	if b == 0 || b >= langNoIndexOffset {
		return b.String()
	}
	l := lang.Elem(int(b))
	if l[3] == 0 {
		return l[:3]
	} else if l[2] == 0 {
		return altLangISO3.Elem(int(l[3]))[:3]
	}
	
	
	return l[0:1] + l[2:4]
}


func (b Language) IsPrivateUse() bool {
	return langPrivateStart <= b && b <= langPrivateEnd
}



func (b Language) SuppressScript() Script {
	if b < langNoIndexOffset {
		return Script(suppressScript[b])
	}
	return 0
}

type Region uint16



func getRegionID(s []byte) (Region, error) {
	if len(s) == 3 {
		if isAlpha(s[0]) {
			return getRegionISO3(s)
		}
		if i, err := strconv.ParseUint(string(s), 10, 10); err == nil {
			return getRegionM49(int(i))
		}
	}
	return getRegionISO2(s)
}



func getRegionISO2(s []byte) (Region, error) {
	i, err := findIndex(regionISO, s, "ZZ")
	if err != nil {
		return 0, err
	}
	return Region(i) + isoRegionOffset, nil
}



func getRegionISO3(s []byte) (Region, error) {
	if tag.FixCase("ZZZ", s) {
		for i := regionISO.Index(s[:1]); i != -1; i = regionISO.Next(s[:1], i) {
			if e := regionISO.Elem(i); e[2] == s[1] && e[3] == s[2] {
				return Region(i) + isoRegionOffset, nil
			}
		}
		for i := 0; i < len(altRegionISO3); i += 3 {
			if tag.Compare(altRegionISO3[i:i+3], s) == 0 {
				return Region(altRegionIDs[i/3]), nil
			}
		}
		return 0, NewValueError(s)
	}
	return 0, ErrSyntax
}

func getRegionM49(n int) (Region, error) {
	if 0 < n && n <= 999 {
		const (
			searchBits = 7
			regionBits = 9
			regionMask = 1<<regionBits - 1
		)
		idx := n >> searchBits
		buf := fromM49[m49Index[idx]:m49Index[idx+1]]
		val := uint16(n) << regionBits 
		i := sort.Search(len(buf), func(i int) bool {
			return buf[i] >= val
		})
		if r := fromM49[int(m49Index[idx])+i]; r&^regionMask == val {
			return Region(r & regionMask), nil
		}
	}
	var e ValueError
	fmt.Fprint(bytes.NewBuffer([]byte(e.v[:])), n)
	return 0, e
}




func normRegion(r Region) Region {
	m := regionOldMap
	k := sort.Search(len(m), func(i int) bool {
		return m[i].From >= uint16(r)
	})
	if k < len(m) && m[k].From == uint16(r) {
		return Region(m[k].To)
	}
	return 0
}

const (
	iso3166UserAssigned = 1 << iota
	ccTLD
	bcp47Region
)

func (r Region) typ() byte {
	return regionTypes[r]
}



func (r Region) String() string {
	if r < isoRegionOffset {
		if r == 0 {
			return "ZZ"
		}
		return fmt.Sprintf("%03d", r.M49())
	}
	r -= isoRegionOffset
	return regionISO.Elem(int(r))[:2]
}




func (r Region) ISO3() string {
	if r < isoRegionOffset {
		return "ZZZ"
	}
	r -= isoRegionOffset
	reg := regionISO.Elem(int(r))
	switch reg[2] {
	case 0:
		return altRegionISO3[reg[3]:][:3]
	case ' ':
		return "ZZZ"
	}
	return reg[0:1] + reg[2:4]
}



func (r Region) M49() int {
	return int(m49[r])
}




func (r Region) IsPrivateUse() bool {
	return r.typ()&iso3166UserAssigned != 0
}

type Script uint16



func getScriptID(idx tag.Index, s []byte) (Script, error) {
	i, err := findIndex(idx, s, "Zzzz")
	return Script(i), err
}



func (s Script) String() string {
	if s == 0 {
		return "Zzzz"
	}
	return script.Elem(int(s))
}


func (s Script) IsPrivateUse() bool {
	return _Qaaa <= s && s <= _Qabx
}

const (
	maxAltTaglen = len("en-US-POSIX")
	maxLen       = maxAltTaglen
)

var (
	
	
	grandfatheredMap = map[[maxLen]byte]int16{
		[maxLen]byte{'a', 'r', 't', '-', 'l', 'o', 'j', 'b', 'a', 'n'}: _jbo, 
		[maxLen]byte{'i', '-', 'a', 'm', 'i'}:                          _ami, 
		[maxLen]byte{'i', '-', 'b', 'n', 'n'}:                          _bnn, 
		[maxLen]byte{'i', '-', 'h', 'a', 'k'}:                          _hak, 
		[maxLen]byte{'i', '-', 'k', 'l', 'i', 'n', 'g', 'o', 'n'}:      _tlh, 
		[maxLen]byte{'i', '-', 'l', 'u', 'x'}:                          _lb,  
		[maxLen]byte{'i', '-', 'n', 'a', 'v', 'a', 'j', 'o'}:           _nv,  
		[maxLen]byte{'i', '-', 'p', 'w', 'n'}:                          _pwn, 
		[maxLen]byte{'i', '-', 't', 'a', 'o'}:                          _tao, 
		[maxLen]byte{'i', '-', 't', 'a', 'y'}:                          _tay, 
		[maxLen]byte{'i', '-', 't', 's', 'u'}:                          _tsu, 
		[maxLen]byte{'n', 'o', '-', 'b', 'o', 'k'}:                     _nb,  
		[maxLen]byte{'n', 'o', '-', 'n', 'y', 'n'}:                     _nn,  
		[maxLen]byte{'s', 'g', 'n', '-', 'b', 'e', '-', 'f', 'r'}:      _sfb, 
		[maxLen]byte{'s', 'g', 'n', '-', 'b', 'e', '-', 'n', 'l'}:      _vgt, 
		[maxLen]byte{'s', 'g', 'n', '-', 'c', 'h', '-', 'd', 'e'}:      _sgg, 
		[maxLen]byte{'z', 'h', '-', 'g', 'u', 'o', 'y', 'u'}:           _cmn, 
		[maxLen]byte{'z', 'h', '-', 'h', 'a', 'k', 'k', 'a'}:           _hak, 
		[maxLen]byte{'z', 'h', '-', 'm', 'i', 'n', '-', 'n', 'a', 'n'}: _nan, 
		[maxLen]byte{'z', 'h', '-', 'x', 'i', 'a', 'n', 'g'}:           _hsn, 

		
		
		[maxLen]byte{'c', 'e', 'l', '-', 'g', 'a', 'u', 'l', 'i', 's', 'h'}: -1, 
		[maxLen]byte{'e', 'n', '-', 'g', 'b', '-', 'o', 'e', 'd'}:           -2, 
		[maxLen]byte{'i', '-', 'd', 'e', 'f', 'a', 'u', 'l', 't'}:           -3, 
		[maxLen]byte{'i', '-', 'e', 'n', 'o', 'c', 'h', 'i', 'a', 'n'}:      -4, 
		[maxLen]byte{'i', '-', 'm', 'i', 'n', 'g', 'o'}:                     -5, 
		[maxLen]byte{'z', 'h', '-', 'm', 'i', 'n'}:                          -6, 

		
		[maxLen]byte{'r', 'o', 'o', 't'}:                                    0,  
		[maxLen]byte{'e', 'n', '-', 'u', 's', '-', 'p', 'o', 's', 'i', 'x'}: -7, 
	}

	altTagIndex = [...]uint8{0, 17, 31, 45, 61, 74, 86, 102}

	altTags = "xtg-x-cel-gaulishen-GB-oxendicten-x-i-defaultund-x-i-enochiansee-x-i-mingonan-x-zh-minen-US-u-va-posix"
)

func grandfathered(s [maxAltTaglen]byte) (t Tag, ok bool) {
	if v, ok := grandfatheredMap[s]; ok {
		if v < 0 {
			return Make(altTags[altTagIndex[-v-1]:altTagIndex[-v]]), true
		}
		t.LangID = Language(v)
		return t, true
	}
	return t, false
}
