




package huff0





type byteReader struct {
	b   []byte
	off int
}


func (b *byteReader) init(in []byte) {
	b.b = in
	b.off = 0
}


func (b byteReader) Int32() int32 {
	v3 := int32(b.b[b.off+3])
	v2 := int32(b.b[b.off+2])
	v1 := int32(b.b[b.off+1])
	v0 := int32(b.b[b.off])
	return (v3 << 24) | (v2 << 16) | (v1 << 8) | v0
}


func (b byteReader) Uint32() uint32 {
	v3 := uint32(b.b[b.off+3])
	v2 := uint32(b.b[b.off+2])
	v1 := uint32(b.b[b.off+1])
	v0 := uint32(b.b[b.off])
	return (v3 << 24) | (v2 << 16) | (v1 << 8) | v0
}


func (b byteReader) remain() int {
	return len(b.b) - b.off
}
