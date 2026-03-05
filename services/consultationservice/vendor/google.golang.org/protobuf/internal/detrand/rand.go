








package detrand

import (
	"encoding/binary"
	"hash/fnv"
	"os"
)



func Disable() {
	randSeed = 0
}


func Bool() bool {
	return randSeed%2 == 1
}


func Intn(n int) int {
	if n <= 0 {
		panic("must be positive")
	}
	return int(randSeed % uint64(n))
}


var randSeed = binaryHash()

func binaryHash() uint64 {
	
	s, err := os.Executable()
	if err != nil {
		return 0
	}
	f, err := os.Open(s)
	if err != nil {
		return 0
	}
	defer f.Close()

	
	const numSamples = 8
	var buf [64]byte
	h := fnv.New64()
	fi, err := f.Stat()
	if err != nil {
		return 0
	}
	binary.LittleEndian.PutUint64(buf[:8], uint64(fi.Size()))
	h.Write(buf[:8])
	for i := int64(0); i < numSamples; i++ {
		if _, err := f.ReadAt(buf[:], i*fi.Size()/numSamples); err != nil {
			return 0
		}
		h.Write(buf[:])
	}
	return h.Sum64()
}
