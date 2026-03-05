package magic

import (
	"bytes"
)

var (
	
	Woff = prefix([]byte("wOFF"))
	
	Woff2 = prefix([]byte("wOF2"))
	
	Otf = prefix([]byte{0x4F, 0x54, 0x54, 0x4F, 0x00})
)


func Ttf(raw []byte, limit uint32) bool {
	if !bytes.HasPrefix(raw, []byte{0x00, 0x01, 0x00, 0x00}) {
		return false
	}
	return !MsAccessAce(raw, limit) && !MsAccessMdb(raw, limit)
}


func Eot(raw []byte, limit uint32) bool {
	return len(raw) > 35 &&
		bytes.Equal(raw[34:36], []byte{0x4C, 0x50}) &&
		(bytes.Equal(raw[8:11], []byte{0x02, 0x00, 0x01}) ||
			bytes.Equal(raw[8:11], []byte{0x01, 0x00, 0x00}) ||
			bytes.Equal(raw[8:11], []byte{0x02, 0x00, 0x02}))
}


func Ttc(raw []byte, limit uint32) bool {
	return len(raw) > 7 &&
		bytes.HasPrefix(raw, []byte("ttcf")) &&
		(bytes.Equal(raw[4:8], []byte{0x00, 0x01, 0x00, 0x00}) ||
			bytes.Equal(raw[4:8], []byte{0x00, 0x02, 0x00, 0x00}))
}
