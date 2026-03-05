



















package atomic

import "time"

//go:generate bin/gen-atomicwrapper -name=Duration -type=time.Duration -wrapped=Int64 -pack=int64 -unpack=time.Duration -cas -swap -json -imports time -file=duration.go


func (d *Duration) Add(delta time.Duration) time.Duration {
	return time.Duration(d.v.Add(int64(delta)))
}


func (d *Duration) Sub(delta time.Duration) time.Duration {
	return time.Duration(d.v.Sub(int64(delta)))
}


func (d *Duration) String() string {
	return d.Load().String()
}
