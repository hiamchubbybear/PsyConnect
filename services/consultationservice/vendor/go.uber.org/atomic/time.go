





















package atomic

import (
	"time"
)


type Time struct {
	_ nocmp 

	v Value
}

var _zeroTime time.Time


func NewTime(val time.Time) *Time {
	x := &Time{}
	if val != _zeroTime {
		x.Store(val)
	}
	return x
}


func (x *Time) Load() time.Time {
	return unpackTime(x.v.Load())
}


func (x *Time) Store(val time.Time) {
	x.v.Store(packTime(val))
}
