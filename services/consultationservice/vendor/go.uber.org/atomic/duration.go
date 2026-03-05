





















package atomic

import (
	"encoding/json"
	"time"
)


type Duration struct {
	_ nocmp 

	v Int64
}

var _zeroDuration time.Duration


func NewDuration(val time.Duration) *Duration {
	x := &Duration{}
	if val != _zeroDuration {
		x.Store(val)
	}
	return x
}


func (x *Duration) Load() time.Duration {
	return time.Duration(x.v.Load())
}


func (x *Duration) Store(val time.Duration) {
	x.v.Store(int64(val))
}


func (x *Duration) CAS(old, new time.Duration) (swapped bool) {
	return x.v.CAS(int64(old), int64(new))
}



func (x *Duration) Swap(val time.Duration) (old time.Duration) {
	return time.Duration(x.v.Swap(int64(val)))
}


func (x *Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.Load())
}


func (x *Duration) UnmarshalJSON(b []byte) error {
	var v time.Duration
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	x.Store(v)
	return nil
}
