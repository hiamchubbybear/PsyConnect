package pflag

import (
	"fmt"
	"io"
	"net"
	"strings"
)


type ipNetSliceValue struct {
	value   *[]net.IPNet
	changed bool
}

func newIPNetSliceValue(val []net.IPNet, p *[]net.IPNet) *ipNetSliceValue {
	ipnsv := new(ipNetSliceValue)
	ipnsv.value = p
	*ipnsv.value = val
	return ipnsv
}



func (s *ipNetSliceValue) Set(val string) error {

	
	rmQuote := strings.NewReplacer(`"`, "", `'`, "", "`", "")

	
	ipNetStrSlice, err := readAsCSV(rmQuote.Replace(val))
	if err != nil && err != io.EOF {
		return err
	}

	
	out := make([]net.IPNet, 0, len(ipNetStrSlice))
	for _, ipNetStr := range ipNetStrSlice {
		_, n, err := net.ParseCIDR(strings.TrimSpace(ipNetStr))
		if err != nil {
			return fmt.Errorf("invalid string being converted to CIDR: %s", ipNetStr)
		}
		out = append(out, *n)
	}

	if !s.changed {
		*s.value = out
	} else {
		*s.value = append(*s.value, out...)
	}

	s.changed = true

	return nil
}


func (s *ipNetSliceValue) Type() string {
	return "ipNetSlice"
}


func (s *ipNetSliceValue) String() string {

	ipNetStrSlice := make([]string, len(*s.value))
	for i, n := range *s.value {
		ipNetStrSlice[i] = n.String()
	}

	out, _ := writeAsCSV(ipNetStrSlice)
	return "[" + out + "]"
}

func ipNetSliceConv(val string) (interface{}, error) {
	val = strings.Trim(val, "[]")
	
	if len(val) == 0 {
		return []net.IPNet{}, nil
	}
	ss := strings.Split(val, ",")
	out := make([]net.IPNet, len(ss))
	for i, sval := range ss {
		_, n, err := net.ParseCIDR(strings.TrimSpace(sval))
		if err != nil {
			return nil, fmt.Errorf("invalid string being converted to CIDR: %s", sval)
		}
		out[i] = *n
	}
	return out, nil
}


func (f *FlagSet) GetIPNetSlice(name string) ([]net.IPNet, error) {
	val, err := f.getFlagType(name, "ipNetSlice", ipNetSliceConv)
	if err != nil {
		return []net.IPNet{}, err
	}
	return val.([]net.IPNet), nil
}



func (f *FlagSet) IPNetSliceVar(p *[]net.IPNet, name string, value []net.IPNet, usage string) {
	f.VarP(newIPNetSliceValue(value, p), name, "", usage)
}


func (f *FlagSet) IPNetSliceVarP(p *[]net.IPNet, name, shorthand string, value []net.IPNet, usage string) {
	f.VarP(newIPNetSliceValue(value, p), name, shorthand, usage)
}



func IPNetSliceVar(p *[]net.IPNet, name string, value []net.IPNet, usage string) {
	CommandLine.VarP(newIPNetSliceValue(value, p), name, "", usage)
}


func IPNetSliceVarP(p *[]net.IPNet, name, shorthand string, value []net.IPNet, usage string) {
	CommandLine.VarP(newIPNetSliceValue(value, p), name, shorthand, usage)
}



func (f *FlagSet) IPNetSlice(name string, value []net.IPNet, usage string) *[]net.IPNet {
	p := []net.IPNet{}
	f.IPNetSliceVarP(&p, name, "", value, usage)
	return &p
}


func (f *FlagSet) IPNetSliceP(name, shorthand string, value []net.IPNet, usage string) *[]net.IPNet {
	p := []net.IPNet{}
	f.IPNetSliceVarP(&p, name, shorthand, value, usage)
	return &p
}



func IPNetSlice(name string, value []net.IPNet, usage string) *[]net.IPNet {
	return CommandLine.IPNetSliceP(name, "", value, usage)
}


func IPNetSliceP(name, shorthand string, value []net.IPNet, usage string) *[]net.IPNet {
	return CommandLine.IPNetSliceP(name, shorthand, value, usage)
}
