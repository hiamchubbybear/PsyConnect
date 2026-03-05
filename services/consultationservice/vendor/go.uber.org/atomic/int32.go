





















package atomic

import (
	"encoding/json"
	"strconv"
	"sync/atomic"
)


type Int32 struct {
	_ nocmp 

	v int32
}


func NewInt32(val int32) *Int32 {
	return &Int32{v: val}
}


func (i *Int32) Load() int32 {
	return atomic.LoadInt32(&i.v)
}


func (i *Int32) Add(delta int32) int32 {
	return atomic.AddInt32(&i.v, delta)
}


func (i *Int32) Sub(delta int32) int32 {
	return atomic.AddInt32(&i.v, -delta)
}


func (i *Int32) Inc() int32 {
	return i.Add(1)
}


func (i *Int32) Dec() int32 {
	return i.Sub(1)
}


func (i *Int32) CAS(old, new int32) (swapped bool) {
	return atomic.CompareAndSwapInt32(&i.v, old, new)
}


func (i *Int32) Store(val int32) {
	atomic.StoreInt32(&i.v, val)
}


func (i *Int32) Swap(val int32) (old int32) {
	return atomic.SwapInt32(&i.v, val)
}


func (i *Int32) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Load())
}


func (i *Int32) UnmarshalJSON(b []byte) error {
	var v int32
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	i.Store(v)
	return nil
}


func (i *Int32) String() string {
	v := i.Load()
	return strconv.FormatInt(int64(v), 10)
}
