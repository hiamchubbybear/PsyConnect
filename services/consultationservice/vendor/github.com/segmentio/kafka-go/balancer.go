package kafka

import (
	"hash"
	"hash/crc32"
	"hash/fnv"
	"math/rand"
	"sort"
	"sync"
)






type Balancer interface {
	
	
	
	
	
	
	
	Balance(msg Message, partitions ...int) (partition int)
}



type BalancerFunc func(Message, ...int) int


func (f BalancerFunc) Balance(msg Message, partitions ...int) int {
	return f(msg, partitions...)
}





type RoundRobin struct {
	ChunkSize int
	
	
	counter uint32

	mutex sync.Mutex
}


func (rr *RoundRobin) Balance(msg Message, partitions ...int) int {
	return rr.balance(partitions)
}

func (rr *RoundRobin) balance(partitions []int) int {
	rr.mutex.Lock()
	defer rr.mutex.Unlock()

	if rr.ChunkSize < 1 {
		rr.ChunkSize = 1
	}

	length := len(partitions)
	counterNow := rr.counter
	offset := int(counterNow / uint32(rr.ChunkSize))
	rr.counter++
	return partitions[offset%length]
}







type LeastBytes struct {
	mutex    sync.Mutex
	counters []leastBytesCounter
}

type leastBytesCounter struct {
	partition int
	bytes     uint64
}


func (lb *LeastBytes) Balance(msg Message, partitions ...int) int {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	
	if len(partitions) != len(lb.counters) {
		lb.counters = lb.makeCounters(partitions...)
	}

	minBytes := lb.counters[0].bytes
	minIndex := 0

	for i, c := range lb.counters[1:] {
		if c.bytes < minBytes {
			minIndex = i + 1
			minBytes = c.bytes
		}
	}

	c := &lb.counters[minIndex]
	c.bytes += uint64(len(msg.Key)) + uint64(len(msg.Value))
	return c.partition
}

func (lb *LeastBytes) makeCounters(partitions ...int) (counters []leastBytesCounter) {
	counters = make([]leastBytesCounter, len(partitions))

	for i, p := range partitions {
		counters[i].partition = p
	}

	sort.Slice(counters, func(i int, j int) bool {
		return counters[i].partition < counters[j].partition
	})
	return
}

var (
	fnv1aPool = &sync.Pool{
		New: func() interface{} {
			return fnv.New32a()
		},
	}
)












type Hash struct {
	rr     RoundRobin
	Hasher hash.Hash32

	
	
	
	lock sync.Mutex
}

func (h *Hash) Balance(msg Message, partitions ...int) int {
	if msg.Key == nil {
		return h.rr.Balance(msg, partitions...)
	}

	hasher := h.Hasher
	if hasher != nil {
		h.lock.Lock()
		defer h.lock.Unlock()
	} else {
		hasher = fnv1aPool.Get().(hash.Hash32)
		defer fnv1aPool.Put(hasher)
	}

	hasher.Reset()
	if _, err := hasher.Write(msg.Key); err != nil {
		panic(err)
	}

	
	
	
	partition := int32(hasher.Sum32()) % int32(len(partitions))
	if partition < 0 {
		partition = -partition
	}

	return int(partition)
}












type ReferenceHash struct {
	rr     randomBalancer
	Hasher hash.Hash32

	
	
	
	lock sync.Mutex
}

func (h *ReferenceHash) Balance(msg Message, partitions ...int) int {
	if msg.Key == nil {
		return h.rr.Balance(msg, partitions...)
	}

	hasher := h.Hasher
	if hasher != nil {
		h.lock.Lock()
		defer h.lock.Unlock()
	} else {
		hasher = fnv1aPool.Get().(hash.Hash32)
		defer fnv1aPool.Put(hasher)
	}

	hasher.Reset()
	if _, err := hasher.Write(msg.Key); err != nil {
		panic(err)
	}

	
	
	
	partition := (int32(hasher.Sum32()) & 0x7fffffff) % int32(len(partitions))
	return int(partition)
}

type randomBalancer struct {
	mock int 
}

func (b randomBalancer) Balance(msg Message, partitions ...int) (partition int) {
	if b.mock != 0 {
		return b.mock
	}
	return partitions[rand.Int()%len(partitions)]
}
















type CRC32Balancer struct {
	Consistent bool
	random     randomBalancer
}

func (b CRC32Balancer) Balance(msg Message, partitions ...int) (partition int) {
	
	
	if len(msg.Key) == 0 && !b.Consistent {
		return b.random.Balance(msg, partitions...)
	}

	idx := crc32.ChecksumIEEE(msg.Key) % uint32(len(partitions))
	return partitions[idx]
}






















type Murmur2Balancer struct {
	Consistent bool
	random     randomBalancer
}

func (b Murmur2Balancer) Balance(msg Message, partitions ...int) (partition int) {
	
	
	if msg.Key == nil && !b.Consistent {
		return b.random.Balance(msg, partitions...)
	}

	idx := (murmur2(msg.Key) & 0x7fffffff) % uint32(len(partitions))
	return partitions[idx]
}



func murmur2(data []byte) uint32 {
	length := len(data)
	const (
		seed uint32 = 0x9747b28c
		
		
		m = 0x5bd1e995
		r = 24
	)

	
	h := seed ^ uint32(length)
	length4 := length / 4

	for i := 0; i < length4; i++ {
		i4 := i * 4
		k := (uint32(data[i4+0]) & 0xff) + ((uint32(data[i4+1]) & 0xff) << 8) + ((uint32(data[i4+2]) & 0xff) << 16) + ((uint32(data[i4+3]) & 0xff) << 24)
		k *= m
		k ^= k >> r
		k *= m
		h *= m
		h ^= k
	}

	
	extra := length % 4
	if extra >= 3 {
		h ^= (uint32(data[(length & ^3)+2]) & 0xff) << 16
	}
	if extra >= 2 {
		h ^= (uint32(data[(length & ^3)+1]) & 0xff) << 8
	}
	if extra >= 1 {
		h ^= uint32(data[length & ^3]) & 0xff
		h *= m
	}

	h ^= h >> 13
	h *= m
	h ^= h >> 15

	return h
}
