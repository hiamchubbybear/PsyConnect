





















package atomic

import (
	"encoding/json"
	"strconv"
	"sync/atomic"
)


type Uintptr struct {
	_ nocmp 

	v uintptr
}


func NewUintptr(val uintptr) *Uintptr {
	return &Uintptr{v: val}
}


func (i *Uintptr) Load() uintptr {
	return atomic.LoadUintptr(&i.v)
}


func (i *Uintptr) Add(delta uintptr) uintptr {
	return atomic.AddUintptr(&i.v, delta)
}


func (i *Uintptr) Sub(delta uintptr) uintptr {
	return atomic.AddUintptr(&i.v, ^(delta - 1))
}


func (i *Uintptr) Inc() uintptr {
	return i.Add(1)
}


func (i *Uintptr) Dec() uintptr {
	return i.Sub(1)
}


func (i *Uintptr) CAS(old, new uintptr) (swapped bool) {
	return atomic.CompareAndSwapUintptr(&i.v, old, new)
}


func (i *Uintptr) Store(val uintptr) {
	atomic.StoreUintptr(&i.v, val)
}


func (i *Uintptr) Swap(val uintptr) (old uintptr) {
	return atomic.SwapUintptr(&i.v, val)
}


func (i *Uintptr) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Load())
}


func (i *Uintptr) UnmarshalJSON(b []byte) error {
	var v uintptr
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	i.Store(v)
	return nil
}


func (i *Uintptr) String() string {
	v := i.Load()
	return strconv.FormatUint(uint64(v), 10)
}
