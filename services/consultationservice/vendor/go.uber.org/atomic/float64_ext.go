



















package atomic

import (
	"math"
	"strconv"
)

//go:generate bin/gen-atomicwrapper -name=Float64 -type=float64 -wrapped=Uint64 -pack=math.Float64bits -unpack=math.Float64frombits -swap -json -imports math -file=float64.go


func (f *Float64) Add(delta float64) float64 {
	for {
		old := f.Load()
		new := old + delta
		if f.CAS(old, new) {
			return new
		}
	}
}


func (f *Float64) Sub(delta float64) float64 {
	return f.Add(-delta)
}
















func (f *Float64) CAS(old, new float64) (swapped bool) {
	return f.v.CAS(math.Float64bits(old), math.Float64bits(new))
}


func (f *Float64) String() string {
	
	return strconv.FormatFloat(f.Load(), 'g', -1, 64)
}
