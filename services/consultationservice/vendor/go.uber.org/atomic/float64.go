





















package atomic

import (
	"encoding/json"
	"math"
)


type Float64 struct {
	_ nocmp 

	v Uint64
}

var _zeroFloat64 float64


func NewFloat64(val float64) *Float64 {
	x := &Float64{}
	if val != _zeroFloat64 {
		x.Store(val)
	}
	return x
}


func (x *Float64) Load() float64 {
	return math.Float64frombits(x.v.Load())
}


func (x *Float64) Store(val float64) {
	x.v.Store(math.Float64bits(val))
}



func (x *Float64) Swap(val float64) (old float64) {
	return math.Float64frombits(x.v.Swap(math.Float64bits(val)))
}


func (x *Float64) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.Load())
}


func (x *Float64) UnmarshalJSON(b []byte) error {
	var v float64
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	x.Store(v)
	return nil
}
