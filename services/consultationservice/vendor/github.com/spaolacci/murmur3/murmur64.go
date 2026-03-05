package murmur3

import (
	"hash"
)


var (
	_ hash.Hash   = new(digest64)
	_ hash.Hash64 = new(digest64)
	_ bmixer      = new(digest64)
)


type digest64 digest128


func New64() hash.Hash64 { return New64WithSeed(0) }


func New64WithSeed(seed uint32) hash.Hash64 {
	d := (*digest64)(New128WithSeed(seed).(*digest128))
	return d
}

func (d *digest64) Sum(b []byte) []byte {
	h1 := d.Sum64()
	return append(b,
		byte(h1>>56), byte(h1>>48), byte(h1>>40), byte(h1>>32),
		byte(h1>>24), byte(h1>>16), byte(h1>>8), byte(h1))
}

func (d *digest64) Sum64() uint64 {
	h1, _ := (*digest128)(d).Sum128()
	return h1
}






func Sum64(data []byte) uint64 { return Sum64WithSeed(data, 0) }






func Sum64WithSeed(data []byte, seed uint32) uint64 {
	d := &digest128{h1: uint64(seed), h2: uint64(seed)}
	d.seed = seed
	d.tail = d.bmix(data)
	d.clen = len(data)
	h1, _ := d.Sum128()
	return h1
}
