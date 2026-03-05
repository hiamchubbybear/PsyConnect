



//go:build amd64 && !purego && gc

package sha3



//go:noescape

func keccakF1600(a *[25]uint64)
