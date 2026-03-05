





package idna



type valueRange struct {
	value  uint16 
	lo, hi byte   
}

type sparseBlocks struct {
	values []valueRange
	offset []uint16
}

var idnaSparse = sparseBlocks{
	values: idnaSparseValues[:],
	offset: idnaSparseOffset[:],
}


var trie = &idnaTrie{}





func (t *sparseBlocks) lookup(n uint32, b byte) uint16 {
	offset := t.offset[n]
	header := t.values[offset]
	lo := offset + 1
	hi := lo + uint16(header.lo)
	for lo < hi {
		m := lo + (hi-lo)/2
		r := t.values[m]
		if r.lo <= b && b <= r.hi {
			return r.value + uint16(b-r.lo)*header.value
		}
		if b < r.lo {
			hi = m
		} else {
			lo = m + 1
		}
	}
	return 0
}
