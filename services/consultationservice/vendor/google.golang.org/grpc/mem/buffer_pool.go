

package mem

import (
	"sort"
	"sync"

	"google.golang.org/grpc/internal"
)



type BufferPool interface {
	
	Get(length int) *[]byte

	
	Put(*[]byte)
}

var defaultBufferPoolSizes = []int{
	256,
	4 << 10,  
	16 << 10, 
	32 << 10, 
	1 << 20,  
}

var defaultBufferPool BufferPool

func init() {
	defaultBufferPool = NewTieredBufferPool(defaultBufferPoolSizes...)

	internal.SetDefaultBufferPoolForTesting = func(pool BufferPool) {
		defaultBufferPool = pool
	}

	internal.SetBufferPoolingThresholdForTesting = func(threshold int) {
		bufferPoolingThreshold = threshold
	}
}




func DefaultBufferPool() BufferPool {
	return defaultBufferPool
}



func NewTieredBufferPool(poolSizes ...int) BufferPool {
	sort.Ints(poolSizes)
	pools := make([]*sizedBufferPool, len(poolSizes))
	for i, s := range poolSizes {
		pools[i] = newSizedBufferPool(s)
	}
	return &tieredBufferPool{
		sizedPools: pools,
	}
}



type tieredBufferPool struct {
	sizedPools   []*sizedBufferPool
	fallbackPool simpleBufferPool
}

func (p *tieredBufferPool) Get(size int) *[]byte {
	return p.getPool(size).Get(size)
}

func (p *tieredBufferPool) Put(buf *[]byte) {
	p.getPool(cap(*buf)).Put(buf)
}

func (p *tieredBufferPool) getPool(size int) BufferPool {
	poolIdx := sort.Search(len(p.sizedPools), func(i int) bool {
		return p.sizedPools[i].defaultSize >= size
	})

	if poolIdx == len(p.sizedPools) {
		return &p.fallbackPool
	}

	return p.sizedPools[poolIdx]
}









type sizedBufferPool struct {
	pool        sync.Pool
	defaultSize int
}

func (p *sizedBufferPool) Get(size int) *[]byte {
	buf := p.pool.Get().(*[]byte)
	b := *buf
	clear(b[:cap(b)])
	*buf = b[:size]
	return buf
}

func (p *sizedBufferPool) Put(buf *[]byte) {
	if cap(*buf) < p.defaultSize {
		
		
		
		return
	}
	p.pool.Put(buf)
}

func newSizedBufferPool(size int) *sizedBufferPool {
	return &sizedBufferPool{
		pool: sync.Pool{
			New: func() any {
				buf := make([]byte, size)
				return &buf
			},
		},
		defaultSize: size,
	}
}

var _ BufferPool = (*simpleBufferPool)(nil)





type simpleBufferPool struct {
	pool sync.Pool
}

func (p *simpleBufferPool) Get(size int) *[]byte {
	bs, ok := p.pool.Get().(*[]byte)
	if ok && cap(*bs) >= size {
		*bs = (*bs)[:size]
		return bs
	}

	
	
	if ok {
		p.pool.Put(bs)
	}

	b := make([]byte, size)
	return &b
}

func (p *simpleBufferPool) Put(buf *[]byte) {
	p.pool.Put(buf)
}

var _ BufferPool = NopBufferPool{}


type NopBufferPool struct{}


func (NopBufferPool) Get(length int) *[]byte {
	b := make([]byte, length)
	return &b
}


func (NopBufferPool) Put(*[]byte) {
}
