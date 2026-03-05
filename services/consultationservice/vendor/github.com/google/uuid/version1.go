



package uuid

import (
	"encoding/binary"
)









func NewUUID() (UUID, error) {
	var uuid UUID
	now, seq, err := GetTime()
	if err != nil {
		return uuid, err
	}

	timeLow := uint32(now & 0xffffffff)
	timeMid := uint16((now >> 32) & 0xffff)
	timeHi := uint16((now >> 48) & 0x0fff)
	timeHi |= 0x1000 

	binary.BigEndian.PutUint32(uuid[0:], timeLow)
	binary.BigEndian.PutUint16(uuid[4:], timeMid)
	binary.BigEndian.PutUint16(uuid[6:], timeHi)
	binary.BigEndian.PutUint16(uuid[8:], seq)

	nodeMu.Lock()
	if nodeID == zeroID {
		setNodeInterface("")
	}
	copy(uuid[10:], nodeID[:])
	nodeMu.Unlock()

	return uuid, nil
}
