
package magic

import (
	"bytes"
	"fmt"
)

type (
	
	
	
	
	Detector func(raw []byte, limit uint32) bool
	xmlSig   struct {
		
		localName []byte
		
		xmlns []byte
	}
)



func prefix(sigs ...[]byte) Detector {
	return func(raw []byte, limit uint32) bool {
		for _, s := range sigs {
			if bytes.HasPrefix(raw, s) {
				return true
			}
		}
		return false
	}
}



func offset(sig []byte, offset int) Detector {
	return func(raw []byte, limit uint32) bool {
		return len(raw) > offset && bytes.HasPrefix(raw[offset:], sig)
	}
}


func ciPrefix(sigs ...[]byte) Detector {
	return func(raw []byte, limit uint32) bool {
		for _, s := range sigs {
			if ciCheck(s, raw) {
				return true
			}
		}
		return false
	}
}
func ciCheck(sig, raw []byte) bool {
	if len(raw) < len(sig)+1 {
		return false
	}
	
	for i, b := range sig {
		db := raw[i]
		if 'A' <= b && b <= 'Z' {
			db &= 0xDF
		}
		if b != db {
			return false
		}
	}

	return true
}



func xml(sigs ...xmlSig) Detector {
	return func(raw []byte, limit uint32) bool {
		raw = trimLWS(raw)
		if len(raw) == 0 {
			return false
		}
		for _, s := range sigs {
			if xmlCheck(s, raw) {
				return true
			}
		}
		return false
	}
}
func xmlCheck(sig xmlSig, raw []byte) bool {
	raw = raw[:min(len(raw), 512)]

	if len(sig.localName) == 0 {
		return bytes.Index(raw, sig.xmlns) > 0
	}
	if len(sig.xmlns) == 0 {
		return bytes.Index(raw, sig.localName) > 0
	}

	localNameIndex := bytes.Index(raw, sig.localName)
	return localNameIndex != -1 && localNameIndex < bytes.Index(raw, sig.xmlns)
}



func markup(sigs ...[]byte) Detector {
	return func(raw []byte, limit uint32) bool {
		if bytes.HasPrefix(raw, []byte{0xEF, 0xBB, 0xBF}) {
			
			
			
			raw = trimLWS(raw[3:])
		} else {
			raw = trimLWS(raw)
		}
		if len(raw) == 0 {
			return false
		}
		for _, s := range sigs {
			if markupCheck(s, raw) {
				return true
			}
		}
		return false
	}
}
func markupCheck(sig, raw []byte) bool {
	if len(raw) < len(sig)+1 {
		return false
	}

	
	for i, b := range sig {
		db := raw[i]
		if 'A' <= b && b <= 'Z' {
			db &= 0xDF
		}
		if b != db {
			return false
		}
	}
	
	if db := raw[len(sig)]; db != ' ' && db != '>' {
		return false
	}

	return true
}



func ftyp(sigs ...[]byte) Detector {
	return func(raw []byte, limit uint32) bool {
		if len(raw) < 12 {
			return false
		}
		for _, s := range sigs {
			if bytes.Equal(raw[8:12], s) {
				return true
			}
		}
		return false
	}
}

func newXMLSig(localName, xmlns string) xmlSig {
	ret := xmlSig{xmlns: []byte(xmlns)}
	if localName != "" {
		ret.localName = []byte(fmt.Sprintf("<%s", localName))
	}

	return ret
}











func shebang(sigs ...[]byte) Detector {
	return func(raw []byte, limit uint32) bool {
		for _, s := range sigs {
			if shebangCheck(s, firstLine(raw)) {
				return true
			}
		}
		return false
	}
}

func shebangCheck(sig, raw []byte) bool {
	if len(raw) < len(sig)+2 {
		return false
	}
	if raw[0] != '#' || raw[1] != '!' {
		return false
	}

	return bytes.Equal(trimLWS(trimRWS(raw[2:])), sig)
}


func trimLWS(in []byte) []byte {
	firstNonWS := 0
	for ; firstNonWS < len(in) && isWS(in[firstNonWS]); firstNonWS++ {
	}

	return in[firstNonWS:]
}


func trimRWS(in []byte) []byte {
	lastNonWS := len(in) - 1
	for ; lastNonWS > 0 && isWS(in[lastNonWS]); lastNonWS-- {
	}

	return in[:lastNonWS+1]
}

func firstLine(in []byte) []byte {
	lineEnd := 0
	for ; lineEnd < len(in) && in[lineEnd] != '\n'; lineEnd++ {
	}

	return in[:lineEnd]
}

func isWS(b byte) bool {
	return b == '\t' || b == '\n' || b == '\x0c' || b == '\r' || b == ' '
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type readBuf []byte

func (b *readBuf) advance(n int) bool {
	if n < 0 || len(*b) < n {
		return false
	}
	*b = (*b)[n:]
	return true
}
