



package language

import (
	"sort"
	"strings"
)



type Builder struct {
	Tag Tag

	private    string 
	variants   []string
	extensions []string
}


func (b *Builder) Make() Tag {
	t := b.Tag

	if len(b.extensions) > 0 || len(b.variants) > 0 {
		sort.Sort(sortVariants(b.variants))
		sort.Strings(b.extensions)

		if b.private != "" {
			b.extensions = append(b.extensions, b.private)
		}
		n := maxCoreSize + tokenLen(b.variants...) + tokenLen(b.extensions...)
		buf := make([]byte, n)
		p := t.genCoreBytes(buf)
		t.pVariant = byte(p)
		p += appendTokens(buf[p:], b.variants...)
		t.pExt = uint16(p)
		p += appendTokens(buf[p:], b.extensions...)
		t.str = string(buf[:p])
		
		
		scan := makeScanner(buf[:p])
		t, _ = parse(&scan, "")
		return t

	} else if b.private != "" {
		t.str = b.private
		t.RemakeString()
	}
	return t
}



func (b *Builder) SetTag(t Tag) {
	b.Tag.LangID = t.LangID
	b.Tag.RegionID = t.RegionID
	b.Tag.ScriptID = t.ScriptID
	
	b.variants = b.variants[:0]
	if variants := t.Variants(); variants != "" {
		for _, vr := range strings.Split(variants[1:], "-") {
			b.variants = append(b.variants, vr)
		}
	}
	b.extensions, b.private = b.extensions[:0], ""
	for _, e := range t.Extensions() {
		b.AddExt(e)
	}
}




func (b *Builder) AddExt(e string) {
	if e[0] == 'x' {
		if b.private == "" {
			b.private = e
		}
		return
	}
	for i, s := range b.extensions {
		if s[0] == e[0] {
			if e[0] == 'u' {
				b.extensions[i] += e[1:]
			}
			return
		}
	}
	b.extensions = append(b.extensions, e)
}





func (b *Builder) SetExt(e string) {
	if e[0] == 'x' {
		b.private = e
		return
	}
	for i, s := range b.extensions {
		if s[0] == e[0] {
			if e[0] == 'u' {
				b.extensions[i] = e + s[1:]
			} else {
				b.extensions[i] = e
			}
			return
		}
	}
	b.extensions = append(b.extensions, e)
}


func (b *Builder) AddVariant(v ...string) {
	for _, v := range v {
		if v != "" {
			b.variants = append(b.variants, v)
		}
	}
}



func (b *Builder) ClearVariants() {
	b.variants = b.variants[:0]
}



func (b *Builder) ClearExtensions() {
	b.private = ""
	b.extensions = b.extensions[:0]
}

func tokenLen(token ...string) (n int) {
	for _, t := range token {
		n += len(t) + 1
	}
	return
}

func appendTokens(b []byte, token ...string) int {
	p := 0
	for _, t := range token {
		b[p] = '-'
		copy(b[p+1:], t)
		p += 1 + len(t)
	}
	return p
}

type sortVariants []string

func (s sortVariants) Len() int {
	return len(s)
}

func (s sortVariants) Swap(i, j int) {
	s[j], s[i] = s[i], s[j]
}

func (s sortVariants) Less(i, j int) bool {
	return variantIndex[s[i]] < variantIndex[s[j]]
}
