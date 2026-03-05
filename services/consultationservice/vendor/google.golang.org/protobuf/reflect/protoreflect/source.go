



package protoreflect

import (
	"strconv"
)


type SourceLocations interface {
	
	Len() int
	
	Get(int) SourceLocation

	
	
	
	
	
	
	ByPath(path SourcePath) SourceLocation

	
	
	
	ByDescriptor(desc Descriptor) SourceLocation

	doNotImplement
}



type SourceLocation struct {
	
	
	Path SourcePath

	
	
	StartLine, StartColumn int
	
	
	
	
	EndLine, EndColumn int

	
	
	LeadingDetachedComments []string
	
	LeadingComments string
	
	TrailingComments string

	
	
	Next int
}






type SourcePath []int32


func (p1 SourcePath) Equal(p2 SourcePath) bool {
	if len(p1) != len(p2) {
		return false
	}
	for i := range p1 {
		if p1[i] != p2[i] {
			return false
		}
	}
	return true
}










func (p SourcePath) String() string {
	b := p.appendFileDescriptorProto(nil)
	for _, i := range p {
		b = append(b, '.')
		b = strconv.AppendInt(b, int64(i), 10)
	}
	return string(b)
}

type appendFunc func(*SourcePath, []byte) []byte

func (p *SourcePath) appendSingularField(b []byte, name string, f appendFunc) []byte {
	if len(*p) == 0 {
		return b
	}
	b = append(b, '.')
	b = append(b, name...)
	*p = (*p)[1:]
	if f != nil {
		b = f(p, b)
	}
	return b
}

func (p *SourcePath) appendRepeatedField(b []byte, name string, f appendFunc) []byte {
	b = p.appendSingularField(b, name, nil)
	if len(*p) == 0 || (*p)[0] < 0 {
		return b
	}
	b = append(b, '[')
	b = strconv.AppendUint(b, uint64((*p)[0]), 10)
	b = append(b, ']')
	*p = (*p)[1:]
	if f != nil {
		b = f(p, b)
	}
	return b
}
