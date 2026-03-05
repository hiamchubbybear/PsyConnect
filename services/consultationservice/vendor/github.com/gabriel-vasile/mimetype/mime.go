package mimetype

import (
	"mime"

	"github.com/gabriel-vasile/mimetype/internal/charset"
	"github.com/gabriel-vasile/mimetype/internal/magic"
)



type MIME struct {
	mime      string
	aliases   []string
	extension string
	
	
	detector magic.Detector
	children []*MIME
	parent   *MIME
}


func (m *MIME) String() string {
	return m.mime
}




func (m *MIME) Extension() string {
	return m.extension
}








func (m *MIME) Parent() *MIME {
	return m.parent
}





func (m *MIME) Is(expectedMIME string) bool {
	
	
	expectedMIME, _, _ = mime.ParseMediaType(expectedMIME)
	found, _, _ := mime.ParseMediaType(m.mime)

	if expectedMIME == found {
		return true
	}

	for _, alias := range m.aliases {
		if alias == expectedMIME {
			return true
		}
	}

	return false
}

func newMIME(
	mime, extension string,
	detector magic.Detector,
	children ...*MIME) *MIME {
	m := &MIME{
		mime:      mime,
		extension: extension,
		detector:  detector,
		children:  children,
	}

	for _, c := range children {
		c.parent = m
	}

	return m
}

func (m *MIME) alias(aliases ...string) *MIME {
	m.aliases = aliases
	return m
}



func (m *MIME) match(in []byte, readLimit uint32) *MIME {
	for _, c := range m.children {
		if c.detector(in, readLimit) {
			return c.match(in, readLimit)
		}
	}

	needsCharset := map[string]func([]byte) string{
		"text/plain": charset.FromPlain,
		"text/html":  charset.FromHTML,
		"text/xml":   charset.FromXML,
	}
	
	ps := map[string]string{}
	if f, ok := needsCharset[m.mime]; ok {
		if cset := f(in); cset != "" {
			ps["charset"] = cset
		}
	}

	return m.cloneHierarchy(ps)
}


func (m *MIME) flatten() []*MIME {
	out := []*MIME{m}
	for _, c := range m.children {
		out = append(out, c.flatten()...)
	}

	return out
}


func (m *MIME) clone(ps map[string]string) *MIME {
	clonedMIME := m.mime
	if len(ps) > 0 {
		clonedMIME = mime.FormatMediaType(m.mime, ps)
	}

	return &MIME{
		mime:      clonedMIME,
		aliases:   m.aliases,
		extension: m.extension,
	}
}



func (m *MIME) cloneHierarchy(ps map[string]string) *MIME {
	ret := m.clone(ps)
	lastChild := ret
	for p := m.Parent(); p != nil; p = p.Parent() {
		pClone := p.clone(nil)
		lastChild.parent = pClone
		lastChild = pClone
	}

	return ret
}

func (m *MIME) lookup(mime string) *MIME {
	for _, n := range append(m.aliases, m.mime) {
		if n == mime {
			return m
		}
	}

	for _, c := range m.children {
		if m := c.lookup(mime); m != nil {
			return m
		}
	}
	return nil
}





func (m *MIME) Extend(detector func(raw []byte, limit uint32) bool, mime, extension string, aliases ...string) {
	c := &MIME{
		mime:      mime,
		extension: extension,
		detector:  detector,
		parent:    m,
		aliases:   aliases,
	}

	mu.Lock()
	m.children = append([]*MIME{c}, m.children...)
	mu.Unlock()
}
