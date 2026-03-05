



package filedesc

import (
	"fmt"
	"math"
	"sort"
	"sync"

	"google.golang.org/protobuf/internal/genid"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/internal/descfmt"
	"google.golang.org/protobuf/internal/errors"
	"google.golang.org/protobuf/internal/pragma"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type FileImports []protoreflect.FileImport

func (p *FileImports) Len() int                            { return len(*p) }
func (p *FileImports) Get(i int) protoreflect.FileImport   { return (*p)[i] }
func (p *FileImports) Format(s fmt.State, r rune)          { descfmt.FormatList(s, r, p) }
func (p *FileImports) ProtoInternal(pragma.DoNotImplement) {}

type Names struct {
	List []protoreflect.Name
	once sync.Once
	has  map[protoreflect.Name]int 
}

func (p *Names) Len() int                            { return len(p.List) }
func (p *Names) Get(i int) protoreflect.Name         { return p.List[i] }
func (p *Names) Has(s protoreflect.Name) bool        { return p.lazyInit().has[s] > 0 }
func (p *Names) Format(s fmt.State, r rune)          { descfmt.FormatList(s, r, p) }
func (p *Names) ProtoInternal(pragma.DoNotImplement) {}
func (p *Names) lazyInit() *Names {
	p.once.Do(func() {
		if len(p.List) > 0 {
			p.has = make(map[protoreflect.Name]int, len(p.List))
			for _, s := range p.List {
				p.has[s] = p.has[s] + 1
			}
		}
	})
	return p
}



func (p *Names) CheckValid() error {
	for s, n := range p.lazyInit().has {
		switch {
		case n > 1:
			return errors.New("duplicate name: %q", s)
		case false && !s.IsValid():
			
			
			return errors.New("invalid name: %q", s)
		}
	}
	return nil
}

type EnumRanges struct {
	List   [][2]protoreflect.EnumNumber 
	once   sync.Once
	sorted [][2]protoreflect.EnumNumber 
}

func (p *EnumRanges) Len() int                             { return len(p.List) }
func (p *EnumRanges) Get(i int) [2]protoreflect.EnumNumber { return p.List[i] }
func (p *EnumRanges) Has(n protoreflect.EnumNumber) bool {
	for ls := p.lazyInit().sorted; len(ls) > 0; {
		i := len(ls) / 2
		switch r := enumRange(ls[i]); {
		case n < r.Start():
			ls = ls[:i] 
		case n > r.End():
			ls = ls[i+1:] 
		default:
			return true
		}
	}
	return false
}
func (p *EnumRanges) Format(s fmt.State, r rune)          { descfmt.FormatList(s, r, p) }
func (p *EnumRanges) ProtoInternal(pragma.DoNotImplement) {}
func (p *EnumRanges) lazyInit() *EnumRanges {
	p.once.Do(func() {
		p.sorted = append(p.sorted, p.List...)
		sort.Slice(p.sorted, func(i, j int) bool {
			return p.sorted[i][0] < p.sorted[j][0]
		})
	})
	return p
}



func (p *EnumRanges) CheckValid() error {
	var rp enumRange
	for i, r := range p.lazyInit().sorted {
		r := enumRange(r)
		switch {
		case !(r.Start() <= r.End()):
			return errors.New("invalid range: %v", r)
		case !(rp.End() < r.Start()) && i > 0:
			return errors.New("overlapping ranges: %v with %v", rp, r)
		}
		rp = r
	}
	return nil
}

type enumRange [2]protoreflect.EnumNumber

func (r enumRange) Start() protoreflect.EnumNumber { return r[0] } 
func (r enumRange) End() protoreflect.EnumNumber   { return r[1] } 
func (r enumRange) String() string {
	if r.Start() == r.End() {
		return fmt.Sprintf("%d", r.Start())
	}
	return fmt.Sprintf("%d to %d", r.Start(), r.End())
}

type FieldRanges struct {
	List   [][2]protoreflect.FieldNumber 
	once   sync.Once
	sorted [][2]protoreflect.FieldNumber 
}

func (p *FieldRanges) Len() int                              { return len(p.List) }
func (p *FieldRanges) Get(i int) [2]protoreflect.FieldNumber { return p.List[i] }
func (p *FieldRanges) Has(n protoreflect.FieldNumber) bool {
	for ls := p.lazyInit().sorted; len(ls) > 0; {
		i := len(ls) / 2
		switch r := fieldRange(ls[i]); {
		case n < r.Start():
			ls = ls[:i] 
		case n > r.End():
			ls = ls[i+1:] 
		default:
			return true
		}
	}
	return false
}
func (p *FieldRanges) Format(s fmt.State, r rune)          { descfmt.FormatList(s, r, p) }
func (p *FieldRanges) ProtoInternal(pragma.DoNotImplement) {}
func (p *FieldRanges) lazyInit() *FieldRanges {
	p.once.Do(func() {
		p.sorted = append(p.sorted, p.List...)
		sort.Slice(p.sorted, func(i, j int) bool {
			return p.sorted[i][0] < p.sorted[j][0]
		})
	})
	return p
}



func (p *FieldRanges) CheckValid(isMessageSet bool) error {
	var rp fieldRange
	for i, r := range p.lazyInit().sorted {
		r := fieldRange(r)
		switch {
		case !isValidFieldNumber(r.Start(), isMessageSet):
			return errors.New("invalid field number: %d", r.Start())
		case !isValidFieldNumber(r.End(), isMessageSet):
			return errors.New("invalid field number: %d", r.End())
		case !(r.Start() <= r.End()):
			return errors.New("invalid range: %v", r)
		case !(rp.End() < r.Start()) && i > 0:
			return errors.New("overlapping ranges: %v with %v", rp, r)
		}
		rp = r
	}
	return nil
}




func isValidFieldNumber(n protoreflect.FieldNumber, isMessageSet bool) bool {
	return protowire.MinValidNumber <= n && (n <= protowire.MaxValidNumber || isMessageSet)
}


func (p *FieldRanges) CheckOverlap(q *FieldRanges) error {
	rps := p.lazyInit().sorted
	rqs := q.lazyInit().sorted
	for pi, qi := 0, 0; pi < len(rps) && qi < len(rqs); {
		rp := fieldRange(rps[pi])
		rq := fieldRange(rqs[qi])
		if !(rp.End() < rq.Start() || rq.End() < rp.Start()) {
			return errors.New("overlapping ranges: %v with %v", rp, rq)
		}
		if rp.Start() < rq.Start() {
			pi++
		} else {
			qi++
		}
	}
	return nil
}

type fieldRange [2]protoreflect.FieldNumber

func (r fieldRange) Start() protoreflect.FieldNumber { return r[0] }     
func (r fieldRange) End() protoreflect.FieldNumber   { return r[1] - 1 } 
func (r fieldRange) String() string {
	if r.Start() == r.End() {
		return fmt.Sprintf("%d", r.Start())
	}
	return fmt.Sprintf("%d to %d", r.Start(), r.End())
}

type FieldNumbers struct {
	List []protoreflect.FieldNumber
	once sync.Once
	has  map[protoreflect.FieldNumber]struct{} 
}

func (p *FieldNumbers) Len() int                           { return len(p.List) }
func (p *FieldNumbers) Get(i int) protoreflect.FieldNumber { return p.List[i] }
func (p *FieldNumbers) Has(n protoreflect.FieldNumber) bool {
	p.once.Do(func() {
		if len(p.List) > 0 {
			p.has = make(map[protoreflect.FieldNumber]struct{}, len(p.List))
			for _, n := range p.List {
				p.has[n] = struct{}{}
			}
		}
	})
	_, ok := p.has[n]
	return ok
}
func (p *FieldNumbers) Format(s fmt.State, r rune)          { descfmt.FormatList(s, r, p) }
func (p *FieldNumbers) ProtoInternal(pragma.DoNotImplement) {}

type OneofFields struct {
	List   []protoreflect.FieldDescriptor
	once   sync.Once
	byName map[protoreflect.Name]protoreflect.FieldDescriptor        
	byJSON map[string]protoreflect.FieldDescriptor                   
	byText map[string]protoreflect.FieldDescriptor                   
	byNum  map[protoreflect.FieldNumber]protoreflect.FieldDescriptor 
}

func (p *OneofFields) Len() int                               { return len(p.List) }
func (p *OneofFields) Get(i int) protoreflect.FieldDescriptor { return p.List[i] }
func (p *OneofFields) ByName(s protoreflect.Name) protoreflect.FieldDescriptor {
	return p.lazyInit().byName[s]
}
func (p *OneofFields) ByJSONName(s string) protoreflect.FieldDescriptor {
	return p.lazyInit().byJSON[s]
}
func (p *OneofFields) ByTextName(s string) protoreflect.FieldDescriptor {
	return p.lazyInit().byText[s]
}
func (p *OneofFields) ByNumber(n protoreflect.FieldNumber) protoreflect.FieldDescriptor {
	return p.lazyInit().byNum[n]
}
func (p *OneofFields) Format(s fmt.State, r rune)          { descfmt.FormatList(s, r, p) }
func (p *OneofFields) ProtoInternal(pragma.DoNotImplement) {}

func (p *OneofFields) lazyInit() *OneofFields {
	p.once.Do(func() {
		if len(p.List) > 0 {
			p.byName = make(map[protoreflect.Name]protoreflect.FieldDescriptor, len(p.List))
			p.byJSON = make(map[string]protoreflect.FieldDescriptor, len(p.List))
			p.byText = make(map[string]protoreflect.FieldDescriptor, len(p.List))
			p.byNum = make(map[protoreflect.FieldNumber]protoreflect.FieldDescriptor, len(p.List))
			for _, f := range p.List {
				
				p.byName[f.Name()] = f
				p.byJSON[f.JSONName()] = f
				p.byText[f.TextName()] = f
				p.byNum[f.Number()] = f
			}
		}
	})
	return p
}

type SourceLocations struct {
	
	
	
	List []protoreflect.SourceLocation

	
	
	
	File protoreflect.FileDescriptor

	once   sync.Once
	byPath map[pathKey]int
}

func (p *SourceLocations) Len() int                              { return len(p.List) }
func (p *SourceLocations) Get(i int) protoreflect.SourceLocation { return p.lazyInit().List[i] }
func (p *SourceLocations) byKey(k pathKey) protoreflect.SourceLocation {
	if i, ok := p.lazyInit().byPath[k]; ok {
		return p.List[i]
	}
	return protoreflect.SourceLocation{}
}
func (p *SourceLocations) ByPath(path protoreflect.SourcePath) protoreflect.SourceLocation {
	return p.byKey(newPathKey(path))
}
func (p *SourceLocations) ByDescriptor(desc protoreflect.Descriptor) protoreflect.SourceLocation {
	if p.File != nil && desc != nil && p.File != desc.ParentFile() {
		return protoreflect.SourceLocation{} 
	}
	var pathArr [16]int32
	path := pathArr[:0]
	for {
		switch desc.(type) {
		case protoreflect.FileDescriptor:
			
			for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
				path[i], path[j] = path[j], path[i]
			}
			return p.byKey(newPathKey(path))
		case protoreflect.MessageDescriptor:
			path = append(path, int32(desc.Index()))
			desc = desc.Parent()
			switch desc.(type) {
			case protoreflect.FileDescriptor:
				path = append(path, int32(genid.FileDescriptorProto_MessageType_field_number))
			case protoreflect.MessageDescriptor:
				path = append(path, int32(genid.DescriptorProto_NestedType_field_number))
			default:
				return protoreflect.SourceLocation{}
			}
		case protoreflect.FieldDescriptor:
			isExtension := desc.(protoreflect.FieldDescriptor).IsExtension()
			path = append(path, int32(desc.Index()))
			desc = desc.Parent()
			if isExtension {
				switch desc.(type) {
				case protoreflect.FileDescriptor:
					path = append(path, int32(genid.FileDescriptorProto_Extension_field_number))
				case protoreflect.MessageDescriptor:
					path = append(path, int32(genid.DescriptorProto_Extension_field_number))
				default:
					return protoreflect.SourceLocation{}
				}
			} else {
				switch desc.(type) {
				case protoreflect.MessageDescriptor:
					path = append(path, int32(genid.DescriptorProto_Field_field_number))
				default:
					return protoreflect.SourceLocation{}
				}
			}
		case protoreflect.OneofDescriptor:
			path = append(path, int32(desc.Index()))
			desc = desc.Parent()
			switch desc.(type) {
			case protoreflect.MessageDescriptor:
				path = append(path, int32(genid.DescriptorProto_OneofDecl_field_number))
			default:
				return protoreflect.SourceLocation{}
			}
		case protoreflect.EnumDescriptor:
			path = append(path, int32(desc.Index()))
			desc = desc.Parent()
			switch desc.(type) {
			case protoreflect.FileDescriptor:
				path = append(path, int32(genid.FileDescriptorProto_EnumType_field_number))
			case protoreflect.MessageDescriptor:
				path = append(path, int32(genid.DescriptorProto_EnumType_field_number))
			default:
				return protoreflect.SourceLocation{}
			}
		case protoreflect.EnumValueDescriptor:
			path = append(path, int32(desc.Index()))
			desc = desc.Parent()
			switch desc.(type) {
			case protoreflect.EnumDescriptor:
				path = append(path, int32(genid.EnumDescriptorProto_Value_field_number))
			default:
				return protoreflect.SourceLocation{}
			}
		case protoreflect.ServiceDescriptor:
			path = append(path, int32(desc.Index()))
			desc = desc.Parent()
			switch desc.(type) {
			case protoreflect.FileDescriptor:
				path = append(path, int32(genid.FileDescriptorProto_Service_field_number))
			default:
				return protoreflect.SourceLocation{}
			}
		case protoreflect.MethodDescriptor:
			path = append(path, int32(desc.Index()))
			desc = desc.Parent()
			switch desc.(type) {
			case protoreflect.ServiceDescriptor:
				path = append(path, int32(genid.ServiceDescriptorProto_Method_field_number))
			default:
				return protoreflect.SourceLocation{}
			}
		default:
			return protoreflect.SourceLocation{}
		}
	}
}
func (p *SourceLocations) lazyInit() *SourceLocations {
	p.once.Do(func() {
		if len(p.List) > 0 {
			
			pathIdxs := make(map[pathKey][]int, len(p.List))
			for i, l := range p.List {
				k := newPathKey(l.Path)
				pathIdxs[k] = append(pathIdxs[k], i)
			}

			
			p.byPath = make(map[pathKey]int, len(p.List))
			for k, idxs := range pathIdxs {
				for i := 0; i < len(idxs)-1; i++ {
					p.List[idxs[i]].Next = idxs[i+1]
				}
				p.List[idxs[len(idxs)-1]].Next = 0
				p.byPath[k] = idxs[0] 
			}
		}
	})
	return p
}
func (p *SourceLocations) ProtoInternal(pragma.DoNotImplement) {}


type pathKey struct {
	arr [16]uint8 
	str string    
}

func newPathKey(p protoreflect.SourcePath) (k pathKey) {
	if len(p) < len(k.arr) {
		for i, ps := range p {
			if ps < 0 || math.MaxUint8 <= ps {
				return pathKey{str: p.String()}
			}
			k.arr[i] = uint8(ps)
		}
		k.arr[len(k.arr)-1] = uint8(len(p))
		return k
	}
	return pathKey{str: p.String()}
}
