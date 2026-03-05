


















package rand

import "sync"



type Source interface {
	Uint64() uint64
	Seed(seed uint64)
}


func NewSource(seed uint64) Source {
	var rng PCGSource
	rng.Seed(seed)
	return &rng
}


type Rand struct {
	src Source

	
	
	
	
	readVal uint64
	
	
	readPos int8
}



func New(src Source) *Rand {
	return &Rand{src: src}
}



func (r *Rand) Seed(seed uint64) {
	if lk, ok := r.src.(*LockedSource); ok {
		lk.seedPos(seed, &r.readPos)
		return
	}

	r.src.Seed(seed)
	r.readPos = 0
}


func (r *Rand) Uint64() uint64 { return r.src.Uint64() }


func (r *Rand) Int63() int64 { return int64(r.src.Uint64() &^ (1 << 63)) }


func (r *Rand) Uint32() uint32 { return uint32(r.Uint64() >> 32) }


func (r *Rand) Int31() int32 { return int32(r.Uint64() >> 33) }


func (r *Rand) Int() int {
	u := uint(r.Uint64())
	return int(u << 1 >> 1) 
}

const maxUint64 = (1 << 64) - 1




func (r *Rand) Uint64n(n uint64) uint64 {
	if n&(n-1) == 0 { 
		if n == 0 {
			panic("invalid argument to Uint64n")
		}
		return r.Uint64() & (n - 1)
	}
	
	
	v := r.Uint64()
	if v > maxUint64-n { 
		ceiling := maxUint64 - maxUint64%n
		for v >= ceiling {
			v = r.Uint64()
		}
	}

	return v % n
}



func (r *Rand) Int63n(n int64) int64 {
	if n <= 0 {
		panic("invalid argument to Int63n")
	}
	return int64(r.Uint64n(uint64(n)))
}



func (r *Rand) Int31n(n int32) int32 {
	if n <= 0 {
		panic("invalid argument to Int31n")
	}
	
	return int32(r.Uint64n(uint64(n)))
}



func (r *Rand) Intn(n int) int {
	if n <= 0 {
		panic("invalid argument to Intn")
	}
	
	return int(r.Uint64n(uint64(n)))
}


func (r *Rand) Float64() float64 {
	
	
	
	
	
	
	
	
	
	
again:
	f := float64(r.Uint64n(1<<53)) / (1 << 53)
	if f == 1.0 {
		goto again 
	}
	return f
}


func (r *Rand) Float32() float32 {
	
	
again:
	f := float32(r.Float64())
	if f == 1 {
		goto again 
	}
	return f
}


func (r *Rand) Perm(n int) []int {
	m := make([]int, n)
	
	
	
	
	
	for i := 0; i < n; i++ {
		j := r.Intn(i + 1)
		m[i] = m[j]
		m[j] = i
	}
	return m
}




func (r *Rand) Shuffle(n int, swap func(i, j int)) {
	if n < 0 {
		panic("invalid argument to Shuffle")
	}

	
	
	
	
	
	
	i := n - 1
	for ; i > 1<<31-1-1; i-- {
		j := int(r.Int63n(int64(i + 1)))
		swap(i, j)
	}
	for ; i > 0; i-- {
		j := int(r.Int31n(int32(i + 1)))
		swap(i, j)
	}
}





func (r *Rand) Read(p []byte) (n int, err error) {
	if lk, ok := r.src.(*LockedSource); ok {
		return lk.Read(p, &r.readVal, &r.readPos)
	}
	return read(p, r.src, &r.readVal, &r.readPos)
}

func read(p []byte, src Source, readVal *uint64, readPos *int8) (n int, err error) {
	pos := *readPos
	val := *readVal
	rng, _ := src.(*PCGSource)
	for n = 0; n < len(p); n++ {
		if pos == 0 {
			if rng != nil {
				val = rng.Uint64()
			} else {
				val = src.Uint64()
			}
			pos = 8
		}
		p[n] = byte(val)
		val >>= 8
		pos--
	}
	*readPos = pos
	*readVal = val
	return
}



var globalRand = New(&LockedSource{src: *NewSource(1).(*PCGSource)})


var _ PCGSource = globalRand.src.(*LockedSource).src





func Seed(seed uint64) { globalRand.Seed(seed) }



func Int63() int64 { return globalRand.Int63() }



func Uint32() uint32 { return globalRand.Uint32() }



func Uint64() uint64 { return globalRand.Uint64() }



func Int31() int32 { return globalRand.Int31() }


func Int() int { return globalRand.Int() }




func Int63n(n int64) int64 { return globalRand.Int63n(n) }




func Int31n(n int32) int32 { return globalRand.Int31n(n) }




func Intn(n int) int { return globalRand.Intn(n) }



func Float64() float64 { return globalRand.Float64() }



func Float32() float32 { return globalRand.Float32() }



func Perm(n int) []int { return globalRand.Perm(n) }




func Shuffle(n int, swap func(i, j int)) { globalRand.Shuffle(n, swap) }




func Read(p []byte) (n int, err error) { return globalRand.Read(p) }









func NormFloat64() float64 { return globalRand.NormFloat64() }








func ExpFloat64() float64 { return globalRand.ExpFloat64() }





type LockedSource struct {
	lk  sync.Mutex
	src PCGSource
}

func (s *LockedSource) Uint64() (n uint64) {
	s.lk.Lock()
	n = s.src.Uint64()
	s.lk.Unlock()
	return
}

func (s *LockedSource) Seed(seed uint64) {
	s.lk.Lock()
	s.src.Seed(seed)
	s.lk.Unlock()
}


func (s *LockedSource) seedPos(seed uint64, readPos *int8) {
	s.lk.Lock()
	s.src.Seed(seed)
	*readPos = 0
	s.lk.Unlock()
}


func (s *LockedSource) Read(p []byte, readVal *uint64, readPos *int8) (n int, err error) {
	s.lk.Lock()
	n, err = read(p, &s.src, readVal, readPos)
	s.lk.Unlock()
	return
}
