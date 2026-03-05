



// +build !darwin,!dragonfly,!freebsd,!linux,!netbsd,!openbsd

package bio

func (r *Reader) sliceOS(length uint64) ([]byte, bool) {
	return nil, false
}
