




package fse





type byteReader struct {
	b   []byte
	off int
}


func (b *byteReader) init(in []byte) {
	b.b = in
	b.off = 0
}


func (b *byteReader) advance(n uint) {
	b.off += int(n)
}


func (b byteReader) Uint32() uint32 {
	b2 := b.b[b.off:]
	b2 = b2[:4]
	v3 := uint32(b2[3])
	v2 := uint32(b2[2])
	v1 := uint32(b2[1])
	v0 := uint32(b2[0])
	return v0 | (v1 << 8) | (v2 << 16) | (v3 << 24)
}


func (b byteReader) unread() []byte {
	return b.b[b.off:]
}


func (b byteReader) remain() int {
	return len(b.b) - b.off
}
