//go:build !race
// +build !race



package encoder

func encodeIntoCheckRace(buf *[]byte, val interface{}, opts Options) error {
	return encodeInto(buf, val, opts)
}
