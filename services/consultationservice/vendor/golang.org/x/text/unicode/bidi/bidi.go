



//go:generate go run gen.go gen_trieval.go gen_ranges.go







package bidi 





import (
	"bytes"
)







type Direction int

const (
	
	
	
	LeftToRight Direction = iota

	
	
	
	RightToLeft

	
	
	Mixed

	
	
	Neutral
)

type options struct {
	defaultDirection Direction
}


type Option func(*options)












func DefaultDirection(d Direction) Option {
	return func(opts *options) {
		opts.defaultDirection = d
	}
}


type Paragraph struct {
	p          []byte
	o          Ordering
	opts       []Option
	types      []Class
	pairTypes  []bracketType
	pairValues []rune
	runes      []rune
	options    options
}










func (p *Paragraph) prepareInput() (n int, err error) {
	p.runes = bytes.Runes(p.p)
	bytecount := 0
	
	p.pairTypes = nil
	p.pairValues = nil
	p.types = nil

	for _, r := range p.runes {
		props, i := LookupRune(r)
		bytecount += i
		cls := props.Class()
		if cls == B {
			return bytecount, nil
		}
		p.types = append(p.types, cls)
		if props.IsOpeningBracket() {
			p.pairTypes = append(p.pairTypes, bpOpen)
			p.pairValues = append(p.pairValues, r)
		} else if props.IsBracket() {
			
			
			p.pairTypes = append(p.pairTypes, bpClose)
			p.pairValues = append(p.pairValues, r)
		} else {
			p.pairTypes = append(p.pairTypes, bpNone)
			p.pairValues = append(p.pairValues, 0)
		}
	}
	return bytecount, nil
}






func (p *Paragraph) SetBytes(b []byte, opts ...Option) (n int, err error) {
	p.p = b
	p.opts = opts
	return p.prepareInput()
}






func (p *Paragraph) SetString(s string, opts ...Option) (n int, err error) {
	p.p = []byte(s)
	p.opts = opts
	return p.prepareInput()
}




func (p *Paragraph) IsLeftToRight() bool {
	return p.Direction() == LeftToRight
}




func (p *Paragraph) Direction() Direction {
	return p.o.Direction()
}






func (p *Paragraph) RunAt(pos int) Run {
	c := 0
	runNumber := 0
	for i, r := range p.o.runes {
		c += len(r)
		if pos < c {
			runNumber = i
		}
	}
	return p.o.Run(runNumber)
}

func calculateOrdering(levels []level, runes []rune) Ordering {
	var curDir Direction

	prevDir := Neutral
	prevI := 0

	o := Ordering{}
	
	
	for i, lvl := range levels {
		if lvl%2 == 0 {
			curDir = LeftToRight
		} else {
			curDir = RightToLeft
		}
		if curDir != prevDir {
			if i > 0 {
				o.runes = append(o.runes, runes[prevI:i])
				o.directions = append(o.directions, prevDir)
				o.startpos = append(o.startpos, prevI)
			}
			prevI = i
			prevDir = curDir
		}
	}
	o.runes = append(o.runes, runes[prevI:])
	o.directions = append(o.directions, prevDir)
	o.startpos = append(o.startpos, prevI)
	return o
}


func (p *Paragraph) Order() (Ordering, error) {
	if len(p.types) == 0 {
		return Ordering{}, nil
	}

	for _, fn := range p.opts {
		fn(&p.options)
	}
	lvl := level(-1)
	if p.options.defaultDirection == RightToLeft {
		lvl = 1
	}
	para, err := newParagraph(p.types, p.pairTypes, p.pairValues, lvl)
	if err != nil {
		return Ordering{}, err
	}

	levels := para.getLevels([]int{len(p.types)})

	p.o = calculateOrdering(levels, p.runes)
	return p.o, nil
}



func (p *Paragraph) Line(start, end int) (Ordering, error) {
	lineTypes := p.types[start:end]
	para, err := newParagraph(lineTypes, p.pairTypes[start:end], p.pairValues[start:end], -1)
	if err != nil {
		return Ordering{}, err
	}
	levels := para.getLevels([]int{len(lineTypes)})
	o := calculateOrdering(levels, p.runes[start:end])
	return o, nil
}




type Ordering struct {
	runes      [][]rune
	directions []Direction
	startpos   []int
}




func (o *Ordering) Direction() Direction {
	return o.directions[0]
}


func (o *Ordering) NumRuns() int {
	return len(o.runes)
}


func (o *Ordering) Run(i int) Run {
	r := Run{
		runes:     o.runes[i],
		direction: o.directions[i],
		startpos:  o.startpos[i],
	}
	return r
}









type Run struct {
	runes     []rune
	direction Direction
	startpos  int
}


func (r *Run) String() string {
	return string(r.runes)
}


func (r *Run) Bytes() []byte {
	return []byte(r.String())
}







func (r *Run) Direction() Direction {
	return r.direction
}



func (r *Run) Pos() (start, end int) {
	return r.startpos, r.startpos + len(r.runes) - 1
}




func AppendReverse(out, in []byte) []byte {
	ret := make([]byte, len(in)+len(out))
	copy(ret, out)
	inRunes := bytes.Runes(in)

	for i, r := range inRunes {
		prop, _ := LookupRune(r)
		if prop.IsBracket() {
			inRunes[i] = prop.reverseBracket(r)
		}
	}

	for i, j := 0, len(inRunes)-1; i < j; i, j = i+1, j-1 {
		inRunes[i], inRunes[j] = inRunes[j], inRunes[i]
	}
	copy(ret[len(out):], string(inRunes))

	return ret
}




func ReverseString(s string) string {
	input := []rune(s)
	li := len(input)
	ret := make([]rune, li)
	for i, r := range input {
		prop, _ := LookupRune(r)
		if prop.IsBracket() {
			ret[li-i-1] = prop.reverseBracket(r)
		} else {
			ret[li-i-1] = r
		}
	}
	return string(ret)
}
