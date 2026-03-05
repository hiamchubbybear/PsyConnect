package kafka

import (
	"sync/atomic"
	"time"
)


type SummaryStats struct {
	Avg   int64 `metric:"avg" type:"gauge"`
	Min   int64 `metric:"min" type:"gauge"`
	Max   int64 `metric:"max" type:"gauge"`
	Count int64 `metric:"count" type:"counter"`
	Sum   int64 `metric:"sum" type:"counter"`
}


type DurationStats struct {
	Avg   time.Duration `metric:"avg" type:"gauge"`
	Min   time.Duration `metric:"min" type:"gauge"`
	Max   time.Duration `metric:"max" type:"gauge"`
	Count int64         `metric:"count" type:"counter"`
	Sum   time.Duration `metric:"sum" type:"counter"`
}





type counter int64

func (c *counter) ptr() *int64 {
	return (*int64)(c)
}

func (c *counter) observe(v int64) {
	atomic.AddInt64(c.ptr(), v)
}

func (c *counter) snapshot() int64 {
	return atomic.SwapInt64(c.ptr(), 0)
}






type gauge int64

func (g *gauge) ptr() *int64 {
	return (*int64)(g)
}

func (g *gauge) observe(v int64) {
	atomic.StoreInt64(g.ptr(), v)
}

func (g *gauge) snapshot() int64 {
	return atomic.LoadInt64(g.ptr())
}






type minimum int64

func (m *minimum) ptr() *int64 {
	return (*int64)(m)
}

func (m *minimum) observe(v int64) {
	for {
		ptr := m.ptr()
		min := atomic.LoadInt64(ptr)

		if min >= 0 && min <= v {
			break
		}

		if atomic.CompareAndSwapInt64(ptr, min, v) {
			break
		}
	}
}

func (m *minimum) snapshot() int64 {
	p := m.ptr()
	v := atomic.LoadInt64(p)
	atomic.CompareAndSwapInt64(p, v, -1)
	if v < 0 {
		v = 0
	}
	return v
}






type maximum int64

func (m *maximum) ptr() *int64 {
	return (*int64)(m)
}

func (m *maximum) observe(v int64) {
	for {
		ptr := m.ptr()
		max := atomic.LoadInt64(ptr)

		if max >= 0 && max >= v {
			break
		}

		if atomic.CompareAndSwapInt64(ptr, max, v) {
			break
		}
	}
}

func (m *maximum) snapshot() int64 {
	p := m.ptr()
	v := atomic.LoadInt64(p)
	atomic.CompareAndSwapInt64(p, v, -1)
	if v < 0 {
		v = 0
	}
	return v
}

type summary struct {
	min   minimum
	max   maximum
	sum   counter
	count counter
}

func makeSummary() summary {
	return summary{
		min: -1,
		max: -1,
	}
}

func (s *summary) observe(v int64) {
	s.min.observe(v)
	s.max.observe(v)
	s.sum.observe(v)
	s.count.observe(1)
}

func (s *summary) observeDuration(v time.Duration) {
	s.observe(int64(v))
}

func (s *summary) snapshot() SummaryStats {
	avg := int64(0)
	min := s.min.snapshot()
	max := s.max.snapshot()
	sum := s.sum.snapshot()
	count := s.count.snapshot()

	if count != 0 {
		avg = int64(float64(sum) / float64(count))
	}

	return SummaryStats{
		Avg:   avg,
		Min:   min,
		Max:   max,
		Count: count,
		Sum:   sum,
	}
}

func (s *summary) snapshotDuration() DurationStats {
	summary := s.snapshot()
	return DurationStats{
		Avg:   time.Duration(summary.Avg),
		Min:   time.Duration(summary.Min),
		Max:   time.Duration(summary.Max),
		Count: summary.Count,
		Sum:   time.Duration(summary.Sum),
	}
}
