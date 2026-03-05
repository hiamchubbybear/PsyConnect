





//go:build go1.16

package idna



func (c info) appendMapping(b []byte, s string) []byte {
	index := int(c >> indexShift)
	if c&xorBit == 0 {
		p := index
		return append(b, mappings[mappingIndex[p]:mappingIndex[p+1]]...)
	}
	b = append(b, s...)
	if c&inlineXOR == inlineXOR {
		
		b[len(b)-1] ^= byte(index)
	} else {
		for p := len(b) - int(xorData[index]); p < len(b); p++ {
			index++
			b[p] ^= xorData[index]
		}
	}
	return b
}
