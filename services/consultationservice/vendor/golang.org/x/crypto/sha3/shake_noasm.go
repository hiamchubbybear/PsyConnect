



//go:build !gc || purego || !s390x

package sha3

func newShake128() *state {
	return newShake128Generic()
}

func newShake256() *state {
	return newShake256Generic()
}
