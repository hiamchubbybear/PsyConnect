



package uuid

import (
	"io"
)













func NewV7() (UUID, error) {
	uuid, err := NewRandom()
	if err != nil {
		return uuid, err
	}
	makeV7(uuid[:])
	return uuid, nil
}




func NewV7FromReader(r io.Reader) (UUID, error) {
	uuid, err := NewRandomFromReader(r)
	if err != nil {
		return uuid, err
	}

	makeV7(uuid[:])
	return uuid, nil
}




func makeV7(uuid []byte) {
	
	_ = uuid[15] 

	t, s := getV7Time()

	uuid[0] = byte(t >> 40)
	uuid[1] = byte(t >> 32)
	uuid[2] = byte(t >> 24)
	uuid[3] = byte(t >> 16)
	uuid[4] = byte(t >> 8)
	uuid[5] = byte(t)

	uuid[6] = 0x70 | (0x0F & byte(s>>8))
	uuid[7] = byte(s)
}





var lastV7time int64

const nanoPerMilli = 1000000




func getV7Time() (milli, seq int64) {
	timeMu.Lock()
	defer timeMu.Unlock()

	nano := timeNow().UnixNano()
	milli = nano / nanoPerMilli
	
	seq = (nano - milli*nanoPerMilli) >> 8
	now := milli<<12 + seq
	if now <= lastV7time {
		now = lastV7time + 1
		milli = now >> 12
		seq = now & 0xfff
	}
	lastV7time = now
	return milli, seq
}
