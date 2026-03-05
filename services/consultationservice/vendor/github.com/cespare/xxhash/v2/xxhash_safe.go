//go:build appengine
// +build appengine



package xxhash


func Sum64String(s string) uint64 {
	return Sum64([]byte(s))
}


func (d *Digest) WriteString(s string) (n int, err error) {
	return d.Write([]byte(s))
}
