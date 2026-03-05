package magic

import (
	"bytes"
	"encoding/binary"
)



func Shp(raw []byte, limit uint32) bool {
	if len(raw) < 112 {
		return false
	}

	if !(binary.BigEndian.Uint32(raw[0:4]) == 9994 &&
		binary.BigEndian.Uint32(raw[4:8]) == 0 &&
		binary.BigEndian.Uint32(raw[8:12]) == 0 &&
		binary.BigEndian.Uint32(raw[12:16]) == 0 &&
		binary.BigEndian.Uint32(raw[16:20]) == 0 &&
		binary.BigEndian.Uint32(raw[20:24]) == 0 &&
		binary.LittleEndian.Uint32(raw[28:32]) == 1000) {
		return false
	}

	shapeTypes := []int{
		0,  
		1,  
		3,  
		5,  
		8,  
		11, 
		13, 
		15, 
		18, 
		21, 
		23, 
		25, 
		28, 
		31, 
	}

	for _, st := range shapeTypes {
		if st == int(binary.LittleEndian.Uint32(raw[108:112])) {
			return true
		}
	}

	return false
}



func Shx(raw []byte, limit uint32) bool {
	return bytes.HasPrefix(raw, []byte{0x00, 0x00, 0x27, 0x0A})
}
