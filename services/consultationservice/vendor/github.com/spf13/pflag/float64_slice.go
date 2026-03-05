package pflag

import (
	"fmt"
	"strconv"
	"strings"
)


type float64SliceValue struct {
	value   *[]float64
	changed bool
}

func newFloat64SliceValue(val []float64, p *[]float64) *float64SliceValue {
	isv := new(float64SliceValue)
	isv.value = p
	*isv.value = val
	return isv
}

func (s *float64SliceValue) Set(val string) error {
	ss := strings.Split(val, ",")
	out := make([]float64, len(ss))
	for i, d := range ss {
		var err error
		out[i], err = strconv.ParseFloat(d, 64)
		if err != nil {
			return err
		}

	}
	if !s.changed {
		*s.value = out
	} else {
		*s.value = append(*s.value, out...)
	}
	s.changed = true
	return nil
}

func (s *float64SliceValue) Type() string {
	return "float64Slice"
}

func (s *float64SliceValue) String() string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = fmt.Sprintf("%f", d)
	}
	return "[" + strings.Join(out, ",") + "]"
}

func (s *float64SliceValue) fromString(val string) (float64, error) {
	return strconv.ParseFloat(val, 64)
}

func (s *float64SliceValue) toString(val float64) string {
	return fmt.Sprintf("%f", val)
}

func (s *float64SliceValue) Append(val string) error {
	i, err := s.fromString(val)
	if err != nil {
		return err
	}
	*s.value = append(*s.value, i)
	return nil
}

func (s *float64SliceValue) Replace(val []string) error {
	out := make([]float64, len(val))
	for i, d := range val {
		var err error
		out[i], err = s.fromString(d)
		if err != nil {
			return err
		}
	}
	*s.value = out
	return nil
}

func (s *float64SliceValue) GetSlice() []string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = s.toString(d)
	}
	return out
}

func float64SliceConv(val string) (interface{}, error) {
	val = strings.Trim(val, "[]")
	
	if len(val) == 0 {
		return []float64{}, nil
	}
	ss := strings.Split(val, ",")
	out := make([]float64, len(ss))
	for i, d := range ss {
		var err error
		out[i], err = strconv.ParseFloat(d, 64)
		if err != nil {
			return nil, err
		}

	}
	return out, nil
}


func (f *FlagSet) GetFloat64Slice(name string) ([]float64, error) {
	val, err := f.getFlagType(name, "float64Slice", float64SliceConv)
	if err != nil {
		return []float64{}, err
	}
	return val.([]float64), nil
}



func (f *FlagSet) Float64SliceVar(p *[]float64, name string, value []float64, usage string) {
	f.VarP(newFloat64SliceValue(value, p), name, "", usage)
}


func (f *FlagSet) Float64SliceVarP(p *[]float64, name, shorthand string, value []float64, usage string) {
	f.VarP(newFloat64SliceValue(value, p), name, shorthand, usage)
}



func Float64SliceVar(p *[]float64, name string, value []float64, usage string) {
	CommandLine.VarP(newFloat64SliceValue(value, p), name, "", usage)
}


func Float64SliceVarP(p *[]float64, name, shorthand string, value []float64, usage string) {
	CommandLine.VarP(newFloat64SliceValue(value, p), name, shorthand, usage)
}



func (f *FlagSet) Float64Slice(name string, value []float64, usage string) *[]float64 {
	p := []float64{}
	f.Float64SliceVarP(&p, name, "", value, usage)
	return &p
}


func (f *FlagSet) Float64SliceP(name, shorthand string, value []float64, usage string) *[]float64 {
	p := []float64{}
	f.Float64SliceVarP(&p, name, shorthand, value, usage)
	return &p
}



func Float64Slice(name string, value []float64, usage string) *[]float64 {
	return CommandLine.Float64SliceP(name, "", value, usage)
}


func Float64SliceP(name, shorthand string, value []float64, usage string) *[]float64 {
	return CommandLine.Float64SliceP(name, shorthand, value, usage)
}
