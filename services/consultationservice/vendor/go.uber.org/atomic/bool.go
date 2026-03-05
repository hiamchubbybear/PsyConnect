





















package atomic

import (
	"encoding/json"
)


type Bool struct {
	_ nocmp 

	v Uint32
}

var _zeroBool bool


func NewBool(val bool) *Bool {
	x := &Bool{}
	if val != _zeroBool {
		x.Store(val)
	}
	return x
}


func (x *Bool) Load() bool {
	return truthy(x.v.Load())
}


func (x *Bool) Store(val bool) {
	x.v.Store(boolToInt(val))
}


func (x *Bool) CAS(old, new bool) (swapped bool) {
	return x.v.CAS(boolToInt(old), boolToInt(new))
}



func (x *Bool) Swap(val bool) (old bool) {
	return truthy(x.v.Swap(boolToInt(val)))
}


func (x *Bool) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.Load())
}


func (x *Bool) UnmarshalJSON(b []byte) error {
	var v bool
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	x.Store(v)
	return nil
}
