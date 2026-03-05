





















package atomic

import (
	"encoding/json"
	"strconv"
	"sync/atomic"
)


type Int64 struct {
	_ nocmp 

	v int64
}


func NewInt64(val int64) *Int64 {
	return &Int64{v: val}
}


func (i *Int64) Load() int64 {
	return atomic.LoadInt64(&i.v)
}


func (i *Int64) Add(delta int64) int64 {
	return atomic.AddInt64(&i.v, delta)
}


func (i *Int64) Sub(delta int64) int64 {
	return atomic.AddInt64(&i.v, -delta)
}


func (i *Int64) Inc() int64 {
	return i.Add(1)
}


func (i *Int64) Dec() int64 {
	return i.Sub(1)
}


func (i *Int64) CAS(old, new int64) (swapped bool) {
	return atomic.CompareAndSwapInt64(&i.v, old, new)
}


func (i *Int64) Store(val int64) {
	atomic.StoreInt64(&i.v, val)
}


func (i *Int64) Swap(val int64) (old int64) {
	return atomic.SwapInt64(&i.v, val)
}


func (i *Int64) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Load())
}


func (i *Int64) UnmarshalJSON(b []byte) error {
	var v int64
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	i.Store(v)
	return nil
}


func (i *Int64) String() string {
	v := i.Load()
	return strconv.FormatInt(int64(v), 10)
}
