package viper

import "github.com/spf13/pflag"



type FlagValueSet interface {
	VisitAll(fn func(FlagValue))
}



type FlagValue interface {
	HasChanged() bool
	Name() string
	ValueString() string
	ValueType() string
}



type pflagValueSet struct {
	flags *pflag.FlagSet
}


func (p pflagValueSet) VisitAll(fn func(flag FlagValue)) {
	p.flags.VisitAll(func(flag *pflag.Flag) {
		fn(pflagValue{flag})
	})
}



type pflagValue struct {
	flag *pflag.Flag
}


func (p pflagValue) HasChanged() bool {
	return p.flag.Changed
}


func (p pflagValue) Name() string {
	return p.flag.Name
}


func (p pflagValue) ValueString() string {
	return p.flag.Value.String()
}


func (p pflagValue) ValueType() string {
	return p.flag.Value.Type()
}
