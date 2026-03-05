

package ast

import (
    `github.com/bytedance/sonic/internal/rt`
    `github.com/bytedance/sonic/internal/native/types`
)


type SearchOptions struct {
    
    ValidateJSON bool

    
    
    CopyReturn bool

    
    
    ConcurrentRead bool
}

type Searcher struct {
    parser Parser
    SearchOptions
}

func NewSearcher(str string) *Searcher {
    return &Searcher{
        parser: Parser{
            s:      str,
            noLazy: false,
        },
        SearchOptions: SearchOptions{
            ValidateJSON: true,
        },
    }
}


func (self *Searcher) GetByPathCopy(path ...interface{}) (Node, error) {
    self.CopyReturn = true
    return self.getByPath(path...)
}





func (self *Searcher) GetByPath(path ...interface{}) (Node, error) {
    return self.getByPath(path...)
}

func (self *Searcher) getByPath(path ...interface{}) (Node, error) {
    var err types.ParsingError
    var start int

    self.parser.p = 0
    start, err = self.parser.getByPath(self.ValidateJSON, path...)
    if err != 0 {
        
        if err == types.ERR_NOT_FOUND {
            return Node{}, ErrNotExist
        }
        if err == types.ERR_UNSUPPORT_TYPE {
            panic("path must be either int(>=0) or string")
        }
        return Node{}, self.parser.syntaxError(err)
    }

    t := switchRawType(self.parser.s[start])
    if t == _V_NONE {
        return Node{}, self.parser.ExportError(err)
    }

    
    var raw string
    if self.CopyReturn {
        raw = rt.Mem2Str([]byte(self.parser.s[start:self.parser.p]))
    } else {
        raw = self.parser.s[start:self.parser.p]
    }
    return newRawNode(raw, t, self.ConcurrentRead), nil
}


func _GetByPath(src string, path ...interface{}) (start int, end int, typ int, err error) {
	p := NewParserObj(src)
	s, e := p.getByPath(false, path...)
	if e != 0 {
		
		if e == types.ERR_NOT_FOUND {
			return -1, -1, 0, ErrNotExist
		}
		if e == types.ERR_UNSUPPORT_TYPE {
			panic("path must be either int(>=0) or string")
		}
		return -1, -1, 0, p.syntaxError(e)
	}

	t := switchRawType(p.s[s])
	if t == _V_NONE {
		return -1, -1, 0, ErrNotExist
	}
    if t == _V_NUMBER {
        p.p = 1 + backward(p.s, p.p-1)
    }
	return s, p.p, int(t), nil
}



func _ValidSyntax(json string) bool {
	p := NewParserObj(json)
    _, e := p.skip()
	if e != 0 {
        return false
    }
   if skipBlank(p.s, p.p) != -int(types.ERR_EOF) {
        return false
   }
   return true
}



func _SkipFast(src string, i int) (int, int, error) {
    p := NewParserObj(src)
    p.p = i
    s, e := p.skipFast()
    if e != 0 {
        return -1, -1, p.ExportError(e)
    }
    t := switchRawType(p.s[s])
	if t == _V_NONE {
		return -1, -1, ErrNotExist
	}
    if t == _V_NUMBER {
        p.p = 1 + backward(p.s, p.p-1)
    }
    return s, p.p, nil
}
