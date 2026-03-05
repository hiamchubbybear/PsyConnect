



package language

import (
	"bytes"
	"errors"
	"fmt"
	"sort"

	"golang.org/x/text/internal/tag"
)



func isAlpha(b byte) bool {
	return b > '9'
}


func isAlphaNum(s []byte) bool {
	for _, c := range s {
		if !('a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9') {
			return false
		}
	}
	return true
}




var ErrSyntax = errors.New("language: tag is not well-formed")



var ErrDuplicateKey = errors.New("language: different values for same key in -u extension")




type ValueError struct {
	v [8]byte
}


func NewValueError(tag []byte) ValueError {
	var e ValueError
	copy(e.v[:], tag)
	return e
}

func (e ValueError) tag() []byte {
	n := bytes.IndexByte(e.v[:], 0)
	if n == -1 {
		n = 8
	}
	return e.v[:n]
}


func (e ValueError) Error() string {
	return fmt.Sprintf("language: subtag %q is well-formed but unknown", e.tag())
}


func (e ValueError) Subtag() string {
	return string(e.tag())
}


type scanner struct {
	b     []byte
	bytes [max99thPercentileSize]byte
	token []byte
	start int 
	end   int 
	next  int 
	err   error
	done  bool
}

func makeScannerString(s string) scanner {
	scan := scanner{}
	if len(s) <= len(scan.bytes) {
		scan.b = scan.bytes[:copy(scan.bytes[:], s)]
	} else {
		scan.b = []byte(s)
	}
	scan.init()
	return scan
}



func makeScanner(b []byte) scanner {
	scan := scanner{b: b}
	scan.init()
	return scan
}

func (s *scanner) init() {
	for i, c := range s.b {
		if c == '_' {
			s.b[i] = '-'
		}
	}
	s.scan()
}


func (s *scanner) toLower(start, end int) {
	for i := start; i < end; i++ {
		c := s.b[i]
		if 'A' <= c && c <= 'Z' {
			s.b[i] += 'a' - 'A'
		}
	}
}

func (s *scanner) setError(e error) {
	if s.err == nil || (e == ErrSyntax && s.err != ErrSyntax) {
		s.err = e
	}
}




func (s *scanner) resizeRange(oldStart, oldEnd, newSize int) {
	s.start = oldStart
	if end := oldStart + newSize; end != oldEnd {
		diff := end - oldEnd
		var b []byte
		if n := len(s.b) + diff; n > cap(s.b) {
			b = make([]byte, n)
			copy(b, s.b[:oldStart])
		} else {
			b = s.b[:n]
		}
		copy(b[end:], s.b[oldEnd:])
		s.b = b
		s.next = end + (s.next - s.end)
		s.end = end
	}
}


func (s *scanner) replace(repl string) {
	s.resizeRange(s.start, s.end, len(repl))
	copy(s.b[s.start:], repl)
}



func (s *scanner) gobble(e error) {
	s.setError(e)
	if s.start == 0 {
		s.b = s.b[:+copy(s.b, s.b[s.next:])]
		s.end = 0
	} else {
		s.b = s.b[:s.start-1+copy(s.b[s.start-1:], s.b[s.end:])]
		s.end = s.start - 1
	}
	s.next = s.start
}


func (s *scanner) deleteRange(start, end int) {
	s.b = s.b[:start+copy(s.b[start:], s.b[end:])]
	diff := end - start
	s.next -= diff
	s.start -= diff
	s.end -= diff
}





func (s *scanner) scan() (end int) {
	end = s.end
	s.token = nil
	for s.start = s.next; s.next < len(s.b); {
		i := bytes.IndexByte(s.b[s.next:], '-')
		if i == -1 {
			s.end = len(s.b)
			s.next = len(s.b)
			i = s.end - s.start
		} else {
			s.end = s.next + i
			s.next = s.end + 1
		}
		token := s.b[s.start:s.end]
		if i < 1 || i > 8 || !isAlphaNum(token) {
			s.gobble(ErrSyntax)
			continue
		}
		s.token = token
		return end
	}
	if n := len(s.b); n > 0 && s.b[n-1] == '-' {
		s.setError(ErrSyntax)
		s.b = s.b[:len(s.b)-1]
	}
	s.done = true
	return end
}



func (s *scanner) acceptMinSize(min int) (end int) {
	end = s.end
	s.scan()
	for ; len(s.token) >= min; s.scan() {
		end = s.end
	}
	return end
}








func Parse(s string) (t Tag, err error) {
	
	if s == "" {
		return Und, ErrSyntax
	}
	defer func() {
		if recover() != nil {
			t = Und
			err = ErrSyntax
			return
		}
	}()
	if len(s) <= maxAltTaglen {
		b := [maxAltTaglen]byte{}
		for i, c := range s {
			
			if 'A' <= c && c <= 'Z' {
				c += 'a' - 'A'
			} else if c == '_' {
				c = '-'
			}
			b[i] = byte(c)
		}
		if t, ok := grandfathered(b); ok {
			return t, nil
		}
	}
	scan := makeScannerString(s)
	return parse(&scan, s)
}

func parse(scan *scanner, s string) (t Tag, err error) {
	t = Und
	var end int
	if n := len(scan.token); n <= 1 {
		scan.toLower(0, len(scan.b))
		if n == 0 || scan.token[0] != 'x' {
			return t, ErrSyntax
		}
		end = parseExtensions(scan)
	} else if n >= 4 {
		return Und, ErrSyntax
	} else { 
		t, end = parseTag(scan, true)
		if n := len(scan.token); n == 1 {
			t.pExt = uint16(end)
			end = parseExtensions(scan)
		} else if end < len(scan.b) {
			scan.setError(ErrSyntax)
			scan.b = scan.b[:end]
		}
	}
	if int(t.pVariant) < len(scan.b) {
		if end < len(s) {
			s = s[:end]
		}
		if len(s) > 0 && tag.Compare(s, scan.b) == 0 {
			t.str = s
		} else {
			t.str = string(scan.b)
		}
	} else {
		t.pVariant, t.pExt = 0, 0
	}
	return t, scan.err
}




func parseTag(scan *scanner, doNorm bool) (t Tag, end int) {
	var e error
	
	t.LangID, e = getLangID(scan.token)
	scan.setError(e)
	scan.replace(t.LangID.String())
	langStart := scan.start
	end = scan.scan()
	for len(scan.token) == 3 && isAlpha(scan.token[0]) {
		
		
		if doNorm {
			lang, e := getLangID(scan.token)
			if lang != 0 {
				t.LangID = lang
				langStr := lang.String()
				copy(scan.b[langStart:], langStr)
				scan.b[langStart+len(langStr)] = '-'
				scan.start = langStart + len(langStr) + 1
			}
			scan.gobble(e)
		}
		end = scan.scan()
	}
	if len(scan.token) == 4 && isAlpha(scan.token[0]) {
		t.ScriptID, e = getScriptID(script, scan.token)
		if t.ScriptID == 0 {
			scan.gobble(e)
		}
		end = scan.scan()
	}
	if n := len(scan.token); n >= 2 && n <= 3 {
		t.RegionID, e = getRegionID(scan.token)
		if t.RegionID == 0 {
			scan.gobble(e)
		} else {
			scan.replace(t.RegionID.String())
		}
		end = scan.scan()
	}
	scan.toLower(scan.start, len(scan.b))
	t.pVariant = byte(end)
	end = parseVariants(scan, end, t)
	t.pExt = uint16(end)
	return t, end
}

var separator = []byte{'-'}



func parseVariants(scan *scanner, end int, t Tag) int {
	start := scan.start
	varIDBuf := [4]uint8{}
	variantBuf := [4][]byte{}
	varID := varIDBuf[:0]
	variant := variantBuf[:0]
	last := -1
	needSort := false
	for ; len(scan.token) >= 4; scan.scan() {
		
		
		v, ok := variantIndex[string(scan.token)]
		if !ok {
			
			
			scan.gobble(NewValueError(scan.token))
			continue
		}
		varID = append(varID, v)
		variant = append(variant, scan.token)
		if !needSort {
			if last < int(v) {
				last = int(v)
			} else {
				needSort = true
				
				
				const maxVariants = 8
				if len(varID) > maxVariants {
					break
				}
			}
		}
		end = scan.end
	}
	if needSort {
		sort.Sort(variantsSort{varID, variant})
		k, l := 0, -1
		for i, v := range varID {
			w := int(v)
			if l == w {
				
				continue
			}
			varID[k] = varID[i]
			variant[k] = variant[i]
			k++
			l = w
		}
		if str := bytes.Join(variant[:k], separator); len(str) == 0 {
			end = start - 1
		} else {
			scan.resizeRange(start, end, len(str))
			copy(scan.b[scan.start:], str)
			end = scan.end
		}
	}
	return end
}

type variantsSort struct {
	i []uint8
	v [][]byte
}

func (s variantsSort) Len() int {
	return len(s.i)
}

func (s variantsSort) Swap(i, j int) {
	s.i[i], s.i[j] = s.i[j], s.i[i]
	s.v[i], s.v[j] = s.v[j], s.v[i]
}

func (s variantsSort) Less(i, j int) bool {
	return s.i[i] < s.i[j]
}

type bytesSort struct {
	b [][]byte
	n int 
}

func (b bytesSort) Len() int {
	return len(b.b)
}

func (b bytesSort) Swap(i, j int) {
	b.b[i], b.b[j] = b.b[j], b.b[i]
}

func (b bytesSort) Less(i, j int) bool {
	for k := 0; k < b.n; k++ {
		if b.b[i][k] == b.b[j][k] {
			continue
		}
		return b.b[i][k] < b.b[j][k]
	}
	return false
}




func parseExtensions(scan *scanner) int {
	start := scan.start
	exts := [][]byte{}
	private := []byte{}
	end := scan.end
	for len(scan.token) == 1 {
		extStart := scan.start
		ext := scan.token[0]
		end = parseExtension(scan)
		extension := scan.b[extStart:end]
		if len(extension) < 3 || (ext != 'x' && len(extension) < 4) {
			scan.setError(ErrSyntax)
			end = extStart
			continue
		} else if start == extStart && (ext == 'x' || scan.start == len(scan.b)) {
			scan.b = scan.b[:end]
			return end
		} else if ext == 'x' {
			private = extension
			break
		}
		exts = append(exts, extension)
	}
	sort.Sort(bytesSort{exts, 1})
	if len(private) > 0 {
		exts = append(exts, private)
	}
	scan.b = scan.b[:start]
	if len(exts) > 0 {
		scan.b = append(scan.b, bytes.Join(exts, separator)...)
	} else if start > 0 {
		
		scan.b = scan.b[:start-1]
	}
	return end
}



func parseExtension(scan *scanner) int {
	start, end := scan.start, scan.end
	switch scan.token[0] {
	case 'u': 
		attrStart := end
		scan.scan()
		for last := []byte{}; len(scan.token) > 2; scan.scan() {
			if bytes.Compare(scan.token, last) != -1 {
				
				p := attrStart + 1
				scan.next = p
				attrs := [][]byte{}
				for scan.scan(); len(scan.token) > 2; scan.scan() {
					attrs = append(attrs, scan.token)
					end = scan.end
				}
				sort.Sort(bytesSort{attrs, 3})
				copy(scan.b[p:], bytes.Join(attrs, separator))
				break
			}
			last = scan.token
			end = scan.end
		}
		
		
		var last, key []byte
		for attrEnd := end; len(scan.token) == 2; last = key {
			key = scan.token
			end = scan.end
			for scan.scan(); end < scan.end && len(scan.token) > 2; scan.scan() {
				end = scan.end
			}
			
			if bytes.Compare(key, last) != 1 || scan.err != nil {
				
				
				p := attrEnd + 1
				scan.next = p
				keys := [][]byte{}
				for scan.scan(); len(scan.token) == 2; {
					keyStart := scan.start
					end = scan.end
					for scan.scan(); end < scan.end && len(scan.token) > 2; scan.scan() {
						end = scan.end
					}
					keys = append(keys, scan.b[keyStart:end])
				}
				sort.Stable(bytesSort{keys, 2})
				if n := len(keys); n > 0 {
					k := 0
					for i := 1; i < n; i++ {
						if !bytes.Equal(keys[k][:2], keys[i][:2]) {
							k++
							keys[k] = keys[i]
						} else if !bytes.Equal(keys[k], keys[i]) {
							scan.setError(ErrDuplicateKey)
						}
					}
					keys = keys[:k+1]
				}
				reordered := bytes.Join(keys, separator)
				if e := p + len(reordered); e < end {
					scan.deleteRange(e, end)
					end = e
				}
				copy(scan.b[p:], reordered)
				break
			}
		}
	case 't': 
		scan.scan()
		if n := len(scan.token); n >= 2 && n <= 3 && isAlpha(scan.token[1]) {
			_, end = parseTag(scan, false)
			scan.toLower(start, end)
		}
		for len(scan.token) == 2 && !isAlpha(scan.token[1]) {
			end = scan.acceptMinSize(3)
		}
	case 'x':
		end = scan.acceptMinSize(1)
	default:
		end = scan.acceptMinSize(2)
	}
	return end
}


func getExtension(s string, p int) (end int, ext string) {
	if s[p] == '-' {
		p++
	}
	if s[p] == 'x' {
		return len(s), s[p:]
	}
	end = nextExtension(s, p)
	return end, s[p:end]
}





func nextExtension(s string, p int) int {
	for n := len(s) - 3; p < n; {
		if s[p] == '-' {
			if s[p+2] == '-' {
				return p
			}
			p += 3
		} else {
			p++
		}
	}
	return len(s)
}
