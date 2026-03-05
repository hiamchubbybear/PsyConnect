





















package atomic

import (
	"encoding/json"
	"strconv"
	"sync/atomic"
)


type Uint32 struct {
	_ nocmp 

	v uint32
}


func NewUint32(val uint32) *Uint32 {
	return &Uint32{v: val}
}


func (i *Uint32) Load() uint32 {
	return atomic.LoadUint32(&i.v)
}


func (i *Uint32) Add(delta uint32) uint32 {
	return atomic.AddUint32(&i.v, delta)
}


func (i *Uint32) Sub(delta uint32) uint32 {
	return atomic.AddUint32(&i.v, ^(delta - 1))
}


func (i *Uint32) Inc() uint32 {
	return i.Add(1)
}


func (i *Uint32) Dec() uint32 {
	return i.Sub(1)
}


func (i *Uint32) CAS(old, new uint32) (swapped bool) {
	return atomic.CompareAndSwapUint32(&i.v, old, new)
}


func (i *Uint32) Store(val uint32) {
	atomic.StoreUint32(&i.v, val)
}


func (i *Uint32) Swap(val uint32) (old uint32) {
	return atomic.SwapUint32(&i.v, val)
}


func (i *Uint32) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Load())
}


func (i *Uint32) UnmarshalJSON(b []byte) error {
	var v uint32
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	i.Store(v)
	return nil
}


func (i *Uint32) String() string {
	v := i.Load()
	return strconv.FormatUint(uint64(v), 10)
}
