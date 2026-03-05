





















package atomic

import (
	"encoding/json"
	"strconv"
	"sync/atomic"
)


type Uint64 struct {
	_ nocmp 

	v uint64
}


func NewUint64(val uint64) *Uint64 {
	return &Uint64{v: val}
}


func (i *Uint64) Load() uint64 {
	return atomic.LoadUint64(&i.v)
}


func (i *Uint64) Add(delta uint64) uint64 {
	return atomic.AddUint64(&i.v, delta)
}


func (i *Uint64) Sub(delta uint64) uint64 {
	return atomic.AddUint64(&i.v, ^(delta - 1))
}


func (i *Uint64) Inc() uint64 {
	return i.Add(1)
}


func (i *Uint64) Dec() uint64 {
	return i.Sub(1)
}


func (i *Uint64) CAS(old, new uint64) (swapped bool) {
	return atomic.CompareAndSwapUint64(&i.v, old, new)
}


func (i *Uint64) Store(val uint64) {
	atomic.StoreUint64(&i.v, val)
}


func (i *Uint64) Swap(val uint64) (old uint64) {
	return atomic.SwapUint64(&i.v, val)
}


func (i *Uint64) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Load())
}


func (i *Uint64) UnmarshalJSON(b []byte) error {
	var v uint64
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	i.Store(v)
	return nil
}


func (i *Uint64) String() string {
	v := i.Load()
	return strconv.FormatUint(uint64(v), 10)
}
