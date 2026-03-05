



package uuid

import "encoding/binary"













func NewV6() (UUID, error) {
	var uuid UUID
	now, seq, err := GetTime()
	if err != nil {
		return uuid, err
	}

	

	binary.BigEndian.PutUint64(uuid[0:], uint64(now))
	binary.BigEndian.PutUint16(uuid[8:], seq)

	uuid[6] = 0x60 | (uuid[6] & 0x0F)
	uuid[8] = 0x80 | (uuid[8] & 0x3F)

	nodeMu.Lock()
	if nodeID == zeroID {
		setNodeInterface("")
	}
	copy(uuid[10:], nodeID[:])
	nodeMu.Unlock()

	return uuid, nil
}
