



















package atomic

import "time"

//go:generate bin/gen-atomicwrapper -name=Time -type=time.Time -wrapped=Value -pack=packTime -unpack=unpackTime -imports time -file=time.go

func packTime(t time.Time) interface{} {
	return t
}

func unpackTime(v interface{}) time.Time {
	if t, ok := v.(time.Time); ok {
		return t
	}
	return time.Time{}
}
