package pflag

import (
	"fmt"
	"strconv"
	"strings"
)


type int64SliceValue struct {
	value   *[]int64
	changed bool
}

func newInt64SliceValue(val []int64, p *[]int64) *int64SliceValue {
	isv := new(int64SliceValue)
	isv.value = p
	*isv.value = val
	return isv
}

func (s *int64SliceValue) Set(val string) error {
	ss := strings.Split(val, ",")
	out := make([]int64, len(ss))
	for i, d := range ss {
		var err error
		out[i], err = strconv.ParseInt(d, 0, 64)
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

func (s *int64SliceValue) Type() string {
	return "int64Slice"
}

func (s *int64SliceValue) String() string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = fmt.Sprintf("%d", d)
	}
	return "[" + strings.Join(out, ",") + "]"
}

func (s *int64SliceValue) fromString(val string) (int64, error) {
	return strconv.ParseInt(val, 0, 64)
}

func (s *int64SliceValue) toString(val int64) string {
	return fmt.Sprintf("%d", val)
}

func (s *int64SliceValue) Append(val string) error {
	i, err := s.fromString(val)
	if err != nil {
		return err
	}
	*s.value = append(*s.value, i)
	return nil
}

func (s *int64SliceValue) Replace(val []string) error {
	out := make([]int64, len(val))
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

func (s *int64SliceValue) GetSlice() []string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = s.toString(d)
	}
	return out
}

func int64SliceConv(val string) (interface{}, error) {
	val = strings.Trim(val, "[]")
	
	if len(val) == 0 {
		return []int64{}, nil
	}
	ss := strings.Split(val, ",")
	out := make([]int64, len(ss))
	for i, d := range ss {
		var err error
		out[i], err = strconv.ParseInt(d, 0, 64)
		if err != nil {
			return nil, err
		}

	}
	return out, nil
}


func (f *FlagSet) GetInt64Slice(name string) ([]int64, error) {
	val, err := f.getFlagType(name, "int64Slice", int64SliceConv)
	if err != nil {
		return []int64{}, err
	}
	return val.([]int64), nil
}



func (f *FlagSet) Int64SliceVar(p *[]int64, name string, value []int64, usage string) {
	f.VarP(newInt64SliceValue(value, p), name, "", usage)
}


func (f *FlagSet) Int64SliceVarP(p *[]int64, name, shorthand string, value []int64, usage string) {
	f.VarP(newInt64SliceValue(value, p), name, shorthand, usage)
}



func Int64SliceVar(p *[]int64, name string, value []int64, usage string) {
	CommandLine.VarP(newInt64SliceValue(value, p), name, "", usage)
}


func Int64SliceVarP(p *[]int64, name, shorthand string, value []int64, usage string) {
	CommandLine.VarP(newInt64SliceValue(value, p), name, shorthand, usage)
}



func (f *FlagSet) Int64Slice(name string, value []int64, usage string) *[]int64 {
	p := []int64{}
	f.Int64SliceVarP(&p, name, "", value, usage)
	return &p
}


func (f *FlagSet) Int64SliceP(name, shorthand string, value []int64, usage string) *[]int64 {
	p := []int64{}
	f.Int64SliceVarP(&p, name, shorthand, value, usage)
	return &p
}



func Int64Slice(name string, value []int64, usage string) *[]int64 {
	return CommandLine.Int64SliceP(name, "", value, usage)
}


func Int64SliceP(name, shorthand string, value []int64, usage string) *[]int64 {
	return CommandLine.Int64SliceP(name, shorthand, value, usage)
}
